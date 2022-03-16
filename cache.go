package main

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/coma-toast/ResumAPI/internal/utils"
	"github.com/coma-toast/ResumAPI/pkg/candidate"
	log "github.com/sirupsen/logrus"
)

// Cache is the cache for storing episode lists
type Cache struct {
	Name           string
	path           string
	candidateState map[string]candidate.Candidate
	mutex          *sync.Mutex
}

// Initialize the cache so it doesn't panic when trying to assign to the map when the map is nil
func (c *Cache) Initialize(path string) error {
	log.Debugf("Initializing caches")
	c.mutex = &sync.Mutex{}
	c.path = path
	cache.candidateState = make(map[string]candidate.Candidate)
	err := utils.ReadJSONFile(fmt.Sprintf("%s/candidate.json", c.path), &cache.candidateState)
	if err != nil {
		return err
	}

	return nil
}

// Set sets the cache
func (c *Cache) SetCandidate(key string, value candidate.Candidate) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	log.Debugf("Setting cache: %s\n", key)
	c.candidateState[key] = value
	err := utils.WriteJSONFile(fmt.Sprintf("%s/candidate.json", c.path), c.candidateState)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) AddCandidate(value candidate.Candidate) (int, error) {
	allCandidates := c.GetAllCandidates()
	id := len(allCandidates) + 1
	value.ID = id
	key := strconv.Itoa(id)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	log.Debugf("Setting cache: %s\n", key)
	c.candidateState[key] = value
	err := utils.WriteJSONFile(fmt.Sprintf("%s/candidate.json", c.path), c.candidateState)
	if err != nil {
		return id, err
	}

	return id, nil
}

// Set sets the cache
func (c *Cache) DeleteCandidate(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	log.Debugf("Deleting cache: %s\n", key)
	if c.IsSet(key) {
		delete(c.candidateState, key)
	}
	err := utils.WriteJSONFile(fmt.Sprintf("%s/candidate.json", c.path), c.candidateState)
	if err != nil {
		return err
	}

	return nil
}

// Get retrieves data from the cache
func (c *Cache) GetCandidate(key string) candidate.Candidate {
	log.Debugf("Getting cache: %s\n", key)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.candidateState[key]
}

func (c *Cache) GetCandidateByID(id int) candidate.Candidate {
	log.Debugf("Getting cache: %s\n", id)
	allCandidates := c.GetAllCandidates()
	for _, candidate := range allCandidates {
		if candidate.ID == id {
			return candidate
		}
	}
	return candidate.Candidate{}
}

// GetAll retrieves all data from the cache
func (c *Cache) GetAllCandidates() map[string]candidate.Candidate {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.candidateState
}

// IsSet determines if the candidate exists already
func (c *Cache) IsSet(key string) bool {
	_, ok := c.candidateState[key]

	return ok
}
