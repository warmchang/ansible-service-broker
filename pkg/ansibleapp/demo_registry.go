package ansibleapp

import (
	"encoding/json"
	"fmt"

	logging "github.com/op/go-logging"
)

var DemoJson = `
[
   {
      "description" : "Note taking web application",
      "parameters" : null,
      "async" : "optional",
      "id" : "c9e3348c-b6fd-45d6-b2ba-a13b41f0b6e3",
      "bindable" : true,
      "name" : "apb/etherpad-ansibleapp"
   },
   {
      "description" : "Wordpress is a website thing",
      "parameters" : null,
      "name" : "apb/wordpress-ansibleapp",
      "id" : "138afbe2-4e33-4faa-b6ae-2e8d61f2bf1f",
      "async" : "optional",
      "bindable" : true
   },
   {
      "parameters" : null,
      "description" : "Another website thing",
      "async" : "required",
      "id" : "507856c7-8e08-460c-a0b2-7d8af8c4318d",
      "bindable" : false,
      "name" : "apb/drupal-ansibleapp"
   },
   {
      "parameters" : null,
      "description" : "ManageIQ",
      "async" : "required",
      "id" : "68e3fc62-22c3-40c8-9be3-43446c2d0cab",
      "bindable" : false,
      "name" : "apb/drupal-ansibleapp"
   },
   {
      "description" : "Blogging and stuff",
      "parameters" : null,
      "async" : "unsupported",
      "id" : "f819ddd5-37d1-4698-b3fb-a3cc99a35d2e",
      "bindable" : false,
      "name" : "apb/ghost-ansibleapp"
   }
]
`

type DemoRegistry struct {
	config RegistryConfig
	log    *logging.Logger
}

func (r *DemoRegistry) Init(config RegistryConfig, log *logging.Logger) error {
	log.Debug("DemoRegistry::Init")
	r.config = config
	r.log = log
	return nil
}

func (r *DemoRegistry) LoadSpecs() ([]*Spec, error) {
	r.log.Debug("DemoRegistry::LoadSpecs")
	specs := loadSpecs([]byte(DemoJson))
	r.log.Debug(fmt.Sprintf("Loaded Specs: %v", specs))
	r.log.Info(fmt.Sprintf("Loaded [ %d ] specs from %s registry", len(specs), r.config.Name))
	return specs, nil
}

func loadDemoSpecs(rawPayload []byte) []*Spec {
	var specs []*Spec
	json.Unmarshal(rawPayload, &specs)
	return specs
}
