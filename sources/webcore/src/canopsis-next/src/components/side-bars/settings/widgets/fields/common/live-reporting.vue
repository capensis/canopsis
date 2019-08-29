<template lang="pug">
  v-container.pa-3(fluid)
    v-layout(align-center, justify-space-between)
      div.subheading {{ $t('settings.liveReporting.title') }}
      v-layout(justify-end)
        v-btn.primary(
          small,
          @click="showEditLiveReportingModal"
        )
          span(v-show="isValueEmpty") {{ $t('common.create') }}
          span(v-show="!isValueEmpty") {{ $t('common.edit') }}
        v-btn.error(
        v-show="!isValueEmpty",
        small,
        @click="clear"
        )
          v-icon delete
</template>

<script>
import { isEmpty } from 'lodash';

import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal';
import formBaseMixin from '@/mixins/form/base';

export default {
  mixins: [
    modalMixin,
    formBaseMixin,
  ],
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    isValueEmpty() {
      return isEmpty(this.value);
    },
  },
  methods: {
    showEditLiveReportingModal() {
      this.showModal({
        name: MODALS.editLiveReporting,
        config: {
          ...this.value,

          action: value => this.updateModel(value),
        },
      });
    },
    clear() {
      this.updateModel({});
    },
  },
};
</script>
