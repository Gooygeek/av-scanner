## av-scanner scan

Scan a file

```
av-scanner scan [flags]
```

### Options

```
      --clamav string            ClamAV server address
                                 Also Controlled by 'AVSCAN_CLAMAV_SERVER' env variable
                                  (default "tcp://127.0.0.1:3310")
      --db string                Path to the database file
                                 Also Controlled by 'AVSCAN_DB_PATH' env variable
                                  (default "av-scanner-results.db")
  -d, --dir string               Directory to watch or sweep for files
                                 Also Controlled by 'AVSCAN_WATCH_DIR' env variable
  -h, --help                     help for scan
      --loglevel string          Log level
                                 Must be one of 'DEBUG', 'INFO', 'WARNING', 'ERROR'
                                 Also Controlled by 'AVSCAN_LOG_LEVEL' env variable
                                  (default "INFO")
  -m, --mode string              Mode of operation
                                 'watch' to use fsnotify, 'sweep' to scan all files in the directory.
                                 Also Controlled by 'AVSCAN_MODE' env variable
                                  (default "sweep")
      --notify-endpoint string   Notification endpoint
                                 Also Controlled by 'AVSCAN_NOTIFY_ENDPOINT' env variable
      --notify-topic string      Notification topic
                                 Also Controlled by 'AVSCAN_NOTIFY_TOPIC' env variable
                                  (default "default")
      --relative-paths           Use relative paths for directories
                                 By default, the paths are converted to absolute
                                 Also Controlled by 'AVSCAN_RELATIVE_PATHS' env variable
                                 
```

### SEE ALSO

* [av-scanner](av-scanner.md)	 - AV-Scanner is a smart file scanner using ClamAV

###### Auto generated by spf13/cobra on 25-Apr-2025
