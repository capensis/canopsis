<template lang="pug">
  v-tabs(color="secondary lighten-1", slider-color="primary", dark, centered)
    v-tab {{ $tc('common.information', 2) }}
    v-tab-item
      v-layout.py-3.secondary.lighten-2
        v-flex(xs12, md8, offset-md2)
          v-card
            v-card-text
              v-data-table(
                :items="info.infos",
                :headers="infosTableHeaders",
                :no-data-text="$t('common.noData')"
              )
                template(#items="{ item }")
                  tr
                    td {{ item.name }}
                    td {{ item.value }}
    v-tab {{ $tc('common.pattern', 2) }}
    v-tab-item
      v-layout.py-3.secondary.lighten-2
        v-flex(xs12, md8, offset-md2)
          v-card
            v-card-text
              dynamic-info-patterns-form(:form="patterns", readonly)
</template>

<script>
import { OLD_PATTERNS_FIELDS, PATTERNS_FIELDS } from '@/constants';

import { filterPatternsToForm } from '@/helpers/forms/filter';

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
      return filterPatternsToForm(
        this.info,
        [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity],
        [OLD_PATTERNS_FIELDS.alarm, OLD_PATTERNS_FIELDS.entity],
      );
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
