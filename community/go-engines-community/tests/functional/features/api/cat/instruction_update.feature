Feature: Instruction update

  Scenario: PUT as unauthorized
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update:
    """
    {
      "name": "test-instruction-to-update-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-description",
      "author": "test-instruction-to-update-author",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-update-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-step-1-operation-1-name",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-to-update-step-1-operation-1-description",
              "jobs": [
                "test-job-to-link-3"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 401

  Scenario: PUT without permissions
    When I am noperms
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update:
    """
    {
      "name": "test-instruction-to-update-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-description",
      "author": "test-instruction-to-update-author",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-update-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-step-1-operation-1-name",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-to-update-step-1-operation-1-description",
              "jobs": [
                "test-job-to-link-3"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 403

  Scenario: PUT a valid instruction without any changes
    When I am admin
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update:
    """
    {
      "name": "test-instruction-to-update-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-description",
      "author": "test-instruction-to-update-author",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-update-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-step-1-operation-1-name",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-to-update-step-1-operation-1-description",
              "jobs": [
                "test-job-to-link-3"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-instruction-to-update",
      "name": "test-instruction-to-update-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-description",
      "author": "test-instruction-to-update-author",
      "enabled": true,
      "rating": 0,
      "avg_complete_time": 10,
      "month_executions": 0,
      "last_executed_by": {
        "_id": "root"
      },
      "last_executed_on": 1596712203,
      "steps": [
        {
          "name": "test-instruction-to-update-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-step-1-operation-1-name",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-to-update-step-1-operation-1-description",
              "jobs": [
                {
                  "_id": "test-job-to-link-3",
                  "name": "test-job-name-to-link-3",
                  "author": "test-author",
                  "config": {
                    "_id": "test-job-config-to-link",
                    "name": "test-job-config-name-to-link",
                    "type": "rundeck",
                    "host": "http://example.com",
                    "auth_token": "test-auth-token"
                  },
                  "job_id": "test-job-id",
                  "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                }
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-step-1-endpoint"
        }
      ]
    }
    """

  Scenario: PUT a valid instruction
    When I am admin
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update:
    """
    {
      "name": "test-instruction-to-update-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-description-changed",
      "author": "test-instruction-to-update-author",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-update-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-step-1-operation-1-name",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-to-update-step-1-operation-1-description",
              "jobs": [
                "test-job-to-link-3",
                "test-job-to-link-2"
              ]
            },
            {
              "name": "test-instruction-to-update-step-1-operation-2-name",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-to-update-step-1-operation-2-description",
              "jobs": [
                "test-job-to-link-3"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-step-1-endpoint"
        },
        {
          "name": "test-instruction-to-update-step-2-name",
          "operations": [
            {
              "name": "test-instruction-to-update-step-2-operation-1-name",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-to-update-step-2-operation-1-description",
              "jobs": [
                "test-job-to-link-3",
                "test-job-to-link-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-instruction-to-update",
      "name": "test-instruction-to-update-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-description-changed",
      "author": "test-instruction-to-update-author",
      "enabled": true,
      "rating": 0,
      "avg_complete_time": 10,
      "month_executions": 0,
      "steps": [
        {
          "name": "test-instruction-to-update-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-step-1-operation-1-name",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-to-update-step-1-operation-1-description",
              "jobs": [
                {
                  "_id": "test-job-to-link-3",
                  "name": "test-job-name-to-link-3",
                  "author": "test-author",
                  "config": {
                    "_id": "test-job-config-to-link",
                    "name": "test-job-config-name-to-link",
                    "type": "rundeck",
                    "host": "http://example.com",
                    "auth_token": "test-auth-token"
                  },
                  "job_id": "test-job-id",
                  "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                },
                {
                  "_id": "test-job-to-link-2",
                  "name": "test-job-name-to-link-2",
                  "author": "test-author",
                  "config": {
                    "_id": "test-job-config-to-link",
                    "name": "test-job-config-name-to-link",
                    "type": "rundeck",
                    "host": "http://example.com",
                    "auth_token": "test-auth-token"
                  },
                  "job_id": "test-job-id",
                  "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                }
              ]
            },
            {
              "name": "test-instruction-to-update-step-1-operation-2-name",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-to-update-step-1-operation-2-description",
              "jobs": [
                {
                  "_id": "test-job-to-link-3",
                  "name": "test-job-name-to-link-3",
                  "author": "test-author",
                  "config": {
                    "_id": "test-job-config-to-link",
                    "name": "test-job-config-name-to-link",
                    "type": "rundeck",
                    "host": "http://example.com",
                    "auth_token": "test-auth-token"
                  },
                  "job_id": "test-job-id",
                  "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                }
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-step-1-endpoint"
        },
        {
          "name": "test-instruction-to-update-step-2-name",
          "operations": [
            {
              "name": "test-instruction-to-update-step-2-operation-1-name",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-to-update-step-2-operation-1-description",
              "jobs": [
                {
                  "_id": "test-job-to-link-3",
                  "name": "test-job-name-to-link-3",
                  "author": "test-author",
                  "config": {
                    "_id": "test-job-config-to-link",
                    "name": "test-job-config-name-to-link",
                    "type": "rundeck",
                    "host": "http://example.com",
                    "auth_token": "test-auth-token"
                  },
                  "job_id": "test-job-id",
                  "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                },
                {
                  "_id": "test-job-to-link-1",
                  "name": "test-job-name-to-link-1",
                  "author": "test-author",
                  "config": {
                    "_id": "test-job-config-to-link",
                    "name": "test-job-config-name-to-link",
                    "type": "rundeck",
                    "host": "http://example.com",
                    "auth_token": "test-auth-token"
                  },
                  "job_id": "test-job-id",
                  "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                }
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-step-2-endpoint"
        }
      ],
      "last_executed_by": {
        "_id": "root"
      },
      "last_executed_on": 1596712203
    }
    """

  Scenario: PUT an invalid instruction, where name already exists
    When I am admin
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update:
    """
    {
      "name": "test-instruction-to-get-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-description-changed",
      "author": "test-instruction-to-update-author",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-update-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-step-1-operation-1-name",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-to-update-step-1-operation-1-description",
              "jobs": [
                "test-job-to-link-3",
                "test-job-to-link-2"
              ]
            },
            {
              "name": "test-instruction-to-update-step-1-operation-2-name",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-to-update-step-1-operation-2-description",
              "jobs": [
                "test-job-to-link-3"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-step-1-endpoint"
        },
        {
          "name": "test-instruction-to-update-step-2-name",
          "operations": [
            {
              "name": "test-instruction-to-update-step-2-operation-1-name",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-to-update-step-2-operation-1-description",
              "jobs": [
                "test-job-to-link-3",
                "test-job-to-link-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "name": "Name already exists"
      }
    }
    """

  Scenario: PUT an invalid instruction, where name already exists
    When I am admin
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update:
    """
    {
      "name": "test-instruction-to-update-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-description-changed",
      "author": "test-instruction-to-update-author",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-to-update-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-step-1-operation-1-name",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-to-update-step-1-operation-1-description",
              "jobs": [
                "test-job-to-link-3",
                "test-job-to-link-2",
                "test-job-to-link-2-NOT-EXIST"
              ]
            },
            {
              "name": "test-instruction-to-update-step-1-operation-2-name",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-to-update-step-1-operation-2-description",
              "jobs": [
                "test-job-to-link-3"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-step-1-endpoint"
        },
        {
          "name": "test-instruction-to-update-step-2-name",
          "operations": [
            {
              "name": "test-instruction-to-update-step-2-operation-1-name",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-to-update-step-2-operation-1-description",
              "jobs": [
                "test-job-to-link-3",
                "test-job-to-link-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "error": "job doesn't exist"
    }
    """

  Scenario: PUT a valid instruction with pbehavior types
    When I am admin
    When I do PUT /api/v4/cat/instructions/test-instruction-with-pbh-to-edit:
    """
    {
      "name": "test-instruction-with-pbh-to-edit-name",
      "description": "test-instruction-with-pbh-to-edit-description",
      "author": "test-instruction-with-pbh-to-edit-author",
      "enabled": true,
      "alarm_patterns": [
        {
          "v": {
            "component": "test-alarm-instruction-with-pbehavior-component-no-match"
          }
        }
      ],
      "steps": [],
      "active_on_pbh": [
        "pbh-type-for-instruction-with-pbehavior-1"
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "name": "test-instruction-with-pbh-to-edit-name",
      "description": "test-instruction-with-pbh-to-edit-description",
      "author": "test-instruction-with-pbh-to-edit-author",
      "enabled": true,
      "alarm_patterns": [
        {
          "v": {
            "component": "test-alarm-instruction-with-pbehavior-component-no-match"
          }
        }
      ],
      "active_on_pbh": [
        "pbh-type-for-instruction-with-pbehavior-1"
      ]
    }
    """

  Scenario: PUT an invalid instruction with pbehavior types, that don't exist
    When I am admin
    When I do PUT /api/v4/cat/instructions/test-instruction-with-pbh-to-edit:
    """
    {
      "name": "test-instruction-with-pbh-to-edit-name",
      "description": "test-instruction-with-pbh-to-edit-description",
      "author": "test-instruction-with-pbh-to-edit-author",
      "enabled": true,
      "alarm_patterns": [
        {
          "v": {
            "component": "test-alarm-instruction-with-pbehavior-component-no-match"
          }
        }
      ],
      "steps": [],
      "active_on_pbh": [
        "pbh-type-for-instruction-with-pbehavior-1",
        "not-exist"
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
          "active_on_pbh": "active_on_pbh doesn't exist"
      }
    }
    """