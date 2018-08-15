<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.filters') }}
    v-container
      v-select(
      :label="$t('settings.selectAFilter')",
      :items="filters",
      :value="value",
      @change="updateField",
      item-text="title",
      item-value="title",
      clearable
      )
    v-container
      v-list
        v-list-tile(@click="")
          v-list-tile-content Some filterSome filterSome filterSome filter
          v-list-tile-action
            div
              v-btn.ma-1(icon, @click="showFilterModal")
                v-icon delete
              v-btn.ma-1(icon)
                v-icon settings
        v-divider
        v-list-tile(@click="")
          v-list-tile-content Test
</template>

<script>
import { MODALS } from '@/constants';
import modalMixin from '@/mixins/modal/modal';

export default {
  mixins: [modalMixin],
  props: {
    filters: {
      type: Array,
      default: () => [],
    },
    value: {
      type: Object,
      default: () => ({}),
    },
  },
  methods: {
    updateField(value) {
      this.$emit('input', this.filters.find(v => v.title === value));
    },
    showFilterModal() {
      this.showModal({
        name: MODALS.confirmation,
      });
    },
  },
};
</script>
