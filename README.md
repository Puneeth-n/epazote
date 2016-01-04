[![Build Status](https://drone.io/github.com/nbari/epazote/status.png)](https://drone.io/github.com/nbari/epazote/latest)
[![Build Status](https://travis-ci.org/nbari/epazote.svg?branch=develop)](https://travis-ci.org/nbari/epazote)

# Epazote 🌿
Microservices supervisor

## Why ?
There are good supervisors,
[daemontools](https://cr.yp.to/daemontools.html),
[runit](http://smarden.org/runit/) just to mention some, on must cases is just
a matter of uploading code to the server, create a run script and you are all
set, your code will start up and live forever, so far so good, but let's face
it, "stuff happens" and here is where Epazote comes into action.

## The problem to solve.
Once your code/application is up and running, it can become idle and
unresponsive, your supervisor will not notice this, since in most of the cases
is just responsible for keeping your app process up and running no matter how it
is behaving, therefore exist the need to monitor the status of the application
and based on the responses take actions.

When doing Continuous Deployment "[CD](https://en.wikipedia.org/wiki/Continuous_delivery)"
if the ping, healthcheck, status, etc; endpoints change, it imply making changes
in order to properly monitor the application, this creates a dependency or extra
task apart from the CD process, therefore exist the need to detect any changes
and automatically apply them upon request.

## How it works.
In its basic way of operation, Epazote periodically checks the services endpoints
"[URLs](https://en.wikipedia.org/wiki/Uniform_Resource_Locator)"
by doing an HTTP GET Request, based on the response Status code, Headers or
either the body, it executes a command.

In most scenarios, is desired to apply a command directly to the application in
cause, like a signal (``kill -HUP``), or either a restart (``sv restart app``),
therefore in this case Epazote and the application should be running on the same
server.

Epazote can also work in a standalone mode by only monitoring and sending alerts
if desired.
