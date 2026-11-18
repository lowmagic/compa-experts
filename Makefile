.PHONY: check standards sync-standards
SHELL := bash
ORG = https://github.com/lowmagic/.github

check:
	golangci-lint run
	go test -race ./...
	go mod tidy -diff
	git ls-files -z '*.sh' | xargs -0 -r shellcheck
	@# A changed package needs a new version, against the merge base with
	@# origin/main; a clone without origin/main fetches it first.
	git rev-parse -q --verify origin/main >/dev/null || git fetch -q origin main:refs/remotes/origin/main
	go run ./cmd/version-check

standards:
	@tmp=$$(mktemp -d) && trap 'rm -rf "$$tmp"' EXIT && \
	git clone -q --depth 1 --branch main $(ORG) "$$tmp" && \
	"$$tmp/code-standards/sync.sh" check "$$tmp"

sync-standards:
	@tmp=$$(mktemp -d) && trap 'rm -rf "$$tmp"' EXIT && \
	git clone -q --depth 1 --branch main $(ORG) "$$tmp" && \
	"$$tmp/code-standards/sync.sh" sync "$$tmp" "$$(git -C "$$tmp" rev-parse HEAD)"
