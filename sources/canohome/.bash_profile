# Get the aliases and functions
if [ -f ~/.bashrc ]; then
	. ~/.bashrc
fi

export PS1="[\u@\h \W]\$ "

export LD_LIBRARY_PATH="$HOME/lib:$LD_LIBRARY_PATH"
export PATH="$HOME/bin:$HOME/sbin:$PATH"
export PYTHONPATH="$HOME/lib/canolibs:$HOME/lib:$HOME/etc:$HOME/etc/tasks.d:$HOME/etc/tasks-cron.d"
export NODE_PATH="$HOME/lib/node_modules"
export UBIK_CONF="$HOME/etc/ubik.conf"

if [ -f ~/bin/activate ]; then
    source ~/bin/activate
fi
