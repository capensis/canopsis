# `postgresql-timescaledb` Ansible role

This is a quick role which adds TimescaleDB to a CentOS 7 host.

Only tested with Ansible 2.8.7. No guarantee.

## Manual installation

Basically, this role is mainly an Ansible adaptation of the following procedure:

```sh
# Add PostgreSQL repos and GPG keys
yum install https://download.postgresql.org/pub/repos/yum/reporpms/EL-7-x86_64/pgdg-redhat-repo-latest.noarch.rpm

# Add TimescaleDB GPG key
cd /etc/pki/rpm-gpg/ && curl -L -o RPM-GPG-KEY-PKGCLOUD-TIMESCALEDB https://packagecloud.io/timescale/timescaledb/gpgkey

# Add TimescaleDB repo
cat > /etc/yum.repos.d/timescale_timescaledb.repo << EOF
[timescale_timescaledb]
name=timescale_timescaledb
baseurl=https://packagecloud.io/timescale/timescaledb/el/\$releasever/\$basearch
repo_gpgcheck=1
# timescaledb doesn't sign all its packages
gpgcheck=0
enabled=1
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-PKGCLOUD-TIMESCALEDB
sslverify=1
sslcacert=/etc/pki/tls/certs/ca-bundle.crt
metadata_expire=300
EOF

# Install TimescaleDB and its PostgreSQL dependencies
yum makecache -y
yum install timescaledb-2-postgresql-13 -y

# Initialize PostgreSQL
postgresql-13-setup initdb

# Let Timescaledb tune itself
# note: the -yes and -pg-config flags are buggy/incomplete, hence the yes(1)/PATH hacks
yes | PATH=/usr/pgsql-13/bin:$PATH timescaledb-tune -color false -pg-config /usr/pgsql-13/bin/pg_config -out-path /var/lib/pgsql/13/data/postgresql.conf -yes

# Disable telemetry by default...
echo "timescaledb.telemetry_level=off" >> /var/lib/pgsql/13/data/postgresql.conf

# Enable and start the service
systemctl enable postgresql-13
systemctl start postgresql-13
```

Main references in building this:

* <https://docs.timescale.com/install/latest/self-hosted/installation-redhat/#setting-up-the-timescaledb-extension> (note that it's quite incomplete as of 2022-01-20)
* <https://github.com/timescale/timescaledb-docker/blob/master/docker-entrypoint-initdb.d/000_install_timescaledb.sh>
* <https://www.digitalocean.com/community/tutorials/how-to-install-and-use-timescaledb-on-centos-7>

## License

© 2022 Capensis, licensed under the AGPLv3.
