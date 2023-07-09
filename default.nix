{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-xMrJhsTo2QZhFOJRsboquIp9J5zv4gF8HjQf1FTXOPc=";
  ldflags = [ "-s" "-w" ];
}
