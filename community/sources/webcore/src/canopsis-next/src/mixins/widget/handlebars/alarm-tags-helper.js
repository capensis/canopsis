import Handlebars from 'handlebars';

import { MAX_LIMIT } from '@/constants';

import { registerHelper, unregisterHelper } from '@/helpers/handlebars';

export const alarmHandlebarsTagsHelper = {
  beforeCreate() {
    registerHelper('tags', () => new Handlebars.SafeString(
      `<c-alarm-tags-chips
        :alarm="alarm"
        :selected-tag="selectedTag"
        inline-count="${MAX_LIMIT}"
        closable-active
        @select="$emit('select:tag', $event)"
        @close="$emit('clear:tag', $event)"
      ></c-alarm-tags-chips>`,
    ));
  },
  destroyed() {
    unregisterHelper('tags');
  },
};
