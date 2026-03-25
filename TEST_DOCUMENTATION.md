# Nitrous Testing Documentation

## Overview

This document provides comprehensive information about the testing infrastructure for the Nitrous platform, including both frontend and backend tests.

## Tech Stack

### Frontend Testing
- **Jest** - Unit test framework
- **React Testing Library** - Component testing utilities
- **Cypress** - End-to-end testing framework
- **React 18.3.1** - Latest React version
- **Next.js 14.2.18** - React framework

### Backend Testing
- **Go Testing** - Built-in Go testing framework
- **Gin Test Router** - HTTP handler testing

## Project Structure

```
nitrous-app/
├── __tests__/                           # Unit test files
│   ├── Nav.test.tsx                    # Navigation component tests
│   ├── Hero.test.tsx                   # Hero section component tests
│   └── api.test.ts                     # API utility function tests
├── cypress/                            # E2E test configuration
│   ├── cypress.config.ts               # Cypress configuration
│   └── e2e/                            # E2E test files
│       ├── home.cy.ts                  # Home page navigation tests
│       └── hero-interactions.cy.ts     # Hero interaction tests
├── jest.config.js                      # Jest configuration
├── jest.setup.js                       # Jest setup file
└── package.json                        # Project dependencies

nitrous-backend/
├── handlers/
│   ├── auth_handlers_test.go
│   ├── admin_management_test.go
│   ├── events_mutations_test.go
│   ├── journeys_teams_test.go
│   ├── orders_reminders_test.go
│   └── test_helpers_test.go
├── middleware/
│   ├── auth_test.go
│   └── admin_test.go
└── utils/
    └── jwt_test.go
```

## Frontend Testing

### Unit Tests

#### Setup
Jest and React Testing Library are configured in [jest.config.js](../../jest.config.js) and [jest.setup.js](../../jest.setup.js).

#### Running Unit Tests

```bash
# Run all tests
npm run test

# Run in watch mode (reruns on file changes)
npm run test:watch

# Generate coverage report
npm run test:coverage
```

#### Test Files and Coverage

##### 1. Nav Component Tests (`__tests__/Nav.test.tsx`)
Tests for the main navigation component.

**Tests:**
- ✅ Renders navigation with logo
- ✅ Renders all navigation links (Live, Events, Teams, Journeys, Merch)
- ✅ Renders Sign In button
- ✅ Displays live events count
- ✅ Renders navigation in correct structure
- ✅ Has correct links href attributes

##### 2. Hero Component Tests (`__tests__/Hero.test.tsx`)
Tests for the hero section component.

**Tests:**
- ✅ Renders hero section
- ✅ Displays main title with NITROUS and FUEL
- ✅ Displays subtitle text
- ✅ Renders action buttons
- ✅ Renders hero navigation cards
- ✅ Renders background image
- ✅ Contains circuit layer elements
- ✅ Renders HUD corners
- ✅ Displays HUD label text

##### 3. API Utility Function Tests (`__tests__/api.test.ts`)
Tests for API client functions in [lib/api.ts](../../lib/api.ts).

**Test Categories:**

**Events API**
- ✅ `getEvents()` - Fetches all events
- ✅ `getEventById(id)` - Fetches event by ID

**Categories API**
- ✅ `getCategories()` - Fetches all categories

**Journeys API**
- ✅ `getJourneys()` - Fetches all journeys

**Merch API**
- ✅ `getMerchItems()` - Fetches merch items

**Authentication API**
- ✅ `register(email, password, name)` - User registration
- ✅ `login(email, password)` - User login
- ✅ `getCurrentUser(token)` - Gets current user info

**Error Handling**
- ✅ Handles API errors properly
- ✅ Handles network errors

### E2E Tests (Cypress)

#### Setup
Cypress is configured in [cypress.config.ts](../../cypress.config.ts) with baseUrl set to `http://localhost:3000`.

#### Running E2E Tests

```bash
# Open Cypress Test Runner (interactive)
npm run cypress

# Run Cypress tests in headless mode
npm run cypress:run
```

#### Test Files

