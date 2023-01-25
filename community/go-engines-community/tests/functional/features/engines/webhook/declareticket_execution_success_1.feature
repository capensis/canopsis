Feature: run a declare ticket rule
  I need to be able to run a declare ticket rule
  Only admin should be able to run a declare ticket rule

  @concurrent
  Scenario: given declare ticket rule and multiple alarms should execute request for each alarms
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type" : "check",
        "state" : 2,
        "connector" : "test-connector-declareticket-execution-1",
        "connector_name" : "test-connector-name-declareticket-execution-1",
        "component" :  "test-component-declareticket-execution-1",
        "resource" : "test-resource-declareticket-execution-1-1",
        "source_type" : "resource"
      },
      {
        "event_type" : "check",
        "state" : 2,
        "connector" : "test-connector-declareticket-execution-1",
        "connector_name" : "test-connector-name-declareticket-execution-1",
        "component" :  "test-component-declareticket-execution-1",
        "resource" : "test-resource-declareticket-execution-1-2",
        "source_type" : "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-execution-1&sort_by=v.resource&sort=asc
    Then the response code should be 200
    When I save response alarmId1={{ (index .lastResponse.data 0)._id }}
    When I save response alarmId2={{ (index .lastResponse.data 1)._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-declareticketrule-declareticket-execution-1-name",
      "system_name": "test-declareticketrule-declareticket-execution-1-system-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"_id\":\"test-ticket-declareticket-execution-1\",\"url\":\"https://test/test-ticket-declareticket-execution-1\",\"name\":\"test-ticket-declareticket-execution-1 {{ `{{ range .Alarms }}{{ .Value.Resource }}{{ end }}`}}\"}",
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
                "test-resource-declareticket-execution-1-1",
                "test-resource-declareticket-execution-1-2"
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
        "comment": "test-comment-declareticket-execution-1-1"
      },
      {
        "_id": "{{ .ruleId }}",
        "alarms": [
          "{{ .alarmId2 }}"
        ],
        "comment": "test-comment-declareticket-execution-1-2"
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
          "comment": "test-comment-declareticket-execution-1-1"
        }
      },
      {
        "status": 200,
        "item": {
          "_id": "{{ .ruleId }}",
          "alarms": [
            "{{ .alarmId2 }}"
          ],
          "comment": "test-comment-declareticket-execution-1-2"
        }
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-execution-1&sort_by=v.resource&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "ticket": {
              "_t": "declareticket",
              "a": "root",
              "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-1-name. Ticket ID: test-ticket-declareticket-execution-1. Ticket URL: https://test/test-ticket-declareticket-execution-1. Ticket name: test-ticket-declareticket-execution-1 test-resource-declareticket-execution-1-1.",
              "ticket": "test-ticket-declareticket-execution-1",
              "ticket_url": "https://test/test-ticket-declareticket-execution-1",
              "ticket_data": {
                "name": "test-ticket-declareticket-execution-1 test-resource-declareticket-execution-1-1"
              },
              "ticket_comment": "test-comment-declareticket-execution-1-1",
              "ticket_rule_id": "{{ .ruleId }}",
              "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-execution-1-name",
              "ticket_system_name": "test-declareticketrule-declareticket-execution-1-system-name"
            },
            "tickets": [
              {
                "_t": "declareticket",
                "a": "root",
                "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-1-name. Ticket ID: test-ticket-declareticket-execution-1. Ticket URL: https://test/test-ticket-declareticket-execution-1. Ticket name: test-ticket-declareticket-execution-1 test-resource-declareticket-execution-1-1.",
                "ticket": "test-ticket-declareticket-execution-1",
                "ticket_url": "https://test/test-ticket-declareticket-execution-1",
                "ticket_data": {
                  "name": "test-ticket-declareticket-execution-1 test-resource-declareticket-execution-1-1"
                },
                "ticket_comment": "test-comment-declareticket-execution-1-1",
                "ticket_rule_id": "{{ .ruleId }}",
                "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-execution-1-name",
                "ticket_system_name": "test-declareticketrule-declareticket-execution-1-system-name"
              }
            ],
            "connector" : "test-connector-declareticket-execution-1",
            "connector_name" : "test-connector-name-declareticket-execution-1",
            "component" :  "test-component-declareticket-execution-1",
            "resource" : "test-resource-declareticket-execution-1-1"
          }
        },
        {
          "v": {
            "ticket": {
              "_t": "declareticket",
              "a": "root",
              "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-1-name. Ticket ID: test-ticket-declareticket-execution-1. Ticket URL: https://test/test-ticket-declareticket-execution-1. Ticket name: test-ticket-declareticket-execution-1 test-resource-declareticket-execution-1-2.",
              "ticket": "test-ticket-declareticket-execution-1",
              "ticket_url": "https://test/test-ticket-declareticket-execution-1",
              "ticket_data": {
                "name": "test-ticket-declareticket-execution-1 test-resource-declareticket-execution-1-2"
              },
              "ticket_comment": "test-comment-declareticket-execution-1-2",
              "ticket_rule_id": "{{ .ruleId }}",
              "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-execution-1-name",
              "ticket_system_name": "test-declareticketrule-declareticket-execution-1-system-name"
            },
            "tickets": [
              {
                "_t": "declareticket",
                "a": "root",
                "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-1-name. Ticket ID: test-ticket-declareticket-execution-1. Ticket URL: https://test/test-ticket-declareticket-execution-1. Ticket name: test-ticket-declareticket-execution-1 test-resource-declareticket-execution-1-2.",
                "ticket": "test-ticket-declareticket-execution-1",
                "ticket_url": "https://test/test-ticket-declareticket-execution-1",
                "ticket_data": {
                  "name": "test-ticket-declareticket-execution-1 test-resource-declareticket-execution-1-2"
                },
                "ticket_comment": "test-comment-declareticket-execution-1-2",
                "ticket_rule_id": "{{ .ruleId }}",
                "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-execution-1-name",
                "ticket_system_name": "test-declareticketrule-declareticket-execution-1-system-name"
              }
            ],
            "connector" : "test-connector-declareticket-execution-1",
            "connector_name" : "test-connector-name-declareticket-execution-1",
            "component" :  "test-component-declareticket-execution-1",
            "resource" : "test-resource-declareticket-execution-1-2"
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
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-1-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-1-name"
      },
      {
        "_t": "declareticket",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-1-name. Ticket ID: test-ticket-declareticket-execution-1. Ticket URL: https://test/test-ticket-declareticket-execution-1. Ticket name: test-ticket-declareticket-execution-1 test-resource-declareticket-execution-1-1."
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
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-1-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-1-name"
      },
      {
        "_t": "declareticket",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-1-name. Ticket ID: test-ticket-declareticket-execution-1. Ticket URL: https://test/test-ticket-declareticket-execution-1. Ticket name: test-ticket-declareticket-execution-1 test-resource-declareticket-execution-1-2."
      }
    ]
    """

  @concurrent
  Scenario: given declare ticket rule and multiple alarms should execute one request for all alarms
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "state" : 2,
      "connector" : "test-connector-declareticket-execution-2",
      "connector_name" : "test-connector-name-declareticket-execution-2",
      "component" :  "test-component-declareticket-execution-2",
      "resource" : "test-resource-declareticket-execution-2-1",
      "source_type" : "resource"
    }
    """
    When I wait 1s
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "state" : 2,
      "connector" : "test-connector-declareticket-execution-2",
      "connector_name" : "test-connector-name-declareticket-execution-2",
      "component" :  "test-component-declareticket-execution-2",
      "resource" : "test-resource-declareticket-execution-2-2",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-execution-2&sort_by=v.resource&sort=asc
    Then the response code should be 200
    When I save response alarmId1={{ (index .lastResponse.data 0)._id }}
    When I save response alarmId2={{ (index .lastResponse.data 1)._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-declareticketrule-declareticket-execution-2-name",
      "system_name": "test-declareticketrule-declareticket-execution-2-system-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"_id\":\"test-ticket-declareticket-execution-2\",\"url\":\"https://test/test-ticket-declareticket-execution-2\",\"name\":\"test-ticket-declareticket-execution-2 {{ `{{ range .Alarms }}{{ .Value.Resource }} {{ end }}`}}\"}",
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
                "test-resource-declareticket-execution-2-1",
                "test-resource-declareticket-execution-2-2"
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
          "{{ .alarmId1 }}",
          "{{ .alarmId2 }}"
        ],
        "comment": "test-comment-declareticket-execution-2"
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
            "{{ .alarmId1 }}",
            "{{ .alarmId2 }}"
          ],
          "comment": "test-comment-declareticket-execution-2"
        }
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-execution-2&sort_by=v.resource&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "ticket": {
              "_t": "declareticket",
              "a": "root",
              "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-2-name. Ticket ID: test-ticket-declareticket-execution-2. Ticket URL: https://test/test-ticket-declareticket-execution-2. Ticket name: test-ticket-declareticket-execution-2 test-resource-declareticket-execution-2-2 test-resource-declareticket-execution-2-1 .",
              "ticket": "test-ticket-declareticket-execution-2",
              "ticket_url": "https://test/test-ticket-declareticket-execution-2",
              "ticket_data": {
                "name": "test-ticket-declareticket-execution-2 test-resource-declareticket-execution-2-2 test-resource-declareticket-execution-2-1 "
              },
              "ticket_comment": "test-comment-declareticket-execution-2",
              "ticket_rule_id": "{{ .ruleId }}",
              "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-execution-2-name",
              "ticket_system_name": "test-declareticketrule-declareticket-execution-2-system-name"
            },
            "tickets": [
              {
                "_t": "declareticket",
                "a": "root",
                "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-2-name. Ticket ID: test-ticket-declareticket-execution-2. Ticket URL: https://test/test-ticket-declareticket-execution-2. Ticket name: test-ticket-declareticket-execution-2 test-resource-declareticket-execution-2-2 test-resource-declareticket-execution-2-1 .",
                "ticket": "test-ticket-declareticket-execution-2",
                "ticket_url": "https://test/test-ticket-declareticket-execution-2",
                "ticket_data": {
                  "name": "test-ticket-declareticket-execution-2 test-resource-declareticket-execution-2-2 test-resource-declareticket-execution-2-1 "
                },
                "ticket_comment": "test-comment-declareticket-execution-2",
                "ticket_rule_id": "{{ .ruleId }}",
                "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-execution-2-name",
                "ticket_system_name": "test-declareticketrule-declareticket-execution-2-system-name"
              }
            ],
            "connector" : "test-connector-declareticket-execution-2",
            "connector_name" : "test-connector-name-declareticket-execution-2",
            "component" :  "test-component-declareticket-execution-2",
            "resource" : "test-resource-declareticket-execution-2-1"
          }
        },
        {
          "v": {
            "ticket": {
              "_t": "declareticket",
              "a": "root",
              "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-2-name. Ticket ID: test-ticket-declareticket-execution-2. Ticket URL: https://test/test-ticket-declareticket-execution-2. Ticket name: test-ticket-declareticket-execution-2 test-resource-declareticket-execution-2-2 test-resource-declareticket-execution-2-1 .",
              "ticket": "test-ticket-declareticket-execution-2",
              "ticket_url": "https://test/test-ticket-declareticket-execution-2",
              "ticket_data": {
                "name": "test-ticket-declareticket-execution-2 test-resource-declareticket-execution-2-2 test-resource-declareticket-execution-2-1 "
              },
              "ticket_comment": "test-comment-declareticket-execution-2",
              "ticket_rule_id": "{{ .ruleId }}",
              "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-execution-2-name",
              "ticket_system_name": "test-declareticketrule-declareticket-execution-2-system-name"
            },
            "tickets": [
              {
                "_t": "declareticket",
                "a": "root",
                "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-2-name. Ticket ID: test-ticket-declareticket-execution-2. Ticket URL: https://test/test-ticket-declareticket-execution-2. Ticket name: test-ticket-declareticket-execution-2 test-resource-declareticket-execution-2-2 test-resource-declareticket-execution-2-1 .",
                "ticket": "test-ticket-declareticket-execution-2",
                "ticket_url": "https://test/test-ticket-declareticket-execution-2",
                "ticket_data": {
                  "name": "test-ticket-declareticket-execution-2 test-resource-declareticket-execution-2-2 test-resource-declareticket-execution-2-1 "
                },
                "ticket_comment": "test-comment-declareticket-execution-2",
                "ticket_rule_id": "{{ .ruleId }}",
                "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-execution-2-name",
                "ticket_system_name": "test-declareticketrule-declareticket-execution-2-system-name"
              }
            ],
            "connector" : "test-connector-declareticket-execution-2",
            "connector_name" : "test-connector-name-declareticket-execution-2",
            "component" :  "test-component-declareticket-execution-2",
            "resource" : "test-resource-declareticket-execution-2-2"
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
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-2-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-2-name"
      },
      {
        "_t": "declareticket",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-2-name. Ticket ID: test-ticket-declareticket-execution-2. Ticket URL: https://test/test-ticket-declareticket-execution-2. Ticket name: test-ticket-declareticket-execution-2 test-resource-declareticket-execution-2-2 test-resource-declareticket-execution-2-1 ."
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
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-2-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-2-name"
      },
      {
        "_t": "declareticket",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-2-name. Ticket ID: test-ticket-declareticket-execution-2. Ticket URL: https://test/test-ticket-declareticket-execution-2. Ticket name: test-ticket-declareticket-execution-2 test-resource-declareticket-execution-2-2 test-resource-declareticket-execution-2-1 ."
      }
    ]
    """

  @concurrent
  Scenario: given declare ticket rule with multiple webhooks should execute all
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "state" : 2,
      "connector" : "test-connector-declareticket-execution-3",
      "connector_name" : "test-connector-name-declareticket-execution-3",
      "component" :  "test-component-declareticket-execution-3",
      "resource" : "test-resource-declareticket-execution-3",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-execution-3
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-declareticketrule-declareticket-execution-3-name",
      "system_name": "test-declareticketrule-declareticket-execution-3-system-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"name\":\"test-ticket-declareticket-execution-3 {{ `{{ range .Alarms }}{{ .Value.Resource }}{{ end }}`}}\"}",
            "method": "POST",
            "auth": {
              "username": "test",
              "password": "test"
            }
          },
          "stop_on_fail": true
        },
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"_id\":\"test-ticket-declareticket-execution-3\",\"url\":\"https://test/test-ticket-declareticket-execution-3\",\"name\":\"{{ `{{ .Response.name }}`}}\"}",
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
                "test-resource-declareticket-execution-3"
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
          "{{ .alarmId }}"
        ],
        "comment": "test-comment-declareticket-execution-3"
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
            "{{ .alarmId }}"
          ],
          "comment": "test-comment-declareticket-execution-3"
        }
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-execution-3 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "ticket": {
              "_t": "declareticket",
              "a": "root",
              "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-3-name. Ticket ID: test-ticket-declareticket-execution-3. Ticket URL: https://test/test-ticket-declareticket-execution-3. Ticket name: test-ticket-declareticket-execution-3 test-resource-declareticket-execution-3.",
              "ticket": "test-ticket-declareticket-execution-3",
              "ticket_url": "https://test/test-ticket-declareticket-execution-3",
              "ticket_data": {
                "name": "test-ticket-declareticket-execution-3 test-resource-declareticket-execution-3"
              },
              "ticket_comment": "test-comment-declareticket-execution-3",
              "ticket_rule_id": "{{ .ruleId }}",
              "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-execution-3-name",
              "ticket_system_name": "test-declareticketrule-declareticket-execution-3-system-name"
            },
            "tickets": [
              {
                "_t": "declareticket",
                "a": "root",
                "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-3-name. Ticket ID: test-ticket-declareticket-execution-3. Ticket URL: https://test/test-ticket-declareticket-execution-3. Ticket name: test-ticket-declareticket-execution-3 test-resource-declareticket-execution-3.",
                "ticket": "test-ticket-declareticket-execution-3",
                "ticket_url": "https://test/test-ticket-declareticket-execution-3",
                "ticket_data": {
                  "name": "test-ticket-declareticket-execution-3 test-resource-declareticket-execution-3"
                },
                "ticket_comment": "test-comment-declareticket-execution-3",
                "ticket_rule_id": "{{ .ruleId }}",
                "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-execution-3-name",
                "ticket_system_name": "test-declareticketrule-declareticket-execution-3-system-name"
              }
            ],
            "connector" : "test-connector-declareticket-execution-3",
            "connector_name" : "test-connector-name-declareticket-execution-3",
            "component" :  "test-component-declareticket-execution-3",
            "resource" : "test-resource-declareticket-execution-3"
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
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-3-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-3-name"
      },
      {
        "_t": "webhookstart",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-3-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-3-name"
      },
      {
        "_t": "declareticket",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-3-name. Ticket ID: test-ticket-declareticket-execution-3. Ticket URL: https://test/test-ticket-declareticket-execution-3. Ticket name: test-ticket-declareticket-execution-3 test-resource-declareticket-execution-3."
      }
    ]
    """

  @concurrent
  Scenario: given failed webhook should execute next webhook
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "state" : 2,
      "connector" : "test-connector-declareticket-execution-4",
      "connector_name" : "test-connector-name-declareticket-execution-4",
      "component" :  "test-component-declareticket-execution-4",
      "resource" : "test-resource-declareticket-execution-4",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-execution-4
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-declareticketrule-declareticket-execution-4-name",
      "system_name": "test-declareticketrule-declareticket-execution-4-system-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/not-exist",
            "payload": "{\"name\":\"test-ticket-declareticket-execution-4 {{ `{{ range .Alarms }}{{ .Value.Resource }}{{ end }}`}}\"}",
            "method": "POST",
            "auth": {
              "username": "test",
              "password": "test"
            }
          },
          "stop_on_fail": false
        },
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"_id\":\"test-ticket-declareticket-execution-4\",\"url\":\"https://test/test-ticket-declareticket-execution-4\",\"name\":\"test-ticket-declareticket-execution-4 {{ `{{ range .Alarms }}{{ .Value.Resource }}{{ end }}`}}\"}",
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
                "test-resource-declareticket-execution-4"
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
          "{{ .alarmId }}"
        ]
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
            "{{ .alarmId }}"
          ]
        }
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-execution-4 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "ticket": {
              "_t": "declareticket"
            },
            "tickets": [
              {
                "_t": "declareticket"
              }
            ],
            "connector" : "test-connector-declareticket-execution-4",
            "connector_name" : "test-connector-name-declareticket-execution-4",
            "component" :  "test-component-declareticket-execution-4",
            "resource" : "test-resource-declareticket-execution-4"
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
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-4-name"
      },
      {
        "_t": "webhookfail",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-4-name. Fail reason: url {{ .dummyApiURL }}/not-exist not found."
      },
      {
        "_t": "webhookstart",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-4-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-4-name"
      },
      {
        "_t": "declareticket",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-4-name. Ticket ID: test-ticket-declareticket-execution-4. Ticket URL: https://test/test-ticket-declareticket-execution-4. Ticket name: test-ticket-declareticket-execution-4 test-resource-declareticket-execution-4."
      }
    ]
    """

  @concurrent
  Scenario: given failed webhook shouldn't execute next webhook
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "state" : 2,
      "connector" : "test-connector-declareticket-execution-5",
      "connector_name" : "test-connector-name-declareticket-execution-5",
      "component" :  "test-component-declareticket-execution-5",
      "resource" : "test-resource-declareticket-execution-5",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-execution-5
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-declareticketrule-declareticket-execution-5-name",
      "system_name": "test-declareticketrule-declareticket-execution-5-system-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/not-exist",
            "payload": "{\"_id\":\"test-ticket-declareticket-execution-5\",\"url\":\"https://test/test-ticket-declareticket-execution-5\",\"name\":\"test-ticket-declareticket-execution-5 {{ `{{ range .Alarms }}{{ .Value.Resource }}{{ end }}`}}\"}",
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
        },
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"_id\":\"test-ticket-declareticket-execution-5\",\"url\":\"https://test/test-ticket-declareticket-execution-5\",\"name\":\"test-ticket-declareticket-execution-5 {{ `{{ range .Alarms }}{{ .Value.Resource }}{{ end }}`}}\"}",
            "method": "POST",
            "auth": {
              "username": "test",
              "password": "test"
            }
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
                "test-resource-declareticket-execution-5"
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
          "{{ .alarmId }}"
        ],
        "comment": "test-comment-declareticket-execution-5"
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
            "{{ .alarmId }}"
          ],
          "comment": "test-comment-declareticket-execution-5"
        }
      }
    ]
    """
    When I save request:
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
    When I do POST /api/v4/alarm-details until response code is 207 and response array key "0.data.steps.data" contains only:
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
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-5-name"
      },
      {
        "_t": "webhookfail",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-5-name. Fail reason: url {{ .dummyApiURL }}/not-exist not found."
      },
      {
        "_t": "declareticketfail",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-5-name. Fail reason: url {{ .dummyApiURL }}/not-exist not found."
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-execution-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "tickets": [
              {
                "_t": "declareticketfail",
                "a": "root",
                "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-5-name. Fail reason: url {{ .dummyApiURL }}/not-exist not found.",
                "ticket_comment": "test-comment-declareticket-execution-5",
                "ticket_rule_id": "{{ .ruleId }}",
                "ticket_rule_name": "Ticket declaration rule: test-declareticketrule-declareticket-execution-5-name",
                "ticket_system_name": "test-declareticketrule-declareticket-execution-5-system-name"
              }
            ],
            "connector" : "test-connector-declareticket-execution-5",
            "connector_name" : "test-connector-name-declareticket-execution-5",
            "component" :  "test-component-declareticket-execution-5",
            "resource" : "test-resource-declareticket-execution-5"
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
    Then the response key "data.0.v.ticket" should not exist
    Then the response key "data.0.v.tickets.0.ticket" should not exist
    Then the response key "data.0.v.tickets.0.ticket_url" should not exist
    Then the response key "data.0.v.tickets.0.ticket_data" should not exist

  @concurrent
  Scenario: given failed webhook should add fail reason to step message
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "state" : 2,
      "connector" : "test-connector-declareticket-execution-6",
      "connector_name" : "test-connector-name-declareticket-execution-6",
      "component" :  "test-component-declareticket-execution-6",
      "resource" : "test-resource-declareticket-execution-6",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-execution-6
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-declareticketrule-declareticket-execution-6-name",
      "system_name": "test-declareticketrule-declareticket-execution-6-system-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/auth-request",
            "method": "POST"
          },
          "stop_on_fail": false
        },
        {
          "request": {
            "url": "{{ .dummyApiURL }}/not-exist",
            "method": "POST"
          },
          "stop_on_fail": false
        },
        {
          "request": {
            "url": "http://not-exist.com",
            "method": "POST",
            "timeout": {
              "value": 1,
              "unit": "s"
            }
          },
          "stop_on_fail": false
        },
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"name\":\"{{ `{{ .Response.not_exist_field }}` }}\"}",
            "method": "POST"
          },
          "stop_on_fail": false
        },
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"_id\":\"test-ticket-declareticket-execution-6\",\"url\":\"https://test/test-ticket-declareticket-execution-6\",\"name\":\"test-ticket-declareticket-execution-6 {{ `{{ range .Alarms }}{{ .Value.Resource }}{{ end }}`}}\"}",
            "method": "POST"
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
                "test-resource-declareticket-execution-6"
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
          "{{ .alarmId }}"
        ]
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
            "{{ .alarmId }}"
          ]
        }
      }
    ]
    """
    When I save request:
    """json
    [
      {
        "_id": "{{ .alarmId }}",
        "steps": {
          "page": 1,
          "limit": 20
        }
      }
    ]
    """
    When I do POST /api/v4/alarm-details until response code is 207 and response array key "0.data.steps.data" contains only:
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
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-6-name"
      },
      {
        "_t": "webhookfail",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-6-name. Fail reason: url {{ .dummyApiURL }}/webhook/auth-request is unauthorized."
      },
      {
        "_t": "webhookstart",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-6-name"
      },
      {
        "_t": "webhookfail",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-6-name. Fail reason: url {{ .dummyApiURL }}/not-exist not found."
      },
      {
        "_t": "webhookstart",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-6-name"
      },
      {
        "_t": "webhookfail",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-6-name. Fail reason: url POST http://not-exist.com cannot be connected."
      },
      {
        "_t": "webhookstart",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-6-name"
      },
      {
        "_t": "webhookfail",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-6-name. Fail reason: invalid template in Payload."
      },
      {
        "_t": "webhookstart"
      },
      {
        "_t": "webhookcomplete"
      },
      {
        "_t": "declareticket"
      }
    ]
    """

  @concurrent
  Scenario: given incorrect declare ticket parameters should add fail reason to step message
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "state" : 2,
      "connector" : "test-connector-declareticket-execution-7",
      "connector_name" : "test-connector-name-declareticket-execution-7",
      "component" :  "test-component-declareticket-execution-7",
      "resource" : "test-resource-declareticket-execution-7",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-execution-7
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-declareticketrule-declareticket-execution-7-name",
      "system_name": "test-declareticketrule-declareticket-execution-7-system-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"_id\":\"test-ticket-declareticket-execution-7\",\"url\":\"https://test/test-ticket-declareticket-execution-7\",\"name\":\"test-ticket-declareticket-execution-7 {{ `{{ range .Alarms }}{{ .Value.Resource }}{{ end }}`}}\"}",
            "method": "POST"
          },
          "declare_ticket": {
            "ticket_id": "not_exist_field",
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
                "test-resource-declareticket-execution-7"
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
          "{{ .alarmId }}"
        ]
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
            "{{ .alarmId }}"
          ]
        }
      }
    ]
    """
    When I save request:
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
    When I do POST /api/v4/alarm-details until response code is 207 and response array key "0.data.steps.data" contains only:
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
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-7-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-7-name"
      },
      {
        "_t": "declareticketfail",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-7-name. Fail reason: ticket_id is emtpy, response has nothing in not_exist_field."
      }
    ]
    """

  @concurrent
  Scenario: given not JSON response for declare ticket webhook should add fail reason to step message
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "state" : 2,
      "connector" : "test-connector-declareticket-execution-8",
      "connector_name" : "test-connector-name-declareticket-execution-8",
      "component" :  "test-component-declareticket-execution-8",
      "resource" : "test-resource-declareticket-execution-8",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-execution-8
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-declareticketrule-declareticket-execution-8-name",
      "system_name": "test-declareticketrule-declareticket-execution-8-system-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "not JSON",
            "method": "POST"
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
                "test-resource-declareticket-execution-8"
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
          "{{ .alarmId }}"
        ]
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
            "{{ .alarmId }}"
          ]
        }
      }
    ]
    """
    When I save request:
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
    When I do POST /api/v4/alarm-details until response code is 207 and response array key "0.data.steps.data" contains only:
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
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-8-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-8-name"
      },
      {
        "_t": "declareticketfail",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-8-name. Fail reason: response of POST {{ .dummyApiURL }}/webhook/request is not valid JSON."
      }
    ]
    """

  @concurrent
  Scenario: given not JSON response for not declare ticket webhook should return ok
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "state" : 2,
      "connector" : "test-connector-declareticket-execution-9",
      "connector_name" : "test-connector-name-declareticket-execution-9",
      "component" :  "test-component-declareticket-execution-9",
      "resource" : "test-resource-declareticket-execution-9",
      "source_type" : "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-execution-9
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-declareticketrule-declareticket-execution-9-name",
      "system_name": "test-declareticketrule-declareticket-execution-9-system-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "not JSON",
            "method": "POST"
          },
          "stop_on_fail": true
        },
        {
          "request": {
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"_id\":\"test-ticket-declareticket-execution-9\",\"url\":\"https://test/test-ticket-declareticket-execution-9\",\"name\":\"test-ticket-declareticket-execution-9 {{ `{{ range .Alarms }}{{ .Value.Resource }}{{ end }}`}}\"}",
            "method": "POST"
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
                "test-resource-declareticket-execution-9"
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
          "{{ .alarmId }}"
        ]
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
            "{{ .alarmId }}"
          ]
        }
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-declareticket-execution-9 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "ticket": {
              "_t": "declareticket"
            },
            "connector" : "test-connector-declareticket-execution-9",
            "connector_name" : "test-connector-name-declareticket-execution-9",
            "component" :  "test-component-declareticket-execution-9",
            "resource" : "test-resource-declareticket-execution-9"
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
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-9-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-9-name"
      },
      {
        "_t": "webhookstart",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-9-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-9-name"
      },
      {
        "_t": "declareticket",
        "a": "root",
        "user_id": "root",
        "m": "Ticket declaration rule: test-declareticketrule-declareticket-execution-9-name. Ticket ID: test-ticket-declareticket-execution-9. Ticket URL: https://test/test-ticket-declareticket-execution-9. Ticket name: test-ticket-declareticket-execution-9 test-resource-declareticket-execution-9."
      }
    ]
    """
