# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
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
import sys
import os
import ConfigParser

import gridfs

from bson import objectid

from pymongo import Connection
from pymongo import ASCENDING
from pymongo import DESCENDING

from canopsis.mongo.core import CanopsisSONManipulator
from canopsis.old.account import Account
from canopsis.old.record import Record

from operator import itemgetter

from urlparse import urlparse

CONFIG = ConfigParser.RawConfigParser()
CONFIG.read(os.path.join(sys.prefix, 'etc', 'cstorage.conf'))


class Storage(object):
    def __init__(
        self,
        account, namespace='object',
        logging_level=logging.ERROR,
        mongo_uri=None,
        mongo_host="127.0.0.1",
        mongo_port=27017,
        mongo_userid=None,
        mongo_password=None,
        mongo_db='canopsis',
        mongo_autoconnect=True,
        groups=[],
        mongo_safe=True,
        *args, **kwargs
    ):

        super(Storage, self).__init__(*args, **kwargs)

        self.logger = logging.getLogger('Storage')
        self.logger.setLevel(logging_level)

        try:
            self.mongo_uri = CONFIG.get('master', 'db_uri')

        except ConfigParser.Error:
            self.mongo_uri = mongo_uri

        try:
            self.mongo_host = CONFIG.get("master", "host")

        except ConfigParser.Error:
            self.mongo_host = mongo_host

        try:
            self.mongo_port = CONFIG.getint("master", "port")

        except ConfigParser.Error:
            self.mongo_port = mongo_port

        try:
            self.mongo_db = CONFIG.get("master", "db")

        except ConfigParser.Error:
            self.mongo_db = mongo_db

        try:
            self.mongo_userid = CONFIG.get("master", "userid")

        except ConfigParser.Error:
            self.mongo_userid = mongo_userid

        try:
            self.mongo_password = CONFIG.get("master", "password")

        except ConfigParser.Error:
            self.mongo_password = mongo_password

        try:
            self.fetch_limit = int(CONFIG.get("master", "fetch_limit"))

        except ConfigParser.Error:
            self.fetch_limit = 10000

        try:
            self.no_count_limit = int(CONFIG.get("master", "no_count_limit"))
        except ConfigParser.Error:
            self.no_count_limit = 200000

        self.mongo_safe = mongo_safe

        self.account = account
        self.root_account = Account(user="root", group="root")

        self.namespace = namespace
        self.backend = None

        self.gridfs_namespace = "binaries"

        self.logger.debug("Object initialised.")

        self.backend = {}
        self.connected = False

        if mongo_autoconnect:
            self.connect()

    def clean_id(self, _id):
        return _id

    def make_mongofilter(self, account):
        Read_mfilter = {}
        Write_mfilter = {}

        if account._id != "account.root" and account.group != "group.CPS_root" and not 'group.CPS_root' in account.groups:
            Read_mfilter = { '$or': [
                {'aaa_owner': account._id, 'aaa_access_owner': 'r'},
                {'aaa_group': account.group, 'aaa_access_group': 'r'},
                {'aaa_group': {'$in': account.groups}, 'aaa_access_group': 'r'},
                {'aaa_admin_group':account.group},
                {'aaa_admin_group':{'$in': account.groups}},
                {'aaa_access_unauth': 'r'}
            ] }

            Write_mfilter = { '$or': [
                {'aaa_owner': account._id, 'aaa_access_owner': 'w'},
                {'aaa_group': account.group, 'aaa_access_group': 'w'},
                {'aaa_group': {'$in': account.groups}, 'aaa_access_group': 'w'},
                {'aaa_admin_group':account.group},
                {'aaa_admin_group':{'$in': account.groups}},
                {'aaa_access_unauth': 'w'}
            ] }

            if account.user != "anonymous":
                Read_mfilter['$or'].append({'aaa_access_other': 'r'})
                Write_mfilter['$or'].append({'aaa_access_other': 'w'})

        return (Read_mfilter, Write_mfilter)

    @property
    def uri(self):
        if self.mongo_uri:
            uri = self.mongo_uri

        else:
            uri = '%s:%s' % (self.mongo_host, self.mongo_port)

            if self.mongo_userid is not None:
                if self.mongo_password is None:
                    uri = '%s@%s' % (self.mongo_userid, uri)

                else:
                    uri = '%s:%s@%s' % (self.mongo_userid, self.mongo_password, uri)

                uri = '%s/%s' % (uri, self.mongo_db)

            uri = 'mongodb://%s' % uri

        return uri

    def connect(self):
        if self.connected:
            return True

        self.conn = Connection(self.uri, safe=True)
        self.db = self.conn[self.mongo_db]

        manipulators = self.db.incoming_manipulators
        manipulators += self.db.outgoing_manipulators

        for manipulator in manipulators:
            if isinstance(manipulator, CanopsisSONManipulator):
                break

        else:
            self.db.add_son_manipulator(
                CanopsisSONManipulator('_id')
            )

        try:
            self.gridfs_namespace = CONFIG.get("master", "gridfs_namespace")
        except:
            pass

        self.fs = gridfs.GridFS(self.db, self.gridfs_namespace)

        self.connected = True

        self.logger.debug("Connected %s" % id(self))

    def disconnect(self):
        if self.connected:
            self.conn.fsync()
            del self.conn
            self.connected = False

    def check_connected(self):
        """
        Check if self is connected to db.
        """
        if not self.connected:
            raise Exception("CSTORAGE is not connected %s" % id(self))

    def get_backend(self, namespace=None):
        self.check_connected()

        if not namespace:
            namespace = self.namespace

        try:
            backend = self.backend[namespace]
            self.logger.debug("Use %s collection" % namespace)

            return backend
        except:
            self.backend[namespace] = self.db[namespace]
            self.logger.debug("Connected to %s collection." % namespace)
            return self.backend[namespace]


    def update(self, _id, data, namespace=None, account=None):
        self.check_connected()

        if not isinstance(data, dict):
            raise Exception('Invalid data, must be a dict ...')

        data['crecord_write_time'] = int(time.time())

        # Check if record exist
        count = self.count({'_id': _id}, namespace=namespace, account=account, for_write=True)
        if count:
            backend = self.get_backend(namespace)
            backend.update({ '_id': self.clean_id(_id) }, { "$set": data });
        else:
            raise KeyError("'%s' not found ..." % _id)

    def put(self, _record_or_records, account=None, namespace=None, mset=False):
        self.check_connected()

        if not account:
            account = self.account

        records = []
        return_ids = []


        if isinstance(_record_or_records, Record):
            records = [_record_or_records]
        elif isinstance(_record_or_records, list):
            records = _record_or_records
        else:
            self.logger.error("Invalid record type")

        backend = self.get_backend(namespace)

        self.logger.debug("Put %s record(s) ..." % len(records))
        for record in records:
            _id = record._id

            access = False

            new_record = True
            oldrecord = None
            if _id:
                try:
                    oldrecord = self.get(
                        _id,
                        namespace=namespace,
                        account=self.root_account
                    )
                    new_record = False
                except:
                    pass

            if not new_record:
                self.logger.debug("Check rights of %s" % _id)
                if account.user == 'root':
                    access = True
                else:
                    access = oldrecord.check_write(account)
            else:
                # New record
                # Todo: check if account have right to create record ...
                access = True

            if not access:
                self.logger.error("Puts: Access denied ...")
                raise ValueError("Access denied")

            if new_record:
                # Insert new record
                self.logger.debug("Insert new record")

                try:
                    # Check if record have binary and store in grid fs
                    if getattr(record, 'binary', None):
                        record.data['binary_id'] = self.put_binary(
                            record.binary,
                            record.data['file_name'],
                            record.data['content_type']
                        )

                    record.write_time = int(time.time())
                    data = record.dump()
                    data['crecord_creation_time'] = record.write_time

                    # Del it if 'None'
                    if not data['_id']:
                        del data['_id']

                    if not _id:
                        _id = backend.insert(
                            data,
                            safe=self.mongo_safe,
                            w=1
                        )
                    else:
                        backend.update(
                            {'_id': _id},
                            data,
                            safe=self.mongo_safe,
                            upsert=True
                        )

                    self.logger.debug("Inserted (_id: '{}'')".format(_id))

                except Exception as err:
                    self.logger.error(
                        "Impossible to store !\nReason: {}".format(err)
                    )
                    self.logger.debug("Dump:\n {}".format(record.dump()))
                    raise ValueError("Impossible to insert {}".format(err))

                record._id = _id
                return_ids.append(_id)
            else:
                # Update record
                self.logger.debug("Update record '%s'" % _id)

                try:
                    # Check if record have binary and store in grid fs
                    if getattr(record, 'binary', None):
                        record.data['binary_id'] = self.put_binary(
                            record.binary,
                            record.data['file_name'],
                            record.data['content_type']
                        )

                    record.write_time = int(time.time())
                    data = record.dump()

                    del data['_id']
                    _id = self.clean_id(_id)

                    if mset:
                        ret = backend.update(
                            {'_id': _id},
                            {"$set": data},
                            upsert=True,
                            safe=self.mongo_safe
                        )
                    else:
                        ret = backend.update(
                            {'_id': _id},
                            data,
                            upsert=True,
                            safe=self.mongo_safe
                        )

                    if self.mongo_safe:
                        if ret['updatedExisting']:
                            self.logger.debug(
                                "Updated (_id: '{}')".format(_id)
                            )
                        else:
                            self.logger.debug("Saved (_id: '{}')".format(_id))

                except Exception as err:
                    self.logger.error(
                        "Impossible to store !\nReason: {}".format(err))
                    self.logger.debug("Record dump:\n{}".format(record.dump()))
                    raise ValueError("Impossible to store ({})".format(err))

                record._id = _id
                return_ids.append(_id)

        if len(return_ids) == 1:
            return return_ids[0]
        else:
            return return_ids

    def find_one(self, *args, **kargs):
        return self.find(one=True, *args, **kargs)

    def count(self, *args, **kargs):
        return self.find(count=True, *args, **kargs)

    def find(self, mfilter={}, mfields=None, account=None, namespace=None, one=False, count=False, sort=None, limit=0, offset=0, for_write=False, ignore_bin=True, raw=False, with_total=False):
        self.check_connected()

        if not account:
            account = self.account

        if isinstance(sort, basestring):
            sort = [(sort, 1)]

        # Clean Id
        if mfilter.get('_id', None):
            mfilter['_id'] = self.clean_id(mfilter['_id'])

        if one:
            sort = [('timestamp', -1)]

        self.logger.debug("Find records from mfilter" )

        (Read_mfilter, Write_mfilter) = self.make_mongofilter(account)

        if for_write:
            if Write_mfilter:
                mfilter = { '$and': [ mfilter, Write_mfilter ] }
        else:
            if Read_mfilter:
                mfilter = { '$and': [ mfilter, Read_mfilter ] }

        self.logger.debug(" + fields : %s" % mfields)
        self.logger.debug(" + mfilter: %s" % mfilter)

        backend = self.get_backend(namespace)

        if one:
            raw_records = backend.find_one(mfilter, fields=mfields, safe=self.mongo_safe)
            if raw_records:
                raw_records = [ raw_records ]
            else:
                raw_records = []
        else:

            count_limit_reached = backend.count() > self.no_count_limit

            if count_limit_reached:

                if limit == 0:
                    limit = self.fetch_limit

                if limit > 1:
                    #change limit artificially to fetch one more result if possible
                    limit += 1

            if sort is None:
                raw_records = backend.find(mfilter, fields=mfields, safe=self.mongo_safe, skip=offset, limit=limit)
            else:
                raw_records = backend.find(mfilter, fields=mfields, safe=self.mongo_safe, skip=offset, limit=limit, sort=sort)



            """
                Because mongo counts computation time is not acceptable, total is equal
                to the element fetched count (can be limit or less before it is artificially changed) OR
                total is offset + limit events and possibly + 1 if limit is reached
                (when +1 , this means some other records are availables)
            """

            if count_limit_reached:
                #When count limit reached, then count is done as described upper
                raw_records = list(raw_records)

                total = len(raw_records) + offset

                if limit > 1:
                    raw_records = raw_records[:limit -1]

            else:
                #Otherwise, count is done on the collection with given filter.
                total = raw_records.count()
                raw_records = list(raw_records)

            # process limit, offset and sort independently of pymongo because sort does not use index
            if count:
                return total



        records=[]

        if not mfields:
            for raw_record in raw_records:
                try:
                    # Remove binary (base64)
                    if ignore_bin and raw_record.get('media_bin', None):
                        del raw_record['media_bin']

                    if not raw:
                        records.append(Record(raw_record=raw_record))
                    else:
                        records.append(raw_record)

                except Exception as err:
                    ## Not record format ..
                    self.logger.error("Impossible parse record ('%s') !" % err)
        else:
            records = raw_records

        self.logger.debug("Found %s record(s)" % len(records))

        if one:
            if len(records) > 0:
                return records[0]
            else:
                return None
        else:
            if with_total: # returns the couple of records, total
                return records, total

            return records

    def get(self, _id_or_ids, account=None, namespace=None, mfields=None, ignore_bin=True):
        self.check_connected()

        if not account:
            account = self.account

        dolist = False
        if isinstance(_id_or_ids, list):
            _ids = _id_or_ids
            dolist = True
        else:
            _ids = [ _id_or_ids ]

        backend = self.get_backend(namespace)

        self.logger.debug(" + Get record(s) '%s'" % _ids)
        if not len(_ids):
            self.logger.debug("   + No ids")
            return []

        self.logger.debug("   + fields : %s" % mfields)

        self.logger.debug("   + Clean ids")
        _ids = [self.clean_id(_id) for _id in _ids]

        #Build basic filter
        (Read_mfilter, Write_mfilter) = self.make_mongofilter(account)

        if len(_ids) == 1:
            mfilter = {'_id': _ids[0]}
        else:
            mfilter = {'_id': {'$in': _ids }}

        mfilter = { '$and': [ mfilter, Read_mfilter ] }

        #self.logger.debug("   + mfilter: %s" % mfilter)
        records = []
        try:
            if len(_ids) == 1:
                raw_record = backend.find_one(mfilter, fields=mfields, safe=self.mongo_safe)

                # Remove binary (base64)
                if ignore_bin and raw_record and raw_record.get('media_bin', None):
                    del raw_record['media_bin']

                if raw_record and mfields:
                    records.append(raw_record)
                elif raw_record:
                    records.append(Record(raw_record=raw_record))
            else:
                raw_records = backend.find(mfilter, fields=mfields, safe=self.mongo_safe)

                if mfields:
                    records = [raw_record for raw_record in raw_records]
                else:
                    for raw_record in raw_records:
                        # Remove binary (base64)
                        if ignore_bin and raw_record.get('media_bin', None):
                            del raw_record['media_bin']

                        records.append(Record(raw_record=raw_record))

        except Exception as err:
            self.logger.error("Impossible get record '%s' !\nReason: %s" % (_ids, err))

        self.logger.debug(" + Found %s records" % len(records))
        if not len(records):
            raise KeyError("'%s' not found ..." % _ids)

        if len(_ids) == 1 and not dolist:
            return records[0]
        else:
            return records

    def remove(self, _id_or_ids, account=None, namespace=None):
        self.check_connected()

        if not account:
            account = self.account

        _ids = []

        if isinstance(_id_or_ids, Record):
            _ids = [ _id_or_ids._id ]
        elif isinstance(_id_or_ids, list):
            if len(_id_or_ids) > 0:
                if isinstance(_id_or_ids[0], Record):
                    for record in _id_or_ids:
                        _ids.append(record._id)
                else:
                    _ids = _id_or_ids
        else:
            _ids = [ _id_or_ids ]

        backend = self.get_backend(namespace)

        self.logger.debug("Remove %s record(s) ..." % len(_ids))
        for _id in _ids:
            self.logger.debug(" + Remove record '%s'" % _id)

            oid = self.clean_id(_id)
            if account.user == 'root':
                access = True
            else:
                try:
                    oldrecord = self.get(oid, account=account)
                except Exception as err:
                    raise ValueError("Access denied or id not found")

                access = oldrecord.check_write(account)

            if access:
                try:
                    backend.remove({'_id': oid}, safe=self.mongo_safe)
                except Exception as err:
                    self.logger.error("Impossible remove record '%s' !\nReason: %s" % (_id, err))

                self.logger.debug(" + Success removed")
            else:
                self.logger.error("Remove: Access denied ...")
                raise ValueError("Access denied ...")

    def map_reduce(self, mfilter_or_ids, mmap, mreduce, account=None, namespace=None):
        self.check_connected()

        if not account:
            account = self.account

        if   isinstance(mfilter_or_ids, dict):
            # mfilter
            mfilter = mfilter_or_ids
        elif isinstance(mfilter_or_ids, list):
            #ids
            mfilter = {'_id': {'$in': mfilter_or_ids }}
        else:
                self.logger.error("Invalid filter")
                raise ValueError("Invalid filter")

        backend = self.get_backend(namespace)

        (Read_mfilter, Write_mfilter) = self.make_mongofilter(account)
        #mfilter = dict(mfilter.items() + Read_mfilter.items())
        if Read_mfilter != {}:
            mfilter = { '$and': [ mfilter, Read_mfilter ] }

        output = {}
        if backend.find(mfilter).count() > 0:
            result = backend.map_reduce(mmap, mreduce, "mapreduce", query=mfilter)
            for doc in result.find():
                output[doc['_id']] = doc['value']
        else:
            self.logger.debug("Nor record matching filter")

        return output


    def drop_namespace(self, namespace):
        self.check_connected()

        self.db.drop_collection(namespace)

    def get_namespace_size(self, namespace=None):
        self.check_connected()

        if not namespace:
            namespace = self.namespace

        try:
            return self.db.command("collstats", namespace)['size']
        except:
            return 0

    def recursive_get(self, record, depth=0,account=None, namespace=None):
        self.check_connected()

        depth += 1
        childs = record.children
        if len(childs) == 0:
            return

        record.children = []

        for child in childs:
            # HACK: fix root_directory in UI !!!
            try:
                rec = self.get(child,account=account,namespace=namespace)
                self.recursive_get(rec, depth,account=account,namespace=namespace)
                record.children.append(rec)
            except Exception as err:
                self.logger.debug(err)


    def get_record_childs(self, record,account=None, namespace=None):
        self.check_connected()

        child_ids = record.children
        if len(child_ids) == 0:
            return []

        records = []
        for _id in child_ids:
            records.append(self.get(str(_id),account=account, namespace=namespace))

        return records


    def print_record_tree(self, record, depth=0):
        self.check_connected()

        depth+=1

        childs = record.children_record
        if len(childs) == 0:
            return

        if depth == 1:
            print "|-> " + str(record.name)

        for child in childs:

            prefix = ""
            for i in range(depth):
                prefix += "  "
            prefix += "|"
            for i in range(depth):
                prefix += "--"
            print prefix + "> " + str(child.name)

            self.print_record_tree(child, depth)


    def get_childs_of_parent(self, record_or_id, rtype=None, account=None):
        self.check_connected()

        if isinstance(record_or_id, Record):
            _id = record_or_id._id
        else:
            _id = record_or_id

        mfilter = {'parent': _id}

        if rtype:
            mfilter['crecord_type'] = rtype

        return self.find(mfilter, account=account)

    def get_parents_of_child(self, record_or_id, rtype=None, account=None):
        self.check_connected()

        if isinstance(record_or_id, Record):
            _id = record_or_id._id
        else:
            _id = record_or_id

        mfilter = {'child': _id}

        if rtype:
            mfilter['crecord_type'] = rtype

        return self.find(mfilter, account=account)

    def is_parent(self, parent_record, child_record):
        self.check_connected()

        if str(child_record._id) in parent_record.children:
            return True
        else:
            return False

    def put_binary(self, data, file_name, content_type):
        self.check_connected()

        bin_id = self.fs.put(data, filename=file_name, content_type=content_type)
        return bin_id

    def get_binary(self, _id):

        self.check_connected()

        cursor = self.fs.get(_id)

        return cursor.read()

    def remove_binary(self, _id):
        self.check_connected()

        try:
            self.fs.delete(_id)
        except Exception as err:
            self.logger.error(u'Error when remove binarie', err)

    def check_binary(self, _id):
        self.check_connected()

        return self.fs.exists(_id)

    def __del__(self):
        self.logger.debug("Object deleted. (namespace: %s)" % self.namespace)

## Cache storage
STORAGES = {}
def get_storage(namespace='object', account=None, logging_level=logging.INFO):
    global STORAGES
    if namespace not in STORAGES:
        if not account:
            account = Account()

        STORAGES[namespace] = Storage(account, namespace=namespace, logging_level=logging_level)

    return STORAGES[namespace]
