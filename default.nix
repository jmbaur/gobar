{ buildGoModule, CGO_ENABLED ? 0, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.4";
  src = ./.;
  vendorSha256 = "sha256-4Q05zdNSdpC+dB2lBrfkycSUzGbruUoJ42yd/ZW8HVI=";
  inherit CGO_ENABLED;
  ldflags = [ "-s" "-w" ];
}
