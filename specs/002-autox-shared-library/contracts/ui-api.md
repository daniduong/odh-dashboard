# AutoX UI API Contract

**Package**: `@odh-dashboard/autox`  
**Language**: TypeScript 5.x  
**Runtime**: React 18  
**Status**: Draft  
**Version**: 1.0.0-alpha  
**Date**: 2026-04-23

---

## Overview

This document defines the public API contract for the AutoX UI (frontend) shared library. All hooks, components, types, and utilities exported by AutoX UI are guaranteed to maintain backward compatibility within the same major version. Breaking changes will only occur with major version bumps.

**Target Consumers**: AutoML Frontend, AutoRAG Frontend (monorepo-internal only)

**Import Path Convention**:

```typescript
// Hooks
import { usePipelineRuns, useNotification, useS3ListFiles } from '@odh-dashboard/autox/hooks';

// Components
import { FileExplorer, PipelineRunsTable } from '@odh-dashboard/autox/components';

// Types
import type { PipelineRun, S3Object, RuntimeStateKF } from '@odh-dashboard/autox/types';

// Utilities
import { isPipelineRunning, isValidK8sName } from '@odh-dashboard/autox/utils';
```

**Module Federation**: AutoX is configured as a **shared dependency** (not a remote) with `singleton: true` to ensure single runtime instance.

---

## API Stability Guarantees

### Stable API (Backward Compatible)

- **Hook signatures**: Function parameters, return types, React Query keys
- **Component props**: Prop names, types, required/optional status
- **Type definitions**: Interface fields, enum values
- **Utility function signatures**: Parameters, return types

### Unstable API (May Change)

- **Internal helpers**: Non-exported functions
- **Component internal state**: Implementation details not exposed via props
- **CSS class names**: Styling internals (use PatternFly design tokens instead)

### Breaking Changes Policy

- **Synchronized versioning**: AutoX, AutoML, and AutoRAG are always updated together in the same PR
- **No independent versioning**: Breaking changes are coordinated and applied atomically
- **CI type-checking**: All packages are type-checked together (`tsc --noEmit`) to catch breaking changes before merge
- **Migration period**: N/A (synchronized updates mean no migration period needed)

---

## 1. Types (`@odh-dashboard/autox/types`)

### 1.1 Pipeline Types

```typescript
// RuntimeStateKF enum
export enum RuntimeStateKF {
  PENDING = 'PENDING',
  RUNNING = 'RUNNING',
  SUCCEEDED = 'SUCCEEDED',
  FAILED = 'FAILED',
  CANCELED = 'CANCELED',
  SKIPPED = 'SKIPPED',
}

// PipelineRun type
export type PipelineRun = {
  id: string;
  name: string;
  state: RuntimeStateKF;
  createdAt: string; // ISO 8601 timestamp
  finishedAt?: string; // ISO 8601 timestamp
  parameters: Record<string, unknown>;
  pipelineType: string;
  pipelineVersion: string;
  error?: {
    message: string;
    code?: string;
  };
};

// PipelineDefinition type
export type PipelineDefinition = {
  name: string;
  displayName: string;
  description?: string;
  yamlPath: string;
};

// DiscoveredPipeline type
export type DiscoveredPipeline = {
  id: string;
  versionId: string;
  name: string;
  definitionName: string;
};
```

**Stability**: ✅ Stable  
**Extractable**: 100% identical (186 LOC)

**Validation**:
- `id` must be non-empty UUID string
- `state` must be valid `RuntimeStateKF` enum value
- `createdAt`/`finishedAt` must be ISO 8601 timestamps
- `parameters` keys must match pipeline definition

---

### 1.2 S3 Types

```typescript
// S3Credentials type
export type S3Credentials = {
  accessKeyId: string;
  secretAccessKey: string;
  region: string;
  endpoint?: string;
};

// S3Object type
export type S3Object = {
  key: string;
  size: number;
  lastModified: string; // ISO 8601 timestamp
  etag: string;
};

// S3ListResponse type
export type S3ListResponse = {
  objects: S3Object[];
  prefix: string;
  isTruncated: boolean;
  nextContinuationToken?: string;
};
```

**Stability**: ✅ Stable  
**Extractable**: 100% identical (20 LOC)

