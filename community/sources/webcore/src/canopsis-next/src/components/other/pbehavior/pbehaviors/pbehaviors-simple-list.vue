<template>
  <v-layout column>
    <v-layout justify-space-between align-center>
      <v-flex v-if="shownUserTimezone" xs5>
        <c-timezone-field
          v-model="timezone"
          server
        />
      </v-flex>
      <v-layout class="gap-2" justify-end>
        <c-action-fab-btn
          v-if="addable"
          :tooltip="$t('modals.createPbehavior.create.title')"
          icon="add"
          color="primary"
          left
          @click="showCreatePbehaviorModal"
        />
        <c-action-fab-btn
          :tooltip="$t('modals.pbehaviorsCalendar.title')"
          icon="calendar_today"
          color="secondary"
          left
          @click="showPbehaviorsCalendarModal"
        />
      </v-layout>
    </v-layout>
    <c-advanced-data-table
      :items="pbehaviors"
      :headers="headers"
      :loading="pending"
      :dense="dense"
    >
      <template #enabled="{ item }">
        <c-enabled :value="item.enabled" />
      </template>
      <template #tstart="{ item }">
        {{ formatIntervalDate(item, 'tstart') }}
      </template>
      <template #tstop="{ item }">
        {{ formatIntervalDate(item, 'tstop') }}
      </template>
      <template #rrule_end="{ item }">
        {{ formatRruleEndDate(item) }}
      </template>
      <template #rrule="{ item }">
        <v-icon>{{ item.rrule ? 'check' : 'clear' }}</v-icon>
      </template>
      <template #icon="{ item }">
        <v-icon color="primary">
          {{ item.type.icon_name }}
        </v-icon>
      </template>
      <template #status="{ item }">
        <v-icon :color="item.is_active_status ? 'primary' : 'error'">
          $vuetify.icons.settings_sync
        </v-icon>
      </template>
      <template #actions="{ item }">
        <pbehavior-actions
          :pbehavior="item"
          :removable="removable"
          :updatable="updatable"
          @refresh="fetchList"
        />
      </template>
    </c-advanced-data-table>
  </v-layout>
</template>

<script>
import {
  computed,
  ref,
  inject,
  onMounted,
  onBeforeUnmount,
} from 'vue';

import { MODALS } from '@/constants';

import Observer from '@/services/observer';

import { createEntityIdPatternByValue } from '@/helpers/entities/pattern/form';

import { useModals } from '@/hooks/modals';
import { useI18n } from '@/hooks/i18n';
import { usePbehavior } from '@/hooks/store/modules/pbehavior';
import { usePendingWithLocalQuery } from '@/hooks/query/shared';

import { usePbehaviorDateFormat } from '@/components/other/pbehavior/pbehaviors/hooks/pbehavior-date-format';

import PbehaviorActions from './partials/pbehavior-actions.vue';

export default {
  components: { PbehaviorActions },
  props: {
    entity: {
      type: Object,
      required: true,
    },
    removable: {
      type: Boolean,
      default: false,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
    dense: {
      type: Boolean,
      default: false,
    },
    addable: {
      type: Boolean,
      default: false,
    },
    withActiveStatus: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const modals = useModals();
    const { t, tc } = useI18n();
    const { fetchPbehaviorsByEntityIdWithoutStore } = usePbehavior();
    const {
      timezone,
      shownUserTimezone,
      formatIntervalDate,
      formatRruleEndDate,
    } = usePbehaviorDateFormat();

    inject('$system');

    const periodicRefresh = inject('$periodicRefresh', new Observer());

    const pbehaviors = ref([]);

    const headers = computed(() => {
      const result = [
        { text: t('common.name'), value: 'name' },
        { text: t('common.author'), value: 'author.display_name' },
        { text: t('pbehavior.isEnabled'), value: 'enabled' },
        { text: t('pbehavior.begins'), value: 'tstart' },
        { text: t('pbehavior.ends'), value: 'tstop' },
        { text: t('pbehavior.rruleEnd'), value: 'rrule_end' },
        { text: t('common.recurrence'), value: 'rrule' },
        { text: t('common.type'), value: 'type.name' },
        { text: t('common.reason'), value: 'reason.name' },
        { text: tc('common.icon', 1), value: 'icon' },
      ];

      if (props.withActiveStatus) {
        result.push({ text: t('common.status'), value: 'status', sortable: false });
      }

      if (props.updatable || props.removable) {
        result.push({ text: t('common.actionsLabel'), value: 'actions', sortable: false });
      }

      return result;
    });

    const {
      pending,
      fetchHandlerWithQuery: fetchList,
    } = usePendingWithLocalQuery({
      fetchHandler: async () => {
        try {
          pbehaviors.value = await fetchPbehaviorsByEntityIdWithoutStore({
            id: props.entity._id,
            params: {
              with_flags: true,
            },
          });
        } catch (err) {
          console.warn(err);
        }
      },
    });

    /**
     * Opens the pbehaviors calendar modal for a specific entity.
     */
    const showPbehaviorsCalendarModal = () => modals.show({
      name: MODALS.pbehaviorsCalendar,
      config: {
        title: t('modals.pbehaviorsCalendar.entity.title', { name: props.entity.name }),
        entityId: props.entity._id,
      },
    });

    /**
     * Opens the creation pbehavior modal with pre-configured entity settings.
     * Automatically refreshes the pbehavior list after successful submission.
     */
    const showCreatePbehaviorModal = () => modals.show({
      name: MODALS.pbehaviorPlanning,
      config: {
        entityPattern: createEntityIdPatternByValue(props.entity._id),
        entities: [props.entity],
        afterSubmit: fetchList,
      },
    });

    onMounted(() => {
      fetchList();

      periodicRefresh.register(fetchList);
    });

    onBeforeUnmount(() => periodicRefresh.unregister(fetchList));

    return {
      timezone,
      shownUserTimezone,
      pending,
      pbehaviors,
      headers,

      formatIntervalDate,
      formatRruleEndDate,
      fetchList,
      showPbehaviorsCalendarModal,
      showCreatePbehaviorModal,
    };
  },
};
</script>
