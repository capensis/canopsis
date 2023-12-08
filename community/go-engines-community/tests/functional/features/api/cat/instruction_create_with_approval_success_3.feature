Feature: instruction approval creation
  I need to be able to create an instruction with approval

  @concurrent
  Scenario: only an author should be able to update a dismissed instruction
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "approval": {
        "user": "user-to-instruction-approve-1",
        "comment": "test-instruction-to-create-with-approval-third-1-comment"
      },
      "type": 0,
      "name": "test-instruction-to-create-with-approval-third-1-name",
      "description": "test-instruction-to-create-with-approval-third-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-third-1-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-third-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-third-1-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-third-1-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I am role-to-instruction-approve-1
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}/approval:
    """json
    {
      "approve": false,
      "comment": "test-instruction-to-create-with-approval-third-1-dismiss-comment"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/instructions/{{ .instructionID }}/approval
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "original instruction not found"
    }
    """
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}:
    """json
    {
      "approval": {
        "user": "user-to-instruction-approve-2",
        "comment": "test-instruction-to-create-with-approval-third-1-comment-updated"
      },
      "type": 0,
      "name": "test-instruction-to-create-with-approval-third-1-name-updated",
      "description": "test-instruction-to-create-with-approval-third-1-description-updated",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-third-1-step-1-updated",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-third-1-step-1-operation-1-updated",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-third-1-step-1-operation-1-description-updated"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-third-1-step-1-endpoint-updated"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "type": 0,
      "status": 3,
      "name": "test-instruction-to-create-with-approval-third-1-name-updated",
      "description": "test-instruction-to-create-with-approval-third-1-description-updated",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-third-1-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-third-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-third-1-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-third-1-step-1-endpoint"
        }
      ],
      "approval": {
        "user": {
          "_id": "user-to-instruction-approve-1",
          "name": "user-to-instruction-approve-1"
        },
        "requested_by": {
          "_id": "root",
          "name": "root"
        },
        "comment": "test-instruction-to-create-with-approval-third-1-comment",
        "dismissed_by": {
          "_id": "user-to-instruction-approve-1",
          "name": "user-to-instruction-approve-1"
        },
        "dismiss_comment": "test-instruction-to-create-with-approval-third-1-dismiss-comment"
      }
    }
    """
    When I am admin
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}:
    """json
    {
      "approval": {
        "user": "user-to-instruction-approve-2",
        "comment": "test-instruction-to-create-with-approval-third-1-comment-updated"
      },
      "type": 0,
      "name": "test-instruction-to-create-with-approval-third-1-name-updated",
      "description": "test-instruction-to-create-with-approval-third-1-description-updated",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-third-1-step-1-updated",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-third-1-step-1-operation-1-updated",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-third-1-step-1-operation-1-description-updated"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-third-1-step-1-endpoint-updated"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "type": 0,
      "status": 1,
      "name": "test-instruction-to-create-with-approval-third-1-name-updated",
      "description": "test-instruction-to-create-with-approval-third-1-description-updated",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-third-1-step-1-updated",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-third-1-step-1-operation-1-updated",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-third-1-step-1-operation-1-description-updated"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-third-1-step-1-endpoint-updated"
        }
      ],
      "approval": {
        "user": {
          "_id": "user-to-instruction-approve-2",
          "name": "user-to-instruction-approve-2"
        },
        "requested_by": {
          "_id": "root",
          "name": "root"
        },
        "comment": "test-instruction-to-create-with-approval-third-1-comment-updated",
        "dismissed_by": {
          "_id": "user-to-instruction-approve-1",
          "name": "user-to-instruction-approve-1"
        },
        "dismiss_comment": "test-instruction-to-create-with-approval-third-1-dismiss-comment"
      }
    }
    """
