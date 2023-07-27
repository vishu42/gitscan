package main

import (
	"fmt"
	"os"

	"github.com/vishu42/gitscan/pkg"
)

func main() {
	// get command line arguments
	// os.Args[0] is the name of the program
	repo := os.Args[1]

	workdir := "/Users/vishaltewatia/github/vishu42/gitscan/temp/Devops-Node"

	// clone repo
	err := pkg.CloneRepo(repo, "/Users/vishaltewatia/github/vishu42/gitscan/temp/")
	if err != nil {
		panic(err)
	}

	c, err := pkg.GetAllCommits(workdir)
	if err != nil {
		panic(err)
	}

	var allTokens []pkg.Token

	for _, commit := range c {
		// for each commit, list all files
		files, err := pkg.GetCommitFiles(workdir, commit)
		if err != nil {
			panic(err)
		}

		// for each file, get file contents
		for _, file := range files {
			cr, err := pkg.GetFileContent(workdir, commit, file)
			if err != nil {
				panic(err)
			}

			// scan file contents
			// new scan

			tokens, err := pkg.ScanFile(cr, commit, file)
			if err != nil {
				panic(err)
			}

			allTokens = append(allTokens, tokens...)
		}
	}

	tp := pkg.GetPossibleTokenPairs(allTokens)
	var verifiedTokenPairs []pkg.TokenPair

	// validate key pairs
	for _, p := range tp {
		err := pkg.VerifyAWSKey(p.KeyID, p.KeySecret)
		if err != nil {
			fmt.Printf("%s,%s,%s,%s,%s\n", p.Commit, p.File, p.KeyID, p.KeySecret, "invalid")
			continue
		}
		verifiedTokenPairs = append(verifiedTokenPairs, p)
	}

	// print verified key pairs
	fmt.Println("Verified key pairs:")
	for _, p := range verifiedTokenPairs {
		fmt.Printf("%s,%s,%s,%s\n", p.Commit, p.File, p.KeyID, p.KeySecret)
	}
}
