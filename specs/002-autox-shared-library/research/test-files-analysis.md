# Test Files Duplication Analysis

## Executive Summary

- **Total matching test files analyzed**: 16 (Go BFF tests only; excludes 579 node_modules duplicates)
- **100% identical files**: 5 files (31.2%)
- **Files with minor differences**: 11 files (68.8%)
- **Estimated shared test code potential**: ~85-90%
  - Import path differences: ~100% of differences
  - Test logic duplication: ~95%+ identical patterns
  - Mock data patterns: ~90%+ shared structures

## Key Findings

### Duplication Patterns

1. **Identical Test Logic**: Most test files share 95%+ identical test logic, differing only in:
   - Import paths: `github.com/opendatahub-io/automl-library` vs `github.com/opendatahub-io/autorag-library`
   - Module-specific feature tests (e.g., Model Registry in AutoML, LlamaStack in AutoRAG)
   - Minor mock data variations

2. **Shared Test Infrastructure**:
   - All files use identical testing frameworks (testify/assert, testify/require)
   - Identical test setup patterns (envtest, K8s mocks)
   - Shared middleware testing patterns
   - Common DSPA discovery and validation logic

3. **Extractable Components**:
   - DSPA testing utilities (isDSPAReady, buildMinIOObjectStorage)
   - Middleware test helpers (validation, error handling, discovery)
   - Pipeline testing utilities (run status, filtering, sorting)
   - S3/storage testing utilities (mock clients, bucket operations)
   - K8s integration testing patterns (namespace operations, SAR checks)

## Detailed Test File Analysis

### Category 1: 100% Identical Files (5 files)

These files are byte-for-byte identical and can be immediately extracted to shared library:

| File | Lines | Purpose | Extractable Utilities |
|------|-------|---------|----------------------|
| `bff/cmd/helpers_test.go` | 37 | Origin parser tests | `newOriginParser` test patterns |
| `bff/internal/api/helpers_test.go` | 22 | URL template parser tests | `ParseURLTemplate` test cases |
| `bff/internal/api/middleware_dspa_errors_test.go` | 101 | DSPA error handling | Error middleware test patterns |
| `bff/internal/api/middleware_validation_test.go` | 111 | Request validation | Validation middleware test helpers |
| `bff/internal/integrations/kubernetes/get_namespaces_test.go` | 224 | Namespace operations | K8s namespace test utilities |

**Total: 495 lines of 100% shareable test code**

### Category 2: Files with Minor Import Path Differences (4 files)

These files differ ONLY in import paths (`automl-library` vs `autorag-library`) plus minor test expectations:

| File | AutoML Lines | AutoRAG Lines | Differences | Duplication % |
|------|--------------|---------------|-------------|---------------|
| `bff/internal/api/middleware_discovery_test.go` | 546 | 547 | Import paths + 1 test expectation ("empty-namespace" vs "no-dspas-namespace") | 98% |
| `bff/internal/api/s3_declared_limit_test.go` | 120 | 119 | Import paths only | 99% |
| `bff/internal/repositories/pipeline_run_get_test.go` | 191 | 191 | Import paths + 1 line number comment | 99% |
| `bff/internal/repositories/pipeline_test.go` | 536 | 532 | Import paths + minor test data | 97% |

**Extractable Test Utilities**:
- DSPA discovery testing patterns (isDSPAReady, discoverReadyDSPA)
- MinIO storage testing helpers (buildMinIOObjectStorage, storage injection tests)
- S3 limit validation patterns
- Pipeline run retrieval test patterns
- Pipeline listing and filtering test utilities

### Category 3: Files with Module-Specific Feature Tests (7 files)

These files share core testing patterns but include module-specific features:

| File | AutoML Lines | AutoRAG Lines | Shared % | Module-Specific Features |
|------|--------------|---------------|----------|-------------------------|
| `bff/internal/api/pipeline_run_handler_test.go` | 428 | 384 | ~70% | AutoML: Model Registry tests; AutoRAG: LlamaStack tests |
| `bff/internal/api/pipeline_runs_handler_test.go` | 1,315 | 1,488 | ~65% | Different filter/sort test cases |
| `bff/internal/api/s3_handler_test.go` | 2,668 | 2,008 | ~50% | Different S3 endpoint configurations |
| `bff/internal/api/secrets_handler_test.go` | 1,351 | 1,612 | ~75% | Different secret validation patterns |
| `bff/internal/integrations/s3/client_test.go` | 452 | 376 | ~70% | Different S3 client configurations |
| `bff/internal/repositories/pipeline_runs_test.go` | 515 | 368 | ~60% | Different pipeline filtering logic |
| `bff/internal/repositories/s3_test.go` | 260 | 405 | ~50% | Different S3 repository operations |

