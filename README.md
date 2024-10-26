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



