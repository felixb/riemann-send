package main

import (
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/amir/raidman"
)

func TestMetric_parseTags_empty(t *testing.T) {
	e := raidman.Event{}
	tags := ""
	parseTags(&e, &tags)
	assert.Nil(t, e.Tags)
}

func TestMetric_parseTags_empty_append(t *testing.T) {
	e := raidman.Event{
		Tags: []string{"some", "tag"},
	}
	tags := ""
	parseTags(&e, &tags)
	assert.Equal(t, []string{"some", "tag"}, e.Tags)
}

func TestMetric_parseTags_args(t *testing.T) {
	e := raidman.Event{}
	tags := "some,other"
	parseTags(&e, &tags)
	assert.Equal(t, []string{"some", "other"}, e.Tags)
}

func TestMetric_parseTags_args_append(t *testing.T) {
	e := raidman.Event{
		Tags: []string{"some", "tag"},
	}
	tags := "some,other"
	parseTags(&e, &tags)
	assert.Equal(t, []string{"some", "tag", "other"}, e.Tags)
}

func TestMetric_parseAttributes_invalid(t *testing.T) {
	e := raidman.Event{}
	attributes := "invalid"
	assert.Error(t, parseAttributes(&e, &attributes))
}

func TestMetric_parseAttributes_empty(t *testing.T) {
	e := raidman.Event{}
	attributes := ""
	assert.Nil(t, parseAttributes(&e, &attributes))
	assert.Nil(t, e.Attributes)
}

func TestMetric_parseAttributes_empty_append(t *testing.T) {
	e := raidman.Event{
		Attributes: map[string]string{"foo": "bar"},
	}
	attributes := ""
	assert.Nil(t, parseAttributes(&e, &attributes))
	assert.Equal(t, map[string]string{"foo": "bar"}, e.Attributes)
}

func TestMetric_parseAttributes_args(t *testing.T) {
	e := raidman.Event{}
	attributes := "key:value,keyyy:valueee"
	assert.Nil(t, parseAttributes(&e, &attributes))
	assert.Equal(t, map[string]string{"key": "value", "keyyy": "valueee"}, e.Attributes)
}

func TestMetric_parseAttributes_args_append(t *testing.T) {
	e := raidman.Event{
		Attributes: map[string]string{"foo": "bar"},
	}
	attributes := "key:value,foo:barrr"
	assert.Nil(t, parseAttributes(&e, &attributes))
	assert.Equal(t, map[string]string{"foo": "barrr", "key": "value"}, e.Attributes)
}

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