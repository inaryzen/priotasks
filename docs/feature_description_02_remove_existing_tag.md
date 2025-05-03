# Feature Description Document - F02

## Overview
This feature allows to remove the tag that already existing in the system.

## Requirements
### Functional Requirements
- A new inline button "Delete" is added to "Tags List" element in "Task Edit" modal window
- The button is visualized as a "Delete" emoji
- When user hover their mouse cursor over the button, the button is slightly highlighted (or react in any other reasonable way to emphasise interactivity. At your discrecion)
- When the user clicks the button, the respective tag is completely removed from the system and from the "Tags List" element

### Non-Functional Requirements
- When the tag is removed, all the corresponding records from "tasks" and "TasksTags" db tables must be removed