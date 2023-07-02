{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-C/6QUPurX4prRAv4NspjS0obh+UOe024/6R0GF9uSPM=";
  ldflags = [ "-s" "-w" ];
}
