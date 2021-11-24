Feature: get test suite's details
  I need to be able to get test suite's details
  Only admin should be able to get test suite's details

  Scenario: GET unauthorized
    When I do GET /api/v4/cat/junit/test-suites/test-suite-to-get-details-1-1/details
    Then the response code should be 401

  Scenario: GET without permissions
    When I am noperms
    When I do GET /api/v4/cat/junit/test-suites/test-suite-to-get-details-1-1/details
    Then the response code should be 403

  Scenario: GET details success, page = 1
    When I am admin
    When I do GET /api/v4/cat/junit/test-suites/test-suite-to-get-details-1-1/details
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "data": [
        {
          "_id": "test-suite-to-get-details-1-test-case-1",
          "name": "test case 1",
          "status": 1,
          "time": 0,
          "message": "test case 1 message",
          "description": "test case 1 description",
          "file": "test case 1 file",
          "classname": "test case 1 classname",
          "description": "test case 1 description",
          "file": "test case 1 file",
          "classname": "test case 1 classname",
          "screenshots": [
            "screenshot-1",
            "screenshot-2"
          ],
          "videos": [
            "video-1"
          ]
        },
        {
          "_id": "test-suite-to-get-details-1-test-case-2",
          "name": "test case 2",
          "line": 2,
          "status": 0,
          "time": 3,
          "message": "test case 2 message",
          "description": "test case 2 description",
          "file": "test case 2 file",
          "classname": "test case 2 classname"
        },
        {
          "_id": "test-suite-to-get-details-1-test-case-3",
          "name": "test case 3",
          "status": 0,
          "time": 6,
          "message": "test case 3 message",
          "description": "test case 3 description",
          "file": "test case 3 file",
          "classname": "test case 3 classname"
        },
        {
          "_id": "test-suite-to-get-details-1-test-case-4",
          "name": "test case 4",
          "status": 0,
          "time": 3,
          "message": "test case 4 message",
          "description": "test case 4 description",
          "file": "test case 4 file",
          "classname": "test case 4 classname"
        },
        {
          "_id": "test-suite-to-get-details-1-test-case-5",
          "name": "test case 5",
          "line": 5,
          "status": 2,
          "time": 3,
          "message": "test case 5 message",
          "description": "test case 5 description",
          "file": "test case 5 file",
          "classname": "test case 5 classname"
        },
        {
          "_id": "test-suite-to-get-details-1-test-case-6",
          "name": "test case 6",
          "line": 6,
          "status": 3,
          "time": 5,
          "message": "test case 6 message",
          "description": "test case 6 description",
          "file": "test case 6 file",
          "classname": "test case 6 classname"
        },
        {
          "_id": "test-suite-to-get-details-1-test-case-7",
          "name": "test case 7",
          "line": 7,
          "status": 2,
          "time": 1,
          "message": "test case 7 message",
          "description": "test case 7 description",
          "file": "test case 7 file",
          "classname": "test case 7 classname"
        },
        {
          "_id": "test-suite-to-get-details-1-test-case-8",
          "name": "test case 8",
          "status": 0,
          "time": 1,
          "message": "test case 8 message",
          "description": "test case 8 description",
          "file": "test case 8 file",
          "classname": "test case 8 classname"
        },
        {
          "_id": "test-suite-to-get-details-1-test-case-9",
          "name": "test case 9",
          "status": 1,
          "time": 0,
          "message": "test case 9 message",
          "description": "test case 9 description",
          "file": "test case 9 file",
          "classname": "test case 9 classname"
        },
        {
          "_id": "test-suite-to-get-details-1-test-case-10",
          "name": "test case 10",
          "line": 10,
          "status": 3,
          "time": 0.1,
          "message": "test case 10 message",
          "description": "test case 10 description",
          "file": "test case 10 file",
          "classname": "test case 10 classname"
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 2,
        "total_count": 15
      }
    }
    """

  Scenario: GET details success, page = 2
    When I am admin
    When I do GET /api/v4/cat/junit/test-suites/test-suite-to-get-details-1-1/details?page=2
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "data": [
        {
          "_id": "test-suite-to-get-details-1-test-case-11",
          "name": "test case 11",
          "status": 0,
          "time": 2,
          "message": "test case 11 message",
          "description": "test case 11 description",
          "file": "test case 11 file",
          "classname": "test case 11 classname"
        },
        {
          "_id": "test-suite-to-get-details-1-test-case-12",
          "name": "test case 12",
          "status": 0,
          "time": 7,
          "message": "test case 12 message",
          "description": "test case 12 description",
          "file": "test case 12 file",
          "classname": "test case 12 classname"
        },
        {
          "_id": "test-suite-to-get-details-1-test-case-13",
          "name": "test case 13",
          "status": 0,
          "time": 3,
          "message": "test case 13 message",
          "description": "test case 13 description",
          "file": "test case 13 file",
          "classname": "test case 13 classname"
        },
        {
          "_id": "test-suite-to-get-details-1-test-case-14",
          "name": "test case 14",
          "status": 0,
          "time": 1,
          "message": "test case 14 message",
          "description": "test case 14 description",
          "file": "test case 14 file",
          "classname": "test case 14 classname"
        },
        {
          "_id": "test-suite-to-get-details-1-test-case-15",
          "name": "test case 15",
          "status": 0,
          "time": 3,
          "message": "test case 15 message",
          "description": "test case 15 description",
          "file": "test case 15 file",
          "classname": "test case 15 classname"
        }
      ],
      "meta": {
        "page": 2,
        "per_page": 10,
        "page_count": 2,
        "total_count": 15
      }
    }
    """

  Scenario: GET details success, sort by status
    When I am admin
    When I do GET /api/v4/cat/junit/test-suites/test-suite-to-get-details-1-1/details?sort_by=status&sort=desc
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "data": [
        {
          "_id": "test-suite-to-get-details-1-test-case-10",
          "name": "test case 10",
          "line": 10,
          "status": 3,
          "time": 0.1,
          "message": "test case 10 message",
          "description": "test case 10 description",
          "file": "test case 10 file",
          "classname": "test case 10 classname"
        },
        {
          "_id": "test-suite-to-get-details-1-test-case-6",
          "name": "test case 6",
          "line": 6,
          "status": 3,
          "time": 5,
          "message": "test case 6 message",
          "description": "test case 6 description",
          "file": "test case 6 file",
          "classname": "test case 6 classname"
        },
        {
          "_id": "test-suite-to-get-details-1-test-case-5",
          "name": "test case 5",
          "line": 5,
          "status": 2,
          "time": 3,
          "message": "test case 5 message",
          "description": "test case 5 description",
          "file": "test case 5 file",
          "classname": "test case 5 classname"
        },
        {
          "_id": "test-suite-to-get-details-1-test-case-7",
          "name": "test case 7",
          "line": 7,
          "status": 2,
          "time": 1,
          "message": "test case 7 message",
          "description": "test case 7 description",
          "file": "test case 7 file",
          "classname": "test case 7 classname"
        },
        {
          "_id": "test-suite-to-get-details-1-test-case-1",
          "name": "test case 1",
          "description": "test case 1 description",
          "file": "test case 1 file",
          "classname": "test case 1 classname",
          "status": 1,
          "time": 0,
          "message": "test case 1 message",
          "description": "test case 1 description",
          "file": "test case 1 file",
          "classname": "test case 1 classname",
          "screenshots": [
            "screenshot-1",
            "screenshot-2"
          ],
          "videos": [
            "video-1"
          ]
        },
        {
          "_id": "test-suite-to-get-details-1-test-case-9",
          "name": "test case 9",
          "status": 1,
          "time": 0,
          "message": "test case 9 message",
          "description": "test case 9 description",
          "file": "test case 9 file",
          "classname": "test case 9 classname"
        },
        {
          "_id": "test-suite-to-get-details-1-test-case-11",
          "name": "test case 11",
          "status": 0,
          "time": 2,
          "message": "test case 11 message",
          "description": "test case 11 description",
          "file": "test case 11 file",
          "classname": "test case 11 classname"
        },
        {
          "_id": "test-suite-to-get-details-1-test-case-12",
          "name": "test case 12",
          "status": 0,
          "time": 7,
          "message": "test case 12 message",
          "description": "test case 12 description",
          "file": "test case 12 file",
          "classname": "test case 12 classname"
        },
        {
          "_id": "test-suite-to-get-details-1-test-case-13",
          "name": "test case 13",
          "status": 0,
          "time": 3,
          "message": "test case 13 message",
          "description": "test case 13 description",
          "file": "test case 13 file",
          "classname": "test case 13 classname"
        },
        {
          "_id": "test-suite-to-get-details-1-test-case-14",
          "name": "test case 14",
          "status": 0,
          "time": 1,
          "message": "test case 14 message",
          "description": "test case 14 description",
          "file": "test case 14 file",
          "classname": "test case 14 classname"
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 2,
        "total_count": 15
      }
    }
    """

  Scenario: GET details, test-suite not found
    When I am admin
    When I do GET /api/v4/cat/junit/test-suites/test-suite-to-get-details-1-1-not-found/details
    Then the response code should be 404
