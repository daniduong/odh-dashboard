package fake

import (
	"context"
	"encoding/json"
	"strings"

	plsvc "github.com/opendatahub-io/odh-dashboard/packages/autox-core/services/pipelines"
)

const (
	fakeTabularPipelineID   = "33dc7341-9341-4a9a-85e2-ba786f2ebce6"
	fakeTabularVersionID    = "96f2b632-faa4-4526-b97b-1a4c2c6ad753"
	fakeTabularPipelineName = "autogluon-tabular-training-pipeline"

	fakeTimeSeriesPipelineID   = "cccccccc-0001-0001-0001-cccccccccccc"
	fakeTimeSeriesVersionID    = "dddddddd-0001-0001-0001-dddddddddddd"
	fakeTimeSeriesPipelineName = "autogluon-timeseries-training-pipeline"

	fakeExperimentID = "e3e5bccd-bc4d-4b70-a735-7645c6258950"
	fakeRunID        = "c377d1a1-15fb-48a3-88c6-8831fa88f19d"
	fakeArtifactID   = "5537ec33-96c2-4b8b-89b2-521e20c1232d"
	fakeCreatedAt    = "2026-05-28T19:02:14Z"
	fakeScheduledAt  = "2026-05-28T19:02:14Z"
	fakeFinishedAt   = "2026-05-28T19:05:07Z"
)

// PipelinesClient is a fake implementation of pipelines.Client for local development.
type PipelinesClient struct{}

var _ plsvc.Client = (*PipelinesClient)(nil)

