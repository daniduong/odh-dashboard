package fake

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	s3svc "github.com/opendatahub-io/odh-dashboard/packages/autox-core/services/s3"
)

// S3Client is a fake implementation of s3.Client for local development.
type S3Client struct{}

var _ s3svc.Client = (*S3Client)(nil)

func (c *S3Client) GetObject(_ context.Context, _ s3svc.ConnectionOptions, input s3svc.GetObjectInput) (io.ReadCloser, string, error) {
	body := fakeObjectBody(input.Key)
	contentType := "application/octet-stream"
	if strings.HasSuffix(input.Key, ".json") {
		contentType = "application/json"
	} else if strings.HasSuffix(input.Key, ".csv") {
		contentType = "text/csv"
	}
	return io.NopCloser(strings.NewReader(body)), contentType, nil
}

func (c *S3Client) DownloadObject(_ context.Context, _ s3svc.ConnectionOptions, input s3svc.DownloadObjectInput) (io.ReadCloser, string, error) {
	body := fakeObjectBody(input.Key)
	contentType := "application/octet-stream"
	if strings.HasSuffix(input.Key, ".json") {
		contentType = "application/json"
	}
	return io.NopCloser(strings.NewReader(body)), contentType, nil
}

func (c *S3Client) UploadObject(_ context.Context, _ s3svc.ConnectionOptions, input s3svc.UploadObjectInput) error {
	// Drain body to simulate upload.
	_, _ = io.Copy(io.Discard, input.Body)
	return nil
}

func (c *S3Client) ListObjects(_ context.Context, _ s3svc.ConnectionOptions, input s3svc.ListObjectsInput) (*s3svc.ListObjectsResponse, error) {
	objects, prefixes := fakeObjectListing(input.Prefix, input.Delimiter)
	return &s3svc.ListObjectsResponse{
		Name:           input.Bucket,
		Prefix:         input.Prefix,
		Delimiter:      input.Delimiter,
		IsTruncated:    false,
		KeyCount:       int32(len(objects) + len(prefixes)),
		MaxKeys:        input.Limit,
		Contents:       objects,
		CommonPrefixes: prefixes,
	}, nil
}

func (c *S3Client) ObjectExists(_ context.Context, _ s3svc.ConnectionOptions, input s3svc.ObjectExistsInput) (bool, error) {
	for _, key := range fakeObjectKeys() {
		if key == input.Key {
			return true, nil
		}
	}
	return false, nil
}

// fakeObjectKeys returns the full set of fake S3 object keys.
func fakeObjectKeys() []string {
	runID := fakeRunID
	artifactID := "5537ec33-96c2-4b8b-89b2-521e20c1232d"
	base := fmt.Sprintf("%s/%s/autogluon-models-training/%s/models_artifact/", fakeTabularPipelineName, runID, artifactID)
	return []string{
		base + "LightGBM_BAG_L2_FULL/model.json",
		base + "ExtraTreesEntr_BAG_L2_FULL/model.json",
		base + "WeightedEnsemble_L3_FULL/model.json",
		base + "leaderboard.csv",
	}
}

// fakeObjectListing returns objects and common prefixes for the given S3 prefix.
func fakeObjectListing(prefix, delimiter string) ([]s3svc.ObjectInfo, []s3svc.CommonPrefix) {
	allKeys := fakeObjectKeys()

	if delimiter == "" {
		// Flat listing — return all keys under prefix.
		var objects []s3svc.ObjectInfo
		for i, key := range allKeys {
			if strings.HasPrefix(key, prefix) {
				objects = append(objects, s3svc.ObjectInfo{
					Key:          key,
					Size:         int64(1024 * (i + 1)),
					ETag:         fmt.Sprintf("\"fake-etag-%02d\"", i),
					StorageClass: "STANDARD",
					LastModified: fakeFinishedAt,
				})
			}
		}
		return objects, nil
	}

	// Hierarchical listing — group by delimiter.
	seen := map[string]bool{}
	var objects []s3svc.ObjectInfo
	var prefixes []s3svc.CommonPrefix

	for i, key := range allKeys {
		if !strings.HasPrefix(key, prefix) {
			continue
		}
		rest := key[len(prefix):]
		idx := strings.Index(rest, delimiter)
		if idx >= 0 {
			// Virtual folder
			folder := prefix + rest[:idx+1]
			if !seen[folder] {
				seen[folder] = true
				prefixes = append(prefixes, s3svc.CommonPrefix{Prefix: folder})
			}
		} else {
			// Direct object
			objects = append(objects, s3svc.ObjectInfo{
				Key:          key,
				Size:         int64(1024 * (i + 1)),
				ETag:         fmt.Sprintf("\"fake-etag-%02d\"", i),
				StorageClass: "STANDARD",
				LastModified: fakeFinishedAt,
			})
		}
	}
	return objects, prefixes
}

// fakeObjectBody returns realistic fake content for the given key.
func fakeObjectBody(key string) string {
	if strings.HasSuffix(key, "leaderboard.csv") {
		return "model,score_val,score_test,pred_time_test,fit_time\n" +
			"WeightedEnsemble_L3_FULL,0.9312,0.9254,0.023,142.5\n" +
			"LightGBM_BAG_L2_FULL,0.9187,0.9141,0.018,38.2\n" +
			"ExtraTreesEntr_BAG_L2_FULL,0.9043,0.8998,0.021,25.7\n"
	}
	if strings.Contains(key, "WeightedEnsemble") {
		return marshalJSON(map[string]any{
			"model_type": "WeightedEnsemble",
			"weights":    []float64{0.6, 0.25, 0.15},
			"base_models": []string{
				"LightGBM_BAG_L2_FULL",
				"ExtraTreesEntr_BAG_L2_FULL",
				"RandomForest_BAG_L2_FULL",
			},
			"num_classes": 2,
		})
	}
	if strings.Contains(key, "LightGBM") {
		return marshalJSON(map[string]any{
			"model_type":   "LightGBM",
			"num_trees":    250,
			"learning_rate": 0.05,
			"num_leaves":   63,
			"feature_importance": map[string]float64{
				"feature_1": 0.32,
				"feature_2": 0.28,
				"feature_3": 0.19,
			},
		})
	}
	return marshalJSON(map[string]any{
		"model_type": "ExtraTrees",
		"n_estimators": 300,
		"max_features": "sqrt",
	})
}

func marshalJSON(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(b)
}
