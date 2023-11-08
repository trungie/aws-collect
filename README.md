AWS collect, runs aws commands on multiple accounts and regions and collects the information, and outputs it via text, table or json

aws-vault exec security -- aws-collect rds describe-instances

Should run `aws rds describe-instances` on each of the accounts and regions listed in the ~/.aws-collect.yml 

- assume the role into the account first
- run the command and collect results
