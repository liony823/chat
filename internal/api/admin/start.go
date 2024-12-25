package admin

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/liony823/tools/errs"
	"github.com/liony823/tools/mw"
	"github.com/liony823/tools/utils/datautil"
	_ "github.com/openimsdk/chat/cmd/api/admin-api/docs" // 导入swagger docs
	chatmw "github.com/openimsdk/chat/internal/api/mw"
	"github.com/openimsdk/chat/internal/api/util"
	"github.com/openimsdk/chat/pkg/common/config"
	"github.com/openimsdk/chat/pkg/common/imapi"
	"github.com/openimsdk/chat/pkg/common/kdisc"
	adminclient "github.com/openimsdk/chat/pkg/protocol/admin"
	chatclient "github.com/openimsdk/chat/pkg/protocol/chat"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	ApiConfig config.API
	Discovery config.Discovery
	Share     config.Share
}

func Start(ctx context.Context, index int, config *Config) error {
	if len(config.Share.ChatAdmin) == 0 {
		return errs.New("share chat admin not configured")
	}
	apiPort, err := datautil.GetElemByIndex(config.ApiConfig.Api.Ports, index)
	if err != nil {
		return err
	}
	client, err := kdisc.NewDiscoveryRegister(&config.Discovery)
	if err != nil {
		return err
	}

	chatConn, err := client.GetConn(ctx, config.Share.RpcRegisterName.Chat, grpc.WithTransportCredentials(insecure.NewCredentials()), mw.GrpcClient())
	if err != nil {
		return err
	}
	adminConn, err := client.GetConn(ctx, config.Share.RpcRegisterName.Admin, grpc.WithTransportCredentials(insecure.NewCredentials()), mw.GrpcClient())
	if err != nil {
		return err
	}
	chatClient := chatclient.NewChatClient(chatConn)
	adminClient := adminclient.NewAdminClient(adminConn)
	im := imapi.New(config.Share.OpenIM.ApiURL, config.Share.OpenIM.Secret, config.Share.OpenIM.AdminUserID)
	base := util.Api{
		ImUserID:        config.Share.OpenIM.AdminUserID,
		ProxyHeader:     config.Share.ProxyHeader,
		ChatAdminUserID: config.Share.ChatAdmin[0],
		BasicAuthUser:   config.Share.BasicAuth.Username,
		BasicAuthPass:   config.Share.BasicAuth.Password,
	}
	adminApi := New(chatClient, adminClient, im, &base)
	mwApi := chatmw.New(adminClient)
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	engine.Use(gin.Recovery(), mw.CorsHandler(), mw.GinParseOperationID(), mw.GinAdminBasicAuth(config.Share.BasicAuth.Username, config.Share.BasicAuth.Password, config.Share.BasicAuth.Secret), mwApi.OperationLog())
	SetAdminRoute(engine, adminApi, mwApi)
	return engine.Run(fmt.Sprintf(":%d", apiPort))
}

