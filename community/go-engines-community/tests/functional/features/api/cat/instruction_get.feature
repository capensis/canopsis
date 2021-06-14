Feature: get a instruction
  I need to be able to get a instruction
  Only admin should be able to get a instruction

  Scenario: given get all request should return instructions
    When I am admin
    When I do GET /api/v4/cat/instructions?search=test-instruction-to-get-name
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "data": [
        {
          "_id": "test-instruction-to-get",
          "name": "test-instruction-to-get-name",
          "entity_patterns": [
            {
              "name": "test filter"
            }
          ],
          "alarm_patterns": null,
          "description": "test-instruction-to-get-description",
          "author": "test-instruction-to-get-author",
          "enabled": true,
          "rating": 0,
          "comments": [],
          "avg_complete_time": 10,
          "month_executions": 0,
          "steps": [
            {
              "name": "test-instruction-to-get-step-1-name",
              "operations": [
                {
                  "name": "test-instruction-to-get-step-1-operation-1-name",
                  "time_to_complete": {"seconds": 1, "unit":"s"},
                  "description": "test-instruction-to-get-step-1-operation-1-description",
                  "jobs": [
                    {
                      "_id": "test-job-to-test-instruction-to-get-step-1-operation-1-1",
                      "name": "test-job-to-test-instruction-to-get-step-1-operation-1-1-name",
                      "author": "test_author",
                      "config": {
                        "_id": "test-job-config-to-link",
                        "name": "test-job-config-name-to-link",
                        "type": "rundeck",
                        "host": "http://example.com",
                        "author": "test-author",
                        "auth_token": "test-auth-token"
                      },
                      "job_id": "test-job-id",
                      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                    },
                    {
                      "_id": "test-job-to-test-instruction-to-get-step-1-operation-1-2",
                      "name": "test-job-to-test-instruction-to-get-step-1-operation-1-2-name",
                      "author": "test_author",
                      "config": {
                        "_id": "test-job-config-to-link",
                        "name": "test-job-config-name-to-link",
                        "type": "rundeck",
                        "host": "http://example.com",
                        "author": "test-author",
                        "auth_token": "test-auth-token"
                      },
                      "job_id": "test-job-id",
                      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                    }
                  ]
                },
                {
                  "name": "test-instruction-to-get-step-1-operation-2-name",
                  "time_to_complete": {"seconds": 3, "unit":"s"},
                  "description": "test-instruction-to-get-step-1-operation-2-description",
                  "jobs": [
                    {
                      "_id": "test-job-to-test-instruction-to-get-step-1-operation-2-1",
                      "name": "test-job-to-test-instruction-to-get-step-1-operation-2-1-name",
                      "author": "test_author",
                      "config": {
                        "_id": "test-job-config-to-link",
                        "name": "test-job-config-name-to-link",
                        "type": "rundeck",
                        "host": "http://example.com",
                        "author": "test-author",
                        "auth_token": "test-auth-token"
                      },
                      "job_id": "test-job-id",
                      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                    }
                  ]
                }
              ],
              "stop_on_fail": true,
              "endpoint": "test-instruction-to-get-step-1-endpoint"
            },
            {
              "name": "test-instruction-to-get-step-2-name",
              "operations": [
                {
                  "name": "test-instruction-to-get-step-2-operation-1-name",
                  "time_to_complete": {"seconds": 6, "unit":"s"},
                  "description": "test-instruction-to-get-step-2-operation-1-description",
                  "jobs": [
                    {
                      "_id": "test-job-to-test-instruction-to-get-step-2-operation-1-1",
                      "name": "test-job-to-test-instruction-to-get-step-2-operation-1-1-name",
                      "author": "test_author",
                      "config": {
                        "_id": "test-job-config-to-link",
                        "name": "test-job-config-name-to-link",
                        "type": "rundeck",
                        "host": "http://example.com",
                        "author": "test-author",
                        "auth_token": "test-auth-token"
                      },
                      "job_id": "test-job-id",
                      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                    },
                    {
                      "_id": "test-job-to-test-instruction-to-get-step-2-operation-1-2",
                      "name": "test-job-to-test-instruction-to-get-step-2-operation-1-2-name",
                      "author": "test_author",
                      "config": {
                        "_id": "test-job-config-to-link",
                        "name": "test-job-config-name-to-link",
                        "type": "rundeck",
                        "host": "http://example.com",
                        "author": "test-author",
                        "auth_token": "test-auth-token"
                      },
                      "job_id": "test-job-id",
                      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                    }
                  ]
                }
              ],
              "stop_on_fail": true,
              "endpoint": "test-instruction-to-get-step-2-endpoint"
            }
          ],
          "last_executed_by": {
            "_id": "root",
            "username": "root"
          },
          "last_executed_on": 1596712203,
          "last_modified": 1596712203
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

  Scenario: given get all request should return instructions with flags
    When I am admin
    When I do GET /api/v4/cat/instructions?search=test-instruction-to-get-name&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-instruction-to-get",
          "deletable": true,
          "running": false
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

  Scenario: GET a instruction but unauthorized
    When I do GET /api/v4/cat/instructions/test-instruction-to-get
    Then the response code should be 401

  Scenario: GET a instruction but without permissions
    When I am noperms
    When I do GET /api/v4/cat/instructions/test-instruction-to-get
    Then the response code should be 403

  Scenario: Get a instruction with success
    When I am admin
    When I do GET /api/v4/cat/instructions/test-instruction-to-get
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "_id": "test-instruction-to-get",
      "name": "test-instruction-to-get-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "alarm_patterns": null,
      "description": "test-instruction-to-get-description",
      "author": "test-instruction-to-get-author",
      "enabled": true,
      "rating": 0,
      "comments": [],
      "avg_complete_time": 10,
      "month_executions": 0,
      "steps": [
      {
        "name": "test-instruction-to-get-step-1-name",
        "operations": [
          {
             "name": "test-instruction-to-get-step-1-operation-1-name",
             "time_to_complete": {"seconds": 1, "unit":"s"},
             "description": "test-instruction-to-get-step-1-operation-1-description",
             "jobs": [
               {
                 "_id": "test-job-to-test-instruction-to-get-step-1-operation-1-1",
                 "name": "test-job-to-test-instruction-to-get-step-1-operation-1-1-name",
                 "author": "test_author",
                 "config": {
                   "_id": "test-job-config-to-link",
                   "name": "test-job-config-name-to-link",
                   "type": "rundeck",
                   "host": "http://example.com",
                   "author": "test-author",
                   "auth_token": "test-auth-token"
                 },
                 "job_id": "test-job-id",
                 "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
               },
               {
                 "_id": "test-job-to-test-instruction-to-get-step-1-operation-1-2",
                 "name": "test-job-to-test-instruction-to-get-step-1-operation-1-2-name",
                 "author": "test_author",
                 "config": {
                   "_id": "test-job-config-to-link",
                   "name": "test-job-config-name-to-link",
                   "type": "rundeck",
                   "host": "http://example.com",
                   "author": "test-author",
                   "auth_token": "test-auth-token"
                 },
                 "job_id": "test-job-id",
                 "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
               }
             ]
            },
           {
             "name": "test-instruction-to-get-step-1-operation-2-name",
             "time_to_complete": {"seconds": 3, "unit":"s"},
             "description": "test-instruction-to-get-step-1-operation-2-description",
             "jobs": [
               {
                 "_id": "test-job-to-test-instruction-to-get-step-1-operation-2-1",
                 "name": "test-job-to-test-instruction-to-get-step-1-operation-2-1-name",
                 "author": "test_author",
                 "config": {
                   "_id": "test-job-config-to-link",
                   "name": "test-job-config-name-to-link",
                   "type": "rundeck",
                   "host": "http://example.com",
                   "author": "test-author",
                   "auth_token": "test-auth-token"
                 },
                 "job_id": "test-job-id",
                 "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
               }
             ]
           }
        ],
        "stop_on_fail": true,
        "endpoint": "test-instruction-to-get-step-1-endpoint"
      },
      {
       "name": "test-instruction-to-get-step-2-name",
       "operations": [
         {
           "name": "test-instruction-to-get-step-2-operation-1-name",
           "time_to_complete": {"seconds": 6, "unit":"s"},
           "description": "test-instruction-to-get-step-2-operation-1-description",
           "jobs": [
             {
               "_id": "test-job-to-test-instruction-to-get-step-2-operation-1-1",
               "name": "test-job-to-test-instruction-to-get-step-2-operation-1-1-name",
               "author": "test_author",
               "config": {
                 "_id": "test-job-config-to-link",
                 "name": "test-job-config-name-to-link",
                 "type": "rundeck",
                 "host": "http://example.com",
                 "author": "test-author",
                 "auth_token": "test-auth-token"
               },
               "job_id": "test-job-id",
               "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
             },
             {
               "_id": "test-job-to-test-instruction-to-get-step-2-operation-1-2",
               "name": "test-job-to-test-instruction-to-get-step-2-operation-1-2-name",
               "author": "test_author",
               "config": {
                 "_id": "test-job-config-to-link",
                 "name": "test-job-config-name-to-link",
                 "type": "rundeck",
                 "host": "http://example.com",
                 "author": "test-author",
                 "auth_token": "test-auth-token"
               },
               "job_id": "test-job-id",
               "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
             }
           ]
         }
       ],
       "stop_on_fail": true,
       "endpoint": "test-instruction-to-get-step-2-endpoint"
      }
      ],
      "last_executed_by": {
        "_id": "root",
        "username": "root"
      },
      "last_executed_on": 1596712203,
      "last_modified": 1596712203
    }
    """

  Scenario: Get a instruction with running executions
    When I am admin
    When I do GET /api/v4/cat/instructions?search=test-instruction-execution-running&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-instruction-execution-running",
          "deletable": false,
          "running": true
        }
      ]
    }
    """

  Scenario: Get a instruction with not found response
    When I am admin
    When I do GET /api/v4/cat/instructions/test-not-found
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """