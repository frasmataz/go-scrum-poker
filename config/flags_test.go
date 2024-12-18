package config

import (
	"reflect"
	"testing"
)

func TestGetConfigFromFlags(t *testing.T) {
	goodInputTestCases := map[string]struct {
		args []string
		want *Config
	}{
		"defaults": {
			[]string{},
			&Config{
				Host:      "localhost",
				Port:      "8080",
				DebugMode: false,
			},
		},
		"ip and port": {
			[]string{
				"-host", "127.0.0.1",
				"-port", "8090",
			},
			&Config{
				Host:      "127.0.0.1",
				Port:      "8090",
				DebugMode: false,
			},
		},
		"hostname with dots": {
			[]string{
				"-host", "my.example.com",
				"-port", "80",
			},
			&Config{
				Host:      "my.example.com",
				Port:      "80",
				DebugMode: false,
			},
		},
		"debugmode implicit true": {
			[]string{
				"-debugMode",
			},
			&Config{
				Host:      "localhost",
				Port:      "8080",
				DebugMode: true,
			},
		},
		"debugmode explicit true": {
			[]string{
				"-debugMode=true",
			},
			&Config{
				Host:      "localhost",
				Port:      "8080",
				DebugMode: true,
			},
		},
		"debugmode explicit false": {
			[]string{
				"-debugMode=false",
			},
			&Config{
				Host:      "localhost",
				Port:      "8080",
				DebugMode: false,
			},
		},
	}
	for name, goodInputTestCase := range goodInputTestCases {
		t.Run(name, func(t *testing.T) {
			got, err := GetConfigFromFlags(goodInputTestCase.args)
			if err != nil {
				t.Errorf("Error running GetConfigFromFlags(): %s", err)
			}

			if !reflect.DeepEqual(got, goodInputTestCase.want) {
				t.Errorf("GetConfigFromFlags() = %v, want %v", got, goodInputTestCase.want)
			}
		})
	}

	badInputTestCases := map[string]struct {
		args []string
	}{
		"letters in port": {
			[]string{
				"-port", "abc",
			},
		},
		"special chars in port 1": {
			[]string{
				"-port", "!£$%^&",
			},
		},
		"special chars in port 2": {
			[]string{
				"-port", "123*",
			},
		},
		"special chars in port 3": {
			[]string{
				"-port", "\"123\"",
			},
		},
		"negative port": {
			[]string{
				"-port", "-1",
			},
		},
		"zero port": {
			[]string{
				"-port", "0",
			},
		},
		"port over 65535": {
			[]string{
				"-port", "65536",
			},
		},
	}
	for name, badInputTestCase := range badInputTestCases {
		t.Run(name, func(t *testing.T) {
			_, err := GetConfigFromFlags(badInputTestCase.args)
			if err == nil {
				t.Errorf("Did not receive an error as expected from GetConfigFromFlags(): %v", badInputTestCase.args)
			}
		})
	}
}
