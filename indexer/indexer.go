package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Set the Maildir directory
	maildir := "../enron_mail_20110402/maildir"

	// Create a slice to store the emails
	emails := make([]string, 0)

	fmt.Println("Looking for emails in " + maildir)

	// Walk the Maildir directory tree and list the files
	err := filepath.Walk(maildir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			emails = append(emails, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(len(emails), "emails found")
}
