package cmd

import (
	"os"
	"path"

	"github.com/simplesurance/sisubuild/cfg"
	"github.com/simplesurance/sisubuild/discover"
	"github.com/simplesurance/sisubuild/sblog"
)

// Ctx stores supporting informations that are required by commands
type ctx struct {
	RepositoryRoot    string
	RepositoryCfg     *cfg.Repository
	RepositoryCfgPath string
}

// InitCtx returns an initialized Ctx. It terminates the application on errors.
func mustInitCtx() *ctx {
	var err error
	var ctx ctx

	ctx.RepositoryRoot, err = discover.RepositoryRoot(cfg.RepositoryFile)
	if err != nil {
		if os.IsNotExist(err) {
			sblog.Fatalf("could not find repository root config file "+
				"ensure the file '%s' exist in the root",
				cfg.RepositoryFile)
		}
		sblog.Fatal("finding repository root config file failed: ", err)
	}

	sblog.Debugf("repository root found: %v", ctx.RepositoryRoot)

	ctx.RepositoryCfgPath = path.Join(ctx.RepositoryRoot, cfg.RepositoryFile)

	ctx.RepositoryCfg, err = cfg.RepositoryFromFile(ctx.RepositoryCfgPath)
	if err != nil {
		sblog.Fatal("reading repository config failed: ", err)
	}

	if err = ctx.RepositoryCfg.Validate(); err != nil {
		sblog.Fatalf("validating repository config (%s) failed: %s",
			ctx.RepositoryCfgPath, err)
	}

	return &ctx
}
