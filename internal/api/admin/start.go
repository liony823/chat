package admin

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	disetcd "github.com/openimsdk/chat/pkg/common/kdisc/etcd"
	adminclient "github.com/openimsdk/chat/pkg/protocol/admin"
	chatclient "github.com/openimsdk/chat/pkg/protocol/chat"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/discovery/etcd"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/mw"
	"github.com/openimsdk/tools/system/program"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/runtimeenv"
	clientv3 "go.etcd.io/etcd/client/v3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	*config.AllConfig

	RuntimeEnv string
	ConfigPath string
}

func Start(ctx context.Context, index int, config *Config) error {
	config.RuntimeEnv = runtimeenv.PrintRuntimeEnvironment()

	if len(config.Share.ChatAdmin) == 0 {
		return errs.New("share chat admin not configured")
	}
	apiPort, err := datautil.GetElemByIndex(config.AdminAPI.Api.Ports, index)
	if err != nil {
		return err
	}
	client, err := kdisc.NewDiscoveryRegister(&config.Discovery, config.RuntimeEnv)
	if err != nil {
		return err
	}

	chatConn, err := client.GetConn(ctx, config.Discovery.RpcService.Chat, grpc.WithTransportCredentials(insecure.NewCredentials()), mw.GrpcClient())
	if err != nil {
		return err
	}
	adminConn, err := client.GetConn(ctx, config.Discovery.RpcService.Admin, grpc.WithTransportCredentials(insecure.NewCredentials()), mw.GrpcClient())
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
	engine.Use(gin.Recovery(), mw.CorsHandler(), mw.GinParseOperationID(), mw.GinAdminBasicAuth(config.Share.BasicAuth.Username, config.Share.BasicAuth.Password, config.Share.BasicAuth.Secret))
	SetAdminRoute(engine, adminApi, mwApi, config, client)

	if config.Discovery.Enable == kdisc.ETCDCONST {
		cm := disetcd.NewConfigManager(client.(*etcd.SvcDiscoveryRegistryImpl).GetClient(), config.GetConfigNames())
		cm.Watch(ctx)
	}
	var (
		netDone = make(chan struct{}, 1)
		netErr  error
	)
	server := http.Server{Addr: fmt.Sprintf(":%d", apiPort), Handler: engine}
	go func() {
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			netErr = errs.WrapMsg(err, fmt.Sprintf("api start err: %s", server.Addr))
			netDone <- struct{}{}
		}
	}()
	shutdown := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			return errs.WrapMsg(err, "shutdown err")
		}
		return nil
	}
	disetcd.RegisterShutDown(shutdown)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	select {
	case <-sigs:
		program.SIGTERMExit()
		if err := shutdown(); err != nil {
			return err
		}
	case <-netDone:
		close(netDone)
		return netErr
	}
	return nil
}

