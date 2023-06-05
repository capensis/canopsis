<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ $t('modals.selectExceptionsLists.title') }}
      template(#text="")
        choose-exceptions-lists(v-model="exceptions")
      template(#actions="")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(type="submit") {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';

import ChooseExceptionsLists from '@/components/other/pbehavior/exceptions/partials/choose-exceptions-lists.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.selectExceptionsLists,
  components: {
    ChooseExceptionsLists,
    ModalWrapper,
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      exceptions: this.modal.config.exceptions ? cloneDeep(this.modal.config.exceptions) : [],
    };
  },
  methods: {
    submit() {
      if (this.config.action) {
        this.config.action(this.exceptions);
      }

      this.$modals.hide();
    },
  },
};
</script>
