package key

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
)

type RsaGen struct {
	keyLenght  int
	privateKey []byte
	publicKey  []byte
}

func NewRsaGen(kLen int) *RsaGen {
	return &RsaGen{keyLenght: kLen}
}

func (r *RsaGen) GenKey() error {
	pik, err := rsa.GenerateKey(rand.Reader, r.keyLenght)
	if err != nil {
		return err
	}
	block := &pem.Block{
		Type:  "private",
		Bytes: x509.MarshalPKCS1PrivateKey(pik),
	}
	file, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}

	r.privateKey, err = ioutil.ReadFile("private.pem")
	if err != nil {
		return err
	}

	pbk := &pik.PublicKey
	der, err := x509.MarshalPKIXPublicKey(pbk)
	if err != nil {
		return err
	}

	block = &pem.Block{
		Type:  "public",
		Bytes: der,
	}

	file, err = os.Create("public.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	r.publicKey, err = ioutil.ReadFile("public.pem")
	if err != nil {
		return err
	}
	return nil
}

func (r *RsaGen) Encrypt(data []byte) ([]byte, error) {
	block, _ := pem.Decode(r.publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

func (r *RsaGen) Decrypt(data []byte) ([]byte, error) {
	block, _ := pem.Decode(r.privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, data)
}
