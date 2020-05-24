package tfexec

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func (t *Terraform) buildTerraformCmd(args ...string) exec.Cmd {
	allArgs := []string{"terraform"}
	allArgs = append(allArgs, args...)
	allArgs = append(allArgs, "-no-color")

	var env []string
	for _, e := range os.Environ() {
		env = append(env, e)
	}

	env = append(env, "TF_LOG=") // so logging can't pollute our stderr output
	env = append(env, "TF_INPUT=0")

	return exec.Cmd{
		Path: t.execPath,
		Env:  t.Env,
		Args: allArgs,
		Dir:  t.workingDir,
	}
}

func (t *Terraform) InitCmd(args ...string) exec.Cmd {
	allArgs := []string{"init"}
	allArgs = append(allArgs, args...)

	return t.buildTerraformCmd(allArgs...)
}

func (t *Terraform) ApplyCmd(opts ...ApplyOption) exec.Cmd {
	c := &defaultApplyOptions

	for _, o := range opts {
		o.configureApply(c)
	}

	args := []string{"apply", "-auto-approve", "-input=false"}

	// string args: only pass if set
	if c.Backup != "" {
		args = append(args, "-backup="+c.Backup)
	}
	if c.LockTimeout != "" {
		args = append(args, "-lock-timeout="+c.LockTimeout)
	}
	if c.State != "" {
		args = append(args, "-state="+c.State)
	}
	if c.StateOut != "" {
		args = append(args, "-state-out="+c.StateOut)
	}
	if c.VarFile != "" {
		args = append(args, "-var-file="+c.VarFile)
	}

	// boolean and numerical args: always pass
	args = append(args, "-lock="+strconv.FormatBool(c.Lock))

	args = append(args, "-parallelism="+fmt.Sprint(c.Parallelism))
	args = append(args, "-refresh="+strconv.FormatBool(c.Refresh))

	// string slice args: pass as separate args
	if c.Targets != nil {
		for _, ta := range c.Targets {
			args = append(args, "-target="+ta)
		}
	}

	if c.Vars != nil {
		for _, v := range c.Vars {
			args = append(args, "-var '"+v+"'")
		}
	}

	return t.buildTerraformCmd(args...)
}

func (t *Terraform) ShowCmd(args ...string) exec.Cmd {
	allArgs := []string{"show", "-json"}
	allArgs = append(allArgs, args...)

	return t.buildTerraformCmd(allArgs...)
}

func (t *Terraform) ProvidersSchemaCmd(args ...string) exec.Cmd {
	allArgs := []string{"providers", "schema", "-json"}
	allArgs = append(allArgs, args...)

	return t.buildTerraformCmd(allArgs...)
}