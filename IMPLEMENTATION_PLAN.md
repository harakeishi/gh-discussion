# Implementation Plan for GitHub Discussion Extension

This document outlines how to extend the current `gh-discussion` project to support
searching discussions within a repository and fetching a discussion's content.

## Features

1. **Search discussions** in a specified repository with filters for:
   - Creation time range (start and/or end date)
   - Author username
   - Keyword query
   - *At least one of the above filters must be provided.*
   - The repository (`OWNER/REPO`) is required.
2. **View discussion**: retrieve and display the full content of a discussion by
   its number.

## CLI Design

The extension will provide two subcommands:

```bash
$ gh discussion search [flags]
$ gh discussion view <number> --repo OWNER/REPO
```

### `search` flags
- `--repo <OWNER/REPO>` (required)
- `--start <YYYY-MM-DD>` (optional)
- `--end <YYYY-MM-DD>` (optional)
- `--user <username>` (optional)
- `--keyword <words>` (optional)
- Pagination will be supported automatically using GraphQL cursors.

### `view` arguments and flags
- `<number>`: discussion number to fetch.
- `--repo <OWNER/REPO>` (required)

## Implementation Steps

1. **Set up command framework**
   - Introduce `spf13/cobra` for subcommand handling or use `gh`'s helper
     packages from `github.com/cli/go-gh/v2`. Create root command
     `discussion` with `search` and `view` subcommands.

2. **Search Discussions**
   - Construct a search query string using provided flags.
     - Example: `repo:OWNER/REPO user:alice created:2024-01-01..2024-02-01 bug`
   - Use the GraphQL `search` API with `type: DISCUSSION` to retrieve
     discussion metadata (number, title, author, createdAt, url).
   - Validate that at least one of `--start`, `--end`, `--user`, or `--keyword`
     is set. Return an error otherwise.
   - Paginate results until the user-defined limit or until no more pages.
   - Output a table or list of discussions to stdout.

3. **View Discussion**
   - Query `repository.discussion(number: X)` via GraphQL to obtain
     `title`, `author`, `createdAt`, and `body`/`bodyHTML`.
   - Display the content in the terminal using plain text or render HTML via
     `gh`'s formatter if available.

4. **Error Handling & Authentication**
   - Leverage `go-gh` to obtain an authenticated GraphQL client.
   - Surface API errors in a user-friendly manner.

5. **Testing**
   - Add unit tests for query construction and flag validation.
   - Consider integration tests using GitHub's GraphQL API if tokens are
     available, otherwise mock responses.

6. **Documentation**
   - Update `README.md` with usage examples for the new commands.
   - Document required permissions (the extension relies on the user's
     GitHub CLI authentication).

## Future Enhancements
- Support filtering by labels or categories.
- Allow output in JSON for scripting purposes.
- Add more formatting options when viewing discussions.

