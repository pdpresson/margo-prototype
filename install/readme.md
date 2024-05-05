# This is setup to use three machines
1. VM for the Workload Orchestration Software
    - Expected to be running a Kubernetes distribution
    - uses wos.local for the hostname
1. VM for the device runnng Workload Orchestration Software Agent
    - Expected to be running a Kubernetes distribution
    - uses edge.local for the hostname
1. VM (or host) running the GOGs code repository services
    - use gogs for the hostname

# Configure the GOGs local repository service
1. This can be run either on the host machine or a VM in the same network as the WOS and device VMs
1. Create the following folders under the ./.docker folder
    - ./gogs/app
    - ./gogs/logs
    - ./gogs/db_data
1. Copy the ./config/gogs/app.ini ./.docker/gogs/app/gogs/conf/app.ini
1. Run the ./.docker/docker-compose.yaml file to start the services
1. Create an admin account. This will be used in the WOS configmap file
1. Create a token for the admin account. This will be used in the WOS configmap file
1. Create a repo to simulate an application vendor (e.g. hello-world-app)
1. Put the ./examples/app-description.yml file in the repo

# Configure WOS VM
1. Add an entry for the gogs server to the /etc/hosts file
   > \<Gogs VM/Host IP address\>  gogs
1. Update the ./WOS/values.yaml file
    - Set _workload-orchestration-service.deviceRepo.username_ to the admin user you created
    - Set _workload-orchestration-service.deviceRepo.username_ to the admin's password you created
    - Set _workload-orchestration-service.deviceRepo.tokenName_ to the token name you created
    - Set _workload-orchestration-service.deviceRepo.token_ to the toke value that was created for you
    - Set _workload-orchestration-service.deviceRepo.hostAliases_ IP address to the IP of your GOGs server
    - Set _gitops-pullservice.gitRepos.hostAliases_ IP address to the IP of your GOGs server
    - Set _gitops-pushservice.deviceRepo.hostAlias_ IP address to the IP of your GOGs server
1. Install the WOS helm chart
    > helm install my-wos ghcr.io/pdpresson//charts/ephemeral-workload-orchestration:0.0.1 --create-namespace -n wos ./WOS/values.yaml

# Configure Device VM
1. Add an entry for the gogs server to the /etc/hosts file
    > \<Gogs VM/Host IP adresss\>  gogs
1. Update the ./WOSA/values.yaml file
    - Set gitops-client.config.deviceRepo.hostAlias IP address to the IP of your GOGs server
1. Install the WOSA helm chart
    > helm install my-wosa ghcr.io/pdpresson/charts/ephemeral-workload-orchestration-agent:0.0.1 --create-namespace -n wosa ./WOSA/values.yaml

# Configure host machine
1. Add entries for the VMs to your hosts file (/etc/hosts or c:\windows\system32\drivers\etc\hosts)
    > \<Gogs VM/Host IP address\>  gogs
    >
    > \<Device VM IP address\>  edge.local
    >
    > \<WOS VM IP address\>  wos.local

# Register the App Vendor's Repo
1. in the browser open http://wos.local/orchestration-portal/apprepos
1. Fill in the URL and branch for the repo 
1. Add the repo
1. Wait for 1-2 minute for the application to be added to the app catalog
1. Check to see if it's been added by clicking the app catalog link at the top

# Upload device description
1. Use curl to upload the device description files to the WOS
    > curl --header "Content-Type: application/json" --request POST --data "$(cat ./examples/deviceDescription.json)" http://wos.local/orchestration-service/devices
1. This will return the URL for the device repo
1. On the device VM update the gitops-client configmap
    - Set _gitops-client.deviceRepo.reportUrl_ to the URL return from the curl command
    - Set _gitops-client.deviceRepo.branch_ to the brach return from the curl command (should be master)
1. Delete the gitops-client pod

# Add the app to the device
1. in the browser open http://wos.local/orchestration-portal/appcatalog
1. Click the button to install the application
1. Fill in the Greeting and Targets
1. Select the device you added in the previous step
1. Click the button to install the application on the device
1. Wait for a couple of minutes
1. Open then hello-world applications by navigationg to http://edge.local/hello 