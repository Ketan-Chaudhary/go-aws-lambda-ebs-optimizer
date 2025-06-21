# AWS EBS Cost Optimizer (Go)

A cost-saving AWS Lambda function written in **Go** that automatically detects **unattached EBS volumes** and deletes **associated snapshots**. It notifies you via **Amazon SNS**.

This version uses the **Go custom runtime** by packaging a `bootstrap` binary.

---

## 📌 Features

✅ Detects unattached EBS volumes
✅ Sends alerts via SNS
✅ Deletes unused snapshots
✅ Scheduled via **Amazon EventBridge**
✅ Written in Go with a custom runtime (`bootstrap`)

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

```
Amazon EventBridge (Daily Schedule)
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

---

## 🚀 Deploy via AWS Console (Step-by-step)

### 1️⃣ Build Go Binary for Linux

From your project directory, run:

```bash
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
```

⚠️ Lambda requires the file to be named **`bootstrap`** when using a custom runtime.

---

### 2️⃣ Zip the Binary

```bash
zip bootstrap.zip bootstrap
```

---

### 3️⃣ Go to Lambda Console

Open the AWS Console → Lambda → Create function

Choose: **Author from scratch**

🔽 Fill in the following:

* **Function name:** `ebs-cost-optimizer`
* **Runtime:** ⚠️ Choose **“Provide your own bootstrap (Custom runtime)”**
* **Architecture:** `x86_64`
* **Execution role:** Choose an existing role with EC2, SNS, and STS permissions

Click **Create Function**

---

### 4️⃣ Upload Deployment Package

Under **Code** → Click **Upload from → .zip file**
Upload your `bootstrap.zip`
Click **Save** or **Deploy**

---

### 5️⃣ Add Environment Variable

Go to **Configuration → Environment Variables** → Add:

| Key             | Value                        |
| --------------- | ---------------------------- |
| `SNS_TOPIC_ARN` | `arn:aws:sns:::<your-topic>` |

Click **Save**

---

### 6️⃣ Set Up Test Event (Optional)

Click **Test → Configure test event**
Use an empty JSON:

```json
{}
```

Name it: `test-event`
Click **Test** to verify output in the logs

---

### ⏰ Set Up Daily Schedule (EventBridge)

To schedule this Lambda to run daily:

1. Go to **Amazon EventBridge → Scheduler → Create schedule**
2. Choose:

   * **Schedule type:** Recurring schedule
   * **Fixed rate:** `1 day`
     or

     * **Cron expression:** `cron(0 5 * * ? *)`
3. **Target:**

   * Select **Lambda function**
   * Choose your function: `ebs-cost-optimizer`

Click **Create Schedule**

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

* **Unattached EBS volumes** → Sends SNS notification
* **Snapshots of those volumes** → Silently deletes

> **Note:** It does **not delete** the EBS volumes themselves.

---

## 📁 Directory Structure

```
aws-ebs-cost-optimizer/
├── main.go               # Lambda handler
├── bootstrap             # Built Go binary (for Lambda)
├── bootstrap.zip         # Deployment package
├── go.mod
├── README.md
```

---

## 🧠 Why Use Custom Runtime?

Since AWS deprecated official support for Go (Go1.x) as a managed runtime in 2023, a **custom runtime** is required. By naming your binary `bootstrap`, you gain:

* Full control
* Minimal overhead
* Long-term flexibility

---

## 📊 Future Enhancements

* [ ] Delete unused volumes (optional toggle)
* [ ] Generate cost savings report
* [ ] Track history in DynamoDB
* [ ] Send daily summary to Slack or email

