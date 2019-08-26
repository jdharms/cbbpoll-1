package app

import (
	"fmt"
	"io/ioutil"

	"github.com/dgrijalva/jwt-go"

	models "github.com/r-cbb/cbbpoll/pkg"
)

func createJWT(u models.User) (string, error) {
	var claims jwt.MapClaims = make(map[string]interface{})
	claims["name"] = u.Nickname
	claims["admin"] = u.IsAdmin

	alg := jwt.GetSigningMethod("RS256")
	keytext, err := ioutil.ReadFile("jwtRS256.key")
	if err != nil {
		fmt.Println("couldn't read from secret file")
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keytext)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("couldn't parse key from byte")
	}

	token := jwt.NewWithClaims(alg, claims)

	out, err := token.SignedString(key)
	if err != nil {
		fmt.Printf("Error signing token: %v", err)
	}
	return out, nil
}
