Feature: update connector alarm
  I need to be able to update connector alarm which is created by idle rule

  Scenario: given connector alarm should enrich entity infos
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """
    {
      "name": "test-idlerule-axe-idlerule-connector-1-name",
      "author": "test-idlerule-axe-idlerule-connector-1-author",
      "type": "entity",
      "enabled": true,
      "priority": 60,
      "duration": {
        "seconds": 3,
        "unit": "s"
      },
      "entity_patterns": [
        {
          "_id": "test-connector-axe-idlerule-connector-1/test-connector-name-axe-idlerule-connector-1"
        }
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type": "enrichment",
      "patterns": [{
        "source_type": "connector",
        "connector": "test-connector-axe-idlerule-connector-1"
      }],
      "external_data": {
        "entity": {
          "type": "entity"
        }
      },
      "actions": [
        {
          "type": "copy",
          "from": "ExternalData.entity",
          "to": "Entity"
        }
      ],
      "on_success": "pass",
      "on_failure": "pass",
      "priority": 1,
      "enabled": true,
      "author": "test-eventfilter-axe-idlerule-connector-1-author",
      "description": "test-eventfilter-axe-idlerule-connector-1-description"
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type": "enrichment",
      "patterns": [{
        "source_type": "connector",
        "connector": "test-connector-axe-idlerule-connector-1"
      }],
      "actions": [
        {
          "type": "set_entity_info_from_template",
          "name": "client",
          "description": "Client",
          "value": "{{ `{{ .Event.ConnectorName }}` }}"
        }
      ],
      "priority": 2,
      "on_success": "pass",
      "on_failure": "pass",
      "enabled": true,
      "author": "test-eventfilter-axe-idlerule-connector-1-author",
      "description": "test-eventfilter-axe-idlerule-connector-1-description"
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """
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
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-connector-axe-idlerule-connector-1/test-connector-name-axe-idlerule-connector-1"}]}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "entity": {
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

  Scenario: given connector alarm should apply pbehavior on it
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """
    {
      "name": "test-idlerule-axe-idlerule-connector-2-name",
      "author": "test-idlerule-axe-idlerule-connector-2-author",
      "type": "entity",
      "enabled": true,
      "priority": 60,
      "duration": {
        "seconds": 3,
        "unit": "s"
      },
      "entity_patterns": [
        {
          "_id": "test-connector-axe-idlerule-connector-2/test-connector-name-axe-idlerule-connector-2"
        }
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "author": "root",
      "name": "test-pbehavior-axe-idlerule-connector-2",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "_id": "test-connector-axe-idlerule-connector-2/test-connector-name-axe-idlerule-connector-2"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I send an event:
    """
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-connector-2",
      "connector_name": "test-connector-name-axe-idlerule-connector-2",
      "source_type": "component",
      "component":  "test-component-axe-idlerule-connector-2",
      "state": 1,
      "output": "test-output-axe-idlerule-connector-2",
      "long_output": "test-long-output-axe-idlerule-connector-2",
      "author": "test-author-axe-idlerule-connector-2"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-connector-axe-idlerule-connector-2/test-connector-name-axe-idlerule-connector-2"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-axe-idlerule-connector-2",
            "connector_name": "test-connector-name-axe-idlerule-connector-2",
            "pbehavior_info": {
              "canonical_type": "maintenance",
              "name": "test-pbehavior-axe-idlerule-connector-2",
              "reason": "Test Engine",
              "type": "test-maintenance-type-to-engine",
              "type_name": "Engine maintenance"
            },
            "steps": [
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
                "a": "system",
                "m": "Pbehavior test-pbehavior-axe-idlerule-connector-2. Type: Engine maintenance. Reason: Test Engine"
              }
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

  Scenario: given connector alarm should update service alarm
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """
    {
      "name": "test-idlerule-axe-idlerule-connector-3-name",
      "author": "test-idlerule-axe-idlerule-connector-3-author",
      "type": "entity",
      "enabled": true,
      "priority": 60,
      "duration": {
        "seconds": 3,
        "unit": "s"
      },
      "entity_patterns": [
        {
          "_id": "test-connector-axe-idlerule-connector-3/test-connector-name-axe-idlerule-connector-3"
        }
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    When I do POST /api/v4/entityservices:
    """
    {
      "name": "test-service-axe-idlerule-connector-3",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
      "enabled": true,
      "impact_level": 1,
      "entity_patterns": [
        {
          "_id": "test-connector-axe-idlerule-connector-3/test-connector-name-axe-idlerule-connector-3"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I wait the next periodical process
    When I send an event:
    """
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-connector-3",
      "connector_name": "test-connector-name-axe-idlerule-connector-3",
      "source_type": "component",
      "component":  "test-component-axe-idlerule-connector-3",
      "state": 1,
      "output": "test-output-axe-idlerule-connector-3",
      "long_output": "test-long-output-axe-idlerule-connector-3",
      "author": "test-author-axe-idlerule-connector-3"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-connector-axe-idlerule-connector-3/test-connector-name-axe-idlerule-connector-3"}]}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "entity": {
            "impact": ["{{ .serviceID }}"]
          },
          "v": {
            "connector": "test-connector-axe-idlerule-connector-3",
            "connector_name": "test-connector-name-axe-idlerule-connector-3"
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
    When I do GET /api/v4/alarms?filter={"$and":[{"entity._id":"{{ .serviceID }}"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "entity": {
            "depends": ["test-connector-axe-idlerule-connector-3/test-connector-name-axe-idlerule-connector-3"]
          },
          "v": {
            "component": "{{ .serviceID }}",
            "connector": "service",
            "connector_name": "service",
            "state": {
              "val": 3
            },
            "status": {
              "val": 1
            },
            "steps": [
              {
                "_t": "stateinc",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "service.service",
                "m": "All: 1; Alarms: 1; Acknowledged: 0; NotAcknowledged: 1; StateCritical: 1; StateMajor: 0; StateMinor: 0; StateInfo: 0; Pbehaviors: map[];",
                "val": 1
              }
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

  Scenario: given connector alarm should apply scenario on it
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """
    {
      "name": "test-idlerule-axe-idlerule-connector-4-name",
      "author": "test-idlerule-axe-idlerule-connector-4-author",
      "type": "entity",
      "enabled": true,
      "priority": 60,
      "duration": {
        "seconds": 3,
        "unit": "s"
      },
      "entity_patterns": [
        {
          "_id": "test-connector-axe-idlerule-connector-4/test-connector-name-axe-idlerule-connector-4"
        }
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    When I do POST /api/v4/scenarios:
    """
    {
      "name": "test-scenario-axe-idlerule-connector-4-name",
      "author": "test-scenario-axe-idlerule-connector-4-author",
      "enabled": true,
      "priority": 60,
      "triggers": ["create"],
      "actions": [
        {
          "entity_patterns": [
            {
              "_id": "test-connector-axe-idlerule-connector-4/test-connector-name-axe-idlerule-connector-4"
            }
          ],
          "type": "assocticket",
          "parameters": {
            "author": "test-author-axe-idlerule-connector-4",
            "output": "test-output-axe-idlerule-connector-4",
            "ticket": "test-ticket-axe-idlerule-connector-4"
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        },
        {
          "entity_patterns": [
            {
              "_id": "test-connector-axe-idlerule-connector-4/test-connector-name-axe-idlerule-connector-4"
            }
          ],
          "type": "ack",
          "parameters": {
            "author": "test-author-axe-idlerule-connector-4",
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
    """
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-connector-4",
      "connector_name": "test-connector-name-axe-idlerule-connector-4",
      "source_type": "component",
      "component":  "test-component-axe-idlerule-connector-4",
      "state": 1,
      "output": "test-output-axe-idlerule-connector-4",
      "long_output": "test-long-output-axe-idlerule-connector-4",
      "author": "test-author-axe-idlerule-connector-4"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-connector-axe-idlerule-connector-4/test-connector-name-axe-idlerule-connector-4"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-axe-idlerule-connector-4",
            "connector_name": "test-connector-name-axe-idlerule-connector-4",
            "steps": [
              {
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "assocticket",
                "a": "test-author-axe-idlerule-connector-4",
                "m": "test-ticket-axe-idlerule-connector-4"
              },
              {
                "_t": "ack",
                "a": "test-author-axe-idlerule-connector-4",
                "m": "test-output-axe-idlerule-connector-4"
              }
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

  Scenario: given connector alarm should apply dynamic infos on it
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """
    {
      "name": "test-idlerule-axe-idlerule-connector-5-name",
      "author": "test-idlerule-axe-idlerule-connector-5-author",
      "type": "entity",
      "enabled": true,
      "priority": 60,
      "duration": {
        "seconds": 3,
        "unit": "s"
      },
      "entity_patterns": [
        {
          "_id": "test-connector-axe-idlerule-connector-5/test-connector-name-axe-idlerule-connector-5"
        }
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/dynamic-infos:
    """
    {
      "name": "test-dynamicinfos-axe-idlerule-connector-5-name",
      "description": "test-dynamicinfos-axe-idlerule-connector-5-description",
      "author": "test-dynamicinfos-axe-idlerule-connector-5-author",
      "disable_during_periods": [],
      "entity_patterns": [
        {
          "_id": "test-connector-axe-idlerule-connector-5/test-connector-name-axe-idlerule-connector-5"
        }
      ],
      "infos": [
        {"name":"test-info-axe-idlerule-connector-5-name", "value":"test-info-axe-idlerule-connector-5-value"}
      ]
    }
    """
    Then the response code should be 201
    When I save response dynamicInfosRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-connector-5",
      "connector_name": "test-connector-name-axe-idlerule-connector-5",
      "source_type": "component",
      "component":  "test-component-axe-idlerule-connector-5",
      "state": 2,
      "output": "test-output-axe-idlerule-connector-5",
      "long_output": "test-long-output-axe-idlerule-connector-5",
      "author": "test-author-axe-idlerule-connector-5"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-connector-axe-idlerule-connector-5/test-connector-name-axe-idlerule-connector-5"}]}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-axe-idlerule-connector-5",
            "connector_name": "test-connector-name-axe-idlerule-connector-5",
            "infos": {
              "{{ .dynamicInfosRuleID }}": {
                "test-info-axe-idlerule-connector-5-name": "test-info-axe-idlerule-connector-5-value"
              }
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

  Scenario: given connector alarm should not create meta alarm
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """
    {
      "name": "test-idlerule-axe-idlerule-connector-6-name",
      "author": "test-idlerule-axe-idlerule-connector-6-author",
      "type": "entity",
      "enabled": true,
      "priority": 44,
      "duration": {
        "seconds": 3,
        "unit": "s"
      },
      "entity_patterns": [
        {
          "_id": "test-connector-axe-idlerule-connector-6/test-connector-name-axe-idlerule-connector-6"
        }
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-metaalarmrule-action-correlation-1",
      "type": "attribute",
      "config": {
        "alarm_patterns": [
          {
            "v": {
              "connector": "test-connector-axe-idlerule-connector-6"
            }
          }
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
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-connector-6",
      "connector_name": "test-connector-name-axe-idlerule-connector-6",
      "source_type": "component",
      "component":  "test-component-axe-idlerule-connector-6",
      "state": 1,
      "output": "test-output-axe-idlerule-connector-6",
      "long_output": "test-long-output-axe-idlerule-connector-6",
      "author": "test-author-axe-idlerule-connector-6"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true
    Then the response code should be 200
    When I save response metalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-component-axe-idlerule-connector-6"}]}
    Then the response code should be 200
    Then the response body should contain:
    """
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
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-connector-axe-idlerule-connector-6/test-connector-name-axe-idlerule-connector-6"}]}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-axe-idlerule-connector-6",
            "connector_name": "test-connector-name-axe-idlerule-connector-6",
            "parents": []
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

  Scenario: given connector alarm should update alarm on event
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """
    {
      "name": "test-idlerule-axe-idlerule-connector-7-name",
      "author": "test-idlerule-axe-idlerule-connector-7-author",
      "type": "entity",
      "enabled": true,
      "priority": 44,
      "duration": {
        "seconds": 3,
        "unit": "s"
      },
      "entity_patterns": [
        {
          "_id": "test-connector-axe-idlerule-connector-7/test-connector-name-axe-idlerule-connector-7"
        }
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-connector-7",
      "connector_name": "test-connector-name-axe-idlerule-connector-7",
      "source_type": "component",
      "component":  "test-component-axe-idlerule-connector-7",
      "state": 1,
      "output": "test-output-axe-idlerule-connector-7",
      "long_output": "test-long-output-axe-idlerule-connector-7",
      "author": "test-author-axe-idlerule-connector-7"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "event_type": "ack",
      "connector": "test-connector-axe-idlerule-connector-7",
      "connector_name": "test-connector-name-axe-idlerule-connector-7",
      "source_type": "connector",
      "output": "test-output-axe-idlerule-connector-7",
      "author": "test-author-axe-idlerule-connector-7"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-connector-axe-idlerule-connector-7/test-connector-name-axe-idlerule-connector-7"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-axe-idlerule-connector-7",
            "connector_name": "test-connector-name-axe-idlerule-connector-7",
            "steps": [
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