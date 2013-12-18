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

import os, sys, json, logging, time
import hashlib
from datetime import datetime

from pyperfstore2.store import store
import pyperfstore2.utils as utils

class manager(object):
	def __init__(self, retention=0, dca_min_length = 250, logging_level=logging.INFO, cache=True, **kwargs):
		self.logger = logging.getLogger('manager')
		self.logger.setLevel(logging_level)

		# Store
		self.store = store(logging_level=logging_level, **kwargs)
		
		self.dca_min_length = dca_min_length
		
		# Seconds
		self.retention = retention
		
		# Cache
		self.cache = cache
		self.cache_max_size = 5000
		self.cache_size = 0
		self.md5_cache = {}
		
		self.fields_map = {
				'retention':	('r', self.retention),
				'type':			('t', 'GAUGE'),
				'unit':			('u', None),
				'min':			('mi', None),
				'max':			('ma', None),
				'thd_warn':		('tw', None),
				'thd_crit':		('tc', None)
		}
				
	def gen_id(self, name):
		return hashlib.md5(name).hexdigest()
		
	def get_id(self, _id=None, name=None):
		if not _id and not name:
			raise Exception('Invalid args')
		
		if not _id:
			if self.cache:
				if self.cache_size <= self.cache_max_size:
					_id = self.md5_cache.get(name, None)
					if not _id:
						_id = self.gen_id(name)
						self.md5_cache[name] = _id
						self.cache_size += 1
		
		if not _id:
			_id = self.gen_id(name)
			
		return _id
		
	def get_data(self, _id):
		data = self.store.redis.lrange(_id, 0, -1)
		def cleanPoint(p):
			p[0] = int(p[0])
			try:
				p[1] = int(p[1])
			except:
				p[1] = float(p[1])
			return p

		data = [ cleanPoint(p.split('|')) for p in data ]
		return data

	def get_meta(self, _id=None, name=None, raw=False, mfields=None):
		_id = self.get_id(_id, name)
		
		meta_data = self.store.get(_id, mfields=mfields)
		if not meta_data:
			return None

		if not mfields or mfields.get('d', False):
			meta_data['d'] = self.get_data(_id)
		
		# Uncompress fields name
		if not raw:
			meta_data = self.uncompress_meta_fields(meta_data)
		
		return meta_data
		
	def uncompress_meta_fields(self, meta_data):
		for field in self.fields_map:
			value = meta_data.get(self.fields_map[field][0], self.fields_map[field][1])
			meta_data[field] = value
			try:
				del meta_data[self.fields_map[field][0]]
			except:
				pass
				
		return meta_data
	
	def compress_meta_fields(self, meta_data):
		# Compress fields name
		def set_meta(meta, field, new_field, default):
			# Check if field is already compressed
			data = meta.get(new_field, None)
			
			if data != None:
				return meta
			else:
				data = meta.get(field, default)
				if data != None:
					meta[new_field] = data
				try:
					del meta[field]
				except:
					pass
				return meta

		for field in self.fields_map:
			meta_data = set_meta(meta_data, field, self.fields_map[field][0], self.fields_map[field][1])
			
		return meta_data

	def push(self, value, _id=None, name=None, timestamp=None, meta_data={}):
		_id = self.get_id(_id, name)
						
		if not timestamp:
			timestamp = int(time.time())
			
		point = (timestamp, value)
		
		meta_data = meta_data.copy()
		
		if meta_data:
			meta_data = self.compress_meta_fields(meta_data)
		self.store.push(_id=_id, point=point, meta_data=meta_data)
		
	def find(self, _id=None, name=None, mfilter=None, limit=0, skip=0, data=True, sort=None):
		mfields = None
		if _id or name:
			_id = self.get_id(_id, name)
		if _id:
			if mfilter:
				mfilter={"$and":[{'_id': _id}, mfilter]}
			else:
				mfilter={'_id': _id}
		else:
			if not mfilter:
				mfilter = {}
		
		if not data:
			mfields = { 'd': 0 }
		
		return self.store.find(mfilter=mfilter, limit=limit, skip=skip, mfields=mfields, sort=sort)
		
	def get_points(self, _id=None, name=None, tstart=None, tstop=None, raw=False, return_meta=False, add_prev_point=False, add_next_point=False):
		_id = self.get_id(_id, name)
		if tstop == None:
			tstop = int(time.time())
		if tstart == None:
			tstart = tstop
		self.logger.debug("Get points: %s (%s -> %s)" % (_id, datetime.utcfromtimestamp(tstart), datetime.utcfromtimestamp(tstop)))
		points = []
		
		dca = self.get_meta(_id=_id)
		if not dca :
			raise Exception('Invalid _id, not found %s' % _id)
		
		plain_fts = None
		plain_lts = None
		
		if dca.get('d', False):
			plain_fts = dca['d'][0][0]
			plain_lts = dca['d'][-1][0]
		self.logger.debug(" + plain_fts: %s" % plain_fts)
		self.logger.debug(" + plain_lts: %s" % plain_lts)
		
		## Check Compressed DCA
		if not plain_fts or tstart < plain_fts:
			self.logger.debug(" + Search in compressed DCA")
			bin_ids = []
			
			for bin_meta in dca.get('c', []):
				fts = bin_meta[0]
				lts = bin_meta[1]
				bin_id = bin_meta[2]
				self.logger.debug(" + Parse DCA:\t\t%s (%s -> %s)" % (bin_id, datetime.utcfromtimestamp(fts), datetime.utcfromtimestamp(lts)))
				if tstart == tstop and tstart >= fts and tstart <= lts:
					bin_ids.append(bin_id)
					self.logger.debug("   + Append")
				elif utils.intersection([fts, lts], [tstart, tstop]):
					bin_ids.append(bin_id)
					self.logger.debug("   + Append")
			for bin_id in bin_ids:
				data = self.store.get_bin(_id=bin_id)
				points += utils.uncompress(data)
					
		## Check Plain DCA
		self.logger.debug(" + Search in plain DCA")
		if plain_fts and plain_lts:
			if tstart == tstop and tstart >= plain_fts and tstart <= plain_lts:
				points += dca['d']
				self.logger.debug("   + Append")
			elif tstart == tstop and tstart >= plain_lts:
				points += dca['d']
				self.logger.debug("   + Append")
			elif utils.intersection([plain_fts, plain_lts], [tstart, tstop]):
				points += dca['d']
				self.logger.debug("   + Append")
		## Drop data of meta
		del dca['d']
		
		self.logger.debug(" + len(points):  %s" % len(points))
		
		## Sort and Split Points
		points = sorted(points, key=lambda point: point[0])
		if add_prev_point or add_next_point:
			self.logger.debug("Find previous and next points")
			
			prev_index = None
			next_index = None
			rpoints = []
			i=0
			for point in points:
				if point[0] >= tstart and point[0] <= tstop:
					rpoints.append(point)
	
				if not prev_index and point[0] >= tstart:
					prev_index = i - 1
	
				i+=1
				
				if point[0] > tstop:
					break
					
			next_index = i-1
			
			# If tstop > last point
			if add_prev_point and prev_index == None and next_index != None and next_index < len(points):
				prev_index  = next_index
			self.logger.debug(" + len(points):  %s" % len(points))
			self.logger.debug(" + len(rpoints): %s" % len(rpoints))
			self.logger.debug(" + prev_index:   %s" % prev_index)
			self.logger.debug(" + next_index:   %s" % next_index)
			
			# Add points
			if add_prev_point and prev_index != None and prev_index >= 0:
				self.logger.debug("   + Add prev")
				rpoints.insert(0, points[prev_index])
				if (dca['type'] == 'DERIVE' and tstart == tstop and prev_index-1 >= 0):
					self.logger.debug("   + Add prev for DERIVE")
					rpoints.insert(0, points[prev_index-1])
			if add_next_point and next_index != None and next_index < len(points):
				self.logger.debug("   + Add next")
				rpoints.append(points[next_index])

			self.logger.debug(" + len(rpoints): %s" % len(rpoints))

			points = rpoints
		else:
			points = [ point for point in points if point[0] >= tstart and point[0] <= tstop ]
		if raw:
			if not return_meta:
				return points
			else:
				return (dca, points)
			
		#parse_dst
		dtype = dca.get('type', None)
		if dtype:
			points = utils.parse_dst(points,dtype)

		if not return_meta:
			return points
		else:
			return (dca, points)
	
	def get_last_point(self, *args, **kargs):
		return self.get_point(*args, ts=None, **kargs)
	
	def get_point(self, _id=None, ts=None, name=None, return_meta=False):
		_id = self.get_id(_id, name)
		
		meta = None
		point = None
		
		if not ts:
			dca = self.get_meta(_id=_id)
			points = dca.get('d', [])
		else:
			(meta, points) = self.get_points(_id=_id, tstart=ts, tstop=ts, add_prev_point=True, return_meta=True)	
		
		if len(points):
			point = points.pop()
			
		if not meta and return_meta:
			meta = self.get_meta(_id=_id)
			
		if return_meta:
			return (meta, point)
		else:
			return point
	
	def rotateAll(self, concurrency=1):
		t = time.time()

		if concurrency > 1:
			from multiprocessing import Pool

		self.logger.info("Rotate All DCA")
		_ids = []
		self.logger.info(" + Get all keys")
		keys = self.store.redis.keys('*')

		try:
			keys.remove("perfstore2:rotate:plan")
		except:
			pass
			
		self.logger.info(" + Check length (%s keys)" % len(keys))
		for key in keys:
			self.store.redis_pipe.llen(key)
	
		result = self.store.redis_pipe.execute()

		for index, key in enumerate(keys):
			if result[index] >= self.dca_min_length:
				_ids.append(key)

		if not _ids:
			self.logger.info("Nothing to do")
			return

		if concurrency <= 1:
			for _id in _ids:
				self.rotate(_id)
		else:
			_ids = split_list(_ids, wanted_parts=concurrency)

			p = Pool(concurrency)
			p.map(rotate_process, _ids)

		t = time.time() - t
		self.logger.info("All perfdata was rotate, elapsed: %.3f seconds" % t)
		
	def rotate(self, _id=None, name=None):
		t = time.time()
		try:
			_id = self.get_id(_id, name)
		except:
			_id = None			
	
		if not _id:
			self.logger.info("Nothing to do")
			return

		self.logger.debug("Start rotation of %s" % _id)

		self.logger.debug(" + DCA: %s" % _id)

		points = self.get_data(_id)
		
		if not points:
			self.logger.debug("No points, Nothing to do")
			return

		fts = points[0][0]
		lts = points[-1][0]
						
		self.logger.debug("  + Compress %s -> %s" % (fts, lts))
		
		data = utils.compress(points)

		try:
			bin_id = "%s%s" % (_id, lts)
			self.logger.debug("   + Store in binary record")
			self.store.create_bin(_id=bin_id, data=data)
			
			self.logger.debug("   + Add bin_id in meta and clean meta")
			##ofts = dca.get('fts', fts)
			##self.store.update(_id=_id, mset={'fts': ofts, 'd': []}, mpush={'c': (fts, lts, bin_id)})
			
			self.store.update(_id=_id, mpush={'c': (fts, lts, bin_id)})
			self.store.redis.delete(_id)
			
		except Exception,err:
			self.logger.warning('Impossible to rotate %s: %s' % (_id, err))
				
		#else:
		#	self.logger.debug("  + Not enough point in DCA")
		#	ofts = dca.get('fts', fts)
		#	self.store.update(_id=_id, mset={'fts': ofts})

		t = time.time() - t
		self.logger.debug(" + Rotation of '%s' done in %.3f seconds" % (_id, t))
			
	def cleanAll(self, timestamp=None):
		return self.clean(timestamp=timestamp)
	
	def clean(self, _id=None, name=None, timestamp=None):
		try:
			_id = self.get_id(_id, name)
		except:
			_id = None
		
		if not timestamp:
			timestamp = int(time.time())
		
		cleaned = 0
		
		self.logger.debug("Remove DCA when 'fts' is older than %s:" % timestamp)
		
		if _id:
			meta = self.find(limit=1, mfilter={'_id': _id, 'fts': {'$lt': timestamp}})
			if meta:
				metas = [meta]
			else:
				metas = []
			nb_metas = len(metas)
		else:
			metas = self.find(limit=0, mfilter={'fts': {'$lt': timestamp}})
			nb_metas = metas.count()

		if not nb_metas:
			self.logger.debug("   + Nothing to clean")
			return cleaned
		else:
			self.logger.debug("   + Start cleanning of %s metas" % nb_metas)
					
		for meta in metas:
			meta_id = meta['_id']
			self.logger.debug("   + Clean meta '%s'" % meta_id)
										
			## Clean binaries
			bin_fts = None
			for dca_meta in meta['c']:
				fts = dca_meta[0]
				lts = dca_meta[1]
				bin_id = dca_meta[2]
				
				# check lts
				if  lts  <= timestamp:
					self.logger.debug("     + Remove binarie DCA '%s'" %  bin_id)
					self.store.grid.delete(bin_id)
					
					# Remove dca meta
					self.store.update(_id=meta_id, mpop={ 'c' : -1  })
				else:
					if not bin_fts:
						bin_fts = fts
			
			## Clean plain
			plain_fts = None
			points = self.get_data(meta_id)
			if len(points):
				if points[0][0] <= timestamp:
					for point in points:
						fts = point[0]
						if fts <= timestamp:
							self.logger.debug("     + Remove plain point")
						else:
							points.append(point)
							if not plain_fts:
								plain_fts = fts
					if points:
						points = ["%s|%s" % (p[0], p[1]) for p in points]
						self.store.redis.lset(meta_id, points)
			
			## Set new fts
			fts = plain_fts
			if bin_fts < plain_fts:
				fts = bin_fts
			
			self.store.update(_id=meta_id, mset={'fts': fts})
			cleaned += 1
						
		return cleaned
	
	def update(self, _id=None, name=None, data=None):
		_id = self.get_id(_id, name)
		self.store.update(_id, mset=data)
	
	def remove(self, _id=None, name=None, purge=True):
		ids = []
		if isinstance(_id, list):
			ids = [self.get_id(iid, name) for iid in _id]
		else:
			ids = [self.get_id(_id, name)]

		if (len(ids) > 1):
			self.logger.debug("Remove: %s DCA" % len(ids))
		else:
			self.logger.debug("Remove: %s" % ids)
		
		dcas = []
		bin_dcas = []

		for _id in ids:
			self.store.redis_pipe.delete(_id)
			dca = self.get_meta(_id=_id, raw=True, mfields={'c': 1})
			if dca:
				dcas.append(dca)
				binaries = dca.get('c', [])
				if len(binaries):
					for bin_dca in binaries:
						bin_dcas.append(bin_dca[2])

		self.store.sync()

		self.logger.debug(" + %s Meta DCA Found" % len(dcas))
		self.logger.debug(" + %s Binaries Found" % len(bin_dcas))


		if len(bin_dcas):
			self.logger.debug("Remove Compressed Binaries ...")
			self.store.db[self.store.mongo_collection+"_bin.chunks"].remove({'files_id': {'$in': bin_dcas}})
			self.store.db[self.store.mongo_collection+"_bin.files"].remove({'_id': {'$in': bin_dcas}})

		if len(dcas):
			self.logger.debug("Remove Meta and Plains DCA ...")
			if len(dcas) == 1:
				self.store.remove(_id=dcas[0]['_id'])
			else:
				self.store.remove(mfilter={'_id': {'$in': [ dca['_id'] for dca in dcas]}})

	def showStats(self):
		metas = self.find(limit=0)
		mcount = metas.count()
		size = self.store.size()
		
		self.logger.info("Metas:       %s" % mcount)
		if mcount:
			self.logger.info("Size/metric: %.3f KB" % ((float(size)/mcount)/1024.0))
		self.logger.info("Total size:  %.3f MB" % (size/1024.0/1024.0))
	
	def showAll(self):
		metas = self.find(limit=0)
		for meta in metas:
			self.show(meta=self.uncompress_meta_fields(meta))
			
	def show(self, _id=None, name=None, meta=None):
		if not meta:
			_id = self.get_id(_id, name)
			meta = self.get_meta(_id=_id)
		else:
			_id = meta['_id']
		
		if meta and _id:
			self.logger.info("Metadata:'%s'" % meta['_id'])
			for key in meta:
				if key != '_id' and key != 'c' and key != 'd' and key != 'nc':
					self.logger.info(" + %s: %s" % (key, meta[key]))
			
			self.logger.info(" + Compressed DCA: %s" % len(meta.get('c', [])))
			self.logger.info(" + Next Clean: %s" % meta.get('nc', None) )


def split_list(alist, wanted_parts=1):
    length = len(alist)
    return [ alist[i*length // wanted_parts: (i+1)*length // wanted_parts] 
             for i in range(wanted_parts) ]

def rotate_process(_ids):
	import pyperfstore2
	manager = pyperfstore2.manager()

	for _id in _ids:
		manager.rotate(_id=_id)