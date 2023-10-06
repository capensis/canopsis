Feature: instruction approval creation
  I need to be able to create an instruction with approval

  @concurrent
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

  @concurrent
  Scenario: given create request with approval request with role should return ok
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

  @concurrent
  Scenario: the approver should be able to get approval data by username
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-3-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-create-with-approval-3-pattern"
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
              "value": "test-instruction-to-create-with-approval-3-pattern"
            }
          }
        ]
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
              "description": "test-instruction-to-create-with-approval-3-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-3-step-1-endpoint"
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
                "value": "test-instruction-to-create-with-approval-3-pattern"
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
                "value": "test-instruction-to-create-with-approval-3-pattern"
              }
            }
          ]
        ],
        "name": "test-instruction-to-create-with-approval-3-name",
        "description": "test-instruction-to-create-with-approval-3-description",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-to-create-with-approval-3-step-1",
            "operations": [
              {
                "name": "test-instruction-to-create-with-approval-3-step-1-operation-1",
                "time_to_complete": {
                    "value": 1,
                    "unit": "s"
                },
                "description": "test-instruction-to-create-with-approval-3-step-1-operation-1-description"
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-to-create-with-approval-3-step-1-endpoint"
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

  @concurrent
  Scenario: the approver should be able to get approval data by role
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
        "role": {
          "_id": "role-to-instruction-approve-1",
          "name": "role-to-instruction-approve-1"
        },
        "requested_by": "root"
      }
    }
    """

  @concurrent
  Scenario: should be possible for any user to get waiting for approval created instruction
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
    When I am admin
    When I do GET /api/v4/cat/instructions/{{ .instructionID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "{{ .instructionID }}",
      "type": 0,
      "status": 1,
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
  Scenario: the users that didn't request the approval can update only name/description/enabled
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
      "name": "test-instruction-to-create-with-approval-6-name-changed",
      "description": "test-instruction-to-create-with-approval-6-description-changed",
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
      "name": "test-instruction-to-create-with-approval-6-name-changed",
      "description": "test-instruction-to-create-with-approval-6-description-changed",
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

  @concurrent
  Scenario: the user that requested the approval can update any field
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
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-with-approval-7-name-updated",
      "description": "test-instruction-to-create-with-approval-7-description-updated",
      "enabled": false,
      "timeout_after_execution": {
        "value": 11,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-with-approval-7-step-1-updated",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-7-step-1-operation-1-updated",
              "time_to_complete": {"value": 5, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-7-step-1-operation-1-description-updated"
            }
          ],
          "stop_on_fail": false,
          "endpoint": "test-instruction-to-create-with-approval-7-step-1-endpoint-updated"
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
      "name": "test-instruction-to-create-with-approval-7-name-updated",
      "description": "test-instruction-to-create-with-approval-7-description-updated",
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
          "name": "test-instruction-to-create-with-approval-7-step-1-updated",
          "operations": [
            {
              "name": "test-instruction-to-create-with-approval-7-step-1-operation-1-updated",
              "time_to_complete": {
                  "value": 5,
                  "unit": "s"
              },
              "description": "test-instruction-to-create-with-approval-7-step-1-operation-1-description-updated"
            }
          ],
          "stop_on_fail": false,
          "endpoint": "test-instruction-to-create-with-approval-7-step-1-endpoint-updated"
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
      "name": "test-instruction-to-create-with-approval-9-name",
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
              "description": "test-instruction-to-create-with-approval-9-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-with-approval-9-step-1-endpoint"
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
      "name": "test-instruction-to-create-with-approval-9-name",
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
              "time_to_complete": {"value": 5, "unit":"s"},
              "description": "test-instruction-to-create-with-approval-9-step-1-operation-1-description"
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
