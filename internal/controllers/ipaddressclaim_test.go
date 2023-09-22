package controllers

// . "github.com/onsi/ginkgo/v2"
// . "github.com/onsi/gomega"
// . "sigs.k8s.io/controller-runtime/pkg/envtest/komega"

// var _ = Describe("IPAddressClaimReconciler", func() {
// 	Context("When a new IPAddressClaim is created", func() {
// 		It("should ignore the claim if it doesn't reference a infoblox pool", func() {
// 			// TODO: figure out how to do that properly
// 		})

// 		When("the referenced pool exists", func() {
// 			const poolName = "test-pool"

// 			BeforeEach(func() {
// 				pool := v1alpha1.InfobloxIPPool{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name:      poolName,
// 						Namespace: "default",
// 					},
// 					Spec: v1alpha1.InfobloxIPPoolSpec{
// 						First:   "10.0.0.1",
// 						Last:    "10.0.0.254",
// 						Prefix:  24,
// 						Gateway: "10.0.0.2",
// 					},
// 				}
// 				Expect(k8sClient.Create(context.Background(), &pool)).To(Succeed())
// 				Eventually(Get(&pool)).Should(Succeed())
// 			})

// 			AfterEach(func() {
// 				pool := v1alpha1.InfobloxIPPool{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name:      poolName,
// 						Namespace: "default",
// 					},
// 				}
// 				Expect(k8sClient.Delete(context.Background(), &pool)).To(Succeed())
// 				Eventually(Get(&pool)).Should(Not(Succeed()))
// 			})

// 			It("should allocate an Address from the Pool", func() {
// 				claim := clusterv1.IPAddressClaim{
// 					TypeMeta: metav1.TypeMeta{
// 						APIVersion: "ipam.cluster.x-k8s.io/v1alpha1",
// 						Kind:       "IPAddressClaim",
// 					},
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name:      "test",
// 						Namespace: "default",
// 					},
// 					Spec: clusterv1.IPAddressClaimSpec{
// 						PoolRef: corev1.TypedLocalObjectReference{
// 							APIGroup: pointer.String("ipam.cluster.x-k8s.io"),
// 							Kind:     "InfobloxIPPool",
// 							Name:     poolName,
// 						},
// 					},
// 				}
// 				address := clusterv1.IPAddress{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name:      "test",
// 						Namespace: "default",
// 					},
// 					Spec: clusterv1.IPAddressSpec{},
// 				}

// 				desired := clusterv1.IPAddress{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name:       "test",
// 						Namespace:  "default",
// 						Finalizers: []string{ProtectAddressFinalizer},
// 						OwnerReferences: []metav1.OwnerReference{
// 							{
// 								APIVersion:         "ipam.cluster.x-k8s.io/v1alpha1",
// 								BlockOwnerDeletion: pointer.Bool(true),
// 								Controller:         pointer.Bool(true),
// 								Kind:               "IPAddressClaim",
// 								Name:               "test",
// 							},
// 							{
// 								APIVersion:         "ipam.cluster.x-k8s.io/v1alpha1",
// 								BlockOwnerDeletion: pointer.Bool(true),
// 								Controller:         pointer.Bool(false),
// 								Kind:               "InfobloxIPPool",
// 								Name:               "test-pool",
// 							},
// 						},
// 					},
// 					Spec: clusterv1.IPAddressSpec{
// 						ClaimRef: corev1.LocalObjectReference{
// 							Name: "test",
// 						},
// 						PoolRef: corev1.TypedLocalObjectReference{
// 							APIGroup: pointer.String("ipam.cluster.x-k8s.io"),
// 							Kind:     "InfobloxIPPool",
// 							Name:     poolName,
// 						},
// 						Address: "10.0.0.1",
// 						Prefix:  24,
// 						Gateway: "10.0.0.2",
// 					},
// 				}