func SetAdminRoute(router gin.IRouter, admin *Api, mw *chatmw.MW) {
	router.GET("/ping", admin.AdminPing)
	adminRouterGroup := router.Group("/account")
	adminRouterGroup.POST("/login", admin.AdminLogin)
	adminRouterGroup.POST("/google_auth", mw.CheckAdmin, admin.GetGoogleAuth)
	adminRouterGroup.POST("/verify_google_auth", mw.CheckAdmin, admin.VerifyGoogleAuth)
	adminRouterGroup.POST("/update", mw.CheckAdmin, admin.AdminUpdateInfo)
	adminRouterGroup.POST("/info", mw.CheckAdmin, admin.AdminInfo)
	adminRouterGroup.POST("/change_password", mw.CheckAdmin, admin.ChangeAdminPassword)
	adminRouterGroup.POST("/add_admin", mw.CheckAdmin, admin.AddAdminAccount)
	adminRouterGroup.POST("/add_user", mw.CheckAdmin, admin.AddUserAccount)
	adminRouterGroup.POST("/del_admin", mw.CheckAdmin, admin.DelAdminAccount)
	adminRouterGroup.POST("/search_admin", mw.CheckAdmin, admin.SearchAdminAccount)

	operationLogGroup := router.Group("/operation_log", mw.CheckAdmin)
	operationLogGroup.POST("/search", admin.SearchOperationLog)

	importGroup := router.Group("/user/import")
	importGroup.POST("/json", mw.CheckAdmin, admin.ImportUserByJson)
	importGroup.POST("/xlsx", mw.CheckAdmin, admin.ImportUserByXlsx)
	importGroup.GET("/xlsx", admin.BatchImportTemplate)

	allowRegisterGroup := router.Group("/user/allow_register", mw.CheckAdmin)
	allowRegisterGroup.POST("/get", admin.GetAllowRegister)
	allowRegisterGroup.POST("/set", admin.SetAllowRegister)

	defaultRouter := router.Group("/default", mw.CheckAdmin)
	defaultUserRouter := defaultRouter.Group("/user")
	defaultUserRouter.POST("/add", admin.AddDefaultFriend)
	defaultUserRouter.POST("/del", admin.DelDefaultFriend)
	defaultUserRouter.POST("/find", admin.FindDefaultFriend)
	defaultUserRouter.POST("/search", admin.SearchDefaultFriend)
	defaultGroupRouter := defaultRouter.Group("/group")
	defaultGroupRouter.POST("/add", admin.AddDefaultGroup)
	defaultGroupRouter.POST("/del", admin.DelDefaultGroup)
	defaultGroupRouter.POST("/find", admin.FindDefaultGroup)
	defaultGroupRouter.POST("/search", admin.SearchDefaultGroup)

	invitationCodeRouter := router.Group("/invitation_code", mw.CheckAdmin)
	invitationCodeRouter.POST("/add", admin.AddInvitationCode)
	invitationCodeRouter.POST("/gen", admin.GenInvitationCode)
	invitationCodeRouter.POST("/del", admin.DelInvitationCode)
	invitationCodeRouter.POST("/search", admin.SearchInvitationCode)

	forbiddenRouter := router.Group("/forbidden", mw.CheckAdmin)
	ipForbiddenRouter := forbiddenRouter.Group("/ip")
	ipForbiddenRouter.POST("/add", admin.AddIPForbidden)
	ipForbiddenRouter.POST("/del", admin.DelIPForbidden)
	ipForbiddenRouter.POST("/search", admin.SearchIPForbidden)
	userForbiddenRouter := forbiddenRouter.Group("/user")
	userForbiddenRouter.POST("/add", admin.AddUserIPLimitLogin)
	userForbiddenRouter.POST("/del", admin.DelUserIPLimitLogin)
	userForbiddenRouter.POST("/search", admin.SearchUserIPLimitLogin)

	appletRouterGroup := router.Group("/applet", mw.CheckAdmin)
	appletRouterGroup.POST("/add", admin.AddApplet)
	appletRouterGroup.POST("/del", admin.DelApplet)
	appletRouterGroup.POST("/update", admin.UpdateApplet)
	appletRouterGroup.POST("/search", admin.SearchApplet)

	blockRouter := router.Group("/block", mw.CheckAdmin)
	blockRouter.POST("/add", admin.BlockUser)
	blockRouter.POST("/del", admin.UnblockUser)
	blockRouter.POST("/search", admin.SearchBlockUser)

	userRouter := router.Group("/user", mw.CheckAdmin)
	userRouter.POST("/password/reset", admin.ResetUserPassword)

	initGroup := router.Group("/client_config", mw.CheckAdmin)
	initGroup.POST("/get", admin.GetClientConfig)
	initGroup.POST("/set", admin.SetClientConfig)
	initGroup.POST("/del", admin.DelClientConfig)
	initGroup.POST("/list", admin.GetListClientConfig)

	smsConfigGroup := router.Group("/sms_config", mw.CheckAdmin)
	smsConfigGroup.POST("/set", admin.SetSmsConfig)
	smsConfigGroup.POST("/get", admin.GetSmsConfig)

	bucketConfigGroup := router.Group("/bucket_config", mw.CheckAdmin)
	bucketConfigGroup.POST("/set", admin.SetBucketConfig)
	bucketConfigGroup.POST("/get", admin.GetBucketConfig)

	signinConfigGroup := router.Group("/signin_config", mw.CheckAdmin)
	signinConfigGroup.POST("/set", admin.SetSigninConfig)
	signinConfigGroup.POST("/get", admin.GetSigninConfig)

	adminMenuGroup := router.Group("/admin_menu", mw.CheckAdmin)
	adminMenuGroup.POST("/create", admin.CreateAdminMenu)
	adminMenuGroup.POST("/update", admin.UpdateAdminMenu)
	adminMenuGroup.POST("/delete", admin.DeleteAdminMenu)
	adminMenuGroup.POST("/take", admin.TakeAdminMenu)
	adminMenuGroup.POST("/list", admin.ListAdminMenu)
	adminMenuGroup.POST("/user_menu", admin.ListAdminUserMenu)
	adminMenuGroup.POST("/user_menu/assign", admin.AssignAdminUserMenu)
	adminMenuGroup.POST("/user_menu/get", admin.GetAdminUserMenu)

	statistic := router.Group("/statistic", mw.CheckAdmin)
	statistic.POST("/login_record", admin.LoginRecord)
	statistic.POST("/new_user_count", admin.NewUserCount)
	statistic.POST("/login_user_count", admin.LoginUserCount)

	applicationGroup := router.Group("application")
	applicationGroup.POST("/add_version", mw.CheckAdmin, admin.AddApplicationVersion)
	applicationGroup.POST("/update_version", mw.CheckAdmin, admin.UpdateApplicationVersion)
	applicationGroup.POST("/delete_version", mw.CheckAdmin, admin.DeleteApplicationVersion)
	applicationGroup.POST("/latest_version", admin.LatestApplicationVersion)
	applicationGroup.POST("/page_versions", admin.PageApplicationVersion)
}
