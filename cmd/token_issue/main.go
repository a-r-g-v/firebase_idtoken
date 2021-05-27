package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"golang.org/x/xerrors"

	"github.com/a-r-g-v/firebase_idtoken/firebase"
)

func init() {
}

func main() {
	cfg, err := parseArgv()
	if err != nil {
		log.Printf("[ERROR] parse argv failed. err: %+v\n", err)
		printUsage()
		return
	}

	ctx := context.Background()
	ctx, canceler := signal.NotifyContext(ctx, os.Interrupt)
	defer canceler()

	r, err := firebase.NewClient(cfg.apiKey).VerifyCustomToken(ctx, firebase.VerifyCustomTokenRequest{
		Token:             cfg.customToken,
		ReturnSecureToken: true,
	})
	if err != nil {
		log.Fatalf("[ERROR] VerifyCustomToken failed. err: %+v\n", err)
	}
	log.Print(r.IdToken)
}

func printUsage() {
	fmt.Printf(`
Command Usage
	Args:
		%s <apiKey> <firebaseCustomToken>

		- apiKey -> google cloud platform api key
		- firebaseCustomToken -> already issued firebase custom token that you want to convert id token
	
	Output:
		firebase id token
	`, os.Args[0])
}

type config struct {
	apiKey      string
	customToken string
}

func parseArgv() (*config, error) {
	if len(os.Args) != 3 {
		return nil, xerrors.Errorf("len(os.Args) must be 3. but %d", len(os.Args))
	}
	apiKey := os.Args[1]
	customToken := os.Args[2]
	return &config{
		apiKey:      apiKey,
		customToken: customToken,
	}, nil
}
