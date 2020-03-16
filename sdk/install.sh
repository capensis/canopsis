#!/usr/bin/env bash

set -u

# WARNING: REQUIRES /bin/sh
#
# Copyright:: Copyright (c) 2019 Capensis
# License:: CC BY-SA 3.0
#
#
#                   Vous êtes autorisé à :
#
#Partager — copier, distribuer et communiquer le matériel par tous moyens 
#            et sous tous formats 
#
#
# Adapter — remixer, transformer et créer à partir du matériel
#           pour toute utilisation, y compris commerciale. 
#
#
#

CPS_SRC="$( cd "$( dirname "${BASH_SOURCE[0]}" )/../" >/dev/null 2>&1 && pwd )"
CPS_DEPS_ROOT="${CPS_SRC}/deps"
CPS_HOME=/opt/canopsis
CPS_SDK="${CPS_SRC}/sdk"
CPS_SDK_BIN="${CPS_SDK}/bin"

# Color values
COL_NC='\e[0m' # No Color
COL_LIGHT_GREEN='\e[1;32m'
COL_LIGHT_RED='\e[1;31m'
TICK="[${COL_LIGHT_GREEN}✓${COL_NC}]"
CROSS="[${COL_LIGHT_RED}✗${COL_NC}]"
INFO="[i]"
# shellcheck disable=SC2034
DONE="${COL_LIGHT_GREEN} done!${COL_NC}"
OVER="\\r\\033[K"
C_MODE=0

# Default Canopsis repository
branch="develop"
project="canopsis"
canopsis_repositories="canopsis-connectors/connector-send_event2canopsis.git canopsis-connectors/connector-libs.git canopsis/go-engines.git"
canopsis_repositories_cat="cat/canopsis-cat.git cat/connector-email2canopsis.git cat/go-engines-cat.git cat/connector-snmp2canopsis.git cat/itop2canopsis.git canopsis-connectors/connector-shinkenenterprise2canopsis.git "
https_url="https://"
ssh_url="ssh://git@"
gitlab_url="git.canopsis.net"
go_version="1.12.9"

debug_mode=0

# Check if command exists - return 0 if it exists, else 1 if not exists 
exists() {
  if command -v $1 >/dev/null 2>&1
  then
    return 0
  else
    return 1
  fi
}

# Display message to report bug
report_bug() {
  echo "Version: $version"
  echo ""
  echo "Please file a Bug Report at https://git.canopsis.net/canopsis/canopsis/issues/new"
  echo "Alternatively, feel free to open a Support Ticket at https://community.capensis.org"
  echo ""
  echo "Please include as many details about the problem as possible i.e., how to reproduce"
  echo "the problem (if possible), type of the Operating System and its version, etc.,"
  echo "and any other relevant details that might help us with troubleshooting."
  echo ""
}


# Install Docker
common_docker_linux() {
 if ! exists "docker"
 then
     echo "docker stuff with curl..."
     curl -fsSL https://download.docker.com/linux/$platform/gpg | sudo apt-key add -
     detect_error "error in download docker l 62 in install.sh"
     curl -fsSL get.docker.com | CHANNEL=stable sudo -E /bin/sh
     detect_error "error in download docker l 64 in install.sh"
 fi
 
 is_docker_group=$(id ${SUDO_USER}|grep "\(docker\)")
 if test $? -ne 0
 then
   #usermod -aG docker $SUDO_USER
   #newgrp docker
   sudo gpasswd -a $USER docker
 fi
   
 
 sudo systemctl start docker
 detect_error "error in docker starting"
 sudo systemctl enable docker
 if ! exists "docker-compose"
 then
     sudo curl -L "https://github.com/docker/compose/releases/download/1.24.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
     detect_error "error in download docker compose docker-compose-$(uname -s)-$(uname -m) l 81"
     sudo chmod +x /usr/local/bin/docker-compose
 fi
}



# Init symlinks
init_sdk(){
    # Shell Env vars
    echo 'PATH="${GOROOT:-/usr/local/go}/bin:${HOME}/golibs:'"${CPS_SDK_BIN}"':${PATH}"' | sudo tee /etc/profile.d/canopsis-sdk.sh > /dev/null
    detect_error "error during profile.d installation"

    # Auto-completion
    sudo cp ${CPS_SDK_BIN}/canopsis_completion /etc/bash_completion.d/
    detect_error "error in bash_completion installation"
}

