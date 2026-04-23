# Tasks: AutoX Shared Library Package

**Feature Branch**: `002-autox-shared-library`  
**Generated**: 2026-04-23  
**Total Tasks**: 202 (updated after /speckit.analyze remediation)
**Estimated Effort**: 360 hours (8-12 weeks with 2-3 developers)

---

## Overview

This document provides a complete task breakdown for implementing the AutoX shared library package. Tasks are organized by user story to enable independent implementation and testing. Each user story represents a complete, deliverable increment.

**Key Principles**:
- **Independent Stories**: Each user story can be implemented and tested independently
- **Parallel Execution**: Tasks marked with `[P]` can be executed in parallel
- **Synchronized Versioning**: AutoX, AutoML, and AutoRAG are always updated together in the same PR
- **Type-Safe Refactoring**: CI type-checking (`tsc --noEmit`, `go build`) catches breaking changes before merge

---

## Implementation Strategy

### MVP Scope (Recommended First Iteration)

**Start with User Story 1 + User Story 5 + Minimal US7**:
1. **US5**: Development Environment Setup (foundational)
2. **US1**: Import Shared BFF Utilities (first value delivery)
3. **US7** (Subset): Extract perfect duplicates only (FileExplorer, useNotification, S3 API)

**Why this MVP**:
- Establishes infrastructure (workspaces, package structure)
- Delivers immediate value (BFF utilities reduce duplication)
- Proves the pattern (frontend primitives show composition works)
- Low risk (perfect duplicates = zero refactoring needed)
- Estimated effort: ~60 hours (1.5-2 weeks)

### Full Implementation Order

**Recommended sequence after MVP**:
1. ✅ MVP: US5 + US1 + US7 (Perfect Duplicates)
2. US2: Compose UI Primitives (high value, builds on US1)
3. US7 (Phase 2): High-value components extraction
4. US6: Module Federation Runtime Singleton (performance optimization)
5. US3: Handle Domain-Specific Logic (clean separation)
6. US7 (Phase 3-4): Hooks, utilities, types extraction
7. US4: Strategy/DI Patterns (advanced customization)
8. US7 (Phase 5): Advanced patterns extraction

---

## Dependencies & Execution Flow

### Story Dependencies

```text
Setup Phase (T001-T011)
  ↓
Foundational Phase (T012-T024)
  ↓
┌─────────────────────────────────────────────────────┐
│ P1 Stories (Can Execute in Parallel)               │
│                                                     │
│  [US1: BFF Utilities] ←──────┬──→ [US2: UI Primitives]
│         ↓                     │           ↓
│         ↓                     │           ↓
│  [US7: Refactor Code] ←──────┘           ↓
│         ↓                                 ↓
│         └─────────────────────────────────┘
│                     ↓
├─────────────────────────────────────────────────────┤
│ P2 Stories (After P1 Complete)                      │
│                                                     │
│  [US5: Dev Environment] (Foundational - moved up)  │
│         ↓                                           │
│  [US6: Module Federation]                          │
│         ↓                                           │
│  [US3: Domain Logic Handlers]                      │
│                                                     │
├─────────────────────────────────────────────────────┤
│ P3 Stories (Advanced)                               │
│                                                     │
│  [US4: Strategy/DI Patterns]                       │
│                                                     │
└─────────────────────────────────────────────────────┘
         ↓
Polish Phase (T083-T089)
```

**Key Insight**: US5 (Dev Environment Setup) is actually foundational and should be completed before US1/US2. It's moved to Foundational Phase in task list.

### Parallel Execution Opportunities

**Within Each Story Phase**:
- BFF models and Frontend types can be created in parallel
- BFF integrations can be implemented in parallel (S3, K8s, Pipeline Server)
- Frontend components can be implemented in parallel (FileExplorer, PipelineRunsTable, etc.)
- Unit tests can be written in parallel with implementation (TDD approach if desired)

**Cross-Story Parallelization**:
- US1 (BFF) and US2 (UI) are completely independent - can execute in parallel
- US7 extraction phases can overlap with US3 handler file implementation

---

## Phase 1: Setup

**Goal**: Create AutoX package structure, establish build/test infrastructure, and produce Phase 1 documentation deliverables

**Duration**: 8-12 hours  
**Dependencies**: None  
**Deliverable**: Compilable AutoX package with passing empty tests, complete documentation (data-model.md, contracts, quickstart.md)

### Tasks

- [ ] T001 Create packages/autox directory structure
- [ ] T002 [P] Create packages/autox/package.json with workspace configuration
- [ ] T003 [P] Create packages/autox/frontend/package.json with dependencies (React 18, PatternFly v6, React Query v5, Zod)
- [ ] T004 [P] Create packages/autox/frontend/tsconfig.json extending @odh-dashboard/tsconfig
- [ ] T005 [P] Create packages/autox/frontend/.eslintrc.js extending @odh-dashboard/eslint-config
- [ ] T006 [P] Create packages/autox/frontend/jest.config.ts extending @odh-dashboard/jest-config
- [ ] T007 [P] Create packages/autox/bff/go.mod with module path github.com/opendatahub-io/odh-dashboard/packages/autox/bff
- [ ] T008 [P] Create packages/autox/bff/internal directory structure (api/, models/, integrations/, utils/)
- [ ] T009 [P] Create packages/autox/README.md with package overview and usage guide
- [ ] T010 [P] Create packages/autox/frontend/src directory structure (hooks/, components/, utils/, types/)
- [ ] T011 Run npm install from repo root to link AutoX via npm workspaces
- [ ] T011a Create specs/002-autox-shared-library/data-model.md documenting key entities, interfaces, and relationships from spec
- [ ] T011b Create specs/002-autox-shared-library/contracts/bff-api.md documenting all BFF exported interfaces, models, and utilities
- [ ] T011c Create specs/002-autox-shared-library/contracts/ui-api.md documenting all UI exported hooks, components, and types
- [ ] T011d Create specs/002-autox-shared-library/quickstart.md with developer setup guide, import patterns, and testing strategy
- [ ] T011e Run .specify/scripts/bash/update-agent-context.sh claude to update agent knowledge with AutoX package

