<template>
  <widget-settings-item :title="$t('settings.treeOfDependenciesSettings')">
    <v-layout>
      <v-radio-group
        v-field="value"
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
import { TREE_OF_DEPENDENCIES_SHOW_TYPES } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: Number,
      default: TREE_OF_DEPENDENCIES_SHOW_TYPES.custom,
    },
  },
  computed: {
    types() {
      return Object.values(TREE_OF_DEPENDENCIES_SHOW_TYPES).map(value => ({
        value,
        label: this.$t(`entity.treeOfDependenciesShowTypes.${value}`),
      }));
    },
  },
};
</script>
