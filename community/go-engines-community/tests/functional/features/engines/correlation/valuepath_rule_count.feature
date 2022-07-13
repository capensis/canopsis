Feature: correlation feature - valuegroup rule with threshold count
  Scenario: given meta alarm rule and events should create 2 separate metaalarms
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "description" : "test-correlation-valuegroup-1",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-valuegroup-1"
            }
          }
        ]
      ],
      "enabled" : true,
      "external_data" : {},
      "config": {
        "actions": [
          {
            "type" : "set_entity_info_from_template",
            "name" : "infoenrich",
            "value" : "{{ `{{.Event.ExtraInfos.infoenrich}}` }}",
            "description" : "infoenrich"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "priority" : 10001,
      "type" : "enrichment"
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-valuegroup-correlation-1",
      "type": "valuegroup",
      "config": {
        "time_interval": {
          "value": 10,
          "unit": "s"
        },
        "threshold_count": 2,
        "value_paths": [
          "entity.infos.infoenrich.value"
        ]
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-valuegroup-1",
      "connector_name": "test-valuegroup-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-1",
      "resource": "test-valuegroup-correlation-1-resource-1",
      "infoenrich": "1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-1",
      "connector_name": "test-valuegroup-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-1",
      "resource": "test-valuegroup-correlation-1-resource-2",
      "infoenrich": "2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I do GET /api/v4/alarms?search={{ .metaAlarmRuleID }}&active_columns[]=v.meta&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    When I send an event:
    """
    {
      "connector": "test-valuegroup-1",
      "connector_name": "test-valuegroup-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-1",
      "resource": "test-valuegroup-correlation-1-resource-3",
      "infoenrich": "1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-1",
      "connector_name": "test-valuegroup-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-1",
      "resource": "test-valuegroup-correlation-1-resource-4",
      "infoenrich": "1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I wait 1s
    When I send an event:
    """
    {
      "connector": "test-valuegroup-1",
      "connector_name": "test-valuegroup-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-1",
      "resource": "test-valuegroup-correlation-1-resource-5",
      "infoenrich": "2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search={{ .metaAlarmRuleID }}&active_columns[]=v.meta&correlation=true&sort_by=t&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-valuegroup-correlation-1"
          }
        },
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-valuegroup-correlation-1"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      },
      {
        "_id": "{{ (index .lastResponse.data 1)._id }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-valuegroup-1",
                  "connector_name": "test-valuegroup-1-name",
                  "component":  "test-valuegroup-correlation-1",
                  "resource": "test-valuegroup-correlation-1-resource-1"
                }
              },
              {
                "v": {
                  "connector": "test-valuegroup-1",
                  "connector_name": "test-valuegroup-1-name",
                  "component":  "test-valuegroup-correlation-1",
                  "resource": "test-valuegroup-correlation-1-resource-3"
                }
              },
              {
                "v": {
                  "connector": "test-valuegroup-1",
                  "connector_name": "test-valuegroup-1-name",
                  "component":  "test-valuegroup-correlation-1",
                  "resource": "test-valuegroup-correlation-1-resource-4"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 3
            }
          }
        }
      },
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-valuegroup-1",
                  "connector_name": "test-valuegroup-1-name",
                  "component":  "test-valuegroup-correlation-1",
                  "resource": "test-valuegroup-correlation-1-resource-2"
                }
              },
              {
                "v": {
                  "connector": "test-valuegroup-1",
                  "connector_name": "test-valuegroup-1-name",
                  "component":  "test-valuegroup-correlation-1",
                  "resource": "test-valuegroup-correlation-1-resource-5"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 2
            }
          }
        }
      }
    ]
    """

  Scenario: given meta alarm rule and events should create 4 separate metaalarms
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "description" : "test-correlation-valuegroup-2",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-valuegroup-2"
            }
          }
        ]
      ],
      "enabled" : true,
      "external_data" : {},
      "config": {
        "actions": [
          {
            "type" : "set_entity_info_from_template",
            "name" : "infoenrich2",
            "value" : "{{ `{{.Event.ExtraInfos.infoenrich2}}` }}",
            "description" : "infoenrich2"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "priority" : 10001,
      "type" : "enrichment"
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "description" : "test-correlation-valuegroup-2",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-valuegroup-2"
            }
          }
        ]
      ],
      "enabled" : true,
      "external_data" : {},
      "config": {
        "actions": [
          {
            "type" : "set_entity_info_from_template",
            "name" : "infoenrich3",
            "value" : "{{ `{{.Event.ExtraInfos.infoenrich3}}` }}",
            "description" : "infoenrich3"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "priority" : 10001,
      "type" : "enrichment"
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-valuegroup-correlation-1",
      "type": "valuegroup",
      "config": {
        "time_interval": {
          "value": 10,
          "unit": "s"
        },
        "threshold_count": 2,
        "value_paths": [
          "entity.infos.infoenrich2.value",
          "entity.infos.infoenrich3.value"
        ]
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-valuegroup-2",
      "connector_name": "test-valuegroup-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-2",
      "resource": "test-valuegroup-correlation-2-resource-1",
      "infoenrich2": "1",
      "infoenrich3": "1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-2",
      "connector_name": "test-valuegroup-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-2",
      "resource": "test-valuegroup-correlation-2-resource-2",
      "infoenrich2": "1",
      "infoenrich3": "2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-2",
      "connector_name": "test-valuegroup-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-2",
      "resource": "test-valuegroup-correlation-2-resource-3",
      "infoenrich2": "2",
      "infoenrich3": "1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-2",
      "connector_name": "test-valuegroup-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-2",
      "resource": "test-valuegroup-correlation-2-resource-4",
      "infoenrich2": "2",
      "infoenrich3": "2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I do GET /api/v4/alarms?search={{ .metaAlarmRuleID }}&active_columns[]=v.meta&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    When I send an event:
    """
    {
      "connector": "test-valuegroup-2",
      "connector_name": "test-valuegroup-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-2",
      "resource": "test-valuegroup-correlation-2-resource-5",
      "infoenrich2": "1",
      "infoenrich3": "1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait 1s
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-2",
      "connector_name": "test-valuegroup-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-2",
      "resource": "test-valuegroup-correlation-2-resource-6",
      "infoenrich2": "1",
      "infoenrich3": "2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait 1s
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-2",
      "connector_name": "test-valuegroup-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-2",
      "resource": "test-valuegroup-correlation-2-resource-7",
      "infoenrich2": "2",
      "infoenrich3": "1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait 1s
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-2",
      "connector_name": "test-valuegroup-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-2",
      "resource": "test-valuegroup-correlation-2-resource-8",
      "infoenrich2": "2",
      "infoenrich3": "2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search={{ .metaAlarmRuleID }}&active_columns[]=v.meta&correlation=true&sort_by=t&sort=asc
    Then the response code should be 200
    """
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-valuegroup-correlation-2"
          }
        },
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-valuegroup-correlation-2"
          }
        },
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-valuegroup-correlation-2"
          }
        },
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-valuegroup-correlation-2"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 4
      }
    }
    """
    When I save response metaAlarmID1={{ (index .lastResponse.data 0)._id }}
    When I save response metaAlarmID2={{ (index .lastResponse.data 1)._id }}
    When I save response metaAlarmID3={{ (index .lastResponse.data 2)._id }}
    When I save response metaAlarmID4={{ (index .lastResponse.data 3)._id }}
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID1 }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      },
      {
        "_id": "{{ .metaAlarmID2 }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      },
      {
        "_id": "{{ .metaAlarmID3 }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      },
      {
        "_id": "{{ .metaAlarmID4 }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-valuegroup-2",
                  "connector_name": "test-valuegroup-2-name",
                  "component":  "test-valuegroup-correlation-2",
                  "resource": "test-valuegroup-correlation-2-resource-1"
                }
              },
              {
                "v": {
                  "connector": "test-valuegroup-2",
                  "connector_name": "test-valuegroup-2-name",
                  "component":  "test-valuegroup-correlation-2",
                  "resource": "test-valuegroup-correlation-2-resource-5"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 2
            }
          }
        }
      },
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-valuegroup-2",
                  "connector_name": "test-valuegroup-2-name",
                  "component":  "test-valuegroup-correlation-2",
                  "resource": "test-valuegroup-correlation-2-resource-2"
                }
              },
              {
                "v": {
                  "connector": "test-valuegroup-2",
                  "connector_name": "test-valuegroup-2-name",
                  "component":  "test-valuegroup-correlation-2",
                  "resource": "test-valuegroup-correlation-2-resource-6"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 2
            }
          }
        }
      },
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-valuegroup-2",
                  "connector_name": "test-valuegroup-2-name",
                  "component":  "test-valuegroup-correlation-2",
                  "resource": "test-valuegroup-correlation-2-resource-3"
                }
              },
              {
                "v": {
                  "connector": "test-valuegroup-2",
                  "connector_name": "test-valuegroup-2-name",
                  "component":  "test-valuegroup-correlation-2",
                  "resource": "test-valuegroup-correlation-2-resource-7"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 2
            }
          }
        }
      },
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-valuegroup-2",
                  "connector_name": "test-valuegroup-2-name",
                  "component":  "test-valuegroup-correlation-2",
                  "resource": "test-valuegroup-correlation-2-resource-4"
                }
              },
              {
                "v": {
                  "connector": "test-valuegroup-2",
                  "connector_name": "test-valuegroup-2-name",
                  "component":  "test-valuegroup-correlation-2",
                  "resource": "test-valuegroup-correlation-2-resource-8"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 2
            }
          }
        }
      }
    ]
    """
    When I send an event:
    """
    {
      "connector": "test-valuegroup-2",
      "connector_name": "test-valuegroup-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-2",
      "resource": "test-valuegroup-correlation-2-resource-9",
      "infoenrich2": "1",
      "infoenrich3": "1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait 1s
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-2",
      "connector_name": "test-valuegroup-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-2",
      "resource": "test-valuegroup-correlation-2-resource-10",
      "infoenrich2": "1",
      "infoenrich3": "2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait 1s
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-2",
      "connector_name": "test-valuegroup-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-2",
      "resource": "test-valuegroup-correlation-2-resource-11",
      "infoenrich2": "2",
      "infoenrich3": "1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait 1s
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-2",
      "connector_name": "test-valuegroup-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-2",
      "resource": "test-valuegroup-correlation-2-resource-12",
      "infoenrich2": "2",
      "infoenrich3": "2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait 1s
    When I wait the end of 2 events processing
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID1 }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      },
      {
        "_id": "{{ .metaAlarmID2 }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      },
      {
        "_id": "{{ .metaAlarmID3 }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      },
      {
        "_id": "{{ .metaAlarmID4 }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-valuegroup-2",
                  "connector_name": "test-valuegroup-2-name",
                  "component":  "test-valuegroup-correlation-2",
                  "resource": "test-valuegroup-correlation-2-resource-1"
                }
              },
              {
                "v": {
                  "connector": "test-valuegroup-2",
                  "connector_name": "test-valuegroup-2-name",
                  "component":  "test-valuegroup-correlation-2",
                  "resource": "test-valuegroup-correlation-2-resource-5"
                }
              },
              {
                "v": {
                  "connector": "test-valuegroup-2",
                  "connector_name": "test-valuegroup-2-name",
                  "component":  "test-valuegroup-correlation-2",
                  "resource": "test-valuegroup-correlation-2-resource-9"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 3
            }
          }
        }
      },
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-valuegroup-2",
                  "connector_name": "test-valuegroup-2-name",
                  "component":  "test-valuegroup-correlation-2",
                  "resource": "test-valuegroup-correlation-2-resource-10"
                }
              },
              {
                "v": {
                  "connector": "test-valuegroup-2",
                  "connector_name": "test-valuegroup-2-name",
                  "component":  "test-valuegroup-correlation-2",
                  "resource": "test-valuegroup-correlation-2-resource-2"
                }
              },
              {
                "v": {
                  "connector": "test-valuegroup-2",
                  "connector_name": "test-valuegroup-2-name",
                  "component":  "test-valuegroup-correlation-2",
                  "resource": "test-valuegroup-correlation-2-resource-6"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 3
            }
          }
        }
      },
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-valuegroup-2",
                  "connector_name": "test-valuegroup-2-name",
                  "component":  "test-valuegroup-correlation-2",
                  "resource": "test-valuegroup-correlation-2-resource-11"
                }
              },
              {
                "v": {
                  "connector": "test-valuegroup-2",
                  "connector_name": "test-valuegroup-2-name",
                  "component":  "test-valuegroup-correlation-2",
                  "resource": "test-valuegroup-correlation-2-resource-3"
                }
              },
              {
                "v": {
                  "connector": "test-valuegroup-2",
                  "connector_name": "test-valuegroup-2-name",
                  "component":  "test-valuegroup-correlation-2",
                  "resource": "test-valuegroup-correlation-2-resource-7"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 3
            }
          }
        }
      },
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-valuegroup-2",
                  "connector_name": "test-valuegroup-2-name",
                  "component":  "test-valuegroup-correlation-2",
                  "resource": "test-valuegroup-correlation-2-resource-12"
                }
              },
              {
                "v": {
                  "connector": "test-valuegroup-2",
                  "connector_name": "test-valuegroup-2-name",
                  "component":  "test-valuegroup-correlation-2",
                  "resource": "test-valuegroup-correlation-2-resource-4"
                }
              },
              {
                "v": {
                  "connector": "test-valuegroup-2",
                  "connector_name": "test-valuegroup-2-name",
                  "component":  "test-valuegroup-correlation-2",
                  "resource": "test-valuegroup-correlation-2-resource-8"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 3
            }
          }
        }
      }
    ]
    """

  Scenario: given meta alarm rule with threshold count and events should create 2 meta alarms because of 2 separate time intervals
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "description" : "test-correlation-valuegroup-3",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-valuegroup-3"
            }
          }
        ]
      ],
      "enabled" : true,
      "external_data" : {},
      "config": {
        "actions": [
          {
            "type" : "set_entity_info_from_template",
            "name" : "infoenrich4",
            "value" : "{{ `{{.Event.ExtraInfos.infoenrich4}}` }}",
            "description" : "infoenrich4"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "priority" : 10001,
      "type" : "enrichment"
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-valuegroup-correlation-3",
      "type": "valuegroup",
      "config": {
        "time_interval": {
          "value": 3,
          "unit": "s"
        },
        "threshold_count": 2,
        "value_paths": [
          "entity.infos.infoenrich4.value"
        ]
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-valuegroup-3",
      "connector_name": "test-valuegroup-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-3",
      "resource": "test-valuegroup-correlation-3-resource-1",
      "infoenrich4": "1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-3",
      "connector_name": "test-valuegroup-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-3",
      "resource": "test-valuegroup-correlation-3-resource-2",
      "infoenrich4": "1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I wait 4s
    When I send an event:
    """
    {
      "connector": "test-valuegroup-3",
      "connector_name": "test-valuegroup-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-3",
      "resource": "test-valuegroup-correlation-3-resource-3",
      "infoenrich4": "1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-3",
      "connector_name": "test-valuegroup-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-3",
      "resource": "test-valuegroup-correlation-3-resource-4",
      "infoenrich4": "1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search={{ .metaAlarmRuleID }}&active_columns[]=v.meta&correlation=true&sort_by=t&sort=asc
    Then the response code should be 200
    """
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-valuegroup-correlation-3"
          }
        },
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-valuegroup-correlation-3"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      },
      {
        "_id": "{{ (index .lastResponse.data 1)._id }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-valuegroup-3",
                  "connector_name": "test-valuegroup-3-name",
                  "component":  "test-valuegroup-correlation-3",
                  "resource": "test-valuegroup-correlation-3-resource-1"
                }
              },
              {
                "v": {
                  "connector": "test-valuegroup-3",
                  "connector_name": "test-valuegroup-3-name",
                  "component":  "test-valuegroup-correlation-3",
                  "resource": "test-valuegroup-correlation-3-resource-2"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 2
            }
          }
        }
      },
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-valuegroup-3",
                  "connector_name": "test-valuegroup-3-name",
                  "component":  "test-valuegroup-correlation-3",
                  "resource": "test-valuegroup-correlation-3-resource-3"
                }
              },
              {
                "v": {
                  "connector": "test-valuegroup-3",
                  "connector_name": "test-valuegroup-3-name",
                  "component":  "test-valuegroup-correlation-3",
                  "resource": "test-valuegroup-correlation-3-resource-4"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 2
            }
          }
        }
      }
    ]
    """

  Scenario: given meta alarm rule with threshold count and events should create one single meta alarms because first group didn't reached threshold
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "description" : "test-correlation-valuegroup-4",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-valuegroup-4"
            }
          }
        ]
      ],
      "enabled" : true,
      "external_data" : {},
      "config": {
        "actions": [
          {
            "type" : "set_entity_info_from_template",
            "name" : "infoenrich5",
            "value" : "{{ `{{.Event.ExtraInfos.infoenrich5}}` }}",
            "description" : "infoenrich5"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "priority" : 10001,
      "type" : "enrichment"
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-valuegroup-correlation-4",
      "type": "valuegroup",
      "config": {
        "time_interval": {
          "value": 3,
          "unit": "s"
        },
        "threshold_count": 2,
        "value_paths": [
          "entity.infos.infoenrich5.value"
        ]
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-valuegroup-4",
      "connector_name": "test-valuegroup-4-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-4",
      "resource": "test-valuegroup-correlation-4-resource-1",
      "infoenrich5": "1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I wait 4s
    When I send an event:
    """
    {
      "connector": "test-valuegroup-4",
      "connector_name": "test-valuegroup-4-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-4",
      "resource": "test-valuegroup-correlation-4-resource-2",
      "infoenrich5": "1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-4",
      "connector_name": "test-valuegroup-4-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-4",
      "resource": "test-valuegroup-correlation-4-resource-3",
      "infoenrich5": "1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search={{ .metaAlarmRuleID }}&active_columns[]=v.meta&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-valuegroup-correlation-4"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-valuegroup-4",
                  "connector_name": "test-valuegroup-4-name",
                  "component":  "test-valuegroup-correlation-4",
                  "resource": "test-valuegroup-correlation-4-resource-2"
                }
              },
              {
                "v": {
                  "connector": "test-valuegroup-4",
                  "connector_name": "test-valuegroup-4-name",
                  "component":  "test-valuegroup-correlation-4",
                  "resource": "test-valuegroup-correlation-4-resource-3"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 2
            }
          }
        }
      }
    ]
    """

  Scenario: given meta alarm rule with threshold count and events should create one single meta alarm without first alarm, because interval shifting
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "description" : "test-correlation-valuegroup-5",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-valuegroup-5"
            }
          }
        ]
      ],
      "enabled" : true,
      "external_data" : {},
      "config": {
        "actions": [
          {
            "type" : "set_entity_info_from_template",
            "name" : "infoenrich6",
            "value" : "{{ `{{.Event.ExtraInfos.infoenrich6}}` }}",
            "description" : "infoenrich6"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "priority" : 10001,
      "type" : "enrichment"
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-valuegroup-correlation-5",
      "type": "valuegroup",
      "config": {
        "time_interval": {
          "value": 5,
          "unit": "s"
        },
        "threshold_count": 3,
        "value_paths": [
          "entity.infos.infoenrich6.value"
        ]
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-valuegroup-5",
      "connector_name": "test-valuegroup-5-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-5",
      "resource": "test-valuegroup-correlation-5-resource-1",
      "infoenrich6": "1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I wait 3s
    When I send an event:
    """
    {
      "connector": "test-valuegroup-5",
      "connector_name": "test-valuegroup-5-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-5",
      "resource": "test-valuegroup-correlation-5-resource-2",
      "infoenrich6": "1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I wait 3s
    When I send an event:
    """
    {
      "connector": "test-valuegroup-5",
      "connector_name": "test-valuegroup-5-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-5",
      "resource": "test-valuegroup-correlation-5-resource-3",
      "infoenrich6": "1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-5",
      "connector_name": "test-valuegroup-5-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-5",
      "resource": "test-valuegroup-correlation-5-resource-4",
      "infoenrich6": "1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search={{ .metaAlarmRuleID }}&active_columns[]=v.meta&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-valuegroup-correlation-5"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-valuegroup-5",
                  "connector_name": "test-valuegroup-5-name",
                  "component":  "test-valuegroup-correlation-5",
                  "resource": "test-valuegroup-correlation-5-resource-2"
                }
              },
              {
                "v": {
                  "connector": "test-valuegroup-5",
                  "connector_name": "test-valuegroup-5-name",
                  "component":  "test-valuegroup-correlation-5",
                  "resource": "test-valuegroup-correlation-5-resource-3"
                }
              },
              {
                "v": {
                  "connector": "test-valuegroup-5",
                  "connector_name": "test-valuegroup-5-name",
                  "component":  "test-valuegroup-correlation-5",
                  "resource": "test-valuegroup-correlation-5-resource-4"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 3
            }
          }
        }
      }
    ]
    """

  Scenario: given meta alarm rule and events shouldn't create metaalarm if empty valuepath
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "description" : "test-correlation-valuegroup-6",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-valuegroup-6"
            }
          }
        ]
      ],
      "enabled" : true,
      "external_data" : {},
      "config": {
        "actions": [
          {
            "type" : "set_entity_info_from_template",
            "name" : "infoenrich7",
            "value" : "{{ `{{.Event.ExtraInfos.infoenrich7}}` }}",
            "description" : "infoenrich7"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "priority" : 10001,
      "type" : "enrichment"
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-valuegroup-correlation-6",
      "type": "valuegroup",
      "config": {
        "time_interval": {
          "value": 10,
          "unit": "s"
        },
        "threshold_count": 2,
        "value_paths": [
          "entity.infos.infoenrich7.value"
        ]
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-valuegroup-6",
      "connector_name": "test-valuegroup-6-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-6",
      "resource": "test-valuegroup-correlation-6-resource-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-6",
      "connector_name": "test-valuegroup-6-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-6",
      "resource": "test-valuegroup-correlation-6-resource-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I do GET /api/v4/alarms?search={{ .metaAlarmRuleID }}&active_columns[]=v.meta&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    When I send an event:
    """
    {
      "connector": "test-valuegroup-6",
      "connector_name": "test-valuegroup-6-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-6",
      "resource": "test-valuegroup-correlation-6-resource-3",
      "infoenrich7": "",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-6",
      "connector_name": "test-valuegroup-6-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-correlation-6",
      "resource": "test-valuegroup-correlation-6-resource-4",
      "infoenrich7": "",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I do GET /api/v4/alarms?search={{ .metaAlarmRuleID }}&active_columns[]=v.meta&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
