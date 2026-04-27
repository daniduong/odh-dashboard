# `autox-core` BFF Architecture & Refactoring Plan

## Context

**Objective**: Design and build autox-core as a shared BFF library to enable massive refactoring and deduplication of automl and autorag packages.

**Current State**:
- AutoML and AutoRAG have ~80-90% duplicate code in their BFF layers
- Both packages use identical patterns for Kubernetes, Pipelines, and S3 integrations
- Business logic is scattered across handlers, middleware, and repositories
- New pipelines client created per request (wasteful)
- Some handlers and middleware contain business logic (should be thin)

**Goal**: 
- Extract common building blocks into autox-core/bff
- Establish clean architectural boundaries (handlers → services → clients)
- Enable code reuse while allowing domain-specific customization
- Fix architectural issues during the migration

---

## Architectural Decisions

### Core Principles

1. **Services-Only Library**: `autox-core` exports service layer with clean request/response types, NOT an HTTP framework
2. **Explicit Client Parameters**: Client objects (http.Client, S3Client) passed as function parameters, NOT via context.WithValue. Entities (identity, namespace, DSPA) stored in context for middleware flow.
3. **Domain Models**: Services return business entities, handlers transform to JSON responses
4. **Feature-Based Structure**: Organize by integration (kubernetes/, pipelines/, storage/) not by layer
5. **Consumer-Owned App**: AutoML/AutoRAG own their App struct and initialization logic

### Specific Decisions

| Aspect | Decision | Rationale |
|--------|----------|-----------|
| **K8s Factory** | Single unified factory | Cleaner API, autox-core handles factory selection internally |
| **Middleware** | Services only, no middleware exports | Consumers write their own thin middleware that calls autox-core services |
| **Client Passing** | Explicit parameters (no client objects in context) | Entities (identity, DSPA) go in context; clients passed as function parameters |
| **HTTP Client** | Shared http.Client in App, passed to services | Efficient connection pooling, consumers control configuration |
| **Caching** | Include in autox-core | Both automl/autorag need identical pipeline/DSPA discovery caching |
| **Request Types** | In autox-core | Consistent API across consumers, defines service contracts |
| **Handler Layer** | Service layer only | Consumers write thin handlers, autox-core focuses on business logic |
| **Discovery** | DSPA & pipeline discovery in autox-core | Common pattern with caching, core business logic |
| **Validation** | Domain-specific validators | Both generic (DNS-1123) and domain (ValidatePipelineRunRequest) |
| **Migration** | Architecture-first | Build complete autox-core, then migrate automl, then autorag |
| **Return Types** | Domain models | Clean separation between business logic and API contracts |
| **Auth** | Accept RequestIdentity as parameter | Consumer middleware extracts, services receive clean input |
| **Errors** | Custom error types | Easier HTTP status mapping (NotFoundError, ForbiddenError, etc) |
| **Directory** | By feature (kubernetes/, pipelines/) | Cohesive per-feature packages |

---

## `autox-core` Directory Structure

```
packages/autox-core/bff/
├── kubernetes/              # Kubernetes integration
│   ├── models.go           # Domain models (Namespace, Secret, RequestIdentity)
│   ├── client.go           # K8s client interface
│   ├── factory.go          # Unified factory (handles static vs token auth)
│   ├── service.go          # KubernetesService (business logic)
│   ├── validators.go       # Input validation (namespace format, RBAC checks)
│   └── errors.go           # Custom errors (NotFoundError, ForbiddenError)
│
├── pipelines/              # Kubeflow Pipelines integration
│   ├── models.go           # Domain models (PipelineRun, Pipeline, DSPA)
│   ├── client.go           # HTTP client for KFP API
│   ├── service.go          # PipelinesService (discovery, caching, CRUD)
│   ├── discovery.go        # DSPA discovery, pipeline discovery with LRU cache
│   ├── validators.go       # Request validation
│   └── errors.go           # Custom errors
│
├── storage/                # S3/MinIO integration
│   ├── models.go           # Domain models (S3Object, Bucket)
│   ├── client.go           # S3 client interface
│   ├── service.go          # StorageService
│   ├── validators.go       # Path validation, SSRF protection
│   └── errors.go           # Custom errors
│
└── common/                 # Shared utilities
    ├── errors/             # Base error types and helpers
    ├── validation/         # Generic validators (DNS-1123, URLs)
    └── cache/              # LRU cache implementation
```

**Key Pattern**: Each feature package is self-contained with models, client, service, validators, and errors.

---

## Service Layer API Design

### Example: Kubernetes Service

```go
package kubernetes

// Request types define service contracts
type ListNamespacesRequest struct {
    Identity *RequestIdentity  // Explicit auth parameter
}

type GetSecretRequest struct {
    Identity  *RequestIdentity
    Namespace string
    Name      string
}

// Domain models (returned by services)
type Namespace struct {
    Name        string
    DisplayName string
    Labels      map[string]string
}

type Secret struct {
    Name      string
    Namespace string
    Data      map[string][]byte
}

// Service interface
type KubernetesService interface {
    // List accessible namespaces with RBAC filtering
    ListNamespaces(ctx context.Context, req ListNamespacesRequest) ([]Namespace, error)
    
    // Get secret with validation
    GetSecret(ctx context.Context, req GetSecretRequest) (*Secret, error)
    
    // Check cluster admin permission
    IsClusterAdmin(ctx context.Context, identity *RequestIdentity) (bool, error)
    
    // Check namespace-level permission
    CanAccessNamespace(ctx context.Context, identity *RequestIdentity, namespace string) (bool, error)
}

// Implementation holds dependencies
type kubernetesService struct {
    clientFactory KubernetesClientFactory
    validator     *Validator
    logger        *slog.Logger
}

// Consumer creates service with dependencies
func NewKubernetesService(
    clientFactory KubernetesClientFactory,
    logger *slog.Logger,
) KubernetesService {
    return &kubernetesService{
        clientFactory: clientFactory,
        validator:     NewValidator(),
        logger:        logger,
    }
}
```

### Example: Pipelines Service

