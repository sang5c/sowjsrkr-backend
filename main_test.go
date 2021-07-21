package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	split := strings.SplitN("1,2,3", ",", 2)
	for _, v := range split {
		fmt.Println(v)

	}
	fmt.Println(split)
}
