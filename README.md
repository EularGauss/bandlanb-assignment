# Architecture
![Architecture Image](https://github.com/user-attachments/assets/6a4d0595-e27f-4a16-b3f5-65a5fe797258)





# User profile
For the sake of completeness and complete view of the design, customer profile is supposed to be implemented with this schema

## User Table

| Column Name      | Data Type   | Constraints                  | Description                        |
|------------------|-------------|------------------------------|------------------------------------|
| id               | INT         | PRIMARY KEY, AUTOINCREMENT   | The unique identifier for the customer |
| name          | VARCHAR(255)        | NULL                     | Name or description of the customer |

-------------------

# Schemas

## Post
| Column Name | Data Type    | Constraints                   | Description                             |
|-------------|--------------|-------------------------------|-----------------------------------------|
| id          | INT          | PRIMARY KEY, AUTOINCREMENT    | The unique identifier for the post      |
| caption     | VARCHAR(255) | NOT NULL                      | Name or description of the customer     |
| image       | VARCHAR(255) | NULL                          | Name or description of the customer     |
| user_id     | INT          | Foreign Key to User Table     | The unique identifier for the customer  |
| created_at  | DATETIME          | Default CURRENT_TIMESTAMP     | The timestamp when the post was created |
| updated_at  | DATETIME          | Default CURRENT_TIMESTAMP ON  UPDATE CURRENT_TIMESTAMP | The timestamp when the post was updated |

## Comment
| Column Name | Data Type    | Constraints                   | Description                             |
|-------------|--------------|-------------------------------|-----------------------------------------|
| id          | INT          | PRIMARY KEY, AUTOINCREMENT    | The unique identifier for the comment   |
| comment     | VARCHAR(255) | NOT NULL                      | Name or description of the comment      |
| post_id     | INT          | Foreign Key to Post Table     | The unique identifier for the post      |
| user_id     | INT          | Foreign Key to User Table     | The unique identifier for the customer  |
| created_at  | DATETIME          | Default CURRENT_TIMESTAMP     | The timestamp when the comment was created |
| updated_at  | DATETIME          | Default CURRENT_TIMESTAMP ON  UPDATE CURRENT_TIMESTAMP | The timestamp when the comment was updated |


# User interactions
## User create a post
POST /posts
```json
{
    "caption": "This is a caption",
    "image": "true"
}
```
Headers
```headers
Authorization
JWT <token> (includes user-id)
```
```response
{
    "id": 1,
    "caption": "This is a caption",
    "image": <presigned-url>
}
```
* User immediately gets a presigned url to upload the image
* Client side handles the logic to upload the image to the presigned url
* Image is uploaded using multipart/form-data

User requests a post
GET /posts
```request
{
    "id": <id of the post in case of specific post retrieval>
    "number_of_posts": <number of posts to fetch>
    "cursor": <id of the last post fetched>
}
```
```response
[{
    "id": 1,
    "caption": "This is a caption",
    "image": <presigned-url>
}]
```

## User comments on a post
POST /api/posts/{postId}/comments
```json
{
    "comment": "This is a comment"
}
```


# Handling uploads
* In order to handle the high bandwidth requirement of the image uploads, I decided to use s3 presigned URLs.
* In order to control the allowed file types to upload, we need to setup the bucket policy to now allow upload more than 100 MB
* Still customer can provide corrupt file or file with wrong extension.
* We can implement a logic in client side to pass on some part of the file being uploaded
* Server side will hash the input (MD5 hash) and set that in the metadata of the presigned URL
* Worker group will recalculate the hash before processing the file and compare it with the hash in the metadata
* If they don't match delete the file and return an error
* We can also do signature verification on the file to make sure it's not tampered with

# production deployment
* we need to add user authentication and authorization in existing workflow
* Lambda function is required to be limited in terms of capacity so that in case the file is corrupt and has worm embedded, it doesn't affect the system
* Posts and comments are required to be horizontally scalable

