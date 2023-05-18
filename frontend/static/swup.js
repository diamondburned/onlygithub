import Swup from "https://cdn.jsdelivr.net/npm/swup@3.0/+esm";
import Swupscroll from "https://cdn.jsdelivr.net/npm/@swup/scroll-plugin@2.0/+esm";
import enhance from "./enhance.js";

const swup = new Swup({
  containers: ["main"],
  animationSelector: "main",
  plugins: [new Swupscroll()],
});

let previousScripts = [];
let scriptCounter = 0;

function init() {
  enhance();

  previousScripts.forEach((script) => script.remove());

  const main = document.querySelector("main");
  const scriptMetas = [...main.querySelectorAll(`meta[name="swup-script"]`)];

  previousScripts = scriptMetas.map((meta) => {
    const script = document.createElement("script");
    script.type = "module";
    // Hack to force the browser to refetch the script.
    script.src = meta.content + `?${scriptCounter++}`;
    script.async = true;
    script.onload = () =>
      console.debug(`Loaded script ${meta.content} for page ${main.id}`);
    script.onerror = (err) =>
      console.log(`Error loading script for page ${main.id}`, err);
    document.body.appendChild(script);
    return script;
  });
}

swup.on("contentReplaced", () => init());
init();
