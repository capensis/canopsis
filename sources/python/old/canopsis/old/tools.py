#!/usr/bin/env python
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

from re import compile as re_compile, search

from os import listdir
from os.path import expanduser

from logging import getLogger

from math import pow

from sys import path

legend = ['ok', 'warning', 'critical', 'unknown']

logger = getLogger('ctools')

internal_metrics = [
    'cps_state',
    'cps_statechange',
    'cps_statechange_nok',
    'cps_statechange_0',
    'cps_statechange_1',
    'cps_statechange_2',
    'cps_statechange_3',
    'cps_statechange_hard',
    'cps_statechange_soft',
    'cps_statechange_component',
    'cps_statechange_service',
    'cps_statechange_service_from_component',
    'cps_evt_per_sec',
    'cps_sec_per_evt',
    'cps_queue_size',
    'cps_sel_state_0',
    'cps_sel_state_1',
    'cps_sel_state_2',
    'cps_sel_state_3'
]


def metric_is_internal(metric_name):
    return metric_name.startswith('cps_')


def calcul_pct(data, total=None):
    if not total:
        ## Get total
        total = 0
        for key in data.keys():
            value = data[key]
            total += value
    ## Calc pct
    data_pct = {}
    for key in data.keys():
        value = data[key]
        data_pct[key] = round(((float(value) * 100) / float(total)), 2)

    ## Fix empty value
    for key in legend:
        try:
            value = data_pct[key]
        except:
            data_pct[key] = 0

    return data_pct

RE_PERF_DATA = re_compile(
    "('?([0-9A-Za-z/\\\:\.%%\-{}\?\[\]_ ]*)'?=(\-?[0-9.,]*)(([A-Za-z%%/]*))(;@?(\-?[0-9.,]*):?)?(;@?(\-?[0-9.,]*):?)?(;@?(\-?[0-9.,]*):?)?(;@?(\-?[0-9.,]*):?)?(;? ?))")


def parse_perfdata(perf_data_raw):
        # 'label'=value[UOM];[warn];[crit];[min];[max]
        #   load1=0.440     ;5.000 ;10.000;0    ;

        logger.debug("Parse: %s" % perf_data_raw)

        perfs = RE_PERF_DATA.split(perf_data_raw)

        perf_data_array = []
        perf_data = {}
        i = 0
        for info in perfs:
            if info == '':
                info = None

            #logger.debug(" + %s '%s'" % (i, info))
            try:
                if info and i == 2:
                    perf_data['metric'] = info
                elif info and i == 3:
                    perf_data['value'] = info.replace(',', '.')
                elif info and i == 4:
                    perf_data['unit'] = info
                elif info and i == 7:
                    perf_data['warn'] = info.replace(',', '.')
                elif info and i == 9:
                    perf_data['crit'] = info.replace(',', '.')
                elif info and i == 11:
                    perf_data['min'] = info.replace(',', '.')
                elif info and i == 13:
                    perf_data['max'] = info.replace(',', '.')

                i += 1
                if i is 15:
                    try:
                        perf_data_clean = {}
                        for key in perf_data.keys():
                            if perf_data[key]:
                                try:
                                    perf_data_clean[key] = float(
                                        perf_data[key])
                                except:
                                    if key == 'metric' or key == 'unit':
                                        perf_data_clean[key] = perf_data[key]
                                    else:
                                        logger.debug(
                                            "Invalid value, '%s' = '%s'" % (
                                                key, perf_data[key]))

                        if 'value' in perf_data_clean \
                                and 'metric' in perf_data_clean:
                            perf_data_array.append(perf_data_clean)

                        if not perf_data_clean.get('unit', None):
                            # split: g[in_bps]= ...
                            metric_ori = perf_data_clean['metric']
                            if metric_ori[len(metric_ori) - 1] == ']':
                                metric_ori = metric_ori[:len(metric_ori) - 1]
                                metric = metric_ori.split('[', 1)
                                if len(metric) == 2:
                                    perf_data_clean['metric'] = metric[0]
                                    perf_data_clean['unit'] = metric[1]

                        logger.debug(" + %s" % perf_data_clean)

                    except Exception as err:

                        logger.error("perf_data: Raw: %s" % perf_data_raw)
                        logger.error(
                            "perf_data: Impossible to clean '%s': %s" % (
                                perf_data, err))

                    perf_data = {}
                    i = 0

            except Exception as err:
                logger.error(
                    "perf_data: Invalid metric %s: %s (%s)" % (i, info, err))

        return perf_data_array


def dynmodloads(_path=".", subdef=False, pattern=".*"):
    loaded = {}
    _path = expanduser(_path)
    logger.debug("Append path '%s' ..." % _path)
    path.append(_path)

    try:
        for mfile in listdir(_path):
            try:
                ext = mfile.split(".")[1]
                name = mfile.split(".")[0]

                if name != "." and ext == "py" and name != '__init__':
                    logger.info("Load '%s' ..." % name)
                    try:

                        module = __import__(name)
                        loaded[name] = module

                        if subdef:
                            alldefs = dir(module)
                            for mydef in alldefs:
                                if mydef not in [
                                        "__builtins__",
                                        "__doc__",
                                        "__file__",
                                        "__name__",
                                        "__package__"]:
                                    if search(pattern, mydef):
                                        logger.debug(
                                            " + From %s import %s ..." % (
                                                name, mydef))
                                        exec("loaded[mydef] = module.%s" %
                                            mydef)

                        logger.debug(" + Success")
                    except Exception as err:
                        logger.error("\t%s" % err)
            except:
                pass
    except Exception as err:
        logger.error(err)

    return loaded


def Str2Number(string):
    value = float(string)

    if int(value) == value:
        value = int(value)

    return value


def uniq(alist):
    """
    remove duplicate entry
    """
    set = {}
    return [set.setdefault(e, e) for e in alist if e not in set]


def cleanTimestamp(timestamp):
    if len(str(timestamp)) > 12:
        return int(timestamp) / 1000
    else:
        return int(timestamp)


def roundSignifiantDigit(value, sig):
    mult = pow(10, sig)
    value = round(value * mult)
    value = value / mult
    return value
