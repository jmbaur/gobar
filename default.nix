{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.6";
  src = ./.;
  vendorSha256 = "sha256-kJlmi+DxOuC1GuJ1XOrdMIvjW3jYD4OzuNZ66EXibRA=";
  ldflags = [ "-s" "-w" ];
}
