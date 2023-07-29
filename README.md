# GITSCAN

## Requirements
git version 2.32.0 
aws-cli/2.13.3
go version go1.20.2

## Compilation and Usage

```bash
# run the following command to build the binary
go build
```

```bash
# run the following command to run the binary
usage: ./gitscan <git repository url>
```

## Problem

To scan the git repository for sensitive information including the commit history.

## Solution

ALGORITHM:

1. Clone the repository
2. Get all the commits using `git log --all --pretty=format:"%H"`
3. Get the file names of all the files in the repository using `git ls-tree --name-only -r <commit>`
4. Get the content of the file using `git show <commit>:<file>`
5. Tokenize the content of the file, and check for the sensitive information using regex. 
   1. Regex used for key id is `AKIA[0-9A-Z]{16}`
   2. Regex used for secret key is `^[a-zA-Z0-9+=/]{40}$` - which simply checks for 40 character base 64 string. This document https://aws.amazon.com/blogs/security/a-safer-way-to-distribute-aws-credentials-to-ec2/ suggests that the secret key is 40 character base 64 string excluding the first and last char. However the valid key does have first and last char as base64. So I have used the regex which checks for 40 character base 64 string. It seems that the key specified in the document is not valid anymore.
6. If the file contains the aws key id or aws key secret, then add the file name and the commit id, line number and the token to the result.
7. Iterate through all the tokens, and produce all possible combinations of aws key id and aws key secret.
8. If the combination is valid, then add the file name and the commit id, and the combination to the result.

## Design Decisions

- I have used git cli to get the commit history, file names and the file content. I could have explored the git api, but considering the time constraint, I have used the git cli straight away.
- I have used aws cli to verify the aws key id and aws key secret for the same reason as above.


## Contact Information
Name: Vishal Tewatia

Email: tewatiavishal3@gmail.com

Phone: +91-8395955922

Github: https://github.com/vishu42

Portfolio: https://vishaltewatia.com

## Optional Enhancements

Have not implemented any optional enhancements due to time constraint. I need to report to office tomorrow so i have to submit the assignment today itself. I was able to spend only 1 day on this.
