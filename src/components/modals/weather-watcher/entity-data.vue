<template lang="pug">
  v-expansion-panel
    v-expansion-panel-content(hide-actions)
      div(slot="header", :class="stateColorClass")
        span {{ watchedEntity.name }}
        div.actions-button-wrapper
          v-btn
            v-icon local_play
          v-btn
            v-icon pause
      v-card
        template(v-for="attribute in Object.keys(attributes)")
          v-card-text
            attribute-block
              template(slot="name")
                span {{ $t(`modals.watcherData.${attribute}`) }}
              template(slot="content")
                span {{ attributes[attribute] }}
          v-divider
        v-card-text
          attribute-block
            template(slot="name")
              span {{ $t('modals.watcherData.ticketing') }}
            template(slot="content")
              v-icon local_play
        v-divider
</template>

<script>
import moment from 'moment';

import { ENTITIES_STATES } from '@/constants';

import AttributeBlock from './attribute-block.vue';

export default {
  components: {
    AttributeBlock,
  },
  props: {
    watcher: {
      type: Object,
      required: true,
    },
    watchedEntity: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      attributes: {
        criticity: this.watchedEntity.criticity,
        organization: this.watchedEntity.org,
        numberOk: this.watchedEntity.stats ? this.watchedEntity.stats.ok : this.$t('modals.watcherData.noData'),
        numberKo: this.watchedEntity.stats ? this.watchedEntity.stats.ko : this.$t('modals.watcherData.noData'),
        state: this.watchedEntity.state.val,
      },
    };
  },
  computed: {
    stateColorClass() {
      if (this.isUnderActivePBehavior) {
        return 'color-pbehavior';
      }

      const classes = {
        [ENTITIES_STATES.ok]: 'color-ok',
        [ENTITIES_STATES.minor]: 'color-minor',
        [ENTITIES_STATES.major]: 'color-major',
        [ENTITIES_STATES.critical]: 'color-critical',
      };

      return classes[this.attributes.state];
    },
    isUnderActivePBehavior() {
      if (!this.watchedEntity.pbehavior || !this.watchedEntity.pbehavior.length) {
        return false;
      }

      let underPBehavior = false;

      this.watchedEntity.pbehavior.forEach((pBehavior) => {
        if (pBehavior.isActive) {
          if (pBehavior.tstart <= moment().unix() && moment().unix() <= pBehavior.tstop) {
            underPBehavior = true;
          }
        }
      });

      return underPBehavior;
    },
  },
};
</script>

<style scoped>
  .expansion-panel__header {
    padding: 12px 12px !important;
  }

  .expansion-panel__header > div:first-child {
    padding: 15px;
  }

  .color-ok {
    background-color: #43A047;
  }

  .color-minor {
    background-color: #FDD835;
  }

  .color-major {
    background-color: #FB8C00;
  }

  .color-critical {
    background-color: #E53935;
  }

  .color-pbehavior {
    background-color: #BDBDBD;
  }

  .attribute, .divider {
    width: 100%;
  }

  .btn {
    margin: 0;
    max-width: 40px;
    min-width: 30px;
  }

  .actions-button-wrapper {
    float: right;
    padding: 0;
  }
</style>
