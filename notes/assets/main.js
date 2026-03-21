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
  const [__, setRoute] = useContext(Route);
  if (files === null) {
    fetch("/api/notes/notes", {
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
            updateFiles(res.files);
          }
        });
      }
    });

    return html`<p>Loading...</p>`;
  } else {
    const rows = files.map((file) =>
      html`<tr><td><a href="#" onClick=${(_) => setRoute("/editor" + file)}>${file}</a></td></tr>`);
    return html`
    <table>
      ${rows}
    </table>`;
  }
}

function Editor(_) {
  const [route, __] = useContext(Route);
  const [noteContents, updateNoteContents] = useState(null);
  const remotePath = route.match(/^\/editor(.*)$/)[1]
  const [path, updatePath] = useState(remotePath);

  function submit() {
    // TODO update this
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
  if (noteContents === null) {
    fetch("/api/notes/note" + remotePath, {headers: {"Content-Type": "application/json"}}).then(function(res) {
      if (!res.ok) {
        error(`Bad response: ${res}`);
      } else {
        res.json().then((res) => {
          // TODO validate schema
          if (res.content === undefined) {
            error("Malformed response!");
          }
          updateNoteContents(res.content);
        }).catch(function (err) {
          console.error("caught", err);
        });
      }
    });
    return html`<p>Loading <code>${remotePath}</code>...</p>`;
  } else {
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
  }
}

function Body(props) {
  const [route, _] = useContext(Route);

  if (route.startsWith("/editor")) {
    return Editor(props);
  } else if (route == "/browser") {
    return html`<${Browser} />`;
  } else {
    return html`<h1>Error invalid route <code>${route}</code></h1>`;
  }
}

function Nav() {
  const [route, setRoute] = useContext(Route);
  let editorElement, browserElement;
  if (route.startsWith("/editor")) {
    editorElement = html`<a class="secondary">Editor</a>`;
    browserElement = html`<a href="#" onclick=${() => setRoute("/browser")}>Browser</a>`;
  } else if (route == "/browser") {
    editorElement = html`<a href="#" onclick=${() => setRoute("/editor")}>Editor</a>`;
    browserElement = html`<a class="secondary">Browser</a>`;
  } else {
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
  const [route, setRoute] = useState("/browser");
  const routeTuple = [route, setRoute];

  return html`
    <${Route.Provider} value=${routeTuple}>
    <header class="container-fluid">
      <${Nav} />
    </header>
    <main class="container-fluid">
      <${Body} />
    </main>
    </${Route.Provider}>
  `;
}

addEventListener("DOMContentLoaded", (_) =>
  render(html`<${App} />`, document.getElementById("root")));
