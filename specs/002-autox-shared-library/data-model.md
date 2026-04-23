# AutoX Shared Library - Data Model & Entities

**Date**: 2026-04-23  
**Phase**: 1 (Design)  
**Status**: Complete

---

## Overview

This document defines the key entities, interfaces, and data structures exported by the AutoX shared library package. All entities are designed to be domain-agnostic, with consumers (AutoML, AutoRAG) providing domain-specific implementations through handler files or strategy patterns.

---

## Entity Categories

### 1. BFF Layer Entities

#### 1.1 Core Models

**Purpose**: Shared data structures representing common domain concepts across AutoML and AutoRAG.

| Entity | Description | Extractable From | LOC |
|--------|-------------|------------------|-----|
| `HealthCheck` | Health check response model | 100% identical in both packages | 10 |
| `Namespace` | Kubernetes namespace representation | 100% identical | 21 |
| `User` | User identity model | 100% identical | 6 |
| `Groups` | User group models for RBAC | 100% identical | 50 |
| `RBACTypes` | Role-based access control types | 100% identical | 27 |
| `Secret` | Generic secret model | 100% identical | 23 |
| `S3Credentials` | S3 credential model | 100% identical | 30 |
| `PipelineServer` | Pipeline server configuration | 100% identical | 93 |
| `PipelineRun` | Base pipeline run model | 85% shared | 191 |

**Example Structure** (Go):

```go
// internal/models/namespace.go
package models

type Namespace struct {
    Name        string `json:"name"`
    DisplayName string `json:"displayName,omitempty"`
    Status      string `json:"status"`
}

// internal/models/pipeline_runs.go
package models

type PipelineRun struct {
    ID              string                 `json:"id"`
    Name            string                 `json:"name"`
    State           string                 `json:"state"`
    CreatedAt       string                 `json:"createdAt"`
    FinishedAt      string                 `json:"finishedAt,omitempty"`
    Parameters      map[string]interface{} `json:"parameters"`
    PipelineType    string                 `json:"pipelineType"`
    PipelineVersion string                 `json:"pipelineVersion"`
}
```

**Validation Rules**:
- All timestamps must be ISO 8601 format
- `State` must be one of: PENDING, RUNNING, SUCCEEDED, FAILED, CANCELED
- `Name` must match Kubernetes DNS-1123 subdomain (lowercase alphanumeric, hyphens)

**Relationships**:
- `PipelineRun` → `PipelineServer` (runs execute on a specific server)
- `User` → `Groups` (many-to-many)
- `S3Credentials` → `Secret` (S3 creds stored as Kubernetes secrets)

---

#### 1.2 Integration Interfaces

**Purpose**: Define contracts for external service interactions, allowing AutoML and AutoRAG to inject custom implementations.

##### Kubernetes Client Interface

```go
// internal/integrations/kubernetes/interface.go
package kubernetes

import (
    "context"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
    
    // Port forwarding
    PortForward(ctx context.Context, namespace, podName string, localPort, remotePort int) error
}
```

**Extraction**: 100% identical across AutoML and AutoRAG (`internal/integrations/kubernetes/interface.go`, 800 LOC total)

**Implementation Note**: AutoX provides a default implementation using `k8s.io/client-go`. Consumers can inject custom implementations for mocking or testing.

##### S3 Client Interface

```go
// internal/integrations/s3/interface.go
package s3

import (
    "context"
    "io"
)

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

type S3Object struct {
    Key          string
    Size         int64
    LastModified string
    ETag         string
}
```

**Extraction**: 95% identical (AutoML has additional `GetCSVSchema` method, which will be handled in AutoML handler)

**Validation Rules**:
- Bucket names must be DNS-compliant (lowercase, no underscores)
- Keys must not start with `/`
- Size must be non-negative

##### Pipeline Server Client Interface

