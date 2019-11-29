
# Test cases for alarms list widget  

## Searching

* [ ] **The search with magnifier button displays relevant results.**
    1. To input some expected data to the search line.  
    2. To click on button with magnifier.
        
    **Expected result**: the search displays relevant results.

* [ ] **The search with button "Enter" displays relevant results.**
    1. To input some expected data to the search line.
    2. To click on button "Enter".
        
    **Expected result**: the search displays relevant results.

* [ ] **The empty search shows no results.**
    1. To click on button "Enter" with empty search field.
        
    **Expected result**: nothing changed.

* [ ] **The button with cross cancels current search.**
    1. To input some expected data to the search line.
    2. To click on button "Enter".
    3. To click on button with cross.
        
    **Expected result**: the current search is canceled.

* [ ] **The click on the button with question mark shows pop-up with additional information.**
    1. To place a cursor to the button with question mark.
        
    **Expected result**: pop-up with additional information appears.

* [ ] **Removing a cursor from pop-up with additional information makes it disappear.**
    1. To place a cursor to the button with question mark.
    2. To remove a cursor from pop-up with additional information.
        
    **Expected result**: pop-up with additional information disappears.
        
## WIP: Pagination

* [ ] **Right arrow opens the next page.**
    1. To click on the right arrow at the top of the page.
        
    **Expected result**: the next page opens.

* [ ] **Left arrow opens the previous page.**
    1. To click on the right arrow at the top of the page.
    2. To click on the left arrow at the top of the page.
        
    **Expected result**: the previous page opens.

* [ ] **Right arrow at the bottom of the widget opens the next page.**
    1. To click on the right arrow at the bottom of the widget.
        
    **Expected result**: the next page opens.

* [ ] **Left arrow at the bottom of the widget opens the previous page.**
    1. To click on the right arrow at the bottom of the widget.
    2. To click on the left arrow at the top of the page.
        
    **Expected result**: the previous page opens.

* [ ] **The button with page number at the bottom of the widget leads to the selected page.**
    1. To click on any button with page number at the bottom of the widget.
        
    **Expected result**: the selected page is shown.

* [ ] **5 alarms can be shown on the page.**
    1. To select “5” in the dropdown at the bottom of the widget.
        
    **Expected result**: 5 alarms are shown.

* [ ] **10 alarms can be shown on the page.**
    1. To select “10” in the dropdown at the bottom of the widget.
        
    **Expected result**: 10 alarms are shown.

## WIP: Filters

* [ ] **A filter can be selected.**
    1. To click on the field "Select a filter".
    2. To click on any filter.
        
    **Expected result**: the selected filter appears in the field "Select a filter", the search shows relevant results.

* [ ] **A selection of filter can be changed.**
    1. To click on the field "Select a filter".
    2. To click on any filter.
    3. To click on the field "Select a filter" again.
    4. To select another filter.
        
    **Expected result**: the last selected filter appears in the field "Select a filter", the search shows relevant results.

* [ ] **The button with cross cancels the selection of filters.**
    1. To select any filter.
    2. To click on the button with cross.
        
    **Expected result**: the selection of filter is canceled (the field "Select a filter" became empty, the table returns to default view).

* [ ] **The "AND" option of "Mix filters" works correctly.**
    1. To select any filter.
    2. To turn on "Mix filters".
    3. To select "AND" option in check-box.
    4. To select the second filter.
        
    **Expected result**: the results in table satisfy the conditions of both filters at the same time.

* [ ] **The "OR" option of "Mix filters" works correctly.**
    1. To select any filter.
    2. To turn on "Mix filters".
    3. To select OR option in check-box.
    4. To select the second filter.
        
    **Expected result**: the results in table satisfy the condition of at least one filter.

* [ ] **Filter can be deleted.**
    1. To click on the triangular button to the right of "Select a filter" field.
    2. To click on the trash bin icon on appeared modal window.
    3. To click "Yes" on the modal window "Are you sure?".
        
    **Expected result**: the filter is deleted.

* [ ] **The filter can be changed.**
    1. To click on the triangular button to the right of "Select a filter" field.
    2. To click on the pencil icon on appeared modal window.
    3. To change any parameter of the filter.
    4. To press the button "Submit".
        
    **Expected result**: the filter was changed.

