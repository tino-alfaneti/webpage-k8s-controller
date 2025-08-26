package manager

import (
	"context"
	webpagecontroller "github/tino-alfaneti/webpagecontroller/controller/api/v1"

	"github.com/go-logr/logr"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// define reconciler
type WebReconciler struct {
 client.Client
 scheme *runtime.Scheme
 log    logr.Logger
}

var (
 scheme = runtime.NewScheme()
)

// initiate the program by creating the scheme
func init() {
 utilruntime.Must(webpagecontroller.AddToScheme(scheme))
}

func initLogger() logr.Logger {
    logger := zap.New()
    ctrl.SetLogger(logger)
    log := ctrl.Log.WithName("main")
    log.Info("set up manager")
    return log

}

func Start() {
  
  log := initLogger()

  // create manager
  mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
    Scheme: scheme,
  })

  if err != nil {
      log.Error(err, "unable to create manager")
      os.Exit(1)
  }

  wr := WebReconciler{
      Client: mgr.GetClient(),
      scheme: mgr.GetScheme(),
      log:    log.WithName("web-reconciler"),
  }

  err = wr.SetupWithManager(mgr)
  if err != nil {
      log.Error(err, "unable to create controller")
      os.Exit(1)
  }


  ctx := ctrl.SetupSignalHandler()
  if err = mgr.Start(ctx); err != nil {
      log.Error(err, "problem running manager")
      os.Exit(1)
    }
    
  }

  func (wr *WebReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
      For(&webpagecontroller.WebPage{}).
      Complete(wr)
  }

  func (wr *WebReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    wr.log.Info("reconciling webpage", "name", req.Name, "namespace", req.Namespace)

    wp := webpagecontroller.WebPage{}

    err := wr.Client.Get(ctx, req.NamespacedName, &wp)
    if err != nil && kerr.IsNotFound(err) {
      return ctrl.Result{}, nil
    } else if err != nil {
      return ctrl.Result{}, err
    }

    wr.log.Info("webpage reconciled")

    return ctrl.Result{}, nil
}