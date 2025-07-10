package main

import (
	"fmt"
	"github.com/msteinert/pam"
	"golang.org/x/term"
	"os"
	"os/user"
	"path/filepath"
	"syscall"
)

const (
	_configPath = "/etc/jus.toml"

	_cantLoadConfig = "can't load config /etc/jus.toml"
	_unknownUser    = "user %s is not in jus.toml\n"
	_unknownAsUser  = "unknown <as> user %s\n"
	_unknownError   = "unkown error:%s\n"

	_pleasePassword = "password: "

	_pamInitError         = "can't init pam backend:%s\n"
	_pamMessageStyleError = "can't set pam message style:%s\n"
	_pamCloseError        = "can't close pam backend:%s\n"
	_pamAuthError         = "auth failed:%s\n"
	_pamAccCheck          = "acc validation failed:%s\n "
	_nothing              = "nothing to do"
	_cantFind             = "can't find %s\n"

	_cantUid  = "can't set uid:%s\n"
	_cantGid  = "can't set gid:%s\n"
	_cantCall = "can't call:%s\n"
)

var (
	_stdPaths = []string{"/bin/", "/sbin/", "/usr/sbin", "/usr/bin/", "/usr/local/bin"}
	_config   = Config{}
	_permit   = Permit{}
)

// CheckUser find current user in /etc/jus.toml
func CheckUser() {
	lookup, err := user.Current()
	if err != nil {
		fmt.Printf(_unknownError, err)
		os.Exit(1)
	}
	launcherUsername := lookup.Username

	for _, perm := range _config.Permits {
		if perm.User == launcherUsername {
			_permit = perm
		}
	}

	if _permit.User == "" {
		fmt.Printf(_unknownUser, launcherUsername)
		os.Exit(1)
	}

	if _permit.Paths == nil {
		_permit.Paths = _stdPaths
	}
}

// PasswordRead read pass with term.ReadPassword
func PasswordRead(prompt string) ([]byte, error) {
	fmt.Fprint(os.Stderr, prompt)
	defer fmt.Fprintf(os.Stderr, "\n")

	return term.ReadPassword(int(syscall.Stdin))
}

// PasswordCheck checks further access via PAM
func PasswordCheck() {
	pass, err := PasswordRead(_pleasePassword)
	if err != nil {
		fmt.Printf(_unknownError, err)
		os.Exit(1)
	}
	defer func() {
		for i := range pass {
			pass[i] = 0
		}
	}()

	trans, err := pam.StartFunc("system-auth", _permit.User, func(s pam.Style, msg string) (string, error) {
		switch s {
		case pam.PromptEchoOff:
			return string(pass), nil
		case pam.PromptEchoOn:
			return _permit.User, nil
		case pam.ErrorMsg:
			fmt.Println(msg)
			return "", nil
		case pam.TextInfo:
			fmt.Println(msg)
			return "", nil
		default:
			return "", fmt.Errorf(_pamMessageStyleError, msg)
		}
	})
	if err != nil {
		fmt.Printf(_pamInitError, err)
		os.Exit(1)
	}
	defer func() {
		if err = trans.CloseSession(pam.Silent); err != nil {
			fmt.Printf(_pamCloseError, err)
			os.Exit(1)
		}
	}()
	err = trans.Authenticate(pam.Silent)
	if err != nil {
		fmt.Printf(_pamAuthError, err)
		os.Exit(1)
	}
	err = trans.AcctMgmt(pam.Silent)
	if err != nil {
		fmt.Printf(_pamAccCheck, err)
		os.Exit(1)
	}
}

// FindUtility find commands in paths
func FindUtility(data string, paths []string) (string, bool) {
	for _, path := range paths {
		finPath := filepath.Join(path, data)

		_, err := os.Stat(finPath)
		if err == nil {
			return finPath, true
		}
	}
	return "", false

}
