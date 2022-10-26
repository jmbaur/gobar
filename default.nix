{ buildGoModule, CGO_ENABLED ? 0, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.3";
  src = ./.;
  vendorSha256 = "sha256-4N5pm67xebRVNmwVFUgevhgWHlhdxipeY5f6CX0MKfg=";
  inherit CGO_ENABLED;
  ldflags = [ "-s" "-w" ];
}