**Validation**:
- `accessKeyId` and `secretAccessKey` must not be empty
- `size` must be non-negative
- `key` must not start with `/`

---

### 1.3 Topology Types

```typescript
// TopologyNode type
export type TopologyNode = {
  id: string;
  label: string;
  type: 'task' | 'artifact';
  status?: RuntimeStateKF;
  data?: Record<string, unknown>;
};

// TopologyEdge type
export type TopologyEdge = {
  id: string;
  source: string;
  target: string;
};

// TopologyGraph type
export type TopologyGraph = {
  nodes: TopologyNode[];
  edges: TopologyEdge[];
};
```

**Stability**: ✅ Stable  
**Extractable**: 100% identical (46 LOC)

**Usage**: Pipeline DAG rendering, visualization

---

### 1.4 Common Types

```typescript
// Namespace type
export type Namespace = {
  name: string;
  displayName?: string;
  status: string;
};

// Secret type
export type Secret = {
  name: string;
  namespace: string;
  type: string;
  keys: string[];
  labels?: Record<string, string>;
  annotations?: Record<string, string>;
};

// User type
export type User = {
  id: string;
  username: string;
  email?: string;
  groups: string[];
};
```

**Stability**: ✅ Stable

---

## 2. Hooks (`@odh-dashboard/autox/hooks`)

### 2.1 usePipelineRuns

```typescript
import type { PipelineRun } from '../types';
import type { UseQueryResult } from '@tanstack/react-query';

export type UsePipelineRunsOptions = {
  namespace?: string;
  pipelineVersionId?: string;
  enabled?: boolean;
};

export function usePipelineRuns(
  options?: UsePipelineRunsOptions
): UseQueryResult<PipelineRun[], Error>;
```

**Stability**: ✅ Stable  
**Extractable**: 100% identical (104 LOC implementation)

**React Query Key**: `['pipelineRuns', namespace, pipelineVersionId]`

**Usage Example**:

```typescript
// AutoML component
import { usePipelineRuns } from '@odh-dashboard/autox/hooks';

function MyComponent() {
  const { data: runs, isLoading, error } = usePipelineRuns({
    namespace: 'my-namespace',
    pipelineVersionId: 'abc-123',
    enabled: true,
  });

  if (isLoading) return <Spinner />;
  if (error) return <Alert variant="danger">{error.message}</Alert>;
  
  return <PipelineRunsList runs={runs} />;
}
```

**Return Type**:
- `data`: `PipelineRun[]` or `undefined`
- `isLoading`: `boolean`
- `isFetching`: `boolean`
- `error`: `Error` or `null`
- `refetch`: `() => void`

**Behavior**:
- Auto-fetches when `enabled: true` and `namespace` is provided
- Refetches on window focus (React Query default)
- Caches results for 5 minutes
- Retries failed requests 3 times

---

### 2.2 usePipelineRunQuery

```typescript
export type UsePipelineRunQueryOptions = {
  runId?: string;
  namespace?: string;
  enabled?: boolean;
  pollingInterval?: number; // milliseconds
};

export function usePipelineRunQuery(
  options: UsePipelineRunQueryOptions
): UseQueryResult<PipelineRun, Error>;
```

**Stability**: ✅ Stable  
**Extractable**: 100% identical (36 LOC)

**React Query Key**: `['pipelineRun', runId, namespace]`

**Usage Example**:

```typescript
import { usePipelineRunQuery } from '@odh-dashboard/autox/hooks';

function RunDetails({ runId, namespace }: { runId: string; namespace: string }) {
  const { data: run } = usePipelineRunQuery({
    runId,
    namespace,
    pollingInterval: 5000, // Poll every 5 seconds while running
  });

  return <div>Status: {run?.state}</div>;
}
```

**Behavior**:
- Polls at specified interval if run is in RUNNING or PENDING state
- Stops polling when run reaches terminal state (SUCCEEDED, FAILED, CANCELED)

---

### 2.3 useTerminatePipelineRun

```typescript
import type { UseMutationResult } from '@tanstack/react-query';

export function useTerminatePipelineRun(
  namespace: string
): UseMutationResult<void, Error, string>;
```

**Stability**: ✅ Stable  
**Extractable**: 100% identical (22 LOC)

