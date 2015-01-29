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

import unittest

from canopsis.old.tools import calcul_pct
from canopsis.old.tools import parse_perfdata

import time

import logging

logging.basicConfig(level=logging.DEBUG)


def check_perfdata(perf_data_raw):
    perf_data = parse_perfdata(perf_data_raw)

    try:
        perf_data[0]['value']
        perf_data[0]['metric']
    except Exception as err:
        raise Exception('Impossible to parse perfdata:\n + err: %s\n + raw: %s\n + result: %s' % (err, perf_data_raw, perf_data))

    return perf_data


class KnownValues(unittest.TestCase):
    def setUp(self):
        pass

    def test_01_Perfdata(self):
        result = [{'warn': 5.0, 'crit': 10.0, 'metric': 'load1', 'value': 0.44, 'min': 0.0}, {'warn': 4.0, 'crit': 6.0, 'metric': 'load5', 'value': 0.19, 'min': 0.0}, {'warn': 3.0, 'crit': 4.0, 'metric': 'load15', 'value': 0.13, 'min': 0.0}]
        perf_data = "load1=0.440;5.000;10.000;0; load5=0.190;4.000;6.000;0; load15=0.130;3.000;4.000;0;"
        perf_data = check_perfdata(perf_data)
        if perf_data != result:
            print(perf_data)
            raise Exception('[1] Error in perfdata parsing ...')

        result = [ {'metric': 'load1', 'value': 0.440} ]
        perf_data = "load1=0.440"
        perf_data = check_perfdata(perf_data)
        if perf_data != result:
            print(perf_data)
            raise Exception('[2] Error in perfdata parsing ...')

        result = [{'min': 0.0, 'max': 100.0, 'metric': 'ok', 'value': 100.0, 'warn': 98.0, 'crit': 95.0, 'unit': '%'}, {'min': 0.0, 'max': 100.0, 'metric': 'warn', 'value': 0.0, 'warn': 0.0, 'crit': 0.0, 'unit': '%'}, {'min': 0.0, 'max': 100.0, 'metric': 'crit', 'value': 0.0, 'warn': 0.0, 'crit': 0.0, 'unit': '%'}]
        perf_data = "'ok'=100.0%;98;95;0;100 'warn'=0%;0;0;0;100 'crit'=0%;0;0;0;100"
        perf_data = check_perfdata(perf_data)
        if perf_data != result:
            print(perf_data)
            raise Exception('[3] Error in perfdata parsing ...')

        result = [{'min': 0.0, 'max': 100.0, 'metric': 'C:/ Used', 'value': 100.0, 'warn': 98.0, 'crit': 95.0, 'unit': '%'}, {'min': 0.0, 'max': 100.0, 'metric': 'warn: /ing', 'value': 0.0, 'warn': 0.0, 'crit': 0.0, 'unit': '%'}, {'min': 0.0, 'max': 100.0, 'metric': 'D:\\ Used', 'value': 0.0, 'warn': 0.0, 'crit': 0.0, 'unit': '%'}]
        perf_data = "'C:/ Used'=100.0%;98;95;0;100; 'warn: /ing'=0%;0;0;0;100; 'D:\ Used'=0%;0;0;0;100;"
        perf_data = check_perfdata(perf_data)
        if perf_data != result:
            print(perf_data)
            raise Exception('[4] Error in perfdata parsing ...')

        perf_data = "/=541MB;906;956;0;1007 /home=62MB;452;477;0;503 /tmp=38MB;906;956;0;1007 /var=943MB;1813;1914;0;2015 /usr=2249MB;5441;5743;0;6046 /opt=68MB;112;118;0;125 /boot=25MB;90;95;0;100 /backup=4410MB;7256;7659;0;8063 /mnt/NAS_OPU_LIVRAISON=35603MB;36300;38317;0;40334 /products/admin=595MB;906;956;0;1007 /products/oracle/10.2.0=146MB;7256;7659;0;8063 /products/agtgrid=851MB;6349;6702;0;7055 /app/PPOPVGL=10966MB;27213;28725;0;30237"
        perf_data = check_perfdata(perf_data)

        perf_data = "/=217MB;460;486;0;512 /usr=2757MB;3052;3222;0;3392 /var=397MB;921;972;0;1024 /tmp=437MB;2275;2401;0;2528 /home=243MB;460;486;0;512 /admin=0MB;115;121;0;128 /opt=489MB;691;729;0;768 /var/adm/ras/livedump=4MB;460;486;0;512 /products=2007MB;9216;9728;0;10240 /products/logs=0MB;921;972;0;1024 /products/admin=216MB;921;972;0;1024 /data/PPOPBGL/oracle/PPOPBGL1/nfsbck=32226MB;46771;49369;0;51968 /data/PPOPBGL/oracle/PPOPBGL2/nfsbck=32226MB;46771;49369;0;51968 /app/P1TECGL/P1OPBGLDB01=409MB;921;972;0;1024 /app/admin=0MB;921;972;0;1024 /app/agtgrid=2433MB;2764;2918;0;3072 /data/PPOPBGL/oracle/PPOPBGL1/admin= /data/PPOPBGL/oracle/PPOPBGL1/backup=55MB;1843;1945;0;2048 /data/PPOPBGL/oracle/PPOPBGL1/dbs/archivelog=6605MB;74628;78774;0;82920 /data/PPOPBGL/oracle/PPOPBGL1/dbs/data01=9263MB;94924;100198;0;105472 /data/PPOPBGL/oracle/PPOPBGL1/dbs/flash_recovery_area=0MB;3600;3800;0;4000 /data/PPOPBGL/oracle/PPOPBGL1/dbs/index0= /data/PPOPBGL/oracle/PPOPBGL1/dbs/redolog01="
        perf_data = check_perfdata(perf_data)

        perf_data = "load1=0.040;3.000;4.000;0; load5=0.070;4.000;5.000;0; load15=0.010;4.000;5.000;0;"
        perf_data = check_perfdata(perf_data)

        perf_data = "'C:\ %'=67%;10;5; 'C:\'=6.65G;2;1;0;20; 'E:\ %'=31%;10;5; 'E:\'=10870.0M;1567.45;783.72;0;15674.52; 'J:\ %'=88%;10;5; 'J:\'=2.57G;2;1;0;20.00; 'T:\ %'=100%;10;5; 'T:\'=49.53M;512.3;256.19;0;5123.90;"
        perf_data = check_perfdata(perf_data)

        perf_data = "'C:\ Used'=55.4%;90;95; 'S:\ Used'=27.2%;90;95; 'E:\ Used'=24.6%;90;95; 'Q:\ Used'=1.1%;90;95; 'F:\ Used'=83.3%;90;95; 'F:\Journaux\ Used'=0.5%;90;95; 'F:\Data\ Used'=2.0%;90;95; 'F:\App\ Used'=1.4%;90;95; 'G:\ Used'=62.6%;90;95; 'G:\Journaux\ Used'=21.6%;90;95; 'G:\App\ Used'=7.8%;90;95; '\\?\Volume{c75733e9-4327-11e0-8596-0010184d9c22}\ Used'=35.9%;90;95; 'H:\ Used'=76.6%;90;95; 'H:\Journaux\ Used'=7.2%;90;95; '\\?\Volume{b7cdcbde-4fe5-11e0-a4e6-0010184d9c22}\ Used'=49.6%;90;95; 'H:\App\ Used'=9.9%;90;95; 'I:\ Used'=46.7%;90;95; 'I:\Journaux\ Used'=8.7%;90;95; '\\?\Volume{b7cdcbe1-4fe5-11e0-a4e6-0010184d9c22}\ Used'=99.0%;90;95; 'I:\App\ Used'=10.8%;90;95; 'J:\ Used'=75.2%;90;95; 'J:\Journaux\ Used'=7.0%;90;95; 'J:\App\ Used'=10.0%;90;95; 'J:\Data\ Used'=74.2%;90;95; 'K:\ Used'=77.3%;90;95; 'K:\Journaux\ Used'=8.2%;90;95; 'K:\Data\ Used'=84.6%;90;95; 'K:\App\ Used'=10.2%;90;95; 'G:\Data\ Used'=25.7%;90;95; 'H:\Data\ Used'=41.4%;90;95; 'I:\Data\ Used'=46.1%;90;95;"
        perf_data = check_perfdata(perf_data)

        perf_data = "'C_Used'=55.4%;90;95;"
        perf_data = check_perfdata(perf_data)
        if perf_data[0]['metric'] != 'C_Used':
            raise Exception('[5] Error in perfdata parsing ...')

        perf_data = "redo_log_file_switch_interval=6s;360:;60: /data/XXXX/oracle/XXXXX/dbs/archivelog=2059MB;30686;32391;0;34096"
        perf_data = check_perfdata(perf_data)
        if perf_data[0]['warn'] != 360 or perf_data[0]['crit'] != 60 or perf_data[1]['metric'] != '/data/XXXX/oracle/XXXXX/dbs/archivelog' or perf_data[1]['warn'] != 30686 or perf_data[1]['crit'] != 32391 :
            raise Exception('[6] Error in perfdata parsing ...')

        perf_data = "test=6s;360:120;@60:23;90;1000"
        perf_data = check_perfdata(perf_data)

        perf_data = "test=6s;360;~:60;90;1000"
        perf_data = check_perfdata(perf_data)

        perf_data = "io_aborted=0;;;; io_busresets=0;;;; io_read_latency=0ms;;;;, io_write_latency=0ms;;;; io_kernel_latency=0ms;;;; io_device_latency=0ms;;;; io_queue_latency=0ms;;;;"
        perf_data = check_perfdata(perf_data)

        perf_data = "g[in_bps]=15793;7000000000;8000000000;0;10000000000 g[out_bps]=116;7000000000;8000000000;0;10000000000 c[in_error]=50;1;2;0;10000"
        perf_data = check_perfdata(perf_data)

        if perf_data[0]['metric'] != 'g' or perf_data[0]['value'] != 15793.0 or perf_data[0]['unit'] != "in_bps" or perf_data[0]['warn'] != 7000000000.0:
            raise Exception('Invalid parsing: %s' % perf_data)

        perf_data = "offset=-0.000641s;-60.000000;120.000000;"
        perf_data = check_perfdata(perf_data)

        if perf_data[0]['value'] != -0.000641 or perf_data[0]['warn'] != -60.0 or perf_data[0]['crit'] != 120.0:
            raise Exception('Invalid parsing: %s' % perf_data)

    def test_02_calcul_pct(self):
        result = {'unknown': 23.01, 'warning': 41.0, 'ok': 26.55, 'critical': 9.44}
        data = {'ok': 90, 'warning': 139, 'critical': 32, 'unknown': 78}
        pct = calcul_pct(data)
        if pct != result:
            raise Exception('Error in pct calculation ...')

    def test_03_perf(self):
        perf_datas = [
            "/=541MB;906;956;0;1007 /home=62MB;452;477;0;503 /tmp=38MB;906;956;0;1007 /var=943MB;1813;1914;0;2015 /usr=2249MB;5441;5743;0;6046 /opt=68MB;112;118;0;125 /boot=25MB;90;95;0;100 /backup=4410MB;7256;7659;0;8063 /mnt/NAS_OPU_LIVRAISON=35603MB;36300;38317;0;40334 /products/admin=595MB;906;956;0;1007 /products/oracle/10.2.0=146MB;7256;7659;0;8063 /products/agtgrid=851MB;6349;6702;0;7055 /app/PPOPVGL=10966MB;27213;28725;0;30237",
            "redo_log_file_switch_interval=6s;360:;60: /data/XXXX/oracle/XXXXX/dbs/archivelog=2059MB;30686;32391;0;34096",
            "'C:\ Used'=55.4%;90;95; 'S:\ Used'=27.2%;90;95; 'E:\ Used'=24.6%;90;95; 'Q:\ Used'=1.1%;90;95; 'F:\ Used'=83.3%;90;95; 'F:\Journaux\ Used'=0.5%;90;95; 'F:\Data\ Used'=2.0%;90;95; 'F:\App\ Used'=1.4%;90;95; 'G:\ Used'=62.6%;90;95; 'G:\Journaux\ Used'=21.6%;90;95; 'G:\App\ Used'=7.8%;90;95; '\\?\Volume{c75733e9-4327-11e0-8596-0010184d9c22}\ Used'=35.9%;90;95; 'H:\ Used'=76.6%;90;95; 'H:\Journaux\ Used'=7.2%;90;95; '\\?\Volume{b7cdcbde-4fe5-11e0-a4e6-0010184d9c22}\ Used'=49.6%;90;95; 'H:\App\ Used'=9.9%;90;95; 'I:\ Used'=46.7%;90;95; 'I:\Journaux\ Used'=8.7%;90;95; '\\?\Volume{b7cdcbe1-4fe5-11e0-a4e6-0010184d9c22}\ Used'=99.0%;90;95; 'I:\App\ Used'=10.8%;90;95; 'J:\ Used'=75.2%;90;95; 'J:\Journaux\ Used'=7.0%;90;95; 'J:\App\ Used'=10.0%;90;95; 'J:\Data\ Used'=74.2%;90;95; 'K:\ Used'=77.3%;90;95; 'K:\Journaux\ Used'=8.2%;90;95; 'K:\Data\ Used'=84.6%;90;95; 'K:\App\ Used'=10.2%;90;95; 'G:\Data\ Used'=25.7%;90;95; 'H:\Data\ Used'=41.4%;90;95; 'I:\Data\ Used'=46.1%;90;95;",
            "load1=0.040;3.000;4.000;0; load5=0.070;4.000;5.000;0; load15=0.010;4.000;5.000;0;"
        ]

        print("Start test:")
        start = time.time()
        for perf_data in perf_datas:
            check_perfdata(perf_data)
        stop = time.time()

        print(" + Elapsed: %.5f" % (stop - start))


if __name__ == "__main__":
    unittest.main(verbosity=1)
