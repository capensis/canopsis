Feature: pause a instruction execution
  I need to be able to pause instruction operation
  Only admin should be able to pause a instruction

  Scenario: given running instruction should pause execution
    When I am admin
    When I do PUT /api/v4/cat/executions/test-instruction-execution-resume-1/resume
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "description": "test-instruction-execution-resume-1-description",
      "name": "test-instruction-execution-resume-1-name",
      "status": 0,
      "steps": [
        {
          "endpoint": "test-instruction-execution-resume-1-step-1-endpoint",
          "name": "test-instruction-execution-resume-1-step-1-name",
          "time_to_complete": {"value": 4, "unit":"s"},
          "operations": [
            {
              "started_at": 1597905434,
              "completed_at": 1597905437,
              "name": "test-instruction-execution-resume-1-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-resume-1-step-1-operation-1-description",
              "jobs": []
            },
            {
              "completed_at": 0,
              "name": "test-instruction-execution-resume-1-step-1-operation-2-name",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-resume-1-step-1-operation-2-description",
              "jobs": []
            }
          ]
        },
        {
          "name": "test-instruction-execution-resume-1-step-2-name",
          "time_to_complete": {"value": 6, "unit":"s"},
          "endpoint": "test-instruction-execution-resume-1-step-2-endpoint",
          "operations": [
            {
              "started_at": 0,
              "completed_at": 0,
              "name": "test-instruction-execution-resume-1-step-2-operation-1-name",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "",
              "jobs": []
            }
          ]
        }
      ]
    }
    """

  Scenario: given unauth request should not allow access
    When I do PUT /api/v4/cat/executions/notexist/resume
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/executions/notexist/resume
    Then the response code should be 403
