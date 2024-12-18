package config

import (
	"errors"
	"flag"
	"regexp"
	"strconv"
)

type Config struct {
	Host      string
	Port      string
	DebugMode bool
}

func GetConfigFromFlags(args []string) (*Config, error) {
	flagset := flag.NewFlagSet("", flag.ExitOnError)
	hostPtr := flagset.String("host", "localhost", "Hostname to host on")
	portPtr := flagset.String("port", "8080", "Port to host on")
	debugPtr := flagset.Bool("debugMode", false, "Enables gin debug mode")

	flagset.Parse(args)

	err := validatePort(*portPtr)
	if err != nil {
		return nil, err
	}

	return &Config{
		Host:      *hostPtr,
		Port:      *portPtr,
		DebugMode: *debugPtr,
	}, nil
}

func validatePort(port string) error {
	re, err := regexp.Compile("^[0-9]+$")
	if err != nil {
		return err
	}
	if !re.MatchString(port) {
		return errors.New("port is not a whole, positive number")
	}

	if port == "0" {
		return errors.New("port cannot be zero")
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return err
	}
	if portInt > 65535 {
		return errors.New("port cannot be higher than 65535")
	}

	return nil
}
