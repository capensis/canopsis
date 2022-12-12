Feature: Create a instruction
  I need to be able to create a instruction
  Only admin should be able to create a instruction

  Scenario: given create manual instruction request should return ok
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-create-1-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-create-1-pattern"
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
              "value": "test-instruction-to-create-1-pattern"
            }
          },
          {
            "field": "v.creation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263992,
                "to": 1605264992
              }
            }
          },
          {
            "field": "v.activation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263992,
                "to": 1605264992
              }
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263992,
                "to": 1605264992
              }
            }
          }
        ]
      ],
      "description": "test-instruction-to-create-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "active_on_pbh": ["test-default-maintenance-type"],
      "disabled_on_pbh": ["test-default-pause-type"],
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
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-create-1-pattern"
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
              "value": "test-instruction-to-create-1-pattern"
            }
          },
          {
            "field": "v.creation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263992,
                "to": 1605264992
              }
            }
          },
          {
            "field": "v.activation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263992,
                "to": 1605264992
              }
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263992,
                "to": 1605264992
              }
            }
          }
        ]
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
      "active_on_pbh": ["test-default-maintenance-type"],
      "disabled_on_pbh": ["test-default-pause-type"],
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
                },
                {
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
                },
                {
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
              ],
              "name": "test-instruction-to-create-1-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"}
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/instructions/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "type": 0,
      "status": 0,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-create-1-pattern"
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
              "value": "test-instruction-to-create-1-pattern"
            }
          },
          {
            "field": "v.creation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263992,
                "to": 1605264992
              }
            }
          },
          {
            "field": "v.activation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263992,
                "to": 1605264992
              }
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263992,
                "to": 1605264992
              }
            }
          }
        ]
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
      "active_on_pbh": ["test-default-maintenance-type"],
      "disabled_on_pbh": ["test-default-pause-type"],
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
                },
                {
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
                },
                {
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
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "description": "test-instruction-to-create-2-description",
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
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "type": 1,
      "status": 0,
      "author": {
        "_id": "root",
        "name": "root"
      },
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
      "description": "test-instruction-to-create-2-description",
      "enabled": true,
      "priority": 21,
      "triggers": ["create"],
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
    When I do GET /api/v4/cat/instructions/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "type": 1,
      "status": 0,
      "author": {
        "_id": "root",
        "name": "root"
      },
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
      "description": "test-instruction-to-create-2-description",
      "enabled": true,
      "priority": 21,
      "triggers": ["create"],
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

  Scenario: given create request with custom id should return ok
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-create-4-id",
      "type": 0,
      "name": "test-instruction-to-create-4-name",
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
              "description": "test-instruction-to-create-4-step-1-operation-1-description"
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
        "triggers": "Triggers is missing.",
        "timeout_after_execution": "TimeoutAfterExecution is missing."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "jobs": [],
      "triggers": []
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "jobs": "Jobs is missing.",
        "triggers": "Triggers is missing."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "triggers": ["notexist"]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "triggers.0": "Triggers[0] must be one of [create statedec stateinc changestate unsnooze activate pbhenter pbhleave]."
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

  Scenario: given create request with not exist pbehavior type should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "active_on_pbh": ["notexist"]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "active_on_pbh": "ActiveOnPbh doesn't exist."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "disabled_on_pbh": ["notexist"]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "disabled_on_pbh": "DisabledOnPbh doesn't exist."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "active_on_pbh": ["test-default-maintenance-type", "notexist"]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "active_on_pbh": "ActiveOnPbh doesn't exist."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "disabled_on_pbh": ["test-default-maintenance-type", "notexist"]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "disabled_on_pbh": "DisabledOnPbh doesn't exist."
      }
    }
    """

  Scenario: given create request with invalid patterns should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "entity_pattern": [
        []
      ],
      "alarm_pattern": [
        []
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "entity_pattern": "EntityPattern is invalid entity pattern.",
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-create-5-pattern"
            }
          },
          {
            "field": "last_event_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
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
              "value": "test-instruction-to-create-5-pattern"
            }
          },
          {
            "field": "v.last_event_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "entity_pattern": "EntityPattern is invalid entity pattern.",
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-create-5-pattern"
            }
          },
          {
            "field": "v.last_update_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-create-5-pattern"
            }
          },
          {
            "field": "v.resolved",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-create-5-pattern"
            }
          },
          {
            "field": "v.created_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-create-5-pattern"
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "corporate_alarm_pattern": "test-pattern-not-exist",
      "type": 1,
      "name": "test-instruction-to-create-3-name",
      "description": "test-instruction-to-create-3-description",
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
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "corporate_alarm_pattern": "CorporateAlarmPattern doesn't exist."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "corporate_entity_pattern": "test-pattern-not-exist",
      "type": 1,
      "name": "test-instruction-to-create-3-name",
      "description": "test-instruction-to-create-3-description",
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
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "corporate_entity_pattern": "CorporateEntityPattern doesn't exist."
      }
    }
    """
    
  Scenario: given create simplified manual instruction request should return ok
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 2,
      "name": "test-instruction-to-create-6-name",
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "description": "test-instruction-to-create-6-description",
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
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "type": 2,
      "status": 0,
      "author": {
        "_id": "root",
        "name": "root"
      },
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
      "description": "test-instruction-to-create-6-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "name": "test-instruction-to-create-6-name",
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
    When I do GET /api/v4/cat/instructions/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "type": 2,
      "status": 0,
      "author": {
        "_id": "root",
        "name": "root"
      },
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
      "description": "test-instruction-to-create-6-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "name": "test-instruction-to-create-6-name",
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
    
  Scenario: given create simplified manual instruction request with manual instruction fields should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 2,
      "name": "test-instruction-to-create-7-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-create-7-pattern"
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
              "value": "test-instruction-to-create-7-pattern"
            }
          },
          {
            "field": "v.creation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263992,
                "to": 1605264992
              }
            }
          },
          {
            "field": "v.activation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263992,
                "to": 1605264992
              }
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263992,
                "to": 1605264992
              }
            }
          }
        ]
      ],
      "description": "test-instruction-to-create-7-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "active_on_pbh": ["test-default-maintenance-type"],
      "disabled_on_pbh": ["test-default-pause-type"],
      "steps": [
        {
          "name": "test-instruction-to-create-7-step-1",
          "operations": [
            {
              "name": "test-instruction-to-create-7-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-create-7-step-1-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            },
            {
              "name": "test-instruction-to-create-7-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-to-create-7-step-1-operation-2-description",
              "jobs": [
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-7-step-1-endpoint"
        },
        {
          "name": "test-instruction-to-create-7-step-2",
          "operations": [
            {
              "name": "test-instruction-to-create-7-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-to-create-7-step-2-operation-1-description",
              "jobs": [
                "test-job-to-instruction-edit-1",
                "test-job-to-instruction-edit-2"
              ]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-create-7-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "jobs": "Jobs is missing.",
        "steps": "Steps is not empty."
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
