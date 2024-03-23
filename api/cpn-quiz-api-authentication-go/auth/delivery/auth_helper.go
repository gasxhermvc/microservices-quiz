package delivery

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/jwk"
	config "github.com/spf13/viper"
)

func getKey(token *jwt.Token) (interface{}, error) {
	endpoint := config.GetString("cpn.quiz.sso.endpoint")
	env := os.Args[1]
	if env == "dev" {
		endpoint = "http://cpn-quiz-keycloak:8080"
	}
	fmt.Println(endpoint)
	keySet, err := jwk.Fetch(context.Background(), endpoint+"/realms/cpn-quiz/protocol/openid-connect/certs")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	keyID, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("expecting JWT header to have a key ID in the kid field")
	}

	key, found := keySet.LookupKeyID(keyID)
	if !found {
		return nil, fmt.Errorf("unable to find key %q", keyID)
	}

	var pubkey interface{}
	if err := key.Raw(&pubkey); err != nil {
		return nil, fmt.Errorf("unable to get the public key. Error: %s", err.Error())
	}

	return pubkey, nil
}
