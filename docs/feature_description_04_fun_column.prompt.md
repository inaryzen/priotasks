# Task Field Enhancement: Fun Factor Implementation

## Overview
Add a new "Fun" field to tasks to quantify user enjoyment level with corresponding visual indicators.

## Technical Specifications

### Data Model Changes
- Field Name: `Fun`
- Type: Enum
- Valid Values: `S`, `M`, `L`, `XL`
- Default Value: `M`

### Visual Requirements
- Display 🍀 emoji indicator next to Fun value
- Emojis mapping:
  - S - 🍀
  - M - 🍀🍀
  - L -🍀🍀🍀
  - XL - 🍀🍀🍀🍀

### UI Updates
1. Task Table:
   - Add "Fun" column after "Value" column
   - Display format: `[Value] - [Emoji]`

2. Task Modal:
   - Reorganize form layout:
     1. Title field (top)
     2. Group all select elements below title:
        - Priority
        - Impact
        - Cost
        - Fun (new)
     3. Remaining fields

### Business Logic
- Update Value calculation formula:
  - Use Priority multipliers for Fun values
  - Multiplier mapping:
    - S: 0.75x
    - M: 1.0x
    - L: 1.25x
    - XL: 1.5x

## Implementation Steps
1. Add database migration in `/db/dbtasks.go`
2. Create value recalculation migration in `services/migration.go`
3. Update task model and DTO
4. Implement UI changes
5. Update value calculation logic
6. Add validation rules

## Validation Criteria
- Fun field is required for all tasks
- Only predefined enum values are accepted
- Value calculation correctly incorporates Fun multiplier
- UI displays correct emoji for each Fun value
- All existing tasks have valid Fun values after migration