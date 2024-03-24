{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorHash = "sha256-tigiZ0/ypAHBLFbI9CYzGtN2OsKuRQpLwpVmWVr1kDQ=";
  ldflags = [ "-s" "-w" ];
}
