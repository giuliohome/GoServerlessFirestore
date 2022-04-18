// Package p contains an HTTP Cloud Function.
package p

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
 	"time"
 	"strings"
	 b64 "encoding/base64"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func createClient(ctx context.Context) *firestore.Client {
	// Sets your Google Cloud Platform project ID.
	projectID := "my-cloud-giulio"

	// [END firestore_setup_client_create]
	// Override with -project flags
	// flag.StringVar(&projectID, "project", projectID, "The Google Cloud Platform project ID.")
	// flag.Parse()

	// [START firestore_setup_client_create]
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	// Close client when done with
	// defer client.Close()
	return client
}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}


// HelloWorld prints the JSON encoded "message" field in the body
// of the request or "Hello, World!" if there isn't one.
func HelloWorld(w http.ResponseWriter, r *http.Request) {
        // Set CORS headers for the preflight request
        if r.Method == http.MethodOptions {
        	w.Header().Set("Access-Control-Allow-Origin", "*")
        	w.Header().Set("Access-Control-Allow-Methods", "POST")
        	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        	w.Header().Set("Access-Control-Max-Age", "3600")
        	w.WriteHeader(http.StatusNoContent)
        	return
	}
        // Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization")
	w.Header().Set("Content-Type", "application/json")

	var d struct {
		Message string `json:"message"`
		Operation string `json:"operation"`
		Written string `json:"written"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		switch err {
		case io.EOF:
			fmt.Fprint(w, "Hello World!")
			return
		default:
			log.Printf("json.NewDecoder: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

	ctx := context.Background()
	client := createClient(ctx)
	defer client.Close()	
	
	authuserinfo := string(r.Header.Get("X-Endpoint-Api-Userinfo"))
	userinfo, _ := b64.StdEncoding.DecodeString(authuserinfo)
	var authperm struct {
		Iss string `json:"iss"`
		Sub string `json:"sub"`
		Aud []string `json:"aud"`
		Iat int `json:"iat"`
		Exp int `json:"exp"`
		Azp string `json:"azp"`
		Scope string `json:"scope"`
		Permissions []string `json:"permissions"`
	}
	if permerr := json.Unmarshal([]byte(string(userinfo)+"}"), &authperm); permerr != nil {
		switch permerr {
		case io.EOF:
			fmt.Fprint(w, "no auth permissions")
			return
		default:
			fmt.Fprint(w, permerr.Error() + " from " + string(userinfo) )
			return
		}
	}
	
	if !contains(authperm.Permissions,"read:items") {
		fmt.Fprint(w, "no read:items in auth permissions")
		return
	}

	if d.Operation == "query" {
		iter := client.Collection("giuliohome").OrderBy("lastupdate", firestore.Desc).Limit(30).Documents(ctx)
		textarea := "user perm info " + strings.Join(authperm.Permissions,",") + "\n"
		for {
			doci, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				fmt.Fprint(w, html.EscapeString(err.Error()))
				return
			}
			wres := doci.Ref.ID + "\n" 
			textarea += wres
		}
		fmt.Fprint(w, html.EscapeString(textarea))
		return
	}
	if d.Operation == "write" {
		doc := make(map[string]interface{})
		t := time.Now()
		doc["text"] = d.Written
		doc["lastupdate"] = t
		_, werr := client.Collection("giuliohome").Doc(d.Message).Set(ctx, doc)
		if werr != nil {
			fmt.Fprint(w, html.EscapeString(werr.Error()))
			return
		}
		wres := fmt.Sprintf("Document %v saved: %#v\n time: %v", 
			d.Message, d.Written, t)

		fmt.Fprint(w, html.EscapeString(wres))
		return
	}

	if d.Message == "" {
		fmt.Fprint(w, "Hello World!")
		return
	}
	
	dsnap, err := client.Collection("giuliohome").Doc(d.Message).Get(ctx)
	if err != nil {
		fmt.Fprint(w, html.EscapeString(err.Error()))
		return
	}
	m := dsnap.Data()

	res := fmt.Sprintf("Document data: %#v\n time: %v", 
		m["text"], m["lastupdate"])

	fmt.Fprint(w, html.EscapeString(res))
}

