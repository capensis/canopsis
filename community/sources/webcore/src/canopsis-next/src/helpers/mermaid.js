import mermaid from 'mermaid';

import uid from '@/helpers/uid';

/**
 * Render mermaid diagram
 *
 * @param {string} value
 * @param {Object} [options = {}]
 * @return {Promise<string>}
 */
export const renderMermaid = (value = '', options = {}) => {
  mermaid.initialize({ ...options, startOnLoad: false });

  /**
   * mermaid.render is dirty function, because function removes svg on each re-render.
   */
  const id = `id_${uid()}`;

  try {
    const svg = mermaid.render(id, value);

    mermaid.mermaidAPI.reinitialize();

    return svg;
  } catch (e) {
    /**
     * Mermaid doesn't clear a body if error is throw
     */
    const element = document.getElementById(`d${id}`);

    if (element) {
      element.remove();
    }

    throw e;
  }
};

/**
 * Check is mermaid diagram valid
 *
 * @param {string} value
 * @return {boolean}
 */
export const validateMermaidDiagram = (value) => {
  mermaid.parse(value);
};
