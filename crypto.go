package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"strings"

	crypto "github.com/foursee/shellgameCrypto"
	"golang.org/x/crypto/openpgp"
)

func myCryptoKey() (privKey, pubKey string, err error) {
	if config().DeviceIdentity.PublicKey == "" || config().DeviceIdentity.PrivateKey == "" {
		return generateCryptoKey()
	}
	return config().DeviceIdentity.PrivateKey, config().DeviceIdentity.PublicKey, nil
}

func generateCryptoKey() (string, string, error) {
	log.Output(0, fmt.Sprintf("Generating %v-bit RSA keypair...", keyBits()))
	user, _ := user.Current()
	name := user.Username
	comment := "PolyRythm generated keypair"
	hostname, _ := os.Hostname()
	email := fmt.Sprintf("%s@%s", name, hostname)
	privKey, pubKey, err := crypto.GenerateRSAKeyPair(keyBits(), name, comment, email)
	check(err)
	config().DeviceIdentity.PrivateKey = privKey
	config().DeviceIdentity.PublicKey = pubKey
	config().save()
	return privKey, pubKey, err
}

func base64reader(s string) io.Reader {
	return base64.NewDecoder(base64.StdEncoding, strings.NewReader(s))
}

func base64keyRing(s string) (openpgp.EntityList, error) {
	return openpgp.ReadKeyRing(base64reader(s))
}
