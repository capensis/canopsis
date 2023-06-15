Feature: get service entities with assigned instructions
  I need to be able get service entities with assigned instructions

  @concurrent
  Scenario: given manual instruction execution should return assigned instructions in weather services entities API
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-get-assigned-instruction-in-weather-api-1-1-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-alarm-weather-widget-instructions-resource-1-1",
                "test-alarm-weather-widget-instructions-resource-1-2"
              ]
            }
          }
        ]
      ],
      "description": "test-instruction-get-assigned-instruction-in-weather-api-1-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-get-assigned-instruction-in-weather-api-1-1-step-1",
          "operations": [
            {
              "name": "test-instruction-get-assigned-instruction-in-weather-api-1-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-get-assigned-instruction-in-weather-api-1-1-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-get-assigned-instruction-in-weather-api-1-1-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionId1={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-get-assigned-instruction-in-weather-api-1-2-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-alarm-weather-widget-instructions-resource-1-2",
                "test-alarm-weather-widget-instructions-resource-1-3"
              ]
            }
          }
        ]
      ],
      "description": "test-instruction-get-assigned-instruction-in-weather-api-1-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-get-assigned-instruction-in-weather-api-1-2-step-1",
          "operations": [
            {
              "name": "test-instruction-get-assigned-instruction-in-weather-api-1-2-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-get-assigned-instruction-in-weather-api-1-2-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-get-assigned-instruction-in-weather-api-1-2-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionId2={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-alarm-weather-widget-instructions-connector-1",
      "connector_name": "test-alarm-weather-widget-instructions-connectorname-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-alarm-weather-widget-instructions-component-1",
      "resource": "test-alarm-weather-widget-instructions-resource-1-1",
      "state": 1,
      "output": "test-alarm-weather-widget-instructions-output-1"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-alarm-weather-widget-instructions-connector-1",
      "connector_name": "test-alarm-weather-widget-instructions-connectorname-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-alarm-weather-widget-instructions-component-1",
      "resource": "test-alarm-weather-widget-instructions-resource-1-2",
      "state": 1,
      "output": "test-alarm-weather-widget-instructions-output-1"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-alarm-weather-widget-instructions-connector-1",
      "connector_name": "test-alarm-weather-widget-instructions-connectorname-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-alarm-weather-widget-instructions-component-1",
      "resource": "test-alarm-weather-widget-instructions-resource-1-3",
      "state": 1,
      "output": "test-alarm-weather-widget-instructions-output-1"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-alarm-weather-widget-instructions-connector-1",
      "connector_name": "test-alarm-weather-widget-instructions-connectorname-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-alarm-weather-widget-instructions-component-1",
      "resource": "test-alarm-weather-widget-instructions-resource-1-4",
      "state": 1,
      "output": "test-alarm-weather-widget-instructions-output-1"
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-entity-instruction-weather-services-1",
      "name": "test-entity-instruction-weather-services-1",
      "output_template": "test-entity-instruction-weather-services-1",
      "category": "",
      "impact_level": 1,
      "enabled": true,
      "sli_avail_state": 1,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-alarm-weather-widget-instructions-resource-1-1",
                "test-alarm-weather-widget-instructions-resource-1-2",
                "test-alarm-weather-widget-instructions-resource-1-3",
                "test-alarm-weather-widget-instructions-resource-1-4"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response serviceId={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "component": "{{ .serviceId }}"
      },
      {
        "event_type": "activate",
        "component": "{{ .serviceId }}"
      }
    ]
    """
    When I do GET /api/v4/weather-services/test-entity-instruction-weather-services-1?with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-weather-widget-instructions-resource-1-1/test-alarm-weather-widget-instructions-component-1",
          "assigned_instructions": [
            {
              "_id": "{{ .instructionId1 }}"
            }
          ]
        },
        {
          "_id": "test-alarm-weather-widget-instructions-resource-1-2/test-alarm-weather-widget-instructions-component-1",
          "assigned_instructions": [
            {
              "_id": "{{ .instructionId1 }}"
            },
            {
              "_id": "{{ .instructionId2 }}"
            }
          ]
        },
        {
          "_id": "test-alarm-weather-widget-instructions-resource-1-3/test-alarm-weather-widget-instructions-component-1",
          "assigned_instructions": [
            {
              "_id": "{{ .instructionId2 }}"
            }
          ]
        },
        {
          "_id": "test-alarm-weather-widget-instructions-resource-1-4/test-alarm-weather-widget-instructions-component-1",
          "assigned_instructions": []
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

  @concurrent
  Scenario: given auto instruction execution should return auto instruction flags in weather API
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "name": "test-instruction-get-assigned-instruction-in-weather-api-2-1-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-alarm-weather-widget-instructions-resource-2"
            }
          }
        ]
      ],
      "description": "test-instruction-get-assigned-instruction-in-weather-api-2-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-7"
        }
      ],
      "triggers": ["create"]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "name": "test-instruction-get-assigned-instruction-in-weather-api-2-2-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-alarm-weather-widget-instructions-resource-2"
            }
          }
        ]
      ],
      "description": "test-instruction-get-assigned-instruction-in-weather-api-2-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-7"
        }
      ],
      "triggers": ["create"]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-alarm-weather-widget-instructions-connector-2",
      "connector_name": "test-alarm-weather-widget-instructions-connectorname-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-alarm-weather-widget-instructions-component-2",
      "resource": "test-alarm-weather-widget-instructions-resource-2",
      "state": 1,
      "output": "test-alarm-weather-widget-instructions-output-2"
    }
    """
    When I do GET /api/v4/alarms?search=test-alarm-weather-widget-instructions-resource-2&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 2
        }
      ]
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-entity-instruction-weather-services-2",
      "name": "test-entity-instruction-weather-services-2",
      "output_template": "test-entity-instruction-weather-services-2",
      "category": "",
      "impact_level": 1,
      "enabled": true,
      "sli_avail_state": 1,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-alarm-weather-widget-instructions-resource-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response serviceId={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "component": "{{ .serviceId }}"
      },
      {
        "event_type": "activate",
        "component": "{{ .serviceId }}"
      }
    ]
    """
    When I do GET /api/v4/weather-services/test-entity-instruction-weather-services-2?with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-weather-widget-instructions-resource-2/test-alarm-weather-widget-instructions-component-2",
          "instruction_execution_icon": 2
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
    When I wait 6s
    When I do GET /api/v4/weather-services/test-entity-instruction-weather-services-2?with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-weather-widget-instructions-resource-2/test-alarm-weather-widget-instructions-component-2",
          "instruction_execution_icon": 12
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
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "trigger",
        "connector": "test-alarm-weather-widget-instructions-connector-2",
        "connector_name": "test-alarm-weather-widget-instructions-connectorname-2",
        "component": "test-alarm-weather-widget-instructions-component-2",
        "resource": "test-alarm-weather-widget-instructions-resource-2",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-alarm-weather-widget-instructions-connector-2",
        "connector_name": "test-alarm-weather-widget-instructions-connectorname-2",
        "component": "test-alarm-weather-widget-instructions-component-2",
        "resource": "test-alarm-weather-widget-instructions-resource-2",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-alarm-weather-widget-instructions-connector-2",
        "connector_name": "test-alarm-weather-widget-instructions-connectorname-2",
        "component": "test-alarm-weather-widget-instructions-component-2",
        "resource": "test-alarm-weather-widget-instructions-resource-2",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-alarm-weather-widget-instructions-connector-2",
        "connector_name": "test-alarm-weather-widget-instructions-connectorname-2",
        "component": "test-alarm-weather-widget-instructions-component-2",
        "resource": "test-alarm-weather-widget-instructions-resource-2",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-alarm-weather-widget-instructions-connector-2",
        "connector_name": "test-alarm-weather-widget-instructions-connectorname-2",
        "component": "test-alarm-weather-widget-instructions-component-2",
        "resource": "test-alarm-weather-widget-instructions-resource-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/weather-services/test-entity-instruction-weather-services-2?with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-weather-widget-instructions-resource-2/test-alarm-weather-widget-instructions-component-2",
          "instruction_execution_icon": 10
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
  Scenario: given manual instruction execution should return manual instruction flags in weather API
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-get-assigned-instruction-in-weather-api-3-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-alarm-weather-widget-instructions-resource-3"
            }
          }
        ]
      ],
      "description": "test-instruction-get-assigned-instruction-in-weather-api-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-get-assigned-instruction-in-weather-api-3-step-1",
          "operations": [
            {
              "name": "test-instruction-get-assigned-instruction-in-weather-api-3-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-get-assigned-instruction-in-weather-api-3-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-get-assigned-instruction-in-weather-api-3-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionId={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-alarm-weather-widget-instructions-connector-3",
      "connector_name": "test-alarm-weather-widget-instructions-connectorname-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-alarm-weather-widget-instructions-component-3",
      "resource": "test-alarm-weather-widget-instructions-resource-3",
      "state": 1,
      "output": "test-alarm-weather-widget-instructions-output-3"
    }
    """
    When I do GET /api/v4/alarms?search=test-alarm-weather-widget-instructions-resource-3
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-entity-instruction-weather-services-3",
      "name": "test-entity-instruction-weather-services-3",
      "output_template": "test-entity-instruction-weather-services-3",
      "category": "",
      "impact_level": 1,
      "enabled": true,
      "sli_avail_state": 1,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-alarm-weather-widget-instructions-resource-3"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response serviceId={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "component": "{{ .serviceId }}"
      },
      {
        "event_type": "activate",
        "component": "{{ .serviceId }}"
      }
    ]
    """
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmId }}",
      "instruction": "{{ .instructionId }}"
    }
    """
    Then the response code should be 200
    When I save response executionId={{ .lastResponse._id }}
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructionstarted",
      "connector": "test-alarm-weather-widget-instructions-connector-3",
      "connector_name": "test-alarm-weather-widget-instructions-connectorname-3",
      "component": "test-alarm-weather-widget-instructions-component-3",
      "resource": "test-alarm-weather-widget-instructions-resource-3",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/weather-services/test-entity-instruction-weather-services-3?with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 1
        }
      ]
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .executionId }}/next-step
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "instructioncompleted",
      "connector": "test-alarm-weather-widget-instructions-connector-3",
      "connector_name": "test-alarm-weather-widget-instructions-connectorname-3",
      "component": "test-alarm-weather-widget-instructions-component-3",
      "resource": "test-alarm-weather-widget-instructions-resource-3",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/weather-services/test-entity-instruction-weather-services-3?with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 1
        }
      ]
    }
    """
    When I do GET /api/v4/weather-services/test-entity-instruction-weather-services-3?with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 11
        }
      ]
    }
    """

  @concurrent
  Scenario: given auto failed instruction execution should return auto instruction flags in weather API
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "name": "test-instruction-get-assigned-instruction-in-weather-api-4-1-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-alarm-weather-widget-instructions-resource-4"
            }
          }
        ]
      ],
      "description": "test-instruction-get-assigned-instruction-in-weather-api-4-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-6"
        }
      ],
      "triggers": ["create"],
      "priority": 2000
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "name": "test-instruction-get-assigned-instruction-in-weather-api-4-2-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-alarm-weather-widget-instructions-resource-4"
            }
          }
        ]
      ],
      "description": "test-instruction-get-assigned-instruction-in-weather-api-4-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 3,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-7"
        }
      ],
      "triggers": ["create"],
      "priority": 2001
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-alarm-weather-widget-instructions-connector-4",
      "connector_name": "test-alarm-weather-widget-instructions-connectorname-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-alarm-weather-widget-instructions-component-4",
      "resource": "test-alarm-weather-widget-instructions-resource-4",
      "state": 1,
      "output": "test-alarm-weather-widget-instructions-output-4"
    }
    """
    When I do GET /api/v4/alarms?search=test-alarm-weather-widget-instructions-resource-4&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "instruction_execution_icon": 2
        }
      ]
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-entity-instruction-weather-services-4",
      "name": "test-entity-instruction-weather-services-4",
      "output_template": "test-entity-instruction-weather-services-4",
      "category": "",
      "impact_level": 1,
      "enabled": true,
      "sli_avail_state": 1,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-alarm-weather-widget-instructions-resource-4"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response serviceId={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "component": "{{ .serviceId }}"
      },
      {
        "event_type": "activate",
        "component": "{{ .serviceId }}"
      }
    ]
    """
    When I do GET /api/v4/weather-services/test-entity-instruction-weather-services-4?with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-weather-widget-instructions-resource-4/test-alarm-weather-widget-instructions-component-4",
          "instruction_execution_icon": 2
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
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "trigger",
        "connector": "test-alarm-weather-widget-instructions-connector-4",
        "connector_name": "test-alarm-weather-widget-instructions-connectorname-4",
        "component": "test-alarm-weather-widget-instructions-component-4",
        "resource": "test-alarm-weather-widget-instructions-resource-4",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-alarm-weather-widget-instructions-connector-4",
        "connector_name": "test-alarm-weather-widget-instructions-connectorname-4",
        "component": "test-alarm-weather-widget-instructions-component-4",
        "resource": "test-alarm-weather-widget-instructions-resource-4",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-alarm-weather-widget-instructions-connector-4",
        "connector_name": "test-alarm-weather-widget-instructions-connectorname-4",
        "component": "test-alarm-weather-widget-instructions-component-4",
        "resource": "test-alarm-weather-widget-instructions-resource-4",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-alarm-weather-widget-instructions-connector-4",
        "connector_name": "test-alarm-weather-widget-instructions-connectorname-4",
        "component": "test-alarm-weather-widget-instructions-component-4",
        "resource": "test-alarm-weather-widget-instructions-resource-4",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/weather-services/test-entity-instruction-weather-services-4?with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-weather-widget-instructions-resource-4/test-alarm-weather-widget-instructions-component-4",
          "instruction_execution_icon": 6
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
    When I do GET /api/v4/weather-services/test-entity-instruction-weather-services-4?with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-weather-widget-instructions-resource-4/test-alarm-weather-widget-instructions-component-4",
          "instruction_execution_icon": 3
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
