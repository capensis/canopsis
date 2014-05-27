#!/usr/bin/env python
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

__all__ = ('Configurable')

import ConfigParser
from os.path import expanduser

import logging


class Configurable(object):
	"""
	Manages configuration from a dictionary of
	parsers by options by sections to parse from a configuration file.
	"""

	DEFAULT_CONFIGURATION_FILE = '~/etc/conf.conf'

	MASTER = 'MASTER'

	AUTO_CONF = 'auto_conf'
	LOGGING_LEVEL = 'logging_level'

	PARSING_RULES = {
		MASTER:
			{
				AUTO_CONF: bool,
				LOGGING_LEVEL: str
			}
	}

	def __init__(self,
		configuration_file=DEFAULT_CONFIGURATION_FILE, auto_conf=True,
		logging_level=logging.INFO, _ready_to_conf=True,
		*args, **kwargs):

		super(Configurable, self).__init__(*args, **kwargs)

		self.configuration_file = configuration_file

		self.auto_conf = auto_conf

		self.logging_level = logging_level

		self.logger = logging.getLogger(type(self).__name__)

		if _ready_to_conf and self.auto_conf:
			self.apply_configuration(
				parsing_rules=self.get_parsing_rules(),
				configuration_file=self.configuration_file)

	def get_parsing_rules(self, *args, **kwargs):
		"""
		Get the right structure to use to parse a configuration file.

		:rtype: list of dict(section, dict(option, parser(option)))
		"""

		result = [Configurable.PARSING_RULES]

		return result

	def apply_configuration(self, parsing_rules=None,
		configuration_file=None, *args, **kwargs):
		"""
		Apply configuration on a destination in 3 phases:

		1. Get options from input init configuration_file which match with self
		sections and options.
		2. convert options with related types.
		3. set destination attributes with options values where attribute names
		are option names in minuscule.

		:param parsing_rules: struct to use
		in order to get values from configuration_file. If None, use
		self.get_parsing_rules.
		:param configuration_file: configuration file.

		:type parsing_rules: dict
		:type configuration_file: str
		"""

		if parsing_rules is None:
			parsing_rules = self.get_parsing_rules()

		if configuration_file is None:
			configuration_file = self.configuration_file

		if isinstance(configuration_file, str):
			configuration_file = [configuration_file]

		config = ConfigParser.RawConfigParser()
		config.read(expanduser(configuration_file))

		parameters = dict()
		error_parameters = dict()

		log_message = '{0}/{1}'.format(configuration_file, '{0}/{1}')

		for parsing_rule in parsing_rules:

			for section, parsers_by_options in parsing_rule.iteritems():

				if config.has_section(section):

					for option, parser in parsers_by_options.iteritems():

						option_log_message = '{0} = {1}'.format(
							log_message.format(section, option), '{0}')

						if config.has_option(section, option):

							option_value = config.get(section, option)

							try:  # parsing option
								value = parser(option_value)

							except Exception as e:
								error_parameters[option] = e
								parameters.pop(option, None)
								error_message = option_log_message.format(e)
								self.logger.error(error_message)

							else:
								parameters[option] = value
								error_parameters.pop(option, None)
								info_message = option_log_message.format(value)
								self.logger.info(info_message)

		self._set_parameters(parameters, error_parameters)

	def _set_parameters(self, parameters, error_parameters, *args, **kwargs):

		logging_level = parameters.get(Configurable.LOGGING_LEVEL)
		if logging_level is not None:
			self.logger.setLevel(logging_level)

		auto_conf = parameters.get(Configurable.AUTO_CONF)
		if auto_conf is not None:
			self.auto_conf = auto_conf
