```markdown
# ai-api-gateway Development Patterns

> Auto-generated skill from repository analysis

## Overview
This skill teaches the core development patterns and conventions used in the `ai-api-gateway` TypeScript codebase. You'll learn how to structure files, write imports/exports, follow commit message conventions, and understand the testing approach. These patterns ensure consistency and maintainability across the project.

## Coding Conventions

### File Naming
- Use **camelCase** for file names.
  - Example: `apiGateway.ts`, `userService.ts`

### Import Style
- Use **alias imports** to reference modules.
  - Example:
    ```typescript
    import apiService from '@/services/apiService';
    ```

### Export Style
- Use **default exports** for modules.
  - Example:
    ```typescript
    const apiGateway = { ... };
    export default apiGateway;
    ```

### Commit Messages
- Follow the **Conventional Commits** specification.
- Use the `feat` prefix for new features.
  - Example:  
    ```
    feat: add support for multiple AI providers in gateway
    ```

## Workflows

### Feature Development
**Trigger:** When adding a new feature to the gateway  
**Command:** `/feature-dev`

1. Create a new file using camelCase naming.
2. Implement the feature, using alias imports as needed.
3. Export your module as default.
4. Write or update corresponding test files (`*.test.*`).
5. Commit your changes using the `feat:` prefix and a concise description.

### Testing
**Trigger:** When validating code changes  
**Command:** `/run-tests`

1. Locate or create test files matching the `*.test.*` pattern.
2. Run the test suite using the project's preferred test runner (framework not specified).
3. Ensure all tests pass before merging changes.

## Testing Patterns

- Test files follow the `*.test.*` naming convention (e.g., `apiGateway.test.ts`).
- The specific testing framework is not detected; check the project documentation or package.json for details.
- Place tests alongside implementation files or in a dedicated test directory.

## Commands
| Command        | Purpose                                         |
|----------------|-------------------------------------------------|
| /feature-dev   | Start the feature development workflow           |
| /run-tests     | Run all tests in the codebase                   |
```