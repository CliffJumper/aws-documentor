# AWS Documentor - Usage Examples

## Quick Start Examples

### Example 1: Scan and View JSON Output
```bash
./aws-documentor -region us-east-1
```

**Output:**
```
Scanning AWS region: us-east-1

Scanning VPCs...
Found 2 VPCs:
{
  "vpc_id": "vpc-0abc123def456",
  "cidr_block": "10.0.0.0/16",
  "state": "available",
  "is_default": false,
  "tags": {
    "Name": "Production VPC",
    "Environment": "prod"
  }
}
---
...
```

### Example 2: Generate Infrastructure Diagram
```bash
./aws-documentor -region us-west-2 -diagram -json=false
```

**Output:**
```
Scanning AWS region: us-west-2

Scanning VPCs...
Found 1 VPCs

Scanning Subnets...
Found 6 Subnets

Scanning Route Tables...
Found 3 Route Tables

Scanning Security Groups...
Found 8 Security Groups

Scanning Internet Gateways...
Found 1 Internet Gateways

Scanning NAT Gateways...
Found 2 NAT Gateways

Scanning Transit Gateways...
Found 1 Transit Gateways

Scanning Transit Gateway Attachments...
Found 3 Transit Gateway Attachments

VPC infrastructure scan complete!

Generating draw.io diagram...
Diagram saved to: vpc-diagram.drawio
You can open this file in draw.io (https://app.diagrams.net)
```

### Example 3: Multi-Region Documentation
```bash
#!/bin/bash
# Document all regions

REGIONS="us-east-1 us-west-2 eu-west-1 ap-southeast-1"

for region in $REGIONS; do
  echo "Documenting $region..."

  # Generate diagram
  ./aws-documentor -region $region -diagram -json=false
  mv vpc-diagram.drawio "diagrams/vpc-${region}.drawio"

  # Save JSON data
  ./aws-documentor -region $region -diagram=false > "data/vpc-${region}.json"

  echo "Completed $region"
  echo "---"
done

echo "All regions documented!"
```

### Example 4: Security Audit Workflow
```bash
# Scan infrastructure
./aws-documentor > vpc-audit.json

# Find all public subnets
cat vpc-audit.json | jq 'select(.map_public_ip_on_launch == true)' > public-subnets.json

# Find security groups allowing 0.0.0.0/0
cat vpc-audit.json | jq '.rules[]? | select(.cidr_block == "0.0.0.0/0")' > open-security-groups.json

# Generate visual diagram
./aws-documentor -diagram -json=false

# Review diagram in draw.io for network topology
```

### Example 5: Automated Daily Documentation
```bash
#!/bin/bash
# Add to cron: 0 2 * * * /path/to/daily-vpc-doc.sh

DATE=$(date +%Y-%m-%d)
OUTPUT_DIR="/var/vpc-docs/$DATE"

mkdir -p $OUTPUT_DIR

# Scan and save
./aws-documentor -diagram > "$OUTPUT_DIR/vpc-data.json"
mv vpc-diagram.drawio "$OUTPUT_DIR/vpc-diagram.drawio"

# Upload to S3
aws s3 cp "$OUTPUT_DIR/" "s3://my-bucket/vpc-docs/$DATE/" --recursive

# Keep last 30 days
find /var/vpc-docs -type d -mtime +30 -exec rm -rf {} \;

echo "VPC documentation complete for $DATE"
```

## Use Case Scenarios

### Scenario 1: New Team Member Onboarding
```bash
# Generate comprehensive documentation
./aws-documentor -region us-east-1 -diagram

# Share the diagram and JSON with new team member
# They can open vpc-diagram.drawio in draw.io to visualize the infrastructure
```

**Benefits:**
- Visual understanding of VPC architecture
- Detailed JSON for specific resource lookups
- No AWS console access needed initially

### Scenario 2: Security Compliance Audit
```bash
# Generate infrastructure snapshot
./aws-documentor -region us-east-1 > compliance-report-$(date +%Y%m%d).json

# Create visual diagram for auditor
./aws-documentor -region us-east-1 -diagram -json=false

# Extract security groups for review
cat compliance-report-*.json | jq '[.group_id, .group_name, .rules[]]' > security-groups.json
```

**Audit deliverables:**
- Complete infrastructure inventory (JSON)
- Visual network diagram (draw.io)
- Security group analysis

### Scenario 3: Disaster Recovery Documentation
```bash
# Document all critical regions
for region in us-east-1 us-west-2; do
  ./aws-documentor -region $region -diagram -json > dr-${region}.json
  mv vpc-diagram.drawio dr-${region}.drawio
done

# Store in secure location
tar czf dr-documentation-$(date +%Y%m%d).tar.gz dr-*.json dr-*.drawio
```

### Scenario 4: Architecture Review
```bash
# Generate current state diagram
./aws-documentor -diagram -json=false

# Open in draw.io
# Review:
# - Subnet distribution across AZs
# - NAT Gateway placement
# - Security group configurations
# - Transit Gateway connections
```