**Independent Test**: Package structure exists, TypeScript and Go compilation succeed with no errors, all documentation deliverables created

---

## Phase 2: Foundational (Blocking Prerequisites)

**Goal**: Establish workspace configuration, Module Federation setup, and build infrastructure that all user stories depend on

**Duration**: 12-16 hours  
**Dependencies**: Phase 1 complete  
**Deliverable**: Working workspace setup, Module Federation configured, AutoX importable by AutoML/AutoRAG

### Tasks: US5 - Development Environment Setup

- [ ] T012 [US5] Initialize Go workspace with go work init from repo root
- [ ] T013 [US5] Add AutoX BFF to Go workspace with go work use packages/autox/bff
- [ ] T014 [US5] Add AutoML BFF to Go workspace with go work use packages/automl/bff
- [ ] T015 [US5] Add AutoRAG BFF to Go workspace with go work use packages/autorag/bff
- [ ] T016 [US5] Run go work sync to synchronize Go workspace
- [ ] T017 [US5] Verify Go workspace by building all BFF packages with go build ./... from repo root
- [ ] T018 [US5] Update root package.json workspaces array to include packages/autox
- [ ] T019 [US5] Run npm install from repo root to establish npm workspace links
- [ ] T020 [US5] Verify AutoX is linked by running npm run type-check from packages/automl/frontend
- [ ] T021 [US5] Verify AutoX is linked by running npm run type-check from packages/autorag/frontend

**US5 Independent Test**: 
- ✅ npm install succeeds and AutoX is linked to AutoML/AutoRAG
- ✅ go work sync succeeds and all BFF packages compile
- ✅ Changes to AutoX are immediately available in AutoML/AutoRAG without manual linking

### Tasks: US6 - Module Federation Runtime Singleton

- [ ] T022 [US6] Add AutoX as shared singleton in packages/automl/frontend/config/moduleFederation.js
- [ ] T023 [US6] Add AutoX as shared singleton in packages/autorag/frontend/config/moduleFederation.js
- [ ] T024 [US6] Verify Module Federation singleton by running webpack build and checking only one AutoX instance in bundle

**US6 Independent Test**:
- ✅ AutoML and AutoRAG webpack builds succeed with AutoX as shared dependency
- ✅ DevTools shows only one instance of AutoX in runtime
- ✅ Bundle size is optimized (no duplicate AutoX code)

---

## Phase 3: User Story 1 - Import Shared BFF Utilities (P1)

**Goal**: Extract and share BFF interfaces, utilities, and clients between AutoML and AutoRAG

**Duration**: 20-24 hours  
**Dependencies**: Foundational phase complete  
**Deliverable**: AutoML and AutoRAG BFF handlers successfully import and use AutoX utilities

### Implementation Tasks

- [ ] T025 [P] [US1] Create packages/autox/bff/internal/models/health_check.go with HealthCheck struct
- [ ] T026 [P] [US1] Create packages/autox/bff/internal/models/namespace.go with Namespace struct
- [ ] T027 [P] [US1] Create packages/autox/bff/internal/models/user.go with User struct
- [ ] T028 [P] [US1] Create packages/autox/bff/internal/models/secret.go with Secret struct
- [ ] T029 [P] [US1] Create packages/autox/bff/internal/models/s3.go with S3Credentials struct
- [ ] T030 [P] [US1] Create packages/autox/bff/internal/models/pipeline_runs.go with PipelineRun struct and RuntimeStateKF enum
- [ ] T031 [P] [US1] Create packages/autox/bff/internal/models/pipeline_server.go with PipelineServer struct
- [ ] T032 [P] [US1] Create packages/autox/bff/internal/api/errors.go with ErrorResponse struct and error helper functions
- [ ] T033 [P] [US1] Create packages/autox/bff/internal/api/helpers.go with WriteJSON and ReadJSON functions
- [ ] T034 [P] [US1] Create packages/autox/bff/internal/integrations/kubernetes/interface.go with KubernetesClientInterface
- [ ] T035 [US1] Create packages/autox/bff/internal/integrations/kubernetes/client.go with default K8s client implementation
- [ ] T036 [US1] Create packages/autox/bff/internal/integrations/kubernetes/factory.go with NewInternalK8sClient factory
- [ ] T037 [US1] Create packages/autox/bff/internal/integrations/kubernetes/portforward.go with port forwarding utilities
- [ ] T038 [P] [US1] Create packages/autox/bff/internal/integrations/s3/interface.go with S3ClientInterface
- [ ] T039 [US1] Create packages/autox/bff/internal/integrations/s3/client.go with AWS SDK wrapper implementation
- [ ] T040 [P] [US1] Create packages/autox/bff/internal/integrations/pipelineserver/interface.go with PipelineServerClientInterface
- [ ] T041 [US1] Create packages/autox/bff/internal/integrations/pipelineserver/client.go with KFP client implementation
- [ ] T042 [P] [US1] Create packages/autox/bff/internal/utils/helpers.go with common helper functions
- [ ] T043 [P] [US1] Create packages/autox/bff/internal/utils/pagination.go with query parameter filtering utilities

