package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	fmt.Println(len(paths), "being parsed")
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

func uploadEmails(maildir string) error {
	emails, err := emailMessages(maildir)
	if err != nil {
		return err
	}

	// Create the payload in the correct format
	payload := map[string]interface{}{
		"index":   "email",
		"records": emails,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Send the POST request
	url := "http://localhost:4080/api/_bulkv2"
	req, err := http.NewRequest("POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth("admin", "Complexpass#123")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Print the response status and body
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	responseBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(responseBody))

	return nil
}

func main() {
	fmt.Println("Starting...")
	start := time.Now()
	maildir := "../enron_mail_20110402/maildir"
	err := uploadEmails(maildir)
	if err != nil {
		fmt.Println("Error uploading emails: ", err)
		return
	}
	fmt.Println("Emails uploaded successfully")
	elapsed := time.Since(start)
	fmt.Printf("Time taken: %s\n", elapsed)
}
