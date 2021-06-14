Feature: instruction execution should be added to alarm steps
  I need to be able to see instruction execution steps in alarm timeline.

  Scenario: given instruction should add instruction start step to alarm steps
    Given I am admin
    When I send an event:
    """
    {
      "connector" : "test-connector-axe-api-instruction-1",
      "connector_name" : "test-connector-name-axe-api-instruction-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-axe-api-instruction-1",
      "resource" : "test-resource-axe-api-instruction-1",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-api-instruction-1
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/instructions:
    """
    {
      "name": "test-instruction-axe-api-instruction-1-name",
      "entity_patterns": [
        {
          "_id": "test-resource-axe-api-instruction-1/test-component-axe-api-instruction-1"
         }
      ],
      "description": "test-instruction-axe-api-instruction-1-description",
      "author": "test-instruction-axe-api-instruction-1-author",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-axe-api-instruction-1-step-1",
          "operations": [
            {
              "name": "test-instruction-axe-api-instruction-1-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-axe-api-instruction-1-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-axe-api-instruction-1-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-api-instruction-1&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "steps": [
              {},
              {},
              {
                "_t": "instructionstart",
                "a": "root",
                "m": "Instruction test-instruction-axe-api-instruction-1-name"
              }
            ]
          }
        }
      ]
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/cancel
    Then the response code should be 204
    When I wait the end of event processing

  Scenario: given instruction should add instruction complete step to alarm steps
    Given I am admin
    When I send an event:
    """
    {
      "connector" : "test-connector-axe-api-instruction-2",
      "connector_name" : "test-connector-name-axe-api-instruction-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-axe-api-instruction-2",
      "resource" : "test-resource-axe-api-instruction-2",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-api-instruction-2
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/instructions:
    """
    {
      "name": "test-instruction-axe-api-instruction-2-name",
      "entity_patterns": [
        {
          "_id": "test-resource-axe-api-instruction-2/test-component-axe-api-instruction-2"
        }
      ],
      "description": "test-instruction-axe-api-instruction-2-description",
      "author": "test-instruction-axe-api-instruction-2-author",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-axe-api-instruction-2-step-1",
          "operations": [
            {
              "name": "test-instruction-axe-api-instruction-2-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-axe-api-instruction-2-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-axe-api-instruction-2-step-1-endpoint"
        },
        {
          "name": "test-instruction-axe-api-instruction-2-step-2",
          "operations": [
            {
              "name": "test-instruction-axe-api-instruction-2-step-2-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-axe-api-instruction-2-step-2-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-axe-api-instruction-2-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-api-instruction-2&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "steps": [
              {},
              {},
              {
                "_t": "instructionstart",
                "a": "root",
                "m": "Instruction test-instruction-axe-api-instruction-2-name"
              },
              {
                "_t": "instructioncomplete",
                "a": "root",
                "m": "Instruction test-instruction-axe-api-instruction-2-name"
              }
            ]
          }
        }
      ]
    }
    """

  Scenario: given paused instruction by request should add instruction pause step to alarm steps
    Given I am admin
    When I send an event:
    """
    {
      "connector" : "test-connector-axe-api-instruction-3",
      "connector_name" : "test-connector-name-axe-api-instruction-3",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-axe-api-instruction-3",
      "resource" : "test-resource-axe-api-instruction-3",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-api-instruction-3
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/instructions:
    """
    {
      "name": "test-instruction-axe-api-instruction-3-name",
      "entity_patterns": [
        {
          "_id": "test-resource-axe-api-instruction-3/test-component-axe-api-instruction-3"
        }
      ],
      "description": "test-instruction-axe-api-instruction-3-description",
      "author": "test-instruction-axe-api-instruction-3-author",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-axe-api-instruction-3-step-1",
          "operations": [
            {
              "name": "test-instruction-axe-api-instruction-3-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-axe-api-instruction-3-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-axe-api-instruction-3-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/pause
    Then the response code should be 204
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-api-instruction-3&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "steps": [
              {},
              {},
              {
                "_t": "instructionstart",
                "a": "root",
                "m": "Instruction test-instruction-axe-api-instruction-3-name"
              },
              {
                "_t": "instructionpause",
                "a": "root",
                "m": "Instruction test-instruction-axe-api-instruction-3-name"
              }
            ]
          }
        }
      ]
    }
    """

  Scenario: given instruction should add instruction resume step to alarm steps
    Given I am admin
    When I send an event:
    """
    {
      "connector" : "test-connector-axe-api-instruction-4",
      "connector_name" : "test-connector-name-axe-api-instruction-4",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-axe-api-instruction-4",
      "resource" : "test-resource-axe-api-instruction-4",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-api-instruction-4
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/instructions:
    """
    {
      "name": "test-instruction-axe-api-instruction-4-name",
      "entity_patterns": [
        {
          "_id": "test-resource-axe-api-instruction-4/test-component-axe-api-instruction-4"
        }
      ],
      "description": "test-instruction-axe-api-instruction-4-description",
      "author": "test-instruction-axe-api-instruction-4-author",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-axe-api-instruction-4-step-1",
          "operations": [
            {
              "name": "test-instruction-axe-api-instruction-4-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-axe-api-instruction-4-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-axe-api-instruction-4-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I wait the end of event processing
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/pause
    Then the response code should be 204
    When I wait the end of event processing
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/resume
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-api-instruction-4&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "steps": [
              {},
              {},
              {
                "_t": "instructionstart",
                "a": "root",
                "m": "Instruction test-instruction-axe-api-instruction-4-name"
              },
              {
                "_t": "instructionpause",
                "a": "root",
                "m": "Instruction test-instruction-axe-api-instruction-4-name"
              },
              {
                "_t": "instructionresume",
                "a": "root",
                "m": "Instruction test-instruction-axe-api-instruction-4-name"
              }
            ]
          }
        }
      ]
    }
    """
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/cancel
    Then the response code should be 204
    When I wait the end of event processing

  Scenario: given instruction should add instruction abort step to alarm steps
    Given I am admin
    When I send an event:
    """
    {
      "connector" : "test-connector-axe-api-instruction-5",
      "connector_name" : "test-connector-name-axe-api-instruction-5",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-axe-api-instruction-5",
      "resource" : "test-resource-axe-api-instruction-5",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-api-instruction-5
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/instructions:
    """
    {
      "name": "test-instruction-axe-api-instruction-5-name",
      "entity_patterns": [
        {
          "_id": "test-resource-axe-api-instruction-5/test-component-axe-api-instruction-5"
        }
      ],
      "description": "test-instruction-axe-api-instruction-5-description",
      "author": "test-instruction-axe-api-instruction-5-author",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-axe-api-instruction-5-step-1",
          "operations": [
            {
              "name": "test-instruction-axe-api-instruction-5-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-axe-api-instruction-5-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-axe-api-instruction-5-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/cancel
    Then the response code should be 204
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-api-instruction-5&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "steps": [
              {},
              {},
              {
                "_t": "instructionstart",
                "a": "root",
                "m": "Instruction test-instruction-axe-api-instruction-5-name"
              },
              {
                "_t": "instructionabort",
                "a": "root",
                "m": "Instruction test-instruction-axe-api-instruction-5-name"
              }
            ]
          }
        }
      ]
    }
    """

  Scenario: given instruction should add instruction fail step to alarm steps
    Given I am admin
    When I send an event:
    """
    {
      "connector" : "test-connector-axe-api-instruction-6",
      "connector_name" : "test-connector-name-axe-api-instruction-6",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-axe-api-instruction-6",
      "resource" : "test-resource-axe-api-instruction-6",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-api-instruction-6
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/instructions:
    """
    {
      "name": "test-instruction-axe-api-instruction-6-name",
      "entity_patterns": [
        {
          "_id": "test-resource-axe-api-instruction-6/test-component-axe-api-instruction-6"
        }
      ],
      "description": "test-instruction-axe-api-instruction-6-description",
      "author": "test-instruction-axe-api-instruction-6-author",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-axe-api-instruction-6-step-1",
          "operations": [
            {
              "name": "test-instruction-axe-api-instruction-6-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-axe-api-instruction-6-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-axe-api-instruction-6-step-1-endpoint"
        },
        {
          "name": "test-instruction-axe-api-instruction-6-step-2",
          "operations": [
            {
              "name": "test-instruction-axe-api-instruction-6-step-2-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-axe-api-instruction-6-step-2-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-axe-api-instruction-6-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step:
    """
    {
      "failed": true
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-api-instruction-6&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "steps": [
              {},
              {},
              {
                "_t": "instructionstart",
                "a": "root",
                "m": "Instruction test-instruction-axe-api-instruction-6-name"
              },
              {
                "_t": "instructionfail",
                "a": "root",
                "m": "Instruction test-instruction-axe-api-instruction-6-name"
              }
            ]
          }
        }
      ]
    }
    """

  Scenario: given paused instruction by ping timeout should add instruction pause step to alarm steps
    Given I am admin
    When I send an event:
    """
    {
      "connector" : "test-connector-axe-api-instruction-7",
      "connector_name" : "test-connector-name-axe-api-instruction-7",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-axe-api-instruction-7",
      "resource" : "test-resource-axe-api-instruction-7",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-api-instruction-7
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/instructions:
    """
    {
      "name": "test-instruction-axe-api-instruction-7-name",
      "filter":{
        "$and":[
           {
              "name": "test filter"
           }
        ]
      },
      "description": "test-instruction-axe-api-instruction-7-description",
      "author": "test-instruction-axe-api-instruction-7-author",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-axe-api-instruction-7-step-1",
          "operations": [
            {
              "name": "test-instruction-axe-api-instruction-7-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-axe-api-instruction-7-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-axe-api-instruction-7-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    Then I wait 5s
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-api-instruction-7&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "steps": [
              {},
              {},
              {
                "_t": "instructionstart",
                "a": "root",
                "m": "Instruction test-instruction-axe-api-instruction-7-name"
              },
              {
                "_t": "instructionpause",
                "a": "system",
                "m": "Instruction test-instruction-axe-api-instruction-7-name"
              }
            ]
          }
        }
      ]
    }
    """

  Scenario: given aborted execution by instruction update should add instruction abort step to alarm steps
    Given I am admin
    When I send an event:
    """
    {
      "connector" : "test-connector-axe-api-instruction-8",
      "connector_name" : "test-connector-name-axe-api-instruction-8",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-axe-api-instruction-8",
      "resource" : "test-resource-axe-api-instruction-8",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-api-instruction-8
    When the response code should be 200
    When the response body should contain:
    """
    {
      "meta": {
        "total_count": 1
      }
    }
    """
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/instructions:
    """
    {
      "name": "test-instruction-axe-api-instruction-8-name",
      "filter":{
        "$and":[
           {
              "name": "test filter"
           }
        ]
      },
      "description": "test-instruction-axe-api-instruction-8-description",
      "author": "test-instruction-axe-api-instruction-8-author",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-axe-api-instruction-8-step-1",
          "operations": [
            {
              "name": "test-instruction-axe-api-instruction-8-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-axe-api-instruction-8-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-axe-api-instruction-8-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/executions:
    """
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}:
    """
    {
      "name": "test-instruction-axe-api-instruction-8-name",
      "filter":{
        "$and":[
           {
              "name": "test filter"
           }
        ]
      },
      "description": "test-instruction-axe-api-instruction-8-description",
      "author": "test-instruction-axe-api-instruction-8-author",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-axe-api-instruction-8-step-1",
          "operations": [
            {
              "name": "test-instruction-axe-api-instruction-8-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-axe-api-instruction-8-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-axe-api-instruction-8-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-api-instruction-8&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "steps": [
              {},
              {},
              {
                "_t": "instructionstart",
                "a": "root",
                "m": "Instruction test-instruction-axe-api-instruction-8-name"
              },
              {
                "_t": "instructionabort",
                "a": "system",
                "m": "Instruction test-instruction-axe-api-instruction-8-name"
              }
            ]
          }
        }
      ]
    }
    """
