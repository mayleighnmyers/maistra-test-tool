package app

import (
	"github.com/maistra/maistra-test-tool/pkg/util/env"
	"github.com/maistra/maistra-test-tool/pkg/util/oc"
	"github.com/maistra/maistra-test-tool/pkg/util/test"
)

type nginx struct {
	ns   string
	mTLS bool
}

var _ App = &nginx{}

func Nginx(ns string) App {
	return &nginx{ns: ns, mTLS: false}
}

func NginxWithMTLS(ns string) App {
	return &nginx{ns: ns, mTLS: true}
}

func (a *nginx) Name() string {
	return "nginx"
}

func (a *nginx) Namespace() string {
	return a.ns
}

func (a *nginx) Install(t test.TestHelper) {
	t.T().Helper()
	oc.CreateGenericSecretFromFiles(t, a.Namespace(),
		"nginx-ca-certs",
		"example.com.crt="+nginxServerCACertFile)

	if a.mTLS {
		oc.CreateTLSSecret(t, a.Namespace(), "nginx-server-certs", meshExtServerCertKeyFile, meshExtServerCertFile)
		oc.CreateConfigMapFromFiles(t, a.Namespace(),
			"nginx-configmap",
			"nginx.conf="+nginxConfMTlsFile)
	} else {
		oc.CreateTLSSecret(t, a.Namespace(), "nginx-server-certs", nginxServerCertKeyFile, nginxServerCertFile)
		oc.CreateConfigMapFromFiles(t, a.Namespace(),
			"nginx-configmap",
			"nginx.conf="+nginxConfFile)
	}

	oc.ApplyFile(t, a.Namespace(), nginxYamlFile)
}

func (a *nginx) Uninstall(t test.TestHelper) {
	t.T().Helper()
	oc.DeleteFile(t, a.Namespace(), nginxYamlFile)
	oc.DeleteSecret(t, a.Namespace(), "nginx-server-certs")
	oc.DeleteSecret(t, a.Namespace(), "nginx-ca-certs")
	oc.DeleteConfigMap(t, a.Namespace(), "nginx-configmap")
}

func (a *nginx) WaitReady(t test.TestHelper) {
	t.T().Helper()
	oc.WaitDeploymentRolloutComplete(t, a.ns, "my-nginx")
}

var (
	rootDir                  = env.GetRootDir()
	nginxYamlFile            = rootDir + "/pkg/app/yaml/nginx.yaml"
	nginxConfMTlsFile        = rootDir + "/pkg/app/yaml/nginx_mesh_external_ssl.conf"
	nginxConfFile            = rootDir + "/pkg/app/yaml/nginx.conf"
	nginxServerCertKeyFile   = rootDir + "/sampleCerts/nginx.example.com/nginx.example.com.key"
	nginxServerCertFile      = rootDir + "/sampleCerts/nginx.example.com/nginx.example.com.crt"
	nginxServerCACertFile    = rootDir + "/sampleCerts/nginx.example.com/example.com.crt"
	meshExtServerCertKeyFile = rootDir + "/sampleCerts/nginx.example.com/my-nginx.mesh-external.svc.cluster.local.key"
	meshExtServerCertFile    = rootDir + "/sampleCerts/nginx.example.com/my-nginx.mesh-external.svc.cluster.local.crt"
)
