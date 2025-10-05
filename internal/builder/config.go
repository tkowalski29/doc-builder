package builder

import "context"

type Config struct {
	Prefix      string
	Engine      string
	SearchPath  string
	DocDir      string
	TempDirName string
	Verbose     bool
}

type Builder struct {
	cfg Config
}

func New(cfg Config) *Builder {
	return &Builder{cfg: cfg}
}

func (b *Builder) Run(ctx context.Context) error {
	if err := b.validateConfig(); err != nil {
		return err
	}

	env, err := b.prepareEnvironment()
	if err != nil {
		return err
	}

	if err := b.prepareWorkspace(env); err != nil {
		return err
	}

	menuRecords := make([]menuRecord, 0, 128)
	recordSet := make(map[string]struct{})

	prefixedRecords, prefCount, err := b.collectPrefixedDocs(env, recordSet)
	if err != nil {
		return err
	}
	menuRecords = append(menuRecords, prefixedRecords...)

	existingRecords, existingCount, err := b.collectExistingDocs(env, recordSet)
	if err != nil {
		return err
	}
	menuRecords = append(menuRecords, existingRecords...)

	if len(menuRecords) == 0 {
		return errNoSources
	}

	if err := b.writeMenuIndex(env, menuRecords); err != nil {
		return err
	}

	if err := b.generateConfig(env, menuRecords); err != nil {
		return err
	}

	if err := b.installDependencies(ctx, env); err != nil {
		return err
	}

	if err := b.buildSite(ctx, env); err != nil {
		return err
	}

	if err := b.publishDist(env); err != nil {
		return err
	}

	b.printSummary(env, prefCount, existingCount, len(menuRecords))
	return nil
}
