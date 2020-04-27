<template lang="pug">
  div
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.playlists') }}
    .white
      v-data-table(
        :headers="headers",
        :items="playlists",
        :loading="pending",
        :rows-per-page-items="$config.PAGINATION_PER_PAGE_VALUES",
        item-key="_id"
      )
        template(slot="items", slot-scope="props")
          tr
            td {{ props.item.name }}
            td
              enabled-column(:value="props.item.fullscreen")
            td
              enabled-column(:value="props.item.enabled")
            td {{ props.item.interval | interval }}
            td
              v-btn.ma-0(
                v-if="hasUpdateAnyPlaylistAccess",
                data-test="editButton",
                icon,
                @click="$emit('edit', props.item)"
              )
                v-icon edit
              v-btn(
                v-if="hasCreateAnyPlaylistAccess",
                icon,
                small,
                @click.stop="$emit('duplicate', props.item)"
              )
                v-icon file_copy
              v-btn.ma-0(
                v-if="hasDeleteAnyPlaylistAccess",
                data-test="deleteButton",
                icon,
                @click="$emit('delete', props.item._id)"
              )
                v-icon(color="error") delete
</template>

<script>
import rightsTechnicalPlaylistMixin from '@/mixins/rights/technical/playlist';

import EnabledColumn from '@/components/tables/enabled-column.vue';

export default {
  filters: {
    interval(interval) {
      return interval && interval.interval && `${interval.interval}${interval.unit}`;
    },
  },
  components: {
    EnabledColumn,
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
