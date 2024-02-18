{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorHash = "sha256-S1zxgZYaq52khFl0RB7nA7SRgtBslOo9IWVeOq4EHTU=";
  ldflags = [ "-s" "-w" ];
}
