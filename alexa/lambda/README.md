## Install

1. Make python file golbally readable: ```chmod u=rwx,go=r <file>```
2. Zip file: ```zip handler.zip lambda/handler.py```
3. Upload function:
```
aws lambda create-function \
  --function-name roku-control \
  --runtime python2.7 --handler handler.handle --zip-file fileb://handler.zip \
  --role arn:aws:iam::<ACCOUNT-ID>:role/lambda_basic_execution
```
4. Test function: ```aws lambda invoke --function-name roku-control output.txt```

