Feature: update alarm on idle rule
  I need to be able to update alarm on idle rule

  @concurrent
  Scenario: given idle rule and no events for alarm should update alarm
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "_id": "test-idlerule-axe-idlerule-1",
      "name": "test-idlerule-axe-idlerule-1-name",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 40,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-axe-idlerule-1"
            }
          }
        ]
      ],
      "comment": "test-idlerule-axe-idlerule-1-comment",
      "operation": {
        "type": "assocticket",
        "parameters": {
          "ticket": "test-idlerule-axe-idlerule-1-ticket",
          "ticket_system_name": "test-idlerule-axe-idlerule-1-system-name",
          "ticket_url": "test-idlerule-axe-idlerule-1-ticket-url",
          "ticket_data": {
            "ticket_param_1": "ticket_value_1"
          }
        }
      }
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-1",
      "connector_name": "test-connector-name-axe-idlerule-1",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-1",
      "resource": "test-resource-axe-idlerule-1",
      "state": 2,
      "output": "test-output-axe-idlerule-1"
    }
    """
    When I wait 3s
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-1",
      "connector_name": "test-connector-name-axe-idlerule-1",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-1",
      "resource": "test-resource-axe-idlerule-1",
      "state": 2,
      "output": "test-output-axe-idlerule-1"
    }
    """
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-1",
            "connector": "test-connector-axe-idlerule-1",
            "connector_name": "test-connector-name-axe-idlerule-1",
            "resource": "test-resource-axe-idlerule-1",
            "state": {
              "val": 2
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
                "val": 2
              },
              {
                "_t": "statusinc",
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
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "assocticket",
      "connector": "test-connector-axe-idlerule-1",
      "connector_name": "test-connector-name-axe-idlerule-1",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-1",
      "resource": "test-resource-axe-idlerule-1"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-1",
            "connector": "test-connector-axe-idlerule-1",
            "connector_name": "test-connector-name-axe-idlerule-1",
            "resource": "test-resource-axe-idlerule-1",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "tickets": [
              {
                "_t": "assocticket",
                "a": "system",
                "initiator": "system",
                "m": "Idle rule: test-idlerule-axe-idlerule-1-name. Ticket ID: test-idlerule-axe-idlerule-1-ticket. Ticket URL: test-idlerule-axe-idlerule-1-ticket-url. Ticket ticket_param_1: ticket_value_1.",
                "ticket": "test-idlerule-axe-idlerule-1-ticket",
                "ticket_rule_id": "test-idlerule-axe-idlerule-1",
                "ticket_rule_name": "Idle rule: test-idlerule-axe-idlerule-1-name",
                "ticket_system_name": "test-idlerule-axe-idlerule-1-system-name",
                "ticket_url": "test-idlerule-axe-idlerule-1-ticket-url",
                "ticket_data": {
                  "ticket_param_1": "ticket_value_1"
                },
                "ticket_comment": "test-idlerule-axe-idlerule-1-comment"
              }
            ],
            "ticket": {
              "_t": "assocticket",
              "a": "system",
              "initiator": "system",
              "m": "Idle rule: test-idlerule-axe-idlerule-1-name. Ticket ID: test-idlerule-axe-idlerule-1-ticket. Ticket URL: test-idlerule-axe-idlerule-1-ticket-url. Ticket ticket_param_1: ticket_value_1.",
              "ticket": "test-idlerule-axe-idlerule-1-ticket",
              "ticket_rule_id": "test-idlerule-axe-idlerule-1",
              "ticket_rule_name": "Idle rule: test-idlerule-axe-idlerule-1-name",
              "ticket_system_name": "test-idlerule-axe-idlerule-1-system-name",
              "ticket_url": "test-idlerule-axe-idlerule-1-ticket-url",
              "ticket_data": {
                "ticket_param_1": "ticket_value_1"
              },
              "ticket_comment": "test-idlerule-axe-idlerule-1-comment"
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
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "assocticket",
                "a": "system",
                "user_id": "",
                "initiator": "system",
                "m": "Idle rule: test-idlerule-axe-idlerule-1-name. Ticket ID: test-idlerule-axe-idlerule-1-ticket. Ticket URL: test-idlerule-axe-idlerule-1-ticket-url. Ticket ticket_param_1: ticket_value_1.",
                "ticket": "test-idlerule-axe-idlerule-1-ticket",
                "ticket_rule_id": "test-idlerule-axe-idlerule-1",
                "ticket_rule_name": "Idle rule: test-idlerule-axe-idlerule-1-name",
                "ticket_system_name": "test-idlerule-axe-idlerule-1-system-name",
                "ticket_url": "test-idlerule-axe-idlerule-1-ticket-url",
                "ticket_data": {
                  "ticket_param_1": "ticket_value_1"
                },
                "ticket_comment": "test-idlerule-axe-idlerule-1-comment"
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
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-1",
      "connector_name": "test-connector-name-axe-idlerule-1",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-1",
      "resource": "test-resource-axe-idlerule-1",
      "state": 2,
      "output": "test-output-axe-idlerule-1"
    }
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "assocticket",
      "connector": "test-connector-axe-idlerule-1",
      "connector_name": "test-connector-name-axe-idlerule-1",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-1",
      "resource": "test-resource-axe-idlerule-1"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-1",
            "connector": "test-connector-axe-idlerule-1",
            "connector_name": "test-connector-name-axe-idlerule-1",
            "resource": "test-resource-axe-idlerule-1",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "tickets": [
              {
                "_t": "assocticket",
                "a": "system",
                "initiator": "system",
                "m": "Idle rule: test-idlerule-axe-idlerule-1-name. Ticket ID: test-idlerule-axe-idlerule-1-ticket. Ticket URL: test-idlerule-axe-idlerule-1-ticket-url. Ticket ticket_param_1: ticket_value_1.",
                "ticket": "test-idlerule-axe-idlerule-1-ticket",
                "ticket_rule_id": "test-idlerule-axe-idlerule-1",
                "ticket_rule_name": "Idle rule: test-idlerule-axe-idlerule-1-name",
                "ticket_system_name": "test-idlerule-axe-idlerule-1-system-name",
                "ticket_url": "test-idlerule-axe-idlerule-1-ticket-url",
                "ticket_data": {
                  "ticket_param_1": "ticket_value_1"
                },
                "ticket_comment": "test-idlerule-axe-idlerule-1-comment"
              },
              {
                "_t": "assocticket",
                "a": "system",
                "initiator": "system",
                "m": "Idle rule: test-idlerule-axe-idlerule-1-name. Ticket ID: test-idlerule-axe-idlerule-1-ticket. Ticket URL: test-idlerule-axe-idlerule-1-ticket-url. Ticket ticket_param_1: ticket_value_1.",
                "ticket": "test-idlerule-axe-idlerule-1-ticket",
                "ticket_rule_id": "test-idlerule-axe-idlerule-1",
                "ticket_rule_name": "Idle rule: test-idlerule-axe-idlerule-1-name",
                "ticket_system_name": "test-idlerule-axe-idlerule-1-system-name",
                "ticket_url": "test-idlerule-axe-idlerule-1-ticket-url",
                "ticket_data": {
                  "ticket_param_1": "ticket_value_1"
                },
                "ticket_comment": "test-idlerule-axe-idlerule-1-comment"
              }
            ],
            "ticket": {
              "_t": "assocticket",
              "a": "system",
              "initiator": "system",
              "m": "Idle rule: test-idlerule-axe-idlerule-1-name. Ticket ID: test-idlerule-axe-idlerule-1-ticket. Ticket URL: test-idlerule-axe-idlerule-1-ticket-url. Ticket ticket_param_1: ticket_value_1.",
              "ticket": "test-idlerule-axe-idlerule-1-ticket",
              "ticket_rule_id": "test-idlerule-axe-idlerule-1",
              "ticket_rule_name": "Idle rule: test-idlerule-axe-idlerule-1-name",
              "ticket_system_name": "test-idlerule-axe-idlerule-1-system-name",
              "ticket_url": "test-idlerule-axe-idlerule-1-ticket-url",
              "ticket_data": {
                "ticket_param_1": "ticket_value_1"
              },
              "ticket_comment": "test-idlerule-axe-idlerule-1-comment"
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
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "assocticket",
                "a": "system",
                "user_id": "",
                "initiator": "system",
                "m": "Idle rule: test-idlerule-axe-idlerule-1-name. Ticket ID: test-idlerule-axe-idlerule-1-ticket. Ticket URL: test-idlerule-axe-idlerule-1-ticket-url. Ticket ticket_param_1: ticket_value_1.",
                "ticket": "test-idlerule-axe-idlerule-1-ticket",
                "ticket_rule_id": "test-idlerule-axe-idlerule-1",
                "ticket_rule_name": "Idle rule: test-idlerule-axe-idlerule-1-name",
                "ticket_system_name": "test-idlerule-axe-idlerule-1-system-name",
                "ticket_url": "test-idlerule-axe-idlerule-1-ticket-url",
                "ticket_data": {
                  "ticket_param_1": "ticket_value_1"
                },
                "ticket_comment": "test-idlerule-axe-idlerule-1-comment"
              },
              {
                "_t": "assocticket",
                "a": "system",
                "user_id": "",
                "initiator": "system",
                "m": "Idle rule: test-idlerule-axe-idlerule-1-name. Ticket ID: test-idlerule-axe-idlerule-1-ticket. Ticket URL: test-idlerule-axe-idlerule-1-ticket-url. Ticket ticket_param_1: ticket_value_1.",
                "ticket": "test-idlerule-axe-idlerule-1-ticket",
                "ticket_rule_id": "test-idlerule-axe-idlerule-1",
                "ticket_rule_name": "Idle rule: test-idlerule-axe-idlerule-1-name",
                "ticket_system_name": "test-idlerule-axe-idlerule-1-system-name",
                "ticket_url": "test-idlerule-axe-idlerule-1-ticket-url",
                "ticket_data": {
                  "ticket_param_1": "ticket_value_1"
                },
                "ticket_comment": "test-idlerule-axe-idlerule-1-comment"
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

  @concurrent
  Scenario: given idle rule and no update for alarm should update alarm
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "_id": "test-idlerule-axe-idlerule-2",
      "name": "test-idlerule-axe-idlerule-2-name",
      "type": "alarm",
      "alarm_condition": "last_update",
      "enabled": true,
      "priority": 41,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-axe-idlerule-2"
            }
          }
        ]
      ],
      "operation": {
        "type": "assocticket",
        "parameters": {
          "ticket": "test-idlerule-axe-idlerule-2-ticket"
        }
      }
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-2",
      "connector_name": "test-connector-name-axe-idlerule-2",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-2",
      "resource": "test-resource-axe-idlerule-2",
      "state": 1,
      "output": "test-output-axe-idlerule-2"
    }
    """
    When I wait 3s
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-2",
      "connector_name": "test-connector-name-axe-idlerule-2",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-2",
      "resource": "test-resource-axe-idlerule-2",
      "state": 2,
      "output": "test-output-axe-idlerule-2"
    }
    """
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-2",
            "connector": "test-connector-axe-idlerule-2",
            "connector_name": "test-connector-name-axe-idlerule-2",
            "resource": "test-resource-axe-idlerule-2",
            "state": {
              "val": 2
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
                "val": 1
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "stateinc",
                "val": 2
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
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "assocticket",
      "connector": "test-connector-axe-idlerule-2",
      "connector_name": "test-connector-name-axe-idlerule-2",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-2",
      "resource": "test-resource-axe-idlerule-2"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-2",
            "connector": "test-connector-axe-idlerule-2",
            "connector_name": "test-connector-name-axe-idlerule-2",
            "resource": "test-resource-axe-idlerule-2",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "tickets": [
              {
                "_t": "assocticket",
                "a": "system",
                "initiator": "system",
                "m": "Idle rule: test-idlerule-axe-idlerule-2-name. Ticket ID: test-idlerule-axe-idlerule-2-ticket.",
                "ticket": "test-idlerule-axe-idlerule-2-ticket",
                "ticket_rule_id": "test-idlerule-axe-idlerule-2",
                "ticket_rule_name": "Idle rule: test-idlerule-axe-idlerule-2-name"
              }
            ],
            "ticket": {
              "_t": "assocticket",
              "a": "system",
              "initiator": "system",
              "m": "Idle rule: test-idlerule-axe-idlerule-2-name. Ticket ID: test-idlerule-axe-idlerule-2-ticket.",
              "ticket": "test-idlerule-axe-idlerule-2-ticket",
              "ticket_rule_id": "test-idlerule-axe-idlerule-2",
              "ticket_rule_name": "Idle rule: test-idlerule-axe-idlerule-2-name"
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
                "val": 1
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "assocticket",
                "a": "system",
                "user_id": "",
                "initiator": "system",
                "m": "Idle rule: test-idlerule-axe-idlerule-2-name. Ticket ID: test-idlerule-axe-idlerule-2-ticket.",
                "ticket": "test-idlerule-axe-idlerule-2-ticket",
                "ticket_rule_id": "test-idlerule-axe-idlerule-2",
                "ticket_rule_name": "Idle rule: test-idlerule-axe-idlerule-2-name"
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
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-2",
      "connector_name": "test-connector-name-axe-idlerule-2",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-2",
      "resource": "test-resource-axe-idlerule-2",
      "state": 2,
      "output": "test-output-axe-idlerule-2"
    }
    """
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-2",
            "connector": "test-connector-axe-idlerule-2",
            "connector_name": "test-connector-name-axe-idlerule-2",
            "resource": "test-resource-axe-idlerule-2",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "tickets": [
              {
                "_t": "assocticket",
                "a": "system",
                "initiator": "system",
                "m": "Idle rule: test-idlerule-axe-idlerule-2-name. Ticket ID: test-idlerule-axe-idlerule-2-ticket.",
                "ticket": "test-idlerule-axe-idlerule-2-ticket",
                "ticket_rule_id": "test-idlerule-axe-idlerule-2",
                "ticket_rule_name": "Idle rule: test-idlerule-axe-idlerule-2-name"
              }
            ],
            "ticket": {
              "_t": "assocticket",
              "a": "system",
              "initiator": "system",
              "m": "Idle rule: test-idlerule-axe-idlerule-2-name. Ticket ID: test-idlerule-axe-idlerule-2-ticket.",
              "ticket": "test-idlerule-axe-idlerule-2-ticket",
              "ticket_rule_id": "test-idlerule-axe-idlerule-2",
              "ticket_rule_name": "Idle rule: test-idlerule-axe-idlerule-2-name"
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
                "val": 1
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "assocticket",
                "a": "system",
                "user_id": "",
                "initiator": "system",
                "m": "Idle rule: test-idlerule-axe-idlerule-2-name. Ticket ID: test-idlerule-axe-idlerule-2-ticket.",
                "ticket": "test-idlerule-axe-idlerule-2-ticket",
                "ticket_rule_id": "test-idlerule-axe-idlerule-2",
                "ticket_rule_name": "Idle rule: test-idlerule-axe-idlerule-2-name"
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
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-2",
      "connector_name": "test-connector-name-axe-idlerule-2",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-2",
      "resource": "test-resource-axe-idlerule-2",
      "state": 3,
      "output": "test-output-axe-idlerule-2"
    }
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "assocticket",
      "connector": "test-connector-axe-idlerule-2",
      "connector_name": "test-connector-name-axe-idlerule-2",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-2",
      "resource": "test-resource-axe-idlerule-2"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-2",
            "connector": "test-connector-axe-idlerule-2",
            "connector_name": "test-connector-name-axe-idlerule-2",
            "resource": "test-resource-axe-idlerule-2",
            "state": {
              "val": 3
            },
            "status": {
              "val": 1
            },
            "tickets": [
              {
                "_t": "assocticket",
                "a": "system",
                "initiator": "system",
                "m": "Idle rule: test-idlerule-axe-idlerule-2-name. Ticket ID: test-idlerule-axe-idlerule-2-ticket.",
                "ticket": "test-idlerule-axe-idlerule-2-ticket",
                "ticket_rule_id": "test-idlerule-axe-idlerule-2",
                "ticket_rule_name": "Idle rule: test-idlerule-axe-idlerule-2-name"
              },
              {
                "_t": "assocticket",
                "a": "system",
                "initiator": "system",
                "m": "Idle rule: test-idlerule-axe-idlerule-2-name. Ticket ID: test-idlerule-axe-idlerule-2-ticket.",
                "ticket": "test-idlerule-axe-idlerule-2-ticket",
                "ticket_rule_id": "test-idlerule-axe-idlerule-2",
                "ticket_rule_name": "Idle rule: test-idlerule-axe-idlerule-2-name"
              }
            ],
            "ticket": {
              "_t": "assocticket",
              "a": "system",
              "initiator": "system",
              "m": "Idle rule: test-idlerule-axe-idlerule-2-name. Ticket ID: test-idlerule-axe-idlerule-2-ticket.",
              "ticket": "test-idlerule-axe-idlerule-2-ticket",
              "ticket_rule_id": "test-idlerule-axe-idlerule-2",
              "ticket_rule_name": "Idle rule: test-idlerule-axe-idlerule-2-name"
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
                "val": 1
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "assocticket",
                "a": "system",
                "user_id": "",
                "initiator": "system",
                "m": "Idle rule: test-idlerule-axe-idlerule-2-name. Ticket ID: test-idlerule-axe-idlerule-2-ticket.",
                "ticket": "test-idlerule-axe-idlerule-2-ticket",
                "ticket_rule_id": "test-idlerule-axe-idlerule-2",
                "ticket_rule_name": "Idle rule: test-idlerule-axe-idlerule-2-name"
              },
              {
                "_t": "stateinc",
                "val": 3
              },
              {
                "_t": "assocticket",
                "a": "system",
                "user_id": "",
                "initiator": "system",
                "m": "Idle rule: test-idlerule-axe-idlerule-2-name. Ticket ID: test-idlerule-axe-idlerule-2-ticket.",
                "ticket": "test-idlerule-axe-idlerule-2-ticket",
                "ticket_rule_id": "test-idlerule-axe-idlerule-2",
                "ticket_rule_name": "Idle rule: test-idlerule-axe-idlerule-2-name"
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
  Scenario: given idle rule and no events for resource should create alarm
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-3-name",
      "type": "entity",
      "enabled": true,
      "priority": 42,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-axe-idlerule-3"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I save response createTimestamp={{ now }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-3",
      "connector_name": "test-connector-name-axe-idlerule-3",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-3",
      "resource": "test-resource-axe-idlerule-3",
      "state": 0,
      "output": "test-output-axe-idlerule-3"
    }
    """
    When I wait 3s
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-3",
      "connector_name": "test-connector-name-axe-idlerule-3",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-3",
      "resource": "test-resource-axe-idlerule-3",
      "state": 0,
      "output": "test-output-axe-idlerule-3"
    }
    """
    When I wait 3s
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "activate",
      "connector": "test-connector-axe-idlerule-3",
      "connector_name": "test-connector-name-axe-idlerule-3",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-3",
      "resource": "test-resource-axe-idlerule-3"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-3",
            "connector": "test-connector-axe-idlerule-3",
            "connector_name": "test-connector-name-axe-idlerule-3",
            "resource": "test-resource-axe-idlerule-3",
            "state": {
              "val": 3,
              "a": "system",
              "user_id": "",
              "initiator": "system",
              "m": "Idle rule test-idlerule-axe-idlerule-3-name"
            },
            "status": {
              "val": 5,
              "a": "system",
              "user_id": "",
              "initiator": "system",
              "m": "Idle rule test-idlerule-axe-idlerule-3-name"
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
    When I save response idleSince={{ (index .lastResponse.data 0).entity.idle_since }}
    Then the difference between idleSince createTimestamp is in range 8,10
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
                "val": 3,
                "a": "system",
                "user_id": "",
                "initiator": "system",
                "m": "Idle rule test-idlerule-axe-idlerule-3-name"
              },
              {
                "_t": "statusinc",
                "val": 5,
                "a": "system",
                "user_id": "",
                "initiator": "system",
                "m": "Idle rule test-idlerule-axe-idlerule-3-name"
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
    When I do GET /api/v4/entities?search=test-resource-axe-idlerule-3&no_events=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-resource-axe-idlerule-3",
          "idle_since": {{ .idleSince }}
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
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-3",
      "connector_name": "test-connector-name-axe-idlerule-3",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-3",
      "resource": "test-resource-axe-idlerule-3",
      "state": 3,
      "output": "test-output-axe-idlerule-3"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-3",
            "connector": "test-connector-axe-idlerule-3",
            "connector_name": "test-connector-name-axe-idlerule-3",
            "resource": "test-resource-axe-idlerule-3",
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
    Then the response key "data.0.entity.idle_since" should not exist
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
                "_t": "statusdec",
                "a": "test-connector-axe-idlerule-3.test-connector-name-axe-idlerule-3",
                "initiator": "external",
                "m": "test-output-axe-idlerule-3",
                "val": 1
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
    When I do GET /api/v4/entities?search=test-resource-axe-idlerule-3&no_events=true
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    When I do GET /api/v4/entities?search=test-resource-axe-idlerule-3
    Then the response code should be 200
    Then the response key "data.0.idle_since" should not exist
    When I do DELETE /api/v4/idle-rules/{{ .ruleID }}
    Then the response code should be 204

  @concurrent
  Scenario: given idle rule and no events for component should create alarm
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-4-name",
      "type": "entity",
      "enabled": true,
      "priority": 43,
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
              "value": "test-component-axe-idlerule-4"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-4",
      "connector_name": "test-connector-name-axe-idlerule-4",
      "source_type": "component",
      "component":  "test-component-axe-idlerule-4",
      "state": 0,
      "output": "test-output-axe-idlerule-4"
    }
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "activate",
      "connector": "test-connector-axe-idlerule-4",
      "connector_name": "test-connector-name-axe-idlerule-4",
      "source_type": "component",
      "component":  "test-component-axe-idlerule-4"
    }
    """
    When I do GET /api/v4/alarms?search=test-component-axe-idlerule-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-4",
            "connector": "test-connector-axe-idlerule-4",
            "connector_name": "test-connector-name-axe-idlerule-4",
            "state": {
              "val": 3,
              "a": "system",
              "initiator": "system",
              "m": "Idle rule test-idlerule-axe-idlerule-4-name"
            },
            "status": {
              "val": 5,
              "a": "system",
              "initiator": "system",
              "m": "Idle rule test-idlerule-axe-idlerule-4-name"
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
                "val": 3,
                "a": "system",
                "initiator": "system",
                "m": "Idle rule test-idlerule-axe-idlerule-4-name"
              },
              {
                "_t": "statusinc",
                "val": 5,
                "a": "system",
                "initiator": "system",
                "m": "Idle rule test-idlerule-axe-idlerule-4-name"
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
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-4",
      "connector_name": "test-connector-name-axe-idlerule-4",
      "source_type": "component",
      "component":  "test-component-axe-idlerule-4",
      "state": 0,
      "output": "test-output-axe-idlerule-4"
    }
    """
    When I do GET /api/v4/alarms?search=test-component-axe-idlerule-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-4",
            "connector": "test-connector-axe-idlerule-4",
            "connector_name": "test-connector-name-axe-idlerule-4",
            "state": {
              "val": 0
            },
            "status": {
              "val": 0
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
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 5
              },
              {
                "_t": "statedec",
                "val": 0
              },
              {
                "_t": "statusdec",
                "val": 0
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
    When I do DELETE /api/v4/idle-rules/{{ .ruleID }}
    Then the response code should be 204

  @concurrent
  Scenario: given idle rule and no events for connector should create alarm
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-5-name",
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
              "value": "test-connector-name-axe-idlerule-5"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-5",
      "connector_name": "test-connector-name-axe-idlerule-5",
      "source_type": "component",
      "component":  "test-component-axe-idlerule-5",
      "state": 0,
      "output": "test-output-axe-idlerule-5"
    }
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "activate",
      "connector": "test-connector-axe-idlerule-5",
      "connector_name": "test-connector-name-axe-idlerule-5",
      "source_type": "connector"
    }
    """
    When I do GET /api/v4/alarms?search=test-connector-axe-idlerule-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-axe-idlerule-5",
            "connector_name": "test-connector-name-axe-idlerule-5",
            "state": {
              "val": 3,
              "a": "system",
              "initiator": "system",
              "m": "Idle rule test-idlerule-axe-idlerule-5-name"
            },
            "status": {
              "val": 5,
              "a": "system",
              "initiator": "system",
              "m": "Idle rule test-idlerule-axe-idlerule-5-name"
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
                "val": 3,
                "a": "system",
                "initiator": "system",
                "m": "Idle rule test-idlerule-axe-idlerule-5-name"
              },
              {
                "_t": "statusinc",
                "val": 5,
                "a": "system",
                "initiator": "system",
                "m": "Idle rule test-idlerule-axe-idlerule-5-name"
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
    When I wait 1s
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-5",
      "connector_name": "test-connector-name-axe-idlerule-5",
      "source_type": "component",
      "component":  "test-component-axe-idlerule-5",
      "state": 0,
      "output": "test-output-axe-idlerule-5"
    }
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "noevents",
      "connector": "test-connector-axe-idlerule-5",
      "connector_name": "test-connector-name-axe-idlerule-5",
      "source_type": "connector"
    }
    """
    When I do GET /api/v4/alarms?search=test-connector-axe-idlerule-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-axe-idlerule-5",
            "connector_name": "test-connector-name-axe-idlerule-5",
            "state": {
              "val": 0
            },
            "status": {
              "val": 0
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
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 5
              },
              {
                "_t": "statedec",
                "val": 0
              },
              {
                "_t": "statusdec",
                "val": 0
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
    When I do DELETE /api/v4/idle-rules/{{ .ruleID }}
    Then the response code should be 204

  @concurrent
  Scenario: given idle rule and no events for alarm should apply most priority rule only once
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-6-1-name",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 46,
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
              "value": "test-resource-axe-idlerule-6"
            }
          }
        ]
      ],
      "operation": {
        "type": "ack",
        "parameters": {
          "output": "test-idlerule-axe-idlerule-6-1-output"
        }
      }
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/idle-rules:
    """json
    {
      "_id": "test-idlerule-axe-idlerule-6-2",
      "name": "test-idlerule-axe-idlerule-6-2-name",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 45,
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
              "value": "test-resource-axe-idlerule-6"
            }
          }
        ]
      ],
      "operation": {
        "type": "assocticket",
        "parameters": {
          "ticket": "test-idlerule-axe-idlerule-6-2-ticket"
        }
      }
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-6",
      "connector_name": "test-connector-name-axe-idlerule-6",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-6",
      "resource": "test-resource-axe-idlerule-6",
      "state": 2,
      "output": "test-output-axe-idlerule-6"
    }
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "assocticket",
      "connector": "test-connector-axe-idlerule-6",
      "connector_name": "test-connector-name-axe-idlerule-6",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-6",
      "resource": "test-resource-axe-idlerule-6"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-6
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-6",
            "connector": "test-connector-axe-idlerule-6",
            "connector_name": "test-connector-name-axe-idlerule-6",
            "resource": "test-resource-axe-idlerule-6",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "tickets": [
              {
                "_t": "assocticket",
                "a": "system",
                "initiator": "system",
                "ticket": "test-idlerule-axe-idlerule-6-2-ticket"
              }
            ],
            "ticket": {
              "_t": "assocticket",
              "a": "system",
              "initiator": "system",
              "ticket": "test-idlerule-axe-idlerule-6-2-ticket"
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
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "assocticket",
                "a": "system",
                "user_id": "",
                "initiator": "system",
                "ticket": "test-idlerule-axe-idlerule-6-2-ticket"
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
    When I wait 5s
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-6
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
                "ticket": "test-idlerule-axe-idlerule-6-2-ticket"
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
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-6",
      "connector_name": "test-connector-name-axe-idlerule-6",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-6",
      "resource": "test-resource-axe-idlerule-6",
      "state": 2,
      "output": "test-output-axe-idlerule-6"
    }
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "assocticket",
      "connector": "test-connector-axe-idlerule-6",
      "connector_name": "test-connector-name-axe-idlerule-6",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-6",
      "resource": "test-resource-axe-idlerule-6"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-6
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
                "ticket": "test-idlerule-axe-idlerule-6-2-ticket"
              },
              {
                "_t": "assocticket",
                "ticket": "test-idlerule-axe-idlerule-6-2-ticket"
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

  @concurrent
  Scenario: given idle rule and no events for resource should apply most priority rule only once
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-7-1-name",
      "type": "entity",
      "enabled": true,
      "priority": 48,
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
              "value": "test-resource-axe-idlerule-7"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response rule1ID={{ .lastResponse._id }}
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-7-2-name",
      "type": "entity",
      "enabled": true,
      "priority": 47,
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
              "value": "test-resource-axe-idlerule-7"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response rule2ID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-7",
      "connector_name": "test-connector-name-axe-idlerule-7",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-7",
      "resource": "test-resource-axe-idlerule-7",
      "state": 0,
      "output": "test-output-axe-idlerule-7"
    }
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "activate",
      "connector": "test-connector-axe-idlerule-7",
      "connector_name": "test-connector-name-axe-idlerule-7",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-7",
      "resource":  "test-resource-axe-idlerule-7"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-7",
            "connector": "test-connector-axe-idlerule-7",
            "connector_name": "test-connector-name-axe-idlerule-7",
            "resource": "test-resource-axe-idlerule-7",
            "state": {
              "val": 3
            },
            "status": {
              "val": 5
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
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 5,
                "a": "system",
                "initiator": "system"
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
    When I wait 5s
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-7
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
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-7",
      "connector_name": "test-connector-name-axe-idlerule-7",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-7",
      "resource": "test-resource-axe-idlerule-7",
      "state": 0,
      "output": "test-output-axe-idlerule-7"
    }
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "noevents",
      "connector": "test-connector-axe-idlerule-7",
      "connector_name": "test-connector-name-axe-idlerule-7",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-7",
      "resource": "test-resource-axe-idlerule-7"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-7
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
                "val": 5,
                "a": "system",
                "initiator": "system"
              },
              {
                "_t": "statedec",
                "val": 0
              },
              {
                "_t": "statusdec",
                "val": 0
              },
              {
                "_t": "stateinc",
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 5,
                "a": "system",
                "initiator": "system"
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
    When I do DELETE /api/v4/idle-rules/{{ .rule1ID }}
    Then the response code should be 204
    When I do DELETE /api/v4/idle-rules/{{ .rule2ID }}
    Then the response code should be 204

  @concurrent
  Scenario: given idle rule and no events for resource should update existed alarm
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-8-name",
      "type": "entity",
      "enabled": true,
      "priority": 49,
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
              "value": "test-resource-axe-idlerule-8"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I save response createTimestamp={{ now }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-8",
      "connector_name": "test-connector-name-axe-idlerule-8",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-8",
      "resource": "test-resource-axe-idlerule-8",
      "state": 1,
      "output": "test-output-axe-idlerule-8"
    }
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "noevents",
      "connector": "test-connector-axe-idlerule-8",
      "connector_name": "test-connector-name-axe-idlerule-8",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-8",
      "resource": "test-resource-axe-idlerule-8"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-8",
            "connector": "test-connector-axe-idlerule-8",
            "connector_name": "test-connector-name-axe-idlerule-8",
            "resource": "test-resource-axe-idlerule-8",
            "state": {
              "val": 3
            },
            "status": {
              "val": 5
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
    When I save response idleSince={{ (index .lastResponse.data 0).entity.idle_since }}
    Then the difference between idleSince createTimestamp is in range 3,5
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
                "val": 1,
                "a": "test-connector-axe-idlerule-8.test-connector-name-axe-idlerule-8",
                "initiator": "external"
              },
              {
                "_t": "statusinc",
                "val": 1,
                "a": "test-connector-axe-idlerule-8.test-connector-name-axe-idlerule-8",
                "initiator": "external"
              },
              {
                "_t": "stateinc",
                "val": 3,
                "a": "system",
                "initiator": "system"
              },
              {
                "_t": "statusinc",
                "val": 5,
                "a": "system",
                "initiator": "system"
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
    When I do GET /api/v4/entities?search=test-resource-axe-idlerule-8&no_events=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-resource-axe-idlerule-8",
          "idle_since": {{ .idleSince }}
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
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-8",
      "connector_name": "test-connector-name-axe-idlerule-8",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-8",
      "resource": "test-resource-axe-idlerule-8",
      "state": 1,
      "output": "test-output-axe-idlerule-8"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-8",
            "connector": "test-connector-axe-idlerule-8",
            "connector_name": "test-connector-name-axe-idlerule-8",
            "resource": "test-resource-axe-idlerule-8",
            "state": {
              "val": 1
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
    Then the response key "data.0.entity.idle_since" should not exist
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
                "val": 1
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "stateinc",
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 5
              },
              {
                "_t": "statedec",
                "val": 1
              },
              {
                "_t": "statusdec",
                "val": 1
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
    When I do GET /api/v4/entities?search=test-resource-axe-idlerule-8&no_events=true
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    When I do GET /api/v4/entities?search=test-resource-axe-idlerule-8
    Then the response code should be 200
    Then the response key "data.0.idle_since" should not exist
    When I do DELETE /api/v4/idle-rules/{{ .ruleID }}
    Then the response code should be 204

  @concurrent
  Scenario: given idle rule and no events for alarm and entity should apply most priority alarm rule
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "_id": "test-idlerule-axe-idlerule-9-1",
      "name": "test-idlerule-axe-idlerule-9-1-name",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 50,
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
              "value": "test-resource-axe-idlerule-9"
            }
          }
        ]
      ],
      "operation": {
        "type": "assocticket",
        "parameters": {
          "ticket": "test-idlerule-axe-idlerule-9-1-ticket"
          }
        }
      }
    }
    """
    Then the response code should be 201
    Then I save response rule1ID={{ .lastResponse._id }}
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-9-2-name",
      "type": "entity",
      "enabled": true,
      "priority": 51,
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
              "value": "test-resource-axe-idlerule-9"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response rule2ID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-9",
      "connector_name": "test-connector-name-axe-idlerule-9",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-9",
      "resource": "test-resource-axe-idlerule-9",
      "state": 2,
      "output": "test-output-axe-idlerule-9"
    }
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "assocticket",
      "connector": "test-connector-axe-idlerule-9",
      "connector_name": "test-connector-name-axe-idlerule-9",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-9",
      "resource": "test-resource-axe-idlerule-9"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-9
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-9",
            "connector": "test-connector-axe-idlerule-9",
            "connector_name": "test-connector-name-axe-idlerule-9",
            "resource": "test-resource-axe-idlerule-9",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "tickets": [
              {
                "_t": "assocticket",
                "a": "system",
                "initiator": "system",
                "ticket": "test-idlerule-axe-idlerule-9-1-ticket"
              }
            ],
            "ticket": {
              "_t": "assocticket",
              "a": "system",
              "initiator": "system",
              "ticket": "test-idlerule-axe-idlerule-9-1-ticket"
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
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "assocticket",
                "a": "system",
                "user_id": "",
                "initiator": "system",
                "ticket": "test-idlerule-axe-idlerule-9-1-ticket"
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
    When I wait 5s
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-9
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
                "ticket": "test-idlerule-axe-idlerule-9-1-ticket"
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
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-9",
      "connector_name": "test-connector-name-axe-idlerule-9",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-9",
      "resource": "test-resource-axe-idlerule-9",
      "state": 2,
      "output": "test-output-axe-idlerule-9"
    }
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "assocticket",
      "connector": "test-connector-axe-idlerule-9",
      "connector_name": "test-connector-name-axe-idlerule-9",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-9",
      "resource": "test-resource-axe-idlerule-9"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-9
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
                "ticket": "test-idlerule-axe-idlerule-9-1-ticket"
              },
              {
                "_t": "assocticket",
                "ticket": "test-idlerule-axe-idlerule-9-1-ticket"
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
    When I do DELETE /api/v4/idle-rules/{{ .rule1ID }}
    Then the response code should be 204
    When I do DELETE /api/v4/idle-rules/{{ .rule2ID }}
    Then the response code should be 204