{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorHash = "sha256-fYypaBKogkaUUs/uTKW4iaTUIFQX3N+LvyyXHgO8MSs=";
  ldflags = [ "-s" "-w" ];
}
