```


User = {
    'rights': ...,               // Map of type Rights
    'groups': ...,               // List of strings (groups names)
    'profile': ...,              // String of profile name (Admin, Root, Manager, ...)
    'contact': {                 // Map of contact informations
        'mail': ...,
        'phone_number': ...,
        ...
        }
    'name': ...,                 // String of user's name
    '_id': ...                   // uniq id
}

Group = {
    'name': ...,                 // String of group's name
    'members': ...,              // List of strings (members names)
    'rights': ...                // Map of type Rights
}

Rghts = {
    object_id...: {             // Right on the object with the identifier id
        'right': ...,           // 1 == Read, 2 == Update, 4 == Create, 8 == Delete
        'desc': ...,            // Short desc of the right
        'context': ...          // Time period
        }
}
```
