from celery.task import task
from celerylibs import decorators
from subprocess import Popen
from tempfile import mkdtemp
from ubik import core as ubik_api
import logging, os, shutil

home_path = os.path.expanduser('~/')
backup_path = '%s/var/backups' % home_path

@task
@decorators.log_task
def mongo(host='localhost'):
	logger = logging.getLogger()
	logger.debug('Mongo Backup start:')
	logger.debug(' + Host  : %s' % host)
	logger.debug(' + Backup path: %s' % backup_path)

	archive_name = 'backup_mongodb'
	tmp_dir = '%s/%s' % (backup_path, archive_name)

	if os.path.exists(tmp_dir):
		logger.debug(' + Remove old temp dir')
		shutil.rmtree(tmp_dir)

	os.makedirs(tmp_dir)

	logger.debug(' + Launch mongodump')
	#mongodump_cmd = 'mongodump --host %s --out - | gzip -c -9 > %s/%s.gz' % (host, backup_path, archive_name)
	mongodump_cmd = 'mongodump --host %s --out %s/' % (host, tmp_dir)
	logger.debug(' + Command: %s' % mongodump_cmd)

	dump_output = Popen(mongodump_cmd, shell=True)
	dump_output.wait()

	logger.debug(' + Make archive')
	shutil.make_archive('%s/%s' % (backup_path, archive_name), 'zip', tmp_dir)	

	logger.debug(' + Remove temp dir')
	shutil.rmtree(tmp_dir)

	logger.debug('Mongo Backup finished')

@task
@decorators.log_task
def config():
	logger = logging.getLogger()
	logger.debug('Config Backup start:')

	archive_name = 'backup_config'
	tmp_dir = '%s/%s' % (backup_path, archive_name)

	if os.path.exists(tmp_dir):
		logger.debug(' + Remove old temp dir')
		shutil.rmtree(tmp_dir)

	logger.debug(' + List all packages')
	lines = []
	for package in ubik_api.db.get_installed():
		lines.append(package.name)
		lines.append('\n')
	lines = lines[:-1]
	f = open('%s/etc/.packages' % home_path, 'w')
	f.writelines(lines)
	f.close()

	logger.debug(' + Copy config files')
	shutil.copytree('%s/etc' % home_path, '%s/' % tmp_dir)
	
	logger.debug(' + Make archive')
	shutil.make_archive('%s/%s' % (backup_path, archive_name), 'zip', tmp_dir)	

	logger.debug(' + Remove temp dir')
	shutil.rmtree(tmp_dir)

	logger.debug('Config Backup finished')
