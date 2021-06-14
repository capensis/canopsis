Feature: move a instruction execution to next operation
  I need to be able to run next instruction operation
  Only admin should be able to run a instruction

  Scenario: given running instruction should start next operation of instruction
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """
    {
      "name": "test-instruction-execution-next-1-name",
      "alarm_patterns": [
        {
          "_id": "test-instruction-execution-next-1"
        }
      ],
      "description": "test-instruction-execution-next-1-description",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-execution-next-1-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-next-1-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-execution-next-1-step-1-operation-1-description connector {{ `{{ .Alarm.Value.Connector }}` }} entity {{ `{{ .Entity.ID }}` }}"
            },
            {
              "name": "test-instruction-execution-next-1-step-1-operation-2",
              "time_to_complete": {"seconds": 3, "unit":"s"},
              "description": "test-instruction-execution-next-1-step-1-operation-2-description connector {{ `{{ .Alarm.Value.Connector }}` }} entity {{ `{{ .Entity.ID }}` }}"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-next-1-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-next-1-step-2",
          "operations": [
            {
              "name": "test-instruction-execution-next-1-step-2-operation-1",
              "time_to_complete": {"seconds": 6, "unit":"s"},
              "description": "test-instruction-execution-next-1-step-2-operation-1-description connector {{ `{{ .Alarm.Value.Connector }}` }} entity {{ `{{ .Entity.ID }}` }}"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-next-1-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """
    {
      "alarm": "test-instruction-execution-next-1",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "status": 0,
      "name": "test-instruction-execution-next-1-name",
      "description": "test-instruction-execution-next-1-description",
      "steps": [
        {
          "name": "test-instruction-execution-next-1-step-1",
          "time_to_complete": {"seconds": 4, "unit":"s"},
          "completed_at": 0,
          "failed_at": 0,
          "operations": [
            {
              "name": "test-instruction-execution-next-1-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-execution-next-1-step-1-operation-1-description connector test-instruction-execution-next-connector entity test-instruction-execution-next-resource-1/test-instruction-execution-next-component"
            },
            {
              "completed_at": 0,
              "name": "test-instruction-execution-next-1-step-1-operation-2",
              "time_to_complete": {"seconds": 3, "unit":"s"},
              "description": "test-instruction-execution-next-1-step-1-operation-2-description connector test-instruction-execution-next-connector entity test-instruction-execution-next-resource-1/test-instruction-execution-next-component",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-next-1-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-next-1-step-2",
          "time_to_complete": {"seconds": 6, "unit":"s"},
          "completed_at": 0,
          "failed_at": 0,
          "operations": [
            {
              "started_at": 0,
              "completed_at": 0,
              "name": "test-instruction-execution-next-1-step-2-operation-1",
              "time_to_complete": {"seconds": 6, "unit":"s"},
              "description": "",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-next-1-step-2-endpoint"
        }
      ]
    }
    """
    Then the response key "steps.0.operations.0.started_at" should not be "0"
    Then the response key "steps.0.operations.0.completed_at" should not be "0"
    Then the response key "steps.0.operations.1.started_at" should not be "0"

  Scenario: given running instruction with last operation started should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """
    {
      "name": "test-instruction-execution-next-2-name",
      "alarm_patterns": [
        {
          "_id": "test-instruction-execution-next-2"
        }
      ],
      "description": "test-instruction-execution-next-2-description",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-execution-next-2-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-next-2-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-execution-next-2-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-next-2-step-1-operation-2",
              "time_to_complete": {"seconds": 3, "unit":"s"},
              "description": "test-instruction-execution-next-2-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-next-2-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-next-2-step-2",
          "operations": [
            {
              "name": "test-instruction-execution-next-2-step-2-operation-1",
              "time_to_complete": {"seconds": 6, "unit":"s"},
              "description": "test-instruction-execution-next-2-step-2-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-next-2-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """
    {
      "alarm": "test-instruction-execution-next-2",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "error": "execution cannot be moved to next operation from last step operation"
    }
    """

  Scenario: given unauth request should not allow access
    When I do PUT /api/v4/cat/executions/test-instruction-execution-running/next
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/executions/test-instruction-execution-running/next
    Then the response code should be 403