#!/usr/bin/env python3
import subprocess
import os
import re
import sys
import time
import argparse

TARGET_VERSION = "4.0.7"
MANAGEMENT_SUFFIX = "-management"
TARGET_TAG = f"{TARGET_VERSION}{MANAGEMENT_SUFFIX}"
INTERMEDIATE_VERSION = "3.13"
INTERMEDIATE_TAG = f"{INTERMEDIATE_VERSION}{MANAGEMENT_SUFFIX}"

# Control colored output
COLOUR_PRINT = 1

# Unicode symbols for messages
TICK = "\u2713"  # ✓
CROSS = "\u2717"  # ✗
WARNING_RUNE = "\u26A0"  # ⚠
INFO_RUNE = "\u2139"  # ℹ

ENV_FILE = ".env"

def print_error(message):
    """
    Prints an error message in red with a cross symbol if COLOUR_PRINT is enabled.
    """
    if COLOUR_PRINT:
        print(f"\033[91m[ERROR] {message} {CROSS}\033[0m")
    else:
        print_error(f"{message}")

def print_warning(message):
    """
    Prints a warning message in yellow with a warning symbol if COLOUR_PRINT is enabled.
    """
    if COLOUR_PRINT:
        print(f"\033[93m[WARNING] {message} {WARNING_RUNE}\033[0m")
    else:
        print(f"[WARNING] {message}")

def print_success(message):
    """
    Prints a success message in green with a tick symbol if COLOUR_PRINT is enabled.
    """
    if COLOUR_PRINT:
        print(f"\033[92m[SUCCESS] {message} {TICK}\033[0m")
    else:
        print_success(f"{message}")

def print_info(message):
    """
    Prints an informational message in blue with an info symbol if COLOUR_PRINT is enabled.
    """
    if COLOUR_PRINT:
        print(f"\033[94m[INFO] {message} {INFO_RUNE}\033[0m")
    else:
        print(f"[INFO] {message}")