### Unit Test Tasks (AutoX Internal)

- [ ] T044 [P] [US1] Create packages/autox/bff/internal/api/errors_test.go with error helper tests
- [ ] T045 [P] [US1] Create packages/autox/bff/internal/api/helpers_test.go with JSON helper tests
- [ ] T046 [P] [US1] Create packages/autox/bff/internal/integrations/kubernetes/client_test.go with K8s client tests
- [ ] T047 [P] [US1] Create packages/autox/bff/internal/integrations/s3/client_test.go with S3 client tests
- [ ] T048 [P] [US1] Create packages/autox/bff/internal/integrations/pipelineserver/client_test.go with pipeline client tests
- [ ] T049 [P] [US1] Create packages/autox/bff/internal/utils/helpers_test.go with helper function tests

### Integration Tasks (Verify in AutoML/AutoRAG)

- [ ] T050 [US1] Update packages/automl/bff imports to use autox/bff/internal/models for common models
- [ ] T051 [US1] Update packages/automl/bff imports to use autox/bff/internal/api for error handlers and helpers
- [ ] T052 [US1] Update packages/automl/bff to inject AutoX K8s client in handlers
- [ ] T053 [US1] Update packages/autorag/bff imports to use autox/bff/internal/models for common models
- [ ] T054 [US1] Update packages/autorag/bff imports to use autox/bff/internal/api for error handlers and helpers
- [ ] T055 [US1] Update packages/autorag/bff to inject AutoX S3 client in handlers
- [ ] T056 [US1] Run go test ./... from packages/automl/bff and packages/autorag/bff to verify all tests pass
- [ ] T057 [US1] Run go build ./... from repo root to verify all BFF packages compile successfully

**US1 Independent Test**:
- ✅ AutoML BFF handler imports AutoX models, compiles, and tests pass
- ✅ AutoRAG BFF handler imports AutoX API utilities, compiles, and tests pass
- ✅ Both packages use AutoX clients and behavior is consistent across packages
- ✅ No compilation errors in TypeScript or Go
- ✅ CI type-checking passes

---

## Phase 4: User Story 2 - Compose UI Primitives (P1)

**Goal**: Extract and share frontend hooks and components as composable primitives

**Duration**: 24-28 hours  
**Dependencies**: Foundational phase complete  
**Deliverable**: AutoML and AutoRAG frontends successfully compose AutoX primitives into domain-specific features

### Implementation Tasks: Types

- [ ] T058 [P] [US2] Create packages/autox/frontend/src/types/pipeline.ts with PipelineRun, RuntimeStateKF, PipelineDefinition types
- [ ] T059 [P] [US2] Create packages/autox/frontend/src/types/s3.ts with S3Credentials, S3Object, S3ListResponse types
- [ ] T060 [P] [US2] Create packages/autox/frontend/src/types/topology.ts with TopologyNode, TopologyEdge, TopologyGraph types
- [ ] T061 [P] [US2] Create packages/autox/frontend/src/types/common.ts with Namespace, Secret, User types
- [ ] T062 [P] [US2] Create packages/autox/frontend/src/types/index.ts barrel export for all types

### Implementation Tasks: Hooks

- [ ] T063 [P] [US2] Create packages/autox/frontend/src/hooks/usePipelineRuns.ts with React Query implementation
- [ ] T064 [P] [US2] Create packages/autox/frontend/src/hooks/usePipelineRunQuery.ts with polling support
- [ ] T065 [P] [US2] Create packages/autox/frontend/src/hooks/useTerminatePipelineRun.ts mutation hook
- [ ] T066 [P] [US2] Create packages/autox/frontend/src/hooks/useRetryPipelineRun.ts mutation hook
- [ ] T067 [P] [US2] Create packages/autox/frontend/src/hooks/useNotification.ts with toast notification wrapper
- [ ] T068 [P] [US2] Create packages/autox/frontend/src/hooks/useS3ListFiles.ts with Zod validation
- [ ] T069 [P] [US2] Create packages/autox/frontend/src/hooks/useS3FileUpload.ts mutation hook
- [ ] T070 [P] [US2] Create packages/autox/frontend/src/hooks/useNamespaces.ts query hook
- [ ] T071 [P] [US2] Create packages/autox/frontend/src/hooks/index.ts barrel export for all hooks

### Implementation Tasks: Components

- [ ] T072 [US2] Create packages/autox/frontend/src/components/FileExplorer/FileExplorer.tsx with breadcrumbs, search, pagination
- [ ] T073 [US2] Create packages/autox/frontend/src/components/PipelineRunsTable/PipelineRunsTable.tsx with sortable columns and actions
- [ ] T074 [US2] Create packages/autox/frontend/src/components/EmptyStates/PipelineEmptyState.tsx with conditional messaging
- [ ] T075 [US2] Create packages/autox/frontend/src/components/ToastNotification/ToastNotification.tsx (used by useNotification)
- [ ] T076 [US2] Create packages/autox/frontend/src/components/index.ts barrel export for all components

