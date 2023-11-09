#!/bin/bash

# This script is a series of curl requests to populate the data in this survey.

# Add demographic question set
curl -X POST "http://localhost:8080/api/v1/question-sets" \
-H "Content-Type: application/json" \
-d '{
    "technology_name": "demographic",
    "questions": [
        {
            "category": "personal",
            "text": "Full Name:"
        },
        {
            "category": "professional",
            "text": "Current Job Title/Role:"
        },
        {
            "category": "education",
            "text": "List any Certifications or Training Programs you have completed in the IT field"
        },
        {
            "category": "organization",
            "text": "Size of your organization:",
            "options": {
                "multi_answer": false,
                "possible_options": [
                    "Small (1-50 employees)",
                    "Medium (51-200 employees)",
                    "Large (201+ employees)"
                ]
            }
        },
        {
            "category": "industry",
            "text": "Industry of your organization:",
            "options": {
                "multi_answer": false,
                "possible_options": [
                    "Finance",
                    "Healthcare",
                    "Education",
                    "Technology",
                    "Government",
                    "Non-profit",
                    "Other (please specify)"
                ]
            }
        }
    ]
}'

# Add AWS question set

curl --location 'http://localhost:8080/api/v1/question-sets' \
--header 'Content-Type: application/json' \
--data '{
    "technology_name": "AWS",
    "questions": [
        {
            "category": "background_knowledge",
            "text": "Can you describe what background you have in information technology, programming, or cloud computing? If you'\''re a beginner, that'\''s also ok!"
        },
        {
            "category": "background_knowledge",
            "text": "What prior experience have you had with AWS?"
        },
        {
            "category": "job_role_responsibilities",
            "text": "What is your job title?"
        },
        {
            "category": "job_role_responsibilities",
            "text": "Can you provide a brief description of your job role? What do you do on a day-to-day basis?"
        },
        {
            "category": "current_skill_level",
            "text": "How would you rate your current understanding of cloud computing concepts?",
            "options": {
                "multi_answer": false,
                "possible_options": [
                    "I can spell AWS",
                    "I sell AWS, but don'\''t use it",
                    "I run a small production workload on AWS",
                    "I run a large production workload on AWS",
                    "I run multiple production workloads on AWS"
                ]
            }
        },
        {
            "category": "current_skill_level",
            "text": "Are you familiar with any programming languages? If yes, please specify."
        },
        {
            "category": "current_skill_level",
            "text": "Have you had any previous training or experience with AWS? If yes, please specify the areas (e.g., EC2, S3, Lambda, etc.)."
        },
        {
            "category": "learning_objectives",
            "text": "What are your learning objectives for this training? (Check all that apply)",
            "options": {
                "multi_answer": true,
                "possible_options": [
                    "Understanding the fundamentals of AWS",
                    "Learning how to architect on AWS",
                    "Developing applications on AWS",
                    "Learning how to operate cloud infrastructure on AWS",
                    "Others (please specify)"
                ]
            }
        },
        {
            "category": "learning_objectives",
            "text": "Are there specific AWS services or features you are particularly interested in learning about?"
        },
        {
            "category": "learning_objectives",
            "text": "How do you plan to apply the skills acquired from this training in your work? Are you working on any specific projects or workloads?"
        },
        {
            "category": "workload_profiling",
            "text": "What types of workloads are you currently managing or planning to manage on AWS? (e.g., web applications, data analytics, etc.)"
        },
        {
            "category": "workload_profiling",
            "text": "Can you describe the current or planned architecture of your AWS workloads?"
        },
        {
            "category": "workload_profiling",
            "text": "Are there any performance, security, or cost-optimization requirements for your AWS workloads?"
        },
        {
            "category": "workload_profiling",
            "text": "Are you using or planning to use any automation or Infrastructure as Code (IaC) tools for managing your AWS workloads?"
        },
        {
            "category": "workload_profiling",
            "text": "Are you interested in learning about best practices for monitoring and optimizing AWS workloads?"
        },
        {
            "category": "workload_profiling",
            "text": "What challenges, if any, are you currently facing or anticipate facing with managing workloads on AWS?"
        }
    ]
}
'
