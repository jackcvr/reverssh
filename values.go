package main

import (
	"fmt"
	"strconv"
	"strings"
)

type BindAddress []string

func (b *BindAddress) String() string {
	return strings.Join(*b, ", ")
}

func (b *BindAddress) Set(value string) error {
	*b = append(*b, value)
	return nil
}

type Ports []int

func (p *Ports) String() string {
	return fmt.Sprintf("%v", *p)
}

func (p *Ports) Set(s string) error {
	*p = Ports{}
	ps := strings.Split(s, ",")
	for _, v := range ps {
		port, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		*p = append(*p, port)
	}
	return nil
}
