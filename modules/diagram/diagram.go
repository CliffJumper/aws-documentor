// Package diagram provides functionality for generating draw.io diagrams from AWS VPC infrastructure data
package diagram

import (
	"encoding/xml"
	"fmt"
	"strings"

	"aws-documentor/modules/vpc"
)

// DrawIO represents the root structure of a draw.io XML file
type DrawIO struct {
	XMLName xml.Name `xml:"mxfile"`
	Host    string   `xml:"host,attr"`
	Version string   `xml:"version,attr"`
	Type    string   `xml:"type,attr"`
	Diagram Diagram  `xml:"diagram"`
}

// Diagram represents a diagram within the draw.io file
type Diagram struct {
	Name         string       `xml:"name,attr"`
	ID           string       `xml:"id,attr"`
	MxGraphModel MxGraphModel `xml:"mxGraphModel"`
}

// MxGraphModel represents the graph model containing all shapes and connections
type MxGraphModel struct {
	Grid      int     `xml:"grid,attr"`
	GridSize  int     `xml:"gridSize,attr"`
	Page      int     `xml:"page,attr"`
	PageScale float64 `xml:"pageScale,attr"`
	Root      Root    `xml:"root"`
}

// Root contains all cells (shapes, connections, etc.)
type Root struct {
	Cells []Cell `xml:"mxCell"`
}

// Cell represents a shape, connection, or container in the diagram
type Cell struct {
	ID       string    `xml:"id,attr"`
	Value    string    `xml:"value,attr,omitempty"`
	Style    string    `xml:"style,attr,omitempty"`
	Parent   string    `xml:"parent,attr,omitempty"`
	Vertex   string    `xml:"vertex,attr,omitempty"`
	Edge     string    `xml:"edge,attr,omitempty"`
	Geometry *Geometry `xml:"mxGeometry,omitempty"`
}

// Geometry defines the position and size of a cell
type Geometry struct {
	X      float64 `xml:"x,attr,omitempty"`
	Y      float64 `xml:"y,attr,omitempty"`
	Width  float64 `xml:"width,attr,omitempty"`
	Height float64 `xml:"height,attr,omitempty"`
	As     string  `xml:"as,attr"`
}

// DiagramGenerator generates draw.io diagrams from VPC data
type DiagramGenerator struct {
	cellIDCounter int
}

// NewDiagramGenerator creates a new diagram generator
func NewDiagramGenerator() *DiagramGenerator {
	return &DiagramGenerator{
		cellIDCounter: 2, // Start at 2 (0 and 1 are reserved for root cells)
	}
}

// nextID generates the next unique cell ID
func (dg *DiagramGenerator) nextID() string {
	id := fmt.Sprintf("cell-%d", dg.cellIDCounter)
	dg.cellIDCounter++
	return id
}

// GenerateVPCDiagram creates a comprehensive VPC architecture diagram
func (dg *DiagramGenerator) GenerateVPCDiagram(
	vpcs []vpc.VPCInfo,
	subnets []vpc.SubnetInfo,
	routeTables []vpc.RouteTableInfo,
	securityGroups []vpc.SecurityGroupInfo,
	internetGateways []vpc.InternetGatewayInfo,
	natGateways []vpc.NatGatewayInfo,
	transitGateways []vpc.TransitGatewayInfo,
	tgwAttachments []vpc.TransitGatewayAttachmentInfo,
) (string, error) {
	// Create base structure
	drawio := DrawIO{
		Host:    "app.diagrams.net",
		Version: "21.0.0",
		Type:    "device",
		Diagram: Diagram{
			Name: "AWS VPC Infrastructure",
			ID:   "vpc-diagram",
			MxGraphModel: MxGraphModel{
				Grid:      1,
				GridSize:  10,
				Page:      1,
				PageScale: 1,
				Root: Root{
					Cells: []Cell{
						{ID: "0"},
						{ID: "1", Parent: "0"},
					},
				},
			},
		},
	}

	// Build diagram cells
	var cells []Cell

	// Generate VPC containers with their contents
	xOffset := 50.0
	for _, v := range vpcs {
		vpcCells := dg.generateVPCContainer(v, subnets, internetGateways, natGateways, xOffset, 50)
		cells = append(cells, vpcCells...)
		xOffset += 1200 // Space between VPCs
	}

	// Generate Transit Gateway section if present
	if len(transitGateways) > 0 {
		tgwCells := dg.generateTransitGatewaySection(transitGateways, tgwAttachments, vpcs, 50, xOffset+100)
		cells = append(cells, tgwCells...)
	}

	// Add all cells to the root
	drawio.Diagram.MxGraphModel.Root.Cells = append(drawio.Diagram.MxGraphModel.Root.Cells, cells...)

	// Marshal to XML
	output, err := xml.MarshalIndent(drawio, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal diagram XML: %w", err)
	}

	return xml.Header + string(output), nil
}

