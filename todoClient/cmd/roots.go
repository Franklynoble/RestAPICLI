package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "todoClient",
	Short: "A Todo api client",

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },

	// Use:
	// "pScan",
	// Short: "A brief description of your application",
	// Long: `A longer description that spans multiple lines and likely contains
	// examples and usage of using your application. For example:
	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },

}

var hostsCmd = &cobra.Command{
	Use: "hosts",

	Short: "Manage the hosts list",

	Long: `Manages the hosts lists for pScan
	
	 Add hosts with the add command
	 Delete hosts with the delete command
	 List hosts with the list command.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hosts called")
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.todoClient.yaml)")

	rootCmd.PersistentFlags().String("api-root", "http://localhost:8080", "Todo API URL")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("TODO")

	viper.BindPFlag("api-root", rootCmd.PersistentFlags().Lookup("api-root"))

	rootCmd.AddCommand(hostsCmd)
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hostsCmd.PersistentFlags().String("foo", "", "A help for foo")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hostsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
