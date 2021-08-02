### Prometheus
```
helm install prom prometheus-community/kube-prometheus-stack -f prometheus.yaml --atomic
helm install nginx ingress-nginx/ingress-nginx -f nginx-ingress.yaml --atomic

kubectl port-forward service/prom-grafana 9000:80
kubectl port-forward service/prom-kube-prometheus-stack-prometheus 9090

helm install postgre-exporter-prom prometheus-community/prometheus-postgres-exporter -f postgres-exporter.yaml --atomic
```

### Install
```
skaffold run
```
or
```
helm install arch-homework ./app-charts -f ./app-charts/values.yaml
```