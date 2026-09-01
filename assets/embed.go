package assets

import "embed"

//go:embed "html"
var HTML embed.FS

//go:embed "emails"
var Emails embed.FS

//go:embed "static"
var Static embed.FS

//go:embed "migrations"
var Migrations embed.FS
