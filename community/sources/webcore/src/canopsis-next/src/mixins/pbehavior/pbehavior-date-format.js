import { DATETIME_FORMATS } from '@/constants';

import { convertDateToTimezoneDateString, getLocalTimezone } from '@/helpers/date/date';

export const pbehaviorsDateFormatMixin = {
  data() {
    return {
      timezone: getLocalTimezone(),
    };
  },
  methods: {
    formatIntervalDate(item, field) {
      const date = item[field];
      const format = item.rrule ? DATETIME_FORMATS.medium : DATETIME_FORMATS.long;

      return convertDateToTimezoneDateString(date, this.timezone, format);
    },

    formatRruleEndDate(item) {
      if (!item.rrule) {
        return '-';
      }

      return item.rrule_end
        ? convertDateToTimezoneDateString(item.rrule_end, this.timezone, DATETIME_FORMATS.long)
        : this.$t('common.undefined');
    },
  },
};
