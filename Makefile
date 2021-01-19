# Make sure we pick up any local overrides.
-include .makerc

.PHONY: deploy
deploy:
	$(MAKE) -C signer-ca deploy-e2e
	$(MAKE) deploy-istio
	$(MAKE) deploy-apps

.PHONY: deploy-apps
deploy-apps:
	kubectl label namespace default istio-injection=enabled
	kubectl apply -n default -f apps/
	kubectl apply -f mtls.yaml

.PHONY: deploy-istio
deploy-istio:
	kubectl create namespace istio-system
	kubectl create secret generic external-ca-cert -n istio-system --from-file=root-cert.pem=signer-ca/$(E2E_PKI)/tls.crt
	istioctl install -f istio.yaml -y

.PHONY: uninstall
uninstall:
	kubectl delete -f apps/ --ignore-not-found
	kubectl label namespace default istio-injection-
	kubectl delete -f mtls.yaml --ignore-not-found
	kubectl delete -f authz.yaml --ignore-not-found
	kustomize build signer-ca/${E2E_PKI} | kubectl delete --ignore-not-found -f -
	istioctl x uninstall --purge -y
	kubectl delete namespace istio-system --ignore-not-found
