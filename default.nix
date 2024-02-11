{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorHash = "sha256-RbuL8P/IkFw07kcZUlii34FYzgF9U8tW8PhfQgWX4Zk=";
  ldflags = [ "-s" "-w" ];
}
