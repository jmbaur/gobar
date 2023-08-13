{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-LcLhGEDlBM1Wr8Cdejpvu7StT506I+3nvjmO6RW1SeE=";
  ldflags = [ "-s" "-w" ];
}