// generateVPCContainer creates a VPC container with subnets and gateways
func (dg *DiagramGenerator) generateVPCContainer(
	vpcInfo vpc.VPCInfo,
	allSubnets []vpc.SubnetInfo,
	allIGWs []vpc.InternetGatewayInfo,
	allNGWs []vpc.NatGatewayInfo,
	x, y float64,
) []Cell {
	var cells []Cell

	// Get subnets for this VPC
	var vpcSubnets []vpc.SubnetInfo
	for _, subnet := range allSubnets {
		if subnet.VpcID == vpcInfo.VpcID {
			vpcSubnets = append(vpcSubnets, subnet)
		}
	}

	// Get IGWs for this VPC
	var vpcIGWs []vpc.InternetGatewayInfo
	for _, igw := range allIGWs {
		if igw.VpcID == vpcInfo.VpcID {
			vpcIGWs = append(vpcIGWs, igw)
		}
	}

	// Get NAT Gateways for this VPC
	var vpcNGWs []vpc.NatGatewayInfo
	for _, ngw := range allNGWs {
		if ngw.VpcID == vpcInfo.VpcID {
			vpcNGWs = append(vpcNGWs, ngw)
		}
	}

	// Separate public and private subnets
	var publicSubnets []vpc.SubnetInfo
	var privateSubnets []vpc.SubnetInfo
	for _, subnet := range vpcSubnets {
		if subnet.MapPublicIpOnLaunch {
			publicSubnets = append(publicSubnets, subnet)
		} else {
			privateSubnets = append(privateSubnets, subnet)
		}
	}

	// Calculate VPC container size based on content
	maxSubnets := len(publicSubnets)
	if len(privateSubnets) > maxSubnets {
		maxSubnets = len(privateSubnets)
	}

	vpcWidth := 250.0 + float64(maxSubnets)*240.0 // IGW space + subnet width * count
	vpcHeight := 400.0 // Fixed height for two rows of subnets

	// Create VPC container with AWS VPC style
	vpcID := dg.nextID()
	vpcName := getResourceName(vpcInfo.Tags, vpcInfo.VpcID)
	vpcLabel := fmt.Sprintf("VPC\n%s\n%s", vpcName, vpcInfo.CidrBlock)

	vpcCell := Cell{
		ID:    vpcID,
		Value: escapeXML(vpcLabel),
		Style: "points=[[0,0],[0.25,0],[0.5,0],[0.75,0],[1,0],[1,0.25],[1,0.5],[1,0.75],[1,1],[0.75,1],[0.5,1],[0.25,1],[0,1],[0,0.75],[0,0.5],[0,0.25]];outlineConnect=0;gradientColor=none;html=1;whiteSpace=wrap;fontSize=12;fontStyle=0;container=1;pointerEvents=0;collapsible=0;recursiveResize=0;shape=mxgraph.aws4.group;grIcon=mxgraph.aws4.group_vpc2;strokeColor=#8C4FFF;fillColor=none;verticalAlign=top;align=left;spacingLeft=30;fontColor=#AAB7B8;dashed=0;",
		Parent: "1",
		Vertex: "1",
		Geometry: &Geometry{
			X:      x,
			Y:      y,
			Width:  vpcWidth,
			Height: vpcHeight,
			As:     "geometry",
		},
	}
	cells = append(cells, vpcCell)

	// Add Internet Gateways (vertical stack on the left)
	igwY := 40.0
	for _, igw := range vpcIGWs {
		igwCell := dg.createInternetGatewayCell(igw, vpcID, 20, igwY)
		cells = append(cells, igwCell)
		igwY += 90
	}

	// Add public subnets horizontally (top row)
	subnetX := 150.0
	subnetY := 40.0
	for _, subnet := range publicSubnets {
		subnetCells := dg.createSubnetCell(subnet, vpcID, subnetX, subnetY)
		cells = append(cells, subnetCells...)

		// Check if this subnet has a NAT Gateway
		for _, ngw := range vpcNGWs {
			if ngw.SubnetID == subnet.SubnetID {
				ngwCell := dg.createNATGatewayCell(ngw, subnet.SubnetID, 40, 50)
				cells = append(cells, ngwCell)
			}
		}

		subnetX += 240.0 // Move right for next subnet
	}

	// Add private subnets horizontally (bottom row)
	subnetX = 150.0
	subnetY = 220.0 // Below public subnets
	for _, subnet := range privateSubnets {
		subnetCells := dg.createSubnetCell(subnet, vpcID, subnetX, subnetY)
		cells = append(cells, subnetCells...)

		subnetX += 240.0 // Move right for next subnet
	}

	return cells
}

