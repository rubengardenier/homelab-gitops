package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

type NodeInfo struct {
	Name         string  `json:"name"`
	Status       string  `json:"status"`
	Role         string  `json:"role"`
	IP           string  `json:"ip"`
	CPUCores     string  `json:"cpuCores"`
	MemoryTotal  string  `json:"memoryTotal"`
	CPUUsage     string  `json:"cpuUsage"`
	MemoryUsage  string  `json:"memoryUsage"`
	CPUPercent   float64 `json:"cpuPercent"`
	MemoryPercent float64 `json:"memoryPercent"`
}

type ClusterStats struct {
	TotalPods     int `json:"totalPods"`
	RunningPods   int `json:"runningPods"`
	TotalServices int `json:"totalServices"`
	TotalNodes    int `json:"totalNodes"`
	ReadyNodes    int `json:"readyNodes"`
}

type VersionInfo struct {
	Kubernetes string `json:"kubernetes"`
	Talos      string `json:"talos"`
}

type DashboardData struct {
	Nodes    []NodeInfo   `json:"nodes"`
	Stats    ClusterStats `json:"stats"`
	Versions VersionInfo  `json:"versions"`
	Updated  string       `json:"updated"`
}

var (
	clientset        *kubernetes.Clientset
	metricsClientset *metricsv.Clientset
)

func init() {
	var config *rest.Config
	var err error

	// Try in-cluster config first
	config, err = rest.InClusterConfig()
	if err != nil {
		// Fall back to kubeconfig for local development
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Printf("Warning: Could not create Kubernetes config: %v", err)
			return
		}
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Warning: Could not create Kubernetes client: %v", err)
		return
	}

	metricsClientset, err = metricsv.NewForConfig(config)
	if err != nil {
		log.Printf("Warning: Could not create metrics client: %v", err)
	}
}

func getNodes(ctx context.Context) ([]NodeInfo, error) {
	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Get node metrics if available
	var nodeMetrics map[string]struct {
		CPU    string
		Memory string
	}
	if metricsClientset != nil {
		metrics, err := metricsClientset.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})
		if err == nil {
			nodeMetrics = make(map[string]struct {
				CPU    string
				Memory string
			})
			for _, m := range metrics.Items {
				nodeMetrics[m.Name] = struct {
					CPU    string
					Memory string
				}{
					CPU:    m.Usage.Cpu().String(),
					Memory: m.Usage.Memory().String(),
				}
			}
		}
	}

	var result []NodeInfo
	for _, node := range nodes.Items {
		info := NodeInfo{
			Name:        node.Name,
			CPUCores:    node.Status.Capacity.Cpu().String(),
			MemoryTotal: node.Status.Capacity.Memory().String(),
		}

		// Get node IP
		for _, addr := range node.Status.Addresses {
			if addr.Type == "InternalIP" {
				info.IP = addr.Address
				break
			}
		}

		// Get node role
		if _, ok := node.Labels["node-role.kubernetes.io/control-plane"]; ok {
			info.Role = "control-plane"
		} else {
			info.Role = "worker"
		}

		// Get node status
		info.Status = "NotReady"
		for _, cond := range node.Status.Conditions {
			if cond.Type == "Ready" && cond.Status == "True" {
				info.Status = "Ready"
				break
			}
		}

		// Add metrics if available
		if nodeMetrics != nil {
			if m, ok := nodeMetrics[node.Name]; ok {
				info.CPUUsage = m.CPU
				info.MemoryUsage = m.Memory

				// Calculate percentages
				cpuCapacity := node.Status.Capacity.Cpu().MilliValue()
				memCapacity := node.Status.Capacity.Memory().Value()

				if metrics, err := metricsClientset.MetricsV1beta1().NodeMetricses().Get(ctx, node.Name, metav1.GetOptions{}); err == nil {
					cpuUsage := metrics.Usage.Cpu().MilliValue()
					memUsage := metrics.Usage.Memory().Value()

					if cpuCapacity > 0 {
						info.CPUPercent = float64(cpuUsage) / float64(cpuCapacity) * 100
					}
					if memCapacity > 0 {
						info.MemoryPercent = float64(memUsage) / float64(memCapacity) * 100
					}
				}
			}
		}

		result = append(result, info)
	}

	return result, nil
}

func getVersions(ctx context.Context) VersionInfo {
	var versions VersionInfo

	// Get Kubernetes version
	serverVersion, err := clientset.Discovery().ServerVersion()
	if err == nil {
		versions.Kubernetes = serverVersion.GitVersion
	}

	// Get Talos version from first node's label
	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err == nil && len(nodes.Items) > 0 {
		if talosVersion, ok := nodes.Items[0].Labels["node.talos.dev/os-version"]; ok {
			versions.Talos = talosVersion
		}
	}

	return versions
}

func getClusterStats(ctx context.Context) (ClusterStats, error) {
	var stats ClusterStats

	// Get pods
	pods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err == nil {
		stats.TotalPods = len(pods.Items)
		for _, pod := range pods.Items {
			if pod.Status.Phase == "Running" {
				stats.RunningPods++
			}
		}
	}

	// Get services
	services, err := clientset.CoreV1().Services("").List(ctx, metav1.ListOptions{})
	if err == nil {
		stats.TotalServices = len(services.Items)
	}

	// Get nodes
	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err == nil {
		stats.TotalNodes = len(nodes.Items)
		for _, node := range nodes.Items {
			for _, cond := range node.Status.Conditions {
				if cond.Type == "Ready" && cond.Status == "True" {
					stats.ReadyNodes++
					break
				}
			}
		}
	}

	return stats, nil
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	if clientset == nil {
		http.Error(w, "Kubernetes client not initialized", http.StatusServiceUnavailable)
		return
	}

	ctx := r.Context()

	nodes, err := getNodes(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stats, err := getClusterStats(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	versions := getVersions(ctx)

	data := DashboardData{
		Nodes:    nodes,
		Stats:    stats,
		Versions: versions,
		Updated:  "just now",
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(data)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// noCacheFileServer wraps http.FileServer to add cache-control headers
func noCacheFileServer(fs http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		fs.ServeHTTP(w, r)
	})
}

func main() {
	// API endpoints
	http.HandleFunc("/api/cluster", apiHandler)
	http.HandleFunc("/health", healthHandler)

	// Serve static files with no-cache headers
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = "./static"
	}
	http.Handle("/", noCacheFileServer(http.FileServer(http.Dir(staticDir))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s", port)
	log.Printf("Serving static files from %s", staticDir)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
