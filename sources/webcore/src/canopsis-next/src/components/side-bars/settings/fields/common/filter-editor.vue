<template lang="pug">
  settings-button-field(
    :title="$t('settings.liveReporting.title')",
    :isEmpty="isValueEmpty",
    @create="openFilterModal",
    @edit="openFilterModal",
    @delete="deleteFilter"
  )
    .subheading(slot="title") {{ $t('settings.filterEditor') }}
      .font-italic.caption.ml-1(v-show="!required") ({{ $t('common.optional') }})
</template>

<script>
import { isEmpty } from 'lodash';

import { MODALS } from '@/constants';

import SettingsButtonField from '../partials/button-field.vue';

export default {
  components: { SettingsButtonField },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    hiddenFields: {
      type: Array,
      default: () => [],
    },
    required: {
      type: Boolean,
      default: false,
    },
    existingTitles: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    isValueEmpty() {
      return isEmpty(this.value);
    },
  },
  methods: {
    openFilterModal() {
      this.$modals.show({
        name: MODALS.createFilter,
        config: {
          title: this.$t('modals.filter.create.title'),
          filter: this.value,
          hiddenFields: this.hiddenFields,
          existingTitles: this.existingTitles,
          action: filterObject => this.$emit('input', filterObject),
        },
      });
    },
    deleteFilter() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => this.$emit('input', {}),
        },
      });
    },
  },
};
</script>
