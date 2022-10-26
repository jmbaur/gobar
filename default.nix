{ buildGoModule, CGO_ENABLED ? 0, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.3";
  src = ./.;
  vendorSha256 = "sha256-UhquUYw+45anj8CEKWYVIb42Gk1j3hQtUDGZRcU+2zI=";
  inherit CGO_ENABLED;
  ldflags = [ "-s" "-w" ];
}
