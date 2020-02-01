package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestReadRandomLineFromFile(t *testing.T) {
	errorsLogged := []string{}
	logErr = func(format string, args ...interface{}) {
		if len(args) > 0 {
			errorsLogged = append(errorsLogged, fmt.Sprintf(format, args))
		} else {
			errorsLogged = append(errorsLogged, format)
		}
	}

	type Test struct {
		Name          string
		ExpectedError string
		toLower       bool
		Before        func()
		After         func()
		EnvName       string
	}

	prefix := "/synonym/ar/"
	baseEndpoint := "https://synonyms.reverso.net"

	defaultTextDir := "../ar_synonyms/arabic.txt"
	testTable := []Test{
		Test{
			Name:          "ARABIC_WORD_LIST_PATH not set",
			ExpectedError: "ARABIC_WORD_LIST_PATH was not set",
			EnvName:       "ARABIC_WORD_LIST_PATH",
			toLower:       true,
			Before: func() {
				os.Setenv("ARABIC_WORD_LIST_PATH", "")
			},
			After: func() {
				os.Setenv("ARABIC_WORD_LIST_PATH", defaultTextDir)
			},
		},
		Test{
			Name:          "gets random word succesfully",
			EnvName:       "ARABIC_WORD_LIST_PATH",
			ExpectedError: "",
			toLower:       true,
			Before: func() {
				os.Setenv("ARABIC_WORD_LIST_PATH", defaultTextDir)
			},
			After: func() {},
		},
		Test{
			Name:          "no such path",
			ExpectedError: "open this/does/not/exist: no such file or directory",
			EnvName:       "ARABIC_WORD_LIST_PATH",
			toLower:       true,
			Before: func() {
				os.Setenv("ARABIC_WORD_LIST_PATH", "this/does/not/exist")
			},
			After: func() {
				os.Setenv("ARABIC_WORD_LIST_PATH", defaultTextDir)
			},
		},
	}

	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			test.Before()
			w, err := ReadRandomLineFromFile(test.EnvName, baseEndpoint, prefix, test.toLower)
			// positive tests
			if test.ExpectedError == "" {
				assert.NotEqual(t, "", w)
				assert.Equal(t, nil, err)
			} else {
				assert.Equal(t, "", w)
				assert.NotEqual(t, nil, err)
				assert.Equal(t, test.ExpectedError, err.Error())
			}
			test.After()
		})
	}

}
