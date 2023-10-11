Feature: instruction approval creation
  I need to be able to create an instruction with approval
  
  @concurrent
  Scenario: the users that didn't request the approval can update only name/description/enabled
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-second-1-name",
      "description": "test-instruction-to-create-with-approval-second-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-second-1-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-second-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-second-1-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-second-1-step-1-endpoint"
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
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-second-1-name-changed",
      "description": "test-instruction-to-create-with-approval-second-1-description-changed",
      "enabled": false,
      "timeout_after_execution": {
        "value": 12,
        "unit": "m"
      },
      "steps": [
        {
          "name": "new step",
          "operations": [
            {
              "name": "new operation",
              "time_to_complete": {"value": 55, "unit":"s"},
              "description": "new operation"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "new endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "user": "user-to-instruction-approve-1"
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "{{ .instructionID }}",
      "type": 0,
      "status": 1,
      "name": "test-instruction-to-create-with-approval-second-1-name-changed",
      "description": "test-instruction-to-create-with-approval-second-1-description-changed",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": false,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-second-1-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-second-1-step-1-operation-1",
              "time_to_complete": {
                  "value": 1,
                  "unit": "s"
              },
              "description": "test-instruction-to-create-with-approval-second-1-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-second-1-step-1-endpoint"
        }
      ]
    }
    """

  @concurrent
  Scenario: the user that requested the approval can update any field
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-second-2-name",
      "description": "test-instruction-to-create-with-approval-second-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-second-2-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-second-2-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-second-2-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-second-2-step-1-endpoint"
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
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-second-2-name-updated",
      "description": "test-instruction-to-create-with-approval-second-2-description-updated",
      "enabled": false,
      "timeout_after_execution": {
        "value": 11,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-second-2-step-1-updated",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-second-2-step-1-operation-1-updated",
              "time_to_complete": {"value": 5, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-second-2-step-1-operation-1-description-updated"
            }
          ],
          "stop_on_fail": false,
          "endpoint": "test-instruction-to-create-with-approval-second-2-step-1-endpoint-updated"
        }
      ],
      "approval": {
        "comment": "test comment",
        "user": "user-to-instruction-approve-1"
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "{{ .instructionID }}",
      "type": 0,
      "status": 1,
      "name": "test-instruction-to-create-with-approval-second-2-name-updated",
      "description": "test-instruction-to-create-with-approval-second-2-description-updated",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": false,
      "timeout_after_execution": {
        "value": 11,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-second-2-step-1-updated",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-second-2-step-1-operation-1-updated",
              "time_to_complete": {
                  "value": 5,
                  "unit": "s"
              },
              "description": "test-instruction-to-create-with-approval-second-2-step-1-operation-1-description-updated"
            }
          ],
          "stop_on_fail": false,
          "endpoint": "test-instruction-to-create-with-approval-second-2-step-1-endpoint-updated"
        }
      ],
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "user-to-instruction-approve-1",
          "name": "user-to-instruction-approve-1"
        },
        "requested_by": "root"
      }
    }
    """

  @concurrent
  Scenario: the users that didn't request the approval can't change approver
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-second-3-name",
      "description": "test-instruction-to-create-with-approval-second-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-second-3-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-second-3-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-second-3-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-second-3-step-1-endpoint"
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
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-second-3-name",
      "description": "test-instruction-to-create-with-approval-second-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-second-3-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-second-3-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-second-3-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-second-3-step-1-endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "user": "user-to-instruction-approve-2"
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "user-to-instruction-approve-1",
          "name": "user-to-instruction-approve-1"
        },
        "requested_by": "root"
      }
    }
    """

  @concurrent
  Scenario: the user that requested the approval can change approver username
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-second-4-name",
      "description": "test-instruction-to-create-with-approval-second-4-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-second-4-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-second-4-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-second-4-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-second-4-step-1-endpoint"
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
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-second-4-name",
      "description": "test-instruction-to-create-with-approval-second-4-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-second-4-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-second-4-step-1-operation-1",
              "time_to_complete": {"value": 5, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-second-4-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": false,
          "endpoint": "new endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "user": "user-to-instruction-approve-2"
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "user-to-instruction-approve-2",
          "name": "user-to-instruction-approve-2"
        },
        "requested_by": "root"
      }
    }
    """

  @concurrent
  Scenario: the user that requested the approval can change approval from username to role
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-second-5-name",
      "description": "test-instruction-to-create-with-approval-second-5-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-second-5-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-second-5-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-second-5-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-second-5-step-1-endpoint"
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
    When I am admin
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-second-5-name",
      "description": "test-instruction-to-create-with-approval-second-5-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-second-5-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-second-5-step-1-operation-1",
              "time_to_complete": {"value": 5, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-second-5-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": false,
          "endpoint": "new endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "role": "role-to-instruction-approve-2"
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "approval": {
        "comment": "test comment",
        "role": {
          "_id": "role-to-instruction-approve-2",
          "name": "role-to-instruction-approve-2"
        },
        "requested_by": "root"
      }
    }
    """

  @concurrent
  Scenario: the user that requested the approval can change approval from role to username
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-second-6-name",
      "description": "test-instruction-to-create-with-approval-second-6-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-second-6-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-second-6-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-second-6-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-second-6-step-1-endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "role": "role-to-instruction-approve-2"
      }
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I am admin
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-second-6-name",
      "description": "test-instruction-to-create-with-approval-second-6-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-second-6-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-second-6-step-1-operation-1",
              "time_to_complete": {"value": 5, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-second-6-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": false,
          "endpoint": "new endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "user": "user-to-instruction-approve-2"
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "user-to-instruction-approve-2",
          "name": "user-to-instruction-approve-2"
        },
        "requested_by": "root"
      }
    }
    """

  @concurrent
  Scenario: the user that didn't request the approval can't cancel approval
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-second-7-name",
      "description": "test-instruction-to-create-with-approval-second-7-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-second-7-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-second-7-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-second-7-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-second-7-step-1-endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "role": "role-to-instruction-approve-2"
      }
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I am manager
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-second-7-name",
      "description": "test-instruction-to-create-with-approval-second-7-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-second-7-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-second-7-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-second-7-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-second-7-step-1-endpoint"
        }
      ],
      "approval": null
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "approval": {
        "comment": "test comment",
        "role": {
          "_id": "role-to-instruction-approve-2",
          "name": "role-to-instruction-approve-2"
        },
        "requested_by": "root"
      }
    }
    """

  @concurrent
  Scenario: the user that request the approval can cancel approval
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-second-8-name",
      "description": "test-instruction-to-create-with-approval-second-8-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-second-8-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-second-8-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-second-8-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-second-8-step-1-endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "role": "role-to-instruction-approve-2"
      }
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I am admin
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-second-8-name",
      "description": "test-instruction-to-create-with-approval-second-8-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-second-8-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-second-8-step-1-operation-1",
              "time_to_complete": {"value": 5, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-second-8-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": false,
          "endpoint": "new endpoint"
        }
      ],
      "approval": null
    }
    """
    Then the response code should be 200
    Then the response key "approval" should not exist

  @concurrent
  Scenario: only the user from approval should be able to approve by username
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-second-9-name",
      "description": "test-instruction-to-create-with-approval-second-9-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-second-9-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-second-9-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-second-9-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-second-9-step-1-endpoint"
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
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}/approval:
    """json
    {
      "approve": true
    }
    """
    Then the response code should be 403
    Then the response body should be:
    """json
    {
      "error": "user is not assigned to approval"
    }
    """
    When I am role-to-instruction-approve-2
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}/approval:
    """json
    {
      "approve": true
    }
    """
    Then the response code should be 403
    Then the response body should be:
    """json
    {
      "error": "user is not assigned to approval"
    }
    """
    When I am role-to-instruction-approve-1
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}/approval:
    """json
    {
      "approve": true
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/instructions/{{ .instructionID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "{{ .instructionID }}",
      "type": 0,
      "status": 0,
      "name": "test-instruction-to-create-with-approval-second-9-name",
      "description": "test-instruction-to-create-with-approval-second-9-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-second-9-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-second-9-step-1-operation-1",
              "time_to_complete": {
                  "value": 1,
                  "unit": "s"
              },
              "description": "test-instruction-to-create-with-approval-second-9-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-second-9-step-1-endpoint"
        }
      ]
    }
    """
    Then the response key "approval" should not exist
