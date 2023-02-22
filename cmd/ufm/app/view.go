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
	"strings"

	"github.com/spf13/cobra"

	"github.com/openbce/kperf/pkg/ufm"
)

type viewCmdOptions struct {
	PkeyStr string
	ufm.IBNetwork
}

var viewCmdOpt = viewCmdOptions{}

// viewCmd represents the list command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View the detail of a IB network in UFM",
	Long:  `View the detail of a IB network in UFM`,
	Run: func(cmd *cobra.Command, args []string) {
		ufmClient, err := ufm.NewUFM()
		if err != nil {
			fmt.Printf("Failed to connect to UFM: %v\n", err)
			os.Exit(1)
		}

		pkey, err := ufm.ParsePkey(viewCmdOpt.PkeyStr)
		if err != nil {
			fmt.Printf("Failed to get IB network in UFM: %v\n", err)
			os.Exit(1)
		}

		ib, ufmErr := ufmClient.GetIBNetwork(pkey)
		if ufmErr != nil {
			fmt.Printf("Failed to get IB network in UFM: %v\n", ufmErr)
			os.Exit(1)
		}

		ibPorts, ufmErr := ufmClient.ListPort(ufm.Filter{GUIDs: ib.GUIDs})

		fmt.Printf("%-15s: %s\n", "Name", ib.Name)
		fmt.Printf("%-15s: 0x%x\n", "Pkey", ib.PKey)
		fmt.Printf("%-15s: %t\n", "IPoIB", ib.IPOverIB)
		fmt.Printf("%-15s: %t\n", "Sharp", ib.EnableSharp)
		fmt.Printf("%-15s: %d\n", "MTU", ib.MTU)
		fmt.Printf("%-15s: %.2f\n", "Rate Limit", ib.RateLimit)
		fmt.Printf("%-15s: %d\n", "Service Level", ib.ServiceLevel)
		fmt.Printf("%-15s: %s\n", "GUIDs", strings.Join(ib.GUIDs, ","))
		fmt.Printf("%-15s:\n", "Ports")
		if len(ibPorts) == 0 {
			fmt.Printf("    %-20s%-20s%-20s%-10s%-20s%-20s\n", "Name", "GUID", "SystemID", "LID", "LogicalState", "PhysicalState")
			for _, p := range ibPorts {
				fmt.Printf("    %-20s%-20s%-20s%-10d%-20s%-20s\n", p.Name, p.GUID, p.SystemID, p.LID, p.LogicalState, p.PhysicalState)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(viewCmd)

	viewCmd.Flags().StringVar(&viewCmdOpt.PkeyStr, "pkey", "", "The pkeys for IB network.")
	createCmd.MarkFlagRequired("pkey")
}