### Implementation Tasks: Utilities

- [ ] T077 [P] [US2] Create packages/autox/frontend/src/utils/pipelineRunState.ts with isPipelineRunning, isPipelineRunComplete helpers
- [ ] T078 [P] [US2] Create packages/autox/frontend/src/utils/validation.ts with isValidK8sName, isValidS3BucketName helpers
- [ ] T079 [P] [US2] Create packages/autox/frontend/src/utils/topology.ts with createTopologyNode, createTopologyEdge helpers
- [ ] T080 [P] [US2] Create packages/autox/frontend/src/utils/index.ts barrel export for all utilities

### Unit Test Tasks (AutoX Internal)

- [ ] T081 [P] [US2] Create packages/autox/frontend/src/hooks/__tests__/usePipelineRuns.spec.ts
- [ ] T082 [P] [US2] Create packages/autox/frontend/src/hooks/__tests__/useNotification.spec.ts
- [ ] T083 [P] [US2] Create packages/autox/frontend/src/hooks/__tests__/useS3ListFiles.spec.ts
- [ ] T084 [P] [US2] Create packages/autox/frontend/src/components/__tests__/FileExplorer.spec.tsx
- [ ] T085 [P] [US2] Create packages/autox/frontend/src/components/__tests__/PipelineRunsTable.spec.tsx
- [ ] T086 [P] [US2] Create packages/autox/frontend/src/utils/__tests__/validation.spec.ts

### Integration Tasks (Verify in AutoML/AutoRAG)

- [ ] T087 [US2] Update packages/automl/frontend to import and use usePipelineRuns hook from AutoX
- [ ] T088 [US2] Update packages/automl/frontend to import PipelineRunsTable component from AutoX
- [ ] T089 [US2] Update packages/automl/frontend to compose FileExplorer with AutoML-specific file actions
- [ ] T090 [US2] Update packages/autorag/frontend to import and use useS3ListFiles hook from AutoX
- [ ] T091 [US2] Update packages/autorag/frontend to compose FileExplorer with AutoRAG-specific navigation
- [ ] T092 [US2] Run npm run type-check from packages/automl/frontend and packages/autorag/frontend
- [ ] T093 [US2] Run npm test from packages/autox/frontend to verify all AutoX tests pass
- [ ] T093a [US2] Generate coverage report (npm test -- --coverage) and verify ≥80% overall coverage and 100% for critical utilities (fail if thresholds not met)
- [ ] T094 [US2] Start dev server (npm run dev) and manually verify AutoML/AutoRAG features work with AutoX primitives
- [ ] T094a [US2] Verify SC-003 threshold: Calculate UI LOC reduction percentage and confirm ≥50% (compare AutoML+AutoRAG frontend before/after extraction to AutoX frontend)

**US2 Independent Test**:
- ✅ AutoML component imports AutoX hook, composes it, and renders correctly
- ✅ AutoRAG component uses AutoX FileExplorer as building block and functionality works
- ✅ Both packages achieve their specific UI while sharing base logic from AutoX
- ✅ TypeScript compilation succeeds with no errors
- ✅ All AutoX unit tests pass
- ✅ UI duplication reduction meets ≥50% threshold (SC-003)

---

## Phase 5: User Story 7 - Refactor Existing Duplicate Code (P1 - Phase 1: Perfect Duplicates)

**Goal**: Migrate perfect duplicates (100% identical code) from AutoML/AutoRAG to AutoX

**Duration**: 16-20 hours  
**Dependencies**: US1 and US2 complete (infrastructure in place)  
**Deliverable**: Perfect duplicate code removed from AutoML/AutoRAG, functionality preserved

### Extraction Tasks: BFF Perfect Duplicates

- [ ] T095 [P] [US7] Extract cmd/helpers.go (100 LOC) from AutoML/AutoRAG to packages/autox/bff/cmd/helpers.go
- [ ] T096 [P] [US7] Extract core models (210 LOC total) - already done in US1, verify completeness
- [ ] T097 [P] [US7] Extract K8s portforward (240 LOC) - already done in US1, verify completeness
- [ ] T098 [P] [US7] Extract API helpers (289 LOC) - already done in US1, verify completeness

### Extraction Tasks: Frontend Perfect Duplicates

- [ ] T099 [US7] Extract FileExplorer component (1,880 LOC total) - already created in US2, verify completeness
- [ ] T100 [US7] Extract ToastNotification component (107 LOC) - already created in US2, verify completeness
- [ ] T101 [US7] Extract Topology components (221 LOC) to packages/autox/frontend/src/components/Topology/
- [ ] T102 [US7] Extract useNotification hook (167 LOC) - already created in US2, verify completeness
- [ ] T103 [US7] Extract S3 API (175 LOC) - already created in US2 hooks, verify completeness

### Cleanup Tasks

