Feature: update a pbehavior
  I need to be able to patch a pbehavior field individually
  Only admin should be able to patch a pbehavior

  @concurrent
  Scenario: given update request with start, stop and rrule should update rrule end
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-third-1:
    """json
    {
      "rrule": "FREQ=DAILY;UNTIL=20221108T103000Z",
      "tstart": {{ parseTimeTz "08-10-2022 10:00" }},
      "tstop": {{ parseTimeTz "08-10-2022 11:00" }}
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "rrule_end": {{ parseTimeTz "08-11-2022 10:00" }}
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-third-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "rrule_end": {{ parseTimeTz "08-11-2022 10:00" }}
    }
    """
