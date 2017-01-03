package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewMetric_nil(t *testing.T) {
	m := NewMetric(nil)
	v, err := m.Parse()
	assert.Nil(t, err)
	assert.Nil(t, v)
}

func TestNewMetric_0(t *testing.T) {
	m := NewMetric(0)
	v, err := m.Parse()
	assert.Nil(t, err)
	assert.Equal(t, int64(0), v)
}

func TestNewMetric_1(t *testing.T) {
	m := NewMetric(1)
	v, err := m.Parse()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), v)
}

func TestNewMetric_1_5(t *testing.T) {
	m := NewMetric(1.5)
	v, err := m.Parse()
	assert.Nil(t, err)
	assert.Equal(t, float64(1.5), v)
}

func TestMetric_Set_0(t *testing.T) {
	m := NewMetric(nil)
	m.Set("0")
	v, err := m.Parse()
	assert.Nil(t, err)
	assert.Equal(t, int64(0), v)
}

func TestMetric_Set_1(t *testing.T) {
	m := NewMetric(nil)
	m.Set("1")
	v, err := m.Parse()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), v)
}

func TestMetric_Set_1_5(t *testing.T) {
	m := NewMetric(nil)
	m.Set("1.5")
	v, err := m.Parse()
	assert.Nil(t, err)
	assert.Equal(t, float64(1.5), v)
}