# AutoX BFF API Contract

**Package**: `@odh-dashboard/autox/bff`  
**Language**: Go 1.24.3  
**Status**: Draft  
**Version**: 1.0.0-alpha  
**Date**: 2026-04-23

---

## Overview

This document defines the public API contract for the AutoX BFF (Backend For Frontend) shared library. All interfaces, types, and functions exported by AutoX BFF are guaranteed to maintain backward compatibility within the same major version. Breaking changes will only occur with major version bumps.

**Target Consumers**: AutoML BFF, AutoRAG BFF (monorepo-internal only)

**Import Path Convention**:

```go
import (
    "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/models"
    "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/integrations/kubernetes"
    "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/integrations/s3"
    "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/integrations/pipelineserver"
    "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/repositories"
    "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/api"
    "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/utils"
)
```

**Note**: While these are `internal/` packages (not `pkg/`), they are intended for consumption by AutoML and AutoRAG within the same Go module (monorepo). The `internal/` directory prevents external packages outside the monorepo from importing AutoX.

---

## API Stability Guarantees

### Stable API (Backward Compatible)

- **Models** (`internal/models/*`): All struct fields, JSON tags, validation
- **Interfaces** (`internal/integrations/*/interface.go`): Method signatures, parameters, return types
- **Repositories** (`internal/repositories/*`): Public methods, data structures
- **API Utilities** (`internal/api/*`): Error response format, helper function signatures

### Unstable API (May Change)

- **Internal helpers**: Unexported functions (lowercase names)
- **Test utilities**: Mock implementations, test fixtures
- **Implementation details**: Caching strategies, internal data structures not exposed via interfaces

### Breaking Changes Policy

- **Synchronized versioning**: AutoX, AutoML, and AutoRAG are always updated together in the same PR
- **No independent versioning**: Breaking changes are coordinated and applied atomically
- **CI type-checking**: All packages are compiled together (`go build`) to catch breaking changes before merge
- **Migration period**: N/A (synchronized updates mean no migration period needed)

---

## 1. Models (`internal/models`)

### 1.1 Core Models

#### HealthCheck

```go
package models

// HealthCheck represents the health check response
type HealthCheck struct {
    Status string `json:"status"` // "healthy" or "unhealthy"
    Time   int64  `json:"timestamp"`
}
```

**Stability**: ✅ Stable  
**Usage**: Health check endpoint responses

---

#### Namespace

```go
package models

// Namespace represents a Kubernetes namespace
type Namespace struct {
    Name        string `json:"name"`
    DisplayName string `json:"displayName,omitempty"`
    Status      string `json:"status"`
}
```

**Stability**: ✅ Stable  
**Validation**: 
- `Name` must match RFC 1123 DNS subdomain (lowercase alphanumeric, hyphens, max 63 chars)
- `Status` is "Active" or "Terminating"

**Usage**: Namespace listing, filtering

---

#### User

```go
package models

// User represents an authenticated user
type User struct {
    ID       string   `json:"id"`
    Username string   `json:"username"`
    Email    string   `json:"email,omitempty"`
    Groups   []string `json:"groups"`
}
```

**Stability**: ✅ Stable  
**Usage**: User identity extraction from request headers

---

#### Secret

```go
package models

// Secret represents a Kubernetes secret with metadata
type Secret struct {
    Name        string            `json:"name"`
    Namespace   string            `json:"namespace"`
    Type        string            `json:"type"` // "s3", "generic", "opaque"
    Keys        []string          `json:"keys"`
    Labels      map[string]string `json:"labels,omitempty"`
    Annotations map[string]string `json:"annotations,omitempty"`
}
```

**Stability**: ✅ Stable  
**Usage**: Secret filtering, credential retrieval

---

#### S3Credentials

```go
package models

// S3Credentials represents S3 access credentials
type S3Credentials struct {
    AccessKeyID     string `json:"accessKeyId"`
    SecretAccessKey string `json:"secretAccessKey"`
    Endpoint        string `json:"endpoint,omitempty"`
    Region          string `json:"region"`
    Bucket          string `json:"bucket,omitempty"`
}
```

**Stability**: ✅ Stable  
**Validation**: 
- `AccessKeyID` and `SecretAccessKey` must not be empty
- `Region` must not be empty
- `Endpoint` is optional (defaults to AWS S3)

**Usage**: S3 client configuration

---

#### PipelineRun

