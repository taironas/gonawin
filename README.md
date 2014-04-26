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
*   `cd $GOPATH/src/github.com/santiaago/purple-wing/purple-wing`
*   `goapp serve`

#### Run App

    > cd $GOPATH/src/github.com/santiaago/purple-wing/purple-wing
    > goapp serve

#### Access app from smartphone

Get your ip in this case `192.168.1.X`

    > ifconfig
    ...
    inet 192.168.1.X netmask 0xffffff00 broadcast 192.168.1.255
    ...
    
Run the app with the `-host` parameter

    > goapp serve -host=0.0.0.0
    INFO     2014-04-26 10:25:33,296 devappserver2.py:764] Skipping SDK update check.
    WARNING  2014-04-26 10:25:33,299 api_server.py:374] Could not initialize images API; you are likely missing the Python "PIL" module.
    INFO     2014-04-26 10:25:33,302 api_server.py:171] Starting API server at: http://localhost:53542
    INFO     2014-04-26 10:25:33,305 dispatcher.py:182] Starting module "default" running at: http://0.0.0.0:8080
    INFO     2014-04-26 10:25:33,307 admin_server.py:117] Starting admin server at: http://localhost:8000

access from your smartphone on `http://192.168.1.X:8080/ng`
    
#### Formatting

    go fmt ..\purple-wing\...

#### Deployment

    go_appengine\appcfg.py update purple-wing
    
#### Documentation

    godoc -http=:6060
