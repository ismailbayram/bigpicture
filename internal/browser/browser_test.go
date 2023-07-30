package browser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBrowser(t *testing.T) {
	browser := NewBrowser(LangGo, nil, nil)
	assert.NotNil(t, browser)

	browser = NewBrowser(LangPy, nil, nil)
	assert.NotNil(t, browser)
}
