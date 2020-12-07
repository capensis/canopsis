<template lang="pug">
  div.white
    v-data-table(
      :headers="headers",
      :items="playlists",
      :loading="pending",
      :rows-per-page-items="$config.PAGINATION_PER_PAGE_VALUES",
      item-key="_id",
      expand
    )
      template(slot="items", slot-scope="props")
        tr(@click="props.expanded = !props.expanded")
          td {{ props.item.name }}
          td
            enabled-column(:value="props.item.fullscreen")
          td
            enabled-column(:value="props.item.enabled")
          td {{ props.item.interval | interval }}
          td
            action-btn(:tooltip="$t('common.play')")
              v-btn.mx-1(
                slot="button",
                :to="{ name: 'playlist', params: { id: props.item._id, userAction: true }, query: { autoplay: true } }",
                icon
              )
                v-icon play_arrow
            action-btn(:tooltip="$t('common.copyLink')")
              v-btn.mx-1(
                slot="button",
                v-clipboard:copy="getPlaylistRoute(props.item)",
                v-clipboard:success="onSuccessCopied",
                v-clipboard:error="onErrorCopied",
                icon,
                @click.stop
              )
                v-icon content_copy
            action-btn(
              v-if="hasCreateAnyPlaylistAccess",
              type="duplicate",
              @click="$emit('duplicate', props.item)"
            )
            action-btn(
              v-if="hasUpdateAnyPlaylistAccess",
              type="edit",
              @click="$emit('edit', props.item)"
            )
            action-btn(
              v-if="hasDeleteAnyPlaylistAccess",
              type="delete",
              @click="$emit('delete', props.item._id)"
            )
      template(slot="expand", slot-scope="{ item }")
        playlist-list-expand-item(:playlist="item")
</template>

<script>
import { getApplicationHost } from '@/helpers/router';

import rightsTechnicalPlaylistMixin from '@/mixins/rights/technical/playlist';

import EnabledColumn from '@/components/tables/enabled-column.vue';
import ActionBtn from '@/components/tables/action-btn.vue';

import PlaylistListExpandItem from './playlists-list-expand-item.vue';

export default {
  filters: {
    interval(interval) {
      return interval && interval.interval && `${interval.interval}${interval.unit}`;
    },
  },
  components: {
    EnabledColumn,
    ActionBtn,
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
          sortable: false,
        },
      ];
    },
  },
  methods: {
    getPlaylistRoute({ _id }) {
      const { href } = this.$router.resolve({ name: 'playlist', params: { id: _id }, query: { autoplay: true } });

      return `${getApplicationHost()}${href}`;
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
