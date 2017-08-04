package common

import (
	"fmt"
	"sync"

	log "github.com/micro/go-log"

	"github.com/spf13/viper"
)

// Configuration configuration which read config from env
type Configuration struct {
	viper.Viper
	initialized bool
	m           sync.Mutex
	objects     map[string]interface{}

	Namespace string
	Name      string

	DatabaseDriver     string
	DatabaseDatasource string
}

// Load load configuration
func (c *Configuration) Load(name string) error {
	if name == "" {
		name = "micro"
	}

	c.SetConfigName("config")                       // name of config file (without extension)
	c.AddConfigPath(fmt.Sprintf("/etc/%s/", name))  // path to look for the config file in
	c.AddConfigPath(fmt.Sprintf("$HOME/.%s", name)) // call multiple times to add many search paths
	c.AddConfigPath(".")                            // optionally look for config in the working directory

	c.AutomaticEnv()

	err := c.ReadInConfig() // Find and read the config file
	if err != nil {         // Handle errors reading the config file
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			log.Log(err)
		} else {
			return fmt.Errorf("Fatal error config file: %s", err)
		}
	}

	c.Name = c.GetString("name")
	c.Namespace = c.GetString("namespace")

	// Let's set a default namespace because a lot of people don't care what it actually is
	if c.Namespace == "" {
		c.Namespace = "com.xxxxx"
	}

	c.initialized = true

	return nil
}

// SetObject set object with config
func (c *Configuration) SetObject(key string, object interface{}) {
	c.m.Lock()
	c.objects[key] = object
	c.m.Unlock()
}

// GetObject get object from config
func (c *Configuration) GetObject(key string) (interface{}, error) {
	object, ok := c.objects[key]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return object, nil
}

// GetBrokerTopic get path for topic
func (c *Configuration) GetBrokerTopic(topic string) string {
	return fmt.Sprintf("topic.%s.%s", c.Namespace, topic)
}

// GetServiceName get service name
func (c *Configuration) GetServiceName(name string) string {
	return fmt.Sprintf("%s.%s", c.Namespace, name)
}

// IsInitialized check configuration is initialized
func (c *Configuration) IsInitialized() bool {
	return c.initialized
}

var (
	config *Configuration
	once   sync.Once
)

// GetConfiguration return the Singleton of configuration
func GetConfiguration() *Configuration {
	once.Do(func() {
		v := viper.New()
		m := make(map[string]interface{})
		config = &Configuration{Viper: *v, objects: m}
	})

	return config
}
