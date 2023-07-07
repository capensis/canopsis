<template lang="pug">
  v-list-group
    template(#activator="")
      v-list-tile {{ $t('settings.columnsSettings.title') }}
    v-container
      v-layout(column)
        c-enabled-field(
          v-model="value.draggable",
          :label="$t('settings.columnsSettings.dragging')"
        )
        c-enabled-field(
          v-model="value.resizable",
          :label="$t('settings.columnsSettings.resizing')"
        )
        v-radio-group.mt-0(
          v-if="value.resizable",
          v-field="value.cells_content_behavior",
          :label="$t('settings.columnsSettings.cellsContentBehavior')",
          name="opened",
          hide-details
        )
          v-radio(
            v-for="type in types",
            :key="type.value",
            :label="type.label",
            :value="type.value",
            color="primary"
          )
</template>

<script>
import { ALARMS_RESIZING_CELLS_CONTENTS_BEHAVIORS } from '@/constants';

export default {
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
