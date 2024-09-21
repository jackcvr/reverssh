package main

import (
	"fmt"
	"strconv"
	"strings"
)

type StringList []string

func (l *StringList) String() string {
	return strings.Join(*l, ",")
}

func (l *StringList) Set(value string) error {
	*l = append(*l, value)
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
