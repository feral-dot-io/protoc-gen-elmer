{
  description = "handles incoming sflow packets";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-21.11-small";

  outputs = { self, nixpkgs }:
    let pkgs = import nixpkgs { system = "x86_64-linux"; };
    in {
      devShell.x86_64-linux = with pkgs;
        mkShell { nativeBuildInputs = [ go protobuf ]; };
    };
}
