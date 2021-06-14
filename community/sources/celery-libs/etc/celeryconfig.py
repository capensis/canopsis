from celerylibs import listing

BROKER_HOST 			= "localhost"
BROKER_PORT 			= 5672
BROKER_USER 			= "guest"
BROKER_PASSWORD			= "guest"
BROKER_VHOST 			= "canopsis"
CELERY_RESULT_BACKEND		= "amqp"
CELERY_IMPORTS 			= listing.tasks('~/etc/tasks.d')

# informations here http://celery.github.com/celery/configuration.html#id1
CELERY_TASK_RESULT_EXPIRES	= 1800

CELERYD_LOG_LEVEL		= 'INFO'

CELERYD_TASK_TIME_LIMIT		= 1800
CELERYD_CONCURRENCY		= 5