```go
package pipelines

type ListPipelineRunsRequest struct {
    Identity          *RequestIdentity
    Namespace         string
    PipelineVersionID string
    PageSize          int32
    PageToken         string
}

type CreatePipelineRunRequest struct {
    Identity     *RequestIdentity
    Namespace    string
    PipelineID   string
    VersionID    string
    Name         string
    Parameters   map[string]string
}

// Domain model
type PipelineRun struct {
    ID          string
    Name        string
    Status      string
    CreatedAt   time.Time
    FinishedAt  *time.Time
    Parameters  map[string]string
}

type PipelinesService interface {
    // Discover DSPA in namespace (with caching)
    DiscoverDSPA(ctx context.Context, req DiscoverDSPARequest) (*DSPA, error)
    
    // Discover pipelines by name prefix (with caching)
    DiscoverPipelines(ctx context.Context, req DiscoverPipelinesRequest) (map[string]*DiscoveredPipeline, error)
    
    // List pipeline runs
    ListPipelineRuns(ctx context.Context, client *http.Client, baseURL string, req ListPipelineRunsRequest) ([]PipelineRun, string, error)
    
    // Create pipeline run
    CreatePipelineRun(ctx context.Context, client *http.Client, baseURL string, req CreatePipelineRunRequest) (*PipelineRun, error)
}

type pipelinesService struct {
    dspaCache      *cache.LRUCache  // 5min TTL, 1000 entries
    pipelineCache  *cache.LRUCache  // 5min TTL, 1000 entries
    validator      *Validator
    logger         *slog.Logger
}

func NewPipelinesService(logger *slog.Logger) PipelinesService {
    return &pipelinesService{
        dspaCache:     cache.NewLRU(1000, 5*time.Minute),
        pipelineCache: cache.NewLRU(1000, 5*time.Minute),
        validator:     NewValidator(),
        logger:        logger,
    }
}
```

**Key Points**:
- Services receive `*http.Client` and `baseURL` explicitly (no baked-in client)
- Caching is internal to services
- Request structs make API clear and future-proof
- Domain models are clean business entities

---

## Consumer App Pattern (AutoML/AutoRAG)

### Three-Layer Architecture (Domain Service Layer is OPTIONAL)

```
┌─────────────────────────────────────────────────────────────┐
│ Handler Layer (packages/automl/bff/internal/api/)          │
│ - HTTP request/response handling                           │
│ - Parameter extraction, JSON marshaling                     │
│ - Error mapping (domain errors → HTTP status codes)         │
└─────────────────────────────────────────────────────────────┘
         ↓                                    ↓
   (complex ops)                         (simple ops)
         ↓                                    ↓
┌─────────────────────────┐                  │
│ Domain Service Layer    │                  │
│ (OPTIONAL - only if     │                  │
│  domain logic exists)   │                  │
│                         │                  │
│ - Domain validation     │                  │
│ - Parameter building    │                  │
│ - Multi-step orchestration                │
│ - Composes autox-core   │                  │
└─────────────────────────┘                  │
         ↓                                    ↓
         └────────────────┬───────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│ `autox-core` Service Layer (packages/autox-core/bff/)        │
│ - Common Kubernetes operations (RBAC, secrets, namespaces)  │
│ - Common Pipelines operations (DSPA discovery, CRUD runs)   │
│ - Common Storage operations (S3 upload/download/list)       │
│ - Reusable across AutoML, AutoRAG, future packages          │
└─────────────────────────────────────────────────────────────┘
                          ↓
                 External Services
         (K8s API, KFP API, S3/MinIO, etc)
```

**Key Principle: Domain Service Layer is OPTIONAL**

- **When NO domain-specific logic** → Handler calls autox-core directly
- **When domain-specific logic EXISTS** → Handler calls domain service → domain service calls autox-core

**Examples of NO Domain Logic** (call autox-core directly):
- List pipeline runs (just pagination, no AutoML-specific params)
- Get pipeline run by ID (just fetch and return)
- List namespaces (just RBAC filtering)
- Get secret (just retrieve and return)
- Terminate pipeline run (just terminate, no orchestration)

**Examples of Domain Logic** (use domain service):
- Create AutoML pipeline run (validate problem type, build AutoML params, validate dataset exists)
- Upload dataset (validate CSV schema, add AutoML metadata, create S3 object)
- Register model (validate model artifact, interact with model registry, create pipeline run)
- Create AutoRAG optimization run (validate vector store config, build RAG params)

**Example Flow: Create AutoML Pipeline Run**

1. **Handler** (`CreatePipelineRunHandler`)
   - Receives HTTP POST request
   - Parses JSON body → `CreateAutoMLRunAPIRequest`
   - Extracts entities from context (identity, namespace, DSPA, discoveredPipeline)
   - Calls `app.autoMLPipelineService.CreateAutoMLPipelineRun()`

2. **AutoML Service** (`AutoMLPipelineService.CreateAutoMLPipelineRun`)
   - Validates AutoML-specific params (problem_type, max_trials, etc)
   - Calls `storageService.ObjectExists()` to verify dataset exists
   - Builds AutoML-specific pipeline parameters map
   - Calls `pipelinesService.CreatePipelineRun()` with built params

3. **`autox-core` Service** (`PipelinesService.CreatePipelineRun`)
   - Validates generic pipeline params (namespace, pipelineID, etc)
   - Makes HTTP POST to KFP API `/apis/v2beta1/runs`
   - Parses KFP response → domain model `PipelineRun`
   - Returns domain model

4. **Response flows back up**:
   - `autox-core` returns `*PipelineRun` to AutoML service
   - AutoML service logs and returns `*PipelineRun` to handler
   - Handler maps `PipelineRun` → `CreateRunAPIResponse` JSON
   - Handler writes HTTP 201 response

**Example Flow: List Pipeline Runs (No Domain Logic)**

1. **Handler** (`PipelineRunsHandler`)
   - Receives HTTP GET request
   - Parses query params (pageSize, pageToken)
   - Extracts entities from context
   - **Calls autox-core directly** (no AutoML service needed)
   - `app.pipelinesService.ListPipelineRuns()`

