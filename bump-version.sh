#!/bin/bash

# Get the latest tag
latest_tag=$(git describe --tags `git rev-list --tags --max-count=1` 2>/dev/null)

if [ -z "$latest_tag" ]; then
    echo "No existing tags found. Creating initial version v0.1.0"
    new_tag="v0.1.0"
else
    # Extract version components
    version=${latest_tag#v}  # Remove 'v' prefix
    IFS='.' read -r major minor patch <<< "$version"
    
    # Increment minor version
    new_minor=$((minor + 1))
    
    # Create new version string
    new_tag="v$major.$new_minor.0"
fi

# Create and push new tag
echo "Creating new tag: $new_tag"
git tag -a "$new_tag" -m "Release $new_tag"
git push origin "$new_tag"

echo "Successfully created and pushed tag $new_tag"