```go
// internal/integrations/pipelineserver/interface.go
package pipelineserver

import (
    "context"
    "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/models"
)

type PipelineServerClientInterface interface {
    // Pipeline operations
    GetPipeline(ctx context.Context, pipelineID string) (*models.Pipeline, error)
    ListPipelines(ctx context.Context, filter string) ([]*models.Pipeline, error)
    UploadPipeline(ctx context.Context, name string, yamlContent []byte) (*models.Pipeline, error)
    
    // Pipeline run operations
    CreatePipelineRun(ctx context.Context, req *CreatePipelineRunRequest) (*models.PipelineRun, error)
    GetPipelineRun(ctx context.Context, runID string) (*models.PipelineRun, error)
    ListPipelineRuns(ctx context.Context, filter string) ([]*models.PipelineRun, error)
    TerminatePipelineRun(ctx context.Context, runID string) error
    RetryPipelineRun(ctx context.Context, runID string) (*models.PipelineRun, error)
    
    // Health
    HealthCheck(ctx context.Context) error
}

type CreatePipelineRunRequest struct {
    PipelineVersionID string                 `json:"pipelineVersionId"`
    DisplayName       string                 `json:"displayName"`
    Parameters        map[string]interface{} `json:"parameters"`
}
```

**Extraction**: 100% identical interface definition, 77% identical implementation (630 LOC)

**Validation Rules**:
- `DisplayName` must not be empty
- `PipelineVersionID` must be valid UUID
- `Parameters` must match pipeline definition schema (validation deferred to pipeline server)

---

#### 1.3 Repository Patterns

**Purpose**: Encapsulate business logic and data access patterns, providing base implementations with extension points.

##### Pipeline Repository

```go
// internal/repositories/pipeline.go
package repositories

import (
    "context"
    ps "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/integrations/pipelineserver"
)

type DiscoveredPipeline struct {
    ID              string
    VersionID       string
    Name            string
    DefinitionName  string // e.g., "time-series", "rag-pipeline"
}

type PipelineDefinition struct {
    Name     string
    YAMLPath string
}

type PipelineRepository interface {
    // Discovery and caching
    DiscoverNamedPipelines(
        client ps.PipelineServerClientInterface,
        ctx context.Context,
        namespace string,
        pipelineServerBaseURL string,
        definitions map[string]string,
    ) (map[string]*DiscoveredPipeline, error)
    
    EnsurePipeline(
        client ps.PipelineServerClientInterface,
        ctx context.Context,
        namespace string,
        pipelineServerBaseURL string,
        def PipelineDefinition,
    ) (*DiscoveredPipeline, error)
    
    InvalidateCache(pipelineServerBaseURL, namespace string)
}
```

**Extraction**: 97% identical (1,069 LOC total)

**Key Pattern**: AutoX provides caching infrastructure and discovery logic. AutoML and AutoRAG provide their own `definitions` map with domain-specific pipeline names.

**Example Usage** (Consumer):

```go
// AutoML handler
func (h *Handler) CreateTimeSeriesPipeline(ctx context.Context) {
    definitions := map[string]string{
        "time-series": "/pipelines/time_series_pipeline.yaml",
        "tabular":     "/pipelines/tabular_pipeline.yaml",
    }
    
    pipelines, err := h.pipelineRepo.DiscoverNamedPipelines(
        h.psClient, ctx, namespace, baseURL, definitions,
    )
    // Use discovered pipelines
}
```

##### Secret Repository

```go
// internal/repositories/secret.go
package repositories

import (
    "context"
    "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/integrations/kubernetes"
    "github.com/opendatahub-io/odh-dashboard/packages/autox/bff/internal/models"
)

type SecretRepository interface {
    GetFilteredSecrets(
        client kubernetes.KubernetesClientInterface,
        ctx context.Context,
        namespace string,
        identity *kubernetes.RequestIdentity,
        secretType string,
    ) ([]models.Secret, error)
}
```

**Extraction**: 95% shared (204-239 LOC)

**Customization Point**: AutoRAG has additional LlamaStack secret filtering. This is handled via handler file in AutoRAG, not in AutoX.

**Handler Pattern** (AutoRAG-specific):

```go
// AutoRAG handler
func (h *Handler) GetLlamaStackSecrets(ctx context.Context, namespace string) ([]models.Secret, error) {
    // Call shared repository for base filtering
    allSecrets, err := h.secretRepo.GetFilteredSecrets(h.k8sClient, ctx, namespace, identity, "all")
    if err != nil {
        return nil, err
    }
    
    // AutoRAG-specific filtering (not in AutoX)
    return filterLlamaStackSecrets(allSecrets), nil
}

func filterLlamaStackSecrets(secrets []models.Secret) []models.Secret {
    // AutoRAG-specific logic
}
```

