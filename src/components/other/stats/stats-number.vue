<template lang="pug">
  v-container(fluid)
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
              v-list-tile-content
                ellipsis(:text="props.item.entity.name")
              v-list-tile-content.align-end
                v-layout(align-center)
                  template(v-if="widget.parameters.yesNoMode")
                    v-chip(v-if="props.item.value === 0", :style="{ backgroundColor: widget.parameters.statColors.ok }")
                      div.body-1 {{ $t('common.no') }}
                    v-chip(v-else, :style="{ backgroundColor: widget.parameters.statColors.critical }")
                      div.body-1 {{ $t('common.yes') }}
                  v-chip(v-else, :style="{ backgroundColor: getCriticity(props.item.value) }", small)
                    div.body-1 {{ props.item.value }}
                  div.caption
                    template(v-if="props.item.trend >= 0") + {{ props.item.trend }}
                    template(v-else) - {{ props.item.trend }}
</template>

<script>
import Ellipsis from '@/components/tables/ellipsis.vue';
import entitiesStatsMixin from '@/mixins/entities/stats';
import sideBarMixin from '@/mixins/side-bar/side-bar';
import widgetQueryMixin from '@/mixins/widget/query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import { SIDE_BARS } from '@/constants';

export default {
  components: {
    Ellipsis,
  },
  mixins: [
    entitiesStatsMixin,
    sideBarMixin,
    widgetQueryMixin,
    entitiesUserPreferenceMixin,
  ],
  props: {
    rowId: {
      type: String,
      required: true,
    },
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
          rowId: this.rowId,
        },
      });
    },
    async fetchList() {
      const query = { ...this.query };

      this.stats = await this.fetchStatValuesWithoutStore({
        params: query,
      });
    },
    getCriticity(value) {
      if (value > this.widget.parameters.criticityLevels.minor) {
        return this.widget.parameters.statColors.minor;
      } else if (value > this.widget.parameters.criticityLevels.major) {
        return this.widget.parameters.statColors.major;
      } else if (value > this.widget.parameters.criticityLevels.critical) {
        return this.widget.parameters.statColors.critical;
      }
      return this.widget.parameters.statColors.ok;
    },
  },
};
</script>

