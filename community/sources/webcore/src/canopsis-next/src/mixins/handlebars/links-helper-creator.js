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
      const category = hash.category ? `'${hash.category}'` : undefined;

      return new Handlebars.SafeString(`
          <c-links-list v-if="links && hasAccessForLinks" :links="links" :category="${category}" />
        `);
    });
  },
  beforeUnmount() {
    unregisterHelper('links');
  },
});
