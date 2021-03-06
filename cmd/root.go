// Copyright © 2017 Alexander Pinnecke <alexander.pinnecke@googlemail.com>

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/apinnecke/go-exitcontext"
	"github.com/gordonklaus/portaudio"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/StageAutoControl/controller/pkg/artnet"
	"github.com/StageAutoControl/controller/pkg/disk"
)

var (
	logLevel          string
	logger            *logrus.Entry
	storagePath       string
	storage           *disk.Storage
	loader            *disk.Loader
	controller        artnet.Controller
	disableController bool
	ctx               context.Context
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "controller",
	Short: "Stage automatic controlling, triggering state changes.",
	Long:  `Automatic stage controlling, including midi and DMX, by analyzing audio signals and pre defined light scenes`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ctx = exitcontext.New()
		logger = createLogger(logLevel)
		storage = createStorage(logger, storagePath)
		controller = createController(logger, disableController)
		loader = disk.NewLoader(storage)

		if err := portaudio.Initialize(); err != nil {
			logger.Fatalf("failed to initialize portaudio: %v", err)
		}
		go func() {
			<-ctx.Done()
			terminateAudio()
		}()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		terminateAudio()
	},
}

func terminateAudio() {
	if err := portaudio.Terminate(); err != nil {
		logger.Errorf("failed to terminate portaudio: %v", err)
	}
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Which log level to use")
	RootCmd.PersistentFlags().BoolVar(&disableController, "disable-controller", false, "Disable the controller, e.g. when not on an artnet network")
	RootCmd.PersistentFlags().StringVarP(&storagePath, "storage-path", "s", "/var/controller/data", "path where the storage should store the data")
}
