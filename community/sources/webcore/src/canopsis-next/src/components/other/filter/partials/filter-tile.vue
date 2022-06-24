<template lang="pug">
  v-list-tile.pa-0
    v-layout(row, align-center)
      v-flex(xs12)
        v-list-tile-content {{ filter.title }}
      v-list-tile-action(v-if="editable")
        v-layout(row, align-center)
          v-tooltip(v-if="filter.old_mongo_query", top)
            template(#activator="{ on, attrs }")
              v-icon.cursor-pointer.mr-2(v-on="on", v-bind="attrs", color="grey") error_outline
            span {{ $t('filter.oldPattern') }}
          c-action-btn(type="edit", @click="$emit('edit')")
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