- [ ] T104 [US7] Remove duplicate BFF cmd/helpers.go from packages/automl/bff and packages/autorag/bff
- [ ] T105 [US7] Remove duplicate FileExplorer from packages/automl/frontend and packages/autorag/frontend
- [ ] T106 [US7] Remove duplicate ToastNotification from packages/automl/frontend and packages/autorag/frontend
- [ ] T107 [US7] Remove duplicate Topology components from packages/automl/frontend and packages/autorag/frontend
- [ ] T108 [US7] Remove duplicate useNotification from packages/automl/frontend and packages/autorag/frontend
- [ ] T109 [US7] Remove duplicate S3 API hooks from packages/automl/frontend and packages/autorag/frontend

### Verification Tasks

- [ ] T110 [US7] Run go test ./... from packages/automl/bff and verify all tests pass with AutoX imports
- [ ] T111 [US7] Run go test ./... from packages/autorag/bff and verify all tests pass with AutoX imports
- [ ] T112 [US7] Run npm test from packages/automl/frontend and verify all tests pass with AutoX imports
- [ ] T113 [US7] Run npm test from packages/autorag/frontend and verify all tests pass with AutoX imports
- [ ] T114 [US7] Run git diff to verify LOC reduction (~3,000 LOC removed from AutoML/AutoRAG combined)

**US7 Phase 1 Independent Test**:
- ✅ Perfect duplicate code (FileExplorer, useNotification, S3 API, etc.) removed from both packages
- ✅ AutoML and AutoRAG import from AutoX and tests pass
- ✅ Manual verification: Features work identically to before refactoring
- ✅ Codebase reduced by ~3,000 LOC

---

## Phase 6: User Story 3 - Handle Domain-Specific Logic (P2)

**Goal**: Establish patterns for domain-specific customization using handler files

**Duration**: 12-16 hours  
**Dependencies**: US1, US2 complete  
**Deliverable**: Clear examples of handler file pattern for AutoML and AutoRAG domain logic

### Implementation Tasks

- [ ] T115 [P] [US3] Create packages/automl/bff/internal/handlers/automl_validation.go with time-series specific validation
- [ ] T116 [P] [US3] Create packages/automl/bff/internal/handlers/model_registry_handler.go leveraging AutoX base utilities
- [ ] T117 [P] [US3] Create packages/autorag/bff/internal/handlers/rag_validation.go with RAG-specific validation
- [ ] T118 [P] [US3] Create packages/autorag/bff/internal/handlers/llamastack_handler.go with LlamaStack secret filtering
- [ ] T119 [P] [US3] Create packages/automl/frontend/src/handlers/timeSeriesValidation.ts extending AutoX validation utilities
- [ ] T120 [P] [US3] Create packages/autorag/frontend/src/handlers/ragParameterParser.ts extending AutoX parser utilities
- [ ] T121 [US3] Document handler file pattern in packages/autox/README.md with examples
- [ ] T122 [US3] Document when to use handler files vs AutoX in quickstart.md (decision criteria: 20-80% code reuse)

### Verification Tasks

- [ ] T123 [US3] Verify AutoML handler uses AutoX validation interface while implementing domain-specific rules
- [ ] T124 [US3] Verify AutoRAG handler extends AutoX parser without modifying AutoX code
- [ ] T125 [US3] Verify neither AutoML nor AutoRAG pollutes AutoX with package-specific logic

**US3 Independent Test**:
- ✅ AutoML handler file leverages AutoX validation interface while implementing custom validation logic
- ✅ AutoRAG handler extends AutoX parser for specific data formats without modifying AutoX
- ✅ Documentation clearly explains when to use handler files vs AutoX shared code

---

## Phase 7: User Story 7 - Refactor Code (P2 - Phase 2: High-Value Components)

**Goal**: Extract high-value components and integration clients with 95%+ similarity

**Duration**: 32-40 hours  
**Dependencies**: US1, US2, US3 complete  
**Deliverable**: Major components and clients extracted, significant LOC reduction

### Extraction Tasks: BFF Integration Clients

- [ ] T126 [US7] Extract S3 client full implementation (900 LOC) to packages/autox/bff/internal/integrations/s3/
- [ ] T127 [US7] Extract K8s integration full implementation (800 LOC) to packages/autox/bff/internal/integrations/kubernetes/
- [ ] T128 [US7] Extract Pipeline server client full implementation (630 LOC) to packages/autox/bff/internal/integrations/pipelineserver/
- [ ] T129 [US7] Update AutoML BFF to use full AutoX S3 client (remove AutoML-specific GetCSVSchema to handler)
- [ ] T130 [US7] Update AutoRAG BFF to use full AutoX K8s client
- [ ] T131 [US7] Update AutoML/AutoRAG to use full AutoX Pipeline server client

### Extraction Tasks: Frontend High-Value Components

- [ ] T132 [P] [US7] Extract Configure pages common logic (1,704 LOC shared) to packages/autox/frontend/src/components/ConfigureForm/
- [ ] T133 [P] [US7] Extract Leaderboard common logic (1,590 LOC shared) to packages/autox/frontend/src/components/Leaderboard/
- [ ] T134 [P] [US7] Extract Connection modals (433 LOC shared) to packages/autox/frontend/src/components/ConnectionModal/
- [ ] T135 [P] [US7] Extract Runs tables (123 LOC shared) to packages/autox/frontend/src/components/RunsTable/
- [ ] T136 [US7] Refactor AutoML ConfigurePage to compose AutoX ConfigureForm primitive
- [ ] T137 [US7] Refactor AutoRAG ConfigurePage to compose AutoX ConfigureForm primitive
- [ ] T138 [US7] Refactor AutoML Leaderboard to compose AutoX Leaderboard primitive with domain-specific columns
- [ ] T139 [US7] Refactor AutoRAG Leaderboard to compose AutoX Leaderboard primitive with RAG-specific metrics

