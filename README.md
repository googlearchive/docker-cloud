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
I'm glad you asked.
### Getting the source ###
```
git clone https://github.com/brendandburns/docker-cloud.git
```

### Building ###
```
go get github.com/brendandburns/docker-cloud
```

### Running the proxy ###
There are different instructions for different cloud providers.

#### Google Compute Engine ####
If you don't already have a [Google Cloud Project](http://cloud.google.com), you can get one on the [Google Cloud Console](http://cloud.google.com/console)

Assuming `$GOPATH/bin` is in your `$PATH`:
```
docker-cloud --project <your-google-cloud-project-here>
```

The first time you run the proxy, you will receive a URL and be prompted to visit it to obtain an
authorization code.  Once you do this, run the proxy again:

```
docker-cloud --project <your-project-here> --code <auth-code-here>
```

The code will be cached, and you should never have to do go through that step again.

### Connecting docker to the proxy ###
Use the "-H" flag on your docker client to connect to the proxy:
```
docker -H tcp://localhost:8080 run ehazlett/tomcat7
```

