<template>
  <widget-settings-item :title="$t('settings.filterOnOpenResolved')">
    <v-layout>
      <v-radio-group
        v-model="localValue"
        class="mt-0"
        name="opened"
        hide-details
        mandatory
      >
        <v-radio
          v-for="type in types"
          :key="type.value"
          :label="type.label"
          :value="type.value"
          color="primary"
        />
      </v-radio-group>
    </v-layout>
  </widget-settings-item>
</template>

<script>
import { computed } from 'vue';

import { ALARMS_OPENED_VALUES } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModelField } from '@/hooks/form/model-field';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: Boolean,
      required: false,
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
