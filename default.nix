{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-y6SROP4Rx46Ndyjrw90uBR+VlQRsX9wyYqW+mFR4ry0=";
  ldflags = [ "-s" "-w" ];
}
