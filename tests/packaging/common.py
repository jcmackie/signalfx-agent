import re
import subprocess
import threading
import time
from contextlib import contextmanager
from pathlib import Path

import docker
from tests.helpers import fake_backend
from tests.helpers.util import get_docker_client, get_host_ip, pull_from_reader_in_background, retry, run_container
from tests.paths import REPO_ROOT_DIR

PACKAGING_DIR = REPO_ROOT_DIR / "packaging"
DEPLOYMENTS_DIR = REPO_ROOT_DIR / "deployments"
INSTALLER_PATH = DEPLOYMENTS_DIR / "installer/install.sh"
RPM_OUTPUT_DIR = PACKAGING_DIR / "rpm/output/x86_64"
DEB_OUTPUT_DIR = PACKAGING_DIR / "deb/output"
DOCKERFILES_DIR = Path(__file__).parent.joinpath("images").resolve()

INIT_SYSV = "sysv"
INIT_UPSTART = "upstart"
INIT_SYSTEMD = "systemd"

AGENT_YAML_PATH = "/etc/signalfx/agent.yaml"
PIDFILE_PATH = "/var/run/signalfx-agent.pid"

BASIC_CONFIG = """
monitors:
  - type: collectd/signalfx-metadata
  - type: collectd/cpu
  - type: collectd/uptime
"""


def build_base_image(name, path=DOCKERFILES_DIR, dockerfile=None, buildargs=None):
    client = get_docker_client()
    dockerfile = dockerfile or Path(path) / f"Dockerfile.{name}"
    image, _ = client.images.build(
        path=str(path), dockerfile=str(dockerfile), pull=True, rm=True, forcerm=True, buildargs=buildargs
    )

    return image.id


LOG_COMMAND = {
    INIT_SYSV: "cat /var/log/signalfx-agent.log",
    INIT_UPSTART: "cat /var/log/signalfx-agent.log",
    INIT_SYSTEMD: "journalctl -u signalfx-agent",
}


def get_agent_logs(container, init_system):
    try:
        _, output = container.exec_run(LOG_COMMAND[init_system])
    except docker.errors.APIError as e:
        print("Error getting agent logs: %s" % e)
        return ""
    return output


def get_deb_package_to_test():
    return get_package_to_test(DEB_OUTPUT_DIR, "deb")


def get_rpm_package_to_test():
    return get_package_to_test(RPM_OUTPUT_DIR, "rpm")


def get_package_to_test(output_dir, extension):
    pkgs = list(Path(output_dir).glob(f"*.{extension}"))
    if not pkgs:
        raise AssertionError(f"No .{extension} files found in {output_dir}")

    if len(pkgs) > 1:
        raise AssertionError(f"More than one .{extension} file found in {output_dir}")

    return pkgs[0]


# Run an HTTPS proxy inside the container with socat so that our fake backend
# doesn't have to worry about HTTPS.  The cert file must be trusted by the
# container running the agent.
# This is pretty hacky but docker makes it hard to communicate from a container
# back to the host machine (and we don't want to use the host network stack in
# the container due to init systems).  The idea is to bind mount a shared
# folder from the test host to the container that two socat instances use to
# communicate using a file to make the bytes flow between the HTTPS proxy and
# the fake backend.
@contextmanager
def socat_https_proxy(container, target_host, target_port, source_host, bind_addr):
    cert = "/%s.cert" % source_host
    key = "/%s.key" % source_host

    socat_bin = DOCKERFILES_DIR / "socat"
    stopped = False
    socket_path = "/tmp/scratch/%s-%s" % (source_host, container.id[:12])

    # Keep the socat instance in the container running across container
    # restarts
    def keep_running_in_container(cont, sock):
        while not stopped:
            try:
                cont.exec_run(
                    [
                        "socat",
                        "-v",
                        "OPENSSL-LISTEN:443,cert=%s,key=%s,verify=0,bind=%s,fork" % (cert, key, bind_addr),
                        "UNIX-CONNECT:%s" % sock,
                    ]
                )
            except docker.errors.APIError:
                print("socat died, restarting...")
                time.sleep(0.1)

    threading.Thread(target=keep_running_in_container, args=(container, socket_path)).start()

    proc = subprocess.Popen(
        [socat_bin, "-v", "UNIX-LISTEN:%s,fork" % socket_path, "TCP4:%s:%d" % (target_host, target_port)],
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
    )

    get_local_out = pull_from_reader_in_background(proc.stdout)

    try:
        yield
    finally:
        stopped = True
        # The socat instance in the container will die with the container
        proc.kill()
        print(get_local_out())


@contextmanager
def run_init_system_image(
    base_image,
    with_socat=True,
    path=DOCKERFILES_DIR,
    dockerfile=None,
    ingest_host="ingest.us0.signalfx.com",  # Whatever value is used here needs a self-signed cert in ./images/certs/
    api_host="api.us0.signalfx.com",  # Whatever value is used here needs a self-signed cert in ./images/certs/
    command=None,
    buildargs=None,
):  # pylint: disable=too-many-arguments
    image_id = retry(lambda: build_base_image(base_image, path, dockerfile, buildargs), docker.errors.BuildError)
    print("Image ID: %s" % image_id)
    if with_socat:
        backend_ip = "127.0.0.1"
    else:
        backend_ip = get_host_ip()
    with fake_backend.start(ip_addr=backend_ip) as backend:
        container_options = {
            # Init systems running in the container want permissions
            "privileged": True,
            "volumes": {
                "/sys/fs/cgroup": {"bind": "/sys/fs/cgroup", "mode": "ro"},
                "/tmp/scratch": {"bind": "/tmp/scratch", "mode": "rw"},
            },
            "extra_hosts": {
                # Socat will be running on localhost to forward requests to
                # these hosts to the fake backend
                ingest_host: backend.ingest_host,
                api_host: backend.api_host,
            },
        }

        if command:
            container_options["command"] = command

        with run_container(image_id, wait_for_ip=True, **container_options) as cont:
            if with_socat:
                # Proxy the backend calls through a fake HTTPS endpoint so that we
                # don't have to change the default configuration default by the
                # package.  The base_image used should trust the self-signed certs
                # default in the images dir so that the agent doesn't throw TLS
                # verification errors.
                with socat_https_proxy(
                    cont, backend.ingest_host, backend.ingest_port, ingest_host, "127.0.0.1"
                ), socat_https_proxy(cont, backend.api_host, backend.api_port, api_host, "127.0.0.2"):
                    yield [cont, backend]
            else:
                yield [cont, backend]


def is_agent_running_as_non_root(container):
    code, output = container.exec_run("pgrep -u signalfx-agent signalfx-agent")
    print("pgrep check: %s" % output)
    return code == 0


def get_agent_version(cont):
    code, output = cont.exec_run("signalfx-agent -version")
    output = output.decode("utf-8").strip()
    assert code == 0, "command 'signalfx-agent -version' failed:\n%s" % output
    match = re.match("^.+?: (.+)?,", output)
    assert match and match.group(1).strip(), "failed to parse agent version from command output:\n%s" % output
    return match.group(1).strip()


def get_win_agent_version(agent_path=r"C:\Program Files\SignalFx\SignalFxAgent\bin\signalfx-agent.exe"):
    proc = subprocess.run(agent_path + " -version", stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    output = proc.stdout.decode("utf-8")
    assert proc.returncode == 0, "command '%s -version' failed:\n%s" % (agent_path, output)
    match = re.match("^.+?: (.+)?,", output)
    assert match and match.group(1).strip(), "failed to parse agent version from command output:\n%s" % output
    return match.group(1).strip()
