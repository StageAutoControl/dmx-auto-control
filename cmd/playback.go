// Copyright © 2017 Alexander Pinnecke <alexander.pinnecke@googlemail.com>

package cmd

import (
	"fmt"
	"os"

	"github.com/StageAutoControl/controller/cntl"
	"github.com/StageAutoControl/controller/cntl/playback"
	"github.com/StageAutoControl/controller/cntl/transport"
	"github.com/StageAutoControl/controller/cntl/waiter"
	"github.com/StageAutoControl/controller/database/files"
	"github.com/spf13/cobra"
)

const (
	playbackTypeSong    = "song"
	playbackTypeSetList = "setlist"
)

var (
	transportTypes = []string{
		transport.TYPE_BUFFER,
		transport.TYPE_VISUALIZER,
		transport.TYPE_ARTNET,
		transport.TYPE_MIDI,
	}
	usedTransports    []string
	viualizerEndpoint string
	midiDeviceID      string

	waiterTypes = []string{
		waiter.TYPE_NONE,
		waiter.TYPE_AUDIO,
	}
	usedWaiters []string
)

// playbackCmd represents the playback command
var playbackCmd = &cobra.Command{
	Use:   "playback [song|setlist] song-valid-uuid-1",
	Short: "Plays a given Song or SetList by id",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			cmd.Usage()
			os.Exit(1)
		}

		var loader cntl.Loader
		switch loaderType {
		case directoryLoader:
			fmt.Printf("Loading data directoy %q ... \n", dataDir)
			loader = files.New(dataDir)
		case databaseLoader:
			//loader = database.New(),
			fmt.Println("Database loader is not yet supported.")
			os.Exit(1)

		default:
			fmt.Printf("Loader %q is not supported. Choose one of %s \n", loader, loaders)
			os.Exit(1)
		}

		data, err := loader.Load()
		if err != nil {
			fmt.Printf("Failed to load data from %q: %v \n", loaderType, err)
			os.Exit(1)
		}

		var writers []playback.TransportWriter
		for _, transportType := range usedTransports {
			switch transportType {
			case transport.TYPE_BUFFER:
				writers = append(writers, transport.NewBuffer(os.Stdout))
				break

			case transport.TYPE_VISUALIZER:
				w, err := transport.NewVisualizer(viualizerEndpoint)
				if err != nil {
					fmt.Printf("Unable to connect to the visualizer: %v \n", err)
					os.Exit(1)
				}

				writers = append(writers, w)
				break

			case transport.TYPE_ARTNET:
				w, err := transport.NewArtNet("stage-auto-control")
				if err != nil {
					fmt.Printf("Unable to connect to the visualizer: %v \n", err)
					os.Exit(1)
				}

				writers = append(writers, w)
				break

			case transport.TYPE_MIDI:
				w, err := transport.NewMIDI(midiDeviceID)
				if err != nil {
					fmt.Printf("Unable to connect to midi device: %v \n", err)
					os.Exit(1)
				}

				writers = append(writers, w)
				break

			default:
				fmt.Printf("Transport %q is not supported. \n", transportType)
				os.Exit(1)
			}
		}

		var waiters []playback.Waiter
		for _, waiterType := range usedWaiters {
			switch waiterType {
			case waiter.TYPE_NONE:

				waiters = append(waiters, waiter.NewNone())
			}
		}

		player := playback.NewPlayer(data, writers, waiters)

		switch args[0] {
		case playbackTypeSong:
			songId := args[1]

			if err = player.PlaySong(songId); err != nil {
				panic(err)
			}

			break
		case playbackTypeSetList:
			setListId := args[2]

			if err = player.PlaySetList(setListId); err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(playbackCmd)

	playbackCmd.PersistentFlags().StringSliceVarP(&usedTransports, "transport", "t", []string{transport.TYPE_BUFFER}, fmt.Sprintf("Which usedTransports to use from %s.", transportTypes))
	playbackCmd.PersistentFlags().StringVar(&viualizerEndpoint, "visualizer-endpoint", "localhost:1337", "Endpoint of the visualizer backend if visualizer transport is chosen.")
	playbackCmd.PersistentFlags().StringVarP(&midiDeviceID, "midi-device-id", "m", "", "DeviceID of MIDI output to use (On empty string the default device is used)")

	playbackCmd.PersistentFlags().StringSliceVarP(&usedWaiters, "wait-for", "w", []string{waiter.TYPE_NONE}, fmt.Sprintf("Wait for a specific signal before playing a song (required to be used on stage, otherwise the next song would start immediately), one of %s", waiterTypes))
}
