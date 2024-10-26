package app

type S3Config struct {
	Region                 string
	Bucket                 string
	Key                    string
	PreSignedURLExpiration int // in minutes
}

func DefaultS3Config() *S3Config {
	return &S3Config{
		Region:                 "your-region",
		Bucket:                 "your-bucket",
		Key:                    "your-key",
		PreSignedURLExpiration: 15,
	}
}
