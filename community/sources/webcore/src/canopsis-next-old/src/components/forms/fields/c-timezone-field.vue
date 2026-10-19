<template>
  <v-autocomplete
    v-field="value"
    v-bind="$attrs"
    :items="timezones"
    :label="label || $t('common.timezone')"
    :disabled="disabled"
    :name="name"
  />
</template>

<script>
import { uniqBy } from 'lodash';
import { computed, inject } from 'vue';

import { getLocalTimezone, getTimezones } from '@/helpers/date/date';

import { useInfo } from '@/hooks/store/modules/info';
import { useI18n } from '@/hooks/i18n';

export default {
  inject: ['$validator'],
  inheritAttrs: false,
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: getLocalTimezone(),
    },
    label: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'timezone',
    },
    server: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const system = inject('$system');

    const { t } = useI18n();
    const { userTimezones } = useInfo();

    const timezones = computed(() => {
      if (!props.server) {
        return getTimezones();
      }

      const preparedUserTimezones = userTimezones.value.map((timezone) => {
        const text = timezone === system.timezone ? `${timezone} (${t('common.server')})` : timezone;

        return {
          text,
          value: timezone,
        };
      });

      const localeTimezone = getLocalTimezone();

      return uniqBy([
        { text: `${localeTimezone} (${t('common.local')})`, value: localeTimezone },
        { text: `${system.timezone} (${t('common.server')})`, value: system.timezone },
        ...preparedUserTimezones,
      ], 'value');
    });

    return {
      timezones,
    };
  },
};
</script>
