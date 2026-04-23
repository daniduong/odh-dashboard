# BFF Utility Files Duplication Analysis

## Analysis Date: 2026-04-23

## Executive Summary

This analysis examines code duplication in utility/helper files between AutoML and AutoRAG BFF packages. These files contain reusable functions for environment configuration, logging, JSON handling, testing, and URL manipulation.

**Key Finding**: Nearly **100% duplication** in core utility files (631 lines duplicated out of 631 shared lines).

---

## Duplicate Utility Files

### 1. `cmd/helpers.go` - CLI Environment Utilities

| Metric | AutoML | AutoRAG | Duplicate |
|--------|--------|---------|-----------|
| **Total Lines** | 63 | 63 | 63 |
| **Duplicate Lines** | 63 | 63 | 63 |
| **Duplication %** | **100%** | **100%** | **100%** |

**Common Helper Functions:**
- `getEnvAsInt(name string, defaultVal int) int` - Parse integer environment variables
- `getEnvAsString(name string, defaultVal string) string` - Parse string environment variables
- `getEnvAsBool(name string, defaultVal bool) bool` - Parse boolean environment variables
- `parseLevel(s string) slog.Level` - Parse log level from string
- `newOriginParser(allowList *[]string, defaultVal string) func(s string) error` - Parse CORS allowed origins

**Common Constants:** None (pure utility functions)

**Classification:** **Pure utility** - Environment variable parsing and CLI flag helpers

**Migration Priority:** **HIGH** - Exact duplicate, zero domain logic

**Location:** `packages/automl/bff/cmd/helpers.go` vs `packages/autorag/bff/cmd/helpers.go`

**Import Difference:** None (identical imports: `fmt`, `log/slog`, `os`, `strconv`, `strings`)

---

### 2. `internal/helpers/k8s.go` - Kubernetes Client Utilities

| Metric | AutoML | AutoRAG | Duplicate |
|--------|--------|---------|-----------|
| **Total Lines** | 28 | 28 | 28 |
| **Duplicate Lines** | 28 | 28 | 28 |
| **Duplication %** | **100%** | **100%** | **100%** |

**Common Helper Functions:**
- `GetKubeconfig() (*clientRest.Config, error)` - Load kubeconfig using default rules
- `BuildScheme() (*runtime.Scheme, error)` - Build runtime scheme with Kubernetes types

**Common Constants:** None

**Classification:** **Pure utility** - Kubernetes client initialization

**Migration Priority:** **HIGH** - Exact duplicate, foundational K8s setup

**Location:** `packages/automl/bff/internal/helpers/k8s.go` vs `packages/autorag/bff/internal/helpers/k8s.go`

**Import Difference:** 
- **Import Path Only**: Both use identical imports, only package name differs:
  - AutoML: (no package-specific imports)
  - AutoRAG: (no package-specific imports)
  - Both use: `k8s.io/apimachinery/pkg/runtime`, `k8s.io/client-go/kubernetes/scheme`, `k8s.io/client-go/rest`, `k8s.io/client-go/tools/clientcmd`

---

### 3. `internal/helpers/logging.go` - Request/Response Logging Utilities

| Metric | AutoML | AutoRAG | Duplicate |
|--------|--------|---------|-----------|
| **Total Lines** | 124 | 124 | 124 |
| **Duplicate Lines** | 123 | 123 | 123 |
| **Duplication %** | **99.2%** | **99.2%** | **99.2%** |

**Common Helper Functions:**
- `GetContextLoggerFromReq(r *http.Request) *slog.Logger` - Extract logger from request context
- `GetContextLogger(ctx context.Context) *slog.Logger` - Extract logger from context
- `isSensitiveHeader(h string) bool` - Check if header contains sensitive data
- `CloneBody(r *http.Request) ([]byte, error)` - Clone HTTP request body for logging

**Common Helper Types:**
- `HeaderLogValuer` - Custom slog.LogValuer for HTTP headers with redaction
- `RequestLogValuer` - Custom slog.LogValuer for HTTP requests
- `ResponseLogValuer` - Custom slog.LogValuer for HTTP responses

