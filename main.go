package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/yargevad/filepathx"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var version string

func main() {
	fmt.Printf("* exim-spool-helper v%s - https://github.com/rwasef1830/exim-spool-helper\n", version)

	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get working dir\n")
		os.Exit(255)
	}

	spoolFileSpec := flag.String("input", workingDir, "input folder to recurse through or spool file")
	backupDir := flag.String("backup", "", "backup folder")
	okToOverwrite := flag.Bool("overwrite", false, "confirm overwriting conflicting backup files")

	if backupStat, err := os.Stat(*backupDir); err == nil {
		if !backupStat.IsDir() {
			fmt.Printf("Backup dir exists and is a file. Cannot proceed.\n")
			os.Exit(255)
		} else if !*okToOverwrite {
			fmt.Printf("Backup dir exists! Rename backup dir or relaunch with overwrite flag to confirm.\n")
			os.Exit(255)
		} else {
			fmt.Printf("Backup dir exists! Conflicting files will be overwritten!\n")
		}
	} else if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(*backupDir, 750); err != nil {
			fmt.Printf("Failed to create backup dir '%s'. Cannot proceed: %s\n", *backupDir, err)
			os.Exit(255)
		}
	} else {
		fmt.Printf("Failed to stat backup dir '%s'. Cannot proceed: %s\n", *backupDir, err)
		os.Exit(255)
	}

	var spoolFileStat os.FileInfo
	if spoolFileStat, err = os.Stat(*spoolFileSpec); err != nil && !os.IsNotExist(err) {
		fmt.Printf("Failed to stat spool filespec '%s'. Cannot proceed: %s\n", *spoolFileSpec, err)
		os.Exit(255)
	}

	var spoolMatches []string
	if spoolFileStat != nil {
		spoolMatches = []string{*spoolFileSpec}
	} else if spoolMatches, err = filepathx.Glob(filepath.Join(*spoolFileSpec, "**", "*-H")); err != nil {
		fmt.Printf("Failed to glob filespec '%s'. Cannot proceed: %s\n", *spoolFileSpec, err)
		os.Exit(255)
	}

	var count int
	for i, spoolFilePath := range spoolMatches {
		spoolFileRelativePath := strings.Replace(spoolFilePath, *spoolFileSpec, "", 1)
		spoolFileBackupFilePath := filepath.Join(*backupDir, spoolFileRelativePath)
		spoolFileBackupDirPath := filepath.Dir(spoolFileBackupFilePath)

		if err := os.MkdirAll(spoolFileBackupDirPath, 0755); err != nil {
			fmt.Printf("Failed to create backup path '%s'. Skipping '%s'. Error: %s\n",
				spoolFileBackupDirPath,
				spoolFilePath,
				err)
			continue
		}

		if err := copyFile(spoolFilePath, spoolFileBackupFilePath); err != nil {
			fmt.Printf("Failed to copy spool file '%s'. Skipping. Error: %s\n",
				spoolFilePath,
				err)
			continue
		}

		if err := processSpoolHeaderFile(spoolFilePath); err != nil {
			fmt.Printf("Failed to process spool file '%s'. Skipping. Error: %s\n",
				spoolFilePath,
				err)
			continue
		}

		fmt.Printf("%s\n", spoolFilePath)
		count = i + 1
	}

	if count == 0 {
		fmt.Printf("* No spool header files were found.\n")
		os.Exit(0)
	} else {
		fmt.Printf("* Processed %d spool header file(s).\n", count)
		os.Exit(0)
	}
}

func copyFile(srcPath string, dstPath string) error {
	// open files r and w
	r, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer func(r *os.File) {
		_ = r.Close()
	}(r)

	w, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer func(w *os.File) {
		_ = w.Close()
	}(w)

	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}

	return nil
}
