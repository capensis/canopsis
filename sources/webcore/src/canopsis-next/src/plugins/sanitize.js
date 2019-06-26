import { OPTIONS_SANITIZE_TEXT_EDITOR } from '@/constants';
import sanitizeHTML from 'sanitize-html';

export default {
  install(Vue) {
    Object.defineProperties(Vue.prototype, {
      $sanitize: {
        get() {
          return (HTML, options = OPTIONS_SANITIZE_TEXT_EDITOR) => {
            try {
              return sanitizeHTML(HTML, options);
            } catch (err) {
              console.warn(err);

              return '';
            }
          };
        },
      },
    });
  },
};
