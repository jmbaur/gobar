{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorHash = "sha256-F55gV4+QziHV4Gqg/7QJXSeOsUAVedcTYzj0KnUEjis=";
  ldflags = [ "-s" "-w" ];
}
