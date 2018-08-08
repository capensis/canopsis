<template lang="pug">
  v-expansion-panel
    v-expansion-panel-content.weather-watcher-entity-expansion-panel(hide-actions)
      .pa-2(slot="header", :class="stateColorClass")
        span.pl-1.white--text.subheading.entity-title {{ entity.name }}
        div.actions-button-wrapper
          v-btn(fab, small)
            v-icon local_play
          v-btn(v-show="!hasActivePbehavior", fab, small)
            v-icon pause
          v-btn(v-show="hasActivePbehavior", fab, small)
            v-icon play_arrow
      v-card
        template(v-for="attribute in Object.keys(attributes)")
          v-card-text
            v-layout(row, wrap)
              v-flex.text-md-right(xs4)
                b {{ $t(`modals.weatherWatcher.${attribute}`) }}
              v-flex.pl-2(xs8)
                span {{ attributes[attribute] }}
          v-divider
        v-card-text
          v-layout(row, wrap)
            v-flex.text-md-right(xs4)
              b {{ $t('modals.weatherWatcher.ticketing') }}
            v-flex.pl-2(xs8)
              v-icon local_play
        v-divider
</template>

<script>
import { ENTITIES_STATES } from '@/constants';

export default {
  props: {
    entity: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      attributes: {
        criticity: this.entity.criticity,
        organization: this.entity.org,
        numberOk: this.entity.stats ? this.entity.stats.ok : this.$t('modals.weatherWatcher.noData'),
        numberKo: this.entity.stats ? this.entity.stats.ko : this.$t('modals.weatherWatcher.noData'),
        state: this.entity.state.val,
      },
    };
  },
  computed: {
    stateColorClass() {
      if (this.hasActivePbehavior) {
        return 'grey lighten-1';
      }

      const classes = {
        [ENTITIES_STATES.ok]: 'green darken-1',
        [ENTITIES_STATES.minor]: 'yellow darken-1',
        [ENTITIES_STATES.major]: 'orange darken-1',
        [ENTITIES_STATES.critical]: 'red darken-1',
      };

      return classes[this.attributes.state];
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
    }
  }
</style>

<style scoped>
  .actions-button-wrapper {
    float: right;
  }
</style>