func SetAdminRoute(router gin.IRouter, admin *Api, mw *chatmw.MW, cfg *Config, client discovery.SvcDiscoveryRegistry) {
	router.GET("/ping", admin.AdminPing)
	adminRouterGroup := router.Group("/account")
	adminRouterGroup.POST("/login", admin.AdminLogin)
	adminRouterGroup.POST("/google_auth", mw.CheckAdmin, admin.GetGoogleAuth)
	adminRouterGroup.POST("/verify_google_auth", mw.CheckAdmin, admin.VerifyGoogleAuth)
	adminRouterGroup.POST("/update", mw.CheckAdmin, admin.AdminUpdateInfo)
	adminRouterGroup.POST("/info", mw.CheckAdmin, admin.AdminInfo)
	adminRouterGroup.POST("/change_password", mw.CheckAdmin, admin.ChangeAdminPassword)
	adminRouterGroup.POST("/add_admin", mw.CheckAdmin, mw.OperationLog, admin.AddAdminAccount)
	adminRouterGroup.POST("/add_user", mw.CheckAdmin, mw.OperationLog, admin.AddUserAccount)
	adminRouterGroup.POST("/del_admin", mw.CheckAdmin, mw.OperationLog, admin.DelAdminAccount)
	adminRouterGroup.POST("/search_admin", mw.CheckAdmin, admin.SearchAdminAccount)

	operationLogGroup := router.Group("/operation_log", mw.CheckAdmin)
	operationLogGroup.POST("/search", admin.SearchOperationLog)

	importGroup := router.Group("/user/import")
	importGroup.POST("/json", mw.CheckAdmin, mw.OperationLog, admin.ImportUserByJson)
	importGroup.POST("/xlsx", mw.CheckAdmin, mw.OperationLog, admin.ImportUserByXlsx)
	importGroup.GET("/xlsx", admin.BatchImportTemplate)

	allowRegisterGroup := router.Group("/user/allow_register", mw.CheckAdmin)
	allowRegisterGroup.POST("/get", admin.GetAllowRegister)
	allowRegisterGroup.POST("/set", admin.SetAllowRegister)

	defaultRouter := router.Group("/default", mw.CheckAdmin)
	defaultUserRouter := defaultRouter.Group("/user")
	defaultUserRouter.POST("/add", mw.OperationLog, admin.AddDefaultFriend)
	defaultUserRouter.POST("/del", mw.OperationLog, admin.DelDefaultFriend)
	defaultUserRouter.POST("/find", mw.OperationLog, admin.FindDefaultFriend)
	defaultUserRouter.POST("/search", admin.SearchDefaultFriend)
	defaultGroupRouter := defaultRouter.Group("/group")
	defaultGroupRouter.POST("/add", mw.OperationLog, admin.AddDefaultGroup)
	defaultGroupRouter.POST("/del", mw.OperationLog, admin.DelDefaultGroup)
	defaultGroupRouter.POST("/find", mw.OperationLog, admin.FindDefaultGroup)
	defaultGroupRouter.POST("/search", admin.SearchDefaultGroup)

	invitationCodeRouter := router.Group("/invitation_code", mw.CheckAdmin)
	invitationCodeRouter.POST("/add", mw.OperationLog, admin.AddInvitationCode)
	invitationCodeRouter.POST("/gen", mw.OperationLog, admin.GenInvitationCode)
	invitationCodeRouter.POST("/del", mw.OperationLog, admin.DelInvitationCode)
	invitationCodeRouter.POST("/search", admin.SearchInvitationCode)

	forbiddenRouter := router.Group("/forbidden", mw.CheckAdmin)
	ipForbiddenRouter := forbiddenRouter.Group("/ip")
	ipForbiddenRouter.POST("/add", mw.OperationLog, admin.AddIPForbidden)
	ipForbiddenRouter.POST("/del", mw.OperationLog, admin.DelIPForbidden)
	ipForbiddenRouter.POST("/search", admin.SearchIPForbidden)
	userForbiddenRouter := forbiddenRouter.Group("/user")
	userForbiddenRouter.POST("/add", mw.OperationLog, admin.AddUserIPLimitLogin)
	userForbiddenRouter.POST("/del", mw.OperationLog, admin.DelUserIPLimitLogin)
	userForbiddenRouter.POST("/search", admin.SearchUserIPLimitLogin)

	appletRouterGroup := router.Group("/applet", mw.CheckAdmin)
	appletRouterGroup.POST("/add", mw.OperationLog, admin.AddApplet)
	appletRouterGroup.POST("/del", mw.OperationLog, admin.DelApplet)
	appletRouterGroup.POST("/update", mw.OperationLog, admin.UpdateApplet)
	appletRouterGroup.POST("/search", admin.SearchApplet)

	blockRouter := router.Group("/block", mw.CheckAdmin)
	blockRouter.POST("/add", mw.OperationLog, admin.BlockUser)
	blockRouter.POST("/del", mw.OperationLog, admin.UnblockUser)
	blockRouter.POST("/search", admin.SearchBlockUser)

	userRouter := router.Group("/user", mw.CheckAdmin, mw.OperationLog)
	userRouter.POST("/password/reset", admin.ResetUserPassword)

	initGroup := router.Group("/client_config", mw.CheckAdmin)
	initGroup.POST("/get", admin.GetClientConfig)
	initGroup.POST("/set", mw.OperationLog, admin.SetClientConfig)
	initGroup.POST("/del", mw.OperationLog, admin.DelClientConfig)
	initGroup.POST("/list", admin.GetListClientConfig)

	smsConfigGroup := router.Group("/sms_config", mw.CheckAdmin)
	smsConfigGroup.POST("/set", mw.OperationLog, admin.SetSmsConfig)
	smsConfigGroup.POST("/get", admin.GetSmsConfig)

	bucketConfigGroup := router.Group("/bucket_config", mw.CheckAdmin)
	bucketConfigGroup.POST("/set", mw.OperationLog, admin.SetBucketConfig)
	bucketConfigGroup.POST("/get", admin.GetBucketConfig)

	signinConfigGroup := router.Group("/signin_config", mw.CheckAdmin)
	signinConfigGroup.POST("/set", mw.OperationLog, admin.SetSigninConfig)
	signinConfigGroup.POST("/get", admin.GetSigninConfig)

	adminMenuGroup := router.Group("/admin_menu", mw.CheckAdmin)
	adminMenuGroup.POST("/create", mw.OperationLog, admin.CreateAdminMenu)
	adminMenuGroup.POST("/update", mw.OperationLog, admin.UpdateAdminMenu)
	adminMenuGroup.POST("/delete", mw.OperationLog, admin.DeleteAdminMenu)
	adminMenuGroup.POST("/take", admin.TakeAdminMenu)
	adminMenuGroup.POST("/list", admin.ListAdminMenu)
	adminMenuGroup.POST("/user_menu", admin.ListAdminUserMenu)
	adminMenuGroup.POST("/user_menu/assign", mw.OperationLog, admin.AssignAdminUserMenu)
	adminMenuGroup.POST("/user_menu/get", admin.GetAdminUserMenu)

	statistic := router.Group("/statistic", mw.CheckAdmin)
	statistic.POST("/login_record", admin.LoginRecord)
	statistic.POST("/new_user_count", admin.NewUserCount)
	statistic.POST("/login_user_count", admin.LoginUserCount)

	applicationGroup := router.Group("application")
	applicationGroup.POST("/add_version", mw.CheckAdmin, mw.OperationLog, admin.AddApplicationVersion)
	applicationGroup.POST("/update_version", mw.CheckAdmin, mw.OperationLog, admin.UpdateApplicationVersion)
	applicationGroup.POST("/delete_version", mw.CheckAdmin, mw.OperationLog, admin.DeleteApplicationVersion)
	applicationGroup.POST("/latest_version", admin.LatestApplicationVersion)
	applicationGroup.POST("/page_versions", admin.PageApplicationVersion)

	var etcdClient *clientv3.Client
	if cfg.Discovery.Enable == kdisc.ETCDCONST {
		etcdClient = client.(*etcd.SvcDiscoveryRegistryImpl).GetClient()
	}
	cm := NewConfigManager(cfg.AllConfig, etcdClient, cfg.ConfigPath, cfg.RuntimeEnv)
	{
		configGroup := router.Group("/config", mw.CheckAdmin)
		configGroup.POST("/get_config_list", cm.GetConfigList)
		configGroup.POST("/get_config", cm.GetConfig)
		configGroup.POST("/set_config", cm.SetConfig)
		configGroup.POST("/reset_config", cm.ResetConfig)
	}
	{
		router.POST("/restart", mw.CheckAdmin, cm.Restart)
	}
}
