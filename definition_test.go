package proxy

import (
	"io/ioutil"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateDefinition(t *testing.T) {
	url, err := url.Parse("http://localhost:2020")
	require.NoError(t, err)
	c := Configuration{
		Metadata: Metadata{
			Favicon:     "http://example.com/favicon.ico",
			Name:        "Example Search",
			Description: "Example description",
		},
		Host: url,
	}

	fixture, err := ioutil.ReadFile("testdata/definition.xml")
	require.NoError(t, err, "failed to load fixture")

	definition, err := generateDefinition(c)
	require.NoError(t, err, "failed to generate definition")

	assert.Equal(t, fixture, definition)
}
