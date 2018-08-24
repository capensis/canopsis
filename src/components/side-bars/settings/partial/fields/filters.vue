<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.filters') }}
    v-container
      v-select(
      :label="$t('settings.selectAFilter')",
      :items="filters",
      :value="value",
      @change="updateSelectedField",
      item-text="title",
      item-value="title",
      clearable
      )
      v-list
        v-list-tile(v-for="(filter, index) in filters", :key="filter.title", @click="")
          v-list-tile-content {{ filter.title }}
          v-list-tile-action
            div
              v-btn.ma-1(icon, @click="showEditFilterModal(index)")
                v-icon settings
              v-btn.ma-1(icon, @click="showDeleteFilterModal(index)")
                v-icon delete
      v-btn(color="success", @click.prevent="showCreateFilterModal") {{ $t('common.add') }}
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
    updateSelectedField(title) {
      this.$emit('input', this.filters.find(v => v.title === title));
    },

    showCreateFilterModal() {
      this.showModal({
        name: MODALS.createFilter,
        config: {
          title: 'modals.filter.create.title',
          action: newFilter => this.$emit('update:filters', [...this.filters, newFilter]),
        },
      });
    },

    showEditFilterModal(index) {
      const filter = this.filters[index];

      this.showModal({
        name: MODALS.createFilter,
        config: {
          title: 'modals.filter.edit.title',
          filter,
          action: (newFilter) => {
            if (this.value.title === filter.title) {
              this.$emit('input', newFilter);
            }

            this.$emit('update:filters', [
              ...this.filters.map((v, i) => (index === i ? newFilter : v)),
            ]);
          },
        },
      });
    },

    showDeleteFilterModal(index) {
      const filter = this.filters[index];

      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => {
            if (this.value.title === filter.title) {
              this.$emit('input', {});
            }

            this.$emit('update:filters', ...this.filters.filter((v, i) => index !== i));
          },
        },
      });
    },
  },
};
</script>
