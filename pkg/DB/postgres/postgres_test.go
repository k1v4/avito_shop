package postgres

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMaxPoolSize(t *testing.T) {
	p := &Postgres{}

	MaxPoolSize(10)(p)

	assert.Equal(t, 10, p.maxPoolSize, "MaxPoolSize should set maxPoolSize to 10")
}

func TestConnAttempts(t *testing.T) {
	p := &Postgres{}

	ConnAttempts(5)(p)

	assert.Equal(t, 5, p.connAttempts, "ConnAttempts should set connAttempts to 5")
}

func TestConnTimeout(t *testing.T) {
	p := &Postgres{}

	timeout := 100 * time.Second
	ConnTimeout(timeout)(p)

	assert.Equal(t, timeout, p.connTimeout, "ConnTimeout should set connTimeout to 10 seconds")
}

func TestMultipleOptions(t *testing.T) {
	p := &Postgres{}

	MaxPoolSize(20)(p)
	ConnAttempts(3)(p)
	ConnTimeout(25 * time.Second)(p)

	assert.Equal(t, 20, p.maxPoolSize, "MaxPoolSize should set maxPoolSize to 20")
	assert.Equal(t, 3, p.connAttempts, "ConnAttempts should set connAttempts to 3")
	assert.Equal(t, 25*time.Second, p.connTimeout, "ConnTimeout should set connTimeout to 5 seconds")
}
