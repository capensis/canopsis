Feature: cancel a instruction execution
  I need to be able to cancel a instruction execution
  Only admin should be able to cancel a instruction execution

  Scenario: given instruction should cancel running executions
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-execution-cancel-1-name",
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-instruction-execution-cancel-resource-1"
            }
          }
        ]
      ],
      "description": "test-instruction-execution-cancel-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-execution-cancel-1-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-cancel-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-cancel-1-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-cancel-1-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-cancel-1-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-cancel-1-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-instruction-execution-cancel-1",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/cancel
    Then the response code should be 204

  Scenario: given instruction should not cancel completed executions
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-execution-cancel-2-name",
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-instruction-execution-cancel-resource-2"
            }
          }
        ]
      ],
      "description": "test-instruction-execution-cancel-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-execution-cancel-2-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-cancel-2-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-cancel-2-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-cancel-2-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-cancel-2-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-cancel-2-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-instruction-execution-cancel-2",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/cancel
    Then the response code should be 404

  Scenario: given unauth request should not allow access
    When I do PUT /api/v4/cat/executions/notexist/cancel
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/executions/notexist/cancel
    Then the response code should be 403
