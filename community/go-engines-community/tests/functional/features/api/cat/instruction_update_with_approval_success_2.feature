Feature: instruction approval update
  I need to be able to update an instruction with approval

  @concurrent
  Scenario: PUT a valid instruction with approval with username request should return ok and valid approval response for auto instructions
    When I am manager
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-second-1:
    """json
    {
      "name": "test-instruction-to-update-with-approval-second-1-name",
      "description": "test-instruction-to-update-with-approval-second-1-description",
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
          "stop_on_fail": false,
          "job": "test-job-to-instruction-edit-1"
        },
        {
          "stop_on_fail": false,
          "job": "test-job-to-instruction-edit-2"
        }
      ],
      "approval": {
        "user": "user-to-instruction-approve-1",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-second-1/approval
    Then the response code should be 403
    Then the response body should be:
    """json
    {
      "error": "user is not assigned to approval"
    }
    """
    When I am role-to-instruction-approve-1
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-second-1/approval
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
        "requested_by": {
          "_id": "manageruser",
          "name": "manageruser"
        }
      },
      "original": {
        "_id": "test-instruction-to-update-with-approval-second-1",
        "type": 1,
        "status": 0,
        "name": "test-instruction-to-update-with-approval-second-1-name",
        "description": "test-instruction-to-update-with-approval-second-1-description",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "enabled": true,
        "timeout_after_execution": {
          "value": 2,
          "unit": "s"
        },
        "triggers": [
          {
            "type": "create"
          }
        ],
        "jobs": [
          {
            "job": {
              "_id": "test-job-to-instruction-edit-1",
              "name": "test-job-to-instruction-edit-1-name",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "config": {
                "_id": "test-job-config-to-edit-instruction",
                "name": "test-job-config-to-edit-instruction-name",
                "type": "rundeck",
                "host": "http://example.com",
                "author": {
                  "_id": "root",
                  "name": "root"
                },
                "auth_token": "test-auth-token"
              },
              "job_id": "test-job-to-instruction-edit-1-external-id",
              "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
            },
            "stop_on_fail": false
          }
        ],
        "created": 1596712203,
        "last_modified": 1596712203
      },
      "updated": {
        "type": 1,
        "status": 2,
        "name": "test-instruction-to-update-with-approval-second-1-name",
        "description": "test-instruction-to-update-with-approval-second-1-description",
        "author": {
          "_id": "manageruser",
          "name": "manageruser"
        },
        "enabled": true,
        "timeout_after_execution": {
          "value": 10,
          "unit": "m"
        },
        "triggers": [
          {
            "type": "create"
          }
        ],
        "jobs": [
          {
            "job": {
              "_id": "test-job-to-instruction-edit-1",
              "name": "test-job-to-instruction-edit-1-name",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "config": {
                "_id": "test-job-config-to-edit-instruction",
                "name": "test-job-config-to-edit-instruction-name",
                "type": "rundeck",
                "host": "http://example.com",
                "author": {
                  "_id": "root",
                  "name": "root"
                },
                "auth_token": "test-auth-token"
              },
              "job_id": "test-job-to-instruction-edit-1-external-id",
              "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
            },
            "stop_on_fail": false
          },
          {
            "job": {
              "_id": "test-job-to-instruction-edit-2",
              "name": "test-job-to-instruction-edit-2-name",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "config": {
                "_id": "test-job-config-to-edit-instruction",
                "name": "test-job-config-to-edit-instruction-name",
                "type": "rundeck",
                "host": "http://example.com",
                "author": {
                  "_id": "root",
                  "name": "root"
                },
                "auth_token": "test-auth-token"
              },
              "job_id": "test-job-to-instruction-edit-2-external-id",
              "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
            },
            "stop_on_fail": false
          }
        ]
      }
    }
    """
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-second-1/approval:
    """json
    {
      "approve": true
    }
    """
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-second-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "type": 1,
      "status": 0,
      "name": "test-instruction-to-update-with-approval-second-1-name",
      "description": "test-instruction-to-update-with-approval-second-1-description",
      "author": {
        "_id": "manageruser",
        "name": "manageruser"
      },
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "triggers": [
        {
          "type": "create"
        }
      ],
      "jobs": [
        {
          "job": {
            "_id": "test-job-to-instruction-edit-1",
            "name": "test-job-to-instruction-edit-1-name",
            "author": {
              "_id": "root",
              "name": "root"
            },
            "config": {
              "_id": "test-job-config-to-edit-instruction",
              "name": "test-job-config-to-edit-instruction-name",
              "type": "rundeck",
              "host": "http://example.com",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "auth_token": "test-auth-token"
            },
            "job_id": "test-job-to-instruction-edit-1-external-id",
            "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
          },
          "stop_on_fail": false
        },
        {
          "job": {
            "_id": "test-job-to-instruction-edit-2",
            "name": "test-job-to-instruction-edit-2-name",
            "author": {
              "_id": "root",
              "name": "root"
            },
            "config": {
              "_id": "test-job-config-to-edit-instruction",
              "name": "test-job-config-to-edit-instruction-name",
              "type": "rundeck",
              "host": "http://example.com",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "auth_token": "test-auth-token"
            },
            "job_id": "test-job-to-instruction-edit-2-external-id",
            "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
          },
          "stop_on_fail": false
        }
      ]
    }
    """
    When I do GET /api/v4/cat/instruction-stats?search=test-instruction-to-update-with-approval-second-1&from=1000000000&to=2000000000
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "last_executed_on": 1596712203,
          "avg_complete_time": 10,
          "rating": 3.5
        }
      ]
    }
    """

  @concurrent
  Scenario: only an author should be able to update a dismissed instruction
    When I am role-to-instruction-approve-1
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-second-2/approval:
    """json
    {
      "approve": false,
      "comment": "test-instruction-to-update-with-approval-second-2-dismiss-comment"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-second-2/approval
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "updated instruction not found"
    }
    """
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-second-2-updated/approval
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "original instruction not found"
    }
    """
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-second-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-update-with-approval-second-2",
      "name": "test-instruction-to-update-with-approval-second-2-name",
      "description": "test-instruction-to-update-with-approval-second-2-description"
    }
    """
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-second-2-updated
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """
    When I do GET /api/v4/cat/instructions?search=test-instruction-to-update-with-approval-second-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-instruction-to-update-with-approval-second-2",
          "name": "test-instruction-to-update-with-approval-second-2-name",
          "description": "test-instruction-to-update-with-approval-second-2-description"
        }
      ]
    }
    """
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-second-2:
    """json
    {
      "name": "test-instruction-to-update-with-approval-second-2-name-updated",
      "description": "test-instruction-to-update-with-approval-second-2-description-updated",
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
          "stop_on_fail": false,
          "job": "test-job-to-instruction-edit-1"
        },
        {
          "stop_on_fail": false,
          "job": "test-job-to-instruction-edit-2"
        }
      ],
      "approval": {
        "user": "user-to-instruction-approve-2",
        "comment": "test-instruction-to-update-with-approval-second-2-comment-updated"
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "status": 0,
      "name": "test-instruction-to-update-with-approval-second-2-name-updated",
      "description": "test-instruction-to-update-with-approval-second-2-description-updated",
      "enabled": true,
      "triggers": [
        {
          "type": "create"
        }
      ],
      "timeout_after_execution": {
        "value": 2,
        "unit": "s"
      },
      "jobs": [
        {
          "stop_on_fail": false,
          "job": {
            "_id": "test-job-to-instruction-edit-1"
          }
        }
      ]
    }
    """
    When I am admin
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-second-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "status": 4,
      "_id": "test-instruction-to-update-with-approval-second-2",
      "name": "test-instruction-to-update-with-approval-second-2-name-updated",
      "description": "test-instruction-to-update-with-approval-second-2-description-updated",
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
          "stop_on_fail": false,
          "job": {
            "_id": "test-job-to-instruction-edit-1"
          }
        },
        {
          "job": {
            "_id": "test-job-to-instruction-edit-2"
          }
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
        "comment": "test-instruction-to-update-with-approval-second-2-comment",
        "dismissed_by": {
          "_id": "user-to-instruction-approve-1",
          "name": "user-to-instruction-approve-1"
        },
        "dismiss_comment": "test-instruction-to-update-with-approval-second-2-dismiss-comment"
      }
    }
    """
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-second-2:
    """json
    {
      "name": "test-instruction-to-update-with-approval-second-2-name-updated",
      "description": "test-instruction-to-update-with-approval-second-2-description-updated",
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
          "stop_on_fail": false,
          "job": "test-job-to-instruction-edit-1"
        },
        {
          "stop_on_fail": false,
          "job": "test-job-to-instruction-edit-2"
        }
      ],
      "approval": {
        "user": "user-to-instruction-approve-2",
        "comment": "test-instruction-to-update-with-approval-second-2-comment-updated"
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "status": 2,
      "name": "test-instruction-to-update-with-approval-second-2-name-updated",
      "description": "test-instruction-to-update-with-approval-second-2-description-updated",
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
          "stop_on_fail": false,
          "job": {
            "_id": "test-job-to-instruction-edit-1"
          }
        },
        {
          "job": {
            "_id": "test-job-to-instruction-edit-2"
          }
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
        "comment": "test-instruction-to-update-with-approval-second-2-comment-updated",
        "dismissed_by": {
          "_id": "user-to-instruction-approve-1",
          "name": "user-to-instruction-approve-1"
        },
        "dismiss_comment": "test-instruction-to-update-with-approval-second-2-dismiss-comment"
      }
    }
    """
