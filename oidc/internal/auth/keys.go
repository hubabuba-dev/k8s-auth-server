package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"

	"errors"
)

type KeyManager struct {
	keys    map[string]*rsa.PrivateKey
	current string
	mu      sync.RWMutex
}

type TokenSigner interface {
	GetCurrentKey() (*rsa.PrivateKey, string, error)
}

func CreateKeyManager(cron string) (*KeyManager, gocron.Scheduler, error) {
	km := &KeyManager{
		keys: make(map[string]*rsa.PrivateKey),
	}

	s, _ := gocron.NewScheduler()

	log.Println("first rotation")
	err := km.RotateKeys()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to genereate key: %w", err)
	}

	log.Println("Generate cronjob")
	j, err := s.NewJob(gocron.CronJob(cron, false), gocron.NewTask(km.RotateKeys))

	if err != nil {
		return nil, nil, fmt.Errorf("failed to genereate key: %w", err)
	}

	s.Start()

	log.Printf("Job created %v", j.ID())

	return km, s, nil
}

func (km *KeyManager) RotateKeys() error {
	log.Println("Start rotation")

	km.mu.Lock()

	defer km.mu.Unlock()

	newKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	keyID := fmt.Sprintf("oidc-key-%d", time.Now().Unix())

	km.keys[keyID] = newKey

	km.current = keyID

	if len(km.keys) > 3 {
		err := km.removeOldKeys()
		if err != nil {
			return err
		}
	}

	return nil
}

func (km *KeyManager) removeOldKeys() error {
	log.Println("Removing key")
	km.mu.Lock()

	defer km.mu.Unlock()
	oldestId := ""
	for id := range km.keys {
		if oldestId == "" || id < oldestId {
			oldestId = id
		}
	}
	if oldestId == "" {
		return errors.New("keys not found. nothing to delete")
	}
	delete(km.keys, oldestId)
	return nil
}

func (km *KeyManager) GetCurrentKey() (*rsa.PrivateKey, string, error) {
	if len(km.keys) == 0 {
		return nil, "", errors.New("keys not found")
	}

	if km.current == "" {
		return nil, "", errors.New("current undefiended")
	}

	return km.keys[km.current], km.current, nil
}
