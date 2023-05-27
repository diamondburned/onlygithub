{ pkgs ? import <nixpkgs> {} }:

let
	lib = pkgs.lib;
	src = import ./nix/sources.nix {};

	fetchPatchFromGitHub = { owner, repo, rev, sha256 }:
		pkgs.fetchpatch {
			url = "https://github.com/${owner}/${repo}/commit/${rev}.patch";
			inherit sha256;
		};

	useForVersion = src: version: hash:
		if (src.version == version) then hash else
			builtins.trace "${hash}: no value for newer version ${src.version}" lib.fakeSha256;
	useForRev = src: rev: hash:
		if (lib.hasPrefix rev src.rev) then hash else
			builtins.trace "${hash}: no value for newer revision ${src.rev}" lib.fakeSha256;
	shortRev = src: "${src.branch}-${builtins.substring 0 7 src.rev}";
in

let
	pkgs = import (src.nixpkgs) {};

	fetchPatchFromGitHub = { owner, repo, rev ? null, pr ? null, sha256 }:
		assert lib.assertMsg (rev != null || pr != null) "rev or pr must be set";
		pkgs.fetchpatch {
			url =
				if rev != null then
					"https://github.com/${owner}/${repo}/commit/${rev}.patch"
				else
					"https://github.com/${owner}/${repo}/pull/${pr}.patch";
			inherit sha256;
		};

	goapi-gen = pkgs.buildGoModule {
		name = "goapi-gen";
		src = src.goapi-gen;
		version = shortRev src.goapi-gen;
		vendorSha256 = useForRev src.goapi-gen "8ca3a5b"
			"1dknfg3w97421c8dnld5kvx0psicvmxr7wzkhqipaxplcg3cqrr9";
	};

	sqlc = pkgs.buildGoModule {
		name = "sqlc";
		src = src.sqlc;
		version = shortRev src.sqlc;
		vendorSha256 = useForRev src.sqlc "2adb18f"
			"sha256-CoJokasqaCK4lKN9A65JU2T00cWR0YclrowNDoe+q1c=";
		doCheck = false;
		proxyVendor = true;
		subPackages = [ "cmd/sqlc" ];
	};

	nixos-shell = pkgs.buildGoModule {
		name = "nixos-shell";
		src = src.nixos-shell;
		version = shortRev src.nixos-shell;
		vendorSha256 = useForRev src.nixos-shell "e238cb5"
			"0gjj1zn29vyx704y91g77zrs770y2rakksnn9dhg8r6na94njh5a";
	};

	genqlient = pkgs.buildGoModule {
		name = "genqlient";
		src = src.genqlient;
		version = shortRev src.genqlient;
		vendorSha256 = useForRev src.genqlient "677fa94"
			"150kwgywpivkc7q901ygdjjw8fgncwgcmkjj4lbrvkik7ynpm9dn";
		subPackages = [ "." ];
	};

	templ = pkgs.buildGoModule {
		name = "templ";
		src = src.templ;
		version = shortRev src.templ;
		vendorSha256 = useForRev src.templ "acea959"
			"sha256-GE471JVtlIOpH3hun/1Kae16/MDmUYV/w8dgV9sagB4=";
		subPackages = [ "cmd/templ" ];
	};

	dart-sass = import src.nix-dart-sass {
		inherit pkgs;
		sha256 = "0vdqcqkdbk1n71lbjkmravpw43h8lxc8dgk6sanlscnm98nlgc01";
		version = "1.62.1";
	};
in

pkgs.mkShell {
	buildInputs = with pkgs; [
		go
		gopls
		gotools
		air
		niv
		goapi-gen
		templ
		sqlc
		deno
		esbuild
		dart-sass
		genqlient
		pgformatter
		nixos-shell
		nodePackages.prettier # for scss mostly
	];

	shellHook = ''
		export DATABASE_URL="sqlite://./onlygithub.dev.db"
	'';
}
