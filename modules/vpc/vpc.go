// Package vpc provides functionality for scanning and retrieving AWS VPC and subnet information
package vpc

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// VPCInfo contains comprehensive information about an AWS VPC
type VPCInfo struct {
	VpcID               string            `json:"vpc_id"`                 // Unique identifier for the VPC
	CidrBlock           string            `json:"cidr_block"`             // Primary CIDR block assigned to the VPC
	State               string            `json:"state"`                  // Current state of the VPC (available, pending)
	IsDefault           bool              `json:"is_default"`             // Whether this is the default VPC for the region
	DhcpOptionsID       string            `json:"dhcp_options_id"`        // ID of the DHCP options set associated with the VPC
	InstanceTenancy     string            `json:"instance_tenancy"`       // Tenancy of instances launched into the VPC (default, dedicated, host)
	Tags                map[string]string `json:"tags"`                   // Key-value tags associated with the VPC
	AssociateCidrBlocks []string         `json:"associate_cidr_blocks"`  // Additional CIDR blocks associated with the VPC
}

// SubnetInfo contains comprehensive information about an AWS subnet
type SubnetInfo struct {
	SubnetID                    string            `json:"subnet_id"`                      // Unique identifier for the subnet
	VpcID                       string            `json:"vpc_id"`                         // ID of the VPC that contains this subnet
	CidrBlock                   string            `json:"cidr_block"`                     // CIDR block assigned to the subnet
	AvailabilityZone            string            `json:"availability_zone"`              // Availability zone where the subnet is located
	AvailabilityZoneID          string            `json:"availability_zone_id"`           // Unique ID of the availability zone
	State                       string            `json:"state"`                          // Current state of the subnet (available, pending)
	MapPublicIpOnLaunch         bool              `json:"map_public_ip_on_launch"`        // Whether instances launched in this subnet receive a public IP
	AssignIpv6AddressOnCreation bool              `json:"assign_ipv6_address_on_creation"` // Whether instances receive an IPv6 address on creation
	DefaultForAz                bool              `json:"default_for_az"`                 // Whether this is the default subnet for the availability zone
	Tags                        map[string]string `json:"tags"`                           // Key-value tags associated with the subnet
}

// RouteInfo contains information about an individual route in a route table
type RouteInfo struct {
	DestinationCidrBlock   string `json:"destination_cidr_block"`    // CIDR block for the route destination
	DestinationIpv6Block   string `json:"destination_ipv6_block"`    // IPv6 CIDR block for the route destination
	GatewayID              string `json:"gateway_id"`                // ID of the internet gateway or VPC gateway
	InstanceID             string `json:"instance_id"`               // ID of a NAT instance
	NatGatewayID           string `json:"nat_gateway_id"`            // ID of a NAT gateway
	NetworkInterfaceID     string `json:"network_interface_id"`      // ID of the network interface
	TransitGatewayID       string `json:"transit_gateway_id"`        // ID of the transit gateway
	VpcPeeringConnectionID string `json:"vpc_peering_connection_id"` // ID of the VPC peering connection
	State                  string `json:"state"`                     // State of the route (active, blackhole)
	Origin                 string `json:"origin"`                    // How the route was created (CreateRouteTable, CreateRoute, EnableVgwRoutePropagation)
}

// RouteTableInfo contains comprehensive information about an AWS route table
type RouteTableInfo struct {
	RouteTableID     string              `json:"route_table_id"`     // Unique identifier for the route table
	VpcID            string              `json:"vpc_id"`             // ID of the VPC that contains this route table
	Routes           []RouteInfo         `json:"routes"`             // List of routes in the route table
	SubnetIDs        []string            `json:"subnet_ids"`         // IDs of subnets explicitly associated with this route table
	IsMainRouteTable bool                `json:"is_main_route_table"` // Whether this is the main route table for the VPC
	Tags             map[string]string   `json:"tags"`               // Key-value tags associated with the route table
}

