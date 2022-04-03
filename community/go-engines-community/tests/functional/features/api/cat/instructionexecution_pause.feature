Feature: pause a instruction execution
  I need to be able to pause instruction operation
  Only admin should be able to pause a instruction

  Scenario: given running instruction should pause execution
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-execution-pause-1-name",
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-instruction-execution-pause-resource-1"
            }
          }
        ]
      ],
      "description": "test-instruction-execution-pause-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-execution-pause-1-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-pause-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-pause-1-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-pause-1-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-pause-1-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-pause-1-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-pause-1-step-2",
          "operations": [
            {
              "name": "test-instruction-execution-pause-1-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-pause-1-step-2-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-pause-1-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-instruction-execution-pause-1",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/pause
    Then the response code should be 204

  Scenario: given unauth request should not allow access
    When I do PUT /api/v4/cat/executions/notexist/pause
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/executions/notexist/pause
    Then the response code should be 403
