{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.7";
  src = ./.;
  vendorSha256 = "sha256-cKLPEtHBiDhTPcy/g2kmgJWetFTLI4Po/15fo/c1lWE=";
  ldflags = [ "-s" "-w" ];
}
