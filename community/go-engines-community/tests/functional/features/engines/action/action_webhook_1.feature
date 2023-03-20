Feature: execute action on trigger
  I need to be able to trigger action on event

  @concurrent
  Scenario: given scenario and check event should update alarm
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-1-name",
      "priority": 10008,
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
                  "value": "test-component-action-webhook-1"
                }
              }
            ]
          ],
          "type": "webhook",
          "comment": "test-scenario-action-webhook-1-comment",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/auth-request",
              "auth": {
                "username": "test",
                "password": "test"
              },
              "headers": {"Content-Type": "application/json"},
              "payload": "{\"_id\":\"test-ticket-action-webhook-1\",\"url\":\"https://test/test-ticket-action-webhook-1\",\"name\":\"test-ticket-action-webhook-1 {{ `{{ .Entity.Name }}`}}\"}"
            },
            "ticket_system_name": "test-scenario-action-webhook-1-system-name",
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "_id",
              "ticket_url": "url",
              "name": "name"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    Then I save response scenarioId={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-action-webhook-1",
      "connector_name" : "test-connector-name-action-webhook-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-webhook-1",
      "resource" : "test-resource-action-webhook-1",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-action-webhook-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "tickets": [
              {
                "_t": "declareticket",
                "a": "system",
                "m": "Scenario: test-scenario-action-webhook-1-name. Ticket ID: test-ticket-action-webhook-1. Ticket URL: https://test/test-ticket-action-webhook-1. Ticket name: test-ticket-action-webhook-1 test-resource-action-webhook-1.",
                "ticket": "test-ticket-action-webhook-1",
                "ticket_url": "https://test/test-ticket-action-webhook-1",
                "ticket_data": {
                  "name": "test-ticket-action-webhook-1 test-resource-action-webhook-1"
                },
                "ticket_rule_id": "{{ .scenarioId }}",
                "ticket_rule_name": "Scenario: test-scenario-action-webhook-1-name",
                "ticket_system_name": "test-scenario-action-webhook-1-system-name",
                "ticket_comment": "test-scenario-action-webhook-1-comment"
              }
            ],
            "ticket": {
              "_t": "declareticket",
              "a": "system",
              "m": "Scenario: test-scenario-action-webhook-1-name. Ticket ID: test-ticket-action-webhook-1. Ticket URL: https://test/test-ticket-action-webhook-1. Ticket name: test-ticket-action-webhook-1 test-resource-action-webhook-1.",
              "ticket": "test-ticket-action-webhook-1",
              "ticket_url": "https://test/test-ticket-action-webhook-1",
              "ticket_data": {
                "name": "test-ticket-action-webhook-1 test-resource-action-webhook-1"
              },
              "ticket_rule_id": "{{ .scenarioId }}",
              "ticket_rule_name": "Scenario: test-scenario-action-webhook-1-name",
              "ticket_system_name": "test-scenario-action-webhook-1-system-name",
              "ticket_comment": "test-scenario-action-webhook-1-comment"
            },
            "connector": "test-connector-action-webhook-1",
            "connector_name": "test-connector-name-action-webhook-1",
            "component": "test-component-action-webhook-1",
            "resource": "test-resource-action-webhook-1"
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
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-1-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-1-name"
      },
      {
        "_t": "declareticket",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-1-name. Ticket ID: test-ticket-action-webhook-1. Ticket URL: https://test/test-ticket-action-webhook-1. Ticket name: test-ticket-action-webhook-1 test-resource-action-webhook-1."
      }
    ]
    """

  @concurrent
  Scenario: given scenario and check event should emit declare ticket trigger
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-2-1-name",
      "priority": 10010,
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
                  "value": "test-component-action-webhook-2"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request",
              "payload": "{\"_id\":\"test-ticket-action-webhook-2\",\"url\":\"https://test/test-ticket-action-webhook-2\",\"name\":\"test-ticket-action-webhook-2 {{ `{{ .Entity.Name }}`}}\"}"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "_id",
              "name": "name"
            }
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
      "name": "test-scenario-action-webhook-2-2-name",
      "priority": 10012,
      "enabled": true,
      "triggers": ["declareticketwebhook"],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-component-action-webhook-2"
                }
              }
            ]
          ],
          "type": "ack",
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
      "connector" : "test-connector-action-webhook-2",
      "connector_name" : "test-connector-name-action-webhook-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-webhook-2",
      "resource" : "test-resource-action-webhook-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-action-webhook-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "ticket": {
              "_t": "declareticket"
            },
            "ack": {
              "_t": "ack"
            },
            "connector": "test-connector-action-webhook-2",
            "connector_name": "test-connector-name-action-webhook-2",
            "component": "test-component-action-webhook-2",
            "resource": "test-resource-action-webhook-2"
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
        "a": "system",
        "user_id": ""
      },
      {
        "_t": "webhookcomplete",
        "a": "system",
        "user_id": ""
      },
      {
        "_t": "declareticket",
        "a": "system",
        "user_id": ""
      },
      {
        "_t": "ack",
        "a": "system",
        "user_id": ""
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
  Scenario: given webhook scenario to test response and header templates
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-4-name",
      "priority": 10015,
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
                  "value": "test-component-action-webhook-4"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request",
              "payload": "{\"url\":\"test-ticket-action-webhook-4-1-url\"}",
              "headers": {
                "X-Go-Test": "test-ticket-action-webhook-4-1-header"
              }
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
                  "value": "test-component-action-webhook-4"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request",
              "payload": "{\"_id\":\"test-ticket-action-webhook-4\",\"url\":\"test-ticket-action-webhook-4-2 {{ `{{ .Response.url }}` }}\",\"header\":\"test-ticket-action-webhook-4-2-header {{ `{{ index .Header \"X-Go-Test\" }}` }}\"}"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "_id",
              "ticket_url": "url",
              "header": "header"
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
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-action-webhook-4",
      "connector_name" : "test-connector-name-action-webhook-4",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-webhook-4",
      "resource" : "test-resource-action-webhook-4",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-action-webhook-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "tickets": [
              {
                "_t": "declareticket",
                "a": "system",
                "ticket": "test-ticket-action-webhook-4",
                "ticket_url": "test-ticket-action-webhook-4-2 test-ticket-action-webhook-4-1-url",
                "ticket_data": {
                  "header": "test-ticket-action-webhook-4-2-header test-ticket-action-webhook-4-1-header"
                }
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

  @concurrent
  Scenario: given scenarios with 2 actions and webhook should be able to use additional data in template
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-6-1-name",
      "priority": 10019,
      "enabled": true,
      "triggers": [
        "create"
      ],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-component-action-webhook-6"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "forward_author": true,
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request",
              "payload": "{\"_id\":\"test-ticket-action-webhook-6-1-action-1\",\"url\":\"https://test/test-ticket-action-webhook-6-1-action-1\",\"name\":\"test-scenario-action-webhook-6-1-action-1 {{`[{{ .AdditionalData.Trigger }}] [{{ .AdditionalData.Author }}] [{{ .AdditionalData.Initiator }}] [{{ .AdditionalData.User }}]`}}\"}"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "_id",
              "ticket_url": "url",
              "name": "name"
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
                  "value": "test-component-action-webhook-6"
                }
              }
            ]
          ],
          "drop_scenario_if_not_matched": false,
          "emit_trigger": true,
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request",
              "payload": "{\"_id\":\"test-ticket-action-webhook-6-1-action-2\",\"url\":\"https://test/test-ticket-action-webhook-6-1-action-2\",\"name\":\"test-scenario-action-webhook-6-1-action-2 {{`[{{ .AdditionalData.Trigger }}] [{{ .AdditionalData.Author }}] [{{ .AdditionalData.Initiator }}] [{{ .AdditionalData.User }}]`}}\"}"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "_id",
              "ticket_url": "url",
              "name": "name"
            }
          }
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-6-2-name",
      "priority": 10022,
      "enabled": true,
      "triggers": [
        "declareticketwebhook"
      ],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-component-action-webhook-6"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "author": "test-scenario-action-webhook-6-2-action-1-author",
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request",
              "payload": "{\"_id\":\"test-ticket-action-webhook-6-2-action-1\",\"url\":\"https://test/test-ticket-action-webhook-6-2-action-1\",\"name\":\"test-scenario-action-webhook-6-2-action-1 {{`[{{ .AdditionalData.Trigger }}] [{{ .AdditionalData.Author }}] [{{ .AdditionalData.Initiator }}] [{{ .AdditionalData.User }}]`}}\"}"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "_id",
              "ticket_url": "url",
              "name": "name"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-6-3-name",
      "priority": 10024,
      "enabled": true,
      "triggers": [
        "create"
      ],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-component-action-webhook-6"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "author": "test-scenario-action-webhook-6-3-action-1-author",
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request",
              "payload": "{\"_id\":\"test-ticket-action-webhook-6-3-action-1\",\"url\":\"https://test/test-ticket-action-webhook-6-3-action-1\",\"name\":\"{{ `{{ $testVar := .AdditionalData.Output }}test-scenario-action-webhook-6-3-action-1 [{{$testVar}}]` }}\"}"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "_id",
              "ticket_url": "url",
              "name": "name"
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
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-action-webhook-6",
      "connector_name" : "test-connector-name-action-webhook-6",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-webhook-6",
      "resource" : "test-resource-action-webhook-6",
      "state" : 2,
      "output" : "noveo alarm",
      "initiator": "user"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-action-webhook-6
    Then the response code should be 200
    Then the response array key "data.0.v.tickets" should contain:
    """json
    [
      {
        "_t": "declareticket",
        "ticket_data": {
          "name": "test-scenario-action-webhook-6-1-action-1 [create] [root] [user] [root]"
        }
      },
      {
        "_t": "declareticket",
        "ticket_data": {
          "name": "test-scenario-action-webhook-6-1-action-2 [create] [system] [user] []"
        }
      },
      {
        "_t": "declareticket",
        "ticket_data": {
          "name": "test-scenario-action-webhook-6-2-action-1 [declareticketwebhook] [test-scenario-action-webhook-6-2-action-1-author] [user] []"
        }
      },
      {
        "_t": "declareticket",
        "ticket_data": {
          "name": "test-scenario-action-webhook-6-3-action-1 [noveo alarm]"
        }
      }
    ]
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1,
          "limit": 20
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
        "m": "Scenario: test-scenario-action-webhook-6-1-name"
      },
      {
        "_t": "webhookstart",
        "a": "test-scenario-action-webhook-6-3-action-1-author",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-6-3-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Scenario: test-scenario-action-webhook-6-1-name"
      },
      {
        "_t": "declareticket",
        "a": "root",
        "user_id": "root",
        "m": "Scenario: test-scenario-action-webhook-6-1-name. Ticket ID: test-ticket-action-webhook-6-1-action-1. Ticket URL: https://test/test-ticket-action-webhook-6-1-action-1. Ticket name: test-scenario-action-webhook-6-1-action-1 [create] [root] [user] [root]."
      },
      {
        "_t": "webhookcomplete",
        "a": "test-scenario-action-webhook-6-3-action-1-author",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-6-3-name"
      },
      {
        "_t": "declareticket",
        "a": "test-scenario-action-webhook-6-3-action-1-author",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-6-3-name. Ticket ID: test-ticket-action-webhook-6-3-action-1. Ticket URL: https://test/test-ticket-action-webhook-6-3-action-1. Ticket name: test-scenario-action-webhook-6-3-action-1 [noveo alarm]."
      },
      {
        "_t": "webhookstart",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-6-1-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-6-1-name"
      },
      {
        "_t": "declareticket",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-6-1-name. Ticket ID: test-ticket-action-webhook-6-1-action-2. Ticket URL: https://test/test-ticket-action-webhook-6-1-action-2. Ticket name: test-scenario-action-webhook-6-1-action-2 [create] [system] [user] []."
      },
      {
        "_t": "webhookstart",
        "a": "test-scenario-action-webhook-6-2-action-1-author",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-6-2-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "test-scenario-action-webhook-6-2-action-1-author",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-6-2-name"
      },
      {
        "_t": "declareticket",
        "a": "test-scenario-action-webhook-6-2-action-1-author",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-6-2-name. Ticket ID: test-ticket-action-webhook-6-2-action-1. Ticket URL: https://test/test-ticket-action-webhook-6-2-action-1. Ticket name: test-scenario-action-webhook-6-2-action-1 [declareticketwebhook] [test-scenario-action-webhook-6-2-action-1-author] [user] []."
      }
    ]
    """
    Then the response array key "0.data.steps.data" should contain in order:
    """json
    [
      {
        "_t": "declareticket",
        "m": "Scenario: test-scenario-action-webhook-6-1-name. Ticket ID: test-ticket-action-webhook-6-1-action-1. Ticket URL: https://test/test-ticket-action-webhook-6-1-action-1. Ticket name: test-scenario-action-webhook-6-1-action-1 [create] [root] [user] [root]."
      },
      {
        "_t": "declareticket",
        "m": "Scenario: test-scenario-action-webhook-6-1-name. Ticket ID: test-ticket-action-webhook-6-1-action-2. Ticket URL: https://test/test-ticket-action-webhook-6-1-action-2. Ticket name: test-scenario-action-webhook-6-1-action-2 [create] [system] [user] []."
      },
      {
        "_t": "declareticket",
        "m": "Scenario: test-scenario-action-webhook-6-2-name. Ticket ID: test-ticket-action-webhook-6-2-action-1. Ticket URL: https://test/test-ticket-action-webhook-6-2-action-1. Ticket name: test-scenario-action-webhook-6-2-action-1 [declareticketwebhook] [test-scenario-action-webhook-6-2-action-1-author] [user] []."
      }
    ]
    """

  @concurrent
  Scenario: given webhook scenario to test multiple response templates
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-7-name",
      "priority": 10026,
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
                  "value": "test-component-action-webhook-7"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-action-webhook-7/test-component-action-webhook-7"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request",
              "payload": "{\"_id\":\"test-ticket-action-webhook-7-1\"}"
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
                  "value": "test-component-action-webhook-7"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-action-webhook-7/test-component-action-webhook-7"
                }
              }
            ]
          ],
          "type": "ack",
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
                  "value": "test-component-action-webhook-7"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-action-webhook-7/test-component-action-webhook-7"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request",
              "payload": "{\"_id\":\"test-ticket-action-webhook-7-2\"}"
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
                  "value": "test-component-action-webhook-7"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-action-webhook-7/test-component-action-webhook-7"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request",
              "payload": "{\"_id\":\"test-ticket-action-webhook-7\",\"name1\":\"{{ `{{index .ResponseMap \"0._id\"}}`}}\",\"name2\":\"{{ `{{index .ResponseMap \"1._id\"}}`}}\",\"name3\":\"{{ `{{ .Response._id }}`}}\"}"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "_id",
              "name1": "name1",
              "name2": "name2",
              "name3": "name3"
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
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-action-webhook-7",
      "connector_name" : "test-connector-name-action-webhook-7",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-webhook-7",
      "resource" : "test-resource-action-webhook-7",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-action-webhook-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "tickets": [
              {
                "_t": "declareticket",
                "a": "system",
                "ticket_data": {
                  "name1": "test-ticket-action-webhook-7-1",
                  "name2": "test-ticket-action-webhook-7-2",
                  "name3": "test-ticket-action-webhook-7-2"
                }
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

  @concurrent
  Scenario: given webhook scenario to test document with a document with arrays in response
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-8-name",
      "priority": 10028,
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
                  "value": "test-component-action-webhook-8"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-action-webhook-8/test-component-action-webhook-8"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request",
              "payload": "{\"array\":[{\"elem1\":\"test1\",\"elem2\":\"test2\"},{\"elem1\":\"test3\",\"elem2\":\"test4\"}]}"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "array.1.elem1",
              "test_val": "array.0.elem1"
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
                  "value": "test-component-action-webhook-8"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-action-webhook-8/test-component-action-webhook-8"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request",
              "payload": "{\"_id\":\"test-ticket-action-webhook-8\",\"name1\":\"{{ `{{index .ResponseMap \"0.array.0.elem1\"}}` }}\",\"name2\":\"{{ `{{index .ResponseMap \"0.array.1.elem2\"}}` }}\"}"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "_id",
              "name1": "name1",
              "name2": "name2"
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
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-action-webhook-8",
      "connector_name" : "test-connector-name-action-webhook-8",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-webhook-8",
      "resource" : "test-resource-action-webhook-8",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-action-webhook-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "tickets": [
              {
                "_t": "declareticket",
                "a": "system",
                "ticket": "test3",
                "ticket_data": {
                  "test_val": "test1"
                }
              },
              {
                "_t": "declareticket",
                "a": "system",
                "ticket_data": {
                  "name1": "test1",
                  "name2": "test4"
                }
              }
            ],
            "connector": "test-connector-action-webhook-8",
            "connector_name": "test-connector-name-action-webhook-8",
            "component": "test-component-action-webhook-8",
            "resource": "test-resource-action-webhook-8"
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

  @concurrent
  Scenario: given webhook scenario where the webhook response is an array
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-9-name",
      "priority": 10030,
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
                  "value": "test-component-action-webhook-9"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-action-webhook-9/test-component-action-webhook-9"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request",
              "payload": "[{\"elem1\":\"test1\",\"elem2\":\"test2\"},{\"elem1\":\"test3\",\"elem2\":\"test4\"}]"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "1.elem1",
              "test_val": "0.elem1"
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
                  "value": "test-component-action-webhook-9"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-action-webhook-9/test-component-action-webhook-9"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request",
              "payload": "{\"_id\":\"test-ticket-action-webhook-9\",\"name1\":\"{{ `{{index .ResponseMap \"0.0.elem1\"}}` }}\",\"name2\":\"{{ `{{index .ResponseMap \"0.1.elem2\"}}` }}\"}"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "_id",
              "name1": "name1",
              "name2": "name2"
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
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-action-webhook-9",
      "connector_name" : "test-connector-name-action-webhook-9",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-webhook-9",
      "resource" : "test-resource-action-webhook-9",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-action-webhook-9
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "tickets": [
              {
                "_t": "declareticket",
                "a": "system",
                "ticket": "test3",
                "ticket_data": {
                  "test_val": "test1"
                }
              },
              {
                "_t": "declareticket",
                "a": "system",
                "ticket_data": {
                  "name1": "test1",
                  "name2": "test4"
                }
              }
            ],
            "connector": "test-connector-action-webhook-9",
            "connector_name": "test-connector-name-action-webhook-9",
            "component": "test-component-action-webhook-9",
            "resource": "test-resource-action-webhook-9"
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

  @concurrent
  Scenario: given scenario should use author from parameters
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-3-name",
      "priority": 10014,
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
                  "value": "test-component-action-webhook-3"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "author": "test-scenario-action-webhook-3-action-1-author {{ `{{ .Alarm.Value.Resource }}` }}",
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request",
              "payload": "{\"_id\":\"test-ticket-action-webhook-3\",\"url\":\"https://test/test-ticket-action-webhook-3\",\"name\":\"test-ticket-action-webhook-3 {{ `{{ .Entity.Name }}`}}\"}"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "_id"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": true
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-action-webhook-3",
      "connector_name" : "test-connector-name-action-webhook-3",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-webhook-3",
      "resource" : "test-resource-action-webhook-3",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-action-webhook-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "tickets": [
              {
                "_t": "declareticket",
                "a": "test-scenario-action-webhook-3-action-1-author test-resource-action-webhook-3"
              }
            ],
            "ticket": {
              "_t": "declareticket",
              "a": "test-scenario-action-webhook-3-action-1-author test-resource-action-webhook-3"
            },
            "connector": "test-connector-action-webhook-3",
            "connector_name": "test-connector-name-action-webhook-3",
            "component": "test-component-action-webhook-3",
            "resource": "test-resource-action-webhook-3"
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
        "a": "test-scenario-action-webhook-3-action-1-author test-resource-action-webhook-3",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-3-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "test-scenario-action-webhook-3-action-1-author test-resource-action-webhook-3",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-3-name"
      },
      {
        "_t": "declareticket",
        "a": "test-scenario-action-webhook-3-action-1-author test-resource-action-webhook-3",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-3-name. Ticket ID: test-ticket-action-webhook-3."
      }
    ]
    """

  @concurrent
  Scenario: given scenario should forward author from trigger event
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-5-name",
      "priority": 10011,
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
                  "value": "test-component-action-webhook-5"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "forward_author": true,
            "author": "test-scenario-action-webhook-5-action-1-author {{ `{{ .Alarm.Value.Resource }}` }}",
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request",
              "payload": "{\"_id\":\"test-ticket-action-webhook-5\",\"url\":\"https://test/test-ticket-action-webhook-5\",\"name\":\"test-ticket-action-webhook-5 {{ `{{ .Entity.Name }}`}}\"}"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "_id"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": true
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-action-webhook-5",
      "connector_name" : "test-connector-name-action-webhook-5",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-webhook-5",
      "resource" : "test-resource-action-webhook-5",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-action-webhook-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "tickets": [
              {
                "_t": "declareticket",
                "a": "root"
              }
            ],
            "ticket": {
              "_t": "declareticket",
              "a": "root"
            },
            "connector": "test-connector-action-webhook-5",
            "connector_name": "test-connector-name-action-webhook-5",
            "component": "test-component-action-webhook-5",
            "resource": "test-resource-action-webhook-5"
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
        "m": "Scenario: test-scenario-action-webhook-5-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Scenario: test-scenario-action-webhook-5-name"
      },
      {
        "_t": "declareticket",
        "a": "root",
        "user_id": "root",
        "m": "Scenario: test-scenario-action-webhook-5-name. Ticket ID: test-ticket-action-webhook-5."
      }
    ]
    """
