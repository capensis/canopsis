Feature: execute action on trigger
  I need to be able to trigger action on event

  Scenario: given scenario and check event should update alarm
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-1-name",
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
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .apiURL }}/api/v4/scenarios",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {"Content-Type": "application/json"},
              "payload": "{\"name\":\"test-scenario-action-webhook-1 {{ `{{ .Entity.ID }}` }}\",\"enabled\":true,\"triggers\":[\"create\"],\"actions\":[{\"entity_pattern\":[[{\"field\":\"name\",\"cond\":{\"type\": \"eq\", \"value\": \"test-scenario-action-webhook-1\"}}]],\"type\":\"ack\",\"drop_scenario_if_not_matched\":false,\"emit_trigger\":false}]}"
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
    When I wait the end of event processing
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
                "ticket_data": {
                  "scenario_name": "test-scenario-action-webhook-1 test-resource-action-webhook-1/test-component-action-webhook-1"
                }
              }
            ],
            "ticket": {
              "_t": "declareticket",
              "a": "system",
              "ticket_data": {
                "scenario_name": "test-scenario-action-webhook-1 test-resource-action-webhook-1/test-component-action-webhook-1"
              }
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
    When I save response ticketID={{ (index (index .lastResponse.data 0).v.tickets 0).ticket }}
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
        "m": "Scenario test-scenario-action-webhook-1-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-1-name. Ticket ID: {{ .ticketID }}. Ticket scenario_name: test-scenario-action-webhook-1 test-resource-action-webhook-1/test-component-action-webhook-1."
      },
      {
        "_t": "declareticket",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-1-name. Ticket ID: {{ .ticketID }}. Ticket scenario_name: test-scenario-action-webhook-1 test-resource-action-webhook-1/test-component-action-webhook-1."
      }
    ]
    """

  Scenario: given scenario and check event should emit declare ticket trigger
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-2-1-name",
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
            "forward_author": true,
            "author": "test-scenario-action-webhook-2-action-1-author {{ `{{ .Alarm.Value.Resource }}` }}",
            "request": {
              "method": "POST",
              "url": "{{ .apiURL }}/api/v4/scenarios",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {"Content-Type": "application/json"},
              "payload": "{\"name\":\"test-scenario-action-webhook-2 {{ `{{ .Alarm.Value.Resource }}` }}\",\"enabled\":true,\"triggers\":[\"create\"],\"actions\":[{\"entity_pattern\":[[{\"field\":\"name\",\"cond\":{\"type\": \"eq\", \"value\": \"test-scenario-action-webhook-2\"}}]],\"type\":\"ack\",\"drop_scenario_if_not_matched\":false,\"emit_trigger\":false}]}"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "_id",
              "scenario_name": "name"
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
    When I send an event:
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
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-action-webhook-2
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
                "a": "root",
                "ticket_data": {
                  "scenario_name": "test-scenario-action-webhook-2 test-resource-action-webhook-2"
                }
              }
            ],
            "ticket": {
              "_t": "declareticket",
              "a": "root",
              "ticket_data": {
                "scenario_name": "test-scenario-action-webhook-2 test-resource-action-webhook-2"
              }
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
    When I save response ticketID={{ (index (index .lastResponse.data 0).v.tickets 0).ticket }}
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
        "m": "Scenario test-scenario-action-webhook-2-1-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Scenario: test-scenario-action-webhook-2-1-name. Ticket ID: {{ .ticketID }}. Ticket scenario_name: test-scenario-action-webhook-2 test-resource-action-webhook-2."
      },
      {
        "_t": "declareticket",
        "a": "root",
        "user_id": "root",
        "m": "Scenario: test-scenario-action-webhook-2-1-name. Ticket ID: {{ .ticketID }}. Ticket scenario_name: test-scenario-action-webhook-2 test-resource-action-webhook-2."
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

  Scenario: given scenario and ack resources event should update resource alarms
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-3-name",
      "enabled": true,
      "triggers": ["ack"],
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "eq",
                  "value": "test-component-action-webhook-3"
                }
              }
            ]
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
              "payload": "{\"name\":\"test-scenario-action-webhook-3\",\"enabled\":true,\"triggers\":[\"create\"],\"actions\":[{\"entity_pattern\":[[{\"field\":\"name\",\"cond\":{\"type\": \"eq\", \"value\": \"test-scenario-action-webhook-3\"}}]],\"type\":\"ack\",\"drop_scenario_if_not_matched\":false,\"emit_trigger\":false}]}"
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
      "connector" : "test-connector-action-webhook-3",
      "connector_name" : "test-connector-name-action-webhook-3",
      "source_type" : "component",
      "event_type" : "check",
      "component" :  "test-component-action-webhook-3",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I send an event:
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
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-action-webhook-3",
      "connector_name" : "test-connector-name-action-webhook-3",
      "source_type" : "component",
      "event_type" : "ack",
      "component" :  "test-component-action-webhook-3",
      "ack_resources": true,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-component-action-webhook-3&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "tickets": [
              {
                "_t": "declareticket"
              }
            ],
            "ticket": {
              "_t": "declareticket"
            },
            "ack": {
              "_t": "ack"
            },
            "connector": "test-connector-action-webhook-3",
            "connector_name": "test-connector-name-action-webhook-3",
            "component": "test-component-action-webhook-3"
          }
        },
        {
          "v": {
            "tickets": [
              {
                "_t": "declareticket"
              }
            ],
            "ticket": {
              "_t": "declareticket"
            },
            "ack": {
              "_t": "ack"
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
        "total_count": 2
      }
    }
    """
    When I save response ticketID={{ (index (index .lastResponse.data 0).v.tickets 0).ticket }}
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
        "_t": "ack",
        "a": "root",
        "user_id": "root"
      },
      {
        "_t": "webhookstart",
        "a": "system",
        "user_id": "",
        "m": "Scenario test-scenario-action-webhook-3-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-3-name. Ticket ID: {{ .ticketID }}. Ticket scenario_name: test-scenario-action-webhook-3."
      },
      {
        "_t": "declareticket",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-3-name. Ticket ID: {{ .ticketID }}. Ticket scenario_name: test-scenario-action-webhook-3."
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
        "_t": "ack",
        "a": "root",
        "user_id": "root"
      },
      {
        "_t": "declareticket",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-3-name. Ticket ID: {{ .ticketID }}. Ticket scenario_name: test-scenario-action-webhook-3."
      }
    ]
    """

  Scenario: given webhook scenario to test response and header templates
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-4-name",
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
              "method": "GET",
              "url": "{{ .apiURL }}/api/v4/pbehavior-types/test-default-active-type",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {"Content-Type": "application/json"}
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
              "url": "{{ .apiURL }}/api/v4/scenarios",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {"Content-Type": "{{ `{{index .Header \"Content-Type\"}}` }}"},
              "payload": "{\"name\":\"test-scenario-action-webhook-4-webhook {{ `{{index .Response \"icon_name\"}}` }}\",\"enabled\":true,\"triggers\":[\"create\"],\"actions\":[{\"alarm_pattern\":[[{\"field\":\"v.resource\",\"cond\":{\"type\": \"eq\", \"value\": \"{{ `{{index .Response \"_id\"}}` }}\"}}]],\"type\":\"ack\",\"drop_scenario_if_not_matched\":false,\"emit_trigger\":false}]}"
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
    When I wait the end of event processing
    When I do GET /api/v4/scenarios?search=test-scenario-action-webhook-4-webhook
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-scenario-action-webhook-4-webhook test-active-icon",
          "actions": [
            {
              "alarm_pattern": [
                [
                  {
                    "field": "v.resource",
                    "cond": {
                      "type": "eq",
                      "value": "test-default-active-type"
                    }
                  }
                ]
              ]
            }
          ]
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

  Scenario: given scenarios with 2 actions and webhook should be able to use additional data in template
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-6-1-name",
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
              "url": "{{ .apiURL }}/api/v4/scenarios",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {
                "Content-Type": "application/json"
              },
              "payload": "{\"name\": \"{{ `test-scenario-action-webhook-6-1-action-1 [{{ .AdditionalData.AlarmChangeType }}] [{{ .AdditionalData.Author }}] [{{ .AdditionalData.Initiator }}] [{{ .AdditionalData.User }}]` }}\", \"enabled\":true,\"triggers\":[\"create\"],\"actions\":[{\"alarm_pattern\":[[{\"field\":\"v.resource\",\"cond\":{\"type\": \"eq\", \"value\": \"test-scenario-action-webhook-6-alarm\"}}]],\"type\":\"ack\",\"drop_scenario_if_not_matched\":false,\"emit_trigger\":false}]}"
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
              "url": "{{ .apiURL }}/api/v4/scenarios",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {
                "Content-Type": "application/json"
              },
              "payload": "{\"name\": \"{{ `test-scenario-action-webhook-6-1-action-2 [{{ .AdditionalData.AlarmChangeType }}] [{{ .AdditionalData.Author }}] [{{ .AdditionalData.Initiator }}] [{{ .AdditionalData.User }}]` }}\", \"enabled\":true,\"triggers\":[\"create\"],\"actions\":[{\"alarm_pattern\":[[{\"field\":\"v.resource\",\"cond\":{\"type\": \"eq\", \"value\": \"test-scenario-action-webhook-6-alarm\"}}]],\"type\":\"ack\",\"drop_scenario_if_not_matched\":false,\"emit_trigger\":false}]}"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "_id",
              "scenario_name_2": "name"
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
              "url": "{{ .apiURL }}/api/v4/scenarios",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {
                "Content-Type": "application/json"
              },
              "payload": "{\"name\": \"{{ `test-scenario-action-webhook-6-2-action-1 [{{ .AdditionalData.AlarmChangeType }}] [{{ .AdditionalData.Author }}] [{{ .AdditionalData.Initiator }}] [{{ .AdditionalData.User }}]` }}\", \"enabled\":true,\"triggers\":[\"create\"],\"actions\":[{\"alarm_pattern\":[[{\"field\":\"v.resource\",\"cond\":{\"type\": \"eq\", \"value\": \"test-scenario-action-webhook-6-alarm\"}}]],\"type\":\"ack\",\"drop_scenario_if_not_matched\":false,\"emit_trigger\":false}]}"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "_id",
              "scenario_name_3": "name"
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
              "url": "{{ .apiURL }}/api/v4/scenarios",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {
                "Content-Type": "application/json"
              },
              "payload": "{\"name\": \"{{ `{{ $testVar := .AdditionalData.Output }}test-scenario-action-webhook-6-3-action-1 [{{$testVar}}]` }}\", \"enabled\":true,\"triggers\":[\"create\"],\"actions\":[{\"alarm_pattern\":[[{\"field\":\"v.resource\",\"cond\":{\"type\": \"eq\", \"value\": \"test-scenario-action-webhook-6-alarm\"}}]],\"type\":\"ack\",\"drop_scenario_if_not_matched\":false,\"emit_trigger\":false}]}"
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
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-action-webhook-6
    Then the response code should be 200
    Then the response array key "data.0.v.tickets" should contain:
    """json
    [
      {
        "_t": "declareticket",
        "ticket_data": {
          "scenario_name": "test-scenario-action-webhook-6-1-action-1 [create] [root] [user] [root]"
        }
      },
      {
        "_t": "declareticket",
        "ticket_data": {
          "scenario_name_2": "test-scenario-action-webhook-6-1-action-2 [declareticketwebhook] [system] [user] []"
        }
      },
      {
        "_t": "declareticket",
        "ticket_data": {
          "scenario_name_3": "test-scenario-action-webhook-6-2-action-1 [declareticketwebhook] [test-scenario-action-webhook-6-2-action-1-author] [user] []"
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
        "m": "Scenario test-scenario-action-webhook-6-1-name"
      },
      {
        "_t": "webhookstart",
        "a": "test-scenario-action-webhook-6-3-action-1-author",
        "user_id": "",
        "m": "Scenario test-scenario-action-webhook-6-3-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "root",
        "user_id": "root"
      },
      {
        "_t": "declareticket",
        "a": "root",
        "user_id": "root"
      },
      {
        "_t": "webhookcomplete",
        "a": "test-scenario-action-webhook-6-3-action-1-author",
        "user_id": ""
      },
      {
        "_t": "webhookstart",
        "a": "system",
        "user_id": "",
        "m": "Scenario test-scenario-action-webhook-6-1-name"
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
        "_t": "webhookstart",
        "a": "test-scenario-action-webhook-6-2-action-1-author",
        "user_id": "",
        "m": "Scenario test-scenario-action-webhook-6-2-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "system",
        "user_id": ""
      },
      {
        "_t": "declareticket",
        "a": "test-scenario-action-webhook-6-2-action-1-author",
        "user_id": ""
      }
    ]
    """
    Then the response array key "0.data.steps.data" should contain in order:
    """json
    [
      {
        "_t": "declareticket",
        "a": "root",
        "user_id": "root"
      },
      {
        "_t": "declareticket",
        "a": "system",
        "user_id": ""
      },
      {
        "_t": "declareticket",
        "a": "test-scenario-action-webhook-6-2-action-1-author",
        "user_id": ""
      }
    ]
    """
    When I do GET /api/v4/scenarios?search=test-scenario-action-webhook-6-1-action-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-scenario-action-webhook-6-1-action-1 [create] [root] [user] [root]",
          "enabled": true,
          "triggers": [
            "create"
          ],
          "actions": [
            {
              "alarm_pattern": [
                [
                  {
                    "field": "v.resource",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-action-webhook-6-alarm"
                    }
                  }
                ]
              ]
            }
          ]
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 1
      }
    }
    """
    When I do GET /api/v4/scenarios?search=test-scenario-action-webhook-6-1-action-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-scenario-action-webhook-6-1-action-2 [declareticketwebhook] [system] [user] []",
          "enabled": true,
          "triggers": [
            "create"
          ],
          "actions": [
            {
              "alarm_pattern": [
                [
                  {
                    "field": "v.resource",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-action-webhook-6-alarm"
                    }
                  }
                ]
              ]
            }
          ]
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 1
      }
    }
    """
    When I do GET /api/v4/scenarios?search=test-scenario-action-webhook-6-2-action-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-scenario-action-webhook-6-2-action-1 [declareticketwebhook] [test-scenario-action-webhook-6-2-action-1-author] [user] []",
          "enabled": true,
          "triggers": [
            "create"
          ],
          "actions": [
            {
              "alarm_pattern": [
                [
                  {
                    "field": "v.resource",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-action-webhook-6-alarm"
                    }
                  }
                ]
              ]
            }
          ]
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 1
      }
    }
    """
    When I do GET /api/v4/scenarios?search=test-scenario-action-webhook-6-3-action-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-scenario-action-webhook-6-3-action-1 [noveo alarm]",
          "enabled": true,
          "triggers": [
            "create"
          ],
          "actions": [
            {
              "alarm_pattern": [
                [
                  {
                    "field": "v.resource",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-action-webhook-6-alarm"
                    }
                  }
                ]
              ]
            }
          ]
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 1
      }
    }
    """

  Scenario: given webhook scenario to test multiple response templates
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-7-name",
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
              "method": "GET",
              "url": "{{ .apiURL }}/api/v4/pbehavior-types/test-default-active-type",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {"Content-Type": "application/json"}
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
              "method": "GET",
              "url": "{{ .apiURL }}/api/v4/pbehavior-types/test-default-maintenance-type",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {"Content-Type": "application/json"}
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
              "url": "{{ .apiURL }}/api/v4/scenarios",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {"Content-Type": "{{ `{{index .Header \"Content-Type\"}}` }}"},
              "payload": "{\"name\":\"test-scenario-action-webhook-7-webhook {{ `{{index .ResponseMap \"0._id\"}}` }}\",\"enabled\":true,\"triggers\":[\"create\"],\"actions\":[{\"entity_pattern\":[[{\"field\":\"_id\",\"cond\":{\"type\":\"eq\",\"value\":\"{{ `{{index .ResponseMap \"1._id\"}}` }}\"}}]],\"type\":\"ack\",\"drop_scenario_if_not_matched\":false,\"emit_trigger\":false}]}"
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
    When I wait the end of event processing
    When I do GET /api/v4/scenarios?search=test-scenario-action-webhook-7-webhook
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-scenario-action-webhook-7-webhook test-default-active-type",
          "actions": [
            {
              "entity_pattern": [
                [
                  {
                    "field": "_id",
                    "cond": {
                      "type": "eq",
                      "value": "test-default-maintenance-type"
                    }
                  }
                ]
              ]
            }
          ]
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

  Scenario: given webhook scenario to test document with a document with arrays in response
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-8-name",
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
              "method": "GET",
              "url": "http://localhost:3000/webhook/document_with_array",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {"Content-Type": "application/json"}
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
              "url": "{{ .apiURL }}/api/v4/scenarios",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {"Content-Type": "application/json"},
              "payload": "{\"name\":\"test-scenario-action-webhook-8-webhook {{ `{{index .ResponseMap \"0.array.0.elem1\"}}` }}\",\"enabled\":true,\"triggers\":[\"create\"],\"actions\":[{\"entity_pattern\":[[{\"field\":\"_id\",\"cond\":{\"type\":\"eq\",\"value\":\"{{ `{{index .ResponseMap \"0.array.1.elem2\"}}` }}\"}}]],\"type\":\"ack\",\"drop_scenario_if_not_matched\":false,\"emit_trigger\":false}]}"
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
    When I wait the end of event processing
    When I do GET /api/v4/scenarios?search=test-scenario-action-webhook-8-webhook
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-scenario-action-webhook-8-webhook test1",
          "actions": [
            {
              "entity_pattern": [
                [
                  {
                    "field": "_id",
                    "cond": {
                      "type": "eq",
                      "value": "test4"
                    }
                  }
                ]
              ]
            }
          ]
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
              }
            ],
            "ticket": {
              "_t": "declareticket",
              "a": "system",
              "ticket": "test3",
              "ticket_data": {
                "test_val": "test1"
              }
            },
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
        "m": "Scenario test-scenario-action-webhook-8-name"
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
        "_t": "webhookstart",
        "a": "system",
        "user_id": "",
        "m": "Scenario test-scenario-action-webhook-8-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "system",
        "user_id": ""
      }
    ]
    """

  Scenario: given webhook scenario where the webhook response is an array
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-action-webhook-9-name",
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
              "method": "GET",
              "url": "http://localhost:3000/webhook/response_is_array",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {"Content-Type": "application/json"}
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
              "url": "{{ .apiURL }}/api/v4/scenarios",
              "auth": {
                "username": "root",
                "password": "test"
              },
              "headers": {"Content-Type": "application/json"},
              "payload": "{\"name\":\"test-scenario-action-webhook-9-webhook {{ `{{index .ResponseMap \"0.0.elem1\"}}` }}\",\"enabled\":true,\"triggers\":[\"create\"],\"actions\":[{\"entity_pattern\":[[{\"field\":\"_id\",\"cond\":{\"type\":\"eq\",\"value\":\"{{ `{{index .ResponseMap \"0.1.elem2\"}}` }}\"}}]],\"type\":\"ack\",\"drop_scenario_if_not_matched\":false,\"emit_trigger\":false}]}"
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
    When I wait the end of event processing
    When I do GET /api/v4/scenarios?search=test-scenario-action-webhook-9-webhook
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-scenario-action-webhook-9-webhook test1",
          "actions": [
            {
              "entity_pattern": [
                [
                  {
                    "field": "_id",
                    "cond": {
                      "type": "eq",
                      "value": "test4"
                    }
                  }
                ]
              ]
            }
          ]
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
              }
            ],
            "ticket": {
              "_t": "declareticket",
              "a": "system",
              "ticket": "test3",
              "ticket_data": {
                "test_val": "test1"
              }
            },
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
        "m": "Scenario test-scenario-action-webhook-9-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-9-name. Ticket ID: test3. Ticket test_val: test1."
      },
      {
        "_t": "declareticket",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-9-name. Ticket ID: test3. Ticket test_val: test1."
      },
      {
        "_t": "webhookstart",
        "a": "system",
        "user_id": "",
        "m": "Scenario test-scenario-action-webhook-9-name"
      },
      {
        "_t": "webhookcomplete",
        "a": "system",
        "user_id": "",
        "m": "Scenario: test-scenario-action-webhook-9-name"
      }
    ]
    """
