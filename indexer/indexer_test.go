package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestEmailPaths(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, _ := ioutil.TempDir("", "maildir")
	defer os.RemoveAll(tempDir)

	// Create some test email files in the directory
	ioutil.WriteFile(tempDir+"/email1.txt", []byte("test email 1"), 0644)
	ioutil.WriteFile(tempDir+"/email2.txt", []byte("test email 2"), 0644)

	// Test emailPaths function
	emails, err := emailPaths(tempDir)
	if err != nil {
		t.Errorf("emailPaths returned an error: %v", err)
	}
	if len(emails) != 2 {
		t.Errorf("emailPaths did not return the correct number of emails. Got: %d, expected: 2", len(emails))
	}
}

func TestParseEmail(t *testing.T) {
	// Create a test email file
	tempFile, _ := ioutil.TempFile("", "email")
	defer os.Remove(tempFile.Name())
	ioutil.WriteFile(tempFile.Name(), []byte("From: test@example.com\nSubject: Test Email\n\nThis is a test email"), 0644)

	// Test parseEmail function
	email, err := parseEmail(tempFile.Name())
	if err != nil {
		t.Errorf("parseEmail returned an error: %v", err)
	}
	if email.Header["From"][0] != "test@example.com" {
		t.Errorf("parseEmail did not return the correct 'From' header. Got: %s, expected: test@example.com", email.Header["From"][0])
	}
	if email.Header["Subject"][0] != "Test Email" {
		t.Errorf("parseEmail did not return the correct 'Subject' header. Got: %s, expected: Test Email", email.Header["Subject"][0])
	}
	if email.Body != "This is a test email" {
		t.Errorf("parseEmail did not return the correct email body. Got: %s, expected: This is a test email", email.Body)
	}
}

func TestCreateBatches(t *testing.T) {
	paths := []string{"path1", "path2", "path3", "path4", "path5"}

	// Test createBatches function
	batches, err := createBatches(paths, 5)
	if err != nil {
		t.Errorf("createBatches returned an error: %v", err)
	}
	if len(batches) != 1 {
		t.Errorf("createBatches did not return the correct number of batches. Got: %d, expected: 1", len(batches))
	}
	if len(batches[0]) != 5 {
		t.Errorf("createBatches did not return the correct number of paths in the first batch. Got: %d, expected: 5", len(batches[0]))
	}
}

func TestParseEmailBatch(t *testing.T) {
	// Create a batch of test email files
	tempDir, _ := ioutil.TempDir("", "maildir")
	defer os.RemoveAll(tempDir)
	ioutil.WriteFile(tempDir+"/email1.txt", []byte("From: test@example.com\nSubject: Test Email 1\n\ntest email 1"), 0644)
	ioutil.WriteFile(tempDir+"/email2.txt", []byte("From: test@example.com\nSubject: Test Email 2\n\ntest email 2"), 0644)
	batch := []string{tempDir + "/email1.txt", tempDir + "/email2.txt"}

	// Test parseEmailBatch function
	emails, err := parseEmailBatch(batch)
	if err != nil {
		t.Errorf("parseEmailBatch returned an error: %v", err)
	}
	if len(emails) != 2 {
		t.Errorf("parseEmailBatch did not return the correct number of emails. Got: %d, expected: 2", len(emails))
	}
	if emails[0].Body != "test email 1" {
		t.Errorf("parseEmailBatch did not return the correct email body for email 1. Got: %s, expected: test email 1", emails[0].Body)
	}
	if emails[1].Body != "test email 2" {
		t.Errorf("parseEmailBatch did not return the correct email body for email 2. Got: %s, expected: test email 2", emails[1].Body)
	}
}

func TestEmailPathsInvalidDir(t *testing.T) {
	// Pass an invalid directory path to the emailPaths function
	emails, err := emailPaths("/invalid/directory")

	// Check if the function returns an error
	if err == nil {
		t.Errorf("emailPaths did not return an error for an invalid directory")
	}

	// Check if the function returns an empty list of email files
	if emails != nil {
		t.Errorf("emailPaths returned a non-nil list of emails for an invalid directory")
	}
}

func TestEmailPathsEmptyDir(t *testing.T) {
	// Create a temporary empty directory
	tempDir, _ := ioutil.TempDir("", "maildir")
	defer os.RemoveAll(tempDir)

	// Pass the empty directory to the emailPaths function
	emails, err := emailPaths(tempDir)

	// Check if the function returns an error
	if err != nil {
		t.Errorf("emailPaths returned an error for an empty directory: %v", err)
	}

	// Check if the function returns an empty list of email files
	if len(emails) != 0 {
		t.Errorf("emailPaths returned a non-empty list of emails for an empty directory")
	}
}

func TestParseEmailMalformedFile(t *testing.T) {
	// Create a test email file with malformed headers
	tempFile, _ := ioutil.TempFile("", "email")
	defer os.Remove(tempFile.Name())
	ioutil.WriteFile(tempFile.Name(), []byte("From test@example.com\nSubject: Test Email\n\nThis is a test email"), 0644)

	// Test parseEmail function
	_, err := parseEmail(tempFile.Name())

	// Check if the function returns an error with a message indicating that the email file is malformed
	if err == nil || !strings.Contains(err.Error(), "malformed") {
		t.Errorf("parseEmail did not return an error for a malformed email file")
	}
}

func TestUploadBatchError(t *testing.T) {
	// Test uploadBatch function with an invalid server URL
	emails := []*EmailJson{
		{
			Header: map[string][]string{"From": {"test@example.com"}},
			Body:   "test email 1",
		},
	}
	url := "http://invalid.server"
	err := uploadBatch(url, emails)
	if err == nil {
		t.Error("uploadBatch did not return an error for an invalid server URL")
	}
}

func TestCreateBatchesLargeBatchSize(t *testing.T) {
	paths := []string{"path1", "path2"}

	// Test createBatches function
	batches, err := createBatches(paths, 5)
	if err != nil {
		t.Errorf("createBatches returned an error: %v", err)
	}

	// Check if the function returns a single batch containing all the email files
	if len(batches) != 1 {
		t.Errorf("createBatches did not return a single batch")
	}
	if len(batches[0]) != 2 {
		t.Errorf("createBatches did not return all the email files in the first batch. Got: %d, expected: 2", len(batches[0]))
	}
}

func TestParseEmailBatchMalformedHeaders(t *testing.T) {
	// Create a batch of test email files with malformed headers
	tempDir, _ := ioutil.TempDir("", "maildir")
	defer os.RemoveAll(tempDir)
	ioutil.WriteFile(tempDir+"/email1.txt", []byte("From: test@example.com\nSubject: Test Email 1\n\ntest email 1"), 0644)
	ioutil.WriteFile(tempDir+"/email2.txt", []byte("From: test@example.com\nSubject: Test Email 2\n\ntest email 2"), 0644)
	ioutil.WriteFile(tempDir+"/email3.txt", []byte("From test@example.com\nSubject: Test Email 3\n\ntest email 3"), 0644)
	batch := []string{tempDir + "/email1.txt", tempDir + "/email2.txt", tempDir + "/email3.txt"}

	// Test parseEmailBatch function
	emails, _ := parseEmailBatch(batch)

	// Check if the function returns the correct number of emails after ignoring the malformed headers
	if len(emails) != 2 {
		t.Errorf("parseEmailBatch did not return the correct number of emails after ignoring malformed headers. Got: %d, expected: 2", len(emails))
	}
}

func BenchmarkProcessMaildir(b *testing.B) {
	maildir := "../enron_mail_20110402/maildir"
	processMaildir(maildir)
}
