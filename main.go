package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"aws-documentor/modules/diagram"
	"aws-documentor/modules/vpc"
)

func main() {
	// Parse command-line flags
	region := flag.String("region", "", "AWS region to scan (optional, uses default config if not specified)")
	generateDiagram := flag.Bool("diagram", false, "Generate draw.io diagram file (saves to vpc-diagram.drawio)")
	outputJSON := flag.Bool("json", true, "Output JSON data to stdout (default: true)")
	flag.Parse()

	ctx := context.Background()

	// Load AWS config with optional region override
	var cfg aws.Config
	var err error
	if *region != "" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(*region))
		fmt.Printf("Scanning AWS region: %s\n\n", *region)
	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
		fmt.Printf("Scanning AWS region: %s (from default config)\n\n", cfg.Region)
	}
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	scanner := vpc.NewScanner(cfg)

	fmt.Println("Scanning VPCs...")
	vpcs, err := scanner.GetVPCs(ctx)
	if err != nil {
		log.Fatalf("Failed to get VPCs: %v", err)
	}

	if *outputJSON {
		fmt.Printf("Found %d VPCs:\n", len(vpcs))
		for _, v := range vpcs {
			vpcJSON, _ := json.MarshalIndent(v, "", "  ")
			fmt.Printf("%s\n", vpcJSON)
			fmt.Println("---")
		}
	} else {
		fmt.Printf("Found %d VPCs\n", len(vpcs))
	}

	fmt.Println("\nScanning Subnets...")
	subnets, err := scanner.GetSubnets(ctx)
	if err != nil {
		log.Fatalf("Failed to get subnets: %v", err)
	}

	if *outputJSON {
		fmt.Printf("Found %d Subnets:\n", len(subnets))
		for _, s := range subnets {
			subnetJSON, _ := json.MarshalIndent(s, "", "  ")
			fmt.Printf("%s\n", subnetJSON)
			fmt.Println("---")
		}
	} else {
		fmt.Printf("Found %d Subnets\n", len(subnets))
	}

	fmt.Println("\nScanning Route Tables...")
	routeTables, err := scanner.GetRouteTables(ctx)
	if err != nil {
		log.Fatalf("Failed to get route tables: %v", err)
	}

	if *outputJSON {
		fmt.Printf("Found %d Route Tables:\n", len(routeTables))
		for _, rt := range routeTables {
			routeTableJSON, _ := json.MarshalIndent(rt, "", "  ")
			fmt.Printf("%s\n", routeTableJSON)
			fmt.Println("---")
		}
	} else {
		fmt.Printf("Found %d Route Tables\n", len(routeTables))
	}

	fmt.Println("\nScanning Security Groups...")
	securityGroups, err := scanner.GetSecurityGroups(ctx)
	if err != nil {
		log.Fatalf("Failed to get security groups: %v", err)
	}

	if *outputJSON {
		fmt.Printf("Found %d Security Groups:\n", len(securityGroups))
		for _, sg := range securityGroups {
			sgJSON, _ := json.MarshalIndent(sg, "", "  ")
			fmt.Printf("%s\n", sgJSON)
			fmt.Println("---")
		}
	} else {
		fmt.Printf("Found %d Security Groups\n", len(securityGroups))
	}

	fmt.Println("\nScanning Internet Gateways...")
	internetGateways, err := scanner.GetInternetGateways(ctx)
	if err != nil {
		log.Fatalf("Failed to get internet gateways: %v", err)
	}

	if *outputJSON {
		fmt.Printf("Found %d Internet Gateways:\n", len(internetGateways))
		for _, igw := range internetGateways {
			igwJSON, _ := json.MarshalIndent(igw, "", "  ")
			fmt.Printf("%s\n", igwJSON)
			fmt.Println("---")
		}
	} else {
		fmt.Printf("Found %d Internet Gateways\n", len(internetGateways))
	}

	fmt.Println("\nScanning NAT Gateways...")
	natGateways, err := scanner.GetNatGateways(ctx)
	if err != nil {
		log.Fatalf("Failed to get NAT gateways: %v", err)
	}

	if *outputJSON {
		fmt.Printf("Found %d NAT Gateways:\n", len(natGateways))
		for _, ngw := range natGateways {
			ngwJSON, _ := json.MarshalIndent(ngw, "", "  ")
			fmt.Printf("%s\n", ngwJSON)
			fmt.Println("---")
		}
	} else {
		fmt.Printf("Found %d NAT Gateways\n", len(natGateways))
	}

	fmt.Println("\nScanning Transit Gateways...")
	transitGateways, err := scanner.GetTransitGateways(ctx)
	if err != nil {
		log.Fatalf("Failed to get transit gateways: %v", err)
	}

	if *outputJSON {
		fmt.Printf("Found %d Transit Gateways:\n", len(transitGateways))
		for _, tgw := range transitGateways {
			tgwJSON, _ := json.MarshalIndent(tgw, "", "  ")
			fmt.Printf("%s\n", tgwJSON)
			fmt.Println("---")
		}
	} else {
		fmt.Printf("Found %d Transit Gateways\n", len(transitGateways))
	}

	fmt.Println("\nScanning Transit Gateway Attachments...")
	tgwAttachments, err := scanner.GetTransitGatewayAttachments(ctx)
	if err != nil {
		log.Fatalf("Failed to get transit gateway attachments: %v", err)
	}

	if *outputJSON {
		fmt.Printf("Found %d Transit Gateway Attachments:\n", len(tgwAttachments))
		for _, attachment := range tgwAttachments {
			attachmentJSON, _ := json.MarshalIndent(attachment, "", "  ")
			fmt.Printf("%s\n", attachmentJSON)
			fmt.Println("---")
		}
	} else {
		fmt.Printf("Found %d Transit Gateway Attachments\n", len(tgwAttachments))
	}

	fmt.Println("\nVPC infrastructure scan complete!")

	// Generate diagram if requested
	if *generateDiagram {
		fmt.Println("\nGenerating draw.io diagram...")
		diagramGen := diagram.NewDiagramGenerator()

		diagramXML, err := diagramGen.GenerateVPCDiagram(
			vpcs,
			subnets,
			routeTables,
			securityGroups,
			internetGateways,
			natGateways,
			transitGateways,
			tgwAttachments,
		)
		if err != nil {
			log.Fatalf("Failed to generate diagram: %v", err)
		}

		// Write diagram to file
		filename := "vpc-diagram.drawio"
		err = os.WriteFile(filename, []byte(diagramXML), 0644)
		if err != nil {
			log.Fatalf("Failed to write diagram file: %v", err)
		}

		fmt.Printf("Diagram saved to: %s\n", filename)
		fmt.Println("You can open this file in draw.io (https://app.diagrams.net)")
	}
}