---

#### 1.4 API Utilities

**Purpose**: Provide HTTP request/response handling, error formatting, middleware.

##### Error Response Models

```go
// internal/api/errors.go
package api

type ErrorResponse struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

type ErrorHandler interface {
    BadRequestResponse(w http.ResponseWriter, r *http.Request, err error)
    UnauthorizedResponse(w http.ResponseWriter, r *http.Request, message string)
    ForbiddenResponse(w http.ResponseWriter, r *http.Request, message string)
    NotFoundResponse(w http.ResponseWriter, r *http.Request, resource string)
    ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error)
}
```

**Extraction**: 99% identical (181 LOC)

##### Request/Response Helpers

```go
// internal/api/helpers.go
package api

import (
    "encoding/json"
    "net/http"
)

// WriteJSON writes a JSON response with status code and optional headers
func WriteJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
    output, err := json.Marshal(data)
    if err != nil {
        return err
    }
    
    for key, value := range headers {
        w.Header()[key] = value
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    w.Write(output)
    return nil
}

// ReadJSON reads and parses JSON request body
func ReadJSON(r *http.Request, dst interface{}) error {
    decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields()
    return decoder.Decode(dst)
}
```

**Extraction**: 100% identical (108 LOC)

---

### 2. Frontend Layer Entities

#### 2.1 Core Types

**Purpose**: TypeScript types representing domain concepts shared across AutoML and AutoRAG.

##### Pipeline Types

```typescript
// frontend/src/types/pipeline.ts

export enum RuntimeStateKF {
  PENDING = 'PENDING',
  RUNNING = 'RUNNING',
  SUCCEEDED = 'SUCCEEDED',
  FAILED = 'FAILED',
  CANCELED = 'CANCELED',
  SKIPPED = 'SKIPPED',
}

export type PipelineRun = {
  id: string;
  name: string;
  state: RuntimeStateKF;
  createdAt: string;
  finishedAt?: string;
  parameters: Record<string, unknown>;
  pipelineType: string;
  pipelineVersion: string;
  error?: {
    message: string;
    code?: string;
  };
};

export type PipelineDefinition = {
  name: string;
  displayName: string;
  description?: string;
  yamlPath: string;
};

export type DiscoveredPipeline = {
  id: string;
  versionId: string;
  name: string;
  definitionName: string;
};
```

**Extraction**: 100% identical (186 LOC)

**Validation Rules**:
- `id` must be non-empty UUID string
- `state` must be valid `RuntimeStateKF` enum value
- `createdAt` must be ISO 8601 timestamp
- `parameters` keys must match pipeline definition

##### S3 Types

```typescript
// frontend/src/types/s3.ts

export type S3Credentials = {
  accessKeyId: string;
  secretAccessKey: string;
  region: string;
  endpoint?: string;
};

export type S3Object = {
  key: string;
  size: number;
  lastModified: string;
  etag: string;
};

export type S3ListResponse = {
  objects: S3Object[];
  prefix: string;
  isTruncated: boolean;
  nextContinuationToken?: string;
};
```

**Extraction**: 100% identical (20 LOC)

##### Topology Types

```typescript
// frontend/src/types/topology.ts

export type TopologyNode = {
  id: string;
  label: string;
  type: 'task' | 'artifact';
  status?: RuntimeStateKF;
  data?: Record<string, unknown>;
};

export type TopologyEdge = {
  id: string;
  source: string;
  target: string;
};

export type TopologyGraph = {
  nodes: TopologyNode[];
  edges: TopologyEdge[];
};
```

**Extraction**: 100% identical (46 LOC)

---

#### 2.2 Shared Hooks

**Purpose**: Encapsulate data fetching and state management logic using React Query.

##### usePipelineRuns

