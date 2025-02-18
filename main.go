/*
   Copyright 2020 Docker Compose CLI authors

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

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/docker/compose-switch/redirect"
	"github.com/spf13/cobra"
)

func main() {
	configFile := redirect.GetConfigFile()
	err := configFile.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load docker configuration: %s\n", err)
		os.Exit(1)
	}

	root := &cobra.Command{
		DisableFlagParsing: true,
		Use:                "docker-compose",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 && args[0] == "disable-v2" {
				configFile.SetFeature("composeV2", redirect.Disabled)
				err := configFile.Write()
				return err
			}
			if len(args) > 0 && args[0] == "enable-v2" {
				configFile.SetFeature("composeV2", redirect.Enabled)
				err := configFile.Write()
				return err
			}

			useV2, ok := configFile.GetFeature(redirect.ComposeV2)

			if ok && strings.HasSuffix(useV2, redirect.Enabled) {
				compose := redirect.Convert(args)
				redirect.RunComposeV2(compose)
				return nil
			}
			redirect.RunComposeV1(args)
			return nil
		},
	}

	err = root.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
