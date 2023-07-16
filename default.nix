{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-VG2l7yB5M8tk9Mh9s01r2+lKPY+ASZ0DiltUTKYPHrU=";
  ldflags = [ "-s" "-w" ];
}
