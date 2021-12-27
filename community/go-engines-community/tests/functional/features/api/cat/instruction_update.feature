Feature: Instruction update

  Scenario: PUT as unauthorized
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update
    Then the response code should be 401

  Scenario: PUT without permissions
    When I am noperms
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update
    Then the response code should be 403

  Scenario: PUT a valid instruction without any changes
    When I am admin
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update:
    """json
    {
      "name": "test-instruction-to-update-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-2"
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
    """json
    {
      "_id": "test-instruction-to-update",
      "type": 0,
      "status": 0,
      "name": "test-instruction-to-update-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-step-1-operation-1-description",
              "jobs": [
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
                    "auth_token": "test-auth-token"
                  },
                  "job_id": "test-job-to-instruction-edit-2-external-id",
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
    """json
    {
      "name": "test-instruction-to-update-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-description-changed",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            },
            {
              "name": "test-instruction-to-update-step-1-operation-2-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-step-1-operation-2-description",
              "jobs": [
                "test-job-to-instruction-edit-2"
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
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-step-2-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-2",
                "test-job-to-instruction-edit-1"
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
    """json
    {
      "_id": "test-instruction-to-update",
      "type": 0,
      "status": 0,
      "name": "test-instruction-to-update-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-to-update-description-changed",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "steps": [
        {
          "name": "test-instruction-to-update-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-step-1-operation-1-description",
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
                    "auth_token": "test-auth-token"
                  },
                  "job_id": "test-job-to-instruction-edit-2-external-id",
                  "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                }
              ]
            },
            {
              "name": "test-instruction-to-update-step-1-operation-2-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-step-1-operation-2-description",
              "jobs": [
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
                    "auth_token": "test-auth-token"
                  },
                  "job_id": "test-job-to-instruction-edit-2-external-id",
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
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-step-2-operation-1-description",
              "jobs": [
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
                    "auth_token": "test-auth-token"
                  },
                  "job_id": "test-job-to-instruction-edit-2-external-id",
                  "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                },
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
                  "author": {
                    "_id": "test-user-author-1-id",
                    "name": "test-user-author-1-username"
                  },
                  "job_id": "test-job-to-instruction-edit-1-external-id",
                  "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                }
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-step-2-endpoint"
        }
      ]
    }
    """


  Scenario: PUT an invalid instruction, where name already exists
    When I am admin
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update:
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

  Scenario: PUT an invalid instruction, where job doesn't exist
    When I am admin
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update:
    """json
    {
      "name": "test-instruction-to-update-name",
      "steps": [
        {
          "operations": [
            {
              "jobs": [
                "test-job-to-instruction-edit-2",
                "test-job-to-instruction-edit-2",
                "test-job-to-instruction-edit-2-NOT-EXIST"
              ]
            },
            {
              "jobs": [
                "test-job-to-instruction-edit-2"
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

  Scenario: PUT a valid instruction with pbehavior types
    When I am admin
    When I do PUT /api/v4/cat/instructions/test-instruction-with-pbh-to-edit:
    """json
    {
      "name": "test-instruction-with-pbh-to-edit-name",
      "description": "test-instruction-with-pbh-to-edit-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "alarm_patterns": [
        {
          "v": {
            "component": "test-alarm-instruction-with-pbehavior-component-no-match"
          }
        }
      ],
      "steps": [
        {
          "name": "test-instruction-with-pbh-to-edit-step-1-name",
          "endpoint": "test-instruction-with-pbh-to-edit-step-1-endpoint",
          "stop_on_fail": true,
          "operations": [
            {
              "name": "test-instruction-with-pbh-to-edit-step-1-operation-1-name",
              "description": "test-instruction-with-pbh-to-edit-step-1-operation-1-description",
              "time_to_complete": {
                "value": 10,
                "unit": "m"
              }
            }
          ]
        }
      ],
      "active_on_pbh": [
        "pbh-type-for-instruction-with-pbehavior-1"
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-instruction-with-pbh-to-edit-name",
      "active_on_pbh": [
        "pbh-type-for-instruction-with-pbehavior-1"
      ]
    }
    """

  Scenario: PUT an invalid instruction with pbehavior types, that don't exist
    When I am admin
    When I do PUT /api/v4/cat/instructions/test-instruction-with-pbh-to-edit:
    """json
    {
      "active_on_pbh": [
        "pbh-type-for-instruction-with-pbehavior-1",
        "not-exist"
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
          "active_on_pbh": "active_on_pbh doesn't exist."
      }
    }
    """
