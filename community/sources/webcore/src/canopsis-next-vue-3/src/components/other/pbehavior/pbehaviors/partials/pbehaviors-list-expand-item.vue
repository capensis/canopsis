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

    <c-layout class="pa-3 secondary lighten-2">
      <v-card class="pa-3 pbehaviors-list-expand-item__card">
        <v-window v-model="activeTab" class="pt-2">
          <v-window-item value="patterns">
            <div class="pbehaviors-list-expand-item__content">
              <pbehavior-patterns-form
                :form="patterns"
                :pbehavior-id="pbehavior._id"
                pbehavior-counter-type
                readonly
              />
            </div>
          </v-window-item>

          <v-window-item value="entities">
            <pbehavior-entities :pbehavior="pbehavior" />
          </v-window-item>

          <v-window-item value="comments">
            <pbehavior-comments :comments="pbehavior.comments" />
          </v-window-item>

          <v-window-item v-if="pbehavior.rrule" value="recurrence">
            <pbehavior-recurrence-rule :pbehavior="pbehavior" />
          </v-window-item>
        </v-window>
      </v-card>
    </c-layout>
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

    const patterns = computed(() =>
      filterPatternsToForm(props.pbehavior, [PATTERNS_FIELDS.entity]),
    );

    return {
      t,
      activeTab,
      patterns,
    };
  },
};
</script>

<style lang="scss" scoped>
.pbehaviors-list-expand-item__card {
  width: 100%;
}

.pbehaviors-list-expand-item__content {
  width: min(100%, 900px);
  margin: 0 auto;
}
</style>