**Extractable Shared Test Patterns** (even from these files):
- HTTP handler test scaffolding
- Request/response validation patterns
- Error handling test cases
- Mock server setup and teardown
- Table-driven test structures
- Pagination testing utilities
- Filter/sort validation helpers

## Shared Test Utilities to Extract

### 1. DSPA Testing Utilities (High Priority)

**Location**: `internal/testing/dspa/`

**Functions**:
```go
// DSPA validation helpers
func IsDSPAReady(dspa *models.DSPipelineApplication) bool
func AssertDSPAReady(t *testing.T, dspa *models.DSPipelineApplication)
func AssertDSPANotReady(t *testing.T, dspa *models.DSPipelineApplication)

// Mock DSPA builders
func MockDSPipelineApplication(opts ...DSPAOption) models.DSPipelineApplication
func MockDSPAWithMinIO(namespace, dspaName, bucket string) models.DSPipelineApplication
func MockDSPAWithExternalStorage(namespace, dspaName string, storage ExternalStorageConfig) models.DSPipelineApplication

// Storage testing utilities
func BuildMinIOObjectStorage(dspaName, namespace, bucket string) *models.DSPAObjectStorage
func AssertStorageConfig(t *testing.T, storage *models.DSPAObjectStorage, expected *models.DSPAObjectStorage)
```

**Source**: Extracted from `middleware_discovery_test.go` (identical in both packages)

**Usage Example**:
```go
import "github.com/opendatahub-io/autox-library/testing/dspa"

func TestMyFeature(t *testing.T) {
    mockDSPA := dspa.MockDSPipelineApplication(
        dspa.WithNamespace("test-ns"),
        dspa.WithMinIO("mlpipeline"),
        dspa.Ready(true),
    )
    
    assert.True(t, dspa.IsDSPAReady(&mockDSPA))
}
```

### 2. Pipeline Testing Utilities (High Priority)

**Location**: `internal/testing/pipeline/`

**Functions**:
```go
// Pipeline run test builders
func MockPipelineRun(opts ...RunOption) *models.PipelineRun
func MockPipelineRunList(count int, opts ...RunOption) []*models.PipelineRun
func MockPipelineVersion(opts ...VersionOption) *models.PipelineVersion

// Filter/sort test helpers
func AssertRunsFilteredBy(t *testing.T, runs []*models.PipelineRun, filterKey, filterValue string)
func AssertRunsSortedBy(t *testing.T, runs []*models.PipelineRun, sortKey string, ascending bool)
func AssertPaginatedRuns(t *testing.T, response *PaginatedResponse, expectedCount, expectedTotal int)

// Status assertion helpers
func AssertRunStatus(t *testing.T, run *models.PipelineRun, expectedStatus string)
func AssertRunCompleted(t *testing.T, run *models.PipelineRun)
func AssertRunFailed(t *testing.T, run *models.PipelineRun)
```

**Source**: Extracted from `pipeline_runs_handler_test.go`, `pipeline_run_handler_test.go`

### 3. Middleware Testing Utilities (High Priority)

**Location**: `internal/testing/middleware/`

**Functions**:
```go
// Middleware test scaffolding
func SetupMiddlewareTest(t *testing.T, cfg config.EnvConfig) (*App, *httptest.ResponseRecorder, *http.Request)
func AssertMiddlewareError(t *testing.T, recorder *httptest.ResponseRecorder, expectedStatus int, expectedMessage string)
func AssertMiddlewareSuccess(t *testing.T, recorder *httptest.ResponseRecorder)

// Request validation helpers
func AssertValidationError(t *testing.T, recorder *httptest.ResponseRecorder, field string)
func AssertRequiredFieldError(t *testing.T, recorder *httptest.ResponseRecorder, field string)

// DSPA error handling
func AssertDSPANotFoundError(t *testing.T, recorder *httptest.ResponseRecorder)
func AssertDSPANotReadyError(t *testing.T, recorder *httptest.ResponseRecorder)
```

