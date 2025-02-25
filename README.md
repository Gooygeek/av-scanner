# av-scanner

A simple antivirus file scanner.

---

This tool uses clamav to explicitly scan a directory and scan for viruses

These scans can be either once-off, periodically via cron, or watch for new files.

It saves the scan results in a database to prevent rescanning and sends notificaitons of the scan results.

## How To Use

The most basic usage is to trigger from the commandline:

```shell
av-scanner scan --dir some_directory
```

Subsiquent reruns will only scan new or updated files.

You can also trigger once and have it scan new files as they are created (**experimental!**)

```shell
av-scanner scan --dir some_directory --mode watch
```

Additional commands other than `scan` and their applicable settings can be found in the [docs-md](./docs-md) folder.

### Environment Variables

Some of the flags can also be configured as environment variables (although the flags will take precendence):

- `AVSCAN_MODE`: Mode of operation (`sweep` or `watch`)
- `AVSCAN_WATCH_DIR`: Directory to watch or sweep for files
- `AVSCAN_DB_PATH`: Path to the database file
- `AVSCAN_CLAMAV_SERVER`: ClamAV server address
- `AVSCAN_LOG_LEVEL`: Log level
- `AVSCAN_NOTIFY_ENDPOINT`: Notification endpoint
- `AVSCAN_NOTIFY_TOPIC`: Notification topic

## Shell Completion

The av-scanner has a built-in shell completion feature. It works by calling the `completion` command and passing the shell type as an argument.

This will generate a completion script that can be sourced in the shell.

For example, to generate a bash completion script, run:

```shell
source <(av-scanner completion bash)
```

See the documentation for more information on how to use the completion feature or run:

```shell
av-scanner completion --help
```

## Set Up for Automatic Scanning

### Periodic Rescanning

To set up periodic rescanning using a cron job, follow these steps:

1. Open your crontab file for editing:
  ```shell
  crontab -e
  ```

2. Add a new line to schedule the `av-scanner` to run at your desired interval. For example, to run the scanner every day at midnight, add the following line:
  ```shell
  0 0 * * * /usr/local/bin/av-scanner scan --dir /path/to/your_directory
  ```

### A systemd service

The systemd service mode is **experimental**. It uses the `watch` mode to run a single persistent instance that scans files as they are created.

To set up the scanner as a systemd service, follow these steps:

1. Create a new service file for the `av-scanner`:
  ```shell
  sudo vi /etc/systemd/system/av-scanner.service
  ```

2. Add the following content to the service file:

  Modify the flags of the command as needed, as well as any cross-service dependancies

  ```ini
  [Unit]
  Description=AV Scanner Service
  After=network.target

  [Service]
  ExecStart=/usr/local/bin/av-scanner scan --dir /path/to/your_directory --mode watch
  Restart=always

  [Install]
  WantedBy=multi-user.target
  ```


3. Reload the systemd daemon to apply the changes:
  ```shell
  sudo systemctl daemon-reload
  ```

4. Enable the `av-scanner` service to start on boot:
  ```shell
  sudo systemctl enable av-scanner
  ```

5. Start the `av-scanner` service:
  ```shell
  sudo systemctl start av-scanner
  ```

## Tech Stack

* golang: Main programming language
* clamav: Virus scanning
* sqlite: Database
* ntfy: Nofitication endpoint
* file: Determining MIME type
* cron: Periodic rescanning
* systemd: Starting as a service
