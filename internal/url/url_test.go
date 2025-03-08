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
		u, err := New("http://example.org/index.html")

		_, err = u.Request()
		if err != nil {
			fmt.Println(err)
		}

		assert.Nil(t, err)
	})
}
