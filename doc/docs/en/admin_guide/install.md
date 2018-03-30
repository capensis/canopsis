# Install from sources

## Requirements

Install requirements with `root` user.

Canopsis can be installed on the following systems :

-   Debian 8, 9
-   RedHat / CentOS 7

### Debian Like:

```bash
apt-get update
apt-get install sudo git-core libcurl4-gnutls-dev libncurses5-dev
```

### Redhat Like:

Disable some services

:   We don't provide any SELinux context so it's better to disable it.
    Feel free to help us writing one ! You can see
    [Setup firewall page](firewall-rules.md) to configure Iptables

```bash
## Disable SELinux and Firewall
setenforce 0
chkconfig iptables off
chkconfig ip6tables off
chkconfig qpidd off
service iptables stop
service ip6tables stop
service qpidd stop
```

Iptables an qpidd may not be available on RedHat/CentOS 7. Take a look
at firewalld

Install some packages

```bash
yum install wget make redhat-lsb gcc gcc-c++ zlib-devel ncurses-devel git
```

## Download sources

Clone git repository:

```bash
git clone https://git.canopsis.net/canopsis/canopsis.git
cd canopsis
```

Available branches are:

-   master: the more stable version of canopsis
-   develop: futur released version, with incomplete features. This
    branch can carry bugs.

```bash
git submodule update --init
```

## Build and install

```bash
sudo ./build-install.sh
```

If build failed, you can see logs in `log/` directory.

Note that install dir will be /opt/canopsis by default. You can change
it by editing SOURCE_PATH/sources/canohome/lib/common.sh

## Start Canopsis

Log in `canopsis` and start it:

```bash
sudo su - canopsis
hypcontrol start
```

## Check installation

You can verify installation (in `canopsis` environment):

```bash
unittest.sh
```

And you can also check needed services (in `canopsis` environment):

```bash
hypcontrol status
```

## Troubleshooting

During some occasions, you could encounter some funny errors. Please
have a look at [Troubleshooting page](common-problems.md).


# Uninstall

In your Canopsis sources folder :

```bash
sudo ./build-install.sh -c
```
