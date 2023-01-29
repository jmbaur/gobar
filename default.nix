{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.6";
  src = ./.;
  vendorSha256 = "sha256-CDXQmy+kp/VtMt6nkce08nUd/6FFxuA0jWo9SkXO2ig=";
  ldflags = [ "-s" "-w" ];
}
