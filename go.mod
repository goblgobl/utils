module src.goblgobl.com/utils

go 1.21

toolchain go1.21.1

// replace src.goblgobl.com/tests => ../tests

require (
	github.com/goccy/go-json v0.10.2
	github.com/google/uuid v1.3.1
	github.com/jackc/pgx/v5 v5.4.3
	github.com/valyala/fasthttp v1.49.0
	golang.org/x/crypto v0.13.0
	golang.org/x/sync v0.1.0
	src.goblgobl.com/sqlite v0.0.4
	src.goblgobl.com/tests v0.0.8
)

require (
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/klauspost/compress v1.16.3 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/text v0.13.0 // indirect
)
