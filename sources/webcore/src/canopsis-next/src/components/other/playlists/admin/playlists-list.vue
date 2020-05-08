<template lang="pug">
  div
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.playlists') }}
    .white
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
              v-btn.ma-0(
                :to="{ name: 'playlist', params: { id: props.item._id }, query: { autoplay: true } }",
                icon
              )
                v-icon play_arrow
              v-btn.ma-0(
                v-if="hasUpdateAnyPlaylistAccess",
                icon,
                @click.stop="$emit('edit', props.item)"
              )
                v-icon edit
              v-btn.ma-0(
                v-if="hasDeleteAnyPlaylistAccess",
                icon,
                @click.stop="$emit('delete', props.item._id)"
              )
                v-icon(color="error") delete
        template(slot="expand", slot-scope="{ item }")
          playlist-list-expand-item(:playlist="item")
</template>

<script>
import rightsTechnicalPlaylistMixin from '@/mixins/rights/technical/playlist';

import EnabledColumn from '@/components/tables/enabled-column.vue';

import PlaylistListExpandItem from './playlists-list-expand-item.vue';

export default {
  filters: {
    interval(interval) {
      return interval && interval.interval && `${interval.interval}${interval.unit}`;
    },
  },
  components: {
    EnabledColumn,
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
};
</script>
