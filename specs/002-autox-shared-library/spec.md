# Feature Specification: AutoX Shared Library Package

**Feature Branch**: `002-autox-shared-library`  
**Created**: 2026-04-22  
**Status**: Draft  
**Input**: User description: "create a shared core library called autox that lives as package under odh-dashboard/packages which exports common code for both ui and bff for the automl and autorag packages and refactor them to consume autox

at the bff layer, all shared logic (interfaces, utilities, clients, etc.) should be consumed from autox, where logic diverges, this should be handled by handler files (eg. for domain specific parsing/validation), if this is insufficient and behaviour needs to be customized further, leverage strategy/DI patterns, prioritize the most reusable logic like low-level interfaces, utilities, clients, etc. first and gradually move towards higher-level like services (for exporting a full service, research is required to ensure there is high duplication currently and for the forseeable future, and whether it would be worth introducing the complexity of DI), leverage go workspaces for local development

at the ui layer, all shared logic (interfaces, hooks, components, etc.) should also be consumed from autox, use the composition pattern where autox provides low-level flexible primitives while automl and autorag compose them to create higher-level components, avoid monolithic shared components with extensive prop variants, leverage a module federation singleton at runtime, leverage npm workspaces for local development"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Import Shared BFF Utilities (Priority: P1)

As a developer working on the AutoML or AutoRAG BFF, I want to import shared interfaces, utilities, and clients from the AutoX package so that I don't duplicate code and maintain consistency across both packages.

**Why this priority**: This is the foundation of the shared library - without it, there's no value. It eliminates the most common source of duplication and inconsistency.

**Independent Test**: Can be fully tested by creating a simple shared utility in AutoX, importing it in AutoML BFF, and verifying it works correctly. Delivers immediate value by reducing duplication.

**Acceptance Scenarios**:

1. **Given** AutoX exports a shared interface or utility, **When** an AutoML BFF handler imports it, **Then** the import succeeds and the utility functions correctly
2. **Given** AutoX exports a client wrapper, **When** AutoRAG BFF imports and uses it, **Then** the client operations work as expected
3. **Given** both AutoML and AutoRAG import the same utility from AutoX, **When** they use it, **Then** behavior is consistent across both packages

---

### User Story 2 - Compose UI Primitives (Priority: P1)

As a developer building UI features in AutoML or AutoRAG, I want to import low-level primitive hooks and components from AutoX and compose them into feature-specific components, so that I can build consistent UIs without duplicating base functionality.

**Why this priority**: UI consistency and avoiding duplication is as critical as BFF sharing. This directly impacts development velocity and user experience consistency.

**Independent Test**: Can be fully tested by creating a primitive hook in AutoX, composing it into an AutoML-specific component, and verifying the composed component works. Delivers value by enabling consistent UI patterns.

**Acceptance Scenarios**:

1. **Given** AutoX provides a primitive hook, **When** AutoML composes it into a feature component, **Then** the composed component works and the hook's functionality is available
2. **Given** AutoX provides a low-level UI component, **When** AutoRAG uses it as a building block, **Then** the component renders correctly and props are properly typed
3. **Given** both AutoML and AutoRAG compose the same primitive differently, **When** they render, **Then** each achieves its specific UI while sharing the base logic

---

### User Story 3 - Handle Domain-Specific Logic (Priority: P2)

As a developer implementing AutoML-specific or AutoRAG-specific features, I want to use handler files for domain-specific parsing and validation while leveraging shared AutoX utilities, so that I can customize behavior without polluting the shared library.

**Why this priority**: This ensures the shared library stays clean and focused. Without proper separation, AutoX would become a dumping ground for all logic, defeating its purpose.

**Independent Test**: Can be tested by creating a shared validation utility in AutoX, then implementing domain-specific validators in AutoML and AutoRAG handler files that extend or customize the base utility. Delivers value by maintaining clean boundaries.

**Acceptance Scenarios**:

1. **Given** AutoX provides a generic validation interface, **When** AutoML creates a handler file with custom validation logic, **Then** the handler can leverage the interface while implementing domain-specific rules
2. **Given** a shared parser utility exists in AutoX, **When** AutoRAG extends it in a handler for specific data formats, **Then** the extension works without modifying AutoX
3. **Given** both packages have different parsing requirements, **When** they implement handlers, **Then** neither pollutes AutoX with package-specific logic

---

### User Story 4 - Strategy/DI Patterns for Customization (Priority: P3)

As a developer needing to customize shared service behavior, I want to use strategy patterns or dependency injection so that AutoML and AutoRAG can provide different implementations of the same interface without forking the shared code.

**Why this priority**: This is needed for advanced scenarios where handler files are insufficient. While important for flexibility, it's not needed for initial MVP functionality.

**Independent Test**: Can be tested by defining a strategy interface in AutoX, implementing different strategies in AutoML and AutoRAG, and injecting them at runtime. Delivers value for complex customization needs.

**Acceptance Scenarios**:

1. **Given** AutoX defines a strategy interface, **When** AutoML provides a custom implementation, **Then** the AutoX service can use the AutoML strategy at runtime
2. **Given** a service in AutoX accepts injected dependencies, **When** AutoRAG injects its own implementations, **Then** the service uses the AutoRAG-specific behavior
3. **Given** both packages inject different strategies, **When** each runs, **Then** behavior is customized correctly without modifying AutoX

---

### User Story 5 - Development Environment Setup (Priority: P0 - Foundational)

As a developer setting up the repository locally, I want to use npm workspaces for UI and Go workspaces for BFF so that all packages (AutoX, AutoML, AutoRAG) are properly linked and I can develop with live changes without manual linking.

**Why this priority**: This is a foundational blocking dependency for US1/US2 (P1 stories). While ranked P0 for dependency ordering, it delivers lower immediate business value than P1 stories. Without workspace setup, P1 stories cannot function. Developer experience is critical for productivity - without proper workspace setup, development becomes painful with manual linking and build steps.

**Independent Test**: Can be tested by running `npm install` at repo root and verifying all UI packages are linked, then using Go workspace commands to verify BFF packages are linked. Delivers value through improved DX.

**Acceptance Scenarios**:

1. **Given** the repo is cloned, **When** `npm install` is run at the root, **Then** AutoX, AutoML, and AutoRAG UI packages are linked via npm workspaces
2. **Given** the Go workspace is configured, **When** a developer runs `go work use`, **Then** AutoX, AutoML, and AutoRAG BFF packages are linked
3. **Given** a change is made to AutoX, **When** AutoML or AutoRAG is built, **Then** the change is reflected without manual linking

---

### User Story 6 - Module Federation Runtime Singleton (Priority: P2)

As a user of the dashboard, I want AutoX to be loaded as a Module Federation singleton at runtime so that only one instance of shared code exists in the browser, reducing bundle size and ensuring consistent state.

**Why this priority**: Performance and correctness are important, but the feature can work without this optimization. It becomes critical as the shared library grows.

**Independent Test**: Can be tested by configuring AutoX as a Module Federation shared singleton, loading AutoML and AutoRAG, and verifying in DevTools that only one instance of AutoX is loaded. Delivers value through reduced bundle size.

**Acceptance Scenarios**:

1. **Given** AutoX is configured as a Module Federation singleton, **When** both AutoML and AutoRAG load, **Then** only one instance of AutoX exists in the runtime
2. **Given** the user navigates between AutoML and AutoRAG features, **When** the app loads, **Then** bundle size is optimized with no duplicate AutoX code

---

### User Story 7 - Refactor Existing Duplicate Code (Priority: P1)

As a maintainer, I want to identify and refactor existing duplicate code between AutoML and AutoRAG into AutoX so that the codebase is DRY and future changes only need to happen in one place.