### Cleanup & Verification Tasks

- [ ] T140 [US7] Remove duplicate Configure pages from AutoML and AutoRAG
- [ ] T141 [US7] Remove duplicate Leaderboard components from AutoML and AutoRAG
- [ ] T142 [US7] Remove duplicate Connection modals from AutoML and AutoRAG
- [ ] T143 [US7] Run go test ./... from all BFF packages and verify tests pass
- [ ] T144 [US7] Run npm test from all frontend packages and verify tests pass
- [ ] T145 [US7] Manual browser testing of AutoML/AutoRAG Configure pages with AutoX primitives
- [ ] T146 [US7] Manual browser testing of AutoML/AutoRAG Leaderboards with AutoX primitives
- [ ] T147 [US7] Run git diff to verify LOC reduction (~7,000 LOC removed)
- [ ] T147a [US7] Verify SC-002 threshold: Calculate BFF LOC reduction percentage and confirm ≥60% (compare AutoML+AutoRAG BFF before/after extraction to AutoX BFF)

**US7 Phase 2 Independent Test**:
- ✅ High-value components (Configure pages, Leaderboards, modals) extracted and working
- ✅ Integration clients (S3, K8s, Pipeline Server) fully extracted and functional
- ✅ AutoML and AutoRAG compose AutoX primitives correctly
- ✅ All tests pass, codebase reduced by ~7,000 LOC
- ✅ BFF duplication reduction meets ≥60% threshold (SC-002)

---

## Phase 8: User Story 4 - Strategy/DI Patterns for Customization (P3)

**Goal**: Implement strategy pattern and dependency injection for advanced customization scenarios

**Duration**: 16-20 hours  
**Dependencies**: US1, US2, US3 complete  
**Deliverable**: Working examples of strategy/DI pattern for AutoML and AutoRAG

### Precondition Validation Tasks

- [ ] T147b [US4] Validate service extraction candidates (PipelineRepository, SecretRepository) meet ≥80% duplication threshold per FR-007 research requirement before implementing strategy/DI pattern

### Implementation Tasks

- [ ] T148 [P] [US4] Create packages/autox/bff/internal/repositories/pipeline_repository.go with PipelineRepository interface
- [ ] T149 [US4] Implement packages/autox/bff/internal/repositories/pipeline_discovery.go with default discovery strategy
- [ ] T150 [US4] Create packages/autox/bff/internal/repositories/secret_repository.go with SecretRepository interface
- [ ] T151 [US4] Implement packages/autox/bff/internal/repositories/secret_filtering.go with base filtering strategy
- [ ] T152 [P] [US4] Create packages/automl/bff/internal/strategies/automl_pipeline_strategy.go implementing PipelineRepository
- [ ] T153 [P] [US4] Create packages/autorag/bff/internal/strategies/autorag_pipeline_strategy.go implementing PipelineRepository
- [ ] T154 [US4] Update packages/automl/bff/cmd/main.go to inject AutoML pipeline strategy into AutoX repository
- [ ] T155 [US4] Update packages/autorag/bff/cmd/main.go to inject AutoRAG pipeline strategy into AutoX repository
- [ ] T156 [US4] Document strategy/DI pattern in packages/autox/README.md with interface definition examples
- [ ] T157 [US4] Document when to use strategy/DI vs handler files in quickstart.md

### Verification Tasks

- [ ] T158 [US4] Verify AutoML pipeline strategy is injected and AutoX repository uses it correctly
- [ ] T159 [US4] Verify AutoRAG pipeline strategy is injected with different behavior
- [ ] T160 [US4] Verify both packages customize behavior without modifying AutoX code
- [ ] T161 [US4] Run integration tests to verify strategy injection works at runtime

**US4 Independent Test**:
- ✅ AutoX defines strategy interface (e.g., PipelineRepository)
- ✅ AutoML provides custom implementation and AutoX service uses it at runtime
- ✅ AutoRAG injects different implementation and service behavior is customized
- ✅ No changes needed to AutoX when strategies change

---

## Phase 9: User Story 7 - Refactor Code (Phases 3-4: Hooks, Utilities, Types)

**Goal**: Complete extraction of remaining shared code (hooks, utilities, types)

**Duration**: 32-40 hours  
**Dependencies**: All previous US7 phases complete  
**Deliverable**: Near-complete extraction (~13,300 LOC extracted across phases 3-4)

### Extraction Tasks: Hooks & State (Phase 3)

- [ ] T162 [P] [US7] Extract usePipelineDefinitions hook (already done in US2, verify completeness)
- [ ] T163 [P] [US7] Extract create pipeline run mutation hook to packages/autox/frontend/src/hooks/useCreatePipelineRun.ts
- [ ] T164 [P] [US7] Extract pipeline queries (180 LOC) to packages/autox/frontend/src/hooks/usePipelineQuery.ts
- [ ] T165 [P] [US7] Extract Results context generator (100 LOC) to packages/autox/frontend/src/contexts/ResultsContext.tsx
- [ ] T166 [P] [US7] Extract Pipeline API utilities (110 LOC) to packages/autox/frontend/src/api/pipelineApi.ts
- [ ] T167 [P] [US7] Extract error handling utilities (88 LOC) to packages/autox/frontend/src/utils/errorHandling.ts

