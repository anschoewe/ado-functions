# Overview
This project consumes webhook requests from Azure DevOps.   The code runs inside the Azure Functions framework.
[Here](https://learn.microsoft.com/en-us/azure/azure-functions/functions-custom-handlers) is a good technical description of how to run Azure Functions with Golang code inside.
<br/>
Of note, the instructions for creating the Azure Function App, in Azure, are flawed.  By default, the Visual Studio Code extension for Azure Fucntions will create an Azure Function where the OS of the underlying Azure App Service is Windows.  But because we're running a Custom Handler in an Azure Function with Golang code, we need to use an App Service with a Linux OS.  You do this by running the Command Pallet with this command `Azure Functions: Create function app in Azure...(Advanced)`.  You can find more complete instructions [here](https://learn.microsoft.com/en-us/azure/azure-functions/functions-develop-vs-code?tabs=csharp#enable-publishing-with-advanced-create-options).

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
Now you can enter the ngrok URL into the Azure DevOps Webhooks. For example, `https://777a-75-166-52-40.ngrok.io/api/repoevent?eventType=pushed`.  When Azure DevOps invokes the webhook, ngrok will listen on that end point and forward the request to `http://localhost:7071/api/repoevent?eventType=pushed`

# Deploy to Azure Functions in Cloud
First, build the app for the OS and architecture in Azure App Service. It needs to run in an Linux App Service.

```
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" handler.go
```

Then in Visual Studio Code, upload your local Azure Function workspace to the cloud. It will overwrite whatever you have deployed there.