**Mutation Key**: `['terminatePipelineRun']`

**Usage Example**:

```typescript
import { useTerminatePipelineRun } from '@odh-dashboard/autox/hooks';

function TerminateButton({ runId, namespace }: { runId: string; namespace: string }) {
  const { mutate: terminate, isPending } = useTerminatePipelineRun(namespace);

  return (
    <Button
      onClick={() => terminate(runId)}
      isDisabled={isPending}
    >
      Terminate
    </Button>
  );
}
```

**Return Type**:
- `mutate`: `(runId: string) => void`
- `isPending`: `boolean`
- `error`: `Error` or `null`

---

### 2.4 useRetryPipelineRun

```typescript
export function useRetryPipelineRun(
  namespace: string
): UseMutationResult<PipelineRun, Error, string>;
```

**Stability**: ✅ Stable  
**Extractable**: 100% identical (22 LOC)

**Usage Example**:

```typescript
import { useRetryPipelineRun } from '@odh-dashboard/autox/hooks';

function RetryButton({ runId, namespace }: { runId: string; namespace: string }) {
  const { mutate: retry, isPending } = useRetryPipelineRun(namespace);

  return (
    <Button
      onClick={() => retry(runId, {
        onSuccess: (newRun) => console.log('Created new run:', newRun.id),
      })}
      isDisabled={isPending}
    >
      Retry
    </Button>
  );
}
```

---

### 2.5 useNotification

```typescript
export type NotificationVariant = 'success' | 'danger' | 'warning' | 'info';

export type NotificationOptions = {
  title: string;
  message?: string;
  variant?: NotificationVariant;
  autoHideDuration?: number;
};

export function useNotification(): {
  notify: (options: NotificationOptions) => void;
};
```

**Stability**: ✅ Stable  
**Extractable**: 100% identical (167 LOC including implementation)

**Usage Example**:

```typescript
import { useNotification } from '@odh-dashboard/autox/hooks';

function MyForm() {
  const { notify } = useNotification();

  const handleSubmit = async () => {
    try {
      await submitForm();
      notify({
        title: 'Success',
        message: 'Form submitted successfully',
        variant: 'success',
      });
    } catch (err) {
      notify({
        title: 'Error',
        message: err.message,
        variant: 'danger',
      });
    }
  };

  return <Button onClick={handleSubmit}>Submit</Button>;
}
```

**Behavior**:
- Displays toast notification in top-right corner
- Auto-hides after `autoHideDuration` (default: 5000ms)
- Stacks multiple notifications

---

### 2.6 useS3ListFiles

```typescript
export type UseS3ListFilesOptions = {
  namespace: string;
  bucket: string;
  prefix?: string;
  recursive?: boolean;
  enabled?: boolean;
};

export function useS3ListFiles(
  options: UseS3ListFilesOptions
): UseQueryResult<S3ListResponse, Error>;
```

**Stability**: ✅ Stable  
**Extractable**: 95% identical (46 LOC)

**React Query Key**: `['s3Files', namespace, bucket, prefix, recursive]`

**Usage Example**:

```typescript
import { useS3ListFiles } from '@odh-dashboard/autox/hooks';

function S3FileBrowser({ namespace, bucket }: { namespace: string; bucket: string }) {
  const [prefix, setPrefix] = React.useState('');
  
  const { data, isLoading } = useS3ListFiles({
    namespace,
    bucket,
    prefix,
    recursive: false,
  });

  return (
    <FileExplorer
      files={data?.objects || []}
      currentPath={prefix}
      onNavigate={setPrefix}
    />
  );
}
```

**Validation**: Uses Zod schema to validate API response structure

---

### 2.7 useS3FileUpload

```typescript
export type S3UploadVariables = {
  bucket: string;
  key: string;
  file: File;
};

export function useS3FileUpload(
  namespace: string
): UseMutationResult<void, Error, S3UploadVariables>;
```

**Stability**: ✅ Stable  
**Extractable**: 100% identical (22 LOC)

**Usage Example**:

