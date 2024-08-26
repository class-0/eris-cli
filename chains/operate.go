package chains

import (
	"fmt"
	"strings"

	"github.com/eris-ltd/eris-cli/services"

	. "github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/eris-ltd/common"
	"github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/spf13/cobra"
)

//----------------------------------------------------------------------

func Start(cmd *cobra.Command, args []string) {
	IfExit(checkChainGiven(args))
	IfExit(StartChainRaw(args[0]))
}

func Logs(cmd *cobra.Command, args []string) {
	IfExit(checkChainGiven(args))
	IfExit(LogsChainRaw(args[0]))
}

func Exec(cmd *cobra.Command, args []string) {
	IfExit(checkChainGiven(args))
	srv := args[0]
	args = args[1:]
	if len(args) == 1 {
		args = strings.Split(args[0], " ")
	}
	IfExit(ExecChainRaw(srv, args, cmd.Flags().Lookup("interactive").Changed))
}

func Kill(cmd *cobra.Command, args []string) {
	IfExit(checkChainGiven(args))
	IfExit(KillChainRaw(args[0]))
}

//----------------------------------------------------------------------

func StartChainRaw(chainName string) error {
	chain, err := LoadChainDefinition(chainName)
	if err != nil {
		return err
	}
	if IsChainRunning(chain) {
		logger.Infoln("Chain already started. Skipping.")
	} else {
		err := services.StartServiceByService(chain.Service, chain.Operations)
		if err != nil {
			return err
		}
	}
	return nil
}

func LogsChainRaw(chainName string) error {
	chain, err := LoadChainDefinition(chainName)
	if err != nil {
		return err
	}
	err = services.LogsServiceByService(chain.Service, chain.Operations)
	if err != nil {
		return err
	}
	return nil
}

func ExecChainRaw(name string, args []string, attach bool) error {
	chain, err := LoadChainDefinition(name)
	if err != nil {
		return err
	}

	if IsChainExisting(chain) {
		logger.Infoln("Chain exists.")
		return services.ExecServiceByService(chain.Service, chain.Operations, args, attach)
	} else {
		return fmt.Errorf("Chain does not exist. Please start the chain container with eris chains start %s.\n", name)
	}

	return nil
}

func KillChainRaw(chainName string) error {
	chain, err := LoadChainDefinition(chainName)
	if err != nil {
		return err
	}

	if IsChainRunning(chain) {
		err := services.KillServiceByService(chain.Service, chain.Operations)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Chain not currently running. Skipping.")
	}
	return nil
}
