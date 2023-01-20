Feature: instruction approval creation
  I need to be able to create an instruction with approval

  Scenario: given create request with approval request with user should return ok
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-1-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-create-with-approval-1-pattern"
            }
          }
        ]
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-create-with-approval-1-pattern"
            }
          }
        ]
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
        "user": "user-to-instruction-approve-1",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "status": 1,
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
    When I do GET /api/v4/cat/instructions?search=test-instruction-to-create-with-approval-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "type": 0,
          "status": 1,
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-instruction-to-create-with-approval-1-pattern"
                }
              }
            ]
          ],
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-instruction-to-create-with-approval-1-pattern"
                }
              }
            ]
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
                    {
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
              "_id": "user-to-instruction-approve-1",
              "name": "user-to-instruction-approve-1"
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

  Scenario: given create request with approval request with user or role should return ok
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-2-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-create-with-approval-2-pattern"
            }
          }
        ]
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-create-with-approval-2-pattern"
            }
          }
        ]
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
        "role": "role-to-instruction-approve-1",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "status": 1,
      "approval": {
        "comment": "test comment",
        "role": {
          "_id": "role-to-instruction-approve-1",
          "name": "role-to-instruction-approve-1"
        },
        "requested_by": "root"
      }
    }
    """
    When I do GET /api/v4/cat/instructions?search=test-instruction-to-create-with-approval-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "type": 0,
          "status": 1,
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-instruction-to-create-with-approval-2-pattern"
                }
              }
            ]
          ],
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-instruction-to-create-with-approval-2-pattern"
                }
              }
            ]
          ],
          "name": "test-instruction-to-create-with-approval-2-name",
          "description": "test-instruction-to-create-with-approval-2-description",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "enabled": true,
          "steps": [
            {
              "name": "test-instruction-to-create-with-approval-2-step-1",
              "operations": [
                {
                  "name": "test-instruction-to-create-with-approval-2-step-1-operation-1",
                  "time_to_complete": {
                      "value": 1,
                      "unit": "s"
                  },
                  "description": "test-instruction-to-create-with-approval-2-step-1-operation-1-description",
                  "jobs": [
                    {
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
                    {
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
                    }
                  ]
                }
              ],
              "stop_on_fail": true,
              "endpoint": "test-instruction-to-create-with-approval-2-step-1-endpoint"
            }
          ],
          "approval": {
            "comment": "test comment",
            "role": {
              "_id": "role-to-instruction-approve-1",
              "name": "role-to-instruction-approve-1"
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

  Scenario: given create request with approval request with not existed username should return ok
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

  Scenario: given create request with valid approval request with existed name should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-check-unique-name",
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
              "description": "test-instruction-to-create-with-approval-3-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-3-step-1-endpoint"
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

  Scenario: the approver should be able to get approval data by username
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-4-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-create-with-approval-4-pattern"
            }
          }
        ]
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-create-with-approval-4-pattern"
            }
          }
        ]
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
              "description": "test-instruction-to-create-with-approval-4-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-4-step-1-endpoint"
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
    When I am role-to-instruction-approve-1
    When I do GET /api/v4/cat/instructions/{{ .instructionID }}/approval
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "original": {
        "_id": "{{ .instructionID }}",
        "type": 0,
        "status": 1,
        "entity_pattern": [
          [
            {
              "field": "name",
              "cond": {
                "type": "eq",
                "value": "test-instruction-to-create-with-approval-4-pattern"
              }
            }
          ]
        ],
        "alarm_pattern": [
          [
            {
              "field": "v.component",
              "cond": {
                "type": "eq",
                "value": "test-instruction-to-create-with-approval-4-pattern"
              }
            }
          ]
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
                "description": "test-instruction-to-create-with-approval-4-step-1-operation-1-description"
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-to-create-with-approval-4-step-1-endpoint"
          }
        ]
      },
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

  Scenario: the approver should be able to get approval data by role
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-10-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-create-with-approval-10-pattern"
            }
          }
        ]
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-create-with-approval-10-pattern"
            }
          }
        ]
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
              "description": "test-instruction-to-create-with-approval-10-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-10-step-1-endpoint"
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
    When I am role-to-instruction-approve-1
    When I do GET /api/v4/cat/instructions/{{ .instructionID }}/approval
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "original": {
        "_id": "{{ .instructionID }}",
        "type": 0,
        "status": 1,
        "entity_pattern": [
          [
            {
              "field": "name",
              "cond": {
                "type": "eq",
                "value": "test-instruction-to-create-with-approval-10-pattern"
              }
            }
          ]
        ],
        "alarm_pattern": [
          [
            {
              "field": "v.component",
              "cond": {
                "type": "eq",
                "value": "test-instruction-to-create-with-approval-10-pattern"
              }
            }
          ]
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
                "description": "test-instruction-to-create-with-approval-10-step-1-operation-1-description"
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
          "_id": "role-to-instruction-approve-1",
          "name": "role-to-instruction-approve-1"
        },
        "requested_by": "root"
      }
    }
    """

  Scenario: the user, which is not set in approval should receive 403
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-11-name",
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
              "description": "test-instruction-to-create-with-approval-11-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-6-step-1-endpoint"
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

  Scenario: the user with a role, which is not set in approval should receive 403
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-12-name",
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
              "description": "test-instruction-to-create-with-approval-12-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-12-step-1-endpoint"
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

  Scenario: if no approval return 404
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-23-name",
      "description": "test-instruction-to-create-with-approval-23-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-23-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-23-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-23-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-23-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/cat/instructions/{{ .lastResponse._id }}/approval
    Then the response code should be 404

  Scenario: should be possible for any user to get waiting for approval created instruction
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-13-name",
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
              "description": "test-instruction-to-create-with-approval-13-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-13-step-1-endpoint"
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
    When I do GET /api/v4/cat/instructions/{{ .instructionID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "{{ .instructionID }}",
      "type": 0,
      "status": 1,
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
              "description": "test-instruction-to-create-with-approval-13-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-13-step-1-endpoint"
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
    When I am manager
    When I do GET /api/v4/cat/instructions/{{ .instructionID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "{{ .instructionID }}",
      "type": 0,
      "status": 1,
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
              "description": "test-instruction-to-create-with-approval-13-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-13-step-1-endpoint"
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
    When I am role-to-instruction-approve-1
    When I do GET /api/v4/cat/instructions/{{ .instructionID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "{{ .instructionID }}",
      "type": 0,
      "status": 1,
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
              "description": "test-instruction-to-create-with-approval-13-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-13-step-1-endpoint"
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

  Scenario: the users that didn't request the approval can update only name/description/enabled
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-14-name",
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
              "description": "test-instruction-to-create-with-approval-14-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-14-step-1-endpoint"
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
      "name": "test-instruction-to-create-with-approval-14-name-changed",
      "description": "test-instruction-to-create-with-approval-14-description-changed",
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
      "name": "test-instruction-to-create-with-approval-14-name-changed",
      "description": "test-instruction-to-create-with-approval-14-description-changed",
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
          "name": "test-instruction-to-create-with-approval-14-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-14-step-1-operation-1",
              "time_to_complete": {
                  "value": 1,
                  "unit": "s"
              },
              "description": "test-instruction-to-create-with-approval-14-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-14-step-1-endpoint"
        }
      ]
    }
    """

  Scenario: the user that requested the approval can update any field
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-15-name",
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
              "description": "test-instruction-to-create-with-approval-15-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-15-step-1-endpoint"
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
      "name": "test-instruction-to-create-with-approval-15-name-updated",
      "description": "test-instruction-to-create-with-approval-15-description-updated",
      "enabled": false,
      "timeout_after_execution": {
        "value": 11,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-15-step-1-updated",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-15-step-1-operation-1-updated",
              "time_to_complete": {"value": 5, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-15-step-1-operation-1-description-updated"
            }
          ],
          "stop_on_fail": false,
          "endpoint": "test-instruction-to-create-with-approval-15-step-1-endpoint-updated"
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
      "name": "test-instruction-to-create-with-approval-15-name-updated",
      "description": "test-instruction-to-create-with-approval-15-description-updated",
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
          "name": "test-instruction-to-create-with-approval-15-step-1-updated",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-15-step-1-operation-1-updated",
              "time_to_complete": {
                  "value": 5,
                  "unit": "s"
              },
              "description": "test-instruction-to-create-with-approval-15-step-1-operation-1-description-updated"
            }
          ],
          "stop_on_fail": false,
          "endpoint": "test-instruction-to-create-with-approval-15-step-1-endpoint-updated"
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

  Scenario: the users that didn't request the approval can't change approver
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-17-name",
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
              "description": "test-instruction-to-create-with-approval-17-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-17-step-1-endpoint"
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
      "name": "test-instruction-to-create-with-approval-17-name",
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
              "description": "test-instruction-to-create-with-approval-17-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-17-step-1-endpoint"
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

  Scenario: the user that requested the approval can change approver username
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-18-name",
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
              "description": "test-instruction-to-create-with-approval-18-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-18-step-1-endpoint"
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
      "name": "test-instruction-to-create-with-approval-18-name",
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
              "description": "test-instruction-to-create-with-approval-18-step-1-operation-1-description"
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

  Scenario: the user that requested the approval can change approval from username to role
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-19-name",
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
              "description": "test-instruction-to-create-with-approval-19-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-19-step-1-endpoint"
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
      "name": "test-instruction-to-create-with-approval-19-name",
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
              "description": "test-instruction-to-create-with-approval-19-step-1-operation-1-description"
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

  Scenario: the user that requested the approval can change approval from role to username
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-20-name",
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
              "description": "test-instruction-to-create-with-approval-20-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-20-step-1-endpoint"
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
      "name": "test-instruction-to-create-with-approval-20-name",
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
              "description": "test-instruction-to-create-with-approval-20-step-1-operation-1-description"
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

  Scenario: the user that didn't request the approval can't cancel approval
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-21-name",
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
              "description": "test-instruction-to-create-with-approval-21-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-21-step-1-endpoint"
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
      "name": "test-instruction-to-create-with-approval-21-name",
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
              "description": "test-instruction-to-create-with-approval-21-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-21-step-1-endpoint"
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

  Scenario: the user that request the approval can cancel approval
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-22-name",
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
              "description": "test-instruction-to-create-with-approval-22-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-22-step-1-endpoint"
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
      "name": "test-instruction-to-create-with-approval-22-name",
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
              "description": "test-instruction-to-create-with-approval-22-step-1-operation-1-description"
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

  Scenario: only the user from approval should be able to approve by username
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-5-name",
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
              "description": "test-instruction-to-create-with-approval-5-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-5-step-1-endpoint"
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
              "description": "test-instruction-to-create-with-approval-5-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-5-step-1-endpoint"
        }
      ]
    }
    """
    Then the response key "approval" should not exist

  Scenario: only the user from approval should be able to approve by role
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-6-name",
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
              "description": "test-instruction-to-create-with-approval-6-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-6-step-1-endpoint"
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
              "description": "test-instruction-to-create-with-approval-6-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-6-step-1-endpoint"
        }
      ]
    }
    """
    Then the response key "approval" should not exist

  Scenario: only the user from approval should be able to dismiss by username
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-7-name",
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
              "description": "test-instruction-to-create-with-approval-7-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-7-step-1-endpoint"
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
    When I am role-to-instruction-approve-2
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}/approval:
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
    When I am role-to-instruction-approve-1
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}/approval:
    """json
    {
      "approve": false
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/instructions/{{ .instructionID }}
    Then the response code should be 404

  Scenario: only the user from approval should be able to dismiss by role
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-8-name",
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
              "description": "test-instruction-to-create-with-approval-8-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-8-step-1-endpoint"
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
    When I am role-to-instruction-approve-1
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}/approval:
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
    When I am role-to-instruction-approve-2
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}/approval:
    """json
    {
      "approve": false
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/instructions/{{ .instructionID }}
    Then the response code should be 404

  Scenario: given create request with approval request for auto instruction with user or role should return ok
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "name": "test-instruction-to-create-with-approval-16-name",
      "description": "test-instruction-to-create-with-approval-16-description",
      "enabled": true,
      "priority": 100,
      "triggers": ["create"],
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
          "_id": "{{ .instructionID }}",
          "type": 1,
          "status": 1,
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
          "requested_by": "root"
        }
      },
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
