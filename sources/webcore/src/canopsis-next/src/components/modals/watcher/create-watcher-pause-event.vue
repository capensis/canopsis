<template lang="pug">
  modal-wrapper
    template(slot="title")
      span {{ $t('modals.createPause.title') }}
    template(slot="text")
      watcher-pause-event-form(v-model="form")
    template(slot="actions")
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
      v-btn.primary(@click="submit") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import WatcherPauseEventForm from '@/components/other/service-weather/watcher-pause-event-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createWatcherPauseEvent,
  $_veeValidate: {
    validator: 'new',
  },
  components: { WatcherPauseEventForm, ModalWrapper },
  mixins: [modalInnerMixin],
  data() {
    return {
      form: {
        comment: '',
        reason: '',
      },
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(this.form);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
