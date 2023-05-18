#!/usr/bin/env -S deno run -A
import * as esbuild from "https://deno.land/x/esbuild@v0.17.18/mod.js";
import { denoPlugins } from "https://deno.land/x/esbuild_deno_loader@0.7.0/mod.ts";

const [entrypoint, outfile] = Deno.args;
if (!entrypoint || !outfile) {
  throw new Error("Usage: bundle.js <entrypoint> <outfile>");
}

await esbuild.build({
  plugins: [...denoPlugins()],
  entryPoints: [entrypoint],
  outfile,
  format: "esm",
  bundle: true,
  minify: true,
  sourcemap: true,
});
esbuild.stop();
