// Copyright © 2020 Dmitry Stoletov <info@imega.ru>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:     "daemon-gen",
	Short:   "Generate code microservice",
	Example: "daemon-gen -config /path/to/microservice/configuration.yaml",
	Long:    ``,
	Run:     run,
}

type generator struct {
	Config config
}

// Execute .
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
	}

	v := viper.New()

	v.SetConfigFile(args[0])

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("failed to read config: %s\n", err)
		os.Exit(1)
	}

	var g generator

	if err := v.Unmarshal(&g); err != nil {
		fmt.Printf("failed to parse config: %s\n", err)
		os.Exit(1)
	}

	f, err := os.Create("main.go")
	if err != nil {
		fmt.Printf("failed to create main.go: %s\n", err)
		os.Exit(1)
	}

	tpl := template.Must(template.New("main").Parse(indent(tplMain())))
	if err := tpl.Execute(f, g.Config); err != nil {
		fmt.Printf("failed to create main.go: %s\n", err)
		os.Exit(1)
	}
}

func indent(s string) string {
	return strings.ReplaceAll(s, "    ", "\t")
}
