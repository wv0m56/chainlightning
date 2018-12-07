# Chainlightning

Disclaimer: probably alpha software.

Chainlightning is a cache server which is:

* chainable: you can put multiple layers of chainlighting servers in front of your origin server, where each one will act as an origin and is able to relay expiry values to the layer in front of it.
* simplified: you don't write your own get and set logic.
* origin controlled: control the lifetime of cached values from the origin server, and dozens or hundreds of chainlighting servers will invalidate your cached values when it is time, without you having to write further logic.
* (todo) purgable: manually purge cached values if your business logic requires it.

If everything works as intended, its mechanics should approximate the one described in this [video](https://www.facebook.com/Engineering/videos/live-video-solutions:-solving-the/10153675295382200/).

## Example

A sample origin server is provided in [samplebackend/timestampserver](https://github.com/wv0m56/chainlightning/tree/master/samplebackend/timestampserver). Build and run it.

Next, build chainlightning itself and run it with the default configuration:
```
chainlightning -c config.toml.example
```

Curl to chainlightning multiple times at random intervals and observe how it and the origin server behave:
```
curl localhost:11000/timeout/11111 -v
```