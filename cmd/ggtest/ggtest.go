package main

import (
	"fmt"
	"github.com/jbrukh/ggit/util"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("go", "test", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Creating test case repos...\n")

	// build the test case repos
	err := util.CreateRepoTestCases()
	defer func() {
		fmt.Println("\nRemoving test case repos.")
		util.RemoveRepoTestCases()
	}()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Done.")

	fmt.Println("Running tests...\n")
	cmd.Run()
}
