<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ title }}
      template(slot="text")
        playlist-form(v-model="form", :groups="availableGroups")
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
import { createNamespacedHelpers } from 'vuex';

import { MODALS } from '@/constants';

import { getDefaultPlaylist } from '@/helpers/entities';
import { playlistToForm, formToPlaylist } from '@/helpers/forms/playlist';

import modalInnerMixin from '@/mixins/modal/inner';
import entitiesViewGroupMixin from '@/mixins/entities/view/group';
import rightsEntitiesGroupMixin from '@/mixins/rights/entities/group';
import rightsEntitiesPlaylistTabMixin from '@/mixins/rights/entities/playlist-tab';
import submittableMixin from '@/mixins/submittable';

import PlaylistForm from '@/components/other/playlists/form/playlist-form.vue';
import ProgressOverlay from '@/components/layout/progress/progress-overlay.vue';

import ModalWrapper from '../modal-wrapper.vue';

const { mapGetters: mapEntitiesGetters } = createNamespacedHelpers('entities');

export default {
  name: MODALS.createPlaylist,

  $_veeValidate: {
    validator: 'new',
  },

  components: { PlaylistForm, ProgressOverlay, ModalWrapper },
  mixins: [
    modalInnerMixin,
    entitiesViewGroupMixin,
    rightsEntitiesGroupMixin,
    rightsEntitiesPlaylistTabMixin,
    submittableMixin(),
  ],
  data() {
    const { playlist } = this.modal.config;

    return {
      form: playlist ? playlistToForm(playlist) : getDefaultPlaylist(),
    };
  },
  computed: {
    ...mapEntitiesGetters(['getList']),

    title() {
      const type = this.modal.config.playlist ? 'edit' : 'create';

      return this.$t(`modals.createPlaylist.${type}.title`);
    },
  },
  async mounted() {
    const { playlist } = this.config;

    if (!this.groupsPending) {
      await this.fetchGroupsList();
    }

    if (playlist && playlist.tabs_list.length) {
      this.form.tabs_list = this.getAvailableTabsByIds(playlist.tabs_list);
    }
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(formToPlaylist(this.form));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
