<template>
  <v-row class="theme-form__colors" no-gutters>
    <v-col cols="12">
      <theme-color-picker-field
        v-for="item in items"
        v-model="fields[item.value].value"
        :key="item.value"
        :label="t(`theme.extraDetails.${item.value}`)"
        :icon="item.icon"
      />
    </v-col>
  </v-row>
</template>

<script>
import { ALARM_LIST_ACTIONS_TYPES } from '@/constants';

import { getAlarmActionIcon } from '@/helpers/entities/alarm/icons';

import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import ThemeColorPickerField from './fields/theme-color-picker-field.vue';

const ITEMS = [
  { value: 'ack', icon: getAlarmActionIcon(ALARM_LIST_ACTIONS_TYPES.ack) },
  { value: 'canceled', icon: getAlarmActionIcon(ALARM_LIST_ACTIONS_TYPES.fastCancel) },
  { value: 'children', icon: 'center_focus_strong' },
  { value: 'parent', icon: 'center_focus_weak' },
  { value: 'comment', icon: getAlarmActionIcon(ALARM_LIST_ACTIONS_TYPES.comment) },
  { value: 'snooze', icon: getAlarmActionIcon(ALARM_LIST_ACTIONS_TYPES.snooze) },
  { value: 'tickets', icon: getAlarmActionIcon(ALARM_LIST_ACTIONS_TYPES.declareTicket) },
  { value: 'cancel_delay', icon: '$vuetify.icons.close_with_delay' },
];

export default {
  components: { ThemeColorPickerField },

  setup() {
    const { t } = useI18n();
    const fields = useFormModels(
      Object.fromEntries(
        ITEMS.map(({ value }) => [value, `colors.icons.${value}`]),
      ),
    );

    return { fields, items: ITEMS, t };
  },
};
</script>
