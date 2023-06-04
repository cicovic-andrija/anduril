package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/cicovic-andrija/anduril/anduril"
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
	encrypted    bool
	password     string
}

func main() {
	options := parseOptions()
	template := openTemplate(options.templatePath)
	fillTemplate(template, options.profile)
	writeTemplate(template, options.outPath)
}

func fillTemplate(template *anduril.Config, profile string) {
	switch profile {
	case DevProfile:
		template.HTTPS.Network.IPAcceptHost = "localhost"
		template.HTTPS.Network.TCPPort = 8080
		template.HTTPS.Network.TLSCertPath = "tlspublic.crt"
		template.HTTPS.Network.TLSKeyPath = "tlsprivate.key"
		template.HTTPS.AllowOnlyGETRequests = false
	case ProdProfile:
	}
}

func writeTemplate(template *anduril.Config, path string) {
	bytes, err := json.Marshal(template)
	if err != nil {
		die(err.Error())
	}
	err = os.WriteFile(path, bytes, 0644)
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

func parseOptions() (opt options) {
	flag.StringVar(&opt.templatePath, "template", "", "config template path")
	flag.StringVar(&opt.outPath, "to", "", "output file path")
	flag.StringVar(&opt.profile, "profile", "dev", "environment profile")
	flag.BoolVar(&opt.encrypted, "encrypted", true, "config file encryption")
	flag.StringVar(&opt.password, "password", "", "ssh key password")
	flag.Parse()

	if opt.templatePath == "" {
		die("template not provided")
	}

	if !(opt.profile == DevProfile || opt.profile == ProdProfile) {
		die("invalid profile")
	}

	if !opt.encrypted && opt.password != "" {
		warn("unencrypted config will contain a password")
	}

	return
}

func warn(message string) {
	fmt.Printf("warning: %s!", message)
}

func die(message string) {
	fmt.Println(message)
	os.Exit(0)
}
