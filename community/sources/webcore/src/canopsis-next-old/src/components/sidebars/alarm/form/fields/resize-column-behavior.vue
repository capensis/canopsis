<template>
  <widget-settings-item :title="$t('settings.columnsSettings.title')">
    <v-layout column>
      <c-enabled-field
        v-if="draggable"
        v-field="value.draggable"
        :label="$t('settings.columnsSettings.dragging')"
      />
      <c-enabled-field
        v-if="resizable"
        v-field="value.resizable"
        :label="$t('settings.columnsSettings.resizing')"
      />
      <v-expand-transition>
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
      </v-expand-transition>
    </v-layout>
  </widget-settings-item>
</template>

<script>
import { computed } from 'vue';

import { RESIZING_CELLS_CONTENTS_BEHAVIORS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    draggable: {
      type: Boolean,
      default: false,
    },
    resizable: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const { t } = useI18n();

    const types = computed(() => Object.values(RESIZING_CELLS_CONTENTS_BEHAVIORS).map(value => ({
      value,
      label: t(`settings.columnsSettings.cellsContentBehaviors.${value}`),
    })));

    return {
      types,
    };
  },
};
</script>
