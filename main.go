package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const ShellToUse = "/opt/homebrew/bin/zsh"

type Device struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type AlfredListItem struct {
	Title        string `json:"title"`
	Arg          string `json:"arg"`
	AutoComplete string `json:"autocomplete"`
}

type AlfredPayload struct {
	Items []AlfredListItem `json:"items"`
}

func Shellout(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func main() {
	out, _, err := Shellout("SwitchAudioSource -a -t output -f json")
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	data := strings.Split(out, "\n")
	devices := make([]AlfredListItem, len(data))

	for index, device := range data {
		if device == "" {
			continue
		}
		var d Device
		err = json.Unmarshal([]byte(device), &d)
		if err != nil {
			log.Printf("error: %v\n", err)
		}
		devices[index] = AlfredListItem{Title: d.Name, AutoComplete: d.Name, Arg: fmt.Sprintf("%s,%s", d.ID, d.Name)}
	}

	payload, err := json.Marshal(AlfredPayload{Items: devices})

	if err != nil {
		log.Printf("error: %v\n", err)
	}

	fmt.Print(string(payload))

}