// SecurityGroupRule contains information about a security group rule
type SecurityGroupRule struct {
	IsEgress       bool   `json:"is_egress"`        // Whether this is an egress rule (true) or ingress rule (false)
	IpProtocol     string `json:"ip_protocol"`      // IP protocol (tcp, udp, icmp, or protocol number)
	FromPort       int32  `json:"from_port"`        // Start of port range (or ICMP type)
	ToPort         int32  `json:"to_port"`          // End of port range (or ICMP code)
	CidrBlock      string `json:"cidr_block"`       // CIDR block for the rule
	Ipv6CidrBlock  string `json:"ipv6_cidr_block"`  // IPv6 CIDR block for the rule
	GroupID        string `json:"group_id"`         // ID of referenced security group
	GroupOwnerID   string `json:"group_owner_id"`   // AWS account ID that owns the referenced security group
	PrefixListID   string `json:"prefix_list_id"`   // ID of the prefix list
	Description    string `json:"description"`      // Description of the rule
}

// SecurityGroupInfo contains comprehensive information about an AWS security group
type SecurityGroupInfo struct {
	GroupID     string              `json:"group_id"`     // Unique identifier for the security group
	GroupName   string              `json:"group_name"`   // Name of the security group
	Description string              `json:"description"`  // Description of the security group
	VpcID       string              `json:"vpc_id"`       // ID of the VPC that contains this security group
	OwnerID     string              `json:"owner_id"`     // AWS account ID that owns the security group
	Rules       []SecurityGroupRule `json:"rules"`        // List of all rules (ingress and egress) in the security group
	Tags        map[string]string   `json:"tags"`         // Key-value tags associated with the security group
}

// InternetGatewayInfo contains information about an AWS internet gateway
type InternetGatewayInfo struct {
	InternetGatewayID string            `json:"internet_gateway_id"` // Unique identifier for the internet gateway
	State             string            `json:"state"`               // State of the internet gateway (available, attached, detached, etc.)
	VpcID             string            `json:"vpc_id"`              // ID of the VPC this gateway is attached to (empty if detached)
	Tags              map[string]string `json:"tags"`                // Key-value tags associated with the internet gateway
}

// NatGatewayInfo contains information about an AWS NAT gateway
type NatGatewayInfo struct {
	NatGatewayID         string            `json:"nat_gateway_id"`          // Unique identifier for the NAT gateway
	SubnetID             string            `json:"subnet_id"`               // ID of the subnet the NAT gateway is in
	VpcID                string            `json:"vpc_id"`                  // ID of the VPC that contains this NAT gateway
	State                string            `json:"state"`                   // State of the NAT gateway (pending, failed, available, deleting, deleted)
	ConnectivityType     string            `json:"connectivity_type"`       // Connectivity type (public, private)
	PrivateIp            string            `json:"private_ip"`              // Private IP address of the NAT gateway
	PublicIp             string            `json:"public_ip"`               // Public IP address of the NAT gateway (if applicable)
	AllocationID         string            `json:"allocation_id"`           // ID of the Elastic IP address allocation
	NetworkInterfaceID   string            `json:"network_interface_id"`    // ID of the network interface for the NAT gateway
	CreatedTime          string            `json:"created_time"`            // Time when the NAT gateway was created
	Tags                 map[string]string `json:"tags"`                    // Key-value tags associated with the NAT gateway
}

// TransitGatewayInfo contains information about an AWS Transit Gateway
type TransitGatewayInfo struct {
	TransitGatewayID     string            `json:"transit_gateway_id"`      // Unique identifier for the transit gateway
	State                string            `json:"state"`                   // State of the transit gateway (pending, available, modifying, deleting, deleted)
	OwnerID              string            `json:"owner_id"`                // AWS account ID that owns the transit gateway
	Description          string            `json:"description"`             // Description of the transit gateway
	CreationTime         string            `json:"creation_time"`           // Time when the transit gateway was created
	DefaultRouteTableID  string            `json:"default_route_table_id"`  // ID of the default route table
	PropagationRouteTableID string         `json:"propagation_route_table_id"` // ID of the default propagation route table
	AmazonSideAsn        int64             `json:"amazon_side_asn"`         // Private Autonomous System Number (ASN) for the Amazon side of the BGP session
	AutoAcceptSharedAttachments string     `json:"auto_accept_shared_attachments"` // Whether to auto-accept shared attachments
	DefaultRouteTableAssociation string    `json:"default_route_table_association"` // Whether to auto-associate with default route table
	DefaultRouteTablePropagation string    `json:"default_route_table_propagation"` // Whether to auto-propagate to default route table
	DnsSupport               string        `json:"dns_support"`             // Whether DNS support is enabled
	MulticastSupport         string        `json:"multicast_support"`       // Whether multicast support is enabled
	Tags                     map[string]string `json:"tags"`                // Key-value tags associated with the transit gateway
}

