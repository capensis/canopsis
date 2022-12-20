<template lang="pug">
  v-list-tile.pa-0
    v-layout(row, align-center)
      v-flex(xs12)
        v-layout(row, align-center)
          v-icon.draggable.ml-0.mr-3.action-drag-handler(v-if="editable", small) drag_indicator
          v-list-tile-content {{ filter.title }}
      v-list-tile-action(v-if="editable")
        v-layout(row, align-center)
          c-action-btn(
            type="edit",
            :badge-value="!!filter.old_mongo_query",
            :badge-tooltip="$t('pattern.oldPatternTooltip')",
            @click="$emit('edit')"
          )
          c-action-btn(type="delete", @click="$emit('delete')")
</template>

<script>
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
};
</script>