2. **`autox-core` Service** (`PipelinesService.ListPipelineRuns`)
   - Makes HTTP GET to KFP API `/apis/v2beta1/runs`
   - Parses response → `[]PipelineRun`
   - Returns domain models

3. **Handler maps and returns**:
   - Maps `[]PipelineRun` → `PipelineRunsResponse` JSON
   - Writes HTTP 200 response

### App Structure

```go
// packages/automl/bff/internal/api/app.go
package api

type App struct {
    config     config.EnvConfig
    logger     *slog.Logger
    httpClient *http.Client  // Shared HTTP client for pipelines
    
    // `autox-core` services (low-level primitives)
    // Handlers call these DIRECTLY for simple operations (list, get, etc)
    kubernetesService k8s.KubernetesService
    pipelinesService  ps.PipelinesService
    storageService    storage.StorageService
    
    // AutoML-specific services (OPTIONAL - only created when domain logic exists)
    // Handlers call these for complex operations with domain-specific logic
    autoMLPipelineService *services.AutoMLPipelineService  // Composes pipelines + storage
    modelRegistryService  *services.ModelRegistryService   // Composes modelregistry + pipelines
    
    // AutoML-specific config
    autoMLPipelinePrefix string
}

func NewApp(cfg config.EnvConfig, logger *slog.Logger) (*App, error) {
    // Create shared HTTP client (connection pooling, timeouts)
    httpClient := &http.Client{
        Timeout: 30 * time.Second,
        Transport: &http.Transport{
            MaxIdleConns:        100,
            MaxIdleConnsPerHost: 10,
            IdleConnTimeout:     90 * time.Second,
        },
    }
    
    // Initialize K8s factory (autox-core handles static vs token)
    k8sFactory := k8s.NewUnifiedFactory(cfg.AuthMethod, cfg.MockK8sClient, logger)
    
    // Create autox-core services
    kubernetesService := k8s.NewKubernetesService(k8sFactory, logger)
    pipelinesService := ps.NewPipelinesService(logger)
    storageService := storage.NewStorageService(logger)
    
    return &App{
        config:               cfg,
        logger:               logger,
        httpClient:           httpClient,
        kubernetesService:    kubernetesService,
        pipelinesService:     pipelinesService,
        storageService:       storageService,
        autoMLPipelinePrefix: cfg.AutoMLPipelinePrefix,
    }, nil
}
```

### Domain Service Layer (AutoML-Specific Business Logic)

```go
// packages/automl/bff/internal/services/automl_pipeline_service.go
package services

// AutoML-specific pipeline service (composes autox-core)
type AutoMLPipelineService struct {
    pipelinesService pipelines.PipelinesService  // `autox-core` service
    storageService   storage.StorageService      // `autox-core` service
    logger           *slog.Logger
}

// Domain-specific request (AutoML params)
type CreateAutoMLPipelineRunRequest struct {
    Identity           *kubernetes.RequestIdentity
    Namespace          string
    DSPA               *pipelines.DSPA
    DiscoveredPipeline *pipelines.DiscoveredPipeline
    
    // AutoML-specific fields
    DatasetPath        string
    TargetColumn       string
    ProblemType        string  // classification, regression, time-series
    MaxTrials          int
    TrainingBudget     string
    S3Credentials      *S3Credentials
}

func (s *AutoMLPipelineService) CreateAutoMLPipelineRun(
    ctx context.Context,
    httpClient *http.Client,
    req CreateAutoMLPipelineRunRequest,
) (*pipelines.PipelineRun, error) {
    // 1. AutoML-specific validation
    if err := s.validateAutoMLParams(req); err != nil {
        return nil, err
    }
    
    // 2. Validate dataset exists in S3 (using autox-core storage service)
    exists, err := s.storageService.ObjectExists(ctx, s3Client, storage.ObjectExistsRequest{
        Bucket: extractBucket(req.DatasetPath),
        Key:    extractKey(req.DatasetPath),
    })
    if err != nil || !exists {
        return nil, errors.New("dataset not found")
    }
    
    // 3. Build AutoML-specific pipeline parameters
    pipelineParams := map[string]string{
        "dataset_path":    req.DatasetPath,
        "target_column":   req.TargetColumn,
        "problem_type":    req.ProblemType,
        "max_trials":      strconv.Itoa(req.MaxTrials),
        "training_budget": req.TrainingBudget,
        // ... more AutoML-specific params
    }
    
    // 4. Delegate to autox-core pipelines service
    run, err := s.pipelinesService.CreatePipelineRun(
        ctx,
        httpClient,
        req.DSPA.BaseURL,
        pipelines.CreatePipelineRunRequest{
            Identity:     req.Identity,
            Namespace:    req.Namespace,
            PipelineID:   req.DiscoveredPipeline.PipelineID,
            VersionID:    req.DiscoveredPipeline.VersionID,
            Name:         generateAutoMLRunName(req),
            Parameters:   pipelineParams,
        },
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create AutoML pipeline run: %w", err)
    }
    
    s.logger.Info("Created AutoML pipeline run",
        "run_id", run.ID,
        "namespace", req.Namespace,
        "problem_type", req.ProblemType,
    )
    
    return run, nil
}

func (s *AutoMLPipelineService) validateAutoMLParams(req CreateAutoMLPipelineRunRequest) error {
    // AutoML-specific validation logic
    if req.ProblemType != "classification" && req.ProblemType != "regression" && req.ProblemType != "time-series" {
        return &common.ValidationError{
            Field:   "problem_type",
            Message: "must be classification, regression, or time-series",
        }
    }
    if req.MaxTrials < 1 || req.MaxTrials > 100 {
        return &common.ValidationError{
            Field:   "max_trials",
            Message: "must be between 1 and 100",
        }
    }
    return nil
}
```

### Handler Pattern (Thin, HTTP Concerns Only)

