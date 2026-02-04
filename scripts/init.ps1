# Initialize Go Microservice Template (PowerShell)
# Usage: .\scripts\init.ps1 -Module "github.com/your-org/service" -ServiceName "service"

param(
    [Parameter(Mandatory=$true)]
    [string]$Module,
    
    [Parameter(Mandatory=$true)]
    [string]$ServiceName
)

$ErrorActionPreference = "Stop"

$CurrentModule = "github.com/minisource/template_go"
$ProjectRoot = Split-Path -Parent $PSScriptRoot

Write-Host "Initializing Go Microservice Template" -ForegroundColor Green
Write-Host "  Module: $Module"
Write-Host "  Service: $ServiceName"
Write-Host ""

if ($Module -eq "github.com/your-org/your-service") {
    Write-Host "Error: Please provide your actual module name" -ForegroundColor Red
    exit 1
}

# Step 1: Replace module in go.mod
Write-Host "[1/5] Replacing module name in go.mod..." -ForegroundColor Yellow
$goModPath = Join-Path $ProjectRoot "src\go.mod"
(Get-Content $goModPath) -replace [regex]::Escape($CurrentModule), $Module | Set-Content $goModPath

# Step 2: Replace imports in all Go files
Write-Host "[2/5] Replacing imports in all Go files..." -ForegroundColor Yellow
Get-ChildItem -Path (Join-Path $ProjectRoot "src") -Filter "*.go" -Recurse | ForEach-Object {
    (Get-Content $_.FullName) -replace [regex]::Escape($CurrentModule), $Module | Set-Content $_.FullName
}

# Step 3: Update config files
Write-Host "[3/5] Updating configuration files..." -ForegroundColor Yellow
Get-ChildItem -Path (Join-Path $ProjectRoot "src\config") -Filter "*.yml" | ForEach-Object {
    (Get-Content $_.FullName) -replace "DiviPay", $ServiceName | Set-Content $_.FullName
}

# Step 4: Update Docker files
Write-Host "[4/5] Updating Docker files..." -ForegroundColor Yellow
$dockerfilePath = Join-Path $ProjectRoot "src\Dockerfile"
if (Test-Path $dockerfilePath) {
    (Get-Content $dockerfilePath) -replace "backend", $ServiceName | Set-Content $dockerfilePath
}

$composePath = Join-Path $ProjectRoot "docker\docker-compose.yml"
if (Test-Path $composePath) {
    (Get-Content $composePath) -replace "backend", $ServiceName | Set-Content $composePath
}

# Step 5: Run go mod tidy
Write-Host "[5/5] Running go mod tidy..." -ForegroundColor Yellow
Push-Location (Join-Path $ProjectRoot "src")
try {
    & go mod tidy
} finally {
    Pop-Location
}

Write-Host ""
Write-Host "Project initialized successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:"
Write-Host "  1. Update src\config\config-development.yml with your database settings"
Write-Host "  2. Update src\docs\docs.go with your API info"
Write-Host "  3. Run 'make run' or 'cd src; go run ./cmd/main.go' to start the server"
Write-Host ""
Write-Host "Useful commands:"
Write-Host "  make run        - Run the application"
Write-Host "  make test       - Run tests"
Write-Host "  make swagger    - Generate API docs"
Write-Host "  make docker-run - Run with Docker"
