{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-9MoWBxL6qIZ4DFvpdwihktr9i3dS/L1hH+9W82cuM4w=";
  ldflags = [ "-s" "-w" ];
}
