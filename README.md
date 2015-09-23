![status: inactive](https://img.shields.io/badge/status-inactive-red.svg)

This project is no longer actively developed or maintained.  

For more information about Compute Engine, refer to our [documentation](https://cloud.google.com/compute).

Docker Cloud
============

What is it?
------------
Docker Cloud is a proxy for the Docker API which automatically creates and destroys cloud virtual machines to run
your docker containers.

Why would I want to do that?
------------
If you are running Docker on OS X or Windows, there is no longer any need to install a virtualization layer like
vagrant on your machine.  You can simply run it in the cloud.  Additionally, if you want to easily turn up and
down containers into a cloud workspace that lasts longer than your laptop, this is also straightforward.

What clouds does it work on?
------------
For now only [Google Compute Engine](https://cloud.google.com/products/compute-engine), but the code
is factored in such a way to make it easy to add other cloud providers.

Sounds great!  How do I use it?
------------

```
go get github.com/GoogleCloudPlatform/docker-cloud
```

If you don't already have a [Google Cloud Project](http://cloud.google.com), you can get one on the [Google Cloud Console](http://cloud.google.com/console)

Create a new **Client ID** for **Installed Application** in the APIs/credentials section.

```
docker-cloud auth -project <your-google-cloud-project-here> -id <your-credentials-client-id> -secret <your-credentials-secret>
# follow the instructions to authorize the client
```

Once the authorization is completed, you can start the proxy server. If you don't specify any project ID, it'll use the project ID provided during authorization.

```
docker-cloud start [-project=<your-google-cloud-project-here>]
```

### Connecting docker to the proxy ###
Use the `-H` flag on your docker client to connect to the proxy:
```
docker -H tcp://localhost:8080 run ehazlett/tomcat7
```



How can I contribute?
------------
I'm glad you asked.
### Getting the source ###
```
git clone https://github.com/GoogleCloudPlatform/docker-cloud.git
```

### Setting up Go ###
If you have not installed Go language yet, [install Go with this instruction](http://golang.org/doc/install).

Add a work directory for Go code, add it to `$GOPATH`, and add `$GOPATH/bin` to `$PATH`.

```
mkdir $HOME/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$PATH
```

If you haven't, you need to install [Mercurial (hg)](http://mercurial.selenic.com/) too.

### Building ###

```
go build
```

