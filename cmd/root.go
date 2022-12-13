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
	"fmt"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/pepabo/cwlf/datasource"
	"github.com/pepabo/cwlf/datasource/fake"
	"github.com/pepabo/cwlf/datasource/local"
	"github.com/pepabo/cwlf/datasource/s3"
	"github.com/pepabo/cwlf/filter"
	"github.com/pepabo/cwlf/parser"
	"github.com/pepabo/cwlf/parser/rdsaudit"
	"github.com/spf13/cobra"
)

const defaultRegion = "ap-northeast-1"

var rootCmd = &cobra.Command{
	Use:   "cwlf [DATASOURCE_DSN]",
	Short: "CloudWatch Logs Filter",
	Long:  `CloudWatch Logs Filter.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dsn := args[0]
		ctx := context.Background()
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return err
		}
		if cfg.Region == "" {
			cfg.Region = defaultRegion
		}

		if err != nil {
			return err
		}
		u, err := url.Parse(dsn)
		if err != nil {
			return err
		}

		// datasource
		var d datasource.Datasource
		switch u.Scheme {
		case "s3":
			d, err = s3.New(cfg, dsn)
			if err != nil {
				return err
			}
		case "local":
			d, err = local.New(dsn)
			if err != nil {
				return err
			}
		case "fake":
			d, err = fake.New(dsn)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsuppoted scheme: %s", dsn)
		}

		// parser
		var p parser.Parser
		p = rdsaudit.New()

		// filter
		f := filter.New([]string{})

		for e := range f.Filter(ctx, p.Parse(ctx, d.Fetch(ctx))) {
			fmt.Println(string(e.LogEvent.Raw))
		}

		if err := d.Err(); err != nil {
			return err
		}

		if err := p.Err(); err != nil {
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
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
