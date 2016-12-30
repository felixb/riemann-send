package main

import (
	"fmt"
	"strconv"
)

type Metric string

func (n *Metric) Parse() (interface{}, error) {
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
