{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-QAxZphIMQvyXzZGEo6+8CySvLqXPZCljy7BVeW2V3uA=";
  ldflags = [ "-s" "-w" ];
}
