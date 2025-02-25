#!/bin/bash

WATCH_DIR="$2"
SCANNER_ROOT_DIR="/VIRUS"
QUARANTEEN_DIR=$SCANNER_ROOT_DIR/quaranteen
LOG_DIR=$SCANNER_ROOT_DIR/log
# Set the maximum file size to scan in KB.
# This limits the files that is passed to clamav.
# Clamav is still set to the max limit,
#     this accounts for extracting archives and different algorithms used to determine filesize.
MAX_FILE_SIZE=$((1024*1024*1024/2)) # 500 MB

NTFY_ENDPOINT="https://ntfy.sh"
NTFY_TOPIC="1234-av-alerts"


# Create directories
mkdir -p $QUARANTEEN_DIR
mkdir -p $LOG_DIR


# Function to check if a file is still being written
wait_for_complete_write() {
    local file="$1"

    while true; do
        # Get the current size of the file
        FILESIZE=$(stat --format=%s "$file")

        # Wait for a short interval
        sleep 1

        # Get the new size after waiting
        NEWFILESIZE=$(stat --format=%s "$file")

        # If the size hasn't changed, the writing is likely done
        if [ "$FILESIZE" -eq "$NEWFILESIZE" ]; then
            break
        fi
    done
}
export -f wait_for_complete_write


# Send notification
notify() {
    TAG="$1"
    TITLE="$2"
    PRIORITY="$3"
    MESSAGE="$4"
    curl -H "Tags: ${TAG}" \
         -H "Title: ${TITLE}" \
         -H "Priority: ${PRIORITY}" \
         -d "${MESSAGE}" \
         ${NTFY_ENDPOINT}/${NTFY_TOPIC}
}
export -f notify


# Perform the full scan cycle on a file
scan() {
    NEWFILE=$1

    # Skip folders, only scan files
    if [ -d "${NEWFILE}" ] ; then
        echo "Seems to be a folder, skipping"
    else
        echo "Waiting for file write to complete..."
        wait_for_complete_write "$NEWFILE"
        echo "done"

        # Get file type (mime)
        mime=$(file --mime-type --mime-encoding -b "$NEWFILE")

        # Get size
        size=$(stat -c%s "$NEWFILE")

        severity_level="Low"

        filename_base64=$(basename "$NEWFILE" | base64 --wrap=0 | head -c 100)

        # Based on size, check the mime type or scan
        if ((size > MAX_FILE_SIZE)); then
            scan_attempted=false
            if [[ $mime != video/* ]]; then
                severity_level="High"
                scan_detection=unknown
            else
                severity_level="Medium"
                scan_detection=unknown
            fi
            echo "Severity: ${severity_level}" >> "$LOG_DIR/$filename_base64"
            echo "Virus detected: ${scan_detection}" >> "$LOG_DIR/$filename_base64"
            echo "Detected MIME type: ${mime}" >> "$LOG_DIR/$filename_base64"
            echo "Size: ${size}" >> "$LOG_DIR/$filename_base64"
            echo "Scan Attempted: ${scan_attempted}" >> "$LOG_DIR/$filename_base64"
        else
            scan_attempted=true
            echo "scanning..."
            clamscan -r --move=$QUARANTEEN_DIR --max-filesize=4000M --max-scansize=4000M "$NEWFILE" | tee "$LOG_DIR/$filename_base64"
            echo "done"
            # Determine if a virus was found
            if [ "x$(grep "Infected files: 0" "$LOG_DIR/$filename_base64")" = "x" ]; then
                scan_detection=true
                severity_level="High"

                # Strip all permissions
                chmod 000 "$QUARANTEEN_DIR/$(basename "$NEWFILE")"
            else
                scan_detection=false
            fi
        fi

        # Compile results
        TITLE="File Scanned: $NEWFILE"
        MESSAGE="""
Severity: ${severity_level}
Virus detected: ${scan_detection}
Detected MIME type: ${mime}
Size: ${size}
Scan Attempted: ${scan_attempted}
$(if [ "x${scan_attempted}" != "x" ]; then echo "Scan Log:"; cat "$LOG_DIR/$filename_base64"; fi)
"""
        if [ "${severity_level}" = "High" ]; then
            TAG="bangbang"
            PRIORITY="urgent"
        elif [ "${severity_level}" = "Medium" ]; then
            TAG="warning"
            PRIORITY="default"
        elif [ "${severity_level}" = "Low" ]; then
            TAG="information_source"
            PRIORITY="default"
        fi

        # Send a notification
        notify "${TAG}" "${TITLE}" "${PRIORITY}" "${MESSAGE}"
    fi
}
export -f scan


# Hash and check against the existing scan logs.
# If the file has not previously been scanned, do so
check_and_scan() {
    FILE="$1"

    filename_base64=$(basename "$FILE" | base64 --wrap=0 | head -c 100)
    if [ "x$(ls -l "${LOG_DIR}" | grep "${filename_base64}")" = "x" ]; then
        echo "File has not been scanned previously. Scanning..."
        scan "${FILE}"
        echo "Scan complete"
    else
        echo "File has been scanned previously"
    fi
}
export -f check_and_scan


if [ "$#" -eq 2 ]; then
    if [ "$1" = "--sweep" ]; then
        # Scan any file that doesn't already have a scan log
        echo "Starting sweep..."
        fd . ${WATCH_DIR} --type=file | while read FILE
        do
            $(realpath "$0") --check-and-scan "${FILE}"
        done
        echo "Sweep complete"
    fi
    if [ "$1" = "--check-and-scan" ]; then
        # Scan the file if it hasn't been already
        echo "Checking $2..."
        check_and_scan "$2"
        echo "Check complete"
    fi
    if [ "$1" = "--watch" ]; then
        # Watch the directory for new files

        # Need to create a file to check the the detection and scanning is working
        # This file needs to be created AFTER the filesystem watcher is started
        # But the inotifywait loop captures the script and no commands run after it
        # Thus the need to 'schedule' the file creation
        $(sleep 10; touch "$WATCH_DIR/selftest-$(date +"%Y-%m-%dT%H-%M-%S").txt") &

        echo "Watching for new files..."
        inotifywait -m -r -e create -e moved_to --format '%w%f' "$WATCH_DIR" | while read NEWFILE
        do
            echo "New file detected: $NEWFILE"
            scan "${NEWFILE}"
        done
        echo "Watching complete"
    fi
else
    echo "No command given"
fi
