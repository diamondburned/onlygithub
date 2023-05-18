import Swup from "https://cdn.jsdelivr.net/npm/swup@3.0/+esm";
import Swupscroll from "https://cdn.jsdelivr.net/npm/@swup/scroll-plugin@2.0/+esm";
import enhance from "./enhance.ts";

const swup = new Swup({
  containers: ["main"],
  animationSelector: "main",
  plugins: [new Swupscroll()],
});

swup.on("contentReplaced", () => enhance());
enhance();
