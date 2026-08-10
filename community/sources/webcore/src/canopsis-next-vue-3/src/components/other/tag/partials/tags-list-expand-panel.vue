<template>
  <v-tabs
    v-model="activeTab"
    bg-color="secondary lighten-1"
    slider-color="primary"
    align-tabs="center"
  >
    <v-tab :value="TABS.patterns">
      {{ $t('common.pattern', 2) }}
    </v-tab>

    <v-window v-model="activeTab" class="bg-secondary-lighten-2">
      <v-window-item :value="TABS.patterns">
        <v-row class="py-3 ma-0">
          <v-col cols="12" md="8" offset-md="2">
            <v-card>
              <v-card-text>
                <tag-patterns-form :form="patterns" readonly />
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-window-item>
    </v-window>
  </v-tabs>
</template>

<script>
import TagPatternsForm from '@/components/other/tag/form/tag-patterns-form.vue';
import { tagFilterPatternsToForm } from '@/helpers/entities/tag/form';

const TABS = {
  patterns: 'patterns',
};

export default {
  components: {
    TagPatternsForm,
  },

  props: {
    tag: {
      type: Object,
      default: () => ({}),
    },
  },

  data() {
    return {
      TABS,
      activeTab: TABS.patterns,
    };
  },

  computed: {
    patterns() {
      return tagFilterPatternsToForm(this.tag);
    },
  },
};
</script>