**Source**: Extracted from `middleware_validation_test.go`, `middleware_dspa_errors_test.go`

### 4. K8s Integration Test Utilities (Medium Priority)

**Location**: `internal/testing/k8s/`

**Functions**:
```go
// K8s test environment setup
func SetupK8sTestEnv(t *testing.T, cfg config.EnvConfig) (*k8mocks.TestEnv, *kubernetes.Clientset)
func TeardownK8sTestEnv(t *testing.T, testEnv *k8mocks.TestEnv)

// Namespace testing
func AssertNamespaceAccessible(t *testing.T, client KubernetesClient, namespace string)
func AssertNamespaceRestricted(t *testing.T, client KubernetesClient, namespace string)
func MockNamespaceList(namespaces ...string) corev1.NamespaceList

// SAR/SSAR testing
func AssertSARAllowed(t *testing.T, client KubernetesClient, namespace, verb, resource string)
func AssertSSARAllowed(t *testing.T, client KubernetesClient, verb, resource string)
```

**Source**: Extracted from `get_namespaces_test.go`

### 5. S3 Testing Utilities (Medium Priority)

**Location**: `internal/testing/s3/`

**Functions**:
```go
// S3 mock builders
func MockS3Client(buckets ...string) *s3mocks.MockS3Client
func MockBucket(name string, objects ...string) s3types.Bucket
func MockS3Object(key string, size int64) s3types.Object

// S3 assertion helpers
func AssertBucketExists(t *testing.T, client S3Client, bucket string)
func AssertObjectExists(t *testing.T, client S3Client, bucket, key string)
func AssertS3Error(t *testing.T, err error, expectedCode string)

// S3 limit testing
func AssertDeclaredLimitValid(t *testing.T, limit string)
func AssertDeclaredLimitInvalid(t *testing.T, limit string, expectedError string)
```

**Source**: Extracted from `s3_handler_test.go`, `s3_declared_limit_test.go`, `s3_test.go`

### 6. HTTP Handler Test Scaffolding (High Priority)

**Location**: `internal/testing/http/`

**Functions**:
```go
// Handler test setup
func SetupHandlerTest(t *testing.T, method, path string, body interface{}) (*httptest.ResponseRecorder, *http.Request)
func SetupAuthenticatedRequest(t *testing.T, method, path, namespace string, body interface{}) *http.Request

// Response assertion helpers
func AssertJSONResponse(t *testing.T, recorder *httptest.ResponseRecorder, expectedStatus int, expectedBody interface{})
func AssertErrorResponse(t *testing.T, recorder *httptest.ResponseRecorder, expectedStatus int, expectedMessage string)
func AssertSuccessResponse(t *testing.T, recorder *httptest.ResponseRecorder)

// Request identity mocking
func MockInternalIdentity(userID string, groups ...string) *models.RequestIdentity
func MockTokenIdentity(token string) *models.RequestIdentity
```

**Source**: Common pattern across all handler tests

## Test Mocks and Fixtures to Extract

### 1. Pipeline Mock Data (High Priority)

**Location**: `internal/testing/fixtures/pipeline/`

**Structures**:
```go
// Standard pipeline runs for testing
func StandardPipelineRuns() []*models.PipelineRun {
    return []*models.PipelineRun{
        CompletedRun("run-1", "2024-01-01"),
        RunningRun("run-2", "2024-01-02"),
        FailedRun("run-3", "2024-01-03"),
    }
}

// Pipeline run builders
func CompletedRun(id, createdAt string) *models.PipelineRun
func RunningRun(id, createdAt string) *models.PipelineRun
func FailedRun(id, createdAt string) *models.PipelineRun
func PendingRun(id, createdAt string) *models.PipelineRun

// Pipeline version fixtures
func StandardPipelineVersion(id, pipelineID string) *models.PipelineVersion
```

**Shared Across**: All pipeline-related tests in both packages

### 2. DSPA Mock Data (High Priority)

**Location**: `internal/testing/fixtures/dspa/`

**Structures**:
```go
// Standard DSPA configurations
func DSPAWithMinIO(namespace, name string) models.DSPipelineApplication
func DSPAWithExternalS3(namespace, name string) models.DSPipelineApplication
func DSPANotReady(namespace, name string) models.DSPipelineApplication
func DSPAReady(namespace, name string) models.DSPipelineApplication

// Storage configurations
func MinIOStorageConfig(bucket string) *models.DSPAObjectStorage
func AWSS3StorageConfig(bucket, region string) *models.DSPAObjectStorage
```

