# Repository Structure

## Purpose
To define the standard folder structure for the repository and individual modules, ensuring navigability and predictability.

## Observed Evidence
- **Root Level**: `cmd/`, `config/`, `module/`, `docs/`.
- **Module Level**: `module/<name>/` containing `delivery/`, `domain/`, `repository/`, `usecase/`.
- **Delivery Level**: `delivery/http/`, `delivery/messaging/`.

## Rules (MUST)
- **MUST** place all feature code under `module/<feature_name>/`.
- **MUST** structure each module with the following subdirectories:
  - `domain/`: Interfaces and domain entities.
  - `usecase/`: Business logic implementation.
  - `repository/`: Data access implementation (e.g., `postgresql/`).
  - `delivery/`: Transport layers (e.g., `http/`, `messaging/`).
- **MUST** place application entry points in `cmd/<command_name>.go`.
- **MUST** place configuration structs in `config/`.

## Conventions (SHOULD)
- **SHOULD** use lowercase, snake_case for directory and file names.
- **SHOULD** group related files by functionality (e.g., `create_user.go`, `get_user.go`) rather than giant single files.

## Open Questions
- None.
