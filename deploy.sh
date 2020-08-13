# λ deployment script
LAMBDA_NAME="slackApp"

echo "STEP 0: >>> Test λ $LAMBDA_NAME >>>"
TEST=$(go test ./... | grep FAIL | head -n 1)


if [[ -z "$TEST" ]]
then
    echo "Test is [PASSED]"
else
    echo "Test $TEST is [FAILED]"
    exit 1
fi

echo "STEP 1: >>> Build λ $LAMBDA_NAME >>>"
go build -o /tmp/$LAMBDA_NAME

echo "STEP 2: >>> Zipping λ $LAMBDA_NAME >>>"
zip -j /tmp/$LAMBDA_NAME.zip /tmp/$LAMBDA_NAME

echo "STEP 3: >>> Updating lambda λ $LAMBDA_NAME code"
aws lambda update-function-code --function-name $LAMBDA_NAME --zip-file fileb:///tmp/$LAMBDA_NAME.zip