<template lang="pug">
  v-expansion-panel
    v-expansion-panel-content.weather-watcher-entity-expansion-panel(hide-actions)
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
        template(v-for="(attribute, attributeKey) in attributes")
          v-card-text
            v-layout(row, wrap)
              v-flex.text-md-right(xs4)
                b {{ $t(`modals.watcher.${attributeKey}`) }}:
              v-flex.pl-2(xs8)
                span {{ attribute }}
          v-divider
        v-card-text
          v-layout(row, wrap)
            v-flex.text-md-right(xs4)
              b {{ $t('modals.watcher.ticketing') }}:
            v-flex.pl-2(xs8)
              v-icon local_play
        v-divider
</template>

<script>
import get from 'lodash/get';
import pick from 'lodash/pick';
import mapValues from 'lodash/mapValues';

import { WATCHER_STATES_COLORS, WATCHER_PBEHAVIOR_COLOR } from '@/constants';

export default {
  props: {
    entity: {
      type: Object,
      required: true,
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
        return WATCHER_PBEHAVIOR_COLOR;
      }

      return WATCHER_STATES_COLORS[this.attributes.state];
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
  },
};
</script>

<style lang="scss">
  .weather-watcher-entity-expansion-panel {
    .expansion-panel__header {
      padding: 0;

      .entity-title {
        line-height: 52px;
      }

      .actions-button-wrapper {
        float: right;
      }
    }
  }
</style>
