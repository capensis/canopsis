<template lang="pug">
  div
    v-card
      v-card-title
        v-layout(justify-center)
          h2 {{ widget.parameters.stat.title }}
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
                  v-chip(:style="{ backgroundColor: getChipColor(props.item.value) }")
                    div.body-1 {{ getChipText(props.item.value) }}
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
    };
  },
  computed: {
    getChipColor() {
      return (value) => {
        const { yesNoMode, criticityLevels, statColors } = this.widget.parameters;

        if (yesNoMode) {
          return value === 0 ? statColors.ok : statColors.critical;
        }

        if (value < criticityLevels.minor) {
          return statColors.ok;
        } else if (value < criticityLevels.major) {
          return statColors.minor;
        } else if (value < criticityLevels.critical) {
          return statColors.major;
        }

        return statColors.critical;
      };
    },
    getChipText() {
      return (value) => {
        const { yesNoMode } = this.widget.parameters;

        if (yesNoMode) {
          return value === 0 ? this.$t('common.no') : this.$t('common.yes');
        }

        return value;
      };
    },
  },
  methods: {
    showSettings() {
      this.showSideBar({
        name: this.$constants.SIDE_BARS.statsNumberSettings,
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
  },
};
</script>

