{
  description = "protoc Elm plugin";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-22.05-small";

  outputs = { self, nixpkgs }:
    let pkgs = import nixpkgs { system = "x86_64-linux"; };
    in {
      devShell.x86_64-linux = with pkgs;
        mkShell {
          nativeBuildInputs = [
            go_1_18
            go-protobuf
            elmPackages.elm
            elmPackages.elm-format
            elmPackages.elm-live
            elmPackages.elm-test
            protobuf
            protoc-gen-twirp
          ];
        };
    };
}
