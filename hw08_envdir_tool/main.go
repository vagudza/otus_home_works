package main

import "os"

func main() {
	// os.Args = [
	//    "./go-envdir",  // [0] — this program name
	//    "/path/to/envdir",  // [1] — envdir directory
	//    "echo",  // [2] — command to run
	//    "Hello World"  // [3] — command arguments
	// ]

	if len(os.Args) < 3 {
		_, _ = os.Stderr.WriteString("Usage: go-envdir <envdir> <cmd> [args...]\n")
		// https://www.unix.com/man_page/debian/8/envdir/
		// envdir exits 111 if it has trouble reading d, if it runs out of memory for environment variables,
		// or if it cannot run child.  Otherwise its exit code is the same as that of child
		os.Exit(111)
	}

	envDir := os.Args[1]
	cmd := os.Args[2:]

	env, err := ReadDir(envDir)
	if err != nil {
		_, _ = os.Stderr.WriteString("Error reading env dir: " + err.Error() + "\n")
		os.Exit(111)
	}

	exitCode := RunCmd(cmd, env)
	os.Exit(exitCode)
}
