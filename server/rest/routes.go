package rest

func (r *REST) InitRoutes() {
	router := r.router
	router.GET("/metrics", r.middlewares.GinMetricsHandler())
	router.Use(r.middlewares.ErrorHandler())
	m := r.middlewares

	{
		v1 := router.Group("/v1")
		auth := v1.Group("/auth")
		{
			session := auth.Group("/session")
			{
				session.GET("/send-verify-code", r.session.SendVerifyCode)
				session.GET("/verify", r.session.Verify)
				session.PUT("/refresh", r.session.Refresh)
				session.GET("/logout", r.session.Logout)
			}

			web := auth.Group("/web")
			{
				web.POST("/sign-in", r.web.SignIn, m.SetToken())
				web.POST("/sign-up", r.web.SignUp, m.SetToken())

				gmail := web.Group("/gmail")
				{
					gmail.GET("/call-back", r.web.GoogleCallBack)
					gmail.GET("/:state/redirect", r.web.GoogleRedirect)
				}
			}

			mobile := auth.Group("/mobile")
			{
				mobile.POST("/sign-in")

				gmail := mobile.Group("/gmail")
				{
					gmail.POST("")
					gmail.PUT("")
				}
			}

			user := auth.Group("/user")
			{
				user.GET("/:user_idcode",
					m.VerifySession,
					r.user.GetByUserIDCode)
			}
		}
	}
}
