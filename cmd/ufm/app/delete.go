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

type deleteCmdOptions struct {
	ufm.IBNetwork
	PKeyStr string
}

var deleteCmdOpt = deleteCmdOptions{}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete IB network from UFM",
	Long:  `Delete IB network from UFM`,
	Run: func(cmd *cobra.Command, args []string) {
		ufmClient, err := ufm.NewUFM()
		if err != nil {
			fmt.Printf("Failed to connect to UFM: %v\n", err)
			os.Exit(1)
		}

		pkey, err := ufm.ParsePkey(deleteCmdOpt.PKeyStr)
		if err != nil {
			fmt.Printf("Failed to delete IB network from UFM: %v\n", err)
			os.Exit(1)
		}

		if err := ufmClient.DeleteIBNetwork(pkey); err != nil {
			fmt.Printf("Failed to delete IB network from UFM: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVar(&deleteCmdOpt.PKeyStr, "pkey", "", "The pkeys of IB network.")
	createCmd.MarkFlagRequired("pkey")
}
