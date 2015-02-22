package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
)

func generateKeyPair() {
	publicKeyCurve := elliptic.P256()

	privateKey := new(ecdsa.PrivateKey)
	privateKey, err := ecdsa.GenerateKey(publicKeyCurve, rand.Reader)

	if err != nil {
		log.Fatal(err)
	}

	privateKeyRaw, err := x509.MarshalECPrivateKey(privateKey)

	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("key", privateKeyRaw, 0644)

	if err != nil {
		log.Fatal(err)
	}

	publicKeyRaw, _ := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)

	err = ioutil.WriteFile("key.pub", publicKeyRaw, 0644)

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	if len(flag.Args()) == 1 {
		if flag.Args()[0] == "keygen" {
			generateKeyPair()
		}
	}
}
