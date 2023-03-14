Feature: execute action on trigger
  I need to be able to trigger action on event

  @concurrent
  Scenario: given scenario with failed declare ticket should add steps
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-second-1-1-name",
      "priority": 10091,
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
                  "value": "test-component-action-webhook-second-1"
                }
              }
            ]
          ],
          "type": "webhook",
          "comment": "test-scenario-action-webhook-second-1-1-comment",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/auth-request"
            },
            "ticket_system_name": "test-scenario-action-webhook-second-1-1-system-name",
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
    Then I save response scenarioId1={{ .lastResponse._id }}
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-second-1-2-name",
      "priority": 10092,
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
                  "value": "test-component-action-webhook-second-1"
                }
              }
            ]
          ],
          "type": "webhook",
          "comment": "test-scenario-action-webhook-second-1-2-comment",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "http://not-exist.com",
              "timeout": {
                "value": 1,
                "unit": "s"
              }
            },
            "ticket_system_name": "test-scenario-action-webhook-second-1-2-system-name",
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
    Then I save response scenarioId2={{ .lastResponse._id }}
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-second-1-3-name",
      "priority": 10094,
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
                  "value": "test-component-action-webhook-second-1"
                }
              }
            ]
          ],
          "type": "webhook",
          "comment": "test-scenario-action-webhook-second-1-3-comment",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request",
              "payload": "{\"name\":\"test-scenario-action-webhook-second-1-3\"}"
            },
            "ticket_system_name": "test-scenario-action-webhook-second-1-3-system-name",
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
    Then I save response scenarioId3={{ .lastResponse._id }}
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-second-1-4-name",
      "priority": 10095,
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
                  "value": "test-component-action-webhook-second-1"
                }
              }
            ]
          ],
          "type": "webhook",
          "comment": "test-scenario-action-webhook-second-1-4-comment",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request",
              "payload": "test-scenario-action-webhook-second-1-4"
            },
            "ticket_system_name": "test-scenario-action-webhook-second-1-4-system-name",
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
    Then I save response scenarioId4={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-action-webhook-second-1",
      "connector_name" : "test-connector-name-action-webhook-second-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-webhook-second-1",
      "resource" : "test-resource-action-webhook-second-1",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-action-webhook-second-1
    Then the response code should be 200
    Then the response array key "data.0.v.tickets" should contain only:
    """json
    [
      {
        "_t": "declareticketfail",
        "a": "system",
        "m": "Scenario: test-scenario-action-webhook-second-1-1-name. Fail reason: url {{ .dummyApiURL }}/webhook/auth-request is unauthorized.",
        "ticket_rule_id": "{{ .scenarioId1 }}",
        "ticket_rule_name": "Scenario: test-scenario-action-webhook-second-1-1-name",
        "ticket_system_name": "test-scenario-action-webhook-second-1-1-system-name",
        "ticket_comment": "test-scenario-action-webhook-second-1-1-comment"
      },
      {
        "_t": "declareticketfail",
        "a": "system",
        "m": "Scenario: test-scenario-action-webhook-second-1-2-name. Fail reason: url POST http://not-exist.com cannot be connected.",
        "ticket_rule_id": "{{ .scenarioId2 }}",
        "ticket_rule_name": "Scenario: test-scenario-action-webhook-second-1-2-name",
        "ticket_system_name": "test-scenario-action-webhook-second-1-2-system-name",
        "ticket_comment": "test-scenario-action-webhook-second-1-2-comment"
      },
      {
        "_t": "declareticketfail",
        "a": "system",
        "m": "Scenario: test-scenario-action-webhook-second-1-3-name. Fail reason: ticket_id is emtpy, response has nothing in _id.",
        "ticket_rule_id": "{{ .scenarioId3 }}",
        "ticket_rule_name": "Scenario: test-scenario-action-webhook-second-1-3-name",
        "ticket_system_name": "test-scenario-action-webhook-second-1-3-system-name",
        "ticket_comment": "test-scenario-action-webhook-second-1-3-comment"
      },
      {
        "_t": "declareticketfail",
        "a": "system",
        "m": "Scenario: test-scenario-action-webhook-second-1-4-name. Fail reason: response of POST {{ .dummyApiURL }}/webhook/request is not valid JSON.",
        "ticket_rule_id": "{{ .scenarioId4 }}",
        "ticket_rule_name": "Scenario: test-scenario-action-webhook-second-1-4-name",
        "ticket_system_name": "test-scenario-action-webhook-second-1-4-system-name",
        "ticket_comment": "test-scenario-action-webhook-second-1-4-comment"
      }
    ]
    """
    Then the response key "data.0.v.ticket" should not exist
    Then the response key "data.0.v.tickets.0.ticket" should not exist
    Then the response key "data.0.v.tickets.1.ticket" should not exist
    Then the response key "data.0.v.tickets.2.ticket" should not exist
    Then the response key "data.0.v.tickets.3.ticket" should not exist
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
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-second-1-1-name"
      },
      {
        "_t": "webhookfail",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-second-1-1-name. Fail reason: url {{ .dummyApiURL }}/webhook/auth-request is unauthorized."
      },
      {
        "_t": "declareticketfail",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-second-1-1-name. Fail reason: url {{ .dummyApiURL }}/webhook/auth-request is unauthorized."
      },
      {
        "_t": "webhookstart",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-second-1-2-name"
      },
      {
        "_t": "webhookfail",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-second-1-2-name. Fail reason: url POST http://not-exist.com cannot be connected."
      },
      {
        "_t": "declareticketfail",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-second-1-2-name. Fail reason: url POST http://not-exist.com cannot be connected."
      },
      {
        "_t": "webhookstart",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-second-1-3-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-second-1-3-name"
      },
      {
        "_t": "declareticketfail",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-second-1-3-name. Fail reason: ticket_id is emtpy, response has nothing in _id."
      },
      {
        "_t": "webhookstart",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-second-1-4-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-second-1-4-name"
      },
      {
        "_t": "declareticketfail",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-second-1-4-name. Fail reason: response of POST {{ .dummyApiURL }}/webhook/request is not valid JSON."
      }
    ]
    """

  @concurrent
  Scenario: given scenario with failed webhook should add steps
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-second-2-1-name",
      "priority": 10096,
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
                  "value": "test-component-action-webhook-second-2"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/auth-request"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    Then I save response scenarioId1={{ .lastResponse._id }}
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-second-2-2-name",
      "priority": 10097,
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
                  "value": "test-component-action-webhook-second-2"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "http://not-exist.com",
              "timeout": {
                "value": 1,
                "unit": "s"
              }
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    Then I save response scenarioId2={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-action-webhook-second-2",
      "connector_name" : "test-connector-name-action-webhook-second-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-webhook-second-2",
      "resource" : "test-resource-action-webhook-second-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-action-webhook-second-2
    Then the response code should be 200
    Then the response key "data.0.v.ticket" should not exist
    Then the response key "data.0.v.tickets" should not exist
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
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-second-2-1-name"
      },
      {
        "_t": "webhookfail",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-second-2-1-name. Fail reason: url {{ .dummyApiURL }}/webhook/auth-request is unauthorized."
      },
      {
        "_t": "webhookstart",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-second-2-2-name"
      },
      {
        "_t": "webhookfail",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-second-2-2-name. Fail reason: url POST http://not-exist.com cannot be connected."
      }
    ]
    """

  @concurrent
  Scenario: given scenario with emtpy response webhook should add steps
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-second-3-name",
      "priority": 10098,
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
                  "value": "test-component-action-webhook-second-3"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request"
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
      "connector" : "test-connector-action-webhook-second-3",
      "connector_name" : "test-connector-name-action-webhook-second-3",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-webhook-second-3",
      "resource" : "test-resource-action-webhook-second-3",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-action-webhook-second-3
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
        "m": "Scenario: test-scenario-action-webhook-second-3-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-second-3-name"
      }
    ]
    """

  @concurrent
  Scenario: given scenario with emtpy response webhook should add ticket
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-second-4-name",
      "priority": 10099,
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
                  "value": "test-component-action-webhook-second-4"
                }
              }
            ]
          ],
          "type": "webhook",
          "comment": "test-scenario-action-webhook-second-4-comment",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request"
            },
            "ticket_system_name": "test-scenario-action-webhook-second-4-system-name",
            "declare_ticket": {
              "empty_response": true
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
      "connector" : "test-connector-action-webhook-second-4",
      "connector_name" : "test-connector-name-action-webhook-second-4",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-webhook-second-4",
      "resource" : "test-resource-action-webhook-second-4",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-action-webhook-second-4
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
                "m": "Scenario: test-scenario-action-webhook-second-4-name. Ticket ID: N/A.",
                "ticket": "N/A",
                "ticket_rule_id": "{{ .scenarioId }}",
                "ticket_rule_name": "Scenario: test-scenario-action-webhook-second-4-name",
                "ticket_system_name": "test-scenario-action-webhook-second-4-system-name",
                "ticket_comment": "test-scenario-action-webhook-second-4-comment"
              }
            ],
            "ticket": {
              "_t": "declareticket",
              "a": "system",
              "m": "Scenario: test-scenario-action-webhook-second-4-name. Ticket ID: N/A.",
              "ticket": "N/A",
              "ticket_rule_id": "{{ .scenarioId }}",
              "ticket_rule_name": "Scenario: test-scenario-action-webhook-second-4-name",
              "ticket_system_name": "test-scenario-action-webhook-second-4-system-name",
              "ticket_comment": "test-scenario-action-webhook-second-4-comment"
            },
            "connector": "test-connector-action-webhook-second-4",
            "connector_name": "test-connector-name-action-webhook-second-4",
            "component": "test-component-action-webhook-second-4",
            "resource": "test-resource-action-webhook-second-4"
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
        "m": "Scenario: test-scenario-action-webhook-second-4-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-second-4-name"
      },
      {
        "_t": "declareticket",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-second-4-name. Ticket ID: N/A."
      }
    ]
    """

  @concurrent
  Scenario: given scenario with empty ticket_id add ticket
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-second-5-name",
      "priority": 10100,
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
                  "value": "test-component-action-webhook-second-5"
                }
              }
            ]
          ],
          "type": "webhook",
          "comment": "test-scenario-action-webhook-second-5-comment",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/auth-request",
              "auth": {
                "username": "test",
                "password": "test"
              },
              "payload": "{}"
            },
            "ticket_system_name": "test-scenario-action-webhook-second-5-system-name",
            "declare_ticket": {}
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
      "connector" : "test-connector-action-webhook-second-5",
      "connector_name" : "test-connector-name-action-webhook-second-5",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-action-webhook-second-5",
      "resource" : "test-resource-action-webhook-second-5",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-action-webhook-second-5
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
                "m": "Scenario: test-scenario-action-webhook-second-5-name. Ticket ID: N/A.",
                "ticket": "N/A",
                "ticket_rule_id": "{{ .scenarioId }}",
                "ticket_rule_name": "Scenario: test-scenario-action-webhook-second-5-name",
                "ticket_system_name": "test-scenario-action-webhook-second-5-system-name",
                "ticket_comment": "test-scenario-action-webhook-second-5-comment"
              }
            ],
            "ticket": {
              "_t": "declareticket",
              "a": "system",
              "m": "Scenario: test-scenario-action-webhook-second-5-name. Ticket ID: N/A.",
              "ticket": "N/A",
              "ticket_rule_id": "{{ .scenarioId }}",
              "ticket_rule_name": "Scenario: test-scenario-action-webhook-second-5-name",
              "ticket_system_name": "test-scenario-action-webhook-second-5-system-name",
              "ticket_comment": "test-scenario-action-webhook-second-5-comment"
            },
            "connector": "test-connector-action-webhook-second-5",
            "connector_name": "test-connector-name-action-webhook-second-5",
            "component": "test-component-action-webhook-second-5",
            "resource": "test-resource-action-webhook-second-5"
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
        "m": "Scenario: test-scenario-action-webhook-second-5-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-second-5-name"
      },
      {
        "_t": "declareticket",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-second-5-name. Ticket ID: N/A."
      }
    ]
    """
