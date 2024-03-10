{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorHash = "sha256-DThgbKF0788KzcdIOLrIfOfiPjw2nFJ0/G4urd3bkI4=";
  ldflags = [ "-s" "-w" ];
}
