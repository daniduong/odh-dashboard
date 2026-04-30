# Frontend Utils and Types Duplication Analysis

**Analysis Date**: 2026-04-23
**Scope**: AutoML and AutoRAG frontend utility functions and type definitions

## Executive Summary

This analysis compares all utility functions, type definitions, and constants between the AutoML and AutoRAG frontend packages to identify code duplication and opportunities for a shared library.

**Key Findings**:
- **100% identical files**: 6 files (pipeline.ts, topology.ts, topology utils, schema.ts, secretValidation.ts, pipelineServerEmptyState.ts)
- **Highly similar files**: 2 files (utils.ts ~40% duplicate, const.ts ~70% duplicate)
- **Domain-specific files**: 2 files (columnUtils.ts, autoragPattern.ts)
- **Total duplication**: ~800 lines of pure duplicate code

## Detailed File-by-File Analysis

### 1. Type Definitions

#### pipeline.ts (100% Duplicate)

| Metric | Value |
|--------|-------|
| **Location (AutoML)** | `packages/automl/frontend/src/app/types/pipeline.ts` |
| **Location (AutoRAG)** | `packages/autorag/frontend/src/app/types/pipeline.ts` |
| **Total Lines** | 186 |
| **Duplicate Lines** | 186 (100%) |
| **Duplication %** | 100% |

**Content**:
- KFP v2beta1 types for pipeline visualization
- Enums: `RuntimeStateKF`, `ExecutionStateKF`
- Constants: `runtimeStateLabels`
- Types: Pipeline spec, components, executors, tasks, runs
- **Classification**: Pure utility types (should be shared)

**Identical Types**:
```typescript
- RuntimeStateKF (10 states)
- ExecutionStateKF (8 states)
- runtimeStateLabels mapping
- InputOutputArtifactType
- InputOutputDefinitionArtifacts
- InputOutputDefinition
- TaskKF
- DAG
- PipelineComponentKF
- PipelineComponentsKF
- PipelineExecutorKF
- PipelineExecutorsKF
- PipelineSpec
- PlatformSpec
- PipelineSpecVariable
- TaskDetailKF
- RunDetailsKF
- PipelineVersionReferenceKF
- PipelineRunKF
- PipelineVersionKF
```

**Recommendation**: Move to shared library as-is.

---

#### topology.ts (100% Duplicate)

| Metric | Value |
|--------|-------|
| **Location (AutoML)** | `packages/automl/frontend/src/app/types/topology.ts` |
| **Location (AutoRAG)** | `packages/autorag/frontend/src/app/types/topology.ts` |
| **Total Lines** | 46 |
| **Duplicate Lines** | 46 (100%) |
| **Duplication %** | 100% |

**Content**:
- Topology types for pipeline graph visualization
- Integration with PatternFly topology components
- **Classification**: Pure utility types (should be shared)

**Identical Types**:
```typescript
- PipelineTaskRunStatus
- PipelineTaskStep
- PipelineTaskInputOutput
- PipelineTask
- StandardTaskNodeData
- PipelineNodeModelExpanded
```

**Recommendation**: Move to shared library as-is.

---

#### autoragPattern.ts (AutoRAG-Specific)

| Metric | Value |
|--------|-------|
| **Location** | `packages/autorag/frontend/src/app/types/autoragPattern.ts` |
| **Total Lines** | 74 |
| **Duplicate Lines** | 0 |
| **Duplication %** | 0% |

**Content**:
- AutoRAG-specific pattern types
- Evaluation result types
- Score metric structures
- **Classification**: Domain-specific (keep in AutoRAG)

**Recommendation**: Keep in AutoRAG package.

---

### 2. Utility Functions

#### utilities/utils.ts (40% Duplicate)

| Metric | AutoML | AutoRAG | Duplicate |
|--------|--------|---------|-----------|
| **Total Lines** | 204 | 137 | ~80 |
| **Duplication %** | 40% | 58% | - |

**Shared Functions** (100% identical):

1. **`isRunTerminatable(state)`** - Lines 7-18 (both)
   ```typescript
   export const isRunTerminatable = (state: string | undefined): boolean => {
     const s = state?.toUpperCase();
     return (
       s === RuntimeStateKF.RUNNING || s === RuntimeStateKF.PENDING || s === RuntimeStateKF.PAUSED
     );
   };
   ```

