<template>
  <v-list-item class="pa-0">
    <v-layout align-center>
      <v-flex xs12>
        <v-layout align-center>
          <v-icon
            class="draggable ml-0 mr-3 action-drag-handler"
            v-if="editable"
            small
          >
            drag_indicator
          </v-icon>
          <v-list-item-content>{{ filter.title }}</v-list-item-content>
        </v-layout>
      </v-flex>
      <v-list-item-action v-if="editable">
        <v-layout align-center>
          <c-action-btn
            type="edit"
            :badge-value="isOldPattern"
            :badge-tooltip="$t('pattern.oldPatternTooltip')"
            @click="$emit('edit')"
          />
          <c-action-btn
            type="delete"
            @click="$emit('delete')"
          />
        </v-layout>
      </v-list-item-action>
    </v-layout>
  </v-list-item>
</template>

<script>
import { isOldPattern } from '@/helpers/entities/pattern/form';

export default {
  props: {
    filter: {
      type: Object,
      default: () => ({}),
    },
    editable: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isOldPattern() {
      return isOldPattern(this.filter);
    },
  },
};
</script>
