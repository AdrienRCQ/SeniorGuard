package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"
)

type Heartbeat struct {
	AgentID      string    `json:"agent_id"`
	Hostname     string    `json:"hostname"`
	OS           string    `json:"os"`
	Architecture string    `json:"architecture"`
	AgentVersion string    `json:"agent_version"`
	SentAt       time.Time `json:"sent_at"`
}

type AgentConfig struct {
	AgentID string `json:"agent_id"`
}

const configFile = "agent.json"
const heartbeatInterval = 30 * time.Second //30 secondes

func generateAgentID() (string, error) {
	/*
		Generate an UID for the agent
	*/
	randomBytes := make([]byte, 8)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return "SG-" + hex.EncodeToString(randomBytes), nil
}

func loadOrCreateAgentID() (string, error) {
	/*
		Verify if agent config file exist and if it has agentID
	*/
	data, err := os.ReadFile(configFile)
	if err == nil {
		var config AgentConfig
		if err := json.Unmarshal(data, &config); err != nil { //unmarshall unencode json content
			return "", err
		}
		return config.AgentID, nil
	}
	if !os.IsNotExist(err) {
		return "", err
	}

	agentID, err := generateAgentID()
	if err != nil {
		return "", err
	}
	config := AgentConfig{
		AgentID: agentID,
	}
	data, err = json.MarshalIndent(config, "", " ") //encode to json with indentation
	if err != nil {
		return "", err
	}

	err = os.WriteFile(
		configFile,
		data,
		0600,
	)
	if err != nil {
		return "", err
	}

	return agentID, nil
}

func sendHeartbeat(agentID string) error {
	/*
		build the heartbeat request and send it to the console
	*/
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	heartbeat := Heartbeat{
		AgentID:      agentID,
		Hostname:     hostname,
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		AgentVersion: "0.1.1",
		SentAt:       time.Now().UTC(),
	}

	payload, err := json.Marshal(heartbeat) //encode to json
	if err != nil {
		return err
	}

	response, err := http.Post(
		"http://localhost:8000/api/v1/agents/heartbeat",
		"application/json",
		bytes.NewBuffer(payload), //body
	)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"server is unvailable : returned code %s",
			response.Status,
		)
	}
	return nil
}

func main() {
	agentID, err := loadOrCreateAgentID()
	if err != nil {
		fmt.Println("Unable to generate an UID : ", err)
		return
	}
	fmt.Println("SeniorGuard Agent")
	fmt.Println("Agent ID :", agentID)

	ticker := time.NewTicker(heartbeatInterval)

	defer ticker.Stop()
	for {
		err := sendHeartbeat(agentID)

		if err != nil {
			fmt.Println(
				"Heartbeat failed:",
				err,
			)
		} else {
			fmt.Println(
				time.Now().Format(time.RFC3339),
				"- heartbeat sent",
			)
		}

		<-ticker.C
	}
}
