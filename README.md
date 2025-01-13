User Management Mircoservice with 4 API endpoints
1. Signup
2. Login
3. Get User By Id
4. Get All Users

This microservice authenticates other microservices using JWT tokens and AWS Cognito.
Stores all the details in Mongo DB and retrives the details when the API's are hit.
This is a microservice which can be called by other microservices to authorize the Users they are specifying.