```typescript
// frontend/src/hooks/usePipelineRuns.ts

import { useQuery, useMutation } from '@tanstack/react-query';
import { PipelineRun } from '../types/pipeline';

export type UsePipelineRunsOptions = {
  namespace?: string;
  pipelineVersionId?: string;
  enabled?: boolean;
};

export function usePipelineRuns(options: UsePipelineRunsOptions = {}) {
  const { namespace, pipelineVersionId, enabled = true } = options;
  
  return useQuery({
    queryKey: ['pipelineRuns', namespace, pipelineVersionId],
    queryFn: async ({ signal }) => {
      const params = new URLSearchParams();
      if (pipelineVersionId) params.set('pipelineVersionId', pipelineVersionId);
      
      const response = await fetch(
        `/api/pipelines/runs?namespace=${namespace}&${params}`,
        { signal }
      );
      
      if (!response.ok) throw new Error('Failed to fetch pipeline runs');
      return response.json() as Promise<PipelineRun[]>;
    },
    enabled: enabled && !!namespace,
  });
}

export function useTerminatePipelineRun(namespace: string) {
  return useMutation({
    mutationFn: async (runId: string) => {
      const response = await fetch(
        `/api/pipelines/runs/${runId}/terminate?namespace=${namespace}`,
        { method: 'POST' }
      );
      if (!response.ok) throw new Error('Failed to terminate run');
    },
  });
}

export function useRetryPipelineRun(namespace: string) {
  return useMutation({
    mutationFn: async (runId: string) => {
      const response = await fetch(
        `/api/pipelines/runs/${runId}/retry?namespace=${namespace}`,
        { method: 'POST' }
      );
      if (!response.ok) throw new Error('Failed to retry run');
      return response.json() as Promise<PipelineRun>;
    },
  });
}
```

**Extraction**: 100% identical (208 LOC)

**Parameterization**: Uses namespace and filter parameters to support both AutoML and AutoRAG use cases without factories.

##### useNotification

```typescript
// frontend/src/hooks/useNotification.ts

import { useCallback } from 'react';
import { useToast } from '@odh-dashboard/plugin-core';

export type NotificationVariant = 'success' | 'danger' | 'warning' | 'info';

export type NotificationOptions = {
  title: string;
  message?: string;
  variant?: NotificationVariant;
  autoHideDuration?: number;
};

export function useNotification() {
  const { addToast } = useToast();
  
  const notify = useCallback((options: NotificationOptions) => {
    const { title, message, variant = 'info', autoHideDuration = 5000 } = options;
    
    addToast({
      title,
      body: message,
      variant,
      autoHideDuration,
    });
  }, [addToast]);
  
  return { notify };
}
```

**Extraction**: 100% identical (167 LOC including implementation)

##### S3 API Hooks

```typescript
// frontend/src/hooks/useS3.ts

import { useQuery, useMutation } from '@tanstack/react-query';
import { z } from 'zod';
import type { S3Object, S3ListResponse } from '../types/s3';

const S3ObjectSchema = z.object({
  key: z.string(),
  size: z.number().nonnegative(),
  lastModified: z.string(),
  etag: z.string(),
});

const S3ListResponseSchema = z.object({
  objects: z.array(S3ObjectSchema),
  prefix: z.string(),
  isTruncated: z.boolean(),
  nextContinuationToken: z.string().optional(),
});

export type UseS3ListFilesOptions = {
  namespace: string;
  bucket: string;
  prefix?: string;
  recursive?: boolean;
  enabled?: boolean;
};

export function useS3ListFiles(options: UseS3ListFilesOptions) {
  const { namespace, bucket, prefix = '', recursive = false, enabled = true } = options;
  
  return useQuery({
    queryKey: ['s3Files', namespace, bucket, prefix, recursive],
    queryFn: async ({ signal }) => {
      const params = new URLSearchParams({
        namespace,
        bucket,
        prefix,
        recursive: String(recursive),
      });
      
      const response = await fetch(`/api/s3/list?${params}`, { signal });
      if (!response.ok) throw new Error('Failed to list S3 files');
      
      const data = await response.json();
      return S3ListResponseSchema.parse(data);
    },
    enabled: enabled && !!namespace && !!bucket,
  });
}

export function useS3FileUpload(namespace: string) {
  return useMutation({
    mutationFn: async ({ bucket, key, file }: { bucket: string; key: string; file: File }) => {
      const formData = new FormData();
      formData.append('file', file);
      formData.append('namespace', namespace);
      formData.append('bucket', bucket);
      formData.append('key', key);
      
      const response = await fetch('/api/s3/upload', {
        method: 'POST',
        body: formData,
      });
      
      if (!response.ok) throw new Error('Failed to upload file');
    },
  });
}
```

