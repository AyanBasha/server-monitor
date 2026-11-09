package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
)

const dataFile = "./data/services.json"

type savedData struct {
	NextID   int       `json:"next_id"`
	Services []Service `json:"services"`
}

func saveLocked() {
	data, err := json.MarshalIndent(savedData{NextID: nextID, Services: services}, "", "  ")
	if err != nil {
		log.Printf("save: encode failed: %v", err)
		return
	}

	if err := os.MkdirAll(filepath.Dir(dataFile), 0o755); err != nil {
		log.Printf("save: mkdir failed: %v", err)
		return
	}

	//all or nothing
	tmp := dataFile + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		log.Printf("save: write failed: %v", err)
		return
	}
	if err := os.Rename(tmp, dataFile); err != nil {
		log.Printf("save: rename failed: %v", err)
	}
}

func loadServices() error {
	data, err := os.ReadFile(dataFile)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	if err != nil {
		return err
	}

	var saved savedData
	if err := json.Unmarshal(data, &saved); err != nil {
		return err
	}

	servicesMu.Lock()
	defer servicesMu.Unlock()
	services = saved.Services
	nextID = saved.NextID
	return nil
}
