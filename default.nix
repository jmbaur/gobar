{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-xC02r6kMtP3M0e8ASXzMf9yCCEOlXKOz39G5xIE6zYw=";
  ldflags = [ "-s" "-w" ];
}
