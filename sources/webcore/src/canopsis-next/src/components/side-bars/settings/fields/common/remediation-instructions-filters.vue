<template lang="pug">
  v-list-group(data-test="filterOnOpenResolved")
    v-list-tile(slot="activator") {{ $t('settings.remediationInstructionsFilters') }}
    v-container
      v-layout(row, wrap)
        remediation-instructions-filters-list(
          v-field="filters",
          :editable="hasAccessToEditFilter",
          :closable="hasAccessToEditFilter"
        )
      v-layout(v-if="hasAccessToAddFilter", row, wrap)
        v-btn.ml-1(
          color="primary",
          @click="showCreateInstructionsFilterModal"
        ) {{ $t('common.add') }}
</template>

<script>
import { MODALS } from '@/constants';

import uid from '@/helpers/uid';

import formArrayMixin from '@/mixins/form/array';

import RemediationInstructionsFiltersList
  from '@/components/other/remediation/instructions-filter/remediation-instructions-filters-list.vue';

export default {
  components: { RemediationInstructionsFiltersList },
  mixins: [formArrayMixin],
  model: {
    prop: 'filters',
    event: 'input',
  },
  props: {
    filters: {
      type: Array,
      default: () => [],
    },
    hasAccessToAddFilter: {
      type: Boolean,
      default: true,
    },
    hasAccessToEditFilter: {
      type: Boolean,
      default: true,
    },
  },
  methods: {
    showCreateInstructionsFilterModal() {
      this.$modals.show({
        name: MODALS.createRemediationInstructionsFilter,
        config: {
          filters: this.filters,
          action: newFilter => this.addItemIntoArray({ _id: uid(), ...newFilter }),
        },
      });
    },
  },
};
</script>
