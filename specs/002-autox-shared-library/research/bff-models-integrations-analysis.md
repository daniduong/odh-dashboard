# BFF Models & Integrations Duplication Analysis

## Summary
- **Total files analyzed**: 44 files (22 pairs)
- **Average duplication**: ~88%
- **Extraction candidates**: 18 files with 90%+ duplication

## Key Findings

### Models Directory (11 files)
- **8 files are 100% identical**: groups.go, health_check.go, namespace.go, rbac_types.go, s3.go, secret.go, user.go, pipeline_server.go
- **1 file with 99% duplication**: pipeline_runs.go (only differs in product-specific request structs)
- **2 files are product-specific**: model_registry.go (AutoML), lsd_models.go + lsd_vector_stores.go (AutoRAG)

### Integrations Directory (7 integration types)
- **100% identical**: http.go, kubernetes/types.go
- **98% identical**: s3/client.go (AutoML has GetCSVSchema method, AutoRAG does not)
- **95% identical**: pipelineserver/client.go (AutoML has ListPipelines and ListPipelineVersions methods)
- **~95% identical**: kubernetes/* clients (minor differences in error handling)

---

## File-by-File Comparison

### Model Pair 1: groups.go

| Metric | AutoML | AutoRAG |
|--------|--------|---------|
| Total Lines | 51 | 51 |
| Duplicate Lines | 51 | 51 |
| Duplication % | 100% | 100% |
| Shared Structs | K8sObjectMeta, Group, GroupModel | Same |
| Shared Interfaces | N/A | N/A |
| Shared Methods | NewGroup(), NewGroupModel() | Same |

**Extractable Logic:**
- Entire file can be extracted to shared library
- K8sObjectMeta struct (core Kubernetes metadata wrapper)
- Group struct (OpenShift/Kubernetes Group object)
- GroupModel struct (legacy backward compatibility)
- NewGroup() constructor
- NewGroupModel() constructor

---

### Model Pair 2: health_check.go

| Metric | AutoML | AutoRAG |
|--------|--------|---------|
| Total Lines | 11 | 11 |
| Duplicate Lines | 11 | 11 |
| Duplication % | 100% | 100% |
| Shared Structs | SystemInfo, HealthCheckModel | Same |
| Shared Interfaces | N/A | N/A |
| Shared Methods | N/A | N/A |

**Extractable Logic:**
- Entire file can be extracted to shared library
- SystemInfo struct (version information)
- HealthCheckModel struct (health check response)

---

### Model Pair 3: namespace.go

| Metric | AutoML | AutoRAG |
|--------|--------|---------|
| Total Lines | 22 | 22 |
| Duplicate Lines | 22 | 22 |
| Duplication % | 100% | 100% |
| Shared Structs | NamespaceModel | Same |
| Shared Interfaces | N/A | N/A |
| Shared Methods | NewNamespaceModelFromNamespace() | Same |

**Extractable Logic:**
- Entire file can be extracted to shared library
- NamespaceModel struct (namespace representation with display name)
- NewNamespaceModelFromNamespace() factory (extracts OpenShift display name annotation)

---

### Model Pair 4: pipeline_runs.go

| Metric | AutoML | AutoRAG |
|--------|--------|---------|
| Total Lines | 192 | 176 |
| Duplicate Lines | 176 | 176 |
| Duplication % | 99% | 100% |
| Shared Structs | PipelineRun, PipelineVersionReference, PipelineRunsData, KFPipelineRunResponse, KFPipelineRun, RuntimeConfig, RuntimeStatus, ErrorInfo, RunDetails, TaskDetail, ChildTask, CreatePipelineRunKFRequest, KFPipeline, KFPipelineVersion, KFPipelinesResponse, KFPipelineVersionsResponse | Same (except create request) |
| Shared Interfaces | N/A | N/A |
| Shared Methods | N/A | N/A |

**Differences:**
- AutoML: CreateAutoMLRunRequest (lines 111-144) - tabular/timeseries specific fields
- AutoRAG: CreateAutoRAGRunRequest (lines 112-128) - LlamaStack/RAG specific fields
- AutoRAG: TaskDetail has ExecutionID field (line 101) that AutoML does not

**Extractable Logic:**
- All structs except CreateAutoMLRunRequest and CreateAutoRAGRunRequest
- PipelineRun (public API format)
- PipelineVersionReference
- PipelineRunsData (with pagination)
- KFPipelineRunResponse (Kubeflow internal format)
- KFPipelineRun (single run from KFP)
- RuntimeConfig, RuntimeStatus, ErrorInfo
- RunDetails, TaskDetail (with optional ExecutionID), ChildTask
- CreatePipelineRunKFRequest (generic KFP create request)
- KFPipeline, KFPipelineVersion (pipeline discovery)
- KFPipelinesResponse, KFPipelineVersionsResponse (pagination)

**Product-Specific Retain:**
- CreateAutoMLRunRequest in AutoML
- CreateAutoRAGRunRequest in AutoRAG
- ExecutionID field in TaskDetail (can be made optional in shared version)

---

### Model Pair 5: pipeline_server.go

| Metric | AutoML | AutoRAG |
|--------|--------|---------|
| Total Lines | 94 | 94 |
| Duplicate Lines | 94 | 94 |
| Duplication % | 100% | 100% |
| Shared Structs | DSPipelineApplication, DSPipelineApplicationMetadata, DSPipelineApplicationSpec, APIServer, ObjectStorage, ExternalStorage, S3CredentialsSecret, MinioStorage, DSPAObjectStorage, DSPipelineApplicationStatus, DSPipelineApplicationCondition, DSPipelineApplicationComponents, DSPipelineApplicationAPIServerStatus | Same |
| Shared Interfaces | N/A | N/A |
| Shared Methods | N/A | N/A |

**Extractable Logic:**
- Entire file can be extracted to shared library
- DSPipelineApplication CR (Kubeflow DSPipelineApplication custom resource)
- DSPipelineApplicationSpec (pipeline configuration)
- ObjectStorage (S3-compatible storage configuration)
- ExternalStorage, MinioStorage (storage backends)
- S3CredentialsSecret (secret reference)
- DSPAObjectStorage (resolved object-storage configuration for request context)
- DSPipelineApplicationStatus (pipeline readiness status)
- DSPipelineApplicationAPIServerStatus (API server URL)

---

### Model Pair 6: rbac_types.go

| Metric | AutoML | AutoRAG |
|--------|--------|---------|
| Total Lines | 28 | 28 |
| Duplicate Lines | 28 | 28 |
| Duplication % | 100% | 100% |
| Shared Structs | CertificateItem, CertificateList, RoleBinding (alias), RoleBindingList | Same |
| Shared Interfaces | N/A | N/A |
| Shared Methods | N/A | N/A |

**Extractable Logic:**
- Entire file can be extracted to shared library
- CertificateItem (secret/configmap with keys for certificates)
- CertificateList (lists of secrets and configmaps)
- RoleBinding (type alias for k8s.io/api/rbac/v1.RoleBinding)
- RoleBindingList (list of role bindings)

---

### Model Pair 7: s3.go

| Metric | AutoML | AutoRAG |
|--------|--------|---------|
| Total Lines | 31 | 31 |
| Duplicate Lines | 31 | 31 |
| Duplication % | 100% | 100% |
| Shared Structs | S3ObjectInfo, S3CommonPrefix, S3ListObjectsResponse | Same |
| Shared Interfaces | N/A | N/A |
| Shared Methods | N/A | N/A |

**Extractable Logic:**
- Entire file can be extracted to shared library
- S3ObjectInfo (single S3 object metadata)
- S3CommonPrefix (virtual folder prefix)
- S3ListObjectsResponse (BFF's decoupled S3 list response)

---

### Model Pair 8: secret.go

| Metric | AutoML | AutoRAG |
|--------|--------|---------|
| Total Lines | 24 | 24 |
| Duplicate Lines | 24 | 24 |
| Duplication % | 100% | 100% |
| Shared Structs | SecretListItem | Same |
| Shared Interfaces | N/A | N/A |
| Shared Methods | NewSecretListItem() | Same |

**Extractable Logic:**
- Entire file can be extracted to shared library
- SecretListItem (filtered secret with UUID, name, type, data)
- NewSecretListItem() constructor

---

### Model Pair 9: user.go

| Metric | AutoML | AutoRAG |
|--------|--------|---------|
| Total Lines | 7 | 7 |
| Duplicate Lines | 7 | 7 |
| Duplication % | 100% | 100% |
| Shared Structs | User | Same |
| Shared Interfaces | N/A | N/A |
| Shared Methods | N/A | N/A |

**Extractable Logic:**
- Entire file can be extracted to shared library
- User struct (user ID and cluster admin flag)

---

### Model Pair 10: model_registry.go (AutoML only)

| Metric | AutoML | AutoRAG |
|--------|--------|---------|
| Total Lines | 27 | N/A |
| Duplicate Lines | 0 | N/A |
| Duplication % | N/A | N/A |
| Shared Structs | ModelRegistrySpec | N/A |
| Shared Interfaces | N/A | N/A |
| Shared Methods | N/A | N/A |

**Extractable Logic:**
- **Product-specific** - AutoML only
- ModelRegistrySpec (model registry configuration)
- Not extractable to shared library

---

### Model Pair 11: register_model.go (AutoML only)

| Metric | AutoML | AutoRAG |
|--------|--------|---------|
| Total Lines | 51 | N/A |
| Duplicate Lines | 0 | N/A |
| Duplication % | N/A | N/A |
| Shared Structs | RegisterModelRequest, ModelDetails, Author, RegisteredModelResponse, ModelVersionResponse, ArtifactResponse, RegisterModelConflictError | N/A |
| Shared Interfaces | N/A | N/A |
| Shared Methods | N/A | N/A |

**Extractable Logic:**
- **Product-specific** - AutoML only
- Model registry integration types
- Not extractable to shared library

---

### Model Pair 12: lsd_models.go (AutoRAG only)

| Metric | AutoML | AutoRAG |
|--------|--------|---------|
| Total Lines | N/A | 24 |
| Duplicate Lines | N/A | 0 |
| Duplication % | N/A | N/A |
| Shared Structs | N/A | LSDModel, LSDModelsResponse |
| Shared Interfaces | N/A | N/A |
| Shared Methods | N/A | N/A |

**Extractable Logic:**
- **Product-specific** - AutoRAG only
- LlamaStack models integration types
- Not extractable to shared library

---

### Model Pair 13: lsd_vector_stores.go (AutoRAG only)

| Metric | AutoML | AutoRAG |
|--------|--------|---------|
| Total Lines | N/A | 21 |
| Duplicate Lines | N/A | 0 |
| Duplication % | N/A | N/A |
| Shared Structs | N/A | LSDVectorStore, LSDVectorStoresResponse |
| Shared Interfaces | N/A | N/A |
| Shared Methods | N/A | N/A |

**Extractable Logic:**
- **Product-specific** - AutoRAG only
- LlamaStack vector stores integration types
- Not extractable to shared library

---

## Integration Files

### Integration Pair 1: http.go

| Metric | AutoML | AutoRAG |
|--------|--------|---------|
| Total Lines | 18 | 18 |
| Duplicate Lines | 18 | 18 |
| Duplication % | 100% | 100% |
| Shared Structs | ErrorResponse, HTTPError | Same |
| Shared Interfaces | N/A | N/A |
| Shared Methods | Error() string | Same |

**Extractable Logic:**
- Entire file can be extracted to shared library
- ErrorResponse (code and message)
- HTTPError (with status code and error interface implementation)
- Error() method for HTTPError

---

### Integration Pair 2: kubernetes/types.go

| Metric | AutoML | AutoRAG |
|--------|--------|---------|
| Total Lines | 34 | 34 |
| Duplicate Lines | 34 | 34 |
| Duplication % | 100% | 100% |
| Shared Structs | ServiceDetails, RequestIdentity, BearerToken | Same |
| Shared Interfaces | N/A | N/A |
| Shared Methods | NewBearerToken(), String(), Raw() | Same |

**Extractable Logic:**
- Entire file can be extracted to shared library
- ServiceDetails (Kubernetes service metadata)
- RequestIdentity (user ID, groups, token)
- BearerToken (secure token wrapper with redaction)
- NewBearerToken() constructor
- String() redaction method
- Raw() accessor method

---

### Integration Pair 3: s3/client.go

| Metric | AutoML | AutoRAG |
|--------|--------|---------|
| Total Lines | 866 | 450 |
| Duplicate Lines | 450 | 450 |
| Duplication % | 52% | 100% |
| Shared Structs | S3Credentials, ListObjectsOptions, S3ClientInterface (interface), RealS3Client | Same |
| Shared Interfaces | S3ClientInterface | Same |
| Shared Methods | buildS3AWSConfig(), NewRealS3Client(), GetObject(), UploadObject(), ListObjects(), ObjectExists(), cloneDefaultTransport(), isInternalHost(), validateAndNormalizeEndpoint(), validateIPAddress(), isS3ConditionalCreateConflict() | Same |

**Differences:**
- **AutoML-specific** (416 lines):
  - ColumnSchema, CSVSchemaResult structs
  - GetCSVSchema() method (lines 309-472)
  - CSV helper functions: normalizeLineEndings(), inferColumnType(), allValuesMatchType(), isNumber(), isInteger(), isTimestamp(), isBoolean(), collectBooleanValues(), extractFirstLine(), countLines()

**Extractable Logic:**
- Core S3 client (450 lines):
  - S3Credentials struct
  - ListObjectsOptions struct
  - S3ClientInterface interface
  - RealS3Client struct
  - buildS3AWSConfig() - AWS config builder
  - NewRealS3Client() - client factory with endpoint validation
  - GetObject() - download with transfer manager
  - UploadObject() - upload with conditional create
  - ListObjects() - list with pagination
  - ObjectExists() - HEAD request existence check
  - cloneDefaultTransport() - HTTP transport cloning
  - isInternalHost() - in-cluster service detection
  - validateAndNormalizeEndpoint() - SSRF protection
  - validateIPAddress() - IP range validation
  - isS3ConditionalCreateConflict() - error code checking

**Product-Specific Retain:**
- AutoML: CSV schema inference (GetCSVSchema and all CSV helpers)
- AutoRAG: None (uses base client only)

---

### Integration Pair 4: pipelineserver/client.go

| Metric | AutoML | AutoRAG |
|--------|--------|---------|
| Total Lines | 631 | 631 |
| Duplicate Lines | 487 | 487 |
| Duplication % | 77% | 77% |
| Shared Structs | HTTPError, PipelineServerClient (interface), RealPipelineServerClient | Same |
| Shared Interfaces | PipelineServerClient | Same (except AutoML adds ListPipelines and ListPipelineVersions) |
| Shared Methods | NewRealPipelineServerClient(), CreateRun(), GetRun(), GetRunDetails(), GetPipelineVersion() | Same |

**Differences:**
- **AutoML-specific** (144 lines):
  - ListPipelines() method (lines 238-321) - retrieves all pipelines with pagination
  - ListPipelineVersions() method (lines 323-406) - retrieves all versions for a pipeline with pagination
  - Both methods handle pagination and aggregate results

**Extractable Logic:**
- Core pipeline server client (487 lines):
  - HTTPError struct
  - PipelineServerClient interface (base methods)
  - RealPipelineServerClient struct
  - NewRealPipelineServerClient() factory
  - CreateRun() - create pipeline run
  - GetRun() - get single run
  - GetRunDetails() - get run with task details
  - GetPipelineVersion() - get pipeline version with spec
  - HTTP helpers and error handling

**Product-Specific Retain:**
- AutoML: ListPipelines() and ListPipelineVersions() methods
  - These are used for AutoML's pipeline discovery mechanism
  - Can be extracted if both products need pipeline discovery
- AutoRAG: Uses base client methods only

---

### Integration Pair 5: kubernetes/client.go

| Metric | AutoML | AutoRAG |
|--------|--------|---------|
| Total Lines | 258 | 259 |
| Duplicate Lines | ~245 | ~245 |
| Duplication % | 95% | 95% |
| Shared Structs | N/A (mostly interfaces) | Same |
| Shared Interfaces | K8sClient | Same |
| Shared Methods | GetNamespaces(), GetSecret(), GetDSPipelineApplication(), GetServiceDetails(), GetGroups(), GetCertificates(), GetRoleBindings(), CreateRoleBinding() | Same |

**Differences:**
- AutoRAG has errors.go file (13 lines) with ErrNoLlamaStackDistribution
- Minor comment differences

**Extractable Logic:**
- Entire K8sClient interface
- All methods for:
  - Namespace operations (GetNamespaces)
  - Secret operations (GetSecret)
  - DSPipelineApplication operations (GetDSPipelineApplication)
  - Service operations (GetServiceDetails)
  - Group operations (GetGroups)
  - Certificate operations (GetCertificates)
  - RBAC operations (GetRoleBindings, CreateRoleBinding)

---

### Integration Pair 6: kubernetes/factory.go

| Metric | AutoML | AutoRAG |
|--------|--------|---------|
| Total Lines | 41 | 41 |
| Duplicate Lines | 41 | 41 |
| Duplication % | 100% | 100% |
| Shared Structs | N/A | Same |
| Shared Interfaces | K8sClientFactory | Same |
| Shared Methods | NewK8sClient() | Same |

**Extractable Logic:**
- Entire file can be extracted to shared library
- K8sClientFactory interface
- NewK8sClient() factory method

---

### Integration Pair 7: kubernetes/internal_k8s_client.go

| Metric | AutoML | AutoRAG |
|--------|--------|---------|
| Total Lines | 442 | 446 |
| Duplicate Lines | ~438 | ~438 |
| Duplication % | 99% | 98% |
| Shared Structs | InternalK8sClient | Same |
| Shared Interfaces | N/A | N/A |
| Shared Methods | NewInternalK8sClient(), GetNamespaces(), GetSecret(), GetDSPipelineApplication(), GetServiceDetails(), GetGroups(), GetCertificates(), GetRoleBindings(), CreateRoleBinding() | Same |

**Differences:**
- Minor differences in error handling (4-8 lines)

**Extractable Logic:**
- Entire InternalK8sClient implementation
- All method implementations (namespace, secret, DSPA, service, group, certificate, RBAC operations)
- Kubernetes client initialization
- Error handling patterns

---

### Integration Pair 8: kubernetes/shared_k8s_client.go

| Metric | AutoML | AutoRAG |
|--------|--------|---------|
| Total Lines | 135 | 135 |
| Duplicate Lines | 135 | 135 |
| Duplication % | 100% | 100% |
| Shared Structs | N/A (shared helper functions) | Same |
| Shared Interfaces | N/A | N/A |
| Shared Methods | getNamespaces(), getDSPipelineApplication(), getServiceDetails(), getGroups(), getCertificates(), getRoleBindings(), createRoleBinding(), getSecret() | Same |

**Extractable Logic:**
- Entire file can be extracted to shared library
- All shared helper functions for Kubernetes operations
- Common patterns for namespace, secret, DSPA, service, group, certificate, RBAC operations

---

### Integration Pair 9: kubernetes/token_k8s_client.go

| Metric | AutoML | AutoRAG |
|--------|--------|---------|
| Total Lines | 382 | 383 |
| Duplicate Lines | ~378 | ~378 |
| Duplication % | 99% | 99% |
| Shared Structs | TokenK8sClient | Same |
| Shared Interfaces | N/A | N/A |
| Shared Methods | NewTokenK8sClient(), GetNamespaces(), GetSecret(), GetDSPipelineApplication(), GetServiceDetails(), GetGroups(), GetCertificates(), GetRoleBindings(), CreateRoleBinding() | Same |

**Differences:**
- Minor differences in error handling (4-5 lines)

**Extractable Logic:**
- Entire TokenK8sClient implementation
- All method implementations with user token-based authentication
- Kubernetes client initialization with bearer token
- Error handling patterns

---

### Integration Pair 10: kubernetes/portforward.go

| Metric | AutoML | AutoRAG |
|--------|--------|---------|
| Total Lines | ~100 | ~100 |
| Duplicate Lines | ~100 | ~100 |
| Duplication % | 100% | 100% |
| Shared Structs | PortForwarder | Same |
| Shared Interfaces | N/A | N/A |
| Shared Methods | NewPortForwarder(), Start(), Stop(), GetLocalPort() | Same |

**Extractable Logic:**
- Entire file can be extracted to shared library
- PortForwarder implementation
- Port forwarding lifecycle management
- Local port allocation

---

## Product-Specific Integrations (Not Shared)

### AutoML-Specific

#### integrations/modelregistry/httpclient.go (326 lines)
- ModelRegistryClient interface
- RealModelRegistryClient implementation
- Model registration operations
- Not extractable to shared library

#### integrations/modelregistry/mock_client.go (127 lines)
- Mock implementation for testing
- Not extractable to shared library

### AutoRAG-Specific

#### integrations/llamastack/llamastack_client.go (198 lines)
- LlamaStackClient interface
- RealLlamaStackClient implementation
- LlamaStack API operations (models, vector stores)
- Not extractable to shared library

#### integrations/llamastack/llamastack_client_factory.go (29 lines)
- LlamaStackClientFactory interface
- Factory implementation
- Not extractable to shared library

#### integrations/llamastack/errors.go (13 lines)
- LlamaStack-specific error types
- Not extractable to shared library

---

## Extraction Summary

### High-Priority Extractions (100% Duplication)

**Models (8 files - 268 lines total):**
1. groups.go (51 lines) - Kubernetes group types
2. health_check.go (11 lines) - Health check types
3. namespace.go (22 lines) - Namespace types
4. pipeline_server.go (94 lines) - DSPipelineApplication types
5. rbac_types.go (28 lines) - RBAC types
6. s3.go (31 lines) - S3 response types
7. secret.go (24 lines) - Secret list types
8. user.go (7 lines) - User types

**Integrations (3 files - 52 lines total):**
1. http.go (18 lines) - HTTP error types
2. kubernetes/types.go (34 lines) - Kubernetes service/identity types

**Kubernetes Integration (5 files - ~1155 lines total):**
1. kubernetes/factory.go (41 lines) - Client factory
2. kubernetes/shared_k8s_client.go (135 lines) - Shared helpers
3. kubernetes/internal_k8s_client.go (442 lines) - Internal client
4. kubernetes/token_k8s_client.go (382 lines) - Token-based client
5. kubernetes/portforward.go (~100 lines) - Port forwarding
6. kubernetes/client.go (258 lines) - Client interface (95% shared)

### Medium-Priority Extractions (95-99% Duplication)

**Models (1 file - 176 lines base):**
1. pipeline_runs.go - Extract base structs, keep product-specific request types separate

**Integrations:**
1. s3/client.go - Extract base client (450 lines), keep AutoML CSV schema separate
2. pipelineserver/client.go - Extract base client (487 lines), keep AutoML pipeline discovery separate

### Total Extractable Code

- **Models**: ~444 lines (8 files at 100% + 176 lines from pipeline_runs.go)
- **Integrations**: ~2144 lines
  - HTTP/Types: 52 lines (100%)
  - Kubernetes: ~1155 lines (95-100%)
  - S3: 450 lines base (52% of AutoML, 100% of AutoRAG)
  - Pipeline Server: 487 lines base (77% of both)

**Total**: ~2588 lines of highly duplicated code ready for extraction

---

## Recommended Shared Library Structure

```
packages/autox-shared/bff/
├── internal/
│   ├── models/
│   │   ├── groups.go              # 100% shared
│   │   ├── health_check.go        # 100% shared
│   │   ├── namespace.go           # 100% shared
│   │   ├── pipeline_runs.go       # 99% shared (base structs)
│   │   ├── pipeline_server.go     # 100% shared
│   │   ├── rbac_types.go          # 100% shared
│   │   ├── s3.go                  # 100% shared
│   │   ├── secret.go              # 100% shared
│   │   └── user.go                # 100% shared
│   └── integrations/
│       ├── http.go                # 100% shared
│       ├── kubernetes/
│       │   ├── types.go           # 100% shared
│       │   ├── client.go          # 95% shared
│       │   ├── factory.go         # 100% shared
│       │   ├── internal_k8s_client.go  # 99% shared
│       │   ├── token_k8s_client.go     # 99% shared
│       │   ├── shared_k8s_client.go    # 100% shared
│       │   └── portforward.go     # 100% shared
│       ├── s3/
│       │   ├── client.go          # Base S3 client (450 lines)
│       │   ├── client_factory.go  # Factory interface
│       │   └── client_options.go  # S3ClientOptions
│       └── pipelineserver/
│           ├── client.go          # Base pipeline client (487 lines)
│           └── client_factory.go  # Factory interface
```

**Product-specific extensions remain in:**
- `packages/automl/bff/internal/models/` - model_registry.go, register_model.go
- `packages/automl/bff/internal/models/pipeline_runs.go` - CreateAutoMLRunRequest
- `packages/automl/bff/internal/integrations/s3/` - CSV schema inference (GetCSVSchema)
- `packages/automl/bff/internal/integrations/pipelineserver/` - ListPipelines, ListPipelineVersions
- `packages/automl/bff/internal/integrations/modelregistry/` - entire directory
- `packages/autorag/bff/internal/models/` - lsd_models.go, lsd_vector_stores.go
- `packages/autorag/bff/internal/models/pipeline_runs.go` - CreateAutoRAGRunRequest
- `packages/autorag/bff/internal/integrations/llamastack/` - entire directory

---

## Implementation Notes

1. **Import Path**: Shared code will use import path like:
   ```go
   import "github.com/opendatahub-io/odh-dashboard/packages/autox-shared/bff/internal/models"
   import "github.com/opendatahub-io/odh-dashboard/packages/autox-shared/bff/internal/integrations/kubernetes"
   ```

2. **Product Extensions**: Product-specific code can embed or extend shared types:
   ```go
   // In AutoML
   type AutoMLPipelineServerClient struct {
       *shared.RealPipelineServerClient
   }
   
   func (c *AutoMLPipelineServerClient) ListPipelines(ctx context.Context, filter string) (*models.KFPipelinesResponse, error) {
       // AutoML-specific implementation
   }
   ```

3. **CSV Schema in S3**: Consider making GetCSVSchema() optional/extensible:
   ```go
   // Shared library
   type S3ClientInterface interface {
       GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, string, error)
       UploadObject(ctx context.Context, bucket, key string, body io.Reader, contentType string) error
       ListObjects(ctx context.Context, bucket string, options ListObjectsOptions) (*models.S3ListObjectsResponse, error)
       ObjectExists(ctx context.Context, bucket, key string) (bool, error)
   }
   
   // AutoML extension
   type S3ClientWithCSV interface {
       S3ClientInterface
       GetCSVSchema(ctx context.Context, bucket, key string) (CSVSchemaResult, error)
   }
   ```

4. **Pipeline Discovery**: Consider making pipeline discovery methods optional:
   ```go
   // Shared library - base interface
   type PipelineServerClient interface {
       CreateRun(ctx context.Context, req *models.CreatePipelineRunKFRequest) (*models.PipelineRun, error)
       GetRun(ctx context.Context, runID string) (*models.PipelineRun, error)
       GetRunDetails(ctx context.Context, runID string) (*models.PipelineRun, error)
       GetPipelineVersion(ctx context.Context, pipelineID, versionID string) (*models.KFPipelineVersion, error)
   }
   
   // AutoML extension - add discovery methods
   type PipelineServerClientWithDiscovery interface {
       PipelineServerClient
       ListPipelines(ctx context.Context, filter string) (*models.KFPipelinesResponse, error)
       ListPipelineVersions(ctx context.Context, pipelineID string) (*models.KFPipelineVersionsResponse, error)
   }
   ```

5. **Testing**: Mocks and tests will need to be updated to reference shared library types

6. **Migration Path**:
   - Phase 1: Extract 100% identical files (models + basic integrations)
   - Phase 2: Extract kubernetes integration (~1155 lines)
   - Phase 3: Extract base S3 client (450 lines), keep CSV schema in AutoML
   - Phase 4: Extract base pipeline server client (487 lines), keep discovery in AutoML
   - Phase 5: Update all imports in AutoML and AutoRAG to reference shared library
   - Phase 6: Add product-specific extensions where needed

---

## Risk Analysis

### Low Risk (100% Duplication)
- models/groups.go, health_check.go, namespace.go, pipeline_server.go, rbac_types.go, s3.go, secret.go, user.go
- integrations/http.go
- integrations/kubernetes/types.go, factory.go, shared_k8s_client.go, portforward.go

### Medium Risk (95-99% Duplication)
- models/pipeline_runs.go - Need to handle product-specific request types
- integrations/kubernetes/client.go, internal_k8s_client.go, token_k8s_client.go - Minor error handling differences

### Higher Risk (Requires Careful Design)
- integrations/s3/client.go - AutoML has 416 additional lines for CSV schema
- integrations/pipelineserver/client.go - AutoML has 144 additional lines for pipeline discovery

### Testing Requirements
- All shared code must have comprehensive unit tests
- Integration tests must verify both AutoML and AutoRAG still work correctly
- Mock implementations must be updated to reference shared types

---

## Conclusion

This analysis identified **~2588 lines of highly duplicated code** across models and integrations that can be extracted to a shared library. The extraction is feasible with proper interface design to accommodate product-specific extensions. The highest value extractions are:

1. **Models directory** (8-9 files, ~444 lines) - 100% identical, zero risk
2. **Kubernetes integration** (~1155 lines) - 95-100% identical, low risk
3. **S3 base client** (450 lines) - Medium risk, requires extension design
4. **Pipeline server base client** (487 lines) - Medium risk, requires extension design

The recommended approach is phased extraction starting with zero-risk 100% identical files, followed by carefully designed interfaces for the partially shared integrations.
