from celery.task import task
from celerylibs import decorators

from caccount import caccount

@task
@decorators.log_task
def hostname(account=None):
	import socket
	return socket.gethostname()