// createSubnetCell creates a subnet cell with details
func (dg *DiagramGenerator) createSubnetCell(subnet vpc.SubnetInfo, parentID string, x, y float64) []Cell {
	var cells []Cell

	subnetID := dg.nextID()
	subnetName := getResourceName(subnet.Tags, subnet.SubnetID)
	subnetType := "Private subnet"
	subnetStyle := "points=[[0,0],[0.25,0],[0.5,0],[0.75,0],[1,0],[1,0.25],[1,0.5],[1,0.75],[1,1],[0.75,1],[0.5,1],[0.25,1],[0,1],[0,0.75],[0,0.5],[0,0.25]];outlineConnect=0;gradientColor=none;html=1;whiteSpace=wrap;fontSize=12;fontStyle=0;container=1;pointerEvents=0;collapsible=0;recursiveResize=0;shape=mxgraph.aws4.group;grIcon=mxgraph.aws4.group_security_group;grStroke=0;strokeColor=#00A4A6;fillColor=#E6F6F7;verticalAlign=top;align=left;spacingLeft=30;fontColor=#147EBA;dashed=0;"

	if subnet.MapPublicIpOnLaunch {
		subnetType = "Public subnet"
		subnetStyle = "points=[[0,0],[0.25,0],[0.5,0],[0.75,0],[1,0],[1,0.25],[1,0.5],[1,0.75],[1,1],[0.75,1],[0.5,1],[0.25,1],[0,1],[0,0.75],[0,0.5],[0,0.25]];outlineConnect=0;gradientColor=none;html=1;whiteSpace=wrap;fontSize=12;fontStyle=0;container=1;pointerEvents=0;collapsible=0;recursiveResize=0;shape=mxgraph.aws4.group;grIcon=mxgraph.aws4.group_security_group;grStroke=0;strokeColor=#7AA116;fillColor=#F2F6E8;verticalAlign=top;align=left;spacingLeft=30;fontColor=#248814;dashed=0;"
	}

	subnetLabel := fmt.Sprintf("%s\n%s\n%s\nAZ: %s", subnetType, subnetName, subnet.CidrBlock, subnet.AvailabilityZone)

	subnetCell := Cell{
		ID:     subnetID,
		Value:  escapeXML(subnetLabel),
		Style:  subnetStyle,
		Parent: parentID,
		Vertex: "1",
		Geometry: &Geometry{
			X:      x,
			Y:      y,
			Width:  200,
			Height: 140,
			As:     "geometry",
		},
	}
	cells = append(cells, subnetCell)

	return cells
}

// createInternetGatewayCell creates an Internet Gateway cell
func (dg *DiagramGenerator) createInternetGatewayCell(igw vpc.InternetGatewayInfo, parentID string, x, y float64) Cell {
	igwName := getResourceName(igw.Tags, igw.InternetGatewayID)
	igwLabel := fmt.Sprintf("Internet Gateway\n%s", igwName)

	return Cell{
		ID:     dg.nextID(),
		Value:  escapeXML(igwLabel),
		Style:  "sketch=0;outlineConnect=0;fontColor=#232F3E;gradientColor=none;fillColor=#8C4FFF;strokeColor=none;dashed=0;verticalLabelPosition=bottom;verticalAlign=top;align=center;html=1;fontSize=12;fontStyle=0;aspect=fixed;pointerEvents=1;shape=mxgraph.aws4.internet_gateway;",
		Parent: parentID,
		Vertex: "1",
		Geometry: &Geometry{
			X:      x,
			Y:      y,
			Width:  78,
			Height: 78,
			As:     "geometry",
		},
	}
}

// createNATGatewayCell creates a NAT Gateway cell
func (dg *DiagramGenerator) createNATGatewayCell(ngw vpc.NatGatewayInfo, parentID string, x, y float64) Cell {
	ngwName := getResourceName(ngw.Tags, ngw.NatGatewayID)
	ngwLabel := fmt.Sprintf("NAT Gateway\n%s", ngwName)

	return Cell{
		ID:     dg.nextID(),
		Value:  escapeXML(ngwLabel),
		Style:  "sketch=0;outlineConnect=0;fontColor=#232F3E;gradientColor=none;fillColor=#8C4FFF;strokeColor=none;dashed=0;verticalLabelPosition=bottom;verticalAlign=top;align=center;html=1;fontSize=12;fontStyle=0;aspect=fixed;pointerEvents=1;shape=mxgraph.aws4.nat_gateway;",
		Parent: parentID,
		Vertex: "1",
		Geometry: &Geometry{
			X:      x,
			Y:      y,
			Width:  78,
			Height: 78,
			As:     "geometry",
		},
	}
}

