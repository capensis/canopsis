<template lang="pug">
  v-expansion-panel
    v-expansion-panel-content(hide-actions)
      div(slot="header", :class="stateColorClass")
        | {{ watchedEntity.name }}
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
        criticity: 'Good',
        organization: this.watchedEntity.org,
        nombreOk: 'No data',
        nombreKo: 'Some data',
        state: this.watchedEntity.state.val,
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

      return classes[this.attributes.state];
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
