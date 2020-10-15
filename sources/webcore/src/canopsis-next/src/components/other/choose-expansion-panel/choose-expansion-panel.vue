<template lang="pug">
  .choose-expansion-panel
    v-expansion-panel.my-1
      v-expansion-panel-content.grey.darken-2.white--text(lazy)
        div.white--text(slot="header") {{ label }}
        v-card
          chips-list(
            :entities="entities",
            :content-key="contentKey",
            :item-key="itemKey",
            :clearable="clearable",
            @remove="$emit('remove', $event)",
            @clear="$emit('clear')"
          )
          slot
</template>

<script>
import ChipsList from './partials/chips-list.vue';

export default {
  components: { ChipsList },
  props: {
    label: {
      type: String,
      required: true,
    },
    itemKey: {
      type: String,
      default: '_id',
    },
    contentKey: {
      type: String,
      required: false,
    },
    entities: {
      type: Array,
      default: () => [],
    },
    clearable: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style scoped lang="scss">
  .choose-expansion-panel {
    & /deep/ .v-expansion-panel__header .v-icon {
      color: white !important;
    }
  }
</style>
