<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ $t('modals.createPlaylist.manageTabs') }}
      template(slot="text")
        manage-playlist-tabs-form(v-model="form.selectedTabs", :groups="groups")
      template(slot="actions")
        v-btn(
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { MODALS } from '@/constants';

import submittableMixin from '@/mixins/submittable';
import confirmableModalMixin from '@/mixins/confirmable-modal';

import ManagePlaylistTabsForm from '@/components/other/playlists/form/manage-playlist-tabs-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.managePlaylistTabs,

  components: {
    ManagePlaylistTabsForm,
    ModalWrapper,
  },
  mixins: [
    submittableMixin(),
    confirmableModalMixin(),
  ],
  data() {
    return {
      form: {
        selectedTabs: this.modal.config.selectedTabs || [],
      },
    };
  },
  computed: {
    groups() {
      return this.config.groups || [];
    },
  },
  methods: {
    async submit() {
      if (this.config.action) {
        this.config.action(this.form.selectedTabs);

        this.$modals.hide();
      }
    },
  },
};
</script>
