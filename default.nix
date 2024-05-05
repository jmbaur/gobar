{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorHash = "sha256-ig8eP9Sq8DGgX0tRM+J4z3Z9JsTq/yBaGnd2rQepyEU=";
  ldflags = [
    "-s"
    "-w"
  ];
}
