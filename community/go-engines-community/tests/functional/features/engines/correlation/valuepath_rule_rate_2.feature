Feature: correlation feature - valuegroup rule with threshold rate

  @concurrent
  Scenario: given meta alarm and removed child should update meta alarm
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "description": "test-eventfilter-correlation-valuegroup-rate-second-1",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-connector-correlation-valuegroup-rate-second-1"
            }
          }
        ]
      ],
      "config": {
        "actions": [
          {
            "type" : "set_entity_info_from_template",
            "name" : "test-info-correlation-valuegroup-rate-second-1",
            "value" : "{{ `{{ index .Event.ExtraInfos \"test-info-correlation-valuegroup-rate-second-1\" }}` }}"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "type" : "enrichment"
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 0,
        "test-info-correlation-valuegroup-rate-second-1": "test",
        "connector": "test-connector-correlation-valuegroup-rate-second-1",
        "connector_name": "test-connector-name-correlation-valuegroup-rate-second-1",
        "component":  "test-component-correlation-valuegroup-rate-second-1",
        "resource": "test-resource-correlation-valuegroup-rate-second-1-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 0,
        "test-info-correlation-valuegroup-rate-second-1": "test",
        "connector": "test-connector-correlation-valuegroup-rate-second-1",
        "connector_name": "test-connector-name-correlation-valuegroup-rate-second-1",
        "component":  "test-component-correlation-valuegroup-rate-second-1",
        "resource": "test-resource-correlation-valuegroup-rate-second-1-2",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 0,
        "test-info-correlation-valuegroup-rate-second-1": "test",
        "connector": "test-connector-correlation-valuegroup-rate-second-1",
        "connector_name": "test-connector-name-correlation-valuegroup-rate-second-1",
        "component":  "test-component-correlation-valuegroup-rate-second-1",
        "resource": "test-resource-correlation-valuegroup-rate-second-1-3",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 0,
        "test-info-correlation-valuegroup-rate-second-1": "test",
        "connector": "test-connector-correlation-valuegroup-rate-second-1",
        "connector_name": "test-connector-name-correlation-valuegroup-rate-second-1",
        "component":  "test-component-correlation-valuegroup-rate-second-1",
        "resource": "test-resource-correlation-valuegroup-rate-second-1-4",
        "source_type": "resource"
      }
    ]
    """
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-correlation-valuegroup-rate-second-1",
      "type": "valuegroup",
      "output_template": "{{ `{{ .Count }}` }}",
      "config": {
        "time_interval": {
          "value": 10,
          "unit": "s"
        },
        "threshold_rate": 0.5,
        "value_paths": [
          "entity.infos.test-info-correlation-valuegroup-rate-second-1.value"
        ]
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleId={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "test-info-correlation-valuegroup-rate-second-1": "test",
      "connector": "test-connector-correlation-valuegroup-rate-second-1",
      "connector_name": "test-connector-name-correlation-valuegroup-rate-second-1",
      "component":  "test-component-correlation-valuegroup-rate-second-1",
      "resource": "test-resource-correlation-valuegroup-rate-second-1-1",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 3,
      "test-info-correlation-valuegroup-rate-second-1": "test",
      "connector": "test-connector-correlation-valuegroup-rate-second-1",
      "connector_name": "test-connector-name-correlation-valuegroup-rate-second-1",
      "component":  "test-component-correlation-valuegroup-rate-second-1",
      "resource": "test-resource-correlation-valuegroup-rate-second-1-2",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-valuegroup-rate-second-1&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-valuegroup-rate-second-1"
          },
          "v": {
            "output": "2",
            "state": {
              "val": 3
            }
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
    When I save response metaAlarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmId }}",
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
                  "connector": "test-connector-correlation-valuegroup-rate-second-1",
                  "connector_name": "test-connector-name-correlation-valuegroup-rate-second-1",
                  "component":  "test-component-correlation-valuegroup-rate-second-1",
                  "resource": "test-resource-correlation-valuegroup-rate-second-1-1"
                }
              },
              {
                "v": {
                  "connector": "test-connector-correlation-valuegroup-rate-second-1",
                  "connector_name": "test-connector-name-correlation-valuegroup-rate-second-1",
                  "component":  "test-component-correlation-valuegroup-rate-second-1",
                  "resource": "test-resource-correlation-valuegroup-rate-second-1-2"
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
    When I save response childAlarmId2={{ (index (index .lastResponse 0).data.children.data 1)._id }}
    When I do PUT /api/v4/cat/meta-alarms/{{ .metaAlarmId }}/remove:
    """json
    {
      "comment": "test-metaalarmrule-correlation-valuegroup-rate-second-1-comment",
      "alarms": ["{{ .childAlarmId2 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-correlation-valuegroup-rate-second-1&correlation=true&sort_by=v.meta&sort=desc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 1,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-valuegroup-rate-second-1"
          },
          "v": {
            "output": "1",
            "state": {
              "val": 2
            }
          }
        },
        {
          "v": {
            "connector": "test-connector-correlation-valuegroup-rate-second-1",
            "connector_name": "test-connector-name-correlation-valuegroup-rate-second-1",
            "component":  "test-component-correlation-valuegroup-rate-second-1",
            "resource": "test-resource-correlation-valuegroup-rate-second-1-2"
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
        "_id": "{{ .metaAlarmId }}",
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
                  "connector": "test-connector-correlation-valuegroup-rate-second-1",
                  "connector_name": "test-connector-name-correlation-valuegroup-rate-second-1",
                  "component":  "test-component-correlation-valuegroup-rate-second-1",
                  "resource": "test-resource-correlation-valuegroup-rate-second-1-1"
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
        }
      }
    ]
    """

  @concurrent
  Scenario: given meta alarm and removed child should not add child to meta alarm again but add to new meta alarm on change
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "description": "test-eventfilter-correlation-valuegroup-rate-second-2",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-connector-correlation-valuegroup-rate-second-2"
            }
          }
        ]
      ],
      "config": {
        "actions": [
          {
            "type" : "set_entity_info_from_template",
            "name" : "test-info-correlation-valuegroup-rate-second-2",
            "value" : "{{ `{{ index .Event.ExtraInfos \"test-info-correlation-valuegroup-rate-second-2\" }}` }}"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "type" : "enrichment"
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 0,
        "test-info-correlation-valuegroup-rate-second-2": "test-1",
        "connector": "test-connector-correlation-valuegroup-rate-second-2",
        "connector_name": "test-connector-name-correlation-valuegroup-rate-second-2",
        "component":  "test-component-correlation-valuegroup-rate-second-2",
        "resource": "test-resource-correlation-valuegroup-rate-second-2-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 0,
        "test-info-correlation-valuegroup-rate-second-2": "test-1",
        "connector": "test-connector-correlation-valuegroup-rate-second-2",
        "connector_name": "test-connector-name-correlation-valuegroup-rate-second-2",
        "component":  "test-component-correlation-valuegroup-rate-second-2",
        "resource": "test-resource-correlation-valuegroup-rate-second-2-2",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 0,
        "test-info-correlation-valuegroup-rate-second-2": "test-1",
        "connector": "test-connector-correlation-valuegroup-rate-second-2",
        "connector_name": "test-connector-name-correlation-valuegroup-rate-second-2",
        "component":  "test-component-correlation-valuegroup-rate-second-2",
        "resource": "test-resource-correlation-valuegroup-rate-second-2-3",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 0,
        "test-info-correlation-valuegroup-rate-second-2": "test-1",
        "connector": "test-connector-correlation-valuegroup-rate-second-2",
        "connector_name": "test-connector-name-correlation-valuegroup-rate-second-2",
        "component":  "test-component-correlation-valuegroup-rate-second-2",
        "resource": "test-resource-correlation-valuegroup-rate-second-2-4",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 0,
        "test-info-correlation-valuegroup-rate-second-2": "test-2",
        "connector": "test-connector-correlation-valuegroup-rate-second-2",
        "connector_name": "test-connector-name-correlation-valuegroup-rate-second-2",
        "component":  "test-component-correlation-valuegroup-rate-second-2",
        "resource": "test-resource-correlation-valuegroup-rate-second-2-5",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 0,
        "test-info-correlation-valuegroup-rate-second-2": "test-2",
        "connector": "test-connector-correlation-valuegroup-rate-second-2",
        "connector_name": "test-connector-name-correlation-valuegroup-rate-second-2",
        "component":  "test-component-correlation-valuegroup-rate-second-2",
        "resource": "test-resource-correlation-valuegroup-rate-second-2-6",
        "source_type": "resource"
      }
    ]
    """
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-correlation-valuegroup-rate-second-2",
      "type": "valuegroup",
      "output_template": "{{ `{{ .Count }}` }}",
      "config": {
        "time_interval": {
          "value": 3,
          "unit": "s"
        },
        "threshold_rate": 0.5,
        "value_paths": [
          "entity.infos.test-info-correlation-valuegroup-rate-second-2.value"
        ]
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleId={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "test-info-correlation-valuegroup-rate-second-2": "test-1",
      "connector": "test-connector-correlation-valuegroup-rate-second-2",
      "connector_name": "test-connector-name-correlation-valuegroup-rate-second-2",
      "component":  "test-component-correlation-valuegroup-rate-second-2",
      "resource": "test-resource-correlation-valuegroup-rate-second-2-1",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 3,
      "test-info-correlation-valuegroup-rate-second-2": "test-1",
      "connector": "test-connector-correlation-valuegroup-rate-second-2",
      "connector_name": "test-connector-name-correlation-valuegroup-rate-second-2",
      "component":  "test-component-correlation-valuegroup-rate-second-2",
      "resource": "test-resource-correlation-valuegroup-rate-second-2-2",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-valuegroup-rate-second-2&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-valuegroup-rate-second-2"
          },
          "v": {
            "output": "2",
            "state": {
              "val": 3
            }
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
    When I save response metaAlarmId1={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmId1 }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      }
    ]
    """
    Then the response code should be 207
    When I save response childAlarmId2={{ (index (index .lastResponse 0).data.children.data 1)._id }}
    When I do PUT /api/v4/cat/meta-alarms/{{ .metaAlarmId1 }}/remove:
    """json
    {
      "comment": "test-metaalarmrule-correlation-valuegroup-rate-second-2-comment",
      "alarms": ["{{ .childAlarmId2 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-correlation-valuegroup-rate-second-2&correlation=true&sort_by=v.meta&sort=desc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 1,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-valuegroup-rate-second-2"
          },
          "v": {
            "output": "1",
            "state": {
              "val": 2
            }
          }
        },
        {
          "v": {
            "resource": "test-resource-correlation-valuegroup-rate-second-2-2"
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
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "test-info-correlation-valuegroup-rate-second-2": "test-1",
      "connector": "test-connector-correlation-valuegroup-rate-second-2",
      "connector_name": "test-connector-name-correlation-valuegroup-rate-second-2",
      "component":  "test-component-correlation-valuegroup-rate-second-2",
      "resource": "test-resource-correlation-valuegroup-rate-second-2-2",
      "source_type": "resource"
    }
    """
    When I wait 1s
    When I do GET /api/v4/alarms?search=test-resource-correlation-valuegroup-rate-second-2&correlation=true&sort_by=v.meta&sort=desc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 1,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-valuegroup-rate-second-2"
          },
          "v": {
            "output": "1",
            "state": {
              "val": 2
            }
          }
        },
        {
          "v": {
            "resource": "test-resource-correlation-valuegroup-rate-second-2-2"
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
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "test-info-correlation-valuegroup-rate-second-2": "test-1",
      "connector": "test-connector-correlation-valuegroup-rate-second-2",
      "connector_name": "test-connector-name-correlation-valuegroup-rate-second-2",
      "component":  "test-component-correlation-valuegroup-rate-second-2",
      "resource": "test-resource-correlation-valuegroup-rate-second-2-3",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-valuegroup-rate-second-2&correlation=true&sort_by=v.meta&sort=desc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-valuegroup-rate-second-2"
          },
          "v": {
            "output": "2",
            "state": {
              "val": 2
            }
          }
        },
        {
          "v": {
            "resource": "test-resource-correlation-valuegroup-rate-second-2-2"
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
    When I wait 2s
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 3,
      "test-info-correlation-valuegroup-rate-second-2": "test-2",
      "connector": "test-connector-correlation-valuegroup-rate-second-2",
      "connector_name": "test-connector-name-correlation-valuegroup-rate-second-2",
      "component":  "test-component-correlation-valuegroup-rate-second-2",
      "resource": "test-resource-correlation-valuegroup-rate-second-2-2",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "test-info-correlation-valuegroup-rate-second-2": "test-2",
      "connector": "test-connector-correlation-valuegroup-rate-second-2",
      "connector_name": "test-connector-name-correlation-valuegroup-rate-second-2",
      "component":  "test-component-correlation-valuegroup-rate-second-2",
      "resource": "test-resource-correlation-valuegroup-rate-second-2-5",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-valuegroup-rate-second-2&correlation=true&sort_by=t&sort=desc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-valuegroup-rate-second-2"
          },
          "v": {
            "output": "2",
            "state": {
              "val": 3
            }
          }
        },
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-valuegroup-rate-second-2"
          },
          "v": {
            "output": "2",
            "state": {
              "val": 2
            }
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
    When I save response metaAlarmId2={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmId2 }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      },
      {
        "_id": "{{ .metaAlarmId1 }}",
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
                  "resource": "test-resource-correlation-valuegroup-rate-second-2-2"
                }
              },
              {
                "v": {
                  "resource": "test-resource-correlation-valuegroup-rate-second-2-5"
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
                  "resource": "test-resource-correlation-valuegroup-rate-second-2-1"
                }
              },
              {
                "v": {
                  "resource": "test-resource-correlation-valuegroup-rate-second-2-3"
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
