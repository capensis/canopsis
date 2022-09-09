Feature: update meta alarm on action
  I need to be able to update meta alarm on action

  Scenario: given meta alarm and scenario should update meta alarm and update children
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-action-correlation-1",
      "type": "attribute",
      "config": {
        "alarm_patterns": [
          {
            "v": {
              "component": "test-component-action-correlation-1"
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
    """json
    {
      "connector": "test-connector-action-correlation-1",
      "connector_name": "test-connector-name-action-correlation-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-action-correlation-1",
      "resource": "test-resource-action-correlation-1",
      "state": 2,
      "output": "test-output-action-correlation-1"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true
    Then the response code should be 200
    When I save response metalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-correlation-1-name",
      "enabled": true,
      "priority": 81,
      "triggers": ["comment"],
      "actions": [
        {
          "entity_patterns": [
            {
              "_id": "{{ .metalarmEntityID }}"
            }
          ],
          "type": "assocticket",
          "parameters": {
            "output": "test-output-action-correlation-1-{{ `{{ .Alarm.Value.Connector }}` }}",
            "ticket": "test-ticket-action-correlation-1"
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        },
        {
          "entity_patterns": [
            {
              "_id": "{{ .metalarmEntityID }}"
            }
          ],
          "type": "ack",
          "parameters": {
            "output": "test-output-action-correlation-1-{{ `{{ .Alarm.Value.Connector }}` }}"
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
      "connector": "{{ .metaAlarmConnector }}",
      "connector_name": "{{ .metaAlarmConnectorName }}",
      "source_type": "resource",
      "event_type": "comment",
      "component":  "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "output": "test-output-action-correlation-1"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "consequences": {
            "data": [
              {
                "v": {
                  "component": "test-component-action-correlation-1",
                  "connector": "test-connector-action-correlation-1",
                  "connector_name": "test-connector-name-action-correlation-1",
                  "resource": "test-resource-action-correlation-1",
                  "steps": [
                    {
                      "_t": "stateinc",
                      "val": 2
                    },
                    {
                      "_t": "statusinc",
                      "val": 1
                    },
                    {
                      "_t": "metaalarmattach",
                      "val": 0
                    },
                    {
                      "_t": "comment"
                    },
                    {
                      "_t": "assocticket",
                      "a": "system",
                      "m": "test-ticket-action-correlation-1"
                    },
                    {
                      "_t": "ack",
                      "a": "system",
                      "m": "test-output-action-correlation-1-engine"
                    }
                  ]
                }
              }
            ],
            "total": 1
          },
          "metaalarm": true,
          "rule": {
            "name": "test-metaalarmrule-action-correlation-1"
          },
          "v": {
            "children": [
              "test-resource-action-correlation-1/test-component-action-correlation-1"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID }}",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "steps": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "comment"
              },
              {
                "_t": "assocticket",
                "a": "system",
                "m": "test-ticket-action-correlation-1"
              },
              {
                "_t": "ack",
                "a": "system",
                "m": "test-output-action-correlation-1-engine"
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

  Scenario: given meta alarm and scenario with webhook action should update meta alarm and update children
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-action-correlation-2",
      "type": "attribute",
      "config": {
        "alarm_patterns": [
          {
            "v": {
              "component": "test-component-action-correlation-2"
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
    """json
    {
      "connector": "test-connector-action-correlation-2",
      "connector_name": "test-connector-name-action-correlation-2",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-action-correlation-2",
      "resource": "test-resource-action-correlation-2-1",
      "state": 2,
      "output": "test-output-action-correlation-2"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-action-correlation-2",
      "connector_name": "test-connector-name-action-correlation-2",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-action-correlation-2",
      "resource": "test-resource-action-correlation-2-2",
      "state": 1,
      "output": "test-output-action-correlation-2"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true
    Then the response code should be 200
    When I save response metalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I save response metalarmDisplayName={{ (index .lastResponse.data 0).v.display_name }}
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-correlation-2-name",
      "enabled": true,
      "priority": 82,
      "triggers": ["comment"],
      "actions": [
        {
          "entity_patterns": [
            {
              "_id": "{{ .metalarmEntityID }}"
            }
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .apiURL }}/api/v4/scenarios",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {"Content-Type": "application/json"},
              "payload": "{\"name\":\"{{ `Alarm: {{ .Alarm.Value.DisplayName }}; Entity: {{ .Entity.ID }}; Children: {{ range $children := .Children }}{{ $children.EntityID }}, {{ end }}` }}\",\"enabled\":true,\"priority\":83,\"triggers\":[\"create\"],\"actions\":[{\"alarm_patterns\":[{\"_id\":\"test-scenario-action-correlation-1-alarm\"}],\"type\":\"ack\",\"drop_scenario_if_not_matched\":false,\"emit_trigger\":false}]}"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "_id",
              "scenario_name": "name"
            }
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
      "connector": "{{ .metaAlarmConnector }}",
      "connector_name": "{{ .metaAlarmConnectorName }}",
      "source_type": "resource",
      "event_type": "comment",
      "component":  "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "output": "test-output-action-correlation-2"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "consequences": {
            "data": [
              {
                "v": {
                  "component": "test-component-action-correlation-2",
                  "connector": "test-connector-action-correlation-2",
                  "connector_name": "test-connector-name-action-correlation-2",
                  "ticket": {
                    "_t": "declareticket",
                    "data": {
                      "scenario_name": "Alarm: {{ .metalarmDisplayName }}; Entity: {{ .metaAlarmResource }}/metaalarm; Children: test-resource-action-correlation-2-1/test-component-action-correlation-2, test-resource-action-correlation-2-2/test-component-action-correlation-2, "
                    }
                  },
                  "steps": [
                    {
                      "_t": "stateinc"
                    },
                    {
                      "_t": "statusinc",
                      "val": 1
                    },
                    {
                      "_t": "metaalarmattach",
                      "val": 0
                    },
                    {
                      "_t": "comment"
                    },
                    {
                      "_t": "declareticket"
                    }
                  ]
                }
              },
              {
                "v": {
                  "component": "test-component-action-correlation-2",
                  "connector": "test-connector-action-correlation-2",
                  "connector_name": "test-connector-name-action-correlation-2",
                  "ticket": {
                    "_t": "declareticket",
                    "data": {
                      "scenario_name": "Alarm: {{ .metalarmDisplayName }}; Entity: {{ .metaAlarmResource }}/metaalarm; Children: test-resource-action-correlation-2-1/test-component-action-correlation-2, test-resource-action-correlation-2-2/test-component-action-correlation-2, "
                    }
                  },
                  "steps": [
                    {
                      "_t": "stateinc"
                    },
                    {
                      "_t": "statusinc",
                      "val": 1
                    },
                    {
                      "_t": "metaalarmattach",
                      "val": 0
                    },
                    {
                      "_t": "comment"
                    },
                    {
                      "_t": "declareticket"
                    }
                  ]
                }
              }
            ],
            "total": 2
          },
          "metaalarm": true,
          "rule": {
            "name": "test-metaalarmrule-action-correlation-2"
          },
          "v": {
            "children": [
              "test-resource-action-correlation-2-1/test-component-action-correlation-2",
              "test-resource-action-correlation-2-2/test-component-action-correlation-2"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID }}",
            "ticket": {
              "_t": "declareticket",
              "data": {
                "scenario_name": "Alarm: {{ .metalarmDisplayName }}; Entity: {{ .metaAlarmResource }}/metaalarm; Children: test-resource-action-correlation-2-1/test-component-action-correlation-2, test-resource-action-correlation-2-2/test-component-action-correlation-2, "
              }
            },
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "steps": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "comment"
              },
              {
                "_t": "declareticket"
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

  Scenario: given meta alarm and ack event should double ack children if they're already acked
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-metaalarmrule-action-correlation-8",
      "type": "attribute",
      "config": {
        "alarm_patterns": [
          {
            "v": {
              "component": "test-component-action-correlation-8"
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
      "connector": "test-connector-action-correlation-8",
      "connector_name": "test-connector-name-action-correlation-8",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-action-correlation-8",
      "resource": "test-resource-action-correlation-8-1",
      "state": 2,
      "output": "test-output-action-correlation-8",
      "long_output": "test-long-output-action-correlation-8",
      "author": "test-author-action-correlation-8"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-connector-action-correlation-8",
      "connector_name": "test-connector-name-action-correlation-8",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-action-correlation-8",
      "resource": "test-resource-action-correlation-8-2",
      "state": 2,
      "output": "test-output-action-correlation-8",
      "long_output": "test-long-output-action-correlation-8",
      "author": "test-author-action-correlation-8"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-connector-action-correlation-8",
      "connector_name": "test-connector-name-action-correlation-8",
      "source_type": "resource",
      "event_type": "ack",
      "component":  "test-component-action-correlation-8",
      "resource": "test-resource-action-correlation-8-2",
      "output": "previous ack",
      "author": "test-author-action-correlation-8"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true
    Then the response code should be 200
    When I save response metalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I send an event:
    """
    {
      "connector": "{{ .metaAlarmConnector }}",
      "connector_name": "{{ .metaAlarmConnectorName }}",
      "source_type": "resource",
      "event_type": "ack",
      "component":  "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "output": "metaalarm ack",
      "long_output": "test-long-output-action-correlation-8",
      "author": "test-author-action-correlation-8"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "metaalarm": true,
          "rule": {
            "name": "test-metaalarmrule-action-correlation-8"
          },
          "v": {
            "ack": {
              "_t": "ack",
              "a": "test-author-action-correlation-8",
              "m": "metaalarm ack",
              "val": 0
            },
            "children": [
              "test-resource-action-correlation-8-1/test-component-action-correlation-8",
              "test-resource-action-correlation-8-2/test-component-action-correlation-8"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID }}",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "steps": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "ack",
                "a": "test-author-action-correlation-8",
                "m": "metaalarm ack",
                "val": 0
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
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-action-correlation-8-1"}]}&with_steps=true&with_consequences=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "ack": {
              "_t": "ack",
              "a": "test-author-action-correlation-8",
              "m": "metaalarm ack",
              "val": 0
            },
            "children": [],
            "component": "test-component-action-correlation-8",
            "connector": "test-connector-action-correlation-8",
            "connector_name": "test-connector-name-action-correlation-8",
            "initial_long_output": "test-long-output-action-correlation-8",
            "initial_output": "test-output-action-correlation-8",
            "parents": [
              "{{ .metalarmEntityID }}"
            ],
            "resource": "test-resource-action-correlation-8-1",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "steps": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "metaalarmattach",
                "a": "engine.correlation",
                "val": 0
              },
              {
                "_t": "ack",
                "a": "test-author-action-correlation-8",
                "m": "metaalarm ack",
                "val": 0
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
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-action-correlation-8-2"}]}&with_steps=true&with_consequences=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "ack": {
              "_t": "ack",
              "a": "test-author-action-correlation-8",
              "m": "metaalarm ack",
              "val": 0
            },
            "children": [],
            "component": "test-component-action-correlation-8",
            "connector": "test-connector-action-correlation-8",
            "connector_name": "test-connector-name-action-correlation-8",
            "initial_long_output": "test-long-output-action-correlation-8",
            "initial_output": "test-output-action-correlation-8",
            "parents": [
              "{{ .metalarmEntityID }}"
            ],
            "resource": "test-resource-action-correlation-8-2",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "steps": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "metaalarmattach",
                "a": "engine.correlation",
                "val": 0
              },
              {
                "_t": "ack",
                "a": "test-author-action-correlation-8",
                "m": "previous ack",
                "val": 0
              },
              {
                "_t": "ack",
                "a": "test-author-action-correlation-8",
                "m": "metaalarm ack",
                "val": 0
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
