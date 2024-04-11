%define debug_package %{nil} 
Name: canopsis
Version: %{version}
Release: 1%{?dist}
Summary: Canopsis
License: ASL 2.0
Source0: https://git.canopsis.net/canopsis/canopsis-pro/-/archive/%{version}/canopsis.tar.gz

BuildRequires: make >= 3.81, gcc, nodejs, yarn, systemd-rpm-macros

Conflicts: canopsis-pro

Prefix: /usr
Prefix: /etc
Prefix: /opt

%description
Canopsis Community RPM Package

%prep
%setup -q
GOLANG_VERSION=$(grep "^GOLANG_VERSION" community/.env |awk -F '=' '{print $NF}' | sed 's/ //g')
echo "install golang $GOLANG_VERSION."
wget https://go.dev/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
tar -C ~ -xzf go$GOLANG_VERSION.linux-amd64.tar.gz

%build
export PATH=$PATH:~/go/bin
make -C community/go-engines-community VERSION=%{version}
make -C community/sources/webcore/src/canopsis-next VERSION=%{version}

%install
make -C community/go-engines-community DESTDIR=%{buildroot} install
make -C community/go-engines-community DESTDIR=%{buildroot} systemd_install SYSTEMD_UNITS="engine-action engine-axe engine-che engine-fifo engine-pbehavior" SERVICES="canopsis-api"
make -C community/sources/webcore/src/canopsis-next DESTDIR=%{buildroot} install
make -C community/sources/webcore/src/canopsis-next DESTDIR=%{buildroot} nginx_config

%pre
getent passwd canopsis >/dev/null || \
  useradd -r -d /opt/canopsis -s /bin/bash -c "Canopsis user account" canopsis

%preun
%systemd_preun canopsis-engine-go@.service
if [ $1 -eq 0 ]; then
  systemctl disable canopsis-engine-go@
  systemctl stop canopsis-engine-go@
fi
%systemd_preun canopsis-service@canopsis-api.service
if [ $1 -eq 0 ]; then
  systemctl disable canopsis-service@canopsis-api.service
  systemctl stop canopsis-service@canopsis-api.service
fi

%files
%defattr(0755, canopsis, canopsis, 0755)
/opt/canopsis/bin
%defattr(0644, canopsis, canopsis, 0755)
%config(noreplace) /opt/canopsis/etc
%config(noreplace) /opt/canopsis/share/config
%attr(0640, root, canopsis) /opt/canopsis/etc/go-engines-vars.conf
/opt/canopsis/share/config
/opt/canopsis/share/database/fixtures/*yml
/opt/canopsis/share/database/migrations/*.js
/opt/canopsis/share/database/postgres_migrations/*.sql
/opt/canopsis/share/database/tech_postgres_migrations/*.sql
/opt/canopsis/var/lib/upload-files
%defattr(0644, root, root, 0755)
/usr/lib/systemd/system/canopsis-*
/usr/lib/systemd/system/canopsis.service

%post
echo "Modify /opt/canopsis/etc/go-engines-vars.conf to configure Canopsis"
echo "and run canopsis-reconfigure"
echo "After that, you can enable and start services"

%postun
if [ "$1" = "0" ]; then
  userdel canopsis
fi

%clean
make -C community/go-engines-community clean
make -C community/sources/webcore/src/canopsis-next clean

%package webui
Summary: Canopsis WebUI

Requires: nginx >= 1:1.20, nginx < 1:1.21

%description webui
Canopsis WebUI RPM Package

%pre webui

%preun webui

%post webui

%postun webui

%files webui
%defattr(0644, nginx, nginx, 0755)
/opt/canopsis/srv
%defattr(0644, nginx, nginx, 0755)
%config(noreplace) /etc/nginx/conf.d/default.conf
%config(noreplace) /etc/nginx/cors.inc
%config(noreplace) /etc/nginx/https.inc
%config(noreplace) /etc/nginx/rundeck.inc
%config(noreplace) /etc/nginx/resolvers.inc
