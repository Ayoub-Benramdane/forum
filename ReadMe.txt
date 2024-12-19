Forum Project
Overview
This is a web-based forum that allows users to register, create posts, comment, like/dislike content, and filter posts based on categories, user activity, and likes. The project uses SQLite for database management and Docker for containerization.

Features
User Authentication: Register and log in with email and password (password encryption is a bonus task).
Posts & Comments: Registered users can create posts and comments, while all users can view them.
Likes/Dislikes: Registered users can like or dislike posts and comments.
Filters: Users can filter posts by categories, created posts, and liked posts (for registered users).
Database
SQLite stores data for users, posts, comments, categories, and likes/dislikes.
Tech Stack
Go (Golang): Backend development.
SQLite: Database for storing data.
Docker: Containerization for deployment.
bcrypt: Password hashing (bonus task).
UUID: Session ID (bonus task).
Setup Instructions
Build & Run with Docker:

bash
Copy code
docker build -t forum-app .
docker run -p 8080:8080 forum-app
The app will be available at http://localhost:8080.

User Actions:

Register with email, username, and password.
Log in to create posts and interact with content.
Filter posts by categories or user activity.
License
This project is licensed under the MIT License.

