#!/usr/bin/env python
#--------------------------------
# Copyright (c) 2011 "Capensis" [http://www.capensis.com]
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

# MongoDB Operators:
# http://docs.mongodb.org/manual/reference/operator/

def field_check(mfilter, event, key):
	for op in mfilter[key]:
		if op == '$exists':
			#check if key is in event
			if mfilter[key][op]:
				if key not in event:
					return False
			#check if key is not in event
			else:
				if key in event:
					return False

		elif op == '$eq':
			if event[key] != mfilter[key][op]:
				return False

		elif op == '$ne':
			if event[key] == mfilter[key][op]:
				return False

		elif op == '$gt':
			if event[key] <= mfilter[key][op]:
				return False

		elif op == '$gte':
			if event[key] < mfilter[key][op]:
				return False

		elif op == '$lt':
			if event[key] >= mfilter[key][op]:
				return False

		elif op == '$lte':
			if event[key] > mfilter[key][op]:
				return False

		elif op == '$in':
			if event[key] not in mfilter[key][op]:
				return False

		elif op == '$nin':
			if event[key] in mfilter[key][op]:
				return False

		elif op == '$not':
			reverse_mfilter = {}
			reverse_mfilter[key] = mfilter[key][op]

			if field_check(reverse_mfilter, event, key):
				return False

		else:
			if event[key] != mfilter[key]:
				return False

	return True

def check(mfilter, event):
	# For each key of filter
	for key in mfilter:
		if key == '$and':
			# Check match for each elements in the list

			for element in mfilter[key]:
				# If one does not match, then return False
				if not check(element, event):
					return False

		elif key == '$or':
			# Check match for each elements in the list

			for element in mfilter[key]:
				# If one match, then return True
				if check(element, event):
					return True
			# Here nothing matched, then return False
			return False

		elif key == '$nor':
			# Check match for each elements in the list

			for element in mfilter[key]:
				# If one match, then return False
				if check(element, event):
					return False

		# For each other case, just test the equality
		elif key in event:
			if isinstance(mfilter[key], dict):
				if not field_check(mfilter, event, key):
					return False

			else:
			
				if event[key] != mfilter[key]:
					return False

		else:
			return False

	# If we arrive here, everything matched
	return True