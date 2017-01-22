package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/koudaiii/qcli/quay"
)

type AddUserCommand struct {
	Meta
}

type DeleteUserCommand struct {
	Meta
}

func (c *AddUserCommand) Run(args []string) int {
	if err := FlagInit(args); err != nil {
		fmt.Fprintln(os.Stderr, c.Help())
		os.Exit(1)
	}

	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, c.Help())
		os.Exit(1)
	}

	ss := strings.Split(args[0], "/")
	if len(ss) != 2 {
		fmt.Fprintln(os.Stderr, c.Help())
		os.Exit(1)
	}

	repos, err := quay.AddPermission(ss[0], ss[1], "user", args[1], role)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "Added! %v(%v) in quay.io/%v/%v\n", repos.Name, repos.Role, ss[0], ss[1])
	return 0
}

func (c *AddUserCommand) Synopsis() string {
	return fmt.Sprint("Add user in repository")
}

func (c *AddUserCommand) Help() string {
	helpText := `
qcli supported only Quay.io
Usage: add-user
  qcli add-user quay.io/koudaiii/qcli koudaiii --role admin
`
	return strings.TrimSpace(helpText)
}

func (c *DeleteUserCommand) Run(args []string) int {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, c.Help())
		os.Exit(1)
	}

	ss := strings.Split(args[0], "/")
	if len(ss) != 2 {
		fmt.Fprintln(os.Stderr, c.Help())
		os.Exit(1)
	}

	err := quay.DeletePermission(ss[0], ss[1], "user", args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "Deleted! %v in quay.io/%v/%v\n", args[1], ss[0], ss[1])
	return 0
}

func (c *DeleteUserCommand) Synopsis() string {
	return fmt.Sprint("Delete user in repository")
}

func (c *DeleteUserCommand) Help() string {
	helpText := `
qcli supported only Quay.io
Usage: delete-user
  qcli delete-user quay.io/koudaiii/qcli koudaiii
`
	return strings.TrimSpace(helpText)
}
