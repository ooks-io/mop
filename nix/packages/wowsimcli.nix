{
  lib,
  buildGoModule,
  go,
  protobuf,
  protoc-gen-go,
  self,
  ...
}:
buildGoModule rec {
  pname = "wowsims-mop";
  version = "0.22.0";
  src = "${self}";

  vendorHash = "sha256-uv0klkfzFnMP7rZJkujhJOnogooOvgL/BCpxRM0KUL8=";

  # Only build the CLI package, not the web packages
  subPackages = ["cmd/wowsimcli"];

  nativeBuildInputs = [
    go
    protobuf
    protoc-gen-go
  ];

  postPatch = ''
    # Generate protocol buffers before vendoring
    protoc -I=./proto --go_out=./sim/core ./proto/*.proto

    # Temporarily remove web packages that cause vendoring issues
    rm -rf sim/web
  '';

  preBuild = ''
    # Set home directory for build
    export HOME=$TMPDIR

    # Set environment variables
    export VERSION=0.22.0
  '';

  ldflags = [
    "-X 'main.Version=0.22.0'"
    "-s"
    "-w"
  ];

  tags = ["with_db"];

  meta = with lib; {
    description = "WoW Mists of Pandaria Classic DPS Simulator";
    homepage = "https://github.com/wowsims/mop";
    license = licenses.mit;
    maintainers = [];
    platforms = platforms.linux ++ platforms.darwin;
  };
}