### Scenario 5: Cost Optimization Analysis
```bash
# Identify all NAT Gateways
./aws-documentor | jq 'select(.nat_gateway_id != null) | {id: .nat_gateway_id, subnet: .subnet_id, state: .state}'

# Check Transit Gateway usage
./aws-documentor | jq 'select(.transit_gateway_id != null)'

# Review diagram for consolidation opportunities
./aws-documentor -diagram -json=false
```

## Advanced Filtering with jq

### Find all VPCs
```bash
./aws-documentor | jq 'select(.vpc_id != null) | {vpc_id, cidr_block, name: .tags.Name}'
```

### List all subnets with AZ distribution
```bash
./aws-documentor | jq 'select(.subnet_id != null) | {subnet: .subnet_id, cidr: .cidr_block, az: .availability_zone, public: .map_public_ip_on_launch}'
```

### Find security groups with SSH access from anywhere
```bash
./aws-documentor | jq 'select(.rules[]? | select(.from_port == 22 and .cidr_block == "0.0.0.0/0"))'
```

### Count resources by type
```bash
# VPCs
./aws-documentor | jq 'select(.vpc_id != null)' | jq -s 'length'

# Subnets
./aws-documentor | jq 'select(.subnet_id != null)' | jq -s 'length'

# NAT Gateways
./aws-documentor | jq 'select(.nat_gateway_id != null)' | jq -s 'length'
```

### List all route tables and their routes
```bash
./aws-documentor | jq 'select(.route_table_id != null) | {
  id: .route_table_id,
  vpc: .vpc_id,
  main: .is_main_route_table,
  routes: [.routes[] | {dest: .destination_cidr_block, target: (.gateway_id // .nat_gateway_id // .transit_gateway_id)}]
}'
```

## Integration Examples

### Jenkins Pipeline
```groovy
pipeline {
    agent any
    stages {
        stage('Document Infrastructure') {
            steps {
                sh './aws-documentor -region us-east-1 -diagram'
                archiveArtifacts artifacts: 'vpc-diagram.drawio', fingerprint: true
                sh './aws-documentor -region us-east-1 > vpc-data.json'
                archiveArtifacts artifacts: 'vpc-data.json', fingerprint: true
            }
        }
    }
}
```

### GitHub Actions
```yaml
name: VPC Documentation

on:
  schedule:
    - cron: '0 0 * * 0'  # Weekly on Sunday
  workflow_dispatch:

jobs:
  document:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2

      - name: Setup AWS Credentials
        uses: aws-actions/configure-aws-credentials@v1
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-east-1

      - name: Build AWS Documentor
        run: go build -o aws-documentor

      - name: Generate Documentation
        run: |
          ./aws-documentor -diagram > vpc-data.json

      - name: Upload Artifacts
        uses: actions/upload-artifact@v2
        with:
          name: vpc-documentation
          path: |
            vpc-diagram.drawio
            vpc-data.json
```

### Terraform Integration
```bash
# After Terraform apply, document the changes
terraform apply -auto-approve

# Wait for resources to be fully created
sleep 30

# Generate updated documentation
./aws-documentor -region us-east-1 -diagram

# Commit documentation to repo
git add vpc-diagram.drawio
git commit -m "Update VPC documentation after Terraform apply"
git push
```

## Tips and Tricks

### 1. Comparing Infrastructure Over Time
```bash
# Save baseline
./aws-documentor > baseline.json

# After changes, compare
./aws-documentor > current.json
diff baseline.json current.json
```

### 2. Creating Regional Diagrams
```bash
# Create separate diagrams per region
for region in us-east-1 us-west-2; do
  ./aws-documentor -region $region -diagram -json=false
  mv vpc-diagram.drawio vpc-${region}.drawio
done
```

### 3. Filtering Output by VPC
```bash
# Get data for specific VPC
VPC_ID="vpc-0abc123def456"
./aws-documentor | jq "select(.vpc_id == \"$VPC_ID\")"
```

### 4. Quick Resource Count
```bash
./aws-documentor -diagram=false -json=false
```
Shows just the counts of each resource type.

### 5. Automated Alerting
```bash
# Check for overly permissive security groups
OPEN_SG=$(./aws-documentor | jq -r 'select(.rules[]? | select(.cidr_block == "0.0.0.0/0" and .from_port == 22))' | wc -l)

if [ $OPEN_SG -gt 0 ]; then
  echo "WARNING: Found $OPEN_SG security groups with SSH open to 0.0.0.0/0"
  # Send alert
fi
```

## Troubleshooting Common Issues

### Issue: Empty diagram generated
**Solution:**
```bash
# Check if resources exist
./aws-documentor -json=false

# Verify AWS credentials
aws sts get-caller-identity

# Check specific region
./aws-documentor -region us-east-1 -json=false
```

### Issue: Diagram too large/complex
**Solution:**
```bash
# Generate separate diagrams per VPC
# (Future enhancement - currently generates all VPCs in one diagram)

# Open in draw.io and use layers to organize
```

### Issue: Missing resources in scan
**Solution:**
```bash
# Verify IAM permissions
aws iam get-user-policy --user-name YOUR_USER --policy-name YOUR_POLICY

# Test specific API calls
aws ec2 describe-vpcs --region us-east-1
aws ec2 describe-subnets --region us-east-1
```
