Feature: update connector alarm
  I need to be able to update connector alarm which is created by idle rule

  Scenario: given connector alarm should enrich entity infos
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-connector-1-name",
      "type": "entity",
      "enabled": true,
      "priority": 60,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-connector-name-axe-idlerule-connector-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "source_type",
            "cond": {
              "type": "eq",
              "value": "connector"
            }
          },
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-connector-axe-idlerule-connector-1"
            }
          }
        ]
      ],
      "config": {
        "actions": [
          {
            "type": "set_entity_info_from_template",
            "name": "client",
            "description": "Client",
            "value": "{{ `{{ .Event.ConnectorName }}` }}"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "priority": 2,
      "enabled": true,
      "description": "test-eventfilter-axe-idlerule-connector-1-description"
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-connector-1",
      "connector_name": "test-connector-name-axe-idlerule-connector-1",
      "source_type": "component",
      "component":  "test-component-axe-idlerule-connector-1",
      "state": 2,
      "output": "test-output-axe-idlerule-connector-1",
      "long_output": "test-long-output-axe-idlerule-connector-1"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-connector-name-axe-idlerule-connector-1&sort_by=entity._id&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "type": "connector",
            "infos": {
              "client": {
                "description": "Client",
                "name": "client",
                "value": "test-connector-name-axe-idlerule-connector-1"
              }
            }
          },
          "v": {
            "connector": "test-connector-axe-idlerule-connector-1",
            "connector_name": "test-connector-name-axe-idlerule-connector-1"
          }
        },
        {}
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """

  Scenario: given connector alarm should apply pbehavior on it
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-connector-2-name",
      "type": "entity",
      "enabled": true,
      "priority": 60,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-connector-name-axe-idlerule-connector-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-axe-idlerule-connector-2",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "eq",
              "value": "test-connector-axe-idlerule-connector-2/test-connector-name-axe-idlerule-connector-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-connector-2",
      "connector_name": "test-connector-name-axe-idlerule-connector-2",
      "source_type": "component",
      "component":  "test-component-axe-idlerule-connector-2",
      "state": 1,
      "output": "test-output-axe-idlerule-connector-2"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-connector-name-axe-idlerule-connector-2&sort_by=entity._id&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "type": "connector"
          },
          "v": {
            "connector": "test-connector-axe-idlerule-connector-2",
            "connector_name": "test-connector-name-axe-idlerule-connector-2",
            "pbehavior_info": {
              "canonical_type": "maintenance",
              "name": "test-pbehavior-axe-idlerule-connector-2",
              "reason": "test-reason-to-engine",
              "reason_name": "Test Engine",
              "type": "test-maintenance-type-to-engine",
              "type_name": "Engine maintenance"
            }
          }
        },
        {}
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
        "steps": {
          "page": 1
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
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 5
              },
              {
                "_t": "pbhenter",
                "m": "Pbehavior test-pbehavior-axe-idlerule-connector-2. Type: Engine maintenance. Reason: Test Engine."
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

  Scenario: given connector alarm should update service alarm
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-connector-3-name",
      "type": "entity",
      "enabled": true,
      "priority": 60,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-connector-name-axe-idlerule-connector-3"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-service-axe-idlerule-connector-3",
      "output_template": "All: {{ `{{.All}}` }}; Active: {{ `{{.Active}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; AcknowledgedUnderPbh: {{ `{{.AcknowledgedUnderPbh}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateOk: {{ `{{.State.Ok}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }}; UnderPbehavior: {{ `{{.UnderPbehavior}}` }};",
      "enabled": true,
      "impact_level": 1,
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "eq",
              "value": "test-connector-axe-idlerule-connector-3/test-connector-name-axe-idlerule-connector-3"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-connector-3",
      "connector_name": "test-connector-name-axe-idlerule-connector-3",
      "source_type": "component",
      "component":  "test-component-axe-idlerule-connector-3",
      "state": 1,
      "output": "test-output-axe-idlerule-connector-3"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-connector-name-axe-idlerule-connector-3&sort_by=entity._id&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "type": "connector"
          },
          "v": {
            "connector": "test-connector-axe-idlerule-connector-3",
            "connector_name": "test-connector-name-axe-idlerule-connector-3"
          }
        },
        {}
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """
    When I do GET /api/v4/alarms?search={{ .serviceID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "{{ .serviceID }}",
            "connector": "service",
            "connector_name": "service",
            "state": {
              "val": 3
            },
            "status": {
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
        "steps": {
          "page": 1
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
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "a": "service.service",
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Active: 1; Acknowledged: 0; NotAcknowledged: 1; AcknowledgedUnderPbh: 0; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateOk: 0; Pbehaviors: map[]; UnderPbehavior: 0;",
                "val": 1
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

  Scenario: given connector alarm should apply scenario on it
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-connector-4-name",
      "type": "entity",
      "enabled": true,
      "priority": 60,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-connector-name-axe-idlerule-connector-4"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-axe-idlerule-connector-4-name",
      "priority": 10059,
      "enabled": true,
      "priority": 60,
      "triggers": ["create"],
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-connector-name-axe-idlerule-connector-4"
                }
              }
            ]
          ],
          "type": "assocticket",
          "parameters": {
            "ticket": "test-ticket-axe-idlerule-connector-4"
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        },
        {
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-connector-name-axe-idlerule-connector-4"
                }
              }
            ]
          ],
          "type": "ack",
          "parameters": {
            "output": "test-output-axe-idlerule-connector-4"
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-connector-4",
      "connector_name": "test-connector-name-axe-idlerule-connector-4",
      "source_type": "component",
      "component":  "test-component-axe-idlerule-connector-4",
      "state": 1,
      "output": "test-output-axe-idlerule-connector-4"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-connector-name-axe-idlerule-connector-4&sort_by=entity._id&sort=desc
    Then the response code should be 200
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1
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
          "steps": {
            "data": [
              {
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "assocticket",
                "a": "system",
                "ticket": "test-ticket-axe-idlerule-connector-4"
              },
              {
                "_t": "ack",
                "a": "system",
                "m": "test-output-axe-idlerule-connector-4"
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

  Scenario: given connector alarm should apply dynamic infos on it
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-connector-5-name",
      "type": "entity",
      "enabled": true,
      "priority": 60,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-connector-name-axe-idlerule-connector-5"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "name": "test-dynamicinfos-axe-idlerule-connector-5-name",
      "description": "test-dynamicinfos-axe-idlerule-connector-5-description",
      "disable_during_periods": [],
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "eq",
              "value": "test-connector-axe-idlerule-connector-5/test-connector-name-axe-idlerule-connector-5"
            }
          }
        ]
      ],
      "infos": [
        {"name":"test-info-axe-idlerule-connector-5-name", "value":"test-info-axe-idlerule-connector-5-value"}
      ],
      "enabled": true
    }
    """
    Then the response code should be 201
    When I save response dynamicInfosRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-connector-5",
      "connector_name": "test-connector-name-axe-idlerule-connector-5",
      "source_type": "component",
      "component":  "test-component-axe-idlerule-connector-5",
      "state": 2,
      "output": "test-output-axe-idlerule-connector-5"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-connector-name-axe-idlerule-connector-5&sort_by=entity._id&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "type": "connector"
          },
          "v": {
            "connector": "test-connector-axe-idlerule-connector-5",
            "connector_name": "test-connector-name-axe-idlerule-connector-5",
            "infos": {
              "{{ .dynamicInfosRuleID }}": {
                "test-info-axe-idlerule-connector-5-name": "test-info-axe-idlerule-connector-5-value"
              }
            }
          }
        },
        {}
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """

  Scenario: given connector alarm should not create meta alarm
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-connector-6-name",
      "type": "entity",
      "enabled": true,
      "priority": 44,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-connector-name-axe-idlerule-connector-6"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-action-correlation-1",
      "type": "attribute",
      "alarm_pattern": [
        [
          {
            "field": "v.connector",
            "cond": {
              "type": "eq",
              "value": "test-connector-axe-idlerule-connector-6"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-connector-6",
      "connector_name": "test-connector-name-axe-idlerule-connector-6",
      "source_type": "component",
      "component":  "test-component-axe-idlerule-connector-6",
      "state": 1,
      "output": "test-output-axe-idlerule-connector-6"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-component-axe-idlerule-connector-6&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 1,
          "v": {
            "meta": "{{ .metaAlarmRuleID }}"
          }
        }
      ]
    }
    """
    When I save response metalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I do GET /api/v4/alarms?search=test-component-axe-idlerule-connector-6
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-axe-idlerule-connector-6",
            "connector_name": "test-connector-name-axe-idlerule-connector-6",
            "component": "test-component-axe-idlerule-connector-6",
            "parents": [
              "{{ .metalarmEntityID }}"
            ]
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
    When I do GET /api/v4/alarms?search=test-connector-name-axe-idlerule-connector-6&sort_by=entity._id&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "type": "connector"
          },
          "v": {
            "connector": "test-connector-axe-idlerule-connector-6",
            "connector_name": "test-connector-name-axe-idlerule-connector-6",
            "parents": []
          }
        },
        {}
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """

  Scenario: given connector alarm should update alarm on event
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-connector-7-name",
      "type": "entity",
      "enabled": true,
      "priority": 44,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-connector-name-axe-idlerule-connector-7"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-connector-7",
      "connector_name": "test-connector-name-axe-idlerule-connector-7",
      "source_type": "component",
      "component":  "test-component-axe-idlerule-connector-7",
      "state": 1,
      "output": "test-output-axe-idlerule-connector-7"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "event_type": "ack",
      "connector": "test-connector-axe-idlerule-connector-7",
      "connector_name": "test-connector-name-axe-idlerule-connector-7",
      "source_type": "connector",
      "output": "test-output-axe-idlerule-connector-7"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-connector-name-axe-idlerule-connector-7&sort_by=entity._id&sort=desc
    Then the response code should be 200
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1
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
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 5
              },
              {
                "_t": "ack"
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
