Name:           switchboard
Version:        __VERSION__
Release:        1%{?dist}
Summary:        Smart URL router that opens links in different browsers

License:        MIT
URL:            https://github.com/nickygerritsen/switchboard
Source0:        https://github.com/nickygerritsen/switchboard/archive/v__VERSION__.tar.gz

BuildRequires:  golang >= 1.24
BuildRequires:  git

%description
Switchboard routes URLs to different browsers based on configurable
patterns, with support for browser profiles and cross-platform usage.

%prep
%autosetup -n %{name}-__VERSION__

%build
export CGO_ENABLED=0
export GOFLAGS="-mod=readonly -trimpath"
go build -o switchboard -ldflags="-s -w -X main.version=__VERSION__" ./cmd/switchboard

%install
# Install binary
install -Dm755 switchboard %{buildroot}%{_bindir}/switchboard

# Install config example
install -Dm644 config.example.yaml %{buildroot}%{_docdir}/switchboard/config.example.yaml

# Install desktop file
install -Dm644 switchboard.desktop %{buildroot}%{_datadir}/applications/switchboard.desktop

# Install license
install -Dm644 LICENSE %{buildroot}%{_datadir}/licenses/switchboard/LICENSE

%files
%{_bindir}/switchboard
%{_docdir}/switchboard/config.example.yaml
%{_datadir}/applications/switchboard.desktop
%{_datadir}/licenses/switchboard/LICENSE

%changelog
* %(date "+%a %b %d %Y") Nicky Gerritsen <github@accounts.nickygerritsen.nl> - __VERSION__-1
- Release __VERSION__
