<template lang="pug">
  .weather-watcher-entity-expansion-panel
    v-expansion-panel
      v-expansion-panel-content(hide-actions)
        .pa-2(slot="header", :class="entityClass")
          span.pl-1.white--text.subheading.entity-title {{ entity.name }}
          div.actions-button-wrapper
            v-btn(fab, small)
              v-icon local_play
            v-btn(v-show="!hasActivePbehavior", fab, small)
              v-icon pause
            v-btn(v-show="hasActivePbehavior", fab, small)
              v-icon play_arrow
        v-card
          v-card-text(v-html="compiledTemplate")
        v-divider
</template>

<script>
import get from 'lodash/get';
import pick from 'lodash/pick';
import mapValues from 'lodash/mapValues';

import compile from '@/helpers/handlebars';

export default {
  props: {
    entity: {
      type: Object,
      required: true,
    },
    template: {
      type: String,
    },
  },
  data() {
    const mainAttributes = pick(this.entity, [
      'criticity',
      'org',
    ]);

    const infoAttributes = mapValues(pick(this.entity.infos, [
      'scenario_label',
      'scenario_probe_name',
      'scenario_calendar',
    ]), v => v.value);

    return {
      attributes: {
        ...mainAttributes,
        ...infoAttributes,

        numberOk: get(this.entity.stats, 'ok', this.$t('modals.watcher.noData')),
        numberKo: get(this.entity.stats, 'ko', this.$t('modals.watcher.noData')),
        state: this.entity.state.val,
      },
    };
  },
  computed: {
    entityClass() {
      if (this.hasActivePbehavior) {
        return this.$constants.WATCHER_PBEHAVIOR_COLOR;
      }

      return this.$constants.WATCHER_STATES_COLORS[this.attributes.state];
    },

    hasActivePbehavior() {
      if (!this.entity.pbehavior || !this.entity.pbehavior.length) {
        return false;
      }

      return this.entity.pbehavior.filter((value) => {
        const start = value.dtstart * 1000;
        const end = value.dtend * 1000;
        const now = Date.now();

        return start <= now && now < end;
      }).length;
    },

    compiledTemplate() {
      return compile(this.template, { watcher: this.watcher, entity: this.entity });
    },
  },
};
</script>

<style lang="scss" scoped>
  .weather-watcher-entity-expansion-panel /deep/ .v-expansion-panel__header {
    padding: 0;
    height: auto;

    .entity-title {
      line-height: 52px;
    }

    .actions-button-wrapper {
      float: right;
    }
  }
</style>
