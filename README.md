##Contributors :
Santiago (sar)
Remy (rej)
Cabernal (cab)


##Ressources:

###Go:
Set up your work environment:
http://golang.org/doc/code.html

Start by taking A Tour of Go:
http://tour.golang.org/

Build a web application by following the Wiki Tutorial:
http://golang.org/doc/articles/wiki/

Follow pluralsight tutorial:
http://pluralsight.com/training/Courses/TableOfContents/go

###GAE:
You can follow udacity cs253  (but it is in python :s)
https://www.udacity.com/course/cs253

To understand GAE and Go I think this tutorial should help us: 
https://developers.google.com/appengine/docs/go/gettingstarted/

Some Go tutorials:
http://golangtutorials.blogspot.fr/2011/05/table-of-contents.html


###Code guidelines:

From: http://golang.org/doc/effective_go.html#names

The package name should be good: short, concise, evocative. By convention, packages are given lower case, single-word names; there should be no need for underscores or mixedCaps. 

Another convention is that the package name is the base name of its source directory; the package in src/pkg/encoding/base64 is imported as *"encoding/base64"* but has name *base64*, not *encoding_base64* and not *encodingBase64*.

Another short example is *once.Do;* *once.Do(setup)* reads well and would not be improved by writing *once.DoOrWaitUntilDone(setup)*. Long names don't automatically make things more readable. A helpful doc comment can often be more valuable than an extra long name.
