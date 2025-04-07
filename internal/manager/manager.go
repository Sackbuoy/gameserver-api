package manager

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/sackbuoy/gameserver-api/internal/gameserver"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	k8syaml "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/client-go/util/homedir"
)

// Manages Gameserver CRD instances
type Manager struct {
	client *dynamic.DynamicClient
}

func New() (*Manager, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// Try to use in-cluster config first, fall back to kubeconfig file
	config, err := rest.InClusterConfig()
	if err != nil {
		// Create config from kubeconfig file
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	// Create a dynamic client for working with custom resources
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Manager{
		client: client,
	}, nil
}

func (m *Manager) Create(c *gin.Context) {
	var request gameserver.CreateRequest

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	err = json.Unmarshal(jsonData, &request)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	gs, err := gameserver.New(request.Name, request.Namespace, request.GameType)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	yamlData, err := yaml.Marshal(gs)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	err = m.CreateFromYAML(yamlData)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}

func (m *Manager) CreateFromYAML(yamlBytes []byte) error {
	decoder := k8syaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

	obj := &unstructured.Unstructured{}
	_, _, err := decoder.Decode(yamlBytes, nil, obj)
	if err != nil {
		return fmt.Errorf("error decoding YAML: %w", err)
	}

	namespace := obj.GetNamespace()
	if namespace == "" {
		namespace = "default"
	}

	_, err = m.client.Resource(gameserver.GVC).Namespace(namespace).Create(
		context.Background(), obj, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("error creating resource: %w", err)
	}

	return nil
}

// TODO
func (m *Manager) Read(c *gin.Context) {
}

// TODO
func (m *Manager) Update(c *gin.Context) {
}

func (m *Manager) Delete(c *gin.Context) {
	var request gameserver.DeleteRequest

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	err = json.Unmarshal(jsonData, &request)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	err = m.client.Resource(gameserver.GVC).
		Namespace(request.Namespace).
		Delete(c, request.Name, metav1.DeleteOptions{})
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}
