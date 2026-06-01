package fake

import (
	"context"
	"strings"

	plsvc "github.com/opendatahub-io/odh-dashboard/packages/autox-core/services/pipelines"
)

const (
	fakeTabularPipelineID   = "aaaaaaaa-0001-0001-0001-aaaaaaaaaaaa"
	fakeTabularVersionID    = "bbbbbbbb-0001-0001-0001-bbbbbbbbbbbb"
	fakeTabularPipelineName = "autogluon-tabular-training-pipeline"

	fakeTimeSeriesPipelineID   = "cccccccc-0001-0001-0001-cccccccccccc"
	fakeTimeSeriesVersionID    = "dddddddd-0001-0001-0001-dddddddddddd"
	fakeTimeSeriesPipelineName = "autogluon-timeseries-training-pipeline"

	fakeRunID          = "eeeeeeee-0001-0001-0001-eeeeeeeeeeee"
	fakeCreatedAt      = "2024-06-01T10:00:00Z"
	fakeFinishedAt     = "2024-06-01T10:45:00Z"
)

// PipelinesClient is a fake implementation of pipelines.Client for local development.
type PipelinesClient struct{}

var _ plsvc.Client = (*PipelinesClient)(nil)

func (c *PipelinesClient) ListPipelines(_ context.Context, _ string, _ string) (*plsvc.PipelinesResponse, error) {
	pipelines := []plsvc.Pipeline{
		{
			PipelineID:  fakeTabularPipelineID,
			DisplayName: fakeTabularPipelineName,
			Description: "AutoGluon tabular training pipeline (classification and regression)",
			CreatedAt:   fakeCreatedAt,
		},
		{
			PipelineID:  fakeTimeSeriesPipelineID,
			DisplayName: fakeTimeSeriesPipelineName,
			Description: "AutoGluon time-series training pipeline",
			CreatedAt:   fakeCreatedAt,
		},
	}
	return &plsvc.PipelinesResponse{Pipelines: pipelines, TotalSize: int32(len(pipelines))}, nil
}

func (c *PipelinesClient) GetPipelineVersion(_ context.Context, _ string, pipelineID, _ string) (*plsvc.PipelineVersion, error) {
	return fakePipelineVersion(pipelineID), nil
}

func (c *PipelinesClient) ListPipelineVersions(_ context.Context, _ string, pipelineID string) (*plsvc.PipelineVersionsResponse, error) {
	v := fakePipelineVersion(pipelineID)
	return &plsvc.PipelineVersionsResponse{
		PipelineVersions: []plsvc.PipelineVersion{*v},
		TotalSize:        1,
	}, nil
}

func (c *PipelinesClient) CreatePipeline(_ context.Context, _ string, name string) (*plsvc.Pipeline, error) {
	return &plsvc.Pipeline{
		PipelineID:  fakeTabularPipelineID,
		DisplayName: name,
		CreatedAt:   fakeCreatedAt,
	}, nil
}

func (c *PipelinesClient) UploadPipelineVersion(_ context.Context, _ string, pipelineID, versionName string, _ []byte) (*plsvc.PipelineVersion, error) {
	v := fakePipelineVersion(pipelineID)
	v.DisplayName = versionName
	return v, nil
}

func (c *PipelinesClient) CreatePipelineRun(_ context.Context, _ string, input *plsvc.CreatePipelineRunInput) (*plsvc.PipelineRun, error) {
	pipelineID, versionID := fakeTabularPipelineID, fakeTabularVersionID
	if input.PipelineVersionReference != nil && input.PipelineVersionReference.PipelineID != "" {
		pipelineID = input.PipelineVersionReference.PipelineID
		versionID = input.PipelineVersionReference.PipelineVersionID
	}
	return &plsvc.PipelineRun{
		RunID:       fakeRunID,
		DisplayName: input.DisplayName,
		Description: input.Description,
		State:       "RUNNING",
		CreatedAt:   fakeCreatedAt,
		PipelineVersionReference: &plsvc.PipelineVersionReference{
			PipelineID:        pipelineID,
			PipelineVersionID: versionID,
		},
		RuntimeConfig: input.RuntimeConfig,
	}, nil
}

func (c *PipelinesClient) GetPipelineRun(_ context.Context, _ string, runID string) (*plsvc.PipelineRun, error) {
	return fakeCompletedRun(runID), nil
}

func (c *PipelinesClient) ListPipelineRuns(_ context.Context, _ string, _ *plsvc.ListRunsParams) (*plsvc.PipelineRunResponse, error) {
	runs := []plsvc.PipelineRun{
		*fakeCompletedRun(fakeRunID),
	}
	return &plsvc.PipelineRunResponse{Runs: runs, TotalSize: int32(len(runs))}, nil
}

func (c *PipelinesClient) TerminateRun(_ context.Context, _ string, _ string) error { return nil }
func (c *PipelinesClient) RetryRun(_ context.Context, _ string, _ string) error      { return nil }
func (c *PipelinesClient) DeleteRun(_ context.Context, _ string, _ string) error     { return nil }

// fakePipelineVersion returns a fake pipeline version for the given pipeline ID.
func fakePipelineVersion(pipelineID string) *plsvc.PipelineVersion {
	versionID := fakeTabularVersionID
	name := fakeTabularPipelineName
	if strings.Contains(pipelineID, "cccc") {
		versionID = fakeTimeSeriesVersionID
		name = fakeTimeSeriesPipelineName
	}
	return &plsvc.PipelineVersion{
		PipelineVersionID: versionID,
		PipelineID:        pipelineID,
		DisplayName:       name + " v1.0.0",
		Description:       "Version 1.0.0",
		CreatedAt:         fakeCreatedAt,
	}
}

// fakeCompletedRun returns a fake SUCCEEDED pipeline run.
func fakeCompletedRun(runID string) *plsvc.PipelineRun {
	return &plsvc.PipelineRun{
		RunID:       runID,
		DisplayName: "Tabular Training Run",
		State:       "SUCCEEDED",
		CreatedAt:   fakeCreatedAt,
		FinishedAt:  fakeFinishedAt,
		PipelineVersionReference: &plsvc.PipelineVersionReference{
			PipelineID:        fakeTabularPipelineID,
			PipelineVersionID: fakeTabularVersionID,
		},
		RuntimeConfig: &plsvc.RuntimeConfig{
			Parameters: map[string]any{
				"target_column": "label",
				"presets":       "medium_quality",
			},
		},
		StateHistory: []plsvc.RuntimeStatus{
			{UpdateTime: fakeCreatedAt, State: "PENDING"},
			{UpdateTime: "2024-06-01T10:01:00Z", State: "RUNNING"},
			{UpdateTime: fakeFinishedAt, State: "SUCCEEDED"},
		},
	}
}
