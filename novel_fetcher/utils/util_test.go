package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateURL(t *testing.T) {
	provider := []struct {
		parts  []string
		expect string
	}{
		{
			[]string{"http://www.baidu.com", "/data"},
			"http://www.baidu.com/data",
		},
		{
			[]string{"http://www.baidu.com", "/data", "inter"},
			"http://www.baidu.com/datainter",
		},
	}
	for _, col := range provider {
		res := GenerateURL(col.parts...)
		assert.Equal(t, col.expect, res)
	}
}
