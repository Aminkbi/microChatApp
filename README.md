To start the application, follow these steps:

Open a terminal and navigate to the root directory of your project.

Ensure Docker is running on your machine.

Run the following command to start the application with Docker Compose:


`docker-compose up --build`



provide `.env` file including these parameters:


`MONGODB_URI="mongodb://username:password@mongodb:27017/micro-chat?authSource=admin"`

`ENVIRONMENT=dev|prod`

`APP_PORT=8080`

`JWT_SECRET={secret}`

`WEBSOCKET_PORT=4000`

`ASYNQ_PORT=4001`

`MONGO_INITDB_ROOT_USERNAME={mongo-username}`

`MONGO_INITDB_ROOT_PASSWORD={mongo-password}`