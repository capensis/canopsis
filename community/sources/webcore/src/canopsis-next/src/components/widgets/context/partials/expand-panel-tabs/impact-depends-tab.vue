<template lang="pug">
  div
    v-progress-linear.ma-0(:active="pending", height="2", indeterminate)
    div.pa-3
      v-layout
        v-flex(xs6)
          h3.headline.text-xs-center.my-1.white--text {{ $t('context.impacts') }}
          v-container
            v-card
              v-card-text
                v-data-iterator(:items="impact")
                  template(#item="props")
                    v-flex
                      v-card
                        v-card-title {{ props.item }}
                  template(#no-data="")
                    v-flex
                      v-card
                        v-card-title {{ $t('common.noData') }}
        v-flex(xs6)
          h3.headline.text-xs-center.my-1.white--text {{ $t('context.dependencies') }}
          v-container
            v-card
              v-card-text
                v-data-iterator(:items="depends")
                  template(#item="props")
                    v-flex
                      v-card
                        v-card-title {{ props.item }}
                  template(#no-data="")
                    v-flex
                      v-card
                        v-card-title {{ $t('common.noData') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import Observer from '@/services/observer';

const { mapActions } = createNamespacedHelpers('entity');

export default {
  inject: {
    $periodicRefresh: {
      default() {
        return new Observer();
      },
    },
  },
  props: {
    entity: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      pending: false,
      impact: [],
      depends: [],
    };
  },
  mounted() {
    this.fetchList();

    this.$periodicRefresh.register(this.fetchList);
  },
  beforeDestroy() {
    this.$periodicRefresh.unregister(this.fetchList);
  },
  methods: {
    ...mapActions({
      fetchContextEntityContextGraphWithoutStore: 'fetchContextGraphWithoutStore',
    }),

    async fetchList() {
      this.pending = true;

      const { impact, depends } = await this.fetchContextEntityContextGraphWithoutStore({ id: this.entity._id });

      this.impact = impact;
      this.depends = depends;
      this.pending = false;
    },
  },
};
</script>
