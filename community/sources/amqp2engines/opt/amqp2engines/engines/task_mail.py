# -*- coding: utf-8 -*-
#--------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------

from ctaskhandler import TaskHandler
from caccount import caccount
from cstorage import cstorage
from cfile import cfile

from email import Encoders
from email.MIMEBase import MIMEBase
from email.MIMEText import MIMEText
from email.MIMEMultipart import MIMEMultipart
from email.Utils import formatdate

import smtplib
import socket
import time
import re


class engine(TaskHandler):
	etype = 'task_mail'

	def __init__(self, *args, **kwargs):
		super(engine, self).__init__(*args, **kwargs)

	def handle_task(self, job):
		user = job.get('user', 'root')
		group = job.get('group', 'root')

		account = caccount(user=user, group=group)

		recipients = job.get('recipients', None)
		subject = job.get('subject', None)
		body = job.get('body', None)
		attachments = job.get('attachments', None)
		smtp_host = job.get('smtp_host', 'localhost')
		smtp_port = job.get('smtp_port', 25)
		html = job.get('html', False)

		# Execute the task
		return self.sendmail(account, recipients, subject, body, smtp_host, smtp_port, html)

	def sendmail(self, account, recipients, subject, body, attachments, smtp_host, smtp_port, html):
		"""
			account        : caccount or nothing for anon
			recipients     : str("glehee@capensis.fr"), caccount
			                 list of (caccount or string)
			subject        : str("My Subject")
			body           : str("My Body")
			attachments    : cfile, list of cfile
			smtp_host      : str("localhost")
			smtp_port      : int(25)
			html           : allow html into mail body (booleen)
		"""

		# Verify account
		account_firstname = account.firstname
		account_lastname = account.lastname
		account_mail = account.mail

		if not account_mail:
			self.logger.info('No mail adress for this user (Fill the mail account field)')
			account_mail = '{0}@{1}'.format(account.user, socket.gethostname())

		if isinstance(account_mail, (list, tuple)):
			account_mail = account_mail[0]

		if not account_lastname and not account_firstname:
			account_full_mail = '"{0}" <{1}>'.format(
				account_mail.split('@')[0].title(),
				account_mail
			)

		else:
			account_full_mail = account.get_full_mail()

		if not re.match("^[a-zA-Z0-9._%-]+@[a-zA-Z0-9._%-]+.([a-zA-Z]{2,6})?$", str(account_mail)):
			return (
				2,
				'Invalid Email format for sender: {0}'.format(account_mail)
			)

		# Verify recipients
		if not recipients:
			return (2, 'No recipients configured')

		if not isinstance(recipients, list):
			recipients = [recipients]

		dests = []

		for dest in recipients:
			if isinstance(dest, basestring):
				if re.match("^[a-zA-Z0-9._%-]+@[a-zA-Z0-9._%-]+.([a-zA-Z]{2,6})?$", dest):
					dest_mail = dest
					dest_full_mail = '"{0}" <{1}>'.format(
						dest_mail.split('@')[0].title(),
						dest_mail
					)
					dests.append(dest_full_mail)

			else:
				self.logger.error('Ignoring invalid recipient: {0}'.format(dest))

		dests_str = ', '.join(dests)

		# Verify attachments
		if attachments:
			storage = cstorage(account=account, namespace='object')

			if not isinstance(attachments, list):
				attachments = [attachments]

		# Send

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
			return (
				2,
				'Connection to SMTP <{0}:{1}> failed: {2}'.format(smtp_host, smtp_port, err)
			)

		try:
			server = smtplib.SMTP(smtp_host, smtp_port)
			server.sendmail(account_full_mail, dests, msg.as_string())
			server.quit()

		except Exception, err:
			return (
				2,
				"Imposible to send mail: {0}".format(err)
			)