// TransitGatewayAttachmentInfo contains information about a Transit Gateway attachment
type TransitGatewayAttachmentInfo struct {
	AttachmentID         string            `json:"attachment_id"`           // Unique identifier for the attachment
	TransitGatewayID     string            `json:"transit_gateway_id"`      // ID of the transit gateway
	ResourceType         string            `json:"resource_type"`           // Type of resource (vpc, vpn, direct-connect-gateway, peering)
	ResourceID           string            `json:"resource_id"`             // ID of the attached resource
	ResourceOwnerID      string            `json:"resource_owner_id"`       // AWS account ID that owns the resource
	State                string            `json:"state"`                   // State of the attachment (initiating, pendingAcceptance, rollingBack, pending, available, modifying, deleting, deleted, failed, rejected, rejecting, failing)
	Association          map[string]string `json:"association"`             // Route table association information
	CreationTime         string            `json:"creation_time"`           // Time when the attachment was created
	Tags                 map[string]string `json:"tags"`                    // Key-value tags associated with the attachment
}

// Scanner provides methods for retrieving VPC and related AWS networking information
type Scanner struct {
	ec2Client *ec2.Client // AWS EC2 client for making API calls
}

// NewScanner creates a new VPC scanner instance with the provided AWS configuration
// cfg: AWS configuration containing credentials and region information
func NewScanner(cfg aws.Config) *Scanner {
	return &Scanner{
		ec2Client: ec2.NewFromConfig(cfg),
	}
}

// GetVPCs retrieves information about all VPCs in the configured AWS region
// ctx: Context for the request, allowing for timeout and cancellation
// Returns: Slice of VPCInfo structs containing VPC details, or error if the operation fails
func (s *Scanner) GetVPCs(ctx context.Context) ([]VPCInfo, error) {
	// Prepare input for describing all VPCs (no filters applied)
	input := &ec2.DescribeVpcsInput{}

	// Call AWS API to retrieve VPC information
	result, err := s.ec2Client.DescribeVpcs(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe VPCs: %w", err)
	}

	// Process each VPC from the API response
	var vpcs []VPCInfo
	for _, vpc := range result.Vpcs {
		// Extract basic VPC information
		vpcInfo := VPCInfo{
			VpcID:           aws.ToString(vpc.VpcId),
			CidrBlock:       aws.ToString(vpc.CidrBlock),
			State:           string(vpc.State),
			IsDefault:       aws.ToBool(vpc.IsDefault),
			DhcpOptionsID:   aws.ToString(vpc.DhcpOptionsId),
			InstanceTenancy: string(vpc.InstanceTenancy),
			Tags:            convertTags(vpc.Tags),
		}

		// Collect all associated CIDR blocks beyond the primary one
		for _, cidr := range vpc.CidrBlockAssociationSet {
			if cidr.CidrBlock != nil {
				vpcInfo.AssociateCidrBlocks = append(vpcInfo.AssociateCidrBlocks, *cidr.CidrBlock)
			}
		}

		vpcs = append(vpcs, vpcInfo)
	}

	return vpcs, nil
}

