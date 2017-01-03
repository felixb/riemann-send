package main

import (
	"fmt"
	"strconv"
)

type Metric string

func NewMetric(defaultValue interface{}) Metric {
	if defaultValue == nil {
		return Metric("")
	} else {
		switch defaultValue.(type) {
		case int:
			return Metric(fmt.Sprintf("%d", defaultValue))
		case float32:
			return Metric(fmt.Sprintf("%f", defaultValue))
		case float64:
			return Metric(fmt.Sprintf("%f", defaultValue))
		default:
			return Metric("")
		}
	}
}

func (n *Metric) Parse() (interface{}, error) {
	if *n == "" {
		return nil, nil
	}
	if i, err := strconv.ParseInt(n.String(), 10, 32); err == nil {
		return i, nil
	}
	if f, err := strconv.ParseFloat(n.String(), 64); err == nil {
		return f, nil
	}
	return nil, fmt.Errorf("'%s' is not a number", n.String())
}

func (n *Metric) Set(s string) error {
	*n = Metric(s)
	_, err := n.Parse()
	return err
}

func (n *Metric) String() string {
	return string(*n)
}
