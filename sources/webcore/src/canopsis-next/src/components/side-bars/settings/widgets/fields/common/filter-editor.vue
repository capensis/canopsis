<template lang="pug">
  settings-button-field(
  :title="$t('settings.liveReporting.title')",
  :isEmpty="isValueEmpty",
  @create="openFilterModal",
  @edit="openFilterModal",
  @delete="deleteFilter",
  )
    .subheading(slot="title") {{ $t('settings.filterEditor') }}
      .font-italic.caption.ml-1(v-show="!required") ({{ $t('common.optional') }})
</template>

<script>
import { isEmpty } from 'lodash';

import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal';

import SettingsButtonField from '../partials/button-field.vue';

export default {
  components: { SettingsButtonField },
  mixins: [modalMixin],
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
  },
  computed: {
    isValueEmpty() {
      return isEmpty(this.value);
    },
  },
  methods: {
    openFilterModal() {
      this.showModal({
        name: MODALS.createFilter,
        config: {
          title: this.$t('modals.filter.create.title'),
          filter: this.value,
          hiddenFields: this.hiddenFields,
          action: filterObject => this.$emit('input', filterObject),
        },
      });
    },
    deleteFilter() {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => this.$emit('input', {}),
        },
      });
    },
  },
};
</script>