```go
// packages/automl/bff/internal/api/pipeline_runs_handler.go

// LIST handler - calls autox-core directly (no AutoML-specific logic)
func (app *App) PipelineRunsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    ctx := r.Context()
    
    // Extract from context (added by middleware)
    identity := ctx.Value(constants.RequestIdentityKey).(*k8s.RequestIdentity)
    namespace := ctx.Value(constants.NamespaceKey).(string)
    dspa := ctx.Value(constants.DSPAKey).(*pipelines.DSPA)
    discoveredPipeline := ctx.Value(constants.DiscoveredPipelineKey).(*pipelines.DiscoveredPipeline)
    
    // Parse query parameters (HTTP concern)
    pageSize := parseInt32(r.URL.Query().Get("pageSize"), 10)
    pageToken := r.URL.Query().Get("pageToken")
    
    // Simple list - call autox-core directly (no domain logic)
    runs, nextToken, err := app.pipelinesService.ListPipelineRuns(
        ctx, 
        app.httpClient,
        dspa.BaseURL, 
        pipelines.ListPipelineRunsRequest{
            Identity:          identity,
            Namespace:         namespace,
            PipelineVersionID: discoveredPipeline.VersionID,
            PageSize:          pageSize,
            PageToken:         pageToken,
        },
    )
    if err != nil {
        app.handleError(w, err)
        return
    }
    
    // Map to API response
    response := PipelineRunsResponse{
        Runs:          mapPipelineRunsToAPI(runs),
        NextPageToken: nextToken,
    }
    
    app.WriteJSON(w, http.StatusOK, response, nil)
}

// CREATE handler - calls AutoML service (has domain logic)
func (app *App) CreatePipelineRunHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    ctx := r.Context()
    
    // Extract from context
    identity := ctx.Value(constants.RequestIdentityKey).(*k8s.RequestIdentity)
    namespace := ctx.Value(constants.NamespaceKey).(string)
    dspa := ctx.Value(constants.DSPAKey).(*pipelines.DSPA)
    discoveredPipeline := ctx.Value(constants.DiscoveredPipelineKey).(*pipelines.DiscoveredPipeline)
    
    // Parse request body (HTTP concern)
    var apiReq CreateAutoMLRunAPIRequest
    if err := json.NewDecoder(r.Body).Decode(&apiReq); err != nil {
        app.badRequestResponse(w, r, "invalid request body")
        return
    }
    
    // Call AutoML domain service (has AutoML-specific logic)
    run, err := app.autoMLPipelineService.CreateAutoMLPipelineRun(
        ctx,
        app.httpClient,
        services.CreateAutoMLPipelineRunRequest{
            Identity:           identity,
            Namespace:          namespace,
            DSPA:               dspa,
            DiscoveredPipeline: discoveredPipeline,
            DatasetPath:        apiReq.DatasetPath,
            TargetColumn:       apiReq.TargetColumn,
            ProblemType:        apiReq.ProblemType,
            MaxTrials:          apiReq.MaxTrials,
            TrainingBudget:     apiReq.TrainingBudget,
            S3Credentials:      apiReq.S3Credentials,
        },
    )
    if err != nil {
        app.handleError(w, err)
        return
    }
    
    // Map to API response
    response := mapPipelineRunToAPI(run)
    app.WriteJSON(w, http.StatusCreated, response, nil)
}
```

### Middleware Pattern (Thin, Delegates to Services)

```go
// packages/automl/bff/internal/api/middleware.go

func (app *App) AttachNamespace(next httprouter.Handle) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        namespace := r.URL.Query().Get("namespace")
        
        // Validate (could use autox-core validator)
        if err := k8s.ValidateNamespaceName(namespace); err != nil {
            app.badRequestResponse(w, r, err.Error())
            return
        }
        
        // Store in context for handlers
        ctx := context.WithValue(r.Context(), constants.NamespaceKey, namespace)
        next(w, r.WithContext(ctx), ps)
    }
}

func (app *App) RequireAccessToNamespace(next httprouter.Handle) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        ctx := r.Context()
        identity := ctx.Value(constants.RequestIdentityKey).(*k8s.RequestIdentity)
        namespace := ctx.Value(constants.NamespaceKey).(string)
        
        // Call autox-core service for RBAC check
        allowed, err := app.kubernetesService.CanAccessNamespace(ctx, identity, namespace)
        if err != nil {
            app.serverErrorResponse(w, r, err)
            return
        }
        if !allowed {
            app.forbiddenResponse(w, r)
            return
        }
        
        next(w, r.WithContext(ctx), ps)
    }
}

func (app *App) AttachDSPA(next httprouter.Handle) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        ctx := r.Context()
        identity := ctx.Value(constants.RequestIdentityKey).(*k8s.RequestIdentity)
        namespace := ctx.Value(constants.NamespaceKey).(string)
        
        // Call autox-core service for DSPA discovery
        dspa, err := app.pipelinesService.DiscoverDSPA(ctx, pipelines.DiscoverDSPARequest{
            Identity:  identity,
            Namespace: namespace,
        })
        if err != nil {
            if errors.Is(err, pipelines.ErrNoDSPAFound) {
                app.notFoundResponse(w, r)
                return
            }
            app.serverErrorResponse(w, r, err)
            return
        }
        if dspa == nil {
            // Exists but not ready
            app.serviceUnavailableResponse(w, r, "Pipeline server is not ready")
            return
        }
        
        // Store in context for handlers
        ctx = context.WithValue(ctx, constants.DSPAKey, dspa)
        next(w, r.WithContext(ctx), ps)
    }
}
```

**Key Points**:
- Middleware is thin (extract, validate, check, store in context)
- Business logic split between domain services (AutoML/AutoRAG) and autox-core services
- Services called with explicit request structs
- Middleware stores discovered entities (identity, namespace, DSPA) in context for handlers
- Handlers receive entities from context, call services with explicit parameters
- **Client objects (http.Client, S3Client) are NEVER stored in context** - passed as function parameters from App
- **Entity objects (RequestIdentity, Namespace, DSPA, DiscoveredPipeline) ARE stored in context** - flow through middleware chain

### Decision Tree: When to Call `autox-core` Directly vs Domain Service

