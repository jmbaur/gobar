{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorHash = "sha256-/OaIJtg2v7W+0PAvAWN7dEYq3uA01+kTB/kE7zOo4T8=";
  ldflags = [ "-s" "-w" ];
}
