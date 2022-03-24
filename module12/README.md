#使用istio ingress gateway 暴露服务


##创建https证书的secret
kubectl create -n istio-system secret tls cncamp-credential --key=cncamp.io.key --cert=cncamp.io.crt (证书所在namespace需与istio-ingress pod 一致)
##创建namespace demo-istio
kubectl create namespace demo-istio
##部署服务
kubectl apply -f appDeploy.yaml -n demo-istio
##部署istio ingress
kubectl apply -f istio-ingress.yaml -n demo-istio
##访问https服务
curl --resolve gohttpserver.cncamp.io:443:`kubectl get svc -n istio-system istio-ingressgateway   -o jsonpath='{.spec.clusterIP}'`   https://gohttpserver.cncamp.io/healthz -v -k
