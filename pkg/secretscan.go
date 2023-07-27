package pkg

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

const (
	TOKEN_TYPE_AWS_IAM_KEYID     = "aws_iam_keyid"
	TOKEN_TYPE_AWS_IAM_KEYSECRET = "aws_iam_keysecret"
	ErrPathNotAbsolute           = "path is not absolute"
)

type Token struct {
	Token     string
	TokenType string
	LineNo    int
	Commit    string
	File      string
}

type TokenPair struct {
	KeyID     string
	KeySecret string
	Commit    string
	File      string
}

// GetPossibleTokenPairs returns a list of TokenPair that are possible AWS IAM key pairs
// Considerations while scanning for possible key pairs:
// 1. KeyID and KeySecret must be in the same file
// 2. If multiple key pairs are found in the same file, the result will contain all possible combinations
func GetPossibleTokenPairs(t []Token) (result []TokenPair) {
	// create a map of commit -> []Token
	commitMap := make(map[string][]Token)
	for _, token := range t {
		commitMap[token.Commit] = append(commitMap[token.Commit], token)
	}
	// iterate over the map
	for commit, tokens := range commitMap {
		// create a map of keyid -> []Token
		keyIDMap := make(map[string][]Token)
		for _, token := range tokens {
			if token.TokenType == TOKEN_TYPE_AWS_IAM_KEYID {
				keyIDMap[token.Token] = append(keyIDMap[token.Token], token)
			}
		}

		// create a map of keysecret -> []Token
		keySecretMap := make(map[string][]Token)
		for _, token := range tokens {
			if token.TokenType == TOKEN_TYPE_AWS_IAM_KEYSECRET {
				keySecretMap[token.Token] = append(keySecretMap[token.Token], token)
			}
		}

		// iterate through keyIDMap
		for keyID, keyIDTokens := range keyIDMap {
			// iterate through keySecretMap
			for keySecret, keySecretTokens := range keySecretMap {
				// iterate through keyIDTokens
				for _, keyIDToken := range keyIDTokens {
					// iterate through keySecretTokens
					for _, keySecretToken := range keySecretTokens {
						// if keyIDToken and keySecretToken are in the same file
						if keyIDToken.File == keySecretToken.File {
							// append to result
							result = append(result, TokenPair{keyID, keySecret, commit, keyIDToken.File})
						}
					}
				}
			}
		}
	}

	return
}

func IsKey(token string) (isKey bool, keyType string, err error) {
	re := regexp.MustCompile(`AKIA[0-9A-Z]{16}`)
	if re.MatchString(token) {
		// print debug message
		isKey = true
		keyType = TOKEN_TYPE_AWS_IAM_KEYID
		return
	}
	re = regexp.MustCompile(`^[a-zA-Z0-9+=/]{40}$`)
	if re.MatchString(token) {
		isKey = true
		keyType = TOKEN_TYPE_AWS_IAM_KEYSECRET
		return
	}
	return
}

type BareToken struct {
	token     string
	tokenType string
}

func ScanLine(line string) (awsKeys []BareToken, err error) {
	// create a new io.Reader for the line
	stringReader := strings.NewReader(line)
	// new scanner for the line
	scanner := bufio.NewScanner(stringReader)
	// split the line into tokens
	scanner.Split(bufio.ScanWords)
	// scan the line
	for scanner.Scan() {
		// print each token
		token := scanner.Text()
		// print token
		isKey, keyType, err := IsKey(token)
		if err != nil {
			return nil, err
		}
		if isKey {
			awsKeys = append(awsKeys, BareToken{token, keyType})
		}
	}
	return
}

func ScanFile(r io.Reader, commit string, file string) (result []Token, err error) {
	lineNo := 0
	// scan file line by line
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		// print each line
		line := scanner.Text()
		lineNo++
		tokens, err := ScanLine(line)
		if err != nil {
			return nil, err
		}
		for _, token := range tokens {
			result = append(result, Token{token.token, token.tokenType, lineNo, commit, file})
		}

	}
	return
}
