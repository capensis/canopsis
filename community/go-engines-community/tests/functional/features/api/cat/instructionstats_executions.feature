Feature: get instruction statistics
  I need to be able to get instruction statistics
  Only admin should be able to get instruction statistics

  Scenario: given get request should return instruction stats
    When I am admin
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-executions-get-1/executions
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "executed_on": 1618280210,
          "instruction_modified_on": 1596712203,
          "timeout_after_execution": {
            "value": 2,
            "unit": "s"
          },
          "status": 2,
          "alarm": {
            "_id": "test-alarm-to-stats-executions-get-1"
          },
          "duration": 600
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
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-executions-get-1/executions?from=1600000000&to=1700000000
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "executed_on": 1618280210,
          "instruction_modified_on": 1596712203,
          "timeout_after_execution": {
            "value": 2,
            "unit": "s"
          },
          "status": 2,
          "alarm": {
            "_id": "test-alarm-to-stats-executions-get-1"
          },
          "duration": 600
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
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-executions-get-1/executions?from=1500000000&to=1600000000
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-executions-get-1/executions?from=1300000000&to=1400000000
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """

  Scenario: given get request should return instruction stats with resolved alarms
    When I am admin
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-executions-get-with-resolved-alarms-1/executions
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "executed_on": 1618280210,
          "instruction_modified_on": 1596712203,
          "status": 2,
          "alarm": {
            "_id": "test-alarm-to-stats-executions-get-with-resolved-alarms-1"
          },
          "duration": 600
        },
        {
          "executed_on": 1518280210,
          "instruction_modified_on": 1596712203,
          "status": 2,
          "alarm": {
            "_id": "test-alarm-to-stats-executions-get-with-resolved-alarms-resolved"
          },
          "duration": 0
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
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-executions-get-with-resolved-alarms-1/executions?from=1500000000&to=1600000000
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "executed_on": 1518280210,
          "instruction_modified_on": 1596712203,
          "status": 2,
          "alarm": {
            "_id": "test-alarm-to-stats-executions-get-with-resolved-alarms-resolved"
          },
          "duration": 0
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
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-executions-get-with-resolved-alarms-1/executions?from=1600000000&to=1700000000
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "executed_on": 1618280210,
          "instruction_modified_on": 1596712203,
          "status": 2,
          "alarm": {
            "_id": "test-alarm-to-stats-executions-get-with-resolved-alarms-1"
          },
          "duration": 600
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

  Scenario: given get request and user without instruction create permission should return instruction stats
    When I am test-role-to-stats-instruction-get
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-executions-get-1/executions
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "executed_on": 1618280210,
          "instruction_modified_on": 1596712203,
          "status": 2,
          "alarm": {
            "_id": "test-alarm-to-stats-executions-get-1"
          },
          "duration": 600
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

  Scenario: given get request and user without instruction create permission should return not found error
    When I am test-role-to-stats-instruction-get
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-executions-get-2/executions
    Then the response code should be 404

  Scenario: given get request should return only corresponding alarm steps
    When I am admin
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-executions-get-3/executions
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "alarm": {
            "_id": "test-alarm-to-stats-executions-get-3",
            "v": {
              "steps": [
                {
                  "_t": "stateinc",
                  "t": 1618279510,
                  "val": 3
                },
                {
                  "_t": "comment",
                  "t": 1618279511
                },
                {
                  "_t": "instructionstart",
                  "t": 1618279610
                },
                {
                  "_t": "instructioncomplete",
                  "t": 1618280210
                },
                {
                  "_t": "ack",
                  "t": 1618280219
                },
                {
                  "_t": "statedec",
                  "t": 1618280220,
                  "val": 1
                }
              ]
            }
          },
          "executed_on": 1618280210,
          "instruction_modified_on": 1596712203,
          "status": 2,
          "duration": 600
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
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-executions-get-4/executions
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "alarm": {
            "_id": "test-alarm-to-stats-executions-get-4",
            "v": {
              "steps": [
                {
                  "_t": "stateinc",
                  "t": 1618279510,
                  "val": 3
                },
                {
                  "_t": "comment",
                  "t": 1618279511
                },
                {
                  "_t": "instructionstart",
                  "t": 1618279610
                },
                {
                  "_t": "instructioncomplete",
                  "t": 1618280210
                },
                {
                  "_t": "ack",
                  "t": 1618280220
                }
              ]
            }
          },
          "executed_on": 1618280210,
          "instruction_modified_on": 1596712203,
          "status": 2,
          "duration": 600
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
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-executions-get-5/executions
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "alarm": {
            "_id": "test-alarm-to-stats-executions-get-5",
            "v": {
              "steps": [
                {
                  "_t": "stateinc",
                  "t": 1618279510,
                  "val": 3
                },
                {
                  "_t": "instructionstart",
                  "t": 1618279610
                },
                {
                  "_t": "instructioncomplete",
                  "t": 1618280210
                },
                {
                  "_t": "ack",
                  "t": 1618280220
                },
                {
                  "_t": "instructionstart",
                  "t": 1618280230
                },
                {
                  "_t": "instructioncomplete",
                  "t": 1618280630
                },
                {
                  "_t": "comment",
                  "t": 1618280650
                }
              ]
            }
          },
          "executed_on": 1618280630,
          "instruction_modified_on": 1596712203,
          "status": 2,
          "duration": 400
        },
        {
          "alarm": {
            "_id": "test-alarm-to-stats-executions-get-5",
            "v": {
              "steps": [
                {
                  "_t": "stateinc",
                  "t": 1618279510,
                  "val": 3
                },
                {
                  "_t": "instructionstart",
                  "t": 1618279610
                },
                {
                  "_t": "instructioncomplete",
                  "t": 1618280210
                },
                {
                  "_t": "ack",
                  "t": 1618280220
                }
              ]
            }
          },
          "executed_on": 1618280210,
          "instruction_modified_on": 1596712203,
          "status": 2,
          "duration": 600
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
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-executions-get-6/executions
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "alarm": {
            "_id": "test-alarm-to-stats-executions-get-6",
            "v": {
              "steps": [
                {
                  "_t": "statedec",
                  "t": 1618280222,
                  "val": 2
                },
                {
                  "_t": "comment",
                  "t": 1618280223
                },
                {
                  "_t": "instructionstart",
                  "t": 1618280230
                },
                {
                  "_t": "instructioncomplete",
                  "t": 1618280630
                },
                {
                  "_t": "comment",
                  "t": 1618280640
                },
                {
                  "_t": "statedec",
                  "t": 1618280650,
                  "val": 1
                }
              ]
            }
          },
          "executed_on": 1618280630,
          "instruction_modified_on": 1596712203,
          "status": 2,
          "duration": 400
        },
        {
          "alarm": {
            "_id": "test-alarm-to-stats-executions-get-6",
            "v": {
              "steps": [
                {
                  "_t": "stateinc",
                  "t": 1618279510,
                  "val": 3
                },
                {
                  "_t": "instructionstart",
                  "t": 1618279610
                },
                {
                  "_t": "instructioncomplete",
                  "t": 1618280210
                },
                {
                  "_t": "ack",
                  "t": 1618280221
                },
                {
                  "_t": "statedec",
                  "t": 1618280222,
                  "val": 2
                }
              ]
            }
          },
          "executed_on": 1618280210,
          "instruction_modified_on": 1596712203,
          "status": 2,
          "duration": 600
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
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-executions-get-7/executions
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "alarm": {
            "_id": "test-alarm-to-stats-executions-get-7",
            "v": {
              "steps": [
                {
                  "_t": "stateinc",
                  "t": 1618279510,
                  "val": 3
                },
                {
                  "_t": "instructionstart",
                  "t": 1618279610
                },
                {
                  "_t": "instructioncomplete",
                  "t": 1618280210
                },
                {
                  "_t": "ack",
                  "t": 1618280220
                }
              ]
            }
          },
          "executed_on": 1618280210,
          "instruction_modified_on": 1596712203,
          "status": 2,
          "duration": 600
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
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-executions-get-8/executions
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "alarm": {
            "_id": "test-alarm-to-stats-executions-get-7",
            "v": {
              "steps": [
                {
                  "_t": "stateinc",
                  "t": 1618279510,
                  "val": 3
                },
                {
                  "_t": "instructionstart",
                  "t": 1618279610
                },
                {
                  "_t": "instructioncomplete",
                  "t": 1618280210
                },
                {
                  "_t": "ack",
                  "t": 1618280220
                },
                {
                  "_t": "instructionstart",
                  "t": 1618280230
                },
                {
                  "_t": "instructioncomplete",
                  "t": 1618280630
                },
                {
                  "_t": "comment",
                  "t": 1618280650
                }
              ]
            }
          },
          "executed_on": 1618280630,
          "instruction_modified_on": 1596712203,
          "status": 2,
          "duration": 400
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
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-executions-get-9/executions
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "alarm": {
            "_id": "test-alarm-to-stats-executions-get-8",
            "v": {
              "steps": [
                {
                  "_t": "stateinc",
                  "t": 1618279510,
                  "val": 3
                },
                {
                  "_t": "instructionstart",
                  "t": 1618279610
                },
                {
                  "_t": "instructioncomplete",
                  "t": 1618280210
                },
                {
                  "_t": "ack",
                  "t": 1618280220
                },
                {
                  "_t": "statedec",
                  "val": 2
                }
              ]
            }
          },
          "executed_on": 1618280210,
          "instruction_modified_on": 1596712203,
          "status": 2,
          "duration": 600
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
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-executions-get-10/executions
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "alarm": {
            "_id": "test-alarm-to-stats-executions-get-8",
            "v": {
              "steps": [
                {
                  "_t": "statedec",
                  "t": 1618280222,
                  "val": 2
                },
                {
                  "_t": "comment",
                  "t": 1618280223
                },
                {
                  "_t": "instructionstart",
                  "t": 1618280230
                },
                {
                  "_t": "instructioncomplete",
                  "t": 1618280630
                },
                {
                  "_t": "comment",
                  "t": 1618280640
                },
                {
                  "_t": "statedec",
                  "t": 1618280650,
                  "val": 1
                }
              ]
            }
          },
          "executed_on": 1618280630,
          "instruction_modified_on": 1596712203,
          "status": 2,
          "duration": 400
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
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-executions-get-11/executions?show_failed=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "alarm": {
            "_id": "test-alarm-to-stats-executions-get-11-2",
            "v": {
              "steps": [
                {
                  "_t": "stateinc",
                  "t": 1596942720
                },
                {
                  "_t": "statusinc",
                  "t": 1596942720
                },
                {
                  "_t": "instructionstart",
                  "t": 1596942720
                },
                {
                  "_t": "instructioncomplete",
                  "t": 1596942720
                }
              ]
            }
          },
          "duration": 0,
          "status": 4,
          "executed_on": 1618280220,
          "instruction_modified_on": 1596712203
        },
        {
          "alarm": {
            "_id": "test-alarm-to-stats-executions-get-11-1",
            "v": {
              "steps": [
                {
                  "_t": "stateinc",
                  "t": 1618280213
                },
                {
                  "_t": "statusinc",
                  "t": 1618280213
                },
                {
                  "_t": "instructionstart",
                  "t": 1618280213
                },
                {
                  "_t": "instructioncomplete",
                  "t": 1618280218
                }
              ]
            }
          },
          "duration": 5,
          "status": 2,
          "executed_on": 1618280218,
          "instruction_modified_on": 1596712203
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
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-get-1/executions
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
          {
              "duration": 200,
              "status": 2,
              "executed_on": 1618394399,
              "instruction_modified_on": 1617555600
          },
          {
              "duration": 350,
              "status": 2,
              "executed_on": 1618307999,
              "instruction_modified_on": 1617555600
          },
          {
              "duration": 400,
              "status": 2,
              "executed_on": 1618221599,
              "instruction_modified_on": 1617555600
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

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/cat/instruction-stats/notexist/executions
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/instruction-stats/notexist/executions
    Then the response code should be 403

  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/cat/instruction-stats/notexist/executions?from=1500000000&to=1700000000
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  Scenario: given get search request should return instruction stats
    When I am admin
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-executions-get-with-resolved-alarms-1/executions?search=resolved-alarms-resolved
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "executed_on": 1518280210,
          "alarm": {
            "_id": "test-alarm-to-stats-executions-get-with-resolved-alarms-resolved"
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
