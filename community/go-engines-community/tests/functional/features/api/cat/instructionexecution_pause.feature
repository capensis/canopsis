Feature: pause a instruction execution
  I need to be able to pause instruction operation
  Only admin should be able to pause a instruction

  Scenario: given running instruction should pause execution
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """
    {
      "name": "test-instruction-execution-pause-1-name",
      "alarm_patterns": [
        {
          "_id": "test-instruction-execution-pause-1"
        }
      ],
      "description": "test-instruction-execution-pause-1-description",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-execution-pause-1-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-pause-1-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-execution-pause-1-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-pause-1-step-1-operation-2",
              "time_to_complete": {"seconds": 3, "unit":"s"},
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
              "time_to_complete": {"seconds": 6, "unit":"s"},
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
    """
    {
      "alarm": "test-instruction-execution-pause-1",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/pause
    Then the response code should be 204

  Scenario: given running instruction should pause execution after long time of inactivity
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """
    {
      "name": "test-instruction-execution-pause-2-name",
      "alarm_patterns": [
        {
          "_id": "test-instruction-execution-pause-2"
        }
      ],
      "description": "test-instruction-execution-pause-2-description",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-execution-pause-2-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-pause-2-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-execution-pause-2-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-pause-2-step-1-operation-2",
              "time_to_complete": {"seconds": 3, "unit":"s"},
              "description": "test-instruction-execution-pause-2-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-pause-2-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-pause-2-step-2",
          "operations": [
            {
              "name": "test-instruction-execution-pause-2-step-2-operation-1",
              "time_to_complete": {"seconds": 6, "unit":"s"},
              "description": "test-instruction-execution-pause-2-step-2-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-pause-2-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """
    {
      "alarm": "test-instruction-execution-pause-2",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/executions/{{ .lastResponse._id }}
    Then the response code should be 200
    Then I wait 5s
    When I do GET /api/v4/cat/executions/{{ .lastResponse._id }}
    Then the response code should be 404

  Scenario: given running instruction should not pause execution user is active
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """
    {
      "name": "test-instruction-execution-pause-3-name",
      "alarm_patterns": [
        {
          "_id": "test-instruction-execution-pause-3"
        }
      ],
      "description": "test-instruction-execution-pause-3-description",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-execution-pause-3-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-pause-3-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-execution-pause-3-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-pause-3-step-1-operation-2",
              "time_to_complete": {"seconds": 3, "unit":"s"},
              "description": "test-instruction-execution-pause-3-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-pause-3-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-pause-3-step-2",
          "operations": [
            {
              "name": "test-instruction-execution-pause-3-step-2-operation-1",
              "time_to_complete": {"seconds": 6, "unit":"s"},
              "description": "test-instruction-execution-pause-3-step-2-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-pause-3-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """
    {
      "alarm": "test-instruction-execution-pause-3",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I do GET /api/v4/cat/executions/{{ .executionID }}
    Then the response code should be 200
    Then I wait 3s
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/ping
    Then the response code should be 200
    Then I wait 3s
    When I do GET /api/v4/cat/executions/{{ .executionID }}
    Then the response code should be 200

  Scenario: given running instruction should add pause execution to user
    When I do POST /auth:
    """
    {
      "username": "test-user-to-test-paused-executions",
      "password": "test"
    }
    """
    Then the response code should be 200
    When I do POST /api/v4/cat/instructions:
    """
    {
      "name": "test-instruction-execution-pause-4-name",
      "alarm_patterns": [
        {
          "_id": "test-instruction-execution-pause-4"
        }
      ],
      "description": "test-instruction-execution-pause-4-description",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-execution-pause-4-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-pause-4-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-execution-pause-4-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-pause-4-step-1-operation-2",
              "time_to_complete": {"seconds": 3, "unit":"s"},
              "description": "test-instruction-execution-pause-4-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-pause-4-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-pause-4-step-2",
          "operations": [
            {
              "name": "test-instruction-execution-pause-4-step-2-operation-1",
              "time_to_complete": {"seconds": 6, "unit":"s"},
              "description": "test-instruction-execution-pause-4-step-2-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-pause-4-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """
    {
      "alarm": "test-instruction-execution-pause-4",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/executions/{{ .lastResponse._id }}
    Then the response code should be 200
    Then I wait 5s
    When I do GET /api/v4/cat/account/paused-executions
    Then the response code should be 200
    Then the response body should contain:
    """
    [
      {
        "alarm_name": "RC-KC_tW",
        "instruction_name": "test-instruction-execution-pause-4-name"
      }
    ]
    """
    When I do GET /api/v4/cat/account/paused-executions
    Then the response code should be 200
    Then the response body should contain:
    """
    []
    """

  Scenario: given unauth request should not allow access
    When I do PUT /api/v4/cat/executions/test-instruction-execution-running/pause
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/executions/test-instruction-execution-running/pause
    Then the response code should be 403