##### 1. Home Page Navigation (`cypress/e2e/home.cy.ts`)
Tests for home page navigation and visibility.

**Tests:**
- ✅ Loads home page successfully
- ✅ Displays navigation menu with all links
- ✅ Displays Sign In button
- ✅ Displays hero title and subtitle
- ✅ Displays action buttons
- ✅ Navigates to Events page
- ✅ Navigates to Teams page
- ✅ Navigates to Journeys page
- ✅ Navigates to Merch page
- ✅ Displays live events status

##### 2. Hero Section Interactions (`cypress/e2e/hero-interactions.cy.ts`)
Tests for hero section interactivity.

**Tests:**
- ✅ Verifies hero action buttons are clickable
- ✅ Verifies all navigation cards are visible
- ✅ Clicks Ignite Stream button
- ✅ Clicks Explore Events button
- ✅ Navigates through hero nav cards
- ✅ Navigates to live streams from hero
- ✅ Verifies hero section styling elements

### Test Metrics

- **Unit Test Suites:** 3
- **Unit Tests Total:** 26
- **E2E Test Suites:** 2
- **E2E Tests Total:** 15+
- **Components Tested:** 2 (Nav, Hero)
- **API Functions Tested:** 9

## Backend Testing

### Setup
Backend tests use Go's built-in `testing` package with Gin test router patterns.

### Running Backend Tests

```bash
# Run all backend tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run tests for a specific package
go test ./handlers -v

# Run tests with coverage
go test ./... -cover
```

### Test Packages

#### 1. Handlers Tests

**Auth Handlers** (`handlers/auth_handlers_test.go`)
- ✅ `TestRegisterFlow` - User registration flow
- ✅ `TestLoginFlow` - User login flow
- ✅ `TestGetCurrentUserFlow` - Get current user info

**Event Mutations** (`handlers/events_mutations_test.go`)
- ✅ `TestCreateEventEndpoint` - Event creation
- ✅ `TestUpdateEventEndpoint` - Event update
- ✅ `TestDeleteEventEndpoint` - Event deletion

**Admin Management** (`handlers/admin_management_test.go`)
- ✅ `TestCategoryManagementAdminRoutes` - Category CRUD operations
- ✅ `TestJourneyCatalogManagementAdminRoutes` - Journey CRUD operations
- ✅ `TestTeamManagementAdminRoutes` - Team CRUD operations
- ✅ `TestStreamManagementAdminRoutes` - Stream CRUD operations

**Journey and Team Routes** (`handlers/journeys_teams_test.go`)
- ✅ `TestBookJourneyEndpoint` - Book journey
- ✅ `TestFollowTeamEndpoint` - Follow team
- ✅ `TestUnfollowTeamEndpoint` - Unfollow team

**Orders and Reminders** (`handlers/orders_reminders_test.go`)
- ✅ `TestCreateOrderEndpoint` - Create order
- ✅ `TestGetMyOrdersEndpoint` - Get user's orders
- ✅ `TestGetOrderByIDEndpoint` - Get order by ID
- ✅ `TestSetReminderEndpoint` - Create reminder
- ✅ `TestGetMyRemindersEndpoint` - Get user's reminders
- ✅ `TestDeleteReminderEndpoint` - Delete reminder

**Test Helpers** (`handlers/test_helpers_test.go`)
- Utilities for setting up test environment

#### 2. Middleware Tests

**Auth Middleware** (`middleware/auth_test.go`)
- ✅ `TestAuthMiddleware` - JWT authentication validation

**Admin Middleware** (`middleware/admin_test.go`)
- ✅ `TestAdminMiddleware` - Admin role authorization

#### 3. Utils Tests

**JWT Utilities** (`utils/jwt_test.go`)
- ✅ `TestJWTUtility` - JWT token generation and validation

### Test Coverage

- **Total Test Packages:** 3 (handlers, middleware, utils)
- **Test Files:** 8 files
- **Test Functions:** 20+ test functions
- **Coverage:** Auth, Admin, Events, Categories, Journeys, Teams, Streams, Orders, Reminders, JWT

## Frontend-Backend Integration

### API Configuration

The frontend connects to the backend via environment variables in [lib/api.ts](../../lib/api.ts):

