package main

import (
	"fmt"
	"os/exec"
	"testing"
)

func TestCmd(t *testing.T) {
	from_block := "16784638"
	to_block := "16786790"

	cmd := exec.Command("../../bin/pro", "-name", "anyswap", "-from", from_block, "-to", to_block, "-chain", "eth")

	dd, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	if string(dd) == "" {
		println("nil")
	}
	fmt.Println(string(dd))
}
