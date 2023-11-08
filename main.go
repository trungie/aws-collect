package main

import (
    "fmt"
    "os"
    "os/exec"
    "strings"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sts"
)

func main() {
    // Assuming command line arguments are passed in this order:
    // 1. AWS Command (e.g., "rds describe instances")
    // 2. Accounts (comma-separated)
    // 3. Roles (comma-separated)
    // 4. Regions (comma-separated)
    // 5. MFA Serial
    if len(os.Args) < 6 {
        fmt.Println("Usage: aws-collect <aws-command> <accounts> <roles> <regions> <mfa-serial>")
        os.Exit(1)
    }

    awsCommand := os.Args[1]
    accounts := os.Args[2]
    roles := os.Args[3]
    regions := os.Args[4]
    mfaSerial := os.Args[5]

    // Initialize AWS session
    sess := session.Must(session.NewSession())

    for _, account := range strings.Split(accounts, ",") {
        for _, role := range strings.Split(roles, ",") {
            // Prompt for MFA token
            fmt.Printf("Enter MFA token for account %s and role %s: ", account, role)
            var mfaToken string
            fmt.Scanln(&mfaToken)

            // Assume role with MFA
            svc := sts.New(sess)
            params := &sts.AssumeRoleInput{
                RoleArn:         aws.String(fmt.Sprintf("arn:aws:iam::%s:role/%s", account, role)),
                RoleSessionName: aws.String("aws-collect-session"),
                SerialNumber:    aws.String(mfaSerial),
                TokenCode:       aws.String(mfaToken),
            }
            assumedRole, err := svc.AssumeRole(params)
            if err != nil {
                fmt.Printf("Error assuming role: %s\n", err)
                continue
            }

            creds := assumedRole.Credentials
            for _, region := range strings.Split(regions, ",") {
                // Set region and execute AWS command
                os.Setenv("AWS_ACCESS_KEY_ID", *creds.AccessKeyId)
                os.Setenv("AWS_SECRET_ACCESS_KEY", *creds.SecretAccessKey)
                os.Setenv("AWS_SESSION_TOKEN", *creds.SessionToken)
                os.Setenv("AWS_DEFAULT_REGION", region)

                cmd := exec.Command("aws", strings.Split(awsCommand, " ")...)
                cmd.Stdout = os.Stdout
                cmd.Stderr = os.Stderr

                err := cmd.Run()
                if err != nil {
                    fmt.Printf("Error executing command '%s' in region %s: %s\n", awsCommand, region, err)
                }
            }
        }
    }
}

