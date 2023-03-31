<template lang="pug">
  c-advanced-data-table(
    :headers="headers",
    :items="playlists",
    :loading="pending",
    :total-items="totalItems",
    :pagination="pagination",
    advanced-pagination,
    select-all,
    expand,
    search,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(#fullscreen="{ item }")
      c-enabled(:value="item.fullscreen")
    template(#interval="{ item }") {{ item.interval | duration }}
    template(#enabled="{ item }")
      c-enabled(:value="item.enabled")
    template(#actions="{ item }")
      v-layout(row)
        c-action-btn(:tooltip="$t('common.play')")
          template(#button="")
            v-btn.mx-1.ma-0(
              :to="getPlaylistRouteById(item._id, true)",
              icon
            )
              v-icon play_arrow
        c-copy-btn(
          :value="getPlaylistRouteFullUrlById(item._id)",
          :tooltip="$t('common.copyLink')",
          @success="onSuccessCopied",
          @error="onErrorCopied"
        )
        c-action-btn(
          v-if="hasCreateAnyPlaylistAccess",
          type="duplicate",
          @click="$emit('duplicate', item)"
        )
        c-action-btn(
          v-if="hasUpdateAnyPlaylistAccess",
          type="edit",
          @click="$emit('edit', item)"
        )
        c-action-btn(
          v-if="hasDeleteAnyPlaylistAccess",
          type="delete",
          @click="$emit('remove', item._id)"
        )
    template(#expand="{ item }")
      playlist-list-expand-item(:playlist="item")
</template>

<script>
import { APP_HOST } from '@/config';
import { ROUTES_NAMES } from '@/constants';

import { removeTrailingSlashes } from '@/helpers/url';

import { permissionsTechnicalPlaylistMixin } from '@/mixins/permissions/technical/playlist';

import PlaylistListExpandItem from './playlists-list-expand-item.vue';

export default {
  components: {
    PlaylistListExpandItem,
  },
  mixins: [permissionsTechnicalPlaylistMixin],
  props: {
    playlists: {
      type: Array,
      required: true,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    pagination: {
      type: Object,
      required: true,
    },
  },
  computed: {
    headers() {
      return [
        {
          text: this.$t('common.name'),
          value: 'name',
        },
        {
          text: this.$t('common.fullscreen'),
          value: 'fullscreen',
          sortable: false,
        },
        {
          text: this.$t('common.enabled'),
          value: 'enabled',
        },
        {
          text: this.$t('common.interval'),
          value: 'interval',
          sortable: false,
        },
        {
          text: this.$t('common.actionsLabel'),
          value: 'actions',
          sortable: false,
        },
      ];
    },
  },
  methods: {
    getPlaylistRouteById(id, userAction = false) {
      return {
        name: ROUTES_NAMES.playlist,
        params: { id, userAction },
        query: { autoplay: true },
      };
    },

    getPlaylistRouteFullUrlById(id) {
      const { href } = this.$router.resolve(this.getPlaylistRouteById(id));

      return removeTrailingSlashes(`${APP_HOST}${href}`);
    },

    onSuccessCopied() {
      this.$popups.success({ text: this.$t('success.linkCopied') });
    },

    onErrorCopied() {
      this.$popups.success({ text: this.$t('errors.default') });
    },
  },
};
</script>
