package commands

import (
	"fmt"
	"strings"

	act "github.com/eris-ltd/eris-cli/actions"

	. "github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/eris-ltd/common"
	"github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/spf13/cobra"
)

//----------------------------------------------------------------------
// cli definitions

// Primary Actions Sub-Command
// flags to add: --global --project
var Actions = &cobra.Command{
	Use:   "actions",
	Short: "Manage and Perform Structured Actions.",
	Long: `Display and Manage actions for various components of the
Eris platform and for the platform itself.

Actions are bundles of commands which rely upon a project
which is currently in scope or on a global set of actions.
Actions are held in yaml, toml, or json action-definition
files within the action folder in the eris tree (globally
scoped actions) or in a directory pointed to by the
actions field of the currently checked out project
(project scoped actions). Actions are a sequence of
commands which operate in a similar fashion to how a
circle.yml file or a .travis.yml script field may operate.

Actions execute in a series of individual sub-shells ran
on the host. Note actions do not run from inside containers
but can interact with containers either via the installed
eris commands or via the docker cli itself or, indeed, any
other programs installed *on the host*.`,
	Run: func(cmd *cobra.Command, args []string) { cmd.Help() },
}

// Build the actions subcommand
func buildActionsCommand() {
	Actions.AddCommand(actionsNew)
	Actions.AddCommand(actionsImport)
	Actions.AddCommand(actionsList)
	Actions.AddCommand(actionsEdit)
	Actions.AddCommand(actionsDo)
	Actions.AddCommand(actionsExport)
	Actions.AddCommand(actionsRename)
	Actions.AddCommand(actionsRemove)
	addActionsFlags()
}

// Actions Sub-sub-Commands
var actionsImport = &cobra.Command{
	Use:   "import [name] [location]",
	Short: "Import an action definition file from Github or IPFS.",
	Long: `Import an action definition for your platform.

By default, Eris will import from ipfs.

To list known actions use: [eris actions known].`,
	Example: "  eris actions import \"do not use\" ipfs:QmNUhPtuD9VtntybNqLgTTevUmgqs13eMvo2fkCwLLx5MX",
	Run: func(cmd *cobra.Command, args []string) {
		ImportAction(cmd, args)
	},
}

// flags to add: template
var actionsNew = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new action definition file.",
	Long:  `Create a new action definition file optionally from a template.`,
	Run: func(cmd *cobra.Command, args []string) {
		NewAction(cmd, args)
	},
}

var actionsList = &cobra.Command{
	Use:   "ls",
	Short: "List all registered action definition files.",
	Long:  `List all registered action definition files.`,
	Run: func(cmd *cobra.Command, args []string) {
		ListActions(cmd, args)
	},
}

var actionsDo = &cobra.Command{
	Use:   "do [name]",
	Short: "Perform an action.",
	Long:  `Perform an action according to the action definition file.

Actions are used to perform functions which are a
semi-scriptable series of steps. These are general
helper functions.

Actions are a series of commands passed to a series of
*individual* subshells. These actions can take a series
of arguments.

Arguments passed into the shells via the command line
(extra arguments which do not match the name) will be
available to the command steps as $1, $2, $3, etc.

In addition, variables will be populated within the
subshell according to the key:val syntax within the
command line.

The shells will be passed the host's environment as
well as any additional env vars added to the action
definition file.
`,
	Example: `  eris actions do dns register -> will run the ~/.eris/actions/dns_register action def file
  eris actions do dns register name:cutemarm ip:111.111.111.111 -> will populate $name and $ip
  eris actions do dns register cutemarm 111.111.111.111 -> will populate $1 and $2`,
	Run: func(cmd *cobra.Command, args []string) {
		DoAction(cmd, args)
	},
}

var actionsEdit = &cobra.Command{
	Use:   "edit [name]",
	Short: "Edit an action definition file.",
	Long:  `Edit an action definition file in the default editor.`,
	Run: func(cmd *cobra.Command, args []string) {
		EditAction(cmd, args)
	},
}

var actionsExport = &cobra.Command{
	Use:   "export [chainName]",
	Short: "Export an action definition file to IPFS.",
	Long: `Export an action definition file to IPFS.

Command will return a machine readable version of the IPFS hash
`,
	Run: func(cmd *cobra.Command, args []string) {
		ExportAction(cmd, args)
	},
}

var actionsRename = &cobra.Command{
	Use:     "rename [old] [new]",
	Short:   "Rename an action.",
	Long:    `Rename an action.`,
	Example: "  eris actions rename \"old action name\" \"new action name\"",
	Run: func(cmd *cobra.Command, args []string) {
		RenameAction(cmd, args)
	},
}

var actionsRemove = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove an action definition file.",
	Long:  `Remove an action definition file.`,
	Run: func(cmd *cobra.Command, args []string) {
		RmAction(cmd, args)
	},
}

//----------------------------------------------------------------------
// cli flags
func addActionsFlags() {
	actionsRemove.Flags().BoolVarP(&Force, "force", "f", false, "force action without confirming")
	actionsDo.Flags().BoolVarP(&Quiet, "quiet", "q", false, "suppress action output")
	actionsDo.Flags().StringSliceVarP(&ServicesSlice, "services", "s", []string{}, "comma separated list of services to start")
	actionsDo.Flags().StringVarP(&Chain, "chain", "c", "", "run action against a particular chain")
}

//----------------------------------------------------------------------
// cli command wrappers

func ImportAction(cmd *cobra.Command, args []string) {
	if err := checkActionGiven(args); err != nil {
		cmd.Help()
		return
	}
	if len(args) != 2 {
		cmd.Help()
		return
	}

	IfExit(act.ImportActionRaw(args[0], args[1]))
}

func NewAction(cmd *cobra.Command, args []string) {
	if err := checkActionGiven(args); err != nil {
		cmd.Help()
		return
	}

	IfExit(act.NewActionRaw(args))
}

func ListActions(cmd *cobra.Command, args []string) {
	// TODO: add scoping for when projects done.
	actions := act.ListKnownRaw()
	for _, s := range actions {
		logger.Println(strings.Replace(s, "_", " ", -1))
	}
}

func EditAction(cmd *cobra.Command, args []string) {
	if err := checkActionGiven(args); err != nil {
		cmd.Help()
		return
	}

	IfExit(act.EditActionRaw(args))
}

func DoAction(cmd *cobra.Command, args []string) {
	action, actionVars, err := act.LoadActionDefinition(args)
	if err != nil {
		logger.Errorln(err)
		return
	}

	if err := act.MergeStepsAndCLIArgs(action, &actionVars, args); err != nil {
		logger.Errorln(err)
		return
	}

	IfExit(act.DoRaw(action, actionVars, cmd.Flags().Lookup("quiet").Changed))
}

func ExportAction(cmd *cobra.Command, args []string) {
	if err := checkActionGiven(args); err != nil {
		cmd.Help()
		return
	}

	IfExit(act.ExportActionRaw(args))
}

func RenameAction(cmd *cobra.Command, args []string) {
	if err := checkActionGiven(args); err != nil {
		cmd.Help()
		return
	}
	if len(args) != 2 {
		cmd.Help()
		return
	}

	IfExit(act.RenameActionRaw(args[0], args[1]))
}

func RmAction(cmd *cobra.Command, args []string) {
	if err := checkActionGiven(args); err != nil {
		cmd.Help()
		return
	}

	IfExit(act.RmActionRaw(args, cmd.Flags().Lookup("force").Changed))
}

func checkActionGiven(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("No Service Given. Please rerun command with a known service.")
	}
	return nil
}
