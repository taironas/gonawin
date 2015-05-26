# gonawin

A social web application to make "friendly bets" on sport games with your friends.

We believe that todays betting platforms are full of adds that polute the game itself. We want to create a platform that people can use to bet/play with there friends and family or anyone they want without anything in between. There is no money involved.


##### Notes:

As we are in early stages data might flight away from here :-) .

[![baby-gopher](https://raw.github.com/drnic/babygopher-site/gh-pages/images/babygopher-badge.png)](http://www.babygopher.org)

This is a baby-gophers project because `The young are born blind and helpless` - [Wikipedia](http://en.wikipedia.org/wiki/Gopher_(animal)#Pocket_gopher) and [babygopher.org](http://www.babygopher.org/) ^^


#### Contributors :
* [Santiago](https://github.com/santiaago)
* [Remy](https://github.com/rjourde)

#### Third Parties Installation

    go get github.com/garyburd/go-oauth/oauth

#### Third Parties
* [Google App Engine for Go](https://developers.google.com/appengine/docs/go/)
* [Boostrap v3](http://getbootstrap.com/)
* [Angularjs](http://angularjs.org/)
* [go-oauth](http://github.com/garyburd/go-oauth)
* [flags](https://github.com/lipis/flag-icon-css)
* [avatars](https://http://www.tinygraphs.com)
* [Social Buttons for Bootstrap](http://lipis.github.io/bootstrap-social/)
* icons
  * [font-awesome](http://fortawesome.github.io/Font-Awesome/icons/)
  * [glyphicons](http://glyphicons.com/)

#### Installation

* install [go](http://golang.org/doc/install)
* set up your [environement](http://golang.org/doc/code.html)
* install the [go appengine sdk](https://developers.google.com/appengine/downloads)
* set up the appengine [environement](https://developers.google.com/appengine/docs/go/gettingstarted/devenvironment)
*   `go get github.com/garyburd/go-oauth/oauth`
*   `go get github.com/santiaago/gonawin`
*   `cd $GOPATH/src/github.com/santiaago/gonawin/gonawin`
*   `cp example-config.json config.json`
*   add your email to the `config.json` in the `admins` section
*   `goapp serve`

#### Run App

    > cd $GOPATH/src/github.com/santiaago/gonawin/gonawin
    > goapp serve

#### Run App with production datastore backup

Create datastore backup

    > go_appengine\appcfg.py download_data --url=http://www.gonawin.com/_ah/remote_api --filename=bck_gonawin_mmddyyyy

Run local server

    > cd $GOPATH/src/github.com/santiaago/gonawin/gonawin
    > goapp serve

Connect datastore backup to the local server

    > go_appengine\appcfg.py upload_data --url=http://localhost:8080/_ah/remote_api --filename=bck_gonawin_mmddyyyy

### Run App with clean data store:

    > cd $GOPATH/src/github.com/santiaago/gonawin/gonawin
    > dev_appserver.py --clear_datastore=yes .

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

    go fmt ..\gonawin\...

#### Deployment

    goapp deploy


__Note:__ If deployment hangs rollback it by doing:

#####On OSX:

    appcfg rollback ..

#####On Windows:

    python appcfg.py rollback $GOPATH\src\github.com\santiaago\gonawin\gonawin


#### Documentation

    godoc -http=:6060
