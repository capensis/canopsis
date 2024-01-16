<template>
  <widget-settings-item :title="$t('settings.remediationInstructionsFilters')">
    <v-layout wrap>
      <remediation-instructions-filters-list
        v-field="filters"
        :editable="editable"
        :closable="editable"
      />
    </v-layout>
    <v-layout
      v-if="addable"
      wrap
    >
      <v-btn
        class="ml-1"
        color="primary"
        @click="showCreateInstructionsFilterModal"
      >
        {{ $t('common.add') }}
      </v-btn>
    </v-layout>
  </widget-settings-item>
</template>

<script>
import { MODALS } from '@/constants';

import { uid } from '@/helpers/uid';

import { formArrayMixin } from '@/mixins/form';

import RemediationInstructionsFiltersList from '@/components/other/remediation/instructions-filter/remediation-instructions-filters-list.vue';
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem, RemediationInstructionsFiltersList },
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
    addable: {
      type: Boolean,
      default: false,
    },
    editable: {
      type: Boolean,
      default: false,
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
