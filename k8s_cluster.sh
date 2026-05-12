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



echo "Checking dependencies..."

for cmd in kubectl docker kind; do
    if ! command -v "$cmd" >/dev/null 2>&1; then
        echo "❌ $cmd is not installed or not in PATH"
        exit 1
    fi
done

echo "✅ All dependencies are installed"

if [ "$ACTION" = "start" ]; then


    echo "Namespace: $NAMESPACE"

    echo "=================================================="
    echo "Build docker images"

    docker build -t api-gateway:0.0.1 api-gateway
    docker build -t download-service:0.0.1 download-service
    docker build -t pdf-converter-service:0.0.1 pdf-converter-service

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