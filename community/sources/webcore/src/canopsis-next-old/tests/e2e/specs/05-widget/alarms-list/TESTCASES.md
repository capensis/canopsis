
# Test cases for alarms list widget  

## Searching

* [x] **The search with magnifier button displays relevant results.**
    1. To input some expected data to the search line.  
    2. To click on button with magnifier.
        
    **Expected result**: the search displays relevant results.

* [x] **The search with button "Enter" displays relevant results.**
    1. To input some expected data to the search line.
    2. To click on button "Enter".
        
    **Expected result**: the search displays relevant results.

* [x] **The empty search shows no results.**
    1. To click on button "Enter" with empty search field.
        
    **Expected result**: nothing changed.

* [x] **The button with cross cancels current search.**
    1. To input some expected data to the search line.
    2. To click on button "Enter".
    3. To click on button with cross.
        
    **Expected result**: the current search is canceled.

* [x] **The click on the button with question mark shows pop-up with additional information.**
    1. To place a cursor to the button with question mark.
        
    **Expected result**: pop-up with additional information appears.

* [x] **Removing a cursor from pop-up with additional information makes it disappear.**
    1. To place a cursor to the button with question mark.
    2. To remove a cursor from pop-up with additional information.
        
    **Expected result**: pop-up with additional information disappears.
        
## Pagination

* [x] **Right arrow opens the next page.**
    1. To click on the right arrow at the top of the page.
        
    **Expected result**: the next page opens.

* [x] **Left arrow opens the previous page.**
    1. To click on the right arrow at the top of the page.
    2. To click on the left arrow at the top of the page.
        
    **Expected result**: the previous page opens.

* [x] **Right arrow at the bottom of the widget opens the next page.**
    1. To click on the right arrow at the bottom of the widget.
        
    **Expected result**: the next page opens.

* [x] **Left arrow at the bottom of the widget opens the previous page.**
    1. To click on the right arrow at the bottom of the widget.
    2. To click on the left arrow at the top of the page.
        
    **Expected result**: the previous page opens.

* [x] **The button with page number at the bottom of the widget leads to the selected page.**
    1. To click on any button with page number at the bottom of the widget.
        
    **Expected result**: the selected page is shown.

* [x] **5 alarms can be shown on the page.**
    1. To select “5” in the dropdown at the bottom of the widget.
        
    **Expected result**: 5 alarms are shown.

* [x] **10 alarms can be shown on the page.**
    1. To select “10” in the dropdown at the bottom of the widget.
        
    **Expected result**: 10 alarms are shown.

## Filtering

* [x] **A filter can be selected.**
    1. To click on the field "Select a filter".
    2. To click on any filter.
        
    **Expected result**: the selected filter appears in the field "Select a filter", the search shows relevant results.

* [x] **A selection of filter can be changed.**
    1. To click on the field "Select a filter".
    2. To click on any filter.
    3. To click on the field "Select a filter" again.
    4. To select another filter.
        
    **Expected result**: the last selected filter appears in the field "Select a filter", the search shows relevant results.

* [x] **The button with cross cancels the selection of filters.**
    1. To select any filter.
    2. To click on the button with cross.
        
    **Expected result**: the selection of filter is canceled (the field "Select a filter" became empty, the table returns to default view).

* [x] **The "AND" option of "Mix filters" works correctly.**
    1. To select any filter.
    2. To turn on "Mix filters".
    3. To select "AND" option in check-box.
    4. To select the second filter.
        
    **Expected result**: the results in table satisfy the conditions of both filters at the same time.

* [x] **The "OR" option of "Mix filters" works correctly.**
    1. To select any filter.
    2. To turn on "Mix filters".
    3. To select OR option in check-box.
    4. To select the second filter.
        
    **Expected result**: the results in table satisfy the condition of at least one filter.

* [x] **Filter can be deleted.**
    1. To click on the triangular button to the right of "Select a filter" field.
    2. To click on the trash bin icon on appeared modal window.
    3. To click "Yes" on the modal window "Are you sure?".
        
    **Expected result**: the filter is deleted.

