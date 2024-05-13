package hostname

import (
	"context"
	metal3v1 "github.com/metal3-io/cluster-api-provider-metal3/api/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	capv1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	ipamv1 "sigs.k8s.io/cluster-api/exp/ipam/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "hostname")
}

var _ = Describe("determining hostnames", func() {
	testScheme := runtime.NewScheme()
	Expect(ipamv1.AddToScheme(testScheme)).To(Succeed())
	Expect(metal3v1.AddToScheme(testScheme)).To(Succeed())
	Expect(capv1.AddToScheme(testScheme)).To(Succeed())

	Context("metal3", func() {
		When("the owner chain can be resolved", func() {
			var cl client.Client
			var claim ipamv1.IPAddressClaim
			BeforeEach(func() {
				cl = fake.NewClientBuilder().
					WithScheme(testScheme).
					WithObjects(
						&metal3v1.Metal3Data{
							ObjectMeta: metav1.ObjectMeta{
								Name: "data",
								OwnerReferences: []metav1.OwnerReference{
									{
										Name:       "machine",
										Kind:       "Metal3Machine",
										APIVersion: metal3v1.GroupVersion.String(),
									},
								},
							},
						},
						&metal3v1.Metal3Machine{
							ObjectMeta: metav1.ObjectMeta{
								Name: "machine",
								OwnerReferences: []metav1.OwnerReference{
									{
										Name:       "capimachine",
										Kind:       "Machine",
										APIVersion: "cluster.x-k8s.io/v1beta1",
									},
								},
							},
						},
					).
					Build()
				claim = ipamv1.IPAddressClaim{
					ObjectMeta: metav1.ObjectMeta{
						OwnerReferences: []metav1.OwnerReference{
							{
								Name:       "data",
								Kind:       "Metal3Data",
								APIVersion: metal3v1.GroupVersion.String(),
							},
						},
					},
				}
			})
			It("the name of the capi Machine is used as the hostname", func() {
				r := OwnerChainResolver{Client: cl, Chain: []metav1.GroupKind{
					{Group: "infrastructure.cluster.x-k8s.io", Kind: "Metal3Data"},
					{Group: "infrastructure.cluster.x-k8s.io", Kind: "Metal3Machine"},
					{Group: "cluster.x-k8s.io", Kind: "Machine"},
				}}
				Expect(r.GetHostname(context.Background(), &claim)).To(Equal("capimachine"))
			})
		})
	})
	Context("vsphere", func() {
		When("the owner chain can be resolved", func() {
			var cl client.Client
			var claim ipamv1.IPAddressClaim
			BeforeEach(func() {
				cl = fake.NewClientBuilder().
					WithScheme(testScheme).
					WithObjects(
						&capv1.VSphereVM{
							ObjectMeta: metav1.ObjectMeta{
								Name: "vm",
								OwnerReferences: []metav1.OwnerReference{
									{
										Name:       "machine",
										Kind:       "VSphereMachine",
										APIVersion: capv1.GroupVersion.String(),
									},
								},
							},
						},
						&capv1.VSphereMachine{
							ObjectMeta: metav1.ObjectMeta{
								Name: "machine",
								OwnerReferences: []metav1.OwnerReference{
									{
										Name:       "capimachine",
										Kind:       "Machine",
										APIVersion: "cluster.x-k8s.io/v1beta1",
									},
								},
							},
						},
					).
					Build()
				claim = ipamv1.IPAddressClaim{
					ObjectMeta: metav1.ObjectMeta{
						OwnerReferences: []metav1.OwnerReference{
							{
								Name:       "vm",
								Kind:       "VSphereVM",
								APIVersion: capv1.GroupVersion.String(),
							},
						},
					},
				}
			})
			It("the name of the capi Machine is used as the hostname", func() {
				r := OwnerChainResolver{Client: cl, Chain: []metav1.GroupKind{
					{Group: "infrastructure.cluster.x-k8s.io", Kind: "VSphereVM"},
					{Group: "infrastructure.cluster.x-k8s.io", Kind: "VSphereMachine"},
					{Group: "cluster.x-k8s.io", Kind: "Machine"},
				}}
				Expect(r.GetHostname(context.Background(), &claim)).To(Equal("capimachine"))
			})
		})
	})
})
