/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/RayLuxembourg/estruct/internal"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "estruct",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		extensions, _ := cmd.Flags().GetString("extensions")
		name, _ := cmd.Flags().GetString("name")
		output, _ := cmd.Flags().GetString("output")

		regexExtension := fmt.Sprintf(`(.(%s))$`,extensions)
		labels:= make([]internal.Label,0) // or create real labels
		p:= internal.NewConfig(path, regexExtension,labels)
		relativePath := "./src"
		fmt.Println("Processing...")
		datasets, _, _ := p.Init(relativePath)

		os.Remove(name)
		b, _ := json.Marshal(datasets)

		ioutil.WriteFile(name, b, 0666)
		succMsg:= fmt.Sprintf("Output save in %s/%s",output,name)
		fmt.Println(succMsg)

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "default .estructrc.json path (./estructrc.json)")
	currentPath,_:=os.Getwd()
	rootCmd.PersistentFlags().StringP("path", "p", currentPath, "directory to parse")
	rootCmd.PersistentFlags().StringP("extensions", "e", `js|jsx`, "regex to match file extensions to parse")
	rootCmd.PersistentFlags().StringP("output", "o", currentPath, "output path for the generated json")
	rootCmd.PersistentFlags().StringP("name", "n", "structure.json", "output file name")

	viper.BindPFlag("path",rootCmd.PersistentFlags().Lookup("path"))
	viper.BindPFlag("extensions",rootCmd.PersistentFlags().Lookup("extensions"))
	viper.BindPFlag("output",rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("name",rootCmd.PersistentFlags().Lookup("name"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".estruct" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".estruct")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