```go
package models

// RuntimeStateKF represents Kubeflow pipeline run states
type RuntimeStateKF string

const (
    RuntimeStatePending   RuntimeStateKF = "PENDING"
    RuntimeStateRunning   RuntimeStateKF = "RUNNING"
    RuntimeStateSucceeded RuntimeStateKF = "SUCCEEDED"
    RuntimeStateFailed    RuntimeStateKF = "FAILED"
    RuntimeStateCanceled  RuntimeStateKF = "CANCELED"
    RuntimeStateSkipped   RuntimeStateKF = "SKIPPED"
)

// PipelineRun represents a pipeline execution
type PipelineRun struct {
    ID              string                 `json:"id"`
    Name            string                 `json:"name"`
    State           RuntimeStateKF         `json:"state"`
    CreatedAt       string                 `json:"createdAt"` // ISO 8601 timestamp
    FinishedAt      string                 `json:"finishedAt,omitempty"`
    Parameters      map[string]interface{} `json:"parameters"`
    PipelineType    string                 `json:"pipelineType"`    // e.g., "time-series", "rag"
    PipelineVersion string                 `json:"pipelineVersion"` // Pipeline version ID
    Error           *PipelineError         `json:"error,omitempty"`
}

// PipelineError represents an error from a pipeline run
type PipelineError struct {
    Message string `json:"message"`
    Code    string `json:"code,omitempty"`
}
```

**Stability**: ✅ Stable  
**Validation**:
- `ID` must be non-empty UUID string
- `State` must be valid `RuntimeStateKF` value
- `CreatedAt` must be ISO 8601 timestamp
- `Name` must not be empty

**Usage**: Pipeline run CRUD operations, state tracking

---

#### PipelineServer

```go
package models

// PipelineServer represents a KFP server configuration
type PipelineServer struct {
    Name      string `json:"name"`
    Namespace string `json:"namespace"`
    BaseURL   string `json:"baseUrl"`
    Version   string `json:"version,omitempty"`
}
```

**Stability**: ✅ Stable  
**Usage**: Pipeline server discovery, connection management

---

### 1.2 Request/Response Models

#### CreatePipelineRunRequest

```go
package models

// CreatePipelineRunRequest represents a request to create a pipeline run
type CreatePipelineRunRequest struct {
    PipelineVersionID string                 `json:"pipelineVersionId"`
    DisplayName       string                 `json:"displayName"`
    Parameters        map[string]interface{} `json:"parameters"`
}
```

**Stability**: ✅ Stable  
**Validation**:
- `PipelineVersionID` must be valid UUID
- `DisplayName` must not be empty
- `Parameters` must match pipeline definition schema (validated by pipeline server)

---

#### ErrorResponse

```go
package models

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
    Code    string `json:"code"`    // e.g., "VALIDATION_ERROR", "NOT_FOUND"
    Message string `json:"message"` // Human-readable error message
    Details string `json:"details,omitempty"` // Additional error context
}
```

**Stability**: ✅ Stable  
**Usage**: HTTP error responses

---

## 2. Integration Interfaces (`internal/integrations`)

### 2.1 Kubernetes Client Interface

```go
package kubernetes

import (
    "context"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KubernetesClientInterface defines the contract for Kubernetes API operations
type KubernetesClientInterface interface {
    // Namespace operations
    GetNamespace(ctx context.Context, name string) (*corev1.Namespace, error)
    ListNamespaces(ctx context.Context, opts metav1.ListOptions) (*corev1.NamespaceList, error)
    
    // Secret operations
    GetSecret(ctx context.Context, namespace, name string) (*corev1.Secret, error)
    ListSecrets(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.SecretList, error)
    
    // Generic resource operations
    Get(ctx context.Context, namespace, name string, obj interface{}) error
    List(ctx context.Context, namespace string, opts metav1.ListOptions, obj interface{}) error
    Create(ctx context.Context, namespace string, obj interface{}) error
    Update(ctx context.Context, namespace, name string, obj interface{}) error
    Delete(ctx context.Context, namespace, name string) error
}
```

**Stability**: ✅ Stable  
**Extractable**: 100% identical (800 LOC including implementation)

**Default Implementation**: AutoX provides `InternalK8sClient` and `TokenK8sClient` implementations

**Usage Example**:

```go
// AutoML handler
import k8s "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/integrations/kubernetes"

func (h *Handler) GetSecrets(ctx context.Context, namespace string) error {
    secrets, err := h.k8sClient.ListSecrets(ctx, namespace, metav1.ListOptions{
        LabelSelector: "type=s3",
    })
    if err != nil {
        return err
    }
    // Process secrets...
}
```

**Error Handling**: Returns standard Go errors. Caller should check error types for retries.

---

### 2.2 S3 Client Interface

