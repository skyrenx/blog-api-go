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