```
Does the operation need ANY of the following?
├─ Domain-specific validation (e.g., AutoML problem_type enum)
├─ Domain-specific parameter building (e.g., AutoML pipeline params)
├─ Multi-step orchestration (e.g., validate dataset THEN create run)
├─ Integration with domain-specific services (e.g., model registry)
└─ Domain-specific business rules (e.g., AutoML max_trials limits)

    YES → Create/Use Domain Service
    NO  → Call `autox-core` Directly from Handler
```

**Concrete Examples**:

| Operation | Call Pattern | Reason |
|-----------|--------------|--------|
| **List pipeline runs** | Handler → `autox-core` | No domain logic - just pagination and filtering |
| **Get pipeline run by ID** | Handler → `autox-core` | No domain logic - just fetch and return |
| **Terminate pipeline run** | Handler → `autox-core` | No domain logic - simple operation |
| **List namespaces** | Handler → `autox-core` | No domain logic - pure RBAC filtering |
| **Get secret** | Handler → `autox-core` | No domain logic - simple retrieval |
| **List S3 objects** | Handler → `autox-core` | No domain logic - basic listing |
| | | |
| **Create AutoML pipeline run** | Handler → AutoML Service → `autox-core` | ✅ Domain validation (problem_type, max_trials)<br>✅ Parameter building (AutoML-specific params)<br>✅ Multi-step (validate dataset exists, then create) |
| **Upload dataset** | Handler → AutoML Service → `autox-core` | ✅ Domain validation (CSV schema)<br>✅ Domain metadata (AutoML dataset tags) |
| **Register model** | Handler → AutoML Service → `autox-core` + ModelRegistry | ✅ Multi-service orchestration<br>✅ Domain validation (model artifact format) |
| **Create AutoRAG run** | Handler → AutoRAG Service → `autox-core` | ✅ Domain validation (vector store config)<br>✅ Parameter building (RAG-specific params) |

**Rule of Thumb**: 
- **Simple CRUD with no domain logic** → Handler calls autox-core directly
- **Domain validation, param building, or orchestration** → Handler calls domain service

---

## Migration Strategy

### Phase 1: Build Complete `autox-core` Architecture

**Goal**: Establish the ideal architecture before migrating consumers

#### 1.1 Core Infrastructure

**Files to Create**:
```
packages/autox-core/bff/
├── common/
│   ├── errors/
│   │   ├── errors.go          # Base error types (NotFoundError, ForbiddenError, etc)
│   │   └── errors_test.go
│   ├── validation/
│   │   ├── validators.go      # Generic validators (DNS-1123, URL, SSRF)
│   │   └── validators_test.go
│   └── cache/
│       ├── lru.go             # LRU cache with TTL
│       └── lru_test.go
```

**Custom Error Types**:
```go
// common/errors/errors.go
type NotFoundError struct {
    Resource string
    Message  string
}

type ForbiddenError struct {
    Resource string
    Message  string
}

type ValidationError struct {
    Field   string
    Message string
}

// HTTP status mapping helper
func StatusCode(err error) int {
    switch err.(type) {
    case *NotFoundError:
        return http.StatusNotFound
    case *ForbiddenError:
        return http.StatusForbidden
    case *ValidationError:
        return http.StatusBadRequest
    default:
        return http.StatusInternalServerError
    }
}
```

#### 1.2 Kubernetes Integration

**Extract from**: 
- `packages/automl/bff/internal/integrations/kubernetes/`
- `packages/autorag/bff/internal/integrations/kubernetes/`

**Files to Create**:
```
packages/autox-core/bff/kubernetes/
├── models.go              # RequestIdentity, Namespace, Secret
├── client.go              # KubernetesClient interface
├── client_internal.go     # Internal k8s client (service account)
├── client_token.go        # Token-based k8s client (per-user)
├── client_shared.go       # Shared RBAC logic
├── factory.go             # Unified factory (selects static vs token)
├── service.go             # KubernetesService implementation
├── validators.go          # Namespace validation, RBAC checks
├── errors.go              # K8s-specific errors
├── portforward.go         # Port-forward manager (dev mode)
└── kubernetes_test.go
```

**Key Extractions**:
- `KubernetesClientFactory` interface + implementations
- `RequestIdentity` model (UserID, Groups, Token)
- `SubjectAccessReview` (SAR) RBAC checking logic
- Namespace listing with permission filtering
- Secret retrieval with validation
- Unified factory that switches between static/token based on config
- Port-forward manager (SPDY port-forwarding for dev mode)

**New Service Interface**:
```go
type KubernetesService interface {
    ListNamespaces(ctx context.Context, req ListNamespacesRequest) ([]Namespace, error)
    GetSecret(ctx context.Context, req GetSecretRequest) (*Secret, error)
    ListSecrets(ctx context.Context, req ListSecretsRequest) ([]Secret, error)
    IsClusterAdmin(ctx context.Context, identity *RequestIdentity) (bool, error)
    CanAccessNamespace(ctx context.Context, identity *RequestIdentity, namespace string) (bool, error)
}
```

#### 1.3 Pipelines Integration

**Extract from**:
- `packages/automl/bff/internal/integrations/pipelineserver/`
- `packages/automl/bff/internal/repositories/pipeline.go`
- `packages/automl/bff/internal/api/middleware.go` (DSPA discovery)

**Files to Create**:
```
packages/autox-core/bff/pipelines/
├── models.go              # PipelineRun, Pipeline, DSPA, DiscoveredPipeline
├── client.go              # HTTP client for KFP API (list, create, terminate runs)
├── service.go             # PipelinesService implementation
├── discovery.go           # DSPA discovery, pipeline discovery with caching
├── cache.go               # Pipeline cache implementation (LRU, 5min TTL)
├── validators.go          # Request validation
├── errors.go              # Pipeline-specific errors
└── pipelines_test.go
```

**Key Extractions**:
- DSPA discovery logic (list CRDs, check readiness, extract baseURL)
- Pipeline discovery by name prefix (with LRU cache)
- KFP HTTP client (list runs, create run, terminate, retry)
- Cache management (5min TTL, 1000 entry LRU)
- DSPA readiness checking (`isDSPAReady()`)

