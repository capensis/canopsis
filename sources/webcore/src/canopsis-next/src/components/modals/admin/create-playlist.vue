<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ title }}
      template(slot="text")
        playlist-form(v-model="form")
      template(slot="actions")
        v-btn(
          data-test="createPbehaviorCancelButton",
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit",
          data-test="createPbehaviorSubmitButton"
        ) {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { MODALS } from '@/constants';
import { getDefaultPlaylist } from '@/helpers/entities';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import PlaylistForm from '@/components/other/playlist/playlist-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createPlaylist,

  $_veeValidate: {
    validator: 'new',
  },

  components: { PlaylistForm, ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin()],
  data() {
    return {
      form: this.modal.config.playlist ? cloneDeep(this.modal.config.playlist) : getDefaultPlaylist(),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createPlaylist.create.title');
    },
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
