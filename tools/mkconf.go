package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/cicovic-andrija/anduril/anduril"
	"github.com/cicovic-andrija/anduril/service"
	"github.com/cicovic-andrija/libgo/crypto"
	"github.com/cicovic-andrija/libgo/fs"
)

const (
	DevProfile  = "dev"
	ProdProfile = "prod"
)

type options struct {
	templatePath string
	outPath      string
	profile      string
	password     string
	salt         string
	decrypt      bool
}

func parseOptions() (opt options) {
	flag.StringVar(&opt.templatePath, "template", "", "config template path")
	flag.StringVar(&opt.outPath, "to", "", "output file path")
	flag.StringVar(&opt.profile, "profile", "dev", "environment profile")
	flag.StringVar(&opt.password, "password", "", "encryption password")
	flag.StringVar(&opt.salt, "salt", "", "encryption salt")
	flag.BoolVar(&opt.decrypt, "decrypt", false, "decrypt output file to stdout (validation)")
	flag.Parse()

	if opt.templatePath == "" {
		die("template not provided")
	}

	if !(opt.profile == DevProfile || opt.profile == ProdProfile) {
		die("invalid profile")
	}

	if opt.password == "" {
		opt.password = service.Version
	}

	if opt.salt == "" {
		opt.salt = service.Build
	}

	return
}

func main() {
	options := parseOptions()
	template := openTemplate(options.templatePath)
	replaceValues(template, options.profile)
	saveEncrypted(template, options.outPath, options.password, options.salt)
	if options.decrypt {
		decrypt(options.outPath, options.password, options.salt)
	}
}

func replaceValues(template *anduril.Config, profile string) {
	switch profile {
	case DevProfile:
		template.HTTPS.Network.IPAcceptHost = "localhost"
		template.HTTPS.Network.TCPPort = 8080
		template.HTTPS.Network.TLSCertPath = "tlspublic.crt"
		template.HTTPS.Network.TLSKeyPath = "tlsprivate.key"
		template.HTTPS.AllowOnlyGETRequests = false
		template.Repository.SSHAuth.PrivateKeyPath = "gitprivate.key"
		template.Repository.SSHAuth.PrivateKeyPassword = ""
		template.Settings.RepositorySyncPeriod = "10s"
		template.Settings.StaleFileCleanupPeriod = "1h"
	case ProdProfile:
		template.HTTPS.Network.IPAcceptHost = "any"
		template.HTTPS.Network.TCPPort = 443
		template.HTTPS.AllowOnlyGETRequests = true
		template.Settings.RepositorySyncPeriod = "15m"
		template.Settings.StaleFileCleanupPeriod = "24h"
	}
}

func saveEncrypted(config *anduril.Config, outPath string, password string, salt string) {
	bytes, err := json.Marshal(config)
	if err != nil {
		die(err.Error())
	}

	encryptedBytes, err := crypto.EncryptAES256CBCPBKDF2(bytes, password, salt)
	if err != nil {
		die(err.Error())
	}

	fh, err := os.Create(outPath)
	if err != nil {
		die(err.Error())
	}
	defer fh.Close()

	_, err = fh.WriteString(base64.StdEncoding.EncodeToString(encryptedBytes))
	if err != nil {
		die(err.Error())
	}
}

func openTemplate(path string) *anduril.Config {
	fh, err := fs.OpenFile(path)
	if err != nil {
		die(err.Error())
	}
	defer fh.Close()

	config := &anduril.Config{}
	if err := json.NewDecoder(fh).Decode(config); err != nil {
		die(err.Error())
	}

	return config
}

func decrypt(path string, password string, salt string) {
	encodedBytes, err := os.ReadFile(path)
	if err != nil {
		die(err.Error())
	}

	decodedBytes := make([]byte, base64.StdEncoding.DecodedLen(len(encodedBytes)))
	n, err := base64.StdEncoding.Decode(decodedBytes, encodedBytes)
	if err != nil {
		die(err.Error())
	}
	decodedBytes = decodedBytes[:n]

	decryptedBytes, err := crypto.DecryptAES256CBCPBKDF2(decodedBytes, password, salt)
	if err != nil {
		die(err.Error())
	}

	fmt.Printf("%s\n", string(decryptedBytes))
}

func die(message string) {
	fmt.Println(message)
	os.Exit(0)
}