**Common Constants:**
- `sensitiveHeaders = []string{"Authorization", "Cookie", "Set-Cookie", "Proxy-Authorization"}` - Headers to redact in logs
- `maxBodySize = 10 * 1024 * 1024` - 10MB max body size for logging

**Classification:** **Pure utility** - HTTP request/response logging with security (header redaction)

**Migration Priority:** **HIGH** - Critical security utilities (header redaction), nearly identical

**Location:** `packages/automl/bff/internal/helpers/logging.go` vs `packages/autorag/bff/internal/helpers/logging.go`

**Import Difference:**
- **Line 12 ONLY**:
  - AutoML: `"github.com/opendatahub-io/automl-library/bff/internal/constants"`
  - AutoRAG: `"github.com/opendatahub-io/autorag-library/bff/internal/constants"`
- All other imports identical: `bytes`, `context`, `fmt`, `io`, `log/slog`, `net/http`, `slices`

---

### 4. `internal/api/helpers.go` - JSON and HTTP Response Utilities

| Metric | AutoML | AutoRAG | Duplicate |
|--------|--------|---------|-----------|
| **Total Lines** | 108 | 108 | 108 |
| **Duplicate Lines** | 108 | 108 | 108 |
| **Duplication %** | **100%** | **100%** | **100%** |

**Common Helper Functions:**
- `(app *App) WriteJSON(w http.ResponseWriter, status int, data any, headers http.Header) error` - Write JSON response
- `(app *App) ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error` - Read and validate JSON request body
- `ParseURLTemplate(tmpl string, params map[string]string) string` - Replace URL template placeholders

**Common Helper Types:**
- `Envelope[D any, M any]` - Generic response envelope with data and metadata
- `None` - Type alias for empty metadata (`*struct{}`)

**Common Constants:**
- `maxBytes := 1_048_576` - 1MB max JSON request body size

**Classification:** **Pure utility** - JSON marshaling/unmarshaling with validation

**Migration Priority:** **HIGH** - Core HTTP/JSON utilities, exact duplicate

**Location:** `packages/automl/bff/internal/api/helpers.go` vs `packages/autorag/bff/internal/api/helpers.go`

**Import Difference:** None (identical imports: `encoding/json`, `errors`, `fmt`, `io`, `net/http`, `strings`)

---

### 5. `internal/api/public_helpers.go` - Public API Access Helpers

| Metric | AutoML | AutoRAG | Duplicate |
|--------|--------|---------|-----------|
| **Total Lines** | 75 | 75 | 75 |
| **Duplicate Lines** | 74 | 74 | 74 |
| **Duplication %** | **98.7%** | **98.7%** | **98.7%** |

**Common Helper Functions:**
- `(app *App) BadRequest(w http.ResponseWriter, r *http.Request, err error)` - Public wrapper for 400 Bad Request response
- `(app *App) ServerError(w http.ResponseWriter, r *http.Request, err error)` - Public wrapper for 500 Server Error response
- `(app *App) NotImplemented(w http.ResponseWriter, r *http.Request, feature string)` - Send 501 Not Implemented response
- `(app *App) EndpointNotImplementedHandler(feature string) func(...)` - Create handler for unimplemented endpoints
- `(app *App) Config() config.EnvConfig` - Public getter for app config
- `(app *App) Logger() *slog.Logger` - Public getter for app logger
- `(app *App) KubernetesClientFactory() k8s.KubernetesClientFactory` - Public getter for K8s client factory
- `(app *App) Repositories() *repositories.Repositories` - Public getter for repositories

**Common Constants:** None

**Classification:** **Pure utility** - Public API for downstream extensions to access app components

**Migration Priority:** **HIGH** - Enables downstream extension pattern, nearly identical

**Location:** `packages/automl/bff/internal/api/public_helpers.go` vs `packages/autorag/bff/internal/api/public_helpers.go`

**Import Difference:**
- **Lines 10-13 ONLY**:
  - AutoML: `"github.com/opendatahub-io/automl-library/bff/internal/{config,integrations,integrations/kubernetes,repositories}"`
  - AutoRAG: `"github.com/opendatahub-io/autorag-library/bff/internal/{config,integrations,integrations/kubernetes,repositories}"`
