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

def check(mfilter, event):
	# For each key of filter
	for key in mfilter:
		if key == '$and':
			# Check match for each elements in the list

			for element in mfilter['$and']:
				# If one does not match, then return False
				if not check(element, event):
					return False

		elif key == '$or':
			# Check match for each elements in the list

			for element in mfilter['$or']:
				# If one match, then return True
				if check(element, event):
					return True

			# Here nothing matched, then return False
			return False

		# For each other case, just test the equality
		else:
			if isinstance(mfilter[key], dict):
				if '$eq' in mfilter[key]:
					if event[key] != mfilter[key]['$eq']:
						return False

				elif '$gt' in mfilter[key]:
					if event[key] <= mfilter[key]['$gt']:
						return False

				elif '$gte' in mfilter[key]:
					if event[key] < mfilter[key]['$gte']:
						return False

				elif '$lt' in mfilter[key]:
					if event[key] >= mfilter[key]['$lt']:
						return False

				elif '$lte' in mfilter[key]:
					if event[key] > mfilter[key]['$lte']:
						return False

				elif '$in' in mfilter[key]:
					if event[key] not in mfilter[key]['$in']:
						return False

				elif '$nin' in mfilter[key]:
					if event[key] in mfilter[key]['$in']:
						return False

			else:
				if event[key] != mfilter[key]:
					return False

	# If we arrive here, everything matched
	return True