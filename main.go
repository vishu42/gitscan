package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/vishu42/gitscan/pkg"
)

func main() {
	// get command line arguments
	repo := os.Args[1]

	// create a temp directory
	workdir := "/tmp/gitscan-" + strconv.Itoa(rand.Intn(1000000))

	// get the name of the repo
	repoSplited := strings.Split(repo, "/")
	repoName := repoSplited[len(repoSplited)-1]
	clonedir := workdir + "/" + repoName

	log.Println("Workdir:", workdir)

	// create temp workdir
	err := pkg.MkDir(workdir)
	if err != nil {
		panic(err)
	}

	// remove temp workdir after function execution is done
	defer func() {
		err = pkg.RmDir(workdir)
		if err != nil {
			panic(err)
		}
	}()

	// clone repo
	err = pkg.CloneRepo(repo, workdir)
	if err != nil {
		panic(err)
	}

	c, err := pkg.GetAllCommits(clonedir)
	if err != nil {
		panic(err)
	}

	var allTokens []pkg.Token

	for _, commit := range c {
		// for each commit, list all files
		files, err := pkg.GetCommitFiles(clonedir, commit)
		if err != nil {
			panic(err)
		}

		// for each file, get file contents
		for _, file := range files {
			cr, err := pkg.GetFileContent(clonedir, commit, file)
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
