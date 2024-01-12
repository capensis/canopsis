<template>
  <v-tabs
    v-model="activeTab"
    background-color="secondary lighten-1"
    slider-color="primary"
    dark
    centered
  >
    <v-tab>{{ $tc('common.pattern', 2) }}</v-tab>
    <v-tab>{{ $tc('common.entity', 2) }}</v-tab>
    <v-tab>{{ $tc('common.comment', 2) }}</v-tab>
    <v-tab v-if="pbehavior.rrule">
      {{ $t('common.recurrence') }}
    </v-tab>
    <v-tabs-items
      v-model="activeTab"
      mandatory
    >
      <v-tab-item>
        <v-layout
          class="py-3"
        >
          <v-flex
            xs12
            sm10
            offset-sm1
          >
            <v-card>
              <v-card-text>
                <pbehavior-patterns-form
                  :form="patterns"
                  readonly
                />
              </v-card-text>
            </v-card>
          </v-flex>
        </v-layout>
      </v-tab-item>
      <v-tab-item>
        <v-layout
          class="py-3"
        >
          <v-flex
            xs12
            sm10
            offset-sm1
          >
            <v-card>
              <v-card-text>
                <pbehavior-entities :pbehavior="pbehavior" />
              </v-card-text>
            </v-card>
          </v-flex>
        </v-layout>
      </v-tab-item>
      <v-tab-item>
        <v-layout
          class="py-3"
        >
          <v-flex
            xs12
            sm10
            offset-sm1
          >
            <v-card>
              <v-card-text class="pa-0">
                <pbehavior-comments :comments="pbehavior.comments" />
              </v-card-text>
            </v-card>
          </v-flex>
        </v-layout>
      </v-tab-item>
      <v-tab-item v-if="pbehavior.rrule">
        <v-layout
          class="py-3"
        >
          <v-flex
            xs12
            sm10
            offset-sm1
          >
            <v-card>
              <v-card-text>
                <pbehavior-recurrence-rule :pbehavior="pbehavior" />
              </v-card-text>
            </v-card>
          </v-flex>
        </v-layout>
      </v-tab-item>
    </v-tabs-items>
  </v-tabs>
</template>

<script>
import { PATTERNS_FIELDS } from '@/constants';

import { filterPatternsToForm } from '@/helpers/entities/filter/form';

import PbehaviorPatternsForm from '../form/pbehavior-patterns-form.vue';

import PbehaviorComments from './pbehavior-comments.vue';
import PbehaviorRecurrenceRule from './pbehavior-recurrence-rule.vue';
import PbehaviorEntities from './pbehavior-entities.vue';

export default {
  components: {
    PbehaviorPatternsForm,
    PbehaviorComments,
    PbehaviorRecurrenceRule,
    PbehaviorEntities,
  },
  props: {
    pbehavior: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      activeTab: 0,
    };
  },
  computed: {
    patterns() {
      return filterPatternsToForm(this.pbehavior, [PATTERNS_FIELDS.entity]);
    },
  },
};
</script>
