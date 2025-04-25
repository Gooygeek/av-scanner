package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:     "av-scanner",
	Short:   "AV-Scanner is a smart file scanner using ClamAV",
	Version: "0.1.0",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of AV-Scanner",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("AV-Scanner version:", rootCmd.Version)
	},
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a file",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generates Docs",
	Run: func(cmd *cobra.Command, args []string) {
		format := viper.GetString("format")

		if format == "all" || format == "md" {
			err := os.MkdirAll("docs-md", os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
			err = doc.GenMarkdownTree(rootCmd, "docs-md")
			if err != nil {
				log.Fatal(err)
			}
		}

		if format == "all" || format == "yaml" {
			err := os.MkdirAll("docs-yaml", os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
			err = doc.GenYamlTree(rootCmd, "docs-yaml")
			if err != nil {
				log.Fatal(err)
			}
		}

		if format == "all" || format == "man" {
			err := os.MkdirAll("docs-man", os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
			err = doc.GenManTree(rootCmd, &doc.GenManHeader{
				Title:   "AV-Scanner",
				Section: "1",
				Manual:  "AV-Scanner Manual",
			}, "docs-man")
			if err != nil {
				log.Fatal(err)
			}
		}

		if format == "all" || format == "rst" {
			err := os.MkdirAll("docs-rst", os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
			err = doc.GenReSTTree(rootCmd, "docs-rst")
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	// Add scan command
	rootCmd.AddCommand(scanCmd)

	scanCmd.PersistentFlags().StringP("mode", "m", "sweep", "Mode of operation\n'watch' to use fsnotify, 'sweep' to scan all files in the directory.\nAlso Controlled by 'AVSCAN_MODE' env variable\n")
	scanCmd.RegisterFlagCompletionFunc("mode", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"sweep\tRecursivly discover files in a directory and scan all of them", "watch\tWait for a file creation event and then scan that file"}, cobra.ShellCompDirectiveDefault
	})
	viper.BindPFlag("mode", scanCmd.PersistentFlags().Lookup("mode"))
	viper.BindEnv("mode", "AVSCAN_MODE")

	scanCmd.PersistentFlags().StringP("dir", "d", "", "Directory to watch or sweep for files\nAlso Controlled by 'AVSCAN_WATCH_DIR' env variable")
	scanCmd.MarkPersistentFlagRequired("dir")
	viper.BindPFlag("dir", scanCmd.PersistentFlags().Lookup("dir"))
	viper.BindEnv("dir", "AVSCAN_WATCH_DIR")

	scanCmd.PersistentFlags().String("db", "av-scanner-results.db", "Path to the database file\nAlso Controlled by 'AVSCAN_DB_PATH' env variable\n")
	viper.BindPFlag("db", scanCmd.PersistentFlags().Lookup("db"))
	viper.BindEnv("db", "AVSCAN_DB_PATH")

	scanCmd.PersistentFlags().String("clamav", "tcp://127.0.0.1:3310", "ClamAV server address\nAlso Controlled by 'AVSCAN_CLAMAV_SERVER' env variable\n")
	viper.BindPFlag("clamav", scanCmd.PersistentFlags().Lookup("clamav"))
	viper.BindEnv("clamav", "AVSCAN_CLAMAV_SERVER")

	scanCmd.PersistentFlags().String("loglevel", "INFO", "Log level\nMust be one of 'DEBUG', 'INFO', 'WARNING', 'ERROR'\nAlso Controlled by 'AVSCAN_LOG_LEVEL' env variable\n")
	scanCmd.RegisterFlagCompletionFunc("loglevel", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"DEBUG", "INFO", "WARNING", "ERROR"}, cobra.ShellCompDirectiveDefault
	})
	viper.BindPFlag("loglevel", scanCmd.PersistentFlags().Lookup("loglevel"))
	viper.BindEnv("loglevel", "AVSCAN_LOG_LEVEL")

	scanCmd.PersistentFlags().String("notify-endpoint", "", "Notification endpoint\nAlso Controlled by 'AVSCAN_NOTIFY_ENDPOINT' env variable")
	viper.BindPFlag("notify-endpoint", scanCmd.PersistentFlags().Lookup("notify-endpoint"))
	viper.BindEnv("notify-endpoint", "AVSCAN_NOTIFY_ENDPOINT")

	scanCmd.PersistentFlags().String("notify-topic", "default", "Notification topic\nAlso Controlled by 'AVSCAN_NOTIFY_TOPIC' env variable\n")
	viper.BindPFlag("notify-topic", scanCmd.PersistentFlags().Lookup("notify-topic"))
	viper.BindEnv("notify-topic", "AVSCAN_NOTIFY_TOPIC")

	scanCmd.PersistentFlags().Bool("relative-paths", false, "Use relative paths for directories\nBy default, the paths are converted to absolute\nAlso Controlled by 'AVSCAN_RELATIVE_PATHS' env variable\n")
	viper.BindPFlag("relative-paths", scanCmd.PersistentFlags().Lookup("relative-paths"))
	viper.BindEnv("relative-paths", "AVSCAN_RELATIVE_PATHS")

	// Add version command
	rootCmd.AddCommand(versionCmd)

	// Add docs command
	rootCmd.AddCommand(docsCmd)

	docsCmd.PersistentFlags().StringP("format", "f", "all", "Documentation format to generate. Options are: 'all', 'md', 'yaml', 'man', 'rst'")
	docsCmd.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"all", "md", "yaml", "man", "rst"}, cobra.ShellCompDirectiveDefault
	})
	viper.BindPFlag("format", docsCmd.PersistentFlags().Lookup("format"))
}

func initConfig() {
	viper.AutomaticEnv()
}

