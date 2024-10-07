package main

import (
	"fmt"
	"github.com/supply-chain-tools/release-prototype/count"
	"os"
)

func main() {
	repoState, err := count.LoadRepoStateFromCurrentDirectory()
	if err != nil {
		print("failed to load repo: ", err.Error(), "\n")
		os.Exit(1)
	}

	commitCount := count.Commits(repoState)
	fmt.Printf("number of commits: %d\n", commitCount)
}
