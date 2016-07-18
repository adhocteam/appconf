appconf
=======

Application settings configurator.

[https://appconf.adhocteam.us/](https://appconf.adhocteam.us/)

Configures the following applications:

 * Marketplace API
 * Window Shopping
 * Plan Compare 2.0

Environments supported:

 * dev
 * test
 * imp1a
 * imp1b
 * prod

### Theory of operation

Configures apps by setting the key-value pairs as objects in S3. The idea is
that the key-values that can be created, edited, and deleted here will be turned
into environment variables of the application. The mechanism by which that
happens is intended to be that the instance (via a user-data script) or the app
itself pulls down the conf vars from S3 at boot-time. This mechanism is outside
this app -- it merely (at the moment) manages the writing and updating of conf
vars.

Build dependencies
------------------

* **Go 1.5+**
* **TypeScript 1.x** -- `npm install -g typescript`

Runtime dependencies
--------------------

* **AWS credentials** -- add them to the shell environment

Installation
------------

``` shell
$ go get github.com/adhocteam.us/appconf
```

Usage
-----

``` shell
$ cd $GOPATH/github.com/adhocteam.us/appconf
$ # ensure AWS credentials are set in the environment
$ appconf
$ open http://localhost:8080/
```

Building for Linux target
-------------------------

``` shell
$ cd $GOPATH/github.com/adhocteam.us/appconf
$ make rpm
$ scp appconf-1.0-1.x86_64.rpm server:
```
