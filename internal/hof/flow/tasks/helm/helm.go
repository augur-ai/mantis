package helm

import (
	"fmt"

	"cuelang.org/go/cue"
	hofcontext "github.com/opentofu/opentofu/internal/hof/flow/context"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
)

// HelmTask is a task for deploying Helm charts
type HelmTask struct {
}

func NewHelmTask(val cue.Value) (hofcontext.Runner, error) {
	return &HelmTask{}, nil
}

func (t *HelmTask) Run(ctx *hofcontext.Context) (any, error) {
	v := ctx.Value
	chartConfig := v.LookupPath(cue.ParsePath("config"))

	// Extract necessary information from the CUE value
	releaseName, _ := chartConfig.LookupPath(cue.ParsePath("releaseName")).String()
	chartName, _ := chartConfig.LookupPath(cue.ParsePath("chartName")).String()
	namespace, _ := chartConfig.LookupPath(cue.ParsePath("namespace")).String()

	// Initialize Helm client
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, "", nil); err != nil {
		return nil, fmt.Errorf("failed to initialize Helm action config: %v", err)
	}

	if ctx.Plan {
		// Perform a dry-run installation
		client := action.NewInstall(actionConfig)
		client.DryRun = true
		client.ReleaseName = releaseName
		client.Namespace = namespace
		chart, err := loader.Load(chartName)
		if err != nil {
			return nil, fmt.Errorf("failed to load chart: %v", err) // Return nil and error
		}

		_, err = client.Run(chart, nil)
		if err != nil {
			return nil, fmt.Errorf("dry-run failed: %v", err) // Return nil and error
		}

		client.Run(chart, nil)

		return "Helm chart dry-run successful", nil
	} else if ctx.Apply {
		// Perform actual installation
		client := action.NewInstall(actionConfig)
		client.ReleaseName = releaseName
		client.Namespace = namespace
		chart, err := loader.Load(chartName)
		if err != nil {
			return nil, fmt.Errorf("failed to load chart: %v", err) // Return nil and error
		}

		_, err = client.Run(chart, nil)
		if err != nil {
			return nil, fmt.Errorf("apply failed: %v", err) // Return nil and error
		}

		return "Helm chart installed successfully", nil
	} else if ctx.Destroy {
		// Uninstall the Helm release
		client := action.NewUninstall(actionConfig)

		_, err := client.Run(releaseName)
		if err != nil {
			return nil, fmt.Errorf("failed to uninstall Helm release: %v", err)
		}

		return "Helm release uninstalled successfully", nil
	}

	return nil, fmt.Errorf("unknown command. Need to use one of plan/apply/destroy")
}
