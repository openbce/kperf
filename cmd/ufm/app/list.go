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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all IB network in UFM",
	Long:  `List all IB network in UFM`,
	Run: func(cmd *cobra.Command, args []string) {
		ufmClient, err := ufm.NewUFM()
		if err != nil {
			fmt.Printf("Failed to connect to UFM: %v\n", err)
			os.Exit(1)
		}
		ibs, ufmErr := ufmClient.ListIBNetwork()
		if ufmErr != nil {
			fmt.Printf("Failed to list IB network in UFM: %v\n", ufmErr)
			os.Exit(1)
		}

		fmt.Printf("%-15s%-10s%-10s%-10s%-10s%-10s%-10s%-10s\n", "Name", "PKey", "Sharp", "IPoIB", "MTU", "Rate", "Level", "GUID#")
		for _, ib := range ibs {
			fmt.Printf("%-15s0x%04x    %-10t%-10t%-10d%-10.2f%-10d%-10d\n",
				ib.Name,
				ib.PKey,
				ib.EnableSharp,
				ib.IPOverIB,
				ib.MTU,
				ib.RateLimit,
				ib.ServiceLevel,
				len(ib.GUIDs))
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// TODO(k82cn): add filters
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
