{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorHash = "sha256-/O5LNpQqXlRG9LC1B2pD3Qu0c4zOWW9IWeu7kEvQfpg=";
  ldflags = [
    "-s"
    "-w"
  ];
}
