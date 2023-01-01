{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.6";
  src = ./.;
  vendorSha256 = "sha256-R6uO67RFehQ2n1eXRPSATNaldr4sIO8X/wL412uZ0Dg=";
  ldflags = [ "-s" "-w" ];
}
