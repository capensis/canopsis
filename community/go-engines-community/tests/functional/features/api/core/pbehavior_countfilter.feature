Feature: create a PBehavior
  I need to be able to create a PBehavior
  Only admin should be able to create a PBehavior
  Wait 2s before query count to let apply chenges to MaxMatchedItems

  Scenario: POST a valid PBehavior but unauthorized
    When I do POST /api/v4/pbehaviors/count
    Then the response code should be 401

  Scenario: POST a valid PBehavior but without permissions
    When I am noperms
    When I do POST /api/v4/pbehaviors/count
    Then the response code should be 403

  Scenario: POST a valid PBehavior
    When I am admin
    When I do POST /api/v4/pbehaviors/count:
    """
    {
      "filter":{
        "$and":[
           {
              "name": {
                "$in": [
                    "test-case-pbehavior-countfilter-1-0-name",
                    "test-case-pbehavior-countfilter-1-1-name",
                    "test-case-pbehavior-countfilter-1-2-name",
                    "test-case-pbehavior-countfilter-1-3-name"
                ]
              }
           }
        ]
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "over_limit": false,
      "total_count": 4
    }
    """

  Scenario: POST a valid PBehavior
    When I am admin
    When I do POST /api/v4/pbehaviors/count:
    """
    {
      "filter":{
        "$and":[
           {
              "name": {
                "$in": [
                    "test-case-pbehavior-countfilter-1-0-name",
                    "test-case-pbehavior-countfilter-1-1-name",
                    "test-case-pbehavior-countfilter-1-2-name",
                    "test-case-pbehavior-countfilter-1-3-name",
                    "test-pbehavior-countfilter-component"
                ]
              }
           }
        ]
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "over_limit": true,
      "total_count": 5
    }
    """
