Feature: get service entities with assigned instructions
  I need to be able get service entities with assigned instructions

  Scenario: given manual instruction execution should return assigned instructions in weather services entities API
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-get-assigned-instruction-in-weather-api-1-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-alarm-weather-widget-instructions-resource-1",
                "test-alarm-weather-widget-instructions-resource-2"
              ]
            }
          }
        ]
      ],
      "description": "test-instruction-get-assigned-instruction-in-weather-api-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-get-assigned-instruction-in-weather-api-1-step-1",
          "operations": [
            {
              "name": "test-instruction-get-assigned-instruction-in-weather-api-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-get-assigned-instruction-in-weather-api-1-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-get-assigned-instruction-in-weather-api-1-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID1={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-get-assigned-instruction-in-weather-api-2-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-alarm-weather-widget-instructions-resource-2",
                "test-alarm-weather-widget-instructions-resource-3"
              ]
            }
          }
        ]
      ],
      "description": "test-instruction-get-assigned-instruction-in-weather-api-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-get-assigned-instruction-in-weather-api-2-step-1",
          "operations": [
            {
              "name": "test-instruction-get-assigned-instruction-in-weather-api-2-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-get-assigned-instruction-in-weather-api-2-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-get-assigned-instruction-in-weather-api-2-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID2={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-alarm-weather-widget-instructions-connector-1",
      "connector_name": "test-alarm-weather-widget-instructions-connectorname-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-alarm-weather-widget-instructions-component-1",
      "resource": "test-alarm-weather-widget-instructions-resource-1",
      "state": 1,
      "output": "test-alarm-weather-widget-instructions-output-1"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-alarm-weather-widget-instructions-connector-1",
      "connector_name": "test-alarm-weather-widget-instructions-connectorname-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-alarm-weather-widget-instructions-component-1",
      "resource": "test-alarm-weather-widget-instructions-resource-2",
      "state": 1,
      "output": "test-alarm-weather-widget-instructions-output-1"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-alarm-weather-widget-instructions-connector-1",
      "connector_name": "test-alarm-weather-widget-instructions-connectorname-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-alarm-weather-widget-instructions-component-1",
      "resource": "test-alarm-weather-widget-instructions-resource-3",
      "state": 1,
      "output": "test-alarm-weather-widget-instructions-output-1"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-alarm-weather-widget-instructions-connector-1",
      "connector_name": "test-alarm-weather-widget-instructions-connectorname-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-alarm-weather-widget-instructions-component-1",
      "resource": "test-alarm-weather-widget-instructions-resource-4",
      "state": 1,
      "output": "test-alarm-weather-widget-instructions-output-1"
    }
    """
    When I wait the end of event processing
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
                "test-alarm-weather-widget-instructions-resource-1",
                "test-alarm-weather-widget-instructions-resource-2",
                "test-alarm-weather-widget-instructions-resource-3",
                "test-alarm-weather-widget-instructions-resource-4"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services/test-entity-instruction-weather-services-1?with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-weather-widget-instructions-resource-1/test-alarm-weather-widget-instructions-component-1",
          "assigned_instructions": [
            {
              "_id": "{{.instructionID1}}"
            }
          ]
        },
        {
          "_id": "test-alarm-weather-widget-instructions-resource-2/test-alarm-weather-widget-instructions-component-1",
          "assigned_instructions": [
            {
              "_id": "{{.instructionID1}}"
            },
            {
              "_id": "{{.instructionID2}}"
            }
          ]
        },
        {
          "_id": "test-alarm-weather-widget-instructions-resource-3/test-alarm-weather-widget-instructions-component-1",
          "assigned_instructions": [
            {
              "_id": "{{.instructionID2}}"
            }
          ]
        },
        {
          "_id": "test-alarm-weather-widget-instructions-resource-4/test-alarm-weather-widget-instructions-component-1",
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

  Scenario: given auto instruction execution should return auto instruction flags in weather API
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "name": "test-instruction-get-assigned-instruction-in-weather-api-3-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-alarm-weather-widget-instructions-resource-5"
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
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-7"
        }
      ],
      "priority": 30,
      "triggers": ["create"]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "name": "test-instruction-get-assigned-instruction-in-weather-api-4-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-alarm-weather-widget-instructions-resource-5"
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
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-7"
        }
      ],
      "priority": 31,
      "triggers": ["create"]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-alarm-weather-widget-instructions-connector-1",
      "connector_name": "test-alarm-weather-widget-instructions-connectorname-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-alarm-weather-widget-instructions-component-1",
      "resource": "test-alarm-weather-widget-instructions-resource-5",
      "state": 1,
      "output": "test-alarm-weather-widget-instructions-output-1"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-alarm-weather-widget-instructions-resource-5&with_instructions=true
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
              "value": "test-alarm-weather-widget-instructions-resource-5"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services/test-entity-instruction-weather-services-2?with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-weather-widget-instructions-resource-5/test-alarm-weather-widget-instructions-component-1",
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
          "_id": "test-alarm-weather-widget-instructions-resource-5/test-alarm-weather-widget-instructions-component-1",
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
    When I wait the end of 4 events processing
    When I do GET /api/v4/weather-services/test-entity-instruction-weather-services-2?with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-weather-widget-instructions-resource-5/test-alarm-weather-widget-instructions-component-1",
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

  Scenario: given manual instruction execution should return manual instruction flags in weather API
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-get-assigned-instruction-in-weather-api-5-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-alarm-weather-widget-instructions-resource-6"
            }
          }
        ]
      ],
      "description": "test-instruction-get-assigned-instruction-in-weather-api-5-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-get-assigned-instruction-in-weather-api-5-step-1",
          "operations": [
            {
              "name": "test-instruction-get-assigned-instruction-in-weather-api-5-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-get-assigned-instruction-in-weather-api-5-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-get-assigned-instruction-in-weather-api-5-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-alarm-weather-widget-instructions-connector-1",
      "connector_name": "test-alarm-weather-widget-instructions-connectorname-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-alarm-weather-widget-instructions-component-1",
      "resource": "test-alarm-weather-widget-instructions-resource-6",
      "state": 1,
      "output": "test-alarm-weather-widget-instructions-output-1"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-alarm-weather-widget-instructions-resource-6
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
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
              "value": "test-alarm-weather-widget-instructions-resource-6"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of event processing
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
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    When I wait the end of event processing
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

  Scenario: given auto failed instruction execution should return auto instruction flags in weather API
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "name": "test-instruction-get-assigned-instruction-in-weather-api-6-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-alarm-weather-widget-instructions-resource-7"
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
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-6"
        }
      ],
      "priority": 32,
      "triggers": ["create"]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "name": "test-instruction-get-assigned-instruction-in-weather-api-7-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-alarm-weather-widget-instructions-resource-7"
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
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-7"
        }
      ],
      "priority": 33,
      "triggers": ["create"]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-alarm-weather-widget-instructions-connector-1",
      "connector_name": "test-alarm-weather-widget-instructions-connectorname-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-alarm-weather-widget-instructions-component-1",
      "resource": "test-alarm-weather-widget-instructions-resource-7",
      "state": 1,
      "output": "test-alarm-weather-widget-instructions-output-1"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-alarm-weather-widget-instructions-resource-7&with_instructions=true
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
              "value": "test-alarm-weather-widget-instructions-resource-7"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services/test-entity-instruction-weather-services-4?with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-weather-widget-instructions-resource-7/test-alarm-weather-widget-instructions-component-1",
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
    When I wait the end of 4 events processing
    When I do GET /api/v4/weather-services/test-entity-instruction-weather-services-4?with_instructions=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-weather-widget-instructions-resource-7/test-alarm-weather-widget-instructions-component-1",
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
          "_id": "test-alarm-weather-widget-instructions-resource-7/test-alarm-weather-widget-instructions-component-1",
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
