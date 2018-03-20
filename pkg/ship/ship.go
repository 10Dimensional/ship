package ship

import (
	"context"

	"fmt"
	"os"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/mitchellh/cli"
	"github.com/pkg/errors"
	"github.com/replicatedcom/ship/pkg/lifecycle"
	"github.com/replicatedcom/ship/pkg/logger"
	"github.com/replicatedcom/ship/pkg/specs"
	"github.com/replicatedcom/ship/pkg/version"
	"github.com/spf13/viper"
)

// Ship configures an application
type Ship struct {
	Logger kitlog.Logger

	Port           int
	CustomerID     string
	InstallationID string
	PlanOnly       bool

	Resolver   *specs.Resolver
	StudioFile string
	Client     *specs.GraphQLClient
	UI         cli.Ui
}

// ResolverFromViper gets an instance using viper to pull config
func FromViper(v *viper.Viper) (*Ship, error) {
	graphql, err := specs.GraphQLClientFromViper(v)
	if err != nil {
		return nil, errors.Wrap(err, "get graphql client")
	}
	return &Ship{
		Logger:   logger.FromViper(v),
		Resolver: specs.ResolverFromViper(v),
		Client:   graphql,

		Port:           v.GetInt("port"),
		CustomerID:     v.GetString("customer_id"),
		InstallationID: v.GetString("installation_id"),
		StudioFile:     v.GetString("studio_file"),

		UI: &cli.ColoredUi{
			OutputColor: cli.UiColorNone,
			ErrorColor:  cli.UiColorRed,
			WarnColor:   cli.UiColorYellow,
			InfoColor:   cli.UiColorGreen,
			Ui: &cli.BasicUi{
				Reader:      os.Stdin,
				Writer:      os.Stdout,
				ErrorWriter: os.Stderr,
			},
		},
	}, nil
}

// Execute starts ship
func (d *Ship) Execute(ctx context.Context) error {
	debug := level.Debug(kitlog.With(d.Logger, "method", "execute"))

	debug.Log("method", "configure", "phase", "initialize",
		"version", version.GitSHA(),
		"buildTime", version.BuildTime(),
		"buildTimeFallback", version.GetBuild().TimeFallback,
		"customer_id", d.CustomerID,
		"installation_id", d.InstallationID,
		"plan_only", d.PlanOnly,
		"studio_file", d.StudioFile,
		"studio", specs.AllowInlineSpecs,
		"port", d.Port,
	)

	debug.Log("phase", "validate-inputs")

	if d.StudioFile != "" && !specs.AllowInlineSpecs {
		debug.Log("phase", "load-specs", "error", "unsupported studio_file")
		return errors.New("unsupported configuration: studio_file")

	}

	spec, err := d.Resolver.ResolveSpecs(ctx)
	if err != nil {
		return errors.Wrap(err, "resolve specs")
	}

	// execute lifecycle
	lc := &lifecycle.Runner{
		CustomerID:     d.CustomerID,
		InstallationID: d.InstallationID,
		GraphQLClient:  d.Client,
		UI:             d.UI,
		Logger:         d.Logger,
		Spec:           spec,
	}

	lc.Run(ctx)

	return nil

}

func (d *Ship) OnError(err error) {
	d.UI.Error(fmt.Sprintf("There was an unexpected error! %+v", err))
	d.UI.Output("")
	d.UI.Info("There was an error configuring the application. Please re-run with --log_level=debug and include the output in any support inquiries.")
}
