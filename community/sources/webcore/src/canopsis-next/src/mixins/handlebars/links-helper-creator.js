import Handlebars from 'handlebars';
import { get } from 'lodash';

import { registerHelper, unregisterHelper } from '@/helpers/handlebars';

import { authMixin } from '@/mixins/auth';

export const handlebarsLinksHelperCreator = (linksPath, right) => ({
  mixins: [authMixin],
  computed: {
    hasAccessForLinks() {
      return this.checkAccess(right);
    },

    links() {
      return get(this, linksPath);
    },
  },
  beforeCreate() {
    registerHelper('links', ({ hash }) => {
      /**
       * Handlebars has already resolved variables in hash.category. The old Vue
       * binding wrapped this value in quotes, so it was a string literal, not a lookup.
       */
      const categoryAttribute = hash.category
        ? `category="${Handlebars.escapeExpression(String(hash.category))}"`
        : '';

      return new Handlebars.SafeString(`
          <c-links-list v-if="links && hasAccessForLinks" :links="links" ${categoryAttribute} />
        `);
    });
  },
  beforeUnmount() {
    unregisterHelper('links');
  },
});