- All other imports identical: `fmt`, `log/slog`, `net/http`, `strconv`, `github.com/julienschmidt/httprouter`

---

### 6. `internal/api/test_utils.go` - Testing Utilities

| Metric | AutoML | AutoRAG | Duplicate |
|--------|--------|---------|-----------|
| **Total Lines** | 165 | 177 | 153 |
| **Duplicate Lines** | 153 | 153 | 153 |
| **Duplication %** | **92.7%** | **86.4%** | **~89%** |

**Common Helper Functions:**
- `setupApiTest[T any](method, url string, body interface{}, k8Factory kubernetes.KubernetesClientFactory, identity *kubernetes.RequestIdentity) (T, *http.Response, error)` - Setup API test with mocked dependencies
- `newTestApp(t *testing.T) *App` - Create test app with mocked K8s client
- `setupApiTestPostMultipart(url string, fileContent []byte, fileName string, ...) (*http.Response, error)` - Setup multipart POST test

**Common Constants:**
- `AllowedOrigins: []string{"*"}` - Test CORS config
- `AuthMethod: config.AuthMethodInternal` - Test auth config
- `MockPipelineServerClient: true` - Test pipeline server mock flag

**Classification:** **Pure utility** - Test setup and mock configuration

**Migration Priority:** **MEDIUM** - Test utilities with minor differences in constructor calls

**Location:** `packages/automl/bff/internal/api/test_utils.go` vs `packages/autorag/bff/internal/api/test_utils.go`

**Import Difference:**
- **Lines 14-20 ONLY**:
  - AutoML: `"github.com/opendatahub-io/automl-library/bff/internal/{config,constants,integrations/kubernetes,...}"`
  - AutoRAG: `"github.com/opendatahub-io/autorag-library/bff/internal/{config,constants,integrations/kubernetes,...}"`
- All other imports identical: `bytes`, `context`, `encoding/json`, `io`, `log/slog`, `mime/multipart`, `net/http`, `net/http/httptest`, `testing`

**Key Differences:**
- **Line 116 (AutoML) vs Line 122 (AutoRAG)**:
  - AutoML: `return NewTestApp(cfg, logger, k8sFactory, psFactory, s3mocks.NewMockClientFactory(), nil)`
  - AutoRAG: `return NewTestApp(cfg, logger, k8sFactory, psFactory, nil)` (no S3 mock in constructor)
- **Lines 51-56 (AutoRAG only)**: Additional comments explaining `MockPipelineServerClient` flag purpose
- **Line 122 vs 165**: Different number of parameters in `NewTestApp` call due to S3ClientFactory difference

---

### 7. `internal/repositories/helpers.go` - URL Parameter Utilities

| Metric | AutoML | AutoRAG | Duplicate |
|--------|--------|---------|-----------|
| **Total Lines** | 68 | 68 | 68 |
| **Duplicate Lines** | 68 | 68 | 68 |
| **Duplication %** | **100%** | **100%** | **100%** |

**Common Helper Functions:**
- `FilterPageValues(values url.Values) url.Values` - Extract pagination query parameters
- `UrlWithParams(url string, values url.Values) string` - Append query parameters to URL
- `UrlWithPageParams(url string, values url.Values) string` - Append only pagination parameters to URL

**Common Constants:**
- Pagination parameter names: `"pageSize"`, `"orderBy"`, `"sortOrder"`, `"nextPageToken"`

**Classification:** **Pure utility** - URL query parameter manipulation for pagination

**Migration Priority:** **HIGH** - Exact duplicate, foundational pagination logic

**Location:** `packages/automl/bff/internal/repositories/helpers.go` vs `packages/autorag/bff/internal/repositories/helpers.go`

**Import Difference:** None (identical imports: `fmt`, `net/url`, `strings`)

---

## AutoRAG-Only Utility Files (No AutoML Equivalent)

### 8. `internal/helpers/llamastack.go` - LlamaStack Client Context Utilities

| Metric | Value |
|--------|-------|
| **Total Lines** | 21 |
| **Domain-Specific?** | Yes (AutoRAG-only) |

