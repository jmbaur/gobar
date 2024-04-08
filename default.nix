{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorHash = "sha256-yrl+/Fb55ysgHWvNOrfETm4mRnK7e3JfcseFapoTTh4=";
  ldflags = [
    "-s"
    "-w"
  ];
}
