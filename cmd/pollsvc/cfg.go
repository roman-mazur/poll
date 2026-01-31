package main

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

type talkConfig struct {
	mu     sync.Mutex
	talkId string
}

func (tc *talkConfig) Setup(name string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	suffix := sha256.Sum256([]byte(time.Now().String()))
	tc.talkId = name + "-" + hex.EncodeToString(suffix[:])
}

func (tc *talkConfig) CurrentId() string {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	return tc.talkId
}
