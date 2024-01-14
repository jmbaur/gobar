{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorHash = "sha256-k5YELJc49y+I3UDCG7zo6bAvIoL+GKCUxiQFVfJaFro=";
  ldflags = [ "-s" "-w" ];
}
