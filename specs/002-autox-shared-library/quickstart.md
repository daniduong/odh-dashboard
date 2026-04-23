# AutoX Shared Library - Developer Quickstart

**Audience**: Developers working on AutoML or AutoRAG packages  
**Goal**: Get started using AutoX shared library in 15 minutes  
**Date**: 2026-04-23

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Project Structure Overview](#project-structure-overview)
3. [Development Setup](#development-setup)
4. [Frontend Usage](#frontend-usage)
5. [BFF Usage](#bff-usage)
6. [Testing](#testing)
7. [Common Patterns](#common-patterns)
8. [Troubleshooting](#troubleshooting)
9. [Migration Guide](#migration-guide)

---

## Prerequisites

### Required Versions

- **Node.js**: >= 22.0.0
- **npm**: >= 10.0.0
- **Go**: >= 1.24.3
- **Git**: Latest stable

### Knowledge Requirements

- Familiarity with ODH Dashboard monorepo structure
- Basic understanding of npm workspaces and Go workspaces
- React 18 and TypeScript experience (for frontend)
- Go 1.24 experience (for BFF)

---

## Project Structure Overview

```text
packages/
├── autox/                       # Shared library package
│   ├── frontend/                # UI shared code
│   │   └── src/
│   │       ├── hooks/           # React hooks
│   │       ├── components/      # UI primitives
│   │       ├── utils/           # Utilities
│   │       └── types/           # TypeScript types
│   └── bff/                     # BFF shared code
│       └── internal/
│           ├── models/          # Data models
│           ├── integrations/    # External service clients
│           ├── repositories/    # Business logic
│           ├── api/             # HTTP utilities
│           └── utils/           # Shared utilities
│
├── automl/                      # AutoML package (consumer)
│   ├── frontend/                # AutoML UI
│   └── bff/                     # AutoML BFF
│
└── autorag/                     # AutoRAG package (consumer)
    ├── frontend/                # AutoRAG UI
    └── bff/                     # AutoRAG BFF
```

**Key Concept**: AutoX exports low-level primitives. AutoML and AutoRAG compose them into domain-specific features.

---

## Development Setup

### Step 1: Clone and Install

```bash
# Clone repository (if not already cloned)
git clone https://github.com/opendatahub-io/odh-dashboard.git
cd odh-dashboard

# Install all dependencies (sets up npm workspaces)
npm install
```

**What happens**:
- npm workspaces automatically link `@odh-dashboard/autox`, `@odh-dashboard/automl`, `@odh-dashboard/autorag`
- All `node_modules` are hoisted to root (except for package-specific dependencies)
- Local changes to AutoX are immediately available to AutoML/AutoRAG

### Step 2: Set Up Go Workspaces

```bash
# From repository root
go work init

# Add AutoX, AutoML, AutoRAG BFF packages
go work use packages/autox/bff
go work use packages/automl/bff
go work use packages/autorag/bff

# Verify workspace
go work sync
go work vendor  # Optional: vendor dependencies
```

**What happens**:
- Go workspaces link all BFF packages in a single module
- Local changes to AutoX BFF are immediately available to AutoML/AutoRAG BFF
- `go build` and `go test` work across all packages

### Step 3: Verify Setup

**Frontend**:

```bash
# From repository root
cd packages/automl/frontend
npm run type-check  # Should succeed with AutoX imports

cd ../../../packages/autorag/frontend
npm run type-check  # Should succeed with AutoX imports
```

**BFF**:

```bash
# From repository root
cd packages/automl/bff
go build ./...  # Should succeed with AutoX imports

cd ../../autorag/bff
go build ./...  # Should succeed with AutoX imports
```

**Success Indicators**:
- ✅ No TypeScript errors about missing `@odh-dashboard/autox` modules
- ✅ No Go compilation errors about missing AutoX packages
- ✅ All packages compile successfully

---

## Frontend Usage

### Import Patterns

```typescript
// Import hooks
import { usePipelineRuns, useNotification, useS3ListFiles } from '@odh-dashboard/autox/hooks';

// Import components
import { FileExplorer, PipelineRunsTable } from '@odh-dashboard/autox/components';

// Import types
import type { PipelineRun, S3Object, RuntimeStateKF } from '@odh-dashboard/autox/types';

// Import utilities
import { isPipelineRunning, isValidK8sName } from '@odh-dashboard/autox/utils';
```

### Example 1: Using usePipelineRuns Hook

```typescript
// packages/automl/frontend/src/components/RunsList.tsx
import React from 'react';
import { usePipelineRuns } from '@odh-dashboard/autox/hooks';
import { PipelineRunsTable } from '@odh-dashboard/autox/components';
import type { PipelineRun } from '@odh-dashboard/autox/types';

export const AutoMLRunsList: React.FC = () => {
  const namespace = useNamespace(); // AutoML-specific hook
  
  const { data: runs, isLoading, error } = usePipelineRuns({
    namespace,
    enabled: !!namespace,
  });

  if (isLoading) return <Spinner />;
  if (error) return <Alert variant="danger">{error.message}</Alert>;

  return (
    <PipelineRunsTable
      runs={runs || []}
      onRunSelect={(run) => navigate(`/automl/runs/${run.id}`)}
      showActions
    />
  );
};
```

### Example 2: Composing FileExplorer

```typescript
// packages/autorag/frontend/src/components/RAGDataBrowser.tsx
import React from 'react';
import { FileExplorer } from '@odh-dashboard/autox/components';
import { useS3ListFiles } from '@odh-dashboard/autox/hooks';
import type { FileItem } from '@odh-dashboard/autox/types';

export const RAGDataBrowser: React.FC<{ namespace: string; bucket: string }> = ({
  namespace,
  bucket,
}) => {
  const [currentPath, setCurrentPath] = React.useState('/');
  const { data, isLoading } = useS3ListFiles({
    namespace,
    bucket,
    prefix: currentPath,
    recursive: false,
  });

  const files: FileItem[] = data?.objects.map(obj => ({
    name: obj.key.split('/').pop()!,
    type: obj.key.endsWith('/') ? 'directory' : 'file',
    size: obj.size,
    lastModified: obj.lastModified,
  })) || [];

  return (
    <FileExplorer
      files={files}
      currentPath={currentPath}
      onNavigate={setCurrentPath}
      onFileSelect={(file) => handleFileSelect(file)}
      isLoading={isLoading}
      searchable
    />
  );
};
```

### Example 3: Using Validation Utilities

```typescript
// packages/automl/frontend/src/components/CreatePipelineForm.tsx
import React from 'react';
import { isValidK8sName } from '@odh-dashboard/autox/utils';
import { TextInput } from '@patternfly/react-core';

export const PipelineNameInput: React.FC = () => {
  const [name, setName] = React.useState('');
  const [error, setError] = React.useState('');

  const handleChange = (value: string) => {
    setName(value);
    if (!isValidK8sName(value)) {
      setError('Name must be lowercase alphanumeric with hyphens, max 63 chars');
    } else {
      setError('');
    }
  };

  return (
    <TextInput
      value={name}
      onChange={(_, value) => handleChange(value)}
      validated={error ? 'error' : 'default'}
      helperTextInvalid={error}
    />
  );
};
```

### Module Federation Setup

**Add AutoX to webpack config**:

```javascript
// packages/automl/frontend/config/moduleFederation.js
const { dependencies } = require('../package.json');

module.exports = {
  name: 'automl',
  filename: 'remoteEntry.js',
  exposes: {
    './extensions': './src/odh/extensions.ts',
  },
  shared: {
    react: {
      singleton: true,
      requiredVersion: dependencies['react'],
    },
    'react-dom': {
      singleton: true,
      requiredVersion: dependencies['react-dom'],
    },
    // ⭐ Add AutoX as shared singleton
    '@odh-dashboard/autox': {
      singleton: true,
      requiredVersion: dependencies['@odh-dashboard/autox'],
    },
    '@tanstack/react-query': {
      singleton: true,
      requiredVersion: dependencies['@tanstack/react-query'],
    },
  },
};
```

**Why `singleton: true`**: Ensures only one instance of AutoX hooks/state exists at runtime, preventing duplicate React Query cache and hook state issues.

---

## BFF Usage

### Import Patterns

```go
// Import models
import "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/models"

// Import integrations
import k8s "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/integrations/kubernetes"
import s3 "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/integrations/s3"
import ps "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/integrations/pipelineserver"

// Import repositories
import "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/repositories"

// Import API utilities
import "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/api"
```

### Example 1: Using Kubernetes Client

```go
// packages/automl/bff/internal/handlers/secrets_handler.go
package handlers

import (
    "context"
    "net/http"
    
    k8s "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/integrations/kubernetes"
    "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/models"
    "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/api"
)

type SecretsHandler struct {
    k8sClient k8s.KubernetesClientInterface
}

func (h *SecretsHandler) ListSecrets(w http.ResponseWriter, r *http.Request) {
    namespace := r.URL.Query().Get("namespace")
    if namespace == "" {
        api.WriteBadRequest(w, "namespace parameter is required")
        return
    }
    
    secrets, err := h.k8sClient.ListSecrets(r.Context(), namespace, metav1.ListOptions{
        LabelSelector: "type=s3",
    })
    if err != nil {
        api.WriteServerError(w, err)
        return
    }
    
    // Transform to response model
    response := transformToSecretList(secrets)
    api.WriteJSON(w, http.StatusOK, response, nil)
}
```

### Example 2: Using Pipeline Repository

```go
// packages/autorag/bff/internal/handlers/pipelines_handler.go
package handlers

import (
    "context"
    "net/http"
    
    "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/repositories"
    ps "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/integrations/pipelineserver"
    "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/api"
)

type PipelinesHandler struct {
    pipelineRepo repositories.PipelineRepository
    psClient     ps.PipelineServerClientInterface
    baseURL      string
}

func (h *PipelinesHandler) DiscoverRAGPipelines(w http.ResponseWriter, r *http.Request) {
    namespace := r.URL.Query().Get("namespace")
    
    // AutoRAG-specific pipeline definitions
    definitions := map[string]string{
        "rag-pipeline":    "/pipelines/rag_pipeline.yaml",
        "vector-indexing": "/pipelines/vector_indexing.yaml",
    }
    
    pipelines, err := h.pipelineRepo.DiscoverNamedPipelines(
        h.psClient, r.Context(), namespace, h.baseURL, definitions,
    )
    if err != nil {
        api.WriteServerError(w, err)
        return
    }
    
    api.WriteJSON(w, http.StatusOK, pipelines, nil)
}
```

### Example 3: Using S3 Client

```go
// packages/automl/bff/internal/handlers/s3_handler.go
package handlers

import (
    "context"
    "net/http"
    
    s3 "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/integrations/s3"
    "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/api"
)

type S3Handler struct {
    s3Client s3.S3ClientInterface
}

func (h *S3Handler) ListFiles(w http.ResponseWriter, r *http.Request) {
    bucket := r.URL.Query().Get("bucket")
    prefix := r.URL.Query().Get("prefix")
    
    objects, err := h.s3Client.ListObjects(r.Context(), bucket, prefix, false)
    if err != nil {
        api.WriteServerError(w, err)
        return
    }
    
    api.WriteJSON(w, http.StatusOK, objects, nil)
}
```

### Dependency Injection Pattern

**Setup in main.go**:

```go
// packages/automl/bff/cmd/main.go
package main

import (
    "log"
    
    k8s "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/integrations/kubernetes"
    s3 "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/integrations/s3"
    "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/repositories"
    "github.com/opendatahub-io/odh-dashboard/packages/automl/bff/internal/handlers"
)

func main() {
    // Initialize AutoX clients
    k8sClient, err := k8s.NewInternalK8sClient()
    if err != nil {
        log.Fatalf("Failed to create K8s client: %v", err)
    }
    
    s3Client := s3.NewAWSS3Client(s3Config)
    pipelineRepo := repositories.NewPipelineRepository()
    
    // Inject into AutoML handlers
    secretsHandler := handlers.NewSecretsHandler(k8sClient)
    s3Handler := handlers.NewS3Handler(s3Client)
    pipelinesHandler := handlers.NewPipelinesHandler(pipelineRepo, psClient, baseURL)
    
    // Setup routes...
}
```

---

## Testing

### Frontend Testing

#### Unit Tests for Components Using AutoX

```typescript
// packages/automl/frontend/src/components/__tests__/RunsList.spec.tsx
import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { AutoMLRunsList } from '../RunsList';

// Mock AutoX hooks
jest.mock('@odh-dashboard/autox/hooks', () => ({
  usePipelineRuns: jest.fn(),
}));

import { usePipelineRuns } from '@odh-dashboard/autox/hooks';

describe('AutoMLRunsList', () => {
  const queryClient = new QueryClient();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should render pipeline runs', () => {
    (usePipelineRuns as jest.Mock).mockReturnValue({
      data: [
        { id: '1', name: 'Run 1', state: 'SUCCEEDED', createdAt: '2024-01-01' },
      ],
      isLoading: false,
      error: null,
    });

    render(
      <QueryClientProvider client={queryClient}>
        <AutoMLRunsList />
      </QueryClientProvider>
    );

    expect(screen.getByText('Run 1')).toBeInTheDocument();
  });
});
```

#### Integration Tests with AutoX

```typescript
// packages/autorag/frontend/src/__tests__/integration/RAGDataBrowser.spec.tsx
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useS3ListFiles } from '@odh-dashboard/autox/hooks';

// Real AutoX hook (no mocking)
describe('RAGDataBrowser integration', () => {
  it('should fetch S3 files using AutoX hook', async () => {
    const queryClient = new QueryClient();
    
    const { result } = renderHook(
      () => useS3ListFiles({ namespace: 'test', bucket: 'data' }),
      {
        wrapper: ({ children }) => (
          <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
        ),
      }
    );

    await waitFor(() => expect(result.current.isLoading).toBe(false));
    expect(result.current.data).toBeDefined();
  });
});
```

### BFF Testing

#### Unit Tests with AutoX Mocks

```go
// packages/automl/bff/internal/handlers/secrets_handler_test.go
package handlers

import (
    "net/http/httptest"
    "testing"
    
    k8smocks "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/integrations/kubernetes/mocks"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSecretsHandler_ListSecrets(t *testing.T) {
    // Use AutoX mock
    mockK8sClient := k8smocks.NewMockKubernetesClient()
    mockK8sClient.On("ListSecrets", mock.Anything, "test-ns", mock.Anything).Return(
        &corev1.SecretList{
            Items: []corev1.Secret{
                {ObjectMeta: metav1.ObjectMeta{Name: "secret1"}},
            },
        },
        nil,
    )
    
    handler := NewSecretsHandler(mockK8sClient)
    
    req := httptest.NewRequest("GET", "/api/secrets?namespace=test-ns", nil)
    w := httptest.NewRecorder()
    
    handler.ListSecrets(w, req)
    
    if w.Code != 200 {
        t.Errorf("Expected 200, got %d", w.Code)
    }
}
```

#### Integration Tests

```go
// packages/autorag/bff/internal/integration_test.go
package internal_test

import (
    "testing"
    
    "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/repositories"
    ps "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/integrations/pipelineserver"
)

func TestPipelineDiscovery_Integration(t *testing.T) {
    // Real AutoX repository (no mocking)
    pipelineRepo := repositories.NewPipelineRepository()
    psClient := ps.NewKFPClient(testConfig)
    
    definitions := map[string]string{
        "rag-pipeline": "/test/pipelines/rag.yaml",
    }
    
    pipelines, err := pipelineRepo.DiscoverNamedPipelines(
        psClient, ctx, "test-ns", "http://localhost:8080", definitions,
    )
    
    if err != nil {
        t.Fatalf("Discovery failed: %v", err)
    }
    
    if len(pipelines) == 0 {
        t.Error("Expected pipelines to be discovered")
    }
}
```

---

## Common Patterns

### Pattern 1: Domain-Specific Hook Wrapping AutoX

```typescript
// AutoRAG-specific hook
import { usePipelineRuns } from '@odh-dashboard/autox/hooks';
import type { PipelineRun } from '@odh-dashboard/autox/types';

export function useRAGPipelineRuns(namespace: string) {
  const query = usePipelineRuns({
    namespace,
    enabled: !!namespace,
  });

  // Add RAG-specific filtering
  const ragRuns = React.useMemo(() => {
    return query.data?.filter(run => run.pipelineType === 'rag') || [];
  }, [query.data]);

  return {
    ...query,
    data: ragRuns,
  };
}
```

### Pattern 2: Handler Files for Domain Logic

```go
// AutoML-specific validation (handler logic, not in AutoX)
package handlers

import (
    "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/repositories"
    "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/models"
)

func (h *AutoMLHandler) ValidateTimeSeriesParameters(params map[string]interface{}) error {
    // Use AutoX base validation
    if err := validateCommonParameters(params); err != nil {
        return err
    }
    
    // AutoML-specific validation
    if _, ok := params["target_column"]; !ok {
        return errors.New("target_column is required for time series")
    }
    
    return nil
}
```

### Pattern 3: Component Composition

```typescript
// AutoML-specific table extending AutoX primitive
import { PipelineRunsTable } from '@odh-dashboard/autox/components';
import type { PipelineRun } from '@odh-dashboard/autox/types';

export function AutoMLRunsTable({ runs }: { runs: PipelineRun[] }) {
  return (
    <div>
      <PipelineRunsTable
        runs={runs}
        showActions
        onRunSelect={(run) => navigate(`/automl/runs/${run.id}`)}
      />
      {/* Add AutoML-specific metrics summary */}
      <AutoMLMetricsSummary runs={runs} />
    </div>
  );
}
```

---

## Troubleshooting

### Issue: "Cannot find module '@odh-dashboard/autox'"

**Cause**: npm workspaces not set up correctly

**Solution**:

```bash
# From repository root
rm -rf node_modules package-lock.json
rm -rf packages/*/node_modules packages/*/package-lock.json
npm install
```

### Issue: Go import errors for AutoX packages

**Cause**: Go workspaces not initialized

**Solution**:

```bash
# From repository root
go work init
go work use packages/autox/bff
go work use packages/automl/bff
go work use packages/autorag/bff
go work sync
```

### Issue: Module Federation errors (duplicate React instances)

**Cause**: AutoX not configured as shared singleton

**Solution**: Verify `moduleFederation.js` includes:

```javascript
shared: {
  '@odh-dashboard/autox': {
    singleton: true,
    requiredVersion: dependencies['@odh-dashboard/autox'],
  },
}
```

### Issue: TypeScript errors after AutoX update

**Cause**: Cached TypeScript build info

**Solution**:

```bash
# From consumer package (e.g., packages/automl/frontend)
rm -rf .tsbuildinfo
npm run type-check
```

### Issue: Go build errors after AutoX update

**Cause**: Stale Go module cache

**Solution**:

```bash
# From repository root
go clean -modcache
go work sync
go build ./...
```

---

## Migration Guide

### Migrating Duplicate Code to AutoX

#### Step 1: Identify Duplicate Code

**Frontend Example**:

```typescript
// ❌ Before: Duplicate in both AutoML and AutoRAG
// packages/automl/frontend/src/hooks/usePipelineRuns.ts
export function usePipelineRuns(namespace: string) {
  return useQuery({
    queryKey: ['pipelineRuns', namespace],
    queryFn: async () => {
      const response = await fetch(`/api/pipeline-runs?namespace=${namespace}`);
      return response.json();
    },
  });
}

// packages/autorag/frontend/src/hooks/usePipelineRuns.ts
export function usePipelineRuns(namespace: string) {
  return useQuery({
    queryKey: ['pipelineRuns', namespace],
    queryFn: async () => {
      const response = await fetch(`/api/pipeline-runs?namespace=${namespace}`);
      return response.json();
    },
  });
}
```

**BFF Example**:

```go
// ❌ Before: Duplicate in both AutoML and AutoRAG
// packages/automl/bff/internal/models/pipeline_runs.go
type PipelineRun struct {
    ID     string `json:"id"`
    Name   string `json:"name"`
    State  string `json:"state"`
    // ...
}

// packages/autorag/bff/internal/models/pipeline_runs.go
type PipelineRun struct {
    ID     string `json:"id"`
    Name   string `json:"name"`
    State  string `json:"state"`
    // ...
}
```

#### Step 2: Move to AutoX

```bash
# Frontend: Move to AutoX
# (Already done in AutoX - just update imports)

# BFF: Move to AutoX
# (Already done in AutoX - just update imports)
```

#### Step 3: Update Imports

**Frontend**:

```typescript
// ✅ After: Import from AutoX
// packages/automl/frontend/src/components/RunsList.tsx
import { usePipelineRuns } from '@odh-dashboard/autox/hooks';

// packages/autorag/frontend/src/components/RunsList.tsx
import { usePipelineRuns } from '@odh-dashboard/autox/hooks';
```

**BFF**:

```go
// ✅ After: Import from AutoX
// packages/automl/bff/internal/handlers/runs_handler.go
import "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/models"

// packages/autorag/bff/internal/handlers/runs_handler.go
import "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/models"
```

#### Step 4: Delete Duplicate Files

```bash
# Frontend
rm packages/automl/frontend/src/hooks/usePipelineRuns.ts
rm packages/autorag/frontend/src/hooks/usePipelineRuns.ts

# BFF
rm packages/automl/bff/internal/models/pipeline_runs.go
rm packages/autorag/bff/internal/models/pipeline_runs.go
```

#### Step 5: Run Tests

```bash
# Frontend tests
cd packages/automl/frontend && npm test
cd ../../../packages/autorag/frontend && npm test

# BFF tests
cd ../../automl/bff && go test ./...
cd ../../autorag/bff && go test ./...
```

#### Step 6: Verify in Browser (Frontend Only)

```bash
# Start dev server
npm run dev

# Navigate to AutoML and AutoRAG features
# Verify functionality unchanged
```

---

## Quick Reference Card

### Frontend Cheat Sheet

```typescript
// Hooks
import { usePipelineRuns, useNotification, useS3ListFiles } from '@odh-dashboard/autox/hooks';

// Components
import { FileExplorer, PipelineRunsTable } from '@odh-dashboard/autox/components';

// Types
import type { PipelineRun, S3Object } from '@odh-dashboard/autox/types';

// Utils
import { isPipelineRunning, isValidK8sName } from '@odh-dashboard/autox/utils';

// Example usage
const { data, isLoading } = usePipelineRuns({ namespace: 'default' });
```

### BFF Cheat Sheet

```go
// Models
import "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/models"

// Clients
import k8s "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/integrations/kubernetes"
import s3 "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/integrations/s3"
import ps "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/integrations/pipelineserver"

// Repositories
import "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/repositories"

// API utilities
import "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/api"

// Example usage
secrets, err := k8sClient.ListSecrets(ctx, namespace, metav1.ListOptions{})
api.WriteJSON(w, http.StatusOK, secrets, nil)
```

---

## Next Steps

1. ✅ Read [contracts/bff-api.md](./contracts/bff-api.md) for BFF API reference
2. ✅ Read [contracts/ui-api.md](./contracts/ui-api.md) for UI API reference
3. ✅ Review [research.md](./research.md) for extraction roadmap
4. ✅ Start with Phase 1 perfect duplicates (FileExplorer, useNotification, etc.)
5. ✅ Gradually migrate remaining code following 5-phase plan

**Questions?** Open an issue with the `autox` label on GitHub.

---

**Last Updated**: 2026-04-23  
**Maintainers**: ODH Dashboard Purple Team
