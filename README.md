# Architecture
![Architecture Image](https://github.com/user-attachments/assets/8deb041e-f08b-452e-b2d9-61fc4f5f71db)


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
    "comment": "This is a comment",
}
```