* [x] **The filter can be changed.**
    1. To click on the triangular button to the right of "Select a filter" field.
    2. To click on the pencil icon on appeared modal window.
    3. To change any parameter of the filter.
    4. To press the button "Submit".
        
    **Expected result**: the filter was changed.

* [x] **The changed filter works in a new way.**
    1. To change any filter.
    2. To select this filter in the field “Select a filter”.
        
    **Expected result**: the search displays relevant results.

* [x] **A new filter can be created.**
    1. To click on the triangular button to the right of "Select a filter" field.
    2. To press the button "Add" on the appeared modal window.
    3. To select any parameters for filter.
    4. To press the button "Submit".
        
    **Expected result**: the new filter appears.

* [x] **A new filter works correctly.**
    1. To create a new filter.
    2. To select this filter in the field "Select a filter".
        
    **Expected result**: the search displays relevant results.

* [x] **"Live reporting" can be created for period of time selected by user.**
    1. To click on the clock icon in the upper right corner of the page.
    2. To choose “Custom” option in the field "Quick ranges".
    3. To choose start date interval on the pop-up "Live reporting"
    4. To choose end date interval on the pop-up "Live reporting"
    5. To press the button "Apply"
        
    **Expected result**: a custom "Live reporting" is created and relevant results are shown in table.

* [x] **"Live reporting" can be created for determined period of time.**
    1. To click on the clock icon in the upper right corner of the page.
    2. To select any option in the field "Quick ranges" except "Custom".
    3. To press the button "Apply".
        
    **Expected result**: a "live reporting" for determined period of time is created and relevant results are shown in table.

* [x] **"Live reporting" can be deleted.**
    1. To create a "live reporting".
    2. To press the button with cross to delete it.
        
    **Expected result**: the "live reporting" is deleted and relevant results are shown in table.

## Mass actions

* [x] **All elements in the table can be selected.**
    1. To click on check-box at the very first row of the table.
        
    **Expected result**: all elements in table are selected.

* [x] **The only one element in the table can be selected.**
    1. To click on check-box in the row of needed element.
        
    **Expected result**: the needed element is selected.

* [x] **Pressing on button "Periodical behavior" creates periodical behavior.**
    1. To select any element in the table.
    2. To press the button "Periodical behavior".
    3. To input all needed information on modal window "Create periodical behavior".
    4. To press the button "Save changes".
        
    **Expected result**: periodical behavior is created.

* [x] **An acknowledge without a ticket can be created.**
    1. To select any element in the table.
    2. To press the button "Ack".
    3. To fill in the form.
    4. To press the button "Acknowledge".
        
    **Expected result**: periodical behavior is created.

* [x] **An acknowledge with ticket can be created.**
    1. To select any element in the table.
    2. To press the button "Ack".
    3. To fill in the form.
    4. To press the button "Acknowledge and declare ticket".
        
    **Expected result**: an acknowledge with ticket is created.

* [x] **"Fast ask" can be created.**
    1. To select any element in the table.
    2. To press the button "Fask ack".
        
    **Expected result**: the "fast ask" is created.

* [x] **An "ask" can be canceled.**
    1. To select any element in the table.
    2. To press the button "Cancel ack".
        
    **Expected result**: the ask is canceled.

* [x] **An alarm can be canceled.**
    1. To select any element in the table.
    2. To press the button "Cancel alarm".
        
    **Expected result**: the alarm is canceled.

## Columns

* [x] **Elements can be sorted by column.**
    1. To click on a column at the very first row of the table.
        
    **Expected result**: all elements are sorted by column.

## Table row

* [x] **Information pop-up can be shown.**
    1. To click on the icon "i" in the column "Connector".
        
    **Expected result**: information pop-up is shown.

* [x] **Information pop-up can be closed.**
    1. To click on the icon "i" in the column "Connector".
    2. To click on the cross button on the pop-up.
        
    **Expected result**: information pop-up is shown.