2. **`isRunInProgress(state)`** - Lines 18-29 (both)
   ```typescript
   export const isRunInProgress = (state: string | undefined): boolean => {
     const s = state?.toUpperCase();
     return (
       s === RuntimeStateKF.RUNNING || s === RuntimeStateKF.PENDING || s === RuntimeStateKF.CANCELING
     );
   };
   ```

3. **`isRunRetryable(state)`** - Lines 28-37 (both)
   ```typescript
   export const isRunRetryable = (state: string | undefined): boolean => {
     const s = state?.toUpperCase();
     return s === RuntimeStateKF.FAILED || s === RuntimeStateKF.CANCELED;
   };
   ```

4. **`parseErrorStatus(error)`** - Lines 39-55 (both)
   ```typescript
   export function parseErrorStatus(error: Error): number | undefined {
     const match =
       error.message.match(/\bstatus\s+code\s+(\d{3})\b/i) ??
       error.message.match(/\bstatus[:\s]+(\d{3})\b/i) ??
       error.message.match(/\b(403|404|503)\b/);
     if (match) {
       const code = parseInt(match[1], 10);
       return code >= 100 && code < 600 ? code : undefined;
     }
     return undefined;
   }
   ```

5. **`formatMetricValue(value)`** - Lines 123-135 (both)
   ```typescript
   export function formatMetricValue(value: number | string): string {
     if (typeof value === 'string') {
       return value;
     }
     const fixed = value.toFixed(3);
     if ((fixed === '0.000' || fixed === '-0.000') && value !== 0) {
       return value.toExponential(3);
     }
     return fixed;
   }
   ```

6. **`downloadBlob(blob, filename)`** - Lines 193-203 (AutoML), 78-87 (AutoRAG)
   ```typescript
   export function downloadBlob(blob: Blob, filename: string): void {
     const url = URL.createObjectURL(blob);
     const link = document.createElement('a');
     link.href = url;
     link.download = filename;
     document.body.appendChild(link);
     link.click();
     document.body.removeChild(link);
     URL.revokeObjectURL(url);
   }
   ```

**AutoML-Specific Functions**:
- `getTaskType(pipelineRun)` - Extracts task type from pipeline parameters
- `isTabularRun(pipelineRun)` - Checks if task is tabular
- `formatMetricName(key)` - Formats ML metric keys (F₁, MAE, MAPE, etc.)
- `toNumericMetric(value)` - Coerces metric values to numbers
- `getOptimizedMetricForTask(taskType)` - Returns metric for task type
- `computeRankMap(models, taskType)` - Builds leaderboard rankings

**AutoRAG-Specific Functions**:
- `getOptimizedMetricForRAG(pipelineRun)` - Gets RAG optimization metric
- `sanitizeFilename(str)` - Sanitizes filenames for downloads
- `formatPatternName(name)` - Formats pattern names with non-breaking spaces
- `formatMetricName(metricKey)` - Formats RAG metric keys (different mappings)

**Recommendation**:
- Move 6 shared functions to shared library
- Keep domain-specific functions in respective packages
- Consider generic metric formatting utility with plugin system

---

#### topology/utils.ts (100% Duplicate)

| Metric | Value |
|--------|-------|
| **Location (AutoML)** | `packages/automl/frontend/src/app/topology/utils.ts` |
| **Location (AutoRAG)** | `packages/autorag/frontend/src/app/topology/utils.ts` |
| **Total Lines** | 40 |
| **Duplicate Lines** | 40 (100%) |
| **Duplication %** | 100% |

**Content**:
- Canvas text width measurement
- Pipeline node creation utility
- **Classification**: Pure utility (should be shared)

**Functions**:
```typescript
- getCanvasContext(): CanvasRenderingContext2D | null
- measureTextWidth(text: string): number
- createNode(id, label, pipelineTask, runAfterTasks?, runStatus?): PipelineNodeModelExpanded
```

**Recommendation**: Move to shared library as-is.

---

#### utilities/schema.ts (100% Duplicate)

| Metric | Value |
|--------|-------|
| **Location (AutoML)** | `packages/automl/frontend/src/app/utilities/schema.ts` |
| **Location (AutoRAG)** | `packages/autorag/frontend/src/app/utilities/schema.ts` |
| **Total Lines** | 26 |
| **Duplicate Lines** | 26 (100%) |
| **Duplication %** | 100% |