### Extraction Tasks: Utilities & Types (Phase 4)

- [ ] T168 [P] [US7] Verify all types extracted to packages/autox/frontend/src/types/ (394 LOC - done in US2)
- [ ] T169 [P] [US7] Verify all frontend utilities extracted to packages/autox/frontend/src/utils/ (118 LOC - done in US2)
- [ ] T170 [P] [US7] Extract BFF helpers (371 LOC) to packages/autox/bff/internal/utils/helpers.go
- [ ] T171 [P] [US7] Extract repository patterns (150 LOC) to packages/autox/bff/internal/repositories/

### Cleanup & Verification

- [ ] T172 [US7] Remove duplicate hooks from AutoML and AutoRAG frontends
- [ ] T173 [US7] Remove duplicate utilities from AutoML and AutoRAG (BFF and frontend)
- [ ] T174 [US7] Remove duplicate types from AutoML and AutoRAG frontends
- [ ] T175 [US7] Run all tests across AutoX, AutoML, AutoRAG (go test ./..., npm test)
- [ ] T176 [US7] Verify type-checking passes across all packages (tsc --noEmit, go build)
- [ ] T177 [US7] Run git diff to verify cumulative LOC reduction (~13,300 LOC for phases 3-4)

**US7 Phases 3-4 Independent Test**:
- ✅ All shared hooks, utilities, and types extracted to AutoX
- ✅ AutoML and AutoRAG import from AutoX and tests pass
- ✅ TypeScript and Go compilation succeeds
- ✅ Codebase reduced by ~13,300 LOC (cumulative with previous phases)

---

## Phase 10: Polish & Cross-Cutting Concerns

**Goal**: Documentation, CI/CD updates, final verification

**Duration**: 8-12 hours  
**Dependencies**: All user stories complete  
**Deliverable**: Production-ready AutoX package with comprehensive documentation

### Documentation Tasks

- [ ] T178 [P] Update packages/autox/README.md with complete API documentation and usage examples
- [ ] T179 [P] Update contracts/bff-api.md to reflect all extracted BFF APIs
- [ ] T180 [P] Update contracts/ui-api.md to reflect all extracted UI APIs
- [ ] T181 [P] Create MIGRATION.md guide for future package migrations to AutoX pattern
- [ ] T182 Update root BOOKMARKS.md to reference AutoX package documentation

### CI/CD Tasks

- [ ] T183 Update .github/workflows CI configuration to run tests for AutoX package
- [ ] T184 Add type-checking step to CI that verifies all packages compile together (go build ./..., tsc --noEmit)
- [ ] T184a Verify CI catches breaking changes by introducing test breaking change in AutoX (e.g., rename exported interface), confirm AutoML/AutoRAG builds fail in CI, then revert change
- [ ] T185 Add linting step for AutoX (ESLint for frontend, golint for BFF)
- [ ] T186 Verify CI catches breaking changes when AutoX is modified

### Final Verification Tasks

- [ ] T187 Run full test suite across all packages (AutoX, AutoML, AutoRAG)
- [ ] T188 Perform manual E2E testing of AutoML features using AutoX
- [ ] T189 Perform manual E2E testing of AutoRAG features using AutoX
- [ ] T190 Generate bundle size report and verify AutoX singleton reduces duplication
- [ ] T191 Run final git diff to verify total LOC reduction (~31,500 LOC target)

**Polish Phase Test**:
- ✅ All documentation up-to-date and accurate
- ✅ CI/CD runs tests for AutoX and catches breaking changes
- ✅ Manual E2E testing confirms all features work
- ✅ Bundle size analysis shows optimization from singleton pattern

---

## Completion Checklist

### User Story Acceptance Criteria

**US1 - Import Shared BFF Utilities**:
- [x] AutoX exports BFF interfaces, utilities, clients
- [x] AutoML BFF handler imports and uses AutoX utilities
- [x] AutoRAG BFF handler imports and uses AutoX clients
- [x] Behavior is consistent across both packages

**US2 - Compose UI Primitives**:
- [x] AutoX provides primitive hooks (usePipelineRuns, useNotification, etc.)
- [x] AutoX provides low-level components (FileExplorer, PipelineRunsTable)
- [x] AutoML composes AutoX primitives into feature components
- [x] AutoRAG composes AutoX primitives differently for RAG-specific UI

**US3 - Handle Domain-Specific Logic**:
- [x] AutoML handler file extends AutoX validation interface
- [x] AutoRAG handler customizes AutoX parser for specific formats
- [x] Neither package pollutes AutoX with domain-specific logic

**US4 - Strategy/DI Patterns**:
- [x] AutoX defines strategy interfaces (PipelineRepository, etc.)
- [x] AutoML injects custom strategy implementation
- [x] AutoRAG injects different strategy implementation
- [x] Services use injected strategies correctly

**US5 - Development Environment Setup**:
- [x] npm workspaces link AutoX, AutoML, AutoRAG frontends
- [x] Go workspaces link AutoX, AutoML, AutoRAG BFFs
- [x] Changes to AutoX reflected in consumers without manual linking

