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

from __future__ import print_function

from argparse import ArgumentParser

from datetime import datetime as dt

from calendar import timegm

from pyperfstore2.store import store as Store

from sys import getsizeof, maxsize

from collections import OrderedDict

DTFORMAT = '%Y-%m-%d'

SIZE = 'size'  #: perfdata size name.
FTS = 'fts'  #: perfdata first timestamp name.
LTS = 'lts'  #: perfdata last timestamp name.

HUMANREADABLEUNITBYQUANTITY = OrderedDict([
    (10**15, 'P'),
    (10**12, 'T'),
    (10**9, 'G'),
    (10**6, 'M'),
    (10**3, 'K')
])

def gethumanreadablevalue(value):
    """Get the right human readable value.

    :param float value: value to convert to a human readable value.
    :rtype: str"""
    unit = ''

    for bound in HUMANREADABLEUNITBYQUANTITY:
        if value >= bound:
            unit = HUMANREADABLEUNITBYQUANTITY[bound]
            break

    result = '{0}o'.format(unit)

    return result


def displaymetrics(store, metrics=None, beforets=None, humanreadable=False):
    """Display size and time window of metrics.

    :param Store store: pyperfstore2 store.
    :param str metrics: metric name regex.
    :param float beforets: last timestamp from when getting information.
    :param bool humanreadable: if True (False by default), display size
        information with o/Ko/Mo/Go/Po."""

    mfilter = getfilter(metrics, beforets)

    totalsize = 0

    documents = store.collection.find(mfilter)

    documents_by_id = {}

    for document in documents:
        idmetrics = documents_by_id.setdefault(document['_id'], {})

        idmetrics[SIZE] = getsizeof(document) + idmetrics.get(SIZE, 0)
        idmetrics[LTS] = max(document[LTS], idmetrics.get(LTS, 0))
        idmetrics[FTS] = min(document[FTS], idmetrics.get(FTS, maxsize))

    fromts = dt.fromtimestamp

    for _id in documents_by_id:
        info = documents_by_id[_id]

        size = info[SIZE]
        totalsize += size

        if humanreadable:
            size = humanreadable(size)

        print(
            '{0}: size={1}, fts={2}, lst={3}'.format(
                _id, size, fromts(info[FTS]), fromts(info[LTS])
            )
        )

    if humanreadable:
        totalsize = humanreadable(totalsize)

    print('total: count={0}, size={1}'.format(len(documents_by_id), totalsize))


def delmetrics(store, metrics=None, beforets=None):
    """Delete metrics.

    :param str metrics: metric name regex.
    :param float beforets: last timestamp from when deleting information."""

    mfilter = getfilter(metrics, beforets)

    store.collection.remove(mfilter)


def getfilter(metrics=None, beforets=None):
    """Get store filter related to input metrics and beforets.

    :param str metrics: metric name regex.
    :param float beforets: last timestamp from when deleting information.
    :rtype: dict"""

    result = {}

    if metrics is not None:
        result['_id'] = {'$regex': metrics}

    if beforets is not None:
        result['lts'] = {'$lte': beforets}

    return result


def main():
    """Main script process."""

    # build argument parser
    epilog = '\
    How to... \
    \nGet time window and size by metric name in DB about metrics named like "pluton":\
    \n\tperfstore2_cleaner --metrics ".*pluton.*"\
    \nDelete all perfdata before the December 25th 2015 and display size information:\
    \n\tperfstore2_cleaner --before 2015-12-25 --delete'

    parser = ArgumentParser(
        description='Get/delete Canopsis-ficus perfdata.', epilog=epilog
    )
    # set metric names
    parser.add_argument(
        '-m', '--metrics',
        help='metric regex name. Corresponds to entity names. All perfdata are considered if not given.',
        type=str, dest='metrics', required=False
    )
    # set delete
    parser.add_argument(
        '--delete', help='Delete data. Ask before deleting if --force is not given.',
        dest='delete', required=False, action='store_true', default=False
    )
    # set force
    parser.add_argument(
        '--force', help='In addition to the --delete argument, delete perfdata without secondary agreement.',
        required=False, action='store_true', default=False, dest='force'
    )
    # set before
    parser.add_argument(
        '-b', '--before', help='End suppression/information date. Must respect the format Y-m-d',
        dest='before', required=False
    )
    # set human readable
    parser.add_argument(
        '-h', '--humanreadable', help='Human readable value for size properties.',
        dest='humanreadable', required=False, action='store_true', default=False
    )

    # parse arguments
    args = parser.parse_args()

    if args.before:
        before = dt.strptime(args.before, DTFORMAT)

    else:
        before = dt.now()

    beforets = timegm(before.timetuple())

    store = Store()

    displaymetrics(
        store=store, metrics=args.metrics, beforets=beforets,
        humanreadable=args.humanreadable
    )

    if args.delete:
        if args.force or input('Delete all those data ? (y/N)') in ('y', 'Y'):
            delmetrics(store=store, metrics=args.metrics, beforets=beforets)


if __name__ == '__main__':
    main()
