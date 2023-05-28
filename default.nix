{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-C/aTc/gDgd41scYWl3axaXThPpe+hzW9KRWD7qa2ykk=";
  ldflags = [ "-s" "-w" ];
}
