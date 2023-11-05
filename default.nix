{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-ocPUcFILwlesN4EX/JTq1TfdfT5e8ovhEsmqdNsRFfs=";
  ldflags = [ "-s" "-w" ];
}
