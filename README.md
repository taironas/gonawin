# gonawin

A social web application for betting with your friends.

We believe that todays betting platforms are just full of adds that polute the game itself.
We want to create a platform that people can use to bet/play with there friends and family or anyone they want without anything in between.


#### Contributors :
* Santiago (sar)
* Remy (rej)
* Cabernal (cab)

#### Third Parties Installation

    go get github.com/garyburd/go-oauth/oauth
    
#### Third Parties
* [Google App Engine for Go](https://developers.google.com/appengine/docs/go/)
* [Boostrap v3](http://getbootstrap.com/)
* [Angularjs](http://angularjs.org/)
* [go-oauth](github.com/garyburd/go-oauth/oauth)
* [flags](https://github.com/lipis/flag-icon-css)
* [avatars](https://github.com/cupcake/sigil)
* icons
  * [font-awesome](http://fortawesome.github.io/Font-Awesome/icons/)
  * [glyphicons](http://glyphicons.com/)
    
#### Installation

* install [go](http://golang.org/doc/install)
* set up your [environement](http://golang.org/doc/code.html)
* install the [go appengine sdk](https://developers.google.com/appengine/downloads)
* set up the appengine [environement](https://developers.google.com/appengine/docs/go/gettingstarted/devenvironment)
*   `go get github.com/garyburd/go-oauth/oauth`
*   `go get github.com/santiaago/purple-wing`
*   `cd $GOPATH/src/github.com/santiaago/purple-wing`
*   `goapp serve`

#### Run App

    go_appengine\dev_appserver.py purple-wing
    
#### Formatting

    go fmt ..\purple-wing\...

#### Deployment

    go_appengine\appcfg.py update purple-wing
    
#### Documentation

    godoc -http=:6060
