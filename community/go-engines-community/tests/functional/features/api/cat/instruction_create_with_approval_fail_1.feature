Feature: instruction approval creation
  I need to be able to create an instruction with approval

  @concurrent
  Scenario: given create request with approval request with not existed username should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "approval": {
        "user": "notexist"
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "approval.user": "User doesn't have approve rights or doesn't exist."
      }
    }
    """

  @concurrent
  Scenario: given create request with approval request with a username without approve right should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "approval": {
        "user": "nopermsuser"
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "approval.user": "User doesn't have approve rights or doesn't exist."
      }
    }
    """

  @concurrent
  Scenario: given create request with approval request with not existed role should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "approval": {
        "role": "notexist"
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "approval.role": "Role doesn't have approve rights or doesn't exist."
      }
    }
    """

  @concurrent
  Scenario: given create request with approval request with a role without approve right should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "approval": {
        "role": "noperms"
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "approval.role": "Role doesn't have approve rights or doesn't exist."
      }
    }
    """

  @concurrent
  Scenario: given create request with approval request with role and username should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "approval": {
        "user": "root",
        "role": "admin"
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "approval.role": "Can't be present both Role and User."
      }
    }
    """

  @concurrent
  Scenario: given create request with valid approval request with existed name should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-check-unique-name",
      "description": "test-instruction-to-create-with-approval-fail-6-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-fail-6-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-fail-6-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-fail-6-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-fail-6-step-1-endpoint"
        }
      ],
      "approval": {
        "role": "role-to-instruction-approve-1",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "name": "Name already exists."
      }
    }
    """

  @concurrent
  Scenario: the user, which is not set in approval should receive 403
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-fail-7-name",
      "description": "test-instruction-to-create-with-approval-fail-7-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-fail-7-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-fail-7-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-fail-7-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-fail-7-step-1-endpoint"
        }
      ],
      "approval": {
        "user": "user-to-instruction-approve-1",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I am manager
    When I do GET /api/v4/cat/instructions/{{ .instructionID }}/approval
    Then the response code should be 403
    Then the response body should contain:
    """json
    {
      "error": "user is not assigned to approval"
    }
    """

  @concurrent
  Scenario: the user with a role, which is not set in approval should receive 403
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-fail-8-name",
      "description": "test-instruction-to-create-with-approval-fail-8-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-fail-8-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-fail-8-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-fail-8-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-fail-8-step-1-endpoint"
        }
      ],
      "approval": {
        "role": "role-to-instruction-approve-1",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I am manager
    When I do GET /api/v4/cat/instructions/{{ .instructionID }}/approval
    Then the response code should be 403
    Then the response body should contain:
    """json
    {
      "error": "role is not assigned to approval"
    }
    """

  @concurrent
  Scenario: if no approval return 404
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-fail-9-name",
      "description": "test-instruction-to-create-with-approval-fail-9-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-fail-9-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-fail-9-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-fail-9-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-fail-9-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/cat/instructions/{{ .lastResponse._id }}/approval
    Then the response code should be 404
