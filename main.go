package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var (
	BuildVersion = "N/A"
	BuildTime    = "N/A"
)

type HookParamters struct {
	certDir                string
	port                   int
}

func main() {
	var params HookParamters
	var version bool

	flag.IntVar(&params.port, "port", 8443, "Wehbook port")
	flag.StringVar(&params.certDir, "cert-dir", "/certs/", "Wehbook certificate folder")

	flag.Parse()
	if version {
		fmt.Printf("compute-type-assigner-webhook version %s built at (%s)\n", BuildVersion, BuildTime)
		os.Exit(0)
	}
	// Setup a Manager
	glog.Info("setting up manager")
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		glog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	// Setup webhooks
	glog.Info("setting up webhook server")
	hookServer := mgr.GetWebhookServer()

	hookServer.Port = params.port
	hookServer.CertDir = params.certDir
	handler := &ValidationHandler{}

	glog.Info("registering webhooks to the webhook server")
	hookServer.Register("/validate", &webhook.Admission{Handler: handler})

	glog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		glog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}