**US6 - Module Federation Runtime Singleton**:
- [x] AutoX configured as Module Federation shared singleton
- [x] Only one instance of AutoX exists in browser runtime
- [x] Bundle size optimized with no duplicate AutoX code

**US7 - Refactor Existing Duplicate Code**:
- [x] Phase 1: Perfect duplicates extracted (~3,000 LOC)
- [x] Phase 2: High-value components extracted (~7,000 LOC)
- [x] Phase 3: Hooks & state extracted (~2,000 LOC)
- [x] Phase 4: Utilities & types extracted (~1,300 LOC)
- [x] Total: ~13,300 LOC extracted (Phases 1-4)

### Technical Verification

**Build & Type-Checking**:
- [ ] `go build ./...` succeeds from repo root (all BFF packages compile)
- [ ] `tsc --noEmit` succeeds for AutoX, AutoML, AutoRAG frontends
- [ ] `npm run lint` passes for all frontend packages
- [ ] `golint ./...` passes for all BFF packages

**Testing**:
- [ ] `go test ./...` passes for AutoX, AutoML, AutoRAG BFFs
- [ ] `npm test` passes for AutoX, AutoML, AutoRAG frontends
- [ ] All AutoX unit tests pass (100% of utilities, 80%+ overall)
- [ ] Integration tests in AutoML/AutoRAG verify composition works

**Runtime Verification**:
- [ ] `npm run dev` starts successfully
- [ ] AutoML features work correctly with AutoX imports
- [ ] AutoRAG features work correctly with AutoX imports
- [ ] DevTools confirms single AutoX instance in runtime
- [ ] Bundle analysis shows reduced duplication

**Documentation**:
- [ ] packages/autox/README.md is complete and accurate
- [ ] contracts/bff-api.md documents all BFF exports
- [ ] contracts/ui-api.md documents all UI exports
- [ ] quickstart.md has clear setup and usage examples

---

## Metrics & Success Criteria

### Code Reduction Metrics

| Metric | Target | Actual | Status |
|--------|--------|--------|--------|
| Total LOC Extracted | ~31,500 | TBD | ⏳ In Progress |
| BFF Duplication Reduction | 77% | TBD | ⏳ In Progress |
| Frontend Duplication Reduction | 94% | TBD | ⏳ In Progress |
| Test Code Shared | 85% | TBD | ⏳ In Progress |

### Quality Metrics

| Metric | Target | Actual | Status |
|--------|--------|--------|--------|
| AutoX Unit Test Coverage | 80%+ | TBD | ⏳ In Progress |
| AutoX Critical Utils Coverage | 100% | TBD | ⏳ In Progress |
| TypeScript Compilation | 0 errors | TBD | ⏳ In Progress |
| Go Compilation | 0 errors | TBD | ⏳ In Progress |
| Linting Errors | 0 | TBD | ⏳ In Progress |

### Performance Metrics

| Metric | Target | Actual | Status |
|--------|--------|--------|--------|
| Bundle Size Impact | <15KB gzipped | TBD | ⏳ In Progress |
| BFF Utility Overhead | <1ms | TBD | ⏳ In Progress |
| Module Federation Instances | 1 (singleton) | TBD | ⏳ In Progress |

---

## Notes for Implementation

### Parallel Execution Strategy

**Phase 2 (Foundational)**: Execute US5 and US6 in parallel (different files)

**Phase 3-4 (P1 Stories)**: 
- Execute US1 (BFF) and US2 (UI) in parallel (completely independent)
- US7 Phase 1 extraction can happen in parallel with US1/US2 (different files)

**Phase 5-6 (P2 Stories)**:
- US3 can start as soon as US1/US2 complete
- US7 Phase 2 requires US1/US2/US3 complete (depends on composition patterns)

**Phase 7 (P3 Stories)**:
- US4 requires understanding of handler file pattern from US3

### Risk Mitigation

**High-Risk Tasks** (require careful review):
- T129-T131: Extracting full integration clients (ensure domain-specific logic stays in handlers)
- T132-T139: Refactoring Configure pages and Leaderboards (complex components, high LOC)
- T148-T161: Strategy/DI pattern implementation (architectural complexity)

**Recommended Review Points**:
- After US1 complete: Review BFF API contract
- After US2 complete: Review UI API contract  
- After US3 complete: Review handler file pattern documentation
- After US7 Phase 2: Review extraction quality and test coverage
- Before Polish phase: Full code review of AutoX package

### Testing Strategy

**TDD Approach** (Optional):
- Tests not explicitly requested in spec, so not included as separate tasks
- If TDD desired: Write T044-T049 (BFF tests) before T025-T043 (BFF implementation)
- If TDD desired: Write T081-T086 (UI tests) before T058-T080 (UI implementation)

**Integration Testing**:
- Each user story includes integration verification tasks in consumer packages
- Manual browser testing required for US2 (UI primitives) and US7 Phase 2 (high-value components)

---

**Total Tasks**: 202 (includes T094a, T147a, T147b, T184a added via /speckit.analyze remediation)
**Parallelizable Tasks**: 68 (marked with `[P]`)  
**Sequential Dependencies**: Clearly marked in phase dependencies  
**Estimated Completion**: 8-12 weeks with 2-3 developers working full-time

**Next Step**: Begin with Phase 1 (Setup) followed by Phase 2 (Foundational: US5 + US6), then proceed to P1 user stories (US1, US2, US7 Phase 1) which deliver immediate value.
