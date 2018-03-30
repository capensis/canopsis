# Alarms Tray

This document describes UI features of the alarms tray that are required
from a user prospective.

## Contents

### Description

Alarms Tray **MUST** replace the current Events Tray in Canopsis, thanks
to new Alerts engine, that handles events diffently.

The general purpose of Alarms Tray is to display alarms in an html table.

### Columns title

#### Alarm columns

Each alarm **MUST** have its own informations splited in multiple
columns.

All informations contained in an [Alarm]{alarm.md}
**MUST** be displayable in a column :

- uid
- connector
- connector name
- component
- resource
- entity id
- state
- status
- snooze conditions
- ack conditions
- cancelation conditions
- ticket declaration
- output (message)
- open date
- solved date

There are two exceptions :

- extra : extra is a set of custom fields. Each field in this set
  **MUST** have its own displayable column.
- steps : [Alarm]{alarm.md} steps are not supposed to be rendered in a table column.
  It **CAN** be exploited by the timeline widget, but this feature is out of
  the scope of this document.

#### Aggregated columns

Some informations are not stored in an alarm document but **MUST** be
displayable in their own columns :

- linklist
- pbehaviours

Those informations are related to an entity_id rather than an [Alarm]{alarm.md} itself.

### Columns display

The display of each column listed above **MUST** be configurable :

- Each column **CAN** be displayed or not
- Each column **MUST** render data with a specific renderer. This
  rendered **CAN** be configurable (particularly useful for
  [Alarm]{alarm.md} extra fields).

### Filtering

Custom filters **CAN** be applied to select alarms that have to be
displayed or not.

Those filters are mongo filters. A user just need to copy a column title
to name keys in his filter. Those keys **MUST** be translated to match
the underlying data model.

The following columns can *not* be filtered because it would cost too
much to aggregate all values :

- linklist
- pbehaviours

#### Dates

Due to UI live reporting behaviour, filtering alarms by period is not
achieved the same way as described above.

Alarms concerned by a date interval are alarms that have been opened,
were opened or have been closed in this interval.

### Searching

Users **MUST** be able to perform quick searches thanks to a text bar.

If the expression contains only alphanumerical characters (+ eventual
spaces, underscores, minuses), this expression **MUST** be searched on a
list of configurable fields.

Users **CAN** also perform advanced searches with a simple DSL.

#### Search DSL

This DSL (Domain-Specific Language) **MUST** allow the following conditional expressions :

- `FIELD = VALUE` (and `!=`, `<`, `<=`, `>`, `>=`, `CONTAINS`, `LIKE` operators)
- `NOT FIELD CONTAINS VALUE` (comparison negation)
- `FIELD < VALUE AND FIELD2 REGEX VALUE2` (and `OR` condition operator)
- `ALL F1 = V1` (meta-parameter changing search behaviour)

Those quick searches **MUST** apply a filter on the subset of alarms
returned by the permanent filter by default. Quick searches **CAN**
apply a filter on all alarms (ignoring current filter) if the expression
is prefixed by `ALL`.

UI **SHOULD** warn users whose search expressions are not *grammatically* correct.

**Note:** Expressions with parenthesis are not supported.

### Columns sorting

Each column title **CAN** be clicked to toggle sorting :

- The first click **MUST** sort raws with a DESC filter on this field
- The second click **MUST** sort raws with an ASC filter on this field
- The third click **MUST** clear any sorting on this field

This behaviour implies that sorting can be done only for *one field at a
time*.

### Pagination

Results **MUST** be paginated. Alarms tray **MUST** show at any time :

- first and last indexes of displayed alarms
- number of alarms in total
- optionally whether or not total number has been truncated

Those values **MUST** take in account filters, searches and date
intervals.

**Note:** For performance reasons, the number of alarms in total can be
approximated or truncated.