**Extraction**: 100% identical (175 LOC for S3 API file + hook wrappers)

**Validation**: Uses Zod for runtime validation of API responses.

---

#### 2.3 Shared Components

**Purpose**: Provide composable UI primitives that AutoML and AutoRAG can use as building blocks.

##### FileExplorer Component

```typescript
// frontend/src/components/FileExplorer/FileExplorer.tsx

import React from 'react';
import {
  DataList,
  DataListItem,
  DataListCell,
  Button,
  Breadcrumb,
  BreadcrumbItem,
  Pagination,
  SearchInput,
} from '@patternfly/react-core';
import { FolderIcon, FileIcon } from '@patternfly/react-icons';

export type FileItem = {
  name: string;
  type: 'file' | 'directory';
  size?: number;
  lastModified?: string;
};

export type FileExplorerProps = {
  files: FileItem[];
  currentPath: string;
  onNavigate: (path: string) => void;
  onFileSelect?: (file: FileItem) => void;
  isLoading?: boolean;
  pagination?: {
    page: number;
    perPage: number;
    total: number;
    onPageChange: (page: number) => void;
    onPerPageChange: (perPage: number) => void;
  };
  searchable?: boolean;
};

export const FileExplorer: React.FC<FileExplorerProps> = ({
  files,
  currentPath,
  onNavigate,
  onFileSelect,
  isLoading = false,
  pagination,
  searchable = false,
}) => {
  const [searchTerm, setSearchTerm] = React.useState('');
  
  const filteredFiles = React.useMemo(() => {
    if (!searchable || !searchTerm) return files;
    return files.filter(f => f.name.toLowerCase().includes(searchTerm.toLowerCase()));
  }, [files, searchTerm, searchable]);
  
  const pathSegments = currentPath.split('/').filter(Boolean);
  
  return (
    <div className="autox-file-explorer" data-testid="file-explorer">
      <Breadcrumb>
        <BreadcrumbItem onClick={() => onNavigate('/')}>Root</BreadcrumbItem>
        {pathSegments.map((segment, i) => {
          const path = '/' + pathSegments.slice(0, i + 1).join('/');
          return (
            <BreadcrumbItem key={path} onClick={() => onNavigate(path)}>
              {segment}
            </BreadcrumbItem>
          );
        })}
      </Breadcrumb>
      
      {searchable && (
        <SearchInput
          value={searchTerm}
          onChange={(_, value) => setSearchTerm(value)}
          onClear={() => setSearchTerm('')}
          placeholder="Search files..."
        />
      )}
      
      <DataList aria-label="File list" isCompact>
        {filteredFiles.map(file => (
          <DataListItem key={file.name} onClick={() => {
            if (file.type === 'directory') {
              onNavigate(`${currentPath}/${file.name}`);
            } else if (onFileSelect) {
              onFileSelect(file);
            }
          }}>
            <DataListCell>
              {file.type === 'directory' ? <FolderIcon /> : <FileIcon />}
              {file.name}
            </DataListCell>
            {file.size && <DataListCell>{formatBytes(file.size)}</DataListCell>}
            {file.lastModified && <DataListCell>{file.lastModified}</DataListCell>}
          </DataListItem>
        ))}
      </DataList>
      
      {pagination && (
        <Pagination
          page={pagination.page}
          perPage={pagination.perPage}
          itemCount={pagination.total}
          onSetPage={(_, page) => pagination.onPageChange(page)}
          onPerPageSelect={(_, perPage) => pagination.onPerPageChange(perPage)}
        />
      )}
    </div>
  );
};

function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}
```

**Extraction**: 100% identical (1,270 LOC total including sub-components)

**Composition Note**: AutoML and AutoRAG use this primitive directly with their own S3 data sources. No prop variants needed.

