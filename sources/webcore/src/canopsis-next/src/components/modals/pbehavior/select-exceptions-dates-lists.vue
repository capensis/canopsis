<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('modals.selectExceptionsDatesLists.title') }}
      template(slot="text")
        choose-exceptions-dates-lists(v-model="exceptionDates")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(type="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import ChooseExceptionsDatesLists from '@/components/other/pbehavior/dates-exceptions/partials/choose-exceptions-dates-lists.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.selectExceptionsDatesLists,
  components: {
    ChooseExceptionsDatesLists,
    ModalWrapper,
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      exceptionDates: this.modal.config.exceptions,
    };
  },
  methods: {
    submit() {
      if (this.config.action) {
        this.config.action(this.exceptionDates);
      }

      this.$modals.hide();
    },
  },
};
</script>
