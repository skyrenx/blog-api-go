package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/skyrenx/blog-api-go/http/entities"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
)

const (
	REGION = "us-east-1"
)

func GetBlogEntries(page int) ([]entities.BlogEntry, int, error) {
	pageSize := 2

	if page < 1 {
		return nil, 0, fmt.Errorf("failed to get blog entries. requested page should be greater than 0")
	}

	// Get the cluster endpoint from the environment
	clusterEndpoint := os.Getenv("CLUSTER_ENDPOINT")
	_, b := os.LookupEnv("CLUSTER_ENDPOINT")
	fmt.Printf("Cluster endpoint found? %v \nUsing cluster endpoint: %s\n", b, clusterEndpoint)

	ctx := context.Background()

	// Establish connection
	conn, err := getConnection(ctx, clusterEndpoint)
	if err != nil {
		return nil, 0, err
	}
	defer conn.Close(ctx)

	var totalRows int
	query := `SELECT COUNT(*) FROM blog_entries`
	err = conn.QueryRow(ctx, query).Scan(&totalRows)
	if err != nil {
		return nil, 0, err
	}
	totalPages := (totalRows + pageSize - 1) / pageSize
	if page > totalPages {
		return nil, 0, fmt.Errorf(
			"requested page does not exist. Page requested was %v, total pages is %v",
			page, pageSize)
	}

	// Calculate the offset
	offset := (page - 1) * pageSize

	query = `SELECT * FROM blog_entries ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := conn.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	blogEntries, _ := pgx.CollectRows(rows, pgx.RowToStructByName[entities.BlogEntry])
	fmt.Printf("blogEntries: %v\n", blogEntries)
	return blogEntries, totalPages, nil
}

func GetBlogEntryById(id int) (*entities.BlogEntry, error){
		// Get the cluster endpoint from the environment
		clusterEndpoint := os.Getenv("CLUSTER_ENDPOINT")
		_, b := os.LookupEnv("CLUSTER_ENDPOINT")
		fmt.Printf("Cluster endpoint found? %v \nUsing cluster endpoint: %s\n", b, clusterEndpoint)
	
		ctx := context.Background()
	
		// Establish connection
		conn, err := getConnection(ctx, clusterEndpoint)
		if err != nil {
			return nil, err
		}
		defer conn.Close(ctx)

		query := `SELECT * FROM blog_entries WHERE id = $1 `
		rows, err := conn.Query(ctx, query, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get row by id: %v: %w", id, err)
		}
		defer rows.Close()
		blogEntry, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.BlogEntry])
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, fmt.Errorf("no blog entry found with id %v", id)
			}
			return nil, fmt.Errorf("failed to collect row: %w", err)
		}
		return &blogEntry, nil

}

func CreateBlogEntry(entry entities.BlogEntry) error {
	clusterEndpoint := os.Getenv("CLUSTER_ENDPOINT")
	if clusterEndpoint == "" {
		return fmt.Errorf("CLUSTER_ENDPOINT is not set")
	}

	ctx := context.Background()

	// Establish connection
	conn, err := getConnection(ctx, clusterEndpoint)
	if err != nil {
		return fmt.Errorf("failed to establish connection: %w", err)
	}
	defer conn.Close(ctx)

	// Start a transaction
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx) // Rollback on error

	// Step 1: Retrieve the current NextId value
	// Aurora Serverless v2 does not allow unqualified FOR UPDATE on tables without a strict equality predicate on the key.
	var nextId int
	err = tx.QueryRow(ctx, `
    UPDATE blog_entry_sequence 
    SET next_id = next_id + 1 
    RETURNING next_id - 1
`).Scan(&nextId)
	if err != nil {
		return fmt.Errorf("failed to get next_id: %w", err)
	}

	// Step 2: Insert the new BlogEntry using the retrieved NextId
	query := `
		INSERT INTO blog_entries (id, title, content, author, created_at, updated_at, published)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = tx.Exec(ctx, query, nextId, entry.Title, entry.Content, entry.Author, time.Now(), time.Now(), entry.Published)
	if err != nil {
		return fmt.Errorf("failed to insert blog entry: %w", err)
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Printf("Blog entry created with ID: %d\n", nextId)
	return nil
}

func Example() error {
	// Get the cluster endpoint from the environment
	clusterEndpoint := os.Getenv("CLUSTER_ENDPOINT")
	_, b := os.LookupEnv("CLUSTER_ENDPOINT")
	fmt.Printf("Cluster endpoint found? %v \nUsing cluster endpoint: %s\n", b, clusterEndpoint)

	ctx := context.Background()

	// Establish connection
	conn, err := getConnection(ctx, clusterEndpoint)
	if err != nil {
		return err
	}

	query := `SELECT 1`
	_, err = conn.Exec(ctx, query)
	if err != nil {
		return err
	}

	var totalRows int
	query = `SELECT COUNT(*) FROM blog_entries`
	err = conn.QueryRow(ctx, query).Scan(&totalRows)
	if err != nil {
		panic(fmt.Sprintf("Failed to count rows: %v", err))
	}

	fmt.Printf("Total number of rows in blog_entries: %d\n", totalRows)

	//blogEntries := []entities.BlogEntry{}
	// Define the SQL query to insert a new owner record.
	query = `SELECT * FROM blog_entries LIMIT 10`
	rows, err := conn.Query(ctx, query)
	if err != nil {
		panic(fmt.Sprintf("error retrieving data, %v", err.Error()))
	}
	defer rows.Close()

	blogEntries, _ := pgx.CollectRows(rows, pgx.RowToStructByName[entities.BlogEntry])
	fmt.Printf("blogEntries: %v\n", blogEntries)
	// if err != nil || owners[0].Name != "John Doe" || owners[0].City != "Anytown" {
	// 	panic("Error retrieving data")
	// }

	defer conn.Close(ctx)

	return nil
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
	_, err = signer.Presign(req, nil, "dsql", REGION, 15*time.Minute, time.Now())
	if err != nil {
		return "", err
	}

	url := req.URL.String()[len("https://"):]

	return url, nil
}

// Connect to the Aurora DSQL cluster.
// https://docs.aws.amazon.com/aurora-dsql/latest/userguide/SECTION_program-with-go.html
func getConnection(ctx context.Context, clusterEndpoint string) (*pgx.Conn, error) {
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

	// The token expiration time is optional, and the default value 900 seconds
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
		os.Exit(1)
	}

	conn, err := pgx.ConnectConfig(ctx, connConfig)

	return conn, err
}
