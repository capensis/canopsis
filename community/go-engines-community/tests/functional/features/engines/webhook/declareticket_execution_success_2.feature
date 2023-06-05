Feature: run a declare ticket rule
  I need to be able to run a declare ticket rule
  Only admin should be able to run a declare ticket rule

  @concurrent
  Scenario: given declare ticket rule and ticket resources should update resources
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type" : "check",
        "state" : 2,
        "connector" : "test-connector-declareticket-execution-second-1",
        "connector_name" : "test-connector-name-declareticket-execution-second-1",
        "component" :  "test-component-declareticket-execution-second-1",
        "source_type" : "component"
      },
      {
        "event_type" : "check",
        "state" : 2,
        "connector" : "test-connector-declareticket-execution-second-1",
        "connector_name" : "test-connector-name-declareticket-execution-second-1",
        "component" :  "test-component-declareticket-execution-second-1",
        "resource" : "test-resource-declareticket-execution-second-1",
        "source_type" : "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-component-declareticket-execution-second-1&sort_by=v.resource&sort=asc
    Then the response code should be 200
    When I save response componentAlarmId={{ (index .lastResponse.data 0)._id }}
    When I save response resourceAlarmId={{ (index .lastResponse.data 1)._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-declareticketrule-declareticket-execution-second-1-name",
      "system_name": "test-declareticketrule-declareticket-execution-second-1-system-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"_id\":\"test-ticket-declareticket-execution-second-1\",\"url\":\"https://test-host\",\"name\":\"{{ `{{ range .Alarms }}{{ .Entity.Name }}{{ end }}`}}\"}",
            "method": "POST",
            "auth": {
              "username": "test",
              "password": "test"
            }
          },
          "declare_ticket": {
            "ticket_id": "_id",
            "ticket_url": "url",
            "name": "name"
          },
          "stop_on_fail": true
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-component-declareticket-execution-second-1"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleId={{ .lastResponse._id }}
    When I do POST /api/v4/cat/bulk/declare-ticket-executions:
    """json
    [
      {
        "_id": "{{ .ruleId }}",
        "alarms": [
          "{{ .componentAlarmId }}"
        ],
        "comment": "test-comment-declareticket-execution-second-1",
        "ticket_resources": true
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "item": {
          "_id": "{{ .ruleId }}",
          "alarms": [
            "{{ .componentAlarmId }}"
          ],
          "comment": "test-comment-declareticket-execution-second-1",
          "ticket_resources": true
        }
      }
    ]
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "declareticketwebhook",
      "connector" : "test-connector-declareticket-execution-second-1",
      "connector_name" : "test-connector-name-declareticket-execution-second-1",
      "component" :  "test-component-declareticket-execution-second-1",
      "resource" : "test-resource-declareticket-execution-second-1",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-component-declareticket-execution-second-1&sort_by=v.resource&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "ticket": {
              "_t": "declareticket",
              "a": "root",
              "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-1-name. Ticket ID: test-ticket-declareticket-execution-second-1. Ticket URL: https://test-host. Ticket name: test-component-declareticket-execution-second-1.",
              "ticket": "test-ticket-declareticket-execution-second-1",
              "ticket_url": "https://test-host",
              "ticket_data": {
                "name": "test-component-declareticket-execution-second-1"
              },
              "ticket_comment": "test-comment-declareticket-execution-second-1",
              "ticket_rule_id": "{{ .ruleId }}",
              "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-1-name",
              "ticket_system_name": "test-declareticketrule-declareticket-execution-second-1-system-name"
            },
            "tickets": [
              {
                "_t": "declareticket",
                "a": "root",
                "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-1-name. Ticket ID: test-ticket-declareticket-execution-second-1. Ticket URL: https://test-host. Ticket name: test-component-declareticket-execution-second-1.",
                "ticket": "test-ticket-declareticket-execution-second-1",
                "ticket_url": "https://test-host",
                "ticket_data": {
                  "name": "test-component-declareticket-execution-second-1"
                },
                "ticket_comment": "test-comment-declareticket-execution-second-1",
                "ticket_rule_id": "{{ .ruleId }}",
                "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-1-name",
                "ticket_system_name": "test-declareticketrule-declareticket-execution-second-1-system-name"
              }
            ],
            "connector" : "test-connector-declareticket-execution-second-1",
            "connector_name" : "test-connector-name-declareticket-execution-second-1",
            "component" :  "test-component-declareticket-execution-second-1"
          }
        },
        {
          "v": {
            "ticket": {
              "_t": "declareticket",
              "a": "root",
              "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-1-name. Ticket ID: test-ticket-declareticket-execution-second-1. Ticket URL: https://test-host. Ticket name: test-component-declareticket-execution-second-1.",
              "ticket": "test-ticket-declareticket-execution-second-1",
              "ticket_url": "https://test-host",
              "ticket_data": {
                "name": "test-component-declareticket-execution-second-1"
              },
              "ticket_comment": "test-comment-declareticket-execution-second-1",
              "ticket_rule_id": "{{ .ruleId }}",
              "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-1-name",
              "ticket_system_name": "test-declareticketrule-declareticket-execution-second-1-system-name"
            },
            "tickets": [
              {
                "_t": "declareticket",
                "a": "root",
                "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-1-name. Ticket ID: test-ticket-declareticket-execution-second-1. Ticket URL: https://test-host. Ticket name: test-component-declareticket-execution-second-1.",
                "ticket": "test-ticket-declareticket-execution-second-1",
                "ticket_url": "https://test-host",
                "ticket_data": {
                  "name": "test-component-declareticket-execution-second-1"
                },
                "ticket_comment": "test-comment-declareticket-execution-second-1",
                "ticket_rule_id": "{{ .ruleId }}",
                "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-1-name",
                "ticket_system_name": "test-declareticketrule-declareticket-execution-second-1-system-name"
              }
            ],
            "connector" : "test-connector-declareticket-execution-second-1",
            "connector_name" : "test-connector-name-declareticket-execution-second-1",
            "component" :  "test-component-declareticket-execution-second-1",
            "resource" : "test-resource-declareticket-execution-second-1"
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
        "_id": "{{ .componentAlarmId }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ .resourceAlarmId }}",
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
        "_t": "webhookstart",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-1-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-1-name"
      },
      {
        "_t": "declareticket",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-1-name. Ticket ID: test-ticket-declareticket-execution-second-1. Ticket URL: https://test-host. Ticket name: test-component-declareticket-execution-second-1."
      }
    ]
    """
    Then the response array key "1.data.steps.data" should contain only:
    """json
    [
      {
        "_t": "stateinc"
      },
      {
        "_t": "statusinc"
      },
      {
        "_t": "declareticket",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-1-name. Ticket ID: test-ticket-declareticket-execution-second-1. Ticket URL: https://test-host. Ticket name: test-component-declareticket-execution-second-1."
      }
    ]
    """

  @concurrent
  Scenario: given declare ticket rule should run scenario
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-declareticket-execution-second-2-name",
      "priority": 10080,
      "enabled": true,
      "triggers": ["declareticketwebhook"],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.resource",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-declareticket-execution-second-2"
                }
              }
            ]
          ],
          "type": "ack",
          "parameters": {
            "output": "test-scenario-declareticket-execution-second-2-output"
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
      "event_type" : "check",
      "state" : 2,
      "connector" : "test-connector-declareticket-execution-second-2",
      "connector_name" : "test-connector-name-declareticket-execution-second-2",
      "component" :  "test-component-declareticket-execution-second-2",
      "resource" : "test-resource-declareticket-execution-second-2",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-component-declareticket-execution-second-2
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-declareticketrule-declareticket-execution-second-2-name",
      "system_name": "test-declareticketrule-declareticket-execution-second-2-system-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"_id\":\"test-ticket-declareticket-execution-second-2\"}",
            "method": "POST"
          },
          "declare_ticket": {
            "ticket_id": "_id"
          },
          "stop_on_fail": true
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-declareticket-execution-second-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleId={{ .lastResponse._id }}
    When I do POST /api/v4/cat/bulk/declare-ticket-executions:
    """json
    [
      {
        "_id": "{{ .ruleId }}",
        "alarms": [
          "{{ .alarmId }}"
        ],
        "comment": "test-comment-declareticket-execution-second-2"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200
      }
    ]
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "trigger",
      "connector" : "test-connector-declareticket-execution-second-2",
      "connector_name" : "test-connector-name-declareticket-execution-second-2",
      "component" :  "test-component-declareticket-execution-second-2",
      "resource" : "test-resource-declareticket-execution-second-2",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-component-declareticket-execution-second-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "ticket": {
              "_t": "declareticket",
              "a": "root",
              "ticket": "test-ticket-declareticket-execution-second-2"
            },
            "tickets": [
              {
                "_t": "declareticket",
                "a": "root",
                "ticket": "test-ticket-declareticket-execution-second-2"
              }
            ],
            "connector" : "test-connector-declareticket-execution-second-2",
            "connector_name" : "test-connector-name-declareticket-execution-second-2",
            "component" :  "test-component-declareticket-execution-second-2",
            "resource" : "test-resource-declareticket-execution-second-2"
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
        "_id": "{{ .alarmId }}",
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
        "_t": "webhookstart",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-2-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-2-name"
      },
      {
        "_t": "declareticket",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-2-name. Ticket ID: test-ticket-declareticket-execution-second-2."
      },
      {
        "_t": "ack",
        "a": "system",
        "user_id": "",
        "m": "test-scenario-declareticket-execution-second-2-output"
      }
    ]
    """
    Then the response array key "0.data.steps.data" should contain in order:
    """json
    [
     {
        "_t": "declareticket"
      },
      {
        "_t": "ack"
      }
    ]
    """

  @concurrent
  Scenario: given declare ticket rule shouldn't run scenario
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-declareticket-execution-second-3-name",
      "priority": 10081,
      "enabled": true,
      "triggers": ["declareticketwebhook"],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.resource",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-declareticket-execution-second-3"
                }
              }
            ]
          ],
          "type": "ack",
          "parameters": {
            "output": "test-scenario-declareticket-execution-second-3-output"
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
      "event_type" : "check",
      "state" : 2,
      "connector" : "test-connector-declareticket-execution-second-3",
      "connector_name" : "test-connector-name-declareticket-execution-second-3",
      "component" :  "test-component-declareticket-execution-second-3",
      "resource" : "test-resource-declareticket-execution-second-3",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-component-declareticket-execution-second-3
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-declareticketrule-declareticket-execution-second-3-name",
      "system_name": "test-declareticketrule-declareticket-execution-second-3-system-name",
      "enabled": true,
      "emit_trigger": false,
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"_id\":\"test-ticket-declareticket-execution-second-3\"}",
            "method": "POST"
          },
          "declare_ticket": {
            "ticket_id": "_id"
          },
          "stop_on_fail": true
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-declareticket-execution-second-3"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleId={{ .lastResponse._id }}
    When I do POST /api/v4/cat/bulk/declare-ticket-executions:
    """json
    [
      {
        "_id": "{{ .ruleId }}",
        "alarms": [
          "{{ .alarmId }}"
        ],
        "comment": "test-comment-declareticket-execution-second-3"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-component-declareticket-execution-second-3 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "ticket": {
              "_t": "declareticket",
              "a": "root",
              "ticket": "test-ticket-declareticket-execution-second-3"
            },
            "tickets": [
              {
                "_t": "declareticket",
                "a": "root",
                "ticket": "test-ticket-declareticket-execution-second-3"
              }
            ],
            "connector" : "test-connector-declareticket-execution-second-3",
            "connector_name" : "test-connector-name-declareticket-execution-second-3",
            "component" :  "test-component-declareticket-execution-second-3",
            "resource" : "test-resource-declareticket-execution-second-3"
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
    When I wait 100ms
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .alarmId }}",
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
        "_t": "webhookstart",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-3-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-3-name"
      },
      {
        "_t": "declareticket",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-3-name. Ticket ID: test-ticket-declareticket-execution-second-3."
      }
    ]
    """

  @concurrent
  Scenario: given declare ticket rule with empty ticket_id and multiple alarms should execute request for each alarms
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type" : "check",
        "state" : 2,
        "connector" : "test-connector-declareticket-execution-second-4",
        "connector_name" : "test-connector-name-declareticket-execution-second-4",
        "component" :  "test-component-declareticket-execution-second-4",
        "resource" : "test-resource-declareticket-execution-second-4-1",
        "source_type" : "resource"
      },
      {
        "event_type" : "check",
        "state" : 2,
        "connector" : "test-connector-declareticket-execution-second-4",
        "connector_name" : "test-connector-name-declareticket-execution-second-4",
        "component" :  "test-component-declareticket-execution-second-4",
        "resource" : "test-resource-declareticket-execution-second-4-2",
        "source_type" : "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-execution-second-4&sort_by=v.resource&sort=asc
    Then the response code should be 200
    When I save response alarmId1={{ (index .lastResponse.data 0)._id }}
    When I save response alarmId2={{ (index .lastResponse.data 1)._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-declareticketrule-declareticket-execution-second-4-name",
      "system_name": "test-declareticketrule-declareticket-execution-second-4-system-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"_id\":\"test-ticket-declareticket-execution-second-4\",\"url\":\"https://test/test-ticket-declareticket-execution-second-4\",\"name\":\"test-ticket-declareticket-execution-second-4 {{ `{{ range .Alarms }}{{ .Value.Resource }}{{ end }}`}}\"}",
            "method": "POST",
            "auth": {
              "username": "test",
              "password": "test"
            }
          },
          "declare_ticket": {
            "ticket_url": "url",
            "name": "name"
          },
          "stop_on_fail": true
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-declareticket-execution-second-4-1",
                "test-resource-declareticket-execution-second-4-2"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleId={{ .lastResponse._id }}
    When I do POST /api/v4/cat/bulk/declare-ticket-executions:
    """json
    [
      {
        "_id": "{{ .ruleId }}",
        "alarms": [
          "{{ .alarmId1 }}"
        ],
        "comment": "test-comment-declareticket-execution-second-4-1"
      },
      {
        "_id": "{{ .ruleId }}",
        "alarms": [
          "{{ .alarmId2 }}"
        ],
        "comment": "test-comment-declareticket-execution-second-4-2"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "item": {
          "_id": "{{ .ruleId }}",
          "alarms": [
            "{{ .alarmId1 }}"
          ],
          "comment": "test-comment-declareticket-execution-second-4-1"
        }
      },
      {
        "status": 200,
        "item": {
          "_id": "{{ .ruleId }}",
          "alarms": [
            "{{ .alarmId2 }}"
          ],
          "comment": "test-comment-declareticket-execution-second-4-2"
        }
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-execution-second-4&sort_by=v.resource&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "ticket": {
              "_t": "declareticket",
              "a": "root",
              "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-4-name. Ticket ID: N/A. Ticket URL: https://test/test-ticket-declareticket-execution-second-4. Ticket name: test-ticket-declareticket-execution-second-4 test-resource-declareticket-execution-second-4-1.",
              "ticket": "N/A",
              "ticket_url": "https://test/test-ticket-declareticket-execution-second-4",
              "ticket_data": {
                "name": "test-ticket-declareticket-execution-second-4 test-resource-declareticket-execution-second-4-1"
              },
              "ticket_comment": "test-comment-declareticket-execution-second-4-1",
              "ticket_rule_id": "{{ .ruleId }}",
              "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-4-name",
              "ticket_system_name": "test-declareticketrule-declareticket-execution-second-4-system-name"
            },
            "tickets": [
              {
                "_t": "declareticket",
                "a": "root",
                "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-4-name. Ticket ID: N/A. Ticket URL: https://test/test-ticket-declareticket-execution-second-4. Ticket name: test-ticket-declareticket-execution-second-4 test-resource-declareticket-execution-second-4-1.",
                "ticket": "N/A",
                "ticket_url": "https://test/test-ticket-declareticket-execution-second-4",
                "ticket_data": {
                  "name": "test-ticket-declareticket-execution-second-4 test-resource-declareticket-execution-second-4-1"
                },
                "ticket_comment": "test-comment-declareticket-execution-second-4-1",
                "ticket_rule_id": "{{ .ruleId }}",
                "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-4-name",
                "ticket_system_name": "test-declareticketrule-declareticket-execution-second-4-system-name"
              }
            ],
            "connector" : "test-connector-declareticket-execution-second-4",
            "connector_name" : "test-connector-name-declareticket-execution-second-4",
            "component" :  "test-component-declareticket-execution-second-4",
            "resource" : "test-resource-declareticket-execution-second-4-1"
          }
        },
        {
          "v": {
            "ticket": {
              "_t": "declareticket",
              "a": "root",
              "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-4-name. Ticket ID: N/A. Ticket URL: https://test/test-ticket-declareticket-execution-second-4. Ticket name: test-ticket-declareticket-execution-second-4 test-resource-declareticket-execution-second-4-2.",
              "ticket": "N/A",
              "ticket_url": "https://test/test-ticket-declareticket-execution-second-4",
              "ticket_data": {
                "name": "test-ticket-declareticket-execution-second-4 test-resource-declareticket-execution-second-4-2"
              },
              "ticket_comment": "test-comment-declareticket-execution-second-4-2",
              "ticket_rule_id": "{{ .ruleId }}",
              "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-4-name",
              "ticket_system_name": "test-declareticketrule-declareticket-execution-second-4-system-name"
            },
            "tickets": [
              {
                "_t": "declareticket",
                "a": "root",
                "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-4-name. Ticket ID: N/A. Ticket URL: https://test/test-ticket-declareticket-execution-second-4. Ticket name: test-ticket-declareticket-execution-second-4 test-resource-declareticket-execution-second-4-2.",
                "ticket": "N/A",
                "ticket_url": "https://test/test-ticket-declareticket-execution-second-4",
                "ticket_data": {
                  "name": "test-ticket-declareticket-execution-second-4 test-resource-declareticket-execution-second-4-2"
                },
                "ticket_comment": "test-comment-declareticket-execution-second-4-2",
                "ticket_rule_id": "{{ .ruleId }}",
                "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-4-name",
                "ticket_system_name": "test-declareticketrule-declareticket-execution-second-4-system-name"
              }
            ],
            "connector" : "test-connector-declareticket-execution-second-4",
            "connector_name" : "test-connector-name-declareticket-execution-second-4",
            "component" :  "test-component-declareticket-execution-second-4",
            "resource" : "test-resource-declareticket-execution-second-4-2"
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
        "_id": "{{ .alarmId1 }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ .alarmId2 }}",
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
        "_t": "webhookstart",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-4-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-4-name"
      },
      {
        "_t": "declareticket",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-4-name. Ticket ID: N/A. Ticket URL: https://test/test-ticket-declareticket-execution-second-4. Ticket name: test-ticket-declareticket-execution-second-4 test-resource-declareticket-execution-second-4-1."
      }
    ]
    """
    Then the response array key "1.data.steps.data" should contain only:
    """json
    [
      {
        "_t": "stateinc"
      },
      {
        "_t": "statusinc"
      },
      {
        "_t": "webhookstart",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-4-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-4-name"
      },
      {
        "_t": "declareticket",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-second-4-name. Ticket ID: N/A. Ticket URL: https://test/test-ticket-declareticket-execution-second-4. Ticket name: test-ticket-declareticket-execution-second-4 test-resource-declareticket-execution-second-4-2."
      }
    ]
    """
