package result

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func readFileTest(fname string) (content []byte, err error) {
	defer Catch(&err)        // If a try fails, the wrapped error raised is set in the return value here
	f := Try(os.Open(fname)) // Try removes err part of return, and panics with special wrapper if err != nil
	defer f.Close()
	content = Try(io.ReadAll(f))
	return
}

func testPanic() (err error) {
	defer Catch(&err)
	var empty []string
	empty[0] = "hello"
	return
}

func TestError(t *testing.T) {
	_, err := readFileTest("nosuchfile")
	assert.Error(t, err)
}

func TestSuccess(t *testing.T) {
	content, err := readFileTest("/etc/passwd")
	assert.NoError(t, err)
	assert.NotEmpty(t, content)
}

func TestActualPanic(t *testing.T) {
	testPanic()
}
