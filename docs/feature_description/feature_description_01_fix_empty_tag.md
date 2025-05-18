# Feature Description Document - F01

## Overview
It is possible for the user to select an option "Select a tag..." in Tags filter in FilterPanel. When that happens an empty string is sent as a filter for tags. This leads to the situation when an empty string is saved in Settings. This should be prevented on all the levels.

## Requirements
### Functional Requirements
- It must not be possible to select an option "Select a tag..." from the user interface
- Handler level functions should not accept an empty string as a tag
- Service layer functions must not accecpt an empty string as a tag