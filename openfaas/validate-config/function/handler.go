package function

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var (
	infoLog  *log.Logger
	errorLog *log.Logger
)

// This is the main function which handles requests
func Handle(responseWriter http.ResponseWriter, request *http.Request) {
	configureLogging()
	err := validateRequest(request)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(err.Error() + "\n"))
		return
	}
	var commandStdout, commandStderr bytes.Buffer
	commandExitCode, err := executeCommand(&commandStdout, &commandStderr, responseWriter, request)
	handleCommandResult(&commandStdout, &commandStderr, commandExitCode, err, responseWriter, request)
}

// Configure info and error logging
func configureLogging() {
	infoLog = log.New(os.Stdout, "[INFO]  ", log.Lshortfile)
	errorLog = log.New(os.Stderr, "[ERROR] ", log.Lshortfile)
}

// Validate the request
func validateRequest(request *http.Request) error {
	if request.ContentLength == 0 {
		return errors.New("Request body is empty, no configuration provided to validate")
	}
	return nil
}

// Execute the command
func executeCommand(stdout *bytes.Buffer, stderr *bytes.Buffer, responseWriter http.ResponseWriter, request *http.Request) (int, error) {
	command := exec.Command("sourcehawk", "validate-config", "-")
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
			errorLog.Printf("Validate Config exited with code: %d", exitErr.ExitCode())
			return exitErr.ExitCode(), exitErr
		}
		errorMessage := fmt.Sprintf("Validate Config command did not exit cleanly: %s", err.Error())
		errorLog.Printf(errorMessage)
		return -1, errors.New(errorMessage)
	}
	infoLog.Printf("Command exited cleanly with code: %d", command.ProcessState.ExitCode())
	return command.ProcessState.ExitCode(), nil
}

// Handle the command result
func handleCommandResult(commandStdout *bytes.Buffer, commandStderr *bytes.Buffer, commandExitCode int, err error, responseWriter http.ResponseWriter, request *http.Request) {
	if commandExitCode == 1 || err == nil {
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Header().Set("Content-Type", "text/plain")
		responseWriter.Write([]byte(commandStdout.String()))
		return
	}
	errorLog.Printf("Command execution resulted in error with exit code: [%d]", commandExitCode)
	responseWriter.WriteHeader(http.StatusInternalServerError)
	responseWriter.Header().Set("Content-Type", "text/plain")
	responseWriter.Header().Set("Validate-Config-Exit-Code", fmt.Sprintf("%d", commandExitCode))
	responseWriter.Write([]byte(commandStderr.String() + "\n"))
}

// TODO: tests
