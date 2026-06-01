#! /bin/sh


set -euo pipefail

ACTION=""
NAMESPACE="api-media"
DELETE_IMAGES=false
CLUSTER_NAME="default"

while [ $# -gt 0 ]; do
    case "$1" in

        start|delete)
            ACTION="$1"
            shift
            ;;

        --namespace|-n)
            NAMESPACE="$2"
            shift 2
            ;;

        --cluster|-c)
            CLUSTER_NAME="$2"
            shift 2
            ;;

        --no-images)
            DELETE_IMAGES=true
            shift
            ;;

        *)
            echo "Unknown argument: $1"
            exit 1
            ;;
    esac
done


echo "=================================================="
echo "Checking dependencies..."
kubectl_version="1.30.0"

version=$(kubectl version --client 2>/dev/null | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' | head -n1 | sed 's/v//')

if [ -z "$version" ]; then
    echo "kubectl not found or version not readable"
    exit 1
fi

if [ "$(printf '%s\n' "$kubectl_version" "$version" | sort -V | head -n1)" != "$kubectl_version" ]; then
    echo "kubectl version $version is lower than required $kubectl_version"
    exit 1
fi

echo "✅ kubectl version OK ($version >= $kubectl_version)"
echo "=================================================="
kind_version="0.20.0"

version=$(kind version 2>/dev/null | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' | head -n1 | sed 's/v//')

if [ -z "$version" ]; then
    echo "kind not installed or version not readable"
    exit 1
fi

     
if [ "$(printf '%s\n' "$kind_version" "$version" | sort -V | head -n1)" != "$kind_version" ]; then
    echo "kind version $version is lower than required $kind_version"
    exit 1
fi

echo "✅ kind version OK ($version >= $kind_version)"

echo "=================================================="
docker_version="27.0.0"

version=$(docker version --format '{{.Client.Version}}' 2>/dev/null)

if [ -z "$version" ]; then
    echo "docker not installed or version not readable"
    exit 1
fi

if [ "$(printf '%s\n' "$docker_version" "$version" | sort -V | head -n1)" != "$docker_version" ]; then
    echo "docker version $version is lower than required $docker_version"
    exit 1
fi

echo "✅ docker version OK ($version >= $docker_version)"

echo "=================================================="

if [ "$ACTION" = "start" ]; then


    echo "Namespace: $NAMESPACE"

    echo "=================================================="
    echo "Build docker images"
    cd services

    docker build -t api-gateway:0.0.1 api-gateway
    docker build -t download-service:0.0.1 download-service
    docker build -t pdf-converter-service:0.0.1 pdf-converter-service

    cd ..  
    echo "=================================================="
    echo "Create k8s cluster"
    kind create cluster --name "$CLUSTER_NAME" --config iac/kind.yaml

    echo "=================================================="
    echo "Create k8s cluster"

    kind load docker-image api-gateway:0.0.1 --name "$CLUSTER_NAME"
    kind load docker-image download-service:0.0.1 --name "$CLUSTER_NAME"
    kind load docker-image pdf-converter-service:0.0.1 --name "$CLUSTER_NAME"

    echo "=================================================="
    echo "Create namespace"

    kubectl create namespace "$NAMESPACE" || true

    echo "=================================================="
    echo "Deploy applications"

    kubectl apply -f iac/k8s/pvc.yaml -n "$NAMESPACE"
    kubectl apply -f iac/k8s/api-gateway.yaml -n "$NAMESPACE"
    kubectl apply -f iac/k8s/pdf-converter-service.yaml -n "$NAMESPACE"
    kubectl apply -f iac/k8s/download-service.yaml -n "$NAMESPACE"


    echo "=================================================="
    echo "🚀 Cluster created successfully"
    echo "🌐 API Gateway:"
    echo "   http://localhost:8080"
    echo "=================================================="
elif [ "$ACTION" = "delete" ]; then
    echo "Deleting cluster..."
        
    kind delete cluster --name "$CLUSTER_NAME"

    if [ "$DELETE_IMAGES" = true ]; then    

        echo "=================================================="
        echo "Deleting docker images"

        docker rmi api-gateway:0.0.1 || true
        docker rmi download-service:0.0.1 || true
        docker rmi pdf-converter-service:0.0.1 || true
    fi

else

    echo "Usage:"
    echo "  $0 start [--namespace <name>]"
    echo "  $0 delete [--namespace <name>] [--no-images]"
    exit 1

fi