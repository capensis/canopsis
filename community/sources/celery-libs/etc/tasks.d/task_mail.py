from celery.task import task
from caccount import caccount
from cstorage import cstorage
from cfile import cfile
from celerylibs import decorators
import logging

import re
import string
import smtplib
import socket
import sys
import ConfigParser
import os

from email import Encoders
from email.MIMEBase import MIMEBase
from email.MIMEText import MIMEText
from email.MIMEMultipart import MIMEMultipart
from email.Utils import formatdate

logger = logging.getLogger('Mail Task')

@task
@decorators.log_task
def send(account=None, recipients=None, subject=None, body=None, attachments=None, smtp_host=None, smtp_port=None, html=False):

	if smtp_port == None and smtp_host == None:

		logger.info('Will read mail configuration from ~/etc/mail.conf')
		config = ConfigParser.RawConfigParser()
		config.read(os.path.expanduser('~/etc/mail.conf'))

		smtp_host = 'localhost'
		smtp_port = 25
		try:
			smtp_host = config.get('master', 'host')
			smtp_port = int(config.getint('master', 'port'))
		except Exception as e:
			logger.error('Unable to read configuration for mail from file ~/etc/mail.conf')

	"""
		account		: caccount or nothing for anon
		recipients	: str("glehee@capensis.fr"), caccount
					  list of (caccount or string)
		subject		: str("My Subject")
		body		: str("My Body")
		attachments	: cfile, list of cfile
		smtp_host	: str("localhost")
		smtp_port	: int(25)
		html		: allow html into mail body (booleen)
	"""
	###########
	# Account #
	###########
	# Defaults
	account_firstname = ""
	account_lastname = ""
	account_mail = ""
	account_full_mail = ""
	
	if isinstance(account, caccount):
		account_firstname = account.firstname
		account_lastname = account.lastname
		account_mail = account.mail
		if not account_mail:
			logger.info('No mail adress for this user (Fill the mail account field)')
			account_mail = '%s@%s' % (account.user,socket.gethostname())

		if isinstance(account_mail, (list, tuple)):
			account_mail = account_mail[0]

		if not account_lastname and not account_firstname:
			account_full_mail = "\"%s\" <%s>" % (account_mail.split('@')[0].title(), account_mail)	
		else:
			account_full_mail = account.get_full_mail()
		if not re.match("^[a-zA-Z0-9._%-]+@[a-zA-Z0-9._%-]+.([a-zA-Z]{2,6})?$", str(account_mail)):
			raise ValueError('Invalid Email format for sender')
	else:
		raise Exception('Need caccount object in account')
	
	##############
	# Recipients #
	##############
	if not recipients:
		raise ValueError('Give at least one recipient')

	if not isinstance(recipients, list):
		recipients = [recipients]

	dests = []
	for dest in recipients:
		if isinstance(dest, caccount):
			dest_firstname = dest.firstname
			dest_lastname =	dest.lastname
			dest_mail = dest.mail
			dest_full_mail = dest.get_full_mail()

			dests.append(dest_full_mail)	
		elif isinstance(dest, str) or isinstance(dest, unicode):
			if re.match("^[a-zA-Z0-9._%-]+@[a-zA-Z0-9._%-]+.([a-zA-Z]{2,6})?$", dest):
				dest_mail = dest
				dest_full_mail = "\"%s\" <%s>" % (dest_mail.split('@')[0].title(), dest_mail)
				dests.append(dest_full_mail)	
			else:
				raise ValueError('Invalid Email format for recipients')
		else:
			raise ValueError('Invalid Email format for recipients')

	dests_str = ', '.join(dests)

	###############
	# Attachments #
	###############
	if attachments:
		storage = cstorage(account=account, namespace='object')	
		if not isinstance(attachments, list):
			attachments = [attachments]

	########
	# Send #
	########
	logger.debug('-----')
	logger.debug('From: %s' % account_full_mail)
	logger.debug('To  : %s' % dests_str)
	logger.debug('-----')
	logger.debug('Subject: %s' % subject)
	logger.debug('Body   : %s' % body)
	logger.debug('Attach.: %s' % attachments)
	logger.debug('-----')

	msg = MIMEMultipart()
	msg["From"] = account_full_mail
	msg["To"] = dests_str
	msg["Subject"] = subject

	if html:
		msg.attach(MIMEText(body, 'html'))
	else:
		msg.attach(MIMEText(body, 'plain'))

	msg['Date'] = formatdate(localtime=True)

	if attachments:	
		for _file in attachments:
			part = MIMEBase('application', "octet-stream")
			if not isinstance(_file, cfile):
				_file.__class__ = cfile
			#meta_file = _file.get(storage)
			content_file = _file.get(storage)
			part.set_payload(content_file)
			Encoders.encode_base64(part)
			part.add_header('Content-Disposition', 'attachment; filename="%s"' % _file.data['file_name'])
			part.add_header('Content-Type', _file.data['content_type'])
			msg.attach(part)

	s = socket.socket()
	try:
		s.connect((smtp_host, smtp_port)) 
	except Exception as err:
		raise Exception('something\'s wrong with %s:%d. Exception type is %s' % (smtp_host, smtp_port, err))

	try:
		server = smtplib.SMTP(smtp_host, smtp_port)
		server.sendmail(account_full_mail, dests, msg.as_string())
		server.quit()
	except Exception, err:
		return "Error: unable to send email (%s)" % err
