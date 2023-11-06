<template>
  <v-tabs
    background-color="secondary lighten-1"
    slider-color="primary"
    dark
    centered
  >
    <v-tab>{{ $tc('common.pattern', 2) }}</v-tab>
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
    <v-tab>{{ $tc('common.entity', 2) }}</v-tab>
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
    <v-tab>{{ $tc('common.comment', 2) }}</v-tab>
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
    <template v-if="pbehavior.rrule">
      <v-tab>{{ $t('common.recurrence') }}</v-tab>
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
                <pbehavior-recurrence-rule :pbehavior="pbehavior" />
              </v-card-text>
            </v-card>
          </v-flex>
        </v-layout>
      </v-tab-item>
    </template>
  </v-tabs>
</template>

<script>
import { OLD_PATTERNS_FIELDS, PATTERNS_FIELDS } from '@/constants';

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
  computed: {
    patterns() {
      return filterPatternsToForm(
        this.pbehavior,
        [PATTERNS_FIELDS.entity],
        [OLD_PATTERNS_FIELDS.mongoQuery],
      );
    },
  },
};
</script>
