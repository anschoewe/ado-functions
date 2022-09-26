# Overview
This project consumes webhook requests from Azure DevOps.   The code runs inside the Azure Functions framework
[Here](https://learn.microsoft.com/en-us/azure/azure-functions/functions-custom-handlers) is a good technical description of how to run Azure Functions with Golang code inside.

# Local development
```
go build handler.go
func start
```

You can then test the function at http://localhost:7071/api/repoevent

# Test locally with Azure DevOps
If you want to test from actual ADO events, you'll need to expose your locally run Azure function through ngrok.

```
# install ngrok
brew install ngrok/ngrok/ngrok
ngrok config add-authtoken <my-token-from-ngrok-website>

# run ngrok on the same port as Azure Functions is listening, 7071, locally
func start
# in a new Terminal tab...
ngrok http 7071
```

# Deploy to Azure Functions in Cloud
First, build the app for the OS and architecture in Azure App Service. It needs to run in an Linux App Service.

```
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" handler.go
```

Then in Visual Studio Code, upload your local Azure Function workspace to the cloud. It will overwrite whatever you have deployed there.