```typescript
import { useS3FileUpload } from '@odh-dashboard/autox/hooks';

function FileUploader({ namespace, bucket }: { namespace: string; bucket: string }) {
  const { mutate: upload, isPending } = useS3FileUpload(namespace);

  const handleFileChange = (file: File) => {
    upload({
      bucket,
      key: `uploads/${file.name}`,
      file,
    }, {
      onSuccess: () => notify({ title: 'Upload complete', variant: 'success' }),
      onError: (err) => notify({ title: 'Upload failed', message: err.message, variant: 'danger' }),
    });
  };

  return <FileUpload onChange={handleFileChange} isLoading={isPending} />;
}
```

---

### 2.8 useNamespaces

```typescript
export function useNamespaces(): UseQueryResult<Namespace[], Error>;
```

**Stability**: ✅ Stable  
**Extractable**: 100% identical (15 LOC)

**React Query Key**: `['namespaces']`

**Usage Example**:

```typescript
import { useNamespaces } from '@odh-dashboard/autox/hooks';

function NamespaceSelector() {
  const { data: namespaces } = useNamespaces();

  return (
    <Select>
      {namespaces?.map(ns => (
        <SelectOption key={ns.name} value={ns.name}>
          {ns.displayName || ns.name}
        </SelectOption>
      ))}
    </Select>
  );
}
```

---

## 3. Components (`@odh-dashboard/autox/components`)

### 3.1 FileExplorer

```typescript
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

export const FileExplorer: React.FC<FileExplorerProps>;
```

**Stability**: ✅ Stable  
**Extractable**: 100% identical (1,270 LOC total with sub-components)

**Features**:
- Breadcrumb navigation
- Search/filter (optional)
- Pagination (optional)
- Click handlers for files and directories
- Loading state support

**Usage Example**:

```typescript
import { FileExplorer } from '@odh-dashboard/autox/components';
import { useS3ListFiles } from '@odh-dashboard/autox/hooks';

function S3Browser({ namespace, bucket }: { namespace: string; bucket: string }) {
  const [currentPath, setCurrentPath] = React.useState('/');
  const { data, isLoading } = useS3ListFiles({ namespace, bucket, prefix: currentPath });

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
      onFileSelect={(file) => console.log('Selected:', file.name)}
      isLoading={isLoading}
      searchable
    />
  );
}
```

**Accessibility**: 
- Fully keyboard navigable
- ARIA labels on all interactive elements
- Semantic HTML structure

**Styling**: Uses PatternFly DataList component with `autox-file-explorer` class prefix

---

### 3.2 PipelineRunsTable

```typescript
export type PipelineRunsTableProps = {
  runs: PipelineRun[];
  onRunSelect?: (run: PipelineRun) => void;
  showActions?: boolean;
  onTerminate?: (runId: string) => void;
  onRetry?: (runId: string) => void;
};

export const PipelineRunsTable: React.FC<PipelineRunsTableProps>;
```

**Stability**: ✅ Stable  
**Extractable**: 98% identical (123 LOC)

**Features**:
- Sortable columns (Name, Status, Created, Finished)
- Row click handler (optional)
- Action buttons (Terminate, Retry) based on run state
- Status badge with color coding

**Usage Example**:

```typescript
import { PipelineRunsTable } from '@odh-dashboard/autox/components';
import { usePipelineRuns, useTerminatePipelineRun } from '@odh-dashboard/autox/hooks';

function RunsList({ namespace }: { namespace: string }) {
  const { data: runs } = usePipelineRuns({ namespace });
  const { mutate: terminate } = useTerminatePipelineRun(namespace);

  return (
    <PipelineRunsTable
      runs={runs || []}
      onRunSelect={(run) => navigate(`/runs/${run.id}`)}
      showActions
      onTerminate={terminate}
    />
  );
}
```

**Composition Note**: Consumers can extend this with domain-specific columns by wrapping or using as reference.

---

### 3.3 PipelineEmptyState

```typescript
export type PipelineEmptyStateProps = {
  title: string;
  message: string;
  actionLabel?: string;
  onAction?: () => void;
  icon?: React.ComponentType;
};

export const PipelineEmptyState: React.FC<PipelineEmptyStateProps>;
```

**Stability**: ✅ Stable  
**Extractable**: 95% identical (37 LOC)

**Usage Example**:

