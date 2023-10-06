Feature: instruction approval update
  I need to be able to update an instruction with approval

  @concurrent
  Scenario: PUT a valid instruction with approval request with a not found user should return error
    When I am manager
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-second-1:
    """json
    {
      "approval": {
        "user": "approvernotexist",
        "comment": "test comment"
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
  Scenario: PUT a valid instruction with approval request with a username without approve right should return error
    When I am manager
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-second-1:
    """json
    {
      "approval": {
        "user": "nopermsuser",
        "comment": "test comment"
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
  Scenario: PUT a valid instruction with approval request with a not found role should return error
    When I am manager
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-second-1:
    """json
    {
      "approval": {
        "role": "rolenotexist",
        "comment": "test comment"
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
  Scenario: PUT a valid instruction with approval request with a role without approve right should return error
    When I am manager
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-second-1:
    """json
    {
      "approval": {
        "role": "noperms",
        "comment": "test comment"
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
  Scenario: PUT a valid instruction with approval request with a role without approve right should return error
    When I am manager
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-second-1:
    """json
    {
      "approval": {
        "user": "root",
        "role": "admin",
        "comment": "test comment"
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

  @standalone
  Scenario: given update request should only create with approval
    When I am admin
    When I do PUT /api/v4/internal/user_interface:
    """json
    {
      "required_instruction_approve": true
    }
    """
    When I wait the next periodical process
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-fail-6:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-update-with-approval-fail-6-name",
      "description": "test-instruction-to-update-with-approval-fail-6-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-fail-6-step-1",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-fail-6-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-fail-6-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-with-approval-fail-6-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "approval": "Approval is missing."
      }
    }
    """
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-fail-6:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-update-with-approval-fail-6-name",
      "description": "test-instruction-to-update-with-approval-fail-6-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-fail-6-step-1",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-fail-6-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-fail-6-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-with-approval-fail-6-step-1-endpoint"
        }
      ],
      "approval": {
        "role": "role-to-instruction-approve-1",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200
