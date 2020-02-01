package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadRandomLineFromFile(t *testing.T) {
	type Test struct {
		Name          string
		File          string
		ExpectedError string
	}

	defaultTextDir := "../ar_synonyms/arabic.txt"
	testTable := []Test{
		Test{
			Name:          "Bad file",
			File:          "sdfsdf",
			ExpectedError: "ARABIC_WORD_LIST_PATH was not set",
		},
		Test{
			Name:          "gets random word succesfully",
			File:          defaultTextDir,
			ExpectedError: "",
		},
	}

	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			w, err := ReadRandomLineFromFile(test.File)
			if err != nil {
				assert.Equal(t, test.ExpectedError, err.Error())
			} else {
				assert.Equal(t, test.ExpectedError, "")
				assert.NotEqual(t, "", w)
			}

		})
	}

}
