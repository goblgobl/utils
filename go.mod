module src.goblgobl.com/utils

go 1.21

toolchain go1.21.1

// replace src.goblgobl.com/tests => ../tests

require (
	github.com/goccy/go-json v0.10.2
	github.com/google/uuid v1.3.1
	github.com/jackc/pgx/v5 v5.1.1
	github.com/valyala/fasthttp v1.43.0
	golang.org/x/crypto v0.3.0
	golang.org/x/sync v0.1.0
	src.goblgobl.com/sqlite v0.0.4
	src.goblgobl.com/tests v0.0.8
)

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/puddle/v2 v2.1.2 // indirect
	github.com/klauspost/compress v1.15.12 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	golang.org/x/sys v0.2.0 // indirect
	golang.org/x/text v0.4.0 // indirect
)