// GetSubnets retrieves information about all subnets across all VPCs in the configured AWS region
// ctx: Context for the request, allowing for timeout and cancellation
// Returns: Slice of SubnetInfo structs containing subnet details, or error if the operation fails
func (s *Scanner) GetSubnets(ctx context.Context) ([]SubnetInfo, error) {
	// Prepare input for describing all subnets (no filters applied)
	input := &ec2.DescribeSubnetsInput{}

	// Call AWS API to retrieve subnet information
	result, err := s.ec2Client.DescribeSubnets(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe subnets: %w", err)
	}

	// Process each subnet from the API response
	var subnets []SubnetInfo
	for _, subnet := range result.Subnets {
		// Extract subnet information and convert AWS types to our struct format
		subnetInfo := SubnetInfo{
			SubnetID:                    aws.ToString(subnet.SubnetId),
			VpcID:                       aws.ToString(subnet.VpcId),
			CidrBlock:                   aws.ToString(subnet.CidrBlock),
			AvailabilityZone:            aws.ToString(subnet.AvailabilityZone),
			AvailabilityZoneID:          aws.ToString(subnet.AvailabilityZoneId),
			State:                       string(subnet.State),
			MapPublicIpOnLaunch:         aws.ToBool(subnet.MapPublicIpOnLaunch),
			AssignIpv6AddressOnCreation: aws.ToBool(subnet.AssignIpv6AddressOnCreation),
			DefaultForAz:                aws.ToBool(subnet.DefaultForAz),
			Tags:                        convertTags(subnet.Tags),
		}
		subnets = append(subnets, subnetInfo)
	}

	return subnets, nil
}

// GetSubnetsByVPC retrieves information about all subnets within a specific VPC
// ctx: Context for the request, allowing for timeout and cancellation
// vpcID: The unique identifier of the VPC to filter subnets by
// Returns: Slice of SubnetInfo structs for subnets in the specified VPC, or error if the operation fails
func (s *Scanner) GetSubnetsByVPC(ctx context.Context, vpcID string) ([]SubnetInfo, error) {
	// Prepare input with VPC ID filter to retrieve only subnets in the specified VPC
	input := &ec2.DescribeSubnetsInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("vpc-id"), // Filter by VPC ID
				Values: []string{vpcID},
			},
		},
	}

	// Call AWS API to retrieve subnet information for the specific VPC
	result, err := s.ec2Client.DescribeSubnets(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe subnets for VPC %s: %w", vpcID, err)
	}

	// Process each subnet from the API response
	var subnets []SubnetInfo
	for _, subnet := range result.Subnets {
		// Extract subnet information and convert AWS types to our struct format
		subnetInfo := SubnetInfo{
			SubnetID:                    aws.ToString(subnet.SubnetId),
			VpcID:                       aws.ToString(subnet.VpcId),
			CidrBlock:                   aws.ToString(subnet.CidrBlock),
			AvailabilityZone:            aws.ToString(subnet.AvailabilityZone),
			AvailabilityZoneID:          aws.ToString(subnet.AvailabilityZoneId),
			State:                       string(subnet.State),
			MapPublicIpOnLaunch:         aws.ToBool(subnet.MapPublicIpOnLaunch),
			AssignIpv6AddressOnCreation: aws.ToBool(subnet.AssignIpv6AddressOnCreation),
			DefaultForAz:                aws.ToBool(subnet.DefaultForAz),
			Tags:                        convertTags(subnet.Tags),
		}
		subnets = append(subnets, subnetInfo)
	}

	return subnets, nil
}

