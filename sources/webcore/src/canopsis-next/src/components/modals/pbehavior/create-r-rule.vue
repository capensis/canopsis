<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span Create RRule
      template(slot="text")
        r-rule-form(v-model="value")
      template(slot="actions")
        v-btn(
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(type="submit") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { MODALS } from '@/constants';

import RRuleForm from '@/components/forms/rrule.vue';

import modalInnerMixin from '@/mixins/modal/inner';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createRRule,
  components: {
    RRuleForm,
    ModalWrapper,
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      value: this.modal.config.rrule || '',
    };
  },
  methods: {
    submit() {
      if (this.config.action) {
        this.config.action(this.value);
      }

      this.$modals.hide();
    },
  },
};
</script>
