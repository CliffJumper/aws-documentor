package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"

	"aws-documentor/modules/vpc"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	scanner := vpc.NewScanner(cfg)

	fmt.Println("Scanning VPCs...")
	vpcs, err := scanner.GetVPCs(ctx)
	if err != nil {
		log.Fatalf("Failed to get VPCs: %v", err)
	}

	fmt.Printf("Found %d VPCs:\n", len(vpcs))
	for _, v := range vpcs {
		vpcJSON, _ := json.MarshalIndent(v, "", "  ")
		fmt.Printf("%s\n", vpcJSON)
		fmt.Println("---")
	}

	fmt.Println("\nScanning Subnets...")
	subnets, err := scanner.GetSubnets(ctx)
	if err != nil {
		log.Fatalf("Failed to get subnets: %v", err)
	}

	fmt.Printf("Found %d Subnets:\n", len(subnets))
	for _, s := range subnets {
		subnetJSON, _ := json.MarshalIndent(s, "", "  ")
		fmt.Printf("%s\n", subnetJSON)
		fmt.Println("---")
	}

	fmt.Println("\nScanning Route Tables...")
	routeTables, err := scanner.GetRouteTables(ctx)
	if err != nil {
		log.Fatalf("Failed to get route tables: %v", err)
	}

	fmt.Printf("Found %d Route Tables:\n", len(routeTables))
	for _, rt := range routeTables {
		routeTableJSON, _ := json.MarshalIndent(rt, "", "  ")
		fmt.Printf("%s\n", routeTableJSON)
		fmt.Println("---")
	}

	fmt.Println("\nScanning Security Groups...")
	securityGroups, err := scanner.GetSecurityGroups(ctx)
	if err != nil {
		log.Fatalf("Failed to get security groups: %v", err)
	}

	fmt.Printf("Found %d Security Groups:\n", len(securityGroups))
	for _, sg := range securityGroups {
		sgJSON, _ := json.MarshalIndent(sg, "", "  ")
		fmt.Printf("%s\n", sgJSON)
		fmt.Println("---")
	}

	fmt.Println("\nScanning Internet Gateways...")
	internetGateways, err := scanner.GetInternetGateways(ctx)
	if err != nil {
		log.Fatalf("Failed to get internet gateways: %v", err)
	}

	fmt.Printf("Found %d Internet Gateways:\n", len(internetGateways))
	for _, igw := range internetGateways {
		igwJSON, _ := json.MarshalIndent(igw, "", "  ")
		fmt.Printf("%s\n", igwJSON)
		fmt.Println("---")
	}

	fmt.Println("\nScanning NAT Gateways...")
	natGateways, err := scanner.GetNatGateways(ctx)
	if err != nil {
		log.Fatalf("Failed to get NAT gateways: %v", err)
	}

	fmt.Printf("Found %d NAT Gateways:\n", len(natGateways))
	for _, ngw := range natGateways {
		ngwJSON, _ := json.MarshalIndent(ngw, "", "  ")
		fmt.Printf("%s\n", ngwJSON)
		fmt.Println("---")
	}

	fmt.Println("\nScanning Transit Gateways...")
	transitGateways, err := scanner.GetTransitGateways(ctx)
	if err != nil {
		log.Fatalf("Failed to get transit gateways: %v", err)
	}

	fmt.Printf("Found %d Transit Gateways:\n", len(transitGateways))
	for _, tgw := range transitGateways {
		tgwJSON, _ := json.MarshalIndent(tgw, "", "  ")
		fmt.Printf("%s\n", tgwJSON)
		fmt.Println("---")
	}

	fmt.Println("\nScanning Transit Gateway Attachments...")
	tgwAttachments, err := scanner.GetTransitGatewayAttachments(ctx)
	if err != nil {
		log.Fatalf("Failed to get transit gateway attachments: %v", err)
	}

	fmt.Printf("Found %d Transit Gateway Attachments:\n", len(tgwAttachments))
	for _, attachment := range tgwAttachments {
		attachmentJSON, _ := json.MarshalIndent(attachment, "", "  ")
		fmt.Printf("%s\n", attachmentJSON)
		fmt.Println("---")
	}

	fmt.Println("\nVPC infrastructure scan complete!")
}