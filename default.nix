{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.6";
  src = ./.;
  vendorSha256 = "sha256-UCkaznWjl37KkJ9/2SjzG4SB4LjG0NkvKc41XEcDQ0g=";
  ldflags = [ "-s" "-w" ];
}
