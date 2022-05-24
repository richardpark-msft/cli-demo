To run a demo:

In the Azure Portal:

1. Create a Service Bus Namespace
2. Create a queue called 'demo'.
3. Click "Shared Access Policies"
4. Click "RootManageSharedAccessKey"
5. Copy the primary connection string and write it into a .env file with a key called SERVICEBUS_CONNECTION_STRING.

Run (from this folder):

```powershell
# a simple receiver that uses the messages in samples\*.json
.\samples\samples.exe
```

Sending a message:

```powershell
type samples\sample_message_earth.json | .\sb.exe send demo --env
type samples\sample_message_mars.json | .\sb.exe send demo --env
type samples\sample_message_unknown.json | .\sb.exe send demo --env
```

Receiving a message using the CLI:

```powershell
.\sb receive demo --env
```

Peeking a message using the CLI:

```powershell
.\sb peek demo --env
```
