{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-Uo9fzF/qpbTD/jAjnY5BqWx1MToH3VQ5J1WHiLwnsnw=";
  ldflags = [ "-s" "-w" ];
}
