{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorHash = "sha256-sBuOR6VxCyqkvuVU7SZOBIUkoQek7WjuXAx1++JxKyU=";
  ldflags = [
    "-s"
    "-w"
  ];
}
