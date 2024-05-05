repo=ghcr.io/pdpresson
version=0.0.1

services=(gitops_client)

for i in "${services[@]}"
do
    echo "Building image for $i service"
    cd ./$i/
    docker build -t $repo/apps/$i:$version -f ./.docker/dockerfile .
    cd ..
done