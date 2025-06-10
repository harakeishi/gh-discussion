# gh-discussion

`gh-discussion` is a GitHub CLI extension to search and view GitHub Discussions.

## Installation
```
gh extension install harakeishi/gh-discussion
```

## Usage
- Search discussions:
```
gh discussion search --from 2023-01-01 --to 2023-01-31 --keyword bug
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

