package rest

func (r *REST) InitRoutes() {
	router := r.router
	router.GET("/metrics", r.middlewares.GinMetricsHandler())
	router.Use(r.middlewares.ErrorHandler())

	{
		v1 := router.Group("/v1")
		auth := v1.Group("/auth")
		{

			auth.POST("/sign-in", r.auth.SignIn)
			auth.POST("/sign-up", r.auth.SignUp)
			auth.PUT("/refresh", r.auth.Refresh)
			auth.GET("/logout", r.auth.Logout)
			auth.GET("/verify", r.auth.VerifySession)

			google := auth.Group("/google")
			{
				google.GET("/call-back", r.oAuth.GoogleCallBack)
				google.GET("/:state/redirect", r.oAuth.GoogleRedirect)
			}

			user := auth.Group("/user")
			{
				user.GET("/:user_idcode",
					//r.middlewares.VerifySession,
					r.user.GetByUserIDCode)
				user.GET("/send-verify-code", r.user.SendVerifyCode)
			}
		}
	}
}
