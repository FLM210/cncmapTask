# 监听端口自定义方法

## 设置LISTENPORT环境变量
环境变量优先级最高
    
## 使用配置文件http.conf
使用配置文件需加入配置文件路径参数 如:/etc/http.conf

## 若未指定任何参数则默认使用80端口

# In Kubernetes 

>部署应用   
>
kubectl apply -f module8/deploy/deployment.yaml   
>应用配置文件  
>
kubectl apply -f module8/deploy/deployment.yaml  
>配置ingress  
>
kubectl apply -f module8/deploy/deployment.yaml 
kubectl apply -f module8/deploy/deployment.yaml 
