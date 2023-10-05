#!/bin/bash

curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" -H "Origin: https://example-origin.com"  --location 'http://localhost:8080/api/v1/users/185bda41-aed7-427d-a7b7-16b99369e746/messages/answers' \
--header 'Content-Type: application/json' \
--data '[
  {
    "answer": {
      "question": {
        "category": "background_knowledge",
        "text": "Can you describe what background you have in information technology, programming, or cloud computing? If you'\''re a beginner, that'\''s also ok!"
      },
      "technology_name": "AWS",
      "answer": "I have a background in information technology with experience in programming and cloud computing."
    }
  },
  {
    "answer": {
      "question": {
        "category": "background_knowledge",
        "text": "What prior experience have you had with AWS?"
      },
      "technology_name": "AWS",
      "answer": "I have worked on various projects involving AWS services like EC2, S3, and Lambda."
    }
  },
  {
    "answer": {
      "question": {
        "category": "job_role_responsibilities",
        "text": "What is your job title?"
      },
      "technology_name": "AWS",
      "answer": "Cloud Solutions Architect"
    }
  },
  {
    "answer": {
      "question": {
        "category": "job_role_responsibilities",
        "text": "Can you provide a brief description of your job role? What do you do on a day-to-day basis?"
      },
      "technology_name": "AWS",
      "answer": "I design and implement cloud infrastructure solutions, liaising with clients and stakeholders to understand their requirements."
    }
  },
  {
    "answer": {
      "question": {
        "category": "current_skill_level",
        "text": "How would you rate your current understanding of cloud computing concepts?",
        "options": {
          "multi_answer": false,
          "possible_options": ["Beginner", "Intermediate", "Advanced"]
        }
      },
      "technology_name": "AWS",
      "answer": "Advanced"
    }
  },
  {
    "answer": {
      "question": {
        "category": "current_skill_level",
        "text": "Are you familiar with any programming languages? If yes, please specify."
      },
      "technology_name": "AWS",
      "answer": "Yes, I am familiar with Python and JavaScript."
    }
  },
  {
    "answer": {
      "question": {
        "category": "current_skill_level",
        "text": "Have you had any previous training or experience with AWS? If yes, please specify the areas (e.g., EC2, S3, Lambda, etc.)."
      },
      "technology_name": "AWS",
      "answer": "Yes, I have experience with EC2, S3, and Lambda."
    }
  },
  {
    "answer": {
      "question": {
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
      "technology_name": "AWS",
      "answer": "Understanding the fundamentals of AWS, Learning how to operate cloud infrastructure on AWS"
    }
  },
  {
    "answer": {
      "question": {
        "category": "learning_objectives",
        "text": "Are there specific AWS services or features you are particularly interested in learning about?"
      },
      "technology_name": "AWS",
      "answer": "I am interested in learning about AWS Lambda and Amazon S3 in depth."
    }
  },
  {
    "answer": {
      "question": {
        "category": "learning_objectives",
        "text": "Are there any other specific topics or areas of focus you would like this training to cover?"
      },
      "technology_name": "AWS",
      "answer": "I would like to learn about best practices for security and cost optimization on AWS."
    }
  },
  {
    "answer": {
      "question": {
        "category": "workload_profiling",
        "text": "What types of workloads are you currently managing or planning to manage on AWS? (e.g., web applications, data analytics, etc.)"
      },
      "technology_name": "AWS",
      "answer": "Currently managing web applications and data analytics workloads."
    }
  },
  {
    "answer": {
      "question": {
        "category": "workload_profiling",
        "text": "Can you describe the current or planned architecture of your AWS workloads?"
      },
      "technology_name": "AWS",
      "answer": "The current architecture involves using EC2 instances for hosting applications, S3 for storage, and Lambda for serverless computing."
    }
  },
  {
    "answer": {
      "question": {
        "category": "workload_profiling",
        "text": "Are there any performance, security, or cost-optimization requirements for your AWS workloads?"
      },
      "technology_name": "AWS",
      "answer": "Yes, there are requirements for high performance, robust security measures, and optimizing costs by utilizing reserved instances and savings plans."
    }
  },
  {
    "answer": {
      "question": {
        "category": "workload_profiling",
        "text": "Are you using or planning to use any automation or Infrastructure as Code (IaC) tools for managing your AWS workloads?"
      },
      "technology_name": "AWS",
      "answer": "Yes, planning to implement automation using AWS CloudFormation for Infrastructure as Code (IaC)."
    }
  },
  {
    "answer": {
      "question": {
        "category": "workload_profiling",
        "text": "Are you interested in learning about best practices for monitoring and optimizing AWS workloads?"
      },
      "technology_name": "AWS",
      "answer": "Yes, I am keen on learning about best practices for monitoring using CloudWatch and optimizing workloads for better performance."
    }
  },
  {
    "answer": {
      "question": {
        "category": "workload_profiling",
        "text": "What challenges, if any, are you currently facing or anticipate facing with managing workloads on AWS?"
      },
      "technology_name": "AWS",
      "answer": "Some challenges faced are ensuring cost-effectiveness while maintaining high performance and managing complex infrastructures with evolving requirements."
    }
  }
]'