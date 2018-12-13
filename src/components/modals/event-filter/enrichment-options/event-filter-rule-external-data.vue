<template lang="pug">
   v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline Event filter - External datas
    v-card-text
      v-textarea(:value="externalDataValue", @input="checkValidity")
    v-divider
    v-layout.pa-2(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(
      @click.prevent="submit",
      :disabled="error",
      ) {{ error ? $t('errors.JSONNotValid') : $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/modal-inner';

export default {
  name: MODALS.eventFilterRuleExternalData,
  mixins: [modalInnerMixin],
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
    submit() {
      this.config.action(this.newVal);
      this.hideModal();
    },
  },
};
</script>
