/**
 * Create Node Element from html text
 *
 * @param {string} [html]
 * @return {ChildNode}
 */
const htmlToElement = (html) => {
  const template = document.createElement('template');

  template.innerHTML = html;

  return template.content.firstChild;
}

/**
 * Load script
 *
 * @param {string} [script]
 * @return {Promise<unknown>}
 */
const loadScript = (script) => {
  return new Promise((resolve, reject) => {
    const scriptNode = document.createElement('script');
    const sourceScriptNode = htmlToElement(script);

    for (let i = 0; i < sourceScriptNode.attributes.length; i += 1) {
      const attribute = sourceScriptNode.attributes[i];

      scriptNode.setAttribute(attribute.name, attribute.value);
    }

    scriptNode.innerHTML = sourceScriptNode.innerHTML;
    scriptNode.onload = resolve;
    scriptNode.onerror = reject;

    document.head.append(scriptNode);
  });
}

if (Array.isArray(window.EXTERNAL_SCRIPTS)) {
  window.EXTERNAL_SCRIPTS.forEach((script) => {
    if (!script) {
      return;
    }

    loadScript(script)
      .catch(console.error);
  });
}
