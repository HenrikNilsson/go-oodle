module example.com/testapp

go 1.22.2

toolchain go1.24.3

replace github.com/new-world-tools/go-oodle => ../../

require github.com/new-world-tools/go-oodle v0.0.0-00010101000000-000000000000

require (
	github.com/ebitengine/purego v0.8.0-alpha.2.0.20240522163517-88cc57927e42 // indirect
	golang.org/x/sys v0.20.0 // indirect
)
