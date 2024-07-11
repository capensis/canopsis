<template>
  <v-layout
    wrap
    align-center
  >
    <remediation-instructions-filters-list
      :filters="lockedFilters"
      :closable="editable"
      @input="$listeners['update:lockedFilters']"
    />
    <remediation-instructions-filters-list
      :filters="filters"
      :editable="editable"
      :closable="editable"
      @input="$listeners['update:filters']"
    />
    <c-simple-tooltip
      v-if="addable"
      :content="$t('remediation.instructionsFilter.button')"
      bottom
    >
      <template #activator="{ on }">
        <v-btn
          class="ma-0"
          icon
          v-on="on"
          @click="showCreateFilterModal"
        >
          <v-icon :color="buttonIconColor">
            assignment
          </v-icon>
        </v-btn>
      </template>
    </c-simple-tooltip>
  </v-layout>
</template>

<script>
import { MODALS } from '@/constants';

import { uid } from '@/helpers/uid';

import RemediationInstructionsFiltersList from '@/components/other/remediation/instructions-filter/remediation-instructions-filters-list.vue';

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
    addable: {
      type: Boolean,
      default: false,
    },
    editable: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    hasAnyEnabledFilters() {
      return this.filters.length || this.lockedFilters.filter(filter => !filter.disabled).length;
    },
    buttonIconColor() {
      return this.hasAnyEnabledFilters ? 'primary' : '';
    },
  },
  methods: {
    showCreateFilterModal() {
      this.$modals.show({
        name: MODALS.createRemediationInstructionsFilter,
        config: {
          filters: this.filters,
          action: newFilter => this.$emit(
            'update:filters',
            [...this.filters, { _id: uid(), ...newFilter }],
          ),
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
