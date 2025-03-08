package url

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewURL(t *testing.T) {
	t.Run("Succeeded in parsing", func(t *testing.T) {
		u, err := New("http://example.org/index.html")

		assert.Equal(t, "example.org", u.Host())
		assert.Nil(t, err)
	})
	t.Run("Without path", func(t *testing.T) {
		u, err := New("http://example.org")

		assert.Equal(t, "example.org", u.Host())
		assert.Equal(t, "/", u.Path())
		assert.Nil(t, err)
	})
}

func Test_Request(t *testing.T) {
	t.Run("Succeeded in requesting", func(t *testing.T) {
		u, err := New("https://example.org/index.html")

		content, err := u.Request()
		if err != nil {
			fmt.Println(err)
		}

		assert.NotEqual(t, "", content)
		assert.Nil(t, err)
	})
}

func Test_OpenFile(t *testing.T) {
	t.Run("Succeeded in opening file", func(t *testing.T) {
		u, err := New("file:///Users/louna/Louna-iTerm2.json")
		if err != nil {
			fmt.Println(err)
		}

		body, err := Load(u)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(body)
	})
}
