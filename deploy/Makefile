.PHONY: kustomize-apply-latest kustomize-manifest-gen

ENV?=localdev

kustomize-apply-latest:
	kubectl kustomize ./$(ENV) |  kubectl apply -f -
	kubectl -n reference rollout restart deployment golang-reference-api

kustomize-manifest-gen:
	kubectl kustomize ./$(ENV)