##### PipelineRunsTable Component

```typescript
// frontend/src/components/PipelineRunsTable/PipelineRunsTable.tsx

import React from 'react';
import { Table, Thead, Tr, Th, Tbody, Td } from '@patternfly/react-table';
import { Timestamp } from '@patternfly/react-core';
import { RuntimeStateKF, PipelineRun } from '../../types/pipeline';

export type PipelineRunsTableProps = {
  runs: PipelineRun[];
  onRunSelect?: (run: PipelineRun) => void;
  showActions?: boolean;
  onTerminate?: (runId: string) => void;
  onRetry?: (runId: string) => void;
};

export const PipelineRunsTable: React.FC<PipelineRunsTableProps> = ({
  runs,
  onRunSelect,
  showActions = false,
  onTerminate,
  onRetry,
}) => {
  return (
    <Table aria-label="Pipeline runs table" variant="compact">
      <Thead>
        <Tr>
          <Th>Name</Th>
          <Th>Status</Th>
          <Th>Created</Th>
          <Th>Finished</Th>
          {showActions && <Th>Actions</Th>}
        </Tr>
      </Thead>
      <Tbody>
        {runs.map(run => (
          <Tr key={run.id} onClick={() => onRunSelect?.(run)}>
            <Td>{run.name}</Td>
            <Td>{getStatusBadge(run.state)}</Td>
            <Td><Timestamp date={new Date(run.createdAt)} /></Td>
            <Td>{run.finishedAt ? <Timestamp date={new Date(run.finishedAt)} /> : '-'}</Td>
            {showActions && (
              <Td>
                {run.state === RuntimeStateKF.RUNNING && onTerminate && (
                  <Button variant="link" onClick={() => onTerminate(run.id)}>
                    Terminate
                  </Button>
                )}
                {run.state === RuntimeStateKF.FAILED && onRetry && (
                  <Button variant="link" onClick={() => onRetry(run.id)}>
                    Retry
                  </Button>
                )}
              </Td>
            )}
          </Tr>
        ))}
      </Tbody>
    </Table>
  );
};

function getStatusBadge(state: RuntimeStateKF): React.ReactNode {
  // Status badge rendering logic
}
```

**Extraction**: 98% identical (123 LOC for runs table)

---

#### 2.4 Shared Utilities

**Purpose**: Pure functions for data transformation, validation, formatting.

##### Pipeline Run State Utilities

```typescript
// frontend/src/utils/pipelineRunState.ts

import { RuntimeStateKF } from '../types/pipeline';

export function isPipelineRunning(state: RuntimeStateKF): boolean {
  return state === RuntimeStateKF.RUNNING || state === RuntimeStateKF.PENDING;
}

export function isPipelineRunComplete(state: RuntimeStateKF): boolean {
  return state === RuntimeStateKF.SUCCEEDED || 
         state === RuntimeStateKF.FAILED || 
         state === RuntimeStateKF.CANCELED;
}

export function isPipelineRunSuccessful(state: RuntimeStateKF): boolean {
  return state === RuntimeStateKF.SUCCEEDED;
}
```

**Extraction**: 100% identical (25 LOC)

##### Validation Utilities

```typescript
// frontend/src/utils/validation.ts

export function isValidK8sName(name: string): boolean {
  // RFC 1123 DNS subdomain validation
  const regex = /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/;
  return regex.test(name) && name.length <= 63;
}

export function isValidS3BucketName(name: string): boolean {
  const regex = /^[a-z0-9][a-z0-9-]*[a-z0-9]$/;
  return regex.test(name) && name.length >= 3 && name.length <= 63;
}

export function isEmptyValue(value: unknown): boolean {
  return value === null || value === undefined || value === '';
}
```

**Extraction**: 100% identical (17 LOC)

##### Topology Node Creation

```typescript
// frontend/src/utils/topology.ts

import type { TopologyNode, TopologyEdge, RuntimeStateKF } from '../types';

export function createTopologyNode(
  id: string,
  label: string,
  type: 'task' | 'artifact',
  status?: RuntimeStateKF,
  data?: Record<string, unknown>,
): TopologyNode {
  return { id, label, type, status, data };
}

export function createTopologyEdge(
  source: string,
  target: string,
): TopologyEdge {
  return { id: `${source}-${target}`, source, target };
}
```

