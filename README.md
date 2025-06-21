
## AWS EBS Cost Optimizer (Go)

A cost-saving AWS Lambda function written in Go that automatically detects unattached EBS volumes and deletes associated snapshots. It notifies you via Amazon SNS.

This version uses the Go custom runtime by packaging a bootstrap binary.

---

## 📌 Features

* ✅ Detects unattached EBS volumes
* ✅ Sends alerts via SNS
* ✅ Deletes unused snapshots
* ✅ Scheduled via CloudWatch Events
* ✅ Written in Go with a custom runtime (bootstrap)

---

## 🧱 Requirements

* AWS Account
* IAM Role with appropriate permissions
* An SNS Topic for notifications
* Go 1.20+
* AWS CLI (optional)
* Access to AWS Console

---

## ⚙️ Architecture Overview

```text
CloudWatch Event (Daily Schedule)
         |
         v
     AWS Lambda (Go)
         |
         +---> Describe EBS Volumes
         |
         +---> For Unattached Volumes:
                   |
                   +--> Send SNS Notification
                   |
                   +--> Describe Snapshots
                           |
                           +--> Delete if any
```

## 🚀 Deploy via AWS Console (Step-by-step)

This section shows how to deploy using only the AWS Management Console.

### 1️⃣ Build Go Binary for Linux

From your project directory, run:

```bash
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
```

This generates a Linux-compatible executable named bootstrap.

> ⚠️ Lambda requires the file to be named bootstrap when using a custom runtime.

### 2️⃣ Zip the binary

```bash
zip bootstrap.zip bootstrap
```

### 3️⃣ Go to Lambda Console

* Open the AWS Console → Lambda → Create function
* Choose: Author from scratch

🔽 Fill in the following:

* Function name: ebs-cost-optimizer
* Runtime: ⚠️ Choose “Provide your own bootstrap (Custom runtime)”
* Architecture: x86\_64
* Execution role: Choose an existing role with EC2, SNS, and STS permissions

Click Create Function

### 4️⃣ Upload Deployment Package

* Under Code → Click Upload from → .zip file
* Upload your bootstrap.zip
* Click Save or Deploy

### 5️⃣ Add Environment Variable

Go to Configuration → Environment Variables → Add:

| Key             | Value                                            |
| --------------- | ------------------------------------------------ |
| SNS\_TOPIC\_ARN | arn\:aws\:sns:<region>:<account-id>:<topic-name> |

Click Save

### 6️⃣ Set up Test Event (Optional)

* Click Test → Configure test event
* Use an empty JSON: {}

```json
{}
```

* Name it: test-event
* Click Test to verify output in the logs

---

## ⏰ Set Up Daily Schedule (Console)

To schedule this Lambda to run daily:

* Go to Amazon CloudWatch → Rules → Create Rule
* Source: EventBridge (Schedule)

  * Choose Fixed rate (1 day) or cron(0 5 \* \* ? \*)
* Target: Lambda function → Choose your function
* Click Create

Now your function will run daily and handle cleanup automatically.

---

## 🔐 IAM Policy Required

Attach the following permissions to your Lambda role:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ec2:DescribeVolumes",
        "ec2:DescribeSnapshots",
        "ec2:DeleteSnapshot"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": "sns:Publish",
      "Resource": "arn:aws:sns:<region>:<account-id>:<topic-name>"
    },
    {
      "Effect": "Allow",
      "Action": "sts:GetCallerIdentity",
      "Resource": "*"
    }
  ]
}
```

---

## 🧼 What It Cleans

* Unattached EBS volumes → Sends notification
* Snapshots of those volumes → Deletes them silently

Note: It does not delete the EBS volumes themselves.

---

## 📁 Directory Structure

```bash
aws-ebs-cost-optimizer/
├── main.go               # Lambda handler
├── bootstrap             # Built Go binary (for Lambda)
├── bootstrap.zip          # Deployment package
├── go.mod
├── README.md
```

---

## 🧠 Why Use Custom Runtime?

Because AWS does not officially support Go as a managed runtime since 2023 (Go1.x is deprecated), you must use a custom runtime by naming the binary bootstrap.

This gives you full control and flexibility with minimal overhead.

---

## 📊 Future Enhancements

* Delete unused volumes (optional toggle)
* Generate cost savings report
* Track history in DynamoDB
* Send daily summary to Slack or email

---
