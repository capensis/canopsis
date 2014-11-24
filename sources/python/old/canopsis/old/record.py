# -*- coding: utf-8 -*-
#!/usr/bin/env python
# --------------------------------
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

from bson import objectid

from re import compile as re_compile

re_owner = re_compile("^account\..*")
re_group = re_compile("^group\..*")


class Record(object):
    def __init__(
        self, data=None, _id=None, name="noname", owner=None, group=None,
        raw_record=None, record=None, storage=None, account=None,
        admin_group=None, _type='raw'
    ):

        super(Record, self).__init__()

        self.write_time = None

        self.owner = None
        self.group = None

        self.chown(owner)
        self.chgrp(group)

        if data is None:
            data = {}

        self.admin_group = admin_group
        self.type = _type
        self.access_owner = ['r', 'w']
        self.access_group = ['r']
        self.access_other = []
        self.access_unauth = []
        self.name = name
        self.parent = []
        self.children = []
        self.children_record = []
        self._id = _id
        self.enable = True
        self.binary = None

        if account:
            self.chown(account.user)
            self.chgrp(account.group)

        #set admin account
        if not self.admin_group:
            self.admin_group = 'group.CPS_%s_admin' % self.type

        try:
            self._id = data['_id']
            del data['_id']
        except:
            pass

        self.data = data.copy()
        self.storage = storage

        if isinstance(record, Record):
            self.load(record.dump())
        elif raw_record:
            self.load(raw_record)

    def load(self, dump):

        self.owner = str(dump.get('aaa_owner', ''))
        self.group = str(dump.get('aaa_group', ''))
        self.access_owner = dump.get('aaa_access_owner', '')
        self.access_group = dump.get('aaa_access_group', '')
        self.access_other = dump.get('aaa_access_other', '')
        self.access_unauth = dump.get('aaa_access_unauth', '')
        self.type = str(dump.get('crecord_type', ''))
        self.write_time = dump.get('crecord_write_time', '')
        self.name = dump.get('crecord_name', '')
        self.children = dump.get('children', '')
        self.parent = dump.get('parent', '')
        self.enable = dump.get('enable', '')

        if not dump.get('crecord_creation_time', None):
            dump['crecord_creation_time'] = self.write_time

        if 'aaa_admin_group' in dump:
            self.admin_group = str(dump['aaa_admin_group'])
            del dump['aaa_admin_group']
        else:
            self.admin_group = 'group.CPS_%s_admin' % self.type

        #security
        if not self.access_owner:
            self.access_owner = []
        if not self.access_group:
            self.access_group = []
        if not self.access_other:
            self.access_other = []
        if not self.access_unauth:
            self.access_unauth = []

        dump['_id'] = dump.get('_id', None)

        self._id = dump['_id']

        dump.pop('_id', '')
        dump.pop('enable', '')
        dump.pop('aaa_owner', '')
        dump.pop('aaa_group', '')
        dump.pop('aaa_access_owner', '')
        dump.pop('aaa_access_group', '')
        dump.pop('aaa_access_other', '')
        dump.pop('aaa_access_unauth', '')
        dump.pop('crecord_type', '')
        dump.pop('crecord_write_time', '')
        dump.pop('crecord_name', '')
        dump.pop('children', '')
        dump.pop('parent', '')

        self.data = dump.copy()

    def save(self, storage=None):
        if not storage:
            if not self.storage:
                raise Exception('For save you must specify storage')
            else:
                storage = self.storage

        return storage.put(self)

    def get(self, key):
        if hasattr(self, key):
            return getattr(self,key)
        else:
            return None

    def dump(self, json=False):
        dump = self.data.copy()


        dump['_id'] = self.get('_id')
        dump['aaa_owner'] = self.get('owner')
        dump['aaa_group'] = self.get('group')
        dump['aaa_access_owner'] = self.get('access_owner')
        dump['aaa_access_group'] = self.get('access_group')
        dump['aaa_access_other'] = self.get('access_other')
        dump['aaa_access_unauth'] = self.get('access_unauth')
        dump['crecord_type'] = self.get('type')
        dump['crecord_write_time'] = self.get('write_time')
        dump['crecord_name'] = self.get('name')

        if 'enable' not in self.data:
            dump['enable'] = self.enable

        dump['parent'] = self.get('parent')
        dump['children'] = self.get('children')

        dump['aaa_admin_group'] = self.get('admin_group')

        if json:
            # Clean objectid
            for key in dump:
                if isinstance(dump[key], objectid.ObjectId):
                    dump[key] = str(dump[key])

            items = []
            for item in dump['parent']:
                items.append(str(item))
            dump['parent'] = list(items)

            items = []
            for item in dump['children']:
                items.append(str(item))
            dump['children'] = list(items)

        return dump

    def recursive_dump(self, json=False):
        dump = self.dump(json=json)
        dump['children'] = []

        for child in self.children:
            if isinstance(child, Record):
                formated = child.recursive_dump(json=json)
                dump['children'].append(formated)

        return dump

    def cat(self, dump=False):
        for_str = False

        if dump:
            data = self.dump()
        else:
            data = self.data.copy()

        output = ""
        for key in data.keys():
            try:
                output += key + ": " + str(data[key]) + "\n"
            except:
                output += key + ": " + data[key] + "\n"

        if for_str:
            return output
        else:
            print(output)

    def __str__(self):
        return str(self.dump())

    def check_write(self, account):
        return True
        if account:
            if account.user == 'root' or account.group == 'group.CPS_root' \
                    or 'group.CPS_root' in account.groups:
                return True

            elif account._id == self.owner and 'w' in self.access_owner:
                return True
            elif account.group == self.group and 'w' in self.access_group:
                return True
            elif self.group in account.groups and 'w' in self.access_group:
                return True
            elif self.admin_group in account.groups \
                    or self.admin_group == account.group:
                return True
        return False

    def chown(self, owner):
        #if isinstance(owner, caccount):
        #   self.owner = owner.user
        #   self.group = owner.group
        #else:
        #   self.owner=owner
        if owner:
            if re_owner.match(str(owner)):
                self.owner = owner
            else:
                self.owner = "account.%s" % owner

    def chgrp(self, group):
        if group:
            if re_group.match(str(group)):
                self.group = group
            else:
                self.group = "group.%s" % group

    def chmod(self, action):
        ## g+w, g+r, u+r, u+w ...
        # u: user
        # g: group
        # o: other
        # a: anonymous
        if not (len(action) < 3):
            field = action[0]
            way = action[1]
            mod = action[2]
            access = None

            #print "Field:", field, "Way:", way, "Mod:", mod
            if field == 'u':
                access = self.access_owner
            elif field == 'g':
                access = self.access_group
            elif field == 'o':
                access = self.access_other
            elif field == 'a':
                access = self.access_unauth

            #print "Before:", access
            if access is not None:
                if way == '+':
                    if mod not in access:
                        access.append(mod)
                elif way == '-':
                    if mod in access:
                        access.remove(mod)

            #print "After", action ,":", access
        else:
            raise ValueError("Invalid argument ...")

    def add_children(self, record, autosave=True):
        _id = record._id

        if autosave:
            if not _id:
                record.save()
            if not self._id:
                self.save()

        if not _id or not self._id:
            raise ValueError(
                "You must save all records before this operation ...")

        if str(_id) not in self.children:
            self.children.append(str(_id))
            record.parent.append(str(self._id))
            if autosave and self.storage and record.storage:
                self.save()
                record.save()

    def remove_children(self, record, autosave=True):
        _id = record._id

        if autosave:
            if not _id:
                record.save()
            if not self._id:
                self.save()

        if not _id or not self._id:
            raise ValueError(
                "You must save all records before this operation ...")

        if str(_id) in self.children:
            self.children.remove(str(_id))
            record.parent.remove(str(self._id))
            if autosave and self.storage and record.storage:
                self.save()
                record.save()

    def is_parent(self, record):
        if str(record._id) in self.children:
            return True
        else:
            return False

    def is_enable(self):
        return self.enable

    def set_enable(self, autosave=True):
        self.enable = True
        if autosave and self.storage:
            self.save()

    def set_disable(self, autosave=True):
        self.enable = False
        if autosave and self.storage:
            self.save()


def access_to_str(access):
    output = ''

    if 'r' in access:
        output += 'r'
    else:
        output += '-'

    if 'w' in access:
        output += 'w'
    else:
        output += '-'

    return output
