{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorHash = "sha256-QxuC38b24T1F+ARzBinvPP8Lwu+lwnyMviTzffuM2iQ=";
  ldflags = [ "-s" "-w" ];
}
