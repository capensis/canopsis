#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2016 "Capensis" [http://www.capensis.com]
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

from unittest import TestCase, main

from canopsis.stats.manager import Stats


class TestManager(TestCase):
    def setUp(self):
        self.stats = Stats()

    def test__influx_or_regex(self):
        cases = [
            {
                'items': ['a', 'b', 'c'],
                'expected': '/^(a|b|c)$/'
            },
            {
                'items': ['a)$/', 'b/d', 'c|||'],
                'expected': '/^(a\)\$\/|b\/d|c\|\|\|)$/'
            }
        ]

        for case in cases:
            res = self.stats._influx_or_regex(case['items'])

            self.assertEqual(res, case['expected'])

    def test_get_event_stats_zeros(self):
        """
        Tests on Stats.get_event_stats expecting all stats to be zeros
        """
        cases = [
            {
                'tstart': 0,
                'tstop': 0,
                'tags': {},
                'expected': {
                    '__total__': {
                        'stats_count': {
                            'alarms_new': 0,
                            'alarms_ack': 0,
                            'alarms_solved_ack': 0,
                            'alarms_solved_without_ack': 0
                        },
                        'stats_delay': {
                            'ack_delay_min': 0,
                            'ack_delay_avg': 0,
                            'ack_delay_max': 0,
                            'ack_solved_delay_min': 0,
                            'ack_solved_delay_avg': 0,
                            'ack_solved_delay_max': 0,
                            'alarm_solved_delay_min': 0,
                            'alarm_solved_delay_avg': 0,
                            'alarm_solved_delay_max': 0
                        }
                    }
                }
            },
            {
                'tstart': 0,
                'tstop': 0,
                'tags': {
                    'domain': ['__test_d1__', '__test_d2__'],
                    'perimeter': ['__test_p1__']
                },
                'expected': {
                    'domain': [
                        {
                            'name': '__test_d1__',
                            'stats_count': {
                                'alarms_new': 0,
                                'alarms_ack': 0,
                                'alarms_solved_ack': 0,
                                'alarms_solved_without_ack': 0
                            },
                            'stats_delay': {
                                'ack_delay_min': 0,
                                'ack_delay_avg': 0,
                                'ack_delay_max': 0,
                                'ack_solved_delay_min': 0,
                                'ack_solved_delay_avg': 0,
                                'ack_solved_delay_max': 0,
                                'alarm_solved_delay_min': 0,
                                'alarm_solved_delay_avg': 0,
                                'alarm_solved_delay_max': 0
                            }
                        },
                        {
                            'name': '__test_d2__',
                            'stats_count': {
                                'alarms_new': 0,
                                'alarms_ack': 0,
                                'alarms_solved_ack': 0,
                                'alarms_solved_without_ack': 0
                            },
                            'stats_delay': {
                                'ack_delay_min': 0,
                                'ack_delay_avg': 0,
                                'ack_delay_max': 0,
                                'ack_solved_delay_min': 0,
                                'ack_solved_delay_avg': 0,
                                'ack_solved_delay_max': 0,
                                'alarm_solved_delay_min': 0,
                                'alarm_solved_delay_avg': 0,
                                'alarm_solved_delay_max': 0
                            }
                        },
                        {
                            'name': '__total__',
                            'stats_count': {
                                'alarms_new': 0,
                                'alarms_ack': 0,
                                'alarms_solved_ack': 0,
                                'alarms_solved_without_ack': 0
                            },
                            'stats_delay': {
                                'ack_delay_min': 0,
                                'ack_delay_avg': 0,
                                'ack_delay_max': 0,
                                'ack_solved_delay_min': 0,
                                'ack_solved_delay_avg': 0,
                                'ack_solved_delay_max': 0,
                                'alarm_solved_delay_min': 0,
                                'alarm_solved_delay_avg': 0,
                                'alarm_solved_delay_max': 0
                            }
                        }
                    ],
                    'perimeter': [
                        {
                            'name': '__test_p1__',
                            'stats_count': {
                                'alarms_new': 0,
                                'alarms_ack': 0,
                                'alarms_solved_ack': 0,
                                'alarms_solved_without_ack': 0
                            },
                            'stats_delay': {
                                'ack_delay_min': 0,
                                'ack_delay_avg': 0,
                                'ack_delay_max': 0,
                                'ack_solved_delay_min': 0,
                                'ack_solved_delay_avg': 0,
                                'ack_solved_delay_max': 0,
                                'alarm_solved_delay_min': 0,
                                'alarm_solved_delay_avg': 0,
                                'alarm_solved_delay_max': 0
                            }
                        },
                        {
                            'name': '__total__',
                            'stats_count': {
                                'alarms_new': 0,
                                'alarms_ack': 0,
                                'alarms_solved_ack': 0,
                                'alarms_solved_without_ack': 0
                            },
                            'stats_delay': {
                                'ack_delay_min': 0,
                                'ack_delay_avg': 0,
                                'ack_delay_max': 0,
                                'ack_solved_delay_min': 0,
                                'ack_solved_delay_avg': 0,
                                'ack_solved_delay_max': 0,
                                'alarm_solved_delay_min': 0,
                                'alarm_solved_delay_avg': 0,
                                'alarm_solved_delay_max': 0
                            }
                        }
                    ],
                    '__total__': {
                        'stats_count': {
                            'alarms_new': 0,
                            'alarms_ack': 0,
                            'alarms_solved_ack': 0,
                            'alarms_solved_without_ack': 0
                        },
                        'stats_delay': {
                            'ack_delay_min': 0,
                            'ack_delay_avg': 0,
                            'ack_delay_max': 0,
                            'ack_solved_delay_min': 0,
                            'ack_solved_delay_avg': 0,
                            'ack_solved_delay_max': 0,
                            'alarm_solved_delay_min': 0,
                            'alarm_solved_delay_avg': 0,
                            'alarm_solved_delay_max': 0
                        }
                    }
                }
            },
            {
                'tstart': 0,
                'tstop': 0,
                'tags': {'domain': ['__test_d2__', '__test_d1__']},
                'expected': {
                    'domain': [
                        {
                            'name': '__test_d2__',
                            'stats_count': {
                                'alarms_new': 0,
                                'alarms_ack': 0,
                                'alarms_solved_ack': 0,
                                'alarms_solved_without_ack': 0
                            },
                            'stats_delay': {
                                'ack_delay_min': 0,
                                'ack_delay_avg': 0,
                                'ack_delay_max': 0,
                                'ack_solved_delay_min': 0,
                                'ack_solved_delay_avg': 0,
                                'ack_solved_delay_max': 0,
                                'alarm_solved_delay_min': 0,
                                'alarm_solved_delay_avg': 0,
                                'alarm_solved_delay_max': 0
                            }
                        },
                        {
                            'name': '__test_d1__',
                            'stats_count': {
                                'alarms_new': 0,
                                'alarms_ack': 0,
                                'alarms_solved_ack': 0,
                                'alarms_solved_without_ack': 0
                            },
                            'stats_delay': {
                                'ack_delay_min': 0,
                                'ack_delay_avg': 0,
                                'ack_delay_max': 0,
                                'ack_solved_delay_min': 0,
                                'ack_solved_delay_avg': 0,
                                'ack_solved_delay_max': 0,
                                'alarm_solved_delay_min': 0,
                                'alarm_solved_delay_avg': 0,
                                'alarm_solved_delay_max': 0
                            }
                        },
                        {
                            'name': '__total__',
                            'stats_count': {
                                'alarms_new': 0,
                                'alarms_ack': 0,
                                'alarms_solved_ack': 0,
                                'alarms_solved_without_ack': 0
                            },
                            'stats_delay': {
                                'ack_delay_min': 0,
                                'ack_delay_avg': 0,
                                'ack_delay_max': 0,
                                'ack_solved_delay_min': 0,
                                'ack_solved_delay_avg': 0,
                                'ack_solved_delay_max': 0,
                                'alarm_solved_delay_min': 0,
                                'alarm_solved_delay_avg': 0,
                                'alarm_solved_delay_max': 0
                            }
                        }
                    ],
                    '__total__': {
                        'stats_count': {
                            'alarms_new': 0,
                            'alarms_ack': 0,
                            'alarms_solved_ack': 0,
                            'alarms_solved_without_ack': 0
                        },
                        'stats_delay': {
                            'ack_delay_min': 0,
                            'ack_delay_avg': 0,
                            'ack_delay_max': 0,
                            'ack_solved_delay_min': 0,
                            'ack_solved_delay_avg': 0,
                            'ack_solved_delay_max': 0,
                            'alarm_solved_delay_min': 0,
                            'alarm_solved_delay_avg': 0,
                            'alarm_solved_delay_max': 0
                        }
                    }
                }
            }
        ]

        for case in cases:
            res = self.stats.get_event_stats(
                case['tstart'],
                case['tstop'],
                tags=case['tags']
            )

            self.assertEqual(res, case['expected'])

    def test_get_user_stats_zeros(self):
        """
        Tests on Stats.get_user_stats expecting all stats to be zeros
        """
        cases = [
            {
                'tstart': 0,
                'tstop': 0,
                'users': [],
                'tags': {},
                'expected': []
            },
            {
                'tstart': 0,
                'tstop': 0,
                'users': ['__test_u1__', '__test_u2__'],
                'tags': {},
                'expected': [
                    {
                        'author': '__test_u1__',
                        'ack': {
                            'total': 0,
                            'delay_min': 0,
                            'delay_avg': 0,
                            'delay_max': 0
                        },
                        'session': {
                            'duration_min': 0,
                            'duration_avg': 0,
                            'duration_max': 0
                        },
                        'tags': {}
                    },
                    {
                        'author': '__test_u2__',
                        'ack': {
                            'total': 0,
                            'delay_min': 0,
                            'delay_avg': 0,
                            'delay_max': 0
                        },
                        'session': {
                            'duration_min': 0,
                            'duration_avg': 0,
                            'duration_max': 0
                        },
                        'tags': {}
                    }
                ]
            },
            {
                'tstart': 0,
                'tstop': 0,
                'users': ['__test_u2__', '__test_u1__'],
                'tags': {},
                'expected': [
                    {
                        'author': '__test_u2__',
                        'ack': {
                            'total': 0,
                            'delay_min': 0,
                            'delay_avg': 0,
                            'delay_max': 0
                        },
                        'session': {
                            'duration_min': 0,
                            'duration_avg': 0,
                            'duration_max': 0
                        },
                        'tags': {}
                    },
                    {
                        'author': '__test_u1__',
                        'ack': {
                            'total': 0,
                            'delay_min': 0,
                            'delay_avg': 0,
                            'delay_max': 0
                        },
                        'session': {
                            'duration_min': 0,
                            'duration_avg': 0,
                            'duration_max': 0
                        },
                        'tags': {}
                    }
                ]
            },
            {
                'tstart': 0,
                'tstop': 0,
                'users': ['__test_u1__', '__test_u2__'],
                'tags': {
                    'domain': ['__test_d1__', '__test_d2__'],
                    'perimeter': ['__test_p1__']
                },
                'expected': [
                    {
                        'author': '__test_u1__',
                        'ack': {
                            'total': 0,
                            'delay_min': 0,
                            'delay_avg': 0,
                            'delay_max': 0
                        },
                        'session': {
                            'duration_min': 0,
                            'duration_avg': 0,
                            'duration_max': 0
                        },
                        'tags': {
                            'domain': [
                                {
                                    'name': '__test_d1__',
                                    'ack_total': 0
                                },
                                {
                                    'name': '__test_d2__',
                                    'ack_total': 0
                                }
                            ],
                            'perimeter': [
                                {
                                    'name': '__test_p1__',
                                    'ack_total': 0
                                }
                            ]
                        }
                    },
                    {
                        'author': '__test_u2__',
                        'ack': {
                            'total': 0,
                            'delay_min': 0,
                            'delay_avg': 0,
                            'delay_max': 0
                        },
                        'session': {
                            'duration_min': 0,
                            'duration_avg': 0,
                            'duration_max': 0
                        },
                        'tags': {
                            'domain': [
                                {
                                    'name': '__test_d1__',
                                    'ack_total': 0
                                },
                                {
                                    'name': '__test_d2__',
                                    'ack_total': 0
                                }
                            ],
                            'perimeter': [
                                {
                                    'name': '__test_p1__',
                                    'ack_total': 0
                                }
                            ]
                        }
                    }
                ]
            }
        ]

        for case in cases:
            res = self.stats.get_user_stats(
                case['tstart'],
                case['tstop'],
                users=case['users'],
                tags=case['tags']
            )

            self.assertEqual(res, case['expected'])

if __name__ == '__main__':
    main()
