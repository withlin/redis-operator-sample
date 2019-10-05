package redisservice

import (
	"context"
	"fmt"

	"github.com/gomodule/redigo/redis"
	redisv1alpha1 "github.com/ym/redis-operator/pkg/apis/redis/v1alpha1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_redisservice")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new RedisService Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileRedisService{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("redisservice-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource RedisService
	err = c.Watch(&source.Kind{Type: &redisv1alpha1.RedisService{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner RedisService
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &redisv1alpha1.RedisService{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileRedisService implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileRedisService{}

// ReconcileRedisService reconciles a RedisService object
type ReconcileRedisService struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a RedisService object and makes changes based on the state read
// and what is in the RedisService.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileRedisService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling RedisService")

	// Fetch the RedisService instance
	instance := &redisv1alpha1.RedisService{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	//// Define a new Pod object
	//pod := newPodForCR(instance)
	//
	//// Set RedisService instance as the owner and controller
	//if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
	//	return reconcile.Result{}, err
	//}
	//
	//// Check if this Pod already exists
	//found := &corev1.Service{}
	//err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	//if err != nil && errors.IsNotFound(err) {
	//	reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
	//	err = r.client.Create(context.TODO(), pod)
	//	if err != nil {
	//		return reconcile.Result{}, err
	//	}
	//
	//	// Pod created successfully - don't requeue
	//	return reconcile.Result{}, nil
	//} else if err != nil {
	//	return reconcile.Result{}, err
	//}

	// Define a new Service object
	svc := newServiceForCR(instance)

	// Set RedisService instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, svc, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Service already exists
	foundSvc := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: svc.Name, Namespace: svc.Namespace}, foundSvc)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		err = r.client.Create(context.TODO(), svc)
		if err != nil {
			reqLogger.Info("----------------------------->1")
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		reqLogger.Info("----------------------------->2")
		return reconcile.Result{}, nil
	} else if err != nil {
		reqLogger.Info("----------------------------->3")
		return reconcile.Result{}, err
	}
	reqLogger.Info("----------------------------->4")
	stat := newStatefulSetForCR(instance)
	foudState := &appsv1beta1.StatefulSet{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: stat.Name, Namespace: stat.Namespace}, foudState)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new StatefulSet", "StatefulSet.Namespace", stat.Namespace, "StatefulSet.Name", stat.Name)
		err = r.client.Create(context.TODO(), stat)
		if err != nil {
			reqLogger.Info("----------------------------->5")
			return reconcile.Result{}, err
		}
		reqLogger.Info("----------------------------->6")
		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		reqLogger.Info("----------------------------->7")
		return reconcile.Result{}, err
	}
	if err := redisConnect(); err != nil {
		reqLogger.Info("----------------------------->8")
		return reconcile.Result{}, err
	}
	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Service already exists", "Service.Namespace", foundSvc.Namespace, "Service.Name", foundSvc.Name)
	return reconcile.Result{}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *redisv1alpha1.RedisService) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}

//缺少OwnerReference，这个有级联删除的作用
func newServiceForCR(cr *redisv1alpha1.RedisService) *corev1.Service {
	svc := &corev1.Service{}
	svcLab := make(map[string]string)
	svcLab["app"] = "redis"
	svc.ObjectMeta = metav1.ObjectMeta{
		Name:      "redis",
		Labels:    svcLab,
		Namespace: cr.Namespace,
	}
	svc.Spec = corev1.ServiceSpec{
		Ports: []corev1.ServicePort{corev1.ServicePort{Name: "redis",
			Port: 6379}, corev1.ServicePort{Name: "redis-cluster",
			Port: 16379},
		},
		Selector:  svcLab,
		Type:      corev1.ClusterIPNone,
		ClusterIP: "None",
	}
	return svc
}

//缺点:
//1.没有指定redis的密码.
//2.指定volume
//3.没有指定size
//4.没有用OwnerReferences
//5.没有自定义Pod没有指定相关的Affinity，Tolerations，SecurityContext
func newStatefulSetForCR(cr *redisv1alpha1.RedisService) *appsv1beta1.StatefulSet {
	stat := &appsv1beta1.StatefulSet{
		TypeMeta:   v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{Name: "redis", Namespace: cr.Namespace},
		Spec:       appsv1beta1.StatefulSetSpec{},
		Status:     appsv1beta1.StatefulSetStatus{},
	}
	var re int32 = 2
	lab := make(map[string]string)
	lab["app"] = "redis"
	spec := appsv1beta1.StatefulSetSpec{
		Replicas:    &re,
		ServiceName: "redis",
		Selector: &metav1.LabelSelector{
			MatchLabels: lab,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: lab,
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{corev1.Container{Name: "redis", Image: "redis:latest", Ports: []corev1.ContainerPort{corev1.ContainerPort{ContainerPort: 6379, Name: "redis"}, {ContainerPort: 16379, Name: "redis-cluster"}}}},
			},
		},
	}
	stat.Spec = spec
	return stat
}

func redisConnect() error {
	c, err := redis.Dial("tcp", "redis-0.redis:6379")
	if err != nil {
		fmt.Println("--------------->Connect to redis error", err)
		return err
	}
	res, err := c.Do("slaveof", "redis-1.redis", "6379")
	if err != nil {
		fmt.Println("-------------------->slaveof command is excute err:", err)
	}
	fmt.Println("------------------------>res is ", res)
	defer c.Close()
	return nil
}
