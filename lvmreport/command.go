package lvmreport

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/kballard/go-shellquote"
)

var defaultConfig = []string{
	"global/suffix=false",
	"global/units=b",

	"report/binary_values_as_numeric=1",
	"report/output_format=json",
	"report/buffered=false",
	"report/time_format=%s",
}

type Command struct {
	args []string
}

func NewCommand(args []string) *Command {
	return &Command{
		args: args,
	}
}

func (c *Command) completeArgs() []string {
	args := append([]string(nil), c.args...)

	args = append(args,
		"fullreport",
		"--all",
		"--config", strings.Join(defaultConfig, " "),
	)

	for _, g := range AllGroupNames {
		args = append(args,
			"--configreport", g.String(),
			"--options", strings.Join(g.fields(), ","),
		)
	}

	return args
}

func prepareCommand(ctx context.Context, args []string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Env = append(os.Environ(),
		"LVM_SUPPRESS_SYSLOG=1",
	)

	return cmd
}

func runCommand(ctx context.Context, args []string) (*ReportData, error) {
	cmd := prepareCommand(ctx, args)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = os.Stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	var data *ReportData

	if err = cmd.Start(); err != nil {
		err = fmt.Errorf("start failed: %w", err)
	} else {
		r := newReader(stdout)
		r.Decode()

		if err = cmd.Wait(); err == nil {
			data, err = r.Data()
		}
	}

	if err != nil {
		log.Printf("Command failed with %q: %s", err, shellquote.Join(cmd.Args...))
		return nil, err
	}

	return data, nil
}

func (c *Command) String() string {
	return shellquote.Join(c.completeArgs()...)
}

func (c *Command) Run(ctx context.Context) (*ReportData, error) {
	return runCommand(ctx, c.completeArgs())
}
