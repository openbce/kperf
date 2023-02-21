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
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/openbce/kperf/pkg/ufm"
)

// versionCmd represents the list command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the release version of UFM",
	Long:  `Show the release version of UFM`,
	Run: func(cmd *cobra.Command, args []string) {
		ufmClient, err := ufm.NewUFM()
		if err != nil {
			fmt.Printf("Failed to connect to UFM: %v\n", err)
			os.Exit(1)
		}
		ver, ufmErr := ufmClient.Version()
		if ufmErr != nil {
			fmt.Printf("Failed to get version of UFM: %v\n", ufmErr)
			os.Exit(1)
		}
		fmt.Printf("UFM Version: %s\n", ver)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
