/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

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
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	stringLong, stringShort string
)

func snake(k, y int) int {
	x := y - k
	for x < len(stringShort) && y < len(stringLong) && stringLong[y] == stringShort[x] {
		x++
		y++
	}
	return y
}

func max(x, y int) int {
	if x >= y {
		return x
	}
	return y
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sabun",
	Short: "show differences between 2 files",
	Long:  "show differences between 2 files",
	RunE: func(cmd *cobra.Command, args []string) error {
		// validation
		if len(args) != 2 {
			return errors.New("only 2 files are acceptible")
		}

		if len(args[0]) > len(args[1]) {
			stringLong = args[0]
			stringShort = args[1]
		} else {
			stringLong = args[1]
			stringShort = args[0]
		}

		N := len(stringLong)
		M := len(stringShort)
		delta := N - M
		offset := M + 1
		fp := make([]int, M+N+3)
		for i := range fp {
			fp[i] = -1
		}

		for p := 0; p <= M; p++ {
			for k := -p; k < delta; k++ {
				fp[k+offset] = snake(k, max(fp[k-1+offset]+1, fp[k+1+offset]))
			}
			for k := delta + p; k >= delta+1; k-- {
				fp[k+offset] = snake(k, max(fp[k-1+offset]+1, fp[k+1+offset]))
			}
			fp[delta+offset] = snake(delta, max(fp[delta-1+offset]+1, fp[delta+1+offset]))

			for k := -p; k <= delta+p; k++ {
				fmt.Printf("p=%d, k=%d, fp[%d]=%d\n", p, k, k+offset, fp[k+offset])
			}

			if fp[delta+offset] == N {
				fmt.Println("edit distance: ", delta+2*p)
				break
			}
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
