# Pull Request

Thank you for contributing to Atron. Before requesting a review, please
complete the checklist below and read CONTRIBUTING.md.

## Summary

<!-- Describe what this PR does in a few sentences. Link any related issues. -->

## Checklist

- [ ] I have read [CONTRIBUTING.md](../CONTRIBUTING.md).
- [ ] My commit messages follow [Conventional Commits](https://www.conventionalcommits.org/) and include a DCO `Signed-off-by:` trailer.
- [ ] The change is within the current milestone (see docs/PROGRESS.md) and was not approved beyond its scope.
- [ ] I added/updates tests covering the new behavior (table-driven and/or golden as appropriate).
- [ ] Every `.go` file I added/edited contains the SPDX header `// SPDX-License-Identifier: Apache-2.0`.
- [ ] Every exported identifier has a `godoc` comment.
- [ ] No external dependencies were added to the language core (internal/*).
- [ ] Local verification passed:
  - [ ] `gofmt -l .` is empty
  - [ ] `go vet ./...` passes
  - [ ] `go test -race ./...` passes
  - [ ] `make lint` passes (golangci-lint)

## Verification output

<!-- Paste the relevant output of the commands above. -->

## Related issues

<!-- e.g. Closes #12 -->
