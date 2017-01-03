package main

import (
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/amir/raidman"
)

func Test_parseArgs(t *testing.T) {
	// you can only test once. flag is a singleton \o/

	s := ServerConfig{}
	e := raidman.Event{
		Metric: 1,
		Tags: []string{"some-tag", "some-other-tag"},
		Attributes: map[string]string{"some-key": "some-value", "some-other-key": "some-other-value"},
	}

	os.Args = []string{"cmd", "-metric", "1.5", "-tags", "fanzy-tag,FANCY-PANTS", "-attributes", "fanzy-key:value,FANCY-PANDA:eucalyptus"}
	parseArgs(&s, &e)

	assert.Equal(t, float64(1.5), e.Metric)
	assert.Equal(t, []string{"some-tag", "some-other-tag", "fanzy-tag", "FANCY-PANTS"}, e.Tags)
	assert.Equal(t, map[string]string{"some-key": "some-value", "some-other-key": "some-other-value", "fanzy-key":"value", "FANCY-PANDA": "eucalyptus"}, e.Attributes)
}