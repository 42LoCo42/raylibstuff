{
  description = "Playing around with raylib";

  inputs.flake-utils.url = "github:numtide/flake-utils";

  inputs.nixgl.url = "github:guibou/nixGL";
  inputs.nixgl.inputs.nixpkgs.follows = "nixpkgs";
  inputs.nixgl.inputs.flake-utils.follows = "flake-utils";

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true;
          overlays = [ self.inputs.nixgl.overlay ];
        };

        vend = pkgs.buildGoModule rec {
          pname = "vend";
          version = "1.0.3";

          src = pkgs.fetchgit {
            url = "https://github.com/nomad-software/vend";
            rev = "v${version}";
            hash = "sha256-7AdE5qps4OMjaubt9Af6ATaqrV3n73ZuI7zTz7Kgm6w=";
          };

          vendorHash = null;
        };

        buildInputs = with pkgs; [
          libglvnd
          xorg.libX11
          xorg.libXcursor
          xorg.libXi
          xorg.libXinerama
          xorg.libXrandr
        ];

        nativeBuildInputs = with pkgs; [
          makeWrapper
          pkg-config
        ];

        devTools = with pkgs; [
          bashInteractive
          go
          gopls
          pkgs.nixgl.auto.nixGLDefault
          vend
        ];
      in
      {
        defaultPackage = pkgs.buildGoModule {
          pname = "raylibstuff";
          version = "0.0.1";
          src = ./.;

          vendorHash = null;
          inherit buildInputs nativeBuildInputs;
        };

        devShell = pkgs.mkShell {
          packages = devTools ++ buildInputs ++ nativeBuildInputs;
        };
      });
}