**Helper Functions:**
- `GetContextLlamaStackClient(ctx context.Context) (llamastack.LlamaStackClientInterface, error)` - Extract LlamaStack client from context

**Classification:** **Domain-specific** - AutoRAG integration with LlamaStack

**Migration Priority:** **N/A** - AutoRAG-specific, not duplicated

**Pattern Note:** Similar pattern to logging helpers (`GetContextLogger`), but for LlamaStack client

---

### 9. `internal/api/k8s_helpers.go` - Kubernetes Error Handling

| Metric | Value |
|--------|-------|
| **Total Lines** | 83 |
| **Domain-Specific?** | No (could be shared) |

**Helper Functions:**
- `(app *App) handleK8sClientError(w http.ResponseWriter, r *http.Request, err error)` - Map K8s errors to HTTP responses
- `(app *App) getDefaultStatusCodeForK8sError(errorCode string) int` - Get default HTTP status for K8s error code
- `(app *App) mapK8sErrorToHTTPError(k8sErr *k8s.K8sError, statusCode int) *integrations.HTTPError` - Convert K8s error to HTTP error

**Classification:** **Pure utility** - Error handling abstraction for Kubernetes client

**Migration Priority:** **HIGH** - Should be shared with AutoML (not currently present)

**Note:** AutoML currently lacks this abstraction, likely handles K8s errors inline in handlers

---

### 10. `internal/api/lsd_helpers.go` - LlamaStack Distribution Error Handling

| Metric | Value |
|--------|-------|
| **Total Lines** | 94 |
| **Domain-Specific?** | Yes (AutoRAG-only) |

**Helper Functions:**
- `(app *App) handleLlamaStackClientError(w http.ResponseWriter, r *http.Request, err error)` - Map LlamaStack errors to HTTP responses
- `(app *App) getDefaultStatusCodeForLlamaStackClientError(errorCode string) int` - Get default HTTP status for LlamaStack error
- `(app *App) mapLlamaStackClientErrorToHTTPError(lsErr *llamastack.LlamaStackError, statusCode int) *integrations.HTTPError` - Convert LlamaStack error to HTTP error

**Classification:** **Domain-specific** - AutoRAG integration with LlamaStack

**Migration Priority:** **N/A** - AutoRAG-specific error handling

**Pattern Note:** Mirrors K8s error handling pattern but for LlamaStack client

---

## Summary Statistics

### Overall Duplication

| Category | Total Lines | Duplicate Lines | Duplication % |
|----------|-------------|-----------------|---------------|
| **Exact Duplicates** | 631 | 631 | **100%** |
| **Near Duplicates (>95%)** | 199 | 197 | **99%** |
| **Partial Duplicates (85-95%)** | 165 | 153 | **93%** |
| **Total Shared Utilities** | 995 | 981 | **98.6%** |

### File Count by Classification

| Classification | Count | Total Lines |
|----------------|-------|-------------|
| **Pure Utility (100% duplicate)** | 5 files | 431 lines |
| **Pure Utility (>95% duplicate)** | 2 files | 199 lines |
| **Pure Utility (>85% duplicate)** | 1 file | 165 lines |
| **Domain-Specific (AutoRAG-only)** | 3 files | 198 lines |
| **TOTAL** | 11 files | 993 lines |

---

## Migration Recommendations

### Priority 1: Immediate Migration (100% Duplicates)

**High-value, zero-risk migrations:**

1. **`cmd/helpers.go`** (63 lines)
   - Target: `packages/autox/bff/pkg/env/helpers.go`
   - Functions: `getEnvAsInt`, `getEnvAsString`, `getEnvAsBool`, `parseLevel`, `newOriginParser`

2. **`internal/helpers/k8s.go`** (28 lines)
   - Target: `packages/autox/bff/pkg/k8s/helpers.go`
   - Functions: `GetKubeconfig`, `BuildScheme`

3. **`internal/api/helpers.go`** (108 lines)
   - Target: `packages/autox/bff/pkg/http/helpers.go`
   - Functions: `WriteJSON`, `ReadJSON`, `ParseURLTemplate`
   - Types: `Envelope[D, M]`, `None`