// 				Expect(k8sClient.Create(context.Background(), &claim)).To(Succeed())
// 				// Eventually(Object(&claim)).Should(HaveField("Status.Address.Name", Equal(claim.ObjectMeta.Name)))

// 				Eventually(Object(&address)).WithTimeout(time.Second).WithPolling(100 * time.Millisecond).Should(And(
// 					EqualObject(&desired, IgnoreAutogeneratedMetadata, IgnorePaths{
// 						"TypeMeta",
// 						"ObjectMeta.OwnerReferences[0].UID",
// 						"ObjectMeta.OwnerReferences[1].UID",
// 						"Spec.Claim.UID",
// 						"Spec.Pool.UID",
// 					}),
// 				))
// 			})
// 		})

// 		When("the referenced pool uses single ip addresses", func() {
// 			const poolName = "test-pool-single-ip-addresses"
// 			const namespace = "test-single-ip-addresses"
// 			BeforeEach(func() {
// 				nsName := corev1.Namespace{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name: namespace,
// 					},
// 				}
// 				Expect(k8sClient.Create(context.Background(), &nsName)).To(Succeed())

// 				pool := v1alpha1.InfobloxIPPool{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name:      poolName,
// 						Namespace: namespace,
// 					},
// 					Spec: v1alpha1.InfobloxIPPoolSpec{
// 						Addresses: []string{
// 							"10.0.0.50",
// 							"10.0.0.128",
// 						},
// 						Prefix:  24,
// 						Gateway: "10.0.0.1",
// 					},
// 				}
// 				Expect(k8sClient.Create(context.Background(), &pool)).To(Succeed())
// 				Eventually(Get(&pool)).Should(Succeed())
// 			})

// 			AfterEach(func() {
// 				pool := v1alpha1.InfobloxIPPool{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name:      poolName,
// 						Namespace: namespace,
// 					},
// 				}
// 				Expect(k8sClient.Delete(context.Background(), &pool)).To(Succeed())
// 				Eventually(Get(&pool)).Should(Not(Succeed()))
// 			})

// 			It("should allocate an Address from the Pool", func() {
// 				claim1 := clusterv1.IPAddressClaim{
// 					TypeMeta: metav1.TypeMeta{
// 						APIVersion: "ipam.cluster.x-k8s.io/v1alpha1",
// 						Kind:       "IPAddressClaim",
// 					},
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name:      "test-1",
// 						Namespace: namespace,
// 					},
// 					Spec: clusterv1.IPAddressClaimSpec{
// 						PoolRef: corev1.TypedLocalObjectReference{
// 							APIGroup: pointer.String("ipam.cluster.x-k8s.io"),
// 							Kind:     "InfobloxIPPool",
// 							Name:     poolName,
// 						},
// 					},
// 				}

// 				claim2 := clusterv1.IPAddressClaim{
// 					TypeMeta: metav1.TypeMeta{
// 						APIVersion: "ipam.cluster.x-k8s.io/v1alpha1",
// 						Kind:       "IPAddressClaim",
// 					},
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name:      "test-2",
// 						Namespace: namespace,
// 					},
// 					Spec: clusterv1.IPAddressClaimSpec{
// 						PoolRef: corev1.TypedLocalObjectReference{
// 							APIGroup: pointer.String("ipam.cluster.x-k8s.io"),
// 							Kind:     "InfobloxIPPool",
// 							Name:     poolName,
// 						},
// 					},
// 				}

