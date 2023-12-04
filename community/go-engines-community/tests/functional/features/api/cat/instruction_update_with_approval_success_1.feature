Feature: instruction approval update
  I need to be able to update an instruction with approval

  @concurrent
  Scenario: given update request with user approval should return ok
    When I am manager
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-1:
    """json
    {
      "name": "test-instruction-to-update-with-approval-1-name-updated",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-update-with-approval-1-pattern-updated"
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
              "value": "test-instruction-to-update-with-approval-1-pattern-updated"
            }
          }
        ]
      ],
      "description": "test-instruction-to-update-with-approval-1-description-updated",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-1-step-1-name-updated",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-1-step-1-operation-1-name-updated",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-1-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-with-approval-1-step-1-endpoint-updated"
        }
      ],
      "approval": {
        "user": "user-to-instruction-approve-1",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-1/approval
    Then the response code should be 403
    Then the response body should be:
    """json
    {
      "error": "user is not assigned to approval"
    }
    """
    When I am role-to-instruction-approve-1
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-1/approval
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
        "_id": "test-instruction-to-update-with-approval-1",
        "type": 0,
        "status": 0,
        "created": 1596712203,
        "last_modified": 1596712203,
        "entity_pattern": [
          [
            {
              "field": "name",
              "cond": {
                "type": "eq",
                "value": "test-instruction-to-update-with-approval-1-pattern"
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
                "value": "test-instruction-to-update-with-approval-1-pattern"
              }
            }
          ]
        ],
        "name": "test-instruction-to-update-with-approval-1-name",
        "description": "test-instruction-to-update-with-approval-1-description",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-to-update-with-approval-1-step-1-name",
            "operations": [
              {
                "name": "test-instruction-to-update-with-approval-1-step-1-operation-1-name",
                "time_to_complete": {
                  "value": 1,
                  "unit": "s"
                },
                "description": "test-instruction-to-update-with-approval-1-step-1-operation-1-description"
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-to-update-with-approval-1-step-1-endpoint"
          }
        ]
      },
      "updated": {
        "type": 0,
        "status": 2,
        "name": "test-instruction-to-update-with-approval-1-name-updated",
        "entity_pattern": [
          [
            {
              "field": "name",
              "cond": {
                "type": "eq",
                "value": "test-instruction-to-update-with-approval-1-pattern-updated"
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
                "value": "test-instruction-to-update-with-approval-1-pattern-updated"
              }
            }
          ]
        ],
        "description": "test-instruction-to-update-with-approval-1-description-updated",
        "author": {
          "_id": "manageruser",
          "name": "manageruser"
        },
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-to-update-with-approval-1-step-1-name-updated",
            "operations": [
              {
                "name": "test-instruction-to-update-with-approval-1-step-1-operation-1-name-updated",
                "time_to_complete": {"value": 1, "unit":"s"},
                "description": "test-instruction-to-update-with-approval-1-step-1-operation-1-description",
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
                      "auth_token": "test-auth-token"
                    },
                    "job_id": "test-job-to-instruction-edit-1-external-id",
                    "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                  }
                ]
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-to-update-with-approval-1-step-1-endpoint-updated"
          }
        ],
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
        }
      }
    }
    """
    When I do GET /api/v4/cat/instruction-stats?search=test-instruction-to-update-with-approval-1-name&from=1000000000&to=2000000000
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
    When I do GET /api/v4/cat/instructions?search=test-instruction-to-update-with-approval-1-name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-instruction-to-update-with-approval-1",
          "status": 1,
          "name": "test-instruction-to-update-with-approval-1-name",
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
  Scenario: given update request with role approval should return ok
    When I am manager
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-2:
    """json
    {
      "name": "test-instruction-to-update-with-approval-2-name-updated",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-update-with-approval-2-pattern-updated"
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
              "value": "test-instruction-to-update-with-approval-2-pattern-updated"
            }
          }
        ]
      ],
      "description": "test-instruction-to-update-with-approval-2-description-updated",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-2-step-1-name-updated",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-2-step-1-operation-1-name-updated",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-with-approval-2-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-with-approval-2-step-1-endpoint-updated"
        }
      ],
      "approval": {
        "role": "role-to-instruction-approve-2",
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
    When I am role-to-instruction-approve-2
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-2/approval
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
          "_id": "manageruser",
          "name": "manageruser"
        }
      },
      "original": {
        "_id": "test-instruction-to-update-with-approval-2",
        "type": 0,
        "status": 0,
        "entity_pattern": [
          [
            {
              "field": "name",
              "cond": {
                "type": "eq",
                "value": "test-instruction-to-update-with-approval-2-pattern"
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
                "value": "test-instruction-to-update-with-approval-2-pattern"
              }
            }
          ]
        ],
        "name": "test-instruction-to-update-with-approval-2-name",
        "description": "test-instruction-to-update-with-approval-2-description",
        "author": {
          "_id": "root",
          "name": "root"
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
        "name": "test-instruction-to-update-with-approval-2-name-updated",
        "entity_pattern": [
          [
            {
              "field": "name",
              "cond": {
                "type": "eq",
                "value": "test-instruction-to-update-with-approval-2-pattern-updated"
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
                "value": "test-instruction-to-update-with-approval-2-pattern-updated"
              }
            }
          ]
        ],
        "description": "test-instruction-to-update-with-approval-2-description-updated",
        "author": {
          "_id": "manageruser",
          "name": "manageruser"
        },
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-to-update-with-approval-2-step-1-name-updated",
            "operations": [
              {
                "name": "test-instruction-to-update-with-approval-2-step-1-operation-1-name-updated",
                "time_to_complete": {"value": 1, "unit":"s"},
                "description": "test-instruction-to-update-with-approval-2-step-1-operation-1-description",
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
                      "auth_token": "test-auth-token"
                    },
                    "job_id": "test-job-to-instruction-edit-1-external-id",
                    "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                  }
                ]
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-to-update-with-approval-2-step-1-endpoint-updated"
          }
        ],
        "approval": {
          "comment": "test comment",
          "role": {
            "_id": "role-to-instruction-approve-2",
            "name": "role-to-instruction-approve-2"
          },
          "requested_by": {
            "_id": "manageruser",
            "name": "manageruser"
          }
        }
      }
    }
    """
    When I am admin
    When I do GET /api/v4/cat/instruction-stats?search=test-instruction-to-update-with-approval-2&from=1000000000&to=2000000000
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
    When I do GET /api/v4/cat/instructions?search=test-instruction-to-update-with-approval-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-instruction-to-update-with-approval-2",
          "name": "test-instruction-to-update-with-approval-2-name",
          "status": 1,
          "approval": {
            "comment": "test comment",
            "role": {
              "_id": "role-to-instruction-approve-2",
              "name": "role-to-instruction-approve-2"
            },
            "requested_by": {
              "_id": "manageruser",
              "name": "manageruser"
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

  @concurrent
  Scenario: Requester should receive updated version on GET request, other users should receive original version on GET request
    When I am manager
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-3:
    """json
    {
      "name": "test-instruction-to-update-with-approval-3-name",
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
          "endpoint": "new endpoint"
        }
      ],
      "approval": {
        "user": "user-to-instruction-approve-1",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200
    When I am manager
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-update-with-approval-3",
      "type": 0,
      "status": 2,
      "name": "test-instruction-to-update-with-approval-3-name",
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
                    "_id": "root",
                    "name": "root"
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
          "_id": "user-to-instruction-approve-1",
          "name": "user-to-instruction-approve-1"
        },
        "requested_by": {
          "_id": "manageruser",
          "name": "manageruser"
        }
      }
    }
    """
    When I am admin
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-update-with-approval-3",
      "type": 0,
      "status": 0,
      "name": "test-instruction-to-update-with-approval-3-name",
      "description": "test-instruction-to-update-with-approval-3-description",
      "created": 1596712203,
      "last_modified": 1596712203,
      "author": {
        "_id": "root",
        "name": "root"
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

  @concurrent
  Scenario: The users that didn't request the approval can update only name/description/enabled
    When I am manager
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-4:
    """json
    {
      "name": "test-instruction-to-update-with-approval-4-name",
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
        "user": "user-to-instruction-approve-1",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200
    When I am admin
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-4:
    """json
    {
      "name": "test-instruction-to-update-with-approval-4-name-changed",
      "description": "test-instruction-to-update-with-approval-4-description-changed",
      "enabled": false,
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
    When I am role-to-instruction-approve-1
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-4/approval
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
        "_id": "test-instruction-to-update-with-approval-4",
        "type": 0,
        "status": 0,
        "created": 1596712203,
        "name": "test-instruction-to-update-with-approval-4-name-changed",
        "description": "test-instruction-to-update-with-approval-4-description-changed",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "enabled": false,
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
      },
      "updated": {
        "type": 0,
        "status": 2,
        "name": "test-instruction-to-update-with-approval-4-name",
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
                      "_id": "root",
                      "name": "root"
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
            "_id": "user-to-instruction-approve-1",
            "name": "user-to-instruction-approve-1"
          },
          "requested_by": {
            "_id": "manageruser",
            "name": "manageruser"
          }
        }
      }
    }
    """
    When I do GET /api/v4/cat/instruction-stats?search=test-instruction-to-update-with-approval-4&from=1000000000&to=2000000000
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
  Scenario: The requester can update any updated fields
    When I am manager
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-5:
    """json
    {
      "name": "test-instruction-to-update-with-approval-5-name",
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
        "user": "user-to-instruction-approve-1",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-5:
    """json
    {
      "name": "test-instruction-to-update-with-approval-5-name",
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
          "endpoint": "shouldn't be ignored in updated"
        }
      ],
      "approval": {
        "user": "user-to-instruction-approve-1",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200
    When I am role-to-instruction-approve-1
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-5/approval
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
        "_id": "test-instruction-to-update-with-approval-5",
        "type": 0,
        "status": 0,
        "name": "test-instruction-to-update-with-approval-5-name",
        "description": "test-instruction-to-update-with-approval-5-description",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "enabled": true,
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
        "description": "test-instruction-to-update-with-approval-5-description",
        "author":{
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
            "_id": "user-to-instruction-approve-1",
            "name": "user-to-instruction-approve-1"
          },
          "requested_by": {
            "_id": "manageruser",
            "name": "manageruser"
          }
        }
      }
    }
    """

  @concurrent
  Scenario: The users that didn't request the approval couldn't change or remove the approval
    When I am manager
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-6:
    """json
    {
      "name": "test-instruction-to-update-with-approval-6-name",
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
        "user": "user-to-instruction-approve-1",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200
    When I am admin
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-6:
    """json
    {
      "name": "test-instruction-to-update-with-approval-6-name",
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
        "role": "role-to-instruction-approve-1"
      }
    }
    """
    Then the response code should be 200
    When I am role-to-instruction-approve-1
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-6/approval
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
      "updated": {
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
        }
      }
    }
    """
    When I am admin
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-6:
    """json
    {
      "name": "test-instruction-to-update-with-approval-6-name",
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
      ]
    }
    """
    Then the response code should be 200
    When I am role-to-instruction-approve-1
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-6/approval
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
      "updated": {
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
        }
      }
    }
    """

  @concurrent
  Scenario: The requester can update or remove the approval, after removal instruction should be updated
    When I am manager
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-7:
    """json
    {
      "name": "test-instruction-to-update-with-approval-7-name",
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
        "user": "user-to-instruction-approve-1",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200
    When I am manager
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-7:
    """json
    {
      "name": "test-instruction-to-update-with-approval-7-name",
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
        "role": "role-to-instruction-approve-2",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200
    When I am role-to-instruction-approve-1
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-7/approval
    Then the response code should be 403
    When I am role-to-instruction-approve-2
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-7/approval
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
          "_id": "manageruser",
          "name": "manageruser"
        }
      },
      "updated": {
        "approval": {
          "comment": "test comment",
          "role": {
            "_id": "role-to-instruction-approve-2",
            "name": "role-to-instruction-approve-2"
          },
          "requested_by": {
            "_id": "manageruser",
            "name": "manageruser"
          }
        }
      }
    }
    """
    When I am manager
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-7:
    """json
    {
      "name": "test-instruction-to-update-with-approval-7-name",
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
      ]
    }
    """
    Then the response code should be 200
    When I am role-to-instruction-approve-2
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-7/approval
    Then the response code should be 404
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-update-with-approval-7",
      "type": 0,
      "status": 0,
      "name": "test-instruction-to-update-with-approval-7-name",
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
          "endpoint": "new endpoint"
        }
      ]
    }
    """

  @concurrent
  Scenario: Only the user from approval should be able to approve
    When I am manager
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-8:
    """json
    {
      "name": "test-instruction-to-update-with-approval-8-name",
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
          "endpoint": "new endpoint 2"
        }
      ],
      "approval": {
        "role": "role-to-instruction-approve-2",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200
    When I am admin
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-8/approval:
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
    When I am manager
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-8/approval:
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
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-8/approval:
    """json
    {
      "approve": true
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-8/approval
    Then the response code should be 404
    When I do GET /api/v4/cat/instructions/test-instruction-to-update-with-approval-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "type": 0,
      "status": 0,
      "name": "test-instruction-to-update-with-approval-8-name",
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
                    "_id": "root",
                    "name": "root"
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
    When I do GET /api/v4/cat/instruction-stats?search=test-instruction-to-update-with-approval-8&from=1000000000&to=2000000000
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
  Scenario: Only the user from approval should be able to dismiss
    When I am manager
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-9:
    """json
    {
      "name": "test-instruction-to-update-with-approval-9-name",
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
          "endpoint": "new endpoint 3"
        }
      ],
      "approval": {
        "user": "user-to-instruction-approve-2",
        "comment": "test comment"
      }
    }
    """
    Then the response code should be 200
    When I am manager
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-9/approval:
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
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-with-approval-9/approval:
    """json
    {
      "approve": false
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
      "_id": "test-instruction-to-update-with-approval-9",
      "type": 0,
      "status": 0,
      "name": "test-instruction-to-update-with-approval-9-name",
      "description": "test-instruction-to-update-with-approval-9-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-update-with-approval-9-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-with-approval-9-step-1-operation-1-name",
              "time_to_complete": {
                "value": 1,
                "unit": "s"
              },
              "description": "test-instruction-to-update-with-approval-9-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-with-approval-9-step-1-endpoint"
        }
      ]
    }
    """
    When I do GET /api/v4/cat/instruction-stats?search=test-instruction-to-update-with-approval-9&from=1000000000&to=2000000000
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
