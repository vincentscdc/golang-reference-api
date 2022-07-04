.PHONY: kustomize-apply-latest kustomize-manifest-gen

ENV?=localdev

kustomize-apply-latest:
	kubectl kustomize ./$(ENV) |  kubectl apply -f -
	kubectl -n reference rollout restart deployment golang-reference-api

kustomize-manifest-gen:
	kubectl kustomize ./$(ENV)

################
# Port Forward #
################

forward-start: ingress/forward-start postgres/forward-start redis/forward-start rabbitmq/forward-start opentelemetry-collector/forward-start jaeger/forward-start

forward-stop:
	pkill -f port-forward

ingress/forward-start:
	kubectl port-forward -n ingress-nginx svc/ingress-nginx-controller \
		8080:80 \
		8443:443 &

postgres/forward-start:
	kubectl port-forward -n ingress-nginx svc/ingress-nginx-controller \
		5432 \
		5433 &

redis/forward-start:
	kubectl port-forward -n ingress-nginx svc/ingress-nginx-controller \
		6379 &

rabbitmq/forward-start:
	kubectl port-forward -n ingress-nginx svc/ingress-nginx-controller \
		5672 &

opentelemetry-collector/forward-start:
	kubectl port-forward -n ingress-nginx svc/ingress-nginx-controller \
		4317 &

jaeger/forward-start:
	kubectl port-forward -n observability svc/jaeger-query \
		16686 &

##########
# Others #
##########

context-select:
	kubectl config use-context $(K8S_CONTEXT)