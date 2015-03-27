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
from canopsis.old.storage import get_storage
from canopsis.old.account import Account

import collections
import json
from bson.json_util import dumps


class ConsoFactory(object):
    FILE_NAME = "object.json"
    SUBSTRING = "consolidation"
    QUERY_CONSO = "{'crecord_type':'consolidation'}"
    QUERY_COMP = "{{'source_type':'component', 'component':{!r}}}"
    QUERY_RESR = "{{'source_type':'resource', 'component':{!r}}}"
    NAMESPACE = ('objectv1', 'eventsv1')
    OPERATOR = 'and'
    COMP = 'co'
    RESR = 're'
    METR = 'me'
    CONSO_METHOD = ["average", "sum", "delta", "min", "max"]
    # formula
    # consolidation_method -> formula
    # mfilter -> metric
    ORDINAR_KEYS = ("_id", "crecord_write_time", "enable", "loaded", "crecord_name", "crecord_type", "crecord_creation_time")
    STATIC_KEYS = ("max", "warn", "aggregate_interval", "crit", "min", "unit")
    ADVANCED_KEYS = ("metric", "mfilter")

    def __init__(self):
        super(ConsoFactory, self).__init__()
        self.data = self.loads(self.NAMESPACE[0], self.QUERY_CONSO, regex=None)

    @classmethod
    def storage_connection(self, namespace):
        '''
        Get the cconnection to canopsis DataBase.

        :param namespace: specify teh collection name.

        :return: a connection.
        :rtype: Mongo Connection.
        '''
        connection = get_storage(
            namespace=namespace,
            account=Account(
                user="root",
                group="root"
            )
        ).get_backend()
        return connection

    def loads(self, namespace, q=None, regex=None):
        '''
         Serialize cursor data into JSON.
         :param kind: specify which request should be running in the DataBase.

         :return: a dictionnary of elements.
         :rtype: dictionnary.
        '''
        _str = self.dump(namespace, q, regex)
        comp_json = json.loads(_str)
        if len(comp_json) > 0:
            # catch exception here
            res = comp_json
        else:
            res = {}
        return res

    def dump(self, namespace, q=None, regex=None):
        '''
        '''
        return dumps(self.get_data(namespace, q, regex))

    def get_data(self, namespace, q=None, regex=None):
        '''
         Access MongoDB and load topology or events data.

         :param kink: specify which request should be running in the DataBase.
         :return: a cursor of topology or events.
         :rtype: Cursor of elements dictionnary or empty dictionnary.
        '''
        print 'le regex: ', regex
        if regex:
            print 'il match ', regex
            q = q.format(regex)
            q = q.replace("\"", "")
            print 'Q******************* ', q
        json_acceptable = q.replace("'", "\"")

        print json_acceptable
        query = json.loads(json_acceptable)
        cursor = self.storage_connection(namespace).find(query)
        return cursor

    def conso_data_json(self):
        consos = []
        with open(self.FILE_NAME) as f:
            for line in f:
                _dict = json.loads(line)
                if _dict["crecord_type"].lower() == self.SUBSTRING:
                    line = line.replace("$oid", "oid")
                    _dict = json.loads(line)
                    consos.append(_dict)

        return consos

    def factory(self, conso):
        serie = {}
        #print conso
        # retreive data
        for k in self.ORDINAR_KEYS:
            try:
                serie[k] = conso[k]
            except KeyError:
                serie["loaded"] = True
                #raise e

        # Set the static keys
        for k in self.STATIC_KEYS:
            serie[k] = unicode(None)

        # Set the aggregate_interval
        try:
            serie["aggregate_interval"] = conso["aggregation_interval"]
        except KeyError:
            serie["aggregate_interval"] = 600

        # Set the aggregate_method
        try:
            serie["aggregate_method"] = conso["aggregation_method"]
        except KeyError:
            serie["aggregate_method"] = None

        # Convert regex here
        all_metrics = self.convert_regex_to_metrics(conso['mfilter'])
        # Set metrcis data
        #serie["metrics"] = all_metrics
        #formula = self.build_formula(serie['consolidation_method'], all_metrics)
        # Set formula
        #serie['serie'] = formula

        return serie

    def convert_regex_to_metrics(self, mfilter):
        mfilter = json.loads(mfilter)
        component = mfilter[self.OPERATOR][0][self.COMP]
        resource = mfilter[self.OPERATOR][1][self.RESR]
        metric = mfilter[self.OPERATOR][2][self.METR]
        #components, resources, metrics = [], [], []
        if isinstance(component, dict):
            print '---------------> ', component
            if 'regex' in component:
                components = self.run_regex(str(component['regex']), self.COMP)
            else:
                regex = self.clean_regex(str(component))
                components = self.run_regex(regex, self.COMP)
        elif isinstance(component, unicode):
            components = component
        if isinstance(resource, dict):
            if 'regex' in resource:
                resources = self.run_regex(str(resource['regex']), self.RESR)
            else:
                regex = self.clean_regex(str(resource['in']))
                resources = self.run_regex(regex, self.RESR)
        elif isinstance(resource, unicode):
            resources = resource
        if isinstance(metric, dict):
            if 'regex' in metric:
                metrics = self.run_regex(str(metric['regex']), self.METR)
            else:
                regex = self.clean_regex(str(metric['in']))
                metrics = self.run_regex(regex, self.METR)
        elif isinstance(metric, unicode):
            metrics = metric

        all_metrics = self.build_metrics(components, resources, metrics)
        print all_metrics
        print '\n\n\n\n\n'

        return all_metrics

    def run_regex(self, regex, identifier):
        result = []
        if identifier == self.COMP:
            print (regex)
            jsons = self.loads(self.NAMESPACE[1], self.QUERY_COMP, regex)
            result = self.get_att(jsons)
        if identifier == self.RESR:
            jsons = self.loads(self.NAMESPACE[1], self.QUERY_RESR, regex)
            result = self.get_att(jsons)
        if identifier == self.METR:
            pass

        return result

    def get_att(self, jsons):
        data = []
        for d in jsons:
            data.append(d['component'])

        return data

    def clean_regex(self, data):
        data = data.replace("in", "$in")
        data = data.replace("u\'", "'")
        return data

    def build_metrics(self, components, resources, metrics):
        c_metrics = []
        c_metric = ""
        if isinstance(components, collections.Iterable):
            for c in components:
                if isinstance(resources, collections.Iterable):
                    for r in resources:
                        if isinstance(metrics, collections.Iterable):
                            for m in metrics:
                                c_metric = "/" + c + "/" + r + "/" + m
                                c_metrics.append(c_metric)
                        else:
                            c_metric = "/" + c + "/" + r + "/" + metrics
                            c_metrics.append(c_metric)
                else:
                    if isinstance(metrics, collections.Iterable):
                        for m in metrics:
                            c_metric = "/" + c + "/" + resources + "/" + m
                            c_metrics.append(c_metric)
                    else:
                        c_metric = "/" + c + "/" + resources + "/" + metrics
                        c_metrics.append(c_metric)
        else:
            if isinstance(resources, collections.Iterable):
                for r in resources:
                    if isinstance(metrics, collections.Iterable):
                        c_metric = "/" + components + "/" + r + "/" + m
                        c_metrics.append(c_metric)
                    else:
                        c_metric = "/" + components + "/" + r + "/" + metrics
                        c_metrics.append(c_metric)
            else:
                if isinstance(metrics, collections.Iterable):
                    for m in metrics:
                        c_metric = "/" + components + "/" + resources + "/" + m
                        c_metrics.append(c_metric)
                else:
                    c_metric = "/" + components + "/" + resources + "/" + metrics
                    c_metrics.append(c_metric)
        return c_metrics

    def build_formula(self, formula, metrics):
        formula = str(formula) + "(" + ",".join(metrics) + ")"
        return formula

    def build(self):
        #consos = self.conso_data_json()
        consos = self.data
        series = []
        for conso in consos:
            serie = self.factory(conso)
            series.append(serie)

        return series


if __name__ == '__main__':
    c = ConsoFactory()
    c.build()
    """
    for s in c.build():
        print s
    """