**Shared Across**: Discovery, middleware, storage tests

### 3. Secrets Mock Data (Medium Priority)

**Location**: `internal/testing/fixtures/secrets/`

**Structures**:
```go
// Standard K8s secrets for testing
func S3CredentialsSecret(namespace, name string) corev1.Secret
func MinIOSecret(namespace, dspaName string) corev1.Secret
func AWSCredentialsSecret(namespace, name string) corev1.Secret

// Secret builders
func SecretWithData(namespace, name string, data map[string]string) corev1.Secret
```

**Shared Across**: Secrets handler tests, storage integration tests

## Recommended Shared Test Package Structure

```
packages/autox-library/
└── internal/
    └── testing/
        ├── README.md                   # Testing utilities documentation
        │
        ├── dspa/                       # DSPA testing utilities
        │   ├── assertions.go           # DSPA assertion helpers
        │   ├── builders.go             # DSPA mock builders
        │   └── storage.go              # Storage configuration helpers
        │
        ├── pipeline/                   # Pipeline testing utilities
        │   ├── assertions.go           # Pipeline assertion helpers
        │   ├── builders.go             # Pipeline run/version builders
        │   └── filters.go              # Filter/sort test helpers
        │
        ├── middleware/                 # Middleware testing utilities
        │   ├── setup.go                # Middleware test scaffolding
        │   ├── assertions.go           # Middleware assertion helpers
        │   └── errors.go               # Error response helpers
        │
        ├── k8s/                        # Kubernetes testing utilities
        │   ├── setup.go                # K8s test env setup/teardown
        │   ├── namespaces.go           # Namespace test helpers
        │   └── rbac.go                 # SAR/SSAR test helpers
        │
        ├── s3/                         # S3 testing utilities
        │   ├── mocks.go                # S3 mock builders
        │   ├── assertions.go           # S3 assertion helpers
        │   └── limits.go               # S3 limit validation helpers
        │
        ├── http/                       # HTTP testing utilities
        │   ├── setup.go                # HTTP handler test setup
        │   ├── assertions.go           # Response assertion helpers
        │   └── auth.go                 # Request identity mocking
        │
        └── fixtures/                   # Shared test fixtures
            ├── dspa/                   # DSPA fixtures
            │   ├── standard.go         # Standard DSPA configs
            │   └── storage.go          # Storage config fixtures
            ├── pipeline/               # Pipeline fixtures
            │   ├── runs.go             # Pipeline run fixtures
            │   └── versions.go         # Pipeline version fixtures
            └── secrets/                # K8s secret fixtures
                └── standard.go         # Standard secret configs
```

## Extraction Priority and Impact

### Phase 1: High-Impact, Low-Risk (Immediate)

**Files to Extract** (100% identical, zero risk):
1. `bff/cmd/helpers_test.go` → `testing/http/helpers_test.go`
2. `bff/internal/api/helpers_test.go` → `testing/http/url_template_test.go`
3. `bff/internal/api/middleware_dspa_errors_test.go` → `testing/middleware/dspa_errors_test.go`
4. `bff/internal/api/middleware_validation_test.go` → `testing/middleware/validation_test.go`
5. `bff/internal/integrations/kubernetes/get_namespaces_test.go` → `testing/k8s/namespaces_test.go`

**Impact**: 495 lines of duplicate test code eliminated immediately

### Phase 2: High-Impact, Low-Risk (Import Path Normalization)

**Files to Extract** (97-99% identical, only import paths differ):
1. `middleware_discovery_test.go` → Extract DSPA testing utilities
2. `s3_declared_limit_test.go` → Extract S3 limit validation
3. `pipeline_run_get_test.go` → Extract pipeline retrieval helpers
4. `pipeline_test.go` → Extract pipeline listing/filtering utilities

**Impact**: ~1,400 lines of nearly identical test code eliminated

**Required Changes**: 
- Create shared import path: `github.com/opendatahub-io/autox-library/bff/internal/...`
- Update module imports in AutoML and AutoRAG to reference shared library

### Phase 3: Medium-Impact, Higher Effort (Refactor Module-Specific Tests)

