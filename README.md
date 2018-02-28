appconf
=======

Application settings configurator.

This requires a S3 bucket to store the configuration files in, and to run as a
AWS IAM user or role with permission to list, get, put, and delete objects on that bucket.

Optionally provide an AWS KMS key, used to encrypt the configuration files at rest in your S3 bucket. The KMS key must be configured to allow use by the IAM user or role for appconf.

If a KMS key id is provided, you must enable S3's [signature version 4](http://docs.aws.amazon.com/AmazonS3/latest/dev/UsingAWSSDK.html#specify-signature-version).

To enable signature version 4, edit your `~/.aws/config` file:

```
[default]
s3 =
    signature_version = s3v4
```

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

* **Go 1.8+**
* **TypeScript 1.x** -- `npm install -g typescript`

Runtime dependencies
--------------------

* **AWS credentials** -- add them to the shell environment
* inventory.json file -- see [inventory.json.example](inventory.json.example) for a skeleton file

Installation
------------

``` shell
$ go get github.com/adhocteam/appconf
```

Usage
-----

``` shell
$ cd $GOPATH/src/github.com/adhocteam/appconf
$ make app.js
$ go install
$ # ensure AWS credentials are set in the environment or $HOME/.aws/credentials
$ $GOBIN/appconf -l :8081 -bucket s3-bucket-goes-here -inv inventory.json
$ open http://localhost:8080/
```

Building for Linux target
-------------------------

``` shell
$ cd $GOPATH/src/github.com/adhocteam/appconf
$ make rpm
$ scp appconf-1.0-1.x86_64.rpm server:
```

### IAM Policy Examples

#### Appconf Server

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "s3:ListBucket"
            ],
            "Resource": "arn:aws:s3:::s3-bucket-goes-here"
        },
        {
            "Effect": "Allow",
            "Action": [
                "s3:GetObject",
                "s3:PutObject",
                "s3:DeleteObject"
            ],
            "Resource": "arn:aws:s3:::s3-bucket-goes-here/*"
        },
        {
            "Effect": "Allow",
            "Action": "kms:GenerateDataKey",
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "kms:Decrypt"
            ],
            "Resource": [
                "every KMS ARN for each application"
            ]
        }
    ]
}
```

#### Application Server

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "s3:GetObject"
            ],
            "Resource": "arn:aws:s3:::s3-bucket-goes-here/application/environment/*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "kms:Decrypt"
            ],
            "Resource": "KMS ARN for the application"
        }
    ]
}
```

### Extras

#### Example init script to pull configuration for an application

An example init script is included in misc/app_init.sh; it's rough, but might help you get started.

#### Import Helper

An example Ansible script to help import a bunch of ENV files is included in misc/import.yml. You may use this to help import lots of variables across many services quickly.

#### Packer build and Ansible playbooks

The QPPFE build source for building an Appconf AMI is available on github, if you have access: https://github.com/CMSgov/qpp-deploy/tree/master/packer
