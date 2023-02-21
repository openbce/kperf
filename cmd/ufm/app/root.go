/*
Copyright 2023 The openBCE Authors.

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

package app

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ufm",
	Short: "The command line of UFM",
	Long: `ufm-client is a command line for UFM administrator to manage the IB network in UFM; the following environment
values are required to connect to the target UFM:

  UFM_USERNAME=<Username of ufm>
  UFM_PASSWORD=<Password of ufm>
  UFM_ADDRESS=<IP address or hostname of ufm server>
  UFM_PORT=<REST API port of ufm>
  UFM_HTTP_SCHEMA=<http or https>
  UFM_CERTIFICATE=<Certificate of ufm>

`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// TODO(k82cn): add a flag on log level
	// zerolog.SetGlobalLevel(zerolog.InfoLevel)
}
