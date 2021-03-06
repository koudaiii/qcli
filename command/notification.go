package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/koudaiii/qucli/quay"
)

type AddNotificationCommand struct {
	Meta
}

type DeleteNotificationCommand struct {
	Meta
}

type TestNotificationCommand struct {
	Meta
}

func (c *AddNotificationCommand) Run(args []string) int {
	if err := FlagInit(args); err != nil {
		fmt.Fprintln(os.Stderr, c.Help())
		os.Exit(1)
	}

	if len(subcommandArgs) < 1 {
		fmt.Fprintln(os.Stderr, c.Help())
		os.Exit(1)
	}

	ss := strings.Split(subcommandArgs[0], "/")
	if len(ss) != 2 {
		fmt.Fprintln(os.Stderr, c.Help())
		os.Exit(1)
	}

	// if you use add-notification command, you must set 'event' option and 'method' option.
	if event == "" || method == "" {
		fmt.Fprintln(os.Stderr, "if you use add-notification command, you must set 'event' option and 'method' option.")
		fmt.Fprintln(os.Stderr, c.Help())
		os.Exit(1)
	}

	// if you use 'vulnerability_found' event, you need set 'level' option.
	if event == "vulnerability_found" && level == "" {
		fmt.Fprintln(os.Stderr, "if you use 'vulnerability_found' event, you need set 'level' option.")
		fmt.Fprintln(os.Stderr, c.Help())
		os.Exit(1)
	}

	// if you set event excluding 'vulnerability_found' event, you can not set 'level' option.
	if event != "vulnerability_found" && level != "" {
		fmt.Fprintln(os.Stderr, "if you set event excluding 'vulnerability_found' event, you can not set 'level' option.")
		fmt.Fprintln(os.Stderr, c.Help())
		os.Exit(1)
	}

	// if you use 'email' method, you need set 'email' option.
	if method == "email" && email == "" {
		fmt.Fprintln(os.Stderr, "if you use 'email' method, you need set 'email' option.")
		fmt.Fprintln(os.Stderr, c.Help())
		os.Exit(1)
	}

	// if you use 'slack'  or 'webhook' method, you need set 'url' option.
	if (method == "slack" || method == "webhook") && url == "" {
		fmt.Fprintln(os.Stderr, "if you use 'slack'  or 'webhook' method, you need set 'url' option.")
		fmt.Fprintln(os.Stderr, c.Help())
		os.Exit(1)
	}

	// if you use 'repo_push' event, you can not set 'ref' option.
	if event == "repo_push" && ref != "" {
		fmt.Fprintln(os.Stderr, "if you use 'repo_push' event, you can not set 'ref' option.")
		fmt.Fprintln(os.Stderr, c.Help())
		os.Exit(1)
	}

	req := quay.RequestRepositoryNotification{
		Title: title,
		Event: event,
		EventConfig: quay.NotificationEventConfig{
			Level:    level,
			RefRegex: ref,
		},
		Method: method,
		Config: quay.NotificationConfig{
			Email: email,
			URL:   url,
		},
	}

	repos, err := quay.AddRepositoryNotification(ss[0], ss[1], req, hostname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "Added! \t%v\t%v\t%v\t%v\t%v\t%v\tin %v/%v/%v\n", repos.UUID, repos.Title, repos.Event, repos.EventConfig, repos.Method, repos.Config, hostname, ss[0], ss[1])
	return 0
}

func (c *AddNotificationCommand) Synopsis() string {
	return fmt.Sprint("Add notification in repository")
}

func (c *AddNotificationCommand) Help() string {
	helpText := `
qucli supported only Quay.io
Usage: add-notification
  qucli add-notification koudaiii/qucli --event="repo_push" --method="webhook" --url="http://url/goes/here"

Option:
  --event string        set 'evnet'.  ['repo_push', 'build_queued', 'build_start', 'build_success', 'build_failure', 'build_cancelled', 'vulnerability_found'].
  --level string        if you use 'vulnerability_found' evnet, A vulnerability must have a severity of the chosen level (highest level is 0).[0-6]
  --ref string          if you use event excluding 'repo_push' event, an optional regular expression for matching the git branch or tag git ref. If left blank, the notification will fire for all builds.(refs/heads/somebranch)|(refs/tags/sometag)
  --method string       set 'method'.  ['webhook', 'slack', 'email'].
  --email string        if you use 'email' method, set E-mail address. 'test@example.com'.
  --url string          if you use 'webhook' or 'slack' method, set url. 'http://url/goes/here' or 'https://hooks.slack.com/service/{some}/{token}/{here}'.
  --title string        The title for a notification is an optional field for a human-readable title for the notification.
`
	return strings.TrimSpace(helpText)
}

func (c *DeleteNotificationCommand) Run(args []string) int {
	if err := FlagInit(args); err != nil {
		fmt.Fprintln(os.Stderr, c.Help())
		os.Exit(1)
	}

	if len(subcommandArgs) < 2 {
		fmt.Fprintln(os.Stderr, c.Help())
		os.Exit(1)
	}

	ss := strings.Split(subcommandArgs[0], "/")
	if len(ss) != 2 {
		fmt.Fprintln(os.Stderr, c.Help())
		os.Exit(1)
	}

	err := quay.DeleteRepositoryNotification(ss[0], ss[1], subcommandArgs[1], hostname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "Deleted! %v notification in %v/%v/%v\n", subcommandArgs[1], hostname, ss[0], ss[1])
	return 0
}

func (c *DeleteNotificationCommand) Synopsis() string {
	return fmt.Sprint("Delete notification in repository")
}

func (c *DeleteNotificationCommand) Help() string {
	helpText := `
qucli supported only Quay.io
Usage: delete-notification
  qucli delete-notification koudaiii/qucli xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
`
	return strings.TrimSpace(helpText)
}

func (c *TestNotificationCommand) Run(args []string) int {
	if err := FlagInit(args); err != nil {
		fmt.Fprintln(os.Stderr, c.Help())
		os.Exit(1)
	}

	if len(subcommandArgs) < 2 {
		fmt.Fprintln(os.Stderr, c.Help())
		os.Exit(1)
	}

	ss := strings.Split(subcommandArgs[0], "/")
	if len(ss) != 2 {
		fmt.Fprintln(os.Stderr, c.Help())
		os.Exit(1)
	}

	err := quay.TestRepositoryNotification(ss[0], ss[1], subcommandArgs[1], hostname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "Test Notification! %v notification in %v/%v/%v\n", subcommandArgs[1], hostname, ss[0], ss[1])
	return 0
}

func (c *TestNotificationCommand) Synopsis() string {
	return fmt.Sprint("Test notification in repository")
}

func (c *TestNotificationCommand) Help() string {
	helpText := `
qucli supported only Quay.io
Usage: test-notification
  qucli test-notification koudaiii/qucli xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
`
	return strings.TrimSpace(helpText)
}
