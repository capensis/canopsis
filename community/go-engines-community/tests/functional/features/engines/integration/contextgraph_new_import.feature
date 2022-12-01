Feature: new-import entities
  I need to be able to new-import entities

  Scenario: given delete import action should delete component and resources which should update service state
    When I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-entityservice-new-import-partial-6",
      "name": "test-entityservice-new-import-partial-6-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-new-import-partial-6-1"
            }
          }
        ],
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-new-import-partial-6-2"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-new-import-partial-6",
      "connector_name": "test-connector-name-new-import-partial-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-6-1",
      "resource": "test-resource-new-import-partial-6-1",
      "state": 1,
      "output": "test-output-import-partial-6"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-new-import-partial-6",
      "connector_name": "test-connector-name-new-import-partial-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-6-1",
      "resource": "test-resource-new-import-partial-6-2",
      "state": 1,
      "output": "test-output-import-partial-6"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-new-import-partial-6",
      "connector_name": "test-connector-name-new-import-partial-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-6-2",
      "resource": "test-resource-new-import-partial-6-3",
      "state": 3,
      "output": "test-output-import-partial-6"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-new-import-partial-6",
      "connector_name": "test-connector-name-new-import-partial-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-6-2",
      "resource": "test-resource-new-import-partial-6-4",
      "state": 3,
      "output": "test-output-import-partial-6"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-entityservice-new-import-partial-6&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-entityservice-new-import-partial-6",
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
    When I do GET /api/v4/weather-services/test-entityservice-new-import-partial-6?sort_by=name&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-new-import-partial-6-1"
        },
        {
          "_id": "test-component-new-import-partial-6-2"
        },
        {
          "_id": "test-resource-new-import-partial-6-1/test-component-new-import-partial-6-1"
        },
        {
          "_id": "test-resource-new-import-partial-6-2/test-component-new-import-partial-6-1"
        },
        {
          "_id": "test-resource-new-import-partial-6-3/test-component-new-import-partial-6-2"
        },
        {
          "_id": "test-resource-new-import-partial-6-4/test-component-new-import-partial-6-2"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 6
      }
    }
    """
    When I do PUT /api/v4/contextgraph-import-partial?source=test-new-import-partial-6-source:
    """json
    {
      "cis": [
        {
          "action": "delete",
          "name": "test-component-new-import-partial-6-2",
          "type": "component",
          "enabled": true
        }
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
       "status": "done"
    }
    """
    When I wait the end of 4 events processing
    When I do GET /api/v4/alarms?search=test-entityservice-new-import-partial-6&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-entityservice-new-import-partial-6",
            "state": {
              "val": 1
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
    When I do GET /api/v4/weather-services/test-entityservice-new-import-partial-6?sort_by=name&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-new-import-partial-6-1"
        },
        {
          "_id": "test-resource-new-import-partial-6-1/test-component-new-import-partial-6-1"
        },
        {
          "_id": "test-resource-new-import-partial-6-2/test-component-new-import-partial-6-1"
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

  Scenario: given delete import action should delete component and resources which should update metaalarm state
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "_id": "test-metaalarm-new-import-partial-7",
      "name": "test-metaalarm-new-import-partial-7-name",
      "type": "complex",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-new-import-partial-7-1"
            }
          }
        ],
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-new-import-partial-7-2"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 10,
          "unit": "s"
        },
        "threshold_count": 4
      }
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-new-import-partial-7",
      "connector_name": "test-connector-name-new-import-partial-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-7-1",
      "resource": "test-resource-new-import-partial-7-1",
      "state": 1,
      "output": "test-output-import-partial-7"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-new-import-partial-7",
      "connector_name": "test-connector-name-new-import-partial-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-7-1",
      "resource": "test-resource-new-import-partial-7-2",
      "state": 1,
      "output": "test-output-import-partial-7"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-new-import-partial-7",
      "connector_name": "test-connector-name-new-import-partial-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-7-2",
      "resource": "test-resource-new-import-partial-7-3",
      "state": 3,
      "output": "test-output-import-partial-7"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-new-import-partial-7",
      "connector_name": "test-connector-name-new-import-partial-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-7-2",
      "resource": "test-resource-new-import-partial-7-4",
      "state": 3,
      "output": "test-output-import-partial-7"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-new-import-partial-7-1&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-metaalarm-new-import-partial-7-name"
          },
          "v": {
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
                  "component": "test-component-new-import-partial-7-1",
                  "resource": "test-resource-new-import-partial-7-1"
                }
              },
              {
                "v": {
                  "component": "test-component-new-import-partial-7-1",
                  "resource": "test-resource-new-import-partial-7-2"
                }
              },
              {
                "v": {
                  "component": "test-component-new-import-partial-7-2",
                  "resource": "test-resource-new-import-partial-7-3"
                }
              },
              {
                "v": {
                  "component": "test-component-new-import-partial-7-2",
                  "resource": "test-resource-new-import-partial-7-4"
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
        }
      }
    ]
    """
    When I do PUT /api/v4/contextgraph-import-partial?source=test-new-import-partial-7-source:
    """json
    {
      "cis": [
        {
          "action": "delete",
          "name": "test-component-new-import-partial-7-2",
          "type": "component",
          "enabled": true
        }
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
       "status": "done"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-new-import-partial-7-1&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-metaalarm-new-import-partial-7-name"
          },
          "v": {
            "state": {
              "val": 1
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
                  "component": "test-component-new-import-partial-7-1",
                  "resource": "test-resource-new-import-partial-7-1"
                }
              },
              {
                "v": {
                  "component": "test-component-new-import-partial-7-1",
                  "resource": "test-resource-new-import-partial-7-2"
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
    
  Scenario: given delete import action should delete component and resources which should resolve alarm
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "_id": "test-metaalarm-new-import-partial-8",
      "name": "test-metaalarm-new-import-partial-8-name",
      "type": "complex",
      "auto_resolve": true,
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-new-import-partial-8-1"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 10,
          "unit": "s"
        },
        "threshold_count": 2
      }
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-new-import-partial-8",
      "connector_name": "test-connector-name-new-import-partial-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-8-1",
      "resource": "test-resource-new-import-partial-8-1",
      "state": 1,
      "output": "test-output-import-partial-8"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-new-import-partial-8",
      "connector_name": "test-connector-name-new-import-partial-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-8-1",
      "resource": "test-resource-new-import-partial-8-2",
      "state": 1,
      "output": "test-output-import-partial-8"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-new-import-partial-8-1&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-metaalarm-new-import-partial-8-name"
          },
          "v": {
            "state": {
              "val": 1
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
                  "component": "test-component-new-import-partial-8-1",
                  "resource": "test-resource-new-import-partial-8-1"
                }
              },
              {
                "v": {
                  "component": "test-component-new-import-partial-8-1",
                  "resource": "test-resource-new-import-partial-8-2"
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
    When I do PUT /api/v4/contextgraph-import-partial?source=test-new-import-partial-8-source:
    """json
    {
      "cis": [
        {
          "action": "delete",
          "name": "test-component-new-import-partial-8-1",
          "type": "component",
          "enabled": true
        }
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
       "status": "done"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-new-import-partial-8-1&correlation=true until response code is 200 and response key "data.0.v.resolved" is greater or equal than 1
