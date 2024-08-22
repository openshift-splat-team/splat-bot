package util

import (
	"context"

	"github.com/apenella/go-ansible/v2/pkg/execute"
	"github.com/apenella/go-ansible/v2/pkg/playbook"
	"github.com/go-git/go-git/v5"
)

func Clone(url, path string) error {
	if _, err := git.PlainClone(path, false, &git.CloneOptions{URL: url}); err != nil {
		return err
	}

	return nil
}

func ExecuteAnsiblePlaybooks(ctx context.Context, playbooks []string) error {
	// leaving these options here just in case we need them.
	options := &playbook.AnsiblePlaybookOptions{
		AskVaultPassword:  false,
		Check:             false,
		Diff:              false,
		ExtraVars:         nil,
		ExtraVarsFile:     nil,
		FlushCache:        false,
		ForceHandlers:     false,
		Forks:             "",
		Inventory:         "",
		Limit:             "",
		ListHosts:         false,
		ListTags:          false,
		ListTasks:         false,
		ModulePath:        "",
		SkipTags:          "",
		StartAtTask:       "",
		Step:              false,
		SyntaxCheck:       false,
		Tags:              "",
		VaultID:           "",
		VaultPasswordFile: "",
		Verbose:           false,
		VerboseV:          false,
		VerboseVV:         false,
		VerboseVVV:        false,
		VerboseVVVV:       false,
		Version:           false,
		AskPass:           false,
		Connection:        "",
		PrivateKey:        "",
		SCPExtraArgs:      "",
		SFTPExtraArgs:     "",
		SSHCommonArgs:     "",
		SSHExtraArgs:      "",
		Timeout:           0,
		User:              "",
		AskBecomePass:     false,
		Become:            false,
		BecomeMethod:      "",
		BecomeUser:        "",
	}

	apbcmd := playbook.NewAnsiblePlaybookCmd(
		playbook.WithPlaybooks(playbooks...),
		playbook.WithPlaybookOptions(options),
	)

	exec := execute.NewDefaultExecute(
		execute.WithCmd(apbcmd),
		execute.WithErrorEnrich(playbook.NewAnsiblePlaybookErrorEnrich()),
	)

	if err := exec.Execute(ctx); err != nil {
		return err
	}

	return nil
}
