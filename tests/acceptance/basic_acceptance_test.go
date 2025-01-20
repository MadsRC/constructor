package acceptance

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"testing"
)

func ensureBinPath(t *testing.T) string {
	t.Helper()
	path, ok := os.LookupEnv("CONSTRUCTOR_BIN_PATH")
	require.True(t, ok, "CONSTRUCTOR_BIN_PATH must be set")
	require.NotEmpty(t, path, "CONSTRUCTOR_BIN_PATH must not be empty")
	require.FileExists(t, path, "CONSTRUCTOR_BIN_PATH must exist and be a file")

	return path
}

func createModule(t *testing.T, dir string) {
	t.Helper()
	var outb, errb bytes.Buffer
	cmd := exec.Command("go", "mod", "init", "acceptance")
	cmd.Dir = dir
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	require.NoError(t, err)
	require.Equal(t, "go: creating new go.mod: module acceptance\n", errb.String())
}

func TestBasicAcceptance(t *testing.T) {
	binPath := ensureBinPath(t)
	dir := t.TempDir()
	createModule(t, dir)
	var outb, errb bytes.Buffer
	mainF, err := os.Create(dir + "/file.go")
	require.NoError(t, err)
	defer mainF.Close()
	testF, err := os.Create(dir + "/file_test.go")
	require.NoError(t, err)
	defer testF.Close()

	cmd := exec.Command(binPath)
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err = cmd.Run()
	assert.NoError(t, err)
	require.Empty(t, errb.String())

	n, err := mainF.Write(outb.Bytes())
	require.NoError(t, err)
	require.Equal(t, n, len(outb.Bytes()))

	var outb1, errb1 bytes.Buffer
	cmd1 := exec.Command(binPath, "--test")
	cmd1.Stdout = &outb1
	cmd1.Stderr = &errb1
	err = cmd1.Run()
	assert.NoError(t, err)
	require.Empty(t, errb1.String())

	n, err = testF.Write(outb1.Bytes())
	require.NoError(t, err)
	require.Equal(t, n, len(outb1.Bytes()))

	var outb2, errb2 bytes.Buffer
	cmd2 := exec.Command("go", "test", ".")
	cmd2.Dir = dir
	cmd2.Stdout = &outb2
	cmd2.Stderr = &errb2
	err = cmd2.Run()
	assert.NoError(t, err)
	require.Empty(t, errb2.String())
	require.Contains(t, outb2.String(), "ok  \tacceptance")
}

func TestBasicAcceptance_with_output_parameter(t *testing.T) {
	binPath := ensureBinPath(t)
	dir := t.TempDir()
	createModule(t, dir)
	var outb, errb bytes.Buffer

	args := []string{"--output", dir + "/file.go"}
	cmd := exec.Command(binPath, args...)
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	assert.NoError(t, err)
	require.Empty(t, errb.String())
	require.Empty(t, outb.String())
	require.FileExists(t, dir+"/file.go")

	var outb1, errb1 bytes.Buffer
	args = []string{"--output", dir + "/file_test.go", "--test"}
	cmd1 := exec.Command(binPath, args...)
	cmd1.Stdout = &outb1
	cmd1.Stderr = &errb1
	err = cmd1.Run()
	assert.NoError(t, err)
	require.Empty(t, errb1.String())
	require.FileExists(t, dir+"/file_test.go")

	var outb2, errb2 bytes.Buffer
	cmd2 := exec.Command("go", "test", ".")
	cmd2.Dir = dir
	cmd2.Stdout = &outb2
	cmd2.Stderr = &errb2
	err = cmd2.Run()
	assert.NoError(t, err)
	require.Empty(t, errb2.String())
	require.Contains(t, outb2.String(), "ok  \tacceptance")
}
