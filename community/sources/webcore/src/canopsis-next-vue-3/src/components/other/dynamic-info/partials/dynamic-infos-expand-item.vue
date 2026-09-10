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

    <v-window v-model="activeTab" class="bg-secondary-lighten-2">
      <c-window-expand-item value="information">
        <v-data-table
          :items="info.infos"
          :headers="infosTableHeaders"
          :no-data-text="$t('common.noData')"
        />
      </c-window-expand-item>

      <c-window-expand-item value="patterns" cols="8">
        <dynamic-info-patterns-form :form="patterns" readonly />
      </c-window-expand-item>
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