**Content**:
- Zod schema builder with validators and transformers
- **Classification**: Pure utility (should be shared)

**Functions**:
```typescript
- createSchema<Schema>({ schema, validators?, transformers? })
```

**Recommendation**: Move to shared library as-is.

---

#### utilities/secretValidation.ts (100% Duplicate)

| Metric | Value |
|--------|-------|
| **Location (AutoML)** | `packages/automl/frontend/src/app/utilities/secretValidation.ts` |
| **Location (AutoRAG)** | `packages/autorag/frontend/src/app/utilities/secretValidation.ts` |
| **Total Lines** | 25 |
| **Duplicate Lines** | 25 (100%) |
| **Duplication %** | 100% |

**Content**:
- Secret key validation utilities
- **Classification**: Pure utility (should be shared)

**Functions**:
```typescript
- getMissingRequiredKeys(requiredKeys, availableKeys): string[]
- formatMissingKeysMessage(missingKeys): string
```

**Recommendation**: Move to shared library as-is.

---

#### utilities/pipelineServerEmptyState.ts (100% Duplicate)

| Metric | Value |
|--------|-------|
| **Location (AutoML)** | `packages/automl/frontend/src/app/utilities/pipelineServerEmptyState.ts` |
| **Location (AutoRAG)** | `packages/autorag/frontend/src/app/utilities/pipelineServerEmptyState.ts` |
| **Total Lines** | 50 |
| **Duplicate Lines** | 50 (100%) |
| **Duplication %** | 100% |

**Content**:
- Pipeline server error detection utilities
- **Classification**: Pure utility (should be shared)

**Functions**:
```typescript
- shouldShowConfigurePipelineServerEmptyState(error): boolean
- shouldShowPipelineServerNotReady(error): boolean
```

**Recommendation**: Move to shared library as-is.

---

#### utilities/columnUtils.ts (AutoML-Specific)

| Metric | Value |
|--------|-------|
| **Location** | `packages/automl/frontend/src/app/utilities/columnUtils.ts` |
| **Total Lines** | 17 |
| **Duplicate Lines** | 0 |
| **Duplication %** | 0% |

**Content**:
- Data type to acronym conversion (BOOL, INT, DBL, TMSTP, STR)
- **Classification**: Domain-specific (keep in AutoML)

**Recommendation**: Keep in AutoML package.

---

### 3. Constants

#### utilities/const.ts (70% Duplicate)

| Metric | AutoML | AutoRAG | Duplicate |
|--------|--------|---------|-----------|
| **Total Lines** | 59 | 44 | ~30 |
| **Duplication %** | 51% | 68% | - |

**Shared Constants** (100% identical):

```typescript
// Environment config (lines 3-29)
- STYLE_THEME
- DEPLOYMENT_MODE
- DEV_MODE
- POLL_INTERVAL
- KUBEFLOW_USERNAME
- IMAGE_DIR
- LOGO_LIGHT
- MANDATORY_NAMESPACE
- BFF_API_VERSION
- COMPANY_URI

// User-facing strings (lines 32-35)
- FindAdministratorOptions (array of 3 strings)
```

**AutoML-Specific Constants**:
```typescript
- URL_PREFIX = '/automl'
- DEFAULT_PAGE_SIZE = 20
- TASK_TYPE_BINARY = 'binary'
- TASK_TYPE_MULTICLASS = 'multiclass'
- TASK_TYPE_REGRESSION = 'regression'
- TASK_TYPE_TIMESERIES = 'timeseries'
- AUTOML_OPTIMIZED_METRIC_BY_TASK (mapping)
- TASK_TYPE_LABELS (mapping)
```

**AutoRAG-Specific Constants**:
```typescript
- URL_PREFIX = '/autorag'
- OPTIMIZATION_METRIC_LABELS (mapping)
```

**Recommendation**:
- Move 11 shared constants to shared library
- Parameterize URL_PREFIX by package name
- Keep domain-specific constants in respective packages

---

#### topology/const.ts (100% Duplicate)

