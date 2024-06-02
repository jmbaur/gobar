{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorHash = "sha256-V6CDBUQpYndV9NMQ5XzCgxewRNzg7CtwX5RuQlkmFk0=";
  ldflags = [
    "-s"
    "-w"
  ];
}
