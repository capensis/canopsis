from celery.task import task
from celerylibs import decorators
from subprocess import Popen, PIPE
from tempfile import mkdtemp
import logging, os, shutil

@task
@decorators.log_task
def mongo(host='localhost', output='/opt/canopsis/var/backups'):
	logger = logging.getLogger()
	logger.debug('Mongo Backup start')
	logger.debug('Host  : %s' % host)
	logger.debug('Output: %s' % output)

	logger.debug('Create temp dir')
	archive_name = 'backup_mongodb'
	tmp_dir = mkdtemp(prefix='/opt/canopsis/tmp/')
	os.makedirs('%s/%s' % (tmp_dir, archive_name))

	logger.debug('Create output dir if not exists')
	if not os.path.exists(output):
		os.makedirs(output)	


	logger.debug('Launch mongodump')
	mongodump_cmd = '/opt/canopsis/bin/mongodump --host %s --out %s/%s' % (host, tmp_dir, archive_name)
	logger.debug('Command: %s' % mongodump_cmd)
	dump_output = Popen(mongodump_cmd, shell=True)
	dump_output.wait()

	logger.debug('Create archive into %s' % output)

	shutil.make_archive('%s/%s' % (output, archive_name),
						'zip',
						tmp_dir)

	logger.debug('Remove temp dir')
	shutil.rmtree(tmp_dir)
	
	logger.debug('Mongo Backup finished')

@task
@decorators.log_task
def config(output='/opt/canopsis/var/backups'):
	logger = logging.getLogger()
	logger.debug('Config Backup start')
	logger.debug('Output: %s' % output)
	
	logger.debug('Create output dir if not exists')
	if not os.path.exists(output):
		os.makedirs(output)	

	logger.debug('Create temp dir')
	archive_name = 'backup_config'
	tmp_dir = mkdtemp(prefix='/opt/canopsis/tmp/')

	logger.debug('Create file with installed packages')
	export_output = Popen('pkgmgr export', shell=True, stdout=PIPE)
	export_output.wait()
	f = open('/opt/canopsis/etc/.packages', 'w')
	f.writelines(export_output.stdout.read())
	f.close()

	logger.debug('Copy config files into tmp folder')
	shutil.copytree('/opt/canopsis/etc', '%s/%s' % (tmp_dir, archive_name))
	
	logger.debug('Create archive into %s' % output)
	logger.debug('Archive name: %s' % archive_name)
	shutil.make_archive('%s/%s' % (output, archive_name),
						'zip',
						tmp_dir)

	logger.debug('Remove temp dir')
	shutil.rmtree(tmp_dir)

	logger.debug('Config Backup finished')
