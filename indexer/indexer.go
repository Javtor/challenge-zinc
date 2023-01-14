package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	BATCH_SIZE = 10000
	API_URL    = "http://localhost:4080/api/_bulkv2"
	INDEX      = "email"
)

type EmailJson struct {
	Header map[string][]string `json:"header"`
	Body   string              `json:"body"`
}

// emailPaths returns a list of email file paths in a given directory
func emailPaths(maildir string) ([]string, error) {
	// Create a slice to store the emails
	var emails []string

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

// parseEmail reads an email file and returns an EmailJson struct
func parseEmail(filePath string) (*EmailJson, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file %s: %v", filePath, err)
		return nil, err
	}
	defer file.Close()

	msg, err := mail.ReadMessage(file)
	if err != nil {
		if strings.Contains(err.Error(), "malformed MIME header") {
			log.Printf("Email %s ignored due to malformed headers", filePath)
			return nil, fmt.Errorf("email %s ignored due to malformed headers", filePath)
		} else {
			log.Printf("Error reading email %s: %v", filePath, err)
			return nil, err
		}
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, msg.Body)

	if err != nil {
		return nil, err
	}
	emailJson := &EmailJson{
		Header: msg.Header,
		Body:   buf.String(),
	}
	return emailJson, nil
}

// createBatches returns a 2D slice of strings, with each inner slice representing a batch of file paths
func createBatches(paths []string, batchSize int) ([][]string, error) {
	numBatches := len(paths) / batchSize
	if len(paths)%batchSize != 0 {
		numBatches++
	}

	batches := make([][]string, numBatches)
	for i := 0; i < numBatches; i++ {
		start := i * batchSize
		end := start + batchSize
		if end > len(paths) {
			end = len(paths)
		}
		batches[i] = paths[start:end]
	}
	return batches, nil
}

// processBatch reads and parses a batch of email files, then uploads the parsed data to the server
func processBatch(batch []string) error {
	log.Println("Starting processing of batch")
	messages, err := parseEmailBatch(batch)
	if err != nil {
		return err
	}
	log.Println("Finishing processing of batch")
	return uploadBatch(API_URL, messages)
}

// parseEmailBatch reads and parses a batch of email files, returning a slice of EmailJson structs
func parseEmailBatch(batch []string) ([]*EmailJson, error) {
	messages := make([]*EmailJson, len(batch))
	i := 0
	for _, path := range batch {
		msg, err := parseEmail(path)
		if err != nil {
			if strings.Contains(err.Error(), "ignored due to malformed headers") {
				continue
			} else {
				log.Printf("Error parsing email %s: %v", path, err)
				return nil, err
			}
		}
		messages[i] = msg
		i++
	}
	//remove nil entries
	messages = messages[:i]

	return messages, nil
}

// uploadBatch uploads a batch of email data to the server
func uploadBatch(url string, batch []*EmailJson) error {
	payload := map[string]interface{}{
		"index":   INDEX,
		"records": batch,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth("admin", "Complexpass#123")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error uploading batch: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// processMaildir reads email files from a given maildir and uploads them to the server
func processMaildir(maildir string) error {
	paths, err := emailPaths(maildir)
	if err != nil {
		return err
	}
	log.Println(len(paths), "emails found.")

	batches, err := createBatches(paths, BATCH_SIZE)
	if err != nil {
		return err
	}

	totalBatches := len(batches)
	log.Printf("Processing %d batches\n", totalBatches)

	startTime := time.Now()
	lastBatchStart := startTime
	for i, batch := range batches {
		log.Printf("Processing batch %d of %d...", i+1, totalBatches)

		if i > 0 {
			elapsed := time.Since(startTime)
			lastBatchTook := time.Since(lastBatchStart)
			estimatedRemaining := elapsed / time.Duration(i) * time.Duration(totalBatches-i)
			log.Printf("Last batch: %v, Elapsed time: %v, Estimated time remaining: %v\n", lastBatchTook, elapsed, estimatedRemaining)
		} else {
			log.Println()
		}

		lastBatchStart = time.Now()
		if err := processBatch(batch); err != nil {
			return err
		}

	}
	log.Println("Process completed.")
	return nil
}

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <maildir>\n", os.Args[0])
	}
	flag.Parse()
	maildir := flag.Arg(0)
	if maildir == "" {
		fmt.Println("Error: No maildir path provided")
		flag.Usage()
		os.Exit(1)
	}

	log.Println("Starting the program")
	start := time.Now()
	err := processMaildir(maildir)
	if err != nil {
		log.Println("Error processing maildir: ", err)
		os.Exit(1)
	}

	log.Println("Maildir processed successfully")
	elapsed := time.Since(start)
	log.Printf("Time taken: %s\n", elapsed)
}
