package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/koudaiii/qcli/quay"
)

type CreateCommand struct {
	Meta
}

func (c *CreateCommand) Run(args []string) int {
	if err := FlagInit(args); err != nil {
		fmt.Fprintln(os.Stderr, c.Help())
		os.Exit(1)
	}

	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, c.Help())
		os.Exit(1)
	}

	ss := strings.Split(args[0], "/")
	if len(ss) != 3 {
		fmt.Fprintln(os.Stderr, c.Help())
		os.Exit(1)
	}

	repos, err := quay.CreateRepository(ss[1], ss[2], visibility)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "Created! quay.io/%v/%v\n", repos.Namespace, repos.Name)

	return 0
}

func (c *CreateCommand) Synopsis() string {
	return fmt.Sprint("Create repository in Quay")
}

func (c *CreateCommand) Help() string {
	helpText := `
qcli supported only Quay.io
Usage: create
  qcli create quay.io/koudaiii/qcli --visibility private
`
	return strings.TrimSpace(helpText)
}
