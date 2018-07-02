Name:		canopsis-core
Version:	CPS_PKG_TAG
Release:	CPS_PKG_REL%{?dist}
Summary:	Canopsis with CAT package

Group:		Canopsis
License:	AGPLv3 (Canopsis Core), Capensis, all rights reserved (CAT)
URL:		https://git.canopsis.net/canopsis/canopsis

BuildRequires: rsync
Requires: zlib libevent libcurl libtool openssl bzip2 cyrus-sasl openldap libcurl python openldap libxml2 libxslt rsync librsync libacl libxslt libffi xmlsec1 xmlsec1-openssl libtool net-snmp

%description
Canopsis with CAT package.

%install
mkdir -p %{buildroot}/usr/lib/systemd/system
cp -ar /usr/lib/systemd/system/canopsis-* %{buildroot}/usr/lib/systemd/system/

mkdir -p %{buildroot}/opt
rsync -aKSH --delete /opt/canopsis %{buildroot}/opt/

# ensure logs are clean
rm -rf %{buildroot}/opt/canopsis/var/log
mkdir -p %{buildroot}/opt/canopsis/var/log/engines

rm -rf %{buildroot}/opt/canopsis/.npm
rm -rf %{buildroot}/opt/canopsis/pip-selfcheck.json
rm -rf %{buildroot}/opt/canopsis/.ssh/*
rm -rf %{buildroot}/opt/canopsis/.ssh/.gitignore
rm -rf %{buildroot}/opt/canopsis/.viminfo
rm -rf %{buildroot}/opt/canopsis/.gitignore
rm -rf %{buildroot}/opt/canopsis/.rnd
rm -rf %{buildroot}/opt/canopsis/.erlang.cookie
rm -rf %{buildroot}/opt/canopsis/.config
rm -rf %{buildroot}/opt/canopsis/.cache
rm -rf %{buildroot}/opt/canopsis/.bash_history
rm -rf %{buildroot}/opt/canopsis/bin/influx*
rm -rf %{buildroot}/opt/canopsis/bin/rabbit*
rm -rf %{buildroot}/opt/canopsis/Linux*

%pre
getent passwd canopsis > /dev/null 2>&1
if [ ! "$?" = "0" ]; then
	useradd -d /opt/canopsis -M -s /bin/bash canopsis
fi

%files
%defattr(0644, canopsis, canopsis, 0755)
/opt/canopsis/include
/opt/canopsis/lib
/opt/canopsis/lib64
/opt/canopsis/share
/opt/canopsis/tmp
/opt/canopsis/var
/opt/canopsis/.smirc
/opt/canopsis/.vimrc

%defattr(0755, canopsis, canopsis, 0755)
/opt/canopsis/bin/
/opt/canopsis/opt/mongodb/filldb.py

%attr(755, canopsis, canopsis) /opt/canopsis/.bashrc
%attr(755, canopsis, canopsis) /opt/canopsis/.bash_profile

%dir %attr(0755, canopsis, canopsis) /opt/canopsis
%dir %attr(0755, canopsis, canopsis) /opt/canopsis/var/log
%dir %attr(0755, canopsis, canopsis) /opt/canopsis/var/log/engines

%config(noreplace) /opt/canopsis/opt/mongodb/load.d
%config(noreplace) /opt/canopsis/etc

%defattr(0644, root, root, 0755)
/usr/lib/systemd/system/canopsis-*

%changelog

