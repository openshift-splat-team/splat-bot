package controllers

import (
	"bytes"
	"fmt"
	t "text/template"
)

const installConfigTemplate = `apiVersion: v1
baseDomain: vmc-ci.devcluster.openshift.com
compute:
- architecture: amd64
  hyperthreading: Enabled
  name: worker
  platform: 
    vsphere: {}
controlPlane:
  architecture: amd64
  hyperthreading: Enabled
  name: master
  platform:
    vsphere: {}
metadata:
  name: {{ .ClusterName }}
networking:
  machineNetwork:
  - cidr: {{ .MachineNetwork }}
platform:
  vsphere: 
    apiVIPs:
    - {{ .ApiVIP }}
    ingressVIPs:
    - {{ .IngressVIP }}
	vcenters:
	- server: {{ .Server }}
      user: {{ .Username }}
	  password: {{ .Password }}
	  datacenters:
	  - {{ .Datacenter }}
    failureDomains: 
    - name: fd-1
      region: us-west
      server: {{ .Server }}
      topology:
        computeCluster: {{ .ComputeCluster }}
        datacenter: {{ .Datacenter }}
        datastore: {{ .Datastore }}
        networks:
        - {{ .Network }}
      zone: us-west-1a
pullSecret: <your pull secret>
sshKey: |
  <your public key>`

var template *t.Template

func init() {
	var err error
	template, err = t.New("yamlTemplate").Parse(installConfigTemplate)
	if err != nil {
		panic(fmt.Errorf("failed to parse install config template: %v", err))
	}
}

func RenderInstallConfig(config map[string]string) (string, error) {
	var renderedTemplate bytes.Buffer
	err := template.Execute(&renderedTemplate, config)
	if err != nil {
		return "", fmt.Errorf("failed to render install config template: %v", err)
	}
	return fmt.Sprintf("```%s```", renderedTemplate.String()), nil
}
