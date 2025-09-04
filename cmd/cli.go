package main

import (
	"github.com/lazygophers/aiload/internal/api"
	"github.com/lazygophers/aiload/internal/core"
	"github.com/lazygophers/aiload/internal/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/app"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: app.Name,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.SetTrace()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		pterm.Info.Printfln("当前软件版本：%s", app.Version)

		err := state.Load()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		api.RegisteApi(Routes)

		err = core.Load()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		return nil
	},
}

func Run() (err error) {
	cobra.EnablePrefixMatching = true
	cobra.EnableCommandSorting = true
	cobra.EnableTraverseRunHooks = true

	err = rootCmd.Execute()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
