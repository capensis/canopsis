<template>
  <v-list-group>
    <template #activator="">
      <v-list-item>{{ $t('settings.columnsSettings.title') }}</v-list-item>
    </template>
    <v-container>
      <v-layout column="column">
        <c-enabled-field
          v-field="value.draggable"
          :label="$t('settings.columnsSettings.dragging')"
        />
        <c-enabled-field
          v-field="value.resizable"
          :label="$t('settings.columnsSettings.resizing')"
        />
        <v-radio-group
          class="mt-0"
          v-if="value.resizable"
          v-field="value.cells_content_behavior"
          :label="$t('settings.columnsSettings.cellsContentBehavior')"
          name="opened"
          hide-details="hide-details"
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
    </v-container>
  </v-list-group>
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
