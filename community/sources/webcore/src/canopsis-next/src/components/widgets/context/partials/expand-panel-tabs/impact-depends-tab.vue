<template>
  <div>
    <v-progress-linear
      :active="pending"
      class="ma-0"
      height="2"
      indeterminate
    />
    <div class="pa-3">
      <v-layout>
        <v-flex xs6>
          <h3 class="text-h5 text-center my-1 white--text">
            {{ $t('context.impacts') }}
          </h3>
          <v-container>
            <v-card>
              <v-card-text>
                <v-data-iterator :items="impact">
                  <template #item="props">
                    <v-flex>
                      <v-card>
                        <v-card-title>{{ props.item }}</v-card-title>
                      </v-card>
                    </v-flex>
                  </template>
                  <template #no-data="">
                    <v-flex>
                      <v-card>
                        <v-card-title>{{ $t('common.noData') }}</v-card-title>
                      </v-card>
                    </v-flex>
                  </template>
                </v-data-iterator>
              </v-card-text>
            </v-card>
          </v-container>
        </v-flex>
        <v-flex xs6>
          <h3 class="text-h5 text-center my-1 white--text">
            {{ $t('context.dependencies') }}
          </h3>
          <v-container>
            <v-card>
              <v-card-text>
                <v-data-iterator :items="depends">
                  <template #item="props">
                    <v-flex>
                      <v-card>
                        <v-card-title>{{ props.item }}</v-card-title>
                      </v-card>
                    </v-flex>
                  </template>
                  <template #no-data="">
                    <v-flex>
                      <v-card>
                        <v-card-title>{{ $t('common.noData') }}</v-card-title>
                      </v-card>
                    </v-flex>
                  </template>
                </v-data-iterator>
              </v-card-text>
            </v-card>
          </v-container>
        </v-flex>
      </v-layout>
    </div>
  </div>
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