func (c *PipelinesClient) ListPipelines(_ context.Context, _ string, _ string) (*plsvc.PipelinesResponse, error) {
	pipelines := []plsvc.Pipeline{
		{
			PipelineID:  fakeTabularPipelineID,
			DisplayName: fakeTabularPipelineName,
			Description: "End-to-end AutoGluon tabular training pipeline (classification and regression)",
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
		RunID:          fakeRunID,
		DisplayName:    input.DisplayName,
		Description:    input.Description,
		ExperimentID:   fakeExperimentID,
		State:          "RUNNING",
		StorageState:   "AVAILABLE",
		ServiceAccount: "pipeline-runner-dspa",
		CreatedAt:      fakeCreatedAt,
		ScheduledAt:    fakeScheduledAt,
		PipelineVersionReference: &plsvc.PipelineVersionReference{
			PipelineID:        pipelineID,
			PipelineVersionID: versionID,
		},
		RuntimeConfig: input.RuntimeConfig,
		StateHistory: []plsvc.RuntimeStatus{
			{UpdateTime: fakeCreatedAt, State: "PENDING"},
			{UpdateTime: "2026-05-28T19:02:15Z", State: "RUNNING"},
		},
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
func (c *PipelinesClient) RetryRun(_ context.Context, _ string, _ string) error     { return nil }
func (c *PipelinesClient) DeleteRun(_ context.Context, _ string, _ string) error    { return nil }

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

// fakeCompletedRun returns a realistic SUCCEEDED tabular pipeline run.
func fakeCompletedRun(runID string) *plsvc.PipelineRun {
	pipelineSpec, _ := json.Marshal(map[string]any{
		"pipeline_spec": map[string]any{
			"pipelineInfo": map[string]any{
				"name":        fakeTabularPipelineName,
				"description": "End-to-end AutoGluon tabular training pipeline",
			},
			"schemaVersion": "2.1.0",
			"sdkVersion":    "kfp-2.16.1",
		},
	})

	return &plsvc.PipelineRun{
		RunID:          runID,
		DisplayName:    "autox-core - 1",
		ExperimentID:   fakeExperimentID,
		State:          "SUCCEEDED",
		StorageState:   "AVAILABLE",
		ServiceAccount: "pipeline-runner-dspa",
		CreatedAt:      fakeCreatedAt,
		ScheduledAt:    fakeScheduledAt,
		FinishedAt:     fakeFinishedAt,
		PipelineVersionReference: &plsvc.PipelineVersionReference{
			PipelineID:        fakeTabularPipelineID,
			PipelineVersionID: fakeTabularVersionID,
		},
		RuntimeConfig: &plsvc.RuntimeConfig{
			Parameters: map[string]any{
				"label_column":           "Survived",
				"task_type":              "binary",
				"top_n":                  3,
				"train_data_bucket_name": "fake-ml-bucket",
				"train_data_file_key":    "datasets/TitanicFullMF.csv",
				"train_data_secret_name": "fake-s3-credentials",
			},
		},
		PipelineSpec: json.RawMessage(pipelineSpec),
		StateHistory: []plsvc.RuntimeStatus{
			{UpdateTime: fakeCreatedAt, State: "PENDING"},
			{UpdateTime: "2026-05-28T19:02:15Z", State: "RUNNING"},
			{UpdateTime: fakeFinishedAt, State: "SUCCEEDED"},
		},
		RunDetails: &plsvc.RunDetails{
			TaskDetails: []plsvc.TaskDetail{
				fakeTask(runID, "ad183ecb-8a0f-4506-aae4-50b8d799e1e7", "autogluon-tabular-training-pipeline-k8dxc", fakeCreatedAt, fakeCreatedAt, fakeFinishedAt, "SUCCEEDED"),
				fakeTask(runID, "cb08a497-8c79-46b3-b8cb-0bc39c42b93d", "root-driver", fakeCreatedAt, fakeCreatedAt, "2026-05-28T19:02:18Z", "SUCCEEDED"),
				fakeTask(runID, "6ed2c277-f52f-4019-b991-00b9fbedfc17", "root", fakeCreatedAt, "2026-05-28T19:02:24Z", fakeFinishedAt, "SUCCEEDED"),
				fakeTask(runID, "f3c68dd8-f0b0-4118-9d1c-83cb04da1b62", "automl-data-loader-driver", fakeCreatedAt, "2026-05-28T19:02:24Z", "2026-05-28T19:02:28Z", "SUCCEEDED"),
				fakeTask(runID, "04d6a28e-9449-45bc-8994-4e7ff5ecec5c", "automl-data-loader", fakeCreatedAt, "2026-05-28T19:02:34Z", "2026-05-28T19:02:55Z", "SUCCEEDED"),
				fakeTask(runID, "4782bc60-065c-4040-8c44-fed145003c16", "autogluon-models-training-driver", fakeCreatedAt, "2026-05-28T19:02:55Z", "2026-05-28T19:02:59Z", "SUCCEEDED"),
				fakeTask(runID, "03703b7f-689e-4334-b86b-0166c9e846b5", "autogluon-models-training", fakeCreatedAt, "2026-05-28T19:03:05Z", "2026-05-28T19:04:36Z", "SUCCEEDED"),
				fakeTask(runID, "f96c9daf-1527-4e84-a7ef-637bb20b338e", "leaderboard-evaluation-driver", fakeCreatedAt, "2026-05-28T19:04:36Z", "2026-05-28T19:04:40Z", "SUCCEEDED"),
				fakeTask(runID, "b6fbd9b9-1f72-4cd3-b083-3d2fc2a6454a", "leaderboard-evaluation", fakeCreatedAt, "2026-05-28T19:04:46Z", fakeFinishedAt, "SUCCEEDED"),
			},
		},
	}
}

func fakeTask(runID, taskID, displayName, createTime, startTime, endTime, state string) plsvc.TaskDetail {
	return plsvc.TaskDetail{
		RunID:       runID,
		TaskID:      taskID,
		DisplayName: displayName,
		CreateTime:  createTime,
		StartTime:   startTime,
		EndTime:     endTime,
		State:       state,
		StateHistory: []plsvc.RuntimeStatus{
			{UpdateTime: startTime, State: "RUNNING"},
			{UpdateTime: endTime, State: state},
		},
	}
}
