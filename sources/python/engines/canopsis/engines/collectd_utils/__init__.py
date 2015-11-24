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

types = {
    'operations': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'ipt_bytes': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'percent': [
        {
            'max': '100.1',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None,
            'unit': None
        }
    ],
    'vmpage_number': [
        {
            'max': '4294967295',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'current_connection': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'virt_vcpu': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'entropy': [
        {
            'max': '4294967295',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'voltage': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': 'U'
        }
    ],
    'current': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': 'U'
        }
    ],
    'spam_score': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': 'U'
        }
    ],
    'derive': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'disk_ops_comple': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'if_rx_errors': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'current_session': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'tcp_connections': [
        {
            'max': '4294967295',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'disk_octets': [
        {
            'max': None,
            'name': 'read',
            'type': 'DERIVE',
            'min': '0',
            'unit': 'o'
        },
        {
            'max': None,
            'name': 'write',
            'type': 'DERIVE',
            'min': '0',
            'unit': 'o'
        }
    ],
    'charge': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'vs_memory': [
        {
            'max': '9223372036854775807',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'swap': [
        {
            'max': '1099511627776',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'io_octets': [
        {
            'max': None,
            'name': 'rx',
            'type': 'DERIVE',
            'min': '0',
            'unit': 'o'
        },
        {
            'max': None,
            'name': 'tx',
            'type': 'DERIVE',
            'min': '0',
            'unit': 'o'
        }
    ],
    'ping_droprate': [
        {
            'max': '100',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'nfs_procedure': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'dns_answer': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'frequency_offse': [
        {
            'max': '1000000',
            'name': 'value',
            'type': 'GAUGE',
            'min': '-1000000'
        }
    ],
    'node_octets': [
        {
            'max': None,
            'name': 'rx',
            'type': 'DERIVE',
            'min': '0',
            'unit': 'o'
        },
        {
            'max': None,
            'name': 'tx',
            'type': 'DERIVE',
            'min': '0',
            'unit': 'o'
        }
    ],
    'df': [
        {
            'max': '1125899906842623',
            'name': 'used',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        },
        {
            'max': '1125899906842623',
            'name': 'free',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'spam_check': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'mysql_qcache': [
        {
            'max': None,
            'name': 'hits',
            'type': 'COUNTER',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'inserts',
            'type': 'COUNTER',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'not_cached',
            'type': 'COUNTER',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'lowmem_prunes',
            'type': 'COUNTER',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'queries_in_cache',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'if_errors': [
        {
            'max': None,
            'name': 'rx',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'tx',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'http_request_method': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'ps_stacksize': [
        {
            'max': '9223372036854775807',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'nginx_connection': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'disk_ops': [
        {
            'max': None,
            'name': 'read',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'write',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'dns_request': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'memcached_items': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'virt_cpu_total': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'io_packets': [
        {
            'max': None,
            'name': 'rx',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'tx',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'humidity': [
        {
            'max': '100',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'memcached_connection': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'volatile_change': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'cache_size': [
        {
            'max': '4294967295',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'gauge': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': 'U',
            'unit': None
        }
    ],
    'ps_rss': [
        {
            'max': '9223372036854775807',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'ps_cputime': [
        {
            'max': None,
            'name': 'user',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'syst',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'dns_query': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'mysql_commands': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'signal_quality': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'bytes': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'connections': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'frequency': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'total_bytes': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'dns_notify': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'vs_processes': [
        {
            'max': '65535',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'compression_ratio': [
        {
            'max': '2',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'vcpu': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'if_octets': [
        {
            'max': None,
            'name': 'rx',
            'type': 'DERIVE',
            'min': '0',
            'unit': 'o'
        },
        {
            'max': None,
            'name': 'tx',
            'type': 'DERIVE',
            'min': '0',
            'unit': 'o'
        }
    ],
    'fork_rate': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'mysql_octets': [
        {
            'max': None,
            'name': 'rx',
            'type': 'DERIVE',
            'min': '0',
            'unit': 'o'
        },
        {
            'max': None,
            'name': 'tx',
            'type': 'DERIVE',
            'min': '0',
            'unit': 'o'
        }
    ],
    'total_sessions': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'power': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'disk_latency': [
        {
            'max': None,
            'name': 'read',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'write',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'memcached_ops': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'cache_result': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'ps_vm': [
        {
            'max': '9223372036854775807',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'if_packets': [
        {
            'max': None,
            'name': 'rx',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'tx',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'if_tx_errors': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'bitrate': [
        {
            'max': '4294967295',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'ps_disk_octets': [
        {
            'max': None,
            'name': 'read',
            'type': 'DERIVE',
            'min': '0',
            'unit': 'o'
        },
        {
            'max': None,
            'name': 'write',
            'type': 'DERIVE',
            'min': '0',
            'unit': 'o'
        }
    ],
    'compression': [
        {
            'max': None,
            'name': 'uncompressed',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'compressed',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'nginx_requests': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'timeleft': [
        {
            'max': '3600',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'multimeter': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': 'U'
        }
    ],
    'dns_qtype_cache': [
        {
            'max': '4294967295',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'total_time_in_m': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'players': [
        {
            'max': '1000000',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'cpufreq': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'ps_disk_ops': [
        {
            'max': None,
            'name': 'read',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'write',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'routes': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'threads': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'requests': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'arc_ratio': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'load': [
        {
            'max': '100',
            'name': 'shortterm',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        },
        {
            'max': '100',
            'name': 'midterm',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        },
        {
            'max': '100',
            'name': 'longterm',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'apache_connection': [
        {
            'max': '65535',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'total_requests': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'apache_scoreboar': [
        {
            'max': '65535',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'dns_qtype': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'route_etx': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'serial_octets': [
        {
            'max': None,
            'name': 'rx',
            'type': 'DERIVE',
            'min': '0',
            'unit': 'o'
        },
        {
            'max': None,
            'name': 'tx',
            'type': 'DERIVE',
            'min': '0',
            'unit': 'o'
        }
    ],
    'if_dropped': [
        {
            'max': None,
            'name': 'rx',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'tx',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'node_stat': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'dns_zops': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'arc_size': [
        {
            'max': None,
            'name': 'current',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'target',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'minlimit',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'maxlimit',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'temperature': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '-273.15'
        }
    ],
    'total_threads': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'dns_question': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'vmpage_faults': [
        {
            'max': None,
            'name': 'minflt',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'majflt',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'dns_rcode': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'arc_l2_bytes': [
        {
            'max': None,
            'name': 'read',
            'type': 'COUNTER',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'write',
            'type': 'COUNTER',
            'min': '0',
            'unit': None
        }
    ],
    'response_time': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'memory': [
        {
            'max': '281474976710656',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'vmpage_io': [
        {
            'max': None,
            'name': 'in',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'out',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'dns_reject': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'apache_requests': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'files': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'df_complex': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'ps_code': [
        {
            'max': '9223372036854775807',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'email_count': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'route_metric': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'ps_count': [
        {
            'max': '1000000',
            'name': 'processes',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        },
        {
            'max': '1000000',
            'name': 'threads',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'ath_nodes': [
        {
            'max': '65535',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'df_inodes': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'fanspeed': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'apache_bytes': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'dns_response': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'mysql_threads': [
        {
            'max': None,
            'name': 'running',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'connected',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'cached',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'created',
            'type': 'COUNTER',
            'min': '0',
            'unit': None
        }
    ],
    'pg_scan': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'node_tx_rate': [
        {
            'max': '127',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'contextswitch': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'cache_operation': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'total_connection': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'invocations': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'queue_length': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'mysql_handler': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'arc_l2_size': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'vmpage_action': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'ping_stddev': [
        {
            'max': '65535',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'time_offset': [
        {
            'max': '1000000',
            'name': 'value',
            'type': 'GAUGE',
            'min': '-1000000'
        }
    ],
    'ipt_packets': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'arc_counts': [
        {
            'max': None,
            'name': 'demand_data',
            'type': 'COUNTER',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'demand_metadata',
            'type': 'COUNTER',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'prefetch_data',
            'type': 'COUNTER',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'prefetch_metadata',
            'type': 'COUNTER',
            'min': '0',
            'unit': None
        }
    ],
    'links': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'mysql_log_positio': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'time_dispersion': [
        {
            'max': '1000000',
            'name': 'value',
            'type': 'GAUGE',
            'min': '-1000000'
        }
    ],
    'memcached_comman': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'pg_xact': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'voltage_threshold': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': 'U'
        },
        {
            'max': None,
            'name': 'threshold',
            'type': 'GAUGE',
            'min': 'U'
        }
    ],
    'http_requests': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'conntrack': [
        {
            'max': '4294967295',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'file_size': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'apache_idle_worker': [
        {
            'max': '65535',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'total_values': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'disk_time': [
        {
            'max': None,
            'name': 'read',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'write',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'absolute': [
        {
            'max': None,
            'name': 'value',
            'type': 'ABSOLUTE',
            'min': '0',
            'unit': None
        }
    ],
    'latency': [
        {
            'max': '65535',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'uptime': [
        {
            'max': '4294967295',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'total_operation': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'email_check': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'ping': [
        {
            'max': '65535',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'dns_transfer': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'if_multicast': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'delay': [
        {
            'max': '1000000',
            'name': 'value',
            'type': 'GAUGE',
            'min': '-1000000'
        }
    ],
    'protocol_counte': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'pg_numbackends': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'vs_threads': [
        {
            'max': '65535',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'ath_stat': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'pg_db_size': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'dns_update': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'users': [
        {
            'max': '65535',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'pg_n_tup_g': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'mysql_locks': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'fscache_stat': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'swap_io': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'dns_resolver': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'signal_noise': [
        {
            'max': '0',
            'name': 'value',
            'type': 'GAUGE',
            'min': 'U'
        }
    ],
    'snr': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'ps_state': [
        {
            'max': '65535',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'memcached_octets': [
        {
            'max': None,
            'name': 'rx',
            'type': 'DERIVE',
            'min': '0',
            'unit': 'o'
        },
        {
            'max': None,
            'name': 'tx',
            'type': 'DERIVE',
            'min': '0',
            'unit': 'o'
        }
    ],
    'disk_merged': [
        {
            'max': None,
            'name': 'read',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'write',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'pg_n_tup_c': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'dns_octets': [
        {
            'max': None,
            'name': 'queries',
            'type': 'DERIVE',
            'min': '0',
            'unit': 'o'
        },
        {
            'max': None,
            'name': 'responses',
            'type': 'DERIVE',
            'min': '0',
            'unit': 'o'
        }
    ],
    'email_size': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'ps_data': [
        {
            'max': '9223372036854775807',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'http_response_code': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'ps_pagefaults': [
        {
            'max': None,
            'name': 'minflt',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        },
        {
            'max': None,
            'name': 'majflt',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'records': [
        {
            'max': None,
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'signal_power': [
        {
            'max': '0',
            'name': 'value',
            'type': 'GAUGE',
            'min': 'U'
        }
    ],
    'counter': [
        {
            'max': None,
            'name': 'value',
            'type': 'COUNTER',
            'min': 'U'
        }
    ],
    'pg_blks': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'node_rssi': [
        {
            'max': '255',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'cache_ratio': [
        {
            'max': '100',
            'name': 'value',
            'type': 'GAUGE',
            'min': '0',
            'unit': None
        }
    ],
    'irq': [
        {
            'max': None,
            'name': 'value',
            'type': 'DERIVE',
            'min': '0',
            'unit': None
        }
    ],
    'dns_opcode': [
        {
                'max': None,
                'name': 'value',
                'type': 'DERIVE',
                'min': '0',
                'unit': None
        }
    ],
    'cpu': [
        {
                'max': None,
                'name': 'value',
                'type': 'DERIVE',
                'min': '0',
                'unit': None
        }
    ],
    'if_collisions': [
        {
                'max': None,
                'name': 'value',
                'type': 'DERIVE',
                'min': '0',
                'unit': None
        }
    ],
    'current': [
        {
                'max': None,
                'name': 'value',
                'type': 'GAUGE',
                'min': '0',
                'unit': "A"
        }
    ],
    'voltage': [
        {
                'max': None,
                'name': 'value',
                'type': 'GAUGE',
                'min': '0',
                'unit': "V"
        }
    ],
    'temperature': [
        {
                'max': None,
                'name': 'value',
                'type': 'GAUGE',
                'min': '0',
                'unit': "C"
        }
    ]
}
