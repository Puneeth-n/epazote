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

## The problem to solve
Once your site/application is up and running, it can become idle and
unresponsive, your supervisor will not notice this, since in most of the cases
is just responsible for keeping your App process up and running no matter how it
is behaving, therefore exists the need to monitor the status of the application
and based on the responses take actions.

When doing Continuous Deployment "[CD](https://en.wikipedia.org/wiki/Continuous_delivery)"
if the ping, healthcheck, status, etc; endpoints change, it implies making changes
in order to properly monitor the application, this creates a dependency or extra
task apart from the "CD" process, therefore exists the need to detect any changes
and automatically apply them upon request.

## How it works
In its basic way of operation, **Epazote** periodically checks the services endpoints
"[URLs](https://en.wikipedia.org/wiki/Uniform_Resource_Locator)"
by doing an [HTTP GET Request](https://en.wikipedia.org/wiki/Hypertext_Transfer_Protocol#Request_methods),
based on the response [Status code](https://en.wikipedia.org/wiki/List_of_HTTP_status_codes),
[Headers](https://en.wikipedia.org/wiki/List_of_HTTP_header_fields) or
either the
[body](https://en.wikipedia.org/wiki/HTTP_message_body), it executes a command.

In most scenarios, is desired to apply a command directly to the application in
cause, like a signal (``kill -HUP``), or either a restart (``sv restart app``),
therefore in this case **Epazote** and the application should be running on the
same server.

**Epazote** can also work in a standalone mode by only monitoring and sending
alerts if desired.

# How to use it
First you need to install **Epazote**, either you can compile it from
[source](https://github.com/nbari/epazote)
or download a pre-compiled binary matching your operating system.

> To compile from source, after downloading the sources use ``make`` to build the binary

**Epazote** was designed with simplicity in mind, as an easy tool for
[DevOps](https://en.wikipedia.org/wiki/DevOps) and as a complement to
infrastructure orchestration tools like [Ansible](http://www.ansible.com/) and
[SaltStack](http://saltstack.com/), because of this [YAML](http://www.yaml.org/)
is used for the configuration files, avoiding with this, the learn of a new
language or syntax, and simplifying the setup.

## The configuration file

The configuration file ([YAML formated](https://en.wikipedia.org/wiki/YAML))
consists of two parts, a **config** and a **services** (Key-value pairs).

## The config section

The **config** section is composed of three settings:

    - post (Url to post the logs)
    - smtp (Email settings for sending notification)
    - scan (Paths used to find the file 'epazote.yml')

Example:

```yaml
config:
    post: http://domain.tld/get/json/
    smtp:
        username: smtp@domain.tld
        password: password
        server: mail.example.com
        port: 587
        headers:
            from: epazote@domain.tld
            to: team@domain.tld ops@domain.tld etc@domain.tld
            subject: [%s -%s], Service, Status
    scan:
        paths:
            - /arena/home/sites
            - /home/apps
        minutes: 5
```

### config - post

An URL to post all activity related to the services if log is enable on the
service. The **post** setting is optional.

### config - smtp

Required to properly send alerts via email, all fields are required, the
``headers`` section can be extended with any desired key-pair values.

### config - scan

Paths to scan every N ``seconds``, ``minutes`` or ``hours``, a search for
services specified in a file call ``epazote.yml`` is made.

The **scan** setting is optional however is very useful when doing Continues
Deployments. for example if your code is automatically uploaded to the
directory ``/arena/home/sites/application_1`` and your scan paths contain
``/arena/home/sites``, you could simple upload on your application directory a
file named ``epazote.yml`` with the service rules, thus achieving the deployment
of your application and the supervising at the same time.

### config (optional)

As you may notice the ``config`` section contains mainly settings for sending
alerts/notifications apart from the ``scan`` setting, therefore is totally
optional, meaning that **Epazote** can still run and check your services without
the need of the ``config`` section.

If you want to automatically update/load services you will need the
``config - scan`` setting.


## The services section

Services are the main functionality of **Epazote**, is where the URL's and the
rules based on the response are defined, since options vary from service to
service, an example could help better to understand the setup:

```yaml
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

### services - name of service (string)
An unique string that identifies your service, in the above example, there are 3
services named:
 - my service 1
 - other service
 - redirect service

### services - url (string)
URL of the service to supervise

### services - timeout in seconds (int)
Timeout specifies a time limit for the HTTP requests, A value of zero means no
timeout, defaults to 5 seconds.

### services - seconds, minutes, hours
How often to check the service, the options are: (Only one should be used)
 - seconds N
 - minutes N
 - hours N

``N`` should be an integer.

### services - log (bool)
If set to true, it will post all events to the defined ``post`` URL on the
**config** section, defaults to false.

### services - expect
The ``expect`` block options are:
- status (int)
- header (string)
- body   (regular expression)
- if_not (Action block)

### services - expect - status
An Integer representing the expected [HTTP Status Code](https://en.wikipedia.org/wiki/List_of_HTTP_status_codes)

### services - expect - header
A key-value map of expected headers, it can be only one or more.

### services - expect - body
A [regular expression](https://en.wikipedia.org/wiki/Regular_expression) used
to match a string on the body of the site, use full in cases you want to ensure
that the content delivered is always the same or keeps a pattern.

### services - expect (How it works)
The ``expect`` logic tries to implement a
[if-else](https://en.wikipedia.org/wiki/if_else) logic ``status``, ``header``,
``body`` are the **if** and the ``if_not`` block becomes the **else**.

    if
        status
        header
        body
    else:
        if_not

In must cases only one option is required, check on the above example for the service named "redirect service".

In case that more than one option is used, this is the order in how they are evaluated, no meter how they where introduced on the configuration file:

    1. body
    2. status
    3. header

The reason for this order is related to performance, at the end we want to
monitor/supervise the services in an efficient way avoiding to waste extra
resources, in must cases only the HTTP Headers are enough to take an action,
therefore we don't need to read the full body page, because of this if no
``body`` is defined, **Epazote** will only read the Headers saving with this
time and process time.

### services - expect - if_not
``if_not`` is a block with an action of what to do it we don't get what we where
expecting (``expect``). See services - Actions

### services - if_status  & if_header
There maybe cases in where third-party dependencies are down and because of this
your application could not be working properly, for this cases the ``if_status``
and ``if_header`` could be useful.

For example if the database is your application could start responding an status
code 500 or either a custom header and based on does values take execute an
action:

The format for ``if_status`` is a key-pair where key is an int representing an
HTTP status code, and the value an Action option

The format for ``if_header`` is a key-pair where key is a string of something
you could relate/match and has in other if_X conditions, value is an Action.

This are the only ``if's`` and the order of execution:
 1. if_status
 2. if_header
 3. if_not

This means that if a service uses ``if_status`` and ``if_not``, it will
evaluate first the ``if_status`` and execute an Action if required, in case
an ``if_status`` and ``if_header`` are set, same applies, first is evaluated
``if_status``, then ``if_header`` and last ``if_not``.

## services - Actions
An Action has tree options:
 - cmd
 - notify
 - msg

They can be used all together, only one or either none.

### services - Actions - cmd (string)
``cmd`` Contains the command to be executed.

### services - Actions - notify (string)
``notify`` Should contain the email email address or addresses (space separated)
of the recipients that will be notified when the action is executed.

### services - Actions - msg (string)
``msg`` The message to send when the action is executed.

## Extra setup
*green dots give some comfort* because of this when using the ``log`` option an
extra service could be configure as a receiver for all the post that **Epazote**
produce and based on the data obtained create a custom dashboard, something
similar to: https://status.cloud.google.com/ or http://status.aws.amazon.com/
