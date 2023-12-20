package controller

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	webgamev1 "github.com/webgamedevelop/webgame/api/v1"
)

var _ = Describe("Test Webgame controller", func() {
	const (
		timeout             = time.Second * 10
		interval            = time.Millisecond * 250
		namespace           = "webgames"
		webgameInstanceName = "webgame-sample"
	)

	BeforeEach(func() {
		var ns corev1.Namespace
		ns.SetName(namespace)
		_, err := controllerutil.CreateOrUpdate(ctx, k8sClient, &ns, func() error { return nil })
		Expect(err).Should(Succeed())
	})

	Context("webgame controller test", func() {
		It("create webgame sample", func() {
			var err error
			var replicas int32 = 1
			var webgame webgamev1.WebGame
			webgame.SetNamespace(namespace)
			webgame.SetName(webgameInstanceName)
			mutate := func() error {
				webgame.Spec.DisplayName = "test-webgame-instance"
				webgame.Spec.GameType = "2048"
				webgame.Spec.IngressClass = "nginx"
				webgame.Spec.Domain = "localhost"
				webgame.Spec.IndexPage = "index.html"
				webgame.Spec.ServerPort = intstr.IntOrString{Type: intstr.Int, IntVal: int32(80)}
				webgame.Spec.Image = "webgamedevelop/2048:latest"
				webgame.Spec.Replicas = &replicas
				webgame.Spec.ImagePullSecrets = []corev1.LocalObjectReference{{Name: "test-image-pull-secret"}}
				return nil
			}

			// create webgame instance
			_, err = controllerutil.CreateOrUpdate(ctx, k8sClient, &webgame, mutate)
			Expect(err).Should(Succeed())

			// get webgame instance
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, ctrlclient.ObjectKeyFromObject(&webgame), &webgame); err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			// get deployment
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, ctrlclient.ObjectKeyFromObject(&webgame), &appsv1.Deployment{}); err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			// get svc
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, ctrlclient.ObjectKeyFromObject(&webgame), &corev1.Service{}); err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			// get ing
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, ctrlclient.ObjectKeyFromObject(&webgame), &networkingv1.Ingress{}); err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
		})

		It("delete webgame instance", func() {
			var err error
			var webgame webgamev1.WebGame
			webgame.SetNamespace(namespace)
			webgame.SetName(webgameInstanceName)
			err = k8sClient.Delete(ctx, &webgame)
			Expect(err).Should(Succeed())

			// get webgame, expect not found
			// get webgame instance
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, ctrlclient.ObjectKeyFromObject(&webgame), &webgame); err != nil && apierrors.IsNotFound(err) {
					return true
				}
				return false
			}, timeout, interval).Should(BeTrue())
		})
	})
})