init_systemd(){
  # Allow service to run when shell exited
  #loginctl enable-linger ${SUDO_USER}
  sudo loginctl enable-linger ${USER}
  detect_error "error during loginctl enable-linger command"
}


# install_dependencies TYPE
# TYPE is distrib "CentOS", "ubuntu", etc.
install_dependencies() {
  echo "Installing $project dependencies"
  case "$1" in
    CentOS | RedHat)
      echo "with yum..."
      detect_error "error during $2 installation"
      ;;
    ubuntu | debian)
      echo "with apt..."
      #dpkg -i "$2"

      sudo apt-get update > /dev/null 2>&1
      DEBIAN_FRONTEND=noninteractive sudo -E apt-get -yq \
      -o Dpkg::Options::=--force-confold \
      -o Dpkg::Options::=--force-confdef \
      --allow-downgrades --allow-remove-essential --allow-change-held-packages \
      install \
        git \
        apt-transport-https \
        ca-certificates \
        curl \
        software-properties-common \
        python-dev \
        python-pip \
        python-virtualenv \
        python3 \
        python3-venv \
        locales-all \
        libsasl2-dev \
        libxml2-dev \
        libxmlsec1-dev \
        libssl-dev \
        python-ldap \
        libpq-dev \
        libldap2-dev

        detect_error "error in install dependencies"
      ;;
    "bff")
      echo "installing with installp..."
      installp -aXYgd "$2" all
      detect_error "error in installp $2"
      ;;
    "solaris")
      echo "installing with pkgadd..."
      echo "conflict=nocheck" > $tmp_dir/nocheck
      echo "action=nocheck" >> $tmp_dir/nocheck
      echo "mail=" >> $tmp_dir/nocheck
      pkgrm -a $tmp_dir/nocheck -n $project >/dev/null 2>&1 || true
      pkgadd -G -n -d "$2" -a $tmp_dir/nocheck $project
      ;;
    "pkg")
      echo "installing with installer..."
      cd / && /usr/sbin/installer -pkg "$2" -target /
      ;;
    "dmg")
      echo "installing dmg file..."
      hdiutil detach "/Volumes/chef_software" >/dev/null 2>&1 || true
      hdiutil attach "$2" -mountpoint "/Volumes/chef_software"
      cd / && /usr/sbin/installer -pkg `find "/Volumes/chef_software" -name \*.pkg` -target /
      hdiutil detach "/Volumes/chef_software"
      ;;
    "sh" )
      echo "installing with sh..."
      sh "$2"
      ;;
    "p5p" )
      echo "installing p5p package..."
      pkg install -g "$2" $project
      detect_error "error in install $2"
      ;;
    *)
      echo "Unknown OS: $1"
      report_bug
      exit 1
      ;;
  esac
  if test $? -ne 0; then
    echo "Installation failed"
    report_bug
    exit 1
  fi
}



# installation des sources canopsis
install_sources() {

  mkdir -p "${CPS_DEPS_ROOT}"

  # Canopsis CORE deps
  for repository in ${canopsis_repositories}
    do
      dst_dist=$(echo "${repository}" | sed -e 's#^[^/]*/\([^/]*\).git$#\1#')
      dst_namespace=$(echo "${repository}" | sed -e 's#^\([^/]*\)/\([^/]*\).git$#\1#')
      if test ! -d ${CPS_DEPS_ROOT}/${dst_namespace}/${dst_dist}
      then
          git clone ${ssh_url}${gitlab_url}/${repository} ${CPS_DEPS_ROOT}/${dst_namespace}/${dst_dist} -b ${branch}
          detect_error "error in git clone ${gitlab_url}/${repository} ${CPS_DEPS_ROOT}/${dst_namespace}/${dst_dist}"
      fi
    done

  #for CAT repo
  if [ $C_MODE = 1 ];then
    for repository in ${canopsis_repositories_cat}
    do
      dst_dist=$(echo "${repository}" | sed -e 's#^[^/]*/\([^/]*\).git$#\1#')
      dst_namespace=$(echo "${repository}" | sed -e 's#^\([^/]*\)/\([^/]*\).git$#\1#')
      if test ! -d ${CPS_DEPS_ROOT}/${dst_namespace}/${dst_dist}
      then
        git clone ${ssh_url}${gitlab_url}/${repository} ${CPS_DEPS_ROOT}/${dst_namespace}/${dst_dist} -b ${branch}
        detect_error "error in git clone ${gitlab_url}/${repository} ${CPS_DEPS_ROOT}/${dst_namespace}/${dst_dist}"
      fi
  done
  fi

}


