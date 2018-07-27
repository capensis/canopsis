<template lang="pug">
  v-expansion-panel
    v-expansion-panel-content
      div(slot="header", :class="stateColorClass") Watcher name
      v-card
        v-card-text
          attribute-block
            template(slot="name")
              | Actions
            template(slot="content")
              v-btn
                v-icon local_play
              v-btn
                v-icon pause
        v-divider
        template(v-for="attribute in Object.keys(attributes)")
          v-card-text
            attribute-block
              template(slot="name")
                | {{ $t(`modals.watcherData.${attribute}`) }}
              template(slot="content")
                | {{ attributes[attribute] }}
          v-divider
</template>

<script>
import { WATCHER_STATES } from '@/constants';
import AttributeBlock from './attribute-block.vue';

export default {
  components: {
    AttributeBlock,
  },
  data() {
    return {
      attributes: {
        criticity: 'Good',
        organization: 'Noveo',
        nombreOk: 'No data',
        nombreKo: 'Some data',
        state: 3,
      },
    };
  },
  computed: {
    stateColorClass() {
      const classes = {
        [WATCHER_STATES.ok]: 'color-ok',
        [WATCHER_STATES.minor]: 'color-minor',
        [WATCHER_STATES.major]: 'color-major',
        [WATCHER_STATES.critical]: 'color-critical',
      };

      return classes[0];
    },
  },
};
</script>

<style scoped>
  .expansion-panel__header {
    padding: 12px 12px!important;
  }

  .expansion-panel__header div:first-child {
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
    width: 60%;
    margin-left: 40%;
  }

  .btn {
    margin: 0;
    max-width: 40px;
    min-width: 30px;
  }
</style>