| Metric | Value |
|--------|-------|
| **Location (AutoML)** | `packages/automl/frontend/src/app/topology/const.ts` |
| **Location (AutoRAG)** | `packages/autorag/frontend/src/app/topology/const.ts` |
| **Total Lines** | 12 |
| **Duplicate Lines** | 12 (100%) |
| **Duplication %** | 100% |

**Content**:
- Pipeline topology layout constants
- Node dimensions and spacing
- **Classification**: Pure constants (should be shared)

**Constants**:
```typescript
- PIPELINE_LAYOUT = 'PipelineLayout'
- PIPELINE_NODE_SEPARATION_VERTICAL = 20
- PIPELINE_NODE_SEPARATION_HORIZONTAL = 70
- NODE_WIDTH = 200
- NODE_PADDING = 40
- NODE_HEIGHT = 35
- NODE_FONT = '0.875rem RedHatText'
```

**Recommendation**: Move to shared library as-is.

---

#### utilities/routes.ts (Pattern-Based)

| Metric | AutoML | AutoRAG |
|--------|--------|---------|
| **Total Lines** | 5 | 5 |
| **Duplicate Pattern** | 100% | 100% |

**AutoML Routes**:
```typescript
export const rootPathname = '/develop-train/automl';
export const automlExperimentsPathname = `${rootPathname}/experiments`;
export const automlConfigurePathname = `${rootPathname}/configure`;
export const automlResultsPathname = `${rootPathname}/results`;
```

**AutoRAG Routes**:
```typescript
export const rootPathname = '/gen-ai-studio/autorag';
export const autoragExperimentsPathname = `${rootPathname}/experiments`;
export const autoragConfigurePathname = `${rootPathname}/configure`;
export const autoragResultsPathname = `${rootPathname}/results`;
```

**Recommendation**: Create route builder utility in shared library:
```typescript
export function createAutoXRoutes(baseRoot: string, productName: string) {
  const rootPathname = `${baseRoot}/${productName}`;
  return {
    root: rootPathname,
    experiments: `${rootPathname}/experiments`,
    configure: `${rootPathname}/configure`,
    results: `${rootPathname}/results`,
  };
}
```

---

## Summary Tables

### Pure Duplicates (100%)

| File | Lines | Category | Priority |
|------|-------|----------|----------|
| `types/pipeline.ts` | 186 | Types | High |
| `types/topology.ts` | 46 | Types | High |
| `topology/utils.ts` | 40 | Utils | High |
| `topology/const.ts` | 12 | Constants | High |
| `utilities/schema.ts` | 26 | Utils | High |
| `utilities/secretValidation.ts` | 25 | Utils | High |
| `utilities/pipelineServerEmptyState.ts` | 50 | Utils | High |
| **Total** | **385** | - | - |

### Partial Duplicates

| File | Total Lines | Duplicate Lines | % | Priority |
|------|-------------|-----------------|---|----------|
| `utilities/utils.ts` | 204 (AutoML)<br>137 (AutoRAG) | ~80 | 40%/58% | High |
| `utilities/const.ts` | 59 (AutoML)<br>44 (AutoRAG) | ~30 | 51%/68% | Medium |
| `utilities/routes.ts` | 5 (both) | Pattern-based | 100% | Low |

### Domain-Specific (No Duplication)

| File | Lines | Package | Purpose |
|------|-------|---------|---------|
| `types/autoragPattern.ts` | 74 | AutoRAG | RAG pattern types |
| `utilities/columnUtils.ts` | 17 | AutoML | Column type acronyms |

---

## Shared Library Organization

### Recommended Structure

```
packages/autox/src/
├── types/
│   ├── pipeline.ts         # 186 lines (pure duplicate)
│   └── topology.ts         # 46 lines (pure duplicate)
├── utilities/
│   ├── runState.ts         # Run state functions (isRunTerminatable, etc.)
│   ├── errorParsing.ts     # parseErrorStatus
│   ├── metrics.ts          # formatMetricValue
│   ├── download.ts         # downloadBlob
│   ├── schema.ts           # createSchema (Zod utility)
│   ├── secretValidation.ts # getMissingRequiredKeys, formatMissingKeysMessage
│   └── pipelineServer.ts   # shouldShowConfigurePipelineServerEmptyState, etc.
├── topology/
│   ├── const.ts            # Topology constants
│   └── utils.ts            # createNode, measureTextWidth
├── constants/
│   ├── env.ts              # Environment constants
│   └── ui.ts               # FindAdministratorOptions
└── builders/
    └── routes.ts           # createAutoXRoutes
```

