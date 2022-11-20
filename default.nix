{ buildGoModule, CGO_ENABLED ? 0, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.4";
  src = ./.;
  vendorSha256 = "sha256-vDWU3XNEHxBVwxuxuZxGp3m4qxk85e8/Ullw04C7ihI=";
  inherit CGO_ENABLED;
  ldflags = [ "-s" "-w" ];
}
