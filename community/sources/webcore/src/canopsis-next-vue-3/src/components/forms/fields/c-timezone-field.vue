<template>
  <c-autocomplete
    v-bind="$attrs"
    :model-value="modelValue"
    :items="timezones"
    :label="label || t('common.timezone')"
    :disabled="disabled"
    :name="name"
    item-title="text"
    item-value="value"
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script>
import { uniqBy } from 'lodash';
import { computed } from 'vue';

import { getLocalTimezone, getTimezones } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';
import { useInfo } from '@/hooks/store/modules/info';

export default {
  inheritAttrs: false,

  props: {
    modelValue: {
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

  emits: ['update:modelValue'],

  setup(props) {
    const { t } = useI18n();
    const { timezone: systemTimezone, userTimezones } = useInfo();

    const timezones = computed(() => {
      if (!props.server) {
        return getTimezones().map((timezone) => ({ text: timezone, value: timezone }));
      }

      const serverTimezone = systemTimezone.value || getLocalTimezone();
      const preparedUserTimezones = userTimezones.value.map((timezone) => {
        const text = timezone === serverTimezone ? `${timezone} (${t('common.server')})` : timezone;

        return {
          text,
          value: timezone,
        };
      });

      const localTimezone = getLocalTimezone();

      return uniqBy(
        [
          { text: `${localTimezone} (${t('common.local')})`, value: localTimezone },
          { text: `${serverTimezone} (${t('common.server')})`, value: serverTimezone },
          ...preparedUserTimezones,
        ],
        'value',
      );
    });

    return {
      t,
      timezones,
    };
  },
};
</script>
