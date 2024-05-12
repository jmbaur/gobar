{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorHash = "sha256-JXseG/O9kXnBQ3tWjPvrX4L7XRjzuGDTQs87nPgsS9k=";
  ldflags = [
    "-s"
    "-w"
  ];
}
