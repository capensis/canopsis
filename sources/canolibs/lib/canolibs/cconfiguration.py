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

	def __init__(self,
		configuration_file=DEFAULT_CONFIGURATION_FILE, auto_conf=True,
		naming_rule=None, logging_level=logging.INFO, _ready_to_conf=True,
		*args, **kwargs):

		super(Configurable, self).__init__(*args, **kwargs)

		self.configuration_file = configuration_file

		self.auto_conf = auto_conf

		self.logging_level = logging_level

		self.logger = logging.getLogger(type(self).__name__)

		self.naming_rule = naming_rule

		if _ready_to_conf and self.auto_conf:
			self.apply_configuration(
				parsers_by_option_by_section=self.get_parsers_by_option_by_section(),
				configuration_file=self.configuration_file,
				naming_rule=self.naming_rule)

	def _set_logging_level(self, option_value):
		self.logger.setLevel(option_value)
		return option_value

	def get_parsers_by_option_by_section(self):
		"""
		Get the right structure to use to parse a configuration file.

		:rtype: dict
		"""

		result = {
			'OPTIONS':
				{
					'auto_conf': bool,
					'logging_level': self._set_logging_level,
					'naming_rule': lambda option_value: eval(option_value)
				}
		}

		return result

	def apply_configuration(self, parsers_by_option_by_section=None,
		configuration_file=None, naming_rule=None):
		"""
		Apply configuration on a destination in 3 phases:

		1. Get options from input init configuration_file which match with self
		sections and options.
		2. convert options with related types.
		3. set destination attributes with options values where attribute names
		are option names in minuscule.

		:param parsers_by_option_by_section: struct to use
		in order to get values from configuration_file. If None, use
		self.get_parsers_by_option_by_section.
		:param configuration_file: configuration file.
		:param naming_rule: option naming rule for setvalues on self.

		:type parsers_by_option_by_section: dict
		:type configuration_file: str
		:type naming_rule: function or method

		:return: a dictionary of errors or parsers respectively for
		option parsing errors or not found options.
		:rtype: dict
		"""

		if parsers_by_option_by_section is None:
			parsers_by_option_by_section = self.get_parsers_by_option_by_section()

		if configuration_file is None:
			configuration_file = self.configuration_file

		config = ConfigParser.RawConfigParser()
		config.read(expanduser(configuration_file))

		result = parsers_by_option_by_section.copy()

		for section, parsers_by_options in \
			parsers_by_option_by_section.iteritems():

			if config.has_section(section):

				for option, parser in parsers_by_options.iteritems():

					if config.has_option(section, option):

						try:  # parsing option
							option_value = config.get(section, option)
							value = parser(option_value)
						except Exception as e:
							result[section][option] = e
							em = 'Impossible to parse {0}/{1} in {2}: {3}'
							em = em.format(section, option, configuration_file, e)
							self.logger.error(em)
						else:
							try:  # naming option
								name = naming_rule(option)
							except Exception as e:
								result[section][option] = e
								em = 'Impossible to rename option {0}/{1} in \
									{2}: {3}'.format(
										section, option, configuration_file, e)
								self.logger.error(em)
							else:
								try:  # setattr
									setattr(self, name, value)
								except TypeError as e:
									result[section][option] = e
									em = 'Impossible to set attribute {0}.{1} \
									to {2}: {3}'.format(self, name, value, e)
									self.logger.error(em)
								else:
									del result[section][option]

					else:
						wm = 'option {0}/{1} not found in {2}'.format(
							section, option, configuration_file)
						self.logger.warning(wm)

				# if no option, remove the section
				if not result[section]:
					del result[section]

			else:
				wm = 'section {0} not found in {1}'.format(
					section, configuration_file)
				self.logger.warning(wm)

		return result