* [x] **Pressing on item shows details about this item.**
    1. To click on any item in the table.
        
    **Expected result**: the details about this item are shown.

* [x] **Pressing on element hidden details about this element.**
    1. To click on any item in the table.
    2. To click on this item again.
        
    **Expected result**: the details about this item are hidden.

* [x] **Placing a cursor on signs in the column "Extra details" makes information pop-up show.**
    1. To place a cursor on any sign in the column "Extra details" (ack, ticket, canceled, snooze, pbehavior).
        
    **Expected result**: information pop-up is shown.

## Item actions when the "Ack" feature is off.

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

## Item actions when the "Ack" feature is on.

* [x] **The button "Declare ticket" reports an incident.**
    1. To find any element in the table where "Ack" feature is on.
    2. To press the button "Declare ticket" in the column "Actions".
    3. To press the button "Report an incident" on the appeared modal window.
        
    **Expected result**: an incident reported, a notification is shown.

* [x] **The "Associate ticket" feature associates a ticket.**
    1. To find any element in the table where "Ack" feature is on.
    2. To press the button "Associate ticket" in the column "Actions".
    3. To fill in the field "Number of the ticket".
    4. To press the button "Save changes".
        
    **Expected result**: a ticket is associated.

* [x] **The button "Cancel alarm" changes element status to "Canceled".**
    1. To find any element in the table where "Ack" feature is on.
    2. To press the button "Cancel alarm" in the column "Actions".
    3. To fill in the field "Note".
    4. To press the button "Save changes".
        
    **Expected result**: the icon of trash bin appeared in the column "Extra details", the status changed to "Canceled", the list of actions in the column "Actions" is changed.

* [x] **The three dots icon opens a menu with six options (when the "Ack" feature is on).**
    1. To find any element in the table where "Ack" feature is on.
    2. To click on the icon with three dots in the column "Actions".
        
    **Expected result**: the menu with five options is opened.

* [x] **The "Cancel ack" feature cancels the "Acknowledge" mode.**
    1. To find any element in the table where "Ack" feature is on.
    2. To open a menu.
    3. To click on the icon "Cancel ack".
    4. To fill in the field "Note".
    5. To press the button "Save changes".
        
    **Expected result**: a notification is shown, the quantity of options in menu is reduced, the "Ack" mode is off.

* [x] **The criticity of alarm can be changed.**
    1. To find any element in the table where "Ack" feature is on.
    2. To open a menu.
    3. To click on the icon "Change criticity".
    4. To choose between values "minor", "major" and "critical".
    5. To fill in the field "Note".
    6. To press the button "Save changes".
        
    **Expected result**: the criticity is changed to selected value.
    
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

## Edit mode: alarm list settings.
**Pre-condition:** 
* *Open a menu at the bottom right corner of the page*
* *Press the button with pencil*
* *Enter the edit mode*
* *Then click on the icon with three dots*
* *Choose "Edit" in the dropdown*

* [x] **Widget’s size can be changed in mobile version.**
    1. To click on the element “Widget’s size”.
    2. To choose a row.
    3. To select the widget’s size near the icon of mobile phone.
    4. To press the button “Save”.
        
    **Expected result**: the widget’s size is changed in mobile version.

* [x] **Widget’s size can be changed in tablet version.**
    1. To click on the element “Widget’s size”.
    2. To choose a row.
    3. To select the widget’s size near the icon of tablet.
    4. To press the button “Save”.
        
    **Expected result**: the widget’s size is changed in tablet version.

* [x] **Widget’s size can be changed in desktop version.**
    1. To click on the element “Widget’s size”.
    2. To choose a row.
    3. To select the widget’s size near the icon of monitor.
    4. To press the button “Save”.
        
    **Expected result**: the widget’s size is changed in desktop version.

* [x] **Widget’s title can be changed.**
    1. To click on the element “Title (optional)”.
    2. To input a title.
    3. To press the button “Save”.
        
    **Expected result**: the widget’s title is changed.