// GetRouteTables retrieves information about all route tables in the configured AWS region
// ctx: Context for the request, allowing for timeout and cancellation
// Returns: Slice of RouteTableInfo structs containing route table details, or error if the operation fails
func (s *Scanner) GetRouteTables(ctx context.Context) ([]RouteTableInfo, error) {
	// Prepare input for describing all route tables (no filters applied)
	input := &ec2.DescribeRouteTablesInput{}

	// Call AWS API to retrieve route table information
	result, err := s.ec2Client.DescribeRouteTables(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe route tables: %w", err)
	}

	// Process each route table from the API response
	var routeTables []RouteTableInfo
	for _, rt := range result.RouteTables {
		// Extract basic route table information
		routeTableInfo := RouteTableInfo{
			RouteTableID:     aws.ToString(rt.RouteTableId),
			VpcID:            aws.ToString(rt.VpcId),
			IsMainRouteTable: false, // Will be determined by checking associations
			Tags:             convertTags(rt.Tags),
		}

		// Process routes in the route table
		for _, route := range rt.Routes {
			routeInfo := RouteInfo{
				DestinationCidrBlock:   aws.ToString(route.DestinationCidrBlock),
				DestinationIpv6Block:   aws.ToString(route.DestinationIpv6CidrBlock),
				GatewayID:              aws.ToString(route.GatewayId),
				InstanceID:             aws.ToString(route.InstanceId),
				NatGatewayID:           aws.ToString(route.NatGatewayId),
				NetworkInterfaceID:     aws.ToString(route.NetworkInterfaceId),
				TransitGatewayID:       aws.ToString(route.TransitGatewayId),
				VpcPeeringConnectionID: aws.ToString(route.VpcPeeringConnectionId),
				State:                  string(route.State),
				Origin:                 string(route.Origin),
			}
			routeTableInfo.Routes = append(routeTableInfo.Routes, routeInfo)
		}

		// Process subnet associations
		for _, assoc := range rt.Associations {
			if aws.ToBool(assoc.Main) {
				// This is the main route table for the VPC
				routeTableInfo.IsMainRouteTable = true
			} else if assoc.SubnetId != nil {
				// This route table is explicitly associated with a subnet
				routeTableInfo.SubnetIDs = append(routeTableInfo.SubnetIDs, aws.ToString(assoc.SubnetId))
			}
		}

		routeTables = append(routeTables, routeTableInfo)
	}

	return routeTables, nil
}

