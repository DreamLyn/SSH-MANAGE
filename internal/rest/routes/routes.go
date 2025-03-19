package routes

// var (
// 	certificateSvc *certificate.CertificateService
// 	workflowSvc    *workflow.WorkflowService
// 	statisticsSvc  *statistics.StatisticsService
// 	notifySvc      *notify.NotifyService
// )

// func Register(router *router.Router[*core.RequestEvent]) {
// 	certificateRepo := repository.NewCertificateRepository()
// 	certificateSvc = certificate.NewCertificateService(certificateRepo)

// 	workflowRepo := repository.NewWorkflowRepository()
// 	workflowRunRepo := repository.NewWorkflowRunRepository()
// 	workflowSvc = workflow.NewWorkflowService(workflowRepo, workflowRunRepo)

// 	statisticsRepo := repository.NewStatisticsRepository()
// 	statisticsSvc = statistics.NewStatisticsService(statisticsRepo)

// 	settingsRepo := repository.NewSettingsRepository()
// 	notifySvc = notify.NewNotifyService(settingsRepo)

// 	group := router.Group("/api")
// 	group.Bind(apis.RequireSuperuserAuth())
// 	handlers.NewCertificateHandler(group, certificateSvc)
// 	handlers.NewWorkflowHandler(group, workflowSvc)
// 	handlers.NewStatisticsHandler(group, statisticsSvc)
// 	handlers.NewNotifyHandler(group, notifySvc)
// }

// func Unregister() {
// 	if workflowSvc != nil {
// 		workflowSvc.Shutdown(context.Background())
// 	}
// }
