# 1. Set paths
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Definition
$rootPath = Split-Path -Parent $scriptPath
$envFile = Join-Path $rootPath ".env"
$migrationsDir = Join-Path $rootPath "migrations"

# 2. Load .env file
if (Test-Path $envFile) {
    Get-Content $envFile | Where-Object { $_ -match '=' -and $_ -notmatch '^#' } | ForEach-Object {
        $name, $value = $_.Split('=', 2)
        [System.Environment]::SetEnvironmentVariable($name.Trim(), $value.Trim())
    }
} else {
    Write-Error ".env file not found at $envFile"
    exit
}

# 3. Get variables from environment
$user = $env:DB_USER
$pass = $env:DB_PASSWORD
$host_addr = $env:DB_HOST
$port = $env:DB_PORT
$name = $env:DB_NAME

# 4. Construct DSN
$DSN = "postgres://$($user):$($pass)@$($host_addr):$($port)/$($name)?sslmode=disable"

# 5. Run Goose
# $args passes through commands like 'up', 'down', or 'status'
Write-Host "Running goose on $name..." -ForegroundColor Cyan
goose -dir "$migrationsDir" postgres "$DSN" $args