// GetSecurityGroups retrieves information about all security groups in the configured AWS region
// ctx: Context for the request, allowing for timeout and cancellation
// Returns: Slice of SecurityGroupInfo structs containing security group details, or error if the operation fails
func (s *Scanner) GetSecurityGroups(ctx context.Context) ([]SecurityGroupInfo, error) {
	// Prepare input for describing all security groups (no filters applied)
	input := &ec2.DescribeSecurityGroupsInput{}

	// Call AWS API to retrieve security group information
	result, err := s.ec2Client.DescribeSecurityGroups(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe security groups: %w", err)
	}

	// Process each security group from the API response
	var securityGroups []SecurityGroupInfo
	for _, sg := range result.SecurityGroups {
		// Extract basic security group information
		sgInfo := SecurityGroupInfo{
			GroupID:     aws.ToString(sg.GroupId),
			GroupName:   aws.ToString(sg.GroupName),
			Description: aws.ToString(sg.Description),
			VpcID:       aws.ToString(sg.VpcId),
			OwnerID:     aws.ToString(sg.OwnerId),
			Tags:        convertTags(sg.Tags),
		}

		// Process ingress rules
		for _, rule := range sg.IpPermissions {
			// Each rule can have multiple IP ranges/groups, so we create separate rule entries
			for _, ipRange := range rule.IpRanges {
				sgRule := SecurityGroupRule{
					IsEgress:    false,
					IpProtocol:  aws.ToString(rule.IpProtocol),
					FromPort:    aws.ToInt32(rule.FromPort),
					ToPort:      aws.ToInt32(rule.ToPort),
					CidrBlock:   aws.ToString(ipRange.CidrIp),
					Description: aws.ToString(ipRange.Description),
				}
				sgInfo.Rules = append(sgInfo.Rules, sgRule)
			}

			// Process IPv6 ranges
			for _, ipv6Range := range rule.Ipv6Ranges {
				sgRule := SecurityGroupRule{
					IsEgress:      false,
					IpProtocol:    aws.ToString(rule.IpProtocol),
					FromPort:      aws.ToInt32(rule.FromPort),
					ToPort:        aws.ToInt32(rule.ToPort),
					Ipv6CidrBlock: aws.ToString(ipv6Range.CidrIpv6),
					Description:   aws.ToString(ipv6Range.Description),
				}
				sgInfo.Rules = append(sgInfo.Rules, sgRule)
			}

			// Process referenced security groups
			for _, userIdGroupPair := range rule.UserIdGroupPairs {
				sgRule := SecurityGroupRule{
					IsEgress:     false,
					IpProtocol:   aws.ToString(rule.IpProtocol),
					FromPort:     aws.ToInt32(rule.FromPort),
					ToPort:       aws.ToInt32(rule.ToPort),
					GroupID:      aws.ToString(userIdGroupPair.GroupId),
					GroupOwnerID: aws.ToString(userIdGroupPair.UserId),
					Description:  aws.ToString(userIdGroupPair.Description),
				}
				sgInfo.Rules = append(sgInfo.Rules, sgRule)
			}

			// Process prefix lists
			for _, prefixListId := range rule.PrefixListIds {
				sgRule := SecurityGroupRule{
					IsEgress:     false,
					IpProtocol:   aws.ToString(rule.IpProtocol),
					FromPort:     aws.ToInt32(rule.FromPort),
					ToPort:       aws.ToInt32(rule.ToPort),
					PrefixListID: aws.ToString(prefixListId.PrefixListId),
					Description:  aws.ToString(prefixListId.Description),
				}
				sgInfo.Rules = append(sgInfo.Rules, sgRule)
			}
		}

		// Process egress rules (similar structure to ingress)
		for _, rule := range sg.IpPermissionsEgress {
			// Each rule can have multiple IP ranges/groups
			for _, ipRange := range rule.IpRanges {
				sgRule := SecurityGroupRule{
					IsEgress:    true,
					IpProtocol:  aws.ToString(rule.IpProtocol),
					FromPort:    aws.ToInt32(rule.FromPort),
					ToPort:      aws.ToInt32(rule.ToPort),
					CidrBlock:   aws.ToString(ipRange.CidrIp),
					Description: aws.ToString(ipRange.Description),
				}
				sgInfo.Rules = append(sgInfo.Rules, sgRule)
			}

			// Process IPv6 ranges
			for _, ipv6Range := range rule.Ipv6Ranges {
				sgRule := SecurityGroupRule{
					IsEgress:      true,
					IpProtocol:    aws.ToString(rule.IpProtocol),
					FromPort:      aws.ToInt32(rule.FromPort),
					ToPort:        aws.ToInt32(rule.ToPort),
					Ipv6CidrBlock: aws.ToString(ipv6Range.CidrIpv6),
					Description:   aws.ToString(ipv6Range.Description),
				}
				sgInfo.Rules = append(sgInfo.Rules, sgRule)
			}

			// Process referenced security groups
			for _, userIdGroupPair := range rule.UserIdGroupPairs {
				sgRule := SecurityGroupRule{
					IsEgress:     true,
					IpProtocol:   aws.ToString(rule.IpProtocol),
					FromPort:     aws.ToInt32(rule.FromPort),
					ToPort:       aws.ToInt32(rule.ToPort),
					GroupID:      aws.ToString(userIdGroupPair.GroupId),
					GroupOwnerID: aws.ToString(userIdGroupPair.UserId),
					Description:  aws.ToString(userIdGroupPair.Description),
				}
				sgInfo.Rules = append(sgInfo.Rules, sgRule)
			}

			// Process prefix lists
			for _, prefixListId := range rule.PrefixListIds {
				sgRule := SecurityGroupRule{
					IsEgress:     true,
					IpProtocol:   aws.ToString(rule.IpProtocol),
					FromPort:     aws.ToInt32(rule.FromPort),
					ToPort:       aws.ToInt32(rule.ToPort),
					PrefixListID: aws.ToString(prefixListId.PrefixListId),
					Description:  aws.ToString(prefixListId.Description),
				}
				sgInfo.Rules = append(sgInfo.Rules, sgRule)
			}
		}

		securityGroups = append(securityGroups, sgInfo)
	}

	return securityGroups, nil
}

