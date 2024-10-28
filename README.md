# Architecture
![Architecture Image](https://github.com/user-attachments/assets/3cc5ffaf-363b-4ffb-a379-406e2a92bc8e)





# User profile
For the sake of completeness and complete view of the design, customer profile is supposed to be implemented with this schema

## User Table

| Column Name      | Data Type    | Constraints                | Description                        |
|------------------|--------------|----------------------------|------------------------------------|
| id               | TEXT         | PRIMARY KEY   | The unique identifier for the customer |
| name          | VARCHAR(255) | NULL                   | Name or description of the customer |

-------------------

# Schemas

## Post
| Column Name | Data Type    | Constraints                 | Description                             |
|-------------|--------------|-----------------------------|-----------------------------------------|
| id          | TEXT          | PRIMARY KEY    | The unique identifier for the post      |
| caption     | VARCHAR(255) | NOT NULL                    | Name or description of the customer     |
| image       | VARCHAR(255) | NULL                        | Name or description of the customer     |
| user_id     | INT          | Foreign Key to User Table   | The unique identifier for the customer  |
| created_at  | DATETIME          | Default CURRENT_TIMESTAMP   | The timestamp when the post was created |
| updated_at  | DATETIME          | Default CURRENT_TIMESTAMP ON  UPDATE CURRENT_TIMESTAMP | The timestamp when the post was updated |

## Comment
| Column Name | Data Type    | Constraints                                            | Description                             |
|-------------|--------------|--------------------------------------------------------|-----------------------------------------|
| id          | TEXT         | PRIMARY KEY                                            | The unique identifier for the comment   |
| comment     | VARCHAR(255) | NOT NULL                                               | Name or description of the comment      |
| post_id     | INT          | Foreign Key to Post Table                              | The unique identifier for the post      |
| user_id     | INT          | Foreign Key to User Table                              | The unique identifier for the customer  |
| created_at  | DATETIME     | Default CURRENT_TIMESTAMP                              | The timestamp when the comment was created |
| updated_at  | DATETIME     | Default CURRENT_TIMESTAMP ON  UPDATE CURRENT_TIMESTAMP | The timestamp when the comment was updated |


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
```json
{
    "id": 1,
    "presignedURL": <presigned-url>
}
```
* User immediately gets a presigned url to upload the image
* Client side handles the logic to upload the image to the presigned url
* Image is uploaded using multipart/form-data

User requests a post
GET /posts
```request
body
{
    "limit": <number of posts to retrieve>,
    "cursor": <timestamp of the last post>
}
```
```response
    "id": <id of the post in case of specific post retrieval>,
    "caption": <caption attached with the post>,
    "image_url": <jpg image of the URL>,
    "comments": [
        {
            "id": <id of the comment>,
            "comment": <comment attached with the post>
        }
    ],
    "created_at": <timestamp of the last post, used for cursor pagination>
}
```

## User comments on a post
POST /api/posts/{postId}/comments
```request
{
    "comment": "<comment>",
    "postId": "<id of the post on which to add the comment>"
}
```
```response
{
    "id": <id of the comment>,
}
```

# Handling image uploads
* In order to handle the high bandwidth requirement of the image uploads, I decided to use s3 presigned URLs.
* In order to control the allowed file types to upload, we need to setup the bucket policy to now allow upload more than 100 MB
* Still customer can provide corrupt file or file with wrong extension.
* We can implement a logic in client side to pass on some part of the file being uploaded
* Server side will hash the input (MD5 hash) and set that in the metadata of the presigned URL
* Worker group will recalculate the hash before processing the file and compare it with the hash in the metadata
* If they don't match delete the file and return an error
* We can also do signature verification on the file to make sure it's not tampered with
![Upload flow](https://github.com/user-attachments/assets/c024c586-550e-4ce8-b049-210048b4ec7c)

  

# production deployment
* we need to add user authentication and authorization in existing workflow
* Lambda function is required to be limited in terms of capacity so that in case the file is corrupt and has worm embedded, it doesn't affect the system
* Posts and comments are required to be horizontally scalable
* Add Unit test and integration test to the code
* We can use ECS service to deploy our services in containers and add autoscaling rules to scale the services


