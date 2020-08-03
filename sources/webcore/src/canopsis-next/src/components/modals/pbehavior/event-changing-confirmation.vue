<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span Modify
      template(slot="text")
        v-radio-group(
          v-model="value",
          hide-details,
          mandatory
        )
          v-radio(
            :value="0",
            label="Only selected period",
            color="primary"
          )
          v-radio(
            :value="1",
            label="All the periods",
            color="primary"
          )
      template(slot="actions")
        v-btn(depressed, flat, @click="cancel") {{ $t('common.cancel') }}
        v-btn.primary(type="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import ExceptionsDatesLists from '@/components/other/pbehavior/exceptions/exceptions-dates-lists.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.selectExceptionsDatesLists,
  components: {
    ExceptionsDatesLists,
    ModalWrapper,
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      value: 0,
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
