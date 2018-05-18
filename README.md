converts trello cards into jira epics and tasks

Quality of Life tips:
Add following env vars to your bashrc or equivalent
- JIRA_USER
- JIRA_PASS
- TRELLO_APP_KEY
- TRELLO_TOKEN


running:
- `jello --out` - converts trello cards and saves them into JSON file
example single json
```
{
    "story_name": "[8] As a developer, I want a mobile configuration that does not have duplication of values so that I am clear on which values I should use within my mobile client",
    "story_points": "8",
    "labels": [
      "F-MobileClientConfig",
      "UI",
      "CLI",
      "Targeted"
    ],
    "checklists": {
      "Acceptance Criterea": [
        "No unnecessary duplication of configuration values",
        "Documented format"
      ],
      "Tasks": [
        "Review the current config format and remove any unneeded duplication (ensure to share with the other service teams to make sure we are not breaking everything on them)",
        "Document the config format",
        "update the config format once we have interacted with the other stakeholders"
      ]
    },
    "card_url": "https://trello.com/c/P0KhgUBk",
    "attachments": [],
    "description": ""
  },
```

add `jira-host` `--jira-user` `--jira-pass` `--trello-token` flags if needed