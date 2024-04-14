{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorHash = "sha256-i3EeYYR1cYnJKz7rKngrxPK0tm2JtQax4PYI4KHgUuk=";
  ldflags = [
    "-s"
    "-w"
  ];
}
