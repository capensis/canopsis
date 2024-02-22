<template>
  <widget-settings-item :title="$t('settings.columnsSettings.title')">
    <v-layout column>
      <c-enabled-field
        v-field="value.draggable"
        :label="$t('settings.columnsSettings.dragging')"
      />
      <c-enabled-field
        v-field="value.resizable"
        :label="$t('settings.columnsSettings.resizing')"
      />
      <v-radio-group
        v-if="value.resizable"
        v-field="value.cells_content_behavior"
        :label="$t('settings.columnsSettings.cellsContentBehavior')"
        class="mt-0"
        name="opened"
        hide-details
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
import { ALARMS_RESIZING_CELLS_CONTENTS_BEHAVIORS } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
  computed: {
    types() {
      return Object.values(ALARMS_RESIZING_CELLS_CONTENTS_BEHAVIORS).map(value => ({
        value,
        label: this.$t(`settings.columnsSettings.cellsContentBehaviors.${value}`),
      }));
    },
  },
};
</script>
