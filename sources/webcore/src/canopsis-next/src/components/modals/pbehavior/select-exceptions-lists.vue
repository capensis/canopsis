<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ $t('modals.selectExceptionsLists.title') }}
      template(slot="text")
        choose-exceptions-lists(v-model="exceptions")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(type="submit") {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { MODALS } from '@/constants';

import ChooseExceptionsLists from '@/components/other/pbehavior/exceptions/partials/choose-exceptions-lists.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.selectExceptionsLists,
  components: {
    ChooseExceptionsLists,
    ModalWrapper,
  },
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