**Service Interface**:
```go
type PipelinesService interface {
    // DSPA operations
    DiscoverDSPA(ctx context.Context, req DiscoverDSPARequest) (*DSPA, error)
    
    // Pipeline discovery (with caching)
    DiscoverPipelines(ctx context.Context, req DiscoverPipelinesRequest) (map[string]*DiscoveredPipeline, error)
    
    // Pipeline run operations (requires http.Client + baseURL)
    ListPipelineRuns(ctx context.Context, client *http.Client, baseURL string, req ListPipelineRunsRequest) ([]PipelineRun, string, error)
    GetPipelineRun(ctx context.Context, client *http.Client, baseURL string, req GetPipelineRunRequest) (*PipelineRun, error)
    CreatePipelineRun(ctx context.Context, client *http.Client, baseURL string, req CreatePipelineRunRequest) (*PipelineRun, error)
    TerminatePipelineRun(ctx context.Context, client *http.Client, baseURL string, req TerminatePipelineRunRequest) error
    
    // Pipeline operations
    ListPipelines(ctx context.Context, client *http.Client, baseURL string, req ListPipelinesRequest) ([]Pipeline, error)
    CreatePipeline(ctx context.Context, client *http.Client, baseURL string, req CreatePipelineRequest) (*Pipeline, error)
}
```

**Discovery Request Types**:
```go
type DiscoverDSPARequest struct {
    Identity       *RequestIdentity
    Namespace      string
    K8sClient      KubernetesClient  // Passed explicitly for CRD queries
}

type DiscoverPipelinesRequest struct {
    Identity       *RequestIdentity
    Namespace      string
    BaseURL        string  // DSPA base URL
    NamePrefixes   map[string]string  // Pipeline type → name prefix
    HTTPClient     *http.Client
}
```

#### 1.4 Storage Integration

**Extract from**:
- `packages/automl/bff/internal/integrations/s3/`
- `packages/autorag/bff/internal/integrations/s3/`

**Files to Create**:
```
packages/autox-core/bff/storage/
├── models.go              # S3Object, Bucket, CSVSchema
├── client.go              # S3Client interface
├── service.go             # StorageService implementation
├── validators.go          # Path validation, SSRF protection
├── errors.go              # Storage-specific errors
└── storage_test.go
```

**Key Extractions**:
- S3 client interface (GetObject, PutObject, ListObjects)
- SSRF validation (block loopback, allow private ranges)
- Endpoint validation (MinIO, AWS S3)
- CSV schema detection

**Service Interface**:
```go
type StorageService interface {
    GetObject(ctx context.Context, client S3Client, req GetObjectRequest) (*S3Object, error)
    UploadObject(ctx context.Context, client S3Client, req UploadObjectRequest) error
    ListObjects(ctx context.Context, client S3Client, req ListObjectsRequest) ([]S3Object, error)
    GetCSVSchema(ctx context.Context, client S3Client, req GetCSVSchemaRequest) (*CSVSchema, error)
}
```

### Phase 2: Migrate AutoML

**Goal**: Refactor AutoML to use autox-core services

#### 2.1 Update App Initialization

**Changes in** `packages/automl/bff/internal/api/app.go`:
```go
// Before: Multiple factory fields
type App struct {
    kubernetesClientFactory     k8s.KubernetesClientFactory
    pipelineServerClientFactory ps.PipelineServerClientFactory
    s3ClientFactory             s3.S3ClientFactory
    repositories                *repositories.Repositories
}

// After: `autox-core` services
type App struct {
    config            config.EnvConfig
    logger            *slog.Logger
    httpClient        *http.Client  // Shared for pipelines
    
    // `autox-core` services
    kubernetesService kubernetes.KubernetesService
    pipelinesService  pipelines.PipelinesService
    storageService    storage.StorageService
    
    // AutoML-specific
    modelRegistryService *modelregistry.Service
    repositories         *automlRepositories  // AutoML-specific repos
}
```

#### 2.2 Refactor Middleware

**Changes in** `packages/automl/bff/internal/api/middleware.go`:

**Before** (passes client objects via context ❌):
```go
func (app *App) AttachPipelineServerClient(next) httprouter.Handle {
    // ... discover DSPA ...
    // ❌ BAD: Creates and stores client in context
    client := app.pipelineServerClientFactory.CreateClient(baseURL, ...)
    ctx = context.WithValue(ctx, PipelineServerClientKey, client)
    next(w, r.WithContext(ctx), ps)
}
```

**After** (stores entity, passes client as parameter ✅):
```go
func (app *App) AttachDSPA(next) httprouter.Handle {
    ctx := r.Context()
    identity := ctx.Value(RequestIdentityKey).(*kubernetes.RequestIdentity)
    namespace := ctx.Value(NamespaceKey).(string)
    
    // Call autox-core service (service handles discovery logic)
    dspa, err := app.pipelinesService.DiscoverDSPA(ctx, pipelines.DiscoverDSPARequest{
        Identity:  identity,
        Namespace: namespace,
    })
    
    // ✅ GOOD: Store entity object in context (not client)
    ctx = context.WithValue(ctx, DSPAKey, dspa)
    next(w, r.WithContext(ctx), ps)
}
```

#### 2.3 Refactor Handlers

**Changes in** `packages/automl/bff/internal/api/pipeline_runs_handler.go`:

**Before** (repository pattern with client from context):
```go
func (app *App) PipelineRunsHandler(w, r, _) {
    client := ctx.Value(PipelineServerClientKey).(PipelineServerClientInterface)
    runs, err := app.repositories.PipelineRuns.GetPipelineRuns(client, ctx, ...)
}
```

**After** (service pattern with explicit parameters):
```go
func (app *App) PipelineRunsHandler(w, r, _) {
    // Extract from context
    identity := ctx.Value(RequestIdentityKey).(*kubernetes.RequestIdentity)
    dspa := ctx.Value(DSPAKey).(*pipelines.DSPA)
    discoveredPipeline := ctx.Value(DiscoveredPipelineKey).(*pipelines.DiscoveredPipeline)
    
    // Call autox-core service
    runs, nextToken, err := app.pipelinesService.ListPipelineRuns(
        ctx, 
        app.httpClient,  // Shared client
        dspa.BaseURL,
        pipelines.ListPipelineRunsRequest{
            Identity:          identity,
            Namespace:         namespace,
            PipelineVersionID: discoveredPipeline.VersionID,
            PageSize:          pageSize,
            PageToken:         pageToken,
        },
    )
    
    // Map to API response
    response := mapToAPIResponse(runs, nextToken)
    app.WriteJSON(w, http.StatusOK, response, nil)
}
```