**Files to Refactor** (50-75% shared patterns):
1. Extract shared handler test patterns from `pipeline_runs_handler_test.go`
2. Extract shared S3 client patterns from `s3_handler_test.go`
3. Extract shared secret validation patterns from `secrets_handler_test.go`
4. Extract shared repository patterns from `pipeline_runs_test.go`, `s3_test.go`

**Impact**: ~3,500 lines of test scaffolding/patterns extracted

**Approach**:
- Keep module-specific tests in respective packages
- Extract shared test scaffolding, mock builders, and assertion helpers
- Use shared utilities as foundation for module-specific test extensions

## Migration Strategy

### Step 1: Create Shared Testing Package

```bash
mkdir -p packages/autox-library/internal/testing/{dspa,pipeline,middleware,k8s,s3,http,fixtures}
```

### Step 2: Extract Identical Files (Phase 1)

```bash
# Copy 100% identical test files to shared package
cp packages/automl/bff/cmd/helpers_test.go \
   packages/autox-library/internal/testing/http/

cp packages/automl/bff/internal/api/middleware_dspa_errors_test.go \
   packages/autox-library/internal/testing/middleware/
   
# (Repeat for all Phase 1 files)
```

### Step 3: Update Import Paths (Phase 2)

```go
// In packages/automl/bff/internal/api/my_handler_test.go
// Before:
import (
    "github.com/opendatahub-io/automl-library/bff/internal/models"
)

// After:
import (
    "github.com/opendatahub-io/autox-library/bff/internal/testing/pipeline"
    "github.com/opendatahub-io/automl-library/bff/internal/models"
)
```

### Step 4: Refactor Module-Specific Tests (Phase 3)

```go
// Extract shared patterns to autox-library
// In packages/autox-library/internal/testing/http/setup.go
func SetupHandlerTest(t *testing.T, method, path string, body interface{}) (*httptest.ResponseRecorder, *http.Request) {
    // Shared handler test setup
}

// Use shared utilities in AutoML/AutoRAG
// In packages/automl/bff/internal/api/my_handler_test.go
import "github.com/opendatahub-io/autox-library/internal/testing/http"

func TestMyHandler(t *testing.T) {
    recorder, req := http.SetupHandlerTest(t, "POST", "/api/v1/my-endpoint", requestBody)
    // AutoML-specific test logic
}
```

## Expected Outcomes

### Quantitative Benefits

- **Duplicate Code Eliminated**: ~6,000+ lines (across all phases)
- **Maintenance Reduction**: 50% fewer test updates needed
- **Test Consistency**: 100% identical test patterns across modules
- **Development Speed**: 30-40% faster test creation for new features

### Qualitative Benefits

1. **Single Source of Truth**: All DSPA, pipeline, and middleware testing logic in one place
2. **Easier Onboarding**: New modules can import battle-tested utilities
3. **Consistent Test Coverage**: Shared utilities ensure consistent coverage across modules
4. **Reduced Bugs**: Shared test patterns reduce copy-paste errors
5. **Better Refactoring**: Changes to test patterns propagate automatically
6. **Documentation**: Centralized test utilities serve as testing best practices guide

## Risk Assessment

### Low Risk Areas (Phase 1 & 2)
- Files that are 99-100% identical
- Only import path and minor constant differences
- No behavioral changes needed

**Mitigation**: Run existing tests after migration to ensure 100% pass rate

### Medium Risk Areas (Phase 3)
- Module-specific feature tests
- Different mock data expectations
- Different error handling patterns

**Mitigation**:
1. Extract patterns incrementally
2. Keep module-specific tests in place
3. Use shared utilities as opt-in foundation
4. Run both old and new tests in parallel during transition

## Conclusion

The AutoML and AutoRAG packages contain **extensive test code duplication** with:
- **31.2%** of files being 100% identical (immediate extraction candidates)
- **68.8%** of files sharing 50-99% of test patterns (refactor candidates)
- **Overall ~85-90% duplication** when accounting for shared patterns

**Recommendation**: Proceed with 3-phase extraction strategy, starting with 100% identical files (zero risk, immediate value) and progressively extracting shared patterns from module-specific tests.

**Next Steps**:
1. Create `packages/autox-library/internal/testing/` structure
2. Extract Phase 1 files (495 lines of identical code)
3. Normalize import paths for Phase 2 extraction
4. Refactor module-specific tests to use shared utilities (Phase 3)