**Why this priority**: This is the migration path - without refactoring existing code, AutoX provides no immediate value. It's as critical as creating the initial shared utilities.

**Independent Test**: Can be tested by identifying a specific duplicate utility, moving it to AutoX, updating imports in both packages, and verifying tests pass. Delivers value by immediately reducing duplication.

**Acceptance Scenarios**:

1. **Given** duplicate BFF utilities exist in AutoML and AutoRAG, **When** they are refactored into AutoX, **Then** both packages import from AutoX and tests pass
2. **Given** duplicate UI hooks exist, **When** they are moved to AutoX primitives, **Then** both packages compose them and functionality is preserved
3. **Given** shared logic is identified, **When** it's evaluated for extraction, **Then** it meets the criteria (high duplication, stable interface, no excessive customization needed)

**Note**: See `research.md` for exhaustive analysis of ~31,500 LOC extractable code across 5 phases with detailed file-by-file extraction roadmap.

---

### Edge Cases

- **What happens when AutoML needs slightly different behavior than AutoRAG for a shared utility?**
  - Use handler files for small differences (e.g., different validation rules)
  - Use strategy/DI patterns for larger differences (e.g., different processing algorithms)
  - If differences are fundamental, keep separate implementations (not everything should be shared)

- **How do we handle breaking changes in AutoX that affect both AutoML and AutoRAG?**
  - Use synchronized versioning: AutoX, AutoML, and AutoRAG are always updated together in the same PR
  - All three packages are tested together in CI with type-checking that fails on errors (TypeScript `tsc --noEmit` for UI, `go build` for BFF)
  - Breaking changes are coordinated and applied atomically across all affected packages
  - CI type-checking catches breaking changes before PR merge, preventing production issues

- **What if a feature seems shared initially but later diverges significantly?**
  - Start with shared implementation
  - If divergence exceeds 30% of logic, move to handler files
  - If divergence is fundamental, extract back to separate packages with a note
  - Document the decision to help future developers

- **How do we determine when logic should be in AutoX vs handler files?**
  - AutoX: Pure utilities, interfaces, low-level clients (>80% code reuse)
  - Handler files: Domain-specific parsing, validation, custom business logic (20-80% code reuse)
  - Package-specific: Highly specialized logic (<20% code reuse)

- **What if a third package (not AutoML/AutoRAG) needs AutoX functionality?**
  - AutoX is not intended to be a general utility package
  - Other packages (eval-hub, etc.) will not be consuming AutoX
  - Keep AutoX focused on highly-reusable utilities to avoid scope creep

- **How do we decide if a UI component belongs in AutoX or stays package-specific?**
  - Think composable first: Design components as single-responsibility primitives by default, even if only one package uses them initially
  - AutoX components: Composable building blocks with clear boundaries and minimal coupling to package-specific concerns
  - Package-specific: Highly specialized orchestration or components tightly coupled to one package's unique domain logic
  - Red flags for AutoX: Extensive conditional rendering based on props, package-specific business logic embedded, "do-everything" designs

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: AutoX package MUST be created under `odh-dashboard/packages/autox` with both `frontend/` and `bff/` subdirectories following standard monorepo package structure: `autox/frontend/src/{hooks,components,utils,types}` for UI and `autox/bff/pkg/{api,models,integrations,utils}` for BFF (consistent with existing packages like gen-ai, automl, autorag)

- **FR-002**: AutoX BFF MUST export shared interfaces, utilities, and clients that are consumed by AutoML and AutoRAG BFF layers

- **FR-003**: AutoX UI MUST export low-level primitive hooks and components that can be composed by AutoML and AutoRAG UI layers

- **FR-004**: AutoML and AutoRAG MUST handle domain-specific logic through handler files within their own packages, not in AutoX

- **FR-005**: AutoX MUST support strategy patterns and dependency injection for advanced customization scenarios where handler files are insufficient

