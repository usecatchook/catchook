# Contributing to Catchook

Thanks for your interest in contributing to Catchook! 🎉

## Quick Start

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/catchook.git`
3. Create a branch: `git checkout -b your-feature-branch`
4. Make your changes
5. Test your changes: `make test`
6. Submit a pull request

## Development Setup

```bash
# Install dependencies
go mod tidy
cd app && npm install && cd ..

# Start services
docker-compose -f docker-compose.dev.yml up -d

# Run API
make dev-api

# Run Frontend
make dev-app
```

## Code Style

- **Go**: Follow `gofmt` and `golint` standards
- **TypeScript**: Use ESLint configuration
- **Commits**: Use conventional commits (feat:, fix:, docs:)

## Testing

```bash
# Run tests
make test

# Run linting
make lint
```

## Submitting Changes

- Write clear, descriptive commit messages
- Add tests for new features
- Update documentation if needed
- Make sure all tests pass

## Questions?

- Open an issue for bugs or feature requests
- Join discussions in existing issues
- Check the README for architecture details

Thanks for contributing! 🚀