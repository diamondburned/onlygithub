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
			throw "${hash}: no value for newer version ${version}";
	useForRev = src: rev: hash:
		if (lib.hasPrefix rev src.rev) then hash else
			throw "${hash}: no value for newer revision ${rev}";
	shortRev = rev: builtins.substring 0 7 rev;
in

let
	pkgs = import (src.nixpkgs) {};

	goapi-gen = pkgs.buildGoModule {
		name = "goapi-gen";
		src = src.goapi-gen;
		version = shortRev src.goapi-gen.rev;
		vendorSha256 = useForRev src.goapi-gen "8ca3a5b"
			"1dknfg3w97421c8dnld5kvx0psicvmxr7wzkhqipaxplcg3cqrr9";
	};

	sqlc = pkgs.buildGoModule {
		name = "sqlc";
		src = src.sqlc;
		version = src.sqlc.version;
		# vendorSha256 = useForVersion src.sqlc "v1.18.0"
		# 	"0pppq4frcavvsllwrl3z8cpn1j23nkiqzh5h4142ajqrw83qydw0";
		vendorSha256 = useForVersion src.sqlc "v1.17.2"
			"0ih9siizn6nkvm4wi0474wxy323hklkhmdl52pql0qzqanmri4yb";
		doCheck = false;
		proxyVendor = true;
		subPackages = [ "cmd/sqlc" ];
	};

	nixos-shell = pkgs.buildGoModule {
		name = "nixos-shell";
		src = src.nixos-shell;
		version = shortRev src.nixos-shell.rev;
		vendorSha256 = useForRev src.nixos-shell "e238cb5"
			"0gjj1zn29vyx704y91g77zrs770y2rakksnn9dhg8r6na94njh5a";
	};

	genqlient = pkgs.buildGoModule {
		name = "genqlient";
		src = src.genqlient;
		version = shortRev src.genqlient.rev;
		vendorSha256 = useForRev src.genqlient "677fa94"
			"150kwgywpivkc7q901ygdjjw8fgncwgcmkjj4lbrvkik7ynpm9dn";
		subPackages = [ "." ];
	};

	templ = pkgs.buildGoModule {
		name = "templ";
		src = src.templ;
		version = shortRev src.templ.rev;
		vendorSha256 = useForRev src.templ "e1ca5e2"
			"07l03bdmfq67qdzqalg6q3y7mvb99byryvkq3ylq753djpa3nkhq";
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
