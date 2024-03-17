{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorHash = "sha256-lzqHG0nSXpI2nyx7akWRj871eovoXqj6tRdFqgXcBSc=";
  ldflags = [ "-s" "-w" ];
}
