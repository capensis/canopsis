import sanitizeHTML from 'sanitize-html';

export default {
  install(Vue, { defaultOptions = {} } = {}) {
    Object.defineProperty(Vue.prototype, '$sanitize', {
      get() {
        return (HTML, options = defaultOptions) => {
          try {
            return sanitizeHTML(HTML, options);
          } catch (err) {
            console.warn(err);

            return '';
          }
        };
      },
    });
  },
};
