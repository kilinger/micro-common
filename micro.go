package common

import (
	"errors"
	"time"

	log "github.com/micro/go-log"
	micro "github.com/micro/go-micro"
)

type InitFunc func(configuration *Configuration) error

var svc micro.Service

// NewService new service with common initial
func NewService(version, defaultName string, initFunc InitFunc) micro.Service {
	service := micro.NewService(
		micro.Version(version),
		micro.BeforeStart(
			func() error {
				conf := GetConfiguration()

				err := conf.Load(defaultName)
				if err != nil {
					log.Log(err)
					return err
				}

				if !conf.initialized {
					err := errors.New("Configuration not initialized, check your config format")
					log.Log(err)
					return err
				}

				if conf.Name == "" {
					if defaultName == "" {
						err := errors.New("No name in configuration and no default name")
						log.Log(err)
						return err
					}
					conf.Name = defaultName
				}

				svc.Init(micro.Name(conf.Namespace + "." + conf.Name))

				svc.Init(micro.RegisterTTL(
					time.Duration(conf.GetInt("micro_register_ttl")) * time.Second,
				))
				svc.Init(micro.RegisterInterval(
					time.Duration(conf.GetInt("micro_register_interval")) * time.Second,
				))

				return initFunc(conf)
			},
		),
	)

	service.Init()

	svc = service

	return service
}

// NilInit empty InitFunc
func NilInit(conf *Configuration) error {
	return nil
}
