{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.6";
  src = ./.;
  vendorSha256 = "sha256-6o6fNg6rUvRwR9YPJ8Eqkg4zu4KHJLBfg2+clopZKk8=";
  ldflags = [ "-s" "-w" ];
}
