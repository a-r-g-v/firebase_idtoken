package example_test

import (
	"context"
	"log"
	"os"

	"github.com/a-r-g-v/firebase_idtoken/firebase"
)

func ExampleClient_VerifyCustomToken() {
	apiKey := os.Getenv("API_KEY")
	customToken := os.Getenv("FIREBASE_CUSTOM_TOKEN")

	c := firebase.NewClient(apiKey)

	idToken, err := c.VerifyCustomToken(context.Background(), firebase.VerifyCustomTokenRequest{
		Token:             customToken,
		ReturnSecureToken: true,
	})
	if err != nil {
		log.Fatalf("VerifyCustomToken failed. err: %+v", err)
	}
	log.Printf("issued idToken: %s", idToken.IdToken)
}
