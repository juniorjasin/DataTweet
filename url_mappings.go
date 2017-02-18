package main

import(
"fmt"
"net/http"

"github.com/juniorjasin/datatweet/controllers"
)

var pt int = 8888
var port *int = &pt

func mapURLsToControllers(){
	// authentication
	http.HandleFunc("/permission", controllers.RedirectUserToTwitter)
	http.HandleFunc("/maketoken", controllers.GetTwitterToken)
	http.HandleFunc("/favorites", controllers.Favorites)
	http.HandleFunc("/dictionary", controllers.Dictionary)
	u := fmt.Sprintf(":%d", *port)
	fmt.Printf("Listening on '%s'\n", u)
	http.ListenAndServe(u, nil)
}
