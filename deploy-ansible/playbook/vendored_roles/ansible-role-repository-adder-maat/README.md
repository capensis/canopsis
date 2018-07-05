# Ansible-role-repository-adder

Ajouter des dépôts à vos stations !

Requirements
------------

Aucun pré-requis

Role Variables
--------------

Ce rôle s'appuie sur un dictionnaire `repos` qui contient les variables propre à chaque repo que vous souhaitez ajouter.

Vous devez préciser les choses suivantes:
*  `name` : le nom du fichier `.repo` qui sera ajouter dans `/etc/yum.repos.d/`
*  `url`: dans ce champs se trouve soit : 
    - le lien vers le fichier rpm qui installe le dépôt
    - le nom d'un paquet qui installe un dépôt
    - le fichier `.repo` à copier dans le dossier `/etc/yum.repos.d/`
*  `content`: dans ce champs vous pouver insérer le contenu d'un fichier .repo
*  `gpg_key`: un lien vers la clef gpg à importer pour valider le contenu du dépôt.

Example Playbook
----------------
```yaml
- hosts: servers
  vars:
    repos:
      - name: epel
        url: "https://dl.fedoraproject.org/pub/epel/epel-release-latest-{{ ansible_distribution_major_version }}.noarch.rpm"
        gpg_key: "/etc/pki/rpm-gpg/RPM-GPG-KEY-EPEL-{{ ansible_distribution_major_version }}"
      - name: Remi
        url: "http://rpms.remirepo.net/enterprise/remi-release-{{ ansible_distribution_major_version }}.rpm"
        gpg_key: http://rpms.remirepo.net/RPM-GPG-KEY-remi
      - name: sclo
        url: centos-release-scl
        gpg_key: "https://www.centos.org/keys/RPM-GPG-KEY-CentOS-SIG-SCLo"
      - name: Atom
        content: |
          [Atom]
          name=Atom Editor
          baseurl=https://packagecloud.io/AtomEditor/atom/el/7/$basearch
          enabled=1
          gpgcheck=0
          repo_gpgcheck=1
          gpgkey=https://packagecloud.io/AtomEditor/atom/gpgkey
  roles:
     - role: ansible-role-repository-adder
```

Info
-------

Vous trouverez des exemples de dépôts dans le [wiki](https://gitlab.capensis.fr/capensis-ansible/ansible-role-repository-adder/wikis/home)

License
-------

GPLv3

Author Information
------------------

Paul MARCHAND ❤️ <pmarchand@capensis.fr>
