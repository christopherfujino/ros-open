// useCallback is a helper for memoization
import { useContext, useState } from 'https://esm.sh/preact/hooks';
import { createContext, render } from 'https://esm.sh/preact';
import { html } from 'https://esm.sh/htm@3.1.1/preact';

const Route = createContext("");

function error(msg) {
  // TODO have this update UI
  console.error(msg);
}

function Browser(_) {
  const [files, updateFiles] = useState(null);
  if (files === null) {
    fetch("/api/notes", {
      headers: { "Content-Type": "application/json" },
    }).then(function(res) {
      if (!res.ok) {
        // TODO surface this in UI
        error(`Bad response: ${res}`);
      } else {
        res.json().then((res) => {
          if (res.files == null) {
            console.error("API should not return NULL files!")
          } else {
            console.log("Got it!", res.files);
            updateFiles(res.files);
          }
        });
      }
    });

    return html`<p>Loading...</p>`;
  } else {
    return html`
    <table>
      ${files.map((file) => html`<tr><td>${file}</td></tr>`)}
    </table>`;
  }
}

function Body(props) {
  const [route, _] = useContext(Route);

  switch (route) {
    case "/editor":
      const [noteContents, updateNoteContents] = useState(props.note);
      const [path, updatePath] = useState(props.path);

      function submit() {
        fetch("/api/notes/update", {
          method: "UPDATE",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ path: path, contents: noteContents }),
        }).then(function(res) {
          if (!res.ok) {
            // TODO surface this in UI
            console.error("Bad response:", res);
          }
        });
      }
      return html`
        <form spellcheck=${false}>
          <label>
            Path
            <input value=${path} onInput=${(e) => updatePath(e.target.value)} />
          </label>
          <label>
            Contents
            <textarea name="note-entry" onInput=${(e) => updateNoteContents(e.target.value)}>${noteContents}</input>
          </label>
          <input type="button" value="Save" onClick=${submit} />
        </form>`;
    case "/browser":
      return html`<${Browser} />`;
    default:
      return html`<h1>Error invalid route <code>${route}</code></h1>`;
  }
}

function Nav() {
  const [route, setRoute] = useContext(Route);
  let editorElement, browserElement;
  switch (route) {
    case "/editor":
      editorElement = html`<a class="secondary">Editor</a>`;
      browserElement = html`<a href="#" onclick=${() => setRoute("/browser")}>Browser</a>`;
      break;
    case "/browser":
      editorElement = html`<a href="#" onclick=${() => setRoute("/editor")}>Editor</a>`;
      browserElement = html`<a class="secondary">Browser</a>`;
      break;
    default:
      return html`<h1>Error invalid route <code>${route}</code></h1>`;
  }
  return html`
    <nav aria-label="breadcrumb">
      <ul>
        <li><a href="/">ros</a></li>
        <li>${editorElement}</li>
        <li>${browserElement}</li>
      </ul>
    </nav>`;
}

// Create your app
function App(_) {
  const [route, setRoute] = useState("/editor");
  const routeTuple = [route, setRoute];

  return html`
    <${Route.Provider} value=${routeTuple}>
    <header class="container-fluid">
      <${Nav} />
    </header>
    <main class="container-fluid">
      <${Body} note="" path="/foo/bar" />
    </main>
    </${Route.Provider}>
  `;
}

addEventListener("DOMContentLoaded", (_) =>
  render(html`<${App} />`, document.getElementById("root")));
