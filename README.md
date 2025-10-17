# AWS Documentor

A comprehensive AWS VPC infrastructure documentation and visualization tool written in Go. Scans your AWS VPC infrastructure and generates detailed JSON output and draw.io diagrams.

## Features

- **Comprehensive Scanning**: Retrieves detailed information about:
  - VPCs (Virtual Private Clouds)
  - Subnets
  - Route Tables
  - Security Groups
  - Internet Gateways
  - NAT Gateways
  - Transit Gateways
  - Transit Gateway Attachments

- **Visual Diagrams**: Generates draw.io compatible diagrams showing:
  - VPC containers with CIDR blocks
  - Public and private subnets
  - Internet Gateway placement
  - NAT Gateway locations
  - Transit Gateway connections
  - Route table information
  - Security group summaries

- **JSON Output**: Detailed JSON output for programmatic analysis and integration

## Installation

```bash
go build -o aws-documentor
```

## Prerequisites

- Go 1.21 or higher
- AWS credentials configured (via environment variables, AWS credentials file, or IAM role)
- IAM permissions for:
  - `ec2:DescribeVpcs`
  - `ec2:DescribeSubnets`
  - `ec2:DescribeRouteTables`
  - `ec2:DescribeSecurityGroups`
  - `ec2:DescribeInternetGateways`
  - `ec2:DescribeNatGateways`
  - `ec2:DescribeTransitGateways`
  - `ec2:DescribeTransitGatewayAttachments`

## Usage

### Basic scan (JSON output only)
```bash
./aws-documentor
```

### Scan specific region
```bash
./aws-documentor -region us-west-2
```

### Generate draw.io diagram
```bash
./aws-documentor -diagram
```

This creates a file named `vpc-diagram.drawio` that can be opened in [draw.io](https://app.diagrams.net).

### Scan without JSON output (diagram only)
```bash
./aws-documentor -diagram -json=false
```

### Scan specific region and generate diagram
```bash
./aws-documentor -region eu-central-1 -diagram
```

## Command-Line Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-region` | string | (from AWS config) | AWS region to scan |
| `-diagram` | bool | false | Generate draw.io diagram file |
| `-json` | bool | true | Output JSON data to stdout |

## Output

### JSON Output
When `-json=true` (default), the tool outputs detailed JSON for each resource type:
- Resource IDs and names
- CIDR blocks and IP addresses
- States and configurations
- Tags
- Associations and relationships

### Diagram Output
When `-diagram` flag is used, generates `vpc-diagram.drawio` containing:

**VPC Visualization**:
- VPC containers showing CIDR blocks
- Subnets labeled as Public/Private with CIDR and AZ information
- Internet Gateways attached to VPCs
- NAT Gateways positioned in their respective subnets

**Transit Gateway Section**:
- Transit Gateway resources with ASN information
- Attachment details showing resource types and states

**Information Panels**:
- Route tables with route destinations and targets
- Security group summaries with rule counts

## Architecture

```
aws-documentor/
├── main.go                    # Main application entry point
├── modules/
│   ├── vpc/
│   │   └── vpc.go            # VPC scanning and data structures
│   └── diagram/
│       └── diagram.go        # Draw.io diagram generation
├── go.mod                    # Go module definition
└── README.md                 # This file
```

## Draw.io Integration

The generated `.drawio` files use:
- **AWS Architecture Icons**: Official AWS shapes for gateways
- **Color Coding**:
  - Blue: VPC containers
  - Yellow: Subnets
  - Purple: AWS services (IGW, NAT, TGW)
  - Gray: Information panels
- **Hierarchical Layout**: VPCs as containers with nested resources

### Opening Diagrams

1. Go to [https://app.diagrams.net](https://app.diagrams.net)
2. Click "Open Existing Diagram"
3. Select your `vpc-diagram.drawio` file
4. Edit, export to PNG/PDF, or share as needed

## Use Cases

- **Security Auditing**: Review security group rules and network access
- **Documentation**: Generate up-to-date infrastructure diagrams
- **Compliance**: Create audit trails for compliance requirements
- **Cloud Visibility**: Understand complex VPC architectures
- **Incident Response**: Quickly visualize network topology during investigations
- **Architecture Review**: Assess network design and identify improvements

## AWS Authentication

The tool uses the AWS SDK default credential chain:

1. Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
2. Shared credentials file (`~/.aws/credentials`)
3. IAM role (if running on EC2)
4. ECS task role (if running in ECS)

## Examples

### Example 1: Multi-region documentation
```bash
for region in us-east-1 us-west-2 eu-west-1; do
  ./aws-documentor -region $region -diagram -json=false
  mv vpc-diagram.drawio vpc-diagram-$region.drawio
done
```

### Example 2: JSON analysis with jq
```bash
# Count subnets per VPC
./aws-documentor -json | jq -r '.vpc_id' | sort | uniq -c

# Find public subnets
./aws-documentor -json | jq 'select(.map_public_ip_on_launch == true)'
```

### Example 3: Automated documentation pipeline
```bash
#!/bin/bash
# Scan and upload to S3
./aws-documentor -diagram -json > vpc-data.json
aws s3 cp vpc-diagram.drawio s3://my-bucket/docs/
aws s3 cp vpc-data.json s3://my-bucket/data/
```

## Limitations

- Single region per execution (use `-region` flag for different regions)
- Read-only operations (no modifications to AWS infrastructure)
- Diagram layout is automatic (may need manual adjustment for complex topologies)

## Contributing

Contributions welcome! Areas for enhancement:
- Additional AWS service support (Load Balancers, Endpoints, etc.)
- Multi-region diagram generation
- Custom diagram layouts and themes
- Export to other formats (Terraform, CloudFormation)

## License

This is an open-source defensive security tool for AWS infrastructure documentation.

## Troubleshooting

### "Failed to load AWS config"
- Ensure AWS credentials are configured
- Check IAM permissions
- Verify region is valid

### "No resources found"
- Confirm you're scanning the correct region
- Verify VPC resources exist in the account
- Check IAM permissions include all required `ec2:Describe*` actions

### Diagram appears empty
- Ensure the region has VPC resources
- Check that the scan completed successfully
- Verify the `.drawio` file was created (should be > 1KB)

## Version

Current version: 1.0.0
