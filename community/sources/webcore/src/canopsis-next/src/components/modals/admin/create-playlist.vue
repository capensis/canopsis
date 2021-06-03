<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ title }}
      template(slot="text")
        playlist-form(v-model="form", :groups="playlistGroups", :tabs-pending="tabsPending")
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

import { MAX_LIMIT, MODALS } from '@/constants';

import { playlistToForm, formToPlaylist } from '@/helpers/forms/playlist';

import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';
import { permissionsEntitiesPlaylistTabMixin } from '@/mixins/permissions/entities/playlist-tab';
import { submittableMixin } from '@/mixins/submittable';
import { confirmableModalMixin } from '@/mixins/confirmable-modal';
import { validationErrorsMixin } from '@/mixins/form/validation-errors';

import PlaylistForm from '@/components/other/playlists/form/playlist-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

const { mapGetters: mapEntitiesGetters } = createNamespacedHelpers('entities');

export default {
  name: MODALS.createPlaylist,
  $_veeValidate: {
    validator: 'new',
  },
  components: { PlaylistForm, ModalWrapper },
  mixins: [
    entitiesViewGroupMixin,
    permissionsEntitiesPlaylistTabMixin,
    submittableMixin(),
    validationErrorsMixin(),
    confirmableModalMixin(),
  ],
  data() {
    return {
      playlistGroups: [],
      tabsPending: true,
      form: playlistToForm(this.modal.config.playlist),
    };
  },
  computed: {
    ...mapEntitiesGetters(['getList']),

    title() {
      return this.config.title || this.$t('modals.createPlaylist.create.title');
    },
  },
  mounted() {
    this.fetchGroups();
  },
  methods: {
    async fetchGroups() {
      const { playlist } = this.config;

      this.tabsPending = true;

      const { data: groups } = await this.fetchGroupsListWithoutStore({
        params: { limit: MAX_LIMIT, with_views: true },
      });

      this.playlistGroups = groups;

      if (playlist && playlist.tabs_list.length) {
        this.form.tabs_list = this.getAvailableTabsByIds(playlist.tabs_list);
      }

      this.tabsPending = false;
    },

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          if (this.config.action) {
            await this.config.action(formToPlaylist(this.form));
          }

          this.$modals.hide();
        } catch (err) {
          this.setFormErrors(err);
        }
      }
    },
  },
};
</script>