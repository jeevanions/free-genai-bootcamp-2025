As a experience backend api developer in golang, write a backend api based in golang following the backend tech spec @backend-spec.md

Instructions: 

1. Go throug the backend tech spec provided and write a backend api based on it.
2. First create the project structure
3. Create the database schema and migration script
4. Implement one api endpoint at a time.
5. For each api endpoint 
   - Implement different layer of abstractions like data model, routes, services, handlers and repositories.
   - Always adhere to project structure and clean architecture principles.
   - Use proper error handling and logging.
   - Use testing to ensure code quality.
   - Write unit tests for the endpoint just implemented. 
6. Update swagger docs
7. Implement the remaining endpoints.



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