// generateTransitGatewaySection creates Transit Gateway visualization with attachments
func (dg *DiagramGenerator) generateTransitGatewaySection(
	transitGateways []vpc.TransitGatewayInfo,
	tgwAttachments []vpc.TransitGatewayAttachmentInfo,
	vpcs []vpc.VPCInfo,
	x, y float64,
) []Cell {
	var cells []Cell

	for i, tgw := range transitGateways {
		tgwID := dg.nextID()
		tgwName := getResourceName(tgw.Tags, tgw.TransitGatewayID)
		tgwLabel := fmt.Sprintf("Transit Gateway\n%s\nASN: %d", tgwName, tgw.AmazonSideAsn)

		tgwCell := Cell{
			ID:     tgwID,
			Value:  escapeXML(tgwLabel),
			Style:  "sketch=0;outlineConnect=0;fontColor=#232F3E;gradientColor=none;fillColor=#8C4FFF;strokeColor=none;dashed=0;verticalLabelPosition=bottom;verticalAlign=top;align=center;html=1;fontSize=12;fontStyle=0;aspect=fixed;pointerEvents=1;shape=mxgraph.aws4.transit_gateway;",
			Parent: "1",
			Vertex: "1",
			Geometry: &Geometry{
				X:      x,
				Y:      y + float64(i)*150,
				Width:  78,
				Height: 78,
				As:     "geometry",
			},
		}
		cells = append(cells, tgwCell)

		// Add attachment icons
		attachY := y + float64(i)*150 + 100
		for _, attachment := range tgwAttachments {
			if attachment.TransitGatewayID == tgw.TransitGatewayID {
				attachID := dg.nextID()
				attachName := getResourceName(attachment.Tags, attachment.AttachmentID)
				attachLabel := fmt.Sprintf("TGW Attachment\n%s\n%s", attachName, attachment.State)

				attachCell := Cell{
					ID:     attachID,
					Value:  escapeXML(attachLabel),
					Style:  "sketch=0;outlineConnect=0;fontColor=#232F3E;gradientColor=none;fillColor=#8C4FFF;strokeColor=none;dashed=0;verticalLabelPosition=bottom;verticalAlign=top;align=center;html=1;fontSize=12;fontStyle=0;aspect=fixed;pointerEvents=1;shape=mxgraph.aws4.transit_gateway_attachment;",
					Parent: "1",
					Vertex: "1",
					Geometry: &Geometry{
						X:      x + 100,
						Y:      attachY,
						Width:  78,
						Height: 78,
						As:     "geometry",
					},
				}
				cells = append(cells, attachCell)
				attachY += 100
			}
		}
	}

	return cells
}

// getResourceName extracts a friendly name from tags, falling back to the resource ID
func getResourceName(tags map[string]string, resourceID string) string {
	if name, ok := tags["Name"]; ok && name != "" {
		return name
	}
	return resourceID
}

// escapeXML escapes special XML characters for use in cell values
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

// GenerateVPCDetailDiagram creates a detailed diagram for a single VPC
func (dg *DiagramGenerator) GenerateVPCDetailDiagram(
	vpcInfo vpc.VPCInfo,
	subnets []vpc.SubnetInfo,
	routeTables []vpc.RouteTableInfo,
	securityGroups []vpc.SecurityGroupInfo,
	internetGateways []vpc.InternetGatewayInfo,
	natGateways []vpc.NatGatewayInfo,
) (string, error) {
	// Create base structure
	drawio := DrawIO{
		Host:    "app.diagrams.net",
		Version: "21.0.0",
		Type:    "device",
		Diagram: Diagram{
			Name: fmt.Sprintf("VPC Detail: %s", getResourceName(vpcInfo.Tags, vpcInfo.VpcID)),
			ID:   "vpc-detail-diagram",
			MxGraphModel: MxGraphModel{
				Grid:      1,
				GridSize:  10,
				Page:      1,
				PageScale: 1,
				Root: Root{
					Cells: []Cell{
						{ID: "0"},
						{ID: "1", Parent: "0"},
					},
				},
			},
		},
	}

	// Generate VPC container with all details
	cells := dg.generateVPCContainer(vpcInfo, subnets, internetGateways, natGateways, 50, 50)

	// Add route tables information panel
	if len(routeTables) > 0 {
		rtCells := dg.generateRouteTablePanel(routeTables, vpcInfo.VpcID, 1200, 50)
		cells = append(cells, rtCells...)
	}

	// Add security groups information panel
	if len(securityGroups) > 0 {
		sgCells := dg.generateSecurityGroupPanel(securityGroups, vpcInfo.VpcID, 1200, 400)
		cells = append(cells, sgCells...)
	}

	drawio.Diagram.MxGraphModel.Root.Cells = append(drawio.Diagram.MxGraphModel.Root.Cells, cells...)

	// Marshal to XML
	output, err := xml.MarshalIndent(drawio, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal diagram XML: %w", err)
	}

	return xml.Header + string(output), nil
}