* [ ] **A new filter can be created.**
    1. To click on the triangular button to the right of "Select a filter" field.
    2. To press the button "Add" on the appeared modal window.
    3. To select any parameters for filter.
    4. To press the button "Submit".
        
    **Expected result**: the new filter appears.

* [ ] **A new filter works correctly.**
    1. To create a new filter.
    2. To select this filter in the field "Select a filter".
        
    **Expected result**: the search displays relevant results.

* [ ] **"Live reporting" can be created for determined period of time.**
    1. To click on the clock icon in the upper right corner of the page.
    2. To select any option in the field "Quick ranges" except "Custom".
    3. To press the button "Apply".
        
    **Expected result**: a "live reporting" for determined period of time is created and relevant results are shown in table.

* [ ] **"Live reporting" can be deleted.**
    1. To create a "live reporting".
    2. To press the button with cross to delete it.
        
    **Expected result**: the "live reporting" is deleted and relevant results are shown in table.

## Mass actions

* [ ] **All elements in the table can be selected.**
    1. To click on check-box at the very first row of the table.
        
    **Expected result**: all elements in table are selected.

* [ ] **The only one element in the table can be selected.**
    1. To click on check-box in the row of needed element.
        
    **Expected result**: the needed element is selected.

* [ ] **Pressing on button "Periodical behavior" creates periodical behavior.**
    1. To select any element in the table.
    2. To press the button "Periodical behavior".
    3. To input all needed information on modal window "Create periodical behavior".
    4. To press the button "Save changes".
        
    **Expected result**: periodical behavior is created.

* [ ] **An acknowledge without a ticket can be created.**
    1. To select any element in the table.
    2. To press the button "Ack".
    3. To fill in the form.
    4. To press the button "Acknowledge".
        
    **Expected result**: periodical behavior is created.

* [ ] **An acknowledge with ticket can be created.**
    1. To select any element in the table.
    2. To press the button "Ack".
    3. To fill in the form.
    4. To press the button "Acknowledge and declare ticket".
        
    **Expected result**: an acknowledge with ticket is created.

* [ ] **"Fast ask" can be created.**
    1. To select any element in the table.
    2. To press the button "Fask ack".
        
    **Expected result**: the "fast ask" is created.

* [ ] **An "ask" can be canceled.**
    1. To select any element in the table.
    2. To press the button "Cancel ack".
        
    **Expected result**: the ask is canceled.

* [ ] **An alarm can be canceled.**
    1. To select any element in the table.
    2. To press the button "Cancel alarm".
        
    **Expected result**: the alarm is canceled.

## WIP: Columns

* [ ] **Elements can be sorted by column.**
    1. To click on a column at the very first row of the table.
        
    **Expected result**: all elements are sorted by connector name.

## WIP: Table row

* [ ] **Information pop-up can be shown.**
    1. To click on the icon "i" in the column "Connector".
        
    **Expected result**: information pop-up is shown.

* [ ] **Information pop-up can be closed.**
    1. To click on the icon "i" in the column "Connector".
    2. To click on the cross button on the pop-up.
        
    **Expected result**: information pop-up is shown.

* [ ] **Pressing on item shows details about this item.**
    1. To click on any item in the table.
        
    **Expected result**: the details about this item are shown.

* [ ] **Details about item can be closed.**
    1. To click on any item in the table.
    2. To click on this item again.
        
    **Expected result**: the details about this item are hidden.

* [ ] **Placing a cursor on signs in the column "Extra details" makes information pop-up show.**
    1. To place a cursor on any sign in the column "Extra details" (ack, ticket, canceled, snooze, pbehavior).
        
    **Expected result**: information pop-up is shown.

## WIP: Item actions when the "Ack" feature is off.

* [ ] **The "Ack" feature can be added without declaring the ticket.**
    1. To find any element in the table where "Ack" feature is off.
    2. To press on the button "Ack" in the column "Actions".
    3. To fill in the form.
    4. To press the button "Acknowledge".
        
    **Expected result**: the sign "Ack" appeared in the column "Extra details" and more features are shown in the column "Action".

