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

type patchCmdOptions struct {
	ufm.IBNetwork
	PkeyString  string
	FieldStr    string
	StrategyStr string
}

var patchCmdOpt = patchCmdOptions{}

// patchCmd represents the update command
var patchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Patch IB network in UFM",
	Long:  `Patch IB network in UFM`,
	Run: func(cmd *cobra.Command, args []string) {
		ufmClient, err := ufm.NewUFM()
		if err != nil {
			fmt.Printf("Failed to connect to UFM: %v\n", err)
			os.Exit(1)
		}

		field := ufm.ParseField(patchCmdOpt.FieldStr)
		if field == ufm.UnknownField {
			fmt.Printf("Failed to update IB network in UFM: unknown filed (%s)\n", field)
			os.Exit(1)
		}

		op := ufm.ParseStrategy(patchCmdOpt.StrategyStr)
		if op == ufm.UnknownStrategy {
			fmt.Printf("Failed to update IB network in UFM: unknown strategy (%s)\n", op)
			os.Exit(1)
		}

		pkey, err := ufm.ParsePkey(patchCmdOpt.PkeyString)
		if err != nil {
			fmt.Printf("Failed to get IB network in UFM: %v\n", err)
			os.Exit(1)
		}
		patchCmdOpt.IBNetwork.PKey = pkey

		_, ufmErr := ufmClient.GetIBNetwork(pkey)
		if ufmErr != nil {
			fmt.Printf("Failed to get IB network in UFM: %v\n", ufmErr)
			os.Exit(1)
		}

		if ufmErr := ufmClient.Patch(&patchCmdOpt.IBNetwork, field, op); ufmErr != nil {
			fmt.Printf("Failed to update IB network in UFM: %v\n", ufmErr)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(patchCmd)

	patchCmd.Flags().StringVar(&patchCmdOpt.PkeyString, "pkey", "-1", "The pkeys for IB network.")
	createCmd.MarkFlagRequired("pkey")
	patchCmd.Flags().StringVar(&patchCmdOpt.FieldStr, "field", "guid", "The field of IB network to patch, one of 'qos' or 'guid'.")
	createCmd.MarkFlagRequired("field")
	patchCmd.Flags().StringVar(&patchCmdOpt.StrategyStr, "strategy", "add", "The strategy of path, one of 'add', 'delete' or 'set'.")
	createCmd.MarkFlagRequired("strategy")
	patchCmd.Flags().StringSliceVar(&patchCmdOpt.GUIDs, "guids", []string{}, "The GUID list of the IB network.")
	patchCmd.Flags().Int32Var(&patchCmdOpt.MTU, "mtu", 2048, "The MTU of the services, one of 2k or 4k.")
	patchCmd.Flags().BoolVar(&patchCmdOpt.IPOverIB, "ip-over-ib", true, "Enable IPoIB.")
	patchCmd.Flags().BoolVar(&patchCmdOpt.Index0, "index0", false, "Store the PKey at index 0 of the PKey table of the GUID.")
	patchCmd.Flags().Int32Var(&patchCmdOpt.ServiceLevel, "service-level", 0, "The service level of IB network, value can be range from 0-15")
	patchCmd.Flags().Float64Var(&patchCmdOpt.RateLimit, "rate-limit", 2.5, "The rate limit of IB network, can be one of the following: 2.5, 10, 30, 5, 20, 40, 60, 80, 120, 14, 56, 112, 168, 25, 100, 200, or 300")
}
