package service

import (
	"encoding/json"
	"fmt"

	"github.com/cicovic-andrija/libgo/fs"
)

func (env *Environment) ConfigInfo() string {
	return fmt.Sprintf("%s (encrypted=%t)", env.configPath, env.encryptedConfig)
}

func (env *Environment) UnmarshalConfig(v interface{}) error {
	if env.encryptedConfig { // encrypted
		// if fhCiphertext, err := util.OpenFile(env.configPath); err != nil {
		// 	return err
		// } else {
		// 	defer fhCiphertext.Close()
		// 	block, err := aes.NewCipher([]byte("matador"))
		// 	if err != nil {
		// 		return errors.New("failed to obtain a block cipher handler")
		// 	}
		// 	iv := make([]byte, aes.BlockSize)
		// 	if _, err := io.ReadFull(fhCiphertext, iv); err != nil {
		// 		return fmt.Errorf("config file loading failed: %v", err)
		// 	}
		// 	stream := cipher.NewCFBDecrypter(block, iv)
		// }
	} else { // plaintext
		if fhPlaintext, err := fs.OpenFile(env.configPath); err != nil {
			return err
		} else {
			defer fhPlaintext.Close()
			return json.NewDecoder(fhPlaintext).Decode(v)
		}
	}

	return nil
}