def detect_docker_compose_plugin():
    """
    Detects whether the 'docker compose' plugin is available.
    Falls back to 'docker-compose' if the plugin is not available.
    """
    try:
        subprocess.run(["docker", "compose", "version"], check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        return "docker compose"
    except subprocess.CalledProcessError:
        print_warning("'docker compose' plugin not found. Falling back to 'docker-compose'.")
        return "docker-compose"

def get_rabbitmq_volume(compose_cmd, compose_yml):
    """
    Retrieves the volume name used by the RabbitMQ service.
    """
    volume_name = None
    try:
        command = f"{compose_cmd} -f {compose_yml} ps -q rabbitmq | xargs docker inspect --format='{{{{ range .Mounts }}}}{{{{ if eq .Destination \"/var/lib/rabbitmq\" }}}}{{{{ .Name }}}}{{{{ end }}}}{{{{ end }}}}'"
        volume_name = subprocess.run(command, shell=True, check=True, text=True, stdout=subprocess.PIPE, env=os.environ.copy()).stdout.strip()
    except subprocess.CalledProcessError:
        pass
    if not volume_name:
        print_error("Could not determine the RabbitMQ volume name, please check that rabbitmq service is running.")
        sys.exit(1)
    return volume_name

def run_command(command, description):
    """
    Runs a shell command and checks its result.
    If the command fails, it stops the script execution.
    """
    print_info(description)
    try:
        result = subprocess.run(command, shell=True, check=True, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, env=os.environ.copy())
        print_success(description)
        print(result.stdout)
        return result.stdout.strip()
    except subprocess.CalledProcessError as e:
        print_error(description)
        print(e.stderr)
        sys.exit(1)

def backup_rabbitmq_data(compose_cmd, compose_yml, backup_file):
    """
    Backs up the RabbitMQ data volume to a tar.gz file.
    """
    if os.path.exists(backup_file):
        print_warning(f"Backup file {backup_file} already exists. Overwriting...")
    rabbitmq_volume = get_rabbitmq_volume(compose_cmd, compose_yml)
    command = f"docker run --rm -v {rabbitmq_volume}:/_data -v $(pwd):/backup alpine tar czf /backup/{backup_file} -C /_data ."
    run_command(command, f"Backing up RabbitMQ data to {backup_file}")

def verify_backup(backup_file):
    """
    Verifies that the backup file exists and is not empty.
    """
    if not os.path.exists(backup_file):
        print_error(f"Backup file {backup_file} does not exist.")
        sys.exit(1)
    if os.path.getsize(backup_file) == 0:
        print_error(f"Backup file {backup_file} is empty.")
        sys.exit(1)
    print_success(f"Backup file {backup_file} verified.")

def check_rabbitmq_version(intermediate_version):
    """
    Retrieves the RabbitMQ version from .env file and fails when it is not found or it's newer than the target version.
    """
    if not os.path.exists(ENV_FILE):
        print_error(f"{ENV_FILE} not found.")
        sys.exit(1)
    with open(ENV_FILE, "r") as file:
        content = file.read()
    match = re.search(r"^RABBITMQ_TAG=(.*)", content, re.MULTILINE)
    if not match:
        print_error("RabbitMQ tag not found in .env file.")
        sys.exit(1)
    version = match.group(1).removesuffix(MANAGEMENT_SUFFIX)
    if version >= intermediate_version:
        print_error(f"RabbitMQ version {version} already newer than {intermediate_version} in {ENV_FILE} file.")
        sys.exit(1)
    print_success(f"RabbitMQ version {version} verified in {ENV_FILE}.")

def update_env_file(rabbitmq_tag):
    """
    Updates the RabbitMQ tag in the .env file.
    """
    if not os.path.exists(ENV_FILE):
        print_error(f"{ENV_FILE} not found.")
        sys.exit(1)
    with open(ENV_FILE, "r") as file:
        content = file.read()
    updated_content = re.sub(r"^RABBITMQ_TAG=.*", f"RABBITMQ_TAG={rabbitmq_tag}", content, flags=re.MULTILINE)
    if updated_content == content:
        print_error(f"RabbitMQ tag not found in {ENV_FILE}.")
        sys.exit(1)
    with open(ENV_FILE, "w") as file:
        file.write(updated_content)
    print_success(f"Updated RabbitMQ tag to {rabbitmq_tag} in {ENV_FILE}")

def stop_and_remove_rabbitmq(compose_cmd, compose_yml):
    """
    Stops and removes the RabbitMQ container.
    """
    run_command(f"{compose_cmd} -f {compose_yml} stop rabbitmq", "Stopping RabbitMQ container")
    run_command(f"{compose_cmd} -f {compose_yml} rm -f rabbitmq", "Removing RabbitMQ container")

def start_rabbitmq(compose_cmd, compose_yml):
    """
    Starts the RabbitMQ container with the updated version.
    """
    run_command(f"{compose_cmd} -f {compose_yml} up -d rabbitmq", "Starting RabbitMQ container")

def check_rabbitmq_ready(compose_cmd, compose_yml):
    """
    Waits until RabbitMQ is fully ready by polling its health status.
    """
    print_info("Checking RabbitMQ readiness...")
    for attempt in range(30):  # Retry for up to 30 seconds
        try:
            result = subprocess.run(
                f"{compose_cmd} -f {compose_yml} logs rabbitmq",
                shell=True,
                check=True,
                text=True,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                env=os.environ.copy()
            )
            if "time to start rabbitmq:" in result.stdout.lower():
                print_success("RabbitMQ is ready.")
                return
        except subprocess.CalledProcessError:
            pass
        sys.stdout.write(".")  # Print a dot to indicate a retry
        sys.stdout.flush()
        time.sleep(1)
    print_error("RabbitMQ did not become ready in time.")
    sys.exit(1)

def enable_stable_feature_flags(compose_cmd, compose_yml):
    """
    Enables stable feature flags in RabbitMQ.
    """
    command = f"{compose_cmd} -f {compose_yml} exec rabbitmq rabbitmqctl enable_feature_flag all"
    run_command(command, "Enabling all stable feature flags...")

def verify_rabbitmq_version(compose_cmd, compose_yml, expected_version):
    """
    Verifies the RabbitMQ version running in the container.
    """
    command = f"{compose_cmd} -f {compose_yml} exec rabbitmq rabbitmqctl status | grep \"RabbitMQ version\""
    output = run_command(command, "Checking RabbitMQ version")
    if expected_version not in output:
        print_error(f"RabbitMQ version mismatch. Expected: {expected_version}, Found: {output}")
        sys.exit(1)
    print_success(f"RabbitMQ version {expected_version} verified.")

def check_if_any_deprecated_are_used(compose_cmd, compose_yml):
    """
    Checks if any deprecated RabbitMQ features are used.
    """
    command = f"{compose_cmd} -f {compose_yml} exec rabbitmq rabbitmq-diagnostics check_if_any_deprecated_features_are_used"
    output = run_command(command, "Checking for deprecated features")
    if "no deprecated features in use" in output.lower():
        print_success("No deprecated RabbitMQ features found.")
    else:
        print_warning("Deprecated RabbitMQ features are used. Please update your configuration.")
        print(output)

def main():
    parser = argparse.ArgumentParser(description="RabbitMQ Upgrade Script")
    parser.add_argument("--compose-file", default="docker-compose.yml", help="Path to the docker-compose.yml file")
    args = parser.parse_args()

    compose_yml = args.compose_file
    if not os.path.exists(compose_yml):
        print_error(f"file {compose_yml} not found.")
        sys.exit(1)

    compose_cmd = detect_docker_compose_plugin()

    # Phase 1: Upgrade to RabbitMQ intermediate version
    print_info(f"=== Phase 1: Upgrade to RabbitMQ {INTERMEDIATE_VERSION} ===")
    check_rabbitmq_version(INTERMEDIATE_VERSION)
    backup_file_initial = "rabbitmqdata-backup-initial.tar.gz"
    backup_rabbitmq_data(compose_cmd, compose_yml, backup_file_initial)
    verify_backup(backup_file_initial)
    update_env_file(INTERMEDIATE_TAG)
    stop_and_remove_rabbitmq(compose_cmd, compose_yml)
    start_rabbitmq(compose_cmd, compose_yml)
    check_rabbitmq_ready(compose_cmd, compose_yml)
    enable_stable_feature_flags(compose_cmd, compose_yml)
    verify_rabbitmq_version(compose_cmd, compose_yml, INTERMEDIATE_VERSION)

    # Phase 2: Upgrade to RabbitMQ final version
    print_info(f"=== Phase 2: Upgrade to RabbitMQ {TARGET_VERSION} ===")
    backup_file_intermediate = f"rabbitmqdata-backup-{INTERMEDIATE_VERSION}.tar.gz"
    backup_rabbitmq_data(compose_cmd, compose_yml, backup_file_intermediate)
    verify_backup(backup_file_intermediate)
    stop_and_remove_rabbitmq(compose_cmd, compose_yml)
    update_env_file(TARGET_TAG)
    start_rabbitmq(compose_cmd, compose_yml)
    check_rabbitmq_ready(compose_cmd, compose_yml)
    verify_rabbitmq_version(compose_cmd, compose_yml, TARGET_VERSION)
    enable_stable_feature_flags(compose_cmd, compose_yml)
    check_if_any_deprecated_are_used(compose_cmd, compose_yml)

    print_info("RabbitMQ upgrade completed successfully!")

if __name__ == "__main__":
    main()