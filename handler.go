package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type ReqData struct {
	Req struct {
		URL     string            `json:"Url"`
		Method  string            `json:"Method"`
		Query   map[string]string `json:"Query"`
		Headers struct {
			ContentType []string `json:"Content-Type"`
		} `json:"Headers"`
		Params struct {
		} `json:"Params"`
		Body json.RawMessage `json:"Body"`
	} `json:"req"`
}

type InvokeRequest struct {
	Data     ReqData
	Metadata map[string]interface{}
}

type InvokeResponse struct {
	Outputs     map[string]interface{}
	Logs        []string
	ReturnValue interface{}
}

//https://777a-75-166-52-40.ngrok.io/api/repoevent?eventType=pushed
//https://777a-75-166-52-40.ngrok.io/api/repoevent?eventType=pullRequestCreated
//https://777a-75-166-52-40.ngrok.io/api/repoevent?eventType=pullRequestUpdated

func main() {
	customHandlerPort, exists := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if !exists {
		customHandlerPort = "8080"
	} else {
		fmt.Println("FUNCTIONS_CUSTOMHANDLER_PORT exists")
	}

	// http.HandleFunc("/", handler)
	http.HandleFunc("/HttpBranchCreatedTrigger", helloHandler)
	fmt.Println("Go server Listening on: ", customHandlerPort)
	log.Fatal(http.ListenAndServe(":"+customHandlerPort, nil))
}

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
// }

func helloHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("***************Hello World***********")

	var invokeRequest InvokeRequest
	defer req.Body.Close()

	d := json.NewDecoder(req.Body)
	d.Decode(&invokeRequest)
	log.Printf("%s\n", invokeRequest.Data.Req.Headers)
	log.Printf("%s\n", invokeRequest.Data.Req.Method)
	log.Printf("%s\n", invokeRequest.Data.Req.Params)
	log.Printf("%s\n", invokeRequest.Data.Req.Query)
	log.Printf("%s\n", invokeRequest.Data.Req.URL)

	message := "This HTTP triggered function executed successfully. Pass a name in the query string for a personalized response.\n"
	eventType := invokeRequest.Data.Req.Query["eventType"]

	switch eventType {
	case "push":
		pushEventHandler(w, invokeRequest.Data.Req.Body)
	case "pullRequestCreated":
		pullRequestCreatedEventHandler(w, req)
	case "pullRequestUpdated":
		pullRequestUpdatedEventHandler(w, req)
	default:
		unknownEventHandler(w, req)
		return
	}

	fmt.Fprint(w, message)
}

func unknownEventHandler(w http.ResponseWriter, req *http.Request) {
	eventType := req.URL.Query().Get("eventType")
	msg := fmt.Sprintf("Unknown event type: %s", eventType)
	log.Println(msg)
	fmt.Fprint(w, msg)

}

func pullRequestUpdatedEventHandler(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	// log.Printf("Got request body...   %s\n", JsonPrettyPrint(string(reqBody)))

	var pullRequestUpdatedEvent PullRequestUpdated
	err = json.Unmarshal(reqBody, &pullRequestUpdatedEvent)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Got request body...   %v\n", pullRequestUpdatedEvent)
}

func pullRequestCreatedEventHandler(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	// log.Printf("Got request body...   %s\n", JsonPrettyPrint(string(reqBody)))

	var pullRequestCreatedEvent PullRequestCreated
	err = json.Unmarshal(reqBody, &pullRequestCreatedEvent)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Got request body...   %v\n", pullRequestCreatedEvent)
}

func pushEventHandler(w http.ResponseWriter, body json.RawMessage) {
	// defer req.Body.Close()
	// reqBody, err := io.ReadAll(body)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Printf("Got request body...   %s\n", JsonPrettyPrint(string(reqBody)))

	// var codeEvent CodePushedEvent
	// err = json.Unmarshal(reqBody, &codeEvent)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Printf("Got request body...   %v\n", codeEvent)

	// write to storage account queue.
	// 'msg' is the name of the queue's output binding
	// 'res' is the name of the http's output binding

	outputs := make(map[string]interface{})
	outputs["msg"] = string(body)

	resData := make(map[string]interface{})
	resData["body"] = "Webhook event enqueued"
	outputs["res"] = resData

	invokeResponse := InvokeResponse{outputs, nil, nil}
	responseJson, _ := json.Marshal(invokeResponse)

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJson)
}

func JsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "  ")
	if err != nil {
		return in
	}
	return out.String()
}

type CodePushedEvent struct {
	SubscriptionID string `json:"subscriptionId"`
	NotificationID int    `json:"notificationId"`
	ID             string `json:"id"`
	EventType      string `json:"eventType"`
	PublisherID    string `json:"publisherId"`
	Message        struct {
		Text     string `json:"text"`
		HTML     string `json:"html"`
		Markdown string `json:"markdown"`
	} `json:"message"`
	DetailedMessage struct {
		Text     string `json:"text"`
		HTML     string `json:"html"`
		Markdown string `json:"markdown"`
	} `json:"detailedMessage"`
	Resource struct {
		Commits []struct {
			CommitID string `json:"commitId"`
			Author   struct {
				Name  string    `json:"name"`
				Email string    `json:"email"`
				Date  time.Time `json:"date"`
			} `json:"author"`
			Committer struct {
				Name  string    `json:"name"`
				Email string    `json:"email"`
				Date  time.Time `json:"date"`
			} `json:"committer"`
			Comment string `json:"comment"`
			URL     string `json:"url"`
		} `json:"commits"`
		RefUpdates []struct {
			Name        string `json:"name"`
			OldObjectID string `json:"oldObjectId"`
			NewObjectID string `json:"newObjectId"`
		} `json:"refUpdates"`
		Repository struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			URL     string `json:"url"`
			Project struct {
				ID             string `json:"id"`
				Name           string `json:"name"`
				URL            string `json:"url"`
				State          string `json:"state"`
				Visibility     string `json:"visibility"`
				LastUpdateTime string `json:"lastUpdateTime"`
			} `json:"project"`
			DefaultBranch string `json:"defaultBranch"`
			RemoteURL     string `json:"remoteUrl"`
		} `json:"repository"`
		PushedBy struct {
			DisplayName string `json:"displayName"`
			ID          string `json:"id"`
			UniqueName  string `json:"uniqueName"`
		} `json:"pushedBy"`
		PushID int       `json:"pushId"`
		Date   time.Time `json:"date"`
		URL    string    `json:"url"`
	} `json:"resource"`
	ResourceVersion    string `json:"resourceVersion"`
	ResourceContainers struct {
		Collection struct {
			ID string `json:"id"`
		} `json:"collection"`
		Account struct {
			ID string `json:"id"`
		} `json:"account"`
		Project struct {
			ID string `json:"id"`
		} `json:"project"`
	} `json:"resourceContainers"`
	CreatedDate time.Time `json:"createdDate"`
}

type PullRequestCreated struct {
	SubscriptionID string `json:"subscriptionId"`
	NotificationID int    `json:"notificationId"`
	ID             string `json:"id"`
	EventType      string `json:"eventType"`
	PublisherID    string `json:"publisherId"`
	Message        struct {
		Text     string `json:"text"`
		HTML     string `json:"html"`
		Markdown string `json:"markdown"`
	} `json:"message"`
	DetailedMessage struct {
		Text     string `json:"text"`
		HTML     string `json:"html"`
		Markdown string `json:"markdown"`
	} `json:"detailedMessage"`
	Resource struct {
		Repository struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			URL     string `json:"url"`
			Project struct {
				ID             string `json:"id"`
				Name           string `json:"name"`
				URL            string `json:"url"`
				State          string `json:"state"`
				Visibility     string `json:"visibility"`
				LastUpdateTime string `json:"lastUpdateTime"`
			} `json:"project"`
			DefaultBranch string `json:"defaultBranch"`
			RemoteURL     string `json:"remoteUrl"`
		} `json:"repository"`
		PullRequestID int    `json:"pullRequestId"`
		Status        string `json:"status"`
		CreatedBy     struct {
			DisplayName string `json:"displayName"`
			URL         string `json:"url"`
			ID          string `json:"id"`
			UniqueName  string `json:"uniqueName"`
			ImageURL    string `json:"imageUrl"`
		} `json:"createdBy"`
		CreationDate          time.Time `json:"creationDate"`
		Title                 string    `json:"title"`
		Description           string    `json:"description"`
		SourceRefName         string    `json:"sourceRefName"`
		TargetRefName         string    `json:"targetRefName"`
		MergeStatus           string    `json:"mergeStatus"`
		MergeID               string    `json:"mergeId"`
		LastMergeSourceCommit struct {
			CommitID string `json:"commitId"`
			URL      string `json:"url"`
		} `json:"lastMergeSourceCommit"`
		LastMergeTargetCommit struct {
			CommitID string `json:"commitId"`
			URL      string `json:"url"`
		} `json:"lastMergeTargetCommit"`
		LastMergeCommit struct {
			CommitID string `json:"commitId"`
			URL      string `json:"url"`
		} `json:"lastMergeCommit"`
		Reviewers []struct {
			ReviewerURL interface{} `json:"reviewerUrl"`
			Vote        int         `json:"vote"`
			DisplayName string      `json:"displayName"`
			URL         string      `json:"url"`
			ID          string      `json:"id"`
			UniqueName  string      `json:"uniqueName"`
			ImageURL    string      `json:"imageUrl"`
			IsContainer bool        `json:"isContainer"`
		} `json:"reviewers"`
		Commits []struct {
			CommitID string `json:"commitId"`
			URL      string `json:"url"`
		} `json:"commits"`
		URL   string `json:"url"`
		Links struct {
			Web struct {
				Href string `json:"href"`
			} `json:"web"`
			Statuses struct {
				Href string `json:"href"`
			} `json:"statuses"`
		} `json:"_links"`
	} `json:"resource"`
	ResourceVersion    string `json:"resourceVersion"`
	ResourceContainers struct {
		Collection struct {
			ID string `json:"id"`
		} `json:"collection"`
		Account struct {
			ID string `json:"id"`
		} `json:"account"`
		Project struct {
			ID string `json:"id"`
		} `json:"project"`
	} `json:"resourceContainers"`
	CreatedDate time.Time `json:"createdDate"`
}