* [x] **The periodic refresh can be set manually.**
    1. To click on the element “Periodic refresh (optional)”.
    2. To turn on periodic refresh mode.
    3. To set some time for periodic refresh manually.
    4. To press the button “Save”.
        
    **Expected result**: the periodic refresh is set.

## Edit mode: alarm list settings (advanced settings).
**Pre-condition:** 
* *Open a menu at the bottom right corner of the page*
* *Press the button with pencil*
* *Enter the edit mode*
* *Then click on the icon with three dots*
* *Choose "Edit" in the dropdown*
* *Unfold the advanced settings menu*

* [x] **Default sort column can be set for table.**
    1. To click on “Default sort column” element in the “Advanced settings”.
    2. To select “Column name”.
    3. To select an order (asc or desc).
    4. To press the button “Save”.
        
    **Expected result**: the default sort column is set for table.

* [x] **A column can be added to table.**
    1. To click on “Column names” element in the “Advanced settings”.
    2. To press the button “Add”.
    3. To input label, value and to choose column’s place in the table.
    4. To press the button “Save”.
        
    **Expected result**: the new column added to the table.

* [x] **A column’s name can be changed.**
    1. To click on “Column names” element in the “Advanced settings”.
    2. To change any column’s name.
    3. To press the button “Save”.
        
    **Expected result**: a column’s name is changed.

* [x] **A column’s value can be changed.**
    1. To click on “Column names” element in the “Advanced settings”.
    2. To change any column’s value.
    3. To press the button “Save”.
        
    **Expected result**: a column’s value is changed.

* [x] **A column can be deleted from the table.**
    1. To click on “Column names” element in the “Advanced settings”.
    2. To press the cross button on any column’s card.
    3. To press the button “Save”.
        
    **Expected result**: a column is deleted from the table.

* [x] **A column’s card can be moved below.**
    1. To click on “Column names” element in the “Advanced settings”.
    2. To press on button with up arrow at any column’s card.
    3. To press the button “Save”.
        
    **Expected result**: a column’s card is moved below.

* [x] **A column’s card can be moved above.**
    1. To click on “Column names” element in the “Advanced settings”.
    2. To press on button with down arrow at any column’s card (but not at the first card).
    3. To press the button “Save”.
        
    **Expected result**: a column’s card is moved above.

* [x] **HTML mode can be set for column.**
    1. To click on “Column names” element in the “Advanced settings”.
    2. To set an HTML mode for any column.
    3. To press the button “Save”.
        
    **Expected result**: HTML mode is selected for any column.

* [x] **5 can be set as default number of elements per page.**
    1. To click on “Default number of elements/page” element in the “Advanced settings”.
    2. To choose the number 5.
    3. To press the button “Save”.
        
    **Expected result**: 5 is set as default number of elements per page.

* [x] **10 can be set as default number of elements per page.**
    1. To click on “Default number of elements/page” element in the “Advanced settings”.
    2. To choose the number 10.
    3. To press the button “Save”.
        
    **Expected result**: 10 is set as default number of elements per page.

* [x] **Filter on Open/Resolved can be turn off.**
    1. To click on “Filter on Open/Resolved” element in the “Advanced settings”.
    2. To leave the check-box “Open” empty.
    3. To leave the check-box “Resolved” empty.
    4. To press the button “Save”.
        
    **Expected result**: filter on Open/Resolved is turn off.

* [x] **Filter on Open can be set.**
    1. To click on “Filter on Open/Resolved” element in the “Advanced settings”.
    2. To mark the check-box “Open”.
    3. To leave the check-box “Resolved” empty.
    4. To press the button “Save”.
        
    **Expected result**: filter on Open is set.

* [x] **Filter on Resolved can be set.**
    1. To click on “Filter on Open/Resolved” element in the “Advanced settings”.
    2. To mark the check-box “Resolved”.
    3. To leave the check-box “Open” empty.
    4. To press the button “Save”.
        
    **Expected result**: filter on Resolved is set.

