import { escapeRegExp } from 'lodash';
import Handlebars from 'handlebars';

import { MAX_LIMIT } from '@/constants';

import { registerHelper, unregisterHelper } from '@/helpers/handlebars';

export const alarmHandlebarsTagsHelper = {
  beforeCreate() {
    /**
     * @example {{tags 'tag1' 'tag2' regex='tag.*' flags='i'}}
     */
    registerHelper('tags', (...args) => {
      const { hash: { regex = '', flags = '' } = {} } = args.pop() ?? {};
      const escapedArgs = args.filter(Boolean).map(arg => escapeRegExp(arg));
      const nameFilter = escapedArgs.length ? `^(${escapedArgs.join('|')})$` : '';

      return new Handlebars.SafeString(
        `<c-alarm-tags-chips
        :alarm="alarm"
        :selected-tag="selectedTag"
        name-filter="${nameFilter}"
        regex-filter="${regex}"
        regex-filter-flags="${flags}"
        inline-count="${MAX_LIMIT}"
        closable-active
        @select="$emit('select:tag', $event)"
        @close="$emit('clear:tag', $event)"
      ></c-alarm-tags-chips>`,
      );
    });
  },
  destroyed() {
    unregisterHelper('tags');
  },
};
