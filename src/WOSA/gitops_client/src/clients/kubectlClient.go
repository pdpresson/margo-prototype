package clients

import (
	"context"
	"flag"
	"fmt"
	"gitops_client/appConfig"
	"gitops_client/models"
	"os"
	"path/filepath"

	"github.com/go-logr/logr"
	yaml2 "gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

type KubectlClient struct {
	log            *logr.Logger
	appConfig      *appConfig.AppConfig
	restConfig     *rest.Config
	restMapper     *restmapper.DeferredDiscoveryRESTMapper
	dicoveryClient *discovery.DiscoveryClient
	dynamicClient  *dynamic.DynamicClient
}

var (
	decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
)

func (c *KubectlClient) Init() error {
	rest.InClusterConfig()
	c.log.Info("Initializing Kubectl Client")
	kubeConfig := flag.String("kubeconfig", os.Getenv("KUBECONFIG"), "(optional) absolute path to the kubeconfig file")
	flag.Parse()

	var restConfig *rest.Config
	var cErr error
	if c.appConfig.InCluster {
		restConfig, cErr = rest.InClusterConfig()
		if cErr != nil {
			c.log.Error(cErr, "Error getting config from cluster")
			return cErr
		}
	} else {
		restConfig, cErr = clientcmd.BuildConfigFromFlags("", *kubeConfig)
		if cErr != nil {
			c.log.Error(cErr, "Error building config from flags")
			return cErr
		}
	}
	c.restConfig = restConfig

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(restConfig)
	if err != nil {
		c.log.Error(err, "Error creating discovery client")
		return err
	}
	c.dicoveryClient = discoveryClient

	c.restMapper = restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(discoveryClient))

	dynamicClient, err := dynamic.NewForConfig(c.restConfig)
	if err != nil {
		c.log.Error(err, "failed to create a dynmic client")
		return err
	}
	c.dynamicClient = dynamicClient

	return nil
}

func (c *KubectlClient) Apply(desiredState models.DesiredState, kind string) error {
	c.log.Info(fmt.Sprintf("Applying %s resource", kind))

	if desiredState.State != models.New {
		c.log.Info("TODO: Handle more than just new state")
		return nil
	}

	b, err := os.ReadFile(filepath.Join(desiredState.SourceFile))
	if err != nil {
		c.log.Error(err, "Error reading the current state file")
		return err
	}

	var namespace string
	if kind == models.KindSolution {
		solution := &models.Solution{}
		err = yaml2.Unmarshal(b, solution)
		if err != nil {
			c.log.Error(err, "Error unmarshalling solution")
			return err
		}
		namespace = solution.Metadata.Namespace
	} else if kind == models.KindInstance {
		instance := &models.Instance{}
		err = yaml2.Unmarshal(b, instance)
		if err != nil {
			c.log.Error(err, "Error unmarshalling instance")
			return err
		}
		namespace = instance.Metadata.Namespace
	}

	resource := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode(b, nil, resource)
	if err != nil {
		c.log.Error(err, "Error decoding")
		return err
	}

	mapping, err := c.restMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		c.log.Error(err, "Error with REST mapping")
	}

	result, err := c.dynamicClient.Resource(mapping.Resource).Namespace(namespace).Create(context.TODO(), resource, metav1.CreateOptions{})
	if err != nil {
		c.log.Error(err, "Error creating resource")
		return err

	}
	c.log.Info(fmt.Sprintf("Created resource %q.\n", result.GetName()))

	return nil
}

func NewKubectlClient(l *logr.Logger, c *appConfig.AppConfig) *KubectlClient {
	return &KubectlClient{
		log:       l,
		appConfig: c,
	}
}