```typescript
import { PipelineEmptyState } from '@odh-dashboard/autox/components';
import { SearchIcon } from '@patternfly/react-icons';

function NoPipelinesFound() {
  return (
    <PipelineEmptyState
      title="No pipeline runs found"
      message="Start a new pipeline run to see results here"
      actionLabel="Create Run"
      onAction={() => navigate('/create')}
      icon={SearchIcon}
    />
  );
}
```

---

### 3.4 ToastNotification

```typescript
export type ToastNotificationProps = {
  title: string;
  message?: string;
  variant?: 'success' | 'danger' | 'warning' | 'info';
  autoHideDuration?: number;
  onClose?: () => void;
};

export const ToastNotification: React.FC<ToastNotificationProps>;
```

**Stability**: ✅ Stable  
**Extractable**: 100% identical (107 LOC)

**Usage**: Internal to `useNotification` hook. Can be used directly if needed.

**Note**: Consumers typically use `useNotification()` hook instead of rendering this component directly.

---

## 4. Utilities (`@odh-dashboard/autox/utils`)

### 4.1 Pipeline State Utilities

```typescript
import type { RuntimeStateKF } from '../types';

export function isPipelineRunning(state: RuntimeStateKF): boolean;
export function isPipelineRunComplete(state: RuntimeStateKF): boolean;
export function isPipelineRunSuccessful(state: RuntimeStateKF): boolean;
export function isPipelineRunFailed(state: RuntimeStateKF): boolean;
```

**Stability**: ✅ Stable  
**Extractable**: 100% identical (25 LOC)

**Usage Example**:

```typescript
import { isPipelineRunning, isPipelineRunSuccessful } from '@odh-dashboard/autox/utils';

function RunStatus({ run }: { run: PipelineRun }) {
  if (isPipelineRunning(run.state)) {
    return <Spinner size="sm" />;
  }
  
  if (isPipelineRunSuccessful(run.state)) {
    return <CheckCircleIcon color="green" />;
  }
  
  return <ExclamationCircleIcon color="red" />;
}
```

---

### 4.2 Validation Utilities

```typescript
export function isValidK8sName(name: string): boolean;
export function isValidS3BucketName(name: string): boolean;
export function isEmptyValue(value: unknown): boolean;
```

**Stability**: ✅ Stable  
**Extractable**: 100% identical (17 LOC)

**Usage Example**:

```typescript
import { isValidK8sName } from '@odh-dashboard/autox/utils';

function NameInput() {
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

  return <TextInput value={name} onChange={handleChange} validated={error ? 'error' : 'default'} />;
}
```

---

### 4.3 Topology Utilities

```typescript
import type { TopologyNode, TopologyEdge, RuntimeStateKF } from '../types';

export function createTopologyNode(
  id: string,
  label: string,
  type: 'task' | 'artifact',
  status?: RuntimeStateKF,
  data?: Record<string, unknown>,
): TopologyNode;

export function createTopologyEdge(
  source: string,
  target: string,
): TopologyEdge;
```

**Stability**: ✅ Stable  
**Extractable**: 100% identical (39 LOC)

**Usage Example**:

```typescript
import { createTopologyNode, createTopologyEdge } from '@odh-dashboard/autox/utils';

function buildPipelineGraph(pipelineRun: PipelineRun): TopologyGraph {
  const nodes = [
    createTopologyNode('task1', 'Data Ingestion', 'task', pipelineRun.state),
    createTopologyNode('task2', 'Training', 'task', RuntimeStateKF.PENDING),
  ];
  
  const edges = [
    createTopologyEdge('task1', 'task2'),
  ];
  
  return { nodes, edges };
}
```

---

### 4.4 URL Utilities

```typescript
export function buildApiUrl(path: string, params?: Record<string, string>): string;
export function parseQueryParams(search: string): Record<string, string>;
```

**Stability**: ✅ Stable  
**Extractable**: 100% identical (11 LOC)

**Usage Example**:

```typescript
import { buildApiUrl } from '@odh-dashboard/autox/utils';

async function fetchRuns(namespace: string, filter?: string) {
  const url = buildApiUrl('/api/pipeline-runs', { namespace, filter });
  const response = await fetch(url);
  return response.json();
}
```

---

### 4.5 Formatting Utilities

