<template lang="pug">
  tr
    td
      v-checkbox(primary, hide-details, v-model="row.selected")
    td(v-for="column in columns", @click="toggleExpandPanel")
      div(v-if="column.value === 'enabled'")
        enabled-column(:value="row.item.enabled")
      ellipsis(
        v-else,
        :text="row.item | get(column.value, null, '')",
        :maxLetters="column.maxLetters"
      )
    td
      actions-panel(:item="row.item", :isEditingMode="isEditingMode", @action="closeExpandPanel")
</template>

<script>
import { WIDGETS_ACTIONS_TYPES } from '@/constants';

import Ellipsis from '@/components/tables/ellipsis.vue';
import EnabledColumn from '@/components/tables/enabled-column.vue';

import ActionsPanel from '../actions/actions-panel.vue';

export default {
  components: {
    Ellipsis,
    ActionsPanel,
    EnabledColumn,
  },
  props: {
    row: {
      type: Object,
      required: true,
    },
    columns: {
      type: Array,
      required: true,
    },
    isEditingMode: {
      type: Boolean,
      required: true,
    },
  },
  methods: {
    toggleExpandPanel() {
      this.row.expanded = !this.row.expanded;
    },

    closeExpandPanel(type) {
      if (type === WIDGETS_ACTIONS_TYPES.context.editEntity) {
        this.row.expanded = false;
      }
    },
  },
};
</script>
