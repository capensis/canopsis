Feature: create a instruction
  I need to be able to create a instruction
  Only admin should be able to create a instruction

  Scenario: POST a valid instruction but unauthorized
    When I do POST /api/v4/cat/instructions:
    """
    {
      "name": "test-instruction-new-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-new-description",
      "author": "test-instruction-new-author",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-new-step-1",
          "operations": [
            {
              "name": "test-instruction-new-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-new-step-1-operation-1-description",
              "jobs": [
                "test-job-to-test-instruction-new-step-1-operation-1-1",
                "test-job-to-test-instruction-new-step-1-operation-1-2"
              ]
            },
            {
              "name": "test-instruction-new-step-1-operation-2",
              "time_to_complete": {"seconds": 3, "unit":"s"},
              "description": "test-instruction-new-step-1-operation-2-description",
              "jobs": [
                "test-job-to-test-instruction-new-step-1-operation-2-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-new-step-1-endpoint"
        },
        {
          "name": "test-instruction-new-step-2",
          "operations": [
            {
              "name": "test-instruction-new-step-2-operation-1",
              "time_to_complete": {"seconds": 6, "unit":"s"},
              "description": "test-instruction-new-step-2-operation-1-description",
              "jobs": [
                "test-job-to-test-instruction-new-step-2-operation-1-1",
                "test-job-to-test-instruction-new-step-2-operation-1-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-new-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 401

  Scenario: POST a valid instruction but without permissions
    When I am noperms
    When I do POST /api/v4/cat/instructions:
    """
      {
        "name": "test-instruction-new-name",
        "entity_patterns": [
          {
            "name": "test filter"
          }
        ],
        "description": "test-instruction-new-description",
        "author": "test-instruction-new-author",
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-new-step-1",
            "operations": [
              {
                "name": "test-instruction-new-step-1-operation-1",
                "time_to_complete": {"seconds": 1, "unit":"s"},
                "description": "test-instruction-new-step-1-operation-1-description",
                "jobs": [
                  "test-job-to-test-instruction-new-step-1-operation-1-1",
                  "test-job-to-test-instruction-new-step-1-operation-1-2"
                ]
              },
              {
                "name": "test-instruction-new-step-1-operation-2",
                "time_to_complete": {"seconds": 3, "unit":"s"},
                "description": "test-instruction-new-step-1-operation-2-description",
                "jobs": [
                  "test-job-to-test-instruction-new-step-1-operation-2-1"
                ]
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-new-step-1-endpoint"
          },
          {
            "name": "test-instruction-new-step-2",
            "operations": [
              {
                "name": "test-instruction-new-step-2-operation-1",
                "time_to_complete": {"seconds": 6, "unit":"s"},
                "description": "test-instruction-new-step-2-operation-1-description",
                "jobs": [
                  "test-job-to-test-instruction-new-step-2-operation-1-1",
                  "test-job-to-test-instruction-new-step-2-operation-1-2"
                ]
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-new-step-2-endpoint"
          }
        ]
      }
    """
    Then the response code should be 403

  Scenario: POST a valid instruction
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """
    {
      "name": "test-instruction-new-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-new-description",
      "author": "test-instruction-new-author",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-new-step-1",
          "operations": [
            {
              "name": "test-instruction-new-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-new-step-1-operation-1-description",
              "jobs": [
                "test-job-to-test-instruction-new-step-1-operation-1-1",
                "test-job-to-test-instruction-new-step-1-operation-1-2"
              ]
            },
            {
              "name": "test-instruction-new-step-1-operation-2",
              "time_to_complete": {"seconds": 3, "unit":"s"},
              "description": "test-instruction-new-step-1-operation-2-description",
              "jobs": [
                "test-job-to-test-instruction-new-step-1-operation-2-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-new-step-1-endpoint"
        },
        {
          "name": "test-instruction-new-step-2",
          "operations": [
            {
              "name": "test-instruction-new-step-2-operation-1",
              "time_to_complete": {"seconds": 6, "unit":"s"},
              "description": "test-instruction-new-step-2-operation-1-description",
              "jobs": [
                "test-job-to-test-instruction-new-step-2-operation-1-1",
                "test-job-to-test-instruction-new-step-2-operation-1-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-new-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "author": "test-instruction-new-author",
      "avg_complete_time": 0,
      "description": "test-instruction-new-description",
      "enabled": true,
      "last_executed_by": null,
      "last_executed_on": null,
      "month_executions": 0,
      "name": "test-instruction-new-name",
      "rating": 0,
      "steps": [
        {
          "stop_on_fail": true,
          "endpoint": "test-instruction-new-step-1-endpoint",
          "name": "test-instruction-new-step-1",
          "operations": [
            {
              "description": "test-instruction-new-step-1-operation-1-description",
              "jobs": [
                {
                  "_id": "test-job-to-test-instruction-new-step-1-operation-1-1",
                  "author": "test_author",
                  "config": {
                    "_id": "test-job-config-to-link",
                    "auth_token": "test-auth-token",
                    "author": "test-author",
                    "host": "http://example.com",
                    "name": "test-job-config-name-to-link",
                    "type": "rundeck"
                  },
                  "job_id": "test-job-id",
                  "name": "test-job-to-test-instruction-new-step-1-operation-1-1-name",
                  "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                },
                {
                  "_id": "test-job-to-test-instruction-new-step-1-operation-1-2",
                  "author": "test_author",
                  "config": {
                    "_id": "test-job-config-to-link",
                    "auth_token": "test-auth-token",
                    "author": "test-author",
                    "host": "http://example.com",
                    "name": "test-job-config-name-to-link",
                    "type": "rundeck"
                  },
                  "job_id": "test-job-id",
                  "name": "test-job-to-test-instruction-new-step-1-operation-1-2-name",
                  "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                }
              ],
              "name": "test-instruction-new-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"}
            },
            {
              "description": "test-instruction-new-step-1-operation-2-description",
              "jobs": [
                {
                  "_id": "test-job-to-test-instruction-new-step-1-operation-2-1",
                  "author": "test_author",
                  "config": {
                    "_id": "test-job-config-to-link",
                    "auth_token": "test-auth-token",
                    "author": "test-author",
                    "host": "http://example.com",
                    "name": "test-job-config-name-to-link",
                    "type": "rundeck"
                  },
                  "job_id": "test-job-id",
                  "name": "test-job-to-test-instruction-new-step-1-operation-2-1-name",
                  "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                }
              ],
              "name": "test-instruction-new-step-1-operation-2",
              "time_to_complete": {"seconds": 3, "unit":"s"}
            }
          ]
        },
        {
          "stop_on_fail": true,
          "endpoint": "test-instruction-new-step-2-endpoint",
          "name": "test-instruction-new-step-2",
          "operations": [
            {
              "description": "test-instruction-new-step-2-operation-1-description",
              "jobs": [
                {
                  "_id": "test-job-to-test-instruction-new-step-2-operation-1-1",
                  "author": "test_author",
                  "config": {
                    "_id": "test-job-config-to-link",
                    "auth_token": "test-auth-token",
                    "author": "test-author",
                    "host": "http://example.com",
                    "name": "test-job-config-name-to-link",
                    "type": "rundeck"
                  },
                  "job_id": "test-job-id",
                  "name": "test-job-to-test-instruction-new-step-2-operation-1-1-name",
                  "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                },
                {
                  "_id": "test-job-to-test-instruction-new-step-2-operation-1-2",
                  "author": "test_author",
                  "config": {
                    "_id": "test-job-config-to-link",
                    "auth_token": "test-auth-token",
                    "author": "test-author",
                    "host": "http://example.com",
                    "name": "test-job-config-name-to-link",
                    "type": "rundeck"
                  },
                  "job_id": "test-job-id",
                  "name": "test-job-to-test-instruction-new-step-2-operation-1-2-name",
                  "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                }
              ],
              "name": "test-instruction-new-step-2-operation-1",
              "time_to_complete": {"seconds": 6, "unit":"s"}
            }
          ]
        }
      ]
    }
    """

  Scenario: POST a valid instruction with custom id
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """
    {
      "_id": "custom-id",
      "name": "test-instruction-new-name-custom-id-1",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-new-description",
      "author": "test-instruction-new-author",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-new-step-1",
          "operations": [
            {
              "name": "test-instruction-new-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-new-step-1-operation-1-description",
              "jobs": [
                "test-job-to-test-instruction-new-step-1-operation-1-1",
                "test-job-to-test-instruction-new-step-1-operation-1-2"
              ]
            },
            {
              "name": "test-instruction-new-step-1-operation-2",
              "time_to_complete": {"seconds": 3, "unit":"s"},
              "description": "test-instruction-new-step-1-operation-2-description",
              "jobs": [
                "test-job-to-test-instruction-new-step-1-operation-2-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-new-step-1-endpoint"
        },
        {
          "name": "test-instruction-new-step-2",
          "operations": [
            {
              "name": "test-instruction-new-step-2-operation-1",
              "time_to_complete": {"seconds": 6, "unit":"s"},
              "description": "test-instruction-new-step-2-operation-1-description",
              "jobs": [
                "test-job-to-test-instruction-new-step-2-operation-1-1",
                "test-job-to-test-instruction-new-step-2-operation-1-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-new-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/cat/instructions/custom-id
    Then the response code should be 200

  Scenario: POST a valid instruction with custom id that already exist should cause dup error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """
    {
      "_id": "test-instruction-to-update",
      "name": "test-instruction-new-name-custom-id-2",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-new-description",
      "author": "test-instruction-new-author",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-new-step-1",
          "operations": [
            {
              "name": "test-instruction-new-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-new-step-1-operation-1-description",
              "jobs": [
                "test-job-to-test-instruction-new-step-1-operation-1-1",
                "test-job-to-test-instruction-new-step-1-operation-1-2"
              ]
            },
            {
              "name": "test-instruction-new-step-1-operation-2",
              "time_to_complete": {"seconds": 3, "unit":"s"},
              "description": "test-instruction-new-step-1-operation-2-description",
              "jobs": [
                "test-job-to-test-instruction-new-step-1-operation-2-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-new-step-1-endpoint"
        },
        {
          "name": "test-instruction-new-step-2",
          "operations": [
            {
              "name": "test-instruction-new-step-2-operation-1",
              "time_to_complete": {"seconds": 6, "unit":"s"},
              "description": "test-instruction-new-step-2-operation-1-description",
              "jobs": [
                "test-job-to-test-instruction-new-step-2-operation-1-1",
                "test-job-to-test-instruction-new-step-2-operation-1-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-new-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "_id": "ID already exists"
      }
    }
    """

  Scenario: POST a valid instruction
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """
    {
      "name": "test-instruction-new-name-2",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-new-description",
      "author": "test-instruction-new-author",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-new-step-1",
          "operations": [
            {
              "name": "test-instruction-new-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-new-step-1-operation-1-description",
              "jobs": [
                "test-job-to-test-instruction-new-step-1-operation-1-1",
                "test-job-to-test-instruction-new-step-1-operation-1-2"
              ]
            },
            {
              "name": "test-instruction-new-step-1-operation-2",
              "time_to_complete": {"seconds": 3, "unit":"s"},
              "description": "test-instruction-new-step-1-operation-2-description",
              "jobs": [
                "test-job-to-test-instruction-new-step-1-operation-2-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-new-step-1-endpoint"
        },
        {
          "name": "test-instruction-new-step-2",
          "operations": [
            {
              "name": "test-instruction-new-step-2-operation-1",
              "time_to_complete": {"seconds": 6, "unit":"s"},
              "description": "test-instruction-new-step-2-operation-1-description",
              "jobs": [
                "test-job-to-test-instruction-new-step-2-operation-1-1",
                "test-job-to-test-instruction-new-step-2-operation-1-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-new-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/cat/instructions/{{ .lastResponse._id}}
    Then the response code should be 200

  Scenario: POST an invalid instruction where job doesn't exist
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """
      {
        "name": "test-instruction-new-name-3",
        "entity_patterns": [
          {
            "name": "test filter"
          }
        ],
        "description": "test-instruction-new-description",
        "author": "test-instruction-new-author",
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-new-step-1",
            "operations": [
              {
                "name": "test-instruction-new-step-1-operation-1",
                "time_to_complete": {"seconds": 1, "unit":"s"},
                "description": "test-instruction-new-step-1-operation-1-description",
                "jobs": [
                  "test-job-to-test-instruction-new-step-1-operation-1-1-NOT-EXIST",
                  "test-job-to-test-instruction-new-step-1-operation-1-2"
                ]
              },
              {
                "name": "test-instruction-new-step-1-operation-2",
                "time_to_complete": {"seconds": 3, "unit":"s"},
                "description": "test-instruction-new-step-1-operation-2-description",
                "jobs": [
                  "test-job-to-test-instruction-new-step-1-operation-2-1"
                ]
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-new-step-1-endpoint"
          },
          {
            "name": "test-instruction-new-step-2",
            "operations": [
              {
                "name": "test-instruction-new-step-2-operation-1",
                "time_to_complete": {"seconds": 6, "unit":"s"},
                "description": "test-instruction-new-step-2-operation-1-description",
                "jobs": [
                  "test-job-to-test-instruction-new-step-2-operation-1-1",
                  "test-job-to-test-instruction-new-step-2-operation-1-2"
                ]
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-new-step-2-endpoint"
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

  Scenario: POST an invalid instruction where name already exists
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """
      {
        "name": "test-instruction-to-get-name",
        "entity_patterns": [
          {
            "name": "test filter"
          }
        ],
        "description": "test-instruction-new-description",
        "author": "test-instruction-new-author",
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-new-step-1",
            "operations": [
              {
                "name": "test-instruction-new-step-1-operation-1",
                "time_to_complete": {"seconds": 1, "unit":"s"},
                "description": "test-instruction-new-step-1-operation-1-description",
                "jobs": [
                  "test-job-to-test-instruction-new-step-1-operation-1-1",
                  "test-job-to-test-instruction-new-step-1-operation-1-2"
                ]
              },
              {
                "name": "test-instruction-new-step-1-operation-2",
                "time_to_complete": {"seconds": 3, "unit":"s"},
                "description": "test-instruction-new-step-1-operation-2-description",
                "jobs": [
                  "test-job-to-test-instruction-new-step-1-operation-2-1"
                ]
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-new-step-1-endpoint"
          },
          {
            "name": "test-instruction-new-step-2",
            "operations": [
              {
                "name": "test-instruction-new-step-2-operation-1",
                "time_to_complete": {"seconds": 6, "unit":"s"},
                "description": "test-instruction-new-step-2-operation-1-description",
                "jobs": [
                  "test-job-to-test-instruction-new-step-2-operation-1-1",
                  "test-job-to-test-instruction-new-step-2-operation-1-2"
                ]
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-new-step-2-endpoint"
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

  Scenario: POST an invalid instruction with invalid patterns
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """
      {
        "name": "test-instruction-invalid-3",
        "entity_patterns": [
          {
            "name": {
              "regex_match": "name:.*"
            }
          },
          {}
        ],
        "description": "test-instruction-new-description",
        "author": "test-instruction-new-author",
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-new-step-1",
            "operations": [
              {
                "name": "test-instruction-new-step-1-operation-1",
                "time_to_complete": {"seconds": 1, "unit":"s"},
                "description": "test-instruction-new-step-1-operation-1-description",
                "jobs": [
                  "test-job-to-test-instruction-new-step-1-operation-1-1",
                  "test-job-to-test-instruction-new-step-1-operation-1-2"
                ]
              },
              {
                "name": "test-instruction-new-step-1-operation-2",
                "time_to_complete": {"seconds": 3, "unit":"s"},
                "description": "test-instruction-new-step-1-operation-2-description",
                "jobs": [
                  "test-job-to-test-instruction-new-step-1-operation-2-1"
                ]
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-new-step-1-endpoint"
          },
          {
            "name": "test-instruction-new-step-2",
            "operations": [
              {
                "name": "test-instruction-new-step-2-operation-1",
                "time_to_complete": {"seconds": 6, "unit":"s"},
                "description": "test-instruction-new-step-2-operation-1-description",
                "jobs": [
                  "test-job-to-test-instruction-new-step-2-operation-1-1",
                  "test-job-to-test-instruction-new-step-2-operation-1-2"
                ]
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-new-step-2-endpoint"
          }
        ]
      }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "entity_patterns": "entity pattern list contains an empty pattern."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """
      {
        "name": "test-instruction-invalid-4",
        "alarm_patterns": [
          {
            "v": {
              "resource": {
                "regex_match": "name:.*"
              }
            }
          },
          {}
        ],
        "description": "test-instruction-new-description",
        "author": "test-instruction-new-author",
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-new-step-1",
            "operations": [
              {
                "name": "test-instruction-new-step-1-operation-1",
                "time_to_complete": {"seconds": 1, "unit":"s"},
                "description": "test-instruction-new-step-1-operation-1-description",
                "jobs": [
                  "test-job-to-test-instruction-new-step-1-operation-1-1",
                  "test-job-to-test-instruction-new-step-1-operation-1-2"
                ]
              },
              {
                "name": "test-instruction-new-step-1-operation-2",
                "time_to_complete": {"seconds": 3, "unit":"s"},
                "description": "test-instruction-new-step-1-operation-2-description",
                "jobs": [
                  "test-job-to-test-instruction-new-step-1-operation-2-1"
                ]
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-new-step-1-endpoint"
          },
          {
            "name": "test-instruction-new-step-2",
            "operations": [
              {
                "name": "test-instruction-new-step-2-operation-1",
                "time_to_complete": {"seconds": 6, "unit":"s"},
                "description": "test-instruction-new-step-2-operation-1-description",
                "jobs": [
                  "test-job-to-test-instruction-new-step-2-operation-1-1",
                  "test-job-to-test-instruction-new-step-2-operation-1-2"
                ]
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-new-step-2-endpoint"
          }
        ]
      }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "alarm_patterns": "alarm pattern list contains an empty pattern."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """
      {
        "name": "test-instruction-invalid-5",
        "alarm_patterns": [
          {
            "qwe": {
              "resource": {
                "regex_match": "name:.*"
              }
            }
          }
        ],
        "description": "test-instruction-new-description",
        "author": "test-instruction-new-author",
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-new-step-1",
            "operations": [
              {
                "name": "test-instruction-new-step-1-operation-1",
                "time_to_complete": {"seconds": 1, "unit":"s"},
                "description": "test-instruction-new-step-1-operation-1-description",
                "jobs": [
                  "test-job-to-test-instruction-new-step-1-operation-1-1",
                  "test-job-to-test-instruction-new-step-1-operation-1-2"
                ]
              },
              {
                "name": "test-instruction-new-step-1-operation-2",
                "time_to_complete": {"seconds": 3, "unit":"s"},
                "description": "test-instruction-new-step-1-operation-2-description",
                "jobs": [
                  "test-job-to-test-instruction-new-step-1-operation-2-1"
                ]
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-new-step-1-endpoint"
          },
          {
            "name": "test-instruction-new-step-2",
            "operations": [
              {
                "name": "test-instruction-new-step-2-operation-1",
                "time_to_complete": {"seconds": 6, "unit":"s"},
                "description": "test-instruction-new-step-2-operation-1-description",
                "jobs": [
                  "test-job-to-test-instruction-new-step-2-operation-1-1",
                  "test-job-to-test-instruction-new-step-2-operation-1-2"
                ]
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-new-step-2-endpoint"
          }
        ]
      }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "alarm_patterns": "Invalid alarm pattern list."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """
      {
        "name": "test-instruction-invalid-6",
        "entity_patterns": [
          {
            "qwe": {
              "regex_match": "name:.*"
            }
          }
        ],
        "description": "test-instruction-new-description",
        "author": "test-instruction-new-author",
        "enabled": true,
        "steps": [
          {
            "name": "test-instruction-new-step-1",
            "operations": [
              {
                "name": "test-instruction-new-step-1-operation-1",
                "time_to_complete": {"seconds": 1, "unit":"s"},
                "description": "test-instruction-new-step-1-operation-1-description",
                "jobs": [
                  "test-job-to-test-instruction-new-step-1-operation-1-1",
                  "test-job-to-test-instruction-new-step-1-operation-1-2"
                ]
              },
              {
                "name": "test-instruction-new-step-1-operation-2",
                "time_to_complete": {"seconds": 3, "unit":"s"},
                "description": "test-instruction-new-step-1-operation-2-description",
                "jobs": [
                  "test-job-to-test-instruction-new-step-1-operation-2-1"
                ]
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-new-step-1-endpoint"
          },
          {
            "name": "test-instruction-new-step-2",
            "operations": [
              {
                "name": "test-instruction-new-step-2-operation-1",
                "time_to_complete": {"seconds": 6, "unit":"s"},
                "description": "test-instruction-new-step-2-operation-1-description",
                "jobs": [
                  "test-job-to-test-instruction-new-step-2-operation-1-1",
                  "test-job-to-test-instruction-new-step-2-operation-1-2"
                ]
              }
            ],
            "stop_on_fail": true,
            "endpoint": "test-instruction-new-step-2-endpoint"
          }
        ]
      }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "entity_patterns": "Invalid entity pattern list."
      }
    }
    """