type PullRequestUpdated struct {
	SubscriptionID string `json:"subscriptionId"`
	NotificationID int    `json:"notificationId"`
	ID             string `json:"id"`
	EventType      string `json:"eventType"`
	PublisherID    string `json:"publisherId"`
	Message        struct {
		Text     string `json:"text"`
		HTML     string `json:"html"`
		Markdown string `json:"markdown"`
	} `json:"message"`
	DetailedMessage struct {
		Text     string `json:"text"`
		HTML     string `json:"html"`
		Markdown string `json:"markdown"`
	} `json:"detailedMessage"`
	Resource struct {
		Repository struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			URL     string `json:"url"`
			Project struct {
				ID             string `json:"id"`
				Name           string `json:"name"`
				URL            string `json:"url"`
				State          string `json:"state"`
				Visibility     string `json:"visibility"`
				LastUpdateTime string `json:"lastUpdateTime"`
			} `json:"project"`
			DefaultBranch string `json:"defaultBranch"`
			RemoteURL     string `json:"remoteUrl"`
		} `json:"repository"`
		PullRequestID int    `json:"pullRequestId"`
		Status        string `json:"status"`
		CreatedBy     struct {
			DisplayName string `json:"displayName"`
			URL         string `json:"url"`
			ID          string `json:"id"`
			UniqueName  string `json:"uniqueName"`
			ImageURL    string `json:"imageUrl"`
		} `json:"createdBy"`
		CreationDate          time.Time `json:"creationDate"`
		ClosedDate            time.Time `json:"closedDate"`
		Title                 string    `json:"title"`
		Description           string    `json:"description"`
		SourceRefName         string    `json:"sourceRefName"`
		TargetRefName         string    `json:"targetRefName"`
		MergeStatus           string    `json:"mergeStatus"`
		MergeID               string    `json:"mergeId"`
		LastMergeSourceCommit struct {
			CommitID string `json:"commitId"`
			URL      string `json:"url"`
		} `json:"lastMergeSourceCommit"`
		LastMergeTargetCommit struct {
			CommitID string `json:"commitId"`
			URL      string `json:"url"`
		} `json:"lastMergeTargetCommit"`
		LastMergeCommit struct {
			CommitID string `json:"commitId"`
			URL      string `json:"url"`
		} `json:"lastMergeCommit"`
		Reviewers []struct {
			ReviewerURL interface{} `json:"reviewerUrl"`
			Vote        int         `json:"vote"`
			DisplayName string      `json:"displayName"`
			URL         string      `json:"url"`
			ID          string      `json:"id"`
			UniqueName  string      `json:"uniqueName"`
			ImageURL    string      `json:"imageUrl"`
			IsContainer bool        `json:"isContainer"`
		} `json:"reviewers"`
		Commits []struct {
			CommitID string `json:"commitId"`
			URL      string `json:"url"`
		} `json:"commits"`
		URL   string `json:"url"`
		Links struct {
			Web struct {
				Href string `json:"href"`
			} `json:"web"`
			Statuses struct {
				Href string `json:"href"`
			} `json:"statuses"`
		} `json:"_links"`
	} `json:"resource"`
	ResourceVersion    string `json:"resourceVersion"`
	ResourceContainers struct {
		Collection struct {
			ID string `json:"id"`
		} `json:"collection"`
		Account struct {
			ID string `json:"id"`
		} `json:"account"`
		Project struct {
			ID string `json:"id"`
		} `json:"project"`
	} `json:"resourceContainers"`
	CreatedDate time.Time `json:"createdDate"`
}
