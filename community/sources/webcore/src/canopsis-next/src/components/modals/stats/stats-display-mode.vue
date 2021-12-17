<template lang="pug">
  v-form(data-test="statsDisplayModeModal", @submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ $t('settings.statsNumbers.displayMode') }}
      template(slot="text")
        stats-display-mode-form(v-model="form")
      template(slot="actions")
        v-btn(
          data-test="statsDisplayModeCancelButton",
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit",
          data-test="statsDisplayModeSubmitButton"
        ) {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { MODALS, STATS_DISPLAY_MODE, STATS_DISPLAY_MODE_PARAMETERS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import StatsDisplayModeForm from '@/components/widgets/stats/stats-display-mode-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.statsDisplayMode,
  components: { StatsDisplayModeForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const { displayMode } = this.modal.config;
    const defaultDisplayMode = {
      mode: STATS_DISPLAY_MODE.criticity,
      parameters: STATS_DISPLAY_MODE_PARAMETERS,
    };

    return {
      form: cloneDeep(displayMode || defaultDisplayMode),
    };
  },
  methods: {
    async submit() {
      if (this.config.action) {
        await this.config.action(this.form);
      }

      this.$modals.hide();
    },
  },
};
</script>
