{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorHash = "sha256-YOoKqWH0gv8TAJIovRgTKNg/LEx7V4zQIQbIESefbd8=";
  ldflags = [ "-s" "-w" ];
}
