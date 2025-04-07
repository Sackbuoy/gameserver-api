package gameserver

import "k8s.io/apimachinery/pkg/runtime/schema"

var GVC = schema.GroupVersionResource{
	Group:    "goopy.us",
	Version:  "v1",
	Resource: "gameservers",
}

type CreateRequest struct {
	GameType  string `json:"gameType"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type DeleteRequest struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// GameServer defines the structure for a GameServer resource
type GameServer struct {
	APIVersion string         `yaml:"apiVersion" json:"apiVersion"`
	Kind       string         `yaml:"kind" json:"kind"`
	Metadata   Metadata       `yaml:"metadata,omitempty" json:"metadata,omitempty"`
	Spec       GameServerSpec `yaml:"spec" json:"spec"`
}

type Metadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

// GameServerSpec defines the desired state of a GameServer
type GameServerSpec struct {
	GameType string `yaml:"gameType" json:"gameType"`
}

func New(name, namespace, gameType string) (*GameServer, error) {
	return &GameServer{
		Kind:       "GameServer",
		APIVersion: "goopy.us/v1",
		Metadata: Metadata{
			Name:      name,
			Namespace: namespace,
		},
		Spec: GameServerSpec{
			GameType: gameType,
		},
	}, nil
}
