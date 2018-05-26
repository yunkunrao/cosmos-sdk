package tests

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var nameCounter = 0

func getCmd(t *testing.T, command string) (cmd *exec.Cmd) {

	nameCounter++
	tmpfile, err := ioutil.TempFile("", fmt.Sprintf("tempCmd%v.sh", nameCounter))
	require.NoError(t, err)
	path := tmpfile.Name()

	content := []byte(fmt.Sprintf("#!/bin/bash\n%v", command))
	_, err = tmpfile.Write(content)
	require.NoError(t, err)

	cmd = exec.Command("/bin/bash", path)
	return
}

// remove temp bash files created with getCmd
func RemoveTempFiles(t *testing.T) {
	for ; nameCounter >= 0; nameCounter-- {
		tmpfile, err := ioutil.TempFile("", fmt.Sprintf("tempCmd%v.sh", nameCounter))
		require.NoError(t, err)
		os.Remove(tmpfile.Name())
	}
}

// Execute the command, return standard output and error, try a few times if requested
func ExecuteT(t *testing.T, command string) (out string) {
	cmd := getCmd(t, command)

	bz, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	require.NoError(t, err, string(bz))
	out = strings.Trim(string(bz), "\n") //trim any new lines
	time.Sleep(time.Second)

	return out
}

// Asynchronously execute the command, return standard output and error
func GoExecuteT(t *testing.T, command string) (cmd *exec.Cmd, pipeIn io.WriteCloser, pipeOut io.ReadCloser) {
	cmd = getCmd(t, command)
	pipeIn, err := cmd.StdinPipe()
	require.NoError(t, err)
	pipeOut, err = cmd.StdoutPipe()
	require.NoError(t, err)
	cmd.Start()
	time.Sleep(time.Second)
	return cmd, pipeIn, pipeOut
}

//func getCmd(t *testing.T, command string) (*exec.Cmd, shFile) {

////split command into command and args
//split := strings.Split(command, " ")
//require.True(t, len(split) > 0, "no command provided")

//var cmd *exec.Cmd
//if len(split) == 1 {
//cmd = exec.Command(split[0])
//} else {
//cmd = exec.Command(split[0], split[1:]...)
//}
//return cmd
//}
