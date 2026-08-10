<template>
  <widget-settings-item :title="$t('settings.columnsSettings.title')">
    <c-layout column>
      <c-enabled-field
        v-if="draggable"
        :model-value="modelValue.draggable"
        :label="$t('settings.columnsSettings.dragging')"
        @update:model-value="updateField('draggable', $event)"
      />
      <c-enabled-field
        v-if="resizable"
        :model-value="modelValue.resizable"
        :label="$t('settings.columnsSettings.resizing')"
        @update:model-value="updateField('resizable', $event)"
      />
      <v-expand-transition>
        <c-radio-group
          v-if="modelValue.resizable"
          :model-value="modelValue.cells_content_behavior"
          :label="$t('settings.columnsSettings.cellsContentBehavior')"
          class="mt-0"
          name="cells_content_behavior"
          hide-details
          @update:model-value="updateField('cells_content_behavior', $event)"
        >
          <v-radio
            v-for="type in types"
            :key="type.value"
            :label="type.label"
            :value="type.value"
            color="primary"
          />
        </c-radio-group>
      </v-expand-transition>
    </c-layout>
  </widget-settings-item>
</template>

<script>
import { computed } from 'vue';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';
import { RESIZING_CELLS_CONTENTS_BEHAVIORS } from '@/constants';
import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

export default {
  components: { WidgetSettingsItem },

  props: {
    modelValue: {
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

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateField } = useModelField(props, emit);
    const types = computed(() =>
      Object.values(RESIZING_CELLS_CONTENTS_BEHAVIORS).map(value => ({
        value,
        label: t(`settings.columnsSettings.cellsContentBehaviors.${value}`),
      })),
    );

    return {
      types,
      updateField,
    };
  },
};
</script>
