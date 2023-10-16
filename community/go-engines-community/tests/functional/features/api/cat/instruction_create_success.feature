Feature: Create a instruction
  I need to be able to create a instruction
  Only admin should be able to create a instruction

  @concurrent
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

  @concurrent
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
      "triggers": [
        {
          "type": "create"
        }
      ],
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
      "triggers": [
        {
          "type": "create"
        }
      ],
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
      "triggers": [
        {
          "type": "create"
        }
      ],
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

  @concurrent
  Scenario: given create simplified manual instruction request should return ok
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 2,
      "name": "test-instruction-to-create-3-name",
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "description": "test-instruction-to-create-3-description",
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
      "description": "test-instruction-to-create-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "name": "test-instruction-to-create-3-name",
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
      "description": "test-instruction-to-create-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "name": "test-instruction-to-create-3-name",
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

  @concurrent
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

  @concurrent
  Scenario: given create auto instruction request with events count trigger should return ok
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "name": "test-instruction-to-create-5-name",
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "description": "test-instruction-to-create-5-description",
      "enabled": true,
      "triggers": [
        {
          "type": "eventscount",
          "threshold": 3
        }
      ],
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "jobs": [
        {
          "stop_on_fail": true,
          "job": "test-job-to-instruction-edit-1"
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
      "description": "test-instruction-to-create-5-description",
      "enabled": true,
      "triggers": [
        {
          "type": "eventscount",
          "threshold": 3
        }
      ],
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "name": "test-instruction-to-create-5-name",
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
      "description": "test-instruction-to-create-5-description",
      "enabled": true,
      "triggers": [
        {
          "type": "eventscount",
          "threshold": 3
        }
      ],
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "name": "test-instruction-to-create-5-name",
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
        }
      ]
    }
    """
