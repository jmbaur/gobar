{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-AIW/88+iM0MFYmsep4TfXPH4dcnRNMBDM1rxCofKwKg=";
  ldflags = [ "-s" "-w" ];
}
