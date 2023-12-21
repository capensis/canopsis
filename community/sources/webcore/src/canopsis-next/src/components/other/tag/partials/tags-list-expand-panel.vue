<template>
  <v-tabs
    v-model="activeTab"
    background-color="secondary lighten-1"
    slider-color="primary"
    dark
    centered
  >
    <v-tab>{{ $tc('common.pattern', 2) }}</v-tab>
    <v-tabs-items
      v-model="activeTab"
      mandatory
    >
      <v-tab-item>
        <v-layout class="py-3">
          <v-flex
            xs12
            md8
            offset-md2
          >
            <v-card>
              <v-card-text>
                <tag-patterns-form
                  :form="patterns"
                  readonly
                />
              </v-card-text>
            </v-card>
          </v-flex>
        </v-layout>
      </v-tab-item>
    </v-tabs-items>
  </v-tabs>
</template>

<script>
import { tagFilterPatternsToForm } from '@/helpers/entities/tag/form';

import TagPatternsForm from '@/components/other/tag/form/tag-patterns-form.vue';

export default {
  components: { TagPatternsForm },
  props: {
    tag: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      activeTab: 0,
    };
  },
  computed: {
    patterns() {
      return tagFilterPatternsToForm(this.tag);
    },
  },
};
</script>
