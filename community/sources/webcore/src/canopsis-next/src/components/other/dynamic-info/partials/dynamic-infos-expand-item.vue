<template>
  <v-tabs
    background-color="secondary lighten-1"
    slider-color="primary"
    dark
    centered
  >
    <v-tab>{{ $tc('common.information', 2) }}</v-tab>
    <v-tab-item>
      <v-layout class="py-3 secondary lighten-2">
        <v-flex
          xs12
          md8
          offset-md2
        >
          <v-card>
            <v-card-text>
              <v-data-table
                :items="info.infos"
                :headers="infosTableHeaders"
                :no-data-text="$t('common.noData')"
              >
                <template #items="{ item }">
                  <tr>
                    <td>{{ item.name }}</td>
                    <td>{{ item.value }}</td>
                  </tr>
                </template>
              </v-data-table>
            </v-card-text>
          </v-card>
        </v-flex>
      </v-layout>
    </v-tab-item>
    <v-tab>{{ $tc('common.pattern', 2) }}</v-tab>
    <v-tab-item>
      <v-layout class="py-3 secondary lighten-2">
        <v-flex
          xs12
          md8
          offset-md2
        >
          <v-card>
            <v-card-text>
              <dynamic-info-patterns-form
                :form="patterns"
                readonly
              />
            </v-card-text>
          </v-card>
        </v-flex>
      </v-layout>
    </v-tab-item>
  </v-tabs>
</template>

<script>
import { PATTERNS_FIELDS } from '@/constants';

import { filterPatternsToForm } from '@/helpers/entities/filter/form';

import DynamicInfoPatternsForm from '../form/fields/dynamic-info-patterns-form.vue';

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
  computed: {
    patterns() {
      return filterPatternsToForm(this.info, [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity]);
    },

    infosTableHeaders() {
      return [
        { text: this.$t('common.name'), value: 'name' },
        { text: this.$t('common.value'), value: 'value' },
      ];
    },
  },
};
</script>
