{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-Z5HMLtxfDGKU0mJ0jw89tADi6/1hnqbKLUgET0ofqFw=";
  ldflags = [ "-s" "-w" ];
}
