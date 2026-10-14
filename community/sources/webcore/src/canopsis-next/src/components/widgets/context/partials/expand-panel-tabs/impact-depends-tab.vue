<template>
  <div>
    <v-progress-linear
      :active="pending"
      class="ma-0"
      height="2"
      indeterminate
    />
    <div class="pa-3">
      <c-layout>
        <v-col cols="6">
          <h3 class="text-h5 text-center my-1 text-white">
            {{ $t('context.impacts') }}
          </h3>
          <v-container>
            <v-card>
              <v-card-text>
                <c-layout v-if="impact.length" gap="2" column>
                  <v-card v-for="item in impact" :key="item">
                    <v-card-title>{{ item }}</v-card-title>
                  </v-card>
                </c-layout>
                <v-card v-else>
                  <v-card-title>{{ $t('common.noData') }}</v-card-title>
                </v-card>
              </v-card-text>
            </v-card>
          </v-container>
        </v-col>
        <v-col cols="6">
          <h3 class="text-h5 text-center my-1 text-white">
            {{ $t('context.dependencies') }}
          </h3>
          <v-container>
            <v-card>
              <v-card-text>
                <c-layout v-if="depends.length" gap="2" column>
                  <v-card v-for="item in depends" :key="item">
                    <v-card-title>{{ item }}</v-card-title>
                  </v-card>
                </c-layout>
                <v-card v-else>
                  <v-card-title>{{ $t('common.noData') }}</v-card-title>
                </v-card>
              </v-card-text>
            </v-card>
          </v-container>
        </v-col>
      </c-layout>
    </div>
  </div>
</template>

<script>
import { mapActions } from 'pinia';

import Observer from '@/services/observer';

import { useEntityStore } from '@/stores/modules/entity';

export default {
  inject: {
    $periodicRefresh: {
      default: () => new Observer(),
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

  beforeUnmount() {
    this.$periodicRefresh.unregister(this.fetchList);
  },

  methods: {
    ...mapActions(useEntityStore, {
      fetchContextEntityContextGraphWithoutStore: 'fetchContextGraphWithoutStore',
    }),

    async fetchList() {
      this.pending = true;

      const { impact, depends } = await this.fetchContextEntityContextGraphWithoutStore({
        id: this.entity._id,
      });

      this.impact = impact;
      this.depends = depends;
      this.pending = false;
    },
  },
};
</script>
