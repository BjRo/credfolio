import { build } from "esbuild";
import { fileURLToPath } from "url";
import { dirname, resolve } from "path";
import { existsSync, mkdirSync } from "fs";

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const appDir = resolve(__dirname, "..");
const src = resolve(appDir, "next.config.ts");
const out = resolve(appDir, "next.config.mjs");

if (!existsSync(appDir)) {
	mkdirSync(appDir, { recursive: true });
}

await build({
	entryPoints: [src],
	outfile: out,
	bundle: true,
	format: "esm",
	platform: "node",
	target: "node20",
	logLevel: "silent",
});

console.log("[prepare-next-config] generated next.config.mjs from next.config.ts");

