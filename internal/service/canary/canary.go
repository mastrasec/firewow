package canary

import (
	"strings"

	"github.com/google/uuid"
)

type Contract interface {
	AddToken(msg string) (newBody, token string)
	HasLeakage(msg, token string) (leaked bool)
	HandleLeakage(msg, token string) (newMsg string)
	removeToken(msg, token string) (NewMsg string)
}

type Canary struct {
}

func New() *Canary {
	return &Canary{}
}

func (c *Canary) AddToken(body string) (newBody, token string) {
	token = uuid.NewString()
	newBody = strings.Join([]string{token, body}, "")

	return newBody, token
}

func (c *Canary) HasLeakage(msg, token string) (leaked bool) {
	return strings.Contains(msg, token)
}

func (c *Canary) HandleLeakage(msg, token string) (newMsg string) {
	newMsg = c.removeToken(msg, token)

	// TODO: Add side effects.

	return newMsg
}

func (c *Canary) removeToken(msg, token string) (newMsg string) {
	before, after, _ := strings.Cut(msg, token)

	return strings.Join([]string{before, after}, "")
}
