{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.6";
  src = ./.;
  vendorSha256 = "sha256-MIu6rTYJH9X5nt3tQUhbaxgPshxGASu0qC7FN2jMOuw=";
  ldflags = [ "-s" "-w" ];
}
