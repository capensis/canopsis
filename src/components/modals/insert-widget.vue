<template lang="pug">
  v-card
    v-card-title
      v-layout(justify-space-between, align-center)
        h2 Select a widget
        v-btn(@click="hideModal", icon, small)
          v-icon close
    v-card-text
      v-layout(row)
        v-flex(xs12 sm8 offset-sm2)
          v-list
            v-list-tile(
            v-for="widgetType in widgetsTypes",
            @click="selectWidgetType(widgetType.value)",
            :key="widgetType.title"
            )
              v-list-tile-action
                v-icon {{ widgetType.icon }}
              v-list-tile-content
                v-list-tile-title {{ widgetType.title }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import modalMixin from '@/mixins/modal/modal';
import { MODALS } from '@/constants';

const { mapActions } = createNamespacedHelpers('view');

/**
 * Modal to add a time filter on alarm-list
 */
export default {
  name: MODALS.insertWidget,
  mixins: [modalMixin],
  data() {
    return {
      widgetsTypes: [
        { title: 'listalarm', icon: 'list', value: 'listalarm' },
        { title: 'crudcontext', icon: 'list', value: 'crudcontext' },
      ],
    };
  },
  methods: {
    ...mapActions(['addWidget']),
    selectWidgetType(widgetType) {
      console.log(widgetType);
      this.addWidget({ widget: { xtype: widgetType } });
    },
  },
};
</script>