* [x] **Live reporting on resolved filter enabled should be 30 days**
    1. To press the button “Delete” on “Live reporting” section in the “Advanced settings” if it exists.
    2. To click on “Filter on Open/Resolved” element in the “Advanced settings”.
    3. To mark the check-box “Resolved”.
    4. To press the button “Save”.
        
    **Expected result**: the live reporting on the table is 30 days selected.

* [x] **Filter on Open and Resolved can be set.**
    1. To click on “Filter on Open/Resolved” element in the “Advanced settings”.
    2. To mark the check-box “Resolved”.
    3. To mark the check-box “Open”.
    4. To press the button “Save”.
        
    **Expected result**: filter on Open and Resolved is set.

* [x] **Default filter can be created in advanced settings.**
    1. To click on “Filters” element in the “Advanced settings”.
    2. To press the button “Add”.
    3. To input all required information in modal window.
    4. To press the button “Submit”.
        
    **Expected result**: a new default filter is created.

* [x] **Default filter can be edited in advanced settings.**
    1. To click on “Filters” element in the “Advanced settings”.
    2. To create a new default filter.
    3. To press the button “Edit” near the new filter.
    4. To edit some information on appeared modal window.
    5. To press the button “Submit”.
        
    **Expected result**: the filter is modified.

* [x] **Default filter can be deleted in advanced settings.**
    1. To click on “Filters” element in the “Advanced settings”.
    2. To create a new default filter.
    3. To press the button with trash bin.
    4. To press the button “Yes” on appeared modal window “Are you sure?”.
        
    **Expected result**: the filter is deleted.

* [x] **Default filter can be set in advanced settings.**
    1. To click on “Filters” element in the “Advanced settings”.
    2. To create a new default filter.
    3. To choose this filter in the dropdown “Default filter”.
    4. To press the button “Save”.
        
    **Expected result**: the default filter is set in advanced settings.

* [x] **Two default filters can be set with AND-rule.**
    1. To click on “Filters” element in the “Advanced settings”.
    2. To create two new default filters.
    3. To turn on “Mix filters”.
    4. To set the AND-rule for filters.
    5. To select both filters in the dropdown “Default filters”.
    6. To press the button “Save”.
        
    **Expected result**: two default filters are set with AND-rule.

* [x] **Two default filters can be set with OR-rule.**
    1. To click on “Filters” element in the “Advanced settings”.
    2. To create two new default filters.
    3. To turn on “Mix filters”.
    4. To set the OR-rule for filters.
    5. To select both filters in the dropdown “Default filters”.
    6. To press the button “Save”.
        
    **Expected result**: two default filters are set with OR-rule.

* [x] **Live reporting with custom dates can be created.**
    1. To press the button “Create” on “Live reporting” section in the “Advanced settings”.
    2. To choose the “Custom” option in dropdown “Quick ranges”.
    3. To set a date interval.
    4. To press the button “Apply”.
        
    **Expected result**: a live reporting with custom dates is created.

* [x] **Live reporting with determined dates can be created.**
    1. To press the button “Create” on “Live reporting” section in the “Advanced settings”.
    2. To choose any option except “Custom” in dropdown “Quick ranges”.
    3. To press the button “Apply”.
        
    **Expected result**: a live reporting with determined dates is created.

* [x] **Live reporting can be edited.**
    1. To create any live reporting in the “Advanced settings”.
    2. To press the button “Edit”.
    3. To change some settings of the live reporting.
    4. To press the button “Apply”.
        
    **Expected result**: the live reporting is edited.

* [x] **Live reporting can be deleted.**
    1. To create any live reporting in the “Advanced settings”.
    2. To press the button with trash bin.
        
    **Expected result**: the live reporting is deleted.

* [x] **Info popup can be created.**
    1. To press the button “Create/Edit” in the section “Info popup”.
    2. To press the button with plus on the appeared modal window.
    3. To select the column’s name.
    4. To input some text to the text area.
    5. To press the button “Submit”.
    6. To press the button “Submit” on the next modal window.
        
    **Expected result**: info popup is created.

