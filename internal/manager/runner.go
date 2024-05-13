package manager

import (
	"context"
	"go.uber.org/fx"
)

func Run(
	lc fx.Lifecycle,
	manager *Manager,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {

			processor := manager.Processor()
			{
				processor.Queue().Producers().Mails()
				processor.Cache()
			}
			repository := manager.Repository()
			{
				repository.Users()
				repository.Sessions()
				repository.GoogleAPICodes()
				repository.SessionsCache()
				repository.VerificationCodes()
			}
			service := manager.Service()
			{
				service.User()
				service.Auth()
				service.OAuth()
			}

			manager.Server().REST().Run()
			return
		},
		OnStop: func(ctx context.Context) (err error) {
			if err = manager.db.Close(); err != nil {
				return err
			}
			if err = manager.cacheProvider.Close(); err != nil {
				return err
			}
			return nil
		},
	})
}
