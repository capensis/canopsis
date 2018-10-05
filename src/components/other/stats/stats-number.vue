<template lang="pug">
  v-layout
    v-flex(xs2)
      v-btn(icon, @click="showSettings")
        v-icon settings
      v-card
        v-card-title
          v-layout(justify-center)
            h2 {{ statName }}
        v-data-iterator(
          :items="stats",
          content-tag="v-layout",
          rows-per-page-text="",
          row,
          wrap,
        )
          v-flex(
            slot="item",
            slot-scope="props",
            xs12,
          )
            v-list(dense)
              v-list-tile
                v-list-tile-content {{ props.item.entity.name }}
                v-list-tile-content.align-end {{ props.item.value }}
</template>

<script>
import entitiesStatsMixin from '@/mixins/entities/stats';
import sideBarMixin from '@/mixins/side-bar/side-bar';
import widgetQueryMixin from '@/mixins/widget/query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import { SIDE_BARS } from '@/constants';

export default {
  mixins: [
    entitiesStatsMixin,
    sideBarMixin,
    widgetQueryMixin,
    entitiesUserPreferenceMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      stats: [],
      statName: 'Test',
    };
  },
  methods: {
    showSettings() {
      this.showSideBar({
        name: SIDE_BARS.statsNumberSettings,
        config: {
          widget: this.widget,
        },
      });
    },
    async fetchList() {
      const query = { ...this.query };

      this.stats = await this.fetchStatValuesWithoutStore({
        params: query,
      });
    },
  },
};
</script>

