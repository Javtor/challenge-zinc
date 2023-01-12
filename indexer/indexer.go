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

func createBatches(paths []string) ([][]string, error) {
	const BATCH_SIZE = 1000

	numBatches := len(paths) / BATCH_SIZE
	if len(paths)%BATCH_SIZE != 0 {
		numBatches++
	}

	batches := make([][]string, numBatches)
	for i := 0; i < numBatches; i++ {
		start := i * BATCH_SIZE
		end := start + BATCH_SIZE
		if end > len(paths) {
			end = len(paths)
		}
		batches[i] = paths[start:end]
	}
	return batches, nil
}

func processBatch(batch []string) error {
	messages, err := parseEmailBatch(batch)
	if err != nil {
		return err
	}
	return uploadBatch(messages)
}

func parseEmailBatch(batch []string) ([]*EmailJson, error) {
	var messages []*EmailJson
	for _, path := range batch {
		msg, err := parseEmail(path)
		if err != nil {
			if strings.Contains(err.Error(), "ignored due to malformed headers") {
				continue
			} else {
				return nil, err
			}
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func uploadBatch(batch []*EmailJson) error {
	payload := map[string]interface{}{
		"index":   "email",
		"records": batch,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func processMaildir(maildir string) error {
	paths, err := emailPaths(maildir)
	if err != nil {
		return err
	}
	fmt.Println(len(paths), "emails found.")
	batches, err := createBatches(paths)
	if err != nil {
		return err
	}
	fmt.Printf("Processing %d batches\n", len(batches))
	for i, batch := range batches {
		fmt.Printf("Processing batch %d of %d\n", i+1, len(batches))
		if err := processBatch(batch); err != nil {
			return err
		}
	}
	fmt.Println("Process completed.")
	return nil
}

func main() {
	fmt.Println("Starting...")
	start := time.Now()
	maildir := "../enron_mail_20110402/maildir"
	err := processMaildir(maildir)
	if err != nil {
		fmt.Println("Error processing maildir: ", err)
		return
	}
	fmt.Println("Maildir processed successfully")
	elapsed := time.Since(start)
	fmt.Printf("Time taken: %s\n", elapsed)
}
