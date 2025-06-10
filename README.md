# gh-discussion

`gh-discussion` is a GitHub CLI extension to search and view GitHub Discussions.

## Installation
```bash
gh extension install harakeishi/gh-discussion
```

Users on GitHub Enterprise should specify the full repository path including
the host:
```bash
gh extension install github.com/harakeishi/gh-discussion
```

The extension builds automatically via the included `gh-discussion` script.

## Usage
- Search discussions:
```
gh discussion search --repo owner/repo --from 2023-01-01 --to 2023-01-31 --keyword bug
```
- View discussion:
```
gh discussion view <url-or-id>
```

## Development
This project is implemented in Go. Tests can be run with:
```
go test ./...
```

