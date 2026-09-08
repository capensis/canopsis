<template>
  <div>
    <v-tabs
      v-model="activeTab"
      bg-color="secondary-lighten-1"
      color="primary"
      align-tabs="center"
    >
      <v-tab value="information">
        {{ $t('common.information', 2) }}
      </v-tab>
      <v-tab value="patterns">
        {{ $t('common.pattern', 2) }}
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab">
      <v-window-item value="information">
        <c-layout class="py-3 secondary lighten-2">
          <v-row class="dynamic-infos-expand-item__row">
            <v-col cols="12" md="8" offset-md="2">
              <v-card>
                <v-card-text>
                  <v-data-table
                    :items="info.infos"
                    :headers="infosTableHeaders"
                    :no-data-text="$t('common.noData')"
                  />
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>
        </c-layout>
      </v-window-item>

      <v-window-item value="patterns">
        <v-row class="py-3 secondary lighten-2 dynamic-infos-expand-item__row" justify="center">
          <v-col cols="8">
            <v-card>
              <v-card-text>
                <dynamic-info-patterns-form :form="patterns" readonly />
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-window-item>
    </v-window>
  </div>
</template>

<script>
import { PATTERNS_FIELDS } from '@/constants';

import { filterPatternsToForm } from '@/helpers/entities/filter/form';

import DynamicInfoPatternsForm from '../form/dynamic-info-patterns-form.vue';

export default {
  components: {
    DynamicInfoPatternsForm,
  },
  props: {
    info: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      activeTab: 'information',
    };
  },
  computed: {
    patterns() {
      return filterPatternsToForm(this.info, [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity]);
    },

    infosTableHeaders() {
      return [
        { title: this.$t('common.name'), key: 'name' },
        { title: this.$t('common.value'), key: 'value' },
      ];
    },
  },
};
</script>

<style lang="scss" scoped>
.dynamic-infos-expand-item__row {
  width: 100%;
  margin: 0;
}
</style>
