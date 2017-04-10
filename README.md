appconf
=======

Application settings configurator.

[https://appconf.adhocteam.us/](https://appconf.adhocteam.us/)

This requires a S3 bucket to store the configuration files in, and to run as a
AWS IAM user with permission to list, get, put, and delete objects on that bucket.

The S3 bucket structure is:

```
/
  application1/
    dev/
      VAR1
      VAR2
      ...
    test/
      VAR1
      VAR2
      ...
    ...
  application2/
    dev/
      VAR1
      VAR2
      ...
    test/
      VAR1
      VAR2
      ...
    ...
  ...
```

You must provide an inventory file that lists each environment and application you'd like
to configure. See [inventory.json.example](inventory.json.example) for a sample file.

### Theory of operation

Configures apps by setting the key-value pairs as objects in S3. The idea is
that the key-values that can be created, edited, and deleted here will be turned
into environment variables of the application. The mechanism by which that
happens is intended to be that the instance (via a user-data script) or the app
itself pulls down the conf vars from S3 at boot-time. Since these apps are
managed by `runit`, they can use
[`chpst`](http://smarden.org/runit/chpst.8.html)'s env dir facility to create
the environment variables from these files. This mechanism is outside this app
-- the app merely (at the moment) manages the writing and updating of conf vars.

Build dependencies
------------------

* **Go 1.5+**
* **TypeScript 1.x** -- `npm install -g typescript`

Runtime dependencies
--------------------

* **AWS credentials** -- add them to the shell environment
* inventory.json file -- see [inventory.json.example](inventory.json.example) for a skeleton file

Installation
------------

``` shell
$ go get github.com/adhocteam.us/appconf
```

Usage
-----

``` shell
$ cd $GOPATH/github.com/adhocteam.us/appconf -inv /path/to/inventory.json -bucket your-bucket-name -l :8080
$ go install
$ # ensure AWS credentials are set in the environment or $HOME/.aws/credentials
$ $GOBIN/appconf
$ open http://localhost:8080/
```

Building for Linux target
-------------------------

``` shell
$ cd $GOPATH/github.com/adhocteam.us/appconf
$ make rpm
$ scp appconf-1.0-1.x86_64.rpm server:
```
