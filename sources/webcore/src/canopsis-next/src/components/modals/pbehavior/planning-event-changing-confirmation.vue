<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span Modify
      template(slot="text")
        v-radio-group(
          v-model="type",
          hide-details,
          mandatory
        )
          v-radio(
            :value="$constants.PBEHAVIOR_PLANNING_EVENT_CHANGING_TYPES.selected",
            label="Only selected period",
            color="primary"
          )
          v-radio(
            :value="$constants.PBEHAVIOR_PLANNING_EVENT_CHANGING_TYPES.all",
            label="All the periods",
            color="primary"
          )
      template(slot="actions")
        v-btn(depressed, flat, @click="cancel") {{ $t('common.cancel') }}
        v-btn.primary(type="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS, PBEHAVIOR_PLANNING_EVENT_CHANGING_TYPES } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.planningEventChangingConfirmation,
  components: {
    ModalWrapper,
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      type: PBEHAVIOR_PLANNING_EVENT_CHANGING_TYPES.selected,
    };
  },
  methods: {
    submit() {
      if (this.config.action) {
        this.config.action(this.value);
      }

      this.$modals.hide();
    },
    cancel() {
      if (this.config.cancelAction) {
        this.config.cancelAction();
      }

      this.$modals.hide();
    },
  },
};
</script>
