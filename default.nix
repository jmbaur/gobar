{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorHash = "sha256-O5Zjo4LYyIZUPXalrAH9it7GOPKwhkeI3h08EfJI0wk=";
  ldflags = [
    "-s"
    "-w"
  ];
}