// GetInternetGateways retrieves information about all internet gateways in the configured AWS region
// ctx: Context for the request, allowing for timeout and cancellation
// Returns: Slice of InternetGatewayInfo structs containing internet gateway details, or error if the operation fails
func (s *Scanner) GetInternetGateways(ctx context.Context) ([]InternetGatewayInfo, error) {
	// Prepare input for describing all internet gateways (no filters applied)
	input := &ec2.DescribeInternetGatewaysInput{}

	// Call AWS API to retrieve internet gateway information
	result, err := s.ec2Client.DescribeInternetGateways(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe internet gateways: %w", err)
	}

	// Process each internet gateway from the API response
	var internetGateways []InternetGatewayInfo
	for _, igw := range result.InternetGateways {
		// Extract basic internet gateway information
		igwInfo := InternetGatewayInfo{
			InternetGatewayID: aws.ToString(igw.InternetGatewayId),
			Tags:              convertTags(igw.Tags),
		}

		// Determine state and VPC association
		if len(igw.Attachments) > 0 {
			// Internet gateway is attached to a VPC
			attachment := igw.Attachments[0] // IGW can only be attached to one VPC
			igwInfo.State = string(attachment.State)
			igwInfo.VpcID = aws.ToString(attachment.VpcId)
		} else {
			// Internet gateway is not attached
			igwInfo.State = "available"
		}

		internetGateways = append(internetGateways, igwInfo)
	}

	return internetGateways, nil
}

// GetNatGateways retrieves information about all NAT gateways in the configured AWS region
// ctx: Context for the request, allowing for timeout and cancellation
// Returns: Slice of NatGatewayInfo structs containing NAT gateway details, or error if the operation fails
func (s *Scanner) GetNatGateways(ctx context.Context) ([]NatGatewayInfo, error) {
	// Prepare input for describing all NAT gateways (no filters applied)
	input := &ec2.DescribeNatGatewaysInput{}

	// Call AWS API to retrieve NAT gateway information
	result, err := s.ec2Client.DescribeNatGateways(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe NAT gateways: %w", err)
	}

	// Process each NAT gateway from the API response
	var natGateways []NatGatewayInfo
	for _, ngw := range result.NatGateways {
		// Extract basic NAT gateway information
		ngwInfo := NatGatewayInfo{
			NatGatewayID:     aws.ToString(ngw.NatGatewayId),
			SubnetID:         aws.ToString(ngw.SubnetId),
			VpcID:            aws.ToString(ngw.VpcId),
			State:            string(ngw.State),
			ConnectivityType: string(ngw.ConnectivityType),
			Tags:             convertTags(ngw.Tags),
		}

		// Set creation time
		if ngw.CreateTime != nil {
			ngwInfo.CreatedTime = ngw.CreateTime.Format("2006-01-02T15:04:05Z")
		}

		// Process NAT gateway addresses to get IP information
		for _, addr := range ngw.NatGatewayAddresses {
			if addr.NetworkInterfaceId != nil {
				ngwInfo.NetworkInterfaceID = aws.ToString(addr.NetworkInterfaceId)
			}
			if addr.PrivateIp != nil {
				ngwInfo.PrivateIp = aws.ToString(addr.PrivateIp)
			}
			if addr.PublicIp != nil {
				ngwInfo.PublicIp = aws.ToString(addr.PublicIp)
			}
			if addr.AllocationId != nil {
				ngwInfo.AllocationID = aws.ToString(addr.AllocationId)
			}
		}

		natGateways = append(natGateways, ngwInfo)
	}

	return natGateways, nil
}

