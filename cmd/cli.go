package main

import (
	"github.com/lazygophers/aiload/internal/api"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils"
	"github.com/lazygophers/utils/app"
	"github.com/lazygophers/utils/osx"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use: app.Name,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		pterm.EnableStyling()
		pterm.EnableColor()
		pterm.EnableOutput()

		if utils.Ignore(cmd.Flags().GetBool("verbose")) {
			env.SetVerbose(true)
		} else if utils.Ignore(cmd.Flags().GetBool("debug")) {
			pterm.DisableOutput()
			pterm.DisableStyling()
			pterm.DisableColor()

			//log.SetOutput(os.Stdout, log.GetOutputWriterHourly(filepath.Join(runtime.ExecDir(), "")))
			log.SetOutput(os.Stdout)
			log.SetLevel(log.DebugLevel)

			env.SetDebug(true)
		}

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

		err = core.Run()
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

	rootCmd.PersistentFlags().Bool("verbose", false, "详细信息")
	rootCmd.PersistentFlags().Bool("debug", false, "调试模式")

	err = rootCmd.Execute()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
