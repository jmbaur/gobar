{ buildGoModule, CGO_ENABLED ? 0, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.4";
  src = ./.;
  vendorSha256 = "sha256-kusaplLYvse/Ea4OGjwJWhcf3eXutLYFC6gZSjOiwHw=";
  inherit CGO_ENABLED;
  ldflags = [ "-s" "-w" ];
}
