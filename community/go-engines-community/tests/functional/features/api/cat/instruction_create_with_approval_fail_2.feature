Feature: instruction approval creation
  I need to be able to create an instruction with approval

  @standalone
  Scenario: given create request should only create with approval
    When I am admin
    When I do PUT /api/v4/internal/user_interface:
    """json
    {
      "required_instruction_approve": true
    }
    """
    When I wait the next periodical process
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-fail-second-1-name",
      "description": "test-instruction-to-create-with-approval-fail-second-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-fail-second-1-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-fail-second-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-fail-second-1-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-fail-second-1-step-1-endpoint"
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
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-fail-second-1-name",
      "description": "test-instruction-to-create-with-approval-fail-second-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-fail-second-1-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-fail-second-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-fail-second-1-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-fail-second-1-step-1-endpoint"
        }
      ],
      "approval": {
        "role": "role-to-instruction-approve-1",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 201

  @concurrent
  Scenario: given dismiss request with missing comment should return error
    When I am admin
    When I do PUT /api/v4/cat/instructions/test-instruction-not-exist/approval:
    """json
    {
      "approve": false
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "comment": "Comment is required when Approve false is defined."
      }
    }
    """
