package helpers

import (
	"testing"

	"github.com/calculeat/main_rest_api/config"
	"github.com/stretchr/testify/assert"
)

func TestSetColForIntVal(t *testing.T) {

	SetColForIntVal("users", "init_water_id", 1, 4)
}

func TestInitiliazeWaterObject(t *testing.T) {
	a := assert.New(t)
	config.Connect()

	_, err := InitiliazeWaterObject(6)
	if err != nil {
		a.Error(err)
	}

}
