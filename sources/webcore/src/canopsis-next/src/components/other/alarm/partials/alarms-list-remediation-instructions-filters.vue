<template lang="pug">
  v-layout(row, wrap, align-center)
    remediation-instructions-filters-list(
      :filters="lockedFilters",
      @input="$listeners['update:lockedFilters']"
    )
    remediation-instructions-filters-list(
      :filters="filters",
      editable,
      closable,
      @input="$listeners['update:filters']"
    )
    v-tooltip(bottom)
      v-btn(
        slot="activator",
        icon,
        small,
        @click="showCreateFilterModal"
      )
        v-icon(:color="filters.length ? 'primary' : 'black'") adjust
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
