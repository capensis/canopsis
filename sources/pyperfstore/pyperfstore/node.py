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

import logging
import time
import hashlib

from pyperfstore.metric import metric
from pyperfstore.dca import dca
from pyperfstore.pmath import aggregate as pmath_aggregate

def make_metric_id(node_id, metric_dn):
	#return hashlib.md5(node_id+"-"+metric_dn).hexdigest()
	return node_id.replace('.','-') + "-" + hashlib.md5(metric_dn).hexdigest()

class node(object):
	def __init__(self, _id, storage, dn=None, point_per_dca=None, retention=None, rotate_plan=None, logging_level=None):
		self.logger = logging.getLogger('node')
		if logging_level:
			self.logger.setLevel(logging_level)
			
		self.logger.debug("Init node '%s'" % dn)
	
		self._id = _id
		#self._id = hashlib.md5(dn).hexdigest()
		self.dn = None
		self.retention = retention

		self.point_per_dca = point_per_dca
		self.rotate_plan = rotate_plan

		self.storage = storage

		self.metrics = {}
		self.metrics_id = {}

		self.writetime = None

		data = self.storage.get(self._id)
		if data:
			# Load data
			self.load(data)
			
			# Save new configurations
			changed = False
			if dn and dn != self.dn:
				self.dn = dn
				changed = True
				
			if retention and retention != self.retention:
				self.retention = retention
				changed = True
				
			if rotate_plan and rotate_plan != self.rotate_plan:
				self.rotate_plan = rotate_plan
				changed = True
			
			if changed:
				self.logger.debug(" + Save Node with new options" % dn)
				self.save()
				
			if not self.dn and not dn:
				self.dn = _id
		else:
			self.dn = dn
				
	def dump(self):
		dump = {
			#'id':		self._id,
			'dn':		self.dn,
			'retention':	self.retention,
			'point_per_dca':self.point_per_dca,
			'rotate_plan':	self.rotate_plan,
			'metrics':		self.metrics,
			'writetime':	time.time()
		}

		for _id in self.metrics.keys():
			item = self.metrics[_id]
			if isinstance(item ,metric):
				dump['metrics'][_id] = { 'id': item._id, 'dn': item.dn, 'bunit': item.bunit }
			else:
				dump['metrics'][_id] = { 'id': item['id'], 'dn': item['dn'], 'bunit': item['bunit'] }
			
		return dump

	def load(self, data):
		self.logger.debug("Load node '%s'" % self._id)

		#self._id			= data['id']
		self.dn				= data['dn']
		self.retention		= data['retention']
		self.point_per_dca	= data['point_per_dca']
		self.rotate_plan	= data['rotate_plan']
		
		self.writetime		= data['writetime']
		self.metrics = data['metrics']
			
	def save(self):
		dump = self.dump()

		self.logger.debug("Save node '%s'" % self._id)
		self.storage.set(self._id, dump)

	def metric_make_id(self, dn):
		return make_metric_id(self._id, dn)

	def metric_get(self, dn=None, _id=None):
		_id = self.metric_get_id(dn, _id)
			
		try:
			item = self.metrics[_id]
		except:
			self.logger.error("Unknown metric '%s' ... " % dn)
			return None

		if not isinstance(item ,metric):
			## load metric from store
			item = metric(_id=item['id'], node=self, storage=self.storage)
			self.metrics[_id] = item

		return item

	def metric_get_all_dn(self):
		try:
			dump = self.dump()
			dns = [ self.metrics[key]['dn'] for key in dump['metrics'] ]
		except:
			## Old format:
			dns = []
			self.logger.warning("Convert Node in new format !")
			for _id in self.metrics.keys():
				metric = self.metric_get(_id=_id)
				dns.append(metric.dn)
			self.save()
		
		return dns

	def metric_get_id(self, dn=None, _id=None):
		if _id:
			return _id
		
		if not dn:
			return None

		return self.metric_make_id(dn)
		
	def metric_dump(self, dn=None, _id=None):
		_id = self.metric_get_id(dn, _id)
		
		item = self.metrics[_id]

		if isinstance(item ,metric):
			## load metric from store
			return item.dump()

		return self.metric_get(_id=_id).dump()

	def metric_exist(self, dn=None, _id=None):	
		_id = self.metric_get_id(dn, _id)
		
		if not _id:
			return False
			
		try:		
			return self.metrics[_id]
		except:
			return False

	def metric_add(self, dn, bunit=None, dtype=None):
		self.logger.debug("Add metric '%s' (%s)" % (dn, bunit))

		if not self.metric_exist(dn=dn):
			metric_id = self.metric_make_id(dn)

			self.logger.debug(" + Metric ID: '%s'" % metric_id)
			mymetric = metric(
							_id=metric_id,
							dn=dn,
							bunit=bunit,
							dtype=dtype,
							node=self,
							retention=self.retention,
							storage=self.storage,
							point_per_dca=self.point_per_dca,
							rotate_plan=self.rotate_plan,
						)
						
			mymetric.save()
			self.metrics[metric_id]	= mymetric
			self.save()
			
			return mymetric
		else:
			self.logger.debug(" + Metric allready exist")
			_id = self.metric_get_id(dn=dn)
			return metric_get(_id=_id)

	def metric_get_values(self, tstart, tstop=None, aggregate=True, max_points=None, atype=None, dn=None, _id=None, time_interval=None):
		_id = self.metric_get_id(dn, _id)
		if not _id:
			return []
		
		if not tstop:
			tstop = int(time.time())

		tstart = int(tstart)
		tstop = int(tstop)			

		self.logger.debug("Get values in '%s'" % dn)
		
		mode = None
		if time_interval:
			mode = 'by_interval'

		mymetric = self.metric_get(_id=_id)

		if mymetric:
			values = mymetric.get_values(tstart, tstop)
			if aggregate:
				return pmath_aggregate(values, max_points=max_points, atype=atype, time_interval=time_interval, mode=mode)
			else:
				return values
		else:
			return []

	def metric_push_value(self, value, unit=None, timestamp=None, dn=None, _id=None, dtype=None, point_per_dca=None, min_value=None, max_value=None, thld_warn_value=None, thld_crit_value=None):
		
		_id = self.metric_get_id(dn, _id)
		self.logger.debug("Push value on metric '%s'" % dn)
		
		if not self.metric_exist(dn=dn):
			mymetric = self.metric_add(dn=dn, bunit=unit, dtype=dtype)
		else:
			mymetric = self.metric_get(_id=_id)
			
		if not timestamp:
			timestamp = int(time.time())
		else:
			timestamp = int(timestamp)
		
		## re-Set dtype
		if dtype:
			if mymetric.dtype != dtype:
				mymetric.dtype = dtype
				
		## re-Set bunit
		if unit:
			if mymetric.bunit != unit:
				mymetric.bunit = unit
				
		## re-Set min/max
		if min_value:
			mymetric.min_value = min_value
		if max_value:
			mymetric.max_value = max_value
			
		## re-Set Threshold
		if thld_warn_value:
			mymetric.thld_warn_value = thld_warn_value
		if thld_crit_value:
			mymetric.thld_crit_value = thld_crit_value
		
		## re-Set point/dca
		if point_per_dca:
			mymetric.auto_point_per_dca = False
			mymetric.point_per_dca = point_per_dca		

		mymetric.push_value(value=value, timestamp=timestamp)

	def metric_remove(self, dn=None, _id=None):
		self.logger.debug("Remove metric '%s'" % dn)
		
		_id = self.metric_get_id(dn, _id)
		if not _id:
			return None
			
		mymetric = self.metric_get(_id=_id)
		
		dn = mymetric.dn
		
		mymetric.dca_remove_all()
		self.storage.rm(_id)

		del self.metrics[_id]
		del mymetric
	
		self.save()
		
	def metric_remove_all(self):
		for _id in self.metrics.keys():
			item = self.metric_get(_id=_id)
			self.metric_remove(_id=_id)

	def remove(self):
		self.metric_remove_all()

		self.logger.debug("Remove node '%s'" % self._id)
		self.storage.rm(self._id)
		
	def size(self):
		size = self.storage.size(self._id)
		
		for _id in self.metrics.keys():
			item = self.metric_get(_id=_id)
			size += item.size()
			
		return size

	def pretty_print(self):
		print " + Id: %s" % self._id
		print " + Node DN: %s" % self.dn
		print " + Retention: %s" % self.retention
		print " + Rotate_plan: %s" % self.rotate_plan
		print " + Metrics:"

		for _id in self.metrics.keys():

			metric = self.metric_get(_id=_id)

			print "    + %s (%s) (%s)" % (metric.dn, metric.dtype, metric._id)

			item = metric.dca_get(metric.current_dca)
			
			bsize = self.storage.size(item._id) / 1024.0
			print "      + Current DCA (%s -> %s),\tPoints: %s\t%.2f KB" % (item.tstart, item.tstop, item.size, bsize )
			print ""

			if metric.dca_PLAIN:
				for item in metric.dca_PLAIN:
					item = metric.dca_get(item)
					bsize = self.storage.size(item._id) / 1024.0
					print "      + %s DCA (%s -> %s),\tPoints: %s\t%.2f KB" % (item.format, item.tstart, item.tstop, item.size, bsize )

				print ""

			if metric.dca_TSC:
				for item in metric.dca_TSC:
					item = metric.dca_get(item)
					bsize = self.storage.size(item._id) / 1024.0
					print "      + %s DCA (%s -> %s),\tPoints: %s\t%.2f KB" % (item.format, item.tstart, item.tstop, item.size, bsize )

				print ""

			if metric.dca_ZTSC:
				for item in metric.dca_ZTSC:
					item = metric.dca_get(item)
					bsize = self.storage.size(item._id) / 1024.0
					print "      + %s DCA (%s -> %s),\tPoints: %s\t%.2f KB" % (item.format, item.tstart, item.tstop, item.size, bsize )

				print ""