```go
package s3

import (
    "context"
    "io"
)

// S3Object represents an object in S3
type S3Object struct {
    Key          string
    Size         int64
    LastModified string // ISO 8601 timestamp
    ETag         string
}

// S3ClientInterface defines the contract for S3 operations
type S3ClientInterface interface {
    // Object operations
    GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error)
    PutObject(ctx context.Context, bucket, key string, body io.Reader) error
    DeleteObject(ctx context.Context, bucket, key string) error
    ListObjects(ctx context.Context, bucket, prefix string, recursive bool) ([]S3Object, error)
    
    // Bucket operations
    CreateBucket(ctx context.Context, bucket string) error
    DeleteBucket(ctx context.Context, bucket string) error
    BucketExists(ctx context.Context, bucket string) (bool, error)
    
    // Validation
    ValidateCredentials(ctx context.Context) error
}
```

**Stability**: ✅ Stable  
**Extractable**: 95% identical (900 LOC)

**Default Implementation**: AutoX provides `AWSS3Client` wrapping AWS SDK v2

**Usage Example**:

```go
// AutoRAG handler
import s3 "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/integrations/s3"

func (h *Handler) UploadFile(ctx context.Context, bucket, key string, file io.Reader) error {
    err := h.s3Client.PutObject(ctx, bucket, key, file)
    if err != nil {
        return fmt.Errorf("failed to upload to S3: %w", err)
    }
    return nil
}
```

**Error Handling**: Wraps AWS SDK errors with context. Check for `awserr.Error` types for retries.

---

### 2.3 Pipeline Server Client Interface

```go
package pipelineserver

import (
    "context"
    "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/models"
)

// PipelineServerClientInterface defines the contract for KFP operations
type PipelineServerClientInterface interface {
    // Pipeline operations
    GetPipeline(ctx context.Context, pipelineID string) (*models.Pipeline, error)
    ListPipelines(ctx context.Context, filter string) ([]*models.Pipeline, error)
    UploadPipeline(ctx context.Context, name string, yamlContent []byte) (*models.Pipeline, error)
    
    // Pipeline run operations
    CreatePipelineRun(ctx context.Context, req *models.CreatePipelineRunRequest) (*models.PipelineRun, error)
    GetPipelineRun(ctx context.Context, runID string) (*models.PipelineRun, error)
    ListPipelineRuns(ctx context.Context, filter string) ([]*models.PipelineRun, error)
    TerminatePipelineRun(ctx context.Context, runID string) error
    RetryPipelineRun(ctx context.Context, runID string) (*models.PipelineRun, error)
    
    // Health
    HealthCheck(ctx context.Context) error
}
```

**Stability**: ✅ Stable  
**Extractable**: 100% identical interface (630 LOC implementation)

**Default Implementation**: AutoX provides `KFPClient` using KFP REST API

**Usage Example**:

```go
// AutoML handler
import ps "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/integrations/pipelineserver"

func (h *Handler) CreateRun(ctx context.Context, req *models.CreatePipelineRunRequest) error {
    run, err := h.psClient.CreatePipelineRun(ctx, req)
    if err != nil {
        return fmt.Errorf("failed to create pipeline run: %w", err)
    }
    log.Printf("Created pipeline run: %s", run.ID)
    return nil
}
```

**Error Handling**: Returns wrapped KFP API errors. Includes HTTP status codes in error context.

---

## 3. Repository Patterns (`internal/repositories`)

### 3.1 Pipeline Repository

```go
package repositories

import (
    "context"
    ps "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/integrations/pipelineserver"
)

// DiscoveredPipeline represents a pipeline found via discovery
type DiscoveredPipeline struct {
    ID              string
    VersionID       string
    Name            string
    DefinitionName  string // e.g., "time-series", "rag-pipeline"
}

// PipelineDefinition represents a pipeline YAML definition
type PipelineDefinition struct {
    Name     string // Definition name (e.g., "time-series")
    YAMLPath string // Path to YAML file
}

// PipelineRepository handles pipeline discovery and caching
type PipelineRepository interface {
    // DiscoverNamedPipelines finds pipelines by name and caches results
    DiscoverNamedPipelines(
        client ps.PipelineServerClientInterface,
        ctx context.Context,
        namespace string,
        pipelineServerBaseURL string,
        definitions map[string]string, // map[definitionName]yamlPath
    ) (map[string]*DiscoveredPipeline, error)
    
    // EnsurePipeline ensures a pipeline exists (uploads if missing)
    EnsurePipeline(
        client ps.PipelineServerClientInterface,
        ctx context.Context,
        namespace string,
        pipelineServerBaseURL string,
        def PipelineDefinition,
    ) (*DiscoveredPipeline, error)
    
    // InvalidateCache clears cached pipeline data
    InvalidateCache(pipelineServerBaseURL, namespace string)
}
```

