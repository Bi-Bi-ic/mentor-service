package app

import (
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/rgrs-x/service/api/controllers"
	courseControl "github.com/rgrs-x/service/api/controllers/course"
	postsControl "github.com/rgrs-x/service/api/controllers/posts"
	userControl "github.com/rgrs-x/service/api/controllers/user"
)

// SetupRoutes ...
func SetupRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(location.Default())

	apiV1 := router.Group("api")
	{
		apiV1.GET("/user/avatar/:name", controllers.Render)
		apiV1.GET("/partner/avatar/:name", controllers.Render)
		apiV1.GET("/file/:name", controllers.Render)
		//@ api for user version 1.0.0
		apiV1.Use(APIAuthentication())

		// Public Infomations a Client
		apiV1.GET("/user/info/:id", controllers.PublicUserInfo)
		apiV1.GET("/partner/info/:id", controllers.PublicPartnerInfo)

		auth := apiV1.Group("/auth")

		// generation
		{
			auth.POST("/user/sign_up", controllers.CreateUserAccount)
			auth.POST("/user/sign_in", controllers.AuthenticateUser)

			// for only refresh token
			auth.POST("/backend/get-access-token/user", controllers.UserToken)
		}

		// Use for all
		tracking := apiV1.Group("/tracking")
		{
			tracking.POST("/read-post", controllers.ReadPost)
		}

		//Comments
		comment := apiV1.Group("auth/content")
		{
			comment.POST("/:id/comments", GeneralAuthentication(), controllers.CreateComment)
			comment.GET("/course/:id/comments", controllers.GetCourseComments)
			comment.PATCH("/comments/:id", GeneralAuthentication(), controllers.UpdateComment)
			comment.DELETE("/comments/:id", GeneralAuthentication(), controllers.DeleteComment)
		}

		contents := apiV1.Group("/contents")
		{
			contents.GET("/", controllers.Pagination)
			contents.GET("/feature", postsControl.PostFeatureListController)
			contents.GET("/filter", controllers.Filter)
			contents.PUT("/:id/like", controllers.LikePost)
			contents.GET("/post/:id", controllers.GetPost)

			contents.GET("/company/:id", controllers.GetCompanyContents)
		}

		locationsService := apiV1.Group("/location")
		{
			locationsService.GET("/:id", controllers.FindLocation)
		}

		company := apiV1.Group("/company")
		{
			company.POST("/", controllers.CreateCompany)
			company.GET("/", controllers.SwitchGetCompany)
		}

		// User and Partner
		tags := apiV1.Group("/post")
		{
			tags.GET("/tags", controllers.GetAllTags)
		}

		mentor := apiV1.Group("/mentor")
		{
			mentor.PUT("/:id/like", controllers.LikeMentor)
		}

		// Only user
		user := apiV1.Group("/user")
		{
			user.Use(UserAuthentication())

			controllers.UploadPool = make(chan controllers.WorkerMessage, 10)
			go controllers.InitWorker(controllers.UploadPool)

			user.PUT("", controllers.UpdateUserInfo)
			user.GET("/", controllers.GetAuthUserInfo)
			user.POST("/avatar", controllers.UpdateAvatarUser)
			user.POST("/cover", controllers.UpdateUserCover)
			user.POST("/time-line", controllers.CreateUserTimeLine)
			user.PUT("/time-line/:id", controllers.UpdateUserTimeLine)
			user.DELETE("/time-line/:id", controllers.DeleteUserTimeLine)

			user.POST("/course", controllers.CreateCourse)
			user.POST("/course/:id/register", controllers.RegisterCourse)
			user.POST("/course/:id/like", controllers.LikeCourse)
			user.PUT("/course/:id", controllers.UpdateCourse)
			user.GET("/courses/:id", controllers.GetAllCourseByMentorId)
			user.DELETE("/course/:id", controllers.DeleteCourseById)

		}

		courseEntitties := apiV1.Group("/course")
		{
			courseEntitties.GET("/:id", controllers.GetCourseById)
			courseEntitties.GET("/:id/mentees", controllers.GetMenteesFromCourse)
		}

		course := apiV1.Group("/courses-all")
		{
			course.GET("", courseControl.CourseListController)
			course.GET("/feature", courseControl.CourseFeatureListController)
		}

		userList := apiV1.Group("/users")
		{
			userList.GET("", userControl.UserListController)
			userList.GET("/feature", userControl.UserFeatureListController)
		}

		posts := apiV1.Group("/posts")
		{
			posts.GET("", postsControl.ListPosts)
		}

		// contents_recommand ...
		contents_recommand := apiV1.Group("/recommand")
		{
			contents_recommand.Use(UserAuthentication())
			contents_recommand.GET("/", controllers.Pagination)
		}

		// Only admin
		admin := apiV1.Group("/admin")
		{
			admin.POST("/sign_in", controllers.AdminSignIn)
		}

		search := apiV1.Group("/search")
		{
			search.GET("/user", controllers.SearchUser)
			search.GET("/post", controllers.SearchPost)
			search.GET("/course", controllers.SearchCourse)
			search.GET("/all", controllers.SearchAll)
		}

	}

	// render json document for api
	//router.NotFoundHandler = app.NotFoundHandler
	return router
}
