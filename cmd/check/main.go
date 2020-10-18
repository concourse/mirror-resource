package main

import (
	"encoding/json"
	"os"
	"time"

	resource "github.com/concourse/mock-resource"
	"github.com/sirupsen/logrus"
)

type CheckRequest struct {
	Source  resource.Source   `json:"source"`
	Version *resource.Version `json:"version"`
}

type CheckResponse []resource.Version

func main() {
	logrus.SetOutput(os.Stderr)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	decoder := json.NewDecoder(os.Stdin)
	decoder.DisallowUnknownFields()

	var req CheckRequest
	err := decoder.Decode(&req)
	if err != nil {
		logrus.Fatalf("invalid payload: %s", err)
		return
	}

	if req.Source.Log != "" {
		logrus.Info(req.Source.Log)
	}

	if req.Source.CheckDelay != "" {
		delay, err := time.ParseDuration(req.Source.CheckDelay)
		if err != nil {
			logrus.Fatalf("malformed check_delay duration (%s): %s", req.Source.CheckDelay, err)
			return
		}

		time.Sleep(delay)
	}

	if req.Source.CheckFailure != "" {
		logrus.Fatalf("intentionally failing to check: %s", req.Source.CheckFailure)
		return
	}

	response := CheckResponse{}

	if req.Source.ForceVersion != "" {
		response = append(response, resource.Version{
			Version: req.Source.ForceVersion,
		})
	} else if req.Version != nil {
		response = append(response, *req.Version)
	} else if !req.Source.NoInitialVersion {
		response = append(response, resource.Version{
			Version: req.Source.InitialVersion(),
		})
	}

	privileged, err := resource.IsPrivileged()
	if err != nil {
		logrus.Fatalf("could not check privilege: %s", err)
		return
	}

	if privileged {
		for i := range response {
			// must be updated inline
			response[i].Privileged = "true"
		}
	}

	json.NewEncoder(os.Stdout).Encode(response)
}