**Stability**: ✅ Stable  
**Extractable**: 97% identical (1,069 LOC)

**Default Implementation**: AutoX provides `PipelineRepositoryImpl` with LRU caching

**Usage Example**:

```go
// AutoML handler
import (
    "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/repositories"
)

func (h *Handler) DiscoverPipelines(ctx context.Context, namespace string) error {
    // Domain-specific definitions
    definitions := map[string]string{
        "time-series": "/pipelines/time_series_pipeline.yaml",
        "tabular":     "/pipelines/tabular_pipeline.yaml",
    }
    
    pipelines, err := h.pipelineRepo.DiscoverNamedPipelines(
        h.psClient, ctx, namespace, h.pipelineServerBaseURL, definitions,
    )
    if err != nil {
        return err
    }
    
    // Use discovered pipelines
    timeSeriesPipeline := pipelines["time-series"]
    log.Printf("Found time-series pipeline: %s", timeSeriesPipeline.ID)
    return nil
}
```

**Caching Strategy**:
- LRU cache with max 100 entries
- Cache key: `{pipelineServerBaseURL}:{namespace}`
- TTL: 5 minutes (configurable)
- Thread-safe

---

### 3.2 Secret Repository

```go
package repositories

import (
    "context"
    k8s "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/integrations/kubernetes"
    "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/models"
)

// SecretRepository handles secret filtering and retrieval
type SecretRepository interface {
    // GetFilteredSecrets retrieves secrets filtered by type
    GetFilteredSecrets(
        client k8s.KubernetesClientInterface,
        ctx context.Context,
        namespace string,
        identity *k8s.RequestIdentity,
        secretType string, // "s3", "generic", "all"
    ) ([]models.Secret, error)
}
```

**Stability**: ✅ Stable  
**Extractable**: 95% shared (base filtering logic)

**Default Implementation**: AutoX provides `SecretRepositoryImpl`

**Usage Example**:

```go
// AutoRAG handler (with custom filtering)
import (
    "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/repositories"
)

func (h *Handler) GetLlamaStackSecrets(ctx context.Context, namespace string) ([]models.Secret, error) {
    // Call shared repository for base filtering
    allSecrets, err := h.secretRepo.GetFilteredSecrets(
        h.k8sClient, ctx, namespace, h.identity, "all",
    )
    if err != nil {
        return nil, err
    }
    
    // AutoRAG-specific filtering (handler logic, not in AutoX)
    llamaStackSecrets := filterLlamaStackSecrets(allSecrets)
    return llamaStackSecrets, nil
}

// filterLlamaStackSecrets is AutoRAG-specific (not in AutoX)
func filterLlamaStackSecrets(secrets []models.Secret) []models.Secret {
    var result []models.Secret
    for _, secret := range secrets {
        if secret.Labels["app"] == "llamastack" {
            result = append(result, secret)
        }
    }
    return result
}
```

**Customization Pattern**: Consumers call AutoX base filtering, then apply domain-specific filters in handler files.

---

## 4. API Utilities (`internal/api`)

### 4.1 Error Response Helpers

```go
package api

import (
    "net/http"
    "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/models"
)

// ErrorHandler provides standardized error responses
type ErrorHandler interface {
    BadRequestResponse(w http.ResponseWriter, r *http.Request, err error)
    UnauthorizedResponse(w http.ResponseWriter, r *http.Request, message string)
    ForbiddenResponse(w http.ResponseWriter, r *http.Request, message string)
    NotFoundResponse(w http.ResponseWriter, r *http.Request, resource string)
    ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error)
}

// NewErrorHandler creates a default error handler
func NewErrorHandler() ErrorHandler

// Standalone error response functions
func WriteBadRequest(w http.ResponseWriter, message string)
func WriteNotFound(w http.ResponseWriter, resource string)
func WriteServerError(w http.ResponseWriter, err error)
```

**Stability**: ✅ Stable  
**Extractable**: 99% identical (181 LOC)

**Usage Example**:

```go
// AutoML handler
import "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/api"

func (h *Handler) HandleRequest(w http.ResponseWriter, r *http.Request) {
    data, err := h.service.GetData(r.Context())
    if err != nil {
        api.WriteServerError(w, err)
        return
    }
    
    api.WriteJSON(w, http.StatusOK, data, nil)
}
```

**Error Response Format** (Guaranteed):

```json
{
  "code": "VALIDATION_ERROR",
  "message": "Invalid pipeline name",
  "details": "Name must match RFC 1123 DNS subdomain"
}
```