func run() {
	logger = Logger{logLevel: viper.GetString("loglevel")}
	notifyer = Notifyer{endpoint: viper.GetString("notify-endpoint"), topic: viper.GetString("notify-topic")}

	logger.Info("Starting av-scanner...")

	db, err := initDB(viper.GetString("db"))
	if err != nil {
		logger.Fatal("Error initializing database: " + err.Error())
	}
	defer db.Close()

	mode := viper.GetString("mode")
	dirToWatch := viper.GetString("dir")
	useRelativePaths := viper.GetBool("relative-paths")

	if !useRelativePaths {
		absDirToWatch, err := filepath.Abs(dirToWatch)
		if err != nil {
			logger.Fatal("Error converting directory to absolute path: " + err.Error())
		}
		if absDirToWatch != dirToWatch {
			logger.Warning("Relative paths can lead to unintended consiquences. The scanned directory has been converted to: '" + absDirToWatch + "'. This can be overriden with the --relative-paths flag or AVSCAN_RELATIVE_PATHS environment variable.")
		}
		dirToWatch = absDirToWatch
	}

	if mode == "watch" {
		logger.Info("Operating in watch mode")
		watchDirectory(dirToWatch, db)
	} else if mode == "sweep" {
		logger.Info("Operating in sweep mode")
		scanDirectory(dirToWatch, db)
	} else {
		logger.Fatal("Invalid mode specified. Use 'watch' or 'sweep'.")
	}
}

func watchDirectory(dirToWatch string, db *sql.DB) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal("Error creating watcher: " + err.Error())
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					logger.Debug("Watcher event channel closed")
					return
				}
				logger.Debug("Event detected: " + event.String())
				if event.Op&fsnotify.Create == fsnotify.Create {
					logger.Info("New file detected: " + event.Name)
					processFile(event.Name, db)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					logger.Debug("Watcher error channel closed")
					return
				}
				logger.Error("Error while watching: " + err.Error())
			}
		}
	}()

	logger.Info("Creating directory (if required): " + dirToWatch)
	err = os.MkdirAll(dirToWatch, os.ModePerm)
	if err != nil {
		logger.Fatal("Error creating directory: " + err.Error())
	}

	logger.Info("Adding directory to watcher: " + dirToWatch)
	err = watcher.Add(dirToWatch)
	if err != nil {
		logger.Fatal("Error adding directory to watcher: " + err.Error())
	}

	logger.Info("Watching directory: " + dirToWatch)
	<-done
}

func scanDirectory(dirToWatch string, db *sql.DB) {
	err := filepath.Walk(dirToWatch, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			logger.Debug("File Detected: " + path)
			processFile(path, db)
		}
		return nil
	})
	if err != nil {
		logger.Fatal("Error scanning directory: " + err.Error())
	}
}

func processFile(filePath string, db *sql.DB) {
	// Get file size and modified time
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		logger.Error("Error getting file size: " + err.Error())
		return
	}
	lastModified := fileInfo.ModTime()
	logger.Debug("Last modified: " + lastModified.String())
	fileSize := fileInfo.Size()
	logger.Debug(fmt.Sprintf("File size: %d bytes", fileSize))

	// Check if the file was modifed since last scan
	prevScanDate, err := getFilePrevScanDate(db, filePath)
	if err != nil {
		logger.Error("Error getting the last scan date: " + err.Error())
		return
	}
	logger.Debug("Previous scan date: " + prevScanDate.String())
	if prevScanDate.After(lastModified) {
		logger.Debug("File not modified since last scan (" + prevScanDate.String() + "), skipping")
		return
	}

	// Compute Hash
	hash, err := computeSHA256(filePath)
	if err != nil {
		logger.Error("Error computing SHA256: " + err.Error())
		return
	}
	logger.Debug("SHA256 hash: " + hash)

	// Check if file already exists in DB
	exists, err := fileExistsInDB(db, hash)
	if err != nil {
		logger.Error("Error checking file in database: " + err.Error())
		return
	}
	if exists {
		logger.Info("File already processed, skipping")
		return
	}

	// Get file MIME type
	mimeType, err := getFileMimeType(filePath)
	if err != nil {
		logger.Error("Error getting file MIME type: " + err.Error())
		return
	}
	logger.Debug("MIME type: " + mimeType)

	// Scan file with ClamAV
	scanResult := "SKIPPED"
	var scanLog string
	// Skip large files (greater than 500MB)
	if fileSize > 1024*1024*1024*0.5 {
		logger.Info("File too large to scan: " + filePath)
	} else {
		logger.Info("Scanning: " + filePath)
		scanResult, scanLog, err = scanFileWithClamAV(filePath, viper.GetString("clamav"))
		if err != nil {
			logger.Error("Error scanning file with ClamAV: " + err.Error())
			return
		}
	}

	logger.Info(fmt.Sprintf("ClamAV scan result for %s: %s", filePath, scanResult))

	// Notify
	formattedString := fmt.Sprintf("Scan Status: %s\nDetected MIME Type: %s\nSize: %d Bytes\nLog: %s", scanResult, mimeType, int(fileSize), scanLog)
	if scanResult == "FOUND" {
		notifyer.Detection("File Scanned: "+filePath, formattedString)
	} else if scanResult == "SKIPPED" {
		notifyer.Warn("File Scanned: "+filePath, formattedString)
	} else if scanResult == "OK" {
		notifyer.OK("File Scanned: "+filePath, formattedString)
	} else {
		notifyer.Info("File Scanned: "+filePath, formattedString)
	}

	// Save file info to DB
	err = saveFileInfo(db, filePath, hash, scanResult, scanLog)
	if err != nil {
		logger.Error("Error saving file info to database: " + err.Error())
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
