package main

import "strings"

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join([]string(*i), ", ")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