```typescript
const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api'
```

### Running Both Services Locally

1. **Start Backend:**
```bash
cd nitrous-backend
go run main.go
# Backend runs on http://localhost:8080
```

2. **Start Frontend:**
```bash
cd nitrous-app
npm run dev
# Frontend runs on http://localhost:3000
```

3. **Verify Integration:**
- Home page loads at `http://localhost:3000`
- API integration verified in "API Data Check" section on home page
- Shows real counts: Events, Categories, Journeys, Merch Items

### API Endpoints Tested

**Public Endpoints:**
- `GET /api/events` - List events
- `GET /api/categories` - List categories
- `GET /api/journeys` - List journeys
- `GET /api/merch` - List merch items
- `GET /api/teams` - List teams
- `GET /api/streams` - List streams

**Protected Endpoints (Require JWT):**
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `GET /api/auth/me` - Get current user
- `POST /api/journeys/:id/book` - Book journey
- `POST /api/teams/:id/follow` - Follow team
- `POST /api/teams/:id/unfollow` - Unfollow team
- `POST /api/orders` - Create order
- `GET /api/orders` - Get user's orders
- `POST /api/reminders` - Create reminder
- `GET /api/reminders` - Get user's reminders

**Admin Endpoints (Require Admin Role):**
- `POST /api/events` - Create event
- `PUT /api/events/:id` - Update event
- `DELETE /api/events/:id` - Delete event
- `POST /api/categories` - Create category
- `PUT /api/categories/:slug` - Update category
- `DELETE /api/categories/:slug` - Delete category
- `POST /api/journeys` - Create journey
- `PUT /api/journeys/:id` - Update journey
- `DELETE /api/journeys/:id` - Delete journey
- `POST /api/teams` - Create team
- `PUT /api/teams/:id` - Update team
- `DELETE /api/teams/:id` - Delete team
- `POST /api/streams` - Create stream
- `PUT /api/streams/:id` - Update stream
- `DELETE /api/streams/:id` - Delete stream

## Continuous Integration

### Test Commands Summary

**Frontend:**
```bash
npm run test              # Run unit tests
npm run test:watch       # Watch mode
npm run test:coverage    # Coverage report
npm run cypress          # Interactive Cypress
npm run cypress:run      # Headless Cypress
```

**Backend:**
```bash
go test ./...            # Run all tests
go test ./... -v         # Verbose output
go test ./... -cover     # With coverage
```

## Performance

### Test Execution Time
- **Frontend Unit Tests:** ~2-3 seconds
- **Frontend E2E Tests:** ~5-10 seconds (per test suite)
- **Backend Tests:** ~3-4 seconds

### Coverage Goals
- **Frontend:** Aim for 1:1 unit test to function ratio
- **Backend:** Comprehensive coverage of handlers, middleware, and utilities

## Best Practices

### Unit Testing
1. Mock external dependencies (Next.js Link, Image)
2. Test component rendering and user interactions
3. Use semantic queries (getByText, getByRole) over implementation details
4. Test error states and edge cases

### E2E Testing
1. Test complete user journeys
2. Verify API integration without mocking
3. Test navigation and user interactions
4. Keep tests stable and maintainable

### API Testing
1. Mock fetch function for unit tests
2. Test both success and error scenarios
3. Verify token handling for protected routes
4. Test with realistic data structures

## Troubleshooting

### Frontend Tests
- **Warn: Received `true` for non-boolean attribute** - Next.js Image component passes boolean props; use container DOM queries instead
- **Unable to find element** - Text may be split across elements; use flexible matchers or container queries

### Backend Tests
- **Test fails with connection error** - Ensure using test router and in-memory database

### E2E Tests
- **Cypress times out** - Increase timeout in cypress.config.ts
- **Navigation fails** - Verify frontend is running on http://localhost:3000

## Next Steps

- Increase unit test coverage to 80%+
- Add performance testing for E2E flows
- Implement CI/CD pipeline with automated testing
- Add visual regression testing with Cypress
- Expand test scenarios for edge cases and error handling

---

Last Updated: March 25, 2026
