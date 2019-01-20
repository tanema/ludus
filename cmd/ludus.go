package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tanema/ludus/src/build"
)

const version = "ludus 0.0.1"

// LudusCmd is the entrypoint command of ludus
var LudusCmd = &cobra.Command{
	Use:           "ludus",
	Short:         "A love2d build and release tool",
	Long:          `ludus is aimed at making it easy to build love into a distributable package and run that package constantly so that your are aware of all your bugs`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build all versions of love",
	Long:  `build will build windows(32|64), macos, and linux deb files for your current project`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return build.Build()
	},
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "remove all build artifacts",
	Long:  `clean will remove the build directory to clean up any old builds`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return build.Clean()
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run your game",
	Long:  `run will build your game for the current OS and run it`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return build.Run()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of wrp",
	Long:  `All software has versions. This is ludus' version.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("wrp %s %s/%s", version, runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	pwd, _ := os.Getwd()
	defaultTitle := filepath.Base(pwd)
	defaultIdent := strings.ToLower(strings.Replace(defaultTitle, " ", "", -1))

	viper.SetConfigName("ludus")
	viper.AddConfigPath(".")

	LudusCmd.PersistentFlags().StringP("title", "t", defaultTitle, "title of the game")
	LudusCmd.PersistentFlags().StringP("love_version", "l", "", "version of love to build")
	LudusCmd.PersistentFlags().StringP("version", "v", "", "version of your game")
	LudusCmd.PersistentFlags().StringP("author", "a", "", "author name")
	LudusCmd.PersistentFlags().StringP("email", "e", "", "email of the author")
	LudusCmd.PersistentFlags().StringP("homepage", "p", "", "homepage for the game")
	LudusCmd.PersistentFlags().StringP("description", "d", "", "short description of the game")
	LudusCmd.PersistentFlags().StringP("identifier", "i", "com.love."+defaultIdent, "short description of the game")
	LudusCmd.PersistentFlags().StringP("source_directory", "s", ".", "directory that contains your source relative to your config file")
	LudusCmd.PersistentFlags().StringP("build_directory", "b", "build", "directory where builds are outputted relative to your config file")

	bindings := []string{
		"title", "love_version", "version", "author", "email",
		"homepage", "description", "identifier", "source_directory", "build_directory",
	}
	for _, bind := range bindings {
		viper.BindPFlag(bind, LudusCmd.PersistentFlags().Lookup(bind))
	}

	LudusCmd.AddCommand(runCmd, cleanCmd, buildCmd, versionCmd)
}