* [x] **Info popup should displayed on table column.**
    1. To press the button “Create/Edit” in the section “Info popup”.
    2. To press the button with plus on the appeared modal window.
    3. To select the column’s name.
    4. To input some text to the text area.
    5. To press the button “Submit”.
    6. To press the button “Submit” on the next modal window.
    7. To press the button “Save”.
        
    **Expected result**: button with "i" is shown on the column.

* [x] **Info popup can be edited.**
    1. To create any info popup in the section “Info popup”.
    2. To press the button “Create/Edit” in the section “Info popup”.
    3. To press the button with pencil.
    4. To change any settings.
    5. To press the button “Submit”.
    6. To press the button “Submit” on the next modal window.
        
    **Expected result**: info popup is edited.

* [x] **Info popup can be deleted.**
    1. To create any info popup in the section “Info popup”.
    2. To press the button “Create/Edit” in the section “Info popup”.
    3. To press the button with trash bin.
    4. To press the button “Submit”.
        
    **Expected result**: info popup is deleted.

* [x] **More than one info popup can be added.**
    1. To create any info popup in the section “Info popup”.
    2. To press the button “Create/Edit” in the section “Info popup”.
    3. To press the button with plus on the appeared modal window.
    4. To select the column’s name.
    5. To input some text to the text area.
    6. To press the button “Submit”.
    7. To press the button “Submit” on the next modal window.
        
    **Expected result**: the second info popup is created.

* [x] **“More infos” popup can be created.**
    1. To press the button “Create” in the section “More infos popup”.
    2. To input any text to the field.
    3. To press the button “Submit”.
    4. To press the button “Save”.
        
    **Expected result**: “More infos” popup is created.

* [x] **“More infos” modal window is showing.**
    1. Create "More infos popup"
    2. To press the button “Save”.
    3. To press the button "..." on table item actions column
    4. To press "More infos"
        
    **Expected result**: “More infos” popup is shown.

* [x] **“More infos” popup can be edited.**
    1. To create any “More infos” popup in the section “More infos popup”.
    2. To press the button “Edit”.
    3. To change the text in the appeared window.
    4. To press the button “Submit”.
        
    **Expected result**: “More infos” popup is edited.

* [x] **“More infos” modal window is showing with edited data.**
    1. Create "More infos popup"
    2. To press the button “Save”.
    3. To press the button "..." on table item actions column
    4. To press "More infos"
        
    **Expected result**: “More infos” popup is shown with new data.

* [x] **“More infos” popup can be deleted.**
    1. To create any “More infos” popup in the section “More infos popup”.
    2. To press the button with trash bin.
    3. To press the button “Yes” on “Are you sure?” modal window.
        
    **Expected result**: “More infos” pop-up is deleted.

* [x] **“HTML enabled on timeline?” checkbox can be turn on.**
    1. To turn on checkbox “HTML enabled on timeline?” in the advanced settings.
        
    **Expected result**: “HTML enabled on timeline?” checkbox is turn on.

* [x] **“Note field required when ack?” checkbox can be turn on.**
    1. To select section “Ack” in the additional settings.
    2. To turn on checkbox “Note field required when ack?”.
        
    **Expected result**: “Note field required when ack?” checkbox is turn on.

* [x] **“Multiple ack” checkbox can be turn on.**
    1. To select section “Ack” in the additional settings.
    2. To turn on checkbox “Multiple ack”.
        
    **Expected result**: “Multiple ack” checkbox is turn on.

* [x] **Fast-ack output (optional) can be enabled.**
    1. To select section “Ack” in the additional settings.
    2. To select section “Fast-ack output (optional)” in the additional settings.
    3. To turn on checkbox “Enabled”.
        
    **Expected result**: Fast-ack output (optional) is enabled.

* [x] **Text of fast-ack output (optional) can be edited.**
    1. To select section “Ack” in the additional settings.
    2. To select section “Fast-ack output (optional)” in the additional settings.
    3. To change the text of fast-ack output (optional).
        
    **Expected result**: Text of fast-ack output (optional) is edited.
