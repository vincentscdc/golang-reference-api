.PHONY: kustomize-apply-latest kustomize-manifest-gen

ENV?=dev

kustomize-apply-latest:
	kubectl kustomize ./$(ENV) |  kubectl apply -f -
	kubectl -n fiat rollout restart deployment golang-reference-api

kustomize-manifest-gen:
	kubectl kustomize ./$(ENV)
