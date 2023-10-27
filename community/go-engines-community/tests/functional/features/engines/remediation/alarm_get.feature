Feature: update an instruction statistics
  I need to be able to update an instruction statistics

  @concurrent
  Scenario: given get request should return assigned instruction, which have an execution for the alarm
    When I am admin
    When I send an event and wait the end of event processing:
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
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-1
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-1-1"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "connector": "test-connector-to-alarm-instruction-get-1",
      "connector_name": "test-connector-name-to-alarm-instruction-get-1",
      "component": "test-component-to-alarm-instruction-get-1",
      "resource": "test-resource-to-alarm-instruction-get-1",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-1&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-1"
          },
          "assigned_instructions": [
            {
              "_id": "test-instruction-to-alarm-instruction-get-1-1",
              "name": "test-instruction-to-alarm-instruction-get-1-1-name",
              "type": 0,
              "execution": {
                "_id": "{{ .executionID }}",
                "status": 0
              }
            },
            {
              "_id": "test-instruction-to-alarm-instruction-get-1-2",
              "name": "test-instruction-to-alarm-instruction-get-1-2-name",
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
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructioncompleted",
      "connector": "test-connector-to-alarm-instruction-get-1",
      "connector_name": "test-connector-name-to-alarm-instruction-get-1",
      "component": "test-component-to-alarm-instruction-get-1",
      "resource": "test-resource-to-alarm-instruction-get-1",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-1&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-1"
          },
          "assigned_instructions": [
            {
              "_id": "test-instruction-to-alarm-instruction-get-1-1",
              "name": "test-instruction-to-alarm-instruction-get-1-1-name",
              "type": 0,
              "execution": null
            },
            {
              "_id": "test-instruction-to-alarm-instruction-get-1-2",
              "name": "test-instruction-to-alarm-instruction-get-1-2-name",
              "type": 0,
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

  @concurrent
  Scenario: given get request should return alarms with assigned instructions depending from exclude or include instructions fields
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-to-alarm-instruction-get-2",
        "connector_name": "test-connector-name-to-alarm-instruction-get-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-instruction-get-2",
        "resource": "test-resource-to-alarm-instruction-get-2-1",
        "state": 1,
        "output": "test-output-to-alarm-instruction-get-2"
      },
      {
        "connector": "test-connector-to-alarm-instruction-get-2",
        "connector_name": "test-connector-name-to-alarm-instruction-get-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-instruction-get-2",
        "resource": "test-resource-to-alarm-instruction-get-2-2",
        "state": 1,
        "output": "test-output-to-alarm-instruction-get-2"
      },
      {
        "connector": "test-connector-to-alarm-instruction-get-2",
        "connector_name": "test-connector-name-to-alarm-instruction-get-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-instruction-get-2",
        "resource": "test-resource-to-alarm-instruction-get-2-3",
        "state": 1,
        "output": "test-output-to-alarm-instruction-get-2"
      },
      {
        "connector": "test-connector-to-alarm-instruction-get-2",
        "connector_name": "test-connector-name-to-alarm-instruction-get-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-instruction-get-2",
        "resource": "test-resource-to-alarm-instruction-get-2-4",
        "state": 1,
        "output": "test-output-to-alarm-instruction-get-2"
      }
    ]
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-instruction-get-2",
        "connector_name": "test-connector-name-to-alarm-instruction-get-2",
        "component": "test-component-to-alarm-instruction-get-2",
        "resource": "test-resource-to-alarm-instruction-get-2-1",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-instruction-get-2",
        "connector_name": "test-connector-name-to-alarm-instruction-get-2",
        "component": "test-component-to-alarm-instruction-get-2",
        "resource": "test-resource-to-alarm-instruction-get-2-2",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-instruction-get-2",
        "connector_name": "test-connector-name-to-alarm-instruction-get-2",
        "component": "test-component-to-alarm-instruction-get-2",
        "resource": "test-resource-to-alarm-instruction-get-2-3",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-instruction-get-2",
        "connector_name": "test-connector-name-to-alarm-instruction-get-2",
        "component": "test-component-to-alarm-instruction-get-2",
        "resource": "test-resource-to-alarm-instruction-get-2-4",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-2",
        "connector_name": "test-connector-name-to-alarm-instruction-get-2",
        "component": "test-component-to-alarm-instruction-get-2",
        "resource": "test-resource-to-alarm-instruction-get-2-3",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-2",
        "connector_name": "test-connector-name-to-alarm-instruction-get-2",
        "component": "test-component-to-alarm-instruction-get-2",
        "resource": "test-resource-to-alarm-instruction-get-2-3",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-2",
        "connector_name": "test-connector-name-to-alarm-instruction-get-2",
        "component": "test-component-to-alarm-instruction-get-2",
        "resource": "test-resource-to-alarm-instruction-get-2-4",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-2",
        "connector_name": "test-connector-name-to-alarm-instruction-get-2",
        "component": "test-component-to-alarm-instruction-get-2",
        "resource": "test-resource-to-alarm-instruction-get-2-4",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?instructions[]={"include_types":[0,1]}&search=test-resource-to-alarm-instruction-get-2&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-2-1"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-2-2"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-2-3"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-2-4"
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
    When I do GET /api/v4/alarms?instructions[]={"include_types":[0]}&search=test-resource-to-alarm-instruction-get-2&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-2-1"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-2-2"
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
    When I do GET /api/v4/alarms?instructions[]={"include_types":[1]}&search=test-resource-to-alarm-instruction-get-2&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-2-3"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-2-4"
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
    When I do GET /api/v4/alarms?instructions[]={"exclude_types":[0,1]}&search=test-resource-to-alarm-instruction-get-2&sort_by=v.resource&sort=asc
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
    When I do GET /api/v4/alarms?instructions[]={"exclude_types":[0]}&search=test-resource-to-alarm-instruction-get-2&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-2-3"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-2-4"
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
    When I do GET /api/v4/alarms?instructions[]={"exclude_types":[1]}&search=test-resource-to-alarm-instruction-get-2&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-2-1"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-2-2"
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
    When I do GET /api/v4/alarms?instructions[]={"include":["test-instruction-to-alarm-instruction-get-2-1","test-instruction-to-alarm-instruction-get-2-2"]}&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-2-1"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-2-2"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-2-3"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-2-4"
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
    When I do GET /api/v4/alarms?instructions[]={"include":["test-instruction-to-alarm-instruction-get-2-1"]}&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-2-1"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-2-2"
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
    When I do GET /api/v4/alarms?instructions[]={"include_types":[0,1],"exclude":["test-instruction-to-alarm-instruction-get-2-1"]}&search=test-resource-to-alarm-instruction-get-2&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-2-3"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-2-4"
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
    When I do GET /api/v4/alarms?instructions[]={"include_types":[0],"include_types":[1],"exclude":["test-instruction-to-alarm-instruction-get-2-1","test-instruction-to-alarm-instruction-get-2-2"]}&search=test-resource-to-alarm-instruction-get-2&sort_by=v.resource&sort=asc
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

  @concurrent
  Scenario: given get request should not return assigned instruction without patterns for the alarm
    When I am admin
    When I send an event and wait the end of event processing:
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

  @concurrent
  Scenario: given get request should not return alarms by instructions without patterns
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-to-alarm-instruction-get-6",
        "connector_name": "test-connector-name-to-alarm-instruction-get-6",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-instruction-get-6",
        "resource": "test-resource-to-alarm-instruction-get-6",
        "state": 1,
        "output": "test-output-to-alarm-instruction-get-6"
      }
    ]
    """
    When I do GET /api/v4/alarms?instructions[]={"include":["test-instruction-without-patterns"]}&sort_by=v.resource&sort=asc
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
    When I do GET /api/v4/alarms?instructions[]={"exclude":["test-instruction-without-patterns"]}&search=test-resource-to-alarm-instruction-get-6&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-6"
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
  Scenario: given auto instruction execution should return alarm
    When I am admin
    When I send an event and wait the end of event processing:
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
    When I do GET /api/v4/alarms?instructions[]={"include_types":[1],"running":true}&search=test-resource-to-alarm-instruction-get-7 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-7-1"
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
    When I do GET /api/v4/alarms?instructions[]={"include_types":[1],"running":false}&search=test-resource-to-alarm-instruction-get-7
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
    When I do GET /api/v4/alarms?instructions[]={"exclude_types":[1],"running":true}&search=test-resource-to-alarm-instruction-get-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-7-2"
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
    When I do GET /api/v4/alarms?instructions[]={"exclude_types":[1],"running":false}&search=test-resource-to-alarm-instruction-get-7&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-7-1"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-7-2"
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
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-7",
        "connector_name": "test-connector-name-to-alarm-instruction-get-7",
        "component": "test-component-to-alarm-instruction-get-7",
        "resource": "test-resource-to-alarm-instruction-get-7-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-to-alarm-instruction-get-7",
        "connector_name": "test-connector-name-to-alarm-instruction-get-7",
        "component": "test-component-to-alarm-instruction-get-7",
        "resource": "test-resource-to-alarm-instruction-get-7-1",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?instructions[]={"include_types":[1],"running":true}&search=test-resource-to-alarm-instruction-get-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-7-1"
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
    When I do GET /api/v4/alarms?instructions[]={"include_types":[1],"running":false}&search=test-resource-to-alarm-instruction-get-7
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
    When I do GET /api/v4/alarms?instructions[]={"include_types":[1],"running":true}&search=test-resource-to-alarm-instruction-get-7 until response code is 200 and body is:
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
    When I do GET /api/v4/alarms?instructions[]={"include_types":[1],"running":false}&search=test-resource-to-alarm-instruction-get-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-7-1"
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
    When I do GET /api/v4/alarms?instructions[]={"exclude_types":[1],"running":true}&search=test-resource-to-alarm-instruction-get-7&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-7-1"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-7-2"
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
    When I do GET /api/v4/alarms?instructions[]={"exclude_types":[1],"running":false}&search=test-resource-to-alarm-instruction-get-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-7-2"
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
  Scenario: given manual instruction execution should return flags in alarm API
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-to-alarm-instruction-get-8",
        "connector_name": "test-connector-name-to-alarm-instruction-get-8",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-instruction-get-8",
        "resource": "test-resource-to-alarm-instruction-get-8-1",
        "state": 1,
        "output": "test-output-to-alarm-instruction-get-8"
      },
      {
        "connector": "test-connector-to-alarm-instruction-get-8",
        "connector_name": "test-connector-name-to-alarm-instruction-get-8",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-instruction-get-8",
        "resource": "test-resource-to-alarm-instruction-get-8-2",
        "state": 1,
        "output": "test-output-to-alarm-instruction-get-8"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-8-1
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-alarm-instruction-get-8"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "connector": "test-connector-to-alarm-instruction-get-8",
      "connector_name": "test-connector-name-to-alarm-instruction-get-8",
      "component": "test-component-to-alarm-instruction-get-8",
      "resource": "test-resource-to-alarm-instruction-get-8-1",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?instructions[]={"include_types":[0],"running":true}&search=test-resource-to-alarm-instruction-get-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-8-1"
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
    When I do GET /api/v4/alarms?instructions[]={"include_types":[0],"running":false}&search=test-resource-to-alarm-instruction-get-8
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
    When I do GET /api/v4/alarms?instructions[]={"exclude_types":[0],"running":true}&search=test-resource-to-alarm-instruction-get-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-8-2"
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
    When I do GET /api/v4/alarms?instructions[]={"exclude_types":[0],"running":false}&search=test-resource-to-alarm-instruction-get-8&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-8-1"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-8-2"
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
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructioncompleted",
      "connector": "test-connector-to-alarm-instruction-get-8",
      "connector_name": "test-connector-name-to-alarm-instruction-get-8",
      "component": "test-component-to-alarm-instruction-get-8",
      "resource": "test-resource-to-alarm-instruction-get-8-1",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?instructions[]={"include_types":[0],"running":true}&search=test-resource-to-alarm-instruction-get-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-8-1"
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
    When I do GET /api/v4/alarms?instructions[]={"include_types":[0],"running":false}&search=test-resource-to-alarm-instruction-get-8
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
    When I do GET /api/v4/alarms?instructions[]={"include_types":[0],"running":true}&search=test-resource-to-alarm-instruction-get-8 until response code is 200 and body is:
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
    When I do GET /api/v4/alarms?instructions[]={"include_types":[0],"running":false}&search=test-resource-to-alarm-instruction-get-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-8-1"
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
    When I do GET /api/v4/alarms?instructions[]={"exclude_types":[0],"running":true}&search=test-resource-to-alarm-instruction-get-8&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-8-1"
          }
        },
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-8-2"
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
    When I do GET /api/v4/alarms?instructions[]={"exclude_types":[0],"running":false}&search=test-resource-to-alarm-instruction-get-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-8-2"
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
  Scenario: given get request should return assigned simplified manual instructions for the alarm
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-9",
      "connector_name": "test-connector-name-to-alarm-instruction-get-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-9",
      "resource": "test-resource-to-alarm-instruction-get-9",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-9"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-9&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-9"
          },
          "assigned_instructions": [
            {
              "_id": "test-instruction-to-alarm-instruction-get-9-1",
              "name": "test-instruction-to-alarm-instruction-get-9-1-name",
              "type": 2,
              "execution": null
            },
            {
              "_id": "test-instruction-to-alarm-instruction-get-9-2",
              "name": "test-instruction-to-alarm-instruction-get-9-2-name",
              "type": 2,
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
