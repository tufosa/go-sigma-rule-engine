/*
Copyright © 2020 Markus Kont alias013@gmail.com

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
package cmd

import (
	"github.com/markuskont/go-sigma-rule-engine/pkg/sigma/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type counts struct {
	ok, fail, unsupported int
}

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse a ruleset for testing",
	Long:  `Recursively parses a sigma ruleset from filesystem and provides detailed feedback to the user about rule support.`,
	Run:   entrypoint,
}

func entrypoint(cmd *cobra.Command, args []string) {
	files, err := sigma.NewRuleFileList(viper.GetStringSlice("sigma.rules.dir"))
	if err != nil {
		logrus.Fatal(err)
	}
	for _, f := range files {
		logrus.Info(f)
	}
	logrus.Info("Parsing rule yaml files")
	rules, err := sigma.NewRuleList(files, true)
	if err != nil {
		switch err.(type) {
		case sigma.ErrBulkParseYaml:
			logrus.Error(err)
		default:
			logrus.Fatal(err)
		}
	}
	logrus.Infof("Got %d rules from yaml", len(rules))
	logrus.Info("Parsing rules into AST")
	c := &counts{}
	for _, raw := range rules {
		_, err := sigma.NewTree(&raw)
		if err != nil {
			logrus.Errorf("%s: %s", raw.Path, err)
			c.fail++
		} else {
			logrus.Infof("%s: ok", raw.Path)
			c.ok++
		}
	}
	logrus.Infof("OK: %d; FAIL: %d; UNSUPPORTED: %d", c.ok, c.fail, c.unsupported)
	/*
		for _, rule := range rules {
			if val, ok := rule.Detection["condition"].(string); ok {
				logrus.Info(val)
			} else if rule.Multipart {
				logrus.Warnf("%s is multipart", rule.Path)
			} else {
				logrus.Errorf("%s missing condition or not string", rule.Path)
			}
		}
	*/
	/*
		contextLogger := log.WithFields(log.Fields{
			"ok":          r.Total,
			"errors":      len(r.Broken),
			"unsupported": len(r.Unsupported),
		})
		contextLogger.Info("Done")
	*/
}

func init() {
	rootCmd.AddCommand(parseCmd)

	parseCmd.PersistentFlags().StringSlice("sigma-rules-dir", []string{},
		"Directories that contains sigma rules.")
	viper.BindPFlag("sigma.rules.dir", parseCmd.PersistentFlags().Lookup("sigma-rules-dir"))
}
