package tests

import (
	"github.com/bregydoc/freya/freyacon/go"
	a "github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestNewFreya(t *testing.T) {
	assert := a.New(t)
	log.Println(t.Name())
	f, err := freya.NewFreya(&freya.FreyaConfig{})
	assert.NotNil(f, "Freya not support empty config")
	assert.Nil(err, "Error creating freya empty config")
}
