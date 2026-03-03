package reporeviewer

import "embed"

//go:embed migrations/*.sql
var MigrationFiles embed.FS
