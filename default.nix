{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-60zKMAlc6wvz/c8xSzCH8l+zh1v8sdiM8lTCRMp7RCY=";
  ldflags = [ "-s" "-w" ];
}
