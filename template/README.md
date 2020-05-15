# Installation

Download the {{ Type }} runner and configure to connect with your central Drone server using your server address and shared secret.

```console
$ docker run -d \
  --publish=3000:3000 \
  --env=DRONE_DEBUG=true \
  --env=DRONE_RPC_HOST=drone.company.com \
  --env=DRONE_RPC_PROTO=https \
  --env=DRONE_RPC_SECRET=${SECRET} \
  --restart=always \
  --name=drone-runner-{{Type}} {{Image}}
```

# License

This software is licensed under the [Blue Oak Model License 1.0.0](https://spdx.org/licenses/BlueOak-1.0.0.html).

# Notice
<!-- do not remove notice -->

This software relies on the [runner-go](https://github.com/drone/runner-go) module authored by [Drone.IO, Inc](https://github.com/drone) under a [Non-compete](https://github.com/drone/runner-go/blob/master/LICENSE.md) license. This module can be used in any software for free for any permitted purpose in accordance with the license. This module cannot be used is any software that competes with the licensor.
