package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/jubblin/omni-api/docs"
	"github.com/jubblin/omni-api/internal/api/handlers"
	omniclient "github.com/jubblin/omni-api/internal/client"
)

// @title           Talos Omni Control API
// @version         0.0.4
// @description     A REST API to interface with Sidero Omni.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

func main() {
	// Initialize Omni client
	client, err := omniclient.NewOmniClient()
	if err != nil {
		log.Fatalf("Failed to create Omni client: %v", err)
	}
	defer client.Close()

	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Middleware to record metrics
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		handlers.RecordRequest(c.FullPath(), duration, c.Writer.Status())
	})

	// Handlers
	clusterHandler := handlers.NewClusterHandler(client.Omni().State())
	machineHandler := handlers.NewMachineHandler(client.Omni().State())
	machineStatusHandler := handlers.NewMachineStatusHandler(client.Omni().State())
	machineLabelsHandler := handlers.NewMachineLabelsHandler(client.Omni().State())
	machineExtensionsHandler := handlers.NewMachineExtensionsHandler(client.Omni().State())
	machineUpgradeStatusHandler := handlers.NewMachineUpgradeStatusHandler(client.Omni().State())
	machineStatusMetricsHandler := handlers.NewMachineStatusMetricsHandler(client.Omni().State())
	machineSetHandler := handlers.NewMachineSetHandler(client.Omni().State())
	machineSetStatusHandler := handlers.NewMachineSetStatusHandler(client.Omni().State())
	configPatchHandler := handlers.NewConfigPatchHandler(client.Omni().State())
	clusterMachineHandler := handlers.NewClusterMachineHandler(client.Omni().State())
	clusterMachineStatusHandler := handlers.NewClusterMachineStatusHandler(client.Omni().State())
	clusterMachineConfigStatusHandler := handlers.NewClusterMachineConfigStatusHandler(client.Omni().State())
	clusterMachineTalosVersionHandler := handlers.NewClusterMachineTalosVersionHandler(client.Omni().State())
	kubeconfigHandler := handlers.NewKubeconfigHandler(client.Omni().State())
	kubernetesUpgradeHandler := handlers.NewKubernetesUpgradeHandler(client.Omni().State())
	clusterEndpointHandler := handlers.NewClusterEndpointHandler(client.Omni().State())
	etcdBackupHandler := handlers.NewEtcdBackupHandler(client.Omni().State())
	machineClassHandler := handlers.NewMachineClassHandler(client.Omni().State())
	machineSetNodeHandler := handlers.NewMachineSetNodeHandler(client.Omni().State())
	talosUpgradeHandler := handlers.NewTalosUpgradeHandler(client.Omni().State())
	schematicHandler := handlers.NewSchematicHandler(client.Omni().State())
	ongoingTaskHandler := handlers.NewOngoingTaskHandler(client.Omni().State())
	kubernetesStatusHandler := handlers.NewKubernetesStatusHandler(client.Omni().State())
	clusterKubernetesNodesHandler := handlers.NewClusterKubernetesNodesHandler(client.Omni().State())
	kubernetesVersionHandler := handlers.NewKubernetesVersionHandler(client.Omni().State())
	clusterMachineConfigHandler := handlers.NewClusterMachineConfigHandler(client.Omni().State())
	controlPlaneStatusHandler := handlers.NewControlPlaneStatusHandler(client.Omni().State())
	etcdBackupStatusHandler := handlers.NewEtcdBackupStatusHandler(client.Omni().State())
	etcdManualBackupHandler := handlers.NewEtcdManualBackupHandler(client.Omni().State())
	schematicConfigurationHandler := handlers.NewSchematicConfigurationHandler(client.Omni().State())
	extensionsConfigurationHandler := handlers.NewExtensionsConfigurationHandler(client.Omni().State())
	kernelArgsHandler := handlers.NewKernelArgsHandler(client.Omni().State())
	loadBalancerConfigHandler := handlers.NewLoadBalancerConfigHandler(client.Omni().State())
	loadBalancerStatusHandler := handlers.NewLoadBalancerStatusHandler(client.Omni().State())
	exposedServiceHandler := handlers.NewExposedServiceHandler(client.Omni().State())
	machineRequestSetHandler := handlers.NewMachineRequestSetHandler(client.Omni().State())
	clusterDiagnosticsHandler := handlers.NewClusterDiagnosticsHandler(client.Omni().State())
	clusterDestroyStatusHandler := handlers.NewClusterDestroyStatusHandler(client.Omni().State())
	machineSetDestroyStatusHandler := handlers.NewMachineSetDestroyStatusHandler(client.Omni().State())
	clusterWorkloadProxyStatusHandler := handlers.NewClusterWorkloadProxyStatusHandler(client.Omni().State())
	imagePullRequestHandler := handlers.NewImagePullRequestHandler(client.Omni().State())
	imagePullStatusHandler := handlers.NewImagePullStatusHandler(client.Omni().State())
	installationMediaHandler := handlers.NewInstallationMediaHandler(client.Omni().State())
	infraMachineConfigHandler := handlers.NewInfraMachineConfigHandler(client.Omni().State())
	machineConfigDiffHandler := handlers.NewMachineConfigDiffHandler(client.Omni().State())
	healthHandler := handlers.NewHealthHandler(client.Omni().State())
	metricsHandler := handlers.NewMetricsHandler()

	// Health and Metrics routes (outside v1 group for easier access)
	r.GET("/health", healthHandler.GetHealth)
	r.GET("/metrics", metricsHandler.GetMetrics)

	// API Routes
	v1 := r.Group("/api/v1")
	{
		// Cluster routes
		v1.GET("/clusters", clusterHandler.ListClusters)
		v1.GET("/clusters/:id", clusterHandler.GetCluster)
		v1.GET("/clusters/:id/status", clusterHandler.GetClusterStatus)
		v1.GET("/clusters/:id/metrics", clusterHandler.GetClusterMetrics)
		v1.GET("/clusters/:id/bootstrap", clusterHandler.GetClusterBootstrap)
		v1.GET("/clusters/:id/kubeconfig", kubeconfigHandler.GetKubeconfig)
		v1.GET("/clusters/:id/kubernetes-upgrade", kubernetesUpgradeHandler.GetKubernetesUpgradeStatus)
		v1.GET("/clusters/:id/talos-upgrade", talosUpgradeHandler.GetTalosUpgradeStatus)
		v1.GET("/clusters/:id/endpoints", clusterEndpointHandler.GetClusterEndpoints)
		v1.GET("/clusters/:id/kubernetes-status", kubernetesStatusHandler.GetKubernetesStatus)
		v1.GET("/clusters/:id/kubernetes-nodes", clusterKubernetesNodesHandler.ListClusterKubernetesNodes)
		v1.GET("/clusters/:id/kubernetes-nodes/:node", clusterKubernetesNodesHandler.GetClusterKubernetesNode)
		v1.GET("/clusters/:id/controlplane-status", controlPlaneStatusHandler.GetControlPlaneStatus)
		v1.GET("/clusters/:id/diagnostics", clusterDiagnosticsHandler.GetClusterDiagnostics)
		v1.GET("/clusters/:id/destroy-status", clusterDestroyStatusHandler.GetClusterDestroyStatus)
		v1.GET("/clusters/:id/workload-proxy-status", clusterWorkloadProxyStatusHandler.GetClusterWorkloadProxyStatus)
		
		// Machine routes
		v1.GET("/machines", machineHandler.ListMachines)
		v1.GET("/machines/:id", machineHandler.GetMachine)
		v1.GET("/machines/:id/status", machineStatusHandler.GetMachineStatus)
		v1.GET("/machines/:id/labels", machineLabelsHandler.GetMachineLabels)
		v1.GET("/machines/:id/extensions", machineExtensionsHandler.GetMachineExtensions)
		v1.GET("/machines/:id/upgrade-status", machineUpgradeStatusHandler.GetMachineUpgradeStatus)
		v1.GET("/machines/:id/metrics", machineStatusMetricsHandler.GetMachineStatusMetrics)
		v1.GET("/machines/:id/config-diff", machineConfigDiffHandler.GetMachineConfigDiff)
		
		// MachineSet routes
		v1.GET("/machinesets", machineSetHandler.ListMachineSets)
		v1.GET("/machinesets/:id", machineSetHandler.GetMachineSet)
		v1.GET("/machinesets/:id/status", machineSetStatusHandler.GetMachineSetStatus)
		v1.GET("/machinesets/:id/destroy-status", machineSetDestroyStatusHandler.GetMachineSetDestroyStatus)
		
		// MachineSetNode routes
		v1.GET("/machinesetnodes", machineSetNodeHandler.ListMachineSetNodes)
		v1.GET("/machinesetnodes/:id", machineSetNodeHandler.GetMachineSetNode)
		
		// ConfigPatch routes
		v1.GET("/configpatches", configPatchHandler.ListConfigPatches)
		v1.GET("/configpatches/:id", configPatchHandler.GetConfigPatch)
		
		// ClusterMachine routes
		v1.GET("/clustermachines", clusterMachineHandler.ListClusterMachines)
		v1.GET("/clustermachines/:id", clusterMachineHandler.GetClusterMachine)
		v1.GET("/clustermachines/:id/status", clusterMachineStatusHandler.GetClusterMachineStatus)
		v1.GET("/clustermachines/:id/config-status", clusterMachineConfigStatusHandler.GetClusterMachineConfigStatus)
		v1.GET("/clustermachines/:id/talos-version", clusterMachineTalosVersionHandler.GetClusterMachineTalosVersion)
		v1.GET("/clustermachines/:id/config", clusterMachineConfigHandler.GetClusterMachineConfig)
		
		// MachineClass routes
		v1.GET("/machineclasses", machineClassHandler.ListMachineClasses)
		v1.GET("/machineclasses/:id", machineClassHandler.GetMachineClass)
		
		// EtcdBackup routes
		v1.GET("/etcdbackups", etcdBackupHandler.ListEtcdBackups)
		v1.GET("/etcdbackups/:id", etcdBackupHandler.GetEtcdBackup)
		v1.GET("/etcdbackups/:id/status", etcdBackupStatusHandler.GetEtcdBackupStatus)
		v1.GET("/etcd-manual-backups", etcdManualBackupHandler.ListEtcdManualBackups)
		v1.GET("/etcd-manual-backups/:id", etcdManualBackupHandler.GetEtcdManualBackup)
		
		// Schematic routes
		v1.GET("/schematics", schematicHandler.ListSchematics)
		v1.GET("/schematics/:id", schematicHandler.GetSchematic)
		v1.GET("/schematic-configurations", schematicConfigurationHandler.ListSchematicConfigurations)
		v1.GET("/schematic-configurations/:id", schematicConfigurationHandler.GetSchematicConfiguration)
		
		// OngoingTask routes
		v1.GET("/ongoingtasks", ongoingTaskHandler.ListOngoingTasks)
		v1.GET("/ongoingtasks/:id", ongoingTaskHandler.GetOngoingTask)
		
		// Kubernetes Version routes
		v1.GET("/kubernetes-versions", kubernetesVersionHandler.ListKubernetesVersions)
		v1.GET("/kubernetes-versions/:id", kubernetesVersionHandler.GetKubernetesVersion)
		
		// Extensions Configuration routes
		v1.GET("/extensions-configurations", extensionsConfigurationHandler.ListExtensionsConfigurations)
		v1.GET("/extensions-configurations/:id", extensionsConfigurationHandler.GetExtensionsConfiguration)
		
		// Kernel Args routes
		v1.GET("/kernel-args", kernelArgsHandler.ListKernelArgs)
		v1.GET("/kernel-args/:id", kernelArgsHandler.GetKernelArgs)
		
		// Load Balancer routes
		v1.GET("/loadbalancer-configs", loadBalancerConfigHandler.ListLoadBalancerConfigs)
		v1.GET("/loadbalancer-configs/:id", loadBalancerConfigHandler.GetLoadBalancerConfig)
		v1.GET("/loadbalancers/:id/status", loadBalancerStatusHandler.GetLoadBalancerStatus)
		
		// Exposed Service routes
		v1.GET("/exposed-services", exposedServiceHandler.ListExposedServices)
		v1.GET("/exposed-services/:id", exposedServiceHandler.GetExposedService)
		
		// Machine Request Set routes
		v1.GET("/machine-request-sets", machineRequestSetHandler.ListMachineRequestSets)
		v1.GET("/machine-request-sets/:id", machineRequestSetHandler.GetMachineRequestSet)
		
		// Image Pull Request routes
		v1.GET("/image-pull-requests", imagePullRequestHandler.ListImagePullRequests)
		v1.GET("/image-pull-requests/:id", imagePullRequestHandler.GetImagePullRequest)
		v1.GET("/image-pull-requests/:id/status", imagePullStatusHandler.GetImagePullStatus)
		
		// Installation Media routes
		v1.GET("/installation-medias", installationMediaHandler.ListInstallationMedias)
		v1.GET("/installation-medias/:id", installationMediaHandler.GetInstallationMedia)
		
		// Infrastructure Machine Config routes
		v1.GET("/infra-machine-configs", infraMachineConfigHandler.ListInfraMachineConfigs)
		v1.GET("/infra-machine-configs/:id", infraMachineConfigHandler.GetInfraMachineConfig)
	}

	// Redirect root to Swagger UI
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

