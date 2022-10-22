{ buildGoModule, go-tools, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.3";
  CGO_ENABLED = 0;
  src = ./.;
  vendorSha256 = "sha256-UhquUYw+45anj8CEKWYVIb42Gk1j3hQtUDGZRcU+2zI=";
  preCheck = "HOME=/tmp ${go-tools}/bin/staticcheck ./...";
  ldflags = [ "-s" "-w" ];
}
