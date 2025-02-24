PORT=15521

.PHONY: run-batched run-sqlite setup-sqlite

run-batched:
	cd batched-server && go run main.go $(PORT)

run-sqlite: setup-sqlite
	cd sqlite-server && go run main.go $(PORT)

setup-sqlite:
	@if [ ! -f sqlite-server/go.mod ]; then \
		cd sqlite-server && go mod init sqlite-server; \
	fi
	cd sqlite-server && go get modernc.org/sqlite