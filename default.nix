{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-tVu2auU1qfXDazMd/Id8PcIYlzRfQKXWPUICUuycamA=";
  ldflags = [ "-s" "-w" ];
}
