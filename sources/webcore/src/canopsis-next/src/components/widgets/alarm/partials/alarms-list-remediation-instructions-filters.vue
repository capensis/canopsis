<template lang="pug">
  v-layout(row, wrap, align-center)
    remediation-instructions-filters-list(
      :filters="lockedFilters",
      :closable="hasAccessToEditFilter",
      @input="$listeners['update:lockedFilters']"
    )
    remediation-instructions-filters-list(
      :filters="filters",
      :editable="hasAccessToEditFilter",
      :closable="hasAccessToEditFilter",
      @input="$listeners['update:filters']"
    )
    v-tooltip(v-if="hasAccessToUserFilter", bottom)
      v-btn(
        slot="activator",
        icon,
        small,
        @click="showCreateFilterModal"
      )
        v-icon(:color="buttonIconColor") adjust
      span {{ $t('remediationInstructionsFilters.button') }}
</template>

<script>
import { MODALS } from '@/constants';

import uid from '@/helpers/uid';

import RemediationInstructionsFiltersList
  from '@/components/other/remediation/instructions-filter/remediation-instructions-filters-list.vue';

export default {
  components: { RemediationInstructionsFiltersList },
  props: {
    filters: {
      type: Array,
      default: () => [],
    },
    lockedFilters: {
      type: Array,
      default: () => [],
    },
    hasAccessToUserFilter: {
      type: Boolean,
      default: true,
    },
    hasAccessToEditFilter: {
      type: Boolean,
      default: true,
    },
    hasAccessToListFilters: {
      type: Boolean,
      default: true,
    },
  },
  computed: {
    hasAnyEnabledFilters() {
      return this.filters.length || this.lockedFilters.filter(filter => !filter.disabled).length;
    },
    buttonIconColor() {
      return this.hasAnyEnabledFilters ? 'primary' : 'black';
    },
  },
  methods: {
    showCreateFilterModal() {
      this.$modals.show({
        name: MODALS.createRemediationInstructionsFilter,
        config: {
          filters: this.filters,
          action: newFilter => this.$emit('update:filters', [...this.filters, { _id: uid(), ...newFilter }]),
        },
      });
    },
  },
};
</script>

<style lang="scss">
.v-chip__custom-close {
  font-size: 20px;
}
</style>
