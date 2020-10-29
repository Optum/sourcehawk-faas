package function

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type ScanRequest struct {
	githubApiUrl    string
	githubAuthToken string
	githubOrg       string
	githubRepo      string
	githubRef       string
	outputFormat    string
}

var (
	infoLog  *log.Logger
	errorLog *log.Logger
)

// This is the main function which handles requests
func Handle(responseWriter http.ResponseWriter, request *http.Request) {
	configureLogging()
	scanRequest, err := validateAndParseRequest(request)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(err.Error() + "\n"))
		return
	}
	var commandStdout, commandStderr bytes.Buffer
	commandExitCode, err := executeScanCommand(*scanRequest, &commandStdout, &commandStderr, responseWriter, request)
	handleScanCommandResult(&commandStdout, &commandStderr, commandExitCode, err, responseWriter, request)
}

// Configure info and error logging
func configureLogging() {
	infoLog = log.New(os.Stdout, "[INFO]  ", log.Lshortfile)
	errorLog = log.New(os.Stderr, "[ERROR] ", log.Lshortfile)
}

// Validate and parse the request
func validateAndParseRequest(request *http.Request) (*ScanRequest, error) {
	path := request.URL.Path
	if strings.HasPrefix(path, "/") {
		runes := []rune(path)
		path = string(runes[1:])
	}
	pathPieces := strings.Split(path, "/")
	if len(pathPieces) < 2 {
		errorLog.Printf("Request URL is invalid: %s", path)
		return nil, errors.New("Request URL must contain Github Organization and Repository")
	}
	authToken := request.Header.Get("Authorization")
	var output string
	if request.Header.Get("Accept") == "application/json" {
		output = "JSON"
	} else if request.Header.Get("Accept") == "text/markdown" {
		output = "MARKDOWN"
	} else {
		output = "TEXT"
	}
	var githubApiUrl string
	if request.Header.Get("Github-API-URL") == "" {
		githubApiUrl = "https://api.github.com/v3"
	} else {
		githubApiUrl = request.Header.Get("Github-API-URL")
	}
	ref := "main"
	if len(pathPieces) > 2 {
		ref = pathPieces[2]
	}
	scanRequest := ScanRequest{
		githubApiUrl:    githubApiUrl,
		githubAuthToken: authToken,
		githubOrg:       pathPieces[0],
		githubRepo:      pathPieces[1],
		githubRef:       ref,
		outputFormat:    output,
	}
	return &scanRequest, nil
}

// Execute the scan command
func executeScanCommand(scanRequest ScanRequest, stdout *bytes.Buffer, stderr *bytes.Buffer, responseWriter http.ResponseWriter, request *http.Request) (int, error) {
	command := exec.Command("./scan.sh")
	command.Env = createScanCommandEnvironment(scanRequest)
	command.Stdin = request.Body
	command.Stdout = stdout
	command.Stderr = stderr
	if err := command.Start(); err != nil {
		errorMessage := fmt.Sprintf("Unable to start command: %s", err.Error())
		errorLog.Printf(errorMessage)
		return -1, errors.New(errorMessage)
	}
	if err := command.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			errorLog.Printf("Scan exited with code: %d", exitErr.ExitCode())
			return exitErr.ExitCode(), exitErr
		}
		errorMessage := fmt.Sprintf("Scan command did not exit cleanly: %s", err.Error())
		errorLog.Printf(errorMessage)
		return -1, errors.New(errorMessage)
	}
	infoLog.Printf("Command exited cleanly with code: %d", command.ProcessState.ExitCode())
	return command.ProcessState.ExitCode(), nil
}

// Create the scan command environment
func createScanCommandEnvironment(scanRequest ScanRequest) []string {
	commandEnv := os.Environ()
	commandEnv = append(commandEnv, fmt.Sprintf("GITHUB_API_URL=%s", scanRequest.githubApiUrl))
	commandEnv = append(commandEnv, fmt.Sprintf("GITHUB_ORG=%s", scanRequest.githubOrg))
	commandEnv = append(commandEnv, fmt.Sprintf("GITHUB_REPO=%s", scanRequest.githubRepo))
	commandEnv = append(commandEnv, fmt.Sprintf("GITHUB_REF=%s", scanRequest.githubRef))
	commandEnv = append(commandEnv, fmt.Sprintf("OUTPUT_FORMAT=%s", scanRequest.outputFormat))
	if scanRequest.githubAuthToken != "" {
		githubAuthToken := scanRequest.githubAuthToken
		if strings.HasPrefix(scanRequest.githubAuthToken, "token ") {
			runes := []rune(scanRequest.githubAuthToken)
			githubAuthToken = string(runes[6:])
		}
		commandEnv = append(commandEnv, fmt.Sprintf("GITHUB_AUTH_TOKEN=%s", githubAuthToken))
	}
	return commandEnv
}

// Handle the scan command result
func handleScanCommandResult(commandStdout *bytes.Buffer, commandStderr *bytes.Buffer, commandExitCode int, err error, responseWriter http.ResponseWriter, request *http.Request) {
	if commandExitCode == 1 || err == nil {
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Header().Set("Content-Type", request.Header.Get("Accept"))
		responseWriter.Write([]byte(commandStdout.String()))
		return
	}
	handleScanCommandError(commandExitCode, commandStderr, responseWriter, request)
}

// Handle the scan command error
func handleScanCommandError(commandExitCode int, commandStderr *bytes.Buffer, responseWriter http.ResponseWriter, request *http.Request) {
	errorLog.Printf("Command execution resulted in error with exit code: [%d]", commandExitCode)
	if commandExitCode == 61 {
		responseWriter.WriteHeader(http.StatusUnauthorized)
	} else if commandExitCode == 64 {
		responseWriter.WriteHeader(http.StatusNotFound)
	} else if commandExitCode == 65 {
		responseWriter.WriteHeader(http.StatusInternalServerError)
	} else {
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}
	responseWriter.Header().Set("Content-Type", "text/plain") // TODO: dynamic
	responseWriter.Header().Set("Scan-Exit-Code", fmt.Sprintf("%d", commandExitCode))
	responseWriter.Write([]byte(commandStderr.String() + "\n"))
}

// TODO: tests
