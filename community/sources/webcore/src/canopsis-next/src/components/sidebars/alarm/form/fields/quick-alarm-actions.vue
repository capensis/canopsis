<template>
  <widget-settings-item :title="title">
    <v-layout column>
      <span v-if="description">{{ description }}</span>
      <quick-alarm-actions-form
        v-field="value"
        :massive="massive"
      />
    </v-layout>
  </widget-settings-item>
</template>

<script>
import { computed } from 'vue';

import { ALARMS_OPENED_VALUES } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModelField } from '@/hooks/form/model-field';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';
import QuickAlarmActionsForm from '@/components/common/actions-panel/quick-alarm-actions-form.vue';

export default {
  components: { WidgetSettingsItem, QuickAlarmActionsForm },
  props: {
    value: {
      type: Array,
      default: () => [],
    },
    massive: {
      type: Boolean,
      default: false,
    },
    title: {
      type: String,
      default: '',
    },
    description: {
      type: String,
      default: '',
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateModel } = useModelField(props, emit);

    const localValue = computed({
      get: () => String(props.value),
      set: newValue => updateModel({
        true: true,
        false: false,
        null: null,
      }[newValue]),
    });

    const types = computed(() => Object.values(ALARMS_OPENED_VALUES).map(value => ({
      value: String(value),
      label: t(`settings.openedTypes.${value}`),
    })));

    return {
      localValue,
      types,
    };
  },
};
</script>
