{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-cTiT3aAPBJTCV3LqhUo9rCeIFwYEKd2TUERGVhQX04c=";
  ldflags = [ "-s" "-w" ];
}
