Feature: execute action on trigger
  I need to be able to trigger action on event

  Scenario: given scenario and check event should update alarm
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "_id": "test-scenario-action-1",
      "name": "test-scenario-action-1-name",
      "priority": 10002,
      "enabled": true,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-component-action-1"
                }
              }
            ]
          ],
          "type": "assocticket",
          "comment": "test-scenario-action-1-action-1-comment",
          "parameters": {
            "forward_author": false,
            "author": "test-scenario-action-1-action-1-author {{ `{{ .Alarm.Value.Resource }}` }}",
            "ticket": "test-scenario-action-1-action-1-ticket",
            "ticket_system_name": "test-scenario-action-1-action-1-system-name",
            "ticket_url": "test-scenario-action-1-action-1-ticket-url",
            "ticket_data": {
              "ticket_param_1": "ticket_value_1"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        },
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-component-action-1"
                }
              }
            ]
          ],
          "type": "ack",
          "parameters": {
            "forward_author": true,
            "author": "test-scenario-action-1-action-2-author {{ `{{ .Alarm.Value.Resource }}` }}",
            "output": "test-scenario-action-1-action-2-output {{ `{{ .Entity.Name }} {{ .Alarm.Value.State.Value }}` }}"
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        },
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-component-action-1"
                }
              }
            ]
          ],
          "type": "changestate",
          "parameters": {
            "state": 3,
            "forward_author": false,
            "author": "",
            "output": "test-scenario-action-1-action-3-output {{ `{{ .Env.Name }} {{ .Entity.Name }} {{ .Alarm.Value.State.Value }}` }}"
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
    [
      {
        "connector" : "test-connector-action-1",
        "connector_name" : "test-connector-name-action-1",
        "source_type" : "resource",
        "event_type" : "check",
        "component" :  "test-component-action-1",
        "resource" : "test-resource-action-1-1",
        "state" : 2,
        "output" : "test-output-action-1"
      },
      {
        "connector" : "test-connector-action-1",
        "connector_name" : "test-connector-name-action-1",
        "source_type" : "resource",
        "event_type" : "check",
        "component" :  "test-component-action-1",
        "resource" : "test-resource-action-1-2",
        "state" : 1,
        "output" : "test-output-action-1"
      }
    ]
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-component-action-1&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "tickets": [
              {
                "_t": "assocticket",
                "a": "test-scenario-action-1-action-1-author test-resource-action-1-1",
                "m": "Scenario: test-scenario-action-1-name. Ticket ID: test-scenario-action-1-action-1-ticket. Ticket URL: test-scenario-action-1-action-1-ticket-url. Ticket ticket_param_1: ticket_value_1.",
                "ticket": "test-scenario-action-1-action-1-ticket",
                "ticket_rule_id": "test-scenario-action-1",
                "ticket_rule_name": "Scenario: test-scenario-action-1-name",
                "ticket_comment": "test-scenario-action-1-action-1-comment",
                "ticket_system_name": "test-scenario-action-1-action-1-system-name",
                "ticket_url": "test-scenario-action-1-action-1-ticket-url",
                "ticket_data": {
                  "ticket_param_1": "ticket_value_1"
                }
              }
            ],
            "ticket": {
              "_t": "assocticket",
              "a": "test-scenario-action-1-action-1-author test-resource-action-1-1",
              "m": "Scenario: test-scenario-action-1-name. Ticket ID: test-scenario-action-1-action-1-ticket. Ticket URL: test-scenario-action-1-action-1-ticket-url. Ticket ticket_param_1: ticket_value_1.",
              "ticket": "test-scenario-action-1-action-1-ticket",
              "ticket_rule_id": "test-scenario-action-1",
              "ticket_rule_name": "Scenario: test-scenario-action-1-name",
              "ticket_comment": "test-scenario-action-1-action-1-comment",
              "ticket_system_name": "test-scenario-action-1-action-1-system-name",
              "ticket_url": "test-scenario-action-1-action-1-ticket-url",
              "ticket_data": {
                "ticket_param_1": "ticket_value_1"
              }
            },
            "ack": {
              "_t": "ack",
              "a": "root",
              "user_id": "root",
              "m": "test-scenario-action-1-action-2-output test-resource-action-1-1 2"
            },
            "connector": "test-connector-action-1",
            "connector_name": "test-connector-name-action-1",
            "component": "test-component-action-1",
            "resource": "test-resource-action-1-1"
          }
        },
        {
          "v": {
            "tickets": [
              {
                "_t": "assocticket",
                "a": "test-scenario-action-1-action-1-author test-resource-action-1-2",
                "m": "Scenario: test-scenario-action-1-name. Ticket ID: test-scenario-action-1-action-1-ticket. Ticket URL: test-scenario-action-1-action-1-ticket-url. Ticket ticket_param_1: ticket_value_1.",
                "ticket": "test-scenario-action-1-action-1-ticket",
                "ticket_rule_id": "test-scenario-action-1",
                "ticket_rule_name": "Scenario: test-scenario-action-1-name",
                "ticket_comment": "test-scenario-action-1-action-1-comment",
                "ticket_system_name": "test-scenario-action-1-action-1-system-name",
                "ticket_url": "test-scenario-action-1-action-1-ticket-url",
                "ticket_data": {
                  "ticket_param_1": "ticket_value_1"
                }
              }
            ],
            "ticket": {
              "_t": "assocticket",
              "a": "test-scenario-action-1-action-1-author test-resource-action-1-2",
              "m": "Scenario: test-scenario-action-1-name. Ticket ID: test-scenario-action-1-action-1-ticket. Ticket URL: test-scenario-action-1-action-1-ticket-url. Ticket ticket_param_1: ticket_value_1.",
              "ticket": "test-scenario-action-1-action-1-ticket",
              "ticket_rule_id": "test-scenario-action-1",
              "ticket_rule_name": "Scenario: test-scenario-action-1-name",
              "ticket_comment": "test-scenario-action-1-action-1-comment",
              "ticket_system_name": "test-scenario-action-1-action-1-system-name",
              "ticket_url": "test-scenario-action-1-action-1-ticket-url",
              "ticket_data": {
                "ticket_param_1": "ticket_value_1"
              }
            },
            "ack": {
              "_t": "ack",
              "a": "root",
              "user_id": "root",
              "m": "test-scenario-action-1-action-2-output test-resource-action-1-2 1"
            },
            "connector": "test-connector-action-1",
            "connector_name": "test-connector-name-action-1",
            "component": "test-component-action-1",
            "resource": "test-resource-action-1-2"
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
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ (index .lastResponse.data 1)._id }}",
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
                "a": "test-scenario-action-1-action-1-author test-resource-action-1-1",
                "user_id": "",
                "m": "Scenario: test-scenario-action-1-name. Ticket ID: test-scenario-action-1-action-1-ticket. Ticket URL: test-scenario-action-1-action-1-ticket-url. Ticket ticket_param_1: ticket_value_1.",
                "ticket": "test-scenario-action-1-action-1-ticket",
                "ticket_rule_id": "test-scenario-action-1",
                "ticket_rule_name": "Scenario: test-scenario-action-1-name",
                "ticket_comment": "test-scenario-action-1-action-1-comment",
                "ticket_system_name": "test-scenario-action-1-action-1-system-name",
                "ticket_url": "test-scenario-action-1-action-1-ticket-url",
                "ticket_data": {
                  "ticket_param_1": "ticket_value_1"
                }
              },
              {
                "_t": "ack",
                "a": "root",
                "user_id": "root",
                "m": "test-scenario-action-1-action-2-output test-resource-action-1-1 2"
              },
              {
                "_t": "changestate",
                "a": "system",
                "user_id": "",
                "m": "test-scenario-action-1-action-3-output Test test-resource-action-1-1 2",
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
                "_t": "statusinc"
              },
              {
                "_t": "assocticket",
                "a": "test-scenario-action-1-action-1-author test-resource-action-1-2",
                "user_id": "",
                "m": "Scenario: test-scenario-action-1-name. Ticket ID: test-scenario-action-1-action-1-ticket. Ticket URL: test-scenario-action-1-action-1-ticket-url. Ticket ticket_param_1: ticket_value_1.",
                "ticket": "test-scenario-action-1-action-1-ticket",
                "ticket_rule_id": "test-scenario-action-1",
                "ticket_rule_name": "Scenario: test-scenario-action-1-name",
                "ticket_comment": "test-scenario-action-1-action-1-comment",
                "ticket_system_name": "test-scenario-action-1-action-1-system-name",
                "ticket_url": "test-scenario-action-1-action-1-ticket-url",
                "ticket_data": {
                  "ticket_param_1": "ticket_value_1"
                }
              },
              {
                "_t": "ack",
                "a": "root",
                "user_id": "root",
                "m": "test-scenario-action-1-action-2-output test-resource-action-1-2 1"
              },
              {
                "_t": "changestate",
                "a": "system",
                "user_id": "",
                "m": "test-scenario-action-1-action-3-output Test test-resource-action-1-2 1",
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

  Scenario: given scenario and check event should not update alarm
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-negative-1-name",
      "priority": 10003,
      "enabled": true,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-component-action-negative-1-should-not-match"
                }
              }
            ]
          ],
          "type": "assocticket",
          "parameters": {
            "forward_author": false,
            "author": "test-scenario-action-negative-1-action-negative-1-author {{ `{{ .Alarm.Value.Resource }}` }}",
            "ticket": "test-scenario-action-negative-1-action-negative-1-ticket"
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
                  "value": "test-resource-action-negative-1-should-not-match"
                }
              }
            ]
          ],
          "type": "ack",
          "parameters": {
            "forward_author": true,
            "author": "test-scenario-action-negative-1-action-2-author {{ `{{ .Alarm.Value.Resource }}` }}",
            "output": "test-scenario-action-negative-1-action-2-output {{ `{{ .Entity.Name }} {{ .Alarm.Value.State.Value }}` }}"
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
      "connector" : "test-connector-action-negative-1",
      "connector_name" : "test-connector-name-action-negative-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-negative-1",
      "resource" : "test-resource-action-negative-1",
      "state" : 2,
      "output" : "test-output-action-negative-1"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-component-action-negative-1&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-action-negative-1",
            "connector_name": "test-connector-name-action-negative-1",
            "component": "test-component-action-negative-1",
            "resource": "test-resource-action-negative-1"
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
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
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

  Scenario: given delayed scenario and check event should update alarm
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "_id": "test-scenario-action-2",
      "name": "test-scenario-action-2-name",
      "priority": 10004,
      "enabled": true,
      "triggers": ["create"],
      "delay": {
        "value": 5,
        "unit": "s"
      },
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.resource",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-action-2"
                }
              }
            ]
          ],
          "type": "assocticket",
          "parameters": {
            "forward_author": true,
            "author": "test-scenario-action-2-action-2-author {{ `{{ .Alarm.Value.Resource }}` }}",
            "ticket": "test-scenario-action-2-action-1-ticket"
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        },
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-component-action-2"
                }
              }
            ]
          ],
          "type": "ack",
          "parameters": {
            "author": "test-scenario-action-2-action-2-author {{ `{{ .Alarm.Value.Resource }}` }}",
            "output": "test-scenario-action-2-action-2-output {{ `{{ .Entity.Name }}` }}"
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
      "connector" : "test-connector-action-2",
      "connector_name" : "test-connector-name-action-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-2",
      "resource" : "test-resource-action-2",
      "state" : 2,
      "output" : "test-output-action-2"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-action-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "tickets": [
              {
                "_t": "assocticket",
                "a": "root",
                "ticket": "test-scenario-action-2-action-1-ticket"
              }
            ],
            "ticket": {
              "_t": "assocticket",
              "a": "root",
              "ticket": "test-scenario-action-2-action-1-ticket"
            },
            "ack": {
              "_t": "ack",
              "a": "test-scenario-action-2-action-2-author test-resource-action-2",
              "user_id": "",
              "m": "test-scenario-action-2-action-2-output test-resource-action-2"
            },
            "connector": "test-connector-action-2",
            "connector_name": "test-connector-name-action-2",
            "component": "test-component-action-2",
            "resource": "test-resource-action-2"
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
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "assocticket",
                "a": "root",
                "user_id": "root",
                "ticket": "test-scenario-action-2-action-1-ticket"
              },
              {
                "_t": "ack",
                "a": "test-scenario-action-2-action-2-author test-resource-action-2",
                "user_id": "",
                "m": "test-scenario-action-2-action-2-output test-resource-action-2"
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

  Scenario: given scenario with emit trigger and check event should update alarm
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "_id": "test-scenario-action-3",
      "name": "test-scenario-action-3-name-1",
      "priority": 10005,
      "enabled": true,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.resource",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-action-3"
                }
              }
            ]
          ],
          "type": "assocticket",
          "parameters": {
            "ticket": "test-scenario-action-3-ticket"
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": true
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-3-name-2",
      "priority": 10006,
      "enabled": true,
      "triggers": ["assocticket"],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-component-action-3"
                }
              }
            ]
          ],
          "type": "ack",
          "parameters": {
            "forward_author": true,
            "output": "test-output-action-3-{{ `{{ .Alarm.Value.Connector }}` }}"
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
      "connector" : "test-connector-action-3",
      "connector_name" : "test-connector-name-action-3",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-3",
      "resource" : "test-resource-action-3",
      "state" : 2,
      "output" : "test-output-action-3"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-action-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "tickets": [
              {
                "_t": "assocticket",
                "a": "system",
                "ticket": "test-scenario-action-3-ticket"
              }
            ],
            "ticket": {
              "_t": "assocticket",
              "a": "system",
              "ticket": "test-scenario-action-3-ticket"
            },
            "ack": {
              "_t": "ack",
              "a": "root",
              "user_id": "root",
              "m": "test-output-action-3-test-connector-action-3"
            },
            "connector": "test-connector-action-3",
            "connector_name": "test-connector-name-action-3",
            "component": "test-component-action-3",
            "resource": "test-resource-action-3"
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
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "assocticket",
                "a": "system",
                "user_id": "",
                "ticket": "test-scenario-action-3-ticket"
              },
              {
                "_t": "ack",
                "a": "root",
                "user_id": "root",
                "m": "test-output-action-3-test-connector-action-3"
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

  Scenario: given scenario and ack already acked alarm should trigger scenario
    Given I am admin
    When I send an event:
    """
    {
      "connector" : "test-connector-double-ack-alarm-1",
      "connector_name" : "test-connector-name-double-ack-alarm-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-double-ack-alarm-1",
      "resource" : "test-resource-double-ack-alarm-1",
      "state" : 2,
      "output" : "test-double-ack-alarm-1"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """
    {
      "connector" : "test-connector-double-ack-alarm-1",
      "connector_name" : "test-connector-name-double-ack-alarm-1",
      "source_type" : "resource",
      "event_type" : "ack",
      "component" :  "test-component-double-ack-alarm-1",
      "resource" : "test-resource-double-ack-alarm-1",
      "output" : "test-double-ack-alarm-1"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-component-double-ack-alarm-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "ack": {
              "_t": "ack",
              "a": "root",
              "user_id": "root",
              "m": "test-double-ack-alarm-1"
            },
            "connector": "test-connector-double-ack-alarm-1",
            "connector_name": "test-connector-name-double-ack-alarm-1",
            "component": "test-component-double-ack-alarm-1",
            "resource": "test-resource-double-ack-alarm-1",
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
        "total_count": 1
      }
    }
    """
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-double-ack",
      "priority": 10007,
      "enabled": true,
      "triggers": ["ack"],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-component-double-ack-alarm-1"
                }
              }
            ]
          ],
          "type": "changestate",
          "parameters": {
            "state": 3,
            "forward_author": false,
            "author": "",
            "output": "state is changed"
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
      "connector" : "test-connector-double-ack-alarm-1",
      "connector_name" : "test-connector-name-double-ack-alarm-1",
      "source_type" : "resource",
      "event_type" : "ack",
      "component" :  "test-component-double-ack-alarm-1",
      "resource" : "test-resource-double-ack-alarm-1",
      "output" : "test-double-ack-alarm-1"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-component-double-ack-alarm-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "ack": {
              "_t": "ack",
              "a": "root",
              "user_id": "root",
              "m": "test-double-ack-alarm-1"
            },
            "connector": "test-connector-double-ack-alarm-1",
            "connector_name": "test-connector-name-double-ack-alarm-1",
            "component": "test-component-double-ack-alarm-1",
            "resource": "test-resource-double-ack-alarm-1",
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

  Scenario: given scenario with old patterns should update alarm with backward compatibility
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-scenario-backward-compatibility-actions-connector",
      "connector_name" : "test-scenario-backward-compatibility-actions-connector-name",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-scenario-backward-compatibility-actions-component",
      "resource" : "test-scenario-backward-compatibility-actions-resource",
      "state" : 2
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-scenario-backward-compatibility-actions-resource&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "tickets": [
              {
                "_t": "assocticket",
                "a": "test-scenario-backward-compatibility-actions-1-author",
                "ticket": "test-scenario-backward-compatibility-actions-1-ticket"
              }
            ],
            "ticket": {
              "_t": "assocticket",
              "a": "test-scenario-backward-compatibility-actions-1-author",
              "ticket": "test-scenario-backward-compatibility-actions-1-ticket"
            },
            "ack": {
              "_t": "ack",
              "a": "test-scenario-backward-compatibility-actions-1-author",
              "m": "test-scenario-backward-compatibility-actions-1-output"
            },
            "connector" : "test-scenario-backward-compatibility-actions-connector",
            "connector_name" : "test-scenario-backward-compatibility-actions-connector-name",
            "component" :  "test-scenario-backward-compatibility-actions-component",
            "resource" : "test-scenario-backward-compatibility-actions-resource",
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
                "_t": "ack"
              },
              {
                "_t": "assocticket"
              },
              {
                "_t": "changestate"
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
