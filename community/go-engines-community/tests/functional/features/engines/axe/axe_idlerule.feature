Feature: update alarm on idle rule
  I need to be able to update alarm on idle rule

  Scenario: given idle rule and no events for alarm should update alarm
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-1-name",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 40,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "entity_patterns": [
        {
          "_id": "test-resource-axe-idlerule-1/test-component-axe-idlerule-1"
        }
      ],
      "operation": {
        "type": "assocticket",
        "parameters": {
          "ticket": "test-idlerule-axe-idlerule-1-ticket",
          "output": "test-idlerule-axe-idlerule-1-output"
        }
      }
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
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
    When I wait the end of event processing
    When I wait 3s
    When I send an event:
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
    When I wait the end of event processing
    When I wait 3s
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-axe-idlerule-1/test-component-axe-idlerule-1"}]}&with_steps=true
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
            "steps": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
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
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-axe-idlerule-1/test-component-axe-idlerule-1"}]}&with_steps=true
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
            "ticket": {
              "_t": "assocticket",
              "a": "system",
              "m": "test-idlerule-axe-idlerule-1-ticket"
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
                "_t": "assocticket",
                "a": "system",
                "user_id": "",
                "m": "test-idlerule-axe-idlerule-1-ticket"
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
    When I send an event:
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
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-axe-idlerule-1/test-component-axe-idlerule-1"}]}&with_steps=true
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
            "ticket": {
              "_t": "assocticket",
              "a": "system",
              "m": "test-idlerule-axe-idlerule-1-ticket"
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
                "_t": "assocticket",
                "a": "system",
                "user_id": "",
                "m": "test-idlerule-axe-idlerule-1-ticket"
              },
              {
                "_t": "assocticket",
                "a": "system",
                "user_id": "",
                "m": "test-idlerule-axe-idlerule-1-ticket"
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

  Scenario: given idle rule and no update for alarm should update alarm
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-2-name",
      "type": "alarm",
      "alarm_condition": "last_update",
      "enabled": true,
      "priority": 41,
      "duration": {
        "value": 5,
        "unit": "s"
      },
      "entity_patterns": [
        {
          "_id": "test-resource-axe-idlerule-2/test-component-axe-idlerule-2"
        }
      ],
      "operation": {
        "type": "assocticket",
        "parameters": {
          "ticket": "test-idlerule-axe-idlerule-2-ticket",
          "output": "test-idlerule-axe-idlerule-2-output"
        }
      }
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
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
    When I wait the end of event processing
    When I wait 3s
    When I send an event:
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
    When I wait the end of event processing
    When I wait 3s
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-axe-idlerule-2/test-component-axe-idlerule-2"}]}&with_steps=true
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
            "steps": [
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
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-axe-idlerule-2/test-component-axe-idlerule-2"}]}&with_steps=true
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
            "ticket": {
              "_t": "assocticket",
              "a": "system",
              "m": "test-idlerule-axe-idlerule-2-ticket"
            },
            "steps": [
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
                "m": "test-idlerule-axe-idlerule-2-ticket"
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
    When I send an event:
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
    When I wait the end of event processing
    When I wait 3s
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-axe-idlerule-2/test-component-axe-idlerule-2"}]}&with_steps=true
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
            "ticket": {
              "_t": "assocticket",
              "a": "system",
              "m": "test-idlerule-axe-idlerule-2-ticket"
            },
            "steps": [
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
                "m": "test-idlerule-axe-idlerule-2-ticket"
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
    When I send an event:
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
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-axe-idlerule-2/test-component-axe-idlerule-2"}]}&with_steps=true
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
            "ticket": {
              "_t": "assocticket",
              "a": "system",
              "m": "test-idlerule-axe-idlerule-2-ticket"
            },
            "steps": [
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
                "m": "test-idlerule-axe-idlerule-2-ticket"
              },
              {
                "_t": "stateinc",
                "val": 3
              },
              {
                "_t": "assocticket",
                "a": "system",
                "user_id": "",
                "m": "test-idlerule-axe-idlerule-2-ticket"
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
      "entity_patterns": [
        {
          "_id": "test-resource-axe-idlerule-3/test-component-axe-idlerule-3"
        }
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
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
    When I wait the end of event processing
    When I wait 3s
    When I send an event:
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
    When I wait the end of event processing
    When I wait 3s
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-axe-idlerule-3/test-component-axe-idlerule-3"}]}&with_steps=true
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
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-axe-idlerule-3/test-component-axe-idlerule-3"}]}&with_steps=true
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
              "m": "Idle rule test-idlerule-axe-idlerule-3-name"
            },
            "status": {
              "val": 5,
              "a": "system",
              "user_id": "",
              "m": "Idle rule test-idlerule-axe-idlerule-3-name"
            },
            "steps": [
              {
                "_t": "stateinc",
                "val": 3,
                "a": "system",
                "user_id": "",
                "m": "Idle rule test-idlerule-axe-idlerule-3-name"
              },
              {
                "_t": "statusinc",
                "val": 5,
                "a": "system",
                "user_id": "",
                "m": "Idle rule test-idlerule-axe-idlerule-3-name"
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
    When I send an event:
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
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-axe-idlerule-3/test-component-axe-idlerule-3"}]}&with_steps=true
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
                "_t": "statusdec",
                "a": "test-connector-axe-idlerule-3.test-connector-name-axe-idlerule-3",
                "m": "test-output-axe-idlerule-3",
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
    When I do DELETE /api/v4/idle-rules/{{ .ruleID }}
    Then the response code should be 204

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
      "entity_patterns": [
        {
          "_id": "test-component-axe-idlerule-4"
        }
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
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
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-component-axe-idlerule-4"}]}&with_steps=true
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
              "m": "Idle rule test-idlerule-axe-idlerule-4-name"
            },
            "status": {
              "val": 5,
              "a": "system",
              "m": "Idle rule test-idlerule-axe-idlerule-4-name"
            },
            "steps": [
              {
                "_t": "stateinc",
                "val": 3,
                "a": "system",
                "m": "Idle rule test-idlerule-axe-idlerule-4-name"
              },
              {
                "_t": "statusinc",
                "val": 5,
                "a": "system",
                "m": "Idle rule test-idlerule-axe-idlerule-4-name"
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
    When I send an event:
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
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-component-axe-idlerule-4"}]}&with_steps=true
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
                "_t": "statedec",
                "val": 0
              },
              {
                "_t": "statusdec",
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
    When I do DELETE /api/v4/idle-rules/{{ .ruleID }}
    Then the response code should be 204

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
      "entity_patterns": [
        {
          "_id": "test-connector-axe-idlerule-5/test-connector-name-axe-idlerule-5"
        }
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
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
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-connector-axe-idlerule-5/test-connector-name-axe-idlerule-5"}]}&with_steps=true
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
              "m": "Idle rule test-idlerule-axe-idlerule-5-name"
            },
            "status": {
              "val": 5,
              "a": "system",
              "m": "Idle rule test-idlerule-axe-idlerule-5-name"
            },
            "steps": [
              {
                "_t": "stateinc",
                "val": 3,
                "a": "system",
                "m": "Idle rule test-idlerule-axe-idlerule-5-name"
              },
              {
                "_t": "statusinc",
                "val": 5,
                "a": "system",
                "m": "Idle rule test-idlerule-axe-idlerule-5-name"
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
    When I wait 1s
    When I send an event:
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
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-connector-axe-idlerule-5/test-connector-name-axe-idlerule-5"}]}&with_steps=true
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
                "_t": "statedec",
                "val": 0
              },
              {
                "_t": "statusdec",
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
    When I do DELETE /api/v4/idle-rules/{{ .ruleID }}
    Then the response code should be 204

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
      "entity_patterns": [
        {
          "_id": "test-resource-axe-idlerule-6/test-component-axe-idlerule-6"
        }
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
      "name": "test-idlerule-axe-idlerule-6-2-name",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 45,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "entity_patterns": [
        {
          "_id": "test-resource-axe-idlerule-6/test-component-axe-idlerule-6"
        }
      ],
      "operation": {
        "type": "assocticket",
        "parameters": {
          "ticket": "test-idlerule-axe-idlerule-6-2-ticket",
          "output": "test-idlerule-axe-idlerule-6-2-output"
        }
      }
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
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
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-axe-idlerule-6/test-component-axe-idlerule-6"}]}&with_steps=true
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
            "ticket": {
              "_t": "assocticket",
              "a": "system",
              "m": "test-idlerule-axe-idlerule-6-2-ticket"
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
                "_t": "assocticket",
                "a": "system",
                "user_id": "",
                "m": "test-idlerule-axe-idlerule-6-2-ticket"
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
    When I wait 5s
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-axe-idlerule-6/test-component-axe-idlerule-6"}]}&with_steps=true
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
            "steps": [
              {
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "assocticket"
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
    When I send an event:
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
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-axe-idlerule-6/test-component-axe-idlerule-6"}]}&with_steps=true
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
            "steps": [
              {
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "assocticket"
              },
              {
                "_t": "assocticket"
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
      "entity_patterns": [
        {
          "_id": "test-resource-axe-idlerule-7/test-component-axe-idlerule-7"
        }
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
      "entity_patterns": [
        {
          "_id": "test-resource-axe-idlerule-7/test-component-axe-idlerule-7"
        }
      ]
    }
    """
    Then the response code should be 201
    Then I save response rule2ID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
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
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-axe-idlerule-7/test-component-axe-idlerule-7"}]}&with_steps=true
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
            },
            "steps": [
              {
                "_t": "stateinc",
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 5,
                "a": "system"
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
    When I wait 5s
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-axe-idlerule-7/test-component-axe-idlerule-7"}]}&with_steps=true
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
            "steps": [
              {
                "_t": "stateinc",
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 5
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
    When I send an event:
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
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-axe-idlerule-7/test-component-axe-idlerule-7"}]}&with_steps=true
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
            "steps": [
              {
                "_t": "stateinc",
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 5,
                "a": "system"
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
                "a": "system"
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
    When I do DELETE /api/v4/idle-rules/{{ .rule1ID }}
    Then the response code should be 204
    When I do DELETE /api/v4/idle-rules/{{ .rule2ID }}
    Then the response code should be 204

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
      "entity_patterns": [
        {
          "_id": "test-resource-axe-idlerule-8/test-component-axe-idlerule-8"
        }
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
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
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-axe-idlerule-8/test-component-axe-idlerule-8"}]}&with_steps=true
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
            },
            "steps": [
              {
                "_t": "stateinc",
                "val": 1,
                "a": "test-connector-axe-idlerule-8.test-connector-name-axe-idlerule-8"
              },
              {
                "_t": "statusinc",
                "val": 1,
                "a": "test-connector-axe-idlerule-8.test-connector-name-axe-idlerule-8"
              },
              {
                "_t": "stateinc",
                "val": 3,
                "a": "system"
              },
              {
                "_t": "statusinc",
                "val": 5,
                "a": "system"
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
    When I send an event:
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
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-axe-idlerule-8/test-component-axe-idlerule-8"}]}&with_steps=true
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
            },
            "steps": [
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
    When I do DELETE /api/v4/idle-rules/{{ .ruleID }}
    Then the response code should be 204

  Scenario: given idle rule and no events for alarm and entity should apply most priority alarm rule
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-9-1-name",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 50,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "entity_patterns": [
        {
          "_id": "test-resource-axe-idlerule-9/test-component-axe-idlerule-9"
        }
      ],
      "operation": {
        "type": "assocticket",
        "parameters": {
          "ticket": "test-idlerule-axe-idlerule-9-1-ticket",
          "output": "test-idlerule-axe-idlerule-9-1-output"
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
      "entity_patterns": [
        {
          "_id": "test-resource-axe-idlerule-9/test-component-axe-idlerule-9"
        }
      ]
    }
    """
    Then the response code should be 201
    Then I save response rule2ID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
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
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-axe-idlerule-9/test-component-axe-idlerule-9"}]}&with_steps=true
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
            "ticket": {
              "_t": "assocticket",
              "a": "system",
              "m": "test-idlerule-axe-idlerule-9-1-ticket"
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
                "_t": "assocticket",
                "a": "system",
                "user_id": "",
                "m": "test-idlerule-axe-idlerule-9-1-ticket"
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
    When I wait 5s
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-axe-idlerule-9/test-component-axe-idlerule-9"}]}&with_steps=true
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
            "steps": [
              {
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "assocticket"
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
    When I send an event:
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
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-axe-idlerule-9/test-component-axe-idlerule-9"}]}&with_steps=true
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
            "steps": [
              {
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "assocticket"
              },
              {
                "_t": "assocticket"
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
    When I do DELETE /api/v4/idle-rules/{{ .rule1ID }}
    Then the response code should be 204
    When I do DELETE /api/v4/idle-rules/{{ .rule2ID }}
    Then the response code should be 204

  Scenario: given idle rule and no events for alarm and entity should apply most priority entity rule
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-10-1-name",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 53,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "entity_patterns": [
        {
          "_id": "test-resource-axe-idlerule-10/test-component-axe-idlerule-10"
        }
      ],
      "operation": {
        "type": "ack",
        "parameters": {
          "output": "test-idlerule-axe-idlerule-10-1-output"
        }
      }
    }
    """
    Then the response code should be 201
    Then I save response rule1ID={{ .lastResponse._id }}
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-10-2-name",
      "type": "entity",
      "enabled": true,
      "priority": 52,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "entity_patterns": [
        {
          "_id": "test-resource-axe-idlerule-10/test-component-axe-idlerule-10"
        }
      ]
    }
    """
    Then the response code should be 201
    Then I save response rule2ID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-10",
      "connector_name": "test-connector-name-axe-idlerule-10",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-10",
      "resource": "test-resource-axe-idlerule-10",
      "state": 2,
      "output": "test-output-axe-idlerule-10"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-axe-idlerule-10/test-component-axe-idlerule-10"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-10",
            "connector": "test-connector-axe-idlerule-10",
            "connector_name": "test-connector-name-axe-idlerule-10",
            "resource": "test-resource-axe-idlerule-10",
            "state": {
              "val": 3
            },
            "status": {
              "val": 5
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
                "_t": "stateinc",
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 5
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
    When I wait 5s
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-axe-idlerule-10/test-component-axe-idlerule-10"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-10",
            "connector": "test-connector-axe-idlerule-10",
            "connector_name": "test-connector-name-axe-idlerule-10",
            "resource": "test-resource-axe-idlerule-10",
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
                "_t": "stateinc",
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 5
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
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-10",
      "connector_name": "test-connector-name-axe-idlerule-10",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-10",
      "resource": "test-resource-axe-idlerule-10",
      "state": 2,
      "output": "test-output-axe-idlerule-10"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-resource-axe-idlerule-10/test-component-axe-idlerule-10"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-10",
            "connector": "test-connector-axe-idlerule-10",
            "connector_name": "test-connector-name-axe-idlerule-10",
            "resource": "test-resource-axe-idlerule-10",
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
                "_t": "stateinc",
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 5
              },
              {
                "_t": "statedec",
                "val": 2
              },
              {
                "_t": "statusdec",
                "val": 1
              },
              {
                "_t": "stateinc",
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 5
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
    When I do DELETE /api/v4/idle-rules/{{ .rule1ID }}
    Then the response code should be 204
    When I do DELETE /api/v4/idle-rules/{{ .rule2ID }}
    Then the response code should be 204

  Scenario: given idle rule and no events for component which is created by resource event should create alarm
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-11-name",
      "type": "entity",
      "enabled": true,
      "priority": 54,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "entity_patterns": [
        {
          "_id": "test-component-axe-idlerule-11"
        }
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-11",
      "connector_name": "test-connector-name-axe-idlerule-11",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-11",
      "resource":  "test-resource-axe-idlerule-11",
      "state": 0,
      "output": "test-output-axe-idlerule-11"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-component-axe-idlerule-11"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-11",
            "connector": "test-connector-axe-idlerule-11",
            "connector_name": "test-connector-name-axe-idlerule-11",
            "state": {
              "val": 3
            },
            "status": {
              "val": 5
            },
            "steps": [
              {
                "_t": "stateinc",
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 5
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
    When I do DELETE /api/v4/idle-rules/{{ .ruleID }}
    Then the response code should be 204

  Scenario: given idle rule with ok changestate operation should update next alarm
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-12-name",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 40,
      "duration": {
        "value": 2,
        "unit": "s"
      },
      "entity_patterns": [
        {
          "_id": "test-resource-axe-idlerule-12/test-component-axe-idlerule-12"
        }
      ],
      "operation": {
        "type": "changestate",
        "parameters": {
          "state": 0
        }
      }
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-12",
      "connector_name": "test-connector-name-axe-idlerule-12",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-12",
      "resource": "test-resource-axe-idlerule-12",
      "state": 2,
      "output": "test-output-axe-idlerule-12"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-12
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-12",
            "connector": "test-connector-axe-idlerule-12",
            "connector_name": "test-connector-name-axe-idlerule-12",
            "resource": "test-resource-axe-idlerule-12",
            "state": {
              "_t": "changestate",
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
    When I send an event:
    """json
    {
      "event_type": "resolve_close",
      "connector": "test-connector-axe-idlerule-12",
      "connector_name": "test-connector-name-axe-idlerule-12",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-12",
      "resource": "test-resource-axe-idlerule-12",
      "output": "test-output-axe-idlerule-12"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-12",
      "connector_name": "test-connector-name-axe-idlerule-12",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-12",
      "resource": "test-resource-axe-idlerule-12",
      "state": 2,
      "output": "test-output-axe-idlerule-12"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-12
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-12",
            "connector": "test-connector-axe-idlerule-12",
            "connector_name": "test-connector-name-axe-idlerule-12",
            "resource": "test-resource-axe-idlerule-12",
            "state": {
              "_t": "changestate",
              "val": 0
            },
            "status": {
              "val": 0
            }
          }
        },
        {
          "v": {
            "component": "test-component-axe-idlerule-12",
            "connector": "test-connector-axe-idlerule-12",
            "connector_name": "test-connector-name-axe-idlerule-12",
            "resource": "test-resource-axe-idlerule-12",
            "state": {
              "_t": "changestate",
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
        "total_count": 2
      }
    }
    """

  Scenario: given idle rule with cancel operation should update next alarm
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-13-name",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 40,
      "duration": {
        "value": 2,
        "unit": "s"
      },
      "entity_patterns": [
        {
          "_id": "test-resource-axe-idlerule-13/test-component-axe-idlerule-13"
        }
      ],
      "operation": {
        "type": "cancel"
      }
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-13",
      "connector_name": "test-connector-name-axe-idlerule-13",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-13",
      "resource": "test-resource-axe-idlerule-13",
      "state": 2,
      "output": "test-output-axe-idlerule-13"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-13
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-13",
            "connector": "test-connector-axe-idlerule-13",
            "connector_name": "test-connector-name-axe-idlerule-13",
            "resource": "test-resource-axe-idlerule-13",
            "canceled": {
              "_t": "cancel"
            },
            "state": {
              "val": 2
            },
            "status": {
              "val": 4
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
    When I send an event:
    """json
    {
      "event_type": "resolve_cancel",
      "connector": "test-connector-axe-idlerule-13",
      "connector_name": "test-connector-name-axe-idlerule-13",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-13",
      "resource": "test-resource-axe-idlerule-13",
      "output": "test-output-axe-idlerule-13"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-13",
      "connector_name": "test-connector-name-axe-idlerule-13",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-13",
      "resource": "test-resource-axe-idlerule-13",
      "state": 2,
      "output": "test-output-axe-idlerule-13"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-13
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-13",
            "connector": "test-connector-axe-idlerule-13",
            "connector_name": "test-connector-name-axe-idlerule-13",
            "resource": "test-resource-axe-idlerule-13",
            "canceled": {
              "_t": "cancel"
            },
            "state": {
              "val": 2
            },
            "status": {
              "val": 4
            }
          }
        },
        {
          "v": {
            "component": "test-component-axe-idlerule-13",
            "connector": "test-connector-axe-idlerule-13",
            "connector_name": "test-connector-name-axe-idlerule-13",
            "resource": "test-resource-axe-idlerule-13",
            "canceled": {
              "_t": "cancel"
            },
            "state": {
              "val": 2
            },
            "status": {
              "val": 4
            }
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
