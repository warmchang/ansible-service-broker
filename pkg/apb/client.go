package apb

import (
	"encoding/json"
	"fmt"
	"os"

	docker "github.com/fsouza/go-dockerclient"
	logging "github.com/op/go-logging"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	restclient "k8s.io/client-go/rest"

	"github.com/pborman/uuid"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

/*
parameters will be 2 keys

answers {}
kubecfg {}

deprovision - delete the namespace and it tears the whole thing down.

oc delete?


route will be hardcoded, need to determine how to get that from the apb.


need to pass in cert through parameters


First cut might have to pass kubecfg from broker. FIRST SPRINT broker passes username and password.

admin/admin
*/

var DockerSocket = "unix:///var/run/docker.sock"

type ClusterConfig struct {
	Target   string
	User     string
	Password string `yaml:"pass"`
}

type Client struct {
	dockerClient  *docker.Client
	ClusterClient *clientset.Clientset
	RESTClient    restclient.Interface
	log           *logging.Logger
}

func createClientConfigFromFile(configPath string) (*restclient.Config, error) {
	clientConfig, err := clientcmd.LoadFromFile(configPath)
	if err != nil {
		return nil, err
	}

	config, err := clientcmd.NewDefaultClientConfig(*clientConfig, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func NewClient(log *logging.Logger) (*Client, error) {
	dockerClient, err := docker.NewClient(DockerSocket)
	if err != nil {
		log.Error("Could not load docker client")
		return nil, err
	}

	// NOTE: Both the external and internal client object are using the same
	// clientset library. Internal clientset normally uses a different
	// library
	clientConfig, err := restclient.InClusterConfig()
	if err != nil {
		log.Warning("Failed to create a InternalClientSet: %v.", err)

		log.Debug("Checking for a local Cluster Config")
		clientConfig, err = createClientConfigFromFile(homedir.HomeDir() + "/.kube/config")
		if err != nil {
			log.Error("Failed to create LocalClientSet")
			return nil, err
		}
	}

	clientset, err := clientset.NewForConfig(clientConfig)
	if err != nil {
		log.Error("Failed to create LocalClientSet")
		return nil, err
	}

	rest := clientset.CoreV1().RESTClient()

	client := &Client{
		dockerClient:  dockerClient,
		ClusterClient: clientset,
		RESTClient:    rest,
		log:           log,
	}

	return client, nil
}

func (c *Client) RunImage(
	action string,
	clusterConfig ClusterConfig,
	spec *Spec,
	p *Parameters,
) (string, error) {
	params, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	err = c.refreshLoginToken(clusterConfig)
	if err != nil {
		c.log.Error("Error occurred while refreshing login token! Aborting apb run.")
		c.log.Error(err.Error())
		return "", err
	}
	c.log.Notice("Login token successfully refreshed.")

	c.log.Debug("Running OC run...")
	c.log.Debug("clusterConfig:")
	c.log.Debug("target: [ %s ]", clusterConfig.Target)
	c.log.Debug("user: [ %s ]", clusterConfig.User)
	c.log.Debug("password:[ %s ]", clusterConfig.Password)
	c.log.Debug("name:[ %s ]", spec.Name)
	c.log.Debug("image:[ %s ]", spec.Image)
	c.log.Debug("action:[ %s ]", action)
	c.log.Debug("params:[ %s ]", string(params))

	//TODO: This can be cleaned up much more. This Function does a lot more
	// than RunImage originally did.
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("aa-%s", uuid.New()),
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "apb",
					Image: spec.Image,
					Args: []string{
						action,
						"--extra-vars",
						string(params),
					},
					Env: []v1.EnvVar{{
						Name:  "OPENSHIFT_TARGET",
						Value: clusterConfig.Target,
					}, {
						Name:  "OPENSHIFT_USER",
						Value: clusterConfig.User,
					}, {
						Name:  "OPENSHIFT_PASS",
						Value: clusterConfig.Password,
					}},
					ImagePullPolicy: v1.PullAlways,
				},
			},
			RestartPolicy: v1.RestartPolicyNever,
		},
	}
	c.log.Notice(fmt.Sprintf("Creating pod %q in namespace default", pod.Name))
	_, err = c.ClusterClient.CoreV1().Pods("default").Create(pod)

	return pod.Name, err
}

func (c *Client) PullImage(imageName string) error {
	// Under what circumstances does this error out?
	c.dockerClient.PullImage(docker.PullImageOptions{
		Repository:   imageName,
		OutputStream: os.Stdout,
	}, docker.AuthConfiguration{})
	return nil
}

func (c *Client) refreshLoginToken(clusterConfig ClusterConfig) error {
	c.log.Debug("Refreshing login token...")
	c.log.Debug("target: [ %s ]", clusterConfig.Target)
	c.log.Debug("user: [ %s ]", clusterConfig.User)
	c.log.Debug("password:[ %s ]", clusterConfig.Password)

	output, err := runCommand(
		"oc", "login", "--insecure-skip-tls-verify", clusterConfig.Target,
		"-u", clusterConfig.User,
		"-p", clusterConfig.Password,
	)

	if err != nil {
		return err
	}

	c.log.Debug("No error reported after running oc login. Cmd output:")
	c.log.Debug(string(output))
	return nil
}