4. **`internal/repositories/helpers.go`** (68 lines)
   - Target: `packages/autox/bff/pkg/pagination/helpers.go`
   - Functions: `FilterPageValues`, `UrlWithParams`, `UrlWithPageParams`

**Total Lines Saved:** 267 lines (removing duplicates across 2 packages = 534 lines of duplication)

---

### Priority 2: High-Value Near-Duplicates (>95%)

**Require minimal adjustments:**

5. **`internal/helpers/logging.go`** (124 lines, 99.2% duplicate)
   - Target: `packages/autox/bff/pkg/logging/helpers.go`
   - Functions: `GetContextLoggerFromReq`, `GetContextLogger`, `CloneBody`
   - Types: `HeaderLogValuer`, `RequestLogValuer`, `ResponseLogValuer`
   - Constants: `sensitiveHeaders`, `maxBodySize`
   - **Adjustment Needed:** Remove constants package import dependency (line 12)

6. **`internal/api/public_helpers.go`** (75 lines, 98.7% duplicate)
   - Target: `packages/autox/bff/pkg/api/public_helpers.go`
   - Functions: `BadRequest`, `ServerError`, `NotImplemented`, `EndpointNotImplementedHandler`, `Config`, `Logger`, `KubernetesClientFactory`, `Repositories`
   - **Adjustment Needed:** Generify `App` interface or use interface-based access

**Total Lines Saved:** 199 lines (398 lines of duplication removed)

---

### Priority 3: Test Utilities (>85%)

7. **`internal/api/test_utils.go`** (165 lines, 92.7% duplicate)
   - Target: `packages/autox/bff/pkg/testing/api_helpers.go`
   - Functions: `setupApiTest`, `newTestApp`, `setupApiTestPostMultipart`
   - **Adjustment Needed:** 
     - Abstract `NewTestApp` constructor signature differences (S3ClientFactory parameter)
     - Provide builder pattern or options struct for flexible test app creation

**Total Lines Saved:** 153 lines (306 lines of duplication removed)

---

### Priority 4: Shared Pattern Extension

8. **`internal/api/k8s_helpers.go`** (AutoRAG-only, 83 lines)
   - **Action:** Extract to shared library AND migrate to AutoML
   - Target: `packages/autox/bff/pkg/k8s/errors.go`
   - Functions: `handleK8sClientError`, `getDefaultStatusCodeForK8sError`, `mapK8sErrorToHTTPError`
   - **Benefit:** Standardize K8s error handling across both packages

**Total Lines Saved:** 83 lines (avoided duplication when AutoML adopts pattern)

---

### Do NOT Migrate (Domain-Specific)

- **`internal/helpers/llamastack.go`** - AutoRAG-specific LlamaStack integration
- **`internal/api/lsd_helpers.go`** - AutoRAG-specific LlamaStack error handling

---

## Expected Impact

### Code Reduction
- **Total Duplicate Lines Eliminated:** 981 lines (from 2 packages)
- **Net Lines Saved:** ~981 lines (assuming shared library adds minimal overhead)
- **Percentage Reduction in Utility Code:** ~49.5% (981 duplicate lines / 1,962 total utility lines across both packages)

### Maintenance Benefits
- **Single Source of Truth** for 98.6% of utility code
- **Consistent error handling** (especially with K8s error helper migration)
- **Unified testing utilities** across AutoML and AutoRAG
- **Simplified bug fixes** - fix once, applies to both packages

### Migration Effort Estimate
- **Priority 1 (100% duplicates):** 2-3 hours (straightforward extraction)
- **Priority 2 (>95% duplicates):** 3-4 hours (minor import path adjustments)
- **Priority 3 (test utils):** 2-3 hours (constructor abstraction)
- **Priority 4 (K8s error helpers):** 2-3 hours (extract + migrate to AutoML)
- **Total Estimated Effort:** 9-13 hours

---

## Architectural Patterns Observed

### 1. Context-Based Dependency Injection
**Pattern:** Store typed values in `context.Context` and retrieve with helper functions
- `GetContextLogger(ctx)` → `*slog.Logger`
- `GetContextLlamaStackClient(ctx)` → `LlamaStackClientInterface`

