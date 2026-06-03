// Asgard web-admin is a Vue 3 + Vite + TypeScript SPA. It has no Go code,
// but some npm packages (e.g. flatted) ship a Go port under
// node_modules/<pkg>/golang/... that Go's ./... walker would otherwise pick
// up while building/testing the parent Asgard Go module.
//
// Declaring this directory a separate (empty) Go module creates a module
// boundary so `go build ./...` / `go test ./...` from the repo root stops
// at web-admin/. The Go file inside node_modules is unrelated to our build.
module web-admin

go 1.16