// GetTransitGateways retrieves information about all transit gateways in the configured AWS region
// ctx: Context for the request, allowing for timeout and cancellation
// Returns: Slice of TransitGatewayInfo structs containing transit gateway details, or error if the operation fails
func (s *Scanner) GetTransitGateways(ctx context.Context) ([]TransitGatewayInfo, error) {
	// Prepare input for describing all transit gateways (no filters applied)
	input := &ec2.DescribeTransitGatewaysInput{}

	// Call AWS API to retrieve transit gateway information
	result, err := s.ec2Client.DescribeTransitGateways(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe transit gateways: %w", err)
	}

	// Process each transit gateway from the API response
	var transitGateways []TransitGatewayInfo
	for _, tgw := range result.TransitGateways {
		// Extract basic transit gateway information
		tgwInfo := TransitGatewayInfo{
			TransitGatewayID: aws.ToString(tgw.TransitGatewayId),
			State:            string(tgw.State),
			OwnerID:          aws.ToString(tgw.OwnerId),
			Description:      aws.ToString(tgw.Description),
			Tags:             convertTags(tgw.Tags),
		}

		// Set creation time
		if tgw.CreationTime != nil {
			tgwInfo.CreationTime = tgw.CreationTime.Format("2006-01-02T15:04:05Z")
		}

		// Process transit gateway options
		if tgw.Options != nil {
			options := tgw.Options
			tgwInfo.AmazonSideAsn = aws.ToInt64(options.AmazonSideAsn)
			tgwInfo.AutoAcceptSharedAttachments = string(options.AutoAcceptSharedAttachments)
			tgwInfo.DefaultRouteTableAssociation = string(options.DefaultRouteTableAssociation)
			tgwInfo.DefaultRouteTablePropagation = string(options.DefaultRouteTablePropagation)
			tgwInfo.DnsSupport = string(options.DnsSupport)
			tgwInfo.MulticastSupport = string(options.MulticastSupport)
			tgwInfo.DefaultRouteTableID = aws.ToString(options.AssociationDefaultRouteTableId)
			tgwInfo.PropagationRouteTableID = aws.ToString(options.PropagationDefaultRouteTableId)
		}

		transitGateways = append(transitGateways, tgwInfo)
	}

	return transitGateways, nil
}

// GetTransitGatewayAttachments retrieves information about all transit gateway attachments in the configured AWS region
// ctx: Context for the request, allowing for timeout and cancellation
// Returns: Slice of TransitGatewayAttachmentInfo structs containing attachment details, or error if the operation fails
func (s *Scanner) GetTransitGatewayAttachments(ctx context.Context) ([]TransitGatewayAttachmentInfo, error) {
	// Prepare input for describing all transit gateway attachments (no filters applied)
	input := &ec2.DescribeTransitGatewayAttachmentsInput{}

	// Call AWS API to retrieve transit gateway attachment information
	result, err := s.ec2Client.DescribeTransitGatewayAttachments(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe transit gateway attachments: %w", err)
	}

	// Process each attachment from the API response
	var attachments []TransitGatewayAttachmentInfo
	for _, attachment := range result.TransitGatewayAttachments {
		// Extract basic attachment information
		attachmentInfo := TransitGatewayAttachmentInfo{
			AttachmentID:     aws.ToString(attachment.TransitGatewayAttachmentId),
			TransitGatewayID: aws.ToString(attachment.TransitGatewayId),
			ResourceType:     string(attachment.ResourceType),
			ResourceID:       aws.ToString(attachment.ResourceId),
			ResourceOwnerID:  aws.ToString(attachment.ResourceOwnerId),
			State:            string(attachment.State),
			Tags:             convertTags(attachment.Tags),
			Association:      make(map[string]string),
		}

		// Set creation time
		if attachment.CreationTime != nil {
			attachmentInfo.CreationTime = attachment.CreationTime.Format("2006-01-02T15:04:05Z")
		}

		// Process association information
		if attachment.Association != nil {
			assoc := attachment.Association
			attachmentInfo.Association["route_table_id"] = aws.ToString(assoc.TransitGatewayRouteTableId)
			attachmentInfo.Association["state"] = string(assoc.State)
		}

		attachments = append(attachments, attachmentInfo)
	}

	return attachments, nil
}

// convertTags converts AWS tag format to a simple key-value map
// tags: Slice of AWS Tag structs containing Key and Value pointers
// Returns: Map of string keys to string values, skipping any nil keys or values
func convertTags(tags []types.Tag) map[string]string {
	result := make(map[string]string)
	for _, tag := range tags {
		// Only include tags that have both key and value set (not nil)
		if tag.Key != nil && tag.Value != nil {
			result[*tag.Key] = *tag.Value
		}
	}
	return result
}