#### 2.4 Delete Duplicate Code

**New Directory Structure**:
```
packages/automl/bff/internal/
├── api/                           # HTTP layer
│   ├── handlers.go               # Thin handlers (HTTP concerns)
│   ├── middleware.go             # Thin middleware (delegates to services)
│   └── errors.go                 # Error mapping helpers
├── services/                      # AutoML domain services
│   ├── automl_pipeline_service.go  # Composes autox-core pipelines service
│   ├── dataset_service.go          # Composes autox-core storage service
│   └── model_service.go            # Uses modelregistry client + autox-core
├── integrations/                  # AutoML-specific integrations
│   └── modelregistry/             # Model Registry client (domain-specific)
├── models/                        # AutoML API request/response types
└── config/                        # Configuration
```

**Files to Delete** (now in autox-core):
- `internal/integrations/kubernetes/` (all of it)
- `internal/integrations/pipelineserver/` (all of it)
- `internal/integrations/s3/` (all of it)
- `internal/repositories/pipeline.go` (pipeline discovery)
- `internal/repositories/pipeline_runs.go` (basic CRUD)
- `internal/repositories/namespace.go` (basic CRUD)
- `internal/repositories/secret.go` (basic CRUD)
- Parts of `internal/api/middleware.go` (DSPA discovery, pipeline discovery)

**Files to Keep/Create** (automl-specific):
- `internal/integrations/modelregistry/` (domain-specific)
- `internal/services/` (domain services that compose autox-core)
- `internal/models/` (AutoML API models)
- `internal/api/` (handlers, middleware, error helpers)

### Phase 3: Migrate AutoRAG

**Goal**: Same refactoring for AutoRAG

Apply same pattern as AutoML:
1. Update App to use autox-core services
2. Refactor middleware to call services (no context.WithValue for clients)
3. Refactor handlers to call services with explicit requests
4. Delete duplicate code

**AutoRAG-specific to keep**:
- `internal/integrations/llamastack/` (domain-specific)
- AutoRAG pipeline discovery logic (different name prefixes)

---

## Validation Helpers

### Common Validators (in autox-core)

**Generic** (in `common/validation/`):
```go
// DNS-1123 subdomain validation (namespaces)
func ValidateDNS1123Subdomain(value string) error

// DNS-1123 label validation (resource names)
func ValidateDNS1123Label(value string) error

// URL validation with SSRF protection
func ValidateURL(url string, allowPrivate bool) error

// IP validation (block loopback, link-local)
func ValidateIPAddress(ip string) error
```

**Domain-Specific** (in feature packages):

`kubernetes/validators.go`:
```go
func ValidateNamespaceName(name string) error
func ValidateSecretName(name string) error
func ValidateListNamespacesRequest(req *ListNamespacesRequest) error
```

`pipelines/validators.go`:
```go
func ValidateListPipelineRunsRequest(req *ListPipelineRunsRequest) error
func ValidateCreatePipelineRunRequest(req *CreatePipelineRunRequest) error
func ValidatePipelineNamePrefix(prefix string) error
```

`storage/validators.go`:
```go
func ValidateS3Path(path string) error
func ValidateBucketName(name string) error
func ValidateUploadObjectRequest(req *UploadObjectRequest) error
```

---

## Critical Files to Extract

### From AutoML BFF

| Category | File | Priority | Extract To |
|----------|------|----------|------------|
| **K8s Client** | `integrations/kubernetes/factory.go` | P0 | `autox-core/bff/kubernetes/factory.go` |
| **K8s Client** | `integrations/kubernetes/client.go` | P0 | `autox-core/bff/kubernetes/client.go` |
| **K8s Client** | `integrations/kubernetes/internal_k8s_client.go` | P0 | `autox-core/bff/kubernetes/client_internal.go` |
| **K8s Client** | `integrations/kubernetes/token_k8s_client.go` | P0 | `autox-core/bff/kubernetes/client_token.go` |
| **K8s Client** | `integrations/kubernetes/shared_k8s_client.go` | P0 | `autox-core/bff/kubernetes/client_shared.go` |
| **K8s Client** | `integrations/kubernetes/portforward.go` | P0 | `autox-core/bff/kubernetes/portforward.go` |
| **Pipelines** | `integrations/pipelineserver/client.go` | P0 | `autox-core/bff/pipelines/client.go` |
| **Pipelines** | `integrations/pipelineserver/client_factory.go` | P0 | `autox-core/bff/pipelines/` (merge into service) |
| **Pipelines** | `repositories/pipeline.go` | P0 | `autox-core/bff/pipelines/discovery.go` |
| **Pipelines** | `api/middleware.go` (DSPA discovery) | P0 | `autox-core/bff/pipelines/discovery.go` |
| **Storage** | `integrations/s3/client.go` | P1 | `autox-core/bff/storage/client.go` |
| **Validation** | Various DNS-1123 validators | P1 | `autox-core/bff/common/validation/` |

---

## Implementation Tasks

### Phase 1: Foundation (Week 1)

1. ✅ **Create autox-core package structure**
   - Set up `packages/autox-core/bff/` directory
   - Create `go.mod` with Go 1.24.3
   - Set up feature-based directories (kubernetes/, pipelines/, storage/, common/)

2. **Build common infrastructure**
   - Create custom error types (`common/errors/`)
   - Create LRU cache with TTL (`common/cache/`)
   - Create generic validators (`common/validation/`)
   - Write unit tests for common utilities

3. **Extract Kubernetes integration**
   - Extract `RequestIdentity` model
   - Extract unified factory pattern
   - Extract internal/token k8s clients
   - Extract RBAC logic (SubjectAccessReview)
   - Create `KubernetesService` interface and implementation
   - Create request/response types
   - Write unit tests

