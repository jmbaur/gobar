{ buildGoModule, CGO_ENABLED ? 0, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.3";
  src = ./.;
  vendorSha256 = "sha256-TFVPj04mw9YEvR4aHpmeK0GVHdgpEYAGx3DV7ikhTgk=";
  inherit CGO_ENABLED;
  ldflags = [ "-s" "-w" ];
}
