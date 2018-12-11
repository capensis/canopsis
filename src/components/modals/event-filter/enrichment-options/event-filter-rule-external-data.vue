<template lang="pug">
   v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline Event filter - External datas
    v-card-text
      v-textarea(:value="JSON.stringify(config.value, undefined, 4)", @input="checkValidity")
    v-divider
    v-layout.pa-2(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(
      @click.prevent="submit",
      :disabled="error ? true : false",
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
      error: '',
    };
  },
  methods: {
    checkValidity(value) {
      try {
        this.newVal = JSON.parse(value);
        this.error = '';
      } catch (err) {
        this.error = err.message;
      }
    },
    submit() {
      this.config.action(this.newVal);
      this.hideModal();
    },
  },
};
</script>
