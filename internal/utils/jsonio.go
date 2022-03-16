package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"golang.org/x/sys/unix"
)

// ReadFile reads the cache file from disk
func ReadJSONFile(path string, target interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		if err == unix.ENOENT {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				os.Mkdir(path, 0777)
			}
		}
		return err
	}

	err = json.Unmarshal(data, target)
	if err != nil {
		return err
	}

	return nil
}

// WriteFile writes the cache file to disk
func WriteJSONFile(path string, payload interface{}) error {
	data, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
