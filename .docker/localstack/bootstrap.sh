#!/usr/bin/env bash
set -e

echo "🚀 Starting LocalStack bootstrap..."

# =====================
# SNS TOPICS
# =====================
ORDER_CREATED_TOPIC=order-created-event
PAYMENT_PROCESSED_TOPIC=payment-processed-event
KITCHEN_ORDER_FINISHED_TOPIC=kitchen-order-finished-event
ORDER_ERROR_TOPIC=order-error-topic

# =====================
# SQS QUEUES
# =====================
CREATE_PAYMENT_QUEUE=create-payment-queue
CREATE_KITCHEN_ORDER_QUEUE=create-kitchen-order-queue
UPDATE_ORDER_STATUS_QUEUE=update-order-status-queue

KITCHEN_ORDER_ERROR_QUEUE=kitchen-order-order-error-queue
ORDER_ERROR_QUEUE=order-error-queue
PAYMENT_ORDER_ERROR_QUEUE=payment-order-error-queue

echo "▶ Creating SNS topics..."
ORDER_CREATED_TOPIC_ARN=$(awslocal sns create-topic \
  --name "$ORDER_CREATED_TOPIC" \
  --query 'TopicArn' \
  --output text)

PAYMENT_PROCESSED_TOPIC_ARN=$(awslocal sns create-topic \
  --name "$PAYMENT_PROCESSED_TOPIC" \
  --query 'TopicArn' \
  --output text)

KITCHEN_ORDER_FINISHED_TOPIC_ARN=$(awslocal sns create-topic \
  --name "$KITCHEN_ORDER_FINISHED_TOPIC" \
  --query 'TopicArn' \
  --output text)

ORDER_ERROR_TOPIC_ARN=$(awslocal sns create-topic \
  --name "$ORDER_ERROR_TOPIC" \
  --query 'TopicArn' \
  --output text)

echo "▶ Creating SQS queues..."
CREATE_PAYMENT_QUEUE_URL=$(awslocal sqs create-queue \
  --queue-name "$CREATE_PAYMENT_QUEUE" \
  --query 'QueueUrl' \
  --output text)

CREATE_KITCHEN_ORDER_QUEUE_URL=$(awslocal sqs create-queue \
  --queue-name "$CREATE_KITCHEN_ORDER_QUEUE" \
  --query 'QueueUrl' \
  --output text)

UPDATE_ORDER_STATUS_QUEUE_URL=$(awslocal sqs create-queue \
  --queue-name "$UPDATE_ORDER_STATUS_QUEUE" \
  --query 'QueueUrl' \
  --output text)

KITCHEN_ORDER_ERROR_QUEUE_URL=$(awslocal sqs create-queue \
  --queue-name "$KITCHEN_ORDER_ERROR_QUEUE" \
  --query 'QueueUrl' \
  --output text)

ORDER_ERROR_QUEUE_URL=$(awslocal sqs create-queue \
  --queue-name "$ORDER_ERROR_QUEUE" \
  --query 'QueueUrl' \
  --output text)

PAYMENT_ORDER_ERROR_QUEUE_URL=$(awslocal sqs create-queue \
  --queue-name "$PAYMENT_ORDER_ERROR_QUEUE" \
  --query 'QueueUrl' \
  --output text)

echo "▶ Getting Queue ARNs..."
CREATE_PAYMENT_QUEUE_ARN=$(awslocal sqs get-queue-attributes \
  --queue-url "$CREATE_PAYMENT_QUEUE_URL" \
  --attribute-names QueueArn \
  --query 'Attributes.QueueArn' \
  --output text)

CREATE_KITCHEN_ORDER_QUEUE_ARN=$(awslocal sqs get-queue-attributes \
  --queue-url "$CREATE_KITCHEN_ORDER_QUEUE_URL" \
  --attribute-names QueueArn \
  --query 'Attributes.QueueArn' \
  --output text)

UPDATE_ORDER_STATUS_QUEUE_ARN=$(awslocal sqs get-queue-attributes \
  --queue-url "$UPDATE_ORDER_STATUS_QUEUE_URL" \
  --attribute-names QueueArn \
  --query 'Attributes.QueueArn' \
  --output text)

KITCHEN_ORDER_ERROR_QUEUE_ARN=$(awslocal sqs get-queue-attributes \
  --queue-url "$KITCHEN_ORDER_ERROR_QUEUE_URL" \
  --attribute-names QueueArn \
  --query 'Attributes.QueueArn' \
  --output text)

ORDER_ERROR_QUEUE_ARN=$(awslocal sqs get-queue-attributes \
  --queue-url "$ORDER_ERROR_QUEUE_URL" \
  --attribute-names QueueArn \
  --query 'Attributes.QueueArn' \
  --output text)

PAYMENT_ORDER_ERROR_QUEUE_ARN=$(awslocal sqs get-queue-attributes \
  --queue-url "$PAYMENT_ORDER_ERROR_QUEUE_URL" \
  --attribute-names QueueArn \
  --query 'Attributes.QueueArn' \
  --output text)

echo "▶ Subscribing SQS queues to SNS topics..."

# order-created -> create-payment-queue
awslocal sns subscribe \
  --topic-arn "$ORDER_CREATED_TOPIC_ARN" \
  --protocol sqs \
  --notification-endpoint "$CREATE_PAYMENT_QUEUE_ARN"

# payment-processed -> create-kitchen-order-queue
awslocal sns subscribe \
  --topic-arn "$PAYMENT_PROCESSED_TOPIC_ARN" \
  --protocol sqs \
  --notification-endpoint "$CREATE_KITCHEN_ORDER_QUEUE_ARN"

# payment-processed -> update-order-status-queue
awslocal sns subscribe \
  --topic-arn "$PAYMENT_PROCESSED_TOPIC_ARN" \
  --protocol sqs \
  --notification-endpoint "$UPDATE_ORDER_STATUS_QUEUE_ARN"

# kitchen-order-finished -> update-order-status-queue
awslocal sns subscribe \
  --topic-arn "$KITCHEN_ORDER_FINISHED_TOPIC_ARN" \
  --protocol sqs \
  --notification-endpoint "$UPDATE_ORDER_STATUS_QUEUE_ARN"

# order-error -> error queues
awslocal sns subscribe \
  --topic-arn "$ORDER_ERROR_TOPIC_ARN" \
  --protocol sqs \
  --notification-endpoint "$KITCHEN_ORDER_ERROR_QUEUE_ARN"

awslocal sns subscribe \
  --topic-arn "$ORDER_ERROR_TOPIC_ARN" \
  --protocol sqs \
  --notification-endpoint "$ORDER_ERROR_QUEUE_ARN"

awslocal sns subscribe \
  --topic-arn "$ORDER_ERROR_TOPIC_ARN" \
  --protocol sqs \
  --notification-endpoint "$PAYMENT_ORDER_ERROR_QUEUE_ARN"

echo "✅ LocalStack bootstrap finished successfully"
