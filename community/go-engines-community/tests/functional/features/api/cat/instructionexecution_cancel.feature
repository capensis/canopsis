Feature: cancel a instruction execution
  I need to be able to cancel a instruction execution
  Only admin should be able to cancel a instruction execution

  Scenario: given instruction should cancel running executions
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """
    {
      "name": "test-instruction-execution-cancel-1-name",
      "alarm_patterns": [
        {
          "_id": "test-instruction-execution-cancel-1"
        }
      ],
      "description": "test-instruction-execution-cancel-1-description",
      "author": "test-instruction-execution-cancel-1-author",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-execution-cancel-1-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-cancel-1-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-execution-cancel-1-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-cancel-1-step-1-operation-2",
              "time_to_complete": {"seconds": 3, "unit":"s"},
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
    """
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
    """
    {
      "name": "test-instruction-execution-cancel-2-name",
      "alarm_patterns": [
        {
          "_id": "test-instruction-execution-cancel-2"
        }
      ],
      "description": "test-instruction-execution-cancel-2-description",
      "author": "test-instruction-execution-cancel-2-author",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-execution-cancel-2-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-cancel-2-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-execution-cancel-2-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-cancel-2-step-1-operation-2",
              "time_to_complete": {"seconds": 3, "unit":"s"},
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
    """
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

  Scenario: given instruction should cancel running executions on instruction update
    When I am admin
    When I do GET /api/v4/cat/executions/test-instruction-cancel-execution-running-1
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "status": 0
    }
    """
    When I do PUT /api/v4/cat/instructions/test-instruction-cancel-execution-running-1:
    """
    {
      "name": "test-instruction-cancel-execution-running-1-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-cancel-execution-running-1-description",
      "author": "test-instruction-cancel-execution-running-1-author",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-cancel-execution-running-1-step-1",
          "operations": [
            {
              "name": "test-instruction-cancel-execution-running-1-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-cancel-execution-running-1-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-cancel-execution-running-1-step-1-operation-2",
              "time_to_complete": {"seconds": 3, "unit":"s"},
              "description": "test-instruction-cancel-execution-running-1-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-cancel-execution-running-1-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/executions/test-instruction-cancel-execution-running-1
    Then the response code should be 410

  Scenario: given unauth request should not allow access
    When I do PUT /api/v4/cat/executions/test-instruction-execution-running/cancel
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/executions/test-instruction-execution-running/cancel
    Then the response code should be 403