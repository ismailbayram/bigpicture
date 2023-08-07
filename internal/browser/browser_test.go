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

	browser = NewBrowser(LangJava, nil, nil)
	assert.NotNil(t, browser)
}

func TestIsIgnored(t *testing.T) {
	browser := &GoBrowser{
		ignoredPaths: []string{"internal/browser"},
		tree:         nil,
	}

	assert.True(t, isIgnored(browser.ignoredPaths, "./internal/browser/go.go"))
	assert.False(t, isIgnored(browser.ignoredPaths, "./internal/other/other.go"))
}
