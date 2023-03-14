Feature: update meta alarm on action
  I need to be able to update meta alarm on action

  @concurrent
  Scenario: given meta alarm and scenario should update meta alarm and update children
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-action-correlation-1",
      "type": "attribute",
      "auto_resolve": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-action-correlation-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-action-correlation-1",
      "connector": "test-connector-action-correlation-1",
      "connector_name": "test-connector-name-action-correlation-1",
      "component":  "test-component-action-correlation-1",
      "resource": "test-resource-action-correlation-1",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-action-correlation-1&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 1,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-action-correlation-1"
          }
        }
      ],
      "meta": {
        "total_count": 1
      }
    }
    """
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I save response metaalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I do POST /api/v4/scenarios:
    """json
    {
      "_id": "test-scenario-action-correlation-1",
      "name": "test-scenario-action-correlation-1-name",
      "priority": 10045,
      "enabled": true,
      "triggers": ["comment"],
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "{{ .metaAlarmResource }}"
                }
              }
            ]
          ],
          "type": "assocticket",
          "comment": "test-scenario-action-correlation-1-action-1-comment",
          "parameters": {
            "ticket": "test-ticket-action-correlation-1",
            "ticket_system_name": "test-system-name-action-correlation-1",
            "ticket_url": "test-ticket-url-action-correlation-1",
            "ticket_data": {
              "ticket_param_1": "ticket_value_1"
            }
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
                  "value": "{{ .metaAlarmResource }}"
                }
              }
            ]
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
      "event_type": "comment",
      "output": "test-output-action-correlation-1",
      "connector": "{{ .metaAlarmConnector }}",
      "connector_name": "{{ .metaAlarmConnectorName }}",
      "component":  "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "source_type": "resource"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "comment",
        "connector": "{{ .metaAlarmConnector }}",
        "connector_name": "{{ .metaAlarmConnectorName }}",
        "component":  "{{ .metaAlarmComponent }}",
        "resource": "{{ .metaAlarmResource }}",
        "source_type": "resource"
      },
      {
        "event_type": "comment",
        "connector": "test-connector-action-correlation-1",
        "connector_name": "test-connector-name-action-correlation-1",
        "component":  "test-component-action-correlation-1",
        "resource": "test-resource-action-correlation-1",
        "source_type": "resource"
      },
      {
        "event_type": "assocticket",
        "connector": "test-connector-action-correlation-1",
        "connector_name": "test-connector-name-action-correlation-1",
        "component":  "test-component-action-correlation-1",
        "resource": "test-resource-action-correlation-1",
        "source_type": "resource"
      },
      {
        "event_type": "ack",
        "connector": "test-connector-action-correlation-1",
        "connector_name": "test-connector-name-action-correlation-1",
        "component":  "test-component-action-correlation-1",
        "resource": "test-resource-action-correlation-1",
        "source_type": "resource"
      }
    ]
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID }}",
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
                  "component": "test-component-action-correlation-1",
                  "connector": "test-connector-action-correlation-1",
                  "connector_name": "test-connector-name-action-correlation-1",
                  "resource": "test-resource-action-correlation-1"
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ (index (index .lastResponse 0).data.children.data 0)._id }}",
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
                "initiator": "system",
                "m": "Scenario: test-scenario-action-correlation-1-name. Ticket ID: test-ticket-action-correlation-1. Ticket URL: test-ticket-url-action-correlation-1. Ticket ticket_param_1: ticket_value_1.",
                "ticket": "test-ticket-action-correlation-1",
                "ticket_rule_id": "test-scenario-action-correlation-1",
                "ticket_rule_name": "Scenario: test-scenario-action-correlation-1-name",
                "ticket_system_name": "test-system-name-action-correlation-1",
                "ticket_url": "test-ticket-url-action-correlation-1",
                "ticket_data": {
                  "ticket_param_1": "ticket_value_1"
                },
                "ticket_comment": "test-scenario-action-correlation-1-action-1-comment"
              },
              {
                "_t": "ack",
                "a": "system",
                "m": "test-output-action-correlation-1-engine"
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 5
            }
          }
        }
      },
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
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
                "initiator": "system",
                "m": "Scenario: test-scenario-action-correlation-1-name. Ticket ID: test-ticket-action-correlation-1. Ticket URL: test-ticket-url-action-correlation-1. Ticket ticket_param_1: ticket_value_1.",
                "ticket": "test-ticket-action-correlation-1",
                "ticket_rule_id": "test-scenario-action-correlation-1",
                "ticket_rule_name": "Scenario: test-scenario-action-correlation-1-name",
                "ticket_system_name": "test-system-name-action-correlation-1",
                "ticket_url": "test-ticket-url-action-correlation-1",
                "ticket_data": {
                  "ticket_param_1": "ticket_value_1"
                },
                "ticket_comment": "test-scenario-action-correlation-1-action-1-comment"
              },
              {
                "_t": "ack",
                "a": "system",
                "m": "test-output-action-correlation-1-engine"
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 6
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given meta alarm and scenario with webhook action should update meta alarm and update children
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-action-correlation-2",
      "type": "attribute",
      "auto_resolve": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-action-correlation-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-action-correlation-2",
      "connector": "test-connector-action-correlation-2",
      "connector_name": "test-connector-name-action-correlation-2",
      "component":  "test-component-action-correlation-2",
      "resource": "test-resource-action-correlation-2-1",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 1,
      "output": "test-output-action-correlation-2",
      "connector": "test-connector-action-correlation-2",
      "connector_name": "test-connector-name-action-correlation-2",
      "component":  "test-component-action-correlation-2",
      "resource": "test-resource-action-correlation-2-2",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-action-correlation-2&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-action-correlation-2"
          }
        }
      ],
      "meta": {
        "total_count": 1
      }
    }
    """
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I save response metaalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I save response metaalarmDisplayName={{ (index .lastResponse.data 0).v.display_name }}
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I do POST /api/v4/scenarios:
    """json
    {
      "_id": "test-scenario-action-correlation-2",
      "name": "test-scenario-action-correlation-2-name",
      "priority": 10046,
      "enabled": true,
      "triggers": ["comment"],
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "{{ .metaAlarmResource }}"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request",
              "payload": "{\"ticket_id\":\"testticket\",\"ticket_data\":\"testdata\"}"
            },
            "declare_ticket": {
              "ticket_id": "ticket_id",
              "ticket_data": "ticket_data"
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
      "event_type": "comment",
      "output": "test-output-action-correlation-2",
      "connector": "{{ .metaAlarmConnector }}",
      "connector_name": "{{ .metaAlarmConnectorName }}",
      "component":  "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "source_type": "resource"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "comment",
        "connector": "{{ .metaAlarmConnector }}",
        "connector_name": "{{ .metaAlarmConnectorName }}",
        "component":  "{{ .metaAlarmComponent }}",
        "resource": "{{ .metaAlarmResource }}",
        "source_type": "resource"
      },
      {
        "event_type": "comment",
        "connector": "test-connector-action-correlation-2",
        "connector_name": "test-connector-name-action-correlation-2",
        "component":  "test-component-action-correlation-2",
        "resource": "test-resource-action-correlation-2-1",
        "source_type": "resource"
      },
      {
        "event_type": "comment",
        "connector": "test-connector-action-correlation-2",
        "connector_name": "test-connector-name-action-correlation-2",
        "component":  "test-component-action-correlation-2",
        "resource": "test-resource-action-correlation-2-2",
        "source_type": "resource"
      },
      {
        "event_type": "declareticketwebhook",
        "connector": "test-connector-action-correlation-2",
        "connector_name": "test-connector-name-action-correlation-2",
        "component":  "test-component-action-correlation-2",
        "resource": "test-resource-action-correlation-2-1",
        "source_type": "resource"
      },
      {
        "event_type": "declareticketwebhook",
        "connector": "test-connector-action-correlation-2",
        "connector_name": "test-connector-name-action-correlation-2",
        "component":  "test-component-action-correlation-2",
        "resource": "test-resource-action-correlation-2-2",
        "source_type": "resource"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-action-correlation-2",
      "connector": "test-connector-action-correlation-2",
      "connector_name": "test-connector-name-action-correlation-2",
      "component":  "test-component-action-correlation-2",
      "resource": "test-resource-action-correlation-2-3",
      "source_type": "resource"
    }
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "declareticketwebhook",
      "connector": "test-connector-action-correlation-2",
      "connector_name": "test-connector-name-action-correlation-2",
      "component":  "test-component-action-correlation-2",
      "resource": "test-resource-action-correlation-2-3",
      "source_type": "resource"
    }
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID }}",
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
                  "component": "test-component-action-correlation-2",
                  "connector": "test-connector-action-correlation-2",
                  "connector_name": "test-connector-name-action-correlation-2",
                  "resource": "test-resource-action-correlation-2-1",
                  "tickets": [
                    {
                      "_t": "declareticket",
                      "ticket": "testticket",
                      "ticket_rule_id": "test-scenario-action-correlation-2",
                      "ticket_rule_name": "Scenario: test-scenario-action-correlation-2-name",
                      "ticket_meta_alarm_id": "{{ .metaAlarmID }}",
                      "ticket_data": {
                        "ticket_data": "testdata"
                      }
                    }
                  ],
                  "ticket": {
                    "_t": "declareticket",
                    "ticket": "testticket",
                    "ticket_rule_id": "test-scenario-action-correlation-2",
                    "ticket_rule_name": "Scenario: test-scenario-action-correlation-2-name",
                    "ticket_meta_alarm_id": "{{ .metaAlarmID }}",
                    "ticket_data": {
                      "ticket_data": "testdata"
                    }
                  }
                }
              },
              {
                "v": {
                  "component": "test-component-action-correlation-2",
                  "connector": "test-connector-action-correlation-2",
                  "connector_name": "test-connector-name-action-correlation-2",
                  "resource": "test-resource-action-correlation-2-2",
                  "tickets": [
                    {
                      "_t": "declareticket",
                      "ticket": "testticket",
                      "ticket_rule_id": "test-scenario-action-correlation-2",
                      "ticket_rule_name": "Scenario: test-scenario-action-correlation-2-name",
                      "ticket_meta_alarm_id": "{{ .metaAlarmID }}",
                      "ticket_data": {
                        "ticket_data": "testdata"
                      }
                    }
                  ],
                  "ticket": {
                    "_t": "declareticket",
                    "ticket": "testticket",
                    "ticket_rule_id": "test-scenario-action-correlation-2",
                    "ticket_rule_name": "Scenario: test-scenario-action-correlation-2-name",
                    "ticket_meta_alarm_id": "{{ .metaAlarmID }}",
                    "ticket_data": {
                      "ticket_data": "testdata"
                    }
                  }
                }
              },
              {
                "v": {
                  "component": "test-component-action-correlation-2",
                  "connector": "test-connector-action-correlation-2",
                  "connector_name": "test-connector-name-action-correlation-2",
                  "resource": "test-resource-action-correlation-2-3",
                  "tickets": [
                    {
                      "_t": "declareticket",
                      "ticket": "testticket",
                      "ticket_rule_id": "test-scenario-action-correlation-2",
                      "ticket_rule_name": "Scenario: test-scenario-action-correlation-2-name",
                      "ticket_meta_alarm_id": "{{ .metaAlarmID }}",
                      "ticket_data": {
                        "ticket_data": "testdata"
                      }
                    }
                  ],
                  "ticket": {
                    "_t": "declareticket",
                    "ticket": "testticket",
                    "ticket_rule_id": "test-scenario-action-correlation-2",
                    "ticket_rule_name": "Scenario: test-scenario-action-correlation-2-name",
                    "ticket_meta_alarm_id": "{{ .metaAlarmID }}",
                    "ticket_data": {
                      "ticket_data": "testdata"
                    }
                  }
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ (index (index .lastResponse 0).data.children.data 0)._id }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ (index (index .lastResponse 0).data.children.data 1)._id }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ (index (index .lastResponse 0).data.children.data 2)._id }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response array key "0.data.steps.data" should contain only:
    """json
    [
      {
        "_t": "stateinc"
      },
      {
        "_t": "statusinc"
      },
      {
        "_t": "comment"
      },
      {
        "_t": "webhookstart",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-correlation-2-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-correlation-2-name"
      },
      {
        "_t": "declareticket",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-correlation-2-name. Ticket ID: testticket. Ticket ticket_data: testdata."
      }
    ]
    """
    Then the response body should contain:
    """json
    [
      {
        "status": 200
      },
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
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
                "_t": "declareticket",
                "a": "system",
                "user_id": "",
                "m": "Scenario: test-scenario-action-correlation-2-name. Ticket ID: testticket. Ticket ticket_data: testdata."
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 5
            }
          }
        }
      },
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
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
                "_t": "declareticket",
                "a": "system",
                "user_id": "",
                "m": "Scenario: test-scenario-action-correlation-2-name. Ticket ID: testticket. Ticket ticket_data: testdata."
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 5
            }
          }
        }
      },
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
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
                "_t": "declareticket",
                "a": "system",
                "user_id": "",
                "m": "Scenario: test-scenario-action-correlation-2-name. Ticket ID: testticket. Ticket ticket_data: testdata."
              },
              {
                "_t": "comment"
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 5
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given meta alarm and scenario for child alarm should update meta alarm state
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-action-correlation-4",
      "type": "attribute",
      "auto_resolve": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-action-correlation-4"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-action-correlation-4",
      "connector": "test-connector-action-correlation-4",
      "connector_name": "test-connector-name-action-correlation-4",
      "component":  "test-component-action-correlation-4",
      "resource": "test-resource-action-correlation-4",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-action-correlation-4&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 1,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-action-correlation-4"
          }
        }
      ],
      "meta": {
        "total_count": 1
      }
    }
    """
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-correlation-4-name",
      "priority": 10049,
      "enabled": true,
      "triggers": ["comment"],
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-action-correlation-4"
                }
              }
            ]
          ],
          "type": "changestate",
          "parameters": {
            "state": 3
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "comment",
      "output": "test-output-action-correlation-4",
      "connector": "test-connector-action-correlation-4",
      "connector_name": "test-connector-name-action-correlation-4",
      "component":  "test-component-action-correlation-4",
      "resource": "test-resource-action-correlation-4",
      "source_type": "resource"
    }
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID }}",
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
                  "component": "test-component-action-correlation-4",
                  "connector": "test-connector-action-correlation-4",
                  "connector_name": "test-connector-name-action-correlation-4",
                  "resource": "test-resource-action-correlation-4"
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ (index (index .lastResponse 0).data.children.data 0)._id }}",
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
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "stateinc",
                "val": 3
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
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "metaalarmattach"
              },
              {
                "_t": "comment",
                "m": "test-output-action-correlation-4"
              },
              {
                "_t": "changestate",
                "val": 3
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 5
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given meta alarm and scenario with webhook action should trigger both for parent and child
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-action-correlation-5",
      "type": "attribute",
      "auto_resolve": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-action-correlation-5"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-action-correlation-5",
      "connector": "test-connector-action-correlation-5",
      "connector_name": "test-connector-name-action-correlation-5",
      "component":  "test-component-action-correlation-5",
      "resource": "test-resource-action-correlation-5",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-action-correlation-5
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do GET /api/v4/alarms?search=test-resource-action-correlation-5&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 1,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-action-correlation-5"
          }
        }
      ],
      "meta": {
        "total_count": 1
      }
    }
    """
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-correlation-5-name",
      "enabled": true,
      "triggers": ["comment"],
      "priority": 10016,
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "{{ .metaAlarmResource }}"
                }
              }
            ],
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-action-correlation-5"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request",
              "payload": "{\"ticket_id\":\"testticket\",\"ticket_data\":\"testdata\"}"
            },
            "declare_ticket": {
              "ticket_id": "ticket_id",
              "ticket_data": "ticket_data"
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
      "event_type": "comment",
      "output": "test-output-action-correlation-5",
      "connector": "{{ .metaAlarmConnector }}",
      "connector_name": "{{ .metaAlarmConnectorName }}",
      "component":  "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "source_type": "resource"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "comment",
        "connector": "{{ .metaAlarmConnector }}",
        "connector_name": "{{ .metaAlarmConnectorName }}",
        "component":  "{{ .metaAlarmComponent }}",
        "resource": "{{ .metaAlarmResource }}",
        "source_type": "resource"
      },
      {
        "event_type": "comment",
        "connector": "test-connector-action-correlation-5",
        "connector_name": "test-connector-name-action-correlation-5",
        "component":  "test-component-action-correlation-5",
        "resource": "test-resource-action-correlation-5",
        "source_type": "resource"
      },
      {
        "event_type": "declareticketwebhook",
        "connector": "test-connector-action-correlation-5",
        "connector_name": "test-connector-name-action-correlation-5",
        "component":  "test-component-action-correlation-5",
        "resource": "test-resource-action-correlation-5",
        "source_type": "resource"
      }
    ]
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .alarmID }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response array key "0.data.steps.data" should contain only:
    """json
    [
      {
        "_t": "stateinc"
      },
      {
        "_t": "statusinc"
      },
      {
        "_t": "metaalarmattach"
      },
      {
        "_t": "comment"
      },
      {
        "_t": "webhookstart",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-correlation-5-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-correlation-5-name"
      },
      {
        "_t": "declareticket",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-correlation-5-name. Ticket ID: testticket. Ticket ticket_data: testdata."
      },
      {
        "_t": "declareticket",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-correlation-5-name. Ticket ID: testticket. Ticket ticket_data: testdata."
      }
    ]
    """

  @concurrent
  Scenario: given meta alarm and scenario with webhook action should trigger only for parent
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-action-correlation-6",
      "type": "attribute",
      "auto_resolve": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-action-correlation-6"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-action-correlation-6",
      "connector": "test-connector-action-correlation-6",
      "connector_name": "test-connector-name-action-correlation-6",
      "component":  "test-component-action-correlation-6",
      "resource": "test-resource-action-correlation-6",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-action-correlation-6
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do GET /api/v4/alarms?search=test-resource-action-correlation-6&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 1,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-action-correlation-6"
          }
        }
      ],
      "meta": {
        "total_count": 1
      }
    }
    """
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-correlation-6-name",
      "enabled": true,
      "triggers": ["comment"],
      "priority": 10101,
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "{{ .metaAlarmResource }}"
                }
              }
            ],
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-action-correlation-6"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "skip_for_child": true,
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request",
              "payload": "{\"ticket_id\":\"testticket\",\"ticket_data\":\"testdata\"}"
            },
            "declare_ticket": {
              "ticket_id": "ticket_id",
              "ticket_data": "ticket_data"
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
      "event_type": "comment",
      "output": "test-output-action-correlation-6",
      "connector": "{{ .metaAlarmConnector }}",
      "connector_name": "{{ .metaAlarmConnectorName }}",
      "component":  "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "source_type": "resource"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "comment",
        "connector": "{{ .metaAlarmConnector }}",
        "connector_name": "{{ .metaAlarmConnectorName }}",
        "component":  "{{ .metaAlarmComponent }}",
        "resource": "{{ .metaAlarmResource }}",
        "source_type": "resource"
      },
      {
        "event_type": "comment",
        "connector": "test-connector-action-correlation-6",
        "connector_name": "test-connector-name-action-correlation-6",
        "component":  "test-component-action-correlation-6",
        "resource": "test-resource-action-correlation-6",
        "source_type": "resource"
      },
      {
        "event_type": "declareticketwebhook",
        "connector": "test-connector-action-correlation-6",
        "connector_name": "test-connector-name-action-correlation-6",
        "component":  "test-component-action-correlation-6",
        "resource": "test-resource-action-correlation-6",
        "source_type": "resource"
      }
    ]
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .alarmID }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response array key "0.data.steps.data" should contain only:
    """json
    [
      {
        "_t": "stateinc"
      },
      {
        "_t": "statusinc"
      },
      {
        "_t": "metaalarmattach"
      },
      {
        "_t": "comment"
      },
      {
        "_t": "declareticket",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-correlation-6-name. Ticket ID: testticket. Ticket ticket_data: testdata."
      }
    ]
    """

  @concurrent
  Scenario: given meta alarm and scenario with failed webhook action should update only meta alarm
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-action-correlation-7",
      "type": "attribute",
      "auto_resolve": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-action-correlation-7"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-action-correlation-7",
      "connector": "test-connector-action-correlation-7",
      "connector_name": "test-connector-name-action-correlation-7",
      "component":  "test-component-action-correlation-7",
      "resource": "test-resource-action-correlation-7-1",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-action-correlation-7&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 1,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-action-correlation-7"
          }
        }
      ],
      "meta": {
        "total_count": 1
      }
    }
    """
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I save response metaalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I save response metaalarmDisplayName={{ (index .lastResponse.data 0).v.display_name }}
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I do POST /api/v4/scenarios:
    """json
    {
      "_id": "test-scenario-action-correlation-7",
      "name": "test-scenario-action-correlation-7-name",
      "priority": 10013,
      "enabled": true,
      "triggers": ["comment"],
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "{{ .metaAlarmResource }}"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/auth-request",
              "payload": "{\"_id\":\"testticket\"}"
            },
            "declare_ticket": {
              "ticket_id": "_id"
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
      "event_type": "comment",
      "output": "test-output-action-correlation-7",
      "connector": "{{ .metaAlarmConnector }}",
      "connector_name": "{{ .metaAlarmConnectorName }}",
      "component":  "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "source_type": "resource"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "comment",
        "connector": "{{ .metaAlarmConnector }}",
        "connector_name": "{{ .metaAlarmConnectorName }}",
        "component":  "{{ .metaAlarmComponent }}",
        "resource": "{{ .metaAlarmResource }}",
        "source_type": "resource"
      },
      {
        "event_type": "comment",
        "connector": "test-connector-action-correlation-7",
        "connector_name": "test-connector-name-action-correlation-7",
        "component":  "test-component-action-correlation-7",
        "resource": "test-resource-action-correlation-7-1",
        "source_type": "resource"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-action-correlation-7",
      "connector": "test-connector-action-correlation-7",
      "connector_name": "test-connector-name-action-correlation-7",
      "component":  "test-component-action-correlation-7",
      "resource": "test-resource-action-correlation-7-2",
      "source_type": "resource"
    }
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID }}",
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
                  "component": "test-component-action-correlation-7",
                  "connector": "test-connector-action-correlation-7",
                  "connector_name": "test-connector-name-action-correlation-7",
                  "resource": "test-resource-action-correlation-7-1"
                }
              },
              {
                "v": {
                  "component": "test-component-action-correlation-7",
                  "connector": "test-connector-action-correlation-7",
                  "connector_name": "test-connector-name-action-correlation-7",
                  "resource": "test-resource-action-correlation-7-2"
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
    Then the response key "0.data.children.data.0.v.ticket" should not exist
    Then the response key "0.data.children.data.0.v.tickets" should not exist
    Then the response key "0.data.children.data.1.v.ticket" should not exist
    Then the response key "0.data.children.data.1.v.tickets" should not exist
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ (index (index .lastResponse 0).data.children.data 0)._id }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ (index (index .lastResponse 0).data.children.data 1)._id }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response array key "0.data.steps.data" should contain only:
    """json
    [
      {
        "_t": "stateinc"
      },
      {
        "_t": "statusinc"
      },
      {
        "_t": "comment"
      },
      {
        "_t": "webhookstart",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-correlation-7-name"
      },
      {
        "_t": "webhookfail",
        "a": "system",
        "user_id": ""
      },
      {
        "_t": "declareticketfail",
        "a": "system",
        "user_id": ""
      }
    ]
    """
    Then the response body should contain:
    """json
    [
      {
        "status": 200
      },
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
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
      },
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
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
