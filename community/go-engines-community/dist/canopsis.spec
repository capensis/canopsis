%define debug_package %{nil} 
Name: canopsis
Version: %{version}
Release: 3%{?dist}
Summary: Canopsis
License: ASL 2.0
Source0: https://git.canopsis.net/canopsis/canopsis-pro/-/archive/%{version}/canopsis.tar.gz

BuildRequires: make >= 3.81, gcc, nodejs = 2:14.18.3, yarn, systemd-rpm-macros

Conflicts: canopsis-pro

%description
Canopsis Community RPM Package

%prep
%setup -q
echo "install golang"
wget https://go.dev/dl/go1.18.6.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.18.6.linux-amd64.tar.gz

%build
export PATH=$PATH:/usr/local/go/bin
make -C community/go-engines-community VERSION=%{version}
make -C community/sources/webcore/src/canopsis-next VERSION=%{version}

%install
make -C community/go-engines-community DESTDIR=%{buildroot} install
make -C community/go-engines-community DESTDIR=%{buildroot} systemd_install SYSTEMD_UNITS="engine-action engine-axe engine-che engine-fifo engine-pbehavior engine-service" SERVICES="canopsis-api"
make -C community/sources/webcore/src/canopsis-next DESTDIR=%{buildroot} install
make -C community/go-engines-community DESTDIR=%{buildroot} nginx_config

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
/opt/canopsis/share/config
/opt/canopsis/share/database/fixtures/*yml
/opt/canopsis/share/database/migrations/*.js
/opt/canopsis/share/database/postgres_migrations/*.sql
/opt/canopsis/share/database/tech_postgres_migrations/*.sql
%defattr(0644, root, root, 0755)
/usr/lib/systemd/system/canopsis-*
/usr/lib/systemd/system/canopsis.service

%post
echo "Modify /opt/canopsis/etc/go-engines-vars.conf to configure Canopsis"
echo "and run canopsis-reconfigure"
echo "After that, you can enable and start services"

%postun
userdel canopsis

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