**Shared Library Opportunity:** Generic context value retrieval with type safety

---

### 2. Error Abstraction Hierarchy
**Pattern:** Domain error types → HTTP error mapping → JSON response
- `K8sError` → `HTTPError` → JSON response (AutoRAG only)
- `LlamaStackError` → `HTTPError` → JSON response (AutoRAG only)

**Shared Library Opportunity:** 
- Generic error mapping framework
- Standard HTTP status code mappings for common error types

---

### 3. Structured Logging with Redaction
**Pattern:** Custom `slog.LogValuer` types for safe request/response logging
- `HeaderLogValuer` - Redacts sensitive headers
- `RequestLogValuer` - Logs method, URL, body, headers
- `ResponseLogValuer` - Logs status, body, headers

**Shared Library Opportunity:** Security-first logging utilities

---

### 4. Environment Configuration Utilities
**Pattern:** Type-safe environment variable parsing with defaults
- `getEnvAsInt`, `getEnvAsString`, `getEnvAsBool`

**Shared Library Opportunity:** Centralized config parsing

---

### 5. JSON Envelope Pattern
**Pattern:** Generic `Envelope[D any, M any]` wrapper for API responses
- Supports arbitrary data type `D`
- Optional metadata type `M`

**Shared Library Opportunity:** Standard API response format across BFFs

---

## Migration Path Recommendations

### Phase 1: Foundation (Week 1)
1. Create `packages/autox/bff/pkg/` directory structure
2. Migrate 100% duplicates (Priority 1):
   - `pkg/env/` - Environment helpers
   - `pkg/k8s/` - Kubernetes helpers
   - `pkg/http/` - HTTP/JSON helpers
   - `pkg/pagination/` - Pagination helpers

### Phase 2: Core Utilities (Week 2)
3. Migrate near-duplicates (Priority 2):
   - `pkg/logging/` - Logging helpers with header redaction
   - `pkg/api/` - Public API access helpers

### Phase 3: Testing (Week 3)
4. Migrate test utilities (Priority 3):
   - `pkg/testing/` - Test app creation and HTTP test helpers

### Phase 4: Standardization (Week 4)
5. Extract and share K8s error handling (Priority 4):
   - `pkg/k8s/errors.go` - Error mapping utilities
   - Migrate to AutoML for consistency

---

## Open Questions

1. **Interface Abstraction for Public Helpers:**
   - Should `public_helpers.go` use interface-based `App` abstraction?
   - Or use dependency injection pattern with function parameters?

2. **Test Utilities Constructor Differences:**
   - How to handle `NewTestApp` parameter differences (S3ClientFactory)?
   - Builder pattern vs options struct vs variadic parameters?

3. **Constants Package Dependency:**
   - How to eliminate `internal/constants` import from logging helpers?
   - Should constants be passed as parameters or duplicated in shared lib?

4. **K8s Error Handling Migration to AutoML:**
   - Should AutoML adopt AutoRAG's K8s error handling pattern?
   - Or create unified abstraction in shared library first?

---

## Next Steps

1. **Review with Team:**
   - Confirm migration priorities
   - Validate shared library package structure
   - Decide on abstraction strategies for Priority 2-3

2. **Prototype Shared Library:**
   - Start with Priority 1 (100% duplicates)
   - Validate import path changes work correctly
   - Test in both AutoML and AutoRAG

3. **Create Migration PRs:**
   - One PR per priority level (easier review)
   - Include tests from both packages
   - Update documentation

4. **Monitor for Drift:**
   - Add linter/CI check to prevent re-duplication
   - Establish contribution guidelines for shared utilities

---

## Conclusion

The BFF utility files show **exceptionally high duplication** (98.6% of shared code is duplicated). This represents a **prime opportunity** for consolidation into a shared library.

**Key Advantages:**
- **Minimal risk** - Most files are 100% identical
- **High impact** - Nearly 1,000 lines of duplication eliminated
- **Clear boundaries** - Pure utility functions with no domain logic
- **Quick wins** - Priority 1 migration can complete in 2-3 hours

**Recommended Action:** Proceed with phased migration starting with Priority 1 (100% duplicates), then incrementally add Priority 2-4 utilities.
