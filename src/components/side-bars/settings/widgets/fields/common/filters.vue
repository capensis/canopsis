<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.filters') }}
    v-container
      v-select(
      v-show="hasAccessToEditFilter && !hideSelect",
      :label="$t('settings.selectAFilter')",
      :items="filters",
      :value="value",
      @change="$emit('input', $event)",
      item-text="title",
      item-value="filter",
      return-object,
      clearable
      )
      v-list
        v-list-tile(v-for="(filter, index) in filters", :key="filter.title", @click="")
          v-list-tile-content {{ filter.title }}
          v-list-tile-action(v-if="hasAccessToEditFilter")
            div
              v-btn.ma-1(icon, @click="showEditFilterModal(index)")
                v-icon edit
              v-btn.ma-1(icon, @click="showDeleteFilterModal(index)")
                v-icon delete
      v-btn(
      v-if="hasAccessToAddFilter",
      color="primary",
      @click.prevent="showCreateFilterModal"
      ) {{ $t('common.add') }}
</template>

<script>
import { MODALS, USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';
import modalMixin from '@/mixins/modal/modal';

export default {
  mixins: [authMixin, modalMixin],
  props: {
    filters: {
      type: Array,
      default: () => [],
    },
    value: {
      type: Object,
      default: () => ({}),
    },
    hideSelect: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    hasAccessToEditFilter() {
      return this.checkAccess(USERS_RIGHTS.business.alarmList.actions.editFilter);
    },
    hasAccessToAddFilter() {
      return this.checkAccess(USERS_RIGHTS.business.alarmList.actions.addFilter);
    },
  },
  methods: {
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
        name: this.$constants.MODALS.createFilter,
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

            this.$emit('update:filters', this.filters.filter((v, i) => index !== i));
          },
        },
      });
    },
  },
};
</script>
