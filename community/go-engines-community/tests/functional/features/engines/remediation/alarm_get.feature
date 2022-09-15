Feature: update an instruction statistics
  I need to be able to update an instruction statistics

  Scenario: given auto instruction execution should return flags in alarm API
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-1",
      "connector_name": "test-connector-name-to-alarm-instruction-get-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-1",
      "resource": "test-resource-to-alarm-instruction-get-1",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-1"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-1&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_auto_instruction_running": true,
          "is_all_auto_instructions_completed": false,
          "is_auto_instruction_failed": false
        }
      ]
    }
    """
    When I wait the end of 4 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-1&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_auto_instruction_running": true,
          "is_all_auto_instructions_completed": false,
          "is_auto_instruction_failed": false
        }
      ]
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-1&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_auto_instruction_running": false,
          "is_all_auto_instructions_completed": true,
          "is_auto_instruction_failed": false
        }
      ]
    }
    """

  Scenario: given manual instruction execution should return flags in alarm API
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-2",
      "connector_name": "test-connector-name-to-alarm-instruction-get-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-2",
      "resource": "test-resource-to-alarm-instruction-get-2",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-2"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-2
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-2"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-2&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_manual_instruction_running": true,
          "is_manual_instruction_waiting_result": false
        }
      ]
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-2&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_manual_instruction_running": false,
          "is_manual_instruction_waiting_result": true
        }
      ]
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-2&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_manual_instruction_running": false,
          "is_manual_instruction_waiting_result": false
        }
      ]
    }
    """

  Scenario: given auto failed instruction execution should return flags in alarm API
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-3",
      "connector_name": "test-connector-name-to-alarm-instruction-get-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-3",
      "resource": "test-resource-to-alarm-instruction-get-3",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-3"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-3&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_auto_instruction_running": true,
          "is_all_auto_instructions_completed": false,
          "is_auto_instruction_failed": false
        }
      ]
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-3&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_auto_instruction_running": true,
          "is_all_auto_instructions_completed": false,
          "is_auto_instruction_failed": true
        }
      ]
    }
    """
    When I wait the end of 4 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-3&with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_auto_instruction_running": false,
          "is_all_auto_instructions_completed": true,
          "is_auto_instruction_failed": true
        }
      ]
    }
    """

  Scenario: given get request should return assigned instruction for the alarm
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-4",
      "connector_name": "test-connector-name-to-alarm-instruction-get-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-4",
      "resource": "test-resource-to-alarm-instruction-get-4",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-4"
    }
    """
    When I wait the end of 4 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-4&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-4"
          },
          "assigned_instructions": [
            {
              "_id": "test-instruction-to-alarm-instruction-get-4-1",
              "name": "test-instruction-to-alarm-instruction-get-4-1-name",
              "execution": null
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

  Scenario: given get request should return assigned instruction, which have an execution for the alarm
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-5",
      "connector_name": "test-connector-name-to-alarm-instruction-get-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-5",
      "resource": "test-resource-to-alarm-instruction-get-5",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-5"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-5
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-5-1"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-5&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-5"
          },
          "assigned_instructions": [
            {
              "_id": "test-instruction-to-alarm-instruction-get-5-1",
              "name": "test-instruction-to-alarm-instruction-get-5-1-name",
              "execution": {
                "_id": "{{ .executionID }}",
                "status": 0
              }
            },
            {
              "_id": "test-instruction-to-alarm-instruction-get-5-2",
              "name": "test-instruction-to-alarm-instruction-get-5-2-name",
              "execution": null
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
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-5&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-5"
          },
          "assigned_instructions": [
            {
              "_id": "test-instruction-to-alarm-instruction-get-5-1",
              "name": "test-instruction-to-alarm-instruction-get-5-1-name",
              "execution": null
            },
            {
              "_id": "test-instruction-to-alarm-instruction-get-5-2",
              "name": "test-instruction-to-alarm-instruction-get-5-2-name",
              "execution": null
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

  Scenario: given get request should return alarms with assigned instructions depending from exclude or include instructions fields
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-to-alarm-instruction-get-6",
        "connector_name": "test-connector-name-to-alarm-instruction-get-6",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-instruction-get-6",
        "resource": "test-resource-to-alarm-instruction-get-6-1",
        "state": 1,
        "output": "test-output-to-alarm-instruction-get-6"
      },
      {
        "connector": "test-connector-to-alarm-instruction-get-6",
        "connector_name": "test-connector-name-to-alarm-instruction-get-6",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-instruction-get-6",
        "resource": "test-resource-to-alarm-instruction-get-6-2",
        "state": 1,
        "output": "test-output-to-alarm-instruction-get-6"
      },
      {
        "connector": "test-connector-to-alarm-instruction-get-6",
        "connector_name": "test-connector-name-to-alarm-instruction-get-6",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-instruction-get-6",
        "resource": "test-resource-to-alarm-instruction-get-6-3",
        "state": 1,
        "output": "test-output-to-alarm-instruction-get-6"
      },
      {
        "connector": "test-connector-to-alarm-instruction-get-6",
        "connector_name": "test-connector-name-to-alarm-instruction-get-6",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-instruction-get-6",
        "resource": "test-resource-to-alarm-instruction-get-6-4",
        "state": 1,
        "output": "test-output-to-alarm-instruction-get-6"
      }
    ]
    """
    When I wait the end of 10 events processing
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-alarm-instruction-get-6",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-alarm-instruction-get-6"
            }
          },
          {
            "field": "type",
            "cond": {
              "type": "eq",
              "value": "resource"
            }
          }
        ]
      ],
      "instruction_pattern": [
        [
          {
            "field": "instructions",
            "cond": {
              "type": "has_one_of",
              "value": ["manual", "auto"]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/alarms?filter={{ .lastResponse._id }}&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-6-1"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-6-2"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-6-3"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-6-4"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 4
      }
    }
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-alarm-instruction-get-6",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-alarm-instruction-get-6"
            }
          },
          {
            "field": "type",
            "cond": {
              "type": "eq",
              "value": "resource"
            }
          }
        ]
      ],
      "instruction_pattern": [
        [
          {
            "field": "instructions",
            "cond": {
              "type": "has_one_of",
              "value": ["manual"]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/alarms?filter={{ .lastResponse._id }}&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-6-1"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-6-2"
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
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-alarm-instruction-get-6",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-alarm-instruction-get-6"
            }
          },
          {
            "field": "type",
            "cond": {
              "type": "eq",
              "value": "resource"
            }
          }
        ]
      ],
      "instruction_pattern": [
        [
          {
            "field": "instructions",
            "cond": {
              "type": "has_one_of",
              "value": ["auto"]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/alarms?filter={{ .lastResponse._id }}&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-6-3"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-6-4"
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
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-alarm-instruction-get-6",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-alarm-instruction-get-6"
            }
          },
          {
            "field": "type",
            "cond": {
              "type": "eq",
              "value": "resource"
            }
          }
        ]
      ],
      "instruction_pattern": [
        [
          {
            "field": "instructions",
            "cond": {
              "type": "has_not",
              "value": ["manual", "auto"]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/alarms?filter={{ .lastResponse._id }}&sort_by=v.resource&sort=asc
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
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-alarm-instruction-get-6",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-alarm-instruction-get-6"
            }
          },
          {
            "field": "type",
            "cond": {
              "type": "eq",
              "value": "resource"
            }
          }
        ]
      ],
      "instruction_pattern": [
        [
          {
            "field": "instructions",
            "cond": {
              "type": "has_not",
              "value": ["manual"]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/alarms?filter={{ .lastResponse._id }}&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-6-3"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-6-4"
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
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-alarm-instruction-get-6",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-alarm-instruction-get-6"
            }
          },
          {
            "field": "type",
            "cond": {
              "type": "eq",
              "value": "resource"
            }
          }
        ]
      ],
      "instruction_pattern": [
        [
          {
            "field": "instructions",
            "cond": {
              "type": "has_not",
              "value": ["auto"]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/alarms?filter={{ .lastResponse._id }}&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-6-1"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-6-2"
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
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-alarm-instruction-get-6",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-alarm-instruction-get-6"
            }
          },
          {
            "field": "type",
            "cond": {
              "type": "eq",
              "value": "resource"
            }
          }
        ]
      ],
      "instruction_pattern": [
        [
          {
            "field": "instructions",
            "cond": {
              "type": "has_one_of",
              "value": [
                "test-instruction-to-alarm-instruction-get-6-1",
                "test-instruction-to-alarm-instruction-get-6-2"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/alarms?filter={{ .lastResponse._id }}&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-6-1"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-6-2"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-6-3"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-6-4"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 4
      }
    }
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-alarm-instruction-get-6",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-alarm-instruction-get-6"
            }
          },
          {
            "field": "type",
            "cond": {
              "type": "eq",
              "value": "resource"
            }
          }
        ]
      ],
      "instruction_pattern": [
        [
          {
            "field": "instructions",
            "cond": {
              "type": "has_one_of",
              "value": [
                "test-instruction-to-alarm-instruction-get-6-1"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/alarms?filter={{ .lastResponse._id }}&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-6-1"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-6-2"
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
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-alarm-instruction-get-6",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-alarm-instruction-get-6"
            }
          },
          {
            "field": "type",
            "cond": {
              "type": "eq",
              "value": "resource"
            }
          }
        ]
      ],
      "instruction_pattern": [
        [
          {
            "field": "instructions",
            "cond": {
              "type": "has_one_of",
              "value": ["manual", "auto"]
            }
          },
          {
            "field": "instructions",
            "cond": {
              "type": "has_not",
              "value": [
                "test-instruction-to-alarm-instruction-get-6-1"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/alarms?filter={{ .lastResponse._id }}&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-6-3"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-6-4"
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
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-alarm-instruction-get-6",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-alarm-instruction-get-6"
            }
          },
          {
            "field": "type",
            "cond": {
              "type": "eq",
              "value": "resource"
            }
          }
        ]
      ],
      "instruction_pattern": [
        [
          {
            "field": "instructions",
            "cond": {
              "type": "has_one_of",
              "value": ["manual", "auto"]
            }
          },
          {
            "field": "instructions",
            "cond": {
              "type": "has_not",
              "value": [
                "test-instruction-to-alarm-instruction-get-6-1",
                "test-instruction-to-alarm-instruction-get-6-2"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/alarms?filter={{ .lastResponse._id }}&sort_by=v.resource&sort=asc
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

  Scenario: given get request should return assigned instruction with old pattern for the alarm
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-to-alarm-instruction-get-7",
        "connector_name": "test-connector-name-to-alarm-instruction-get-7",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-instruction-get-7",
        "resource": "test-resource-to-alarm-instruction-get-7-1",
        "state": 1,
        "output": "test-output-to-alarm-instruction-get-7"
      },
      {
        "connector": "test-connector-to-alarm-instruction-get-7",
        "connector_name": "test-connector-name-to-alarm-instruction-get-7",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-instruction-get-7",
        "resource": "test-resource-to-alarm-instruction-get-7-2",
        "state": 1,
        "output": "test-output-to-alarm-instruction-get-7"
      }
    ]
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-7&with_instructions=true&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-7-1"
          },
          "assigned_instructions": [
            {
              "_id": "test-instruction-to-alarm-instruction-get-7-1",
              "name": "test-instruction-to-alarm-instruction-get-7-1-name",
              "execution": null
            }
          ]
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-7-2"
          },
          "assigned_instructions": [
            {
              "_id": "test-instruction-to-alarm-instruction-get-7-2",
              "name": "test-instruction-to-alarm-instruction-get-7-2-name",
              "execution": null
            }
          ]
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

  Scenario: given get request should not return assigned instruction without patterns for the alarm
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-8",
      "connector_name": "test-connector-name-to-alarm-instruction-get-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-8",
      "resource": "test-resource-to-alarm-instruction-get-8",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-8"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-8&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-8"
          },
          "assigned_instructions": []
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

  Scenario: given get request should return alarms with assigned instructions with old pattern depending from exclude or include instructions fields
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-to-alarm-instruction-get-9",
        "connector_name": "test-connector-name-to-alarm-instruction-get-9",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-instruction-get-9",
        "resource": "test-resource-to-alarm-instruction-get-9-1",
        "state": 1,
        "output": "test-output-to-alarm-instruction-get-9"
      },
      {
        "connector": "test-connector-to-alarm-instruction-get-9",
        "connector_name": "test-connector-name-to-alarm-instruction-get-9",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-instruction-get-9",
        "resource": "test-resource-to-alarm-instruction-get-9-2",
        "state": 1,
        "output": "test-output-to-alarm-instruction-get-9"
      }
    ]
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-alarm-instruction-get-9",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-alarm-instruction-get-9"
            }
          },
          {
            "field": "type",
            "cond": {
              "type": "eq",
              "value": "resource"
            }
          }
        ]
      ],
      "instruction_pattern": [
        [
          {
            "field": "instructions",
            "cond": {
              "type": "has_one_of",
              "value": [
                "test-instruction-to-alarm-instruction-get-9-1",
                "test-instruction-to-alarm-instruction-get-9-2"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/alarms?filter={{ .lastResponse._id }}&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-9-1"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-9-2"
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
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-alarm-instruction-get-9",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-alarm-instruction-get-9"
            }
          },
          {
            "field": "type",
            "cond": {
              "type": "eq",
              "value": "resource"
            }
          }
        ]
      ],
      "instruction_pattern": [
        [
          {
            "field": "instructions",
            "cond": {
              "type": "has_one_of",
              "value": [
                "test-instruction-to-alarm-instruction-get-9-1"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/alarms?filter={{ .lastResponse._id }}&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-9-1"
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
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-alarm-instruction-get-9",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-alarm-instruction-get-9"
            }
          },
          {
            "field": "type",
            "cond": {
              "type": "eq",
              "value": "resource"
            }
          }
        ]
      ],
      "instruction_pattern": [
        [
          {
            "field": "instructions",
            "cond": {
              "type": "has_one_of",
              "value": [
                "test-instruction-to-alarm-instruction-get-9-2"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/alarms?filter={{ .lastResponse._id }}&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-9-2"
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
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-alarm-instruction-get-9",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-alarm-instruction-get-9"
            }
          },
          {
            "field": "type",
            "cond": {
              "type": "eq",
              "value": "resource"
            }
          }
        ]
      ],
      "instruction_pattern": [
        [
          {
            "field": "instructions",
            "cond": {
              "type": "has_not",
              "value": [
                "test-instruction-to-alarm-instruction-get-9-2"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/alarms?filter={{ .lastResponse._id }}&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-9-1"
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
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-alarm-instruction-get-9",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-alarm-instruction-get-9"
            }
          },
          {
            "field": "type",
            "cond": {
              "type": "eq",
              "value": "resource"
            }
          }
        ]
      ],
      "instruction_pattern": [
        [
          {
            "field": "instructions",
            "cond": {
              "type": "has_not",
              "value": [
                "test-instruction-to-alarm-instruction-get-9-1"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/alarms?filter={{ .lastResponse._id }}&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-9-2"
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

  Scenario: given get request should not return alarms by instructions without patterns
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-to-alarm-instruction-get-10",
        "connector_name": "test-connector-name-to-alarm-instruction-get-10",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-instruction-get-10",
        "resource": "test-resource-to-alarm-instruction-get-10",
        "state": 1,
        "output": "test-output-to-alarm-instruction-get-10"
      }
    ]
    """
    When I wait the end of event processing
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-alarm-instruction-get-10",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-alarm-instruction-get-10"
            }
          },
          {
            "field": "type",
            "cond": {
              "type": "eq",
              "value": "resource"
            }
          }
        ]
      ],
      "instruction_pattern": [
        [
          {
            "field": "instructions",
            "cond": {
              "type": "has_one_of",
              "value": [
                "test-instruction-without-patterns"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/alarms?filter={{ .lastResponse._id }}&sort_by=v.resource&sort=asc
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
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-alarm-instruction-get-10",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-to-alarm-instruction-get-10"
            }
          },
          {
            "field": "type",
            "cond": {
              "type": "eq",
              "value": "resource"
            }
          }
        ]
      ],
      "instruction_pattern": [
        [
          {
            "field": "instructions",
            "cond": {
              "type": "has_not",
              "value": [
                "test-instruction-without-patterns"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/alarms?filter={{ .lastResponse._id }}&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-10"
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

  Scenario: given auto instruction execution should return alarm
    When I am admin
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-alarm-instruction-get-11",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-alarm-instruction-get-11"
            }
          }
        ]
      ],
      "instruction_pattern": [
        [
          {
            "field": "instruction_execution",
            "cond": {
              "type": "exist",
              "value": true
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterId1={{ .lastResponse._id }}
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-alarm-instruction-get-11",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-alarm-instruction-get-11"
            }
          }
        ]
      ],
      "instruction_pattern": [
        [
          {
            "field": "instruction_execution",
            "cond": {
              "type": "exist",
              "value": false
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterId2={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-11",
      "connector_name": "test-connector-name-to-alarm-instruction-get-11",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-11",
      "resource": "test-resource-to-alarm-instruction-get-11",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-11"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={{ .filterId1 }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-11"
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
    When I do GET /api/v4/alarms?filter={{ .filterId2 }}
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
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={{ .filterId1 }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-11"
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
    When I do GET /api/v4/alarms?filter={{ .filterId2 }}
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
    When I do GET /api/v4/alarms?filter={{ .filterId1 }} until response code is 200 and body is:
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
    When I do GET /api/v4/alarms?filter={{ .filterId2 }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-11"
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

  Scenario: given manual instruction execution should return flags in alarm API
    When I am admin
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-alarm-instruction-get-12",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-alarm-instruction-get-12"
            }
          }
        ]
      ],
      "instruction_pattern": [
        [
          {
            "field": "instruction_execution",
            "cond": {
              "type": "exist",
              "value": true
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterId1={{ .lastResponse._id }}
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-alarm-instruction-get-12",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-alarm-instruction-get-12"
            }
          }
        ]
      ],
      "instruction_pattern": [
        [
          {
            "field": "instruction_execution",
            "cond": {
              "type": "exist",
              "value": false
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterId2={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-12",
      "connector_name": "test-connector-name-to-alarm-instruction-get-12",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-12",
      "resource": "test-resource-to-alarm-instruction-get-12",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-12"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-12
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-12"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={{ .filterId1 }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-12"
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
    When I do GET /api/v4/alarms?filter={{ .filterId2 }}
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
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={{ .filterId1 }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-12"
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
    When I do GET /api/v4/alarms?filter={{ .filterId2 }}
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
    When I do GET /api/v4/alarms?filter={{ .filterId1 }} until response code is 200 and body is:
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
    When I do GET /api/v4/alarms?filter={{ .filterId2 }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-12"
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
