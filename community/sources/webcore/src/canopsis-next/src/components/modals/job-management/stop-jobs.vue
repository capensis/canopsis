<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        {{ t('modals.stopJobs.title') }}
      </template>

      <template #text>
        <v-row>
          <v-col
            v-for="warning in warnings"
            :key="warning.key"
            cols="12"
          >
            <v-alert type="warning" variant="tonal">
              {{ warning.text }}
            </v-alert>
          </v-col>

          <v-col
            v-for="item in itemsWithPendingJobs"
            :key="item._id"
            cols="12"
          >
            <div class="stop-jobs__item pa-3">
              <v-row dense>
                <v-col cols="12">
                  <strong>{{ item.rule_name }}</strong>
                </v-col>
                <v-col cols="12">
                  {{ t('modals.stopJobs.pendingCount', item.pending_jobs, {
                    count: item.pending_jobs,
                  }) }}
                </v-col>
                <v-col cols="12">
                  <c-radio-group v-model="choices[item._id]" hide-details>
                    <v-radio
                      :value="true"
                      :label="t(`modals.stopJobs.options.stop.${getStopOptionKey(item)}`)"
                      color="primary"
                    />
                    <v-radio
                      :value="false"
                      :label="t('modals.stopJobs.options.continue')"
                      color="primary"
                    />
                  </c-radio-group>
                </v-col>
              </v-row>
            </div>
          </v-col>
        </v-row>
      </template>

      <template #actions>
        <v-btn :disabled="submitting" variant="text" @click="close">
          {{ t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="submitting"
          :loading="submitting"
          color="primary"
          type="submit"
        >
          {{ t('jobs.actions.stopJob') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script setup>
import { computed, ref } from 'vue';

import { JOB_RULE_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { usePendingHandler } from '@/hooks/query/pending';

const props = defineProps({
  modal: {
    type: Object,
    required: true,
  },
});

const { t } = useI18n();
const { config, close } = useInnerModal(props);
const items = computed(() => config.value.items ?? []);
const jobType = computed(() => config.value.jobType);
const itemsWithPendingJobs = computed(() => items.value.filter(item => item.pending_jobs > 0));
const choices = ref(Object.fromEntries(itemsWithPendingJobs.value.map(item => [item._id, true])));

const warnings = computed(() => {
  const groups = itemsWithPendingJobs.value.reduce((acc, item) => {
    const key = item.rule_type === JOB_RULE_TYPES.scenario ? 'scenario' : jobType.value;
    acc[key] = (acc[key] ?? 0) + item.pending_jobs;

    return acc;
  }, {});

  return Object.entries(groups).map(([key, count]) => ({
    key,
    text: t(`modals.stopJobs.warnings.${key}`, count, { count }),
  }));
});

const getStopOptionKey = item => item.rule_type === JOB_RULE_TYPES.scenario
  ? 'scenario'
  : jobType.value;

const { pending: submitting, handler: submit } = usePendingHandler(async () => {
  await config.value.action?.(choices.value);
  close();
});
</script>

<style lang="scss" scoped>
.stop-jobs__item {
  border: 1px solid rgba(0, 0, 0, 0.12);
  border-radius: 4px;
}
</style>
