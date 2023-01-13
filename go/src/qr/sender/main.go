package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tuvirz/qr/go/src/qr/sender/sender"
)

func main() {
	defaultDstFile, err := defaultDstPath()
	if err != nil {
		fmt.Println("--- ERROR:", err)
	}

	srcFileName := flag.String("src", "big_file.txt", "source file path and name")
	dstFileName := flag.String("dst", defaultDstFile, "target file name")
	flag.Parse()

	err = sender.RunSender(*srcFileName, *dstFileName)
	if err != nil {
		fmt.Println("--- ERROR:", err)
	} else {
		fmt.Println("sender: file was transfered successfully.")
	}
}

func defaultDstPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home dir: %w", err)
	}
	return filepath.Join(homeDir, "Desktop", "qr_code.txt"), nil
}
