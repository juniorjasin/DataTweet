package main

import(
"fmt"
"net/http"

"github.com/juniorjasin/datatweet/controllers"
)

var pt int = 8888
var port *int = &pt

func mapURLsToControllers(){
	http.HandleFunc("/", controllers.RedirectUserToTwitter)
	http.HandleFunc("/maketoken", controllers.GetTwitterToken)
	u := fmt.Sprintf(":%d", *port)
	fmt.Printf("Listening on '%s'\n", u)
	http.ListenAndServe(u, nil)
}
