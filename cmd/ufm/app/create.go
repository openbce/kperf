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

type createCmdOptions struct {
	ufm.IBNetwork
}

var createCmdOpt = createCmdOptions{}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an IB network in UFM",
	Long:  `Create an IB network in UFM`,
	Run: func(cmd *cobra.Command, args []string) {
		ufmClient, err := ufm.NewUFM()
		if err != nil {
			fmt.Printf("Failed to connect to UFM: %v\n", err)
			os.Exit(1)
		}

		if err := ufmClient.CreateIBNetwork(&createCmdOpt.IBNetwork); err != nil {
			fmt.Printf("Failed to create IB network in UFM: %v\n", err)
			os.Exit(1)
		}

		if err := ufmClient.Patch(&createCmdOpt.IBNetwork, ufm.QoSField, ufm.SetStrategy); err != nil {
			fmt.Printf("Failed to update IB network QoS in UFM: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().Int32Var(&createCmdOpt.PKey, "pkey", 0, "The pkeys for IB network.")
	createCmd.MarkFlagRequired("pkey")
	createCmd.Flags().BoolVar(&createCmdOpt.EnableSharp, "enable-sharp", false, "Create sharp allocation accordingly")
	createCmd.Flags().StringSliceVar(&createCmdOpt.GUIDs, "guids", []string{}, "The GUID list of the IB network.")
	createCmd.MarkFlagRequired("guids")
	createCmd.Flags().Int32Var(&createCmdOpt.MTU, "mtu", 2048, "The MTU of the services, one of 2k or 4k.")
	createCmd.Flags().BoolVar(&createCmdOpt.IPOverIB, "ip-over-ib", true, "Enable IPoIB.")
	createCmd.Flags().BoolVar(&createCmdOpt.Index0, "index0", false, "Store the PKey at index 0 of the PKey table of the GUID.")
	createCmd.Flags().Int32Var(&createCmdOpt.ServiceLevel, "service-level", 0, "The service level of IB network, value can be range from 0-15")
	createCmd.Flags().Float64Var(&createCmdOpt.RateLimit, "rate-limit", 2.5, "The rate limit of IB network, can be one of the following: 2.5, 10, 30, 5, 20, 40, 60, 80, 120, 14, 56, 112, 168, 25, 100, 200, or 300")
}
