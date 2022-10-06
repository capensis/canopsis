Feature: update an instruction statistics
  I need to be able to update an instruction statistics


  Scenario: given get request should return assigned simplified manual instructions for the alarm
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-alarm-instruction-get-11",
      "connector_name": "test-connector-name-to-alarm-instruction-get-11",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-alarm-instruction-get-11",
      "resource": "test-resource-to-alarm-instruction-get-11",
      "state": 1,
      "output": "test-output-to-alarm-instruction-get-11"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-instruction-get-11&with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-to-alarm-instruction-get-11"
          },
          "assigned_instructions": [
            {
              "_id": "test-instruction-to-alarm-instruction-get-11-1",
              "name": "test-instruction-to-alarm-instruction-get-11-1-name",
              "execution": null
            },
            {
              "_id": "test-instruction-to-alarm-instruction-get-11-2",
              "name": "test-instruction-to-alarm-instruction-get-11-2-name",
              "execution": null
            }
          ]
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
