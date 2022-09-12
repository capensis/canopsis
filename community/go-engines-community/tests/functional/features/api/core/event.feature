Feature: send an event
  I need to be able to send an event
  Only admin should be able to send an event

  Scenario: POST a valid event but unauthorized
    When I do POST /api/v4/event
    Then the response code should be 401

  Scenario: POST a valid event but without permissions
    When I am noperms
    When I do POST /api/v4/event
    Then the response code should be 403

  Scenario: POST a valid event
    When I am admin
    When I do POST /api/v4/event:
    """json
    {
      "connector": "test_connector",
      "connector_name": "test_connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component": "test_component",
      "state": 1,
      "resource": "test_resource"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "sent_events": [
        {
          "author": "root",
          "connector": "test_connector",
          "connector_name": "test_connectorname",
          "source_type": "resource",
          "event_type": "check",
          "component": "test_component",
          "state": 1,
          "resource": "test_resource"
        }
      ],
      "failed_events": [],
      "retry_events": []
    }
    """

  Scenario: POST a valid events in array
    When I am admin
    When I do POST /api/v4/event:
    """json
    [
      {
        "connector": "test_connector",
        "connector_name": "test_connectorname",
        "source_type": "resource",
        "event_type": "check",
        "component": "test_component_array",
        "state": 1,
        "resource": "test_resource_array"
      },
      {
        "author": "test_author",
        "connector": "test_connector",
        "connector_name": "test_connectorname",
        "source_type": "resource",
        "event_type": "check",
        "component": "test_component_array",
        "state": 1,
        "resource": "test_resource_array2"
      }
    ]
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "sent_events": [
        {
          "author": "root",
          "connector": "test_connector",
          "connector_name": "test_connectorname",
          "source_type": "resource",
          "event_type": "check",
          "component": "test_component_array",
          "state": 1,
          "resource": "test_resource_array"
        },
        {
          "author": "test_author",
          "connector": "test_connector",
          "connector_name": "test_connectorname",
          "source_type": "resource",
          "event_type": "check",
          "component": "test_component_array",
          "state": 1,
          "resource": "test_resource_array2"
        }
      ],
      "failed_events": [],
      "retry_events": []
    }
    """

  Scenario: POST an invalid event without event_type
    When I am admin
    When I do POST /api/v4/event:
    """json
    {
      "connector": "test_connector",
      "connector_name": "test_connectorname",
      "source_type": "resource",
      "component": "test_component",
      "state": 1,
      "resource": "test_resource"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "sent_events": [],
      "failed_events": [
        {
          "connector": "test_connector",
          "connector_name": "test_connectorname",
          "source_type": "resource",
          "component": "test_component",
          "state": 1,
          "resource": "test_resource"
        }
      ],
      "retry_events": []
    }
    """

  Scenario: POST an invalid event when event_type is not string
    When I am admin
    When I do POST /api/v4/event:
    """json
    {
      "connector": "test_connector",
      "connector_name": "test_connectorname",
      "source_type": "resource",
      "component": "test_component",
      "event_type": 123,
      "state": 1,
      "resource": "test_resource"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "sent_events": [],
      "failed_events": [
        {
          "connector": "test_connector",
          "connector_name": "test_connectorname",
          "source_type": "resource",
          "component": "test_component",
          "event_type": 123,
          "state": 1,
          "resource": "test_resource"
        }
      ],
      "retry_events": []
    }
    """

  Scenario: POST an invalid event without state
    When I am admin
    When I do POST /api/v4/event:
    """json
    {
      "connector": "test_connector",
      "connector_name": "test_connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component": "test_component",
      "resource": "test_resource"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "sent_events": [],
      "failed_events": [
        {
          "connector": "test_connector",
          "connector_name": "test_connectorname",
          "source_type": "resource",
          "event_type": "check",
          "component": "test_component",
          "resource": "test_resource"
        }
      ],
      "retry_events": []
    }
    """

  Scenario: POST an invalid event when state is not an int
    When I am admin
    When I do POST /api/v4/event:
    """json
    {
      "connector": "test_connector",
      "connector_name": "test_connectorname",
      "source_type": "resource",
      "event_type": "check",
      "state": "abc",
      "component": "test_component",
      "resource": "test_resource"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "sent_events": [],
      "failed_events": [
        {
          "connector": "test_connector",
          "connector_name": "test_connectorname",
          "source_type": "resource",
          "event_type": "check",
          "state": "abc",
          "component": "test_component",
          "resource": "test_resource"
        }
      ],
      "retry_events": []
    }
    """

  Scenario: POST valid and invalid events in array
    When I am admin
    When I do POST /api/v4/event:
    """json
    [
      {
        "connector": "test_connector",
        "connector_name": "test_connectorname",
        "source_type": "resource",
        "event_type": "check",
        "component": "test_component_array",
        "state": 1,
        "resource": "test_resource_array3"
      },
      {
        "connector": "test_connector",
        "connector_name": "test_connectorname",
        "source_type": "resource",
        "component": "test_component_array",
        "state": 1,
        "resource": "test_resource_array4"
      }
    ]
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "sent_events": [
        {
          "author": "root",
          "connector": "test_connector",
          "connector_name": "test_connectorname",
          "source_type": "resource",
          "event_type": "check",
          "component": "test_component_array",
          "state": 1,
          "resource": "test_resource_array3"
        }
      ],
      "failed_events": [
        {
          "connector": "test_connector",
          "connector_name": "test_connectorname",
          "source_type": "resource",
          "component": "test_component_array",
          "state": 1,
          "resource": "test_resource_array4"
        }
      ],
      "retry_events": []
    }
    """

  Scenario: POST a valid event when long_output is not string should transform to empty string
    When I am admin
    When I do POST /api/v4/event:
    """json
    {
      "connector": "test_connector",
      "connector_name": "test_connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component": "test_component",
      "long_output": 123,
      "state": 1,
      "resource": "test_resource2"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "sent_events": [
        {
          "connector": "test_connector",
          "connector_name": "test_connectorname",
          "source_type": "resource",
          "event_type": "check",
          "component": "test_component",
          "state": 1,
          "long_output": "",
          "resource": "test_resource2",
          "author": "root"
        }
      ],
      "failed_events": [],
      "retry_events": []
    }
    """

  Scenario: POST a valid changestate event should add role to event
    When I am admin
    When I do POST /api/v4/event:
    """json
    {
      "connector": "test_connector",
      "connector_name": "test_connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component": "test_component2",
      "state": 2,
      "resource": "test_resource"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "sent_events": [
        {
          "connector": "test_connector",
          "connector_name": "test_connectorname",
          "source_type": "resource",
          "event_type": "check",
          "component": "test_component2",
          "state": 2,
          "resource": "test_resource",
          "author": "root"
        }
      ],
      "failed_events": [],
      "retry_events": []
    }
    """
    When I do POST /api/v4/event:
    """json
    {
      "connector": "test_connector",
      "connector_name": "test_connectorname",
      "source_type": "resource",
      "event_type": "changestate",
      "component": "test_component2",
      "state": 3,
      "resource": "test_resource"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "sent_events": [
        {
          "author": "root",
          "connector": "test_connector",
          "connector_name": "test_connectorname",
          "source_type": "resource",
          "event_type": "changestate",
          "component": "test_component2",
          "state": 3,
          "role": "admin",
          "resource": "test_resource"
        }
      ],
      "failed_events": [],
      "retry_events": []
    }
    """

  Scenario: POST a valid event with urlencoded content-type
    When I am admin
    When I set header Content-Type=application/x-www-form-urlencoded
    When I do POST /api/v4/event:
    """
    event_type=check&connector=computer24&connector_name=computer24&component=phone&resource=ram&source_type=resource&author=superviseur1&state=2&output=canopsis
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "sent_events": [
        {
          "author": "superviseur1",
          "component": "phone",
          "connector": "computer24",
          "connector_name": "computer24",
          "event_type": "check",
          "output": "canopsis",
          "resource": "ram",
          "source_type": "resource",
          "state": 2
        }
      ],
      "failed_events": [],
      "retry_events": []
    }
    """

  Scenario: POST a valid event with extended urlencoded content-type
    When I am admin
    When I set header Content-Type=application/x-www-form-urlencoded; charset=utf-8
    When I do POST /api/v4/event:
    """
    event_type=check&connector=computer25&connector_name=computer25&component=phone2&resource=cpu&source_type=resource&author=superviseur2&state=2&output=canopsis
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "sent_events": [
        {
          "author": "superviseur2",
          "component": "phone2",
          "connector": "computer25",
          "connector_name": "computer25",
          "event_type": "check",
          "output": "canopsis",
          "resource": "cpu",
          "source_type": "resource",
          "state": 2
        }
      ],
      "failed_events": [],
      "retry_events": []
    }
    """

  Scenario: POST an invalid state event with extended urlencoded content-type
    When I am admin
    When I set header Content-Type=application/x-www-form-urlencoded; charset=utf-8
    When I do POST /api/v4/event:
    """
    event_type=check&connector=computer25&connector_name=computer25&component=phone2&resource=cpu&source_type=resource&author=superviseur2&state=abc&output=canopsis
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "failed_events": [
        {
          "author": "superviseur2",
          "component": "phone2",
          "connector": "computer25",
          "connector_name": "computer25",
          "event_type": "check",
          "output": "canopsis",
          "resource": "cpu",
          "source_type": "resource",
          "state": "abc"
        }
      ],
      "sent_events": [],
      "retry_events": []
    }
    """

  Scenario: POST an invalid event when connector or connector_name is empty
    When I am admin
    When I do POST /api/v4/event:
    """json
    [
      {
        "connector_name": "test_connectorname",
        "source_type": "resource",
        "component": "test_component",
        "event_type": "check",
        "state": 1,
        "resource": "test_resource"
      },
      {
        "connector": "test_connector",
        "source_type": "resource",
        "component": "test_component",
        "event_type": "check",
        "state": 1,
        "resource": "test_resource"
      }
    ]
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "sent_events": [],
      "failed_events": [
        {
          "connector_name": "test_connectorname",
          "source_type": "resource",
          "component": "test_component",
          "event_type": "check",
          "state": 1,
          "resource": "test_resource",
          "author": "root"
        },
        {
          "connector": "test_connector",
          "source_type": "resource",
          "component": "test_component",
          "event_type": "check",
          "state": 1,
          "resource": "test_resource",
          "author": "root"
        }
      ],
      "retry_events": []
    }
    """

  Scenario: POST an invalid event when connector or connector_name or component or resource is not string
    When I am admin
    When I do POST /api/v4/event:
    """json
    [
      {
        "connector": 123,
        "connector_name": "test_connectorname",
        "source_type": "resource",
        "component": "test_component",
        "event_type": "check",
        "state": 1,
        "resource": "test_resource"
      },
      {
        "connector": "test_connector",
        "connector_name": 123,
        "source_type": "resource",
        "component": "test_component",
        "event_type": "check",
        "state": 1,
        "resource": "test_resource"
      },
      {
        "connector": "test_connector",
        "connector_name": "test_connectorname",
        "source_type": "resource",
        "component": 123,
        "event_type": "check",
        "state": 1,
        "resource": "test_resource"
      },
      {
        "connector": "test_connector",
        "connector_name": "test_connectorname",
        "source_type": "resource",
        "component": "test_component",
        "event_type": "check",
        "state": 1,
        "resource": 123
      }
    ]
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "sent_events": [],
      "failed_events": [
        {
          "connector": 123,
          "connector_name": "test_connectorname",
          "source_type": "resource",
          "component": "test_component",
          "event_type": "check",
          "state": 1,
          "resource": "test_resource",
          "author": "root"
        },
        {
          "connector": "test_connector",
          "connector_name": 123,
          "source_type": "resource",
          "component": "test_component",
          "event_type": "check",
          "state": 1,
          "resource": "test_resource",
          "author": "root"
        },
        {
          "connector": "test_connector",
          "connector_name": "test_connectorname",
          "source_type": "resource",
          "component": 123,
          "event_type": "check",
          "state": 1,
          "resource": "test_resource",
          "author": "root"
        },
        {
          "connector": "test_connector",
          "connector_name": "test_connectorname",
          "source_type": "resource",
          "component": "test_component",
          "event_type": "check",
          "state": 1,
          "resource": 123,
          "author": "root"
        }
      ],
      "retry_events": []
    }
    """

  Scenario: POST an event when source_type is empty
    When I am admin
    When I do POST /api/v4/event:
    """json
    [
      {
        "connector": "test_connector",
        "connector_name": "test_connectorname",
        "event_type": "ack"
      },
      {
        "connector": "test_connector",
        "connector_name": "test_connectorname",
        "component": "test_component",
        "event_type": "check",
        "state": 1
      },
      {
        "connector": "test_connector",
        "connector_name": "test_connectorname",
        "component": "test_component",
        "event_type": "check",
        "state": 1,
        "resource": "test_resource"
      }
    ]
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "sent_events": [
        {
          "connector": "test_connector",
          "connector_name": "test_connectorname",
          "event_type": "ack",
          "role": "admin",
          "source_type": "connector",
          "author": "root"
        },
        {
          "connector": "test_connector",
          "connector_name": "test_connectorname",
          "component": "test_component",
          "event_type": "check",
          "source_type": "component",
          "state": 1,
          "author": "root"
        },
        {
          "connector": "test_connector",
          "connector_name": "test_connectorname",
          "component": "test_component",
          "event_type": "check",
          "source_type": "resource",
          "state": 1,
          "resource": "test_resource",
          "author": "root"
        }
      ],
      "failed_events": [],
      "retry_events": []
    }
    """

  Scenario: POST an event when source_type is invalid
    When I am admin
    When I do POST /api/v4/event:
    """json
    [
      {
        "connector": "test_connector",
        "connector_name": "test_connectorname",
        "event_type": "ack",
        "source_type": "resource"
      },
      {
        "connector": "test_connector",
        "connector_name": "test_connectorname",
        "component": "test_component",
        "event_type": "check",
        "state": 1,
        "source_type": "connector"
      },
      {
        "connector": "test_connector",
        "connector_name": "test_connectorname",
        "component": "test_component",
        "event_type": "check",
        "state": 1,
        "resource": "test_resource",
        "source_type": "connector"
      }
    ]
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "sent_events": [
        {
          "connector": "test_connector",
          "connector_name": "test_connectorname",
          "event_type": "ack",
          "role": "admin",
          "source_type": "connector",
          "author": "root"
        },
        {
          "connector": "test_connector",
          "connector_name": "test_connectorname",
          "component": "test_component",
          "event_type": "check",
          "source_type": "component",
          "state": 1,
          "author": "root"
        },
        {
          "connector": "test_connector",
          "connector_name": "test_connectorname",
          "component": "test_component",
          "event_type": "check",
          "source_type": "resource",
          "state": 1,
          "resource": "test_resource",
          "author": "root"
        }
      ],
      "failed_events": [],
      "retry_events": []
    }
    """

  Scenario: POST an invalid event when event_type=check and source_type=connector
    When I am admin
    When I do POST /api/v4/event:
    """json
    [
      {
        "connector": "test_connector",
        "connector_name": "test_connectorname",
        "event_type": "check",
        "state": 1
      }
    ]
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "sent_events": [],
      "failed_events": [
        {
          "connector": "test_connector",
          "connector_name": "test_connectorname",
          "event_type": "check",
          "state": 1,
          "author": "root"
        }
      ],
      "retry_events": []
    }
    """
