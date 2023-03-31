<template lang="pug">
  div.choose-expansion-panel
    v-expansion-panel.my-1
      v-expansion-panel-content.grey.darken-2.white--text(:class="{ error: errors.length }", lazy)
        template(#header="")
          div.white--text {{ label }}
        v-card.pt-1
          v-alert.pa-2.mx-2(type="error", :value="!!errors.length") {{ errors.join(' ') }}
          chips-list(
            :entities="entities",
            :disabled-entities="disabledEntities",
            :existing-entities="existingEntities",
            :content-key="contentKey",
            :item-key="itemKey",
            :clearable="clearable",
            @remove="$listeners.remove",
            @clear="$listeners.clear"
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
    disabledEntities: {
      type: Array,
      default: () => [],
    },
    existingEntities: {
      type: Array,
      default: () => [],
    },
    errors: {
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
    & ::v-deep .v-expansion-panel__header .v-icon {
      color: white !important;
    }
  }
</style>
