Feature: Get alarms
  I need to be able to get a alarms

  @concurrent
  Scenario: given manual metaalarm and alarms should increase events_count and last_event_date in metaalarm
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-correlation-alarm-get-second-1",
      "connector_name": "test-correlation-alarm-get-second-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-correlation-alarm-get-second-1",
      "resource": "test-resource-correlation-alarm-get-second-1-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-second-1-1&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-correlation-alarm-get-second-1-1"
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
    When I save response alarmId1={{ (index .lastResponse.data 0)._id }}
    Then I save response lastEventDate1={{ (index .lastResponse.data 0).v.last_event_date }}
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metaalarm-correlation-alarm-get-second-1",
      "comment": "test-metaalarm-correlation-alarm-get-second-1-comment",
      "alarms": ["{{ .alarmId1 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/cat/manual-meta-alarms?search=test-metaalarm-correlation-alarm-get-second-1 until response code is 200 and body contains:
    """json
    [
      {
        "name": "test-metaalarm-correlation-alarm-get-second-1"
      }
    ]
    """
    Then I save response manualMetaAlarm={{ (index .lastResponse 0)._id }}
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-second-1&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "v": {
            "events_count": 1
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
    Then the response key "data.0.v.last_event_date" should be "{{ .lastEventDate1 }}"
    When I wait 1s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-correlation-alarm-get-second-1",
      "connector_name": "test-correlation-alarm-get-second-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-correlation-alarm-get-second-1",
      "resource": "test-resource-correlation-alarm-get-second-1-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-correlation-alarm-get-second-1",
      "connector_name": "test-correlation-alarm-get-second-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-correlation-alarm-get-second-1",
      "resource": "test-resource-correlation-alarm-get-second-1-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-second-1-2&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-correlation-alarm-get-second-1-2"
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
    When I save response alarmId2={{ (index .lastResponse.data 0)._id }}
    Then I save response lastEventDate2={{ (index .lastResponse.data 0).v.last_event_date }}
    When I wait 1s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-correlation-alarm-get-second-1",
      "connector_name": "test-correlation-alarm-get-second-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-correlation-alarm-get-second-1",
      "resource": "test-resource-correlation-alarm-get-second-1-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-correlation-alarm-get-second-1",
      "connector_name": "test-correlation-alarm-get-second-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-correlation-alarm-get-second-1",
      "resource": "test-resource-correlation-alarm-get-second-1-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-correlation-alarm-get-second-1",
      "connector_name": "test-correlation-alarm-get-second-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-correlation-alarm-get-second-1",
      "resource": "test-resource-correlation-alarm-get-second-1-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-second-1-3&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-correlation-alarm-get-second-1-3"
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
    When I save response alarmId3={{ (index .lastResponse.data 0)._id }}
    Then I save response lastEventDate3={{ (index .lastResponse.data 0).v.last_event_date }}
    When I do PUT /api/v4/cat/manual-meta-alarms/{{ .manualMetaAlarm }}/add:
    """json
    {
      "comment": "test-metaalarm-correlation-alarm-get-second-1-comment",
      "alarms": ["{{ .alarmId2 }}", "{{ .alarmId3 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-second-1&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "v": {
            "events_count": 6
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
    Then the response key "data.0.v.last_event_date" should be "{{ .lastEventDate3 }}"
    When I do PUT /api/v4/cat/manual-meta-alarms/{{ .manualMetaAlarm }}/remove:
    """json
    {
      "comment": "test-metaalarm-correlation-alarm-get-second-1-comment",
      "alarms": ["{{ .alarmId2 }}", "{{ .alarmId3 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-second-1&correlation=true&sort_by=v.resource&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "v": {
            "events_count": 1
          }
        },
        {
        },
        {
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """
    Then the response key "data.0.v.last_event_date" should be "{{ .lastEventDate1 }}"

  @concurrent
  Scenario: given get search correlation request with entity infos should return filtered children
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "event_pattern":[[
        {
          "field": "event_type",
          "cond": {
            "type": "eq",
            "value": "check"
          }
        },
        {
          "field": "component",
          "cond": {
            "type": "eq",
            "value": "test-component-to-alarm-correlation-get-second-2"
          }
        }
      ]],
      "config": {
        "actions": [
          {
            "type": "set_entity_info_from_template",
            "name": "output",
            "description": "Output",
            "value": "{{ `{{ .Event.Output }}` }}"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "description": "test-event-filter-to-alarm-correlation-get-second-2-description",
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-to-alarm-correlation-get-second-2",
        "connector_name": "test-connector-name-to-alarm-correlation-get-second-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-second-2",
        "resource": "test-resource-to-alarm-correlation-get-second-2-1",
        "state": 1,
        "output": "test-output-to-alarm-correlation-get-second-2"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-second-2",
        "connector_name": "test-connector-name-to-alarm-correlation-get-second-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-second-2",
        "resource": "test-resource-to-alarm-correlation-get-second-2-2",
        "state": 1,
        "output": "test-resource-to-alarm-correlation-get-second-2-search"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-second-2",
        "connector_name": "test-connector-name-to-alarm-correlation-get-second-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-second-2",
        "resource": "test-resource-to-alarm-correlation-get-second-2-3",
        "state": 1,
        "output": "test-resource-to-alarm-correlation-get-second-2-search"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-second-2",
        "connector_name": "test-connector-name-to-alarm-correlation-get-second-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-second-2",
        "resource": "test-resource-to-alarm-correlation-get-second-2-4",
        "state": 1,
        "output": "test-resource-to-alarm-correlation-get-second-2-search"
      },
      {
        "connector": "test-connector-to-alarm-correlation-get-second-2",
        "connector_name": "test-connector-name-to-alarm-correlation-get-second-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-correlation-get-second-2",
        "resource": "test-resource-to-alarm-correlation-get-second-2-5",
        "state": 1,
        "output": "test-resource-to-alarm-correlation-get-second-2-search"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-correlation-get-second-2&sort_by=v.resource&sort=asc
    Then the response code should be 200
    When I save response childAlarmID1={{ (index .lastResponse.data 0)._id }}
    When I save response childAlarmID2={{ (index .lastResponse.data 1)._id }}
    When I save response childAlarmID3={{ (index .lastResponse.data 2)._id }}
    When I save response childAlarmID4={{ (index .lastResponse.data 3)._id }}
    When I save response childAlarmID5={{ (index .lastResponse.data 4)._id }}
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metalarm-to-alarm-correlation-get-second-2-1",
      "comment": "test-metalarm-to-alarm-correlation-get-second-2-1-comment",
      "alarms": [
        "{{ .childAlarmID1 }}",
        "{{ .childAlarmID2 }}",
        "{{ .childAlarmID3 }}"
      ]
    }
    """
    Then the response code should be 204
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metalarm-to-alarm-correlation-get-second-2-2",
      "comment": "test-metalarm-to-alarm-correlation-get-second-2-2-comment",
      "alarms": [
        "{{ .childAlarmID4 }}"
      ]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-correlation-get-second-2&correlation=true&sort_by=v.output&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 3
        },
        {
          "children": 1
        },
        {}
      ]
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-alarm-correlation-get-second-2",
      "connector_name": "test-connector-name-to-alarm-correlation-get-second-2",
      "source_type": "resource",
      "event_type": "cancel",
      "component": "test-component-to-alarm-correlation-get-second-2",
      "resource": "test-resource-to-alarm-correlation-get-second-2-3"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-alarm-correlation-get-second-2",
      "connector_name": "test-connector-name-to-alarm-correlation-get-second-2",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component": "test-component-to-alarm-correlation-get-second-2",
      "resource": "test-resource-to-alarm-correlation-get-second-2-3"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-correlation-get-second-2-search&active_columns[]=entity.infos.output.value&correlation=true&sort_by=children&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 3
        },
        {
          "is_meta_alarm": true,
          "children": 1
        },
        {
          "is_meta_alarm": false,
          "v": {
            "resource": "test-resource-to-alarm-correlation-get-second-2-5"
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
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "search": "test-resource-to-alarm-correlation-get-second-2-search",
        "search_by": ["entity.infos.output.value"],
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      },
      {
        "_id": "{{ (index .lastResponse.data 1)._id }}",
        "search": "test-resource-to-alarm-correlation-get-second-2-search",
        "search_by": ["entity.infos.output.value"],
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
                  "resource": "test-resource-to-alarm-correlation-get-second-2-1"
                },
                "filtered": false
              },
              {
                "v": {
                  "resource": "test-resource-to-alarm-correlation-get-second-2-2"
                },
                "filtered": true
              },
              {
                "v": {
                  "resource": "test-resource-to-alarm-correlation-get-second-2-3"
                },
                "filtered": true
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
                  "resource": "test-resource-to-alarm-correlation-get-second-2-4"
                },
                "filtered": true
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
    When I do GET /api/v4/alarms?opened=true&search=test-resource-to-alarm-correlation-get-second-2-search&active_columns[]=entity.infos.output.value&correlation=true&sort_by=children&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 3
        },
        {
          "is_meta_alarm": true,
          "children": 1
        },
        {
          "is_meta_alarm": false,
          "v": {
            "resource": "test-resource-to-alarm-correlation-get-second-2-5"
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
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "opened": true,
        "search": "test-resource-to-alarm-correlation-get-second-2-search",
        "search_by": ["entity.infos.output.value"],
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      },
      {
        "_id": "{{ (index .lastResponse.data 1)._id }}",
        "opened": true,
        "search": "test-resource-to-alarm-correlation-get-second-2-search",
        "search_by": ["entity.infos.output.value"],
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
                  "resource": "test-resource-to-alarm-correlation-get-second-2-1"
                },
                "filtered": false
              },
              {
                "v": {
                  "resource": "test-resource-to-alarm-correlation-get-second-2-2"
                },
                "filtered": true
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
                  "resource": "test-resource-to-alarm-correlation-get-second-2-4"
                },
                "filtered": true
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
