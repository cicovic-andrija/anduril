package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cicovic-andrija/libgo/crypto"
	"github.com/cicovic-andrija/libgo/fs"
)

func (env *Environment) ConfigInfo() string {
	return fmt.Sprintf("%s (encrypted=%t)", env.configPath, env.encryptedConfig)
}

func (env *Environment) UnmarshalConfig(v interface{}) error {
	if env.encryptedConfig {
		encodedBytes, err := os.ReadFile(env.configPath)
		if err != nil {
			return err
		}
		decodedBytes := make([]byte, base64.StdEncoding.DecodedLen(len(encodedBytes)))
		n, err := base64.StdEncoding.Decode(decodedBytes, encodedBytes)
		if err != nil {
			return err
		}
		decodedBytes = decodedBytes[:n]
		decryptedBytes, err := crypto.DecryptAES256CBCPBKDF2(decodedBytes, Version, Build)
		if err != nil {
			return err
		}
		return json.NewDecoder(bytes.NewReader(decryptedBytes)).Decode(v)
	} else { // plaintext
		if fhPlaintext, err := fs.OpenFile(env.configPath); err != nil {
			return err
		} else {
			defer fhPlaintext.Close()
			return json.NewDecoder(fhPlaintext).Decode(v)
		}
	}
}
