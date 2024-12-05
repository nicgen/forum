# Forum Project - README

## Table of Contents
- [How to Use](#how-to-use)
- [Prerequisites](#prerequisites)
- [Environment Variables Configuration](#environment-variables-configuration)
  - [Creating the `.env` File](#creating-the-env-file)
  - [Environment Variables Description](#environment-variables-description)
  - [Important Notes](#important-notes)
- [Building and Running the Application](#building-and-running-the-application)
   - [Clone the Repository](#clone-the-repository)
   - [Build the Docker Image and Run the Application](#build-the-docker-image-and-run-the-application)
   - [Stopping the Application](#stopping-the-application)
   - [Cleaning Up Resources](#cleaning-up-resources)
- [Additional Commands](#additional-commands)
- [Troubleshooting](#troubleshooting)
- [Project Overview](#project-overview)
- [Technologies Used](#technologies-used)
- [Features](#features)
  - [Authentication](#authentication)
  - [User Interactions](#user-interactions)
  - [Post and Comment Moderation](#post-and-comment-moderation)
  - [Image Upload](#image-upload)
  - [Forum Activity Tracking](#forum-activity-tracking)
  - [Likes and Dislikes](#likes-and-dislikes)
  - [User Roles](#user-roles)
  - [Forum Security](#forum-security)
- [Database Schema](#database-schema)
- [Testing](#testing)
- [Contributing](#contributing)

---

## How to Use

This project is a Go application that can be built and run using Docker. The provided Makefile simplifies the process of managing the application. Follow the steps below to get started.

### Prerequisites

- Ensure you have [Docker](https://www.docker.com/get-started) installed on your machine.
- Make sure you have [Make](https://www.gnu.org/software/make/) installed.

### Environment Variables Configuration

To run the application correctly, you need to create a `.env` file in the root directory of the project. This file will contain all the necessary environment variables required for the application to function properly.

#### Creating the `.env` File

1. **Copy the Example File**: Start by copying the `env_example` file to create your `.env` file:

   ```bash
   cp env_example .env
   ```

2. **Edit the `.env` File**: Open the `.env` file in your preferred text editor and fill in the required values for each environment variable. Below is a description of each variable:

   ```plaintext
   # go to https://ntfy.sh/ for generating a token and receive notifications
   NFTY_TOKEN=<token from NTFY>

   # Password for the database
   DB_PASSWORD=<database password>

   # Admin Information
   USER_UUID=<UUID admin>
   EMAIL=<admin email>
   USERNAME=<admin username>
   ADMIN_PASSWORD=<admin user password>

   # Google App Credentials
   GOOGLE_APP_EMAIL=<google email>
   GOOGLE_APP_PASSWORD=<google password>

   # Google Auth Credentials
   GOOGLE_CLIENT_ID=<google client ID>
   GOOGLE_CLIENT_SECRET=<google client secret>

   # GitHub Auth Credentials
   GITHUB_CLIENT_ID=<github client ID>
   GITHUB_CLIENT_SECRET=<github client secret>

   # Discord Auth Credentials
   DISCORD_CLIENT_ID=<discord client ID>
   DISCORD_CLIENT_SECRET=<discord client secret>
   ```

#### Environment Variables Description

- **NFTY_TOKEN**: Obtain a token from [NTFY](https://ntfy.sh/) to receive notifications.
  
- **DB_PASSWORD**: Set the password for your database. Ensure this is secure.

- **Admin Information**:
  - **USER_UUID**: The UUID for the admin user.
  - **EMAIL**: The email address for the admin user.
  - **USERNAME**: The username for the admin user.
  - **ADMIN_PASSWORD**: The password for the admin user.

- **Google App Credentials**:
  - **GOOGLE_APP_EMAIL**: The email associated with your Google app.
  - **GOOGLE_APP_PASSWORD**: The password for your Google app.

- **Google Auth Credentials**:
  - **GOOGLE_CLIENT_ID**: The client ID for Google authentication.
  - **GOOGLE_CLIENT_SECRET**: The client secret for Google authentication.

- **GitHub Auth Credentials**:
  - **GITHUB_CLIENT_ID**: The client ID for GitHub authentication.
  - **GITHUB_CLIENT_SECRET**: The client secret for GitHub authentication.

- **Discord Auth Credentials**:
  - **DISCORD_CLIENT_ID**: The client ID for Discord authentication.
  - **DISCORD_CLIENT_SECRET**: The client secret for Discord authentication.

#### Important Notes

- Ensure that you do not share your `.env` file publicly, as it contains sensitive information.
- After filling in the required values, save the `.env` file and restart your application for the changes to take effect.

#### Instructions

- This README section provides clear instructions on how to create and configure the `.env` file.
- Adjust any specific details or instructions to fit the context of your project.
- You can include additional sections as needed, such as "Troubleshooting" or "Further Reading," depending on your project's requirements.

---

### Building and Running the Application

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/nicgen/forum.git
   cd forum
   ```
2. Build the Docker Image and Run the Application Use the following command to build the Docker image and run the application:
   ```bash
   make
   ```
   This command will:
   - Build the Docker image.
   - Create and start the Docker container.
3. Stopping the Application To stop the running application, use the following command:
   ```bash
   make stop
   ```
4. Cleaning Up Resources If you want to stop the application and remove the container and image, use the following command:
   ```bash
   make clean
   ```
   This command will:
   - Stop the running container (if any).
   - Remove the container.
   - Remove the Docker image.

### Additional Commands

- View Logs: To view the logs of the running container, you can use:
   ```bash
   docker logs forum
   ```
- Access the Container: If you need to access the container's shell, you can use:
   ```bash
   docker exec -it forum /bin/sh
   ```

### Troubleshooting

If you encounter issues with Docker permissions, you may need to run commands with sudo or add your user to the Docker group.
Ensure that the forum.db file is created in the project directory before running the application.

---

## Project Overview

This is a web-based forum designed to facilitate communication between users. It includes features such as posts and comments, user registration, session management, image uploads, moderation capabilities, activity tracking, and security features. The backend of the project is implemented in Go, and it uses SQLite for data storage. The forum is structured to handle various user interactions like liking, disliking, filtering posts, and managing post categories. In addition, security measures such as HTTPS encryption, rate limiting, and password hashing are incorporated to ensure safe usage.

---

## Technologies Used

- **Backend**: Go (Golang)
- **Database**: SQLite3
- **Authentication**: bcrypt, UUID
- **Sessions and Cookies**: Managed through HTTP cookies with an expiration date
- **Image Uploads**: JPEG, PNG, GIF (up to 20 MB per file)
- **Containerization**: Docker
- **HTTPS**: SSL certificates for encrypted communication
- **API Rate Limiting**: Rate-limiting mechanisms implemented
- **Unit Testing**: Go test files for unit testing
- **OAuth Authentication**: Google and GitHub login

---

## Features

### Authentication
- Users can register and log in using their email, username, and password.
- Passwords are securely encrypted using bcrypt.
- Users can log in via **Google** and **GitHub** authentication for easier access.
- Session management is handled using cookies, with each session having an expiration time.

### User Interactions
- Registered users can create posts and comments.
- Posts can be associated with one or more categories (e.g., "Technology", "General", etc.).
- All users (registered and unregistered) can view posts and comments.
- Only registered users can like or dislike posts and comments.

### Post and Comment Moderation
- Moderators can approve or reject posts before they become publicly visible.
- Moderators have the ability to delete posts and comments based on content.
- A moderation report system allows moderators to report inappropriate content to admins.

### Image Upload
- Registered users can upload images (JPEG, PNG, GIF) along with their posts.
- The maximum allowed file size for uploaded images is 20 MB.
- Images are displayed within posts for both registered users and guests.

### Forum Activity Tracking
- Users can view their activity, including:
  - Posts they have created.
  - Posts they have liked or disliked.
  - Comments they have made, including the content of the comments and the posts they were made on.

### Likes and Dislikes
- Users can like or dislike posts and comments.
- The total count of likes and dislikes for posts and comments is visible to all users.

### User Roles
The forum supports multiple user roles with varying access permissions:
- **Guests**: Can only view posts and comments.
- **Users**: Can create posts, comment, and like/dislike posts and comments.
- **Moderators**: Can approve/reject content, delete posts/comments, and report inappropriate content to admins.
- **Administrators**: Can manage user roles, approve reports from moderators, and manage categories.

### Forum Security
- **HTTPS**: The entire forum is accessible over HTTPS with SSL certificates to ensure encrypted communication.
- **Rate Limiting**: Limits the number of requests a user can make to prevent abuse.
- **Password Encryption**: All user passwords are encrypted using bcrypt.
- **Secure Session Cookies**: Session cookies are unique and tamper-proof, storing only an identifier.

---

## Database Schema

The database is structured using SQLite and the following tables are created to manage users, posts, comments, reactions, and more. The schema supports the key features of the forum, including roles, reactions (likes/dislikes), notifications, post categories, and OAuth authentication.

### Tables

1. **User**
   - Stores information about users.
   - `ID` (INTEGER): Unique auto-incremented identifier.
   - `UUID` (VARCHAR(255)): Unique identifier for each user (used for linking user-related data).
   - `Email` (VARCHAR(50)): The user's email (unique).
   - `Username` (VARCHAR(25)): The user's username (unique).
   - `Password` (VARCHAR(100)): The user's password (encrypted).
   - `OAuthID` (VARCHAR(255), UNIQUE): Stores external OAuth identifiers (for Google, GitHub logins).
   - `Role` (TEXT): The role of the user (Admin, User, Moderator, DeleteUser).
   - `IsLogged` (BOOL): Whether the user is currently logged in.
   - `IsDeleted` (BOOL): Marks the user as deleted (soft delete).
   - `IsRequest` (BOOL): Whether the user has requested a special action (such as becoming a moderator).
   - `CreatedAt` (DATETIME): When the user was created.

2. **Categories**
   - Stores the categories available in the forum.
   - `ID` (INTEGER): Unique auto-incremented identifier.
   - `Name` (VARCHAR(50)): The name of the category (unique).

3. **Posts**
   - Stores the forum posts created by users.
   - `ID` (INTEGER): Unique auto-incremented identifier.
   - `User_UUID` (VARCHAR(255)): The UUID of the user who created the post (foreign key referencing User).
   - `Title` (TEXT): The title of the post.
   - `Category_ID` (INTEGER): The category associated with the post (foreign key referencing Categories).
   - `Text` (TEXT): The content of the post.
   - `ImagePath` (TEXT): The path to an image file associated with the post (if any).
   - `Like` (INTEGER): Count of likes for the post.
   - `Dislike` (INTEGER): Count of dislikes for the post.
   - `CreatedAt` (DATETIME): When the post was created.

4. **Post_Categories**
   - Links posts to multiple categories.
   - `ID` (INTEGER): Unique auto-incremented identifier.
   - `Post_ID` (INTEGER): The ID of the post (foreign key referencing Posts).
   - `Categories_ID` (INTEGER): The ID of the category (foreign key referencing Categories).

5. **Comments**
   - Stores comments made by users on posts.
   - `ID` (INTEGER): Unique auto-incremented identifier.
   - `User_UUID` (VARCHAR(255)): The UUID of the user who made the comment (foreign key referencing User).
   - `Post_ID` (INTEGER): The ID of the post being commented on (foreign key referencing Posts).
   - `Text` (TEXT): The content of the comment.
   - `Like` (INTEGER): Count of likes for the comment.
   - `Dislike` (INTEGER): Count of dislikes for the comment.
   - `CreatedAt` (DATETIME): When the comment was created.
   - `UpdatedAt` (DATETIME): When the comment was last updated.

6. **Report**
   - Stores reports made by users or moderators regarding posts.
   - `ID` (INTEGER): Unique auto-incremented identifier.
   - `User_UUID` (VARCHAR(255)): The UUID of the user who made the report (foreign key referencing User).
   - `Username` (VARCHAR(255)): The username of the user making the report.
   - `Post_ID` (INTEGER): The post being reported (foreign key referencing Posts).
   - `Title` (TEXT): The title of the report.
   - `Respons_Text` (TEXT): The admin's response to the report.

7. **Reaction**
   - Stores the reactions (like/dislike) made by users on posts and comments.
   - `ID` (INTEGER): Unique auto-incremented identifier.
   - `Post_ID` (INTEGER): The ID of the post being reacted to (foreign key referencing Posts, nullable).
   - `Comment_ID` (INTEGER): The ID of the comment being reacted to (foreign key referencing Comments, nullable).
   - `User_UUID` (VARCHAR(255)): The UUID of the user who reacted (foreign key referencing User).
   - `Status` (VARCHAR(255)): The type of reaction (e.g., "Like" or "Dislike").
   - **Constraints**: Ensures a reaction is either for a post or a comment, not both.

8. **Notification**
   - Stores notifications for users about reactions to their posts/comments.
   - `ID` (INTEGER): Unique auto-incremented identifier.
   - `User_UUID` (VARCHAR(255)): The UUID of the user receiving the notification (foreign key referencing User).
   - `Reaction_ID` (INTEGER): The ID of the reaction (foreign key referencing Reaction, nullable).
   - `Post_ID` (INTEGER): The ID of the post (foreign key referencing Posts, nullable).
   - `Comment_ID` (INTEGER): The ID of the comment (foreign key referencing Comments, nullable).
   - `CreatedAt` (DATETIME): When the notification was created.
   - `IsRead` (BOOL): Marks the notification as read or unread.

9. **Image**
   - Stores images uploaded by users for posts.
   - `ID` (INTEGER): Unique auto-incremented identifier.
   - `FilePath` (TEXT): The file path of the uploaded image.
   - `Post_ID` (INTEGER): The ID of the post associated with the image (foreign key referencing Posts).

10. **oauth_states**
    - Stores states for OAuth authentication flows (e.g., Google, GitHub).
    - `state` (TEXT): The state value for OAuth flow.
    - `created_at` (DATETIME): When the state was created.

---

## Contributing

Contributions are welcome! If you have suggestions or want to contribute to the project, please fork the repository, create a new branch, and submit a pull request.

If you find any bugs or have feature requests, please open an issue in the repository. Provide as much detail as possible to help us understand the problem or suggestion.

Thank you for your contributions!

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.