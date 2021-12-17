Feature: instruction approval creation
  I need to be able to create an instruction with approval

  Scenario: given create request with approval request with user or role should return ok
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-1-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-1-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-1-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-1-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-1-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-1-step-1-endpoint"
        }
      ],
      "approval": {
        "user": "approveruser",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-1-id",
      "status": 1,
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "root"
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-4-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-4-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-4-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-4-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-4-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-4-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-4-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-4-step-1-endpoint"
        }
      ],
      "approval": {
        "role": "approver",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-4-id",
      "status": 1,
      "approval": {
        "comment": "test comment",
        "role": {
          "_id": "approver",
          "name": "approver"
        },
        "requested_by": "root"
      }
    }
    """
    When I do GET /api/v4/cat/instructions?search=test-instruction-to-create-with-approval
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-instruction-to-create-with-approval-1-id",
          "type": 0,
          "status": 1,
          "alarm_patterns": null,
          "entity_patterns": [
            {
              "name": "test-instruction-to-create-with-approval-1-pattern"
            }
          ],
          "name": "test-instruction-to-create-with-approval-1-name",
          "description": "test-instruction-to-create-with-approval-1-description",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "enabled": true,
          "steps": [
            {
              "name": "test-instruction-to-create-with-approval-1-step-1",
              "operations": [
                {
                  "name": "test-instruction-to-create-with-approval-1-step-1-operation-1",
                  "time_to_complete": {
                      "value": 1,
                      "unit": "s"
                  },
                  "description": "test-instruction-to-create-with-approval-1-step-1-operation-1-description",
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
              "endpoint": "test-instruction-to-create-with-approval-1-step-1-endpoint"
            }
          ],
          "approval": {
            "comment": "test comment",
            "user": {
              "_id": "approveruser",
              "name": "approveruser"
            },
            "requested_by": "root"
          }
        },
        {
          "_id": "test-instruction-to-create-with-approval-4-id",
          "type": 0,
          "status": 1,
          "alarm_patterns": null,
          "entity_patterns": [
            {
              "name": "test-instruction-to-create-with-approval-4-pattern"
            }
          ],
          "name": "test-instruction-to-create-with-approval-4-name",
          "description": "test-instruction-to-create-with-approval-4-description",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "enabled": true,
          "steps": [
            {
              "name": "test-instruction-to-create-with-approval-4-step-1",
              "operations": [
                {
                  "name": "test-instruction-to-create-with-approval-4-step-1-operation-1",
                  "time_to_complete": {
                      "value": 1,
                      "unit": "s"
                  },
                  "description": "test-instruction-to-create-with-approval-4-step-1-operation-1-description",
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
              "endpoint": "test-instruction-to-create-with-approval-4-step-1-endpoint"
            }
          ],
          "approval": {
            "comment": "test comment",
            "role": {
              "_id": "approver",
              "name": "approver"
            },
            "requested_by": "root"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """

  Scenario: given create request with approval request with not existed username should return ok
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-2-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-2-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-2-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-2-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-2-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-2-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-2-step-1-endpoint"
        }
      ],
      "approval": {
        "user": "rootnotexist",
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

  Scenario: given create request with approval request with a username without approve right should return error
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-3-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-3-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-3-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-3-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-3-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-3-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-3-step-1-endpoint"
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
    
  Scenario: given create request with approval request with not existed role should return error
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-5-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-5-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-5-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-5-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-5-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-5-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-5-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-5-step-1-endpoint"
        }
      ],
      "approval": {
        "role": "adminnotexist",
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

  Scenario: given create request with approval request with a role without approve right should return error
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-6-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-6-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-6-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-6-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-6-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-6-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-6-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-6-step-1-endpoint"
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

  Scenario: given create request with approval request with role and username should return error
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-7-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-7-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-7-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-7-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-7-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-7-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-7-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-7-step-1-endpoint"
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

  Scenario: given create request with valid approval request with existed name should return error
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-8-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-1-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-8-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-8-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-8-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-8-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-8-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-8-step-1-endpoint"
        }
      ],
      "approval": {
        "role": "approver",
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

  Scenario: the approver should be able to get approval data by username
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-9-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-9-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-9-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-9-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-9-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-9-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-9-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-9-step-1-endpoint"
        }
      ],
      "approval": {
        "user": "approveruser",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-9-id",
      "status": 1,
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "root"
      }
    }
    """
    When I am authenticated with username "approveruser" and password "test"
    When I do GET /api/v4/cat/instructions/test-instruction-to-create-with-approval-9-id/approval
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "original": {
        "_id": "test-instruction-to-create-with-approval-9-id",
        "type": 0,
        "status": 1,
        "alarm_patterns": null,
        "entity_patterns": [
          {
            "name": "test-instruction-to-create-with-approval-9-pattern"
          }
        ],
        "name": "test-instruction-to-create-with-approval-9-name",
        "description": "test-instruction-to-create-with-approval-9-description",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-to-create-with-approval-9-step-1",
            "operations": [
              {
                "name": "test-instruction-to-create-with-approval-9-step-1-operation-1",
                "time_to_complete": {
                    "value": 1,
                    "unit": "s"
                },
                "description": "test-instruction-to-create-with-approval-9-step-1-operation-1-description",
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
            "endpoint": "test-instruction-to-create-with-approval-9-step-1-endpoint"
          }
        ]
      },
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "root"
      }
    }
    """

  Scenario: the approver should be able to get approval data by role
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-10-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-10-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-10-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-10-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-10-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-10-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-10-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-10-step-1-endpoint"
        }
      ],
      "approval": {
        "role": "approver",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-10-id",
      "status": 1,
      "approval": {
        "comment": "test comment",
        "role": {
          "_id": "approver",
          "name": "approver"
        },
        "requested_by": "root"
      }
    }
    """
    When I am authenticated with username "approveruser" and password "test"
    When I do GET /api/v4/cat/instructions/test-instruction-to-create-with-approval-10-id/approval
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "original": {
        "_id": "test-instruction-to-create-with-approval-10-id",
        "type": 0,
        "status": 1,
        "alarm_patterns": null,
        "entity_patterns": [
          {
            "name": "test-instruction-to-create-with-approval-10-pattern"
          }
        ],
        "name": "test-instruction-to-create-with-approval-10-name",
        "description": "test-instruction-to-create-with-approval-10-description",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-to-create-with-approval-10-step-1",
            "operations": [
              {
                "name": "test-instruction-to-create-with-approval-10-step-1-operation-1",
                "time_to_complete": {
                    "value": 1,
                    "unit": "s"
                },
                "description": "test-instruction-to-create-with-approval-10-step-1-operation-1-description",
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
            "endpoint": "test-instruction-to-create-with-approval-10-step-1-endpoint"
          }
        ]
      },
      "approval": {
        "comment": "test comment",
        "role": {
          "_id": "approver",
          "name": "approver"
        },
        "requested_by": "root"
      }
    }
    """

  Scenario: the user, which is not set in approval should receive 403
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-11-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-11-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-11-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-11-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-11-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-11-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-11-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-11-step-1-endpoint"
        }
      ],
      "approval": {
        "user": "approveruser",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-11-id",
      "status": 1,
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "root"
      }
    }
    """  
    When I am authenticated with username "manageruser" and password "test"
    When I do GET /api/v4/cat/instructions/test-instruction-to-create-with-approval-11-id/approval
    Then the response code should be 403
    Then the response body should contain:
    """json
    {
      "error": "user is not assigned to approval"
    }
    """

  Scenario: the user with a role, which is not set in approval should receive 403
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-12-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-12-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-12-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-12-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-12-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-12-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-12-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-12-step-1-endpoint"
        }
      ],
      "approval": {
        "role": "approver",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-12-id",
      "status": 1,
      "approval": {
        "comment": "test comment",
        "role": {
          "_id": "approver",
          "name": "approver"
        },
        "requested_by": "root"
      }
    }
    """  
    When I am authenticated with username "manageruser" and password "test"
    When I do GET /api/v4/cat/instructions/test-instruction-to-create-with-approval-12-id/approval
    Then the response code should be 403
    Then the response body should contain:
    """json
    {
      "error": "role is not assigned to approval"
    }
    """

  Scenario: if no approval return 404
    When I am authenticated with username "root" and password "test"
    When I do GET /api/v4/cat/instructions/test-instruction-to-get-1/approval
    Then the response code should be 404

  Scenario: should be possible for any user to get waiting for approval created instruction
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-13-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-13-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-13-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-13-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-13-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-13-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-13-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-13-step-1-endpoint"
        }
      ],
      "approval": {
        "user": "approveruser",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-13-id",
      "status": 1,
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "root"
      }
    }
    """  
    When I am authenticated with username "root" and password "test"
    When I do GET /api/v4/cat/instructions/test-instruction-to-create-with-approval-13-id
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-13-id",
      "type": 0,
      "status": 1,
      "alarm_patterns": null,
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-13-pattern"
        }
      ],
      "name": "test-instruction-to-create-with-approval-13-name",
      "description": "test-instruction-to-create-with-approval-13-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-13-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-13-step-1-operation-1",
              "time_to_complete": {
                  "value": 1,
                  "unit": "s"
              },
              "description": "test-instruction-to-create-with-approval-13-step-1-operation-1-description",
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
          "endpoint": "test-instruction-to-create-with-approval-13-step-1-endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "root"
      }
    }
    """
    When I am authenticated with username "manageruser" and password "test"
    When I do GET /api/v4/cat/instructions/test-instruction-to-create-with-approval-13-id
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-13-id",
      "type": 0,
      "status": 1,
      "alarm_patterns": null,
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-13-pattern"
        }
      ],
      "name": "test-instruction-to-create-with-approval-13-name",
      "description": "test-instruction-to-create-with-approval-13-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-13-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-13-step-1-operation-1",
              "time_to_complete": {
                  "value": 1,
                  "unit": "s"
              },
              "description": "test-instruction-to-create-with-approval-13-step-1-operation-1-description",
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
          "endpoint": "test-instruction-to-create-with-approval-13-step-1-endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "root"
      }
    }
    """
    When I am authenticated with username "approveruser" and password "test"
    When I do GET /api/v4/cat/instructions/test-instruction-to-create-with-approval-13-id
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-13-id",
      "type": 0,
      "status": 1,
      "alarm_patterns": null,
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-13-pattern"
        }
      ],
      "name": "test-instruction-to-create-with-approval-13-name",
      "description": "test-instruction-to-create-with-approval-13-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-13-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-13-step-1-operation-1",
              "time_to_complete": {
                  "value": 1,
                  "unit": "s"
              },
              "description": "test-instruction-to-create-with-approval-13-step-1-operation-1-description",
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
          "endpoint": "test-instruction-to-create-with-approval-13-step-1-endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "root"
      }
    }
    """

  Scenario: The users that didn't request the approval can update only name/description/enabled
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-14-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-14-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-14-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-14-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-14-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-14-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-14-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-14-step-1-endpoint"
        }
      ],
      "approval": {
        "user": "approveruser",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-14-id",
      "status": 1,
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "root"
      }
    }
    """
    When I am authenticated with username "manageruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-create-with-approval-14-id:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-14-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-14-name-changed",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-14-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-14-description-changed",
      "enabled": false,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-14-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-14-step-1-operation-1",
              "time_to_complete": {"value": 5, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-14-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": false,
          "endpoint": "new endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "user": "approveruser"
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-14-id",
      "type": 0,
      "status": 1,
      "alarm_patterns": null,
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-14-pattern"
        }
      ],
      "name": "test-instruction-to-create-with-approval-14-name-changed",
      "description": "test-instruction-to-create-with-approval-14-description-changed",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": false,
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-14-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-14-step-1-operation-1",
              "time_to_complete": {
                  "value": 1,
                  "unit": "s"
              },
              "description": "test-instruction-to-create-with-approval-14-step-1-operation-1-description",
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
          "endpoint": "test-instruction-to-create-with-approval-14-step-1-endpoint"
        }
      ]
    }
    """

  Scenario: The user that requested the approval can update any field
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-15-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-15-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-15-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-15-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-15-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-15-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-15-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-15-step-1-endpoint"
        }
      ],
      "approval": {
        "user": "approveruser",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-15-id",
      "status": 1,
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "root"
      }
    }
    """  
    When I do PUT /api/v4/cat/instructions/test-instruction-to-create-with-approval-15-id:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-15-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-15-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-15-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-15-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-15-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-15-step-1-operation-1",
              "time_to_complete": {"value": 5, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-15-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": false,
          "endpoint": "new endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "user": "approveruser"
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-15-id",
      "type": 0,
      "status": 1,
      "alarm_patterns": null,
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-15-pattern"
        }
      ],
      "name": "test-instruction-to-create-with-approval-15-name",
      "description": "test-instruction-to-create-with-approval-15-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-15-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-15-step-1-operation-1",
              "time_to_complete": {
                  "value": 5,
                  "unit": "s"
              },
              "description": "test-instruction-to-create-with-approval-15-step-1-operation-1-description",
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
          "stop_on_fail": false,
          "endpoint": "new endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "root"
      }
    }
    """

  Scenario: The users that didn't request the approval can't change approver
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-17-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-17-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-17-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-17-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-17-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-17-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-17-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-17-step-1-endpoint"
        }
      ],
      "approval": {
        "user": "approveruser",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-17-id",
      "status": 1,
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "root"
      }
    }
    """
    When I am authenticated with username "manageruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-create-with-approval-17-id:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-17-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-17-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-17-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-17-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-17-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-17-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-17-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-17-step-1-endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "user": "approveruser2"
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-17-id",
      "type": 0,
      "status": 1,
      "alarm_patterns": null,
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-17-pattern"
        }
      ],
      "name": "test-instruction-to-create-with-approval-17-name",
      "description": "test-instruction-to-create-with-approval-17-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-17-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-17-step-1-operation-1",
              "time_to_complete": {
                  "value": 1,
                  "unit": "s"
              },
              "description": "test-instruction-to-create-with-approval-17-step-1-operation-1-description",
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
          "endpoint": "test-instruction-to-create-with-approval-17-step-1-endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "root"
      }
    }
    """

  Scenario: The user that requested the approval can change approver username
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-18-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-18-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-18-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-18-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-18-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-18-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-18-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-18-step-1-endpoint"
        }
      ],
      "approval": {
        "user": "approveruser",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-18-id",
      "status": 1,
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "root"
      }
    }
    """  
    When I am authenticated with username "root" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-create-with-approval-18-id:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-18-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-18-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-18-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-18-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-18-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-18-step-1-operation-1",
              "time_to_complete": {"value": 5, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-18-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": false,
          "endpoint": "new endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "user": "approveruser2"
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-18-id",
      "type": 0,
      "status": 1,
      "alarm_patterns": null,
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-18-pattern"
        }
      ],
      "name": "test-instruction-to-create-with-approval-18-name",
      "description": "test-instruction-to-create-with-approval-18-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-18-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-18-step-1-operation-1",
              "time_to_complete": {
                  "value": 5,
                  "unit": "s"
              },
              "description": "test-instruction-to-create-with-approval-18-step-1-operation-1-description",
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
          "stop_on_fail": false,
          "endpoint": "new endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser2",
          "name": "approveruser2"
        },
        "requested_by": "root"
      }
    }
    """

  Scenario: The user that requested the approval can change approval from username to role
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-19-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-19-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-19-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-19-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-19-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-19-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-19-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-19-step-1-endpoint"
        }
      ],
      "approval": {
        "user": "approveruser",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-19-id",
      "status": 1,
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "root"
      }
    }
    """  
    When I am authenticated with username "root" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-create-with-approval-19-id:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-19-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-19-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-19-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-19-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-19-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-19-step-1-operation-1",
              "time_to_complete": {"value": 5, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-19-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": false,
          "endpoint": "new endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "role": "approver2"
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-19-id",
      "type": 0,
      "status": 1,
      "alarm_patterns": null,
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-19-pattern"
        }
      ],
      "name": "test-instruction-to-create-with-approval-19-name",
      "description": "test-instruction-to-create-with-approval-19-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-19-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-19-step-1-operation-1",
              "time_to_complete": {
                  "value": 5,
                  "unit": "s"
              },
              "description": "test-instruction-to-create-with-approval-19-step-1-operation-1-description",
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
          "stop_on_fail": false,
          "endpoint": "new endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "role": {
          "_id": "approver2",
          "name": "approver2"
        },
        "requested_by": "root"
      }
    }
    """

  Scenario: The user that requested the approval can change approval from role to username
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-20-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-20-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-20-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-20-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-20-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-20-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-20-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-20-step-1-endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "role": "approver2"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-20-id",
      "status": 1,
      "approval": {
        "comment": "test comment",
        "role": {
          "_id": "approver2",
          "name": "approver2"
        },
        "requested_by": "root"
      }
    }
    """    
    When I am authenticated with username "root" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-create-with-approval-20-id:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-20-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-20-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-20-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-20-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-20-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-20-step-1-operation-1",
              "time_to_complete": {"value": 5, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-20-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": false,
          "endpoint": "new endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "user": "approveruser2"
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-20-id",
      "type": 0,
      "status": 1,
      "alarm_patterns": null,
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-20-pattern"
        }
      ],
      "name": "test-instruction-to-create-with-approval-20-name",
      "description": "test-instruction-to-create-with-approval-20-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-20-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-20-step-1-operation-1",
              "time_to_complete": {
                  "value": 5,
                  "unit": "s"
              },
              "description": "test-instruction-to-create-with-approval-20-step-1-operation-1-description",
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
          "stop_on_fail": false,
          "endpoint": "new endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser2",
          "name": "approveruser2"
        },
        "requested_by": "root"
      }
    }
    """
    
  Scenario:  The user that didn't request the approval can't cancel approval
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-21-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-21-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-21-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-21-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-21-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-21-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-21-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-21-step-1-endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "role": "approver2"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-21-id",
      "status": 1,
      "approval": {
        "comment": "test comment",
        "role": {
          "_id": "approver2",
          "name": "approver2"
        },
        "requested_by": "root"
      }
    }
    """  
    When I am authenticated with username "manageruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-create-with-approval-21-id:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-21-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-21-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-21-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-21-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-21-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-21-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-21-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-21-step-1-endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "user": "approveruser"
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-21-id",
      "type": 0,
      "status": 1,
      "alarm_patterns": null,
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-21-pattern"
        }
      ],
      "name": "test-instruction-to-create-with-approval-21-name",
      "description": "test-instruction-to-create-with-approval-21-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-21-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-21-step-1-operation-1",
              "time_to_complete": {
                  "value": 1,
                  "unit": "s"
              },
              "description": "test-instruction-to-create-with-approval-21-step-1-operation-1-description",
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
          "endpoint": "test-instruction-to-create-with-approval-21-step-1-endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "role": {
          "_id": "approver2",
          "name": "approver2"
        },
        "requested_by": "root"
      }
    }
    """

  Scenario: The user that request the approval can cancel approval
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-22-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-22-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-22-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-22-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-22-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-22-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-22-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-22-step-1-endpoint"
        }
      ],
      "approval": {
        "comment": "test comment",
        "role": "approver2"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-22-id",
      "status": 1,
      "approval": {
        "comment": "test comment",
        "role": {
          "_id": "approver2",
          "name": "approver2"
        },
        "requested_by": "root"
      }
    }
    """   
    When I am authenticated with username "root" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-create-with-approval-22-id:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-22-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-22-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-22-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-22-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-22-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-22-step-1-operation-1",
              "time_to_complete": {"value": 5, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-22-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": false,
          "endpoint": "new endpoint"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-22-id",
      "type": 0,
      "status": 0,
      "alarm_patterns": null,
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-22-pattern"
        }
      ],
      "name": "test-instruction-to-create-with-approval-22-name",
      "description": "test-instruction-to-create-with-approval-22-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-22-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-22-step-1-operation-1",
              "time_to_complete": {
                  "value": 5,
                  "unit": "s"
              },
              "description": "test-instruction-to-create-with-approval-22-step-1-operation-1-description",
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
          "stop_on_fail": false,
          "endpoint": "new endpoint"
        }
      ]
    }
    """

  Scenario: Only the user from approval should be able to approve by username
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-5-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-5-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-5-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-5-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-5-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-5-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-5-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-5-step-1-endpoint"
        }
      ],
      "approval": {
        "user": "approveruser",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-5-id",
      "status": 1,
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "root"
      }
    }
    """
    When I do PUT /api/v4/cat/instructions/test-instruction-to-create-with-approval-5-id/approval:
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
    When I am authenticated with username "approveruser2" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-create-with-approval-5-id/approval:
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
    When I am authenticated with username "approveruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-create-with-approval-5-id/approval:
    """json
    {
      "approve": true
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/instructions/test-instruction-to-create-with-approval-5-id
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-5-id",
      "type": 0,
      "status": 0,
      "alarm_patterns": null,
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-5-pattern"
        }
      ],
      "name": "test-instruction-to-create-with-approval-5-name",
      "description": "test-instruction-to-create-with-approval-5-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-5-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-5-step-1-operation-1",
              "time_to_complete": {
                  "value": 1,
                  "unit": "s"
              },
              "description": "test-instruction-to-create-with-approval-5-step-1-operation-1-description",
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
          "endpoint": "test-instruction-to-create-with-approval-5-step-1-endpoint"
        }
      ]
    }
    """
    
  Scenario: Only the user from approval should be able to approve by role
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-6-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-6-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-6-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-6-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-6-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-6-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-6-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-6-step-1-endpoint"
        }
      ],
      "approval": {
        "role": "approver2",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-6-id",
      "status": 1,
      "approval": {
        "comment": "test comment",
        "role": {
          "_id": "approver2",
          "name": "approver2"
        },
        "requested_by": "root"
      }
    }
    """
    When I do PUT /api/v4/cat/instructions/test-instruction-to-create-with-approval-6-id/approval:
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
    When I am authenticated with username "approveruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-create-with-approval-6-id/approval:
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
    When I do PUT /api/v4/cat/instructions/test-instruction-to-create-with-approval-6-id/approval:
    """json
    {
      "approve": true
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/instructions/test-instruction-to-create-with-approval-6-id
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-6-id",
      "type": 0,
      "status": 0,
      "alarm_patterns": null,
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-6-pattern"
        }
      ],
      "name": "test-instruction-to-create-with-approval-6-name",
      "description": "test-instruction-to-create-with-approval-6-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-6-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-6-step-1-operation-1",
              "time_to_complete": {
                  "value": 1,
                  "unit": "s"
              },
              "description": "test-instruction-to-create-with-approval-6-step-1-operation-1-description",
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
          "endpoint": "test-instruction-to-create-with-approval-6-step-1-endpoint"
        }
      ]
    }
    """
    
  Scenario: Only the user from approval should be able to dismiss by username
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-7-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-7-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-7-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-7-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-7-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-7-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-7-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-7-step-1-endpoint"
        }
      ],
      "approval": {
        "user": "approveruser",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-7-id",
      "status": 1,
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "root"
      }
    }
    """
    When I do PUT /api/v4/cat/instructions/test-instruction-to-create-with-approval-7-id/approval:
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
    When I do PUT /api/v4/cat/instructions/test-instruction-to-create-with-approval-7-id/approval:
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
    When I am authenticated with username "approveruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-create-with-approval-7-id/approval:
    """json
    {
      "approve": false
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/instructions/test-instruction-to-create-with-approval-7-id
    Then the response code should be 404
    
  Scenario: Only the user from approval should be able to dismiss by role
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-8-id",
      "type": 0,
      "name": "test-instruction-to-create-with-approval-8-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-8-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-8-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-8-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-8-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-8-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-8-step-1-endpoint"
        }
      ],
      "approval": {
        "role": "approver2",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-8-id",
      "status": 1,
      "approval": {
        "comment": "test comment",
        "role": {
          "_id": "approver2",
          "name": "approver2"
        },
        "requested_by": "root"
      }
    }
    """
    When I do PUT /api/v4/cat/instructions/test-instruction-to-create-with-approval-8-id/approval:
    """json
    {
      "approve": false
    }
    """
    Then the response code should be 403
    Then the response body should be:
    """json
    {
      "error": "role is not assigned to approval"
    }
    """
    When I am authenticated with username "approveruser" and password "test"
    When I do PUT /api/v4/cat/instructions/test-instruction-to-create-with-approval-8-id/approval:
    """json
    {
      "approve": false
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
    When I do PUT /api/v4/cat/instructions/test-instruction-to-create-with-approval-8-id/approval:
    """json
    {
      "approve": false
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/instructions/test-instruction-to-create-with-approval-8-id
    Then the response code should be 404
    
  Scenario: given create request with approval request for auto instruction with user or role should return ok
    When I am authenticated with username "root" and password "test"
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-16-id",
      "type": 1,
      "name": "test-instruction-to-create-with-approval-16-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-with-approval-16-pattern"
        }
      ],
      "description": "test-instruction-to-create-with-approval-16-description",
      "enabled": true,
      "priority": 100,
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
        "user": "approveruser",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-with-approval-16-id",
      "status": 1,
      "type": 1,
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "root"
      }
    }
    """
    When I do GET /api/v4/cat/instructions?search=test-instruction-to-create-with-approval-16
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-instruction-to-create-with-approval-16-id",
          "type": 1,
          "status": 1,
          "alarm_patterns": null,
          "entity_patterns": [
            {
              "name": "test-instruction-to-create-with-approval-16-pattern"
            }
          ],
          "name": "test-instruction-to-create-with-approval-16-name",
          "description": "test-instruction-to-create-with-approval-16-description",
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
                  "_id": "test-user-author-1-id",
                  "name": "test-user-author-1-username"
                },
                "config": {
                  "_id": "test-job-config-to-edit-instruction",
                  "auth_token": "test-auth-token",
                  "author": {
                    "_id": "test-user-author-1-id",
                    "name": "test-user-author-1-username"
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
                  "_id": "test-user-author-1-id",
                  "name": "test-user-author-1-username"
                },
                "config": {
                  "_id": "test-job-config-to-edit-instruction",
                  "auth_token": "test-auth-token",
                  "author": {
                    "_id": "test-user-author-1-id",
                    "name": "test-user-author-1-username"
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
              "_id": "approveruser",
              "name": "approveruser"
            },
            "requested_by": "root"
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
    When I am authenticated with username "approveruser" and password "test"
    When I do GET /api/v4/cat/instructions/test-instruction-to-create-with-approval-16-id/approval
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "original": {
        "_id": "test-instruction-to-create-with-approval-16-id",
        "type": 1,
        "status": 1,
        "alarm_patterns": null,
        "entity_patterns": [
          {
            "name": "test-instruction-to-create-with-approval-16-pattern"
          }
        ],
        "name": "test-instruction-to-create-with-approval-16-name",
        "description": "test-instruction-to-create-with-approval-16-description",
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
                "_id": "test-user-author-1-id",
                "name": "test-user-author-1-username"
              },
              "config": {
                "_id": "test-job-config-to-edit-instruction",
                "auth_token": "test-auth-token",
                "author": {
                  "_id": "test-user-author-1-id",
                  "name": "test-user-author-1-username"
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
                "_id": "test-user-author-1-id",
                "name": "test-user-author-1-username"
              },
              "config": {
                "_id": "test-job-config-to-edit-instruction",
                "auth_token": "test-auth-token",
                "author": {
                  "_id": "test-user-author-1-id",
                  "name": "test-user-author-1-username"
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
            "_id": "approveruser",
            "name": "approveruser"
          },
          "requested_by": "root"
        }
      },
      "approval": {
        "comment": "test comment",
        "user": {
          "_id": "approveruser",
          "name": "approveruser"
        },
        "requested_by": "root"
      }
    }
    """
