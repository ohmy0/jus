package main

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	args := os.Args[1:]

	CheckConfigPerm()
	LoadConfig()
	CheckUser()
	PasswordCheck()

	if len(args) == 0 {
		fmt.Println(_nothing)
		os.Exit(0)
	}

	if _, err := os.Stat(args[0]); err != nil {
		full, exists := FindUtility(args[0], _permit.Paths)
		if !exists {
			fmt.Printf(_cantFind, args[0])
			os.Exit(0)
		}
		args[0] = full
	}

	if _permit.As == "" {
		fmt.Printf(_unknownAsUser, _permit.As)
		os.Exit(1)
	}
	finalUser, err := user.Lookup(_permit.As)
	if err != nil {
		fmt.Printf(_unknownAsUser, _permit.As)
		os.Exit(1)
	}

	finalUid, err := strconv.Atoi(finalUser.Uid)
	if err != nil {
		fmt.Printf(_unknownError, err)
		os.Exit(1)
	}
	finalGid, err := strconv.Atoi(finalUser.Gid)
	if err != nil {
		fmt.Printf(_unknownError, err)
		os.Exit(1)
	}

	err = syscall.Setuid(finalUid)
	if err != nil {
		fmt.Printf(_cantUid, err)
		os.Exit(1)
	}
	err = syscall.Setgid(finalGid)
	if err != nil {
		fmt.Printf(_cantGid, err)
		os.Exit(1)
	}

	env := []string{}
	if _permit.KeepEnv {
		for _, e := range os.Environ() {
			if strings.HasPrefix(e, "PATH=") ||
				strings.HasPrefix(e, "TERM=") ||
				strings.HasPrefix(e, "LANG=") {
				env = append(env, e)
			}
		}
	} else {
		env = _stdEnv
	}

	err = syscall.Exec(args[0], args[0:], env)
	if err != nil {
		fmt.Printf(_cantCall, err)
	}
}
