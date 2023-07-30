{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-zunH0UI3V2TuZDS+Djtfb7oMQQGzD8dlKIzPXAv7/CI=";
  ldflags = [ "-s" "-w" ];
}
