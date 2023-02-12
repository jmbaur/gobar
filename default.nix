{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.7";
  src = ./.;
  vendorSha256 = "sha256-ijk0kn++2BxzyNt1qkpyQTWDBKQTZTHq8UPBJBTeg2w=";
  ldflags = [ "-s" "-w" ];
}