* [ ] **The "Acknowledge and Declare Ticket" feature adds acknowledge and declares ticket.**
    1. To find any element in the table where "Ack" feature is off.
    2. To press on the button "Ack" in the column "Actions".
    3. To fill in the form.
    4. To press the button "Acknowledge and Declare Ticket".
        
    **Expected result**: the signs "Ack" and "Declare ticket" appeared in the column "Extra details" and more features are shown in the column "Action".

* [ ] **The "Fast ack" feature adds acknowledge without filling in the form.**
    1. To find any element in the table where "Ack" feature is off.
    2. To press on the button "Fast ack" in the column "Actions".
        
    **Expected result**: the sign "Ack" appeared in the column "Extra details" and more features are shown in the column "Action".

* [ ] **The "Snooze alarm" feature creates a "Snooze" event.**
    1. To find any element in the table where "Ack" feature is off.
    2. To press on the button "Snooze alarm" in the column "Actions".
    3. To choose a duration.
    4. To press "Save changes" button.
        
    **Expected result**: the sign "Snooze alarm" appeared in the column "Extra details".

* [ ] **The three dots icon opens a menu with three options (when the "Ack" feature is off).**
    1. To find any element in the table where "Ack" feature is off.
    2. To click on the icon with three dots in the column "Actions".
        
    **Expected result**: the menu with three options is opened.

* [ ] **The "Report alarm" feature reports an alarm.**
    1. To find any element in the table where "Ack" feature is off.
    2. To open a menu.
    3. To click on the icon "Report alarm".
    4. To fill in the form.
    5. To press the button "Notify".
    6. To wait about 60 seconds.
        
    **Expected result**: the text in the column "Output" is changed and more features are shown in the column "Action".

* [ ] **The "Periodical behavior" feature creates periodical behavior.**
    1. To find any element in the table where "Ack" feature is off.
    2. To open a menu.
    3. To click on the icon "Periodical behavior".
    4. To input all needed information on pop-up "Create periodical behavior".
    5. To press the button "Notify".
    6. To press the button "Save changes".
        
    **Expected result**: periodical behavior is created.

* [ ] **The "List periodic behaviors" feature allows to look through all periodic behaviors.**
    1. To find any element in the table where "Ack" feature is off.
    2. To open a menu.
    3. To click on the icon "List periodic behaviors".
        
    **Expected result**: the pop-up with all periodic behaviors is shown.

## WIP: Item actions when the "Ack" feature is off.

* [ ] **The button "Declare ticket" reports an incident.**
    1. To find any element in the table where "Ack" feature is on.
    2. To press the button "Declare ticket" in the column "Actions".
    3. To press the button "Report an incident" on the appeared modal window.
        
    **Expected result**: an incident reported, a notification is shown.

* [ ] **The "Associate ticket" feature associates a ticket.**
    1. To find any element in the table where "Ack" feature is on.
    2. To press the button "Associate ticket" in the column "Actions".
    3. To fill in the field "Number of the ticket".
    4. To press the button "Save changes".
        
    **Expected result**: a ticket is associated.

* [ ] **The button "Cancel alarm" changes element status to "Canceled".**
    1. To find any element in the table where "Ack" feature is on.
    2. To press the button "Cancel alarm" in the column "Actions".
    3. To fill in the field "Note".
    4. To press the button "Save changes".
        
    **Expected result**: the icon of trash bin appeared in the column "Extra details", the status changed to "Canceled", the list of actions in the column "Actions" is changed.

* [ ] **The three dots icon opens a menu with six options (when the "Ack" feature is on).**
    1. To find any element in the table where "Ack" feature is on.
    2. To click on the icon with three dots in the column "Actions".
        
    **Expected result**: the menu with six options is opened.

* [ ] **The "Cancel ack" feature cancels the "Acknowledge" mode.**
    1. To find any element in the table where "Ack" feature is on.
    2. To open a menu.
    3. To click on the icon "Cancel ack".
    4. To fill in the field "Note".
    5. To press the button "Save changes".
        
    **Expected result**: a notification is shown, the quantity of options in menu is reduced, the "Ack" mode is off.