// 				expectedAddress1 := clusterv1.IPAddress{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name:       "test-1",
// 						Namespace:  namespace,
// 						Finalizers: []string{ProtectAddressFinalizer},
// 						OwnerReferences: []metav1.OwnerReference{
// 							{
// 								APIVersion:         "ipam.cluster.x-k8s.io/v1alpha1",
// 								BlockOwnerDeletion: pointer.Bool(true),
// 								Controller:         pointer.Bool(true),
// 								Kind:               "IPAddressClaim",
// 								Name:               "test-1",
// 							},
// 							{
// 								APIVersion:         "ipam.cluster.x-k8s.io/v1alpha1",
// 								BlockOwnerDeletion: pointer.Bool(true),
// 								Controller:         pointer.Bool(false),
// 								Kind:               "InfobloxIPPool",
// 								Name:               "test-pool-single-ip-addresses",
// 							},
// 						},
// 					},
// 					Spec: clusterv1.IPAddressSpec{
// 						ClaimRef: corev1.LocalObjectReference{
// 							Name: "test-1",
// 						},
// 						PoolRef: corev1.TypedLocalObjectReference{
// 							APIGroup: pointer.String("ipam.cluster.x-k8s.io"),
// 							Kind:     "InfobloxIPPool",
// 							Name:     poolName,
// 						},
// 						Address: "10.0.0.50",
// 						Prefix:  24,
// 						Gateway: "10.0.0.1",
// 					},
// 				}

// 				expectedAddress2 := clusterv1.IPAddress{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name:       "test-2",
// 						Namespace:  namespace,
// 						Finalizers: []string{ProtectAddressFinalizer},
// 						OwnerReferences: []metav1.OwnerReference{
// 							{
// 								APIVersion:         "ipam.cluster.x-k8s.io/v1alpha1",
// 								BlockOwnerDeletion: pointer.Bool(true),
// 								Controller:         pointer.Bool(true),
// 								Kind:               "IPAddressClaim",
// 								Name:               "test-2",
// 							},
// 							{
// 								APIVersion:         "ipam.cluster.x-k8s.io/v1alpha1",
// 								BlockOwnerDeletion: pointer.Bool(true),
// 								Controller:         pointer.Bool(false),
// 								Kind:               "InfobloxIPPool",
// 								Name:               "test-pool-single-ip-addresses",
// 							},
// 						},
// 					},
// 					Spec: clusterv1.IPAddressSpec{
// 						ClaimRef: corev1.LocalObjectReference{
// 							Name: "test-2",
// 						},
// 						PoolRef: corev1.TypedLocalObjectReference{
// 							APIGroup: pointer.String("ipam.cluster.x-k8s.io"),
// 							Kind:     "InfobloxIPPool",
// 							Name:     poolName,
// 						},
// 						Address: "10.0.0.128",
// 						Prefix:  24,
// 						Gateway: "10.0.0.1",
// 					},
// 				}

// 				Expect(k8sClient.Create(context.Background(), &claim1)).To(Succeed())
// 				Expect(k8sClient.Create(context.Background(), &claim2)).To(Succeed())
// 				// Eventually(Object(&claim)).Should(HaveField("Status.Address.Name", Equal(claim.ObjectMeta.Name)))

// 				Eventually(Object(&clusterv1.IPAddress{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name:      expectedAddress1.GetName(),
// 						Namespace: namespace,
// 					},
// 				})).WithTimeout(time.Second).WithPolling(100 * time.Millisecond).Should(
// 					EqualObject(&expectedAddress1, IgnoreAutogeneratedMetadata, IgnorePaths{
// 						"TypeMeta",
// 						"ObjectMeta.OwnerReferences[0].UID",
// 						"ObjectMeta.OwnerReferences[1].UID",
// 						"Spec.Claim.UID",
// 						"Spec.Pool.UID",
// 					}),
// 				)

// 				Eventually(Object(&clusterv1.IPAddress{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name:      expectedAddress2.GetName(),
// 						Namespace: namespace,
// 					},
// 				})).WithTimeout(time.Second).WithPolling(100 * time.Millisecond).Should(
// 					EqualObject(&expectedAddress2, IgnoreAutogeneratedMetadata, IgnorePaths{
// 						"TypeMeta",
// 						"ObjectMeta.OwnerReferences[0].UID",
// 						"ObjectMeta.OwnerReferences[1].UID",
// 						"Spec.Claim.UID",
// 						"Spec.Pool.UID",
// 					}),
// 				)
// 			})
// 		})

// 	})
// })
