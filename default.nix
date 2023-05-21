{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-EuHEBfeT+V0LSuNDSGYqS1llpK25qCY/HzRczYttTv4=";
  ldflags = [ "-s" "-w" ];
}
