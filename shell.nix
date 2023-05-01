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
		vendorSha256 = useForVersion src.sqlc "v1.18.0"
			"0pppq4frcavvsllwrl3z8cpn1j23nkiqzh5h4142ajqrw83qydw0";
		doCheck = false;
		proxyVendor = true;
	};

	nixos-shell = pkgs.buildGoModule {
		name = "nixos-shell";
		src = src.nixos-shell;
		version = shortRev src.nixos-shell.rev;
		vendorSha256 = useForRev src.nixos-shell "e238cb5"
			"0gjj1zn29vyx704y91g77zrs770y2rakksnn9dhg8r6na94njh5a";
	};
in

pkgs.mkShell {
	buildInputs = with pkgs; [
		go
		gopls
		gotools
		niv
		goapi-gen
		sqlc
		pgformatter
		nixos-shell
	];
}