```typescript
export function formatBytes(bytes: number): string;
export function formatTimestamp(timestamp: string): string;
export function formatDuration(startTime: string, endTime?: string): string;
```

**Stability**: ✅ Stable

**Usage Example**:

```typescript
import { formatBytes, formatDuration } from '@odh-dashboard/autox/utils';

function FileDetails({ file }: { file: FileItem }) {
  return (
    <div>
      <div>Size: {formatBytes(file.size)}</div>
    </div>
  );
}

function RunDuration({ run }: { run: PipelineRun }) {
  return <div>Duration: {formatDuration(run.createdAt, run.finishedAt)}</div>;
}
```

---

## 5. Module Federation Configuration

### Webpack Config (Consumer)

**AutoML and AutoRAG must configure AutoX as a shared dependency**:

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
    '@odh-dashboard/autox': {
      singleton: true,
      requiredVersion: dependencies['@odh-dashboard/autox'],
    },
  },
};
```

**Why `singleton: true`**:
- Ensures only one instance of AutoX in browser
- Prevents duplicate React Query state
- Reduces bundle size
- Guarantees hook stability

---

## 6. Version Compatibility

### Current Version: 1.0.0-alpha

**Supported Versions**:
- React: 18.x
- TypeScript: 5.x
- React Query: 5.x
- PatternFly: 6.x

### Deprecation Policy

- **Deprecation notices**: Added as JSDoc comments and in release notes
- **Deprecation period**: Minimum 2 releases (or 1 month)
- **Removal**: Only in major version bumps

**Example**:

```typescript
/**
 * @deprecated Use usePipelineRuns instead.
 * Will be removed in v2.0.0.
 */
export function usePipelines() {
  // ...
}
```

---

## 7. Testing Strategy

### Consumer Testing

**AutoML and AutoRAG must test**:
1. **Import validation**: All AutoX imports work without errors
2. **Hook behavior**: Custom hook usage integrates correctly with AutoX hooks
3. **Component rendering**: Composed components render correctly
4. **Type safety**: TypeScript compilation succeeds

**AutoX provides**:
- Unit tests for all hooks (React Testing Library)
- Component tests (React Testing Library)
- Type definitions with full JSDoc comments
- Mock data factories for testing (`@odh-dashboard/autox/mocks`)

---

## 8. Migration Guide

### From Duplicate Code to AutoX

**Step 1**: Identify duplicate code in AutoML/AutoRAG frontend  
**Step 2**: Find corresponding AutoX export in this contract  
**Step 3**: Update imports:

```typescript
// Before (AutoML internal code)
import { usePipelineRuns } from '../hooks/usePipelineRuns';
import type { PipelineRun } from '../types/pipeline';

// After (AutoX shared code)
import { usePipelineRuns } from '@odh-dashboard/autox/hooks';
import type { PipelineRun } from '@odh-dashboard/autox/types';
```

**Step 4**: Remove duplicate code from consumer package  
**Step 5**: Update Module Federation config to include AutoX as shared singleton  
**Step 6**: Run tests to verify behavior preserved

---

## 9. Composition Patterns

### Hook Composition

```typescript
// AutoRAG-specific hook composing AutoX primitives
import { usePipelineRuns } from '@odh-dashboard/autox/hooks';

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

### Component Composition

```typescript
// AutoML-specific table extending AutoX primitive
import { PipelineRunsTable } from '@odh-dashboard/autox/components';

export function AutoMLRunsTable({ runs }: { runs: PipelineRun[] }) {
  return (
    <div>
      <PipelineRunsTable
        runs={runs}
        showActions
        onRunSelect={(run) => navigate(`/automl/runs/${run.id}`)}
      />
      {/* Add AutoML-specific metrics summary below table */}
      <MetricsSummary runs={runs} />
    </div>
  );
}
```

---

## 10. Contact & Support

**Maintainers**: ODH Dashboard Team  
**Issues**: Report via GitHub Issues with `autox` label  
**Breaking Changes**: Announced in PR description and release notes  
**CI/CD**: All packages type-checked together (`tsc --noEmit`) from repo root

---

**Next**: See [bff-api.md](./bff-api.md) for BFF contract, [../quickstart.md](../quickstart.md) for setup guide.
