module src.goblgobl.com/utils

go 1.21

toolchain go1.21.1

// replace src.goblgobl.com/tests => ../tests

require (
	github.com/goccy/go-json v0.10.2
	github.com/google/uuid v1.5.0
	github.com/jackc/pgx/v5 v5.5.1
	github.com/valyala/fasthttp v1.51.0
	golang.org/x/crypto v0.17.0
	golang.org/x/sync v0.5.0
	src.goblgobl.com/sqlite v0.0.4
	src.goblgobl.com/tests v0.1.0
)

require (
	github.com/andybalholm/brotli v1.0.6 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/klauspost/compress v1.17.4 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
