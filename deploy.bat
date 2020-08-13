set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -o slackApp main.go
%USERPROFILE%\Go\bin\build-lambda-zip.exe -o slackApp.zip slackApp
aws lambda update-function-code --function-name slackApp --zip-file fileb://C:/fun/code/slackApp/slackApp.zip