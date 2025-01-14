User Management Mircoservice with 4 API endpoints
1. Signup
2. Login
3. Get User By Id
4. Get All Users

This microservice authenticates other microservices using JWT tokens and AWS Cognito.

Each Request and Response will have a JWT token in the header so that the microservice will verify it with AWS Cognito.

Stores all the details in Mongo DB and retrives the details when the API's are called.

This microservice works hand in hand with other microservices which need authentication for the users.


