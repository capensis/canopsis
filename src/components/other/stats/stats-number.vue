<template lang="pug">
  v-container(fluid)
    v-card
      v-card-title
        v-layout(justify-center)
          h2 {{ statName }}
        v-btn(icon, @click="showSettings")
          v-icon settings
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
                  v-chip(:style="{ backgroundColor: chipColorAndText(props.item.value).color }")
                    div.body-1 {{ chipColorAndText(props.item.value).text }}
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
    chipColorAndText(value) {
      if (this.widget.parameters.yesNoMode) {
        return {
          text: value === 0 ? this.$t('common.no') : this.$t('common.yes'),
          color: value === 0 ? this.widget.parameters.statColors.ok : this.widget.parameters.statColors.critical,
        };
      }

      let color;
      if (value < this.widget.parameters.criticityLevels.minor) {
        color = this.widget.parameters.statColors.ok;
      } else if (value < this.widget.parameters.criticityLevels.major) {
        color = this.widget.parameters.statColors.minor;
      } else if (value < this.widget.parameters.criticityLevels.critical) {
        color = this.widget.parameters.statColors.major;
      } else {
        color = this.widget.parameters.statColors.critical;
      }

      return {
        text: value,
        color,
      };
    },
  },
};
</script>