- **FR-006**: AutoX MUST prioritize most reusable logic first (≥80% code reuse threshold): low-level interfaces → utilities → clients → services. Extraction follows 5-phase roadmap detailed in research.md: Phase 1 (perfect duplicates), Phase 2 (high-value components), Phase 3 (hooks/state), Phase 4 (utilities/types), Phase 5 (advanced patterns)

- **FR-007**: Before exporting a complete service from AutoX, research MUST be conducted to ensure high duplication (≥80% code reuse) exists currently and for the foreseeable future, and that the complexity of DI is justified

- **FR-008**: AutoX UI MUST be configured as a federated shared dependency (not a remote) in the webpack Module Federation config of both AutoML and AutoRAG with `singleton: true` and `requiredVersion` to ensure single runtime instance and automatic version resolution

- **FR-009**: AutoX MUST NOT contain monolithic shared components with extensive prop variants - all shared components MUST be low-level primitives focused on single responsibility and composability. Components with >10 props or handling >2 unrelated concerns require justification and may indicate the need for decomposition. Components should be composable building blocks, but avoid "kitchen sink" designs (e.g., a single component handling both data fetching AND rendering AND validation AND error handling)

- **FR-010**: Repository MUST use npm workspaces at the root to link AutoX, AutoML, and AutoRAG UI packages for local development

- **FR-011**: Repository MUST use Go workspaces to link AutoX, AutoML, and AutoRAG BFF packages for local development

- **FR-012**: Existing duplicate code in AutoML and AutoRAG MUST be refactored to consume from AutoX where appropriate

- **FR-013**: AutoX MUST have comprehensive unit tests for all exported utilities, hooks, and components. Integration tests verifying AutoX primitives compose correctly MUST be implemented in AutoML and AutoRAG test suites (not in AutoX to avoid circular dependencies)

- **FR-014**: AutoML and AutoRAG MUST update their imports to consume shared logic from AutoX instead of local duplicates

- **FR-015**: AutoX package MUST follow the same structure and conventions as other packages in the monorepo (TypeScript, Go, linting, etc.)

- **FR-016**: AutoX package MUST include comprehensive documentation before implementation begins: data-model.md (key entities and interfaces), contracts/bff-api.md (all BFF exported interfaces, models, utilities), contracts/ui-api.md (all UI exported hooks, components, types), and quickstart.md (developer setup guide, import patterns, testing strategy)

### Key Entities

- **AutoX Package**: The shared library package containing both BFF and UI code for related functionality shared between AutoML and AutoRAG

- **BFF Shared Logic**: Includes interfaces, utilities, and clients used by both AutoML and AutoRAG backend services

- **UI Primitives**: Low-level hooks and components that provide flexible building blocks for composing feature-specific UIs

- **Handler Files**: Package-specific files in AutoML and AutoRAG that handle domain-specific parsing, validation, and business logic

- **Strategy/DI Patterns**: Design patterns used for advanced customization where behavior needs to vary between AutoML and AutoRAG

- **Workspace Configuration**: npm and Go workspace files that link packages for local development

## Success Criteria *(mandatory)*

### Measurable Outcomes

**Note**: Success criteria are categorized as:
- **Buildable/Verifiable** (SC-001, SC-002, SC-003, SC-006, SC-007, SC-008): Can be verified during implementation via automated tests or measurements
- **Post-Launch Outcome Metrics** (SC-004, SC-005, SC-009, SC-010): Require user research, telemetry, or post-deployment measurement

- **SC-001**: Developers can import and use AutoX utilities in AutoML or AutoRAG with zero configuration beyond initial workspace setup

- **SC-002**: Code duplication between AutoML and AutoRAG BFF layers is reduced by at least 60% (measured by: count LOC in AutoML BFF + AutoRAG BFF before extraction, count LOC removed from both packages after extraction to AutoX, calculate reduction % = removed LOC / total original LOC)

