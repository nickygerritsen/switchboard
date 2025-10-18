# Switchboard Installation Script for Windows
# Requires PowerShell 5.1 or later

param(
    [string]$Version = "",
    [switch]$Register = $false
)

$ErrorActionPreference = "Stop"

# Configuration
$Repo = "nickygerritsen/switchboard"
$BinaryName = "switchboard.exe"
$InstallDir = "$env:LOCALAPPDATA\switchboard"

# Print colored output
function Write-Info {
    param([string]$Message)
    Write-Host "==> " -ForegroundColor Green -NoNewline
    Write-Host $Message
}

function Write-Error {
    param([string]$Message)
    Write-Host "Error: " -ForegroundColor Red -NoNewline
    Write-Host $Message
}

function Write-Warning {
    param([string]$Message)
    Write-Host "Warning: " -ForegroundColor Yellow -NoNewline
    Write-Host $Message
}

# Detect architecture
function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    switch ($arch) {
        "AMD64" { return "x86_64" }
        "ARM64" { return "arm64" }
        default {
            Write-Error "Unsupported architecture: $arch"
            Write-Error "Switchboard supports x86_64 and arm64 only."
            exit 1
        }
    }
}

# Get latest version from GitHub
function Get-LatestVersion {
    Write-Info "Fetching latest version..."

    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
        return $response.tag_name
    }
    catch {
        Write-Error "Failed to fetch latest version: $_"
        exit 1
    }
}

# Download file
function Download-File {
    param(
        [string]$Url,
        [string]$Output
    )

    try {
        Invoke-WebRequest -Uri $Url -OutFile $Output -UseBasicParsing
    }
    catch {
        Write-Error "Failed to download from $Url : $_"
        exit 1
    }
}

# Verify checksum
function Test-Checksum {
    param(
        [string]$FilePath,
        [string]$ChecksumsFile
    )

    Write-Info "Verifying checksum..."

    # Read checksums file
    $checksums = Get-Content $ChecksumsFile
    $fileName = Split-Path -Leaf $FilePath

    # Find expected checksum
    $checksumLine = $checksums | Where-Object { $_ -match [regex]::Escape($fileName) }

    if (-not $checksumLine) {
        Write-Warning "Could not find checksum for $fileName"
        return $false
    }

    $expectedChecksum = ($checksumLine -split '\s+')[0]

    # Calculate actual checksum
    $actualChecksum = (Get-FileHash -Path $FilePath -Algorithm SHA256).Hash.ToLower()

    if ($expectedChecksum -eq $actualChecksum) {
        Write-Info "Checksum verified successfully"
        return $true
    }
    else {
        Write-Error "Checksum verification failed!"
        Write-Error "Expected: $expectedChecksum"
        Write-Error "Got:      $actualChecksum"
        return $false
    }
}

# Check if directory is in PATH
function Test-InPath {
    param([string]$Directory)

    $pathDirs = $env:PATH -split ';'
    return $pathDirs -contains $Directory
}

# Add directory to PATH
function Add-ToPath {
    param([string]$Directory)

    Write-Info "Adding $Directory to user PATH..."

    # Get current user PATH
    $userPath = [Environment]::GetEnvironmentVariable("PATH", "User")

    # Add directory if not already present
    if ($userPath -notlike "*$Directory*") {
        $newPath = "$userPath;$Directory"
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")

        # Update current session PATH
        $env:PATH = "$env:PATH;$Directory"

        Write-Info "Added to PATH. You may need to restart your terminal."
        return $true
    }

    return $false
}

# Main installation
function Main {
    Write-Info "Installing Switchboard..."

    # Detect system
    $arch = Get-Architecture
    $os = "Windows"

    Write-Info "Detected OS: $os"
    Write-Info "Detected Architecture: $arch"

    # Get version if not specified
    if (-not $Version) {
        $Version = Get-LatestVersion
    }

    Write-Info "Installing version: $Version"

    # Build download URLs
    $versionNoV = $Version -replace '^v', ''
    $archiveName = "${BinaryName.Replace('.exe', '')}_${versionNoV}_${os}_${arch}.zip"
    $downloadUrl = "https://github.com/$Repo/releases/download/$Version/$archiveName"
    $checksumsUrl = "https://github.com/$Repo/releases/download/$Version/checksums.txt"

    # Create temporary directory
    $tmpDir = New-Item -ItemType Directory -Path (Join-Path $env:TEMP ([System.IO.Path]::GetRandomFileName()))

    try {
        # Download archive and checksums
        Write-Info "Downloading $archiveName..."
        $archivePath = Join-Path $tmpDir $archiveName
        Download-File -Url $downloadUrl -Output $archivePath

        Write-Info "Downloading checksums..."
        $checksumsPath = Join-Path $tmpDir "checksums.txt"
        Download-File -Url $checksumsUrl -Output $checksumsPath

        # Verify checksum
        if (-not (Test-Checksum -FilePath $archivePath -ChecksumsFile $checksumsPath)) {
            Write-Error "Installation aborted due to checksum verification failure"
            exit 1
        }

        # Extract archive
        Write-Info "Extracting archive..."
        Expand-Archive -Path $archivePath -DestinationPath $tmpDir -Force

        # Create installation directory
        if (-not (Test-Path $InstallDir)) {
            New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        }

        # Install binary
        Write-Info "Installing to $InstallDir..."
        $sourcePath = Join-Path $tmpDir $BinaryName
        $targetPath = Join-Path $InstallDir $BinaryName

        Copy-Item -Path $sourcePath -Destination $targetPath -Force

        # Add to PATH if not already
        if (-not (Test-InPath -Directory $InstallDir)) {
            Add-ToPath -Directory $InstallDir
        }

        # Verify installation
        Write-Info "Installation successful!"
        Write-Info "Installed: $targetPath"

        # Show version
        try {
            $installedVersion = & $targetPath --version 2>&1 | Select-Object -First 1
            Write-Info "Version: $installedVersion"
        }
        catch {
            Write-Warning "Could not verify installed version"
        }

        # Optionally register as default browser
        if ($Register) {
            Write-Info "Registering Switchboard as a browser..."
            try {
                & $targetPath register
                Write-Info "Registration successful!"
                Write-Info "You can now set Switchboard as your default browser in Windows Settings"
            }
            catch {
                Write-Warning "Registration failed. You can register manually later with: switchboard register"
            }
        }
        else {
            Write-Host ""
            Write-Info "To register Switchboard as a browser, run: switchboard register"
            Write-Info "Then set it as your default browser in Windows Settings"
        }

        Write-Host ""
        Write-Info "Installation complete!"
    }
    finally {
        # Cleanup
        Remove-Item -Path $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
    }
}

# Run main function
Main