* [ ] **The criticity of alarm can be changed.**
    1. To find any element in the table where "Ack" feature is on.
    2. To open a menu.
    3. To click on the icon "Change criticity".
    4. To choose between values "minor", "major" and "critical".
    5. To fill in the field "Note".
    6. To press the button "Save changes".
        
    **Expected result**: the criticity is changed to selected value.

* [ ] **The "Ack" feature can be added without declaring the ticket.**
    1. To find any element in the table where "Ack" feature is on.
    2. To open a menu.
    3. To press on the button "Ack".
    4. To write a note.
    5. To press the button "Acknowledge".
        
    **Expected result**: the ack event was sent.

* [ ] **The "Acknowledge and Associate Ticket" feature adds acknowledge and associates ticket.**
    1. To find any element in the table where "Ack" feature is on.
    2. To open a menu.
    3. To press on the button "Ack".
    4. To write a ticket.
    5. To write a note.
    6. To press the button "Acknowledge".
        
    **Expected result**: the ack and assocticket events were sent.

* [ ] **The "Acknowledge and Declare Ticket" feature adds acknowledge and declares ticket.**
    1. To find any element in the table where "Ack" feature is on.
    2. To open a menu.
    3. To press on the button "Ack".
    4. To write a note.
    5. To press the button "Acknowledge and Declare Ticket"
        
    **Expected result**: the ack and declareticket events were sent.

* [ ] **The "Acknowledge" feature adds acknowledge by additional modal window.**
    1. To find any element in the table where "Ack" feature is on.
    2. To open a menu.
    3. To press on the button "Ack".
    4. To write a ticket.
    5. To write a note.
    6. To press the button "Acknowledge"
    7. To press the button "Continue" on the appeared modal window.
        
    **Expected result**: the ack event was sent.

* [ ] **The "Acknowledge and Associate ticket" feature adds acknowledge and associates ticket by additional modal window.**
    1. To find any element in the table where "Ack" feature is on.
    2. To open a menu.
    3. To press on the button "Ack".
    4. To write a ticket.
    5. To write a note.
    6. To press the button "Acknowledge"
    7. To press the button "Continue and associate ticket" on the appeared modal window.
        
    **Expected result**: the ack and assocticket events were sent.

* [ ] **The "Snooze alarm" feature creates a "Snooze" event.**
    1. To find any element in the table where "Ack" feature is on.
    2. To open a menu.
    3. To press on the button "Snooze alarm".
    4. To choose a duration.
    5. To press "Save changes" button.
        
    **Expected result**: the sign "Snooze alarm" appeared in the column "Extra details".

* [ ] **The "Periodical behavior" feature creates periodical behavior.**
    1. To find any element in the table where "Ack" feature is on.
    2. To open a menu.
    3. To click on the icon "Periodical behavior".
    4. To input all needed information on modal window "Create periodical behavior".
    5. To press "Save changes" button.
        
    **Expected result**: periodical behavior is created.

* [ ] **The "List periodic behaviors" feature allows to look through all periodic behaviors.**
    1. To find any element in the table where "Ack" feature is on.
    2. To open a menu.
    3. To click on the icon "List periodic behaviors".
        
    **Expected result**: the modal window with all periodic behaviors is shown.

## The item “Actions” when an alarm is canceled.

* [ ] **The “Snooze alarm” feature creates a “Snooze” event.**
    1. To find any canceled alarm in the table or to create a new one.
    2. To press on the button “Snooze alarm”.
    3. To choose a duration.
    4. To press “Save changes” button.
        
    **Expected result**: the sign “Snooze alarm” appeared in the column “Extra details”.

* [ ] **The “Periodical behavior” feature creates periodical behavior.**
    1. To find any canceled alarm in the table or to create a new one.
    2. To click on the icon “Periodical behavior”.
    3. To input all needed information on "Create periodical behavior" modal window.
    4. To press “Save changes” button.
        
    **Expected result**: periodical behavior is created.

* [ ] **The “List periodic behaviors” feature allows to look through all periodic behaviors.**
    1. To find any canceled alarm in the table or to create a new one.
    2. To click on the icon “List periodic behaviors”.
        
    **Expected result**: the modal window with all periodic behaviors is shown.
