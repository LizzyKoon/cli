package cmd

import (
	"github.com/Scalingo/cli/keys"
	"github.com/codegangsta/cli"
)

var (
	ListSSHKeyCommand = cli.Command{
		Name:  "keys",
		Usage: "List your SSH public keys",
		Description: `List all the public SSH keys associated with your account:

    $ scalingo keys
		
    # See also commands 'keys-add' and 'keys-remove'`,

		Action: func(c *cli.Context) {
			err := keys.List()
			if err != nil {
				errorQuit(err)
			}
		},
	}

	AddSSHKeyCommand = cli.Command{
		Name:  "keys-add",
		Usage: "Add a public SSH key to deploy your apps",
		Description: `Add a public SSH key:

    $ scalingo keys-add keyname /path/to/key

    # See also commands 'keys' and 'keys-remove'`,

		Action: func(c *cli.Context) {
			if len(c.Args()) != 2 {
				cli.ShowCommandHelp(c, "keys-add")
				return
			}
			err := keys.Add(c.Args()[0], c.Args()[1])
			if err != nil {
				errorQuit(err)
			}
		},
	}

	RemoveSSHKeyCommand = cli.Command{
		Name:  "keys-remove",
		Usage: "Remove a public SSH key",
		Description: `Remove a public SSH key:

    $ scalingo keys-remove keyname

    # See also commands 'keys' and 'keys-add'`,

		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "keys-remove")
				return
			}
			err := keys.Remove(c.Args()[0])
			if err != nil {
				errorQuit(err)
			}
		},
	}
)
