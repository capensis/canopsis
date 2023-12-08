Feature: instruction approval creation
  I need to be able to create an instruction with approval

  @concurrent
  Scenario: the user that requested the approval can change approval from username to role
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
    When I am admin
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}:
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
              "time_to_complete": {"value": 5, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-second-1-step-1-operation-1-description"
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
        "requested_by": {
          "_id": "root",
          "name": "root"
        }
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
              "time_to_complete": {"value": 5, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-second-2-step-1-operation-1-description"
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
        "requested_by": {
          "_id": "root",
          "name": "root"
        }
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
        "requested_by": {
          "_id": "root",
          "name": "root"
        }
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
      "name": "test-instruction-to-create-with-approval-second-5-name",
      "description": "test-instruction-to-create-with-approval-second-5-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-second-5-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-second-5-step-1-operation-1",
              "time_to_complete": {
                  "value": 1,
                  "unit": "s"
              },
              "description": "test-instruction-to-create-with-approval-second-5-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-second-5-step-1-endpoint"
        }
      ]
    }
    """
    Then the response key "approval" should not exist

  @concurrent
  Scenario: only the user from approval should be able to approve by role
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
        "role": "role-to-instruction-approve-2",
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
      "error": "role is not assigned to approval"
    }
    """
    When I am role-to-instruction-approve-1
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
      "error": "role is not assigned to approval"
    }
    """
    When I am role-to-instruction-approve-2
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
      "name": "test-instruction-to-create-with-approval-second-6-name",
      "description": "test-instruction-to-create-with-approval-second-6-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-second-6-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-second-6-step-1-operation-1",
              "time_to_complete": {
                  "value": 1,
                  "unit": "s"
              },
              "description": "test-instruction-to-create-with-approval-second-6-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-second-6-step-1-endpoint"
        }
      ]
    }
    """
    Then the response key "approval" should not exist

  @concurrent
  Scenario: only the user from approval should be able to dismiss by username
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
      "approve": false,
      "comment": "test-instruction-to-create-with-approval-second-7-dismiss-comment"
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
      "approve": false,
      "comment": "test-instruction-to-create-with-approval-second-7-dismiss-comment"
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
      "approve": false,
      "comment": "test-instruction-to-create-with-approval-second-7-dismiss-comment"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/instructions/{{ .instructionID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "type": 0,
      "status": 3,
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
        "user": {
          "_id": "user-to-instruction-approve-1",
          "name": "user-to-instruction-approve-1"
        },
        "requested_by": {
          "_id": "root",
          "name": "root"
        },
        "comment": "test comment",
        "dismissed_by": {
          "_id": "user-to-instruction-approve-1",
          "name": "user-to-instruction-approve-1"
        },
        "dismiss_comment": "test-instruction-to-create-with-approval-second-7-dismiss-comment"
      }
    }
    """

  @concurrent
  Scenario: only the user from approval should be able to dismiss by role
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
        "role": "role-to-instruction-approve-2",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}/approval:
    """json
    {
      "approve": false,
      "comment": "test-instruction-to-create-with-approval-second-8-dismiss-comment"
    }
    """
    Then the response code should be 403
    Then the response body should be:
    """json
    {
      "error": "role is not assigned to approval"
    }
    """
    When I am role-to-instruction-approve-1
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}/approval:
    """json
    {
      "approve": false,
      "comment": "test-instruction-to-create-with-approval-second-8-dismiss-comment"
    }
    """
    Then the response code should be 403
    Then the response body should be:
    """json
    {
      "error": "role is not assigned to approval"
    }
    """
    When I am role-to-instruction-approve-2
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}/approval:
    """json
    {
      "approve": false,
      "comment": "test-instruction-to-create-with-approval-second-8-dismiss-comment"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/instructions/{{ .instructionID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "type": 0,
      "status": 3,
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
        "role": {
          "_id": "role-to-instruction-approve-2",
          "name": "role-to-instruction-approve-2"
        },
        "requested_by": {
          "_id": "root",
          "name": "root"
        },
        "comment": "test comment",
        "dismissed_by": {
          "_id": "user-to-instruction-approve-2",
          "name": "user-to-instruction-approve-2"
        },
        "dismiss_comment": "test-instruction-to-create-with-approval-second-8-dismiss-comment"
      }

    }
    """

  @concurrent
  Scenario: given create request with approval request for auto instruction with user or role should return ok
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "name": "test-instruction-to-create-with-approval-second-9-name",
      "description": "test-instruction-to-create-with-approval-second-9-description",
      "enabled": true,
      "triggers": [
        {
          "type": "create"
        }
      ],
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "jobs": [
        {
          "stop_on_fail": true,
          "job": "test-job-to-instruction-edit-1"
        },
        {
          "job": "test-job-to-instruction-edit-2"
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
    Then the response body should contain:
    """json
    {
      "_id": "{{ .instructionID }}",
      "status": 1,
      "type": 1,
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "user-to-instruction-approve-1",
          "name": "user-to-instruction-approve-1"
        },
        "requested_by": {
          "_id": "root",
          "name": "root"
        }
      }
    }
    """
    When I do GET /api/v4/cat/instructions?search=test-instruction-to-create-with-approval-second-9
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "{{ .instructionID }}",
          "type": 1,
          "status": 1,
          "name": "test-instruction-to-create-with-approval-second-9-name",
          "description": "test-instruction-to-create-with-approval-second-9-description",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "enabled": true,
          "jobs": [
            {
              "stop_on_fail": true,
              "job": {
                "_id": "test-job-to-instruction-edit-1",
                "author": {
                  "_id": "root",
                  "name": "root"
                },
                "config": {
                  "_id": "test-job-config-to-edit-instruction",
                  "auth_token": "test-auth-token",
                  "author": {
                    "_id": "root",
                    "name": "root"
                  },
                  "host": "http://example.com",
                  "name": "test-job-config-to-edit-instruction-name",
                  "type": "rundeck"
                },
                "job_id": "test-job-to-instruction-edit-1-external-id",
                "name": "test-job-to-instruction-edit-1-name",
                "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
              }
            },
            {
              "job": {
                "_id": "test-job-to-instruction-edit-2",
                "author": {
                  "_id": "root",
                  "name": "root"
                },
                "config": {
                  "_id": "test-job-config-to-edit-instruction",
                  "auth_token": "test-auth-token",
                  "author": {
                    "_id": "root",
                    "name": "root"
                  },
                  "host": "http://example.com",
                  "name": "test-job-config-to-edit-instruction-name",
                  "type": "rundeck"
                },
                "job_id": "test-job-to-instruction-edit-2-external-id",
                "name": "test-job-to-instruction-edit-2-name",
                "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
              }
            }
          ],
          "approval": {
            "comment": "test comment",
            "user": {
              "_id": "user-to-instruction-approve-1",
              "name": "user-to-instruction-approve-1"
            },
            "requested_by": {
              "_id": "root",
              "name": "root"
            }
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I am role-to-instruction-approve-1
    When I do GET /api/v4/cat/instructions/{{ .instructionID }}/approval
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "original": {
        "_id": "{{ .instructionID }}",
        "type": 1,
        "status": 1,
        "name": "test-instruction-to-create-with-approval-second-9-name",
        "description": "test-instruction-to-create-with-approval-second-9-description",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "enabled": true,
        "jobs": [
          {
            "stop_on_fail": true,
            "job": {
              "_id": "test-job-to-instruction-edit-1",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "config": {
                "_id": "test-job-config-to-edit-instruction",
                "auth_token": "test-auth-token",
                "author": {
                  "_id": "root",
                  "name": "root"
                },
                "host": "http://example.com",
                "name": "test-job-config-to-edit-instruction-name",
                "type": "rundeck"
              },
              "job_id": "test-job-to-instruction-edit-1-external-id",
              "name": "test-job-to-instruction-edit-1-name",
              "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
            }
          },
          {
            "job": {
              "_id": "test-job-to-instruction-edit-2",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "config": {
                "_id": "test-job-config-to-edit-instruction",
                "auth_token": "test-auth-token",
                "author": {
                  "_id": "root",
                  "name": "root"
                },
                "host": "http://example.com",
                "name": "test-job-config-to-edit-instruction-name",
                "type": "rundeck"
              },
              "job_id": "test-job-to-instruction-edit-2-external-id",
              "name": "test-job-to-instruction-edit-2-name",
              "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
            }
          }
        ],
        "approval": {
          "comment": "test comment",
          "user": {
            "_id": "user-to-instruction-approve-1",
            "name": "user-to-instruction-approve-1"
          },
          "requested_by": {
            "_id": "root",
            "name": "root"
          }
        }
      },
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "user-to-instruction-approve-1",
          "name": "user-to-instruction-approve-1"
        },
        "requested_by": {
          "_id": "root",
          "name": "root"
        }
      }
    }
    """
