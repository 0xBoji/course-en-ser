#!/bin/bash

# Setup CloudWatch Alarms for Course Enrollment Service
# Run this script to create monitoring alarms

INSTANCE_ID="i-04fc814a42bac1d7a"
REGION="ap-southeast-2"
SNS_TOPIC_ARN="arn:aws:sns:ap-southeast-2:975050162743:course-enrollment-alerts"

echo "🚨 Setting up CloudWatch Alarms for Course Enrollment Service"

# Create SNS Topic for alerts (if not exists)
echo "📧 Creating SNS topic for alerts..."
aws sns create-topic --name course-enrollment-alerts --region $REGION

# High CPU Utilization Alarm
echo "⚠️ Creating High CPU Alarm..."
aws cloudwatch put-metric-alarm \
    --alarm-name "CourseEnrollment-HighCPU" \
    --alarm-description "High CPU utilization on course enrollment instance" \
    --metric-name cpu_usage_user \
    --namespace CourseEnrollment/Production \
    --statistic Average \
    --period 300 \
    --threshold 80 \
    --comparison-operator GreaterThanThreshold \
    --evaluation-periods 2 \
    --alarm-actions $SNS_TOPIC_ARN \
    --dimensions Name=InstanceId,Value=$INSTANCE_ID \
    --region $REGION

# High Memory Utilization Alarm  
echo "💾 Creating High Memory Alarm..."
aws cloudwatch put-metric-alarm \
    --alarm-name "CourseEnrollment-HighMemory" \
    --alarm-description "High memory utilization on course enrollment instance" \
    --metric-name mem_used_percent \
    --namespace CourseEnrollment/Production \
    --statistic Average \
    --period 300 \
    --threshold 85 \
    --comparison-operator GreaterThanThreshold \
    --evaluation-periods 2 \
    --alarm-actions $SNS_TOPIC_ARN \
    --dimensions Name=InstanceId,Value=$INSTANCE_ID \
    --region $REGION

# High Disk Utilization Alarm
echo "💿 Creating High Disk Alarm..."
aws cloudwatch put-metric-alarm \
    --alarm-name "CourseEnrollment-HighDisk" \
    --alarm-description "High disk utilization on course enrollment instance" \
    --metric-name disk_used_percent \
    --namespace CourseEnrollment/Production \
    --statistic Average \
    --period 300 \
    --threshold 90 \
    --comparison-operator GreaterThanThreshold \
    --evaluation-periods 1 \
    --alarm-actions $SNS_TOPIC_ARN \
    --dimensions Name=InstanceId,Value=$INSTANCE_ID \
    --region $REGION

# Application Error Rate Alarm
echo "🔥 Creating Application Error Alarm..."
aws logs put-metric-filter \
    --log-group-name course-enrollment-app \
    --filter-name ErrorCount \
    --filter-pattern "ERROR" \
    --metric-transformations \
        metricName=ApplicationErrors,metricNamespace=CourseEnrollment/Production,metricValue=1 \
    --region $REGION

aws cloudwatch put-metric-alarm \
    --alarm-name "CourseEnrollment-HighErrors" \
    --alarm-description "High error rate in course enrollment application" \
    --metric-name ApplicationErrors \
    --namespace CourseEnrollment/Production \
    --statistic Sum \
    --period 300 \
    --threshold 5 \
    --comparison-operator GreaterThanThreshold \
    --evaluation-periods 1 \
    --alarm-actions $SNS_TOPIC_ARN \
    --treat-missing-data notBreaching \
    --region $REGION

# Instance Status Check Alarm
echo "❤️ Creating Instance Health Alarm..."
aws cloudwatch put-metric-alarm \
    --alarm-name "CourseEnrollment-InstanceStatusCheck" \
    --alarm-description "Instance status check failed" \
    --metric-name StatusCheckFailed_Instance \
    --namespace AWS/EC2 \
    --statistic Maximum \
    --period 60 \
    --threshold 0 \
    --comparison-operator GreaterThanThreshold \
    --evaluation-periods 2 \
    --alarm-actions $SNS_TOPIC_ARN \
    --dimensions Name=InstanceId,Value=$INSTANCE_ID \
    --region $REGION

echo "✅ CloudWatch Alarms setup completed!"
echo ""
echo "📊 Created Alarms:"
echo "  - CourseEnrollment-HighCPU (>80%)"
echo "  - CourseEnrollment-HighMemory (>85%)"
echo "  - CourseEnrollment-HighDisk (>90%)"
echo "  - CourseEnrollment-HighErrors (>5 errors/5min)"
echo "  - CourseEnrollment-InstanceStatusCheck"
echo ""
echo "📧 Subscribe to SNS topic to receive alerts:"
echo "aws sns subscribe --topic-arn $SNS_TOPIC_ARN --protocol email --notification-endpoint your-email@example.com"
