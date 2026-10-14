import { computed, reactive } from 'vue';

import {
  DEFAULT_USER_INTERFACE_DATE_FORMAT,
  USER_INTERFACE_DATE_FORMATS,
} from '@/constants';

import { convertDateToString } from '@/helpers/date/date';
import { registerUserInterfaceDateFormatProvider } from '@/helpers/date/date-format';

import { systemMixin } from '@/mixins/system';

describe('system mixin', () => {
  const date = new Date(2021, 8, 20, 21, 15, 40);

  afterEach(() => {
    registerUserInterfaceDateFormatProvider(() => DEFAULT_USER_INTERFACE_DATE_FORMAT);
  });

  it('reactively applies the date format stored in system', () => {
    const context = {
      dateFormat: undefined,
      timezone: undefined,
    };

    context.system = reactive(systemMixin.data.call(context).system);
    context.setSystemData = systemMixin.methods.setSystemData.bind(context);
    systemMixin.created.call(context);

    const formattedDate = computed(() => convertDateToString(date));

    expect(formattedDate.value).toBe('20/09/2021 21:15:40');

    context.dateFormat = USER_INTERFACE_DATE_FORMATS.weekdayWithDateTime.value;
    systemMixin.methods.syncSystemInfo.call(context);

    expect(formattedDate.value).toBe('Mon 20/09/2021 21:15:40');
  });
});
