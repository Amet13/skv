# PowerShell build script for skv binary - supports local and cross-platform builds
param(
    [Parameter(Position=0)]
    [ValidateSet("host", "all")]
    [string]$Mode = "host"
)

$ErrorActionPreference = "Stop"

# Setup build environment
$RepoRoot = Split-Path -Parent $PSScriptRoot
$DistDir = Join-Path $RepoRoot "dist"
$PgoProfile = $env:PGO_PROFILE
$VersionOverride = $env:VERSION

if (!(Test-Path $DistDir)) {
    New-Item -ItemType Directory -Path $DistDir | Out-Null
}

function Build-Target {
    param(
        [string]$Goos,
        [string]$Goarch
    )

    Write-Host "Building skv for ${Goos}/${Goarch}..."

    # Build configuration
    $BuildFlags = @("-trimpath")
    $LdFlags = @("-s", "-w")

    # Version information
    $Version = if ($VersionOverride) {
        $VersionOverride
    } else {
        try {
            $GitVersion = git describe --tags --dirty --always 2>$null
            if ($LASTEXITCODE -eq 0) { $GitVersion } else { "dev" }
        } catch {
            "dev"
        }
    }

    # Git commit and build timestamp
    $Commit = try {
        git rev-parse --short HEAD 2>$null
        if ($LASTEXITCODE -eq 0) { $Commit } else { "" }
    } catch { "" }

    $Date = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")

    # Inject build-time variables
    $LdFlags += "-X skv/internal/version.Version=$Version"
    $LdFlags += "-X skv/internal/version.Commit=$Commit"
    $LdFlags += "-X skv/internal/version.Date=$Date"

    # Enable PGO if profile is available
    if ($PgoProfile -and (Test-Path $PgoProfile)) {
        $BuildFlags += "-pgo=$PgoProfile"
    }

    # Set environment variables and execute build
    $env:GOOS = $Goos
    $env:GOARCH = $Goarch

    $OutputName = "skv_${Goos}_${Goarch}"
    if ($Goos -eq "windows") {
        $OutputName += ".exe"
    }

    $OutputPath = Join-Path $DistDir $OutputName

    Push-Location $RepoRoot
    try {
        go build $BuildFlags -ldflags ($LdFlags -join " ") -o $OutputPath ./cmd/skv
    } finally {
        Pop-Location
    }
}

# Build execution based on mode
switch ($Mode) {
    "host" {
        # Build for current platform
        $CurrentGoos = go env GOOS
        $CurrentGoarch = go env GOARCH
        Build-Target -Goos $CurrentGoos -Goarch $CurrentGoarch
    }
    "all" {
        # Cross-platform build matrix
        Build-Target -Goos "darwin" -Goarch "arm64"
        Build-Target -Goos "darwin" -Goarch "amd64"
        Build-Target -Goos "linux" -Goarch "amd64"
        Build-Target -Goos "windows" -Goarch "amd64"
    }
}

# Display build artifacts
Write-Host "Artifacts in $DistDir:"
Get-ChildItem $DistDir -File | Sort-Object Name | ForEach-Object { $_.Name }
