Feature: get a instruction
  I need to be able to get a instruction
  Only admin should be able to get a instruction

  Scenario: given get all request should return instructions
    When I am admin
    When I do GET /api/v4/cat/instructions?search=test-instruction-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-instruction-to-get-1",
          "type": 0,
          "status": 0,
          "name": "test-instruction-to-get-1-name",
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-instruction-to-get-1-pattern"
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
                  "value": "test-instruction-to-get-1-pattern"
                }
              }
            ]
          ],
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "description": "test-instruction-to-get-1-description",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "enabled": true,
          "timeout_after_execution": {
            "value": 2,
            "unit": "s"
          },
          "steps": [
            {
              "name": "test-instruction-to-get-1-step-1-name",
              "operations": [
                {
                  "name": "test-instruction-to-get-1-step-1-operation-1-name",
                  "time_to_complete": {"value": 1, "unit":"s"},
                  "description": "test-instruction-to-get-1-step-1-operation-1-description",
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
                        "auth_username": "",
                        "auth_token": "test-auth-token"
                      },
                      "job_id": "test-job-to-instruction-edit-1-external-id",
                      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
                      "query": null,
                      "multiple_executions": false
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
                        "auth_username": "",
                        "auth_token": "test-auth-token"
                      },
                      "job_id": "test-job-to-instruction-edit-2-external-id",
                      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
                      "query": null,
                      "multiple_executions": false
                    }
                  ]
                },
                {
                  "name": "test-instruction-to-get-1-step-1-operation-2-name",
                  "time_to_complete": {"value": 3, "unit":"s"},
                  "description": "test-instruction-to-get-1-step-1-operation-2-description",
                  "jobs": [
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
                        "auth_username": "",
                        "auth_token": "test-auth-token"
                      },
                      "job_id": "test-job-to-instruction-edit-2-external-id",
                      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
                      "query": null,
                      "multiple_executions": false
                    }
                  ]
                }
              ],
              "stop_on_fail": true,
              "endpoint": "test-instruction-to-get-1-step-1-endpoint"
            },
            {
              "name": "test-instruction-to-get-1-step-2-name",
              "operations": [
                {
                  "name": "test-instruction-to-get-1-step-2-operation-1-name",
                  "time_to_complete": {"value": 6, "unit":"s"},
                  "description": "test-instruction-to-get-1-step-2-operation-1-description",
                  "jobs": [
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
                        "auth_username": "",
                        "auth_token": "test-auth-token"
                      },
                      "job_id": "test-job-to-instruction-edit-2-external-id",
                      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
                      "query": null,
                      "multiple_executions": false
                    },
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
                        "auth_username": "",
                        "auth_token": "test-auth-token"
                      },
                      "job_id": "test-job-to-instruction-edit-1-external-id",
                      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
                      "query": null,
                      "multiple_executions": false
                    }
                  ]
                }
              ],
              "stop_on_fail": true,
              "endpoint": "test-instruction-to-get-1-step-2-endpoint"
            }
          ],
          "last_executed_on": 1596712203,
          "created": 1596712203,
          "last_modified": 1596712203
        },
        {
          "_id": "test-instruction-to-get-2",
          "type": 1,
          "status": 0,
          "name": "test-instruction-to-get-2-name",
          "corporate_entity_pattern": "test-pattern-to-filter-edit-2",
          "corporate_entity_pattern_title": "test-pattern-to-filter-edit-2-title",
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-filter-edit-2-pattern"
                }
              }
            ]
          ],
          "corporate_alarm_pattern": "test-pattern-to-filter-edit-1",
          "corporate_alarm_pattern_title": "test-pattern-to-filter-edit-1-title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-filter-edit-1-pattern"
                }
              }
            ]
          ],
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "description": "test-instruction-to-get-2-description",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "enabled": true,
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
                  "auth_username": "",
                  "auth_token": "test-auth-token"
                },
                "job_id": "test-job-to-instruction-edit-1-external-id",
                "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
                "query": null,
                "multiple_executions": false
              },
              "stop_on_fail": true
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
                  "auth_username": "",
                  "auth_token": "test-auth-token"
                },
                "job_id": "test-job-to-instruction-edit-2-external-id",
                "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
                "query": null,
                "multiple_executions": false
              }
            }
          ],
          "priority": 2,
          "triggers": ["create"],
          "timeout_after_execution": {
            "value": 2,
            "unit": "s"
          },
          "last_executed_on": 1596712203,
          "created": 1596712203,
          "last_modified": 1596712203
        },
        {
          "_id": "test-instruction-to-get-3",
          "type": 1,
          "status": 0,
          "name": "test-instruction-to-get-3-name",
          "old_alarm_patterns": [
            {
              "_id": "test-instruction-to-get-3-pattern"
            }
          ],
          "old_entity_patterns": [
            {
              "name": "test-instruction-to-get-3-pattern"
            }
          ],
          "description": "test-instruction-to-get-3-description",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "enabled": true,
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
                  "auth_username": "",
                  "auth_token": "test-auth-token"
                },
                "job_id": "test-job-to-instruction-edit-1-external-id",
                "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
                "query": null,
                "multiple_executions": false
              }
            }
          ],
          "priority": 3,
          "triggers": ["create"],
          "timeout_after_execution": {
            "value": 2,
            "unit": "s"
          },
          "last_executed_on": 1596712203,
          "created": 1596712203,
          "last_modified": 1596712203
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """

  Scenario: given get all request should return instructions with flags
    When I am admin
    When I do GET /api/v4/cat/instructions?search=test-instruction-to-get-1-name&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-instruction-to-get-1",
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

  Scenario: given filter by type request should return instructions
    When I am admin
    When I do GET /api/v4/cat/instructions?search=test-instruction-to-get&type=0
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-instruction-to-get-1"
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
    When I do GET /api/v4/cat/instructions?search=test-instruction-to-get&type=1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-instruction-to-get-2"
        },
        {
          "_id": "test-instruction-to-get-3"
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

  Scenario: GET a instruction but unauthorized
    When I do GET /api/v4/cat/instructions/test-instruction-to-get
    Then the response code should be 401

  Scenario: GET a instruction but without permissions
    When I am noperms
    When I do GET /api/v4/cat/instructions/test-instruction-to-get
    Then the response code should be 403

  Scenario: Get a instruction with success
    When I am admin
    When I do GET /api/v4/cat/instructions/test-instruction-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-instruction-to-get-1",
      "type": 0,
      "status": 0,
      "name": "test-instruction-to-get-1-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-get-1-pattern"
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
              "value": "test-instruction-to-get-1-pattern"
            }
          }
        ]
      ],
      "old_alarm_patterns": null,
      "old_entity_patterns": null,
      "description": "test-instruction-to-get-1-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "timeout_after_execution": {
        "value": 2,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-get-1-step-1-name",
          "operations": [
            {
               "name": "test-instruction-to-get-1-step-1-operation-1-name",
               "time_to_complete": {"value": 1, "unit":"s"},
               "description": "test-instruction-to-get-1-step-1-operation-1-description",
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
                     "auth_username": "",
                     "auth_token": "test-auth-token"
                   },
                   "job_id": "test-job-to-instruction-edit-1-external-id",
                   "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
                   "query": null,
                   "multiple_executions": false
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
                     "auth_username": "",
                     "auth_token": "test-auth-token"
                   },
                   "job_id": "test-job-to-instruction-edit-2-external-id",
                   "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
                   "query": null,
                   "multiple_executions": false
                 }
               ]
              },
             {
               "name": "test-instruction-to-get-1-step-1-operation-2-name",
               "time_to_complete": {"value": 3, "unit":"s"},
               "description": "test-instruction-to-get-1-step-1-operation-2-description",
               "jobs": [
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
                     "auth_username": "",
                     "auth_token": "test-auth-token"
                   },
                   "job_id": "test-job-to-instruction-edit-2-external-id",
                   "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
                   "query": null,
                   "multiple_executions": false
                 }
               ]
             }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-get-1-step-1-endpoint"
        },
        {
         "name": "test-instruction-to-get-1-step-2-name",
         "operations": [
           {
             "name": "test-instruction-to-get-1-step-2-operation-1-name",
             "time_to_complete": {"value": 6, "unit":"s"},
             "description": "test-instruction-to-get-1-step-2-operation-1-description",
             "jobs": [
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
                   "auth_username": "",
                   "auth_token": "test-auth-token"
                 },
                 "job_id": "test-job-to-instruction-edit-2-external-id",
                 "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
                 "query": null,
                 "multiple_executions": false
               },
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
                   "auth_username": "",
                   "auth_token": "test-auth-token"
                 },
                 "job_id": "test-job-to-instruction-edit-1-external-id",
                 "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
                 "query": null,
                 "multiple_executions": false
               }
             ]
           }
         ],
         "stop_on_fail": true,
         "endpoint": "test-instruction-to-get-1-step-2-endpoint"
        }
      ],
      "created": 1596712203,
      "last_modified": 1596712203
    }
    """

  Scenario: Get a instruction with running executions
    When I am admin
    When I do GET /api/v4/cat/instructions?search=test-instruction-execution-running&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """json
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
    """json
    {
      "error": "Not found"
    }
    """
