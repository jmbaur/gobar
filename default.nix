{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorHash = "sha256-cdyjiENv5vlTODBrtPXamN4ne/TxgsXkEg7WWaYxWkk=";
  ldflags = [ "-s" "-w" ];
}
