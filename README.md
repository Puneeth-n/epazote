[![Download](https://api.bintray.com/packages/nbari/epazote/epazote/images/download.svg)](https://bintray.com/nbari/epazote/epazote/_latestVersion)
[![Build Status](https://drone.io/github.com/nbari/epazote/status.png)](https://drone.io/github.com/nbari/epazote/latest)
[![Build Status](https://travis-ci.org/nbari/epazote.svg?branch=develop)](https://travis-ci.org/nbari/epazote)
[![Coverage Status](https://coveralls.io/repos/github/nbari/epazote/badge.svg?branch=master)](https://coveralls.io/github/nbari/epazote?branch=master)

# Epazote 🌿
Automated Microservices Supervisor

**Epazote** automatically update/add services specified in a file call
``epazote.yml``. Periodically checks the defined endpoints and execute recovery
commands in case services responses are not behaving like expected helping with
this to automate actions in order to keep services/applications up and running.

In Continuous Integration/Deployment environments the file ``epazote.yml`` can
dynamically be updated/change without need to restart the supervisor, avoiding
with this an extra dependency on the deployment flow which could imply to
restart the supervisor, in this case **Epazote**.

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
or download a pre-compiled binary matching your operating system from here:
https://dl.bintray.com/nbari/epazote/

 [![Download](https://api.bintray.com/packages/nbari/epazote/epazote/images/download.svg)](https://bintray.com/nbari/epazote/epazote/_latestVersion)

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

The **config** section is composed of:

    - smtp (Email settings for sending notification)
    - scan (Paths used to find the file 'epazote.yml')

Example:

```yaml
config:
    smtp:
        username: smtp@domain.tld
        password: password
        server: mail.example.com
        port: 587
        headers:
            from: epazote@domain.tld
            to: team@domain.tld ops@domain.tld etc@domain.tld
            subject: "[name - status]"
    scan:
        paths:
            - /arena/home/sites
            - /home/apps
        minutes: 5
```

### config - smtp

Required to properly send alerts via email, all fields are required, the
``headers`` section can be extended with any desired key-pair values.

### config - smtp - subject (because exit name output status url)
The subject can be formed by using this keywords: ``because`` ``exit`` ``name``
``output`` ``status`` ``url`` on the previous example, ``subject: [name - status]``
would transform to ``[my service - 500]`` the ``name`` has replaced
by the service name, ``my service`` and ``status`` by the response status code
``500`` in this case.

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
        log: http://monitor.domain.tld
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
                notify: yes

    salt-master:
        test: pgrep -f salt
        if_not:
            cmd: service restart salt_master
            notify: operations@domain.tld
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

### services - log (URL)
An URL to post all events, default disabled.

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
``notify`` Should contain ``yes``, the email email address or addresses (space separated)
of the recipients that will be notified when the action is executed.

If the string is ``yes`` the global recipients will be used.

### services - Actions - msg (string)
``msg`` The message to send when the action is executed.

## services - Test
**Epazote** It is mainly used for HTTP services, for supervising other
applications that don't listen or accept HTTP connections, like a database,
cache engine, etc. There are tools like
[daemontools](https://cr.yp.to/daemontools.html),
[runit](http://smarden.org/runit/) as already mentioned, even so, **Epazote**
can eventually be used to execute an action based on the exit of a command
for example:

```yaml
    salt-master:
        test: pgrep -f salt
        if_not:
            cmd: service restart salt_master
            notify: operations@domain.tld
```

In this case: ``test: pgrep -f salt`` will execute the ``cmd`` on the ``if_not``
block in case the exit code is > 0, from the ``pgrep`` man page:

```txt
EXIT STATUS
     The pgrep and pkill utilities return one of the following values upon exit:

          0       One or more processes were matched.
          1       No processes were matched.
          2       Invalid options were specified on the command line.
          3       An internal error occurred.
```


## Extra setup
*green dots give some comfort* -- Because of this when using the ``log``
option an extra service could be configure as a receiver for all the post
that **Epazote** produce and based on the data obtained create a custom
dashboard, something similar to: https://status.cloud.google.com/ or
http://status.aws.amazon.com/

# Issues
Please report any problem, bug, here: https://github.com/nbari/epazote/issues
