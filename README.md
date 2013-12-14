Docker Cloud
============

What is it?
------------
Docker Cloud is a proxy for the Docker API which automatically creates and destroys cloud virtual machines to run
your docker containers.

Why would I want to do that?
------------
If you are running Docker on OS X or Windows, there is no longer any need to install a virtualization layer like
vagrant on your machine.  You can simply run it in the cloud.

What clouds does it work on?
------------
For now only <a href="https://cloud.google.com/products/compute-engine/">Google Compute Engine</a>, but the code
is factored in such a way to make it easy to add other cloud providers.

Sounds great!  How do I use it?
------------
I'm glad you asked.
### Getting the source ###
<code>
git clone https://github.com/brendandburns/docker-cloud.git
</code>

### Building ###
<code>
<pre>
cd docker-cloud
./build
</code>

### Running the proxy ###
There are different instructions for different cloud providers.

#### Google Compute Engine ####
If you don't already have a Google Cloud Project, you can get one on the <a href="http://cloud.google.com/console">Google Cloud Console</a>

<code>
./docker-proxy --project <your-google-cloud-project-here>
</code>

### Connecting docker to the proxy ###
Use the "-H" flag on your docker client to connect to the proxy:
<code>
docker -H tcp://localhost:8080 run ehazlett/tomcat7
</code>