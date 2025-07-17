# Bitbucket Google Chat Template Example

This shows how the Google Chat template will render Bitbucket pipeline notifications.

## Sample Bitbucket Payload

```json
{
  "type": "build",
  "key": "123",
  "state": "SUCCESSFUL",
  "name": "Model Alias changed",
  "url": "https://bitbucket.org/company/repo/addon/pipelines/home#!/results/123",
  "description": "Build completed successfully with all tests passing"
}
```

## Processed Notification Data

After processing by the notification service, this becomes:

```json
{
  "title": "Bitbucket Build SUCCESSFUL",
  "description": "Build completed successfully with all tests passing",
  "build_key": "123",
  "build_name": "Model Alias changed", 
  "build_state": "SUCCESSFUL",
  "build_url": "https://bitbucket.org/company/repo/addon/pipelines/home#!/results/123",
  "severity": "info",
  "timestamp": 1673980800
}
```

## Expected Google Chat Card Output

The template should render a Google Chat card with:

### Header
- **Title**: "✅ RESOLVED: ✅ Model Alias changed - SUCCESSFUL"
- **Subtitle**: "ℹ️ Severity: INFO (Bitbucket)"

### Body Sections
- **Resource**: "Build #123"
- **Description**: "Build completed successfully with all tests passing"
- **Time**: "2023-01-17 16:00:00" (or actual timestamp)

### Build Details Section
- **Build Details:**
  - • Build Number: 123
  - • Pipeline: Model Alias changed
  - • Status: SUCCESSFUL

### Buttons
- **View Build** (links to the Bitbucket build URL)

## Template Features for Bitbucket

1. **Source Detection**: Detects Bitbucket payloads by presence of `build_key`, `build_state`, and `build_name` fields
2. **Build State Icons**: 
   - ✅ SUCCESSFUL
   - ❌ FAILED  
   - ⏳ INPROGRESS
   - ⏹️ STOPPED
   - 🔧 Default/Unknown
3. **Severity Mapping**:
   - FAILED/STOPPED → critical (🔴)
   - SUCCESSFUL → info (ℹ️) 
   - INPROGRESS → warning (🟡)
4. **Status Mapping**:
   - SUCCESSFUL → RESOLVED (✅)
   - FAILED/STOPPED/INPROGRESS → FIRING (🔥)
5. **Build-Specific Details**: Shows build number, pipeline name, and current status
6. **Contextual Button**: "View Build" instead of generic "Diagnostics"

## Configuration

The template will automatically handle Bitbucket notifications when:
- Google Chat is enabled in the configuration
- The notification contains the required Bitbucket fields
- The template path is correctly set to the updated `ggchat_message.tmpl`

## Testing

To test the Bitbucket notification:

```bash
curl -X POST "http://localhost:8080/api/notifications/bitbucket" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "build",
    "key": "123", 
    "state": "SUCCESSFUL",
    "name": "Model Alias changed",
    "url": "https://bitbucket.org/company/repo/addon/pipelines/home#!/results/123",
    "description": "Build completed successfully"
  }'
```

This will trigger a Google Chat notification using the updated template with Bitbucket-specific formatting.
