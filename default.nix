{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorHash = "sha256-pl5/M4i+OazHQdarV5WEg1Zfvm5nXr2qsarF9Fx9470=";
  ldflags = [ "-s" "-w" ];
}
