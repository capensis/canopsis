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
      v-list
        v-list-tile(v-for="filter in filters", :key="filter.title", @click="")
          v-list-tile-content {{ filter.title }}
          v-list-tile-action
            div
              v-btn.ma-1(icon, @click="showFilterModal(filter)")
                v-icon settings
              v-btn.ma-1(icon, @click="showFilterModal(filter)")
                v-icon delete
      v-btn(color="success", @click.prevent="showFilterModal") Add
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
    showFilterModal(filter) {
      this.showModal({
        name: MODALS.createFilter,
        config: {
          filter,
        },
      });
    },
  },
};
</script>
