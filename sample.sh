#!/bin/bash

# This script is a series of curl requests to populate the data in this survey.

# Add demographic question set
curl -X POST "http://localhost:8080/api/v1/question-sets" \
-H "Content-Type: application/json" \
-d '{
    "technologyName": "demographic",
    "questions": [
        {
            "category": "personal",
            "text": "Full Name:",
            "options": []
        },
        {
            "category": "professional",
            "text": "Current Job Title/Role:",
            "options": []
        },
        {
            "category": "education",
            "text": "List any Certifications or Training Programs you have completed in the IT field",
            "options": []
        },
        {
            "category": "organization",
            "text": "Size of your organization:",
            "options": [
                "Small (1-50 employees)",
                "Medium (51-200 employees)",
                "Large (201+ employees)"
            ]
        },
        {
            "category": "industry",
            "text": "Industry of your organization:",
            "options": [
                "Finance",
                "Healthcare",
                "Education",
                "Technology",
                "Government",
                "Non-profit",
                "Other (please specify)"
            ]
        }
    ]
}'
