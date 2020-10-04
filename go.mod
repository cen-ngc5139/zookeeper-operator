module github.com/ghostbaby/zookeeper-operator

go 1.13

require (
	github.com/anacrolix/log v0.7.0 // indirect
	github.com/crossplane/crossplane-runtime v0.9.0
	github.com/go-logr/logr v0.2.1-0.20200730175230-ee2de8da5be6
	github.com/googleapis/gnostic v0.4.1
	github.com/gorilla/mux v1.7.5-0.20200711200521-98cb6bf42e08
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus-operator/prometheus-operator v0.42.1
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.42.1
	github.com/samuel/go-zookeeper v0.0.0-20200724154423-2164a8ac840e
	github.com/sirupsen/logrus v1.6.0
	github.com/toolkits/file v0.0.0-20160325033739-a5b3c5147e07
	gopkg.in/fatih/set.v0 v0.2.1
	k8s.io/api v0.19.0
	k8s.io/apimachinery v0.19.0
	k8s.io/cli-runtime v0.19.0
	k8s.io/client-go v11.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.6.1-0.20200804124940-17eebbff0d48
)

replace k8s.io/client-go v11.0.0+incompatible => k8s.io/client-go v0.0.0-20200813012017-e7a1d9ada0d5

replace github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring => github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.42.1
