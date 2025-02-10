package main

import (
	"context"
	"github.com/centralmind/gateway/pkg/api"
	"github.com/centralmind/gateway/pkg/config"
	"github.com/doublecloud/transfer/library/go/core/xerrors"
	"github.com/doublecloud/transfer/pkg/abstract"
	"github.com/doublecloud/transfer/pkg/abstract/model"
	"github.com/doublecloud/transfer/pkg/worker/tasks"
	"net/http"
	"os"

	"github.com/doublecloud/transfer/pkg/cobraaux"
	"github.com/spf13/cobra"
)

func StartCommand() *cobra.Command {
	var gatewayParams string
	var port string
	checkCommand := &cobra.Command{
		Use:   "start",
		Short: "Start gateway locally",
		Args:  cobra.MatchAll(cobra.ExactArgs(0)),
		RunE:  run(&gatewayParams, &port),
	}
	checkCommand.Flags().StringVar(&gatewayParams, "config", "./gateway.yaml", "path to yaml file with gateway configuration")
	checkCommand.Flags().StringVar(&port, "port", "9090", "port for gateway")
	return checkCommand
}

func run(configPath *string, port *string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		gwRaw, err := os.ReadFile(*configPath)
		if err != nil {
			return xerrors.Errorf("unable to read yaml config file: %w", err)
		}
		gw, err := config.FromYaml(gwRaw)
		if err != nil {
			return xerrors.Errorf("unable to parse config file: %w", err)
		}

		endpoint, err := gw.Endpoint()
		if err != nil {
			return xerrors.Errorf("endpoint is invalid: %w", err)
		}
		transfer := &model.Transfer{
			Type: abstract.TransferTypeSnapshotOnly,
			Src:  endpoint,
			Dst:  new(model.MockDestination),
		}
		res := tasks.TestEndpoint(context.Background(), &tasks.TestEndpointParams{
			Transfer:             transfer,
			TransformationConfig: nil,
			ParamsSrc: &tasks.EndpointParam{
				Type:  transfer.SrcType(),
				Param: transfer.SrcJSON(),
			},
			ParamsDst: nil,
		}, abstract.NewTestResult())

		mux := http.NewServeMux()
		api.NewAPI(res.Schema, res.Preview).RegisterRoutes(mux)
		return http.ListenAndServe(":8080", mux)
	}
}

func main() {
	rootCommand := &cobra.Command{
		Use:          "gateway",
		Short:        "gateway cli",
		Example:      "./gateway help",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	cobraaux.RegisterCommand(rootCommand, StartCommand())
	err := rootCommand.Execute()
	if err != nil {
		os.Exit(1)
	}
}
