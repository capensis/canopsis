<template>
  <div>
    <v-tabs
      v-model="activeTab"
      bg-color="secondary-lighten-1"
      color="primary"
      align-tabs="center"
    >
      <v-tab value="patterns">
        {{ t('common.pattern', 2) }}
      </v-tab>
      <v-tab value="entities">
        {{ t('common.entity', 2) }}
      </v-tab>
      <v-tab value="comments">
        {{ t('common.comment', 2) }}
      </v-tab>
      <v-tab v-if="pbehavior.rrule" value="recurrence">
        {{ t('common.recurrence') }}
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab" class="bg-secondary-lighten-2">
      <c-window-expand-item value="patterns" cols="8">
        <pbehavior-patterns-form
          :form="patterns"
          :pbehavior-id="pbehavior._id"
          pbehavior-counter-type
          readonly
        />
      </c-window-expand-item>

      <c-window-expand-item value="entities" cols="10">
        <pbehavior-entities :pbehavior="pbehavior" />
      </c-window-expand-item>

      <c-window-expand-item value="comments" cols="10">
        <pbehavior-comments :comments="pbehavior.comments" />
      </c-window-expand-item>

      <c-window-expand-item
        v-if="pbehavior.rrule"
        value="recurrence"
        cols="10"
      >
        <pbehavior-recurrence-rule :pbehavior="pbehavior" />
      </c-window-expand-item>
    </v-window>
  </div>
</template>

<script>
import { computed, ref } from 'vue';

import { PATTERNS_FIELDS } from '@/constants';

import { filterPatternsToForm } from '@/helpers/entities/filter/form';

import { useI18n } from '@/hooks/i18n';

import PbehaviorPatternsForm from '../form/pbehavior-patterns-form.vue';

import PbehaviorComments from './pbehavior-comments.vue';
import PbehaviorEntities from './pbehavior-entities.vue';
import PbehaviorRecurrenceRule from './pbehavior-recurrence-rule.vue';

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

  setup(props) {
    const { t } = useI18n();
    const activeTab = ref('patterns');

    const patterns = computed(() => filterPatternsToForm(props.pbehavior, [PATTERNS_FIELDS.entity]));

    return {
      t,
      activeTab,
      patterns,
    };
  },
};
</script>
