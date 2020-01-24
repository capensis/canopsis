<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('modals.eventFilterRule.externalData') }}
      template(slot="text")
        v-textarea(:value="externalDataValue", @input="checkValidity")
      template(slot="actions")
        v-btn(@click="$modals.hide", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled || error",
          :loading="submitting",
          type="submit"
        ) {{ error ? $t('errors.JSONNotValid') : $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.eventFilterRuleExternalData,
  components: { ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin()],
  data() {
    return {
      newVal: {},
      error: false,
    };
  },
  computed: {
    externalDataValue() {
      try {
        return JSON.stringify(this.config.value, undefined, 4);
      } catch (err) {
        console.error(err);
        return '';
      }
    },
  },
  methods: {
    checkValidity(value) {
      try {
        this.newVal = JSON.parse(value);
        this.error = false;
      } catch (err) {
        this.error = true;
      }
    },
    async submit() {
      await this.config.action(this.newVal);
      this.$modals.hide();
    },
  },
};
</script>
