<template lang="pug">
  v-list-tile.pa-0
    v-layout(row, align-center)
      v-flex(xs12)
        v-layout(row, align-center)
          v-icon.draggable.ml-0.mr-3.action-drag-handler(v-if="editable", small) drag_indicator
          v-list-tile-content {{ filter.title }}
      v-list-tile-action(v-if="editable")
        v-layout(row, align-center)
          v-tooltip(v-if="!hasSomeOnePattern && isOldPatternFormat", top)
            template(#activator="{ on, attrs }")
              v-icon.cursor-pointer.mr-2(v-on="on", v-bind="attrs", color="grey") error_outline
            span {{ $t('filter.oldPattern') }}
          c-action-btn(type="edit", @click="$emit('edit')")
          c-action-btn(type="delete", @click="$emit('delete')")
</template>

<script>
import { PATTERNS_FIELDS } from '@/constants';

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
    hasSomeOnePattern() {
      return Object.values(PATTERNS_FIELDS).some(field => this.filter[field]?.length);
    },

    isOldPatternFormat() {
      return this.filter.old_mongo_query;
    },
  },
};
</script>
