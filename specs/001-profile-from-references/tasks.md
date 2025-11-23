---
description: "Task list for Profile Generation from References"
---

# Tasks: Profile Generation from References

**Input**: Design documents from `/specs/001-profile-from-references/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml

**Tests**: Tests are **NOT OPTIONAL**. The goal is to have strong test coverage to guide our implementation.
- Generate unit level tests where you can.
- Keep tests generally in AAA format (Arrange-Act-Assert).
- Try not to leak too much implementation details into the individual unit tests.
- Prefer unit tests over integration tests.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- **Include exact file paths in descriptions**

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Verify and ensure Go module setup in `apps/backend/go.mod`
- [x] T002 Verify and ensure Next.js app setup in `apps/frontend/package.json`
- [x] T003 [P] Configure environment variables loading in `apps/backend/main.go`
- [x] T004 [P] Setup basic logging configuration in `apps/backend/src/utils/logger.go`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T005 Setup Postgres database connection logic in `apps/backend/src/db/db.go`
- [x] T006 Create database migrations for User, ReferenceLetter, Company, Experience, Skill in `apps/backend/src/db/migrations/000001_init_schema.up.sql`
- [x] T007 [P] Define Go structs for UserProfile, ReferenceLetter, CompanyEntry, WorkExperience in `apps/backend/src/models/types.go`
- [x] T008 [P] Implement PDF text extraction service wrapper (ledongthuc/pdf) in `apps/backend/src/services/extractor/pdf.go`
- [x] T009 [P] Create unit tests for PDF extraction in `apps/backend/src/services/extractor/pdf_test.go`
- [x] T010 [P] Implement OpenAI client wrapper and configuration in `apps/backend/src/services/llm/client.go`
- [x] T011 [P] Create unit tests for OpenAI client in `apps/backend/src/services/llm/client_test.go`
- [x] T012 Setup API router (Chi) and middleware in `apps/backend/src/api/router.go`
- [x] T013 [P] Configure global error handling middleware in `apps/backend/src/api/middleware/errors.go`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Create Profile from References (Priority: P1) 🎯 MVP

**Goal**: Users can upload PDF references and have a profile automatically created.

**Independent Test**: Upload a sample PDF reference letter via the API or UI and verify that `UserProfile` and related tables are populated with extracted data.

### Implementation for User Story 1

- [x] T014 [P] [US1] Define LLM prompts for data extraction in `apps/backend/src/services/llm/prompts.go`
- [x] T015 [US1] Implement local file storage service for uploaded files in `apps/backend/src/services/storage/local.go`
- [x] T016 [P] [US1] Create unit tests for local storage in `apps/backend/src/services/storage/local_test.go`
- [x] T017 [US1] Implement core extraction service (PDF text -> LLM -> JSON) in `apps/backend/src/services/profile/extractor.go`
- [x] T018 [P] [US1] Create unit tests for core extraction service in `apps/backend/src/services/profile/extractor_test.go`
- [x] T019 [US1] Implement logic to merge extracted data into UserProfile (handling name mismatch/updates) in `apps/backend/src/services/profile/service.go`
- [x] T020 [P] [US1] Create unit tests for profile merge logic in `apps/backend/src/services/profile/service_test.go`
- [x] T021 [US1] Create `POST /api/upload` handler in `apps/backend/src/api/handlers/upload.go`
- [x] T022 [US1] Implement frontend API client upload method in `apps/frontend/services/api.ts`
- [x] T023 [P] [US1] Create File Upload component with progress state in `apps/frontend/components/upload/FileUploader.tsx`
- [x] T024 [US1] Create Upload Page with drag-and-drop zone in `apps/frontend/app/upload/page.tsx`

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - View Enhanced Profile (Priority: P1)

**Goal**: Users can view their generated profile including experience, skills, and employer feedback.

**Independent Test**: access `/dashboard` and verify the profile data matches the database content, including "Credibility" section.

### Implementation for User Story 2

- [x] T025 [P] [US2] Implement database queries to fetch full profile with relations in `apps/backend/src/db/queries/profile.go`
- [x] T026 [P] [US2] Create unit/integration tests for profile queries in `apps/backend/src/db/queries/profile_test.go`
- [x] T027 [US2] Create `GET /api/profile` handler in `apps/backend/src/api/handlers/profile.go`
- [x] T028 [US2] Implement frontend API client getProfile method in `apps/frontend/services/api.ts`
- [x] T029 [P] [US2] Create Profile Header component (Name, Summary) in `apps/frontend/components/profile/ProfileHeader.tsx`
- [x] T030 [P] [US2] Create Experience List component with "Verified" badges in `apps/frontend/components/profile/ExperienceList.tsx`
- [x] T031 [P] [US2] Create Credibility/Feedback section component in `apps/frontend/components/profile/CredibilitySection.tsx`
- [x] T032 [P] [US2] Create Skills List component linked to experiences in `apps/frontend/components/profile/SkillsList.tsx`
- [x] T033 [US2] Assemble Dashboard Page in `apps/frontend/app/dashboard/page.tsx`

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Download Standard CV (Priority: P2)

**Goal**: Users can download a PDF version of their profile.

**Independent Test**: Click "Download CV" and verify a valid PDF is downloaded containing profile info.

### Implementation for User Story 3

- [x] T034 [P] [US3] Implement PDF generation service (converting profile struct to PDF layout) in `apps/backend/src/services/generator/cv_pdf.go`
- [x] T035 [P] [US3] Create unit tests for PDF generation in `apps/backend/src/services/generator/cv_pdf_test.go`
- [x] T036 [US3] Create `GET /api/profile/cv` handler in `apps/backend/src/api/handlers/cv.go`
- [x] T037 [US3] Add "Download CV" button to Profile Header in `apps/frontend/components/profile/ProfileHeader.tsx`
- [x] T038 [US3] Implement download trigger in frontend `apps/frontend/app/dashboard/page.tsx`

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: User Story 4 - Tailor Profile to Job Description (Priority: P2)

**Goal**: Users can customize their profile/CV for a specific job description.

**Independent Test**: Submit a Job Description text and verify the returned profile preview highlights relevant skills/experiences.

### Implementation for User Story 4

- [x] T039 [P] [US4] Define LLM prompts for profile tailoring/re-ranking in `apps/backend/src/services/llm/prompts.go`
- [x] T040 [US4] Implement tailoring service logic in `apps/backend/src/services/profile/tailor.go`
- [x] T041 [P] [US4] Create unit tests for tailoring service in `apps/backend/src/services/profile/tailor_test.go`
- [x] T042 [US4] Create `POST /api/profile/tailor` handler in `apps/backend/src/api/handlers/tailor.go`
- [x] T043 [P] [US4] Create Job Description input modal/form in `apps/frontend/components/profile/JobDescriptionInput.tsx`
- [x] T044 [US4] Update Profile View to support "Tailored Mode" (visual highlights) in `apps/frontend/components/profile/DashboardView.tsx`
- [x] T045 [US4] Update CV generation service to support tailored ordering/highlighting in `apps/backend/src/services/generator/cv_pdf.go`
- [x] T046 [P] [US4] Update unit tests for tailored CV generation in `apps/backend/src/services/generator/cv_pdf_test.go`

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T047 Verify all linting rules pass in backend and frontend
- [x] T048 Ensure error messages are user-friendly across UI
- [x] T049 [P] Add basic API documentation (Swagger UI or README)
- [x] T050 Verify docker-compose build works for full stack

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
- **Polish (Final Phase)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Independent.
- **User Story 2 (P1)**: Independent (can use dummy data) but integration requires US1 data.
- **User Story 3 (P2)**: Depends on US1/US2 data structures.
- **User Story 4 (P2)**: Depends on US1/US2 data structures.

### Parallel Opportunities

- Backend handlers and Frontend components for the same story can often be built in parallel once the API contract is defined.
- Independent components (e.g., PDF Extractor vs OpenAI Client) in Foundational phase.

## Implementation Strategy

### MVP First (User Story 1 & 2)

1. Complete Setup + Foundational
2. Implement Upload Flow (US1)
3. Implement Profile View (US2)
4. **Validate**: Can upload a letter and see the profile.

### Incremental Delivery

1. Add Download CV (US3)
2. Add Tailoring (US4)
