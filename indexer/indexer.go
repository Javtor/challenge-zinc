package main

import (
	"fmt"
	"io/ioutil"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type EmailJson struct {
	Header map[string][]string `json:"header"`
	Body   string              `json:"body"`
}

func emailPaths(maildir string) ([]string, error) {
	// Create a slice to store the emails
	emails := make([]string, 0)

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
		return nil, err
	}

	return emails, nil
}

func parseEmail(filePath string) (*EmailJson, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	msg, err := mail.ReadMessage(file)
	if err != nil {
		if strings.Contains(err.Error(), "malformed MIME header") {
			// handle malformed header, maybe just ignoring the email
			return nil, fmt.Errorf("email %s ignored due to malformed headers", filePath)
		} else {
			return nil, err
		}
	}
	body, _ := ioutil.ReadAll(msg.Body)
	emailJson := &EmailJson{
		Header: msg.Header,
		Body:   string(body),
	}
	return emailJson, nil
}

func emailMessages(maildir string) ([]*EmailJson, error) {
	paths, err := emailPaths(maildir)
	if err != nil {
		return nil, err
	}
	paths = paths[:100000]
	var messages []*EmailJson
	for _, path := range paths {
		msg, err := parseEmail(path)
		if err != nil {
			if strings.Contains(err.Error(), "ignored due to malformed headers") {
				// skip emails that were ignored due to malformed headers
				continue
			} else {
				return nil, err
			}
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func main() {
	fmt.Println("Starting...")
	start := time.Now()
	msgs, err := emailMessages("../enron_mail_20110402/maildir")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(msgs), "emails found and parsed")
	elapsed := time.Since(start)
	fmt.Printf("Time taken: %s\n", elapsed)
}
