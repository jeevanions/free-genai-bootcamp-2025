It seems that we are going in circles to fix the tests and the backend code. So here is what we are going to do now.

As a experience backend api developer in golang, I want you to follow below steps to fix the issues with the backend and the tests:

We are going to take one endpoint at a time from the @backend-tech-spec.md , @frontend-tech-spec.md and then

Current Endpoint:  **POST /api/groups/:id/words**
1. Review the database schema, table structure for the current endpoint and make sure it adhere the technical spec.
2. Review the response for the current endpoint and make sure it adhere the technical spec.
3. Review the models for the current endpoint and make sure it adhere the technical spec.
4. Review the handlers for the current endpoint and make sure it adhere the technical spec.
5. Review the routes for the current endpoint and make sure it adhere the technical spec.
6. Review the services for the current endpoint and make sure it adhere the technical spec.
7. Review the repositories for the current endpoint and make sure it adhere the technical spec.
8. Review the migration for the current endpoint and make sure it adhere the technical spec.
9. For all the above fix any issues you find.
10. Verify the different level of abstractions are implemented to the standards.
11. Make sure the models, classes you create are not duplicate and placed in the right project structure.
12. Before creating any class make sure it is not already present to avoid duplication and if present analyse why then resolve.
13. Update the swagger docs to reflect the changes.
14. Do not create additional endpoint.


14. Verify and modify the unit for the current endpoints
15. Run the unit tests only related to that endpoint and fix issues.  Do not run integration tests now.
Finally we fix the integration testing and we are good to go.