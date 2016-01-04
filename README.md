[![Build Status](https://drone.io/github.com/nbari/epazote/status.png)](https://drone.io/github.com/nbari/epazote/latest)
[![Build Status](https://travis-ci.org/nbari/epazote.svg?branch=develop)](https://travis-ci.org/nbari/epazote)

# Epazote 🌿
Microservices supervisor

## Why ?
There are good supervisors,
[daemontools](https://cr.yp.to/daemontools.html),
[runit](http://smarden.org/runit/) just to mention some, on most cases is just
a matter of uploading code to the server, create a run script and you are all
set, your code will start up and live forever, so far so good, but let's face
it, "stuff happens", suddenly the site or application can stop responding
request, display unwanted content, etc. here is where **Epazote** comes into
action.

## The problem to solve.
Once your site/application is up and running, it can become idle and
unresponsive, your supervisor will not notice this, since in most of the cases
is just responsible for keeping your App process up and running no matter how it
is behaving, therefore exist the need to monitor the status of the application
and based on the responses take actions.

When doing Continuous Deployment "[CD](https://en.wikipedia.org/wiki/Continuous_delivery)"
if the ping, healthcheck, status, etc; endpoints change, it imply making changes
in order to properly monitor the application, this creates a dependency or extra
task apart from the CD process, therefore exist the need to detect any changes
and automatically apply them upon request.

## How it works.
In its basic way of operation, **Epazote** periodically checks the services endpoints
"[URLs](https://en.wikipedia.org/wiki/Uniform_Resource_Locator)"
by doing an [HTTP GET Request](https://en.wikipedia.org/wiki/Hypertext_Transfer_Protocol#Request_methods),
based on the response [Status code](https://en.wikipedia.org/wiki/List_of_HTTP_status_codes),
[Headers](https://en.wikipedia.org/wiki/List_of_HTTP_header_fields) or
either the
[body](https://en.wikipedia.org/wiki/HTTP_message_body), it executes a command.

In most scenarios, is desired to apply a command directly to the application in
cause, like a signal (``kill -HUP``), or either a restart (``sv restart app``),
therefore in this case **Epazote** and the application should be running on the same
server.

**Epazote** can also work in a standalone mode by only monitoring and sending alerts
if desired.

# How to use it.
First you need to install **Epazote**, either you can compile it from [source](https://github.com/nbari/epazote)
or download a pre-compiled binary matching your operating system.

> To compile from source, after downloading the sources use ``make`` to build the binary

**Epazote** was designed with simplicity in mind, as an easy tool for
[DevOps](https://en.wikipedia.org/wiki/DevOps) and as a complement to
infrastructure orchestration tools like [Ansible](http://www.ansible.com/) and
[SaltStack](http://saltstack.com/), because of this [YAML](http://www.yaml.org/)
is used for the configuration files, avoiding with this, the learn of a new
language or syntax, and simplifying the setup.

## The config file.

This is an example/structure of the config file:

```yaml
config:
    smtp:
        username: smtp@domain.tld
        password: password
        host: smtp.domain.tld
        port: 587
        tls: true
        headers:
            from: epazote@domain.tld
            to: team@domain.tld ops@domain.tld etc@domain.tld
            subject: [%s -%s], Service, Status
    http:
        host: http://domain.tld/post/here
    scan:
        paths:
            - /arena/home/sites
            - /home/apps
        minutes: 5

services:
    my service 1:
        url: http://myservice.domain.tld/_healthcheck_
        timeout: 5
        seconds: 60
        log: True
        expect:
            status: 200
            header:
                content-type: application/json; charset=UTF-8
            body: find this string on my site
        if_not:
            cmd: sv restart /services/my_service_1
            notify: team@domain.tld
            msg: |
                line 1 bla bla
                line 2
        if_status:
            500:
                cmd: reboot
            404:
                cmd: sv restart /services/cache
                msg: restarting cache
                notify: team@domain.tld x@domain.tld
        if_header:
            x-amqp-kapputt:
                cmd: restart abc
                notify: bunny@domain.tld
                msg: |
                    The rabbit is angry
                    & hungry
            x-db-kapputt:
                cmd: svc restart /services/db

    other service:
        url: http://other-service.domain.tld/ping
        minutes: 3

    redirect service:
        url: http://test.domain.tld/
        hour: 1
        expect:
            status: 302
        if_not:
            cmd: service restart abc
            notify: abc@domain.tld
```
