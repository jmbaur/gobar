{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.6";
  src = ./.;
  vendorSha256 = "sha256-xCF6+Z/gFVKo3GTYNx1sf9DfCssFVwaN49wXdujTzgc=";
  ldflags = [ "-s" "-w" ];
}