### Phase 2: Pipelines Integration (Week 2)

4. **Extract Pipelines integration**
   - Extract DSPA models and discovery logic
   - Extract pipeline discovery with caching
   - Extract KFP HTTP client
   - Create `PipelinesService` interface and implementation
   - Create request/response types
   - Write unit tests

5. **Extract Storage integration**
   - Extract S3 client interface
   - Extract SSRF validation
   - Create `StorageService` interface
   - Create request/response types
   - Write unit tests

### Phase 3: AutoML Migration (Week 3)

6. **Refactor AutoML App**
   - Update `App` struct to use autox-core services
   - Create shared `http.Client` for pipelines
   - Remove old factory fields
   - Update initialization logic

7. **Create AutoML Domain Services**
   - Create `internal/services/` directory
   - Create `AutoMLPipelineService` (composes autox-core pipelines + storage)
   - Create `DatasetService` (composes autox-core storage service)
   - Create `ModelService` (uses modelregistry + autox-core)
   - Define AutoML-specific request types
   - Implement domain validation and orchestration logic

8. **Refactor AutoML Middleware**
   - Update `AttachDSPA` to call autox-core `DiscoverDSPA`
   - Update `AttachDiscoveredPipeline` to call autox-core
   - Remove client creation logic (store DSPA/entities instead)
   - Keep thin middleware that calls services

9. **Refactor AutoML Handlers**
   - **Simple operations with NO domain logic** → call autox-core directly
     - Examples: List runs, Get run, Terminate run, List namespaces, Get secret
   - **Complex operations WITH domain logic** → call AutoML domain service
     - Examples: Create AutoML run, Upload dataset, Register model
   - Pass explicit request structs
   - Remove repository layer (replaced by autox-core + domain services)
   - Map domain models to API responses

10. **Cleanup AutoML**
    - Delete duplicate integration code (k8s, pipelines, s3)
    - Delete old repository layer
    - Keep domain-specific code (modelregistry)
    - Run tests, fix issues

### Phase 4: AutoRAG Migration (Week 4)

1.  **Create AutoRAG Domain Services**
    - Create `internal/services/` directory
    - Create `AutoRAGPipelineService` (composes autox-core pipelines + storage)
    - Create `VectorStoreService` (composes llamastack + autox-core storage)
    - Define AutoRAG-specific request types
    - Implement domain validation and orchestration logic

2.  **Migrate AutoRAG** (same steps as AutoML)
    - Update App to use autox-core services
    - Refactor middleware (thin, delegates to services)
    - Refactor handlers (simple ops → autox-core, complex ops → AutoRAG service)
    - Delete duplicate integration code
    - Keep llamastack integration

3.  **Final Cleanup**
    - Update documentation
    - Add autox-core API documentation
    - Update consumer integration guides
    - Run full test suite

---

## Verification Checklist

After implementation, verify:

### `autox-core` BFF

- [ ] Feature-based directory structure (kubernetes/, pipelines/, storage/)
- [ ] Each feature has: models, client, service, validators, errors
- [ ] Services use explicit request/response types
- [ ] Services return domain models (not JSON response structs)
- [ ] No HTTP framework code (no handlers, no middleware)
- [ ] Custom error types with HTTP status mapping
- [ ] LRU cache implementation with TTL
- [ ] Generic validators (DNS-1123, URL, SSRF)
- [ ] Domain-specific validators per feature
- [ ] Comprehensive unit tests (>80% coverage)

### AutoML/AutoRAG BFF

**App Structure**:
- [ ] App struct owns `http.Client` and autox-core services
- [ ] App struct owns domain-specific services (AutoMLPipelineService, ModelRegistryService)
- [ ] No client factories (use autox-core services)

**Domain Service Layer** (OPTIONAL, only when domain logic exists):
- [ ] `internal/services/` directory created
- [ ] Domain services only created when there's actual domain-specific logic
- [ ] Domain services compose autox-core services (explicit dependencies)
- [ ] Domain services implement domain-specific validation
- [ ] Domain services build domain-specific parameters
- [ ] Domain services orchestrate multi-step operations
- [ ] Domain services receive explicit parameters (no context.WithValue for clients)
- [ ] Domain services return domain models (not JSON responses)

**Handler Layer**:
- [ ] Handlers are thin (<50 lines, HTTP concerns only)
- [ ] **Simple operations with NO domain logic** → call autox-core directly
  - Examples: List runs, Get run, Terminate run, List namespaces, Get secret
- [ ] **Complex operations WITH domain logic** → call domain services
  - Examples: Create AutoML run, Upload dataset, Register model
- [ ] Handlers extract entities from context (identity, namespace, DSPA)
- [ ] Handlers pass App-owned clients to services (http.Client, S3Client)
- [ ] Domain models mapped to API responses in handlers

**Middleware Layer**:
- [ ] Middleware is thin (<30 lines, delegates to services)
- [ ] Middleware calls autox-core services for discovery/validation
- [ ] **No client objects in context** (http.Client, S3Client passed as function parameters)
- [ ] **Entity objects stored in context** (RequestIdentity, Namespace, DSPA, DiscoveredPipeline)

**Code Cleanup**:
- [ ] Duplicate code removed (k8s, pipelines, storage integrations)
- [ ] Old repository layer removed (replaced by autox-core + domain services)
- [ ] Domain-specific code kept (modelregistry for AutoML, llamastack for AutoRAG)
- [ ] Tests pass
- [ ] No regressions in API behavior

---

## Success Metrics

1. **Code Reduction**: 60-70% reduction in automl/autorag BFF code (by removing duplicates)
2. **Single http.Client**: One shared client in App, not per-request
3. **Thin Handlers**: Handlers < 50 lines, delegate to services
4. **Thin Middleware**: Middleware < 30 lines, call services
5. **Service Coverage**: All common operations (k8s, pipelines, storage) in autox-core
6. **Test Coverage**: `autox-core` services >80% unit test coverage
7. **API Stability**: No breaking changes to automl/autorag REST APIs
8. **Performance**: No performance regression (caching improves performance)
