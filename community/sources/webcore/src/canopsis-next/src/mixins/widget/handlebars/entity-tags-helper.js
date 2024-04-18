import Handlebars from 'handlebars';

import { MAX_LIMIT } from '@/constants';

import { registerHelper, unregisterHelper } from '@/helpers/handlebars';

export const entityHandlebarsTagsHelper = {
  beforeCreate() {
    registerHelper('tags', () => new Handlebars.SafeString(
      `<c-entity-tags-chips
        :entity="entity"
        inline-count="${MAX_LIMIT}"
       ></c-entity-tags-chips>`,
    ));
  },
  destroyed() {
    unregisterHelper('tags');
  },
};
