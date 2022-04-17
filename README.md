Aucht
===========

This repository provides an easy way to switch between multiple guest virtual machines when your host is headless.

It provides:

  * A way to start/suspend/kill libvirt domains (virtual machines)

Table of Contents:

  * [About](#about)
  * [Installing and Compiling from Source](#installing-and-compiling-from-source)
  * [Contributing](#contributing)
  * [License](#license)

About
-----

Aucht was born out of a need to switch between multiple guest VM's when doing single GPU passthrough using a headless host. This is useful if you want to use several Windows/Mac/Linux/Other VMs with GPU passthrough on linux but don't want to mess around with binding/unbinding GPU drivers and such.

Installing and Compiling from Source
------------

In order to compile Aucht from source, you'll need the following:

  * [Go](https://golang.org) installed and [configured](https://golang.org/doc/install)
  * A system setup with [QEMU](https://www.qemu.org/)/[libvirt](https://libvirt.org/)
  * At least one VM, preferrably two or more
  * A little patience :)

Contributing
------------

Contributions are always welcome. If you're interested in contributing, send me an email or submit a PR.

License
-------

This project is currently licensed under GPLv3. This means you may use our source for your own project, so long as it remains open source and is licensed under GPLv3.

Please refer to the [license](/LICENSE) file for more information.