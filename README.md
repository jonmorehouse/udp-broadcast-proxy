# UDP Broadcast Proxy

A proxy for broadcasting UDP packets to one ore more upstreams.

## Installation

**udp broadcast proxy** is available both as a binary as well as via `Docker` image

## Usage

```bash
udp-broadcast-proxy --upstreams=:5213 --listen-port=3122
```

## Inspiration

This was originally implemented to be used alongside of [statsd nsq](https://github.com/jonmorehouse/statsd-nsq) so a subset of **DataDog** statsd metrics could be archived as well as sent to **DataDog**.
