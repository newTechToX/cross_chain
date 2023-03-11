package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestCmd(t *testing.T) {
	from_block := "16784638"
	to_block := "16786692"
	cmd := exec.Command("cd ../../")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	cmd = exec.Command("pwd")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	data, err := cmd.Output()
	fmt.Println(string(data))

	cmd = exec.Command(". ../../bin/pro", "-name", "anyswap", "-from", from_block, "-to", to_block, "-chain", "eth")
	data, err = cmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}