### Migration Impact

**Files to Create**: 14 files in shared library

**Files to Update**:
- AutoML: 10 files (remove duplicate code, import from shared)
- AutoRAG: 9 files (remove duplicate code, import from shared)

**Lines of Code Reduction**: ~800 lines

---

## Classification Summary

### Pure Utilities (Move to Shared Library)

**Functions**:
- `isRunTerminatable`, `isRunInProgress`, `isRunRetryable` - Run state checks
- `parseErrorStatus` - HTTP status extraction from error messages
- `formatMetricValue` - Metric value formatting with scientific notation
- `downloadBlob` - Browser download trigger
- `createSchema` - Zod schema builder
- `getMissingRequiredKeys`, `formatMissingKeysMessage` - Secret validation
- `shouldShowConfigurePipelineServerEmptyState`, `shouldShowPipelineServerNotReady` - Error detection
- `createNode`, `measureTextWidth`, `getCanvasContext` - Topology utilities

**Types**:
- All KFP pipeline types (RuntimeStateKF, PipelineSpec, etc.)
- All topology types (PipelineTask, StandardTaskNodeData, etc.)

**Constants**:
- Environment config (STYLE_THEME, POLL_INTERVAL, etc.)
- Topology layout (NODE_WIDTH, NODE_HEIGHT, etc.)
- UI strings (FindAdministratorOptions)

### Domain-Specific (Keep in Packages)

**AutoML**:
- `getTaskType`, `isTabularRun` - Task type detection
- `formatMetricName` - ML metric formatting (F₁, MAE, R², etc.)
- `toNumericMetric`, `getOptimizedMetricForTask`, `computeRankMap` - ML-specific
- `getTypeAcronym` - Column type acronyms
- Task type constants (TASK_TYPE_BINARY, etc.)

**AutoRAG**:
- `getOptimizedMetricForRAG` - RAG metric extraction
- `sanitizeFilename` - Filename sanitization
- `formatPatternName` - Pattern name formatting
- `formatMetricName` - RAG metric formatting (different from AutoML)
- AutoRAG pattern types
- Optimization metric labels

---

## Recommendations

### High Priority (Immediate)

1. **Create shared types package** - 232 lines of pure duplicate type definitions
2. **Extract topology utilities** - 52 lines of pure duplicate topology code
3. **Move validation utilities** - 75 lines of pure duplicate validation code

### Medium Priority (Sprint 2)

1. **Extract run state utilities** - 80 lines of partial duplicates in utils.ts
2. **Centralize environment constants** - 30 lines of duplicate constants
3. **Create route builder pattern** - Remove route duplication

### Low Priority (Future)

1. **Generic metric formatting** - Consider plugin system for different metric types
2. **Shared UI constants** - Extract FindAdministratorOptions and similar strings

---

## Test Coverage Impact

**Files with Tests** (need updates):
- AutoML: `utilities/utils.ts`, `topology/utils.ts`, `utilities/schema.ts`
- AutoRAG: `utilities/utils.ts`, `topology/utils.ts`, `utilities/schema.ts`

**Test Migration**: Move shared function tests to shared library package.

---

## Next Steps

1. Review this analysis with team
2. Create shared library package structure
3. Move pure duplicates first (types, topology, validation)
4. Update imports in AutoML and AutoRAG
5. Move partial duplicates with domain-specific functions remaining
6. Update tests
7. Document shared library API

---

## Appendix: Line Count by Category

| Category | Pure Duplicate | Partial Duplicate | Domain-Specific | Total |
|----------|----------------|-------------------|-----------------|-------|
| **Types** | 232 | 0 | 74 | 306 |
| **Utils** | 141 | 80 | ~150 | ~371 |
| **Constants** | 12 | 30 | ~30 | ~72 |
| **Topology** | 52 | 0 | 0 | 52 |
| **Routes** | 0 | 10 | 0 | 10 |
| **Total** | **437** | **120** | **~254** | **~811** |

**Duplication Rate**: 68.7% (557/811 lines are duplicate or partially duplicate)