create_cps_run_dir() {
  sudo mkdir -p ${CPS_HOME}
  detect_error "error during ${CPS_HOME} creation"
  sudo -E chown ${USER}:${USER} ${CPS_HOME}
  detect_error "error during ${CPS_HOME} chown"
}


# Instal Go Env
install_golang () {

  go_package=go${go_version}.linux-amd64.tar.gz
  if ! exists "go"
  then
      curl https://storage.googleapis.com/golang/${go_package} | sudo tar -C /usr/local -xz
      detect_error "error during GO install"
      export PATH=/usr/local/go/bin:$PATH
  fi

}

# install Node.js Env
install_nodejs () {
  curl -sL https://deb.nodesource.com/setup_10.x | sudo -E bash -
  sudo apt-get install -y nodejs
  detect_error "error in Node.js install"

  curl -sL https://dl.yarnpkg.com/debian/pubkey.gpg | sudo apt-key add -
  echo "deb https://dl.yarnpkg.com/debian/ stable main" | sudo tee /etc/apt/sources.list.d/yarn.list
  
  sudo apt-get update && sudo apt-get install yarn
  detect_error "error in yarn install"
}


disable_output () {
  exec 3<&1
  exec 1>/dev/null
  exec 4<&2
  exec 2>/dev/null
}

enable_output () {
  exec 1<&3
}

enable_stderr () {
  exec 2<&4
}

# $1 message string
# $2 debug mode
print_message_val() {
  if [ $debug_mode = 0 ]; then
    enable_output
  fi
  if test -z ${2:-} || test ${2:-} -ne 1; then
    printf "%b  %b %s\\n" "${OVER}" "${TICK}" "$1"    
  fi
  if [ $debug_mode = 0 ]; then
    disable_output
  fi
}


print_message_err() {
  if [ $debug_mode = 0 ]; then
    enable_output
  fi
  
  if test -z ${2:-} || test ${2:-} -ne 1; then
    printf "%b  %b %s\\n" "${OVER}" "${CROSS}  ${COL_LIGHT_RED}" "$1"    
  fi
  if [ $debug_mode = 0 ]; then
    disable_output
  fi
}


detect_error(){
  if [ $? -ne 0 ] ; then
    print_message_err "$1"
    exit 1
  fi
}


# MAIN SCRIPT
##################
# script_cli_parameters.sh
############
# This section reads the CLI parameters for the install script and translates
#   them to the local parameters to be used later by the script.
#
# Outputs:
# $version: Requested version to be installed.
# $branch: Channel to install the product from
# $project: Project to be installed
# $cmdline_filename: Name of the package downloaded on local disk.
# $cmdline_dl_dir: Name of the directory downloaded package will be saved to on local disk.
# $install_strategy: Method of package installations. default strategy is to always install upon exec. Set to "once" to skip if project is installed
# $download_url_override: Install package downloaded from a direct URL.
# $checksum: SHA256 for download_url_override file (optional)
############


while getopts dv:f:P:d:s:l:ac opt
do
  case "$opt" in

    v)  version="$OPTARG";;
    b)  branch="$OPTARG";;
    f)  cmdline_filename="$OPTARG";;
    P)  project="$OPTARG";;
    s)  install_strategy="$OPTARG";;
    l)  download_url_override="$OPTARG";;
    a)  checksum="$OPTARG";;
    d)  debug_mode=1;;
    c)  C_MODE=1;;
    \?)   # unknown flag
      echo >&2 \
      " ${INFO} ${COL_LIGHT_RED} usage: $0 [-d] [-P project] [-b release_branch] [-v version] [-f filename | -d download_dir] [-s install_strategy] [-l download_url_override] [-a checksum]"
      exit 1;;
  esac
done

shift `expr $OPTIND - 1`

printf "%b  %b %s\\n" "${INFO}" "${COL_LIGHT_GREEN}" "Canopsis SDK Installation starting"