// generateRouteTablePanel creates an information panel for route tables
func (dg *DiagramGenerator) generateRouteTablePanel(routeTables []vpc.RouteTableInfo, vpcID string, x, y float64) []Cell {
	var cells []Cell

	// Filter route tables for this VPC
	var vpcRouteTables []vpc.RouteTableInfo
	for _, rt := range routeTables {
		if rt.VpcID == vpcID {
			vpcRouteTables = append(vpcRouteTables, rt)
		}
	}

	if len(vpcRouteTables) == 0 {
		return cells
	}

	yOffset := y
	for _, rt := range vpcRouteTables {
		rtName := getResourceName(rt.Tags, rt.RouteTableID)
		mainText := ""
		if rt.IsMainRouteTable {
			mainText = " (Main)"
		}

		// Build routes text
		var routesText []string
		for _, route := range rt.Routes {
			dest := route.DestinationCidrBlock
			if dest == "" {
				dest = route.DestinationIpv6Block
			}
			target := route.GatewayID
			if target == "" {
				target = route.NatGatewayID
			}
			if target == "" {
				target = route.TransitGatewayID
			}
			if target == "" {
				target = "local"
			}
			routesText = append(routesText, fmt.Sprintf("  %s → %s", dest, target))
		}

		rtLabel := fmt.Sprintf("Route Table%s\n%s\n%s", mainText, rtName, strings.Join(routesText, "\n"))

		rtCell := Cell{
			ID:     dg.nextID(),
			Value:  escapeXML(rtLabel),
			Style:  "rounded=1;whiteSpace=wrap;html=1;fillColor=#f5f5f5;strokeColor=#666666;fontSize=9;align=left;verticalAlign=top;spacingLeft=5;spacingTop=5;",
			Parent: "1",
			Vertex: "1",
			Geometry: &Geometry{
				X:      x,
				Y:      yOffset,
				Width:  300,
				Height: 100 + float64(len(routesText)*15),
				As:     "geometry",
			},
		}
		cells = append(cells, rtCell)
		yOffset += 120 + float64(len(routesText)*15)
	}

	return cells
}

// generateSecurityGroupPanel creates an information panel for security groups
func (dg *DiagramGenerator) generateSecurityGroupPanel(securityGroups []vpc.SecurityGroupInfo, vpcID string, x, y float64) []Cell {
	var cells []Cell

	// Filter security groups for this VPC
	var vpcSecurityGroups []vpc.SecurityGroupInfo
	for _, sg := range securityGroups {
		if sg.VpcID == vpcID {
			vpcSecurityGroups = append(vpcSecurityGroups, sg)
		}
	}

	if len(vpcSecurityGroups) == 0 {
		return cells
	}

	yOffset := y
	for _, sg := range vpcSecurityGroups {
		sgName := getResourceName(sg.Tags, sg.GroupID)

		// Count ingress/egress rules
		ingressCount := 0
		egressCount := 0
		for _, rule := range sg.Rules {
			if rule.IsEgress {
				egressCount++
			} else {
				ingressCount++
			}
		}

		sgLabel := fmt.Sprintf("Security Group\n%s\n%s\nIngress: %d rules\nEgress: %d rules",
			sgName, sg.GroupName, ingressCount, egressCount)

		sgCell := Cell{
			ID:     dg.nextID(),
			Value:  escapeXML(sgLabel),
			Style:  "rounded=1;whiteSpace=wrap;html=1;fillColor=#fff2cc;strokeColor=#d6b656;fontSize=9;align=left;verticalAlign=top;spacingLeft=5;spacingTop=5;",
			Parent: "1",
			Vertex: "1",
			Geometry: &Geometry{
				X:      x,
				Y:      yOffset,
				Width:  280,
				Height: 100,
				As:     "geometry",
			},
		}
		cells = append(cells, sgCell)
		yOffset += 120
	}

	return cells
}
