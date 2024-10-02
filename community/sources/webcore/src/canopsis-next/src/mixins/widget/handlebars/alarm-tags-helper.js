import { isString, escapeRegExp } from 'lodash';
import Handlebars from 'handlebars';

import { MAX_LIMIT } from '@/constants';

import { registerHelper, unregisterHelper } from '@/helpers/handlebars';

export const alarmHandlebarsTagsHelper = {
  beforeCreate() {
    registerHelper('tags', (...args) => {
      const filterPatterns = args.reduce((acc, arg) => {
        if (!arg) {
          return acc;
        }

        if (isString(arg)) {
          acc.push(escapeRegExp(arg));

          return acc;
        }

        const { regex = '' } = arg?.hash ?? {};

        if (regex) {
          acc.push(`(${regex})`);
        }

        return acc;
      }, []).join('|');

      return new Handlebars.SafeString(
        `<c-alarm-tags-chips
        :alarm="alarm"
        :selected-tag="selectedTag"
        filter-pattern="${filterPatterns}"
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
