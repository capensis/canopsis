Feature: Assigned declare tickets
  I need to be able get assigned declare ticket rules for alarms
  Only admin should be able get assigned declare ticket rules for alarms

  @concurrent
  Scenario: given get unauth request should not allow access
    When I do GET /api/v4/cat/declare-ticket-assigned
    Then the response code should be 401

  @concurrent
  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/declare-ticket-assigned
    Then the response code should be 403

  @concurrent
  Scenario: given get assigned declare ticket rules request should return assigned rules for given alarms
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-assigned-declare-ticket-connector-1",
      "connector_name": "test-assigned-declare-ticket-connector-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-assigned-declare-ticket-component-1",
      "resource": "test-assigned-declare-ticket-resource-1-1",
      "state": 2,
      "output": "test-assigned-declare-ticket-output",
      "long_output": "test-assigned-declare-ticket-long-output",
      "author": "test-assigned-declare-ticket-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-assigned-declare-ticket-connector-1",
      "connector_name": "test-assigned-declare-ticket-connector-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-assigned-declare-ticket-component-1",
      "resource": "test-assigned-declare-ticket-resource-1-2",
      "state": 2,
      "output": "test-assigned-declare-ticket-output",
      "long_output": "test-assigned-declare-ticket-long-output",
      "author": "test-assigned-declare-ticket-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-assigned-declare-ticket-resource-1-1&with_declare_tickets=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-assigned-declare-ticket-resource-1-1"
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
    Then the response key "assigned_declare_ticket_rules" should not exist
    When I save response alarmID1={{ (index .lastResponse.data 0)._id }}
    When I do GET /api/v4/alarms?search=test-assigned-declare-ticket-resource-1-2&with_declare_tickets=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-assigned-declare-ticket-resource-1-2"
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
    Then the response key "assigned_declare_ticket_rules" should not exist
    When I save response alarmID2={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-assigned-declare-ticket-rule-1-1",
      "system_name": "test-assigned-declare-ticket-rule-1-1-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "https://canopsis-test.com",
            "method": "GET",
            "auth": {
              "username": "test",
              "password": "test"
            },
            "skip_verify": true,
            "timeout": {
              "value": 30,
              "unit": "s"
            },
            "retry_count": 3,
            "retry_delay": {
              "value": 1,
              "unit": "s"
            }
          },
          "declare_ticket": {
            "is_regexp": false,
            "ticket_id": "_id",
            "ticket_url": "url",
            "ticket_custom": "custom",
            "empty_response": false
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
              "value": "test-assigned-declare-ticket-resource-1-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID1={{ .lastResponse._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-assigned-declare-ticket-rule-1-2",
      "system_name": "test-assigned-declare-ticket-rule-1-2-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "https://canopsis-test.com",
            "method": "GET",
            "auth": {
              "username": "test",
              "password": "test"
            },
            "skip_verify": true,
            "timeout": {
              "value": 30,
              "unit": "s"
            },
            "retry_count": 3,
            "retry_delay": {
              "value": 1,
              "unit": "s"
            }
          },
          "declare_ticket": {
            "is_regexp": false,
            "ticket_id": "_id",
            "ticket_url": "url",
            "ticket_custom": "custom",
            "empty_response": false
          },
          "stop_on_fail": true
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-assigned-declare-ticket-component-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID2={{ .lastResponse._id }}
    When I do GET /api/v4/alarms?search=test-assigned-declare-ticket-resource-1-1&with_declare_tickets=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-assigned-declare-ticket-resource-1-1"
          },
          "assigned_declare_ticket_rules": [
            {
              "_id": "{{ .ruleID1 }}",
              "name": "test-assigned-declare-ticket-rule-1-1"
            },
            {
              "_id": "{{ .ruleID2 }}",
              "name": "test-assigned-declare-ticket-rule-1-2"
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
    When I do GET /api/v4/alarms?search=test-assigned-declare-ticket-resource-1-2&with_declare_tickets=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-assigned-declare-ticket-resource-1-2"
          },
          "assigned_declare_ticket_rules": [
            {
              "_id": "{{ .ruleID2 }}",
              "name": "test-assigned-declare-ticket-rule-1-2"
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
    When I do GET /api/v4/cat/declare-ticket-assigned?alarms[]={{ .alarmID1 }}&alarms[]={{ .alarmID2 }}
    Then the response code should be 200
    Then the response array key "by_alarms.{{ .alarmID1 }}" should contain:
    """json
    [
      {
        "_id": "{{ .ruleID1 }}",
        "name": "test-assigned-declare-ticket-rule-1-1"
      },
      {
        "_id": "{{ .ruleID2 }}",
        "name": "test-assigned-declare-ticket-rule-1-2"
      }
    ]
    """
    Then the response array key "by_alarms.{{ .alarmID2 }}" should contain:
    """json
    [
      {
        "_id": "{{ .ruleID2 }}",
        "name": "test-assigned-declare-ticket-rule-1-2"
      }
    ]
    """
    Then the response array key "by_rules.{{ .ruleID1 }}.alarms" should contain:
    """json
    [
      "{{ .alarmID1 }}"
    ]
    """
    Then the response array key "by_rules.{{ .ruleID2 }}.alarms" should contain:
    """json
    [
      "{{ .alarmID1 }}",
      "{{ .alarmID2 }}"
    ]
    """

  @concurrent
  Scenario: given get assigned declare ticket rules request should not return assigned rules for given alarms
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-assigned-declare-ticket-connector-2",
      "connector_name": "test-assigned-declare-ticket-connector-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-assigned-declare-ticket-component-2",
      "resource": "test-assigned-declare-ticket-resource-2-1",
      "state": 2,
      "output": "test-assigned-declare-ticket-output",
      "long_output": "test-assigned-declare-ticket-long-output",
      "author": "test-assigned-declare-ticket-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-assigned-declare-ticket-connector-2",
      "connector_name": "test-assigned-declare-ticket-connector-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-assigned-declare-ticket-component-2",
      "resource": "test-assigned-declare-ticket-resource-2-2",
      "state": 2,
      "output": "test-assigned-declare-ticket-output",
      "long_output": "test-assigned-declare-ticket-long-output",
      "author": "test-assigned-declare-ticket-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-assigned-declare-ticket-resource-2-1&with_declare_tickets=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-assigned-declare-ticket-resource-2-1"
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
    Then the response key "assigned_declare_ticket_rules" should not exist
    When I save response alarmID1={{ (index .lastResponse.data 0)._id }}
    When I do GET /api/v4/alarms?search=test-assigned-declare-ticket-resource-2-2&with_declare_tickets=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-assigned-declare-ticket-resource-2-2"
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
    Then the response key "assigned_declare_ticket_rules" should not exist
    When I save response alarmID2={{ (index .lastResponse.data 0)._id }}
    When I do GET /api/v4/cat/declare-ticket-assigned?alarms[]={{ .alarmID1 }}&alarms[]={{ .alarmID2 }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "by_alarms": {
        "{{ .alarmID1 }}": [],
        "{{ .alarmID2 }}": []
      },
      "by_rules": {}
    }
    """

  @concurrent
  Scenario: given get assigned declare ticket rules request should not return assigned disabled rules
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-assigned-declare-ticket-connector-3",
      "connector_name": "test-assigned-declare-ticket-connector-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-assigned-declare-ticket-component-3",
      "resource": "test-assigned-declare-ticket-resource-3-1",
      "state": 2,
      "output": "test-assigned-declare-ticket-output",
      "long_output": "test-assigned-declare-ticket-long-output",
      "author": "test-assigned-declare-ticket-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-assigned-declare-ticket-connector-3",
      "connector_name": "test-assigned-declare-ticket-connector-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-assigned-declare-ticket-component-3",
      "resource": "test-assigned-declare-ticket-resource-3-2",
      "state": 2,
      "output": "test-assigned-declare-ticket-output",
      "long_output": "test-assigned-declare-ticket-long-output",
      "author": "test-assigned-declare-ticket-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-assigned-declare-ticket-resource-3-1&with_declare_tickets=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-assigned-declare-ticket-resource-3-1"
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
    Then the response key "assigned_declare_ticket_rules" should not exist
    When I save response alarmID1={{ (index .lastResponse.data 0)._id }}
    When I do GET /api/v4/alarms?search=test-assigned-declare-ticket-resource-3-2&with_declare_tickets=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-assigned-declare-ticket-resource-3-2"
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
    Then the response key "assigned_declare_ticket_rules" should not exist
    When I save response alarmID2={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-assigned-declare-ticket-rule-3-1",
      "system_name": "test-assigned-declare-ticket-rule-3-1-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "https://canopsis-test.com",
            "method": "GET",
            "auth": {
              "username": "test",
              "password": "test"
            },
            "skip_verify": true,
            "timeout": {
              "value": 30,
              "unit": "s"
            },
            "retry_count": 3,
            "retry_delay": {
              "value": 1,
              "unit": "s"
            }
          },
          "declare_ticket": {
            "is_regexp": false,
            "ticket_id": "_id",
            "ticket_url": "url",
            "ticket_custom": "custom",
            "empty_response": false
          },
          "stop_on_fail": true
        }
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-assigned-declare-ticket-resource-3-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID1={{ .lastResponse._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-assigned-declare-ticket-rule-3-2",
      "system_name": "test-assigned-declare-ticket-rule-3-2-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "https://canopsis-test.com",
            "method": "GET",
            "auth": {
              "username": "test",
              "password": "test"
            },
            "skip_verify": true,
            "timeout": {
              "value": 30,
              "unit": "s"
            },
            "retry_count": 3,
            "retry_delay": {
              "value": 1,
              "unit": "s"
            }
          },
          "declare_ticket": {
            "is_regexp": false,
            "ticket_id": "_id",
            "ticket_url": "url",
            "ticket_custom": "custom",
            "empty_response": false
          },
          "stop_on_fail": true
        }
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-assigned-declare-ticket-component-3"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID2={{ .lastResponse._id }}
    When I do GET /api/v4/alarms?search=test-assigned-declare-ticket-resource-3-1&with_declare_tickets=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-assigned-declare-ticket-resource-3-1"
          },
          "assigned_declare_ticket_rules": [
            {
              "_id": "{{ .ruleID1 }}",
              "name": "test-assigned-declare-ticket-rule-3-1"
            },
            {
              "_id": "{{ .ruleID2 }}",
              "name": "test-assigned-declare-ticket-rule-3-2"
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
    When I save response alarmID1={{ (index .lastResponse.data 0)._id }}
    When I do GET /api/v4/alarms?search=test-assigned-declare-ticket-resource-3-2&with_declare_tickets=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-assigned-declare-ticket-resource-3-2"
          },
          "assigned_declare_ticket_rules": [
            {
              "_id": "{{ .ruleID2 }}",
              "name": "test-assigned-declare-ticket-rule-3-2"
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
    When I do GET /api/v4/cat/declare-ticket-assigned?alarms[]={{ .alarmID1 }}&alarms[]={{ .alarmID2 }}
    Then the response code should be 200
    Then the response array key "by_alarms.{{ .alarmID1 }}" should contain:
    """json
    [
      {
        "_id": "{{ .ruleID1 }}",
        "name": "test-assigned-declare-ticket-rule-3-1"
      },
      {
        "_id": "{{ .ruleID2 }}",
        "name": "test-assigned-declare-ticket-rule-3-2"
      }
    ]
    """
    Then the response array key "by_alarms.{{ .alarmID2 }}" should contain:
    """json
    [
      {
        "_id": "{{ .ruleID2 }}",
        "name": "test-assigned-declare-ticket-rule-3-2"
      }
    ]
    """
    Then the response array key "by_rules.{{ .ruleID1 }}.alarms" should contain:
    """json
    [
      "{{ .alarmID1 }}"
    ]
    """
    Then the response array key "by_rules.{{ .ruleID2 }}.alarms" should contain:
    """json
    [
      "{{ .alarmID1 }}",
      "{{ .alarmID2 }}"
    ]
    """
    When I do PUT /api/v4/cat/declare-ticket-rules/{{ .ruleID2 }}:
    """json
    {
      "name": "test-assigned-declare-ticket-rule-3-2",
      "system_name": "test-assigned-declare-ticket-rule-3-2-name",
      "enabled": false,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "https://canopsis-test.com",
            "method": "GET",
            "auth": {
              "username": "test",
              "password": "test"
            },
            "skip_verify": true,
            "timeout": {
              "value": 30,
              "unit": "s"
            },
            "retry_count": 3,
            "retry_delay": {
              "value": 1,
              "unit": "s"
            }
          },
          "declare_ticket": {
            "is_regexp": false,
            "ticket_id": "_id",
            "ticket_url": "url",
            "ticket_custom": "custom",
            "empty_response": false
          },
          "stop_on_fail": true
        }
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-assigned-declare-ticket-component-3"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/alarms?search=test-assigned-declare-ticket-resource-3-1&with_declare_tickets=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-assigned-declare-ticket-resource-3-1"
          },
          "assigned_declare_ticket_rules": [
            {
              "_id": "{{ .ruleID1 }}",
              "name": "test-assigned-declare-ticket-rule-3-1"
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
    When I save response alarmID1={{ (index .lastResponse.data 0)._id }}
    When I do GET /api/v4/alarms?search=test-assigned-declare-ticket-resource-3-2&with_declare_tickets=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-assigned-declare-ticket-resource-3-2"
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
    Then the response key "assigned_declare_ticket_rules" should not exist
    When I do GET /api/v4/cat/declare-ticket-assigned?alarms[]={{ .alarmID1 }}&alarms[]={{ .alarmID2 }}
    Then the response code should be 200
    Then the response array key "by_alarms.{{ .alarmID1 }}" should contain:
    """json
    [
      {
        "_id": "{{ .ruleID1 }}",
        "name": "test-assigned-declare-ticket-rule-3-1"
      }
    ]
    """
    Then the response array key "by_alarms.{{ .alarmID2 }}" should contain:
    """json
    []
    """
    Then the response array key "by_rules.{{ .ruleID1 }}.alarms" should contain:
    """json
    [
      "{{ .alarmID1 }}"
    ]
    """

  @concurrent
  Scenario: given get assigned declare ticket rules request should return assigned rules by pbh patterns
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-assigned-declare-ticket-connector-4",
      "connector_name": "test-assigned-declare-ticket-connector-4-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-assigned-declare-ticket-component-4",
      "resource": "test-assigned-declare-ticket-resource-4-1",
      "state": 2,
      "output": "test-assigned-declare-ticket-output",
      "long_output": "test-assigned-declare-ticket-long-output",
      "author": "test-assigned-declare-ticket-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-assigned-declare-ticket-connector-4",
      "connector_name": "test-assigned-declare-ticket-connector-4-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-assigned-declare-ticket-component-4",
      "resource": "test-assigned-declare-ticket-resource-4-2",
      "state": 2,
      "output": "test-assigned-declare-ticket-output",
      "long_output": "test-assigned-declare-ticket-long-output",
      "author": "test-assigned-declare-ticket-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-assigned-declare-ticket-resource-4-1&with_declare_tickets=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-assigned-declare-ticket-resource-4-1"
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
    Then the response key "assigned_declare_ticket_rules" should not exist
    When I save response alarmID1={{ (index .lastResponse.data 0)._id }}
    When I do GET /api/v4/alarms?search=test-assigned-declare-ticket-resource-4-2&with_declare_tickets=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-assigned-declare-ticket-resource-4-2"
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
    Then the response key "assigned_declare_ticket_rules" should not exist
    When I save response alarmID2={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-assigned-declare-ticket-pbehavior-4-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-assigned-declare-ticket-resource-4-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response pbehaviorID={{ .lastResponse._id }}
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-assigned-declare-ticket-pbehavior-4-2",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-assigned-declare-ticket-resource-4-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "pbhenter",
        "resource": "test-assigned-declare-ticket-resource-4-1"
      },
      {
        "event_type": "pbhenter",
        "resource": "test-assigned-declare-ticket-resource-4-2"
      }
    ]
    """
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-assigned-declare-ticket-rule-4-1",
      "system_name": "test-assigned-declare-ticket-rule-4-1-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "https://canopsis-test.com",
            "method": "GET",
            "auth": {
              "username": "test",
              "password": "test"
            },
            "skip_verify": true,
            "timeout": {
              "value": 30,
              "unit": "s"
            },
            "retry_count": 3,
            "retry_delay": {
              "value": 1,
              "unit": "s"
            }
          },
          "declare_ticket": {
            "is_regexp": false,
            "ticket_id": "_id",
            "ticket_url": "url",
            "ticket_custom": "custom",
            "empty_response": false
          },
          "stop_on_fail": true
        }
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.id",
            "cond": {
              "type": "eq",
              "value": "{{ .pbehaviorID }}"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID1={{ .lastResponse._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-assigned-declare-ticket-rule-4-2",
      "system_name": "test-assigned-declare-ticket-rule-4-2-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "https://canopsis-test.com",
            "method": "GET",
            "auth": {
              "username": "test",
              "password": "test"
            },
            "skip_verify": true,
            "timeout": {
              "value": 30,
              "unit": "s"
            },
            "retry_count": 3,
            "retry_delay": {
              "value": 1,
              "unit": "s"
            }
          },
          "declare_ticket": {
            "is_regexp": false,
            "ticket_id": "_id",
            "ticket_url": "url",
            "ticket_custom": "custom",
            "empty_response": false
          },
          "stop_on_fail": true
        }
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.canonical_type",
            "cond": {
              "type": "eq",
              "value": "maintenance"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID2={{ .lastResponse._id }}
    When I do GET /api/v4/alarms?search=test-assigned-declare-ticket-resource-4-1&with_declare_tickets=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-assigned-declare-ticket-resource-4-1"
          },
          "assigned_declare_ticket_rules": [
            {
              "_id": "{{ .ruleID1 }}",
              "name": "test-assigned-declare-ticket-rule-4-1"
            },
            {
              "_id": "{{ .ruleID2 }}",
              "name": "test-assigned-declare-ticket-rule-4-2"
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
    When I save response alarmID1={{ (index .lastResponse.data 0)._id }}
    When I do GET /api/v4/alarms?search=test-assigned-declare-ticket-resource-4-2&with_declare_tickets=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-assigned-declare-ticket-resource-4-2"
          },
          "assigned_declare_ticket_rules": [
            {
              "_id": "{{ .ruleID2 }}",
              "name": "test-assigned-declare-ticket-rule-4-2"
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
    When I do GET /api/v4/cat/declare-ticket-assigned?alarms[]={{ .alarmID1 }}&alarms[]={{ .alarmID2 }}
    Then the response code should be 200
    Then the response array key "by_alarms.{{ .alarmID1 }}" should contain:
    """json
    [
      {
        "_id": "{{ .ruleID1 }}",
        "name": "test-assigned-declare-ticket-rule-4-1"
      },
      {
        "_id": "{{ .ruleID2 }}",
        "name": "test-assigned-declare-ticket-rule-4-2"
      }
    ]
    """
    Then the response array key "by_alarms.{{ .alarmID2 }}" should contain:
    """json
    [
      {
        "_id": "{{ .ruleID2 }}",
        "name": "test-assigned-declare-ticket-rule-4-2"
      }
    ]
    """
    Then the response array key "by_rules.{{ .ruleID1 }}.alarms" should contain:
    """json
    [
      "{{ .alarmID1 }}"
    ]
    """
    Then the response array key "by_rules.{{ .ruleID2 }}.alarms" should contain:
    """json
    [
      "{{ .alarmID1 }}",
      "{{ .alarmID2 }}"
    ]
    """
