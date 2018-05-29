<template lang="pug">
  v-card
    v-card-text
      div(v-html="output")
</template>

<script>
import HandleBars from 'handlebars';
import { createNamespacedHelpers } from 'vuex';

const { mapGetters } = createNamespacedHelpers('modal');

export default {
  name: 'more-infos',
  data() {
    return {
      template: '<h1>{{entity.type}}</h1><p>{{alarm.v.connector}}</p>',
    };
  },
  computed: {

    ...mapGetters(['config']),

    output() {
      const output = HandleBars.compile(this.template);
      const context = { alarm: this.config.alarm.props, entity: this.config.alarm.props.entity };
      return output(context);
    },
  },
};
</script>