---

### 4.2 Request/Response Helpers

```go
package api

import (
    "net/http"
)

// WriteJSON writes a JSON response with status code and optional headers
func WriteJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error

// ReadJSON reads and parses JSON request body
func ReadJSON(r *http.Request, dst interface{}) error
```

**Stability**: ✅ Stable  
**Extractable**: 100% identical (108 LOC)

**Usage Example**:

```go
import "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/api"

func (h *Handler) CreatePipeline(w http.ResponseWriter, r *http.Request) {
    var req models.CreatePipelineRunRequest
    if err := api.ReadJSON(r, &req); err != nil {
        api.WriteBadRequest(w, "Invalid request body")
        return
    }
    
    // Process request...
    result := map[string]string{"id": "123"}
    api.WriteJSON(w, http.StatusCreated, result, nil)
}
```

---

### 4.3 Middleware

```go
package api

import (
    "net/http"
)

// CORSMiddleware creates a CORS middleware handler
func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler

// LoggingMiddleware creates a request logging middleware
func LoggingMiddleware(logger Logger) func(http.Handler) http.Handler

// Logger interface for middleware
type Logger interface {
    Info(msg string, fields ...interface{})
    Error(msg string, err error, fields ...interface{})
}
```

**Stability**: ✅ Stable  
**Extractable**: 85% shared (middleware setup)

**Usage Example**:

```go
import "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/api"

func setupRouter() http.Handler {
    router := httprouter.New()
    
    // Apply AutoX middleware
    handler := api.CORSMiddleware([]string{"*"})(router)
    handler = api.LoggingMiddleware(logger)(handler)
    
    return handler
}
```

---

## 5. Utilities (`internal/utils`)

### 5.1 Pagination Helpers

```go
package utils

// PaginationParams represents pagination query parameters
type PaginationParams struct {
    Page    int
    PerPage int
    Offset  int
}

// ParsePaginationParams parses pagination query parameters
func ParsePaginationParams(r *http.Request) PaginationParams

// ApplyPagination applies pagination to a slice
func ApplyPagination[T any](items []T, params PaginationParams) []T
```

**Stability**: ✅ Stable  
**Usage**: Query parameter parsing, list slicing

---

### 5.2 Logging Helpers

```go
package utils

import "log/slog"

// NewLogger creates a structured logger
func NewLogger(level slog.Level) *slog.Logger

// LogError logs an error with context
func LogError(logger *slog.Logger, msg string, err error, fields ...interface{})
```

**Stability**: ✅ Stable  
**Extractable**: 100% identical (125 LOC)

---

## 6. Version Compatibility

### Current Version: 1.0.0-alpha

**Supported Go Versions**: 1.24.3+  
**Supported K8s Client Versions**: 0.31.x  
**Supported AWS SDK Versions**: v2.x

### Deprecation Policy

- **Deprecation notices**: Added as code comments and in release notes
- **Deprecation period**: Minimum 2 releases (or 1 month)
- **Removal**: Only in major version bumps

**Example**:

```go
// Deprecated: Use GetFilteredSecrets instead.
// Will be removed in v2.0.0.
func GetSecrets(namespace string) ([]models.Secret, error) {
    // ...
}
```

---

## 7. Testing Strategy

### Consumer Testing

**AutoML and AutoRAG must test**:
1. **Import validation**: All AutoX imports compile without errors
2. **Interface compliance**: Custom implementations satisfy AutoX interfaces
3. **Integration tests**: End-to-end flows using AutoX utilities

**AutoX provides**:
- Unit tests for all exported functions
- Mock implementations for all interfaces (in `internal/mocks/`)
- Example usage in `internal/examples/` (not part of stable API)

---

## 8. Migration Guide

### From Duplicate Code to AutoX

**Step 1**: Identify duplicate code in AutoML/AutoRAG  
**Step 2**: Find corresponding AutoX export in this contract  
**Step 3**: Update imports:

```go
// Before (AutoML internal code)
import "github.com/opendatahub-io/odh-dashboard/packages/automl/bff/internal/models"

// After (AutoX shared code)
import "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/models"
```

**Step 4**: Remove duplicate code from consumer package  
**Step 5**: Run tests to verify behavior preserved

---

## 9. Contact & Support

**Maintainers**: ODH Dashboard Team  
**Issues**: Report via GitHub Issues with `autox` label  
**Breaking Changes**: Announced in PR description and release notes  
**CI/CD**: All packages compile together (`go build ./...` from repo root)

---

**Next**: See [ui-api.md](./ui-api.md) for frontend contract, [../quickstart.md](../quickstart.md) for setup guide.