- **SC-003**: Code duplication between AutoML and AutoRAG UI layers is reduced by at least 50% (measured by: count LOC in AutoML frontend + AutoRAG frontend before extraction, count LOC removed from both packages after extraction to AutoX, calculate reduction % = removed LOC / total original LOC)

- **SC-004**: Development environment setup time is reduced to under 5 minutes (single `npm install` from repo root links all packages)

- **SC-005**: UI components can be composed from AutoX primitives in under 30 minutes for standard features

- **SC-006**: Only one instance of AutoX loads in the browser when both AutoML and AutoRAG are active (verified via webpack bundle analysis)

- **SC-007**: All AutoX exports have 100% unit test coverage for critical utilities (utilities handling authentication, S3 operations, K8s client interactions, or data validation) and at least 80% for other code. AutoML and AutoRAG each have integration tests covering their composition of AutoX primitives for key user flows

- **SC-008**: Breaking changes in AutoX are caught at compile time in AutoML and AutoRAG (TypeScript/Go type checking), and all three packages are updated together atomically in the same PR with coordinated testing

- **SC-009**: 90% of developers successfully set up local development without manual linking steps (based on team feedback)

- **SC-010**: Changes to AutoX are reflected in AutoML and AutoRAG without rebuilding (via workspace hot-reload for UI, go work sync for BFF)

## Assumptions

- The AutoML and AutoRAG packages already have significant code duplication that justifies creating a shared library
- Both packages will continue to share fundamental patterns (pipeline runs, s3, leaderboard, etc.) for the foreseeable future
- The monorepo already uses Turbo or similar tooling that supports npm workspaces
- Go 1.18+ is available for Go workspace support
- Module Federation is already configured in the monorepo for other packages
- Developers are familiar with composition patterns and strategy/DI concepts
- AutoX, AutoML, and AutoRAG will use synchronized versioning and be updated together atomically (no independent versioning needed)

## Scope Boundaries

### In Scope

- Creating the AutoX package structure (BFF and UI)
- Identifying and refactoring duplicate code from AutoML and AutoRAG
- Configuring npm and Go workspaces
- Setting up Module Federation singleton for AutoX UI
- Creating low-level primitive hooks and components
- Exporting shared BFF interfaces, utilities, and clients
- Documenting patterns for handler files and strategy/DI usage
- Unit testing all AutoX exports
- Adding integration tests in AutoML and AutoRAG to verify AutoX primitives compose correctly

### Out of Scope

- Creating high-level, opinionated components with many variants
- Migrating 100% of code (some package-specific logic will remain)
- Performance optimization beyond Module Federation singleton
- Creating a separate npm package published to registry (monorepo package only)
- Backward compatibility with old import paths (breaking change is acceptable)
- Automated migration tooling (manual refactoring is acceptable)

## Clarifications

### Session 2026-04-22

- Q: What internal directory structure should AutoX follow? → A: Standard monorepo package structure (autox/frontend/src/{hooks,components,utils,types}, autox/bff/pkg/{api,models,utils}) matching existing packages
- Q: How should AutoX UI be configured in Module Federation? → A: AutoX as federated shared dependency in AutoML/AutoRAG webpack config (singleton: true, requiredVersion) - standard shared lib pattern
- Q: What versioning strategy should be used for AutoX and its consumers? → A: Synchronized versioning: AutoX, AutoML, AutoRAG always updated together in same PR, tested together, no independent versions
- Q: Where should integration tests for AutoX live? → A: Integration tests in AutoML/AutoRAG packages that verify AutoX primitives compose correctly - tests at point of use, no circular deps
- Q: How should breaking changes in AutoX interfaces be detected before they break AutoML/AutoRAG in production? → A: CI type-checking all packages together (fail on errors)
- Q: What specific characteristics define a "low-level primitive" component in AutoX versus a "monolithic component with extensive prop variants" that should be avoided? → A: Single responsibility, composable, no strict prop limit
