// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package command

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/go-logr/zapr"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	gglog "github.com/anydotcloud/grm-generate/pkg/log"
)

const (
	appName      = "grm-generate"
	appShortDesc = "grm-generate - generate cloud resource manager packages"
	appLongDesc  = `grm-generate

A tool to generate cloud resource manager packages`
)

var (
	outputOptions = []string{
		"yaml",
		"table",
	}
)

var (
	defaultCachePath string
	optCachePath     string
	optDryRun        bool
	optDebug         bool
	optOutput        string
	log              gglog.Logger
)

var rootCmd = &cobra.Command{
	Use:               appName,
	Short:             appShortDesc,
	Long:              appLongDesc,
	PersistentPreRunE: setupLogger,
}

func init() {
	hd, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("unable to determine $HOME: %s\n", err)
		os.Exit(1)
	}
	defaultCachePath = filepath.Join(hd, ".cache", appName)

	rootCmd.PersistentFlags().BoolVar(
		&optDebug, "debug", false,
		"If true, shows debug log messages",
	)
	rootCmd.PersistentFlags().BoolVar(
		&optDryRun, "dry-run", false,
		"If true, only outputs to stdout",
	)
	rootCmd.PersistentFlags().StringVar(
		&optCachePath, "cache-path", defaultCachePath,
		"Path to directory to store cached files (including clone'd aws-sdk-go repo)",
	)
}

// setupLogger instantiates the package-level logger
func setupLogger(cmd *cobra.Command, args []string) error {
	zc := zap.NewProductionConfig()
	zc.DisableStacktrace = true
	zc.Encoding = "console"
	zc.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zc.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	if optDebug {
		zc.EncoderConfig.EncodeCaller = customCallerEncoder
		zc.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		zc.DisableCaller = true
		zc.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	z, err := zc.Build(zap.AddCallerSkip(1))
	if err != nil {
		return err
	}
	log = gglog.New(zapr.NewLogger(z))
	return nil
}

// customCallerEncoder encodes the caller filepath in a not-too-long,
// not-too-short way.
//
// grm-generate has nested packages that can be like-named and appear in
// ShortCallerEncoder as the same caller. For instance, assume a file
// cmd/grm-generate/command/discover/aws.go and a file pkg/discover/aws.go,
// both callers would appear as "discover/aws.go". This custom encoder simply
// unwinds the package path back to either "cmd" or "pkg" bases in order to
// have more fine-grained (but not the entire file path) caller identifiers.
func customCallerEncoder(
	caller zapcore.EntryCaller,
	enc zapcore.PrimitiveArrayEncoder,
) {
	// NOTE(jaypipes): The zap library is (overly) concerned about allocations
	// and performance. We are not because we're using this in a CLI tool, not
	// a service logger.
	full := caller.String()
	parts := strings.Split(full, string(filepath.Separator))
	for _, search := range []string{"pkg", "cmd"} {
		idx := lo.LastIndexOf(parts, search)
		if idx > -1 {
			enc.AppendString(filepath.Join(parts[idx:]...))
			return
		}
	}
	enc.AppendString(full)
}

// Execute adds all child commands to the root command and sets flags
// appropriately. This is called by main.main(). It only needs to happen once
// to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// newContext returns a context.Context that cancels requests when a SIGTERM or
// SIGINT signal is received and has the console logger cached in a context
// key.
func newContext(
	ctx context.Context,
) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	signalCh := make(chan os.Signal, 1)

	// recreate the context.CancelFunc
	cancelFunc := func() {
		signal.Stop(signalCh)
		cancel()
	}

	// notify on SIGINT or SIGTERM
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-signalCh:
			cancel()
		case <-ctx.Done():
		}
	}()

	// Cache the grm-generate specific logger
	ctx = context.WithValue(ctx, gglog.ContextKey, log)

	return ctx, cancelFunc
}
