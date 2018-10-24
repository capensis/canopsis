# Calendar APIs

## CRUD on calendar entities

This endpoint create/update/delete/returns a list of calendar objects, or
sources events.

#### Url

  `GET` /calendar/uids?limit=0&skip=0&sort=&projection=&with_count=

  `PUT` /calendar/?uid=&eventcategories=&output=&dtstart=&dtend=

  `POST` /calendar/?eventcategories=&output=&dtstart=&dtend=

  `DELETE` /calendar/?uids=

  `GET` /calendar/values?query=&sources=&dtstart=&dtend=

#### Params

- **uids** : a list of calendar uids
- **limit** (optional): max number of elements to get
- **skip** (optional): first element index among searched list
- **sort** (optional): contains a list of couples of field (name, ASC/DESC)
    or field name which denots an implicitelly ASC order
- **projection** (optional): key names to keep from elements
- **with_count** (optional): If True (False by default), add count to the result
- **uid** (optional): the uid of the event (POST only)
- **eventcategories**: eventcategories of the event
- **output**: description of the event
- **dtstart**: beginning date (timestamp)
- **dtend**: ending date (timestamp)
- **query** (optional): vevent information if given
- **sources** (optional): sources from where get values. If None, use all sources

#### GET example

/calendar/de3ea0b3-d6e6-4079-9eba-b57834821996

```{json}
    {
        "total": 1,
        "data": [
            {
                "_id": "de3ea0b3-d6e6-4079-9eba-b57834821996"
                "dtend": 2000000000,
                "dtstart": 0,
                "duration": 0,
                "eventcategories": {},
                "source": null,
                "rrule": null,
                "output": "a calendar event !",
                "uid": "de3ea0b3-d6e6-4079-9eba-b57834821996",
            }
        ],
        "success": true
    }

```

#### PUT example

/calendar/?eventcategories={}&output="a calendar event !"&dtstart=0&dtend=2000000000

```{json}
    {
        "total": 1,
        "data": [
            {
                "_id": "de3ea0b3-d6e6-4079-9eba-b57834821996",
                "dtend": 2000000000,
                "dtstart": 0
                "duration": 0,
                "eventcategories": {},
                "source": null,
                "rrule": null,
                "output": "a calendar event !",
                "uid": "de3ea0b3-d6e6-4079-9eba-b57834821996",
            }
        ],
        "success": true
    }
```

#### POST example

/calendar/?uid=de3ea0b3-d6e6-4079-9eba-b57834821996&eventcategories={}&output="another calendar event !"&dtstart=0&dtend=2000000001

```{json}
    {
        "total": 1,
        "data": [
            {
                "_id": "de3ea0b3-d6e6-4079-9eba-b57834821996",
                "dtend": 2000000001,
                "dtstart": 0
                "duration": 0,
                "eventcategories": {},
                "rrule": null,
                "source": null,
                "output": "another calendar event !",
                "uid": "de3ea0b3-d6e6-4079-9eba-b57834821996",
            }
        ],
        "success": true
    }
```

#### DELETE example

/calendar/?ids=[de3ea0b3-d6e6-4079-9eba-b57834821996]

```{json}
    {
        "total": 1,
        "data": [
            true
        ],
        "success": true
    }
```

#### GET VALUES example

/calendar/values/?query={}&dtstart=0

```{json}
{
    "total": 1,
    "data": [
        {
            "_id": "de3ea0b3-d6e6-4079-9eba-b57834821996"
            "dtend": 2000000000,
            "dtstart": 0,
            "duration": 0,
            "eventcategories": {},
            "rrule": null,
            "source": null,
            "output": "a calendar event !",
            "uid": "de3ea0b3-d6e6-4079-9eba-b57834821996",
        }
    ],
    "success": true
}
```
