package repository

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/skyrenx/blog-api-go/http/entities"
	"github.com/skyrenx/blog-api-go/http/entities/dto"
	"golang.org/x/crypto/bcrypt"
)

const (
	REGION                = "us-east-1"
	TOKEN_EXPIRATION_TIME = 15 //Minutes
)

func GetUserByUsername(username string) (*dto.UserWithoutPassword, error) {
	ctx := context.Background()
	conn, err := getConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	query := `SELECT username, enabled FROM users WHERE username = $1 `
	rows, err := conn.Query(ctx, query, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get row by username: %v: %w", username, err)
	}
	defer rows.Close()
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dto.UserWithoutPassword])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("no blog entry found with username %v", username)
		}
		return nil, fmt.Errorf("failed to collect row: %w", err)
	}
	return &user, nil
}

func RegisterUser(user entities.User) error {
	ctx := context.Background()
	conn, err := getConnection(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	// Insert user into the database
	query := `INSERT INTO users (username, password, enabled) VALUES ($1, $2, $3)`
	_, err = conn.Exec(ctx, query, user.Username, hashedPassword, true)
	if err != nil {
		return err
	}
	return nil
}

func Login(userCredentials entities.User) (*string, error) {
	ctx := context.Background()
	conn, err := getConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	query := `SELECT username, password, enabled FROM users WHERE username = $1 `
	rows, err := conn.Query(ctx, query, userCredentials.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to get row by username: %v: %w", userCredentials.Username, err)
	}
	defer rows.Close()
	foundUser, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.User])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("no blog entry found with username %v", userCredentials.Username)
		}
		return nil, fmt.Errorf("failed to collect row: %w", err)
	}

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(userCredentials.Password)); err != nil {
		return nil, fmt.Errorf("invalid username or password: %v: %w", userCredentials.Username, err)
	}
	// Generate JWT token
	token, err := generateJWT(userCredentials.Username)
	if err != nil {
		return nil, fmt.Errorf("could not generate token: %v: %w", userCredentials.Username, err)
	}
	return &token, nil
}

// hashPassword hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func generateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours

	claims := &entities.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := os.Getenv("JWT_SECRET")

	// Sign the token with the secret key
	return token.SignedString([]byte(jwtSecret))
}

// Connect to the Aurora DSQL cluster.
// https://docs.aws.amazon.com/aurora-dsql/latest/userguide/SECTION_program-with-go.html
// Returns a closer function to close the connection.
func getConnection(ctx context.Context) (*pgx.Conn, error) {
	clusterEndpoint := os.Getenv("CLUSTER_ENDPOINT")

	// Build connection URL
	var sb strings.Builder
	sb.WriteString("postgres://")
	sb.WriteString(clusterEndpoint)
	sb.WriteString(":5432/postgres?user=admin&sslmode=verify-full")
	url := sb.String()

	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	creds, err := sess.Config.Credentials.Get()
	if err != nil {
		return nil, err
	}
	staticCredentials := credentials.NewStaticCredentials(
		creds.AccessKeyID,
		creds.SecretAccessKey,
		creds.SessionToken,
	)

	// The token expiration time is optional, and the default value 900 seconds (15 minutes)
	// If you are not connecting as admin, use DbConnect action instead
	token, err := generateDbConnectAdminAuthToken(staticCredentials, clusterEndpoint)
	if err != nil {
		return nil, err
	}

	connConfig, err := pgx.ParseConfig(url)
	// To avoid issues with parse config set the password directly in config
	connConfig.Password = token
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse config: %v\n", err)
		return nil, err
	}

	conn, err := pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect config: %v\n", err)
		return nil, err
	}
	return conn, nil
}

// generate password token to connect to your Aurora DSQL cluster.
func generateDbConnectAdminAuthToken(creds *credentials.Credentials, clusterEndpoint string) (string, error) {
	// the scheme is arbitrary and is only needed because validation of the URL requires one.
	endpoint := "https://" + clusterEndpoint
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}
	values := req.URL.Query()
	values.Set("Action", "DbConnectAdmin")
	req.URL.RawQuery = values.Encode()

	signer := v4.Signer{
		Credentials: creds,
	}
	_, err = signer.Presign(req, nil, "dsql", REGION, TOKEN_EXPIRATION_TIME*time.Minute, time.Now())
	if err != nil {
		return "", err
	}

	url := req.URL.String()[len("https://"):]

	return url, nil
}
