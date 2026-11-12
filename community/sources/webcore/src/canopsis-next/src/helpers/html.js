import { registerCustomProtocol } from 'linkifyjs';
import linkifyHtmlLib from 'linkify-html';
import sanitizeHtmlLib from 'sanitize-html';

import { DEFAULT_SANITIZE_OPTIONS, DEFAULT_LINKIFY_OPTIONS, LINKIFY_PROTOCOLS } from '@/config';

import { TEMPLATE_VUE_ATTRIBUTES, CUSTOM_ICON_SANITIZE_OPTIONS } from '@/constants/sanitize';

LINKIFY_PROTOCOLS.forEach(registerCustomProtocol);

const VUE_ATTRIBUTE = /^(?:v-|:|\.|@|#|data-v-)/i;
const UNSAFE_ATTRIBUTE = /^(?:on|is$|tag$|ref$|innerhtml$|textcontent$|srcdoc$|slot$|slot-scope$|scope$|inline-template$)/i;

/**
 * Parse in an inert fragment: raw-text and attribute contents are not comments,
 * and SVG/HTML parsing follows the browser that will ultimately render it.
 *
 * @param {string} html - HTML markup to parse without activating its contents.
 * @returns {HTMLTemplateElement} Template containing the parsed fragment.
 */
const parseHtml = (html) => {
  const element = document.createElement('template');
  element.innerHTML = html;

  return element;
};

/**
 * Visit a root node and its descendants, including nested template contents.
 *
 * @param {Node} root - Root of the traversal.
 * @param {function(Node): void} visit - Callback invoked for each node.
 * @returns {void}
 */
const visitNodes = (root, visit) => {
  const nodes = [root];
  while (nodes.length) {
    const node = nodes.pop();
    visit(node);
    nodes.push(...(node.content || node).childNodes);
  }
};

/**
 * Validate Vue expressions against the exact helper bindings and reject DOM overrides.
 * Ordinary attributes still require validation by the HTML sanitizer.
 *
 * @param {string} tag - Element or component tag name.
 * @param {string} name - Attribute name, including any directive prefix.
 * @param {string} value - Decoded attribute value.
 * @returns {boolean} Whether the attribute passes the Vue template policy.
 */
const isAllowedAttribute = (tag, name, value) => {
  if (VUE_ATTRIBUTE.test(name)) {
    const allowed = TEMPLATE_VUE_ATTRIBUTES[tag];

    return Boolean(allowed && Object.prototype.hasOwnProperty.call(allowed, name) && allowed[name] === value);
  }

  return !UNSAFE_ATTRIBUTE.test(name);
};

/**
 * Copy attributes that pass the mandatory Vue template policy.
 *
 * @param {string} tag - Element or component tag name.
 * @param {Object.<string, string>} attributes - Decoded attributes from sanitize-html.
 * @returns {Object.<string, string>} Attributes to pass to the HTML sanitizer.
 */
const filterAttributes = (tag, attributes) => Object.entries(attributes).reduce((acc, [name, value]) => {
  if (isAllowedAttribute(tag, name, value)) {
    acc[name] = value;
  }

  return acc;
}, {});

/**
 * Keep comment contents only in this call's map. Tokens never originate in data-c
 * or another user attribute. Restore before parsing the sanitized string again, so table
 * comments cannot undergo foster parenting as placeholder elements or text.
 *
 * @param {string} html - Source markup containing comments to preserve.
 * @returns {{html: string, restore: function(string): string}} Markup with placeholders
 * and a function that restores only this call's comments in the sanitized string.
 */
const preserveComments = (html) => {
  const element = parseHtml(html);
  const serialized = element.innerHTML;
  let prefix;
  do {
    prefix = `canopsisComment${window.crypto.getRandomValues(new Uint32Array(4)).join('')}`;
  } while (serialized.includes(prefix) || html.includes(prefix));

  const comments = new Map();
  visitNodes(element.content, (node) => {
    if (node.nodeType !== Node.COMMENT_NODE) {
      return;
    }
    const { data } = node;
    if (/--|^>|^->|-$/.test(data)) {
      node.remove();

      return;
    }
    const token = `${prefix}_${comments.size}_end`;
    comments.set(token, `<!--${data}-->`);
    node.replaceWith(document.createTextNode(token));
  });

  return {
    html: element.innerHTML,
    restore: result => [...comments].reduce((value, [token, comment]) => value.split(token).join(comment), result),
  };
};

/**
 * Sanitize source HTML without consuming Handlebars expressions. Custom options
 * cannot replace the mandatory validation of Vue attributes.
 *
 * @param {string} [html=''] - Source HTML, which may still contain Handlebars expressions.
 * @param {Object} [options=DEFAULT_SANITIZE_OPTIONS] - sanitize-html options overriding the defaults.
 * @returns {string} Sanitized HTML with safe comments and approved helper bindings preserved.
 */
export const sanitizeHtml = (html = '', options = DEFAULT_SANITIZE_OPTIONS) => {
  const preserved = preserveComments(html);
  const merged = { ...DEFAULT_SANITIZE_OPTIONS, ...options };
  const transforms = merged.transformTags || {};
  const allowedAttributes = { ...merged.allowedAttributes };
  Object.entries(TEMPLATE_VUE_ATTRIBUTES).forEach(([tag, attributes]) => {
    allowedAttributes[tag] = [...(allowedAttributes[tag] || []), ...Object.keys(attributes)];
  });

  return preserved.restore(sanitizeHtmlLib(preserved.html, {
    ...merged,
    allowedAttributes,
    parser: { ...merged.parser, decodeEntities: true, lowerCaseTags: true, lowerCaseAttributeNames: true },
    transformTags: {
      ...transforms,
      '*': (tag, attributes) => {
        const transform = transforms['*'];
        const transformed = typeof transform === 'function' ? transform(tag, attributes) : { tagName: tag, attribs: attributes };

        return { ...transformed, attribs: filterAttributes(transformed.tagName, transformed.attribs) };
      },
    },
  }));
};

/**
 * The last boundary before Vue 2 compilation, after HBS, linkification, wrapping
 * and HTML normalization. Vue 2 decodes entities before parsing interpolations:
 * split literal opening braces with empty comments instead of entity escaping.
 * The compiler sees separate text nodes, while displayed/copied text is intact.
 *
 * @param {string} [html=''] - Compiled Handlebars HTML, including any wrapper and linkification.
 * @param {Object} [options=DEFAULT_SANITIZE_OPTIONS] - HTML sanitization options.
 * @returns {string} HTML ready for Vue compilation without further normalization or decoding.
 */
export const sanitizeVueTemplate = (html = '', options = DEFAULT_SANITIZE_OPTIONS) => {
  const element = parseHtml(sanitizeHtml(html, options));
  visitNodes(element.content, (node) => {
    if (node.nodeType === Node.ELEMENT_NODE) {
      [...node.attributes].forEach(({ name, value }) => {
        if (!isAllowedAttribute(node.localName, name, value)) {
          node.removeAttribute(name);
        }
      });
    }
    if (node.nodeType !== Node.TEXT_NODE || !node.data.includes('{{')) {
      return;
    }
    const fragment = document.createDocumentFragment();
    node.data.split('{').forEach((part, index) => {
      if (index) {
        fragment.append(document.createTextNode('{'), document.createComment(''));
      }
      fragment.append(document.createTextNode(part));
    });
    node.replaceWith(fragment);
  });

  return element.innerHTML;
};

/**
 * Sanitize uploaded SVG independently of the HTML policy. Only local fragment
 * links and paint references survive. Browser serialization preserves SVG case.
 *
 * @param {string} [html=''] - Uploaded or previously stored SVG markup.
 * @returns {string} Sanitized first SVG element, or an empty string when no SVG is present.
 */
export const sanitizeSvg = (html = '') => {
  const element = parseHtml(html);
  const svg = element.content.querySelector('svg');
  if (!svg) {
    return '';
  }
  const sanitized = sanitizeHtmlLib(svg.outerHTML, {
    ...CUSTOM_ICON_SANITIZE_OPTIONS,
    transformTags: {
      '*': (tagName, attributes) => ({
        tagName,
        attribs: Object.entries(attributes).reduce((acc, [name, value]) => {
          const localLink = !/^(?:href|xlink:href)$/i.test(name) || /^#[\w:.-]+$/.test(value);
          const localPaint = !/url\s*\(/i.test(value) || /^url\(\s*(['"]?)#[\w:.-]+\1\s*\)$/i.test(value);
          if (!VUE_ATTRIBUTE.test(name) && !UNSAFE_ATTRIBUTE.test(name) && localLink && localPaint && !value.includes('\\')) {
            acc[name] = value;
          }

          return acc;
        }, {}),
      }),
    },
  });

  return parseHtml(sanitized).innerHTML;
};

/**
 * Convert all links in html to tag <a>
 *
 * @param {string} [html = '']
 * @param {Object} [options = DEFAULT_LINKIFY_OPTIONS]
 * @return {string}
 */
export const linkifyHtml = (html = '', options = DEFAULT_LINKIFY_OPTIONS) => linkifyHtmlLib(html, options);

/**
 * Normalize HTML using the browser parser, including closing unclosed tags.
 *
 * @param {string} [html = '']
 * @return {string}
 */
export const normalizeHtml = (html = '') => {
  const element = document.createElement('div');

  element.innerHTML = html;

  return element.innerHTML;
};
