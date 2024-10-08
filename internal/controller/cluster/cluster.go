package cluster

import (
	"context"

	resourceClient "github.com/zncdatadev/operator-go/pkg/client"
	"github.com/zncdatadev/operator-go/pkg/reconciler"
	"github.com/zncdatadev/operator-go/pkg/util"
	supersetv1alpha1 "github.com/zncdatadev/superset-operator/api/v1alpha1"
	"github.com/zncdatadev/superset-operator/internal/controller/node"
	corev1 "k8s.io/api/core/v1"
)

var _ reconciler.Reconciler = &Reconciler{}

type Reconciler struct {
	reconciler.BaseCluster[*supersetv1alpha1.SupersetClusterSpec]
	ClusterConfig *supersetv1alpha1.ClusterConfigSpec
}

func NewReconciler(
	client *resourceClient.Client,
	clusterInfo reconciler.ClusterInfo,
	spec *supersetv1alpha1.SupersetClusterSpec,
) *Reconciler {

	return &Reconciler{
		BaseCluster: *reconciler.NewBaseCluster(
			client,
			clusterInfo,
			spec.ClusterOperation,
			spec,
		),
		ClusterConfig: spec.ClusterConfig,
	}

}

func (r *Reconciler) GetImage() *util.Image {
	image := &util.Image{
		Repo:            supersetv1alpha1.DefaultRepository,
		ProductName:     supersetv1alpha1.DefaultProductName,
		PlatformVersion: supersetv1alpha1.DefaultPlatformVersion,
		ProductVersion:  supersetv1alpha1.DefaultProductVersion,
		PullPolicy:      &[]corev1.PullPolicy{corev1.PullIfNotPresent}[0],
	}

	if r.Spec.Image != nil {
		image.Custom = r.Spec.Image.Custom
		image.Repo = r.Spec.Image.Repo
		image.PlatformVersion = r.Spec.Image.PlatformVersion
		image.ProductVersion = r.Spec.Image.ProductVersion
		image.PullPolicy = &r.Spec.Image.PullPolicy
	}
	return image
}

func (r *Reconciler) RegisterResources(ctx context.Context) error {

	node := node.NewReconciler(
		r.Client,
		reconciler.RoleInfo{
			ClusterInfo: r.ClusterInfo,
			RoleName:    "node",
		},
		r.GetClusterOperation(),
		r.ClusterConfig,
		r.GetImage(),
		r.Spec.Node,
	)

	if err := node.RegisterResources(ctx); err != nil {
		return err
	}

	r.AddResource(node)

	return nil

}
