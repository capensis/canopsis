<template>
  <widget-settings-item :title="$t('settings.filterOnOpenResolved')">
    <c-layout>
      <c-radio-group
        v-model="localValue"
        class="mt-0"
        name="opened"
        color="primary"
        hide-details
      >
        <v-radio
          v-for="type in types"
          :key="type.value"
          :label="type.label"
          :value="type.value"
          color="primary"
        />
      </c-radio-group>
    </c-layout>
  </widget-settings-item>
</template>

<script>
import { ALARMS_OPENED_VALUES } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },

  props: {
    modelValue: {
      type: Boolean,
      required: false,
    },
  },

  emits: ['update:modelValue'],

  computed: {
    localValue: {
      get() {
        return String(this.modelValue);
      },
      set(value) {
        this.$emit(
          'update:modelValue',
          {
            true: true,
            false: false,
            null: null,
          }[value],
        );
      },
    },

    types() {
      return Object.values(ALARMS_OPENED_VALUES).map(value => ({
        value: String(value),
        label: this.$t(`settings.openedTypes.${value}`),
      }));
    },
  },
};
</script>
