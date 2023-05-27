import TinyMDE from "https://cdn.jsdelivr.net/npm/tiny-markdown-editor@0.1.5/+esm";
import "https://cdn.jsdelivr.net/npm/highlighted-code@0.3.7/+esm";
import { micromark } from "https://esm.sh/micromark@3";

function growTextarea() {
  document
    .querySelectorAll("textarea.grow")
    .forEach((textarea: HTMLTextAreaElement) => {
      function resize() {
        textarea.style.height = "0";
        textarea.style.height = textarea.scrollHeight + "px";
      }

      textarea.style.resize = "none";
      textarea.addEventListener("input", () => resize());
      resize();
    });
}

function enableMarkdownEditors() {
  const textareas = document.querySelectorAll("textarea.markdown-editor");
  textareas.forEach((textarea: HTMLTextAreaElement) => {
    new TinyMDE.Editor({
      element: textarea,
      content: textarea.value,
    });
  });
}

function renderMarkdown() {
  document.querySelectorAll("div.markdown").forEach((div) => {
    div.innerHTML = micromark(div.innerHTML, "utf8", {
      allowDangerousHtml: div.classList.contains("markdown-unsafe"),
    });
  });
}

const createRelativeTimeFormatter = (style) =>
  new Intl.RelativeTimeFormat(undefined, {
    style,
    numeric: "auto",
  });

const relativeTimes = {
  short: createRelativeTimeFormatter("short"),
  long: createRelativeTimeFormatter("long"),
};

const Duration = {
  second: 1000,
  minute: 60 * 1000,
  hour: 60 * 60 * 1000,
  day: 24 * 60 * 60 * 1000,
};

function updateTimes() {
  document.querySelectorAll("time.relative").forEach((time) => {
    let formatter = relativeTimes.long;
    if (time.classList.contains("short")) {
      formatter = relativeTimes.short;
    }

    const date = new Date(time.getAttribute("datetime")!);

    let delta = date.getTime() - Date.now();
    let unit = "day";
    if (Math.abs(delta) < Duration.minute) {
      delta = Math.round(delta / Duration.second);
      unit = "second";
    } else if (Math.abs(delta) < Duration.hour) {
      delta = Math.round(delta / Duration.minute);
      unit = "minute";
    } else if (Math.abs(delta) < Duration.day) {
      delta = Math.round(delta / Duration.hour);
      unit = "hour";
    }

    time.textContent = formatter.format(delta, unit as any);
  });

  document.querySelectorAll("time.localize").forEach((time) => {
    const date = new Date(time.getAttribute("datetime")!);
    time.textContent = date.toLocaleString(undefined, {
      dateStyle: "full",
      timeStyle: "long",
    });
  });
}

export default function enhance() {
  growTextarea();
  renderMarkdown();
  enableMarkdownEditors();
  updateTimes();
}
