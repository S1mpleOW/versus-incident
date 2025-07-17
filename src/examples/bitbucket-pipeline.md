# Bitbucket Pipeline Integration

This document describes how to integrate Bitbucket Pipelines with the Versus Incident notification system.

## Overview

The Versus Incident system now supports receiving notifications from Bitbucket Pipelines. When a build completes (success, failure, or other states), Bitbucket can send a webhook to notify relevant channels configured in your system.

## Endpoint

**POST** `/api/notifications/bitbucket`

This is a dedicated endpoint for Bitbucket pipeline notifications that validates the payload structure and processes build state information.

## Payload Structure

The endpoint expects the following JSON payload structure:

```json
{
  "type": "build",
  "key": "BUILD_NUMBER",
  "state": "BUILD_STATE", 
  "name": "BUILD_NAME",
  "url": "BUILD_URL",
  "description": "BUILD_DESCRIPTION"
}
```

### Fields

- **type** (required): Must be "build" for build notifications
- **key** (required): Unique identifier for the build (e.g., build number)
- **state** (required): Build state (SUCCESSFUL, FAILED, INPROGRESS, STOPPED)
- **name** (required): Name/title of the build
- **url** (required): URL to the build details page
- **description** (required): Description of the build or changes

## Bitbucket Pipeline Configuration

Add the following step to your `bitbucket-pipelines.yml` to send notifications:

```yaml
pipelines:
  default:
    - step:
        name: Build and Test
        script:
          # Your build steps here
          - echo "Building application..."
          
    - step:
        name: Send Notification
        script:
          - |
            # Determine status based on previous step result
            if [ $BITBUCKET_EXIT_CODE -eq 0 ]; then
              STATUS="SUCCESSFUL"
              DESC="Build completed successfully"
            else
              STATUS="FAILED" 
              DESC="Build failed with errors"
            fi
            
            # Send notification to Versus Incident
            curl -X POST "https://your-versus-incident-host/api/notifications/bitbucket" \
              -H "Content-Type: application/json" \
              -d "{\"type\": \"build\", \"key\": \"$BITBUCKET_BUILD_NUMBER\", \"state\": \"$STATUS\", \"name\": \"Model Alias changed\", \"url\": \"$BUILD_URL\", \"description\": \"$DESC\"}"
        after-script:
          - |
            # Send failure notification if step failed
            if [ $BITBUCKET_EXIT_CODE -ne 0 ]; then
              curl -X POST "https://your-versus-incident-host/api/notifications/bitbucket" \
                -H "Content-Type: application/json" \
                -d "{\"type\": \"build\", \"key\": \"$BITBUCKET_BUILD_NUMBER\", \"state\": \"FAILED\", \"name\": \"Model Alias changed\", \"url\": \"$BUILD_URL\", \"description\": \"Build step failed\"}"
            fi
```

## Using Bitbucket Variables

Bitbucket provides several built-in variables you can use:

- `$BITBUCKET_BUILD_NUMBER`: Unique build number
- `$BITBUCKET_REPO_FULL_NAME`: Repository name
- `$BITBUCKET_BRANCH`: Current branch name
- `$BITBUCKET_COMMIT`: Commit hash
- `$BITBUCKET_PR_ID`: Pull request ID (if applicable)

Example with more context:

```bash
curl -X POST "https://your-versus-incident-host/api/notifications/bitbucket" \
  -H "Content-Type: application/json" \
  -d "{
    \"type\": \"build\",
    \"key\": \"$BITBUCKET_BUILD_NUMBER\",
    \"state\": \"$STATUS\",
    \"name\": \"$BITBUCKET_REPO_FULL_NAME - $BITBUCKET_BRANCH\",
    \"url\": \"https://bitbucket.org/$BITBUCKET_REPO_FULL_NAME/addon/pipelines/home#!/results/$BITBUCKET_BUILD_NUMBER\",
    \"description\": \"Build for commit $BITBUCKET_COMMIT\"
  }"
```

## Build States

The system recognizes the following build states:

- **SUCCESSFUL**: Build completed successfully (mapped to "info" severity)
- **FAILED**: Build failed (mapped to "critical" severity)  
- **INPROGRESS**: Build is currently running (mapped to "warning" severity)
- **STOPPED**: Build was stopped/cancelled (mapped to "critical" severity)

## Notification Channels

The notification will be sent to all configured channels in your Versus Incident configuration:

- Slack
- Microsoft Teams
- Email
- Telegram
- Lark
- Google Chat
- Viber

## Query Parameters

You can override configuration values using query parameters:

```bash
POST /api/notifications/bitbucket?slack_channel=builds&email_to=devops@company.com
```

## Response

### Success Response (201 Created)
```json
{
  "status": "Bitbucket notification sent successfully",
  "build": "123",
  "state": "SUCCESSFUL"
}
```

### Error Response (400 Bad Request)
```json
{
  "error": "Missing required field: state"
}
```

### Error Response (500 Internal Server Error)
```json
{
  "error": "failed to send notifications: provider error"
}
```

## Testing

You can test the endpoint using curl:

```bash
curl -X POST "http://localhost:8080/api/notifications/bitbucket" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "build",
    "key": "123",
    "state": "SUCCESSFUL", 
    "name": "Test Build",
    "url": "https://bitbucket.org/test/builds/123",
    "description": "Test notification"
  }'
```

## Security Considerations

- Consider implementing authentication/authorization for the webhook endpoint
- Use HTTPS for production deployments
- Validate the source IP if Bitbucket provides webhook IP ranges
- Consider rate limiting to prevent abuse
