package fake

import (
	"context"
	"fmt"

	kubernetes "github.com/opendatahub-io/odh-dashboard/packages/autox-core/services/kubernetes"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	fakeNamespace  = "fake-project"
	fakeDSPAName   = "dspa"
	fakeS3Secret   = "ds-pipeline-s3-dspa"
	fakeS3Bucket   = "fake-ml-bucket"
	fakeS3Endpoint = "http://fake-s3.fake-project.svc.cluster.local:9000"
	fakeS3Region   = "us-east-1"
)

// K8sClient is a fake implementation of kubernetes.Client for local development.
type K8sClient struct{}

var _ kubernetes.Client = (*K8sClient)(nil)

func (c *K8sClient) GetNamespaces(_ context.Context) ([]v1.Namespace, error) {
	return []v1.Namespace{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: fakeNamespace,
				Annotations: map[string]string{
					"openshift.io/display-name": "Fake Project",
				},
			},
		},
	}, nil
}

func (c *K8sClient) GetPods(_ context.Context, _ string) (*v1.PodList, error) {
	return &v1.PodList{}, nil
}

func (c *K8sClient) GetSecrets(_ context.Context, _ string) ([]v1.Secret, error) {
	return []v1.Secret{fakeS3CredentialsSecret(fakeNamespace, fakeS3Secret)}, nil
}

func (c *K8sClient) GetSecret(_ context.Context, namespace, secretName string) (*v1.Secret, error) {
	s := fakeS3CredentialsSecret(namespace, secretName)
	return &s, nil
}

func (c *K8sClient) GetUser(_ context.Context) (string, error) {
	return "fake-user@example.com", nil
}

func (c *K8sClient) IsClusterAdmin(_ context.Context) (bool, error) {
	return true, nil
}

func (c *K8sClient) CanAccessResource(_ context.Context, _, _, _, _, _ string) (bool, error) {
	return true, nil
}

// ListResources returns a ready DSPipelineApplication when queried for the DSPA GVR,
// and an empty list for everything else.
func (c *K8sClient) ListResources(_ context.Context, gvr schema.GroupVersionResource, namespace string) (*unstructured.UnstructuredList, error) {
	if gvr.Resource == "datasciencepipelinesapplications" {
		return &unstructured.UnstructuredList{
			Items: []unstructured.Unstructured{fakeDSPA(namespace)},
		}, nil
	}
	return &unstructured.UnstructuredList{}, nil
}

func (c *K8sClient) GetResource(_ context.Context, gvr schema.GroupVersionResource, namespace, name string) (*unstructured.Unstructured, error) {
	if gvr.Resource == "datasciencepipelinesapplications" {
		dspa := fakeDSPA(namespace)
		dspa.SetName(name)
		return &dspa, nil
	}
	return &unstructured.Unstructured{}, nil
}

func (c *K8sClient) CreateResource(_ context.Context, _ schema.GroupVersionResource, _ string, obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	return obj, nil
}

func (c *K8sClient) DiscoverResourceGVR(_ context.Context, group, resource, _ string, knownVersions []string) (schema.GroupVersionResource, error) {
	version := "v1"
	if len(knownVersions) > 0 {
		version = knownVersions[0]
	}
	return schema.GroupVersionResource{Group: group, Version: version, Resource: resource}, nil
}

// fakeDSPA builds a ready DSPipelineApplication unstructured object for the given namespace.
func fakeDSPA(namespace string) unstructured.Unstructured {
	pipelineURL := fmt.Sprintf("http://ds-pipeline-api-%s.%s.svc.cluster.local:8888", fakeDSPAName, namespace)
	dspa := unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "datasciencepipelinesapplications.opendatahub.io/v1",
			"kind":       "DSPipelineApplication",
			"metadata": map[string]any{
				"name":      fakeDSPAName,
				"namespace": namespace,
			},
			"spec": map[string]any{
				"objectStorage": map[string]any{
					"externalStorage": map[string]any{
						"s3CredentialSecret": map[string]any{
							"secretName": fakeS3Secret,
							"accessKey":  "AWS_ACCESS_KEY_ID",
							"secretKey":  "AWS_SECRET_ACCESS_KEY",
						},
						"scheme": "http",
						"host":   "fake-s3.fake-project.svc.cluster.local",
						"port":   "9000",
						"bucket": fakeS3Bucket,
						"region": fakeS3Region,
					},
				},
			},
			"status": map[string]any{
				"conditions": []any{
					map[string]any{"type": "Ready", "status": "True"},
				},
				"components": map[string]any{
					"apiServer": map[string]any{
						"url": pipelineURL,
					},
				},
			},
		},
	}
	return dspa
}

// fakeS3CredentialsSecret returns a v1.Secret pre-populated with fake AWS S3 credentials.
func fakeS3CredentialsSecret(namespace, name string) v1.Secret {
	return v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"AWS_ACCESS_KEY_ID":     []byte("FAKEACCESSKEYID12345"),
			"AWS_SECRET_ACCESS_KEY": []byte("FakeSecretAccessKey/ABCDEFGHIJKLMNOP"),
			"AWS_DEFAULT_REGION":    []byte(fakeS3Region),
			"AWS_S3_ENDPOINT":       []byte(fakeS3Endpoint),
			"AWS_S3_BUCKET":         []byte(fakeS3Bucket),
		},
	}
}
