import { Editor } from "https://cdn.jsdelivr.net/npm/tiny-markdown-editor@0.1.5/+esm";

new Editor({
  textarea: "editor",
  content: `# Hello World

This is a test post. The title will be the first h1, if there is one.
You can also drag and drop images into the editor to upload them!`,
});
