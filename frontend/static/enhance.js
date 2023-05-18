import TinyMDE from "https://cdn.jsdelivr.net/npm/tiny-markdown-editor@0.1.5/+esm";

function growTextarea() {
  document.querySelectorAll("textarea.grow").forEach((textarea) => {
    function resize() {
      textarea.style.height = 0;
      textarea.style.height = textarea.scrollHeight + "px";
    }

    textarea.style = "resize: none;";
    textarea.addEventListener("input", () => resize());
    resize();
  });
}

function enableMarkdownEditors() {
  const textareas = document.querySelectorAll("textarea.markdown-editor");
  textareas.forEach((textarea) => {
    new TinyMDE.Editor({
      element: textarea,
      content: textarea.value,
    });
  });
}

export default function enhance() {
  growTextarea();
  enableMarkdownEditors();
}
