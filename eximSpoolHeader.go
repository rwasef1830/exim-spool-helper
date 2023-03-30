package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var headerPattern, _ = regexp.Compile(`^\s*(?P<ByteCount>\d+)(?P<Flag>[A-Z*]?)\s*(?P<HeaderName>[^:]+):\s*(?P<HeaderValue>.*)`)

func processSpoolHeaderFile(filePath string) error {
	fileName := filepath.Base(filePath)
	writeFilePath := filePath + ".tmp"

	spoolFileHandle, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	defer func(spoolFileHandle *os.File) {
		_ = spoolFileHandle.Close()
	}(spoolFileHandle)

	spoolWriteFileHandle, err := os.OpenFile(writeFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Cannot open '%s' file for writing. Skipping. Error: %s", filePath, err)
		return nil
	}

	defer func(spoolWriteFileHandle *os.File) {
		_ = spoolWriteFileHandle.Close()
	}(spoolWriteFileHandle)

	if err := processSpoolHeaderLines(fileName, spoolFileHandle, spoolWriteFileHandle); err != nil {
		return err
	}

	if err := os.Rename(writeFilePath, fileName); err != nil {
		return err
	}

	return nil
}

func processSpoolHeaderLines(fileName string, reader io.Reader, writer io.Writer) error {
	validLines, err := getProcessedSpoolHeaderLines(fileName, reader)
	if err != nil {
		return err
	}

	var firstProcessed bool
	for _, validLine := range validLines {
		if !firstProcessed {
			firstProcessed = true
		} else {
			if _, err := fmt.Fprint(writer, "\n"); err != nil {
				return err
			}
		}

		if _, err := fmt.Fprint(writer, validLine); err != nil {
			return err
		}
	}

	return nil
}

func getProcessedSpoolHeaderContent(fileName string, reader io.Reader) (string, error) {
	validLines, err := getProcessedSpoolHeaderLines(fileName, reader)
	if err != nil {
		return "", err
	}

	return strings.Join(validLines, "\n"), nil
}

func getProcessedSpoolHeaderLines(fileName string, reader io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	var foundAnyContent bool
	var foundEmptyLine bool
	var foundAnyHeader bool

	var validLines []string
	var lastHeaderByteCount int64
	var lastHeaderFlag string
	var lastHeaderName string
	var lastHeaderValue string

	flushHeader := func() {
		lastHeaderByteCount = int64(len(lastHeaderName) + len(": ") + len(lastHeaderValue) + 1)

		var spacer string
		if lastHeaderFlag == "" {
			spacer = "  "
		} else {
			spacer = " "
		}

		newValidLine := fmt.Sprintf(
			"%03d%s%s%s: %s",
			lastHeaderByteCount,
			lastHeaderFlag,
			spacer,
			lastHeaderName,
			lastHeaderValue)
		validLines = append(validLines, newValidLine)
	}

	var count int
	for scanner.Scan() {
		count++

		logLine := func(message string) {
			fmt.Printf("%s:line %d: %s\n", fileName, count, message)
		}

		line := scanner.Text()
		if line == "" {
			if !foundEmptyLine && foundAnyContent {
				foundEmptyLine = true
				validLines = append(validLines, "")
				continue
			}

			if lastHeaderName == "" {
				logLine("ignoring extra empty line")
				continue
			}
		}

		foundAnyContent = true
		match := headerPattern.FindStringSubmatch(line)

		if len(match) == 0 {
			trimmedLine := strings.Trim(line, "\r\n\t ")

			if lastHeaderName != "" {
				if lastHeaderValue == "" {
					lastHeaderValue = trimmedLine
					continue
				}

				lastHeaderValue += fmt.Sprintf("\n\t%s", trimmedLine)
				continue
			}

			if !foundAnyHeader && foundEmptyLine {
				logLine("Ignoring pre-header area trash")
				continue
			}

			validLines = append(validLines, trimmedLine)
			continue
		}

		foundAnyHeader = true

		if !foundEmptyLine {
			logLine("correcting missing empty line before start of header block")
			validLines = append(validLines, "")
		}

		components := make(map[string]string)
		for i, name := range headerPattern.SubexpNames() {
			if i != 0 && name != "" {
				components[name] = strings.Trim(match[i], "\r\n\t ")
			}
		}

		if lastHeaderName != "" {
			flushHeader()
		}

		_, err := strconv.ParseInt(components["ByteCount"], 10, 32)
		if err != nil {
			logLine(fmt.Sprintf("ignoring corrupted header byte count: %s", components["ByteCount"]))
		}

		lastHeaderFlag = components["Flag"]
		lastHeaderName = strings.Trim(components["HeaderName"], "\r\n\t ")
		lastHeaderValue = strings.Trim(components["HeaderValue"], "\r\n\t ")
	}

	if lastHeaderName != "" {
		flushHeader()
	}

	return validLines, nil
}
