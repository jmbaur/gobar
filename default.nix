{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-B7G1Q5P+9R7MPcAf1UpXMSbjD9ou/hmKtFvT8X+epi0=";
  ldflags = [ "-s" "-w" ];
}
