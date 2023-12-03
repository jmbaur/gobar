{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorHash = "sha256-dnfYF9IpQ01MIYtMjaVKhvkS0V1kfSHbKg39H5MjYK4=";
  ldflags = [ "-s" "-w" ];
}
