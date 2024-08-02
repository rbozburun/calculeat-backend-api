package controller

import (
	"testing"

	"github.com/calculeat/main_rest_api/config"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	a := assert.New()

	// Setup DB
	config.Connect()

}
