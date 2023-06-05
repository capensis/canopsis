Feature: Instruction update

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update
    Then the response code should be 401

  Scenario: given update request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update
    Then the response code should be 403

  Scenario: given update manual instruction request should return ok
    When I am admin
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-1:
    """json
    {
      "name": "test-instruction-to-update-1-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-update-1-pattern"
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
              "value": "test-instruction-to-update-1-pattern"
            }
          }
        ]
      ],
      "description": "test-instruction-to-update-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "active_on_pbh": ["test-default-maintenance-type"],
      "disabled_on_pbh": ["test-default-pause-type"],
      "steps": [
        {
          "name": "test-instruction-to-update-1-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-1-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-1-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            },
            {
              "name": "test-instruction-to-update-1-step-1-operation-2-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-1-step-1-operation-2-description",
              "jobs": [
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-1-step-1-endpoint"
        },
        {
          "name": "test-instruction-to-update-1-step-2-name",
          "operations": [
            {
              "name": "test-instruction-to-update-1-step-2-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-1-step-2-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-2",
                "test-job-to-instruction-edit-1"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-1-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-update-1",
      "type": 0,
      "status": 0,
      "name": "test-instruction-to-update-1-name",
      "created": 1596712203,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-update-1-pattern"
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
              "value": "test-instruction-to-update-1-pattern"
            }
          }
        ]
      ],
      "description": "test-instruction-to-update-1-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "active_on_pbh": ["test-default-maintenance-type"],
      "disabled_on_pbh": ["test-default-pause-type"],
      "steps": [
        {
          "name": "test-instruction-to-update-1-step-1-name",
          "operations": [
            {
              "name": "test-instruction-to-update-1-step-1-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-1-step-1-operation-1-description",
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
                    "auth_token": "test-auth-token"
                  },
                  "job_id": "test-job-to-instruction-edit-2-external-id",
                  "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                }
              ]
            },
            {
              "name": "test-instruction-to-update-1-step-1-operation-2-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-1-step-1-operation-2-description",
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
                    "auth_token": "test-auth-token"
                  },
                  "job_id": "test-job-to-instruction-edit-2-external-id",
                  "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                }
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-1-step-1-endpoint"
        },
        {
          "name": "test-instruction-to-update-1-step-2-name",
          "operations": [
            {
              "name": "test-instruction-to-update-1-step-2-operation-1-name",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-update-1-step-2-operation-1-description",
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
                    "auth_token": "test-auth-token"
                  },
                  "job_id": "test-job-to-instruction-edit-2-external-id",
                  "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
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
                    "auth_token": "test-auth-token"
                  },
                  "job_id": "test-job-to-instruction-edit-1-external-id",
                  "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
                }
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-update-1-step-2-endpoint"
        }
      ]
    }
    """
    When I do GET /api/v4/cat/instruction-stats?search=test-instruction-to-update-1&from=1000000000&to=2000000000
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "last_executed_on": 1596712203,
          "avg_complete_time": 400,
          "rating": 3.5
        }
      ]
    }
    """

  Scenario: given update auto instruction request should return ok
    When I am admin
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-2:
    """json
    {
      "name": "test-instruction-to-update-2-name",
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "description": "test-instruction-to-update-2-description",
      "enabled": true,
      "priority": 21,
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
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-update-2",
      "type": 1,
      "status": 0,
      "name": "test-instruction-to-update-2-name",
      "created": 1596712203,
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "corporate_alarm_pattern_title": "test-pattern-to-rule-edit-1-title",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-rule-edit-1-pattern"
            }
          }
        ]
      ],
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "corporate_entity_pattern_title": "test-pattern-to-rule-edit-2-title",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-rule-edit-2-pattern"
            }
          }
        ]
      ],
      "description": "test-instruction-to-update-2-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "priority": 21,
      "triggers": ["create"],
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
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
      ]
    }
    """
    When I do GET /api/v4/cat/instruction-stats?search=test-instruction-to-update-2&from=1000000000&to=2000000000
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "last_executed_on": 1596712203,
          "avg_complete_time": 400,
          "rating": 3.5
        }
      ]
    }
    """

  Scenario: given update request with already exists name should return error
    When I am admin
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update:
    """json
    {
      "name": "test-instruction-to-check-unique-name"
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

  Scenario: given update request for a instruction with old patterns should return keep old patterns
    When I am admin
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-3:
    """json
    {
      "name": "test-instruction-to-update-3-name",
      "description": "test-instruction-to-update-3-description",
      "enabled": true,
      "priority": 21,
      "triggers": ["create"],
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "jobs": [
        {
          "job": "test-job-to-instruction-edit-1"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-update-3",
      "type": 1,
      "status": 0,
      "name": "test-instruction-to-update-3-name",
      "description": "test-instruction-to-update-3-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "priority": 21,
      "triggers": ["create"],
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "created": 1596712203,
      "old_alarm_patterns": [
        {
          "_id": "test-instruction-to-update-3-pattern"
        }
      ],
      "old_entity_patterns": [
        {
          "name": "test-instruction-to-update-3-pattern"
        }
      ],
      "jobs": [
        {
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
        }
      ]
    }
    """
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-3:
    """json
    {
      "name": "test-instruction-to-update-3-name",
      "description": "test-instruction-to-update-3-description",
      "enabled": true,
      "priority": 21,
      "triggers": ["create"],
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-update-3-pattern"
            }
          }
        ]
      ],
      "jobs": [
        {
          "job": "test-job-to-instruction-edit-1"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-update-3",
      "type": 1,
      "status": 0,
      "name": "test-instruction-to-update-3-name",
      "description": "test-instruction-to-update-3-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "priority": 21,
      "triggers": ["create"],
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "created": 1596712203,
      "old_entity_patterns": null,
      "old_alarm_patterns": [
        {
          "_id": "test-instruction-to-update-3-pattern"
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-update-3-pattern"
            }
          }
        ]
      ],
      "jobs": [
        {
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
        }
      ]
    }
    """
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-3:
    """json
    {
      "name": "test-instruction-to-update-3-name",
      "description": "test-instruction-to-update-3-description",
      "enabled": true,
      "priority": 21,
      "triggers": ["create"],
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-update-3-pattern"
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
              "value": "test-instruction-to-update-3-pattern"
            }
          }
        ]
      ],
      "jobs": [
        {
          "job": "test-job-to-instruction-edit-1"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-update-3",
      "type": 1,
      "status": 0,
      "name": "test-instruction-to-update-3-name",
      "description": "test-instruction-to-update-3-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "priority": 21,
      "triggers": ["create"],
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "created": 1596712203,
      "old_entity_patterns": null,
      "old_alarm_patterns": null,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-update-3-pattern"
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
              "value": "test-instruction-to-update-3-pattern"
            }
          }
        ]
      ],
      "jobs": [
        {
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
        }
      ]
    }
    """

  Scenario: given update simplified manual instruction request should return ok
    When I am admin
    When I do PUT /api/v4/cat/instructions/test-instruction-to-update-4:
    """json
    {
      "name": "test-instruction-to-update-4-name",
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "description": "test-instruction-to-update-4-description",
      "enabled": true,
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
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-instruction-to-update-4",
      "type": 2,
      "status": 0,
      "name": "test-instruction-to-update-4-name",
      "created": 1596712203,
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "corporate_alarm_pattern_title": "test-pattern-to-rule-edit-1-title",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-rule-edit-1-pattern"
            }
          }
        ]
      ],
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "corporate_entity_pattern_title": "test-pattern-to-rule-edit-2-title",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-rule-edit-2-pattern"
            }
          }
        ]
      ],
      "description": "test-instruction-to-update-4-description",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
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
      ]
    }
    """
