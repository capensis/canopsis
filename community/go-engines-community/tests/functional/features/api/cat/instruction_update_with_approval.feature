Feature: instruction approval update
  I need to be able to update an instruction with approval

  Scenario: PUT a valid instruction with approval with username request should return ok
    When I am authenticated with username "manageruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval:
    """json
    {
      "name": "test-instruction-to-update-with-approval-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-with-approval-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "new endpoint"
        }
      ],
      "approval": {
        "user": "approveruser",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval/approval
    Then the response code should be 403
    Then the response body should be:
    """json
    {
      "error": "user is not assigned to approval"
    }
    """
    When I am authenticated with username "approveruser" and password "test"
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval/approval
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "manageruser"
      },
      "original": {
        "_id": "test-instruction-to-update-with-approval",
        "type": 0,
        "status": 0,
        "alarm_patterns": null,
        "entity_patterns": [
          {
            "name": "test filter"
          }
        ],
        "name": "test-instruction-to-update-with-approval-name",
        "description": "test-instruction-to-update-with-approval-description",
        "author": {
          "_id": "test-user-author-1-id",
          "name": "test-user-author-1-username"
        },
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-to-update-with-approval-step-1-name",
            "operations": [
              {
                "name": "test-instruction-to-update-with-approval-step-1-operation-1-name",
                "time_to_complete": {
                  "value": 1,
                  "unit": "s"
                },
                "description": "test-instruction-to-update-with-approval-step-1-operation-1-description"
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-to-update-with-approval-step-1-endpoint"
          }
        ]
      },
      "updated": {
        "type": 0,
        "status": 2,
        "name": "test-instruction-to-update-with-approval-name",
        "entity_patterns": [
          {
            "name": "test filter"
          }
        ],
        "description": "test-instruction-to-update-with-approval-description",
        "author": {
          "_id": "manageruser",
          "name": "manageruser"
        },
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-to-update-with-approval-step-1-name",
            "operations": [
              {
                "name": "test-instruction-to-update-with-approval-step-1-operation-1-name",
                "time_to_complete": {"value": 1, "unit":"s"},
                "description": "test-instruction-to-update-with-approval-step-1-operation-1-description",
                "jobs": [
                  {
                    "_id": "test-job-to-instruction-edit-1",
                    "name": "test-job-to-instruction-edit-1-name",
                    "author": {
                      "_id": "test-user-author-1-id",
                      "name": "test-user-author-1-username"
                    },
                    "config": {
                      "_id": "test-job-config-to-edit-instruction",
                      "name": "test-job-config-to-edit-instruction-name",
                      "type": "rundeck",
                      "host": "http://example.com",
                      "auth_token": "test-auth-token"
                    },
                    "job_id": "test-job-to-instruction-edit-1-external-id",
                    "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                  }
                ]
              }
            ],
            "stop_on_fail": true,
            "endpoint": "new endpoint"
          }
        ],
        "approval": {
          "comment": "test comment",
          "user": {
            "_id": "approveruser",
            "name": "approveruser"
          },
          "requested_by": "manageruser"
        }
      }
    }
    """
    When I am authenticated with username "manageruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-2:
    """json
    {
      "name": "test-instruction-to-update-with-approval-2-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-with-approval-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-2-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-2-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-2-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "new endpoint 2"
        }
      ],
      "approval": {
        "role": "approver2",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-2/approval
    Then the response code should be 403
    Then the response body should be:
    """json
    {
      "error": "role is not assigned to approval"
    }
    """
    When I am authenticated with username "approveruser2" and password "test"
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-2/approval
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "approval": {
        "comment": "test comment",
        "role": {
          "_id": "approver2",
          "name": "approver2"
        },
        "requested_by": "manageruser"
      },
      "original": {
        "_id": "test-instruction-to-update-with-approval-2",
        "type": 0,
        "status": 0,
        "alarm_patterns": null,
        "entity_patterns": [
          {
            "name": "test filter"
          }
        ],
        "name": "test-instruction-to-update-with-approval-2-name",
        "description": "test-instruction-to-update-with-approval-2-description",
        "author": {
          "_id": "test-user-author-1-id",
          "name": "test-user-author-1-username"
        },
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-to-update-with-approval-2-step-1-name",
            "operations": [
              {
                "name": "test-instruction-to-update-with-approval-2-step-1-operation-1-name",
                "time_to_complete": {
                  "value": 1,
                  "unit": "s"
                },
                "description": "test-instruction-to-update-with-approval-2-step-1-operation-1-description"
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-to-update-with-approval-2-step-1-endpoint"
          }
        ]
      },
      "updated": {
        "type": 0,
        "status": 2,
        "name": "test-instruction-to-update-with-approval-2-name",
        "entity_patterns": [
          {
            "name": "test filter"
          }
        ],
        "description": "test-instruction-to-update-with-approval-2-description",
        "author": {
          "_id": "manageruser",
          "name": "manageruser"
        },
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-to-update-with-approval-2-step-1-name",
            "operations": [
              {
                "name": "test-instruction-to-update-with-approval-2-step-1-operation-1-name",
                "time_to_complete": {"value": 1, "unit":"s"},
                "description": "test-instruction-to-update-with-approval-2-step-1-operation-1-description",
                "jobs": [
                  {
                    "_id": "test-job-to-instruction-edit-1",
                    "name": "test-job-to-instruction-edit-1-name",
                    "author": {
                      "_id": "test-user-author-1-id",
                      "name": "test-user-author-1-username"
                    },
                    "config": {
                      "_id": "test-job-config-to-edit-instruction",
                      "name": "test-job-config-to-edit-instruction-name",
                      "type": "rundeck",
                      "host": "http://example.com",
                      "auth_token": "test-auth-token"
                    },
                    "job_id": "test-job-to-instruction-edit-1-external-id",
                    "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                  }
                ]
              }
            ],
            "stop_on_fail": true,
            "endpoint": "new endpoint 2"
          }
        ],
        "approval": {
          "comment": "test comment",
          "role": {
            "_id": "approver2",
            "name": "approver2"
          },
          "requested_by": "manageruser"
        }
      }
    }
    """
    When I am admin
    When I do GET /api/v4/cat/instructions?search=test-instruction-to-update-with-approval
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-instruction-to-update-with-approval",
          "status": 1,
          "approval": {
            "comment": "test comment",
            "user": {
              "_id": "approveruser",
              "name": "approveruser"
            },
            "requested_by": "manageruser"
          }
        },
        {
          "_id": "test-instruction-to-update-with-approval-2",
          "status": 1,
          "approval": {
            "comment": "test comment",
            "role": {
              "_id": "approver2",
              "name": "approver2"
            },
            "requested_by": "manageruser"
          }
        },
        {
          "_id": "test-instruction-to-update-with-approval-3",
          "status": 0
        },
        {
          "_id": "test-instruction-to-update-with-approval-4",
          "status": 0
        },
        {
          "_id": "test-instruction-to-update-with-approval-5",
          "status": 0
        },
        {
          "_id": "test-instruction-to-update-with-approval-6",
          "status": 0
        },
        {
          "_id": "test-instruction-to-update-with-approval-7",
          "status": 0
        },
        {
          "_id": "test-instruction-to-update-with-approval-8",
          "status": 0
        },
        {
          "_id": "test-instruction-to-update-with-approval-9",
          "status": 0
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 9
      }
    }
    """

  Scenario: PUT a valid instruction with approval request with a not found user should return error
    When I am authenticated with username "manageruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-3:
    """json
    {
      "name": "test-instruction-to-update-with-approval-3-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-with-approval-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-3-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-3-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-3-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "new endpoint 3"
        }
      ],
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
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-3/approval
    Then the response code should be 404

  Scenario: PUT a valid instruction with approval request with a username without approve right should return error
    When I am authenticated with username "manageruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-3:
    """json
    {
      "name": "test-instruction-to-update-with-approval-3-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-with-approval-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-3-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-3-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-3-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "new endpoint 3"
        }
      ],
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
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-3/approval
    Then the response code should be 404

  Scenario: PUT a valid instruction with approval request with a not found role should return error
    When I am authenticated with username "manageruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-3:
    """json
    {
      "name": "test-instruction-to-update-with-approval-3-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-with-approval-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-3-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-3-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-3-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "new endpoint 3"
        }
      ],
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
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-3/approval
    Then the response code should be 404

  Scenario: PUT a valid instruction with approval request with a role without approve right should return error
    When I am authenticated with username "manageruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-3:
    """json
    {
      "name": "test-instruction-to-update-with-approval-3-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-with-approval-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-3-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-3-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-3-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "new endpoint 3"
        }
      ],
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
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-3/approval
    Then the response code should be 404

  Scenario: PUT a valid instruction with approval request with a role without approve right should return error
    When I am authenticated with username "manageruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-3:
    """json
    {
      "name": "test-instruction-to-update-with-approval-3-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-with-approval-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-3-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-3-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-3-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "new endpoint 3"
        }
      ],
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
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-3/approval
    Then the response code should be 404

  Scenario: Requester should receive updated version on GET request, other users should receive original version on GET request
    When I am authenticated with username "manageruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-4:
    """json
    {
      "name": "test-instruction-to-update-with-approval-4-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-with-approval-4-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-4-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-4-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-4-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "new endpoint"
        }
      ],
      "approval": {
        "user": "approveruser",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200
    When I am authenticated with username "manageruser" and password "test"
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-update-with-approval-4",
      "type": 0,
      "status": 2,
      "name": "test-instruction-to-update-with-approval-4-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-with-approval-4-description",
      "author": {
        "_id": "manageruser",
        "name": "manageruser"
      },
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-4-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-4-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-4-step-1-operation-1-description",
              "jobs": [
                {
                  "_id": "test-job-to-instruction-edit-1",
                  "name": "test-job-to-instruction-edit-1-name",
                  "author": {
                    "_id": "test-user-author-1-id",
                    "name": "test-user-author-1-username"
                  },
                  "config": {
                    "_id": "test-job-config-to-edit-instruction",
                    "name": "test-job-config-to-edit-instruction-name",
                    "type": "rundeck",
                    "host": "http://example.com",
                    "auth_token": "test-auth-token"
                  },
                  "job_id": "test-job-to-instruction-edit-1-external-id",
                  "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                }
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "new endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "manageruser"
      }
    }
    """
    When I am authenticated with username "root" and password "test"
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-update-with-approval-4",
      "type": 0,
      "status": 0,
      "alarm_patterns": null,
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "name": "test-instruction-to-update-with-approval-4-name",
      "description": "test-instruction-to-update-with-approval-4-description",
      "author": {
        "_id": "test-user-author-1-id",
        "name": "test-user-author-1-username"
      },
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-4-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-4-step-1-operation-1-name",
              "time_to_complete": {
                "value": 1,
                "unit": "s"
              },
              "description": "test-instruction-to-update-with-approval-4-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-with-approval-4-step-1-endpoint"
        }
      ]
    }
    """

  Scenario: The users that didn't request the approval can update only name/description/enabled
    When I am authenticated with username "manageruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-5:
    """json
    {
      "name": "test-instruction-to-update-with-approval-5-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-with-approval-5-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-5-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-5-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-5-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "new endpoint"
        }
      ],
      "approval": {
        "user": "approveruser",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200  
    When I am authenticated with username "root" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-5:
    """json
    {
      "name": "test-instruction-to-update-with-approval-5-name-changed",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-with-approval-5-description-changed",
      "enabled": false,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-5-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-5-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-5-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "should be ignored"
        }
      ]
    }
    """
    Then the response code should be 200
    When I am authenticated with username "approveruser" and password "test"
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-5/approval
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "manageruser"
      },
      "original": {
        "_id": "test-instruction-to-update-with-approval-5",
        "type": 0,
        "status": 0,
        "alarm_patterns": null,
        "entity_patterns": [
          {
            "name": "test filter"
          }
        ],
        "name": "test-instruction-to-update-with-approval-5-name-changed",
        "description": "test-instruction-to-update-with-approval-5-description-changed",
        "author": {
          "_id": "test-user-author-1-id",
          "name": "test-user-author-1-username"
        },
        "enabled": false,
        "steps": [
          {
            "name": "test-instruction-to-update-with-approval-5-step-1-name",
            "operations": [
              {
                "name": "test-instruction-to-update-with-approval-5-step-1-operation-1-name",
                "time_to_complete": {
                  "value": 1,
                  "unit": "s"
                },
                "description": "test-instruction-to-update-with-approval-5-step-1-operation-1-description"
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-to-update-with-approval-5-step-1-endpoint"
          }
        ]
      },
      "updated": {
        "type": 0,
        "status": 2,
        "name": "test-instruction-to-update-with-approval-5-name",
        "entity_patterns": [
          {
            "name": "test filter"
          }
        ],
        "description": "test-instruction-to-update-with-approval-5-description",
        "author": {
          "_id": "manageruser",
          "name": "manageruser"
        },
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-to-update-with-approval-5-step-1-name",
            "operations": [
              {
                "name": "test-instruction-to-update-with-approval-5-step-1-operation-1-name",
                "time_to_complete": {"value": 1, "unit":"s"},
                "description": "test-instruction-to-update-with-approval-5-step-1-operation-1-description",
                "jobs": [
                  {
                    "_id": "test-job-to-instruction-edit-1",
                    "name": "test-job-to-instruction-edit-1-name",
                    "author": {
                      "_id": "test-user-author-1-id",
                      "name": "test-user-author-1-username"
                    },
                    "config": {
                      "_id": "test-job-config-to-edit-instruction",
                      "name": "test-job-config-to-edit-instruction-name",
                      "type": "rundeck",
                      "host": "http://example.com",
                      "auth_token": "test-auth-token"
                    },
                    "job_id": "test-job-to-instruction-edit-1-external-id",
                    "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                  }
                ]
              }
            ],
            "stop_on_fail": true,
            "endpoint": "new endpoint"
          }
        ],
        "approval": {
          "comment": "test comment",
          "user": {
            "_id": "approveruser",
            "name": "approveruser"
          },
          "requested_by": "manageruser"
        }
      }
    }
    """

  Scenario: The requester can update any updated fields
    When I am authenticated with username "manageruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-6:
    """json
    {
      "name": "test-instruction-to-update-with-approval-6-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-with-approval-6-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-6-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-6-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-6-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "new endpoint"
        }
      ],
      "approval": {
        "user": "approveruser",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-6:
    """json
    {
      "name": "test-instruction-to-update-with-approval-6-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-with-approval-6-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-6-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-6-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-6-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "shouldn't be ignored in updated"
        }
      ],
      "approval": {
        "user": "approveruser",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200
    When I am authenticated with username "approveruser" and password "test"
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-6/approval
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "manageruser"
      },
      "original": {
        "_id": "test-instruction-to-update-with-approval-6",
        "type": 0,
        "status": 0,
        "alarm_patterns": null,
        "entity_patterns": [
          {
            "name": "test filter"
          }
        ],
        "name": "test-instruction-to-update-with-approval-6-name",
        "description": "test-instruction-to-update-with-approval-6-description",
        "author": {
          "_id": "test-user-author-1-id",
          "name": "test-user-author-1-username"
        },
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-to-update-with-approval-6-step-1-name",
            "operations": [
              {
                "name": "test-instruction-to-update-with-approval-6-step-1-operation-1-name",
                "time_to_complete": {
                  "value": 1,
                  "unit": "s"
                },
                "description": "test-instruction-to-update-with-approval-6-step-1-operation-1-description"
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-to-update-with-approval-6-step-1-endpoint"
          }
        ]
      },
      "updated": {
        "type": 0,
        "status": 2,
        "name": "test-instruction-to-update-with-approval-6-name",
        "entity_patterns": [
          {
            "name": "test filter"
          }
        ],
        "description": "test-instruction-to-update-with-approval-6-description",
        "author":{
          "_id": "manageruser",
          "name": "manageruser"
        },
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-to-update-with-approval-6-step-1-name",
            "operations": [
              {
                "name": "test-instruction-to-update-with-approval-6-step-1-operation-1-name",
                "time_to_complete": {"value": 1, "unit":"s"},
                "description": "test-instruction-to-update-with-approval-6-step-1-operation-1-description",
                "jobs": [
                  {
                    "_id": "test-job-to-instruction-edit-1",
                    "name": "test-job-to-instruction-edit-1-name",
                    "author": {
                      "_id": "test-user-author-1-id",
                      "name": "test-user-author-1-username"
                    },
                    "config": {
                      "_id": "test-job-config-to-edit-instruction",
                      "name": "test-job-config-to-edit-instruction-name",
                      "type": "rundeck",
                      "host": "http://example.com",
                      "author": {
                        "_id": "test-user-author-1-id",
                        "name": "test-user-author-1-username"
                      },
                      "auth_token": "test-auth-token"
                    },
                    "job_id": "test-job-to-instruction-edit-1-external-id",
                    "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                  }
                ]
              }
            ],
            "stop_on_fail": true,
            "endpoint": "shouldn't be ignored in updated"
          }
        ],
        "approval": {
          "comment": "test comment",
          "user": {
            "_id": "approveruser",
            "name": "approveruser"
          },
          "requested_by": "manageruser"
        }
      }
    }
    """

  Scenario: The users that didn't request the approval couldn't change or remove the approval
    When I am authenticated with username "manageruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-7:
    """json
    {
      "name": "test-instruction-to-update-with-approval-7-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-with-approval-7-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-7-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-7-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-7-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "new endpoint"
        }
      ],
      "approval": {
        "user": "approveruser",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200  
    When I am authenticated with username "root" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-7:
    """json
    {
      "name": "test-instruction-to-update-with-approval-7-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-with-approval-7-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-7-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-7-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-7-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "new endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "role": "approver"
      }
    }
    """
    Then the response code should be 200
    When I am authenticated with username "approveruser" and password "test"
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-7/approval
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "manageruser"
      },
      "original": {
        "_id": "test-instruction-to-update-with-approval-7",
        "type": 0,
        "status": 0,
        "alarm_patterns": null,
        "entity_patterns": [
          {
            "name": "test filter"
          }
        ],
        "name": "test-instruction-to-update-with-approval-7-name",
        "description": "test-instruction-to-update-with-approval-7-description",
        "author": {
          "_id": "test-user-author-1-id",
          "name": "test-user-author-1-username"
        },
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-to-update-with-approval-7-step-1-name",
            "operations": [
              {
                "name": "test-instruction-to-update-with-approval-7-step-1-operation-1-name",
                "time_to_complete": {
                  "value": 1,
                  "unit": "s"
                },
                "description": "test-instruction-to-update-with-approval-7-step-1-operation-1-description"
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-to-update-with-approval-7-step-1-endpoint"
          }
        ]
      },
      "updated": {
        "type": 0,
        "status": 2,
        "name": "test-instruction-to-update-with-approval-7-name",
        "entity_patterns": [
          {
            "name": "test filter"
          }
        ],
        "description": "test-instruction-to-update-with-approval-7-description",
        "author": {
          "_id": "manageruser",
          "name": "manageruser"
        },
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-to-update-with-approval-7-step-1-name",
            "operations": [
              {
                "name": "test-instruction-to-update-with-approval-7-step-1-operation-1-name",
                "time_to_complete": {"value": 1, "unit":"s"},
                "description": "test-instruction-to-update-with-approval-7-step-1-operation-1-description",
                "jobs": [
                  {
                    "_id": "test-job-to-instruction-edit-1",
                    "name": "test-job-to-instruction-edit-1-name",
                    "author": {
                      "_id": "test-user-author-1-id",
                      "name": "test-user-author-1-username"
                    },
                    "config": {
                      "_id": "test-job-config-to-edit-instruction",
                      "name": "test-job-config-to-edit-instruction-name",
                      "type": "rundeck",
                      "host": "http://example.com",
                      "author": {
                        "_id": "test-user-author-1-id",
                        "name": "test-user-author-1-username"
                      },
                      "auth_token": "test-auth-token"
                    },
                    "job_id": "test-job-to-instruction-edit-1-external-id",
                    "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                  }
                ]
              }
            ],
            "stop_on_fail": true,
            "endpoint": "new endpoint"
          }
        ],
        "approval": {
          "comment": "test comment",
          "user": {
            "_id": "approveruser",
            "name": "approveruser"
          },
          "requested_by": "manageruser"
        }
      }
    }
    """
    When I am authenticated with username "root" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-7:
    """json
    {
      "name": "test-instruction-to-update-with-approval-7-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-with-approval-7-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-7-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-7-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-7-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "new endpoint"
        }
      ]
    }
    """
    Then the response code should be 200
    When I am authenticated with username "approveruser" and password "test"
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-7/approval
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "manageruser"
      },
      "original": {
        "_id": "test-instruction-to-update-with-approval-7",
        "type": 0,
        "status": 0,
        "alarm_patterns": null,
        "entity_patterns": [
          {
            "name": "test filter"
          }
        ],
        "name": "test-instruction-to-update-with-approval-7-name",
        "description": "test-instruction-to-update-with-approval-7-description",
        "author": {
          "_id": "test-user-author-1-id",
          "name": "test-user-author-1-username"
        },
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-to-update-with-approval-7-step-1-name",
            "operations": [
              {
                "name": "test-instruction-to-update-with-approval-7-step-1-operation-1-name",
                "time_to_complete": {
                  "value": 1,
                  "unit": "s"
                },
                "description": "test-instruction-to-update-with-approval-7-step-1-operation-1-description"
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-to-update-with-approval-7-step-1-endpoint"
          }
        ]
      },
      "updated": {
        "type": 0,
        "status": 2,
        "name": "test-instruction-to-update-with-approval-7-name",
        "entity_patterns": [
          {
            "name": "test filter"
          }
        ],
        "description": "test-instruction-to-update-with-approval-7-description",
        "author": {
          "_id": "manageruser",
          "name": "manageruser"
        },
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-to-update-with-approval-7-step-1-name",
            "operations": [
              {
                "name": "test-instruction-to-update-with-approval-7-step-1-operation-1-name",
                "time_to_complete": {"value": 1, "unit":"s"},
                "description": "test-instruction-to-update-with-approval-7-step-1-operation-1-description",
                "jobs": [
                  {
                    "_id": "test-job-to-instruction-edit-1",
                    "name": "test-job-to-instruction-edit-1-name",
                    "author": {
                      "_id": "test-user-author-1-id",
                      "name": "test-user-author-1-username"
                    },
                    "config": {
                      "_id": "test-job-config-to-edit-instruction",
                      "name": "test-job-config-to-edit-instruction-name",
                      "type": "rundeck",
                      "host": "http://example.com",
                      "author": {
                        "_id": "test-user-author-1-id",
                        "name": "test-user-author-1-username"
                      },
                      "auth_token": "test-auth-token"
                    },
                    "job_id": "test-job-to-instruction-edit-1-external-id",
                    "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                  }
                ]
              }
            ],
            "stop_on_fail": true,
            "endpoint": "new endpoint"
          }
        ],
        "approval": {
          "comment": "test comment",
          "user": {
            "_id": "approveruser",
            "name": "approveruser"
          },
          "requested_by": "manageruser"
        }
      }
    }
    """

  Scenario: The requester can update or remove the approval, after removal instruction should be updated
    When I am authenticated with username "manageruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-8:
    """json
    {
      "name": "test-instruction-to-update-with-approval-8-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-with-approval-8-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-8-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-8-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-8-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "new endpoint"
        }
      ],
      "approval": {
        "user": "approveruser",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200  
    When I am authenticated with username "manageruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-8:
    """json
    {
      "name": "test-instruction-to-update-with-approval-8-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-with-approval-8-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-8-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-8-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-8-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "new endpoint"
        }
      ],
      "approval": {
        "role": "approver2",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200
    When I am authenticated with username "approveruser" and password "test"
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-8/approval
    Then the response code should be 403
    When I am authenticated with username "approveruser2" and password "test"
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-8/approval
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "approval": {
        "comment": "test comment",
        "role": {
          "_id": "approver2",
          "name": "approver2"
        },
        "requested_by": "manageruser"
      },
      "original": {
        "_id": "test-instruction-to-update-with-approval-8",
        "type": 0,
        "status": 0,
        "alarm_patterns": null,
        "entity_patterns": [
          {
            "name": "test filter"
          }
        ],
        "name": "test-instruction-to-update-with-approval-8-name",
        "description": "test-instruction-to-update-with-approval-8-description",
        "author": {
          "_id": "test-user-author-1-id",
          "name": "test-user-author-1-username"
        },
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-to-update-with-approval-8-step-1-name",
            "operations": [
              {
                "name": "test-instruction-to-update-with-approval-8-step-1-operation-1-name",
                "time_to_complete": {
                  "value": 1,
                  "unit": "s"
                },
                "description": "test-instruction-to-update-with-approval-8-step-1-operation-1-description"
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-to-update-with-approval-8-step-1-endpoint"
          }
        ]
      },
      "updated": {
        "type": 0,
        "status": 2,
        "name": "test-instruction-to-update-with-approval-8-name",
        "entity_patterns": [
          {
            "name": "test filter"
          }
        ],
        "description": "test-instruction-to-update-with-approval-8-description",
        "author": {
          "_id": "manageruser",
          "name": "manageruser"
        },
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-to-update-with-approval-8-step-1-name",
            "operations": [
              {
                "name": "test-instruction-to-update-with-approval-8-step-1-operation-1-name",
                "time_to_complete": {"value": 1, "unit":"s"},
                "description": "test-instruction-to-update-with-approval-8-step-1-operation-1-description",
                "jobs": [
                  {
                    "_id": "test-job-to-instruction-edit-1",
                    "name": "test-job-to-instruction-edit-1-name",
                    "author": {
                      "_id": "test-user-author-1-id",
                      "name": "test-user-author-1-username"
                    },
                    "config": {
                      "_id": "test-job-config-to-edit-instruction",
                      "name": "test-job-config-to-edit-instruction-name",
                      "type": "rundeck",
                      "host": "http://example.com",
                      "author": {
                        "_id": "test-user-author-1-id",
                        "name": "test-user-author-1-username"
                      },
                      "auth_token": "test-auth-token"
                    },
                    "job_id": "test-job-to-instruction-edit-1-external-id",
                    "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                  },
                  {
                    "_id": "test-job-to-instruction-edit-2",
                    "name": "test-job-to-instruction-edit-2-name",
                    "author": {
                      "_id": "test-user-author-1-id",
                      "name": "test-user-author-1-username"
                    },
                    "config": {
                      "_id": "test-job-config-to-edit-instruction",
                      "name": "test-job-config-to-edit-instruction-name",
                      "type": "rundeck",
                      "host": "http://example.com",
                      "author": {
                        "_id": "test-user-author-1-id",
                        "name": "test-user-author-1-username"
                      },
                      "auth_token": "test-auth-token"
                    },
                    "job_id": "test-job-to-instruction-edit-2-external-id",
                    "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                  }
                ]
              }
            ],
            "stop_on_fail": true,
            "endpoint": "new endpoint"
          }
        ],
        "approval": {
          "comment": "test comment",
          "role": {
            "_id": "approver2",
            "name": "approver2"
          },
          "requested_by": "manageruser"
        }
      }
    }
    """
    When I am authenticated with username "manageruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-8:
    """json
    {
      "name": "test-instruction-to-update-with-approval-8-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-with-approval-8-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-8-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-8-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-8-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "new endpoint"
        }
      ]
    }
    """
    Then the response code should be 200
    When I am authenticated with username "approveruser2" and password "test"
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-8/approval
    Then the response code should be 404
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-update-with-approval-8",
      "type": 0,
      "status": 0,
      "name": "test-instruction-to-update-with-approval-8-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-with-approval-8-description",
      "author": {
        "_id": "manageruser",
        "name": "manageruser"
      },
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-8-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-8-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-8-step-1-operation-1-description",
              "jobs": [
                {
                  "_id": "test-job-to-instruction-edit-1",
                  "name": "test-job-to-instruction-edit-1-name",
                  "author": {
                    "_id": "test-user-author-1-id",
                    "name": "test-user-author-1-username"
                  },
                  "config": {
                    "_id": "test-job-config-to-edit-instruction",
                    "name": "test-job-config-to-edit-instruction-name",
                    "type": "rundeck",
                    "host": "http://example.com",
                    "author": {
                      "_id": "test-user-author-1-id",
                      "name": "test-user-author-1-username"
                    },
                    "auth_token": "test-auth-token"
                  },
                  "job_id": "test-job-to-instruction-edit-1-external-id",
                  "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                },
                {
                  "_id": "test-job-to-instruction-edit-2",
                  "name": "test-job-to-instruction-edit-2-name",
                  "author": {
                    "_id": "test-user-author-1-id",
                    "name": "test-user-author-1-username"
                  },
                  "config": {
                    "_id": "test-job-config-to-edit-instruction",
                    "name": "test-job-config-to-edit-instruction-name",
                    "type": "rundeck",
                    "host": "http://example.com",
                    "author": {
                      "_id": "test-user-author-1-id",
                      "name": "test-user-author-1-username"
                    },
                    "auth_token": "test-auth-token"
                  },
                  "job_id": "test-job-to-instruction-edit-2-external-id",
                  "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                }
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "new endpoint"
        }
      ]
    }
    """

  Scenario: Only the user from approval should be able to approve
    When I am authenticated with username "manageruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-9:
    """json
    {
      "name": "test-instruction-to-update-with-approval-9-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-with-approval-9-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-9-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-9-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-9-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "new endpoint 2"
        }
      ],
      "approval": {
        "role": "approver2",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200  
    When I am authenticated with username "root" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-9/approval:
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
    When I am authenticated with username "manageruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-9/approval:
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
    When I am authenticated with username "approveruser2" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-9/approval:
    """json
    {
      "approve": true
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-9/approval
    Then the response code should be 404
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-9
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "type": 0,
      "status": 0,
      "name": "test-instruction-to-update-with-approval-9-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-with-approval-9-description",
      "author": {
        "_id": "manageruser",
        "name": "manageruser"
      },
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-9-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-9-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-9-step-1-operation-1-description",
              "jobs": [
                {
                  "_id": "test-job-to-instruction-edit-1",
                  "name": "test-job-to-instruction-edit-1-name",
                  "author": {
                    "_id": "test-user-author-1-id",
                    "name": "test-user-author-1-username"
                  },
                  "config": {
                    "_id": "test-job-config-to-edit-instruction",
                    "name": "test-job-config-to-edit-instruction-name",
                    "type": "rundeck",
                    "host": "http://example.com",
                    "auth_token": "test-auth-token"
                  },
                  "job_id": "test-job-to-instruction-edit-1-external-id",
                  "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                }
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "new endpoint 2"
        }
      ]
    }
    """
    
  Scenario: Only the user from approval should be able to dismiss
    When I am authenticated with username "manageruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-3:
    """json
    {
      "name": "test-instruction-to-update-with-approval-3-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-with-approval-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-3-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-3-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-3-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "new endpoint 3"
        }
      ],
      "approval": {
        "user": "approveruser2",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200
    When I am authenticated with username "approveruser2" and password "test"
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-3/approval
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser2",
          "name": "approveruser2"
        },
        "requested_by": "manageruser"
      },
      "original": {
        "_id": "test-instruction-to-update-with-approval-3",
        "type": 0,
        "status": 0,
        "alarm_patterns": null,
        "entity_patterns": [
          {
            "name": "test filter"
          }
        ],
        "name": "test-instruction-to-update-with-approval-3-name",
        "description": "test-instruction-to-update-with-approval-3-description",
        "author": {
          "_id": "test-user-author-1-id",
          "name": "test-user-author-1-username"
        },
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-to-update-with-approval-3-step-1-name",
            "operations": [
              {
                "name": "test-instruction-to-update-with-approval-3-step-1-operation-1-name",
                "time_to_complete": {
                  "value": 1,
                  "unit": "s"
                },
                "description": "test-instruction-to-update-with-approval-3-step-1-operation-1-description"
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-to-update-with-approval-3-step-1-endpoint"
          }
        ]
      },
      "updated": {
        "type": 0,
        "status": 2,
        "name": "test-instruction-to-update-with-approval-3-name",
        "entity_patterns": [
          {
            "name": "test filter"
          }
        ],
        "description": "test-instruction-to-update-with-approval-3-description",
        "author": {
          "_id": "manageruser",
          "name": "manageruser"
        },
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-to-update-with-approval-3-step-1-name",
            "operations": [
              {
                "name": "test-instruction-to-update-with-approval-3-step-1-operation-1-name",
                "time_to_complete": {"value": 1, "unit":"s"},
                "description": "test-instruction-to-update-with-approval-3-step-1-operation-1-description",
                "jobs": [
                  {
                    "_id": "test-job-to-instruction-edit-1",
                    "name": "test-job-to-instruction-edit-1-name",
                    "author": {
                      "_id": "test-user-author-1-id",
                      "name": "test-user-author-1-username"
                    },
                    "config": {
                      "_id": "test-job-config-to-edit-instruction",
                      "name": "test-job-config-to-edit-instruction-name",
                      "type": "rundeck",
                      "host": "http://example.com",
                      "auth_token": "test-auth-token"
                    },
                    "job_id": "test-job-to-instruction-edit-1-external-id",
                    "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                  }
                ]
              }
            ],
            "stop_on_fail": true,
            "endpoint": "new endpoint 3"
          }
        ],
        "approval": {
          "comment": "test comment",
          "user": {
            "_id": "approveruser2",
            "name": "approveruser2"
          },
          "requested_by": "manageruser"
        }
      }
    }
    """
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-update-with-approval-3",
      "type": 0,
      "status": 0,
      "alarm_patterns": null,
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "name": "test-instruction-to-update-with-approval-3-name",
      "description": "test-instruction-to-update-with-approval-3-description",
      "author": {
        "_id": "test-user-author-1-id",
        "name": "test-user-author-1-username"
      },
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-3-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-3-step-1-operation-1-name",
              "time_to_complete": {
                "value": 1,
                "unit": "s"
              },
              "description": "test-instruction-to-update-with-approval-3-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-with-approval-3-step-1-endpoint"
        }
      ]
    }
    """
    When I am authenticated with username "manageruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-3/approval:
    """json
    {
      "approve": false
    }
    """
    Then the response code should be 403
    Then the response body should be:
    """json
    {
      "error": "user is not assigned to approval"
    }
    """
    When I am authenticated with username "approveruser2" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-3/approval:
    """json
    {
      "approve": false
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-3/approval
    Then the response code should be 404
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-update-with-approval-3",
      "type": 0,
      "status": 0,
      "alarm_patterns": null,
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "name": "test-instruction-to-update-with-approval-3-name",
      "description": "test-instruction-to-update-with-approval-3-description",
      "author": {
        "_id": "test-user-author-1-id",
        "name": "test-user-author-1-username"
      },
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-3-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-3-step-1-operation-1-name",
              "time_to_complete": {
                "value": 1,
                "unit": "s"
              },
              "description": "test-instruction-to-update-with-approval-3-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-with-approval-3-step-1-endpoint"
        }
      ]
    }
    """

  Scenario: PUT a valid instruction with approval with username request should return ok and valid approval response for auto instructions
    When I am authenticated with username "manageruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-run-auto-instruction-to-approve-update:
    """json
    {
      "name": "test-instruction-to-run-auto-instruction-to-approve-update-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-run-auto-instruction-to-approve-update-pattern"
        }
      ],
      "description": "test-instruction-to-run-auto-instruction-to-approve-update-description",
      "enabled": true,
      "priority": 1000,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "jobs": [
        {
          "stop_on_fail": false,
          "job": "test-job-to-run-auto-instruction-to-approve-update"
        },
        {
          "stop_on_fail": false,
          "job": "test-job-to-run-auto-instruction-to-approve-update-2"
        }
      ],
      "approval": {
        "user": "approveruser",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/instructions/test-instruction-to-run-auto-instruction-to-approve-update/approval
    Then the response code should be 403
    Then the response body should be:
    """json
    {
      "error": "user is not assigned to approval"
    }
    """
    When I am authenticated with username "approveruser" and password "test"
    When I do GET /api/v4/cat/instructions/test-instruction-to-run-auto-instruction-to-approve-update/approval
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "manageruser"
      },
      "original": {
        "_id": "test-instruction-to-run-auto-instruction-to-approve-update",
        "type": 1,
        "status": 0,
        "alarm_patterns": null,
        "entity_patterns": [
          {
            "name": "test-instruction-to-run-auto-instruction-to-approve-update-pattern"
          }
        ],
        "name": "est-instruction-to-run-auto-instruction-to-approve-update-name",
        "description": "test-instruction-to-run-auto-instruction-to-approve-update-description",
        "author": {
          "_id": "test-user-author-1-id",
          "name": "test-user-author-1-username"
        },
        "enabled": true,
        "timeout_after_execution": {
          "value": 2,
          "unit": "s"
        },
        "priority": 12,
        "jobs": [
          {
            "job": {
              "_id": "test-job-to-run-auto-instruction-to-approve-update",
              "name": "test-job-to-run-auto-instruction-to-approve-update-name",
              "author": {
                "_id": "test-user-author-1-id",
                "name": "test-user-author-1-username"
              },
              "config": {
                "_id": "test-job-config-to-run-auto-instruction",
                "name": "test-job-config-to-run-auto-instruction-name",
                "type": "rundeck",
                "host": "http://localhost:3000",
                "author": {
                  "_id": "test-user-author-1-id",
                  "name": "test-user-author-1-username"
                },
                "auth_token": "test-job-config-to-run-auto-instruction-token"
              },
              "job_id": "test-job-http-error",
              "payload": "{\"test-job-to-run-auto-instruction-to-approve-update-key\": \"test-job-to-run-auto-instruction-to-approve-update-val\"}"
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
        "alarm_patterns": null,
        "entity_patterns": [
          {
            "name": "test-instruction-to-run-auto-instruction-to-approve-update-pattern"
          }
        ],
        "name": "test-instruction-to-run-auto-instruction-to-approve-update-name",
        "description": "test-instruction-to-run-auto-instruction-to-approve-update-description",
        "author": {
          "_id": "manageruser",
          "name": "manageruser"
        },
        "enabled": true,
        "timeout_after_execution": {
          "value": 10,
          "unit": "m"
        },
        "priority": 1000,
        "jobs": [
          {
            "job": {
              "_id": "test-job-to-run-auto-instruction-to-approve-update",
              "name": "test-job-to-run-auto-instruction-to-approve-update-name",
              "author": {
                "_id": "test-user-author-1-id",
                "name": "test-user-author-1-username"
              },
              "config": {
                "_id": "test-job-config-to-run-auto-instruction",
                "name": "test-job-config-to-run-auto-instruction-name",
                "type": "rundeck",
                "host": "http://localhost:3000",
                "author": {
                  "_id": "test-user-author-1-id",
                  "name": "test-user-author-1-username"
                },
                "auth_token": "test-job-config-to-run-auto-instruction-token"
              },
              "job_id": "test-job-http-error",
              "payload": "{\"test-job-to-run-auto-instruction-to-approve-update-key\": \"test-job-to-run-auto-instruction-to-approve-update-val\"}"
            },
            "stop_on_fail": false
          },
          {
            "job": {
              "_id": "test-job-to-run-auto-instruction-to-approve-update-2",
              "name": "test-job-to-run-auto-instruction-to-approve-update-2-name",
              "author": {
                "_id": "test-user-author-1-id",
                "name": "test-user-author-1-username"
              },
              "config": {
                "_id": "test-job-config-to-run-auto-instruction",
                "name": "test-job-config-to-run-auto-instruction-name",
                "type": "rundeck",
                "host": "http://localhost:3000",
                "author": {
                  "_id": "test-user-author-1-id",
                  "name": "test-user-author-1-username"
                },
                "auth_token": "test-job-config-to-run-auto-instruction-token"
              },
              "job_id": "test-job-http-error",
              "payload": "{\"test-job-to-run-auto-instruction-to-approve-update-2-key\": \"test-job-to-run-auto-instruction-to-approve-update-2-val\"}"
            },
            "stop_on_fail": false
          }
        ]
      }
    }
    """
    When I do PUT /api/v4/cat/instructions/test-instruction-to-run-auto-instruction-to-approve-update/approval:
    """json
    {
      "approve": true
    }
    """
    When I do GET /api/v4/cat/instructions/test-instruction-to-run-auto-instruction-to-approve-update
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "type": 1,
      "status": 0,
      "alarm_patterns": null,
      "entity_patterns": [
        {
          "name": "test-instruction-to-run-auto-instruction-to-approve-update-pattern"
        }
      ],
      "name": "test-instruction-to-run-auto-instruction-to-approve-update-name",
      "description": "test-instruction-to-run-auto-instruction-to-approve-update-description",
      "author": {
        "_id": "manageruser",
        "name": "manageruser"
      },
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "priority": 1000,
      "jobs": [
        {
          "job": {
            "_id": "test-job-to-run-auto-instruction-to-approve-update",
            "name": "test-job-to-run-auto-instruction-to-approve-update-name",
            "author": {
              "_id": "test-user-author-1-id",
              "name": "test-user-author-1-username"
            },
            "config": {
              "_id": "test-job-config-to-run-auto-instruction",
              "name": "test-job-config-to-run-auto-instruction-name",
              "type": "rundeck",
              "host": "http://localhost:3000",
              "author": {
                "_id": "test-user-author-1-id",
                "name": "test-user-author-1-username"
              },
              "auth_token": "test-job-config-to-run-auto-instruction-token"
            },
            "job_id": "test-job-http-error",
            "payload": "{\"test-job-to-run-auto-instruction-to-approve-update-key\": \"test-job-to-run-auto-instruction-to-approve-update-val\"}"
          },
          "stop_on_fail": false
        },
        {
          "job": {
            "_id": "test-job-to-run-auto-instruction-to-approve-update-2",
            "name": "test-job-to-run-auto-instruction-to-approve-update-2-name",
            "author": {
              "_id": "test-user-author-1-id",
              "name": "test-user-author-1-username"
            },
            "config": {
              "_id": "test-job-config-to-run-auto-instruction",
              "name": "test-job-config-to-run-auto-instruction-name",
              "type": "rundeck",
              "host": "http://localhost:3000",
              "author": {
                "_id": "test-user-author-1-id",
                "name": "test-user-author-1-username"
              },
              "auth_token": "test-job-config-to-run-auto-instruction-token"
            },
            "job_id": "test-job-http-error",
            "payload": "{\"test-job-to-run-auto-instruction-to-approve-update-2-key\": \"test-job-to-run-auto-instruction-to-approve-update-2-val\"}"
          },
          "stop_on_fail": false
        }
      ]
    }
    """
