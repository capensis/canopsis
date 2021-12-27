Feature: create a instruction
  I need to be able to create a instruction
  Only admin should be able to create a instruction

  Scenario: given create manual instruction request should return ok
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-1-name",
      "entity_patterns": [
        {"name": "test-instruction-to-create-1-pattern"
        }
      ],
      "description": "test-instruction-to-create-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-1-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-1-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            },
            {
              "name": "test-instruction-to-create-1-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-to-create-1-step-1-operation-2-description",
              "jobs": [
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-1-step-1-endpoint"
        },
        {
          "name": "test-instruction-to-create-1-step-2",
          "operations": [
            {
              "name": "test-instruction-to-create-1-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-to-create-1-step-2-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-1-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "type": 0,
      "status": 0,
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-1-pattern"
        }
      ],
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test-instruction-to-create-1-description",
      "enabled": true,
      "name": "test-instruction-to-create-1-name",
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-1-step-1-endpoint",
          "name": "test-instruction-to-create-1-step-1",
          "operations": [
            {
              "description": "test-instruction-to-create-1-step-1-operation-1-description",
              "jobs": [
                {
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
                },
                {
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
              ],
              "name": "test-instruction-to-create-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"}
            },
            {
              "description": "test-instruction-to-create-1-step-1-operation-2-description",
              "jobs": [
                {
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
              ],
              "name": "test-instruction-to-create-1-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"}
            }
          ]
        },
        {
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-1-step-2-endpoint",
          "name": "test-instruction-to-create-1-step-2",
          "operations": [
            {
              "description": "test-instruction-to-create-1-step-2-operation-1-description",
              "jobs": [
                {
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
                },
                {
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
              ],
              "name": "test-instruction-to-create-1-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"}
            }
          ]
        }
      ]
    }
    """

  Scenario: given create auto instruction request should return ok
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "name": "test-instruction-to-create-2-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-2-pattern"
        }
      ],
      "description": "test-instruction-to-create-2-description",
      "enabled": true,
      "priority": 21,
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
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-2-pattern"
        }
      ],
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test-instruction-to-create-2-description",
      "enabled": true,
      "priority": 21,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "name": "test-instruction-to-create-2-name",
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
      ]
    }
    """

  Scenario: given create request should return ok to get request
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-3-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-3-pattern"
        }
      ],
      "description": "test-instruction-to-create-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-3-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-3-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-3-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            },
            {
              "name": "test-instruction-to-create-3-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-to-create-3-step-1-operation-2-description",
              "jobs": [
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-3-step-1-endpoint"
        },
        {
          "name": "test-instruction-to-create-3-step-2",
          "operations": [
            {
              "name": "test-instruction-to-create-3-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-to-create-3-step-2-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-3-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "type": 0,
      "status": 0,
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-3-pattern"
        }
      ],
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test-instruction-to-create-3-description",
      "enabled": true,
      "name": "test-instruction-to-create-3-name",
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-3-step-1-endpoint",
          "name": "test-instruction-to-create-3-step-1",
          "operations": [
            {
              "description": "test-instruction-to-create-3-step-1-operation-1-description",
              "jobs": [
                {
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
                },
                {
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
              ],
              "name": "test-instruction-to-create-3-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"}
            },
            {
              "description": "test-instruction-to-create-3-step-1-operation-2-description",
              "jobs": [
                {
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
              ],
              "name": "test-instruction-to-create-3-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"}
            }
          ]
        },
        {
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-3-step-2-endpoint",
          "name": "test-instruction-to-create-3-step-2",
          "operations": [
            {
              "description": "test-instruction-to-create-3-step-2-operation-1-description",
              "jobs": [
                {
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
                },
                {
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
              ],
              "name": "test-instruction-to-create-3-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"}
            }
          ]
        }
      ]
    }
    """

  Scenario: given create request with custom id should return ok
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-4-id",
      "type": 0,
      "name": "test-instruction-to-create-4-name",
      "entity_patterns": [
        {
          "name": "test-instruction-to-create-4-pattern"
        }
      ],
      "description": "test-instruction-to-create-4-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-create-4-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-4-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-4-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-4-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-create-4-id"
    }
    """
    When I do GET /api/v4/cat/instructions/test-instruction-to-create-4-id
    Then the response code should be 200

  Scenario: given create request with custom id that already exist should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-check-unique-id"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "_id": "ID already exists."
      }
    }
    """

  Scenario: given create request with name that already exist should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "name": "test-instruction-to-check-unique-name-name"
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

  Scenario: given invalid create request to create manual instruction should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "description": "Description is missing.",
        "enabled": "Enabled is missing.",
        "name": "Name is missing.",
        "steps": "Steps is missing.",
        "timeout_after_execution": "TimeoutAfterExecution is missing."
      }
    }
    """

  Scenario: given invalid create request to create auto instruction should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "description": "Description is missing.",
        "enabled": "Enabled is missing.",
        "jobs": "Jobs is missing.",
        "name": "Name is missing.",
        "priority": "Priority is missing.",
        "timeout_after_execution": "TimeoutAfterExecution is missing."
      }
    }
    """

  Scenario: given create request with not exist job should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "steps": [
        {
          "operations": [
            {
              "jobs": [
                "test-job-not-exist"
              ]
            }
          ]
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "steps.0.operations.0.jobs": "Jobs doesn't exist."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "jobs": [
        {
          "job": "test-job-not-exist"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "jobs.0.job": "Job doesn't exist."
      }
    }
    """

  Scenario: given create request with invalid patterns should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """
      {
        "entity_patterns": [
          {
            "name": {
              "regex_match": "name:.*"
            }
          },
          {}
        ],
        "alarm_patterns": [
          {
            "v": {
              "resource": {
                "regex_match": "name:.*"
              }
            }
          },
          {}
        ]
      }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "entity_patterns": "entity pattern list contains an empty pattern.",
        "alarm_patterns": "alarm pattern list contains an empty pattern."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """
      {
        "entity_patterns": [
          {
            "qwe": {
              "regex_match": "name:.*"
            }
          }
        ],
        "alarm_patterns": [
          {
            "qwe": {
              "resource": {
                "regex_match": "name:.*"
              }
            }
          }
        ]
      }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
      "entity_patterns": "Invalid entity pattern list.",
        "alarm_patterns": "Invalid alarm pattern list."
      }
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/cat/instructions
    Then the response code should be 401

  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/instructions
    Then the response code should be 403