**Extraction**: 100% identical (39 LOC)

---

## Relationships Between Entities

### BFF Layer

```
PipelineRepository
  ├── depends on → PipelineServerClientInterface (integration)
  ├── returns → DiscoveredPipeline (model)
  └── caches → Pipeline definitions

SecretRepository
  ├── depends on → KubernetesClientInterface (integration)
  ├── returns → Secret (model)
  └── filters by → SecretType (enum)

S3ClientInterface
  ├── implements → AWS SDK wrapper
  ├── returns → S3Object (model)
  └── uses → S3Credentials (model)

PipelineRunRepository
  ├── depends on → PipelineServerClientInterface
  ├── creates → PipelineRun (model)
  └── validates → Parameters (domain-specific)
```

### Frontend Layer

```
usePipelineRuns (hook)
  ├── fetches → PipelineRun[] (type)
  ├── uses → React Query
  └── returns → Query state

FileExplorer (component)
  ├── accepts → FileItem[] (type)
  ├── renders → PatternFly DataList
  └── emits → onNavigate, onFileSelect (callbacks)

useNotification (hook)
  ├── wraps → ODH Toast system
  ├── accepts → NotificationOptions (type)
  └── provides → notify() function

Validation utilities
  ├── validates → K8s names, S3 bucket names
  └── used by → Form components in consumers
```

---

## State Transitions

### Pipeline Run States

```
PENDING → RUNNING → SUCCEEDED
                 → FAILED
                 → CANCELED (via terminate)

FAILED → (retry) → PENDING
```

**Validation**: 
- Only RUNNING runs can be terminated
- Only FAILED runs can be retried
- State transitions are managed by Pipeline Server, not AutoX

### File Explorer Navigation

```
Root (/) → Directory (click) → Subdirectory (click) → ...
                             → File (click) → onFileSelect()

Breadcrumb (click any segment) → Navigate to that path
```

---

## Validation Rules Summary

### BFF Layer

| Entity | Validation Rule | Error Handling |
|--------|-----------------|----------------|
| Namespace name | RFC 1123 DNS subdomain | 400 Bad Request |
| Pipeline run name | Non-empty, max 253 chars | 400 Bad Request |
| S3 bucket name | DNS-compliant, 3-63 chars | 400 Bad Request |
| S3 object key | Non-empty, no leading `/` | 400 Bad Request |
| User groups | Must be array of strings | 400 Bad Request |

### Frontend Layer

| Entity | Validation Rule | UI Feedback |
|--------|-----------------|-------------|
| K8s name input | `isValidK8sName()` | Inline error message |
| S3 bucket input | `isValidS3BucketName()` | Inline error message |
| File upload size | Max size from config | Toast notification |
| Pipeline parameters | Schema from definition | Form validation errors |

---

## Extension Points

### Where Consumers Extend AutoX

1. **Handler Files** (20-80% code divergence):
   - AutoML: Time-series specific validation, tabular pipeline creation
   - AutoRAG: RAG pipeline creation, LlamaStack secret filtering
   - Pattern: Import AutoX utilities, add domain logic

2. **Strategy Pattern** (advanced customization):
   - Pipeline parameter builders (AutoML vs AutoRAG have different schemas)
   - Custom validation rules beyond base validation
   - Pattern: Implement interface, inject at runtime

3. **Component Composition**:
   - AutoML: Composes `FileExplorer` + domain-specific file actions
   - AutoRAG: Composes `PipelineRunsTable` + RAG-specific columns
   - Pattern: Import AutoX primitive, wrap with domain logic

---

## Migration Strategy

### Phase 1: Perfect Duplicates (Low Risk)
- Extract: FileExplorer, useNotification, S3 API, K8s client, pipeline models
- Impact: ~3,000 LOC, zero refactoring needed

### Phase 2-5: Gradual Extraction
- See `research.md` for detailed 5-phase roadmap
- Total: ~31,500 LOC over 8-12 weeks

---

**Next**: See `contracts/` for detailed API specifications and `quickstart.md` for developer setup guide.
