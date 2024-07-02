package main

import (
	"flag"
)

const migrationsPath = "/opt/canopsis/share/database/migrations"
const (
	MigrationExecGoja    = "goja"
	MigrationExecMongosh = "mongosh"
)

type createFlags struct {
	name string
	path string
}

func (f *createFlags) Parse(arguments []string) error {
	flags := flag.NewFlagSet("create: Create a blank migration", flag.ContinueOnError)
	flags.StringVar(&f.path, "path", migrationsPath, "Migration directory")
	flags.StringVar(&f.name, "name", "", "Migration name")
	return flags.Parse(arguments)
}

type upFlags struct {
	to            string
	path          string
	migrationExec string
}

func (f *upFlags) Parse(arguments []string) error {
	flags := flag.NewFlagSet("up: Execute migrations to a specified version or the latest available version", flag.ContinueOnError)
	flags.StringVar(&f.path, "path", migrationsPath, "Migration directory")
	flags.StringVar(&f.to, "to", "", "Migrate to version")
	flag.StringVar(&f.migrationExec, "migration-exec", MigrationExecGoja, "The execution of Mongo migration scripts: "+MigrationExecGoja+" or "+MigrationExecMongosh)

	return flags.Parse(arguments)
}

type downFlags struct {
	to            string
	path          string
	migrationExec string
}

func (f *downFlags) Parse(arguments []string) error {
	flags := flag.NewFlagSet("down: Roll migrations up to a specified version or all tracked versions", flag.ContinueOnError)
	flags.StringVar(&f.path, "path", migrationsPath, "Migration directory")
	flags.StringVar(&f.to, "to", "", "Revert migrations to version")
	flag.StringVar(&f.migrationExec, "migration-exec", MigrationExecGoja, "The execution of Mongo migration scripts: "+MigrationExecGoja+" or "+MigrationExecMongosh)

	return flags.Parse(arguments)
}

type statusFlags struct {
	path string
}

func (f *statusFlags) Parse(arguments []string) error {
	flags := flag.NewFlagSet("status: View the status of migrations", flag.ContinueOnError)
	flags.StringVar(&f.path, "path", migrationsPath, "Migration directory")
	return flags.Parse(arguments)
}

type skipFlags struct {
	path    string
	version string
}

func (f *skipFlags) Parse(arguments []string) error {
	flags := flag.NewFlagSet("skip: Manually add migrations until specified version or all untracked versions to the version table", flag.ContinueOnError)
	flags.StringVar(&f.path, "path", migrationsPath, "Migration directory")
	flags.StringVar(&f.version, "version", "", "Migration version")
	return flags.Parse(arguments)
}
