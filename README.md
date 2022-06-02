# Velocity Automatically Built Docker Image

[![Build Status](https://github.drone.webzyno.com/api/badges/dockcenter/velocity/status.svg)](https://github.drone.webzyno.com/dockcenter/velocity)
[![GitHub](https://img.shields.io/github/license/dockcenter/velocity?color=informational)](https://github.com/dockcenter/velocity/blob/main/LICENSE)

This is a [Velocity](https://velocitypowered.com/) docker image with optimized Java flag provided by official [docs](https://velocitypowered.com/wiki/users/getting-started/).

We use [dedicated CI server](https://github.drone.webzyno.com/dockcenter/velocity) to track Velocity builds and automatically build Docker image.

## What is Velocity?

Velocity is a next-generation Minecraft proxy focused on scalability and flexibility.
It allows server owners to link together multiple Minecraft servers so they may appear as one.

### Blazing fast, extensible, and secure — choose three.

Velocity is blazing fast. 
Fast logins, fast server switches, optimizations to get the most from your server's hardware, and a focus on security means you can finally have plugins, a highly optimized proxy resilient to attacks, and protection against your backend servers being accessed by malicious users — no compromises required.

### Always there for you.
Velocity powers some of the world's largest Minecraft networks along with numerous small networks. 
Velocity can scale to thousands of players per proxy instance. Best of all, it works with Paper, Sponge, Forge, Fabric, and all versions of Minecraft from 1.7.2 to 1.18.1.

For more information, please reach to [Velocity official documentation](https://velocitypowered.com/wiki).

![Velocity](assets/velocity.png)

## How to use this image

### Start a Velocity server

With this image, you can create a new Velocity Minecraft proxy server with one command.
Here is an example:

```bash
sudo docker run -p 25565:25565 dockcenter/velocity
```

While this command will work just fine in many cases, it is only the bare minimum required to start a functional server and can be vastly improved by specifying some options.

## How to extend this image

There are many ways to extend the `dockcenter/velocity` image. Without trying to support every possible use case, here are just a few that we have found useful.

### Environment Variables

The `dockcenter/velocity` image uses several environment variables which are easy to miss.
`JAVA_MEMORY` environment variable is not required, but it is highly recommended to set an appropriate value according to your usage.

#### `JAVA_MEMORY`

This variable is not required, but is highly recommended.
By setting this value, you set the java `-Xms` and `Xmx` flag.
For more information about JVM memory size, refer to this [Oracle guide](https://docs.oracle.com/cd/E21764_01/web.1111/e13814/jvm_tuning.htm#PERFM160).

Default: `512M`

#### `JAVA_FLAGS`

This optional environment variable is used in conjunction with `JAVA_MEMORY` to provide additional java flag.
We use [Velocity officially recommended value](https://velocitypowered.com/wiki/users/getting-started/) as the default value.

Default: `-XX:+UseStringDeduplication -XX:+UseG1GC -XX:G1HeapRegionSize=4M -XX:+UnlockExperimentalVMOptions -XX:+ParallelRefProcEnabled -XX:+AlwaysPreTouch`

### Volume

The server data is stored in `/data` folder, and we create a volume for you.
To use your host directory to store data, please mount volume by adding the following options:

Using volume:
```bash
-v <my_volume_name>:/data
```

Using bind mount:
```bash
-v </path/to/folders>:/data
```

## LICENSE

Be careful using this container image as you must meet the obligations and conditions of the [GPLv3 ](https://github.com/PaperMC/Velocity/blob/dev/3.0.0/LICENSE) provided by the [Velocity](https://github.com/PaperMC/Velocity) development team.

The code for the [project](https://github.com/dockcenter/velocity) that builds the [`dockcenter/velocity`](https://hub.docker.com/r/dockcenter/velocity) image and pushes it to Docker Hub is distributed under the [MIT License](https://github.com/dockcenter/velocity/blob/main/LICENSE).

Please, don't confuse the two licenses.