if test ${debug_mode} -eq 1; then
  set -x
else
  disable_output
fi

if test -d "/opt/$project" && test "x${install_strategy:-}" = "xonce"; then
  echo "$project installation detected"
  echo "install_strategy set to 'once'"
  echo "Nothing to install"
  exit
fi


# platform_detection.sh
############
# This section makes platform detection compatible with omnitruck on the system
#   it runs.
#
# Outputs:
# $platform: Name of the platform.
# $platform_version: Version of the platform.
# $machine: System's architecture.
############

#
# Platform and Platform Version detection
#
# NOTE: This should now match ohai platform and platform_version matching.
# do not invented new platform and platform_version schemas, just make this behave
# like what ohai returns as platform and platform_version for the server.
#
# ALSO NOTE: Do not mangle platform or platform_version here.  It is less error
# prone and more future-proof to do that in the server, and then all omnitruck clients
# will 'inherit' the changes (install.sh is not the only client of the omnitruck
# endpoint out there).
#

machine=`uname -m`
os=`uname -s`

if test -f "/etc/lsb-release" && grep -q DISTRIB_ID /etc/lsb-release && ! grep -q wrlinux /etc/lsb-release; then
  platform=`grep DISTRIB_ID /etc/lsb-release | cut -d "=" -f 2 | tr '[A-Z]' '[a-z]'`
  platform_version=`grep DISTRIB_RELEASE /etc/lsb-release | cut -d "=" -f 2`

  if test "$platform" = "\"cumulus linux\""; then
    platform="cumulus_linux"
  elif test "$platform" = "\"cumulus networks\""; then
    platform="cumulus_networks"
  fi

elif test -f "/etc/debian_version"; then
  platform="debian"
  platform_version=`cat /etc/debian_version`
elif test -f "/etc/Eos-release"; then
  # EOS may also contain /etc/redhat-release so this check must come first.
  platform=arista_eos
  platform_version=`awk '{print $4}' /etc/Eos-release`
  machine="i386"
elif test -f "/etc/redhat-release"; then
  platform=`sed 's/^\(.\+\) release.*/\1/' /etc/redhat-release | tr '[A-Z]' '[a-z]'`
  platform_version=`sed 's/^.\+ release \([.0-9]\+\).*/\1/' /etc/redhat-release`

  if test "$platform" = "xenserver"; then
    # Current XenServer 6.2 is based on CentOS 5, platform is not reset to "el" server should hanlde response
    platform="xenserver"
  else
    # FIXME: use "redhat"
    platform="el"
  fi

elif test -f "/etc/system-release"; then
  platform=`sed 's/^\(.\+\) release.\+/\1/' /etc/system-release | tr '[A-Z]' '[a-z]'`
  platform_version=`sed 's/^.\+ release \([.0-9]\+\).*/\1/' /etc/system-release | tr '[A-Z]' '[a-z]'`
  case $platform in amazon*) # sh compat method of checking for a substring
    platform="el"

    . /etc/os-release
    platform_version=$VERSION_ID
    if test "$platform_version" = "2"; then
      platform_version="7"
    else
      # VERSION_ID will match YYYY.MM for Amazon Linux AMIs
      platform_version="6"
    fi
  esac

# Apple OS X
elif test -f "/usr/bin/sw_vers"; then
  platform="mac_os_x"
  # Matching the tab-space with sed is error-prone
  platform_version=`sw_vers | awk '/^ProductVersion:/ { print $2 }' | cut -d. -f1,2`

  # x86_64 Apple hardware often runs 32-bit kernels (see OHAI-63)
  x86_64=`sysctl -n hw.optional.x86_64`
  if test $x86_64 -eq 1; then
    machine="x86_64"
  fi
elif test -f "/etc/release"; then
  machine=`/usr/bin/uname -p`
  if grep -q SmartOS /etc/release; then
    platform="smartos"
    platform_version=`grep ^Image /etc/product | awk '{ print $3 }'`
  else
    platform="solaris2"
    platform_version=`/usr/bin/uname -r`
  fi
