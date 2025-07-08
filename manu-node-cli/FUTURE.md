# Future Manufacturing Node Fields

## Operational Data
- **Status** - Current state (running, idle, maintenance, error, offline)
- **Efficiency/OEE** - Overall Equipment Effectiveness metrics
- **Current Operation** - What it's doing right now
- **Queue Size** - How many jobs waiting
- **Cycle Time** - Time per operation
- **Setup Time** - Time to switch between operations

## Resource Management
- **Required Skills** - Operator certifications needed
- **Tool Requirements** - Specific tools/fixtures needed
- **Material Compatibility** - What materials it can process
- **Power Requirements** - Electrical needs
- **Consumables** - Things it uses (cutting fluid, filters, etc.)

## Capacity & Constraints
- **Max Capacity** - Units per hour/day
- **Size Constraints** - Max/min part dimensions
- **Weight Limits** - Material handling limits
- **Precision Level** - Tolerance capabilities

## Maintenance & Quality
- **Last Maintenance** - When last serviced
- **Next Maintenance Due** - Scheduled maintenance
- **Maintenance History** - Track of issues/repairs
- **Quality Metrics** - Defect rates, first-pass yield
- **Calibration Status** - For precision equipment

## Integration & Communication
- **PLC/Controller Type** - Hardware details
- **Communication Protocol** - How to connect (OPC-UA, MQTT, etc.)
- **Data Points** - Available sensors/parameters
- **API Endpoints** - For direct integration

## Business Logic
- **Cost Per Hour** - Operating cost
- **Department/Owner** - Who manages it
- **Shift Availability** - When it operates
- **Priority Level** - For scheduling