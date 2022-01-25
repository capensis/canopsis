Name:		canopsis-core
Version:	CANOPSIS_PACKAGE_TAG
Release:	CANOPSIS_PACKAGE_REL%{?dist}
Summary:	Canopsis

Group:		Canopsis
License:	AGPLv3 (Canopsis Community)
URL:		https://git.canopsis.net/canopsis/canopsis-community

BuildRequires: rsync
Requires: bzip2 canopsis-engines-go canopsis-webui cyrus-sasl curl epel-release libacl libcurl libevent libffi librsync libxml2 libxslt net-snmp openldap openssl postgresql-libs python rsync sudo xmlsec1 xmlsec1-openssl zlib

%description
Canopsis Community main package.

%install
mkdir -p %{buildroot}/usr/lib/systemd/system
cp -ar /usr/lib/systemd/system/canopsis-* %{buildroot}/usr/lib/systemd/system/

mkdir -p %{buildroot}/opt
rsync -aKSH --delete /opt/canopsis %{buildroot}/opt/

mkdir -p %{buildroot}/usr/bin
ln -s /opt/canopsis/bin/canoctl %{buildroot}/usr/bin/canoctl

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
rm -rf %{buildroot}/opt/canopsis/bin/rabbit*
rm -rf %{buildroot}/opt/canopsis/Linux*

%pre
getent passwd canopsis > /dev/null 2>&1
if [ ! "$?" = "0" ]; then
	useradd -d /opt/canopsis -M -s /bin/bash canopsis
fi

%files
/opt/canopsis/venv-ansible
/usr/bin/canoctl

%defattr(0644, canopsis, canopsis, 0755)
/opt/canopsis/deploy-ansible
/opt/canopsis/include
/opt/canopsis/lib
/opt/canopsis/lib64
/opt/canopsis/share
/opt/canopsis/tmp
/opt/canopsis/var
/opt/canopsis/.vimrc
/opt/canopsis/VERSION.txt

%defattr(0755, canopsis, canopsis, 0755)
/opt/canopsis/bin/
/opt/canopsis/opt/mongodb/filldb.py

%attr(755, canopsis, canopsis) /opt/canopsis/deploy-ansible/install-self.sh
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
