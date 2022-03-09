Feature: get test suite's gantt intervals
  I need to be able to get test suite's gantt intervals
  Only admin should be able to get test suite's gantt intervals

  Scenario: GET unauthorized
    When I do GET /api/v4/cat/junit/test-suites/test-suite-to-get-gantt-1-3/gantt
    Then the response code should be 401

  Scenario: GET without permissions
    When I am noperms
    When I do GET /api/v4/cat/junit/test-suites/test-suite-to-get-gantt-1-3/gantt
    Then the response code should be 403

  Scenario: GET gantt intervals success
    When I am admin
    When I do GET /api/v4/cat/junit/test-suites/test-suite-to-get-gantt-1-3/gantt?page=1&limit=10
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "from": 0,
          "to": 1.5,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-1",
          "message": "test case 1 message",
          "name": "test case 1",
          "time": 1.5,
          "avg_to": 0,
          "avg_time": 0,
          "avg_status": 0
        },
        {
          "from": 1.5,
          "to": 2.5,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-2",
          "message": "test case 2 message",
          "name": "test case 2",
          "time": 1,
          "avg_to": 0,
          "avg_time": 0,
          "avg_status": 0
        },
        {
          "from": 2.5,
          "to": 4.5,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-3",
          "message": "test case 3 message",
          "name": "test case 3",
          "time": 2,
          "avg_to": 0,
          "avg_time": 0,
          "avg_status": 0
        },
        {
          "from": 4.5,
          "to": 6.1,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-4",
          "message": "test case 4 message",
          "name": "test case 4",
          "time": 1.6,
          "avg_to": 0,
          "avg_time": 0,
          "avg_status": 0
        },
        {
          "from": 6.1,
          "to": 7.5,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-5",
          "message": "test case 5 message",
          "name": "test case 5",
          "time": 1.4,
          "avg_to": 0,
          "avg_time": 0,
          "avg_status": 0
        },
        {
          "from": 7.5,
          "to": 10,
          "status": 3,
          "_id": "test-suite-to-get-gantt-1-test-case-6",
          "message": "test case 6 message",
          "name": "test case 6",
          "time": 2.5,
          "avg_to": 0,
          "avg_time": 0,
          "avg_status": 0
        },
        {
          "from": 10,
          "to": 10.5,
          "status": 2,
          "_id": "test-suite-to-get-gantt-1-test-case-7",
          "message": "test case 7 message",
          "name": "test case 7",
          "time": 0.5,
          "avg_to": 0,
          "avg_time": 0,
          "avg_status": 0
        },
        {
          "from": 10.5,
          "to": 11,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-8",
          "message": "test case 8 message",
          "name": "test case 8",
          "time": 0.5,
          "avg_to": 0,
          "avg_time": 0,
          "avg_status": 0
        },
        {
          "from": 11,
          "to": 11,
          "status": 1,
          "_id": "test-suite-to-get-gantt-1-test-case-9",
          "message": "test case 9 message",
          "name": "test case 9",
          "time": 0,
          "avg_to": 0,
          "avg_time": 0,
          "avg_status": 0
        },
        {
          "from": 11,
          "to": 11.5,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-10",
          "message": "test case 10 message",
          "name": "test case 10",
          "time": 0.5,
          "avg_to": 0,
          "avg_time": 0,
          "avg_status": 0
        }
      ],
      "meta": {
        "time": 20,
        "avg_time": 0,
        "page": 1,
        "per_page": 10,
        "page_count": 2,
        "total_count": 15
      }
    }
    """
    When I do GET /api/v4/cat/junit/test-suites/test-suite-to-get-gantt-1-3/gantt?page=2&limit=10
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "from": 11.5,
          "to": 13,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-11",
          "message": "test case 11 message",
          "name": "test case 11",
          "time": 1.5,
          "avg_to": 0,
          "avg_time": 0,
          "avg_status": 0
        },
        {
          "from": 13,
          "to": 16.5,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-12",
          "message": "test case 12 message",
          "name": "test case 12",
          "time": 3.5,
          "avg_to": 0,
          "avg_time": 0,
          "avg_status": 0
        },
        {
          "from": 16.5,
          "to": 18,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-13",
          "message": "test case 13 message",
          "name": "test case 13",
          "time": 1.5,
          "avg_to": 0,
          "avg_time": 0,
          "avg_status": 0
        },
        {
          "from": 18,
          "to": 18.5,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-14",
          "message": "test case 14 message",
          "name": "test case 14",
          "time": 0.5,
          "avg_to": 0,
          "avg_time": 0,
          "avg_status": 0
        },
        {
          "from": 18.5,
          "to": 20,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-15",
          "message": "test case 15 message",
          "name": "test case 15",
          "time": 1.5,
          "avg_to": 0,
          "avg_time": 0,
          "avg_status": 0
        }
      ],
      "meta": {
        "time": 20,
        "avg_time": 0,
        "page": 2,
        "per_page": 10,
        "page_count": 2,
        "total_count": 15
      }
    }
    """

  Scenario: GET gantt intervals with history for 6 months, last suite
    When I am admin
    When I do GET /api/v4/cat/junit/test-suites/test-suite-to-get-gantt-1-3/gantt?page=1&limit=10&months=6
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "from": 0,
          "to": 1.5,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-1",
          "message": "test case 1 message",
          "name": "test case 1",
          "time": 1.5,
          "avg_to": 1.5,
          "avg_time": 1.5,
          "avg_status": 1
        },
        {
          "from": 1.5,
          "to": 2.5,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-2",
          "message": "test case 2 message",
          "name": "test case 2",
          "time": 1,
          "avg_to": 4.5,
          "avg_time": 3,
          "avg_status": 0
        },
        {
          "from": 2.5,
          "to": 4.5,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-3",
          "message": "test case 3 message",
          "name": "test case 3",
          "time": 2,
          "avg_to": 6.5,
          "avg_time": 4,
          "avg_status": 0
        },
        {
          "from": 4.5,
          "to": 6.1,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-4",
          "message": "test case 4 message",
          "name": "test case 4",
          "time": 1.6,
          "avg_to": 6.43,
          "avg_time": 1.93,
          "avg_status": 0
        },
        {
          "from": 6.1,
          "to": 7.5,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-5",
          "message": "test case 5 message",
          "name": "test case 5",
          "time": 1.4,
          "avg_to": 9.17,
          "avg_time": 3.07,
          "avg_status": 2
        },
        {
          "from": 7.5,
          "to": 10,
          "status": 3,
          "_id": "test-suite-to-get-gantt-1-test-case-6",
          "message": "test case 6 message",
          "name": "test case 6",
          "time": 2.5,
          "avg_to": 11.67,
          "avg_time": 4.17,
          "avg_status": 3
        },
        {
          "from": 10,
          "to": 10.5,
          "status": 2,
          "_id": "test-suite-to-get-gantt-1-test-case-7",
          "message": "test case 7 message",
          "name": "test case 7",
          "time": 0.5,
          "avg_to": 10.83,
          "avg_time": 0.83,
          "avg_status": 2
        },
        {
          "from": 10.5,
          "to": 11,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-8",
          "message": "test case 8 message",
          "name": "test case 8",
          "time": 0.5,
          "avg_to": 11.33,
          "avg_time": 0.83,
          "avg_status": 0
        },
        {
          "from": 11,
          "to": 11,
          "status": 1,
          "_id": "test-suite-to-get-gantt-1-test-case-9",
          "message": "test case 9 message",
          "name": "test case 9",
          "time": 0,
          "avg_to": 11.5,
          "avg_time": 0.5,
          "avg_status": 1
        },
        {
          "from": 11,
          "to": 11.5,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-10",
          "message": "test case 10 message",
          "name": "test case 10",
          "time": 0.5,
          "avg_to": 11.53,
          "avg_time": 0.53,
          "avg_status": 3
        }
      ],
      "meta": {
        "time": 20,
        "avg_time": 30.37,
        "page": 1,
        "per_page": 10,
        "page_count": 2,
        "total_count": 15
      }
    }
    """
    When I do GET /api/v4/cat/junit/test-suites/test-suite-to-get-gantt-1-3/gantt?page=2&limit=10&months=6
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "from": 11.5,
          "to": 13,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-11",
          "message": "test case 11 message",
          "name": "test case 11",
          "time": 1.5,
          "avg_to": 13.5,
          "avg_time": 2,
          "avg_status": 0
        },
        {
          "from": 13,
          "to": 16.5,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-12",
          "message": "test case 12 message",
          "name": "test case 12",
          "time": 3.5,
          "avg_to": 18.25,
          "avg_time": 5.25,
          "avg_status": 0
        },
        {
          "from": 16.5,
          "to": 18,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-13",
          "message": "test case 13 message",
          "name": "test case 13",
          "time": 1.5,
          "avg_to": 19,
          "avg_time": 2.5,
          "avg_status": 0
        },
        {
          "from": 18,
          "to": 18.5,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-14",
          "message": "test case 14 message",
          "name": "test case 14",
          "time": 0.5,
          "avg_to": 18.83,
          "avg_time": 0.83,
          "avg_status": 0
        },
        {
          "from": 18.5,
          "to": 20,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-15",
          "message": "test case 15 message",
          "name": "test case 15",
          "time": 1.5,
          "avg_to": 21,
          "avg_time": 2.5,
          "avg_status": 0
        }
      ],
      "meta": {
        "time": 20,
        "avg_time": 30.37,
        "page": 2,
        "per_page": 10,
        "page_count": 2,
        "total_count": 15
      }
    }
    """

  Scenario: GET gantt intervals with history for 3 months, last suite
    When I am admin
    When I do GET /api/v4/cat/junit/test-suites/test-suite-to-get-gantt-1-3/gantt?page=1&limit=10&months=3
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "from": 0,
          "to": 1.5,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-1",
          "message": "test case 1 message",
          "name": "test case 1",
          "time": 1.5,
          "avg_to": 1.5,
          "avg_time": 1.5,
          "avg_status": 0
        },
        {
          "from": 1.5,
          "to": 2.5,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-2",
          "message": "test case 2 message",
          "name": "test case 2",
          "time": 1,
          "avg_to": 4.5,
          "avg_time": 3,
          "avg_status": 0
        },
        {
          "from": 2.5,
          "to": 4.5,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-3",
          "message": "test case 3 message",
          "name": "test case 3",
          "time": 2,
          "avg_to": 5.5,
          "avg_time": 3,
          "avg_status": 0
        },
        {
          "from": 4.5,
          "to": 6.1,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-4",
          "message": "test case 4 message",
          "name": "test case 4",
          "time": 1.6,
          "avg_to": 5.9,
          "avg_time": 1.4,
          "avg_status": 0
        },
        {
          "from": 6.1,
          "to": 7.5,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-5",
          "message": "test case 5 message",
          "name": "test case 5",
          "time": 1.4,
          "avg_to": 9.2,
          "avg_time": 3.1,
          "avg_status": 0
        },
        {
          "from": 7.5,
          "to": 10,
          "status": 3,
          "_id": "test-suite-to-get-gantt-1-test-case-6",
          "message": "test case 6 message",
          "name": "test case 6",
          "time": 2.5,
          "avg_to": 11.25,
          "avg_time": 3.75,
          "avg_status": 3
        },
        {
          "from": 10,
          "to": 10.5,
          "status": 2,
          "_id": "test-suite-to-get-gantt-1-test-case-7",
          "message": "test case 7 message",
          "name": "test case 7",
          "time": 0.5,
          "avg_to": 10.75,
          "avg_time": 0.75,
          "avg_status": 2
        },
        {
          "from": 10.5,
          "to": 11,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-8",
          "message": "test case 8 message",
          "name": "test case 8",
          "time": 0.5,
          "avg_to": 11.25,
          "avg_time": 0.75,
          "avg_status": 0
        },
        {
          "from": 11,
          "to": 11,
          "status": 1,
          "_id": "test-suite-to-get-gantt-1-test-case-9",
          "message": "test case 9 message",
          "name": "test case 9",
          "time": 0,
          "avg_to": 11.5,
          "avg_time": 0.5,
          "avg_status": 0
        },
        {
          "from": 11,
          "to": 11.5,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-10",
          "message": "test case 10 message",
          "name": "test case 10",
          "time": 0.5,
          "avg_to": 11.75,
          "avg_time": 0.75,
          "avg_status": 0
        }
      ],
      "meta": {
        "time": 20,
        "avg_time": 26.5,
        "page": 1,
        "per_page": 10,
        "page_count": 2,
        "total_count": 15
      }
    }
    """
    When I do GET /api/v4/cat/junit/test-suites/test-suite-to-get-gantt-1-3/gantt?page=2&limit=10&months=3
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "from": 11.5,
          "to": 13,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-11",
          "message": "test case 11 message",
          "name": "test case 11",
          "time": 1.5,
          "avg_to": 13.5,
          "avg_time": 2,
          "avg_status": 0
        },
        {
          "from": 13,
          "to": 16.5,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-12",
          "message": "test case 12 message",
          "name": "test case 12",
          "time": 3.5,
          "avg_to": 16.5,
          "avg_time": 3.5,
          "avg_status": 0
        },
        {
          "from": 16.5,
          "to": 18,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-13",
          "message": "test case 13 message",
          "name": "test case 13",
          "time": 1.5,
          "avg_to": 18.75,
          "avg_time": 2.25,
          "avg_status": 0
        },
        {
          "from": 18,
          "to": 18.5,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-14",
          "message": "test case 14 message",
          "name": "test case 14",
          "time": 0.5,
          "avg_to": 18.75,
          "avg_time": 0.75,
          "avg_status": 0
        },
        {
          "from": 18.5,
          "to": 20,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-15",
          "message": "test case 15 message",
          "name": "test case 15",
          "time": 1.5,
          "avg_to": 20.75,
          "avg_time": 2.25,
          "avg_status": 0
        }
      ],
      "meta": {
        "time": 20,
        "avg_time": 26.5,
        "page": 2,
        "per_page": 10,
        "page_count": 2,
        "total_count": 15
      }
    }
    """

  Scenario: GET gantt intervals with history for 3 months, middle suite
    When I am admin
    When I do GET /api/v4/cat/junit/test-suites/test-suite-to-get-gantt-1-2/gantt?page=1&limit=10&months=3
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "from": 0,
          "to": 0,
          "status": 1,
          "_id": "test-suite-to-get-gantt-1-test-case-1",
          "message": "test case 1 message",
          "name": "test case 1",
          "time": 0,
          "avg_to": 0,
          "avg_time": 0,
          "avg_status": 1
        },
        {
          "from": 0,
          "to": 5,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-2",
          "message": "test case 2 message",
          "name": "test case 2",
          "time": 5,
          "avg_to": 4,
          "avg_time": 4,
          "avg_status": 0
        },
        {
          "from": 5,
          "to": 9,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-3",
          "message": "test case 3 message",
          "name": "test case 3",
          "time": 4,
          "avg_to": 10,
          "avg_time": 5,
          "avg_status": 0
        },
        {
          "from": 9,
          "to": 10.2,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-4",
          "message": "test case 4 message",
          "name": "test case 4",
          "time": 1.2,
          "avg_to": 11.1,
          "avg_time": 2.1,
          "avg_status": 0
        },
        {
          "from": 10.2,
          "to": 15,
          "status": 2,
          "_id": "test-suite-to-get-gantt-1-test-case-5",
          "message": "test case 5 message",
          "name": "test case 5",
          "time": 4.8,
          "avg_to": 14.1,
          "avg_time": 3.9,
          "avg_status": 2
        },
        {
          "from": 15,
          "to": 20,
          "status": 3,
          "_id": "test-suite-to-get-gantt-1-test-case-6",
          "message": "test case 6 message",
          "name": "test case 6",
          "time": 5,
          "avg_to": 20,
          "avg_time": 5,
          "avg_status": 3
        },
        {
          "from": 20,
          "to": 21,
          "status": 2,
          "_id": "test-suite-to-get-gantt-1-test-case-7",
          "message": "test case 7 message",
          "name": "test case 7",
          "time": 1,
          "avg_to": 21,
          "avg_time": 1,
          "avg_status": 2
        },
        {
          "from": 21,
          "to": 22,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-8",
          "message": "test case 8 message",
          "name": "test case 8",
          "time": 1,
          "avg_to": 22,
          "avg_time": 1,
          "avg_status": 0
        },
        {
          "from": 22,
          "to": 22.5,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-9",
          "message": "test case 9 message",
          "name": "test case 9",
          "time": 0.5,
          "avg_to": 22.5,
          "avg_time": 0.5,
          "avg_status": 0
        },
        {
          "from": 22.5,
          "to": 23.5,
          "status": 3,
          "_id": "test-suite-to-get-gantt-1-test-case-10",
          "message": "test case 10 message",
          "name": "test case 10",
          "time": 1,
          "avg_to": 23.05,
          "avg_time": 0.55,
          "avg_status": 3
        }
      ],
      "meta": {
        "time": 33,
        "avg_time": 35.55,
        "page": 1,
        "per_page": 10,
        "page_count": 2,
        "total_count": 15
      }
    }
    """
    When I do GET /api/v4/cat/junit/test-suites/test-suite-to-get-gantt-1-2/gantt?page=2&limit=10&months=3
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "from": 23.5,
          "to": 26,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-11",
          "message": "test case 11 message",
          "name": "test case 11",
          "time": 2.5,
          "avg_to": 25.75,
          "avg_time": 2.25,
          "avg_status": 0
        },
        {
          "from": 26,
          "to": 26,
          "status": 1,
          "_id": "test-suite-to-get-gantt-1-test-case-12",
          "message": "test case 12 message",
          "name": "test case 12",
          "time": 0,
          "avg_to": 33,
          "avg_time": 7,
          "avg_status": 0
        },
        {
          "from": 26,
          "to": 29,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-13",
          "message": "test case 13 message",
          "name": "test case 13",
          "time": 3,
          "avg_to": 29,
          "avg_time": 3,
          "avg_status": 0
        },
        {
          "from": 29,
          "to": 30,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-14",
          "message": "test case 14 message",
          "name": "test case 14",
          "time": 1,
          "avg_to": 30,
          "avg_time": 1,
          "avg_status": 0
        },
        {
          "from": 30,
          "to": 33,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-15",
          "message": "test case 15 message",
          "name": "test case 15",
          "time": 3,
          "avg_to": 33,
          "avg_time": 3,
          "avg_status": 0
        }
      ],
      "meta": {
        "time": 33,
        "avg_time": 35.55,
        "page": 2,
        "per_page": 10,
        "page_count": 2,
        "total_count": 15
      }
    }
    """

  Scenario: GET gantt intervals with history for 3 months, first suite
    When I am admin
    When I do GET /api/v4/cat/junit/test-suites/test-suite-to-get-gantt-1-1/gantt?page=1&limit=10&months=3
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "from": 0,
          "to": 0,
          "status": 1,
          "_id": "test-suite-to-get-gantt-1-test-case-1",
          "message": "test case 1 message",
          "name": "test case 1",
          "time": 0,
          "avg_to": 0,
          "avg_time": 0,
          "avg_status": 1
        },
        {
          "from": 0,
          "to": 3,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-2",
          "message": "test case 2 message",
          "name": "test case 2",
          "time": 3,
          "avg_to": 3,
          "avg_time": 3,
          "avg_status": 0
        },
        {
          "from": 3,
          "to": 9,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-3",
          "message": "test case 3 message",
          "name": "test case 3",
          "time": 6,
          "avg_to": 9,
          "avg_time": 6,
          "avg_status": 0
        },
        {
          "from": 9,
          "to": 12,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-4",
          "message": "test case 4 message",
          "name": "test case 4",
          "time": 3,
          "avg_to": 12,
          "avg_time": 3,
          "avg_status": 0
        },
        {
          "from": 12,
          "to": 15,
          "status": 2,
          "_id": "test-suite-to-get-gantt-1-test-case-5",
          "message": "test case 5 message",
          "name": "test case 5",
          "time": 3,
          "avg_to": 15,
          "avg_time": 3,
          "avg_status": 2
        },
        {
          "from": 15,
          "to": 20,
          "status": 3,
          "_id": "test-suite-to-get-gantt-1-test-case-6",
          "message": "test case 6 message",
          "name": "test case 6",
          "time": 5,
          "avg_to": 20,
          "avg_time": 5,
          "avg_status": 3
        },
        {
          "from": 20,
          "to": 21,
          "status": 2,
          "_id": "test-suite-to-get-gantt-1-test-case-7",
          "message": "test case 7 message",
          "name": "test case 7",
          "time": 1,
          "avg_to": 21,
          "avg_time": 1,
          "avg_status": 2
        },
        {
          "from": 21,
          "to": 22,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-8",
          "message": "test case 8 message",
          "name": "test case 8",
          "time": 1,
          "avg_to": 22,
          "avg_time": 1,
          "avg_status": 0
        },
        {
          "from": 22,
          "to": 22,
          "status": 1,
          "_id": "test-suite-to-get-gantt-1-test-case-9",
          "message": "test case 9 message",
          "name": "test case 9",
          "time": 0,
          "avg_to": 22,
          "avg_time": 0,
          "avg_status": 1
        },
        {
          "from": 22,
          "to": 22.1,
          "status": 3,
          "_id": "test-suite-to-get-gantt-1-test-case-10",
          "message": "test case 10 message",
          "name": "test case 10",
          "time": 0.1,
          "avg_to": 22.1,
          "avg_time": 0.1,
          "avg_status": 3
        }
      ],
      "meta": {
        "time": 38.1,
        "avg_time": 38.1,
        "page": 1,
        "per_page": 10,
        "page_count": 2,
        "total_count": 15
      }
    }
    """
    When I do GET /api/v4/cat/junit/test-suites/test-suite-to-get-gantt-1-1/gantt?page=2&limit=10&months=3
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "from": 22.1,
          "to": 24.1,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-11",
          "message": "test case 11 message",
          "name": "test case 11",
          "time": 2,
          "avg_to": 24.1,
          "avg_time": 2,
          "avg_status": 0
        },
        {
          "from": 24.1,
          "to": 31.1,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-12",
          "message": "test case 12 message",
          "name": "test case 12",
          "time": 7,
          "avg_to": 31.1,
          "avg_time": 7,
          "avg_status": 0
        },
        {
          "from": 31.1,
          "to": 34.1,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-13",
          "message": "test case 13 message",
          "name": "test case 13",
          "time": 3,
          "avg_to": 34.1,
          "avg_time": 3,
          "avg_status": 0
        },
        {
          "from": 34.1,
          "to": 35.1,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-14",
          "message": "test case 14 message",
          "name": "test case 14",
          "time": 1,
          "avg_to": 35.1,
          "avg_time": 1,
          "avg_status": 0
        },
        {
          "from": 35.1,
          "to": 38.1,
          "status": 0,
          "_id": "test-suite-to-get-gantt-1-test-case-15",
          "message": "test case 15 message",
          "name": "test case 15",
          "time": 3,
          "avg_to": 38.1,
          "avg_time": 3,
          "avg_status": 0
        }
      ],
      "meta": {
        "time": 38.1,
        "avg_time": 38.1,
        "page": 2,
        "per_page": 10,
        "page_count": 2,
        "total_count": 15
      }
    }
    """

  Scenario: GET gantt intervals no test cases
    When I am admin
    When I do GET /api/v4/cat/junit/test-suites/d4a1fb11-9825-40cc-9b17-38bb4ac1e3zz/gantt
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "time": 0,
        "avg_time": 0,
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 0
      }
    }
    """

  Scenario: GET gantt intervals, test-suite not found
    When I am admin
    When I do GET /api/v4/cat/junit/test-suites/test-suite-to-get-gantt-1-3-not-found/gantt
    Then the response code should be 404