elif test -f "/etc/SuSE-release"; then
  if grep -q 'Enterprise' /etc/SuSE-release;
  then
      platform="sles"
      platform_version=`awk '/^VERSION/ {V = $3}; /^PATCHLEVEL/ {P = $3}; END {print V "." P}' /etc/SuSE-release`
  else
      platform="suse"
      platform_version=`awk '/^VERSION =/ { print $3 }' /etc/SuSE-release`
  fi
elif test "x$os" = "xFreeBSD"; then
  platform="freebsd"
  platform_version=`uname -r | sed 's/-.*//'`
elif test "x$os" = "xAIX"; then
  platform="aix"
  platform_version="`uname -v`.`uname -r`"
  machine="powerpc"
elif test -f "/etc/os-release"; then
  . /etc/os-release
  if test "x$CISCO_RELEASE_INFO" != "x"; then
    . $CISCO_RELEASE_INFO
  fi

  platform=$ID
  platform_version=$VERSION
fi

if test "x$platform" = "x"; then
  echo "Unable to determine platform version!"
  report_bug
  exit 1
fi


#
# NOTE: platform manging in the install.sh is DEPRECATED
#
# - install.sh should be true to ohai and should not remap
#   platform or platform versions.
#
# - remapping platform and mangling platform version numbers is
#   now the complete responsibility of the server-side endpoints
#

major_version=`echo $platform_version | cut -d. -f1`
case $platform in
  # FIXME: should remove this case statement completely
  "el")
    # FIXME:  "el" is deprecated, should use "redhat"
    platform_version=$major_version
    ;;
  "debian")
    if test "x$major_version" = "x5"; then
      # This is here for potential back-compat.
      # We do not have 5 in versions we publish for anymore but we
      # might have it for earlier versions.
      platform_version="6"
    else
      platform_version=$major_version
    fi
    ;;
  "freebsd")
    platform_version=$major_version
    ;;
  "sles")
    platform_version=$major_version
    ;;
  "suse")
    platform_version=$major_version
    ;;
esac

# normalize the architecture we detected
case $machine in
  "x86_64"|"amd64"|"x64")
    machine="x86_64"
    ;;
  "i386"|"i86pc"|"x86"|"i686")
    machine="i386"
    ;;
  "sparc"|"sun4u"|"sun4v")
    machine="sparc"
    ;;
esac

if test "x$platform_version" = "x"; then
  echo "Unable to determine platform version!"
  report_bug
  exit 1
fi

if test "x$platform" = "xsolaris2"; then
  # hack up the path on Solaris to find wget, pkgadd
  PATH=/usr/sfw/bin:/usr/sbin:$PATH
  export PATH
fi

printf "%b  %b %s\\n" "${INFO}" "${COL_LIGHT_GREEN}" "Current System Detected: $platform $platform_version $machine"

############
# end of platform_detection.sh
############


# All of the download utilities in this script load common proxy env vars.
# If variables are set they will override any existing env vars.
# Otherwise, default proxy env vars will be loaded by the respective
# download utility.

if test "x${https_proxy:-}" != "x"; then
  echo "setting https_proxy: $https_proxy"
  export HTTPS_PROXY=$https_proxy
  export https_proxy=$https_proxy
fi

if test "x${http_proxy:-}" != "x"; then
  echo "setting http_proxy: $http_proxy"
  export HTTP_PROXY=$http_proxy
  export http_proxy=$http_proxy
fi

if test "x${ftp_proxy:-}" != "x"; then
  echo "setting ftp_proxy: $ftp_proxy"
  export FTP_PROXY=$ftp_proxy
  export ftp_proxy=$ftp_proxy
fi

if test "x${no_proxy:-}" != "x"; then
  echo "setting no_proxy: $no_proxy"
  export NO_PROXY=$no_proxy
  export no_proxy=$no_proxy
fi

create_cps_run_dir
print_message_val "'${CPS_HOME}' created"

install_dependencies $platform
print_message_val "Packages dependencies installed"

install_sources
print_message_val "All sources retrieved from GIT repositories"

common_docker_linux
print_message_val "Docker stuff installed"

install_python_virtualenv
print_message_val "Python stuff installed"

install_golang
print_message_val "Golang stuff installed"

install_nodejs
print_message_val "Nodejs stuff installed"

init_sdk
init_systemd
print_message_val "Init successfull"

print_message_val "Canopsis SDK successfully installed"
print_message_val "Restart your terminal to Reload your PATH"

############
# end of install.sh
############
