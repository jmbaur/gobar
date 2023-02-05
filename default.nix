{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.7";
  src = ./.;
  vendorSha256 = "sha256-wVcr4ev4hPVg4oJ5UGm8VeLHg7EHJsu/CnAtbz+/ZxU=";
  ldflags = [ "-s" "-w" ];
}
