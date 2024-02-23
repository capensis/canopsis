<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <playlist-form
          v-model="form"
          :groups="playlistGroups"
          :tabs-pending="tabsPending"
        />
      </template>
      <template #actions="">
        <v-btn
          depressed
          text
          @click="$modals.hide"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="isDisabled"
          :loading="submitting"
          class="primary"
          type="submit"
        >
          {{ $t('common.saveChanges') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MAX_LIMIT, MODALS } from '@/constants';

import { playlistToForm, formToPlaylist } from '@/helpers/entities/playlist/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';
import { permissionsEntitiesPlaylistTabMixin } from '@/mixins/permissions/entities/playlist-tab';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

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
    modalInnerMixin,
    entitiesViewGroupMixin,
    permissionsEntitiesPlaylistTabMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
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
        params: { limit: MAX_LIMIT, with_views: true, with_tabs: true },
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
        if (this.config.action) {
          await this.config.action(formToPlaylist(this.form));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
