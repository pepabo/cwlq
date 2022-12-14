/*
Copyright © 2022 GMO Pepabo, inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"os"

	"github.com/pepabo/cwlq"
	"github.com/pepabo/cwlq/parser"
	"github.com/spf13/cobra"
)

const defaultRegion = "ap-northeast-1"

var (
	parserType string
	filters    []string
)

var rootCmd = &cobra.Command{
	Use:   "cwlq [DATASOURCE_DSN]",
	Short: "cwlq is a tool for querying logs (of Amazon CloudWatch Logs) stored in various datasources",
	Long:  `cwlq is a tool for querying logs (of Amazon CloudWatch Logs) stored in various datasources.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dsn := args[0]
		c, err := cwlq.New(dsn, parserType, filters)
		if err != nil {
			return err
		}
		ctx := context.Background()
		if err := c.Run(ctx); err != nil {
			return err
		}
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&parserType, "parser", "p", parser.RDSAudit, "parser for logs")
	rootCmd.Flags().StringSliceVarP(&filters, "filter", "f", []string{}, "filter for parsed logs. If multiple filters are specified, they are executed under OR conditions.")
}
