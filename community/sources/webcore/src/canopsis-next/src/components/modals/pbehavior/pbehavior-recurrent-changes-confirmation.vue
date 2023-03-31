<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(:close="cancel")
      template(#title="")
        span {{ $t('modals.pbehaviorRecurrentChangesConfirmation.title') }}
      template(#text="")
        v-radio-group(
          v-model="type",
          hide-details,
          mandatory
        )
          v-radio(
            :value="$constants.PBEHAVIOR_PLANNING_EVENT_CHANGING_TYPES.selected",
            :label="$t('modals.pbehaviorRecurrentChangesConfirmation.fields.selected')",
            color="primary"
          )
          v-radio(
            :value="$constants.PBEHAVIOR_PLANNING_EVENT_CHANGING_TYPES.all",
            :label="$t('modals.pbehaviorRecurrentChangesConfirmation.fields.all')",
            color="primary"
          )
      template(#actions="")
        v-btn(depressed, flat, @click="cancel") {{ $t('common.cancel') }}
        v-btn.primary(type="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS, PBEHAVIOR_PLANNING_EVENT_CHANGING_TYPES } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.pbehaviorRecurrentChangesConfirmation,
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
        this.config.action(this.type);
      }

      this.$modals.hide();
    },
    cancel() {
      if (this.config.cancel) {
        this.config.cancel();
      }

      this.$modals.hide();
    },
  },
};
</script>
