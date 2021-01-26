<template lang="pug">
  c-advanced-data-table.white(
    :headers="headers",
    :items="playlists",
    :loading="pending",
    :rows-per-page-items="$config.PAGINATION_PER_PAGE_VALUES",
    expand
  )
    template(slot="fullscreen", slot-scope="props")
      c-enabled(:value="props.item.fullscreen")
    template(slot="interval", slot-scope="props") {{ props.item.interval | interval }}
    template(slot="enabled", slot-scope="props")
      c-enabled(:value="props.item.enabled")
    template(slot="actions", slot-scope="props")
      v-layout(row)
        c-action-btn(:tooltip="$t('common.play')")
          v-btn.mx-1(
            slot="button",
            :to="getPlaylistRouteById(props.item._id, true)",
            icon
          )
            v-icon play_arrow
        c-action-btn(:tooltip="$t('common.copyLink')")
          v-btn.mx-1(
            slot="button",
            v-clipboard:copy="getPlaylistRouteFullUrlById(props.item._id)",
            v-clipboard:success="onSuccessCopied",
            v-clipboard:error="onErrorCopied",
            icon,
            @click.stop
          )
            v-icon content_copy
        c-action-btn(
          v-if="hasCreateAnyPlaylistAccess",
          type="duplicate",
          @click="$emit('duplicate', props.item)"
        )
        c-action-btn(
          v-if="hasUpdateAnyPlaylistAccess",
          type="edit",
          @click="$emit('edit', props.item)"
        )
        c-action-btn(
          v-if="hasDeleteAnyPlaylistAccess",
          type="delete",
          @click="$emit('remove', props.item._id)"
        )
    template(slot="expand", slot-scope="props")
      playlist-list-expand-item(:playlist="props.item")
</template>

<script>
import { APP_HOST } from '@/config';

import rightsTechnicalPlaylistMixin from '@/mixins/rights/technical/playlist';

import PlaylistListExpandItem from './playlists-list-expand-item.vue';

export default {
  filters: {
    interval(interval) {
      return interval && interval.interval && `${interval.interval}${interval.unit}`;
    },
  },
  components: {
    PlaylistListExpandItem,
  },
  mixins: [rightsTechnicalPlaylistMixin],
  props: {
    playlists: {
      type: Array,
      required: true,
    },
    pending: {
      type: Boolean,
      default: false,
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
        },
        {
          text: this.$t('common.enabled'),
          value: 'enabled',
        },
        {
          text: this.$t('common.interval'),
          value: 'interval',
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
        name: 'playlist',
        params: { id, userAction },
        query: { autoplay: true },
      };
    },

    getPlaylistRouteFullUrlById(id) {
      const { href } = this.$router.resolve(this.getPlaylistRouteById(id));

      return `${APP_HOST}${href}`;
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
