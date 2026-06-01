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
	contentType := contentTypeForKey(input.Key)
	return io.NopCloser(strings.NewReader(body)), contentType, nil
}

func (c *S3Client) DownloadObject(_ context.Context, _ s3svc.ConnectionOptions, input s3svc.DownloadObjectInput) (io.ReadCloser, string, error) {
	body := fakeObjectBody(input.Key)
	return io.NopCloser(strings.NewReader(body)), contentTypeForKey(input.Key), nil
}

func (c *S3Client) UploadObject(_ context.Context, _ s3svc.ConnectionOptions, input s3svc.UploadObjectInput) error {
	_, _ = io.Copy(io.Discard, input.Body)
	return nil
}

func (c *S3Client) ListObjects(_ context.Context, _ s3svc.ConnectionOptions, input s3svc.ListObjectsInput) (*s3svc.ListObjectsResponse, error) {
	objects, prefixes := fakeObjectListing(input.Prefix, input.Delimiter)
	keyCount := int32(len(objects) + len(prefixes))
	maxKeys := input.Limit
	if maxKeys == 0 {
		maxKeys = 1000
	}
	return &s3svc.ListObjectsResponse{
		Name:           input.Bucket,
		Prefix:         input.Prefix,
		Delimiter:      input.Delimiter,
		IsTruncated:    false,
		KeyCount:       keyCount,
		MaxKeys:        maxKeys,
		Contents:       objects,
		CommonPrefixes: prefixes,
	}, nil
}

func (c *S3Client) ObjectExists(_ context.Context, _ s3svc.ConnectionOptions, input s3svc.ObjectExistsInput) (bool, error) {
	for _, key := range fakeLeafKeys() {
		if key == input.Key {
			return true, nil
		}
	}
	return false, nil
}

// fakeLeafKeys returns all fully-qualified S3 object keys in the fake dataset.
// These match the real path structure observed from the live API.
func fakeLeafKeys() []string {
	base := fmt.Sprintf("%s/%s/autogluon-models-training/%s/models_artifact/",
		fakeTabularPipelineName, fakeRunID, fakeArtifactID)
	models := []string{"ExtraTreesEntr_BAG_L2_FULL", "LightGBM_BAG_L2_FULL", "WeightedEnsemble_L3_FULL"}
	var keys []string
	for _, m := range models {
		keys = append(keys,
			base+m+"/model.json",
			base+m+"/metrics/metrics.json",
			base+m+"/notebooks/automl_predictor_notebook.ipynb",
		)
	}
	return keys
}

// fakeObjectListing returns objects and common prefixes that match the real app's
// hierarchical S3 listing behaviour. Each folder level returns common_prefixes only;
// leaf files appear only when the prefix matches their immediate parent folder.
func fakeObjectListing(prefix, delimiter string) ([]s3svc.ObjectInfo, []s3svc.CommonPrefix) {
	allKeys := fakeLeafKeys()

	if delimiter == "" {
		// Flat listing — return every key under the prefix.
		var objects []s3svc.ObjectInfo
		for i, key := range allKeys {
			if strings.HasPrefix(key, prefix) {
				objects = append(objects, objectInfo(key, i))
			}
		}
		return objects, nil
	}

	// Hierarchical listing — group by the first delimiter after the prefix.
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
			// Virtual folder — emit the common prefix once.
			folder := prefix + rest[:idx+1]
			if !seen[folder] {
				seen[folder] = true
				prefixes = append(prefixes, s3svc.CommonPrefix{Prefix: folder})
			}
		} else {
			// Direct object in this folder.
			objects = append(objects, objectInfo(key, i))
		}
	}
	return objects, prefixes
}

func objectInfo(key string, i int) s3svc.ObjectInfo {
	return s3svc.ObjectInfo{
		Key:          key,
		Size:         int64(4096 * (i + 1)),
		ETag:         fmt.Sprintf("\"fake-etag-%02d\"", i),
		StorageClass: "STANDARD",
		LastModified: fakeFinishedAt,
	}
}

func contentTypeForKey(key string) string {
	switch {
	case strings.HasSuffix(key, ".json"):
		return "application/json"
	case strings.HasSuffix(key, ".csv"):
		return "text/csv"
	case strings.HasSuffix(key, ".ipynb"):
		return "application/json"
	default:
		return "application/octet-stream"
	}
}

// fakeObjectBody returns realistic fake content for the given key.
func fakeObjectBody(key string) string {
	switch {
	case strings.Contains(key, "metrics.json"):
		return marshalJSON(map[string]any{
			"accuracy":                0.8312,
			"balanced_accuracy":       0.8187,
			"roc_auc":                 0.9043,
			"log_loss":                -0.4021,
			"root_mean_squared_error": -0.2891,
		})
	case strings.Contains(key, "automl_predictor_notebook.ipynb"):
		return marshalJSON(map[string]any{
			"nbformat":       4,
			"nbformat_minor": 5,
			"cells":          []any{},
			"metadata": map[string]any{
				"kernelspec": map[string]any{
					"display_name": "Python 3",
					"language":     "python",
					"name":         "python3",
				},
			},
		})
	case strings.Contains(key, "WeightedEnsemble"):
		return marshalJSON(map[string]any{
			"name":     "WeightedEnsemble_L3_FULL",
			"location": map[string]any{"model_directory": "WeightedEnsemble_L3_FULL"},
			"metrics": map[string]any{
				"test_data": map[string]any{"accuracy": 0.9312, "roc_auc": 0.9743},
			},
		})
	case strings.Contains(key, "LightGBM"):
		return marshalJSON(map[string]any{
			"name":     "LightGBM_BAG_L2_FULL",
			"location": map[string]any{"model_directory": "LightGBM_BAG_L2_FULL"},
			"metrics": map[string]any{
				"test_data": map[string]any{"accuracy": 0.9187, "roc_auc": 0.9621},
			},
		})
	case strings.Contains(key, "ExtraTrees"):
		return marshalJSON(map[string]any{
			"name":     "ExtraTreesEntr_BAG_L2_FULL",
			"location": map[string]any{"model_directory": "ExtraTreesEntr_BAG_L2_FULL"},
			"metrics": map[string]any{
				"test_data": map[string]any{"accuracy": 0.9043, "roc_auc": 0.9487},
			},
		})
	default:
		return "{}"
	}
}

func marshalJSON(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(b)
}
